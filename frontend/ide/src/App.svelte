<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { get } from "svelte/store";
  import WorkspaceSelector from "./lib/WorkspaceSelector.svelte";
  import TitleBar from "./lib/TitleBar.svelte";
  import { Toolbar } from "./lib/toolbar";
  import { WBSListView, WBSGraphView } from "./lib/wbs";
  import GridCanvas from "./lib/grid/GridCanvas.svelte";
  import {
    tasks,
    selectedTask,
    selectedTaskId,
    poolSummaries,
    viewMode,
  } from "./stores";
  import { Logger } from "./services/logger";
  import type { Task, PoolSummary } from "./types";
  // @ts-ignore - Wails自動生成ファイル
  import { ListTasks, GetPoolSummaries } from "../wailsjs/go/main/App";
  import FloatingChatWindow from "./lib/components/chat/FloatingChatWindow.svelte";
  import {
    initExecutionEvents,
    syncExecutionState,
  } from "./stores/executionStore";
  import { initTaskEvents } from "./stores/taskStore";
  import { initChatEvents } from "./stores/chat";
  import { initBacklogEvents, unresolvedCount } from "./stores/backlogStore";
  import BacklogPanel from "./lib/backlog/BacklogPanel.svelte";

  const log = Logger.withComponent("App");

  let workspaceId: string | null = null;
  let interval: ReturnType<typeof setInterval> | null = null;

  // Chat State
  let isChatVisible = true;
  let chatPosition = { x: 0, y: 0 };

  // Backlog State
  let isBacklogVisible = false;

  // 共通のタスクマッピング（Wailsの生データ→UI用Task型）
  const mapBackendTask = (t: any): Task => ({
    id: t.id,
    title: t.title,
    status: t.status as Task["status"],
    poolId: t.poolId,
    createdAt: t.createdAt,
    updatedAt: t.updatedAt,
    startedAt: t.startedAt,
    doneAt: t.doneAt,
    description: t.description,
    dependencies: t.dependencies ?? [],
    parentId: t.parentId,
    wbsLevel: t.wbsLevel,
    phaseName: t.phaseName as Task["phaseName"],
    milestone: t.milestone,
    sourceChatId: t.sourceChatId,
    acceptanceCriteria: t.acceptanceCriteria ?? [],
    attemptCount: t.attemptCount,
    nextRetryAt: t.nextRetryAt,
  });

  onMount(() => {
    // Calculate initial position (Bottom-Right)
    // 600px width, 350px height, 20px padding
    const width = 600;
    const height = 350;
    const padding = 20;
    chatPosition = {
      x: window.innerWidth - width - padding,
      y: window.innerHeight - height - padding,
    };
    // Wails Events 初期化
    initExecutionEvents();
    initTaskEvents();
    initChatEvents();
    initBacklogEvents();
  });

  // タスク一覧を読み込み
  async function loadTasks() {
    if (!workspaceId) return;
    try {
      const result = await ListTasks();
      // Wails生成型からローカル型へ変換
      const taskList: Task[] = (result || []).map(mapBackendTask);
      log.debug("tasks loaded", { count: taskList.length });
      tasks.setTasks(taskList);
    } catch (e) {
      log.error("failed to load tasks", { error: e });
    }
  }

  // Pool別サマリを読み込み
  async function loadPoolSummaries() {
    if (!workspaceId) return;
    try {
      const summaries: PoolSummary[] = await GetPoolSummaries();
      log.debug("pool summaries loaded", { count: summaries?.length ?? 0 });
      poolSummaries.setSummaries(summaries || []);
    } catch (e) {
      log.error("failed to load pool summaries", { error: e });
    }
  }

  // データ読み込み（タスク + Poolサマリ）
  async function loadData() {
    await Promise.all([loadTasks(), loadPoolSummaries()]);
  }

  // Workspace選択時
  function onWorkspaceSelected(event: CustomEvent<string>) {
    // 前のワークスペースのデータをクリア
    tasks.clear();
    selectedTaskId.clear();
    poolSummaries.clear();

    workspaceId = event.detail;
    log.info("workspace selected", { workspaceId });

    loadData();
    // 実行状態をバックエンドと同期
    syncExecutionState();
    // 10秒間隔でポーリング（Wails Events でリアルタイム更新されるためフォールバック）
    interval = setInterval(loadData, 10000);
    log.info("polling started", { interval: 10000 });
  }

  onDestroy(() => {
    if (interval) {
      log.info("polling stopped");
      clearInterval(interval);
    }
  });

  // チャットから生成されたタスクを即時反映
  function onTasksGenerated(event: CustomEvent<{ tasks: any[] }>) {
    const generated = event.detail?.tasks ?? [];
    if (generated.length === 0) return;

    const existing = new Map(get(tasks).map((t) => [t.id, t]));
    for (const raw of generated) {
      const mapped = mapBackendTask(raw);
      if (existing.has(mapped.id)) {
        tasks.updateTask(mapped.id, mapped);
      } else {
        tasks.addTask(mapped);
      }
    }
  }
</script>

<main class="app">
  <TitleBar />
  {#if !workspaceId}
    <WorkspaceSelector on:selected={onWorkspaceSelected} />
  {:else}
    <!-- ツールバー -->
    <Toolbar />

    <!-- メインコンテンツ -->
    <div class="main-content">
      {#if $viewMode === "graph"}
        <!-- Graph モード: GridCanvas で依存グラフ表示 -->
        <div class="canvas-layer">
          <GridCanvas />
        </div>
      {:else}
        <!-- WBS モード: WBSGraphView + WBSListView -->
        <div class="canvas-layer">
          <WBSGraphView />
        </div>
        <div class="list-overlay">
          <WBSListView />
        </div>
      {/if}
    </div>

    <!-- チャットウィンドウ -->
    {#if isChatVisible}
      <FloatingChatWindow
        initialPosition={chatPosition}
        on:close={() => (isChatVisible = false)}
        on:tasksGenerated={onTasksGenerated}
      />
    {/if}

    <!-- チャット再表示ボタン (簡易FAB) -->
    {#if !isChatVisible}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div
        class="chat-fab"
        on:click={() => (isChatVisible = true)}
        on:keydown={(e) => e.key === "Enter" && (isChatVisible = true)}
        role="button"
        tabindex="0"
        aria-label="Open Chat"
      >
        💬
      </div>
    {/if}

    <!-- バックログ表示ボタン -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="backlog-fab"
      class:has-items={$unresolvedCount > 0}
      on:click={() => (isBacklogVisible = !isBacklogVisible)}
      on:keydown={(e) =>
        e.key === "Enter" && (isBacklogVisible = !isBacklogVisible)}
      role="button"
      tabindex="0"
      aria-label="Toggle Backlog"
    >
      {#if $unresolvedCount > 0}
        <span class="backlog-count">{$unresolvedCount}</span>
      {:else}
        &#9776;
      {/if}
    </div>

    <!-- バックログパネル -->
    {#if isBacklogVisible}
      <div class="backlog-sidebar">
        <BacklogPanel />
      </div>
    {/if}
  {/if}
</main>

<style>
  .chat-fab {
    position: fixed;
    bottom: var(--mv-spacing-lg);
    right: var(--mv-spacing-lg);
    width: var(--mv-icon-size-xxxl);
    height: var(--mv-icon-size-xxxl);
    background: var(--mv-color-surface-primary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--mv-shadow-card);
    cursor: pointer;
    z-index: 1000;
    font-size: var(--mv-icon-size-md);
  }
  .chat-fab:hover {
    background: var(--mv-color-surface-hover);
  }

  .app {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--mv-color-surface-app);
    color: var(--mv-color-text-primary);
    font-family: var(--mv-font-sans);
    overflow: hidden;
  }

  .main-content {
    display: block; /* フレックスからブロックへ変更 (絶対配置のコンテナにするため) */
    position: relative;
    flex: 1;
    overflow: hidden;
    background: var(--mv-color-surface-base); /* Canvasの背景色 */
  }

  .canvas-layer {
    position: absolute;
    inset: 0;
    z-index: 1;
  }

  .list-overlay {
    position: absolute;
    inset: var(--mv-spacing-md); /* 少し余白を持たせてフローティング感を出す */
    z-index: 10;
    background: var(--mv-color-surface-primary);
    border-radius: var(--mv-radius-lg);
    box-shadow: var(--mv-shadow-modal);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  /* タスク作成モーダルは削除済み */

  .backlog-fab {
    position: fixed;
    bottom: var(--mv-spacing-lg);

    /* ズームコントロールの右側に配置（約200px右にオフセット） */
    left: var(--mv-backlog-fab-left);
    width: var(--mv-icon-size-xxxl);
    height: var(--mv-icon-size-xxxl);
    background: var(--mv-color-surface-primary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-full);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--mv-shadow-card);
    cursor: pointer;
    z-index: 1000;
    font-size: var(--mv-icon-size-md);
    transition: all var(--mv-transition-hover);
  }

  .backlog-fab:hover {
    background: var(--mv-color-surface-hover);
  }

  .backlog-fab.has-items {
    background: var(--mv-color-status-failed-bg);
    border-color: var(--mv-color-status-failed-text);
  }

  .backlog-count {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-bold);
    color: var(--mv-color-status-failed-text);
  }

  .backlog-sidebar {
    position: fixed;
    top: var(--mv-backlog-sidebar-top);
    left: 0;
    bottom: 0;
    width: var(--mv-backlog-sidebar-width);
    z-index: 100;
    box-shadow: var(--mv-shadow-modal);
  }
</style>
