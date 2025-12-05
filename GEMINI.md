# GEMINI.md

## プロジェクト概要

`multiverse` は、AI ネイティブな開発環境を実現するための統合プラットフォームです。以下のコンポーネント群で構成され、ローカル環境での自律的なソフトウェア開発タスクの実行を可能にします。

1.  **AgentRunner Core (Engine)**: AI エージェントを Docker サンドボックス内で安全に実行・管理するコアエンジン。
2.  **Multiverse Orchestrator (Backend)**: 複数のタスクと Worker を管理し、IDE からのリクエストを処理するオーケストレーション層。
3.  **Multiverse IDE (Frontend)**: 開発者がタスクの作成、実行、監視を行うためのデスクトップアプリケーション (Wails + Svelte)。

## アーキテクチャ

システム全体は 3 層構造になっています。

- **Frontend Layer (IDE)**: ユーザーインターフェース。タスクの定義と監視を担当。Wails で実装。
- **Orchestration Layer**: タスクのキューイング、スケジューリング、永続化（`$HOME/.multiverse`）を担当。
- **Execution Layer (Core)**: 実際のコード生成やテスト実行。Docker コンテナ内で完結。

## ビルドと実行

### パッケージマネージャー

このプロジェクトでは **pnpm** を使用します。npm や yarn は使用しないでください。

### 全体ビルド (IDE)

```bash
wails build
```

### コンポーネント別ビルド

- **AgentRunner Core**:
  ```bash
  go build ./cmd/agent-runner
  ```
- **Orchestrator CLI**:
  ```bash
  go build ./cmd/multiverse-orchestrator
  ```

### フロントエンド開発

```bash
cd frontend/ide
pnpm install
pnpm dev          # 開発サーバー起動
pnpm check        # Svelte 型チェック
pnpm lint         # ESLint (oxlint) チェック
pnpm lint:css     # Stylelint チェック
pnpm test:e2e     # Playwright E2E テスト
pnpm storybook    # Storybook 起動
```

## 開発の規約

- **言語**: コミュニケーションは常に **日本語** で行います。
- **ドキュメント**:
  - `docs/`: 仕様書と設計書。
  - `GEMINI.md`: 各ディレクトリの役割とコンテキスト（本ファイルおよびサブディレクトリ内のファイル）。
  - `CLAUDE.md`: AI アシスタント向けのガイドライン（共存）。

## ディレクトリ構成

- `cmd/`: 各コンポーネントのエントリポイント (`agent-runner`, `multiverse-ide`, `multiverse-orchestrator`)
- `internal/`: 内部ロジック (`core`, `meta`, `worker`, `orchestrator`, `ide`)
- `frontend/`: IDE のフロントエンドコード (Svelte)
- `docs/`: プロジェクト全体のドキュメント

## 重要：作業前の確認事項

各ディレクトリの `CLAUDE.md` および `GEMINI.md` を必ず確認してください。
