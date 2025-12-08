<script lang="ts">
  /**
   * Flex Component
   * Flexbox container helper.
   */
  import Box from "./Box.svelte";

  interface Props {
    direction?: "row" | "column" | "row-reverse" | "column-reverse";
    align?: "start" | "center" | "end" | "stretch" | "baseline";
    justify?: "start" | "center" | "end" | "between" | "around" | "evenly";
    wrap?: "nowrap" | "wrap" | "wrap-reverse";
    gap?: string;
    // Box props
    p?: string;
    m?: string;
    bg?: string | undefined;
    height?: string;
    width?: string;
    grow?: boolean;
    class?: string;
    children?: import("svelte").Snippet;
    onclick?: (event: MouseEvent) => void;
    onkeydown?: (event: KeyboardEvent) => void;
  }

  let {
    direction = "row",
    align = "stretch",
    justify = "start",
    wrap = "nowrap",
    gap = "0",
    p = "0",
    m = "0",
    bg = undefined,
    height = "auto",
    width = "auto",
    grow = false,
    class: className = "",
    children,
    onclick,
    onkeydown,
  }: Props = $props();

  const justifyMap = {
    start: "flex-start",
    center: "center",
    end: "flex-end",
    between: "space-between",
    around: "space-around",
    evenly: "space-evenly",
  };

  const alignMap = {
    start: "flex-start",
    center: "center",
    end: "flex-end",
    stretch: "stretch",
    baseline: "baseline",
  };
</script>

{#if onclick || onkeydown}
  <div
    class="flex {className}"
    class:grow
    style:--flex-dir={direction}
    style:--flex-align={alignMap[align]}
    style:--flex-justify={justifyMap[justify]}
    style:--flex-wrap={wrap}
    style:--flex-gap={gap}
    style:--box-p={p}
    style:--box-m={m}
    style:--box-bg={bg}
    style:--box-h={height}
    style:--box-w={width}
    {onclick}
    {onkeydown}
    role="button"
    tabindex="0"
  >
    {@render children?.()}
  </div>
{:else}
  <div
    class="flex {className}"
    class:grow
    style:--flex-dir={direction}
    style:--flex-align={alignMap[align]}
    style:--flex-justify={justifyMap[justify]}
    style:--flex-wrap={wrap}
    style:--flex-gap={gap}
    style:--box-p={p}
    style:--box-m={m}
    style:--box-bg={bg}
    style:--box-h={height}
    style:--box-w={width}
  >
    {@render children?.()}
  </div>
{/if}

<style>
  .flex {
    display: flex;
    box-sizing: border-box;
    flex-flow: var(--flex-dir) var(--flex-wrap);
    align-items: var(--flex-align);
    justify-content: var(--flex-justify);
    gap: var(--flex-gap);

    padding: var(--box-p);
    margin: var(--box-m);
    background: var(--box-bg, transparent);
    height: var(--box-h);
    width: var(--box-w);
  }

  .grow {
    flex-grow: 1;
  }
</style>
