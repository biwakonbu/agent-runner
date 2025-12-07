<script lang="ts">
  // Wailsのドラッグ領域として機能させる
</script>

<header class="title-bar" style="--wails-draggable:drag">
  <!-- ドラッグ専用領域 (透明) -->
</header>

<style>
  .title-bar {
    height: var(--mv-titlebar-height);
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    user-select: none;
    cursor: default;
    z-index: 9999;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;

    /* Macの信号機ボタン（左上）を避けるためのパディング */
    padding-left: var(--mv-titlebar-padding-left);
  }
</style>
