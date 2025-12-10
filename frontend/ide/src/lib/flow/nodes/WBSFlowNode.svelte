<script lang="ts">
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";
  import { expandedNodes } from "../../../stores/wbsStore";
  import type { WBSNode } from "../../../stores/wbsStore";
  import { phaseToCssClass } from "../../../schemas";
  import WBSStatusBadge from "../../wbs/WBSStatusBadge.svelte";

  // WBSNode uses fixed dimensions in the old graph, but here we can be flexible or fixed.
  // Using the constants from the old utils if available, or defining local defaults.
  // GRAPH_NODE_WIDTH=220, GRAPH_NODE_HEIGHT=80 from WBSGraphNode context.
  const NODE_WIDTH = 220;
  const NODE_HEIGHT = 80;

  interface Props extends NodeProps {
    data: {
      node: WBSNode;
    };
  }

  let { data }: Props = $props();

  let node = $derived(data.node);
  let expanded = $derived($expandedNodes.has(node.id));
  let phaseClass = $derived(phaseToCssClass(node.phaseName));

  function normalizeStatus(status: string): any {
    return status.toLowerCase();
  }

  function handleGenericClick() {
    if (node.children.length > 0) {
      expandedNodes.toggle(node.id);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && node.children.length > 0) {
      expandedNodes.toggle(node.id);
    }
  }
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class="graph-node {phaseClass}"
  style:width="{NODE_WIDTH}px"
  style:height="{NODE_HEIGHT}px"
  onclick={handleGenericClick}
  onkeydown={handleKeydown}
  role="button"
  tabindex="0"
>
  <!-- Svelte Flow Handles (Left/Right for tree layout) -->
  <Handle type="target" position={Position.Left} class="flow-handle" />
  <Handle type="source" position={Position.Right} class="flow-handle" />

  <div class="phase-bar"></div>
  <div class="node-content">
    <div class="node-title" title={node.label}>{node.label}</div>
    <div class="node-meta">
      {#if node.type === "phase"}
        <span class="phase-badge">{node.label}</span>
      {:else if node.task}
        <WBSStatusBadge status={normalizeStatus(node.task.status)} />
      {/if}
      {#if node.children.length > 0}
        <span class="children-count">
          {expanded ? "▼" : "▶"}
          {node.children.length}
        </span>
      {/if}
    </div>
  </div>
</div>

<style>
  /* Handle Styles (Invisible mostly) */
  :global(.flow-handle) {
    background: transparent !important;
    border: none !important;
    width: var(--mv-space-px) !important;
    height: var(--mv-space-px) !important;
    min-width: var(--mv-space-0) !important;
    min-height: var(--mv-space-0) !important;
    top: var(--mv-space-half, 50%) !important;
  }

  .graph-node {
    position: relative; /* Changed from absolute for Flow */
    background: var(--mv-color-surface-node);
    border: var(--mv-border-width-default) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-md);
    cursor: pointer;
    display: flex;
    overflow: hidden;
    transition:
      border-color var(--mv-transition-hover),
      box-shadow var(--mv-transition-hover),
      transform var(--mv-transition-hover);
    box-shadow: var(--mv-shadow-node-glow);
    box-sizing: border-box;
  }

  .graph-node:hover {
    border-color: var(--mv-color-border-focus);
    transform: translateY(-2px);
    box-shadow: var(--mv-shadow-card);
  }

  .graph-node:focus {
    outline: none;
    border-color: var(--mv-color-border-focus);
    box-shadow: var(--mv-shadow-focus);
  }

  /* フェーズバー */
  .phase-bar {
    width: var(--mv-spacing-xxs);
    flex-shrink: 0;
  }

  /* Phase Colors */
  .phase-concept .phase-bar {
    background: var(--mv-primitive-frost-3);
  }

  .phase-design .phase-bar {
    background: var(--mv-primitive-aurora-purple);
  }

  .phase-impl .phase-bar {
    background: var(--mv-primitive-aurora-green);
  }

  .phase-verify .phase-bar {
    background: var(--mv-primitive-aurora-yellow);
  }

  /* ノードコンテンツ */
  .node-content {
    flex: 1;
    padding: var(--mv-spacing-xs) var(--mv-spacing-sm);
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: var(--mv-spacing-xxs);
  }

  .node-title {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-semibold);
    color: var(--mv-color-text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .node-meta {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-xs);
  }

  .phase-badge {
    font-size: var(--mv-font-size-xs);
    padding: 0 var(--mv-spacing-xxs);
    border-radius: var(--mv-radius-sm);
    text-transform: uppercase;
    letter-spacing: var(--mv-letter-spacing-wide);
    background: var(--mv-color-surface-secondary);
    color: var(--mv-color-text-secondary);
  }

  .children-count {
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
    font-family: var(--mv-font-mono);
  }
</style>
