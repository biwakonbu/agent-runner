<script lang="ts">
  import { Handle, Position, type NodeProps } from "@xyflow/svelte";
  import type { WBSNode } from "../../../stores/wbsStore";

  interface Props extends NodeProps {
    data: {
      node: WBSNode;
    };
  }

  let { data }: Props = $props();
  let node = $derived(data.node);

  // Calculate complete percentage
  let percentage = $derived(node?.progress?.percentage ?? 0);
  let isCompleted = $derived(percentage === 100);

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      // Logic for selecting milestone could go here
    }
  }
</script>

<div
  class="milestone-node"
  class:completed={isCompleted}
  role="button"
  tabindex="0"
  onkeydown={handleKeydown}
>
  <Handle type="target" position={Position.Left} class="flow-handle" />
  <Handle type="source" position={Position.Right} class="flow-handle" />

  <div class="milestone-shape">
    <div class="milestone-inner">
      <div class="milestone-icon">◈</div>
    </div>
  </div>
  <div class="milestone-content">
    <div class="milestone-label">{node.label}</div>
    <div class="milestone-progress">
      <div class="progress-bar">
        <div class="progress-fill" style:width="{percentage}%"></div>
      </div>
      <span class="progress-text">{percentage}%</span>
    </div>
  </div>
</div>

<style>
  :global(.flow-handle) {
    background: transparent !important;
    border: none !important;
    width: var(--mv-space-px) !important;
    height: var(--mv-space-px) !important;
    min-width: var(--mv-space-0) !important;
    min-height: var(--mv-space-0) !important;
    top: var(--mv-space-half, 50%) !important;
  }

  .milestone-node {
    display: flex;
    align-items: center;
    gap: var(--mv-spacing-sm);
    padding: var(--mv-spacing-xs);
    background: var(--mv-glass-bg-card);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-lg);
    min-width: var(--mv-space-200, 200px);
    transition:
      transform var(--mv-duration-fast),
      box-shadow var(--mv-duration-fast);
  }

  .milestone-node:hover {
    transform: translateY(-2px);
    box-shadow: var(--mv-shadow-card);
    border-color: var(--mv-primitive-frost-2);
  }

  .milestone-node.completed {
    border-color: var(--mv-primitive-aurora-green);
    box-shadow: var(--mv-shadow-badge-glow-sm) var(--mv-primitive-aurora-green);
  }

  .milestone-shape {
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--mv-space-8);
    height: var(--mv-space-8);
    background: var(--mv-color-surface-secondary);
    border-radius: var(--mv-radius-sm);
    flex-shrink: 0;
    transform: rotate(45deg);
    margin: var(--mv-space-1);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
  }

  .milestone-node.completed .milestone-shape {
    background: var(--mv-primitive-aurora-green);
    border-color: var(--mv-primitive-aurora-green);
    color: var(--mv-primitive-snow-storm-1);
  }

  .milestone-inner {
    transform: rotate(-45deg); /* Counter rotate icon */
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .milestone-icon {
    font-size: var(--mv-font-size-md);
    color: var(--mv-color-text-primary);
  }

  .milestone-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--mv-space-1);
  }

  .milestone-label {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-bold);
    color: var(--mv-color-text-primary);
  }

  .milestone-progress {
    display: flex;
    align-items: center;
    gap: var(--mv-space-2);
  }

  .progress-bar {
    flex: 1;
    height: var(--mv-space-1);
    background: var(--mv-color-surface-tertiary);
    border-radius: var(--mv-radius-xs, 2px);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: var(--mv-primitive-frost-3);
    border-radius: var(--mv-radius-xs, 2px);
    transition: width 0.3s ease;
  }

  .milestone-node.completed .progress-fill {
    background: var(--mv-primitive-aurora-green);
  }

  .progress-text {
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
    font-family: var(--mv-font-mono);
  }
</style>
