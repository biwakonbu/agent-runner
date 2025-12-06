<script lang="ts">
  /**
   * バッジコンポーネント
   * ステータス表示やラベル付けに使用
   */

  export let status:
    | "pending"
    | "ready"
    | "running"
    | "succeeded"
    | "failed"
    | "canceled"
    | "blocked"
    | undefined = undefined;

  export let variant: "default" | "outline" | "glass" = "default";
  export let color:
    | "primary"
    | "secondary"
    | "success"
    | "warning"
    | "danger"
    | "info"
    | "neutral" = "neutral";

  export let size: "small" | "medium" = "medium";
  export let label = "";

  // ステータスからカラーとラベルを自動解決
  const statusConfig: Record<
    NonNullable<typeof status>,
    { color: typeof color; label: string }
  > = {
    pending: { color: "warning", label: "Pending" },
    ready: { color: "info", label: "Ready" },
    running: { color: "success", label: "Running" },
    succeeded: { color: "primary", label: "Succeeded" },
    failed: { color: "danger", label: "Failed" },
    canceled: { color: "neutral", label: "Canceled" },
    blocked: { color: "secondary", label: "Blocked" },
  };

  $: resolvedColor = status ? statusConfig[status].color : color;
  $: resolvedLabel = label || (status ? statusConfig[status].label : "");
</script>

<span
  class="badge variant-{variant} color-{resolvedColor} size-{size}"
  class:pulse={status === "running"}
>
  {#if status === "running"}
    <span class="pulse-dot"></span>
  {/if}
  {resolvedLabel}
  <slot />
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    gap: var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-sm);
    font-family: var(--mv-font-sans);
    font-weight: var(--mv-font-weight-medium);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-widest);
    line-height: 1;
    white-space: nowrap;
  }

  /* サイズ */
  .size-small {
    padding: 2px 6px;
    font-size: var(--mv-font-size-xxs);
  }

  .size-medium {
    padding: 4px 8px;
    font-size: var(--mv-font-size-xs);
  }

  /* カラーマッピング (CSS変数) */
  .color-primary {
    --badge-bg: var(--mv-primitive-frost-3);
    --badge-border: var(--mv-primitive-frost-1);
    --badge-text: var(--mv-primitive-frost-1);
    --badge-glow: var(--mv-color-glow-focus);
  }
  .color-secondary {
    --badge-bg: var(--mv-color-surface-secondary);
    --badge-border: var(--mv-color-border-subtle);
    --badge-text: var(--mv-color-text-secondary);
  }
  .color-success {
    --badge-bg: rgba(163, 190, 140, 0.2);
    --badge-border: var(--mv-primitive-aurora-green);
    --badge-text: var(--mv-primitive-aurora-green);
    --badge-glow: var(--mv-color-glow-running);
  }
  .color-warning {
    --badge-bg: rgba(235, 203, 139, 0.2);
    --badge-border: var(--mv-primitive-aurora-yellow);
    --badge-text: var(--mv-primitive-aurora-yellow);
  }
  .color-danger {
    --badge-bg: rgba(191, 97, 106, 0.2);
    --badge-border: var(--mv-primitive-aurora-red);
    --badge-text: var(--mv-primitive-aurora-red);
  }
  .color-info {
    --badge-bg: rgba(136, 192, 208, 0.2);
    --badge-border: var(--mv-primitive-frost-2);
    --badge-text: var(--mv-primitive-frost-2);
  }
  .color-neutral {
    --badge-bg: var(--mv-color-surface-secondary);
    --badge-border: var(--mv-color-border-subtle);
    --badge-text: var(--mv-color-text-muted);
  }

  /* バリアント */

  /* Default: Flat/Solid-ish */
  .variant-default {
    background: var(--badge-bg);
    border: 1px solid var(--badge-border);
    color: var(--badge-text);
  }

  /* Outline: Transparent bg */
  .variant-outline {
    background: transparent;
    border: 1px solid var(--badge-border);
    color: var(--badge-text);
  }

  /* Glass: Rich transparent Look */
  .variant-glass {
    background: var(
      --badge-bg
    ); /* Use same semi-transparent base but maybe add blur? */
    background: color-mix(
      in srgb,
      var(--badge-bg),
      transparent 50%
    ); /* make it subtler */
    border: 1px solid color-mix(in srgb, var(--badge-border), transparent 30%);
    color: var(--badge-text);
    box-shadow: 0 0 8px var(--badge-bg); /* subtle glow */
    backdrop-filter: blur(4px);
  }

  /* パルスドット */
  .pulse-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: currentColor;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.4;
      transform: scale(0.8);
    }
    100% {
      opacity: 1;
      transform: scale(1);
    }
  }
</style>
