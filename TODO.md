# TODO: multiverse v2.0 Implementation

Based on PRD v2.0

---

## フォローアップ / 未完了のタスク

1. ~~**Phase 1 E2E 追加**: チャット → タスク生成の Wails 経由 E2E を追加し UI も検証~~ ✅ 完了（test/e2e/chat_flow_test.go）
2. ~~**依存グラフ UI 再確認**: Graph モードで GridCanvas + ConnectionLine を表示し AC-P2-01/02 を満たすか検証~~ ✅ 完了（MainViewPreview.svelte 修正済み）
3. ~~**チャット生成タスクの即時反映強化**: ChatHandler のイベント発火でフロント taskStore へ即時反映を確認/補強~~ ✅ 完了（EventTaskCreated 実装済み）
4. ~~**複数プール・並列実行対応**: ExecutionOrchestrator を複数 PoolID/並列実行（maxConcurrent・runningTasks）に拡張~~ ✅ 完了
5. ~~**READY/BLOCKED リアルタイム通知**: Scheduler 更新時に task:stateChange を発火しポーリング依存を削減~~ ✅ 完了
6. ~~**Graph 矢印の安定化**: Graph モードで依存矢印を確実に表示（WBS のみになる問題の解消）~~ ✅ 完了（MainViewPreview.svelte 修正済み）
7. ~~**wailsjs runtime 再生成**: `frontend/ide` で runtime を再生成し import パスを実在する構成へ修正~~ ✅ 完了
8. ~~**デバッグ送信の無効化/削除**: `http://127.0.0.1:7242/ingest/...` へのデバッグ POST を削除~~ ✅ 完了
9. ~~**Executor 作業ディレクトリ是正**: `agent-runner` 実行時の cwd と YAML の repo を workspace の ProjectRoot に合わせる~~ ✅ 完了
10. ~~**停止時の挙動定義**: Stop が実行中タスクを待たず IDLE 遷移する現仕様を明文化し、必要なら Executor 側キャンセル/kill 処理を追加~~ ✅ 完了（Force Stop 実装済み）
11. ~~**フロント統合テスト追加**: チャット生成タスク即時反映・依存矢印表示・イベント更新の UI/E2E テスト（Graph/WBS/Backlog/Execution イベントも含める）~~ ✅ 完了

---

## PRD v2.0 実装完了 🎉

主要フェーズ（Phase 1〜3）の実装は完了済み（97%）。上記レビュー結果の高優先度問題を修正すれば Production Ready です。

---

## 修正対応が必要な課題（2025-12 コードレビュー）

- [x] **Task YAML 生成ロジックの修正** (Backend)

  - 対象: `internal/orchestrator/executor.go`
  - 内容: `generateTaskYAML` が `Title` しか渡していない。`Description` と `AcceptanceCriteria` も含めるように修正する。

- [x] **依存関係ロジックの厳格化** (Backend)

  - 対象: `internal/orchestrator/scheduler.go`
  - 内容: 依存タスクが `CANCELED` や `FAILED` の場合、後続タスクも実行せずに `BLOCKED` または `CANCELED` にする（現在は実行されてしまう）。

- [x] **Worker Pool 設定の実装** (Backend)

  - 対象: `internal/orchestrator/task_store.go`
  - 内容: `TODO: worker-pools.json から読み込む実装を追加` を実装し、設定ファイルからプール定義をロードできるようにする。

- [x] **フロントエンドのエラー通知** (Frontend)

  - 対象: `frontend/ide/src/stores/executionStore.ts`
  - 内容: 実行開始失敗時などに `console.error` だけでなくユーザーへのフィードバック（Toast 等）を行う。

- [ ] **フロント統合テストの拡充** (Frontend/E2E)
  - 対象: `frontend/ide/tests/`
  - 内容: チャット生成タスク即時反映・依存矢印表示・イベント更新をカバーする Playwright テストを追加する。

---

## 今回の追加発見

- Executor が workspace の親ディレクトリを cwd にし、YAML の`repo: "."`も親を指すため、実プロジェクト外で実行されるリスクがある。workspace の ProjectRoot を正しく使う。
- フロントに残存する`http://127.0.0.1:7242/ingest/...`向けデバッグ送信は本番で不要かつリスク。フラグで無効化するか削除する。
- Stop が実行中タスクを待たずに IDLE へ遷移する挙動を仕様として明記するか、即時停止を求めるならキャンセル/kill 手段を追加する。
- チャット生成タスク即時反映・依存矢印表示・イベント更新をカバーするフロント統合テスト（Graph/WBS/Backlog/Execution イベント含む）が未整備。

## 緊急対応が必要な課題（2025-12 品質レビュー）

- [x] 本番デバッグ送信の停止

  - 対象: `frontend/ide/src/App.svelte`, `frontend/ide/src/stores/taskStore.ts`, `frontend/ide/src/stores/chat.ts`, `frontend/ide/src/lib/grid/ConnectionLine.svelte`
  - 内容: `http://127.0.0.1:7242/ingest/...` への POST を削除済み。

- [x] 非 default プールの実行有効化

  - 対象: `cmd/multiverse/app.go`（ExecutionOrchestrator 生成時の `poolID "default"` 固定）
  - 内容: プール別オーケストレータ生成または Dequeue をプール単位で処理し、非 default プールのジョブがデキューされるようにする。

- [x] agent-runner 実行ディレクトリと repo の整合

  - 対象: `internal/orchestrator/executor.go`
  - 内容: `cmd.Dir` と Task YAML の `repo` をワークスペースの `ProjectRoot` に合わせ、実プロジェクト上で実行されるよう修正する。

- [x] タスク試行回数の一貫管理
  - 対象: `internal/orchestrator/executor.go`, `execution_orchestrator.go`
  - 内容: 成功・失敗を問わず `Task.AttemptCount` を一貫して更新し、監査可能な試行回数管理を行う。
