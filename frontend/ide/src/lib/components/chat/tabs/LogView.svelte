<script lang="ts">
  import type { ChatLogEntry } from "../../../../stores/chat";
  import { chatLog } from "../../../../stores/chat";
</script>

<div class="log-placeholder">
  <div class="log-entry system">
    <span class="timestamp">{new Date().toLocaleTimeString()}</span>
    <span class="step">[System]</span>
    <span class="message">Log tab selected.</span>
  </div>
  {#each $chatLog as log}
    <div class="log-entry">
      <span class="timestamp"
        >{new Date(log.timestamp).toLocaleTimeString()}</span
      >
      <span class="step">[{log.step}]</span>
      <span class="message">{log.message}</span>
    </div>
  {/each}
</div>

<style>
  .log-placeholder {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xs);
  }

  .log-entry {
    display: flex;
    gap: var(--mv-spacing-xs);
    font-family: var(--mv-font-mono);
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
    padding: var(--mv-spacing-xxs);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-glass-border-subtle);
  }

  .log-entry .timestamp {
    color: var(--mv-color-text-secondary);
    min-width: var(--mv-chat-timestamp-width);
  }

  .log-entry .step {
    color: var(--mv-primitive-frost-2);
    font-weight: var(--mv-font-weight-bold);
    min-width: var(--mv-chat-step-width);
  }

  .log-entry .message {
    color: var(--mv-color-text-primary);
    word-break: break-all;
  }
</style>
