# TODO: Multiverse IDE Core 永続化 & スケジューラ V2

## 1. データモデル & リポジトリ層 (新規 `internal/orchestrator/persistence`)

- [x] **データモデル定義 (`internal/orchestrator/persistence/models.go`)**

  - [x] `WBS` および `NodeIndex` 構造体
  - [x] `NodeDesign` 構造体 (estimate, suggested_impl 含む)
  - [x] `NodesRuntime` および `NodeRuntime` 構造体
  - [x] `TasksState`, `TaskState` (v2), `QueueMeta` 構造体
  - [x] `AgentsState` および `AgentState` 構造体
  - [x] `Action` 構造体 (共通フィールド + Payload)

- [x] **リポジトリ実装 (`internal/orchestrator/persistence/repo.go`)**

  - [x] `DesignRepository`: `LoadWBS`, `SaveWBS`, `GetNode`, `SaveNode`
  - [x] `StateRepository`: `LoadTasks`, `SaveTasks`, `LoadNodesRuntime`, `SaveNodesRuntime`
  - [x] `HistoryRepository`: `AppendAction`, `ListActions`
  - [x] `WorkspaceRepository`: 全リポジトリのマネージャ、ディレクトリ初期化処理

- [x] **データアクセス・テスト**
  - [x] JSON シリアライズ/デシリアライズのユニットテスト
  - [x] ファイル I/O と原子的更新の統合テスト

## 2. スケジューラロジック V2 (`internal/orchestrator/scheduler_v2.go`)

- [x] **スケジューラのリファクタリング/書き換え**

  - [x] `TaskStore` の代わりに `StateRepository` を使用するよう `Scheduler` 構造体を更新
  - [x] 以下の処理を行う `ScheduleLoop` (または `Step`) を実装:
    1. `state/tasks.json` から保留中(pending)タスクを読み込む
    2. `state/nodes-runtime.json` と `design/` を使用して依存関係をチェック
    3. `state/agents.json` でエージェントの空き状況をチェック
    4. `task.started` アクションを作成
    5. `state/tasks.json` (status=running) と `state/agents.json` を更新
  - [x] `allDependenciesSatisfied` ロジックを新データモデルへ移植

- [x] **アクション & イベント統合**
  - [x] すべての状態変更に先立って `HistoryRepository.AppendAction` が実行されることを保証
  - [ ] UI 通知用の内部イベントを発火 (後で IPC 経由で送信)

## 3. Executor & IPC との統合

- [ ] **Executor 更新**

  - [x] `TaskState` (v2) オブジェクトを受け取るよう `Executor` を修正 (v2 実装済み)
  - [x] 実行完了時、`task.succeeded` または `task.failed` アクションを記録
  - [x] 結果に基づいて `state/tasks.json` と `state/nodes-runtime.json` を更新

- [ ] **IPC Manager 更新**
  - [ ] 生ファイルを書き込むのではなく、アクション (例: `node.created`, `task.created`) を作成するよう IPC ハンドラを適合させる

## 4. 検証 & テスト

- [x] **ユニットテスト**
  - [x] `persistence` パッケージのテスト
  - [x] `scheduler` ロジックのテスト (リポジトリはモック化)
- [ ] **統合テスト**
  - [x] シミュレーションテスト: WBS 作成 -> ノード計画 -> タスクスケジュール -> 実行 -> 状態 & 履歴の検証
