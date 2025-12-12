# PRD: Multiverse IDE Pragmatic MVP（Chat→WBS/Node→AgentRunner 実行）

最終更新: 2025-12-12

## 1. 背景

- AgentRunner Core は `plan_task`/`next_action`/`completion_assessment` と Docker Sandbox を含む実行ループが安定稼働している（`docs/CURRENT_STATUS.md:18`、`docs/design/data-flow.md:13`）。
- Orchestrator は file-based IPC + Scheduler + Executor を備えるが Beta 段階で、WBS/Node 中心の永続化（design/state/history）を前提とした v2 実装が途中（`docs/CURRENT_STATUS.md:20`、`docs/design/orchestrator-persistence-v2.md:11`）。
- 現在の Chat は `decompose` を用いてタスク分解し TaskStore に Task を保存するが、design/state には反映していないため、Scheduler が依存解決できず実行に進まない（`internal/chat/handler.go:128`、`internal/orchestrator/scheduler.go:89`）。

## 2. 目的 / ゴール

MVP の到達点は「IDE のチャット入力から、WBS/ノード計画を生成・永続化し、その計画に基づいて Orchestrator が AgentRunner を起動してタスクを順次完了させ、IDE 上に結果が表示される」こと。

具体的には:

1. チャット入力 → Meta-agent `decompose` → WBS/Node/TaskState の生成と永続化が行われる。
2. `ExecutionOrchestrator` が依存関係を解決し、READY タスクを IPC Queue に流し、`agent-runner` を実行できる。
3. 実行結果で TaskState / NodesRuntime / TaskStore が更新され、IDE が一覧/グラフ表示できる。

## 3. 非ゴール（MVPでは扱わない）

- ログのリアルタイムストリーミング（WebSocket/gRPC などの IPC 強化）。
- マルチノード/リモート Worker プール。
- 高度な承認フローや差分レビュー UI。
- アニメーションや高度な UI エフェクト。UI は「カクつかず安定して操作できる」ことを優先する。

## 4. ユーザーストーリー

- US-1: 開発者は IDE のチャットに要望を入力し、数秒〜数十秒後に WBS/ノードとタスクリストが生成される。
- US-2: 開発者は「Run」操作で計画全体または特定ノードを実行できる。
- US-3: IDE 上で各タスク/ノードのステータス（PENDING/READY/RUNNING/SUCCEEDED/FAILED/BLOCKED）が確認でき、生成・更新されたファイル一覧を参照できる。

## 5. アーキテクチャ方針

### 5.1 計画と実行の真実源

- 計画（WBS/NodeDesign）は `~/.multiverse/workspaces/<id>/design/` を真実源とする（`docs/design/orchestrator-persistence-v2.md:33`）。
- 実行状態（TasksState/NodesRuntime/AgentsState）は `state/` を真実源とする。
- `internal/orchestrator/task_store.go` の TaskStore は IDE 表示と後方互換のため当面併用し、design/state と同期させる。

### 5.2 Planner/TaskBuilder の配置

MVP では **Chat Handler が Planner/TaskBuilder の役割を兼務**する。

- `decompose` 呼び出しは現状のまま Chat Handler が行う。
- `decompose` 結果を design/state/task_store に写像して永続化する。

将来的には Planner を Orchestrator 側に移し、Chat は UI 層へ戻す。

## 6. データモデル（MVPスキーマ）

### 6.1 design/wbs.json

- WBS ルートのみ保持。最低限 `wbs_id`, `project_root`, `root_node_id`, `node_index` を保存する（`internal/orchestrator/persistence/models.go:9`）。

### 6.2 design/nodes/<node-id>.json

- `decompose.phases[].tasks[]` を 1:1 で NodeDesign として保存する。
- NodeDesign.Dependencies は `decompose` の `dependencies` を `node_id` に解決したものを格納する。

主要フィールド:

- `node_id`: UUID または `node-<task-id>` 形式。
- `name`, `summary`: decompose task の `title`/`description`。
- `acceptance_criteria`: decompose task の `acceptance_criteria`。
- `suggested_impl.file_paths/constraints`: decompose の `suggested_impl` から転記。

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

## 7. 主要フロー

### 7.1 Chat → 計画生成

1. IDE Chat が `internal/chat/handler.go` にメッセージを渡す。
2. Handler が `Meta.Decompose` を呼び、DecomposeResponse を得る（現状実装）。
3. Handler が DecomposeResponse を永続化:
   - WBS ルート作成/更新。
   - NodeDesign 作成/更新。
   - NodesRuntime へ `planned` を追加。
   - TasksState へ `pending` の TaskState を追加。
   - TaskStore へ Task を append し、IDE に `task.created` イベントを emit。

### 7.2 Run → 実行

1. IDE が Task の Run を押し、Scheduler が TaskState を READY にし IPC queue に Job を enqueue（`internal/orchestrator/scheduler.go:31`）。
2. ExecutionOrchestrator が 2 秒ポーリングで Job を dequeue し Executor を起動する（`internal/orchestrator/execution_orchestrator.go:190`）。
3. Executor が agent-runner に YAML を渡して実行する（`internal/orchestrator/executor.go:64`）。

### 7.3 結果反映

1. Executor の Attempt 結果で TaskState.Status を `SUCCEEDED/FAILED` に更新。
2. SUCCEEDED の場合 NodeRuntime.Status を `implemented` へ更新。
3. TaskStore の Task も同様に更新し、IDE へ stateChange/artifacts 更新イベントを emit。

（将来拡張）AgentRunner の出力（Task Note や JSON サマリ）から「生成・変更されたファイル一覧」を抽出し、`Artifacts.Files` に保存して IDE で参照できるようにする。MVP では `Artifacts.Files` が空でも許容する。

## 8. UX/性能方針（非リアルタイム前提）

- 画面のカクつきを避けるため、バックエンドからフロントへのイベントは「状態変化単位」に限定し、ログの行単位配信は MVP では行わない。
- Graph/WBS の再レイアウトは Task一覧のバッチ更新後に一度だけ行う。
- 大量タスク生成時は UI 更新をスロットリング（例: 100ms 単位）する。

## 9. MVP 完了条件

- ゴールデン入力（例: 「TODO アプリを作成して」）で、チャット→計画→実行→結果表示がローカルで一気通しで成功する。
- 依存関係を持つタスクが、依存ノード完了後に自動で READY になり実行される。
- IDE で操作中に明確なカクつきやフリーズが起きない。
