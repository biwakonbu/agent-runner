<script lang="ts">
  import {
    taskCountsByStatus,
    poolSummaries,
    viewMode,
    overallProgress,
  } from "../../stores";
  import type { TaskStatus } from "../../types";
  import BrandText from "../components/brand/BrandText.svelte";
  import ProgressBar from "../wbs/ProgressBar.svelte";
  import ExecutionControls from "./ExecutionControls.svelte";
  import Button from "../../design-system/components/Button.svelte";
  import Badge from "../../design-system/components/Badge.svelte";
  import { Network, ListTree } from "lucide-svelte";

  // Badge status type
  type BadgeStatus =
    | "pending"
    | "ready"
    | "running"
    | "succeeded"
    | "failed"
    | "canceled"
    | "blocked";

  // TaskStatus → BadgeStatus マッピング
  const statusToLower: Record<string, BadgeStatus> = {
    PENDING: "pending",
    READY: "ready",
    RUNNING: "running",
    SUCCEEDED: "succeeded",
    FAILED: "failed",
    CANCELED: "canceled",
    BLOCKED: "blocked",
  };

  // ステータスサマリの表示設定
  const statusDisplay: {
    key: TaskStatus;
    label: string;
    showCount: boolean;
    color:
      | "primary"
      | "secondary"
      | "success"
      | "warning"
      | "danger"
      | "info"
      | "neutral";
  }[] = [
    { key: "RUNNING", label: "実行中", showCount: true, color: "success" },
    { key: "PENDING", label: "待機", showCount: true, color: "warning" },
    { key: "FAILED", label: "失敗", showCount: true, color: "danger" },
  ];

  // Pool別サマリがある場合はそれを表示、なければステータス別サマリを表示
  $: hasPoolSummaries = $poolSummaries.length > 0;
  $: isGraphMode = $viewMode === "graph";
</script>

<header class="toolbar">
  <!-- 左側：ブランド -->
  <div class="toolbar-left">
    <BrandText size="sm" />
  </div>

  <!-- 中央：Pool別サマリ or ステータスサマリ -->
  <div class="toolbar-center">
    {#if hasPoolSummaries}
      <!-- Pool別サマリ -->
      <div class="pool-summary">
        {#each $poolSummaries as pool (pool.poolId)}
          <Badge variant="glass" color="neutral" size="medium">
            <span class="pool-name">{pool.poolId}</span>
            <span class="pool-separator">:</span>
            {#if pool.running > 0}
              <span class="pool-stat running">{pool.running} 実行中</span>
            {/if}
            {#if pool.queued > 0}
              <span class="pool-stat queued">{pool.queued} 待機</span>
            {/if}
            {#if pool.failed > 0}
              <span class="pool-stat failed">{pool.failed} 失敗</span>
            {/if}
            {#if pool.running === 0 && pool.queued === 0 && pool.failed === 0}
              <span class="pool-stat idle">{pool.total} タスク</span>
            {/if}
          </Badge>
        {/each}
      </div>
    {:else}
      <!-- フォールバック: ステータス別サマリ -->
      <div class="status-summary">
        {#each statusDisplay as { key, label, showCount, color }}
          {#if showCount && $taskCountsByStatus[key] > 0}
            <Badge variant="glass" {color} {label}>
              <span class="status-count">{$taskCountsByStatus[key]}</span>
            </Badge>
          {/if}
        {/each}
      </div>
    {/if}
  </div>

  <!-- 右側：進捗・ビュー切替 -->
  <div class="toolbar-right">
    <!-- 実行コントロール -->
    <ExecutionControls />

    <!-- 進捗率バー -->
    <div class="progress-section">
      <ProgressBar percentage={$overallProgress.percentage} size="mini" />
      <span class="progress-text">{$overallProgress.percentage}%</span>
    </div>

    <!-- ビュー切替 -->
    <div class="view-toggle">
      <Button
        variant={isGraphMode ? "secondary" : "ghost"}
        size="small"
        on:click={() => viewMode.setGraph()}
        title="グラフビュー"
      >
        <Network size="14" />
        <span>Graph</span>
      </Button>
      <Button
        variant={!isGraphMode ? "secondary" : "ghost"}
        size="small"
        on:click={() => viewMode.setWBS()}
        title="WBSビュー"
      >
        <ListTree size="14" />
        <span>WBS</span>
      </Button>
    </div>
  </div>
</header>
```

<style>
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: var(--mv-layout-toolbar-height);
    padding: 0 var(--mv-spacing-lg);

    /* Crystal HUD Style */
    background: var(--mv-glass-bg);
    backdrop-filter: blur(12px);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-glass-border-subtle);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);

    flex-shrink: 0;
    gap: var(--mv-spacing-md);
    position: relative;
    z-index: 100;
  }

  .toolbar-left,
  .toolbar-center,
  .toolbar-right {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
  }

  .toolbar-left {
    flex: 1;
    justify-content: flex-start;
  }

  .toolbar-center {
    flex: 2;
    justify-content: center;
  }

  .toolbar-right {
    flex: 1;
    justify-content: flex-end;
    gap: var(--mv-spacing-md);
  }

  /* ステータスサマリ */
  .status-summary {
    display: flex;
    gap: var(--mv-spacing-sm);
  }

  .status-count {
    font-weight: var(--mv-font-weight-bold);
    margin-right: var(--mv-spacing-xxs);
  }

  /* Pool別サマリ */
  .pool-summary {
    display: flex;
    gap: var(--mv-spacing-md);
  }

  .pool-name {
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    font-family: var(--mv-font-mono);
  }

  .pool-separator {
    color: var(--mv-color-text-muted);
    margin: 0 var(--mv-spacing-xxs);
  }

  .pool-stat {
    margin-left: var(--mv-spacing-xs);
    font-weight: var(--mv-font-weight-medium);
  }

  .pool-stat.running {
    color: var(--mv-color-status-running-text);
    text-shadow: 0 0 5px var(--mv-color-status-running-glow);
  }

  .pool-stat.queued {
    color: var(--mv-color-status-pending-text);
  }

  .pool-stat.failed {
    color: var(--mv-color-status-failed-text);
  }

  .pool-stat.idle {
    color: var(--mv-color-text-muted);
  }

  /* 進捗バー（ミニ） - styles handled by ProgressBar component */
  .progress-section {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xs);
  }

  .progress-text {
    font-size: var(--mv-font-size-xs);
    font-family: var(--mv-font-mono);
    color: var(--mv-color-text-muted);
    min-width: var(--mv-progress-text-width-sm);
    text-align: right;
  }

  /* ビュー切り替え - Cleaned up to rely on Button component */
  .view-toggle {
    display: flex;
    gap: var(--mv-spacing-xxs);
    /* Remove background/border to let Buttons stand on glass */
    /* background: var(--mv-color-surface-secondary); */
    /* border: var(--mv-border-width-thin) solid var(--mv-color-border-default); */
    /* border-radius: var(--mv-radius-sm); */
    /* padding: 2px; */
  }
</style>
