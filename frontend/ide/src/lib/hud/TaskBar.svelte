<script lang="ts">
  import { windowStore, type WindowId } from "../../stores/windowStore";
  import { unresolvedCount } from "../../stores/backlogStore";
  import { MessageSquare, Cpu, ListTodo, ClipboardList } from "lucide-svelte";

  function toggle(id: WindowId) {
    windowStore.toggle(id);
  }
</script>

<div class="taskbar-container">
  <div class="taskbar glass-panel">
    <!-- Chat Toggle -->
    <button
      class="taskbar-item"
      class:active={$windowStore.chat.isOpen}
      onclick={() => toggle("chat")}
      title="Chat"
      aria-label="Toggle Chat"
    >
      <div class="icon-wrapper">
        <MessageSquare size={20} absoluteStrokeWidth class="icon" />
      </div>
      {#if $windowStore.chat.isOpen}
        <div class="active-glow"></div>
      {/if}
    </button>

    <!-- Process Toggle -->
    <button
      class="taskbar-item"
      class:active={$windowStore.process.isOpen}
      onclick={() => toggle("process")}
      title="Process & Resources"
      aria-label="Toggle Process View"
    >
      <div class="icon-wrapper">
        <Cpu size={20} absoluteStrokeWidth class="icon" />
      </div>
      {#if $windowStore.process.isOpen}
        <div class="active-glow"></div>
      {/if}
    </button>

    <!-- Backlog Toggle -->
    <button
      class="taskbar-item"
      class:active={$windowStore.backlog.isOpen}
      onclick={() => toggle("backlog")}
      title="Backlog"
      aria-label="Toggle Backlog"
    >
      <div class="icon-wrapper">
        <ClipboardList size={20} absoluteStrokeWidth class="icon" />
        {#if $unresolvedCount > 0}
          <span class="badge">{$unresolvedCount}</span>
        {/if}
      </div>
      {#if $windowStore.backlog.isOpen}
        <div class="active-glow"></div>
      {/if}
    </button>
  </div>
</div>

<style>
  .taskbar-container {
    position: fixed;
    bottom: var(--mv-spacing-lg);
    left: 50%; /* Center horizontally */
    transform: translateX(-50%);
    z-index: 2000;
  }

  .taskbar {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm); /* Increased gap */
    padding: var(--mv-spacing-xs) var(--mv-spacing-md);

    /* Sophisticated Glassmorphism */
    background: rgba(15, 23, 42, 0.6); /* Darker base for contrast */
    backdrop-filter: blur(20px) saturate(180%);
    -webkit-backdrop-filter: blur(20px) saturate(180%);

    border: 1px solid rgba(255, 255, 255, 0.08);
    border-top: 1px solid rgba(255, 255, 255, 0.15); /* Highlight top edge */
    border-radius: 9999px; /* Full pill shape */

    box-shadow:
      0 4px 6px -1px rgba(0, 0, 0, 0.1),
      0 2px 4px -1px rgba(0, 0, 0, 0.06),
      0 0 0 1px rgba(255, 255, 255, 0.05) inset; /* Inner ring */

    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .taskbar:hover {
    background: rgba(30, 41, 59, 0.7);
    box-shadow:
      0 10px 15px -3px rgba(0, 0, 0, 0.1),
      0 4px 6px -2px rgba(0, 0, 0, 0.05),
      0 0 0 1px rgba(255, 255, 255, 0.1) inset;
    transform: translateY(-2px);
  }

  .taskbar-item {
    position: relative;
    width: 44px; /* Fixed touch target size */
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    border-radius: 50%; /* Circular buttons */
    cursor: pointer;
    transition: all 0.2s ease;

    color: var(--mv-color-text-secondary);
  }

  .taskbar-item:hover {
    background: rgba(255, 255, 255, 0.1);
    color: var(--mv-color-text-primary);
  }

  .taskbar-item:active {
    transform: scale(0.95);
  }

  .taskbar-item.active {
    color: var(--mv-primitive-frost-2); /* Bright accent color */
    background: rgba(136, 192, 208, 0.15); /* Subtle tint of accent */
    box-shadow: 0 0 0 1px rgba(136, 192, 208, 0.2) inset;
  }

  .active-glow {
    position: absolute;
    bottom: -6px;
    left: 50%;
    transform: translateX(-50%);
    width: 4px;
    height: 4px;
    background-color: var(--mv-primitive-frost-2);
    border-radius: 50%;
    box-shadow: 0 0 8px 2px rgba(136, 192, 208, 0.6);
  }

  .icon-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Icon styling handled by Lucide component classes via global CSS or simpler: */
  :global(.icon) {
    stroke-width: 2px;
    opacity: 0.8;
    transition: opacity 0.2s;
  }

  .taskbar-item:hover :global(.icon),
  .taskbar-item.active :global(.icon) {
    opacity: 1;
  }

  /* Badge styling */
  .badge {
    position: absolute;
    top: -6px;
    right: -8px;

    background: var(--mv-primitive-aurora-red);
    color: white;

    font-size: 10px;
    font-weight: 700;
    line-height: 1;

    padding: 3px 5px;
    border-radius: 10px;
    border: 2px solid rgba(15, 23, 42, 0.8); /* Match backdrop to cut out */

    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    min-width: 16px;
    text-align: center;
  }
</style>
