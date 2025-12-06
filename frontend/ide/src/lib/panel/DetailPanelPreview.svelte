<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Button, Badge } from "../../design-system";
  import { statusLabels, attemptStatusLabels } from "../../types";
  import type { Task, Attempt, AttemptStatus } from "../../types";

  const dispatch = createEventDispatcher<{
    close: void;
    run: void;
  }>();

  // Props（Wails API依存を排除）
  export let task: Task | null = null;
  export let attempts: Attempt[] = [];
  export let isRunning = false;
  export let loadingAttempts = false;

  function handleRun() {
    if (!task || isRunning) return;
    dispatch("run");
  }

  function handleClose() {
    dispatch("close");
  }

  function formatDate(dateString: string | undefined): string {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString("ja-JP");
  }

  function getAttemptStatusClass(status: AttemptStatus): string {
    switch (status) {
      case "SUCCEEDED":
        return "succeeded";
      case "FAILED":
      case "TIMEOUT":
        return "failed";
      case "RUNNING":
      case "STARTING":
        return "running";
      case "CANCELED":
        return "canceled";
      default:
        return "pending";
    }
  }

  $: badgeStatus = task
    ? (task.status.toLowerCase() as
        | "pending"
        | "ready"
        | "running"
        | "succeeded"
        | "failed"
        | "canceled"
        | "blocked")
    : "pending";
  $: canRun = task && task.status !== "RUNNING";
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
          {#if task.status === "RUNNING"}
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

      <!-- 実行履歴セクション -->
      <div class="attempts-section">
        <h4 class="section-title">実行履歴</h4>

        {#if loadingAttempts}
          <p class="loading-text">読み込み中...</p>
        {:else if attempts.length === 0}
          <p class="empty-text">まだ実行履歴がありません</p>
        {:else}
          <ul class="attempts-list">
            {#each attempts as attempt (attempt.id)}
              <li class="attempt-item">
                <div class="attempt-header">
                  <span
                    class="attempt-status status-{getAttemptStatusClass(
                      attempt.status
                    )}"
                  >
                    {attemptStatusLabels[attempt.status]}
                  </span>
                  <span class="attempt-time"
                    >{formatDate(attempt.startedAt)}</span
                  >
                </div>
                {#if attempt.errorSummary}
                  <p class="attempt-error">{attempt.errorSummary}</p>
                {/if}
              </li>
            {/each}
          </ul>
        {/if}
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
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-color-border-subtle);
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

  /* 実行履歴セクション */
  .attempts-section {
    border-top: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    padding-top: var(--mv-spacing-md);
  }

  .loading-text,
  .empty-text {
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-muted);
    margin: 0;
  }

  .attempts-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xs);
  }

  .attempt-item {
    background: var(--mv-color-surface-secondary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    border-radius: var(--mv-radius-sm);
    padding: var(--mv-spacing-sm);
  }

  .attempt-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--mv-spacing-xs);
  }

  .attempt-status {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-semibold);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
  }

  .attempt-status.status-succeeded {
    background: var(--mv-color-status-succeeded-bg);
    color: var(--mv-color-status-succeeded-text);
  }

  .attempt-status.status-failed {
    background: var(--mv-color-status-failed-bg);
    color: var(--mv-color-status-failed-text);
  }

  .attempt-status.status-running {
    background: var(--mv-color-status-running-bg);
    color: var(--mv-color-status-running-text);
  }

  .attempt-status.status-canceled {
    background: var(--mv-color-status-canceled-bg);
    color: var(--mv-color-status-canceled-text);
  }

  .attempt-status.status-pending {
    background: var(--mv-color-status-pending-bg);
    color: var(--mv-color-status-pending-text);
  }

  .attempt-time {
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
    font-family: var(--mv-font-mono);
  }

  .attempt-error {
    margin: var(--mv-spacing-xs) 0 0;
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-status-failed-text);
    background: var(--mv-color-status-failed-bg);
    padding: var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
    word-break: break-word;
    white-space: pre-wrap;
    max-height: var(--mv-container-max-height-error);
    overflow-y: auto;
  }
</style>
