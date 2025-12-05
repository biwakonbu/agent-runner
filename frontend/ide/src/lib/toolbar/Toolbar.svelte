<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '../../design-system';
  import { viewport, zoomPercent, taskCountsByStatus } from '../../stores';
  import type { TaskStatus } from '../../types';

  const dispatch = createEventDispatcher<{
    createTask: void;
  }>();

  // ステータスサマリの表示設定
  const statusDisplay: { key: TaskStatus; label: string; showCount: boolean }[] = [
    { key: 'RUNNING', label: '実行中', showCount: true },
    { key: 'PENDING', label: '待機', showCount: true },
    { key: 'FAILED', label: '失敗', showCount: true },
  ];

  function handleCreateTask() {
    dispatch('createTask');
  }
</script>

<header class="toolbar">
  <!-- 左側：タイトルと操作 -->
  <div class="toolbar-left">
    <h1 class="app-title">multiverse IDE</h1>

    <Button variant="primary" size="small" on:click={handleCreateTask}>
      <span class="btn-icon">+</span>
      新規タスク
    </Button>
  </div>

  <!-- 中央：ステータスサマリ -->
  <div class="toolbar-center">
    <div class="status-summary">
      {#each statusDisplay as { key, label, showCount }}
        {#if showCount && $taskCountsByStatus[key] > 0}
          <div class="status-badge status-{key.toLowerCase()}">
            <span class="status-count">{$taskCountsByStatus[key]}</span>
            <span class="status-label">{label}</span>
          </div>
        {/if}
      {/each}
    </div>
  </div>

  <!-- 右側：ズームコントロール -->
  <div class="toolbar-right">
    <div class="zoom-controls">
      <Button
        variant="ghost"
        size="small"
        on:click={() => viewport.zoomOut()}
        label="−"
      />

      <button
        class="zoom-value"
        on:click={() => viewport.reset()}
        aria-label="ズームリセット"
        title="リセット (0)"
      >
        {$zoomPercent}%
      </button>

      <Button
        variant="ghost"
        size="small"
        on:click={() => viewport.zoomIn()}
        label="+"
      />
    </div>
  </div>
</header>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: var(--mv-layout-toolbar-height);
    padding: 0 var(--mv-spacing-md);
    background: var(--mv-color-surface-primary);
    border-bottom: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    flex-shrink: 0;
  }

  .toolbar-left,
  .toolbar-center,
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-md);
  }

  .toolbar-left {
    flex: 1;
  }

  .toolbar-center {
    flex: 2;
    justify-content: center;
  }

  .toolbar-right {
    flex: 1;
    justify-content: flex-end;
  }

  .app-title {
    font-size: var(--mv-font-size-lg);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    margin: 0;
  }

  .btn-icon {
    font-size: var(--mv-font-size-lg);
    line-height: 1;
  }

  /* ステータスサマリ */
  .status-summary {
    display: flex;
    gap: var(--mv-spacing-sm);
  }

  .status-badge {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xxs);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
    font-size: var(--mv-font-size-xs);
  }

  .status-badge.status-running {
    background: var(--mv-color-status-running-bg);
    color: var(--mv-color-status-running-text);
  }

  .status-badge.status-pending {
    background: var(--mv-color-status-pending-bg);
    color: var(--mv-color-status-pending-text);
  }

  .status-badge.status-failed {
    background: var(--mv-color-status-failed-bg);
    color: var(--mv-color-status-failed-text);
  }

  .status-count {
    font-weight: var(--mv-font-weight-bold);
  }

  .status-label {
    font-weight: var(--mv-font-weight-normal);
  }

  /* ズームコントロール */
  .zoom-controls {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xxs);
    background: var(--mv-color-surface-secondary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-sm);
    padding: var(--mv-spacing-xxs);
  }

  .zoom-value {
    min-width: var(--mv-spacing-xxl);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    font-size: var(--mv-font-size-xs);
    font-family: var(--mv-font-mono);
    color: var(--mv-color-text-secondary);
    background: transparent;
    border: none;
    cursor: pointer;
    text-align: center;
  }

  .zoom-value:hover {
    color: var(--mv-color-text-primary);
  }
</style>
