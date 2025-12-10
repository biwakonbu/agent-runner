<script lang="ts">
  import { stopPropagation } from "svelte/legacy";
  import { X } from "lucide-svelte";

  interface Props {
    id?: string;
    initialPosition?: { x: number; y: number };
    initialSize?: { width: number; height: number };
    title?: string;
    controls?: { close: boolean };
    zIndex?: number;
    header?: import("svelte").Snippet;
    children?: import("svelte").Snippet;
    footer?: import("svelte").Snippet;
    // コールバックプロップ
    onclose?: () => void;
    onclick?: () => void;
    ondragend?: (data: { x: number; y: number }) => void;
    onresizeend?: (data: { width: number; height: number }) => void;
  }

  let {
    id = "unknown",
    initialPosition = { x: 20, y: 20 },
    initialSize = undefined,
    title = "",
    controls = { close: true },
    zIndex = 100,
    header,
    children,
    footer,
    onclose,
    onclick,
    ondragend,
    onresizeend,
  }: Props = $props();

  let position = $state({ x: 0, y: 0 });
  let isDragging = false;
  let isResizing = false;
  let windowEl: HTMLElement | undefined = $state();

  let size = $state<{ width: number; height: number } | undefined>(undefined);

  // Sync initial values when props change
  $effect(() => {
    position = { ...initialPosition };
  });

  $effect(() => {
    size = initialSize;
  });

  function startDrag(e: MouseEvent) {
    if (e.button !== 0) return;
    if ((e.target as HTMLElement).closest(".window-controls")) return;
    if (!windowEl) return;

    isDragging = true;
    (windowEl as HTMLElement).style.cursor = "grabbing";
    window.addEventListener("mouseup", stopDrag);
    onclick?.();
  }

  function onMouseMove(e: MouseEvent) {
    if (isDragging) {
      position = {
        x: position.x + e.movementX,
        y: position.y + e.movementY,
      };
    } else if (isResizing && size) {
      size = {
        width: Math.max(200, size.width + e.movementX),
        height: Math.max(100, size.height + e.movementY),
      };
    }
  }

  function stopDrag() {
    isDragging = false;
    if (windowEl) (windowEl as HTMLElement).style.cursor = "";
    window.removeEventListener("mouseup", stopDrag);
    ondragend?.(position);
  }

  function startResize(e: MouseEvent) {
    if (e.button !== 0) return;
    e.stopPropagation();
    isResizing = true;
    window.addEventListener("mouseup", stopResize);
  }

  function stopResize() {
    isResizing = false;
    window.removeEventListener("mouseup", stopResize);
    if (size) {
      onresizeend?.(size);
    }
  }

  function closeWindow() {
    onclose?.();
  }

  function onWindowClick() {
    onclick?.();
  }
</script>

<svelte:window onmousemove={onMouseMove} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="floating-window"
  style:top="{position.y}px"
  style:left="{position.x}px"
  style:z-index={zIndex}
  style:width={size ? `${size.width}px` : undefined}
  style:height={size ? `${size.height}px` : undefined}
  bind:this={windowEl}
  onmousedown={onWindowClick}
>
  <div class="header" onmousedown={startDrag}>
    <div class="header-content">
      {#if header}{@render header()}{:else}
        <span class="title">{title}</span>
      {/if}
    </div>
    <div class="window-controls">
      {#if controls.close}
        <button
          class="control-btn close"
          onclick={stopPropagation(closeWindow)}
          aria-label="Close"
          type="button"
        >
          <X size={14} strokeWidth={2.5} />
        </button>
      {/if}
    </div>
  </div>

  <div class="content">
    {@render children?.()}
  </div>

  {#if footer}
    <div class="footer">
      {@render footer?.()}
    </div>
  {/if}

  <div
    class="resize-handle"
    onmousedown={startResize}
    role="button"
    tabindex="0"
  ></div>
</div>

<style>
  .floating-window {
    position: fixed;

    /* Default dims if not provided */
    min-width: var(--mv-floating-window-min-width);
    min-height: var(--mv-floating-window-min-height);

    /* Crystal HUD Redesign */
    background: var(--mv-window-bg);
    backdrop-filter: blur(24px) saturate(140%);

    /* Assertive Gradient Border */
    border: var(--mv-border-width-thin) solid var(--mv-window-border);
    border-top: var(--mv-border-width-thin) solid var(--mv-window-border-top);

    border-radius: var(--mv-radius-lg);

    /* Deep Shadow */
    box-shadow: var(--mv-shadow-window);

    display: flex;
    flex-direction: column;

    overflow: hidden;
    transition:
      height 0.2s cubic-bezier(0.16, 1, 0.3, 1),
      box-shadow 0.2s ease;
  }

  .floating-window:focus-within {
    border-color: var(--mv-window-border-focus);
    box-shadow: var(--mv-shadow-window-focus);
  }

  /* Header Area */
  .header {
    height: var(--mv-size-floating-header);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--mv-space-0) var(--mv-spacing-md);
    cursor: grab;
    user-select: none;
    flex-shrink: 0;

    background: linear-gradient(
      to bottom,
      var(--mv-window-header-gradient-start),
      transparent
    );
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-window-border-bottom);
  }

  .header:active {
    cursor: grabbing;
    background: linear-gradient(
      to bottom,
      var(--mv-window-header-gradient-active),
      var(--mv-window-header-gradient-active-end)
    );
  }

  .header-content {
    flex: 1;
    overflow: hidden;
  }

  .title {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    letter-spacing: var(--mv-letter-spacing-widest);
  }

  .window-controls {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xs);
    margin-left: var(--mv-spacing-sm);
  }

  .control-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--mv-icon-size-lg);
    height: var(--mv-icon-size-lg);
    background: transparent;
    border: none;
    border-radius: var(--mv-radius-sm);
    color: var(--mv-color-text-muted);
    cursor: pointer;
    padding: var(--mv-space-0);
    transition: all 0.15s ease;
  }

  .control-btn:hover {
    background: var(--mv-btn-hover-bg);
    color: var(--mv-color-text-primary);
  }

  .control-btn.close:hover {
    background: var(--mv-btn-close-hover-bg);
    color: var(--mv-primitive-aurora-red);
  }

  .content {
    flex: 1;
    min-height: var(--mv-space-0);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    position: relative;

    /* Ensure code remains readable on top of blur */
    background: var(--mv-window-content-bg);
  }

  .footer {
    flex-shrink: 0;
    padding: var(--mv-spacing-sm) var(--mv-spacing-md);
    border-top: var(--mv-border-width-thin) solid var(--mv-window-border-bottom);
    background: var(--mv-window-footer-bg);
  }

  .resize-handle {
    position: absolute;
    bottom: var(--mv-space-0);
    right: var(--mv-space-0);
    width: var(--mv-icon-size-sm);
    height: var(--mv-icon-size-sm);
    cursor: nwse-resize;
    z-index: var(--mv-z-index-base);
  }

  /* Visual indicator for resize handle */
  .resize-handle::after {
    content: "";
    position: absolute;
    bottom: var(--mv-space-1);
    right: var(--mv-space-1);
    width: var(--mv-indicator-size-xs);
    height: var(--mv-indicator-size-xs);
    border-right: var(--mv-border-width-md) solid var(--mv-resize-handle-border);
    border-bottom: var(--mv-border-width-md) solid
      var(--mv-resize-handle-border);
    border-bottom-right-radius: var(--mv-radius-progress);
    transition: border-color 0.2s;
  }

  .resize-handle:hover::after {
    border-color: var(--mv-color-text-secondary);
  }
</style>
