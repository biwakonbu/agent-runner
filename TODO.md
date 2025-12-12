# TODO: Pragmatic MVP 実装手順

最終更新: 2025-12-12

## 0. 前提

- MVP のゴールは `PRD.md` の「MVP 完了条件」参照。
- リアルタイム UX（ログストリーミング等）は後回し。UI のカクつき防止と体感性能を優先する。

## 1. 計画生成の橋渡し（Chat → design/state）

### 1.1 Decompose → WBS/NodeDesign 永続化（DONE）

【目的】Scheduler が参照する `design/` を、チャット入力から生成できる状態にする。

- 対象ファイル:
  - `internal/chat/handler.go`
  - `internal/orchestrator/persistence/repo.go`
  - `internal/orchestrator/persistence/models.go`
- 実装タスク:
  1. Handler に `persistence.WorkspaceRepository` を注入し、チャット処理から design/state へ書き込めるようにする。
  2. `decomposeResp` から WBS ルートを作成/更新し `design/wbs.json` に保存する。
  3. 各 `DecomposedTask` を 1:1 で `NodeDesign` に写像し `design/nodes/<node-id>.json` に保存する。
  4. `DecomposedTask.dependencies` を `node_id` に解決し、`NodeDesign.Dependencies` に格納する。
- 完了条件:
  - `~/.multiverse/workspaces/<id>/design/wbs.json` と `design/nodes/*.json` が生成される。

### 1.2 Decompose → NodesRuntime/TasksState 永続化（DONE）

【目的】ExecutionOrchestrator が読む `state/` を plan と同期させ、実行可能にする。

- 対象ファイル:
  - `internal/chat/handler.go`
  - `internal/orchestrator/persistence/repo.go`
- 実装タスク:
  1. NodeDesign 作成時に `state/nodes-runtime.json` に `NodeRuntime{status:"planned"}` を追加（既存なら更新）。
  2. 各 NodeDesign に対応する `TaskState{kind:"implementation", status:"pending"}` を `state/tasks.json` に追加。
  3. `TaskState.NodeID` と `NodeDesign.NodeID` を一致させる。
- 完了条件:
  - `state/nodes-runtime.json` と `state/tasks.json` に新規エントリが作成される。

### 1.3 TaskStore との同期（DONE）

【目的】IDE 表示用の TaskStore と design/state の整合を取る。

- 対象ファイル:
  - `internal/orchestrator/task_store.go`
  - `internal/chat/handler.go`
- 実装タスク:
  1. TaskStore の Task.ID を NodeDesign/TaskState と同一にする（既存 UUID マッピングを整理）。
  2. `dependencies / wbsLevel / phaseName / suggestedImpl / acceptanceCriteria` を同期して保存する。
- 完了条件:
  - IDE の Task 一覧/グラフが従来どおり表示できる。

## 2. 実行結果の反映（state/design の整合）

### 2.1 AttemptCount / Retry の整理（DONE）

【目的】RetryPolicy が確実に動くよう、試行回数を永続化する。

- 対象ファイル:
  - `internal/orchestrator/execution_orchestrator.go`
  - `internal/orchestrator/persistence/models.go`
- 実装タスク:
  1. `TaskState` に `AttemptCount` フィールドを追加するか、`Inputs["attempt_count"]` を正式仕様として扱う。
  2. `processJob` 開始時に attempt_count をインクリメントして `state/tasks.json` に保存する。
- 完了条件:
  - 連続失敗時に backoff が段階的に伸びる。

### 2.2 Task 成功時の NodesRuntime 更新（DONE）

【目的】ノード依存が解決され、後続タスクが READY になるようにする。

- 対象ファイル:
  - `internal/orchestrator/execution_orchestrator.go`
- 実装タスク:
  1. `attempt.Status == SUCCEEDED` の場合、該当 `node_id` の `NodeRuntime.Status` を `implemented` に更新（無ければ作成）。
  2. 更新後に `state/nodes-runtime.json` を保存する。
- 完了条件:
  - 依存ノード完了後、次ノードのタスクが自動で READY になり実行される。

### 2.3 TaskStore / IDE イベント反映（DONE）

【目的】IDE 表示と実行状態を同期する。

- 対象ファイル:
  - `internal/orchestrator/executor.go`
  - `internal/orchestrator/execution_orchestrator.go`
- 実装タスク:
  1. TaskStore の Task.Status と Artifacts を更新する。
  2. `EventTaskStateChange` など既存イベントで IDE に通知する。
- 完了条件:
  - IDE に SUCCEEDED/FAILED が反映される。

## 3. Executor YAML の最小改善（ハードコード排除）

### 3.1 max_loops / worker kind の受け渡し（DONE）

【目的】MVP でも最低限 Runner 設定を差し替えられるようにする。

- 対象ファイル:
  - `internal/orchestrator/executor.go`
- 実装タスク:
  1. `TaskState.Inputs` などから `runner.meta.max_loops`/`runner.worker.kind` を読んで YAML に反映する。
  2. 値が無い場合は現状デフォルトにフォールバックする。
- 完了条件:
  - 設定変更が破壊的変更なしに反映される。

## 4. IDE の性能チューニング（必要なら）

### 4.1 グラフ再レイアウトのバッチ化

【目的】大量タスク生成・状態変化時のカクつきを抑える。

- 対象ファイル:
  - `frontend/ide/src/stores/*`
  - `frontend/ide/src/lib/flow/UnifiedFlowCanvas.svelte`
- 実装タスク:
  1. stateChange イベントの連続受信時に layout 計算をまとめる。
- 完了条件:
  - タスク大量生成時に明確な UI の遅延が出ない。

## 5. ゴールデンパス検証

### 5.1 自動ゴールデンパス（DONE）

【目的】IDE の入力をモックし、バックエンドの「Chat→計画生成→依存解決→実行→同期」までを自動で検証する。

- 実装:
  - `test/e2e/golden_pass_test.go`
- 検証内容:
  - ChatHandler が `design/` と `state/` を生成すること。
  - ExecutionOrchestrator が依存順にタスクを実行し、`state/` と TaskStore を更新すること。

### 5.2 手動一気通し（必要なら）

- 手順:
  1. IDE 起動 → ワークスペース選択
  2. チャットで簡単な要求を入力
  3. 生成されたタスク（またはノード）を Run
  4. 完了ステータスと成果物を確認
- 観測ポイント:
  - `design/`・`state/`・`tasks/` の 3 層が整合する。
  - 依存順に実行される。

### 5.3 テスト追加（既存パターンに沿う）

- 対象ファイル:
  - `internal/chat/handler_test.go`
  - `internal/orchestrator/*_test.go`
- 追加タスク:
  1. `decompose → design/state 保存` の単体テスト。
  2. `processJob` が `NodesRuntime` を更新するテスト。

## 6. 将来拡張

### 6.1 Artifacts.Files の自動抽出/反映

【目的】実行したタスクが「どのファイルを生成/変更したか」を IDE で追跡できるようにする。

- 対象ファイル（候補）:
  - `internal/note/writer.go`
  - `internal/orchestrator/executor.go`
  - `internal/orchestrator/execution_orchestrator.go`
- 実装タスク:
  1. AgentRunner の Task Note/JSON 出力から変更・生成ファイルを抽出する仕組みを定義。
  2. 抽出結果を `Artifacts.Files` に保存し、TaskStore と state を同期。
  3. IDE の TaskPropPanel で一覧表示。
- 完了条件:
  - タスク完了後にファイル一覧が確認できる。

### 6.2 Meta Protocol のバージョニング導入

【目的】Meta-agent と Core 間のプロトコル互換性を将来にわたって維持する。

- 対象ファイル（候補）:
  - `internal/core/meta/*`
  - `docs/specifications/meta-protocol.md`
- 実装タスク:
  1. YAML メッセージに `protocol_version`（または同等）を追加し、Core 側で解釈する。
  2. バージョン不一致時のフォールバック/警告/拒否方針を定義する。
- 完了条件:
  - プロトコル更新時に旧クライアントが安全に扱える。

### 6.3 追加 Worker 種別のサポート

【目的】`codex-cli` 以外の CLI エージェントを Worker として選択可能にする。

- 対象ファイル（候補）:
  - `internal/orchestrator/executor*.go`
  - `internal/worker/*`
  - `docs/cli-agents/*`
- 実装タスク:
  1. Worker kind と Docker イメージ/起動コマンドの対応表を追加する。
  2. `runner.worker.kind` に応じた選択とエラーハンドリングを実装する。
  3. 各 CLI のナレッジ（`docs/cli-agents/<kind>/...`）とテストを追加する。
- 完了条件:
  - `gemini-cli` / `claude-code-cli` / `cursor-cli` を指定して実行できる。

### 6.4 IPC の WebSocket / gRPC 化

【目的】ファイルポーリング IPC の性能/拡張性の制約を解消する。

- 対象ファイル（候補）:
  - `internal/orchestrator/ipc/*`
  - `frontend/ide/src/*`
- 実装タスク:
  1. Queue/イベント通知を WebSocket か gRPC に置き換える設計を確定する。
  2. 既存 file-based IPC と並行稼働できる移行パスを用意する。
- 完了条件:
  - 大量ジョブ時のポーリング負荷が解消される。

### 6.5 Frontend E2E の安定化

【目的】IDE フロントの E2E が CI で継続的に回る状態にする。

- 対象ファイル（候補）:
  - `frontend/ide/*`
  - `docs/guides/testing.md`
- 実装タスク:
  1. タイムアウト/待機条件/テストデータを見直し安定化する。
  2. 失敗時ログの拡充とリトライ方針を整備する。
- 完了条件:
  - `pnpm test:e2e` が安定して完走する。

### 6.6 Task Note 保存の圧縮

【目的】大きな Task Note/履歴の保存サイズを抑え、読み書き性能を維持する。

- 対象ファイル（候補）:
  - `internal/note/*`
  - `internal/orchestrator/persistence/*`
- 実装タスク:
  1. Task Note の圧縮形式（gzip 等）と保存/読み込み API を定義する。
  2. 既存データとの後方互換を確保する。
- 完了条件:
  - Task Note の保存容量が有意に削減される。
