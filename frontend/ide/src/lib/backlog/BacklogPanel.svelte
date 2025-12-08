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
    height: 100%;

    /* Glassmorphism Background - Darker for sidebar */
    background: rgba(11, 15, 23, 0.65);
    backdrop-filter: blur(24px) saturate(140%);
    -webkit-backdrop-filter: blur(24px) saturate(140%);

    /* Subtle glass border */
    border-left: 1px solid rgba(255, 255, 255, 0.08);

    /* Soft ambient glow */
    box-shadow: -10px 0 30px rgba(0, 0, 0, 0.3);
  }

  /* === Header with HUD styling === */
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: var(--mv-size-floating-header);
    padding: 0 var(--mv-spacing-md);

    background: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.03),
      transparent
    );
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
    color: var(--mv-color-text-secondary);
  }

  /* Icon handled by Lucide classes global or via :global */
  :global(.header-icon) {
    opacity: 0.7;
    margin-top: -1px; /* Visual alignment */
  }

  .panel-header h3 {
    margin: 0;
    font-size: var(--mv-font-size-sm);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--mv-color-text-primary);
  }

  /* === Scrollable Content === */
  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--mv-spacing-md);

    /* Smooth scrollbar */
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.1) transparent;
  }

  .panel-content::-webkit-scrollbar {
    width: 4px;
  }

  .panel-content::-webkit-scrollbar-track {
    background: transparent;
  }

  .panel-content::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
  }

  .panel-content::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.2);
  }

  /* === Backlog List === */
  .backlog-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-md);
  }
</style>
