<script lang="ts">
  import {
    backlogItems,
    unresolvedCount,
    resolveItem,
    deleteItem,
    type BacklogItem,
  } from "../../stores/backlogStore";
  import BacklogItemComponent from "./components/BacklogItem.svelte";
  import ResolveDialog from "./components/ResolveDialog.svelte";
  import EmptyBacklog from "./components/EmptyBacklog.svelte";
  import { ClipboardList } from "lucide-svelte";

  // 解決ダイアログ
  let resolvingItem: BacklogItem | null = $state(null);

  function openResolveDialog(item: BacklogItem) {
    resolvingItem = item;
  }

  function closeResolveDialog() {
    resolvingItem = null;
  }

  async function handleResolve(event: { text: string }) {
    if (!resolvingItem) return;
    try {
      await resolveItem(resolvingItem.id, event.text);
      closeResolveDialog();
    } catch {
      // エラーは store でログ出力済み
    }
  }

  async function handleDelete(item: BacklogItem) {
    if (confirm(`「${item.title}」を削除しますか？`)) {
      try {
        await deleteItem(item.id);
      } catch {
        // エラーは store でログ出力済み
      }
    }
  }
</script>

<aside class="backlog-panel">
  <header class="panel-header">
    <div class="header-title">
      <ClipboardList size={16} class="header-icon" />
      <h3>Backlog ({$unresolvedCount})</h3>
    </div>
  </header>

  <div class="panel-content">
    {#if $backlogItems.length === 0}
      <EmptyBacklog />
    {:else}
      <ul class="backlog-list">
        {#each $backlogItems as item (item.id)}
          <BacklogItemComponent
            {item}
            onresolve={() => openResolveDialog(item)}
            ondelete={() => handleDelete(item)}
          />
        {/each}
      </ul>
    {/if}
  </div>
</aside>

<!-- 解決ダイアログ -->
{#if resolvingItem}
  <ResolveDialog
    item={resolvingItem}
    onclose={closeResolveDialog}
    onconfirm={handleResolve}
  />
{/if}

<style>
  /* === Crystal Glass Panel === */
  .backlog-panel {
    display: flex;
    flex-direction: column;
    height: var(--mv-size-full);

    /* Glassmorphism Background - Darker for sidebar */
    background: var(--mv-glass-bg);
    backdrop-filter: blur(24px) saturate(140%);

    /* Subtle glass border */
    border-left: var(--mv-border-width-thin) solid var(--mv-window-border);

    /* Soft ambient glow */
    box-shadow: var(--mv-shadow-backlog-panel);
  }

  /* === Header with HUD styling === */
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: var(--mv-size-floating-header);
    padding: var(--mv-space-0) var(--mv-spacing-md);

    background: linear-gradient(
      to bottom,
      var(--mv-window-header-gradient-start),
      transparent
    );
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-window-border-bottom);
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
    color: var(--mv-color-text-secondary);
  }

  /* Icon handled by Lucide classes global or via :global */
  :global(.header-icon) {
    opacity: var(--mv-opacity-70);
    margin-top: calc(-1 * var(--mv-space-px));
  }

  .panel-header h3 {
    margin: var(--mv-space-0);
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-count);
    color: var(--mv-color-text-primary);
  }

  /* === Scrollable Content === */
  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--mv-spacing-md);

    /* Smooth scrollbar */
    scrollbar-width: thin;
    scrollbar-color: var(--mv-glass-scrollbar) transparent;
  }

  .panel-content::-webkit-scrollbar {
    width: var(--mv-size-scrollbar);
  }

  .panel-content::-webkit-scrollbar-track {
    background: transparent;
  }

  .panel-content::-webkit-scrollbar-thumb {
    background: var(--mv-glass-scrollbar);
    border-radius: var(--mv-radius-sm);
  }

  .panel-content::-webkit-scrollbar-thumb:hover {
    background: var(--mv-glass-scrollbar-hover);
  }

  /* === Backlog List === */
  .backlog-list {
    list-style: none;
    margin: var(--mv-space-0);
    padding: var(--mv-space-0);
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-md);
  }
</style>
