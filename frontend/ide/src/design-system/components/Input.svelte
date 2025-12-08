<script lang="ts">
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  interface Props {
    type?: "text" | "password" | "search" | "email";
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    disabled?: boolean;
    autofocus?: boolean;
    id?: string;
    oninput?: (e: Event & { currentTarget: HTMLInputElement }) => void;
    onchange?: (e: Event & { currentTarget: HTMLInputElement }) => void;
    onkeydown?: (e: KeyboardEvent) => void;
    onsubmit?: () => void;
    onfocus?: (e: FocusEvent) => void;
    onblur?: (e: FocusEvent) => void;
  }

  let {
    type = "text",
    value = $bindable(""),
    placeholder = "",
    label = "",
    error = "",
    disabled = false,
    autofocus = false,
    id = "",
    oninput,
    onchange,
    onkeydown,
    onsubmit,
    onfocus,
    onblur,
  }: Props = $props();

  let inputRef: HTMLInputElement | undefined = $state(undefined);

  onMount(() => {
    if (autofocus && inputRef) {
      inputRef.focus();
    }
  });

  function handleInput(event: Event & { currentTarget: HTMLInputElement }) {
    value = event.currentTarget.value;
    oninput?.(event);
  }

  function handleChange(event: Event & { currentTarget: HTMLInputElement }) {
    onchange?.(event);
  }

  function handleKeydown(event: KeyboardEvent) {
    onkeydown?.(event);
    if (event.key === "Enter") {
      onsubmit?.();
    }
  }
</script>

<div class="input-wrapper" class:has-error={!!error} class:disabled>
  {#if label}
    <label for={id} class="label">{label}</label>
  {/if}

  <div class="input-container">
    <input
      bind:this={inputRef}
      {id}
      {type}
      {value}
      {placeholder}
      {disabled}
      class="input"
      oninput={handleInput}
      onchange={handleChange}
      onkeydown={handleKeydown}
      {onfocus}
      {onblur}
    />
  </div>

  {#if error}
    <p class="error-message" transition:fade={{ duration: 150 }}>{error}</p>
  {/if}
</div>

<style>
  .input-wrapper {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xs);
    width: 100%;
  }

  .label {
    font-size: var(--mv-font-size-sm);
    color: var(--mv-color-text-secondary);
    font-weight: var(--mv-font-weight-medium);
  }

  .input-container {
    position: relative;
    width: 100%;
  }

  .input {
    width: 100%;
    height: var(--mv-input-height-md);
    padding: 0 var(--mv-spacing-md);
    font-family: var(--mv-font-sans);
    font-size: var(--mv-font-size-md);
    color: var(--mv-color-text-primary);

    background: var(--mv-color-surface-secondary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-md);

    transition: var(--mv-transition-base);
  }

  .input:hover:not(:disabled) {
    border-color: var(--mv-color-border-strong);
    background: var(--mv-color-surface-hover);
  }

  .input:focus {
    outline: none;
    border-color: var(--mv-color-interactive-primary);
    background: var(--mv-color-surface-primary);
    box-shadow: var(--mv-shadow-focus);
  }

  /* Disabled State */
  .disabled .input {
    opacity: 0.6;
    cursor: not-allowed;
    background: var(--mv-color-surface-primary);
  }

  /* Error State */
  .has-error .input {
    border-color: var(--mv-color-status-failed-border);
    background: var(--mv-bg-glow-red-lighter);
  }

  .has-error .input:focus {
    border-color: var(--mv-color-status-failed-border);
    box-shadow: var(--mv-shadow-focus-error);
  }

  .has-error .label {
    color: var(--mv-color-status-failed-text);
  }

  .error-message {
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-status-failed-text);
    margin: 0;
    padding-left: var(--mv-spacing-xs);
  }
</style>
