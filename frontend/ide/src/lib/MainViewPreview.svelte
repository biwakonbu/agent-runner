<script lang="ts">
  /**
   * MainViewPreview - メインビュー全体のプレビューコンポーネント
   *
   * App.svelte のワークスペース選択後の状態を再現
   * Store/Wails依存を排除し、propsのみで動作
   */
  import { createEventDispatcher } from "svelte";
  import ToolbarPreview from "./toolbar/ToolbarPreview.svelte";
  import DetailPanelPreview from "./panel/DetailPanelPreview.svelte";
  import { WBSListView, WBSGraphView } from "./wbs";
  import { Button } from "../design-system";
  import TaskCreatePreview from "./TaskCreatePreview.svelte";
  import FloatingChatWindow from "./components/chat/FloatingChatWindow.svelte";
  import { tasks, selectedTaskId } from "../stores/taskStore";
  import type { Task, TaskStatus, PoolSummary, Attempt } from "../types";

  const dispatch = createEventDispatcher();

  // === Props ===

  // ビュー設定
  export let viewMode: "graph" | "wbs" = "wbs";
  export let zoomPercent = 100;

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
  };

  // 詳細パネル
  export let selectedTask: Task | null = null;
  export let attempts: Attempt[] = [];
  export let isTaskRunning = false;

  // モーダル・チャット
  export let showCreateModal = false;
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

  function handleCreateTask() {
    dispatch("createTask");
  }

  function handleCloseModal() {
    dispatch("closeModal");
  }

  function handleClosePanel() {
    dispatch("closePanel");
  }

  function handleRunTask() {
    dispatch("runTask");
  }

  function handleCloseChat() {
    dispatch("closeChat");
  }

  function handleOpenChat() {
    dispatch("openChat");
  }

  $: isGraphMode = viewMode === "graph";
</script>

<main class="app">
  <!-- ツールバー -->
  <ToolbarPreview
    {viewMode}
    {zoomPercent}
    {overallProgress}
    {poolSummaries}
    {taskCountsByStatus}
    on:createTask={handleCreateTask}
  />

  <!-- メインコンテンツ -->
  <div class="main-content">
    <!-- Graph/WBS ビュー切り替え -->
    {#if isGraphMode}
      <WBSGraphView />
    {:else}
      <WBSListView />
    {/if}

    <!-- 詳細パネル -->
    <DetailPanelPreview
      task={selectedTask}
      {attempts}
      isRunning={isTaskRunning}
      on:close={handleClosePanel}
      on:run={handleRunTask}
    />
  </div>

  <!-- タスク作成モーダル -->
  {#if showCreateModal}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-overlay" on:click={handleCloseModal} role="presentation">
      <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
      <div
        class="modal-content"
        on:click|stopPropagation
        role="dialog"
        aria-modal="true"
        aria-labelledby="create-task-title"
      >
        <header class="modal-header">
          <h2 id="create-task-title">新規タスク作成</h2>
          <Button
            variant="ghost"
            size="small"
            on:click={handleCloseModal}
            label="×"
          />
        </header>
        <TaskCreatePreview />
      </div>
    </div>
  {/if}

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
      💬
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

  .main-content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  /* モーダルオーバーレイ */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: var(--mv-color-surface-overlay);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    background: var(--mv-color-surface-primary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-lg);
    width: 100%;
    max-width: var(--mv-container-max-width-sm);
    max-height: var(--mv-container-max-height-modal);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--mv-spacing-md);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-color-border-subtle);
  }

  .modal-header h2 {
    font-size: var(--mv-font-size-lg);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    margin: 0;
  }
</style>
