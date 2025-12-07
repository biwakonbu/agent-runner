<script lang="ts">
  /**
   * MainViewPreview - メインビュー全体のプレビューコンポーネント
   *
   * App.svelte のワークスペース選択後の状態を再現
   * Store/Wails依存を排除し、propsのみで動作
   */
  import { createEventDispatcher } from "svelte";
  import ToolbarPreview from "./toolbar/ToolbarPreview.svelte";
  import { WBSListView, WBSGraphView } from "./wbs";
  import GridCanvas from "./grid/GridCanvas.svelte";
  import FloatingChatWindow from "./components/chat/FloatingChatWindow.svelte";
  import { tasks, selectedTaskId } from "../stores/taskStore";
  import type { Task, TaskStatus, PoolSummary } from "../types";
  import { MessageSquare } from "lucide-svelte";

  const dispatch = createEventDispatcher();

  // === Props ===

  // ビュー設定
  export let viewMode: "graph" | "wbs" = "wbs";

  // タスクデータ
  export let taskList: Task[] = [];
  export let poolSummaries: PoolSummary[] = [];

  // 進捗
  export let overallProgress = { percentage: 0, completed: 0, total: 0 };

  // ステータス別カウント
  export let taskCountsByStatus: Record<TaskStatus, number> = {
    PENDING: 0,
    READY: 0,
    RUNNING: 0,
    SUCCEEDED: 0,
    COMPLETED: 0,
    FAILED: 0,
    CANCELED: 0,
    BLOCKED: 0,
    RETRY_WAIT: 0,
  };

  // 選択中タスク（ストア同期用。UI描画では未使用）
  export let selectedTask: Task | null = null;

  // モーダル・チャット
  export let showChat = true;
  export let chatPosition = { x: 600, y: 300 };

  // タスクストアを更新
  $: {
    tasks.setTasks(taskList);
    if (selectedTask) {
      selectedTaskId.select(selectedTask.id);
    } else {
      selectedTaskId.clear();
    }
  }

  function handleCloseChat() {
    dispatch("closeChat");
  }

  function handleOpenChat() {
    dispatch("openChat");
  }
</script>

<main class="app">
  <!-- ツールバー -->
  <div class="toolbar-overlay">
    <ToolbarPreview
      {viewMode}
      {overallProgress}
      {poolSummaries}
      {taskCountsByStatus}
    />
  </div>

  <!-- メインコンテンツ -->
  <div class="main-content">
    {#if viewMode === "graph"}
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
  {#if showChat}
    <FloatingChatWindow
      initialPosition={chatPosition}
      on:close={handleCloseChat}
    />
  {/if}

  <!-- チャット再表示ボタン -->
  {#if !showChat}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="chat-fab"
      on:click={handleOpenChat}
      on:keydown={(e) => e.key === "Enter" && handleOpenChat()}
      role="button"
      tabindex="0"
      aria-label="Open Chat"
    >
      <MessageSquare size="24" />
    </div>
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

  /* Toolbar is overlaid logic not strictly in App.svelte?
     Wait, App.svelte puts Toolbar *above* main-content in flex column.
     So keep Toolbar where it is.
     Correcting template structure to match App.svelte (Toolbar NOT overlay).
  */

  .main-content {
    display: block; /* Flex -> Block */
    position: relative;
    flex: 1;
    overflow: hidden;
    background: var(--mv-color-surface-base);
  }

  .canvas-layer {
    position: absolute;
    inset: 0;
    z-index: 1;
  }

  .list-overlay {
    position: absolute;
    inset: var(--mv-spacing-md);
    z-index: 10;
    background: var(--mv-color-surface-primary);
    border-radius: var(--mv-radius-lg);
    box-shadow: var(--mv-shadow-modal);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  /* モーダルは削除済み */
</style>
