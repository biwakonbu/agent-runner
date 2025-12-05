# test/GEMINI.md

このディレクトリには、プロジェクト全体のテストファイルが含まれています。

## テスト戦略（4 段階）

| Stage | 名称             | ディレクトリ        | 依存性       | 実行速度       |
| ----- | ---------------- | ------------------- | ------------ | -------------- |
| 1     | ユニットテスト   | `internal/*/`       | モック       | 高速（<1s）    |
| 2     | 統合テスト       | `test/integration/` | モック       | 中速（<5s）    |
| 3     | Sandbox テスト   | `test/sandbox/`     | Docker       | 遅い（10-30s） |
| 4     | Codex 統合テスト | `test/codex/`       | Docker+Codex | 最遅（1-5m）   |

## ディレクトリ構成

- **`integration/`**: モック統合テスト。FSM フローの完全性検証。
- **`sandbox/`**: Docker サンドボックステスト（`-tags=docker`）。コンテナ管理の検証。
- **`codex/`**: Codex CLI 統合テスト（`-tags=codex`）。実運用シナリオ検証。
- **`e2e/`**: フロントエンド E2E テスト（Playwright）。

## 実行方法

```bash
# ユニットテスト（高速）
go test ./internal/...

# 統合テスト（モック）
go test ./test/integration/...

# Docker サンドボックステスト
go test -tags=docker ./test/sandbox/...

# Codex 統合テスト（実 LLM）
go test -tags=codex ./test/codex/...

# フロントエンド E2E テスト（pnpm 使用）
cd frontend/ide
pnpm test:e2e
```

## 詳細

テストパターン・ベストプラクティスの詳細は `CLAUDE.md` を参照してください。
