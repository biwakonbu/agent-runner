<script lang="ts">
  import { overallProgress, expandedNodes } from "../../stores/wbsStore";
  import { getProgressColor } from "./utils";
  import ProgressBar from "./ProgressBar.svelte";

  function splitProgress(percentage: number) {
    const str = Math.round(percentage).toString();
    if (str.length === 1) return { first: str, rest: "" };
    return { first: str[0], rest: str.slice(1) };
  }

  $: percentage = $overallProgress?.percentage ?? 0;
  $: completed = $overallProgress?.completed ?? 0;
  $: total = $overallProgress?.total ?? 0;
  $: progressParts = splitProgress(percentage);
  $: progressColor = getProgressColor(percentage);
</script>

<header class="wbs-header">
  <div class="header-title">
    <h2>Work Breakdown Structure</h2>
    <span class="task-count">
      {completed} / {total} Tasks Completed
    </span>
  </div>

  <div class="header-progress">
    <ProgressBar {percentage} size="md" />
    <span
      class="progress-percentage"
      style:color={progressColor.fill}
      style:text-shadow={progressColor.textShadowMd}
    >
      <span class="progress-first-digit">{progressParts.first}</span>
      <span class="progress-rest-digits">{progressParts.rest}</span>
      <span class="progress-symbol">%</span>
    </span>
  </div>

  <div class="header-actions">
    <button
      class="action-btn"
      on:click={() => expandedNodes.expandAll()}
      title="Expand All"
    >
      ↕ Expand All
    </button>
    <button
      class="action-btn"
      on:click={() => expandedNodes.collapseAll()}
      title="Collapse All"
    >
      ⇕ Collapse All
    </button>
  </div>
</header>

<style>
  .wbs-header {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-sm);
    padding: var(--mv-spacing-md);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-color-border-subtle);
    background: var(--mv-color-surface-hover);
    flex-shrink: 0;
  }

  .header-title {
    display: flex;
    align-items: baseline;
    gap: var(--mv-spacing-sm);
  }

  .header-title h2 {
    font-size: var(--mv-font-size-lg);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    margin: 0;
  }

  .task-count {
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-muted);
  }

  .header-progress {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
  }

  .progress-percentage {
    display: flex;
    align-items: baseline;
    font-family: var(--mv-font-display); /* Back to Orbitron */
    color: var(--mv-progress-text-color);
    min-width: var(--mv-progress-text-width-md);
    justify-content: flex-end;
    text-shadow: var(--mv-text-shadow-glow);
    line-height: 1;
    font-style: italic;
  }

  .progress-first-digit {
    font-size: calc(var(--mv-font-size-xl) * 1.15); /* Exactly 1.15x */
    font-weight: var(--mv-font-weight-bold);
  }

  .progress-rest-digits {
    font-size: var(--mv-font-size-xl);
    font-weight: var(--mv-font-weight-bold); /* Bold */
  }

  .progress-symbol {
    font-size: var(--mv-font-size-xl);
    font-weight: var(--mv-font-weight-bold);
    margin-left: var(--mv-spacing-xxxs);
    opacity: 0.8;
    font-style: italic;
  }

  .header-actions {
    display: flex;
    gap: var(--mv-spacing-xs);
  }

  .action-btn {
    padding: var(--mv-spacing-xxs) var(--mv-spacing-sm);
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-medium);
    color: var(--mv-color-text-secondary);
    background: var(--mv-color-surface-primary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-sm);
    cursor: pointer;
    transition:
      background-color var(--mv-transition-hover),
      color var(--mv-transition-hover);
  }

  .action-btn:hover {
    background: var(--mv-color-surface-hover);
    color: var(--mv-color-text-primary);
  }
</style>
