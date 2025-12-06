<script lang="ts">
  import { statusToCssClass } from "../../schemas";
  import type { TaskStatus } from "../../schemas";

  export let status: TaskStatus;

  $: statusClass = statusToCssClass(status);
</script>

<span class="status-badge status-{statusClass}">
  {status}
</span>

<style>
  .status-badge {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-bold);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-wide);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
    display: inline-block;
    min-width: 80px;
    text-align: center;
  }

  .status-pending {
    background: var(--mv-color-status-pending-bg);
    color: var(--mv-color-status-pending-text);
    border: 1px solid var(--mv-color-status-pending-border);
  }

  .status-ready {
    background: var(--mv-color-status-ready-bg);
    color: var(--mv-color-status-ready-text);
    border: 1px solid var(--mv-color-status-ready-border);
    box-shadow: 0 0 4px var(--mv-color-status-ready-text);
  }

  .status-running {
    background: var(--mv-color-status-running-bg);
    color: var(--mv-color-status-running-text);
    border: 1px solid var(--mv-color-status-running-border);
    box-shadow: 0 0 6px var(--mv-color-status-running-glow);
    animation: mv-pulse 2s infinite;
  }

  .status-succeeded {
    background: var(--mv-color-status-succeeded-bg);
    color: var(--mv-color-status-succeeded-text);
    border: 1px solid var(--mv-color-status-succeeded-border);
    box-shadow: 0 0 4px var(--mv-color-status-succeeded-border);
  }

  .status-completed {
    background: var(--mv-color-status-completed-bg);
    color: var(--mv-color-status-completed-text);
    border: 1px solid var(--mv-color-status-completed-border);
    box-shadow: 0 0 8px var(--mv-color-status-completed-glow);
  }

  .status-failed {
    background: var(--mv-color-status-failed-bg);
    color: var(--mv-color-status-failed-text);
    border: 1px solid var(--mv-color-status-failed-border);
    box-shadow: 0 0 4px var(--mv-color-status-failed-border);
  }

  .status-canceled {
    background: var(--mv-color-status-canceled-bg);
    color: var(--mv-color-status-canceled-text);
    border: 1px solid var(--mv-color-status-canceled-border);
  }

  .status-blocked {
    background: var(--mv-color-status-blocked-bg);
    color: var(--mv-color-status-blocked-text);
    border: 1px solid var(--mv-color-status-blocked-border);
  }
</style>
