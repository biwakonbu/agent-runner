# cmd/GEMINI.md

このディレクトリには、Multiverse プロジェクトの各アプリケーションのエントリポイント（`main` パッケージ）が含まれています。

## コンポーネント一覧

### 1. `agent-runner/`

- **概要**: AgentRunner Core の CLI。
- **入力**: 標準入力からの Task YAML。
- **役割**: 単一タスクの実行、Docker サンドボックス管理、Meta-agent との通信。
- **使用方法**: `agent-runner < task.yaml`

### 2. `multiverse-ide/`

- **概要**: デスクトップ IDE アプリケーション。
- **技術**: Wails (Go) + Svelte (Frontend)。
- **役割**: ユーザーインターフェース。タスク作成、監視、結果確認。
- **依存**: 内部で `orchestrator` パッケージを使用し、バックグラウンド処理を行う。

### 3. `multiverse-orchestrator/`

- **概要**: タスクオーケストレーション用 CLI / デーモン（**開発中**）。
- **役割**: IPC キューの監視、Worker プールの管理、`agent-runner` の呼び出し。
- **ステータス**: 現在 CLI 実装は未完了。ロジックは `internal/orchestrator` に実装済み。
