# Task Builder & Golden Test 設計書

本ドキュメントは、Multiverse IDE における「チャット入力 → TaskConfig YAML 生成 → AgentRunner 実行 → 結果反映」までの最小パイプラインと、ゴールデンテスト（`TODO アプリを作成して`）の仕様を定義する。

実装時の指示書として利用することを前提とする。

---

## 1. 背景・目的

- ユーザーは **チャット UI から自然文を入力**してタスクを起動する。
- 内部では、その自然文をもとに **TaskConfig YAML** を生成し、それを AgentRunner に渡す。
- AgentRunner は、タスク分析・実装・ファイル生成・検証（テスト等）までを実行し、その結果を IDE に返す。
- Phase 0 のゴールは、以下の 1 本のパイプラインが「ローカルで一気通しで動作すること」である。

> Chat（`TODO アプリを作成して`）  
> → Task Builder（LLM）で TaskConfig YAML を生成  
> → AgentRunner でタスク分析 / 実装 / ファイル生成 / 検証を実行  
> → Orchestrator 経由で結果が IDE に表示される

TODO アプリの仕様・技術スタック・テスト戦略などは **一切固定しない**。  
本ドキュメントの範囲は「パイプラインとしての契約と責務」のみを定義する。

---

## 2. コンポーネントと責務

### 2.1 IDE (Chat Layer)

- ユーザーと対話するフロントエンド。
- ユーザー入力（自然文）を `raw_prompt` として TaskStore に保存する。
- Task の一覧表示、ステータス表示、結果サマリの表示を行う。

### 2.2 Orchestrator

- Workspace / TaskStore / IPC の管理を行うバックエンドコンポーネント。
- 主な責務:
  - Task 作成時に TaskStore レコードを生成。
  - Task 実行要求を IPC queue にジョブとして登録。
  - **Task Builder**（LLM）を呼び出し、`raw_prompt` から TaskConfig YAML を生成。
  - TaskConfig YAML を AgentRunner に渡して実行。
  - AgentRunner の結果を受け取り、TaskAttempt として保存し、IPC 結果を IDE に露出。

### 2.3 Task Builder（CLI プロバイダ）

- Orchestrator から呼び出される CLI ベースのコンポーネント（Codex CLI 等）。
- 入力:
  - Workspace 情報（root_dir 等）
  - ユーザー入力の自然文（`raw_prompt`）
- 出力:
  - AgentRunner に渡す **TaskConfig YAML**（本ドキュメントでフォーマットを定義）。
- 実装:
  - `codex chat` コマンドを実行し、JSON 形式で応答を受信
  - CLI セッションの検証とエラーハンドリングを実装

### 2.4 AgentRunner

- Meta / Worker エージェントのランタイム。
- 入力:
  - TaskConfig YAML（Task Builder の出力をそのまま受け取る）
- 処理:
  - タスク分析・プランニング
  - コード編集・新規ファイル生成
  - 可能な範囲での検証（テスト実行・ビルド・lint 等）
- 出力:
  - Task 実行結果の JSON（タスクサマリ・検証内容・ステータス等）。

### 2.5 TaskStore / Workspace

- ローカルファイルシステム上のメタデータ保存レイヤ。
- ディレクトリ構造（概要）:

```text
~/.multiverse/workspaces/<workspace-id>/
  workspace.json
  tasks/
    <task-id>.jsonl
  ipc/
    queue/
    results/
  logs/
```

---

## 3. データモデル

### 3.1 TaskStore: Task レコード

IDE から作成されるタスクの最小レコード定義。

```jsonc
// ~/.multiverse/workspaces/<workspace-id>/tasks/<task-id>.jsonl
{
  "id": "golden-todo-001",
  "workspace_id": "abcd1234ef56",
  "title": "TODO アプリを作成して",
  "raw_prompt": "TODO アプリを作成して",
  "created_at": "2025-12-07T08:00:00Z"
}
```

- `title`
  - UI 表示用。初期値は `raw_prompt` と同一でよい。
- `raw_prompt`
  - ユーザーがチャットで入力した自然文。
  - Task Builder の入力として利用する。

※ Phase 0 では `test_command` 等は持たない。検証戦略は AgentRunner 側に委譲する。

### 3.2 IPC Queue: Job JSON

IDE からの「実行してほしい」要求は、Orchestrator に対して IPC queue 経由で渡される。

```jsonc
// ~/.multiverse/workspaces/<workspace-id>/ipc/queue/<job-id>.json
{
  "workspace_id": "abcd1234ef56",
  "task_id": "golden-todo-001"
}
```

- Orchestrator は queue ディレクトリをポーリングし、Job を検出して処理する。

### 3.3 TaskConfig YAML（Task Builder 出力 / AgentRunner 入力）

Task Builder により生成され、AgentRunner に渡される YAML の最小スキーマを定義する。

```yaml
task:
  id: "golden-todo-001"
  title: "TODO アプリを作成して"
  instructions: |
    TODO アプリを作成して。
    技術スタックや実装方針、検証方法はあなたの判断に任せます。
    必要に応じて、コードや設定ファイル、テストコードなどを自由に生成してください。
  project:
    root_dir: "/path/to/workspace"

runner:
  meta:
    model: "gpt-4.1"
    temperature: 0.2
  worker:
    type: "docker_codex_cli"
    # 必要に応じて image / mount / env 等を拡張
```

必須フィールド:

- `task.id`（TaskStore の `id` と一致）
- `task.title`
- `task.instructions`
- `task.project.root_dir`
- `runner.meta.model`
- `runner.worker.type`

Task Builder 実装は、このスキーマを満たすように LLM 出力を誘導する。

### 3.4 AgentRunner 結果 JSON

AgentRunner がタスク実行完了時に Orchestrator に返す結果 JSON の最小仕様。

```jsonc
{
  "task_id": "golden-todo-001",
  "status": "succeeded",   // "succeeded" | "failed"
  "summary": "TODO アプリを作成し、基本的な追加・削除・一覧機能と簡単な検証処理を実行しました。",
  "validation": {
    "overall": "passed",   // "passed" | "failed" | "unknown"
    "commands": [
      {
        "command": "npm test",
        "exit_code": 0,
        "duration_ms": 12345
      }
    ]
  },
  "duration_ms": 600000
}
```

- `status`
  - AgentRunner レベルでの成功/失敗。
- `summary`
  - 実装内容の自然文サマリ（IDE 表示用）。
- `validation`
  - AgentRunner 内で実施した検証（テスト / ビルド / lint 等）の概要。
  - Phase 0 では 1 コマンド / 0 コマンドでも可（`commands` は空配列を許容）。
- `duration_ms`
  - 全体の実行時間（任意だが、あると便利）。

Orchestrator は、本 JSON を TaskAttempt（JSONL）に埋め込み、IDE から参照可能にする。

---

## 4. 処理フロー

### 4.1 Chat → Task 作成

1. ユーザーが IDE のチャット欄に以下を入力する。

   > `TODO アプリを作成して`

2. IDE は以下を行う。
   - Workspace を選択中であることを前提に、`workspace_id` を決定。
   - 新規 Task を作成し、`title` と `raw_prompt` に上記の文言を保存。
   - TaskStore の `tasks/<task-id>.jsonl` に Task レコードを append。

3. ユーザーは Task 一覧画面で、`TODO アプリを作成して` タスクを確認できる。

### 4.2 Task 実行要求 → Orchestrator

1. ユーザーが IDE 上で Task の「Run」ボタンを押下。
2. IDE は IPC queue に Job JSON を作成する（3.2 参照）。
3. Orchestrator は queue ディレクトリを監視し、Job を検出。

### 4.3 Task Builder 呼び出し

1. Orchestrator は TaskStore から `task_id` に対応する Task をロード。
2. Orchestrator は CLI プロバイダ（Task Builder）に以下の情報を渡して呼び出す。

   - Workspace 情報（例）
     - `root_dir: "/path/to/workspace"`
   - `raw_prompt: "TODO アプリを作成して"`

3. Task Builder（Codex CLI 等）は、`codex chat` コマンドを実行し、3.3 の TaskConfig スキーマに従う YAML を生成する。
4. Orchestrator は生成された YAML を TaskConfig として検証する。
   - YAML としてパース可能か
   - 必須フィールドが存在するか
5. 検証に失敗した場合、または CLI セッションが無い場合は、その時点で TaskAttempt を `failed` として記録し、結果を IDE に返す。

### 4.4 AgentRunner 実行

1. Orchestrator は検証済み TaskConfig YAML を AgentRunner に渡す（実装としては `agent-runner` サブプロセスの stdin 等）。
2. AgentRunner は内部で以下を行う（振る舞いは AgentRunner 側の設計に従う）:
   - タスク分析・プランニング
   - コード編集・ファイル生成
   - 可能な限りの自己検証（テスト / ビルド / lint 等）
3. 完了時、AgentRunner は 3.4 の JSON を stdout（またはファイル）として出力する。
4. Orchestrator はこの JSON を受け取り、TaskAttempt として TaskStore に追記し、IPC results にも書き出す。

### 4.5 IDE での結果表示

1. IDE は IPC results をポーリング or ファイル監視し、対象 Job の result JSON を検出。
2. Task 一覧画面:
   - 対象 Task のステータスを `SUCCEEDED` / `FAILED` に更新。
3. Task 詳細画面:
   - `status` / `summary` / `validation.overall` / `validation.commands` 等を表示する。

---

## 5. ゴールデンテスト仕様

### 5.1 前提

- ゴールデンテストのユーザー入力は **固定** とする。

  ```text
  TODO アプリを作成して
  ```

- TODO アプリの解釈・技術スタック・設計・テスト戦略に関するルールは **一切課さない**。
- 検証対象は「アプリとして妥当か」ではなく、「パイプラインとして正しく通るか」である。

### 5.2 GT-1: Chat → TaskConfig（Task Builder テスト）

目的:

- `raw_prompt = "TODO アプリを作成して"` から **有効な TaskConfig YAML** が生成されることを確認する。

前提条件:

- Codex CLI がインストールされ、有効なセッションが存在すること

テスト手順（ロジック）:

1. テスト用 Workspace を作成（空 or ほぼ空でよい）。
2. Task を作成し、`raw_prompt = "TODO アプリを作成して"` を設定。
3. Orchestrator 経由で Task Builder（Codex CLI）を呼び出し、TaskConfig YAML を取得。
4. アサーション:
   - YAML としてパース可能。
   - `task.id` が TaskStore の `id` と一致。
   - `task.title` が `TODO アプリを作成して` を含む。
   - `task.instructions` に `TODO アプリを作成して` の文言が含まれる。
   - `task.project.root_dir` が Workspace のパスと一致。
   - `runner.meta.model` / `runner.worker.type` が存在。

### 5.3 GT-2: TaskConfig → AgentRunner（実行テスト）

目的:

- TaskConfig YAML を AgentRunner に渡した際、実装・ファイル生成・自己検証までの処理が完了し、結果 JSON が返ることを確認する。

前提条件:

- Codex CLI がインストールされ、有効なセッションが存在すること
- Docker が起動しており、Codex Worker イメージが利用可能であること

テスト手順（ロジック）:

1. GT-1 で取得した TaskConfig YAML をそのまま AgentRunner に入力。
2. AgentRunner を実行し、結果 JSON（3.4）を取得。
   - AgentRunner は Docker サンドボックス内で Codex CLI を実行
   - Codex CLI セッションが Docker コンテナ内で利用可能であることを確認
3. アサーション:
   - プロセスとして正常終了している（exit code = 0 が望ましいが、結果 JSON の `status` を見て判定）。
   - Workspace ディレクトリ内で 1 つ以上のファイルが新規作成 or 更新されている。
   - 結果 JSON に以下が含まれる:
     - `task_id`（TaskStore の id と一致）
     - `status`（"succeeded" or "failed"）
     - `summary`（非空の文字列）
     - `validation` オブジェクト（存在すればよい。`commands` が空でも許容）

※ Phase 0 の時点では、`status = failed` であっても、「パイプラインとして最後まで処理され、結果が返る」ことを成功条件としてよい。

### 5.4 GT-3: E2E（Chat → TaskConfig → AgentRunner → 結果）

目的:

- IDE チャット入力から結果表示まで、全パスが一気通しで動くことを確認する。

テスト手順（ロジック）:

1. IDE のテストモードで以下を実行する:
   - Chat に `TODO アプリを作成して` を入力し、Task 作成。
   - Task の「Run」ボタンを押下。
2. バックグラウンドで:
   - Orchestrator が Job 処理 → Task Builder → TaskConfig YAML 生成 → AgentRunner 実行 → 結果 JSON 生成。
3. IDE で Task 詳細画面を開き、以下を確認:
   - ステータスが `SUCCEEDED` または `FAILED` のいずれか。
   - summary が表示されている。
   - validation.overall が `passed` / `failed` / `unknown` のいずれか（存在すればよい）。

---

## 6. 実装順序（Phase 0 向け指針）

実装順序の推奨:

1. Workspace / TaskStore / IPC（queue/results）の基盤実装。
2. IDE:
   - Workspace 選択 UI
   - Task 作成 UI（Chat 入力 → TaskStore に `raw_prompt` 保存）
3. Orchestrator:
   - Job queue 処理
   - TaskBuilder 呼び出し（LLM API ラッパ）
   - TaskConfig YAML 検証
4. AgentRunner 連携:
   - TaskConfig YAML を stdin で渡す Executor 実装
   - 結果 JSON の受信と TaskAttempt への保存
5. IDE:
   - IPC results の監視
   - Task ステータスと結果サマリの表示
6. ゴールデンテスト（GT-1 / GT-2 / GT-3）の追加

本設計書は Phase 0 の最小スコープを対象とする。  
Phase 1 以降で、複数エージェント、WorkerPool、シナリオベースの L2 テスト等を拡張するが、それらは別途仕様書で定義する。
