# PRD v2.0: multiverse - チャット駆動AI開発支援プラットフォーム

## 1. プロダクトビジョン

### 1.1 ビジョンステートメント

**multiverse** は、チャットインターフェースを通じて開発者の意図を理解し、
Meta-agentが自律的にタスクを分解・実行・評価する AI 開発支援プラットフォームです。

**コアコンセプト:**
- チャットウィンドウが全ての入力経路（AIとの対話）
- Meta-agentによる徹底的なタスク分解
  - 概念設計 → 実装設計 → 実装計画 → タスクマネジメント → アサイン
- 2D俯瞰UIでタスクグラフを視覚化（有向グラフ）
- WBSはリリースマイルストーンとして別枠管理
- 自律実行（計画→実行まで全自動、一時停止機能あり）

### 1.2 解決する課題

| 現状の課題 | multiverse v2.0 での解決 |
|-----------|-------------------------|
| タスク作成が手動・煩雑 | チャットから自然言語でタスク生成 |
| タスク間依存関係の管理が困難 | 有向グラフで依存関係を可視化 |
| 達成判定が曖昧 | 細分化されたタスクで個別・シンプルな達成判定 |
| 人間の介入が頻繁に必要 | 自律実行ループで人間待ち不要 |
| 問題・検討材料の散逸 | バックログで一元管理 |

### 1.3 ターゲットユーザー

- ソフトウェア開発者（個人・チーム）
- AIアシスタントと協調して開発を進めたいエンジニア
- 複数の並行タスクを俯瞰的に管理したい開発リーダー

---

## 2. 機能要件（フェーズ別）

### Phase 1: チャット → タスク生成（MVP）【優先度: 最高】

#### FR-P1-001: チャット入力UI

- 既存 FloatingChatWindow を拡張
- テキスト入力・送信
- メッセージ履歴表示（user/assistant/system）
- Wails IPC 経由でバックエンドと通信
- タスク生成結果のインライン表示

#### FR-P1-002: ChatHandler（バックエンド新規）

```go
// internal/chat/handler.go
type ChatHandler struct {
    Meta          MetaClient
    TaskStore     *orchestrator.TaskStore
    SessionStore  *ChatSessionStore
}

func (h *ChatHandler) HandleMessage(ctx context.Context, sessionID, message string) (*ChatResponse, error)
```

処理フロー:
1. ユーザーメッセージを ChatSession に保存
2. Meta-agent の `decompose` を呼び出し
3. 生成されたタスクを TaskStore に永続化
4. フロントエンドに結果を返却

#### FR-P1-003: Meta-agent decompose プロトコル

リクエスト:
```yaml
type: decompose
version: 1
payload:
  user_input: "認証機能を実装してほしい"
  context:
    workspace_path: "/path/to/project"
    existing_tasks: [...]
    conversation_history: [...]
```

レスポンス:
```yaml
type: decompose
version: 1
payload:
  understanding: "認証機能の実装を要求..."
  phases:
    - name: "概念設計"
      milestone: "M1-Auth-Design"
      tasks:
        - id: "task-001"
          title: "認証フロー設計"
          description: "..."
          acceptance_criteria: [...]
          dependencies: []
          wbs_level: 1
    - name: "実装設計"
      tasks: [...]
    - name: "実装"
      tasks: [...]
  potential_conflicts:
    - file: "src/auth/login.ts"
      tasks: ["task-004"]
      warning: "既存ファイルを変更"
```

#### FR-P1-004: Task 構造体拡張

```go
type Task struct {
    // 既存
    ID        string     `json:"id"`
    Title     string     `json:"title"`
    Status    TaskStatus `json:"status"`
    PoolID    string     `json:"poolId"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
    StartedAt *time.Time `json:"startedAt,omitempty"`
    DoneAt    *time.Time `json:"doneAt,omitempty"`

    // 新規
    Description        string   `json:"description,omitempty"`
    Dependencies       []string `json:"dependencies,omitempty"`
    ParentID           *string  `json:"parentId,omitempty"`
    WBSLevel           int      `json:"wbsLevel,omitempty"`
    PhaseName          string   `json:"phaseName,omitempty"`
    SourceChatID       *string  `json:"sourceChatId,omitempty"`
    AcceptanceCriteria []string `json:"acceptanceCriteria,omitempty"`
}
```

#### FR-P1-005: ノード表示（GridCanvas拡張）

- 新規タスク生成時のアニメーション
- 依存関係インジケーター（Phase 2 準備）
- フェーズ（概念設計/実装設計/実装）の色分け

---

### Phase 2: 依存関係グラフ・WBS表示【優先度: 高】

#### FR-P2-001: TaskGraphManager

```go
// internal/orchestrator/task_graph.go
type TaskGraphManager struct {
    TaskStore *TaskStore
}

type TaskGraph struct {
    Nodes map[string]*GraphNode
    Edges []TaskEdge
}

func (m *TaskGraphManager) BuildGraph() (*TaskGraph, error)
func (m *TaskGraphManager) GetExecutionOrder() ([]string, error)
func (m *TaskGraphManager) GetBlockedTasks() ([]string, error)
```

#### FR-P2-002: ConnectionLine（依存矢印）

```svelte
<!-- frontend/ide/src/lib/grid/ConnectionLine.svelte -->
<svg class="connection-line">
  <path d={calculatePath(fromNode, toNode)} class={status} />
  <marker id="arrowhead" ... />
</svg>
```

視覚表現:
- 完了した依存: 緑色の実線
- 未完了の依存: オレンジの破線
- ブロック状態: 赤色の太線

#### FR-P2-003: WBS表示モード

- ツールバーに WBS/Graph 切り替えボタン
- WBS ビュー: マイルストーン別のツリー表示
- 折りたたみ/展開機能
- 進捗率表示（完了タスク / 全タスク）

#### FR-P2-004: 依存に基づくスケジューリング

```go
func (s *Scheduler) ScheduleReadyTasks() error {
    for _, task := range s.GetPendingTasks() {
        if s.allDependenciesSatisfied(task) {
            s.ScheduleTask(task.ID)
        }
    }
}
```

---

### Phase 3: 自律実行ループ【優先度: 中】

#### FR-P3-001: ExecutionOrchestrator

```go
type ExecutionOrchestrator struct {
    Scheduler    *Scheduler
    GraphManager *TaskGraphManager
    State        ExecutionState  // IDLE | RUNNING | PAUSED | STOPPED
    PauseSignal  chan struct{}
    ResumeSignal chan struct{}
}

func (e *ExecutionOrchestrator) Start(ctx context.Context) error
func (e *ExecutionOrchestrator) Pause()
func (e *ExecutionOrchestrator) Resume()
```

#### FR-P3-002: リアルタイム進捗表示

```go
// バックエンド
runtime.EventsEmit(ctx, "task:stateChange", TaskStateChangeEvent{...})
```

```typescript
// フロントエンド
runtime.EventsOn('task:stateChange', (event) => {
    tasks.updateTask(event.taskId, { status: event.newStatus });
});
```

#### FR-P3-003: 一時停止・再開機能

- ツールバーに一時停止/再開ボタン
- 一時停止時は実行中タスクを中断せず、新規タスク開始のみ停止

#### FR-P3-004: 自動リトライ/人間判断

```go
type RetryPolicy struct {
    MaxAttempts     int
    BackoffDuration time.Duration
    RequireHuman    bool
}

func (e *ExecutionOrchestrator) HandleFailure(task *Task, err error) {
    if attempt < policy.MaxAttempts {
        // 自動リトライ
    } else if policy.RequireHuman {
        // バックログに追加
    } else {
        // FAILED としてマーク
    }
}
```

#### FR-P3-005: バックログ管理

```go
type BacklogItem struct {
    ID          string      `json:"id"`
    TaskID      string      `json:"taskId"`
    Type        BacklogType `json:"type"`  // FAILURE | QUESTION | BLOCKER
    Description string      `json:"description"`
    Priority    int         `json:"priority"`
    CreatedAt   time.Time   `json:"createdAt"`
    ResolvedAt  *time.Time  `json:"resolvedAt,omitempty"`
}
```

---

## 3. データモデル

### Task（拡張）

| フィールド | 型 | 説明 |
|-----------|-----|------|
| dependencies | []string | 依存タスクIDリスト |
| parentId | *string | 親タスクID（WBS階層用） |
| wbsLevel | int | WBS階層レベル（1=概念設計, 2=実装設計, 3=実装） |
| phaseName | string | フェーズ名 |
| sourceChatId | *string | 生成元チャットセッションID |
| acceptanceCriteria | []string | 達成条件リスト |

### ChatSession

| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | string | セッションID |
| workspaceId | string | ワークスペースID |
| messages | []ChatMessage | メッセージ一覧 |
| createdAt | time.Time | 作成日時 |
| updatedAt | time.Time | 更新日時 |

### ChatMessage

| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | string | メッセージID |
| role | string | user / assistant / system |
| content | string | メッセージ本文 |
| timestamp | time.Time | タイムスタンプ |
| generatedTasks | []string | このメッセージで生成されたタスクID |

### BacklogItem

| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | string | バックログID |
| taskId | string | 関連タスクID |
| type | BacklogType | FAILURE / QUESTION / BLOCKER |
| title | string | タイトル |
| description | string | 説明 |
| priority | int | 優先度 |
| createdAt | time.Time | 作成日時 |
| resolvedAt | *time.Time | 解決日時 |
| resolution | string | 解決方法 |

---

## 4. アーキテクチャ

### 4層構造（維持 + 拡張）

```
┌─────────────────────────────────────────────────────┐
│  multiverse-ide (Desktop UI)                        │
│  - ChatWindow → タスク生成                           │
│  - GridCanvas → 依存グラフ表示                       │
│  - WBSView → マイルストーン表示                      │
│  - BacklogPanel → バックログ管理                     │
└──────────────┬──────────────────────────────────────┘
               │ Wails IPC + Events
┌──────────────▼──────────────────────────────────────┐
│  Orchestrator Layer                                 │
│  - ChatHandler (NEW)                                │
│  - TaskGraphManager (NEW)                           │
│  - ExecutionOrchestrator (NEW)                      │
│  - BacklogStore (NEW)                               │
│  - TaskStore / Scheduler                            │
└──────────────┬──────────────────────────────────────┘
               │
┌──────────────▼──────────────────────────────────────┐
│  AgentRunner Core + Meta-agent                      │
│  - FSM（既存維持）                                   │
│  - decompose プロトコル (NEW)                        │
└──────────────┬──────────────────────────────────────┘
               │
┌──────────────▼──────────────────────────────────────┘
│  Worker (Docker Sandbox)                            │
└─────────────────────────────────────────────────────┘
```

### 新規コンポーネント

| コンポーネント | 場所 | 責務 |
|--------------|------|------|
| ChatHandler | internal/chat/handler.go | チャット入力のMeta-agent転送、タスク生成 |
| TaskGraphManager | internal/orchestrator/task_graph.go | 依存関係グラフの構築・管理 |
| ExecutionOrchestrator | internal/orchestrator/executor.go | 自律実行ループ、一時停止/再開 |
| BacklogStore | internal/orchestrator/backlog.go | 問題・検討材料の永続化 |
| ChatSessionStore | internal/chat/session_store.go | チャット履歴の永続化 |

---

## 5. マイルストーン

### M1: チャット→タスク生成（2週間）

**Week 1:**
- Task 構造体拡張
- Meta-agent decompose プロトコル
- ChatHandler 実装
- ChatSession 永続化

**Week 2:**
- FloatingChatWindow バックエンド連携
- タスク生成結果のUI表示
- E2Eテスト

### M2: 依存グラフ・WBS表示（2週間）

**Week 3:**
- TaskGraphManager
- Scheduler 依存チェック拡張
- ConnectionLine コンポーネント

**Week 4:**
- WBS ツリービュー
- マイルストーン表示
- 進捗率計算

### M3: 自律実行ループ（2週間）

**Week 5:**
- ExecutionOrchestrator
- 一時停止/再開
- Wails Events リアルタイム通知

**Week 6:**
- 自動リトライ
- BacklogStore
- バックログUI

---

## 6. 受け入れ条件

### Phase 1 完了条件

| ID | 条件 |
|----|------|
| AC-P1-01 | チャットからテキストを送信できる |
| AC-P1-02 | Meta-agent がタスク分解を行い、複数タスクが生成される |
| AC-P1-03 | 生成タスクが tasks/*.jsonl に永続化される |
| AC-P1-04 | タスクに依存関係情報が含まれる |
| AC-P1-05 | GridCanvas にノードとして表示される |

### Phase 2 完了条件

| ID | 条件 |
|----|------|
| AC-P2-01 | タスク間依存が矢印で表示される |
| AC-P2-02 | 依存タスク未完了時に BLOCKED 状態になる |
| AC-P2-03 | WBS ビューでツリー表示できる |
| AC-P2-04 | マイルストーン別の進捗率が表示される |

### Phase 3 完了条件

| ID | 条件 |
|----|------|
| AC-P3-01 | 自動実行で依存順にタスクが実行される |
| AC-P3-02 | 一時停止で新規タスク開始が停止する |
| AC-P3-03 | 再開で実行が継続する |
| AC-P3-04 | 失敗時に自動リトライまたはバックログ追加 |

---

## 7. 技術的リスクと対策

| リスク | 影響度 | 対策 |
|--------|--------|------|
| Meta-agent のタスク分解精度が低い | 高 | プロンプトエンジニアリング、人間レビュー機能 |
| 依存関係の循環参照 | 中 | グラフ構築時にサイクル検出 |
| 大量タスク時のUI性能劣化 | 中 | 仮想化描画（可視領域のみレンダリング） |
| ファイルコンフリクト検出漏れ | 高 | Meta-agent に明示的なコンフリクト分析を依頼 |
| 自律実行中のエラー連鎖 | 高 | 失敗回数閾値でループ停止、人間判断モード |

---

## 8. 既存設計との差分

### 削除

- タスク作成ダイアログ（FR-IDE-012）→ チャットに置換

### 維持

- 4層アーキテクチャ
- $HOME/.multiverse/workspaces/ 構造
- JSONL/JSON 永続化形式
- FSM 状態遷移
- Task/Attempt のステータス定義

### 変更

- Task 構造体: 依存関係、WBS、生成元情報追加
- Meta-agent プロトコル: decompose 追加
- Scheduler: 依存チェック追加

---

## 9. 技術スタック

### バックエンド（維持）

| カテゴリ | 技術 | バージョン |
|---------|------|-----------|
| 言語 | Go | 1.23+ |
| デスクトップ | Wails | v2 |
| コンテナ | Docker | - |
| LLM | OpenAI API | - |

### フロントエンド（維持）

| カテゴリ | 技術 | バージョン |
|---------|------|-----------|
| フレームワーク | Svelte | 4 |
| 型安全 | TypeScript | 5 |
| ビルド | Vite | 5 |
| パッケージ管理 | pnpm | - |

### 新規追加

| カテゴリ | 技術 | 用途 |
|---------|------|------|
| グラフ描画 | SVG | 依存関係の矢印描画 |
| リアルタイム通信 | Wails Events | 状態変更通知 |
