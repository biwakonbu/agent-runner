<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button, Badge } from '../../design-system';
  // @ts-ignore - Wails自動生成ファイル
  import { RunTask } from '../../../wailsjs/go/main/App';
  import { selectedTask, selectedTaskId } from '../../stores';
  import { statusLabels } from '../../types';

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  let isRunning = false;

  async function handleRun() {
    if (!$selectedTask || isRunning) return;

    isRunning = true;
    try {
      await RunTask($selectedTask.id);
    } catch (e) {
      console.error('タスク実行エラー:', e);
    } finally {
      isRunning = false;
    }
  }

  function handleClose() {
    selectedTaskId.clear();
    dispatch('close');
  }

  function formatDate(dateString: string | undefined): string {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('ja-JP');
  }

  $: task = $selectedTask;
  $: badgeStatus = task ? task.status.toLowerCase() as 'pending' | 'ready' | 'running' | 'succeeded' | 'failed' | 'canceled' | 'blocked' : 'pending';
  $: canRun = task && task.status !== 'RUNNING';
</script>

<aside class="detail-panel" class:open={!!task}>
  {#if task}
    <!-- ヘッダー -->
    <header class="panel-header">
      <h2 class="panel-title">タスク詳細</h2>
      <Button variant="ghost" size="small" on:click={handleClose} label="×" />
    </header>

    <!-- コンテンツ -->
    <div class="panel-content">
      <!-- タスク名とステータス -->
      <div class="task-header">
        <h3 class="task-title">{task.title}</h3>
        <Badge status={badgeStatus} label={statusLabels[task.status]} />
      </div>

      <!-- アクション -->
      <div class="actions">
        <Button
          variant="primary"
          on:click={handleRun}
          disabled={!canRun}
          loading={isRunning}
          loadingLabel="実行中..."
        >
          {#if task.status === 'RUNNING'}
            実行中
          {:else}
            タスクを実行
          {/if}
        </Button>
      </div>

      <!-- メタ情報 -->
      <div class="meta-section">
        <h4 class="section-title">情報</h4>

        <dl class="meta-list">
          <div class="meta-item">
            <dt>ID</dt>
            <dd class="mono">{task.id}</dd>
          </div>

          <div class="meta-item">
            <dt>Pool</dt>
            <dd class="mono">{task.poolId}</dd>
          </div>

          <div class="meta-item">
            <dt>作成日時</dt>
            <dd>{formatDate(task.createdAt)}</dd>
          </div>

          {#if task.startedAt}
            <div class="meta-item">
              <dt>開始日時</dt>
              <dd>{formatDate(task.startedAt)}</dd>
            </div>
          {/if}

          {#if task.doneAt}
            <div class="meta-item">
              <dt>完了日時</dt>
              <dd>{formatDate(task.doneAt)}</dd>
            </div>
          {/if}
        </dl>
      </div>
    </div>
  {:else}
    <div class="empty-state">
      <p>タスクを選択してください</p>
    </div>
  {/if}
</aside>

<style>
  .detail-panel {
    width: var(--mv-layout-detail-panel-width);
    background: var(--mv-color-surface-primary);
    border-left: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    transform: translateX(100%);
    transition: transform var(--mv-transition-state);
  }

  .detail-panel.open {
    transform: translateX(0);
  }

  /* パネルヘッダー */
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--mv-spacing-md);
    border-bottom: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
  }

  .panel-title {
    font-size: var(--mv-font-size-md);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    margin: 0;
  }

  /* パネルコンテンツ */
  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--mv-spacing-md);
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-lg);
  }

  /* タスクヘッダー */
  .task-header {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xs);
  }

  .task-title {
    font-size: var(--mv-font-size-lg);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    margin: 0;
    word-break: break-word;
  }

  /* アクション */
  .actions {
    display: flex;
    gap: var(--mv-spacing-xs);
  }

  /* メタ情報セクション */
  .meta-section {
    border-top: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    padding-top: var(--mv-spacing-md);
  }

  .section-title {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-muted);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-wider);
    margin: 0 0 var(--mv-spacing-sm);
  }

  .meta-list {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-sm);
    margin: 0;
  }

  .meta-item {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xxs);
  }

  .meta-item dt {
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
  }

  .meta-item dd {
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-primary);
    margin: 0;
    word-break: break-all;
  }

  .meta-item dd.mono {
    font-family: var(--mv-font-mono);
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-secondary);
  }

  /* 空状態 */
  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--mv-color-text-muted);
    font-size: var(--mv-font-size-sm);
  }
</style>
