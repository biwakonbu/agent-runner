# CLI サブスクリプション運用ガイド

このガイドでは、AgentRunner で各種 CLI プロバイダ（Codex, Gemini, Claude 等）を利用するためのセットアップと運用方法について説明します。

## 1. 概要

AgentRunner (v1) は、ローカル環境にインストールされた各社 AI アシスタント CLI ツールを「Worker」として利用します。API キーを直接埋め込むのではなく、CLI ツールが保持する認証セッション（サブスクリプション）を再利用することで、以下のメリットがあります：

- **セキュア**: API キーの管理が不要。
- **コスト効率**: 既存の Pro/Enterprise サブスクリプションを活用可能。
- **一貫性**: 開発者が普段使用している CLI と同じモデル・設定を利用可能。

## 2. 対応プロバイダ状況

| プロバイダ              | 対応状況        | 備考                |
| :---------------------- | :-------------- | :------------------ |
| **Codex CLI** (`codex`) | ✅ **対応済み** | v1 のデフォルト推奨 |
| **Gemini CLI**          | ✅ **対応済み** | `gemini-cli`        |
| **Claude Code**         | 🚧 準備中       | スタブのみ実装      |
| **Cursor CLI**          | 🚧 準備中       | スタブのみ実装      |

## 3. Codex CLI のセットアップ

現在、メインでサポートされている Codex CLI の利用手順です。

### 3.1 インストールとログイン

事前に `codex` コマンドが利用可能である必要があります。

```bash
# 1. ログイン（ブラウザが開きます）
codex login

# 2. 接続確認
codex auth status
# -> "Logged in as user@example.com (Pro)" のように表示されればOK
```

### 3.2 AgentRunner での利用

AgentRunner は実行時に自動的にローカルの CLI セッションを検出します。特別な設定は不要です。

#### 実行モデルの指定

デフォルトでは `agent-runner` 側で設定されたデフォルトモデルが使用されますが、CLI オプションで明示的に指定することも可能です。

```bash
# 特定のモデルを指定して実行（CLI オプションは Task YAML やデフォルトより優先されます）
agent-runner --meta-model="gpt-5.1" < task.yaml
```

## 4. Gemini CLI のセットアップ

### 4.1 前提条件

`gemini` コマンドがパスに通っており、CLI からの実行が可能である必要があります。
（例: Node.js ベースの CLI ツール等）

### 4.2 認証

使用する CLI ツールに合わせて認証を行ってください。一般的には API キーを環境変数 `GOOGLE_API_KEY` に設定するか、`gemini login` 等のコマンドを使用します。

### 4.3 AgentRunner での利用

Task YAML で `runner.worker.kind: "gemini-cli"` を指定します。

```yaml
runner:
  worker:
    kind: "gemini-cli"
    # model: "gemini-1.5-pro" # 任意。デフォルトは gemini-1.5-pro
```

## 5. トラブルシューティング

### Q. "CLI session not found" エラーが出る

**原因**: `codex login` が行われていないか、セッションが切れています。
**対策**: 再度 `codex login` を実行してください。

### Q. コンテナ内で CLI が動かない

**原因**: Docker マウントの設定不備の可能性があります。
**対策**: AgentRunner はデフォルトで `~/.config/codex` 等の認証ディレクトリをコンテナにマウントします。ホスト側のパスが標準と異なる場合、正しく動作しない可能性があります。

## 6. 今後の予定

Claude Code, Cursor などの他プロバイダについても、順次実装を進めていく予定です。
