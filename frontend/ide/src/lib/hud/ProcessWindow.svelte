<script lang="ts">
  import DraggableWindow from "../components/ui/window/DraggableWindow.svelte";
  import ResourceList from "./ResourceList.svelte";
  import type { ResourceNode } from "./types";
  import { windowStore } from "../../stores/windowStore";
  import { Cpu } from "lucide-svelte";

  interface Props {
    resources?: ResourceNode[];
  }

  let { resources = [] }: Props = $props();

  let isOpen = $derived($windowStore.process.isOpen);
  let position = $derived($windowStore.process.position);
  let size = $derived($windowStore.process.size);
  let zIndex = $derived($windowStore.process.zIndex);

  function handleClose() {
    windowStore.close("process");
  }

  function handleMinimize(data: { minimized: boolean }) {
    windowStore.minimize("process", data.minimized);
  }

  function handleDragEnd(data: { x: number; y: number }) {
    windowStore.updatePosition("process", data.x, data.y);
  }

  function handleResizeEnd(data: { width: number; height: number }) {
    windowStore.updateSize("process", data.width, data.height);
  }

  function handleClick() {
    windowStore.bringToFront("process");
  }
</script>

{#if isOpen}
  <DraggableWindow
    id="process"
    initialPosition={position}
    initialSize={size}
    {zIndex}
    onclose={handleClose}
    onminimize={handleMinimize}
    ondragend={handleDragEnd}
    onresizeend={handleResizeEnd}
    onclick={handleClick}
  >
    {#snippet header()}
      <div class="window-header">
        <Cpu size={16} class="header-icon" />
        <span class="header-title">Process & Resources</span>
      </div>
    {/snippet}

    {#snippet children()}
      <div class="resource-window-content">
        <ResourceList {resources} />
      </div>
    {/snippet}
  </DraggableWindow>
{/if}

<style>
  /* Header Styling matching other windows */
  .window-header {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
    color: var(--mv-color-text-secondary);
  }

  :global(.header-icon) {
    opacity: var(--mv-opacity-70);
  }

  .header-title {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    letter-spacing: var(--mv-letter-spacing-widest);
  }

  .resource-window-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--mv-spacing-xs);

    /* Matches ResourceList usage in ProcessHUD but adapted for window */
  }
</style>
