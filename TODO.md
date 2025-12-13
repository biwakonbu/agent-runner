# Pragmatic MVP（完了）+ Quality Hardening（vNext）実装手順

最終更新: 2025-12-14

## 0. 前提

- MVP のゴールは `PRD.md` の「MVP 完了条件」参照。
- リアルタイム UX は最小限（Wails Events による `chat:progress`/`task:stateChange`/`task:log` は実装済み、ログはフロントで最大 1000 行に制限）。UI のカクつき防止と体感性能を優先する（`internal/orchestrator/events.go:34`、`internal/orchestrator/executor.go:121`、`frontend/ide/src/stores/logStore.ts:16`）。

## 0.1 まず読む（PRD と並行で参照する一次ドキュメント）

### 必読（この順で読む）

1. `PRD.md`（到達点と主要フローの定義）
2. `docs/COMPLETE_DOCUMENTATION.md`（全体像・用語・仕様索引）
3. `docs/design/chat-autopilot.md`（会話だけで「計画 → 実行 → 質問 → 継続」する設計）
4. `docs/design/task-execution-and-visual-grouping.md`（タスク実行と多軸グルーピング設計）

### 仕様（必要になった時点で読む）

- `docs/specifications/orchestrator-spec.md`（Orchestrator の仕様）
- `docs/specifications/meta-protocol.md`（plan_task/next_action/completion_assessment/decompose の仕様）
- `docs/specifications/worker-interface.md`（WorkerCall/実行インターフェース）
- `docs/specifications/testing-strategy.md`（テスト戦略）

### 実装参照（設計の根拠として読む）

- `docs/design/data-flow.md`（Core の plan→next_action→assessment ループ）
- `docs/task-builder-and-golden-test-design.md`（MVP の一気通しパイプライン）
- `docs/CURRENT_STATUS.md`（現状の制約・既知課題）

### 「この TODO を更新する時」のルール

- `PRD.md` と矛盾する変更は、先に `PRD.md` を更新してから TODO を直す。
- Chat Autopilot / grouping の設計変更は、先に `docs/design/chat-autopilot.md` と `docs/design/task-execution-and-visual-grouping.md` を更新してから TODO を直す。

## 1. 状態

- MVP は完了（`PRD.md` の 9 章）。
- vNext は “Quality Hardening（100点化）” を最優先で実施する（`PRD.md` の 10〜12 章）。
- `ISSUE.md` は運用上のバックログ置き場として残すが、**設計/DoD/優先度の真実源は `PRD.md`** とする。

---

## 2. 反省の継承（このTODOを実行する前提）

この章は “実装手順に必ず反映するべき反省点” をチェックリスト化する。詳細は `PRD.md` 10 章を正とする。

- 【事実】plan_patch の入力コンテキストを Meta に完全に渡していない（`internal/meta/utils.go`、`PRD.md` 10.1）。
- 【事実】delete(cascade=false) のセマンティクスが未確定で、WBS 孤児が発生し得る（`internal/chat/plan_patch.go`、`PRD.md` 10.2）。
- 【事実】history の順序が設計と逆で、監査/復元性が落ちる（`internal/chat/plan_patch.go`、`PRD.md` 10.3）。
- 【事実】ユニットテストが外部依存/モック不整合で品質ゲートとして機能しない（`internal/meta/*`、`PRD.md` 10.4）。

---

## 3. vNext: Quality Hardening（100点化）実行順序

### 3.1 事前チェック（必須）

- `go test ./...` を実行し、現状の失敗を “再現可能な形” で把握する（例: `internal/meta/client_test.go`）。
- WBS の不変条件（`PRD.md` 11.2）を読み、delete/move の期待動作をチームで合意する。

### 3.2 P0（最優先・これが終わらない限り他へ進まない）

#### 3.2.1 QH-001: plan_patch プロンプトに構造化コンテキストを完全継承

- 目的/受け入れ条件/検証は `PRD.md` 12.1 を正とする。
- 実装後、最低限以下を追加で確認する:
  - `plan_patch` のプロンプトに `existing_wbs.node_index` と `conversation_history` が含まれる。
  - トリミング規則がコードとテストに固定されている（将来の regress を防ぐ）。

#### 3.2.2 QH-003: delete(cascade=false) のセマンティクス確定 + 実装 + テスト

- 決め打ち（案A/案B）は `PRD.md` 12.2 を参照し、PRD に採用方針を明記してから実装する。
- WBS 不変条件テストを必須化する（孤児/重複 children を検知）。

#### 3.2.3 QH-004: history→design/state の順序保証（擬似トランザクション）

- `PRD.md` 12.3 の受け入れ条件を満たすまで完了扱いにしない。

#### 3.2.4 QH-005: meta テストの無外部依存化（品質ゲート復旧）

- `go test ./...` がネットワーク無しで通ることを “Done” の必須条件にする（`PRD.md` 11.3）。

### 3.3 P1（運用品質）

#### 3.3.1 QH-006: 実行ログ UI の運用可能化

- `ISSUE.md:15` の残タスクを “最小導線” に落とす（タスク別フィルタ/クリア/常時表示）。

#### 3.3.2 QH-007: Codex CLI セッション検証と注入仕様

- `ISSUE.md:21` を “実装 + ドキュメント” まで閉じる。

### 3.4 P2（将来拡張：仕様/テストを先に固めてから実装）

以下は `PRD.md` 12.6 を正とし、着手順は運用上の都合で入れ替えてよい。ただし P0 を完了してから着手する。

- QH-008: Artifacts.Files の自動抽出/反映（`ISSUE.md:38`）
- QH-009: Meta Protocol のバージョニング導入（`ISSUE.md:52`）
- QH-010: 追加 Worker 種別のサポート（`ISSUE.md:64`）
- QH-011: IPC の WebSocket / gRPC 化（`ISSUE.md:78`）
- QH-012: Frontend E2E の安定化（`ISSUE.md:90`）
- QH-013: Task Note 保存の圧縮（`ISSUE.md:102`）

---

## 4. 完了判定（このTODOのDefinition of Done）

以下がすべて満たされない限り “100点” としてクローズしない。DoD の正は `PRD.md` 11 章。

- 【事実】`go test ./...` がネットワーク不要で安定して通る。
- 【事実】plan_patch の create/update/move/delete（cascade true/false）で WBS 不変条件を満たすことをテストで担保できる。
- 【事実】history の順序が設計通りで、失敗時の扱いが実装・テスト・ドキュメントで一致している。
