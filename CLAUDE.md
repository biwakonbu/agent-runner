# CLAUDE.md

このファイルはこのリポジトリでClaudeコードを操作する際のガイダンスを提供します。

## プロジェクト概要

**agent-runner** は、開発タスクを自動実行するメタエージェント・オーケストレーションレイヤーです。以下を組み合わせています：
- **Meta-agent**: LLM（GPT-4 Turbo）がタスクを計画・評価
- **AgentRunner Core**: タスク状態と実行フローを管理するエンジン
- **Worker Agents**: 実際の開発作業（コーディング、テスト実行）を行うCLIツール（例：Codex CLI）

Docker サンドボックス環境で安全・再現性高く実行され、全ての実行履歴はMarkdown形式のTask Note（`.agent-runner/task-*.md`）として保存されます。

## プロジェクト構造

```
agent-runner/
├── cmd/
│   └── agent-runner/
│       └── main.go              # CLIエントリポイント
├── internal/
│   ├── core/                    # タスク実行エンジン ★詳細は core/CLAUDE.md を参照
│   │   ├── runner.go            # FSM オーケストレーション
│   │   ├── context.go           # TaskContext・TaskState定義
│   │   └── runner_test.go       # プロパティベーステスト
│   ├── meta/                    # Meta-agent通信層 ★詳細は meta/CLAUDE.md を参照
│   │   ├── client.go            # OpenAI API通信・モック実装
│   │   └── protocol.go          # YAMLプロトコル定義
│   ├── worker/                  # Worker実行・Dockerサンドボックス ★詳細は worker/CLAUDE.md を参照
│   │   ├── executor.go          # Worker CLI実行の抽象化
│   │   └── sandbox.go           # Docker API管理
│   ├── note/                    # Task Note生成
│   │   └── writer.go            # Markdown テンプレート出力
│   └── mock/                    # テスト用モック実装
│       ├── meta.go
│       ├── worker.go
│       └── note.go
├── pkg/
│   └── config/                  # 公開パッケージ（YAML設定）
│       └── config.go            # TaskConfig構造体定義
├── test/                         # ★詳細は test/CLAUDE.md を参照
│   ├── CLAUDE.md                # テスト戦略・実装パターン・精度管理
│   ├── integration/
│   │   └── run_flow_test.go     # Mock統合テスト
│   ├── sandbox/
│   │   ├── Dockerfile.test      # Docker Sandboxテスト用軽量イメージ
│   │   └── sandbox_test.go      # Docker API・コンテナ管理テスト（-tags=docker）
│   └── codex/
│       └── codex_integration_test.go  # Codex統合テスト（-tags=codex）
├── docs/                        # 設計・仕様・開発ガイド
│   ├── CLAUDE.md                # ドキュメント整理ルール
│   ├── AgentRunner-architecture.md
│   ├── agentrunner-spec-v1.md
│   ├── AgentRunner-impl-design-v1.md
│   ├── TESTING.md               # テストベストプラクティス
│   └── CODEX_TEST_README.md     # Codex統合ガイド
├── examples/                    # ★詳細は examples/CLAUDE.md を参照
│   ├── CLAUDE.md                # サンプル・スクリプト管理ガイド
│   ├── tasks/
│   │   ├── sample_task_go.yaml  # Goプロジェクト用サンプル
│   │   └── test_codex_task.yaml # Codex統合テスト定義
│   └── scripts/
│       └── run_codex_test.sh    # Codex統合テスト実行スクリプト
├── sandbox/
│   └── Dockerfile               # Worker実行環境（Codex CLI）
├── CLAUDE.md                    # このファイル（プロジェクトメモリ）
├── GEMINI.md                    # プロジェクト概要と背景（変更禁止）
├── README.md                    # ユーザー向け紹介（変更禁止）
├── go.mod                       # Goモジュール管理
├── go.sum
├── Makefile                     # ビルド・テストターゲット
└── .golangci.yml                # Code quality linter設定
```

### ディレクトリの役割

| ディレクトリ | 役割 | 詳細ドキュメント |
|-----------|------|-----------------|
| `cmd/` | CLIアプリケーションのエントリポイント | main.go のみ |
| `internal/core/` | タスクFSMとオーケストレーション | [internal/core/CLAUDE.md](internal/core/CLAUDE.md) |
| `internal/meta/` | Meta-agent（LLM）との通信層 | [internal/meta/CLAUDE.md](internal/meta/CLAUDE.md) |
| `internal/worker/` | Worker実行とDocker管理 | [internal/worker/CLAUDE.md](internal/worker/CLAUDE.md) |
| `internal/note/` | タスク実行履歴の Markdown出力 | [internal/note/CLAUDE.md](internal/note/CLAUDE.md) |
| `internal/mock/` | テスト用モック実装 | [internal/mock/CLAUDE.md](internal/mock/CLAUDE.md) |
| `pkg/config/` | YAML設定パース（再利用可能） | [pkg/config/CLAUDE.md](pkg/config/CLAUDE.md) |
| `test/` | **4段階テスト戦略**（ユニット→Mock→Docker→Codex） | **[test/CLAUDE.md](test/CLAUDE.md)** |
| `docs/` | 設計・仕様・開発ガイド統合 | [docs/CLAUDE.md](docs/CLAUDE.md) |
| `examples/` | サンプルタスク・実行スクリプト | [examples/CLAUDE.md](examples/CLAUDE.md) |
| `sandbox/` | Docker イメージ定義 | Codex CLI ランタイム |

## よく使うコマンド

### ビルド
```bash
# バイナリをビルド
go build ./cmd/agent-runner

# 出力ファイル名を指定
go build -o agent-runner ./cmd/agent-runner
```

### テスト

**詳細は [test/CLAUDE.md](test/CLAUDE.md) を参照してください。**

```bash
# ユニットテスト（依存なし、高速）
go test ./...

# 全テスト実行（推奨、Docker + Codex CLI 必須）
go test -tags=docker,codex -timeout=15m ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

test/CLAUDE.md には以下の情報が含まれています：
- **4段階テスト戦略**: ユニット → Mock統合 → Docker → Codex
- **実装パターン集**: Table-driven tests、PBT、Mock/Stub、Context-based testing
- **ビルドタグ戦略**: docker/codex タグによる段階的テスト実行
- **完全なコマンドリファレンス**: 日常開発・統合・カバレッジ・デバッグ
- **ベストプラクティス**: 8項目（エラーメッセージ、環境変数、タイムスタンプ等）
- **精度管理手法**: カバレッジ目標、不変条件、テストデータ生成、並行実行
- **トラブルシューティング**: Docker未起動、認証エラー、タイムアウト等

### 実行
```bash
# 基本的な実行（標準入力からYAMLを読み込む）
./agent-runner < task.yaml

# 設定ファイルをパイプ
cat task.yaml | go run cmd/agent-runner/main.go
```

### Docker（Worker実行環境）
```bash
# Codex worker のランタイムをビルド
docker build -t agent-runner-codex:latest sandbox/

# Codex統合テストを実行
./run_codex_test.sh
```

## アーキテクチャ

### 三層構造
```
Meta-agent (LLM)
    ↕ YAML プロトコル ↕
AgentRunner Core (状態 + TaskContext)
    ↕ Exec + Docker ↕
Worker Agents (Codex CLI等)
```

### タスク状態機械（FSM）
```
PENDING → PLANNING → RUNNING → VALIDATING → COMPLETE
                                              ↓
                                            FAILED
```

**主要な遷移**:
- `PENDING → PLANNING`: Meta-agentがPRDから受け入れ条件（AC）を生成
- `PLANNING → RUNNING`: Meta-agentがworker実行を要求
- `RUNNING → VALIDATING`: Worker実行が完了、Meta-agentが完了を評価
- `VALIDATING → COMPLETE/FAILED`: タスク完了判定を決定

### 主要パッケージ

#### `/internal/core` - タスク実行エンジン
- **runner.go**: タスクFSMのオーケストレーション、TaskContextの管理
- **context.go**: TaskContext（PRD、AC、worker実行、meta呼び出し）とTaskStateを定義
- `gopter`を使用した状態不変条件のプロパティベーステスト
- **詳細**: [internal/core/CLAUDE.md](internal/core/CLAUDE.md) - FSM遷移ルール、依存性注入パターン、エラーハンドリング戦略

#### `/internal/meta` - Meta-agent 通信層
- **client.go**: Meta-agent推論用のOpenAI API通信
- **protocol.go**: YAML メッセージ構造（`PlanTaskResponse`、`NextActionResponse`等）
- モック実装に対応（`kind: "mock"`）
- **詳細**: [internal/meta/CLAUDE.md](internal/meta/CLAUDE.md) - YAMLプロトコル仕様、LLM通信、モック実装

#### `/internal/worker` - Worker 実行とサンドボックス管理
- **executor.go**: Worker CLI実行の抽象化
- **sandbox.go**: Docker API管理（コンテナ作成、exec、クリーンアップ）
  - リポジトリは`/workspace/project`にマウント
  - 認証情報を自動マウント（`~/.codex/auth.json`）
- **詳細**: [internal/worker/CLAUDE.md](internal/worker/CLAUDE.md) - Dockerマウント戦略、認証管理、トラブルシューティング

#### `/internal/note` - Task Note 生成
- **writer.go**: TaskContextからMarkdown Task Noteを生成
- Go `text/template`によるテンプレートベース出力
- **詳細**: [internal/note/CLAUDE.md](internal/note/CLAUDE.md) - テンプレート設計、ファイルシステム操作、拡張ガイド

#### `/internal/mock` - テスト用モック実装
- **meta.go, worker.go, note.go**: Function Field Injection パターン
- 依存性注入を実現し、外部システムなしで テストを実行
- **詳細**: [internal/mock/CLAUDE.md](internal/mock/CLAUDE.md) - 設計パターン、テストケース、拡張ガイド

#### `/pkg/config` - 設定管理
- **config.go**: タスク設定構造体（YAMLパース）
- **詳細**: [pkg/config/CLAUDE.md](pkg/config/CLAUDE.md) - YAML スキーマ、バージョニング戦略、拡張ガイド

## タスク YAML 形式

```yaml
version: 1

task:
  id: "TASK-123"              # オプション（省略時は自動生成）
  title: "タスク名"           # オプション
  repo: "."                   # 作業対象リポジトリ（相対パスまたは絶対パス）

  prd:
    path: "./docs/prd.md"     # PRDファイルパス、または
    text: |                   # PRDテキストを直接埋め込み
      要件定義...

  test:
    command: "npm test"       # オプション（自動テストコマンド）
    cwd: "./"                 # オプション（テスト実行ディレクトリ）

runner:
  meta:
    kind: "openai-chat"       # または "mock"
    model: "gpt-4-turbo"      # LLMモデル名

  worker:
    kind: "codex-cli"
    docker_image: "agent-runner-codex:latest"
    max_run_time_sec: 1800    # 実行タイムアウト
    env:
      CODEX_API_KEY: "env:CODEX_API_KEY"  # 環境変数を参照
```

**参考**: [sample_task_go.yaml](sample_task_go.yaml)、[test_codex_task.yaml](test_codex_task.yaml)

## 重要な設計パターン

### 依存性注入（Dependency Injection）
`Runner`構造体はインターフェース（`MetaClient`、`WorkerExecutor`、`NoteWriter`）を受け入れるため、統合依存性なしでモックベースのテストが可能です。

```go
type Runner struct {
    Config *config.TaskConfig
    Meta   MetaClient        // インターフェース
    Worker WorkerExecutor    // インターフェース
    Note   NoteWriter        // インターフェース
}
```

### TaskContext の伝播
実行状態（PRD、受け入れ条件、worker結果、meta呼び出し、タイミング）全てが単一の`TaskContext`構造体で保持され、FSMを通じて伝播します。これにより完全な再現性と監査証跡（Task Noteに記録）が実現されます。

### Sandbox 隔離
- タスク1つあたり1つのDockerコンテナ（`/workspace/project`マウントポイント）
- 完了またはエラー時の自動クリーンアップ
- 環境変数注入と認証情報マウントをサポート

## 既知の課題

### 相対パスの解決
タスク設定で相対パス`"."`を使用するとDockerマウントエラーが発生します：
```
invalid mount path: '.' mount path must be absolute
```

**対応**: 絶対パスを使用するか、worker実行前にリポジトリパスを絶対パスに解決してください（`worker/executor.go`で`filepath.Abs`を使用）。

## 環境変数

```bash
# Meta-agent用（OpenAI API）
export OPENAI_API_KEY="sk-..."

# Worker agent用（auth.json未使用の場合）
export CODEX_API_KEY="..."
```

## サブディレクトリメモリ体系

**★ [docs/CLAUDE.md](docs/CLAUDE.md)** - ドキュメント整理ルール・命名規則・拡張ガイド（メンテナンス責務分担）

各パッケージ内 `CLAUDE.md` は標準フォーマット：**責務** → **主要概念** → **実装パターン** → **テスト戦略** → **拡張・既知問題**

### 実装層パッケージ

| 層 | パッケージ | 責務 | 詳細 |
|----|---------|----|------|
| **Core** | core | FSM・TaskContext・状態遷移 | [CLAUDE.md](internal/core/CLAUDE.md) |
| **Core** | meta | LLM通信・YAMLプロトコル | [CLAUDE.md](internal/meta/CLAUDE.md) |
| **Core** | worker | CLI実行・Docker サンドボックス | [CLAUDE.md](internal/worker/CLAUDE.md) |
| **Util** | note | Task Note生成・テンプレート | [CLAUDE.md](internal/note/CLAUDE.md) |
| **Util** | mock | テストダブル・FuncField注入 | [CLAUDE.md](internal/mock/CLAUDE.md) |
| **Config** | pkg/config | YAML設定スキーマ | [CLAUDE.md](pkg/config/CLAUDE.md) |

### テスト戦略（重要）

**[test/CLAUDE.md](test/CLAUDE.md)** - 4段階テスト戦略、ビルドタグ、精度管理、完全コマンドリファレンス

## ドキュメント体系

### AI開発者向け（メモリ・操作ガイド）
- **[CLAUDE.md](CLAUDE.md)** (このファイル) - プロジェクト全体・操作ガイド
- **[docs/CLAUDE.md](docs/CLAUDE.md)** ★ドキュメント整理ルール・命名規則・メンテナンス責務
- **各internal/*/CLAUDE.md** - パッケージ実装ガイダンス（標準フォーマット統一）
- **[test/CLAUDE.md](test/CLAUDE.md)** - テスト戦略・ビルドタグ・精度管理

### docs/ - 設計・仕様・開発ガイド
| ファイル | 対象 | 用途 |
|---------|------|------|
| [CLAUDE.md](docs/CLAUDE.md) | AI開発者 | ドキュメント整理ルール・命名規則 |
| [AgentRunner-architecture.md](docs/AgentRunner-architecture.md) | アーキテクト | システム設計・モジュール関係 |
| [agentrunner-spec-v1.md](docs/agentrunner-spec-v1.md) | 開発者・PO | MVP仕様・API・スキーマ |
| [AgentRunner-impl-design-v1.md](docs/AgentRunner-impl-design-v1.md) | 実装者 | Go実装設計・パッケージ分割 |
| [TESTING.md](docs/TESTING.md) | テスター・開発者 | テストベストプラクティス |
| [CODEX_TEST_README.md](docs/CODEX_TEST_README.md) | 開発者 | Codex統合テスト実行ガイド |

### エンドユーザー向け（不変）
- **[README.md](README.md)** - プロジェクト紹介（変更禁止）
- **[GEMINI.md](GEMINI.md)** - プロジェクト背景・コンテキスト（変更禁止）

## 開発ノート

- **言語**: コメント・ドキュメント は日本語、コード（関数・変数名）は英語
- **ドキュメント管理**: ★[docs/CLAUDE.md](docs/CLAUDE.md)で命名規則・責務分担を一元管理
- **テスト**: 依存性注入とモックを使用；プロパティベーステストで不変条件を検証
- **テスト自動化**: ビルドタグ（-tags=docker,codex）を使用した段階的テスト実行
- **ロギング**: 構造化ログ（log/slog）導入済み
- **依存関係**: Go 1.23以上、Docker、OpenAI API

## プロジェクトメモリ管理ルール

**CLAUDE.mdには静的な構造・設計情報のみを記載します。**

### 記載すべき内容
- プロジェクト構造とアーキテクチャ
- 設計パターンと実装ガイドライン
- よく使うコマンドと操作方法
- 既知の課題と対応方法
- サブディレクトリメモリ体系

### 記載しない内容
- Phase履歴や開発進捗
- カバレッジ推移や改善履歴
- 完了時の成果テーブル
- 「✅完了」「次フェーズの予定」等の進捗状況

**理由**: 進捗情報は時間とともに陳腐化し、メモリの可読性を低下させます。履歴情報はGitコミットログで管理してください。
