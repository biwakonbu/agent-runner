# TODO: multiverse v2.0 Implementation

Based on PRD v2.0

---

## 進捗サマリ

| Phase | Status | 備考 |
|-------|--------|------|
| Phase 1: チャット→タスク生成 | 🟡 進行中 | Week 1 完了、Week 2 作業中 |
| Phase 2: 依存グラフ・WBS表示 | ⚪ 未着手 | Phase 1 完了後 |
| Phase 3: 自律実行ループ | ⚪ 未着手 | Phase 2 完了後 |

---

## Phase 1: チャット → タスク生成（MVP）

### Week 1: バックエンド実装

#### 1.1 Task 構造体拡張

- [x] `internal/orchestrator/task_store.go`
  - [x] `Description string` フィールド追加
  - [x] `Dependencies []string` フィールド追加
  - [x] `ParentID *string` フィールド追加
  - [x] `WBSLevel int` フィールド追加
  - [x] `PhaseName string` フィールド追加
  - [x] `SourceChatID *string` フィールド追加
  - [x] `AcceptanceCriteria []string` フィールド追加

#### 1.2 Meta-agent decompose プロトコル

- [x] `internal/meta/protocol.go`
  - [x] `DecomposeRequest` 構造体追加
  - [x] `DecomposeResponse` 構造体追加
  - [x] `DecomposedTask` 構造体追加
  - [x] `DecomposedPhase` 構造体追加
- [x] `internal/meta/client.go`
  - [x] `Decompose(ctx, request)` メソッド追加
  - [x] decompose 用システムプロンプト定義

#### 1.3 ChatHandler 実装

- [x] `internal/chat/handler.go` (新規)
  - [x] `ChatHandler` 構造体
  - [x] `HandleMessage()` メソッド
  - [x] Meta-agent 呼び出しロジック
  - [x] タスク生成・保存ロジック
- [x] `internal/chat/session_store.go` (新規)
  - [x] `ChatSession` 構造体
  - [x] `ChatMessage` 構造体
  - [x] JSONL 永続化
- [x] `internal/chat/CLAUDE.md` (新規)

#### 1.4 IDE バックエンド API

- [x] `cmd/multiverse-ide/app.go`
  - [x] `SendChatMessage(sessionID, message string) (*ChatResponse, error)`
  - [x] `GetChatHistory(sessionID string) ([]ChatMessage, error)`
  - [x] `CreateChatSession() (string, error)`
  - [x] ChatHandler 初期化

### Week 2: フロントエンド連携

#### 2.1 チャットUI連携

- [ ] `frontend/ide/src/lib/components/chat/FloatingChatWindow.svelte`
  - [ ] Wails API 呼び出し（SendChatMessage）
  - [ ] 応答メッセージの表示
  - [ ] タスク生成結果のインライン表示
- [ ] `frontend/ide/src/stores/chat.ts`
  - [ ] セッション管理
  - [ ] メッセージ履歴管理
  - [ ] Wails API 連携

#### 2.2 タスク表示更新

- [ ] `frontend/ide/src/stores/taskStore.ts`
  - [ ] 新規タスク追加時の状態更新
  - [ ] 依存関係情報の保持
- [ ] `frontend/ide/src/lib/grid/GridNode.svelte`
  - [ ] フェーズ別色分け（概念設計/実装設計/実装）

#### 2.3 テスト

- [ ] ChatHandler ユニットテスト
- [ ] Meta-agent decompose モックテスト
- [ ] E2E テスト（チャット→タスク生成フロー）

---

## Phase 2: 依存関係グラフ・WBS表示

### Week 3: グラフ管理

#### 3.1 TaskGraphManager

- [ ] `internal/orchestrator/task_graph.go` (新規)
  - [ ] `TaskGraphManager` 構造体
  - [ ] `TaskGraph` 構造体
  - [ ] `GraphNode` 構造体
  - [ ] `TaskEdge` 構造体
  - [ ] `BuildGraph()` メソッド
  - [ ] `GetExecutionOrder()` メソッド（トポロジカルソート）
  - [ ] `GetBlockedTasks()` メソッド
  - [ ] サイクル検出ロジック

#### 3.2 Scheduler 拡張

- [ ] `internal/orchestrator/scheduler.go`
  - [ ] `ScheduleReadyTasks()` メソッド
  - [ ] `allDependenciesSatisfied()` メソッド
  - [ ] BLOCKED 状態の自動設定

#### 3.3 ConnectionLine コンポーネント

- [ ] `frontend/ide/src/lib/grid/ConnectionLine.svelte` (新規)
  - [ ] SVG パス計算
  - [ ] 依存状態による色分け
  - [ ] 矢印マーカー
- [ ] `frontend/ide/src/lib/grid/GridCanvas.svelte`
  - [ ] ConnectionLine のレンダリング

### Week 4: WBS・視覚化

#### 4.1 WBS ビュー

- [ ] `frontend/ide/src/lib/wbs/WBSView.svelte` (新規)
  - [ ] ツリー構造表示
  - [ ] 折りたたみ/展開
  - [ ] マイルストーン表示
- [ ] `frontend/ide/src/lib/wbs/WBSNode.svelte` (新規)
- [ ] `frontend/ide/src/stores/wbsStore.ts` (新規)

#### 4.2 進捗率表示

- [ ] `frontend/ide/src/lib/toolbar/Toolbar.svelte`
  - [ ] 進捗率バー
  - [ ] Graph/WBS 切り替えボタン

---

## Phase 3: 自律実行ループ

### Week 5: 実行オーケストレーション

#### 5.1 ExecutionOrchestrator

- [ ] `internal/orchestrator/executor.go` (拡張)
  - [ ] `ExecutionOrchestrator` 構造体
  - [ ] `ExecutionState` 定義
  - [ ] `Start()` メソッド
  - [ ] `Pause()` メソッド
  - [ ] `Resume()` メソッド
  - [ ] 依存順実行ループ

#### 5.2 リアルタイム通知

- [ ] Wails Events 設定
  - [ ] `task:stateChange` イベント
  - [ ] `execution:stateChange` イベント
- [ ] フロントエンド Events リスナー

#### 5.3 一時停止UI

- [ ] `frontend/ide/src/lib/toolbar/Toolbar.svelte`
  - [ ] 一時停止ボタン
  - [ ] 再開ボタン
  - [ ] 実行状態表示

### Week 6: エラーハンドリング

#### 6.1 自動リトライ

- [ ] `internal/orchestrator/executor.go`
  - [ ] `RetryPolicy` 構造体
  - [ ] `HandleFailure()` メソッド
  - [ ] バックオフロジック

#### 6.2 バックログ管理

- [ ] `internal/orchestrator/backlog.go` (新規)
  - [ ] `BacklogItem` 構造体
  - [ ] `BacklogStore` 構造体
  - [ ] JSONL 永続化

#### 6.3 バックログUI

- [ ] `frontend/ide/src/lib/backlog/BacklogPanel.svelte` (新規)
- [ ] `frontend/ide/src/stores/backlogStore.ts` (新規)

---

## 実装済みファイル一覧

### Phase 1 で作成予定

| ファイル | 種別 | 説明 |
|---------|------|------|
| `internal/chat/handler.go` | 新規 | ChatHandler |
| `internal/chat/session_store.go` | 新規 | ChatSession 永続化 |
| `internal/chat/CLAUDE.md` | 新規 | パッケージドキュメント |

### Phase 2 で作成予定

| ファイル | 種別 | 説明 |
|---------|------|------|
| `internal/orchestrator/task_graph.go` | 新規 | TaskGraphManager |
| `frontend/ide/src/lib/grid/ConnectionLine.svelte` | 新規 | 依存矢印 |
| `frontend/ide/src/lib/wbs/WBSView.svelte` | 新規 | WBS ビュー |
| `frontend/ide/src/lib/wbs/WBSNode.svelte` | 新規 | WBS ノード |
| `frontend/ide/src/stores/wbsStore.ts` | 新規 | WBS 状態管理 |

### Phase 3 で作成予定

| ファイル | 種別 | 説明 |
|---------|------|------|
| `internal/orchestrator/backlog.go` | 新規 | BacklogStore |
| `frontend/ide/src/lib/backlog/BacklogPanel.svelte` | 新規 | バックログ UI |
| `frontend/ide/src/stores/backlogStore.ts` | 新規 | バックログ状態管理 |

---

## 次のアクション

1. **Phase 1 Week 1** から開始
2. まず `internal/orchestrator/task_store.go` の Task 構造体を拡張
3. 次に `internal/meta/protocol.go` に decompose プロトコルを追加
