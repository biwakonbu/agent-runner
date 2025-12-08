<script lang="ts">
  import {
    SvelteFlow,
    Controls,
    MiniMap,
    type Node,
    type Edge,
    Panel,
  } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import { tasks } from "../../stores/taskStore";
  import { viewMode } from "../../stores/wbsStore"; // Import viewMode
  import type { Task } from "../../types";
  import TaskNode from "./nodes/TaskNode.svelte";
  import DependencyEdge from "./edges/DependencyEdge.svelte";
  import { getLayoutedElements, convertTasksToFlowData } from "./dagreLayout";
  import WBSListView from "../wbs/WBSListView.svelte";

  // Custom Node/Edge Types
  const nodeTypes = {
    task: TaskNode,
  };

  const edgeTypes = {
    dependency: DependencyEdge,
  };

  interface Props {
    taskList?: Task[];
  }

  let { taskList = undefined }: Props = $props();

  // State
  let nodes = $state<Node[]>([]);
  let edges = $state<Edge[]>([]);

  // Update flow data when tasks change
  $effect(() => {
    const targetTasks = taskList ?? $tasks;
    if (targetTasks.length > 0) {
      const { nodes: initialNodes, edges: initialEdges } =
        convertTasksToFlowData(targetTasks);
      const { nodes: layoutedNodes, edges: layoutedEdges } =
        getLayoutedElements(initialNodes, initialEdges);
      nodes = layoutedNodes;
      edges = layoutedEdges;
    } else {
      nodes = [];
      edges = [];
    }
  });

  let isWBSMode = $derived($viewMode === "wbs");
</script>

<div class="flow-container" class:wbs-mode={isWBSMode}>
  <!-- Custom Grid Background (below Svelte Flow) -->
  <div class="grid-background">
    <svg class="grid-pattern" width="100%" height="100%">
      <defs>
        <pattern
          id="grid-cross"
          width="200"
          height="140"
          patternUnits="userSpaceOnUse"
        >
          <path
            d="M96 70H104M100 66V74"
            stroke="var(--mv-primitive-aurora-yellow)"
            stroke-width="1"
            opacity="0.15"
          />
        </pattern>
      </defs>
      <rect width="100%" height="100%" fill="url(#grid-cross)" />
    </svg>
  </div>

  <!-- SVG Markers for edges (must be in DOM for marker-end references) -->
  <svg class="markers-defs" width="0" height="0">
    <defs>
      <!-- Source Port (Hollow Circle) -->
      <marker
        id="marker-source"
        markerWidth="8"
        markerHeight="8"
        refX="4"
        refY="4"
        orient="auto"
        markerUnits="userSpaceOnUse"
      >
        <circle
          cx="4"
          cy="4"
          r="2.5"
          fill="var(--mv-color-surface-app)"
          stroke="var(--mv-color-text-muted)"
          stroke-width="1"
        />
      </marker>

      <!-- Terminal: Satisfied (Solid Square) -->
      <marker
        id="marker-terminal-satisfied"
        markerWidth="10"
        markerHeight="10"
        refX="5"
        refY="5"
        orient="auto"
        markerUnits="userSpaceOnUse"
      >
        <rect
          x="2"
          y="2"
          width="6"
          height="6"
          fill="var(--mv-color-status-succeeded-border)"
        />
      </marker>

      <!-- Terminal: Unsatisfied (Solid Diamond) -->
      <marker
        id="marker-terminal-unsatisfied"
        markerWidth="12"
        markerHeight="12"
        refX="6"
        refY="6"
        orient="auto"
        markerUnits="userSpaceOnUse"
      >
        <path
          d="M6 1 L11 6 L6 11 L1 6 Z"
          fill="var(--mv-color-status-blocked-border)"
        />
      </marker>
    </defs>
  </svg>

  <SvelteFlow
    {nodes}
    {edges}
    {nodeTypes}
    {edgeTypes}
    fitView
    minZoom={0.1}
    maxZoom={4}
    defaultEdgeOptions={{ type: "dependency" }}
    nodesDraggable={!isWBSMode}
    nodesConnectable={!isWBSMode}
    elementsSelectable={!isWBSMode}
  >
    <Controls showZoom={true} />
    <MiniMap />

    <Panel position="top-left" class="wbs-panel">
      <!-- WBSListView is always mounted but hidden via CSS when not in WBS mode -->
      <WBSListView />
    </Panel>
  </SvelteFlow>
</div>

<style>
  .flow-container {
    position: relative;
    width: 100%;
    height: 100%;
    background: var(--mv-color-surface-app);
    overflow: hidden;
  }

  /* Grid background - fixed behind Svelte Flow */
  .grid-background {
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 0;
  }

  .grid-pattern {
    width: 100%;
    height: 100%;
  }

  .markers-defs {
    position: absolute;
    pointer-events: none;
  }

  /* stylelint-disable selector-class-pattern -- Svelte Flow library classes */
  :global(.svelte-flow) {
    position: relative;
    z-index: 1;
    background: transparent !important;
  }

  :global(.svelte-flow__pane) {
    background: transparent !important;
  }

  :global(.svelte-flow__viewport) {
    background: transparent !important;
  }

  :global(.svelte-flow__node) {
    border-radius: var(--mv-radius-md);
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
  }

  :global(.svelte-flow__edge-path) {
    stroke: var(--mv-color-border-default);
    stroke-width: 2;
  }

  :global(.svelte-flow__controls) {
    box-shadow: var(--mv-shadow-card);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    border-radius: var(--mv-radius-sm);
    background: var(--mv-color-surface-primary);
  }

  :global(.svelte-flow__controls-button) {
    background: var(--mv-color-surface-primary);
    border-bottom: var(--mv-border-width-thin) solid
      var(--mv-color-border-subtle);
    fill: var(--mv-color-text-secondary);
  }

  :global(.svelte-flow__controls-button:last-child) {
    border-bottom: none;
  }

  :global(.svelte-flow__controls-button:hover) {
    background: var(--mv-color-surface-hover);
    fill: var(--mv-color-text-primary);
  }

  :global(.svelte-flow__minimap) {
    background: var(--mv-color-surface-primary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    border-radius: var(--mv-radius-md);
  }

  /* Layered Mode Styles */

  /* Graph Layer Control */
  .flow-container.wbs-mode :global(.svelte-flow__viewport),
  .flow-container.wbs-mode :global(.svelte-flow__controls),
  .flow-container.wbs-mode :global(.svelte-flow__minimap) {
    opacity: 0.15;
    filter: blur(2px) grayscale(0.5);
    transition: all 0.3s ease;
    pointer-events: none;
  }

  .flow-container:not(.wbs-mode) :global(.svelte-flow__viewport) {
    opacity: 1;
    filter: none;
    transition: all 0.3s ease;
  }

  /* WBS Panel Control */
  :global(.svelte-flow__panel.wbs-panel) {
    margin: 0;
    width: 100%;
    height: 100%;
    pointer-events: none; /* Default container pass-through */
    transition: all 0.3s ease;
    z-index: 1000;
  }

  /* When WBS is NOT active, hide/fade it */
  .flow-container:not(.wbs-mode) :global(.svelte-flow__panel.wbs-panel) {
    opacity: 0;
    pointer-events: none;
    transform: scale(0.98);
  }

  /* When WBS IS active */
  .flow-container.wbs-mode :global(.svelte-flow__panel.wbs-panel) {
    opacity: 1;
    transform: scale(1);
  }

  /* Ensure WBS content is interactive when active */
  .flow-container.wbs-mode :global(.wbs-panel > *) {
    pointer-events: auto;
  }
</style>
