<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let disabled = false;

  const dispatch = createEventDispatcher<{ send: string }>();

  let value = "";

  function handleSend() {
    if (value.trim() && !disabled) {
      dispatch("send", value);
      value = "";
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }
</script>

<div class="chat-input-container">
  <span class="prompt-icon">&gt;</span>
  <div class="input-wrapper">
    <textarea
      bind:value
      placeholder="何か話す..."
      class="transparent-input"
      class:disabled
      on:keydown={handleKeydown}
      rows="2"
      {disabled}
    ></textarea>
  </div>
</div>

<style>
  .chat-input-container {
    display: flex;
    align-items: flex-start; /* Align to top for multi-line */
    gap: var(--mv-spacing-xs);
    padding: var(--mv-spacing-xs);
    background: var(--mv-color-surface-overlay);
    border-top: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
  }

  .prompt-icon {
    color: var(--mv-primitive-frost-1); /* User color */
    font-weight: bold;
    font-family: var(--mv-font-mono);
    margin-top: var(--mv-spacing-xxs); /* Align with first line of text */
  }

  .input-wrapper {
    flex: 1;
  }

  .transparent-input {
    width: 100%;
    background: transparent;
    border: none;
    color: var(--mv-color-text-primary);
    font-family: var(--mv-font-sans);
    font-size: var(--mv-font-size-md);
    outline: none;
    text-shadow: var(--mv-border-width-thin) var(--mv-border-width-thin) var(--mv-border-width-thin) var(--mv-primitive-deep-0);
    resize: none; /* User can't resize manually, fixed to rows */
    display: block;
    line-height: var(--mv-line-height-normal);
  }

  .transparent-input::placeholder {
    color: var(--mv-color-text-disabled);
    opacity: 0.6;
  }

  .transparent-input.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
