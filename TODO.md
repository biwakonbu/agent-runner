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

## M1.2: Orchestrator CLI / Daemon ✅ 完了

- [x] **Orchestrator CLI (`cmd/multiverse-orchestrator`)**
  - [x] Main entrypoint implementation
  - [x] IPC Queue Consumer (Polling/Events)
  - [x] Integrate with `Executor` to run queued tasks

---

## M1.3: E2E Verification ✅ 完了

- [x] **ビルド検証**

  - [x] `wails build` → `build/bin/multiverse-ide.app` 生成成功
  - [x] `agent-runner` バイナリ生成成功

- [x] **テスト検証**
  - [x] Go Unit Tests: 7 packages passed
  - [x] Orchestrator Tests: 2 tests passed
  - [x] Playwright E2E Tests: 2 tests passed (タスクリスト表示、タスク作成フロー)

---

## 次のステップ

### 優先度: 高

1. **手動 E2E フロー確認**: `open build/bin/multiverse-ide.app` で実際にアプリを起動し、タスク作成 → 実行 → 結果確認の動作を検証
2. **Orchestrator 統合テスト**: IDE → Queue → Orchestrator → agent-runner の統合フローテスト

### 優先度: 中

3. **エラーハンドリング改善**: 各レイヤー(IDE, Orchestrator, Core)でのエラー表示とリカバリ
4. **IPC Queue テスト追加**: `internal/orchestrator/ipc` にユニットテスト追加

### 優先度: 低

5. **UI 改善**: タスク実行中のプログレス表示、ログビューア実装
6. **ドキュメント更新**: テスト戦略、デプロイ手順の文書化
