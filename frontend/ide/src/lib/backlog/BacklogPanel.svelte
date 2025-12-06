<script lang="ts">
  import {
    backlogItems,
    unresolvedCount,
    resolveItem,
    deleteItem,
    type BacklogItem,
  } from "../../stores/backlogStore";

  // 解決ダイアログ
  let resolvingItem: BacklogItem | null = null;
  let resolutionText = "";

  function openResolveDialog(item: BacklogItem) {
    resolvingItem = item;
    resolutionText = "";
  }

  function closeResolveDialog() {
    resolvingItem = null;
    resolutionText = "";
  }

  async function handleResolve() {
    if (!resolvingItem) return;
    try {
      await resolveItem(resolvingItem.id, resolutionText || "Resolved");
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

  function getTypeLabel(type: BacklogItem["type"]): string {
    switch (type) {
      case "FAILURE":
        return "失敗";
      case "QUESTION":
        return "質問";
      case "BLOCKER":
        return "ブロッカー";
      default:
        return type;
    }
  }

  function getPriorityLabel(priority: number): string {
    if (priority >= 5) return "最高";
    if (priority >= 4) return "高";
    if (priority >= 3) return "中";
    if (priority >= 2) return "低";
    return "最低";
  }

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString("ja-JP", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<aside class="backlog-panel">
  <header class="panel-header">
    <h3>バックログ ({$unresolvedCount})</h3>
  </header>

  <div class="panel-content">
    {#if $backlogItems.length === 0}
      <div class="empty-state">
        <span class="empty-icon">&#10003;</span>
        <p>バックログは空です</p>
      </div>
    {:else}
      <ul class="backlog-list">
        {#each $backlogItems as item (item.id)}
          <li class="backlog-item" class:failure={item.type === "FAILURE"}>
            <div class="item-header">
              <span class="type-badge {item.type.toLowerCase()}"
                >{getTypeLabel(item.type)}</span
              >
              <span class="priority">{getPriorityLabel(item.priority)}</span>
              <span class="date">{formatDate(item.createdAt)}</span>
            </div>
            <h4 class="item-title">{item.title}</h4>
            <p class="item-description">{item.description}</p>
            <div class="item-actions">
              <button
                class="btn-resolve"
                on:click={() => openResolveDialog(item)}
              >
                解決
              </button>
              <button class="btn-delete" on:click={() => handleDelete(item)}>
                削除
              </button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</aside>

<!-- 解決ダイアログ -->
{#if resolvingItem}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="dialog-overlay" on:click={closeResolveDialog}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="dialog" on:click|stopPropagation>
      <h4>バックログを解決</h4>
      <p class="dialog-item-title">{resolvingItem.title}</p>
      <label>
        解決方法:
        <textarea
          bind:value={resolutionText}
          placeholder="どのように解決したかを入力..."
          rows="3"
        />
      </label>
      <div class="dialog-actions">
        <button class="btn-cancel" on:click={closeResolveDialog}>
          キャンセル
        </button>
        <button class="btn-confirm" on:click={handleResolve}> 解決 </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* === Crystal Glass Panel === */
  .backlog-panel {
    display: flex;
    flex-direction: column;
    height: 100%;

    /* Glassmorphism Background */
    background: var(--mv-glass-bg);
    backdrop-filter: blur(16px);

    /* Subtle glass border */
    border-left: var(--mv-border-width-thin) solid var(--mv-glass-border-strong);

    /* Soft ambient glow */
    box-shadow:
      inset 1px 0 0 var(--mv-glass-border),
      -4px 0 24px rgba(0, 0, 0, 0.3);
  }

  /* === Header with HUD styling === */
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--mv-spacing-sm) var(--mv-spacing-md);
    background: var(--mv-glass-bg-strong);
    border-bottom: var(--mv-border-width-thin) solid var(--mv-glass-border);

    /* Inner glow effect */
    box-shadow:
      inset 0 -1px 0 var(--mv-glass-border-subtle),
      0 1px 4px rgba(0, 0, 0, 0.2);
  }

  .panel-header h3 {
    margin: 0;
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-bold);
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--mv-color-text-secondary);

    /* Glow text effect */
    text-shadow: 0 0 12px rgba(136, 192, 208, 0.3);
  }

  /* === Scrollable Content === */
  .panel-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--mv-spacing-md);

    /* Smooth scrollbar */
    scrollbar-width: thin;
    scrollbar-color: var(--mv-glass-border) transparent;
  }

  .panel-content::-webkit-scrollbar {
    width: 6px;
  }

  .panel-content::-webkit-scrollbar-track {
    background: transparent;
  }

  .panel-content::-webkit-scrollbar-thumb {
    background: var(--mv-glass-border);
    border-radius: 3px;
  }

  .panel-content::-webkit-scrollbar-thumb:hover {
    background: var(--mv-glass-border-strong);
  }

  /* === Empty State with Glow === */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--mv-color-text-muted);
    text-align: center;
    padding: var(--mv-spacing-xl);
  }

  .empty-icon {
    font-size: 3rem;
    margin-bottom: var(--mv-spacing-md);
    color: var(--mv-primitive-aurora-green);

    /* Success glow */
    filter: drop-shadow(0 0 12px rgba(163, 190, 140, 0.6));
    animation: gentle-pulse 3s ease-in-out infinite;
  }

  @keyframes gentle-pulse {
    0%,
    100% {
      opacity: 0.8;
      transform: scale(1);
    }
    50% {
      opacity: 1;
      transform: scale(1.05);
    }
  }

  .empty-state p {
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-muted);
    letter-spacing: 0.05em;
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

  /* === Backlog Item Card === */
  .backlog-item {
    position: relative;

    /* Glass Card */
    background: var(--mv-glass-bg-strong);
    backdrop-filter: blur(8px);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border);
    border-radius: var(--mv-radius-md);
    padding: var(--mv-spacing-md);

    /* Card shadow */
    box-shadow: var(--mv-shadow-card);

    /* Animation */
    transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
    overflow: hidden;
  }

  .backlog-item::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 3px;
    height: 100%;
    background: var(--mv-glass-border);
    opacity: 0.5;
    transition: all 0.25s ease;
  }

  .backlog-item:hover {
    background: var(--mv-glass-hover);
    border-color: var(--mv-glass-border-strong);
    transform: translateX(4px);

    box-shadow:
      var(--mv-shadow-card),
      0 0 20px rgba(136, 192, 208, 0.1);
  }

  .backlog-item:hover::before {
    opacity: 1;
    background: var(--mv-primitive-frost-2);
    box-shadow: 0 0 8px var(--mv-primitive-frost-2);
  }

  /* === Failure Type - Glowing Red Edge === */
  .backlog-item.failure::before {
    background: var(--mv-primitive-aurora-red);
    opacity: 1;
    box-shadow: 0 0 8px var(--mv-primitive-aurora-red);
  }

  .backlog-item.failure:hover::before {
    box-shadow: 0 0 16px var(--mv-primitive-aurora-red);
  }

  /* === Item Header === */
  .item-header {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
    margin-bottom: var(--mv-spacing-sm);
    font-size: var(--mv-font-size-xs);
  }

  /* === Type Badge with Glow === */
  .type-badge {
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
    font-weight: var(--mv-font-weight-bold);
    text-transform: uppercase;
    font-size: var(--mv-font-size-xxs);
    letter-spacing: 0.08em;

    /* Glass effect */
    backdrop-filter: blur(4px);
    border: var(--mv-border-width-thin) solid transparent;
  }

  .type-badge.failure {
    background: rgba(191, 97, 106, 0.2);
    color: var(--mv-primitive-aurora-red);
    border-color: rgba(191, 97, 106, 0.4);
    box-shadow: 0 0 8px rgba(191, 97, 106, 0.3);
  }

  .type-badge.question {
    background: rgba(235, 203, 139, 0.2);
    color: var(--mv-primitive-aurora-yellow);
    border-color: rgba(235, 203, 139, 0.4);
    box-shadow: 0 0 8px rgba(235, 203, 139, 0.3);
  }

  .type-badge.blocker {
    background: rgba(136, 192, 208, 0.2);
    color: var(--mv-primitive-frost-2);
    border-color: rgba(136, 192, 208, 0.4);
    box-shadow: 0 0 8px rgba(136, 192, 208, 0.3);
  }

  /* === Priority Badge === */
  .priority {
    font-size: var(--mv-font-size-xxs);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-muted);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    background: var(--mv-glass-active);
    border-radius: var(--mv-radius-sm);
    letter-spacing: 0.05em;
  }

  /* === Date === */
  .date {
    margin-left: auto;
    font-size: var(--mv-font-size-xxs);
    font-family: var(--mv-font-mono);
    color: var(--mv-color-text-disabled);
    letter-spacing: 0.02em;
  }

  /* === Item Title === */
  .item-title {
    margin: 0 0 var(--mv-spacing-xs);
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    line-height: 1.4;

    /* Subtle glow on text */
    text-shadow: 0 0 20px rgba(236, 239, 244, 0.1);
  }

  /* === Item Description === */
  .item-description {
    margin: 0 0 var(--mv-spacing-sm);
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-secondary);
    line-height: var(--mv-line-height-relaxed);
    opacity: 0.9;
  }

  /* === Error Detail Box === */
  .error-detail {
    margin: var(--mv-spacing-sm) 0;
    padding: var(--mv-spacing-sm);
    background: rgba(191, 97, 106, 0.1);
    border: var(--mv-border-width-thin) solid rgba(191, 97, 106, 0.3);
    border-radius: var(--mv-radius-sm);
    font-size: var(--mv-font-size-xs);
    font-family: var(--mv-font-mono);
    color: var(--mv-primitive-aurora-red);
    overflow-x: auto;
    white-space: pre-wrap;
    word-break: break-all;
  }

  /* === Action Buttons === */
  .item-actions {
    display: flex;
    gap: var(--mv-spacing-sm);
    margin-top: var(--mv-spacing-md);
    padding-top: var(--mv-spacing-sm);
    border-top: var(--mv-border-width-thin) solid var(--mv-glass-border-subtle);
  }

  .btn-resolve,
  .btn-delete {
    padding: var(--mv-spacing-xs) var(--mv-spacing-md);
    border-radius: var(--mv-radius-sm);
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-semibold);
    letter-spacing: 0.05em;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  /* === Resolve Button === */
  .btn-resolve {
    background: rgba(163, 190, 140, 0.2);
    color: var(--mv-primitive-aurora-green);
    border: var(--mv-border-width-thin) solid rgba(163, 190, 140, 0.4);
  }

  .btn-resolve:hover {
    background: rgba(163, 190, 140, 0.3);
    border-color: var(--mv-primitive-aurora-green);
    box-shadow: 0 0 12px rgba(163, 190, 140, 0.4);
    transform: translateY(-1px);
  }

  /* === Delete Button === */
  .btn-delete {
    background: transparent;
    color: var(--mv-color-text-muted);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border);
  }

  .btn-delete:hover {
    background: rgba(191, 97, 106, 0.2);
    color: var(--mv-primitive-aurora-red);
    border-color: rgba(191, 97, 106, 0.5);
    box-shadow: 0 0 12px rgba(191, 97, 106, 0.3);
  }

  /* === Dialog Overlay === */
  .dialog-overlay {
    position: fixed;
    inset: 0;
    background: rgba(46, 52, 64, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  /* === Dialog Box === */
  .dialog {
    background: var(--mv-glass-bg);
    backdrop-filter: blur(24px);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border-strong);
    border-radius: var(--mv-radius-lg);
    padding: var(--mv-spacing-xl);
    min-width: 360px;
    max-width: 480px;

    box-shadow:
      0 0 40px rgba(0, 0, 0, 0.4),
      0 0 60px rgba(136, 192, 208, 0.1),
      inset 0 1px 0 var(--mv-glass-border);
  }

  .dialog h4 {
    margin: 0 0 var(--mv-spacing-md);
    font-size: var(--mv-font-size-md);
    font-weight: var(--mv-font-weight-bold);
    color: var(--mv-color-text-primary);
    text-shadow: 0 0 20px rgba(136, 192, 208, 0.3);
  }

  .dialog-item-title {
    margin: 0 0 var(--mv-spacing-lg);
    padding: var(--mv-spacing-sm) var(--mv-spacing-md);
    background: var(--mv-glass-bg-strong);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border);
    border-radius: var(--mv-radius-sm);
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-secondary);
  }

  .dialog label {
    display: block;
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    margin-bottom: var(--mv-spacing-xs);
  }

  .dialog textarea {
    width: 100%;
    padding: var(--mv-spacing-sm) var(--mv-spacing-md);
    background: var(--mv-glass-bg-dark);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border);
    border-radius: var(--mv-radius-sm);
    color: var(--mv-color-text-primary);
    font-family: var(--mv-font-sans);
    font-size: var(--mv-font-size-sm);
    resize: vertical;
    transition: all 0.2s ease;
  }

  .dialog textarea::placeholder {
    color: var(--mv-color-text-disabled);
    font-style: italic;
  }

  .dialog textarea:focus {
    outline: none;
    border-color: var(--mv-primitive-frost-2);
    box-shadow: 0 0 12px rgba(136, 192, 208, 0.2);
  }

  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--mv-spacing-sm);
    margin-top: var(--mv-spacing-lg);
    padding-top: var(--mv-spacing-md);
    border-top: var(--mv-border-width-thin) solid var(--mv-glass-border-subtle);
  }

  .btn-cancel,
  .btn-confirm {
    padding: var(--mv-spacing-xs) var(--mv-spacing-lg);
    border-radius: var(--mv-radius-sm);
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    letter-spacing: 0.05em;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  .btn-cancel {
    background: transparent;
    color: var(--mv-color-text-muted);
    border: var(--mv-border-width-thin) solid var(--mv-glass-border);
  }

  .btn-cancel:hover {
    background: var(--mv-glass-hover);
    color: var(--mv-color-text-primary);
    border-color: var(--mv-glass-border-strong);
  }

  .btn-confirm {
    background: rgba(163, 190, 140, 0.2);
    color: var(--mv-primitive-aurora-green);
    border: var(--mv-border-width-thin) solid rgba(163, 190, 140, 0.5);
  }

  .btn-confirm:hover {
    background: rgba(163, 190, 140, 0.3);
    border-color: var(--mv-primitive-aurora-green);
    box-shadow: 0 0 16px rgba(163, 190, 140, 0.4);
    transform: translateY(-1px);
  }
</style>
