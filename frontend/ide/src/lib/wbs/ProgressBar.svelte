<script lang="ts">
  import { getProgressColor } from "./utils";

  export let percentage: number = 0;
  export let size: "sm" | "md" | "mini" = "sm";
  export let className: string = "";

  $: progressColor = getProgressColor(percentage);

  // Calculate dynamic shadow and background for the container
  $: containerShadow =
    size === "md"
      ? `0 0 2px ${progressColor.glow}, inset 0 1px 2px rgba(0, 0, 0, 0.2)`
      : `inset 0 1px 2px rgba(0, 0, 0, 0.3)`;
</script>

<div
  class="progress-bar {size} {className}"
  style:box-shadow={containerShadow}
  style:background-color={progressColor.bg}
>
  <div
    class="progress-fill"
    style:width="{percentage}%"
    style:background-color={progressColor.fill}
    style:box-shadow="0 0 6px {progressColor.glow}"
  ></div>
</div>

<style>
  .progress-bar {
    /* Background is now handled inline via style:background-color */
    border-radius: var(--mv-radius-sm);
    overflow: hidden;
    /* box-shadow is now handled inline */
  }

  /* Size variants */
  .progress-bar.sm {
    width: var(--mv-progress-bar-width-sm);
    height: var(--mv-progress-bar-height-sm);
  }

  .progress-bar.mini {
    width: var(--mv-progress-bar-width-mini);
    height: var(--mv-progress-bar-height-sm);
  }

  .progress-bar.md {
    width: 100%; /* Header bar takes full width of container */
    height: var(--mv-progress-bar-height-md);
    /* Background is now handled inline */
    /* box-shadow is now handled inline */
    border: var(--mv-border-panel); /* Header bar had border */
  }

  .progress-fill {
    height: 100%;
    width: var(--progress, 0%);
    background: var(--mv-progress-bar-fill);
    border-radius: var(--mv-radius-sm);
    transition: width var(--mv-duration-slow);
    box-shadow: 0 0 3px var(--mv-progress-bar-fill-glow);
  }
</style>
