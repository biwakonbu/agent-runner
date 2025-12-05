# TODO: multiverse IDE v0.1 Implementation

Based on `PRD.md`.

---

## 進捗サマリ

| Phase                  | Status  | 備考                                                      |
| ---------------------- | ------- | --------------------------------------------------------- |
| M0: 基盤整備           | ✅ 完了 | Backend ロジック実装済み、テスト済み                      |
| M1: IDE v0.1           | ✅ 完了 | Wails ビルド成功、`build/bin/multiverse-ide.app` 生成済み |
| M1.1: Core Integration | ✅ 完了 | Executor 実装、AgentRunner Core 連携完了                  |

---

## M0: 基盤整備（バックエンドのみ） ✅ 完了

- [x] **ディレクトリ構成の作成**
- [x] **Workspace 管理 (Backend)**
- [x] **Task / Attempt 永続化 (Backend)**
- [x] **Orchestrator (Dummy Implementation)**

---

## M1: IDE v0.1 (UI Implementation) ✅ 完了

- [x] **Wails プロジェクト立ち上げ**
- [x] **UI コンポーネント実装**
- [x] **Backend API 実装 (Go)**

---

## M1.1: Core Integration ✅ 完了

- [x] **Executor 実装**

  - [x] `internal/orchestrator/executor.go`: AgentRunner Core 連携
  - [x] `ExecuteTask()`: Task YAML 生成 → agent-runner 実行 → 結果更新

- [x] **IDE 統合**
  - [x] `app.go`: Executor を使用した RunTask 実装
  - [x] バックグラウンド実行対応

---

## 実装済みファイル一覧

### Core Integration (新規)

| ファイル                            | 説明                          |
| ----------------------------------- | ----------------------------- |
| `internal/orchestrator/executor.go` | AgentRunner Core 実行ラッパー |

### Backend (既存)

| ファイル  | 説明                                     |
| --------- | ---------------------------------------- |
| `main.go` | Wails エントリポイント                   |
| `app.go`  | Wails バインディング (Executor 統合済み) |

---

## 次のステップ

1. アプリの起動テスト: `open build/bin/multiverse-ide.app`
2. Task 作成 → Run → ステータス確認の E2E テスト
3. エラーハンドリングの改善
