# examples/ - サンプルタスク定義・実行スクリプト

このディレクトリはプロジェクトの**実行例・テスト用サンプル・参考スクリプト**を管理します。

## ディレクトリ構成

```
examples/
├── CLAUDE.md                    # このファイル（管理ガイド）
├── tasks/
│   ├── sample_task_go.yaml      # Goプロジェクト用サンプルタスク定義
│   ├── test_codex_task.yaml     # Codex統合テスト用タスク定義
│   └── (将来) sample_python.yaml, sample_node.yaml など
└── scripts/
    ├── run_codex_test.sh        # Codex統合テスト実行スクリプト
    └── (将来) debug_sandbox.sh, generate_report.sh など
```

## ファイル責務

### tasks/ - タスク定義サンプル

| ファイル | 用途 | 対象読者 | 更新頻度 |
|---------|------|--------|--------|
| **sample_task_go.yaml** | Go開発タスク実行例・ドキュメント | 開発者・ユーザー | スキーマ変更時 |
| **test_codex_task.yaml** | Codex統合テスト用タスク定義 | テスター・CI/CD | Codex仕様変更時 |

#### sample_task_go.yaml

**目的**：
- プロジェクト README・ドキュメントの参考例
- 実行可能な完全な設定例
- YAML スキーマの説明用

**含まれるもの**：
```yaml
version: 1
task:
  id: "TASK-GO-001"
  title: "Go project task"
  repo: "/path/to/go/project"  # 絶対パス推奨
  prd:
    text: |
      Create calculator.go with Add(), Subtract(), Multiply(), Divide() functions
  test:
    command: "go test ./..."
    cwd: ""
runner:
  meta:
    kind: "openai-chat"  # または "mock"
    model: "gpt-4-turbo"
  worker:
    kind: "codex-cli"
    docker_image: "agent-runner-codex:latest"
    max_run_time_sec: 1800
    env:
      CODEX_API_KEY: "env:CODEX_API_KEY"
```

#### test_codex_task.yaml

**目的**：
- Codex統合テスト実行時のデファクトスタンダード定義
- CI/CD パイプラインで使用
- 簡潔で確実性の高い動作例

**含まれるもの**：
- 単純な calculator.py 生成タスク
- 認証・Docker・Codex CLI の正常動作確認
- ファイル保存確認まで含む

### scripts/ - 実行・管理スクリプト

| ファイル | 用途 | 対象読者 | 使用場面 |
|---------|------|--------|--------|
| **run_codex_test.sh** | Codex統合テスト実行スクリプト | テスター・CI/CD | テスト実行時 |

#### run_codex_test.sh

**機能**：
- Docker イメージビルド
- テスト用タスク YAML の読み込み
- agent-runner コマンド実行
- テスト結果検証

**使用法**：
```bash
# 直接実行
./run_codex_test.sh

# または go test で実行（-tags=codex）
go test -tags=codex -timeout=10m ./test/codex/...
```

## 配置ルール

### サンプルタスク定義をどこに置くか

| パターン | 配置先 | 理由 |
|---------|-------|------|
| **汎用サンプル** | `examples/tasks/sample_*.yaml` | 再利用可能・複数言語対応 |
| **統合テスト定義** | `examples/tasks/test_*.yaml` | テスト用・固定仕様 |
| **単一プロジェクト用** | `examples/tasks/project-{name}.yaml` | プロジェクト別整理 |

### 実行スクリプト配置ルール

| パターン | 配置先 | 理由 |
|---------|-------|------|
| **テスト実行** | `examples/scripts/run_*.sh` | テスト関連・非本番 |
| **デバッグ・開発** | `examples/scripts/debug_*.sh` | 開発者用・非本番 |
| **本番実行** | `cmd/agent-runner/` + ドキュメント | メインツール |

## 更新ルール

### サンプル定義を更新する場合

**トリガー**：
1. **YAML スキーマ変更** → サンプルを同期更新
2. **新言語・新パターン対応** → sample_*.yaml 追加
3. **テスト手法改善** → test_*.yaml 更新

**手順**：
```bash
# 1. サンプル更新
vi examples/tasks/sample_task_go.yaml

# 2. 実行確認
./agent-runner < examples/tasks/sample_task_go.yaml

# 3. ドキュメント更新（変更があれば）
# - docs/CLAUDE.md
# - README.md（サンプル説明）
```

### スクリプト追加時の手順

1. `examples/scripts/{purpose}_{context}.sh` に作成
2. 実行権限付与：`chmod +x examples/scripts/new_script.sh`
3. this CLAUDE.md の scripts セクション更新
4. 必要に応じて docs/ に実行手順ドキュメント追加

## 拡張ガイド

### 新しいサンプルタスク追加

**想定シーン**：Python プロジェクト用サンプル追加

```bash
# 1. ファイル作成
touch examples/tasks/sample_task_python.yaml

# 2. テンプレート（sample_task_go.yaml を参考に）
cat > examples/tasks/sample_task_python.yaml << 'EOF'
version: 1
task:
  id: "TASK-PYTHON-001"
  title: "Python project task"
  repo: "/path/to/python/project"
  prd:
    text: |
      Create utils.py with helper functions
  test:
    command: "pytest tests/"
    cwd: ""
runner:
  meta:
    kind: "openai-chat"
    model: "gpt-4-turbo"
  worker:
    kind: "codex-cli"
    docker_image: "agent-runner-codex:latest"
    max_run_time_sec: 1800
    env:
      CODEX_API_KEY: "env:CODEX_API_KEY"
EOF

# 3. this CLAUDE.md を更新
# - tasks/ テーブルに追加
# - 説明文追加

# 4. 検証（必要に応じて）
./agent-runner < examples/tasks/sample_task_python.yaml
```

### デバッグスクリプト追加例

```bash
# デバッグスクリプト例：sandbox 動作確認
cat > examples/scripts/debug_sandbox.sh << 'EOF'
#!/bin/bash
set -e

echo "Testing Docker Sandbox..."
docker ps -q
echo "✓ Docker daemon active"

echo "Testing image availability..."
docker images | grep agent-runner || echo "⚠ Image not found - run: docker build -t agent-runner-codex:latest sandbox/"

echo "Testing volume mount..."
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT
docker run --rm -v "$TEST_DIR:/workspace/project" alpine echo "✓ Mount successful"

echo "All checks passed!"
EOF

chmod +x examples/scripts/debug_sandbox.sh
```

## 既知問題・注意点

### 相対パス・絶対パス

サンプルで相対パス `.` を使用するとDocker マウントエラーが発生します。

```yaml
repo: "."     # ❌ エラー：invalid mount path
repo: "/home/user/project"  # ✓ 絶対パス推奨
```

### 環境変数参照

```yaml
env:
  CODEX_API_KEY: "env:CODEX_API_KEY"  # ✓ 実行時に環境変数から読み込み
  CODEX_API_KEY: "sk-xxx"             # ❌ 平文で埋め込まない
```

## 関連リンク

- [docs/CODEX_TEST_README.md](../docs/CODEX_TEST_README.md) - Codex統合テスト詳細
- [docs/TESTING.md](../docs/TESTING.md) - テストベストプラクティス
- [CLAUDE.md](../CLAUDE.md) - プロジェクトメモリ
