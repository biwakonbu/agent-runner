<script lang="ts">
  export let size: "sm" | "md" | "lg" | "xl" = "md";

  const fontSizes = {
    sm: "1rem", // 16px
    md: "1.5rem", // 24px
    lg: "3rem", // 48px
    xl: "5rem", // 80px
  };
</script>

<div class="brand-text {size}" style="--font-size: {fontSizes[size]}">
  <span class="text">MULTIVERSE</span>
</div>

<style>
  @import url("https://fonts.googleapis.com/css2?family=Orbitron:wght@600;800&display=swap");

  .brand-text {
    font-family: var(--mv-brand-font-family);
    font-size: var(--font-size);
    font-weight: var(--mv-brand-font-weight);
    line-height: 1;
    letter-spacing: var(--mv-brand-letter-spacing);
    user-select: none;
    display: inline-flex;
    text-transform: uppercase;
  }

  .text {
    background: var(--mv-brand-gradient);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
    filter: drop-shadow(var(--mv-brand-shadow));
    position: relative;
    z-index: 1;
    transition: filter 0.3s ease;
  }

  .brand-text:hover .text {
    filter: drop-shadow(var(--mv-brand-shadow-hover));
  }
</style>
