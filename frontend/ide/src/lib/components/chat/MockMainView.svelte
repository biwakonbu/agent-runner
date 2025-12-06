<script lang="ts">
  import FloatingChatWindow from "./FloatingChatWindow.svelte";
</script>

<div class="main-view">
  <FloatingChatWindow />
</div>

<style>
  .main-view {
    width: 100vw;
    height: 100vh;
    background-color: var(--mv-color-surface-app);
    position: relative;
    overflow: hidden;
    color: var(--mv-color-text-primary);
  }
</style>
