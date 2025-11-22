# TODO: AgentRunner v1 実装タスク

`docs/` ディレクトリの仕様書に基づき、現状の実装で不足している機能や改善が必要な項目をまとめる。

## 1. 実装すべき機能（不足機能）

### Worker パッケージ

- [x] **相対パスの完全な解決**

  - 現状: `repo` が空の場合のみ絶対パス化している。
  - 仕様: `.` や `../foo` などの相対パスが指定された場合も、Docker マウント用に絶対パスに変換する必要がある。
  - 対象: `internal/worker/executor.go`

- [x] **実行タイムアウト制御の確認と実装**
  - 仕様: `runner.worker.max_run_time_sec` (デフォルト 1800s) に基づき、Worker の実行を打ち切る必要がある。
  - 現状: `context.WithTimeout` などで制御されているか確認し、未実装なら実装する。
  - 対象: `internal/worker/executor.go`

### Core パッケージ

- [ ] **テストコマンドの実行（将来的な拡張）**
  - 仕様: `task.test.command` が指定されている場合、v1 では Worker に任せる方針だが、Core 側で実行する基盤も検討が必要（現状はスキップで OK）。

## 2. 改善・リファクタリング

### テスト

- [x] **Docker 統合テストの CI 組み込み**
  - 現状: `make test-worker-coverage` で手動実行。
  - 目標: CI で自動実行できるようにする（GitHub Actions の Service Container 利用など）。
  - 完了: `.github/workflows/ci.yaml` に `test-worker-coverage` ステップを追加し、カバレッジレポートをアーティファクトとして保存するように設定。

### ドキュメント

- [x] **CLAUDE.md の更新**
  - 「既知の課題」として記載されている相対パス問題が解決したら削除する。

## 3. 完了済み（確認用）

- [x] Meta プロトコル (`plan_task`, `next_action`, `completion_assessment`)
- [x] LLM リトライロジック (Exponential Backoff)
- [x] Sandbox 基本機能 (Start/Exec/Stop)
- [x] ImagePull 自動実行
- [x] Codex 認証自動マウント
- [x] `env:` プレフィックス対応
- [x] Task Note 生成
