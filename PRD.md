# PRD: Multiverse IDE Pragmatic MVP（Chat→WBS/Node→AgentRunner 実行）+ Quality Hardening

最終更新: 2025-12-14

## 1. 背景

- AgentRunner Core は `plan_task`/`next_action`/`completion_assessment` と Docker Sandbox を含む実行ループが安定稼働している（`docs/CURRENT_STATUS.md:18`、`docs/design/data-flow.md:13`）。
- Orchestrator は file-based IPC + Scheduler + Executor を備えるが Beta 段階で、WBS/Node 中心の永続化（design/state/history）を前提とした v2 実装が途中（`docs/CURRENT_STATUS.md:20`、`docs/design/orchestrator-persistence-v2.md:11`）。
- 現在の Chat は Meta-agent の `plan_patch`（作成/更新/削除/移動）結果を TaskStore に保存するだけでなく、`design/`（WBS/NodeDesign）と `state/`（NodesRuntime/TasksState）へも永続化するため、Scheduler が依存解決して実行に進める（`internal/chat/handler.go`、`internal/orchestrator/scheduler.go`）。

## 2. 目的 / ゴール

MVP の到達点は「IDE のチャット入力から、WBS/ノード計画を生成・永続化し、その計画に基づいて Orchestrator が AgentRunner を起動してタスクを順次完了させ、IDE 上に結果が表示される」こと。

具体的には:

1. チャット入力 → Meta-agent `plan_patch` → WBS/Node/TaskState の作成/更新/削除/移動が永続化される。
2. `ExecutionOrchestrator` が依存関係を解決し、READY タスクを IPC Queue に流し、`agent-runner` を実行できる。
3. 実行結果で TaskState / NodesRuntime / TaskStore が更新され、IDE が一覧/グラフ表示できる。
4. IDE 上で `milestone/phase/workType/domain` 等の軸で **グルーピング/フィルタリング**できる（最低限 `milestone -> phase -> task` の WBS が成立する、`frontend/ide/src/stores/wbsStore.ts:161`）。

## 3. 非ゴール（MVP では扱わない）

- ログのリアルタイムストリーミングの外部公開（WebSocket/gRPC などの IPC 強化）。※IDE 内は `task:log` を Wails Events で配信する（`internal/orchestrator/executor.go:121`、`internal/orchestrator/events.go:39`）。
- マルチノード/リモート Worker プール。
- 高度な承認フローや差分レビュー UI。
- アニメーションや高度な UI エフェクト。UI は「カクつかず安定して操作できる」ことを優先する。

## 4. ユーザーストーリー

- US-1: 開発者は IDE のチャットに要望を入力し、数秒〜数十秒後に WBS/ノードとタスクリストが生成される。
- US-2: 開発者は **チャットだけで計画生成〜実行開始まで**進められ、必要に応じて停止/一時停止/再開できる（UI の実行操作はフォールバック）。
- US-3: IDE 上で各タスク/ノードのステータス（PENDING/READY/RUNNING/SUCCEEDED/COMPLETED/FAILED/CANCELED/BLOCKED/RETRY_WAIT）が確認でき、生成・更新されたファイル一覧を参照できる（`internal/orchestrator/task_store.go:16`）。
- US-4: IDE 上でタスクが `milestone/phase`（将来: `workType/domain/tags`）で分類され、WBS/Graph の可視化がフラットに潰れない（WBS は `milestone -> phase -> task` 前提、`frontend/ide/src/stores/wbsStore.ts:161`）。
- US-5: 開発者はチャットで「不要タスクを削除」「順序/依存の整理」「フェーズ移動」等を指示でき、既存計画が **重複生成ではなく差分更新**される。

## 5. アーキテクチャ方針

### 5.1 計画と実行の真実源

- 計画（WBS/NodeDesign）は `~/.multiverse/workspaces/<id>/design/` を真実源とする（`docs/design/orchestrator-persistence-v2.md:33`）。
- 実行状態（TasksState/NodesRuntime/AgentsState）は `state/` を真実源とする。
- `internal/orchestrator/task_store.go` の TaskStore は IDE 表示と後方互換のため当面併用し、design/state と同期させる。

### 5.2 Planner/TaskBuilder の配置

MVP では **Chat Handler が Planner/TaskBuilder の役割を兼務**する。

- `plan_patch` 呼び出しは Chat Handler が行う。
- `plan_patch`（create/update/delete/move）結果を design/state/task_store に写像して永続化する。

将来的には Planner を Orchestrator 側に移し、Chat は UI 層へ戻す。

## 6. データモデル（MVP スキーマ）

### 6.1 design/wbs.json

- WBS ルートのみ保持。最低限 `wbs_id`, `project_root`, `root_node_id`, `node_index` を保存する（`internal/orchestrator/persistence/models.go:9`）。

### 6.2 design/nodes/<node-id>.json

- `plan_patch` の `create` を NodeDesign として保存し、`update/move/delete` は NodeDesign/WBS/TaskState に反映する。
- NodeDesign.Dependencies は `plan_patch` の `dependencies` を `node_id` に解決したものを格納する。

主要フィールド:

- `node_id`: UUID または `node-<task-id>` 形式。
- `name`, `summary`: task の `title`/`description`。
- `phase_name`, `milestone`, `wbs_level`: グルーピング/移動のための facet（`frontend/ide/src/stores/wbsStore.ts:161`）。
- `acceptance_criteria`: task の `acceptance_criteria`。
- `suggested_impl.file_paths/constraints`: `suggested_impl` から転記。

### 6.3 state/tasks.json

- 各 NodeDesign に対し少なくとも 1 つの TaskState を作成する。
- TaskState.Kind は MVP では `implementation` 固定とし、将来 `planning`/`test` を追加する。
- TaskState.NodeID が Scheduler の依存解決単位。

### 6.4 state/nodes-runtime.json

- 新規 NodeDesign 作成時に NodeRuntime を `planned` で追加する。
- TaskState が `SUCCEEDED` になったら対応 NodeRuntime.Status を `implemented` に更新する。
  - `test` Kind が追加された場合は `verified` へ更新する。

### 6.5 tasks/<task-id>.jsonl（TaskStore）

- IDE 表示用の `orchestrator.Task` を保存する既存形式を維持。
- NodeDesign/TaskState と同一の `id` を持ち、最低限 `dependencies`, `wbsLevel`, `phaseName`, `suggestedImpl`, `artifacts` を同期する。
- `ListTasks()` は IDE 表示の正規化のため、少なくとも `phaseName/milestone/wbsLevel/dependencies` を返す（`app.go:269`）。

## 7. 主要フロー

### 7.1 Chat → 計画生成

1. IDE Chat が `internal/chat/handler.go` にメッセージを渡す。
2. Handler が `Meta.PlanPatch` を呼び、PlanPatchResponse（operations）を得る。
3. Handler が operations を適用して永続化:
   - `create`: WBS/NodeDesign/NodesRuntime/TasksState を作成し、TaskStore へ append（IDE に `task:created` を emit）。
   - `update`: NodeDesign/TaskStore を更新し、必要なら依存関係を更新する。
   - `move`: WBS の親子/順序（および facet）を更新する。
   - `delete`: WBS/state から除外し、依存関係から参照を除去する（soft delete）。

### 7.2 Run → 実行

1. 自律実行ループは `StartExecution` で開始する（`app.go:472`、`internal/orchestrator/execution_orchestrator.go:80`）。
   - 基本動作: **チャット完了後に自動で開始**する（Chat Autopilot）
   - フォールバック: UI から明示的に開始/停止できる。
2. Scheduler が依存解決し、実行可能タスクを READY→enqueue する（自動: `internal/orchestrator/execution_orchestrator.go:245`、手動: `app.go:377`、`internal/orchestrator/scheduler.go:31`）。
3. ExecutionOrchestrator が 2 秒ポーリングで Job を dequeue し Executor を起動する（`internal/orchestrator/execution_orchestrator.go:190`、`internal/orchestrator/execution_orchestrator.go:256`）。
4. Executor が agent-runner に YAML を stdin 経由で渡して実行する（`internal/orchestrator/executor.go:83`、`internal/orchestrator/executor.go:157`）。

### 7.3 結果反映

1. Executor の Attempt 結果で TaskState.Status を `SUCCEEDED/FAILED` に更新。
2. SUCCEEDED の場合 NodeRuntime.Status を `implemented` へ更新。
3. TaskStore（legacy）の Task も更新し、IDE へ `task:stateChange` を emit。

（将来拡張）AgentRunner の出力（Task Note や JSON サマリ）から「生成・変更されたファイル一覧」を抽出し、`Artifacts.Files` に保存して IDE で参照できるようにする。MVP では `Artifacts.Files` が空でも許容する。

## 8. UX/性能方針（イベント駆動）

- 画面のカクつきを避けるため、状態変化系イベント（`task:created`/`task:stateChange`/`execution:stateChange`/`chat:progress`）の粒度を維持しつつ、ログ系イベント `task:log` はフロント側で最大 1000 行に制限する（`internal/orchestrator/events.go:34`、`internal/orchestrator/executor.go:121`、`frontend/ide/src/stores/logStore.ts:16`）。
- Graph/WBS の再レイアウトは Task 一覧のバッチ更新後に一度だけ行う。
- 大量タスク生成時は UI 更新をスロットリング（例: 100ms 単位）する。

## 9. MVP 完了条件

- ゴールデン入力（例: 「TODO アプリを作成して」）で、チャット → 計画 → 実行 → 結果表示がローカルで一気通しで成功する。
- 依存関係を持つタスクが、依存ノード完了後に自動で READY になり実行される。
- IDE で操作中に明確なカクつきやフリーズが起きない。

---

## 10. 反省（Post-MVP）と原因分析（再発防止の前提）

この章は「今回の実装で露呈した設計/実装/運用の欠陥」を一次ソース付きで列挙し、vNext のタスク設計に **強制的に継承**する。

### 10.1 プロトコル/実装の乖離（Meta plan_patch）

- 【事実】`plan_patch` の入力は「既存タスク要約 + 既存WBS概要 + 会話履歴」を Meta に渡す仕様（`docs/specifications/meta-protocol.md:461`）。
- 【事実】MVP 当時は `PlanPatchRequest` を作っていても、プロンプトで `existing_wbs.node_index` と `conversation_history` を落としていた（`internal/meta/utils.go:153`）。
- 【事実】現状は QH-001 により、`existing_tasks` の facet/依存・`existing_wbs.node_index`・`conversation_history` をプロンプトに含めている（`internal/meta/utils.go:155`）。
- 【結果】“差分更新” の精度を阻害する主要因は解消。ただし巨大WBS/大量タスク時の **決定論トリミング**（11.1/12.1 参照）は未完。
- 【事実】`planPatchSystemPrompt` の例に `status` があるが、`PlanOperation` に `status` は定義されていない（`internal/meta/client.go:198` と `internal/meta/protocol.go:203`）。

### 10.2 WBS 整合性の欠陥（delete/cascade の定義不足）

- 【事実】仕様は `cascade: false` を許容する（`docs/specifications/meta-protocol.md:509`）。
- 【事実】MVP 当時は `cascade=false` で子ノードの再接続（reparent）が無く、孤児が発生し得た（`internal/chat/plan_patch.go:482`、`internal/chat/plan_patch.go:511`）。
- 【事実】現状は QH-003（案A）として「子を親へ繰り上げ（splice）」を実装した（`internal/chat/plan_patch.go:511`、`internal/chat/plan_patch.go:550`）。
- 【結果】孤児リスクは大幅に低減。ただし WBS 不変条件テスト（11.3）が不足している。

### 10.3 監査/復元性（history の順序）

- 【事実】設計は「history append → design/state を atomic write」の順序を要求する（`docs/design/orchestrator-persistence-v2.md:92`）。
- 【事実】MVP 当時は design/state 保存後に history を best-effort append していた（`internal/chat/plan_patch.go:390`）。
- 【事実】現状は QH-004 として、history append を先行させてから design/state を保存している（`internal/chat/plan_patch.go:380`）。
- 【結果】設計意図に沿った順序は満たした。一方で「history append 失敗時の扱い」をエラーにするか継続するか（11.1/12.3）の設計が未確定。

### 10.4 テストの信頼性（外部ネットワーク依存/モック整合）

- 【事実】MVP 当時は `internal/meta` のテストが失敗していた（例: `internal/meta/client_test.go:47`）。
- 【事実】現状は QH-005 として NextAction プロンプトに `WorkerRuns` を含め、モック分岐と整合させた（`internal/meta/openai_provider.go:328`、`internal/meta/mock_adapter.go:69`）。
- 【事実】現状の品質ゲートとして `go test ./...` が通る。
- 【結果】品質ゲートとしては復旧。ただし “プロンプト文字列の偶然” に依存している点は残り、将来的には「構造（payload）ベースのモック」へ移行する余地がある。

### 10.5 ドキュメントの真実源の弱さ（“一次ソースで検証できる”形になっていない）

- 【事実】PRD/仕様は存在するが、実装側で満たすべき **不変条件（invariants）** と検証方法（テスト/コマンド）が PRD の DoD に明示されていない。
- 【結果】MVP 完了後に “品質の穴” が検出され、手戻りが増える。

---

## 11. Quality Hardening の到達点（Definition of Done / 品質ゲート）

この章の DoD は vNext の「完了判定の唯一の基準」とする（“動いた” だけでは完了扱いにしない）。

### 11.1 機能 DoD（仕様適合）

- 【事実】`plan_patch` は入力コンテキスト（既存タスク要約/既存WBS概要/会話履歴）を **構造を保持した形で** Meta に渡す（仕様根拠: `docs/specifications/meta-protocol.md:461`）。
- 【事実】`delete` のセマンティクスが定義され、`cascade=false` でも WBS 不変条件を壊さない。
- 【事実】history は「先に append、後で design/state 更新」を満たし、失敗時の扱い（失敗レコード/ロールバック方針）が定義されている。

### 11.2 データ不変条件（invariants）

- 【事実】WBS:
  - `root_node_id` は必ず `node_index` 内に存在する。
  - 全ノード（root を除く）は `parent_id` を持ち、親の `children` に含まれる。
  - `children` は重複しない（集合性）。
  - delete 後も “孤児ノード” が存在しない。
- 【事実】design/state/task_store:
  - `design/nodes/<id>.json`（NodeDesign）と `state/tasks.json`（TaskState）と `tasks/<id>.jsonl`（TaskStore）が “同一 TaskID” を参照し、facet（phase/milestone/wbsLevel/dependencies）が矛盾しない（`app.go:404` の join が破綻しない）。

### 11.3 テスト/運用 DoD

- 【事実】`go test ./...` が **ネットワーク不要**で安定して通る（ユニットテストは外部 API を叩かない）。
- 【事実】plan_patch の delete/move/update/create について、WBS 不変条件を検証するテストが存在する。
- 【事実】CLI セッション（Codex/Claude）の存在確認とユーザー向けエラーメッセージが整備される（`ISSUE.md:21`）。

---

## 12. 100点タスク設計（vNext: Quality Hardening）

この章は “実装に落とせる粒度” のタスクとして記述する。各タスクは **目的/範囲/受け入れ条件/検証方法/影響範囲** を必須とする。

### 12.1 P0: 仕様適合（plan_patch の入力/出力の整合）

#### QH-001: plan_patch プロンプトに構造化コンテキストを完全継承

- 【目的】Meta が差分更新を正しく行えるよう、仕様通りの入力を失わず渡す。
- 【範囲】`internal/meta/utils.go` の `buildPlanPatchUserPrompt` と関連するプロンプト生成。
- 【受け入れ条件】
  - `existing_tasks` の各要素に `dependencies/phase_name/milestone/wbs_level/parent_id` が含まれる（仕様: `docs/specifications/meta-protocol.md:465`）。
  - `existing_wbs` の `root_node_id` と `node_index` がプロンプトに含まれる（仕様: `docs/specifications/meta-protocol.md:467`）。
  - `conversation_history` がロール/本文を保って含まれる（仕様: `docs/specifications/meta-protocol.md:468`）。
  - 文字数/トークン制限対策（トリミング規則）が明文化され、テストで固定化されている。
    - 会話履歴: 最新 `N=10` 件、各本文は `max=300 chars` に丸める（現実装の decompose 側と同一規約に統一する）。
    - 既存タスク要約: `max=200` 件（超過時は「直近更新順」等の決定論規則で切る）。
    - WBS: `node_index` は全件が原則だが、超過時は “root からの部分木 + 参照される task_id 周辺” の決定論サブセットに落とす（ルールを固定しテストする）。
- 【検証】ユニットテストで “プロンプトに必要フィールドが含まれること” を固定文字列で検証。
- 【一次根拠】欠陥: `internal/meta/utils.go:153`

#### QH-002: plan_patch system prompt と protocol の整合

- 【目的】LLM が返すべき JSON スキーマを “実装が受け取れる形” に固定する。
- 【範囲】`internal/meta/client.go` の `planPatchSystemPrompt`、`internal/meta/protocol.go`。
- 【受け入れ条件】
  - system prompt から `status` の例を削除する、または `PlanOperation` に `status` を追加し適用実装まで含める（どちらかに統一）。
  - `docs/specifications/meta-protocol.md` と矛盾しない。
- 【検証】スキーマテスト（PlanPatchResponse の json unmarshal）+ 既存の chat handler テストを通す。
- 【一次根拠】欠陥: `internal/meta/client.go:198`、`internal/meta/protocol.go:203`

### 12.2 P0: WBS 整合性（delete の仕様確定 + 実装 + テスト）

#### QH-003: delete(cascade=false) のセマンティクス確定（孤児を作らない）

- 【目的】WBS を破壊しない delete を定義し、実装・テストで固定する。
- 【範囲】`internal/chat/plan_patch.go` の delete 処理、WBS 操作関数。
- 【方針（採用）】案A を採用する（UX と “差分編集” の自然さを優先）。
  - `cascade=false` は “削除ノードの子を削除ノードの親へ繰り上げ（順序維持）”。
  - 例: 親Pの children が `[... , X, ...]`、X の children が `[a,b,c]` のとき、X を削除すると P の children は `[..., a, b, c, ...]` となる（X の位置に splice）。
- 【受け入れ条件】
  - delete 後も WBS 不変条件を満たす（`11.2`）。
  - 依存関係から削除対象への参照が除去される（仕様: `docs/specifications/meta-protocol.md:553`）。
- 【検証】WBS 不変条件テスト + plan_patch 適用テスト。
- 【一次根拠】欠陥: `internal/chat/plan_patch.go:473`

### 12.3 P0: 監査/復元（history の順序と失敗時の扱い）

#### QH-004: “擬似トランザクション” を設計に合わせて実装

- 【目的】復元可能性の源泉を history に置く設計を、実装で担保する。
- 【範囲】`internal/chat/plan_patch.go`（plan_patch 適用時の永続化順序）、必要なら `internal/orchestrator/persistence/*`。
- 【受け入れ条件】
  - plan_patch 適用前に history に action が append される（設計: `docs/design/orchestrator-persistence-v2.md:92`）。
  - design/state 保存が失敗した場合の扱い（失敗 action の追記、または idempotent な再適用）が定義され、テストで再現できる。
- 【検証】失敗注入テスト（writeJSON 失敗を擬似化）で history と state の整合を確認。
- 【一次根拠】欠陥: `internal/chat/plan_patch.go:390`

### 12.4 P0: テストの無外部依存化（品質ゲートの復旧）

#### QH-005: meta の mock を “プロンプトの偶然” ではなく “構造” に合わせる

- 【目的】`go test ./...` を安定した品質ゲートに戻す。
- 【範囲】`internal/meta/mock_adapter.go`、`internal/meta/openai_provider.go`、`internal/meta/client_test.go`。
- 【受け入れ条件】
  - mock の分岐条件が “固定文字列の部分一致” ではなく、リクエスト構造（JSON/YAML payload のフィールド）に基づく。
  - NextAction のコンテキストに `worker_runs_count` 等の必要情報を含める、または mock がそれを要求しないよう統一する。
  - `go test ./...` がネットワーク無しで成功する。
- 【一次根拠】欠陥: `internal/meta/mock_adapter.go:72`、`internal/meta/openai_provider.go:334`、失敗: `internal/meta/client_test.go:47`

### 12.5 P1: 実行ログ/セッション（運用の完成度）

#### QH-006: 実行ログ UI を運用可能にする

- 【目的】実行中の状況が “追える/切り替えられる/クリアできる” 状態にする。
- 【範囲】フロントの `task:log` 表示（`ISSUE.md:15`）。
- 【受け入れ条件】タスク別フィルタ、クリア、常時表示（最小導線）が存在する。

#### QH-007: Codex CLI セッション検証と注入仕様の確定

- 【目的】実行環境の失敗を “事前に検知” し、ユーザーに正しく伝える。
- 【範囲】Worker 起動時のセッション検証とドキュメント化（`ISSUE.md:21`）。
- 【受け入れ条件】セッション未検出時に UI/ログへ明確なガイダンスが出る。

### 12.6 P2: 将来拡張（“今回の反省” 由来で仕様/テストを先に固める）

#### QH-008: Artifacts.Files の自動抽出/反映（実行結果の追跡可能性）

- 【目的】「どのタスクがどのファイルを生成/変更したか」を IDE の一次情報として扱える状態にする。
- 【範囲】`ISSUE.md:38` の実装タスク群（executor/note/persistence/UI）。
- 【受け入れ条件】タスク完了後にファイル一覧が IDE で確認できる。

#### QH-009: Meta Protocol のバージョニング導入（将来互換）

- 【目的】プロトコル更新時の破壊的変更を検知し、安全にフォールバックする。
- 【範囲】`ISSUE.md:52`（spec 更新 + core 実装 + テスト）。
- 【受け入れ条件】バージョン不一致時の挙動が仕様/実装/テストで一致している。

#### QH-010: 追加 Worker 種別のサポート（設計先行）

- 【目的】`codex-cli` 以外の CLI エージェントを Worker として選択可能にする。
- 【範囲】`ISSUE.md:64`。
- 【受け入れ条件】`runner.worker.kind` に応じた選択/エラー/ドキュメントが揃う。

#### QH-011: IPC の WebSocket / gRPC 化（性能ボトルネック解消）

- 【目的】ファイルポーリング IPC の負荷/拡張性制約を解消する。
- 【範囲】`ISSUE.md:78`。
- 【受け入れ条件】移行パス（並行稼働）と大規模ジョブ時の負荷低減が確認できる。

#### QH-012: Frontend E2E の安定化（品質ゲートの2本目）

- 【目的】フロントの回帰を CI で継続検出できる状態にする。
- 【範囲】`ISSUE.md:90`。
- 【受け入れ条件】`pnpm test:e2e` が安定して完走する。

#### QH-013: Task Note 保存の圧縮（ストレージ/性能）

- 【目的】大きな履歴の保存サイズを抑え、読み書き性能を維持する。
- 【範囲】`ISSUE.md:102`。
- 【受け入れ条件】後方互換を維持しつつ保存容量が有意に減る。
