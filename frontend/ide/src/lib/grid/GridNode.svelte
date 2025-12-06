<script lang="ts">
  import type { Task, PhaseName } from '../../types';
  import { gridToCanvas } from '../../design-system';
  import { statusToCssClass, statusLabels } from '../../types';
  import { selectedTaskId } from '../../stores';

  // Props
  export let task: Task;
  export let col: number;
  export let row: number;
  export let zoomLevel: number = 1;

  // フェーズ名からCSSクラス名への変換
  function phaseToClass(phase: PhaseName | undefined): string {
    if (!phase) return '';
    const phaseMap: Record<string, string> = {
      '概念設計': 'phase-concept',
      '実装設計': 'phase-design',
      '実装': 'phase-impl',
      '検証': 'phase-verify',
    };
    return phaseMap[phase] || '';
  }

  // フェーズの表示ラベル
  const phaseLabels: Record<string, string> = {
    '概念設計': 'CONCEPT',
    '実装設計': 'DESIGN',
    '実装': 'IMPL',
    '検証': 'VERIFY',
  };

  // キャンバス座標を計算
  $: position = gridToCanvas(col, row);
  $: isSelected = $selectedTaskId === task.id;
  $: statusClass = statusToCssClass(task.status);
  $: phaseClass = phaseToClass(task.phaseName);
  $: hasDependencies = task.dependencies && task.dependencies.length > 0;

  // ズームレベルに応じた表示制御
  $: showTitle = zoomLevel >= 0.4;
  $: showDetails = zoomLevel >= 1.2;

  function handleClick() {
    selectedTaskId.select(task.id);
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleClick();
    }
  }
</script>

<div
  class="node status-{statusClass} {phaseClass}"
  class:selected={isSelected}
  class:has-deps={hasDependencies}
  style="left: {position.x}px; top: {position.y}px;"
  on:click={handleClick}
  on:keydown={handleKeydown}
  role="button"
  tabindex="0"
  aria-label="{task.title} - {statusLabels[task.status]}"
>
  <!-- フェーズインジケーター（左端のバー） -->
  {#if task.phaseName && phaseClass}
    <div class="phase-bar" aria-hidden="true"></div>
  {/if}

  <!-- ステータスインジケーター -->
  <div class="status-indicator">
    <span class="status-dot"></span>
    <span class="status-text">{statusLabels[task.status]}</span>
    {#if task.phaseName}
      <span class="phase-badge">{phaseLabels[task.phaseName] || ''}</span>
    {/if}
  </div>

  <!-- タイトル（ズームレベルに応じて表示） -->
  {#if showTitle}
    <div class="title" title={task.title}>
      {task.title}
    </div>
  {/if}

  <!-- 詳細情報（高ズームレベルで表示） -->
  {#if showDetails}
    <div class="details">
      <span class="pool">{task.poolId}</span>
      {#if hasDependencies}
        <span class="deps-count" title="依存タスク数">
          {task.dependencies?.length || 0}
        </span>
      {/if}
    </div>
  {/if}
</div>

<style>
  .node {
    position: absolute;
    width: var(--mv-grid-cell-width);
    height: var(--mv-grid-cell-height);
    background: var(--mv-color-surface-node);
    border: var(--mv-border-width-default) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-md);
    padding: var(--mv-spacing-xs) var(--mv-spacing-sm);
    cursor: pointer;
    transition: border-color var(--mv-transition-hover),
                box-shadow var(--mv-transition-hover),
                transform var(--mv-transition-hover);
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xxs);
    overflow: hidden;
    user-select: none;
  }

  .node:hover {
    border-color: var(--mv-color-border-strong);
    transform: var(--mv-transform-hover-lift);
  }

  .node:focus {
    outline: none;
    border-color: var(--mv-color-border-focus);
    box-shadow: var(--mv-shadow-focus);
  }

  .node.selected {
    border-color: var(--mv-color-border-focus);
    box-shadow: var(--mv-shadow-selected);
  }

  /* ステータス別スタイル */
  .node.status-pending {
    background: var(--mv-color-status-pending-bg);
    border-color: var(--mv-color-status-pending-border);
  }

  .node.status-ready {
    background: var(--mv-color-status-ready-bg);
    border-color: var(--mv-color-status-ready-border);
  }

  .node.status-running {
    background: var(--mv-color-status-running-bg);
    border-color: var(--mv-color-status-running-border);
    animation: mv-pulse var(--mv-duration-pulse) infinite;
  }

  .node.status-succeeded {
    background: var(--mv-color-status-succeeded-bg);
    border-color: var(--mv-color-status-succeeded-border);
  }

  .node.status-failed {
    background: var(--mv-color-status-failed-bg);
    border-color: var(--mv-color-status-failed-border);
  }

  .node.status-canceled {
    background: var(--mv-color-status-canceled-bg);
    border-color: var(--mv-color-status-canceled-border);
  }

  .node.status-blocked {
    background: var(--mv-color-status-blocked-bg);
    border-color: var(--mv-color-status-blocked-border);
  }

  /* ステータスインジケーター */
  .status-indicator {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xxs);
  }

  .status-dot {
    width: var(--mv-indicator-size-sm);
    height: var(--mv-indicator-size-sm);
    border-radius: var(--mv-radius-full);
    flex-shrink: 0;
  }

  .status-pending .status-dot {
    background: var(--mv-color-status-pending-text);
  }

  .status-ready .status-dot {
    background: var(--mv-color-status-ready-text);
  }

  .status-running .status-dot {
    background: var(--mv-color-status-running-text);
  }

  .status-succeeded .status-dot {
    background: var(--mv-color-status-succeeded-text);
  }

  .status-failed .status-dot {
    background: var(--mv-color-status-failed-text);
  }

  .status-canceled .status-dot {
    background: var(--mv-color-status-canceled-text);
  }

  .status-blocked .status-dot {
    background: var(--mv-color-status-blocked-text);
  }

  .status-text {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-bold);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-wide);
  }

  .status-pending .status-text {
    color: var(--mv-color-status-pending-text);
  }

  .status-ready .status-text {
    color: var(--mv-color-status-ready-text);
  }

  .status-running .status-text {
    color: var(--mv-color-status-running-text);
  }

  .status-succeeded .status-text {
    color: var(--mv-color-status-succeeded-text);
  }

  .status-failed .status-text {
    color: var(--mv-color-status-failed-text);
  }

  .status-canceled .status-text {
    color: var(--mv-color-status-canceled-text);
  }

  .status-blocked .status-text {
    color: var(--mv-color-status-blocked-text);
  }

  /* タイトル */
  .title {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    line-height: var(--mv-line-height-normal);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  /* 詳細情報 */
  .details {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xs);
    margin-top: auto;
  }

  .pool {
    font-size: var(--mv-font-size-xs);
    font-family: var(--mv-font-mono);
    color: var(--mv-color-text-secondary);
    background: var(--mv-color-surface-secondary);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-sm);
  }

  /* 依存タスク数バッジ */
  .deps-count {
    font-size: var(--mv-font-size-xs);
    font-family: var(--mv-font-mono);
    color: var(--mv-color-text-muted);
    background: var(--mv-color-surface-overlay);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-sm);
  }

  .deps-count::before {
    content: '\2192 '; /* 矢印 → */
  }

  /* フェーズバー（左端のカラーインジケーター） */
  .phase-bar {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-md) 0 0 var(--mv-radius-md);
  }

  /* フェーズバッジ */
  .phase-badge {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-bold);
    padding: 0 var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-sm);
    margin-left: auto;
    letter-spacing: var(--mv-letter-spacing-wide);
  }

  /* フェーズ別カラー - 概念設計（青系） */
  .phase-concept .phase-bar {
    background: var(--mv-primitive-frost-3);
  }

  .phase-concept .phase-badge {
    color: var(--mv-primitive-frost-3);
    background: var(--mv-color-surface-secondary);
  }

  /* フェーズ別カラー - 実装設計（紫系） */
  .phase-design .phase-bar {
    background: var(--mv-primitive-aurora-purple);
  }

  .phase-design .phase-badge {
    color: var(--mv-primitive-aurora-purple);
    background: var(--mv-color-surface-secondary);
  }

  /* フェーズ別カラー - 実装（緑系） */
  .phase-impl .phase-bar {
    background: var(--mv-primitive-aurora-green);
  }

  .phase-impl .phase-badge {
    color: var(--mv-primitive-aurora-green);
    background: var(--mv-color-surface-secondary);
  }

  /* フェーズ別カラー - 検証（オレンジ系） */
  .phase-verify .phase-bar {
    background: var(--mv-primitive-aurora-yellow);
  }

  .phase-verify .phase-badge {
    color: var(--mv-primitive-aurora-yellow);
    background: var(--mv-color-surface-secondary);
  }

  /* 依存がある場合のスタイル */
  .has-deps {
    padding-left: var(--mv-spacing-md);
  }
</style>
