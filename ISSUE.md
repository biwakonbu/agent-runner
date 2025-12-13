# ISSUE Log (2025-12-07)

## Open Items

- [x] Meta 層が CLI サブスクリプション未対応（API キー不要方針と不整合）

  - `app.go` の `newMetaClientFromEnv()` と `chat.NewHandler` が HTTP クライアント (OPENAI_API_KEY 前提) を生成し、LLMConfigStore/設定画面を参照しない。Codex / Gemini / Claude Code / Cursor などの CLI セッション再利用方針と乖離。
  - 対応完了（2025-12-14）: Meta プロバイダを `CLIProvider` (汎用) と `OpenAIProvider` に分離し、API キーなしでの利用を可能にした。`app.go` は `claude` または `codex` を自動検出し、サブスクリプションベースの利用を優先するよう改修済み。

- [x] LLM 設定 UI が実行系に反映されず、API キー前提の表示が残存

  - `LLMSettings` は Kind/Model を保存するが、CLI セッション状態を表示できず、`TestLLMConnection` も OpenAI HTTP 前提で CLI セッションを検証できない。API キーは不要なので UI をセッション表示に置換する必要。
  - 対応完了（2025-12-14）: API キーの強制を解除し、`TestLLMConnection` がプロバイダ経由で CLI バージョンチェックを行うように変更済み。設定画面で API キーが空でも CLI があれば接続成功となる。

- [ ] 実行ログ（stdout/stderr）のリアルタイム配信/表示の整備

  - バックエンド: stdout/stderr を逐次読み取り、`task:log` イベントを送出済み（`internal/orchestrator/executor.go:93`、`internal/orchestrator/executor.go:121`、`internal/orchestrator/events.go:39`）。
  - フロント: `task:log` を購読して store に蓄積する実装は存在（`frontend/ide/src/stores/logStore.ts:49`）。
  - 残タスク: UI 上でのタスク別フィルタ/クリア導線/常時表示など、運用可能な表示体験に仕上げる。

- [ ] Codex CLI セッションの存在確認・注入手段が未整備
  - Worker Executor は `codex exec ...` を呼び出すが、セッション有無の検証・警告やコンテナへのセッション注入方法（環境変数/ボリューム）が明確でない。
  - 対応: コンテナ起動時にセッション確認を行い、失敗時は UI に警告を返す。セッションの受け渡し仕様をドキュメント化。

## Deferred (moved from TODO.md, 2025-12-13)

- [ ] 手動一気通し（必要なら）

  - 手順:
    1. IDE 起動 → ワークスペース選択
    2. チャットで簡単な要求を入力
    3. （Chat Autopilot）自動で計画 → 実行に遷移することを確認（未実装の場合はギャップとして残る、`app.go:432-452`）
    4. 完了ステータスと成果物を確認
  - 観測ポイント:
    - `design/`・`state/`・`tasks/` の 3 層が整合する。
    - 依存順に実行される。

- [ ] Artifacts.Files の自動抽出/反映

  - 【目的】実行したタスクが「どのファイルを生成/変更したか」を IDE で追跡できるようにする。
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

- [ ] Meta Protocol のバージョニング導入

  - 【目的】Meta-agent と Core 間のプロトコル互換性を将来にわたって維持する。
  - 対象ファイル（候補）:
    - `internal/core/meta/*`
    - `docs/specifications/meta-protocol.md`
  - 実装タスク:
    1. YAML メッセージに `protocol_version`（または同等）を追加し、Core 側で解釈する。
    2. バージョン不一致時のフォールバック/警告/拒否方針を定義する。
  - 完了条件:
    - プロトコル更新時に旧クライアントが安全に扱える。

- [ ] 追加 Worker 種別のサポート

  - 【目的】`codex-cli` 以外の CLI エージェントを Worker として選択可能にする。
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

- [ ] IPC の WebSocket / gRPC 化

  - 【目的】ファイルポーリング IPC の性能/拡張性の制約を解消する。
  - 対象ファイル（候補）:
    - `internal/orchestrator/ipc/*`
    - `frontend/ide/src/*`
  - 実装タスク:
    1. Queue/イベント通知を WebSocket か gRPC に置き換える設計を確定する。
    2. 既存 file-based IPC と並行稼働できる移行パスを用意する。
  - 完了条件:
    - 大量ジョブ時のポーリング負荷が解消される。

- [ ] Frontend E2E の安定化

  - 【目的】IDE フロントの E2E が CI で継続的に回る状態にする。
  - 対象ファイル（候補）:
    - `frontend/ide/*`
    - `docs/guides/testing.md`
  - 実装タスク:
    1. タイムアウト/待機条件/テストデータを見直し安定化する。
    2. 失敗時ログの拡充とリトライ方針を整備する。
  - 完了条件:
    - `pnpm test:e2e` が安定して完走する。

- [ ] Task Note 保存の圧縮
  - 【目的】大きな Task Note/履歴の保存サイズを抑え、読み書き性能を維持する。
  - 対象ファイル（候補）:
    - `internal/note/*`
    - `internal/orchestrator/persistence/*`
  - 実装タスク:
    1. Task Note の圧縮形式（gzip 等）と保存/読み込み API を定義する。
    2. 既存データとの後方互換を確保する。
  - 完了条件:
    - Task Note の保存容量が有意に削減される。
