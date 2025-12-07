<script lang="ts">
  // No logic needed, just static presentation
</script>

<div class="empty-state">
  <span class="empty-icon">&#10003;</span>
  <p>バックログは空です</p>
</div>

<style>
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
    font-size: var(--mv-font-size-xxxl);
    margin-bottom: var(--mv-spacing-md);
    color: var(--mv-primitive-aurora-green);

    /* Success glow */
    filter: var(--mv-filter-success-glow);
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
    letter-spacing: var(--mv-letter-spacing-count);
  }
</style>
