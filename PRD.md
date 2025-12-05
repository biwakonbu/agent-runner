# PRD: multiverse IDE v0.1（UI 実装開始版）

## 0. 本 PRD の位置づけ

- 本 PRD は、`multiverse` プロジェクトにおける **ローカル IDE（デスクトップ UI）部分の v0.1** を定義する。
- 既存の AgentRunner コア仕様・設計と、multiverse Orchestrator / IDE の設計ドラフトを前提としつつ、
  **UI 実装を具体的に開始できるレベルまで要件・技術スタック・ディレクトリ構成を固める**ことを目的とする。
- 前提としているドキュメント（事実ベース）:

  - AgentRunner コア仕様・設計一式
  - multiverse / Task Orchestrator 統合設計 v0.2
  - multiverse / IDE MVP 設計ドラフト v0.3

以降、

- 「既存仕様に基づく前提」は【既存】
- 本 PRD で新たに定める内容・変更提案は【提案】
  と明示する。

---

## 1. 背景 / 課題

### 1.1 背景【既存】

- AgentRunner は、CLI から Task YAML を受け取り、Meta-agent（LLM）＋ Worker（CLI / コンテナ）をオーケストレーションする実行基盤として設計されている。
- 現状の主な利用想定は「CLI + Task YAML」であり、

  - タスク開始、状態確認、失敗時リトライ、ログ確認などを人間が CLI / ファイルで行う前提になっている。

- Web UI / IDE に関しては「将来の拡張」として設計レベルでは触れられているが、
  具体的な UI 実装・プロセス分離・ストレージ構造はまだ実装されていない。

### 1.2 課題【既存】

- YAML / CLI ベースでは、以下の点で利用性に課題がある:

  - タスクの一覧性が低く、ステータス管理が煩雑
  - 実行履歴や失敗リトライを「人間がファイルを開いて読む」前提
  - 複数 Worker / 複数タスクの並列実行状況を俯瞰しにくい

- IDE / Orchestrator の設計案はあるが、実装に必要な

  - 技術スタックの確定
  - リポジトリ / ディレクトリ構成の更新方針
  - MVP としての機能スコープ
    が明文化されていないため、エージェントに「どこまで作れば良いか」が渡しづらい。

---

## 2. ゴール（この PRD で達成したいこと）【提案】

1. **ローカル IDE（multiverse IDE）から、AgentRunner ベースのタスクを GUI で起動・監視できる状態にする。**
2. **ユーザー（開発者）が YAML を直接触らなくても、UI からタスクを作成・実行・状態確認できる。**
3. **IDE / Orchestrator / Worker / Core のプロセス分離・責務分離方針を維持しつつ、
   UI 実装に必要な技術スタックとディレクトリ構成の「現時点の答え」を固定する。**
4. **今後の自動化（エージェントによる開発・テスト）にとって、
   タスク・ワークスペース周りのメタデータが一貫した構造で管理されるようにする。**

---

## 3. スコープ / 非スコープ【提案】

### 3.1 スコープ（v0.1）

- デスクトップ IDE アプリケーション（multiverse IDE）

  - Workspace（ローカルプロジェクト）選択・登録
  - Task の作成・閲覧・ステータス表示
  - Task 実行（`Run` ボタン） → Orchestrator 経由で Core / Worker を起動
  - TaskAttempt（実行履歴）の一覧・結果概要の表示
  - WorkerPool の基本設定の閲覧（編集は v0.1 では任意）

- ストレージとプロセス構成

  - `$HOME/.multiverse` 以下に Workspace & Task メタデータを管理（既存設計の具現化）
  - Orchestrator ↔ Worker 間のファイルベース IPC（JSON）を前提とした構造の確定

- UI 実装のための技術スタック / ディレクトリ構成の更新

  - Go + Wails + Svelte + TypeScript を前提とした IDE 実装
  - 既存 AgentRunner コアのコードレイアウトとの共存方針

### 3.2 非スコープ（v0.1 でやらない）

- Meta-agent の高度なマルチ meta 構成（子 meta、複数 LLM）実装
- 既存 Core の FSM ロジックや YAML プロトコルの大幅な変更
- タスクの親子構造（ツリー表示）の UI（設計は踏まえるが実装は次フェーズ）
- Web ブラウザ版 IDE（ローカルデスクトップ専用想定）
- タスクノート（.agent-runner/task-xxx.md）のリッチビューア統合（将来拡張）

---

## 4. ターゲットユーザー / ユースケース【提案】

### 4.1 ターゲットユーザー

- ローカル環境でエージェント開発・検証を行うソフトウェアエンジニア
- 自動コーディング / 自動テストなどのタスクを複数並行で動かしたい開発者
- 将来的には CI / チーム開発でも利用されうるが、v0.1 は **ローカル 1 人利用** を主対象とする。

### 4.2 主要ユースケース（MVP）

1. **ローカルプロジェクトを IDE に登録し、タスクを 1 つ作って実行する**

   - プロジェクトディレクトリを選択
   - Task を作成（タイトル、種別、実行対象などを設定）
   - Run ボタンで実行
   - 実行状態（PENDING / RUNNING / SUCCEEDED / FAILED）を UI で確認

2. **失敗したタスクの実行履歴を確認し、再実行する**

   - Task 詳細画面から Attempt の履歴を一覧表示
   - 失敗した Attempt のエラー概要を確認
   - 再実行ボタンで新しい Attempt を起動

3. **複数タスクの並行実行状況をプール単位で把握する**

   - WorkerPool ごとに「Running / Queue 中」のタスク数を表示
   - タスクが溢れていないか、詰まっていないかを俯瞰する

---

## 5. 機能要件（Functional Requirements）

番号付きで定義する。テキストは UI 実装エージェント向けにやや具体的に記述。

### 5.1 Workspace 管理

- **FR-IDE-001**【既存設計の具現化】

  - IDE は「プロジェクトルートの絶対パス」から `workspace-id = sha1(projectRoot)[:12]` を計算し、
    `$HOME/.multiverse/workspaces/<workspace-id>/` をワークスペースディレクトリとして使用する。

- **FR-IDE-002**【提案】

  - IDE は初回オープン時にユーザーに `projectRoot` を選択させる UI を提供する。
  - 既に同じ `projectRoot` に対応する `workspace-id` が存在する場合は、その Workspace を再利用する。

- **FR-IDE-003**【既存＋提案】

  - `workspace.json` は以下の情報を必須として保持する（既存案の具現化）。

    - `version`
    - `projectRoot`
    - `displayName`
    - `createdAt`
    - `lastOpenedAt`
      （フォーマットは IDE MVP ドキュメントに準拠）

### 5.2 Task モデル / 一覧・詳細

- **FR-IDE-010**【既存＋提案】

  - Task は 1 Task 1 ファイル（JSONL）で `tasks/<task-id>.jsonl` に永続化される。
  - ファイル内の **最後の行が最新状態** として扱われる。

- **FR-IDE-011**【提案】

  - IDE は `Task` の最新スナップショットを読み込み、以下の情報を UI に表示する。

    - ID、Title、Status、PoolID、CreatedAt、UpdatedAt、StartedAt、DoneAt

  - TaskStatus は `PENDING / READY / RUNNING / SUCCEEDED / FAILED / CANCELED / BLOCKED` を扱う（既存案踏襲）。

- **FR-IDE-012**【提案】

  - Task 作成ダイアログでは最低限以下を入力できること:

    - Title（必須）
    - Kind（"codegen" / "test" 等のプリセット選択）
    - PoolID（`worker-pools.json` の ID から選択）

  - 作成時に

    - `tasks/<task-id>.jsonl` の最初のレコード
    - `instructions/tasks/<task-id>.md`（空テンプレート）
      を生成する。

### 5.3 Task 実行 / Attempt

- **FR-IDE-020**【既存＋提案】

  - 「Run」ボタン押下時、IDE は Orchestrator に「Task を実行キューに投入する」リクエストを行う。
  - Orchestrator は `ipc/queue/<pool-id>/<job-id>.json` にジョブを追加し、TaskStatus を `READY` に更新する。

- **FR-IDE-021**【既存＋提案】

  - Orchestrator により TaskAttempt が生成され、`attempts/<attempt-id>.json` に保存される。
  - IDE は Task 詳細画面から、指定 TaskID に紐づく Attempt 一覧を取得し表示する。

- **FR-IDE-022**【提案】

  - Attempt 詳細では、最低限以下を表示する:

    - AttemptID
    - Status（STARTING / RUNNING / SUCCEEDED / FAILED / TIMEOUT / CANCELED）
    - StartedAt / FinishedAt
    - ErrorSummary（あれば）

### 5.4 ステータス更新 / ポーリング

- **FR-IDE-030**【提案】

  - IDE は 1〜2 秒間隔で Workspace の Task 状態を再読込し、一覧画面の表示を更新する。
  - 実装は単純なポーリングでよく、WebSocket 等は v0.1 では不要。

- **FR-IDE-031**【提案】

  - Task 一覧において、現在 `RUNNING` or `READY` の Task に対し、Pool ごとのカウントを表示できるようにする

    - UI 上部 or サイドバーに Pool 単位のサマリ
    - 例: `codegen: 2 running / 3 queued`

---

## 6. 技術要件 / 技術スタック

### 6.1 プロセス構成【既存＋提案】

【既存】

- プロセスは IDE(UI) / Orchestrator / Worker / Core（AgentRunner）で **必ず分離** する方針。

【提案 / 明示】

- v0.1 でのプロセスは以下とする:

  - `multiverse-ide`

    - Wails + Go バックエンド
    - Svelte + TypeScript フロントエンド

  - `multiverse-orchestrator`

    - Go 製 CLI / デーモンプロセス
    - ファイルベース IPC で Worker と連携

  - Worker プロセス（複数種）

    - 例: `multiverse-worker-codegen`, `multiverse-worker-test`

  - AgentRunner Core

    - 既存 `agent-runner` CLI 相当
    - Orchestrator から起動される位置づけ

※ Core のバイナリ名や配置は後述ディレクトリ構成に従い整理する。

### 6.2 新技術スタック【提案】

- 言語 / ランタイム

  - Backend: Go（既存 AgentRunner と同一バージョン系を使用）
  - Desktop Shell: Wails（Go ベースのデスクトップアプリフレームワーク）
  - Frontend: Svelte + TypeScript

- ビルド / 配布

  - Wails 標準のビルドチェーンを利用し、単一バイナリ配布を前提とする（Mac を第一ターゲット）

- ストレージ

  - メタデータ: JSON / JSONL / Markdown（既存設計と同じ）
  - 配置: `$HOME/.multiverse/workspaces/...`（既存案踏襲）

### 6.3 ディレクトリ構成（リポジトリ）【提案】

既存 AgentRunner の Go コード構成（`cmd/agent-runner`, `internal/core`, `pkg/config` 等）
と、multiverse IDE / Orchestrator の設計を統合するため、リポジトリのトップレベルを次のように整理する案とする。

```text
multiverse/
  cmd/
    agent-runner/           # 既存 Core CLI（名称は当面維持）
    multiverse-orchestrator/
    multiverse-ide/         # Wails エントリポイント

  internal/
    core/                   # AgentRunner コア（既存）
    meta/                   # Meta-agent 通信（既存）
    worker/                 # Worker 実行・Sandbox（既存）
    note/                   # Task Note 生成（既存）

    orchestrator/           # Orchestrator ドメインロジック（新規）
      scheduler.go
      task_store.go
      ipc/
        filesystem_queue.go

    ide/                    # IDE バックエンドロジック（新規）
      workspace_store.go
      task_service.go

  frontend/
    ide/                    # Svelte + TS フロントエンド（新規）
      src/
        styles/
        ui/
        layout/
        features/
        stores/
        App.svelte
      package.json
      vite.config.ts 等

  docs/
    COMPLETE_DOCUMENTATION.md
    multiverse-architecture-v0.2.md
    multiverse-ide-mvp-v0.3.md
    PRD.md                  # 本ドキュメント
```

- 【事実】既存 AgentRunner の `internal/*` 構成は原則変更しない。
- 【提案】multiverse 関連コードを `cmd/multiverse-*` ＋ `internal/orchestrator`, `internal/ide`, `frontend/ide` として追加する。
- 【提案】IDE のフロントエンドは `frontend/ide` に単独で置き、Wails の build ステップがここを参照する。

---

## 7. データモデル / 永続化方針（要約）

### 7.1 Workspace

- `workspace-id = sha1(projectRoot)[:12]`【既存/確定】
- ディレクトリ:

  - `$HOME/.multiverse/workspaces/<workspace-id>/workspace.json`
  - `config/worker-pools.json`
  - `tasks/*.jsonl`
  - `attempts/*.json`
  - `ipc/queue/*`, `ipc/results/*`
  - `logs/*`

### 7.2 Task / Attempt

- Task

  - 1 Task 1 ファイル JSONL
  - 状態遷移は `PENDING → READY → RUNNING → SUCCEEDED / FAILED` を基本とする。【既存】

- Attempt

  - 1 Attempt 1 ファイル JSON
  - Task の再実行ごとに新しい Attempt が生成される。【既存】

### 7.3 Orchestrator ↔ Worker IPC

- `ipc/queue/<pool-id>/<job-id>.json`（Orchestrator → Worker）
- `ipc/results/<job-id>.json`（Worker → Orchestrator）
- メッセージ本体は JSON 固定【既存】

IDE は IPC ファイルそのものは直接扱わず、
Orchestrator の状態反映として Task / Attempt のファイルを読むのみとする【提案】。

---

## 8. 現状とのギャップと修正方針【提案】

| 項目            | 現状（ドキュメントベース）                                                | ありたい姿（本 PRD 後）                                                     | 修正方針                                                                  |
| --------------- | ------------------------------------------------------------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| UI プロセス     | 「Web UI / IDE は将来拡張」としてのみ記載                                 | `multiverse-ide` として Wails + Svelte の具体実装開始                       | 本 PRD で技術スタックとディレクトリ構成を固定し、エージェントタスク化     |
| ストレージ配置  | AgentRunner 側は `.agent-runner/task-*.md` をリポジトリ直下に作成【既存】 | IDE / Orchestrator は `$HOME/.multiverse/workspaces/...` にメタデータを集約 | Core の Task Note 仕様は維持しつつ、IDE 側のメタは `~/.multiverse` に限定 |
| Orchestrator    | 統合設計 v0.2 / IDE MVP v0.3 に概念レベルで記述                           | `multiverse-orchestrator` CLI として実装、IPC ディレクトリ構造を実際に作成  | `internal/orchestrator` を新設し、既存 Core とは明確に分離                |
| WorkerPool 管理 | ドキュメント上の JSON 設計のみ                                            | IDE から参照・編集できる UI コンポーネントに昇格                            | v0.1 では閲覧のみ最低限実装し、編集は v0.1.1 以降で拡張                   |

---

## 9. マイルストーン / フェーズ分割【提案】

### M0: 基盤整備（バックエンドのみ）

- `~/.multiverse/workspaces/...` ディレクトリ構造の自動生成
- `workspace.json` / `tasks/*.jsonl` / `attempts/*.json` の read/write ライブラリ（Go）
- Orchestrator のダミー実装（queue / results の単純な状態遷移）

### M1: IDE v0.1（この PRD の範囲）

- Wails プロジェクト立ち上げ (`cmd/multiverse-ide` + `frontend/ide`)
- Workspace 選択画面
- Task 一覧 / 詳細画面
- Task 作成ダイアログ
- Run ボタン → Orchestrator 呼び出し（簡易）
- ポーリングによる Task 状態更新
- Attempt 一覧・概要表示

### M1.1: Core 統合

- `multiverse-orchestrator` から AgentRunner Core を呼び出し
- 実際に Worker / Meta を使ったタスクが走るところまで

---

## 10. 受け入れ条件（Acceptance Criteria）【提案】

M1 完了の判断基準として、少なくとも以下を満たすこと。

1. **単一 Workspace / 単一プロジェクトでの操作**

   - 任意のローカルディレクトリを Workspace として登録できる
   - Workspace 一覧に登録済みプロジェクトが表示される

2. **Task 管理**

   - UI から Task を新規作成できる
   - Task 一覧で Status / PoolID / 日付が表示される
   - Task 詳細から、関連する Attempt の一覧が表示される

3. **Task 実行**

   - Run ボタン押下で Orchestrator が起動し、少なくともダミー Worker による SUCCEEDED / FAILED が返る
   - Task の Status が RUNNING → SUCCEEDED / FAILED に UI 上で正しく反映される

4. **永続化**

   - IDE を終了・再起動しても、Workspace / Task / Attempt の情報が保持されている
   - `~/.multiverse/workspaces/<workspace-id>/` 以下のファイル群が正しく更新されている

5. **安定性**

   - 想定外の Worker 終了（results が出ないケース）では Task が FAILED と判定される
   - IDE プロセスが異常終了しても、再起動後に Task 状態が破綻しない（少なくとも読み直しが可能）

---

以上を本 PRD v0.1 とし、
この内容をもとに「multiverse IDE UI 実装タスク」をエージェントに渡せる状態とする。
