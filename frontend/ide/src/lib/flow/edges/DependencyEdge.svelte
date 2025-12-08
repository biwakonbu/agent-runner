<script lang="ts">
  import { type EdgeProps, getSmoothStepPath } from "@xyflow/svelte";

  interface Props extends EdgeProps {
    data?: {
      satisfied?: boolean;
    };
  }

  let {
    id,
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition,
    targetPosition,
    style = "",
    markerEnd,
    data,
  }: Props = $props();

  let edgePath = $derived(
    getSmoothStepPath({
      sourceX,
      sourceY,
      sourcePosition,
      targetX,
      targetY,
      targetPosition,
      borderRadius: 16, // Smoother corners for the grid-aligned look
    })[0]
  );

  let strokeColor = $derived(
    data?.satisfied
      ? "var(--mv-color-status-succeeded-border)"
      : "var(--mv-color-status-blocked-border)"
  );

  let edgeClass = $derived(data?.satisfied ? "satisfied" : "unsatisfied");
  let markerEndId = $derived(
    `marker-terminal-${data?.satisfied ? "satisfied" : "unsatisfied"}`
  );
</script>

<g class="connection-line {edgeClass}">
  <!-- 背景パス（ヒット判定拡大用） -->
  <path
    class="path-hit"
    d={edgePath}
    fill="none"
    stroke="transparent"
    stroke-width="16"
  />

  <!-- メインパス -->
  <path
    class="path-main"
    d={edgePath}
    fill="none"
    stroke={strokeColor}
    marker-end="url(#{markerEndId})"
    marker-start="url(#marker-source)"
  />

  <!-- シグナルパルス（データフロー） -->
  {#if !data?.satisfied}
    <rect width="4" height="4" fill="var(--mv-primitive-aurora-purple)" rx="1">
      <animateMotion
        dur="1.5s"
        repeatCount="indefinite"
        path={edgePath}
        keyPoints="0;1"
        keyTimes="0;1"
        calcMode="linear"
      />
    </rect>
  {/if}
</g>

<style>
  .connection-line {
    pointer-events: none;
  }

  .path-hit {
    stroke-width: 20;
    cursor: pointer;
    pointer-events: stroke;
  }

  .path-main {
    stroke-width: 1.5;
    transition:
      stroke var(--mv-duration-normal),
      stroke-width var(--mv-duration-fast);
    vector-effect: non-scaling-stroke;
  }

  /* 満たされた依存 */
  .satisfied .path-main {
    stroke-opacity: 0.5;
  }

  /* 未満の依存（アクティブ） */
  .unsatisfied .path-main {
    stroke-opacity: 0.9;
    stroke-dasharray: 2 4;
    stroke-linecap: square;
  }

  /* ホバー時 logic would rely on parent or JS state, 
     Svelte Flow selected state handles some, but hover might need 'interactive' handling
     or CSS on g:hover if pointer-events allowed. */
  .connection-line:hover .path-main {
    stroke-width: 2.5;
    stroke-opacity: 1;
    stroke-dasharray: none;
  }

  .satisfied:hover .path-main {
    stroke: var(--mv-color-status-succeeded-text);
  }

  .unsatisfied:hover .path-main {
    stroke: var(--mv-color-status-blocked-text);
  }
</style>
