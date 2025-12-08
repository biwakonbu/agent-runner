<script lang="ts">
  import { stopPropagation } from "svelte/legacy";
  import { Minus, X } from "lucide-svelte";

  interface Props {
    id?: string;
    initialPosition?: { x: number; y: number };
    initialSize?: { width: number; height: number };
    title?: string;
    controls?: { minimize: boolean; close: boolean };
    zIndex?: number;
    header?: import("svelte").Snippet;
    children?: import("svelte").Snippet;
    footer?: import("svelte").Snippet;
    // コールバックプロップ
    onclose?: () => void;
    onminimize?: (data: { minimized: boolean }) => void;
    onclick?: () => void;
    ondragend?: (data: { x: number; y: number }) => void;
    onresizeend?: (data: { width: number; height: number }) => void;
  }

  let {
    id = "unknown",
    initialPosition = { x: 20, y: 20 },
    initialSize = undefined,
    title = "",
    controls = { minimize: true, close: true },
    zIndex = 100,
    header,
    children,
    footer,
    onclose,
    onminimize,
    onclick,
    ondragend,
    onresizeend,
  }: Props = $props();

  let position = $state({ x: 0, y: 0 });
  let isDragging = false;
  let isResizing = false;
  let windowEl: HTMLElement | undefined = $state();
  let isMinimized = $state(false);
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

  function toggleMinimize() {
    isMinimized = !isMinimized;
    onminimize?.({ minimized: isMinimized });
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
  class:minimized={isMinimized}
  style:top="{position.y}px"
  style:left="{position.x}px"
  style:z-index={zIndex}
  style:width={size ? `${size.width}px` : undefined}
  style:height={size && !isMinimized ? `${size.height}px` : undefined}
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
      {#if controls.minimize}
        <button
          class="control-btn"
          onclick={stopPropagation(toggleMinimize)}
          aria-label="Minimize"
          type="button"
        >
          <Minus size={14} strokeWidth={2.5} />
        </button>
      {/if}
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

  {#if !isMinimized}
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
  {/if}
</div>

<style>
  .floating-window {
    position: fixed;

    /* Default dims if not provided */
    min-width: var(--mv-floating-window-min-width);
    min-height: var(--mv-floating-window-min-height);

    /* Crystal HUD Redesign */
    background: rgba(10, 10, 12, 0.8); /* Neutral dark glass, no navy tint */
    backdrop-filter: blur(24px) saturate(140%);
    -webkit-backdrop-filter: blur(24px) saturate(140%);

    /* Assertive Gradient Border */
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-top: 1px solid rgba(255, 255, 255, 0.12); /* Subtle highlight */

    border-radius: var(--mv-radius-lg);

    /* Deep Shadow */
    box-shadow:
      0 20px 25px -5px rgba(0, 0, 0, 0.3),
      0 8px 10px -6px rgba(0, 0, 0, 0.2),
      0 0 0 1px rgba(0, 0, 0, 0.4); /* Definition outline */

    display: flex;
    flex-direction: column;

    overflow: hidden;
    transition:
      height 0.2s cubic-bezier(0.16, 1, 0.3, 1),
      box-shadow 0.2s ease;
  }

  .floating-window:focus-within {
    border-color: rgba(255, 255, 255, 0.15);
    box-shadow:
      0 25px 50px -12px rgba(0, 0, 0, 0.4),
      0 0 0 1px rgba(136, 192, 208, 0.3); /* Subtle blue glow on focus */
  }

  .floating-window.minimized {
    height: var(--mv-size-floating-header) !important;
    overflow: hidden;
    background: rgba(15, 23, 42, 0.9);
  }

  /* Header Area */
  .header {
    height: var(--mv-size-floating-header);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 var(--mv-spacing-md);
    cursor: grab;
    user-select: none;
    flex-shrink: 0;

    background: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.03),
      transparent
    );
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }

  .header:active {
    cursor: grabbing;
    background: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.05),
      rgba(255, 255, 255, 0.01)
    );
  }

  .header-content {
    flex: 1;
    overflow: hidden;
  }

  .title {
    font-size: var(--mv-font-size-sm);
    font-weight: 600;
    color: var(--mv-color-text-primary);
    letter-spacing: 0.01em;
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
    width: 24px;
    height: 24px;
    background: transparent;
    border: none;
    border-radius: 6px;
    color: var(--mv-color-text-muted);
    cursor: pointer;
    padding: 0;
    transition: all 0.15s ease;
  }

  .control-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: var(--mv-color-text-primary);
  }

  .control-btn.close:hover {
    background: rgba(239, 68, 68, 0.2);
    color: #f87171;
  }

  .content {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    position: relative;

    /* Ensure code remains readable on top of blur */
    background: rgba(0, 0, 0, 0.1);
  }

  .footer {
    flex-shrink: 0;
    padding: var(--mv-spacing-sm) var(--mv-spacing-md);
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    background: rgba(0, 0, 0, 0.2);
  }

  .resize-handle {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 16px;
    height: 16px;
    cursor: nwse-resize;
    z-index: 10;
  }

  /* Visual indicator for resize handle */
  .resize-handle::after {
    content: "";
    position: absolute;
    bottom: 4px;
    right: 4px;
    width: 6px;
    height: 6px;
    border-right: 2px solid rgba(255, 255, 255, 0.1);
    border-bottom: 2px solid rgba(255, 255, 255, 0.1);
    border-bottom-right-radius: 2px;
    transition: border-color 0.2s;
  }

  .resize-handle:hover::after {
    border-color: var(--mv-color-text-secondary);
  }
</style>
