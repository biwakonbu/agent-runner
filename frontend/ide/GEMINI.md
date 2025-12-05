# frontend/ide/GEMINI.md

このディレクトリは Multiverse IDE の Web フロントエンドを提供します。

## 技術スタック

- **Svelte 4**: リアクティブ UI フレームワーク
- **TypeScript 5**: 型安全な JavaScript
- **Vite 5**: 高速ビルドツール
- **Wails v2**: Go ↔ Web IPC
- **oxlint**: 高速リンター
- **Storybook 8**: コンポーネントカタログ
- **Playwright**: E2E テスト

## パッケージマネージャー

**pnpm** を使用します。npm や yarn は使用しないでください。

## 開発コマンド

```bash
pnpm install          # 依存パッケージインストール
pnpm dev              # 開発サーバー起動
pnpm build            # 本番ビルド
pnpm check            # Svelte 型チェック
pnpm lint             # ESLint (oxlint) チェック
pnpm lint:css         # Stylelint チェック
pnpm check:all        # 全チェック（型 + lint + knip）
pnpm storybook        # Storybook 起動（http://localhost:6006）
pnpm test:e2e         # Playwright E2E テスト
```

## ディレクトリ構成

- **`src/`**: ソースコード
  - `design-system/`: デザイントークン・基底コンポーネント
  - `stores/`: Svelte Store（状態管理）
  - `types/`: TypeScript 型定義
  - `lib/`: UI コンポーネント
- **`wailsjs/`**: Wails 自動生成バインディング
- **`tests/`**: E2E テスト（Playwright）

## デザインシステム

**テーマ**: Nord Deep（深い背景にパステル UI が輝くデザイン）

詳細は `src/design-system/CLAUDE.md` および `CLAUDE.md` を参照してください。
