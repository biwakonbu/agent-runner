# TODO: multiverse v3.0 - Phase 4 & 4.5 Implementation

Based on PRD v3.0 - Codex CLI 統合と実タスク実行 + Svelte 5 移行

---

## 現在のステータス

| フェーズ  | 内容                             | ステータス |
| --------- | -------------------------------- | ---------- |
| Phase 1   | チャット → タスク生成            | ✅ 完了    |
| Phase 2   | 依存関係グラフ・WBS 表示         | ✅ 完了    |
| Phase 3   | 自律実行ループ                   | ✅ 完了    |
| Phase 4   | CLI セッション統合・実タスク実行 | ✅ 完了    |
| Phase 4.5 | Svelte 5 + Svelte Flow 移行      | ✅ 完了    |

---

## 設計方針（重要・現状差分あり）

> [!IMPORTANT]
> API キーは不要。Codex / Claude Code / Gemini / Cursor など **CLI サブスクリプションセッションを優先利用**する。Meta 層も CLI セッション前提に置き換え、API キー依存を排除する。

**現在のデータフロー（実装ベース）:**

```
Chat → Meta-agent (openai-chat via HTTP + OPENAI_API_KEY) → Task 生成
                                                            ↓
ExecutionOrchestrator → agent-runner → Docker Sandbox → codex CLI（既存セッション想定）
```

---

## 現在の実装メモ（2025-12-07 時点）

### バックエンド

- [x] **LLMConfigStore** (`internal/ide/llm_config.go`)
  - Kind/Model/BaseURL/SystemPrompt を `~/.multiverse/config/llm.json` に永続化
  - 環境変数オーバーライドあり（API キー保存は不要にする方針）
- [x] **App API** (`app.go`)
  - `GetLLMConfig` / `SetLLMConfig` / `TestLLMConnection` を追加
  - ただし **ChatHandler 生成は `newMetaClientFromEnv()` 固定**で LLMConfigStore の設定が Meta 層に反映されない
  - `TestLLMConnection` は OpenAI API キー前提の HTTP 呼び出し（API キー不要の CLI セッション検証に置換予定）
- [x] **AgentToolProvider 基盤** (`internal/agenttools`)
  - 共通 Request/ExecPlan/ProviderConfig と registry を追加
  - Codex CLI プロバイダ実装（exec/chat、model/temperature/max-tokens/flags/env を透過）
  - Gemini / Claude Code / Cursor は stub プロバイダで登録（未実装アラートのみ）
- [x] **Worker Executor**
  - `RunWorker` → `RunWorkerCall` に内部委譲し、AgentToolProvider 経由で ExecPlan を構築して Sandbox.Exec 実行
  - `meta.WorkerCall` に model/flags/env/tool_specific/use_stdin などを拡張し、CLI 切替の土台を用意
  - stdin 実行は未サポート（現在はエラーにする）

### フロントエンド

- [x] **LLMSettings** (`frontend/ide/src/lib/settings/LLMSettings.svelte`)
  - プロバイダ選択、モデル/エンドポイント入力、接続テスト UI
  - API キーは「環境変数に設定済みか」を表示するのみ（保存不可）
- [x] **Toolbar 設定ボタン & モーダル** (`Toolbar.svelte`, `App.svelte`)
  - 設定モーダルから LLMSettings を呼び出し

### ビルド検証

- [x] `go build .`
- [x] `pnpm build`（警告 5 件、エラー 0）
- [x] `pnpm check`

---

## 残りのタスク（優先度順）

### 完了済み（Phase4 実装要点）

- [x] Meta/LLM: LLMConfigStore 経由で `codex-cli` 初期化、接続テストを CLI セッション検証に変更
- [x] Worker: コンテナ起動前に Codex セッション検証を強制し、未ログインなら IDE へエラー通知して中断
- [x] Orchestrator: 実行ログを `task:log` イベントでストリーミング
- [x] UI: LLMSettings を CLI セッション表示に対応（codex-cli 選択可）
- [x] Doc: PRD/TODO/Golden テスト設計を CLI 前提に更新

### Phase 4 完了タスク

- [x] CLI サブスクリプション運用手順を GEMINI.md / CLAUDE.md / guides に追記
- [x] Sandbox Exec で stdin 入力をサポートし、AgentToolProvider の UseStdin を有効化（確認済み）
- [x] Gemini / Claude Code / Cursor の実プロバイダを実装し、registry stub を置換
  - `internal/agenttools/claude.go`
  - `internal/agenttools/cursor.go`
- [x] Meta 層からの WorkerCall 生成で新フィールド（model/flags/env/tool_specific）を活用する経路を整備（確認済み）

### 残タスク（オプション・フォローアップ）

- [ ] E2E: CLI セッション未設定時の IDE 通知を含む回帰テストを追加
- [ ] CLI 未ログイン時の IDE 通知と再試行 UX の改善（案内リンク・ボタン）

---

## 設計上の注意点

### Codex / CLI 統合（現状）

1. **Meta-agent (decompose)**: `internal/meta/client.go` が HTTP で OpenAI Chat Completion を呼び出す（`OPENAI_API_KEY` 必須）。CLI サブスクリプション非対応。
2. **Worker (codex-cli)**: `internal/worker/executor.go` が Docker サンドボックス内で `codex exec ...` を実行。CLI セッション引き継ぎ方法は未整備。

### セッション/環境（現状）

| 項目                    | 用途                                     | 備考                           |
| ----------------------- | ---------------------------------------- | ------------------------------ |
| `MULTIVERSE_META_KIND`  | Meta-agent の種別                        | 現状: mock / openai-chat       |
| `MULTIVERSE_META_MODEL` | Meta-agent のモデル                      | 現状: gpt-5.1                  |
| CLI セッション          | Codex / Claude Code / Gemini / Cursor 等 | **API キー不要。要セッション** |

---

## 次のアクション

1. Meta 層を CLI セッション対応に変更する設計・実装方針を決定（AgentToolProvider と整合）
2. `agent-runner` + worker へ CLI セッションを確実に引き継ぐ仕組みを確認（env/マウント/cli path）
3. `go test ./internal/ide/...` 実行で LLMConfigStore の回帰確認
4. ストリーミングログと CLI ベース接続の E2E テストを追加

---

## 追加で必要な対応（漏れ防止メモ）

- [x] CLI サブスクリプション運用手順のドキュメント化（`docs/guides/cli-subscription.md` 作成済み）
- [ ] CLI 未ログイン時の IDE 通知と再試行 UX の改善（案内リンク・ボタン）

---

## Phase 4.5: Svelte 5 + Svelte Flow 移行

### 背景・目的

現在のグラフノード管理（GridCanvas/WBSGraphView）は手実装で以下の課題がある：

- **大量ノード非対応**: 全ノードを常時描画、2000+ タスクでパフォーマンス劣化
- **レイアウト最適化なし**: 単純な列配置、依存関係を考慮しない
- **保守コスト高**: パン/ズーム/エッジ描画を全て自前実装

**解決策**: Svelte 5 へアップグレードし、Svelte Flow (@xyflow/svelte v1.5+) を導入

### Svelte 5 移行タスク

#### Step 1: 依存パッケージ更新

```bash
cd frontend/ide
pnpm install svelte@^5 @sveltejs/vite-plugin-svelte@^4 --save-dev
```

- [x] svelte: ^4.2.12 → ^5.0.0 (確認済み `^5.0.0`)
- [x] @sveltejs/vite-plugin-svelte: ^3.0.2 → ^4.0.0 (確認済み `^4.0.0`)
- [x] vite: 維持（^5.x）
- [x] typescript: 維持（^5.x）

#### Step 2: 自動移行ツール実行

```bash
npx sv migrate svelte-5
```

**自動変換される内容:**

- `let` → `$state`
- `$:` (派生) → `$derived`
- `export let` → `$props`

**手動変換が必要な内容:**

- `createEventDispatcher` → コールバックプロップ（約 10 ファイル）
- `beforeUpdate`/`afterUpdate` → `$effect.pre`/`$effect`
- 複雑な `$:` の `$effect` vs `$derived` 判別

#### Step 3: createEventDispatcher 置き換え

**対象ファイル（要手動変換）:**

| ファイル                    | dispatch イベント | 変換後                 |
| --------------------------- | ----------------- | ---------------------- |
| `FloatingChatWindow.svelte` | close             | `onClose` コールバック |
| `ChatInput.svelte`          | send              | `onSend` コールバック  |
| `TaskDetail.svelte`         | close             | `onClose` コールバック |
| `Modal.svelte`              | close             | `onClose` コールバック |
| その他約 6 ファイル         | 各種              | 各コールバック         |

**変換例:**

```svelte
// Before (Svelte 4)
<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  function close() { dispatch('close'); }
</script>

// After (Svelte 5)
<script>
  let { onClose } = $props();
  function close() { onClose?.(); }
</script>
```

#### Step 4: テスト実行・修正

- [x] `pnpm check` パス (0 errors, 7 warnings)
- [x] `pnpm build` パス
- [ ] `pnpm test` パス（該当する場合）
- [ ] 手動で全画面動作確認

### Svelte Flow 移行タスク

#### Step 5: パッケージインストール

```bash
cd frontend/ide
pnpm add @xyflow/svelte dagre
pnpm add -D @types/dagre
```

#### Step 6: 新規ファイル作成

```
frontend/ide/src/lib/flow/
├── CLAUDE.md                        # 設計ガイド
├── index.ts                         # エクスポート集約
├── UnifiedFlowCanvas.svelte         # 統合キャンバス
├── nodes/
│   ├── TaskFlowNode.svelte          # タスクノード
│   ├── WBSFlowNode.svelte           # WBS ノード
│   ├── MilestoneFlowNode.svelte     # マイルストーン
│   └── index.ts
├── edges/
│   ├── DependencyEdge.svelte        # 依存エッジ
│   └── index.ts
├── layout/
│   ├── dagreLayout.ts               # Dagre 統合
│   ├── layoutStore.ts               # レイアウト状態
│   └── index.ts
└── utils/
    ├── nodeConverter.ts             # Task → FlowNode 変換
    ├── edgeConverter.ts             # Edge 変換
    └── constants.ts                 # サイズ定数

frontend/ide/src/stores/
└── flowStore.ts                     # Svelte Flow 用ストア
```

#### Step 7: カスタムノード実装

- [x] `TaskNode.svelte` - GridNode.svelte のスタイルを移植（`lib/flow/nodes/TaskNode.svelte`）
- [x] `DependencyEdge.svelte` - ConnectionLine.svelte のスタイルを移植（`lib/flow/edges/DependencyEdge.svelte`）
- [ ] `WBSFlowNode.svelte` - WBSGraphNode.svelte のスタイルを移植（WBS 切り替えは UnifiedFlowCanvas 内部で対応）
- [ ] `MilestoneFlowNode.svelte` - マイルストーン表示（将来対応可）

#### Step 8: Dagre レイアウト統合

- [x] `dagreLayout.ts` - Dagre による自動レイアウト計算（`lib/flow/dagreLayout.ts` 実装済み）
- [ ] `layoutStore.ts` - レイアウト方向（LR/TB）の状態管理（将来対応可）

#### Step 9: UnifiedFlowCanvas 実装

- [x] Svelte Flow のセットアップ
- [x] カスタムノード/エッジタイプ登録
- [x] taskStore/wbsStore との連携
- [x] viewMode 切替対応（WBS モード時にグラフをフェードアウト）

#### Step 10: App.svelte 統合

- [x] GridCanvas → UnifiedFlowCanvas 切替（`App.svelte` で確認済み）
- [x] WBSGraphView → UnifiedFlowCanvas 統合（WBS パネルを Panel として統合）
- [x] Toolbar との連携確認

#### Step 11: パフォーマンステスト

- [ ] 500 ノードで動作確認
- [ ] 2000 ノードで動作確認
- [ ] パン/ズームの滑らかさ確認

#### Step 12: クリーンアップ

- [ ] `frontend/ide/src/lib/grid/` 削除
- [ ] `frontend/ide/src/lib/wbs/WBSGraphView.svelte` 削除
- [ ] `frontend/ide/src/lib/wbs/WBSGraphNode.svelte` 削除
- [ ] `frontend/ide/src/stores/viewportStore.ts` 削除（flowStore に統合）

### 技術メモ

#### Svelte 5 Runes 早見表

| Rune              | 用途             | Svelte 4 相当         |
| ----------------- | ---------------- | --------------------- |
| `$state(value)`   | リアクティブ状態 | `let value`           |
| `$derived(expr)`  | 派生値           | `$: derived = expr`   |
| `$derived.by(fn)` | 複雑な派生       | `$: { ... }`          |
| `$effect(fn)`     | 副作用           | `$: { sideEffect() }` |
| `$props()`        | プロップ受取     | `export let`          |
| `$bindable()`     | bind 可能        | `export let`          |

#### Svelte Flow 基本構成

```svelte
<script>
  import { SvelteFlow, Background, Controls } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';

  import TaskFlowNode from './nodes/TaskFlowNode.svelte';

  const nodeTypes = { task: TaskFlowNode };

  let nodes = $state([...]);
  let edges = $state([...]);
</script>

<SvelteFlow
  {nodes}
  {edges}
  {nodeTypes}
  fitView
  onlyRenderVisibleElements={true}
>
  <Background />
  <Controls />
</SvelteFlow>
```

#### 仮想化（Viewport Culling）

```svelte
<SvelteFlow
  onlyRenderVisibleElements={true}  <!-- 画面外ノードは非描画 -->
  minZoom={0.1}
  maxZoom={3}
/>
```

### 参考リンク

- [Svelte 5 Migration Guide](https://svelte.dev/docs/svelte/v5-migration-guide)
- [sv migrate CLI](https://svelte.dev/docs/cli/sv-migrate)
- [Svelte Flow Docs](https://svelteflow.dev/)
- [Svelte Flow Dagre Example](https://svelteflow.dev/examples/layout/dagre)
- [Svelte 5 Runes](https://svelte.dev/docs/svelte/runes)
