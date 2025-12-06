<script lang="ts">
  import BrandFull from "../components/brand/BrandFull.svelte";
  // Welcome Header - Logo and Title
</script>

<header class="welcome-header">
  <BrandFull size="lg" layout="vertical" />
</header>

<style>
  .welcome-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-bottom: var(--mv-spacing-xl);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-color-border-subtle);
  }
</style>
