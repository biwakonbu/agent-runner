<script lang="ts">
  import { slide } from "svelte/transition";
  import { Brain, Server, Box } from "lucide-svelte";
  import type { ResourceNode, ResourceType, ResourceStatus } from "./types";

  interface Props {
    resources?: ResourceNode[];
  }

  let { resources = $bindable([]) }: Props = $props();

  // Flatten the tree for rendering
  function flatten(
    nodes: ResourceNode[],
    depth = 0
  ): Array<ResourceNode & { depth: number }> {
    let result: Array<ResourceNode & { depth: number }> = [];
    for (const node of nodes) {
      result.push({ ...node, depth });
      if (node.children && node.expanded !== false) {
        // Default to expanded if undefined
        result = result.concat(flatten(node.children, depth + 1));
      }
    }
    return result;
  }

  let flatNodes = $derived(flatten(resources));

  function toggleExpand(nodeId: string) {
    resources = toggleNode(resources, nodeId);
  }

  function toggleNode(nodes: ResourceNode[], id: string): ResourceNode[] {
    return nodes.map((node) => {
      if (node.id === id) {
        return { ...node, expanded: node.expanded === false ? true : false };
      }
      if (node.children) {
        return { ...node, children: toggleNode(node.children, id) };
      }
      return node;
    });
  }

  function getStatusColor(status: ResourceStatus): string {
    switch (status) {
      case "RUNNING":
        return "var(--mv-color-status-success-text)"; // Aurora Green
      case "THINKING":
        return "var(--mv-color-status-info-text)"; // Frost Blue
      case "ERROR":
        return "var(--mv-color-status-failed-text)"; // Aurora Red
      case "PAUSED":
        return "var(--mv-color-status-warning-text)"; // Aurora Yellow
      case "DONE":
        return "var(--mv-color-text-muted)";
      default:
        return "var(--mv-color-text-muted)";
    }
  }

  function getTypeBadgeStyle(type: ResourceType): string {
    switch (type) {
      case "META": // Purple
        return "background: rgba(180, 142, 173, 0.2); color: #B48EAD; border: 1px solid rgba(180, 142, 173, 0.1);";
      case "WORKER": // Green
        return "background: rgba(163, 190, 140, 0.2); color: #A3BE8C; border: 1px solid rgba(163, 190, 140, 0.1);";
      case "CONTAINER": // Orange/Gold
        return "background: rgba(235, 203, 139, 0.2); color: #EBCB8B; border: 1px solid rgba(235, 203, 139, 0.1);";
      default:
        return "background: rgba(255,255,255,0.05); color: #999;";
    }
  }

  // Parse utility for metrics
  function parseMetrics(
    detail?: string
  ): { label: string; percent: number; color: string }[] {
    if (!detail) return [];
    const metrics: { label: string; percent: number; color: string }[] = [];

    // Match "CPU: 12%" or "Mem: 50%" patterns
    const cpuMatch = detail.match(/CPU:\s*(\d+)%/i);
    if (cpuMatch) {
      const val = parseInt(cpuMatch[1], 10);
      metrics.push({
        label: "CPU",
        percent: val,
        color:
          val > 80
            ? "var(--mv-primitive-aurora-red)"
            : "var(--mv-primitive-frost-3)",
      });
    }

    // Additional parsing for logic can go here
    return metrics;
  }
</script>

<div class="resource-list">
  <div class="header-row">
    <div class="col-name">Resource</div>
    <div class="col-type">Type</div>
    <div class="col-status">Status</div>
    <div class="col-activity">Activity / Monitor</div>
  </div>

  <div class="list-body">
    {#each flatNodes as node (node.id)}
      {@const metrics = parseMetrics(node.detail)}
      {@const icon =
        node.type === "META" ? Brain : node.type === "WORKER" ? Server : Box}
      <div
        class="resource-row"
        class:status-running={node.status === "RUNNING"}
        class:status-error={node.status === "ERROR"}
        onclick={() => toggleExpand(node.id)}
        role="button"
        tabindex="0"
        onkeydown={(e) => e.key === "Enter" && toggleExpand(node.id)}
      >
        <!-- Name Column with Indent -->
        <div class="col-name" style:--depth={node.depth}>
          {#if node.children && node.children.length > 0}
            <span class="disclosure-icon"
              >{node.expanded !== false ? "▼" : "▶"}</span
            >
          {:else}
            <span class="disclosure-placeholder"></span>
          {/if}
          <span class="node-name">{node.name}</span>
        </div>

        <!-- Type Column -->
        <div class="col-type">
          <div class="type-badge-pill" style={getTypeBadgeStyle(node.type)}>
            {#if node.type === "META"}
              <Brain size={10} strokeWidth={3} />
            {:else if node.type === "WORKER"}
              <Server size={10} strokeWidth={3} />
            {:else}
              <Box size={10} strokeWidth={3} />
            {/if}
            <span>{node.type}</span>
          </div>
        </div>

        <!-- Status Column -->
        <div class="col-status">
          <div class="status-indicator-wrapper">
            <div
              class="status-dot"
              style:background={getStatusColor(node.status)}
            >
              {#if node.status === "RUNNING" || node.status === "THINKING"}
                <div
                  class="status-pulse"
                  style:border-color={getStatusColor(node.status)}
                ></div>
              {/if}
            </div>
            <span class="status-label" style:color={getStatusColor(node.status)}
              >{node.status}</span
            >
          </div>
        </div>

        <!-- Activity/Monitor Column -->
        <div class="col-activity">
          {#if metrics.length > 0}
            <div class="metrics-grid">
              {#each metrics as metric}
                <div class="metric-item">
                  <span class="metric-label">{metric.label}</span>
                  <div class="metric-bar-bg">
                    <div
                      class="metric-bar-fill"
                      style:width="{metric.percent}%"
                      style:background-color={metric.color}
                    ></div>
                  </div>
                  <span class="metric-val">{metric.percent}%</span>
                </div>
              {/each}
              {#if node.detail && !node.detail.includes("CPU:")}
                <span class="detail-text extra">{node.detail}</span>
              {/if}
            </div>
          {:else}
            <span class="detail-text">{node.detail || "-"}</span>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .resource-list {
    display: flex;
    flex-direction: column;
    width: 100%;
    /* Remove self-contained glass style since it lives in a glass window now */
    background: transparent;
    border: none;
    border-radius: 0;
    font-family: var(--mv-font-sans);
    overflow: hidden;
  }

  .header-row {
    display: grid;
    grid-template-columns: 2.5fr 120px 120px 4fr;
    padding: var(--mv-space-2) var(--mv-space-3);
    background: rgba(255, 255, 255, 0.03);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    font-size: var(--mv-font-size-xs);
    font-weight: 600;
    color: var(--mv-color-text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    align-items: center;
  }

  /* Header alignment specific */
  .header-row > .col-type,
  .header-row > .col-status {
    display: flex;
    justify-content: center;
  }

  .list-body {
    max-height: var(--mv-space-96);
    overflow-y: auto;
  }

  .resource-row {
    display: grid;
    grid-template-columns: 2.5fr 120px 120px 4fr;
    padding: var(--mv-space-2) 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.02);
    align-items: center;
    cursor: pointer;
    transition: background 0.1s ease;
    font-size: var(--mv-font-size-sm);
  }

  .resource-row:hover {
    background: rgba(255, 255, 255, 0.04);
  }

  .resource-row.status-running {
    background: linear-gradient(
      90deg,
      rgba(163, 190, 140, 0.05) 0%,
      transparent 100%
    );
  }

  .resource-row.status-error {
    background: linear-gradient(
      90deg,
      rgba(191, 97, 106, 0.05) 0%,
      transparent 100%
    );
  }

  .resource-row:last-child {
    border-bottom: none;
  }

  .col-name {
    display: flex;
    align-items: center;
    gap: var(--mv-space-1-5);
    color: var(--mv-color-text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    padding-left: calc(var(--depth, 0) * var(--mv-space-5) + var(--mv-space-3));
  }

  .disclosure-icon,
  .disclosure-placeholder {
    width: var(--mv-space-3);
    font-size: var(--mv-font-size-2xs);
    color: var(--mv-color-text-muted);
    text-align: center;
  }

  .node-name {
    font-weight: var(--mv-font-weight-medium);
  }

  .col-type {
    display: flex;
    align-items: center;
    justify-content: center; /* Center content in body */
  }

  .type-badge-pill {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 10px;
    padding: 2px 8px;
    border-radius: 99px; /* Pill shape */
    font-weight: 700;
    letter-spacing: 0.05em;
    min-width: 80px;
    justify-content: center;
  }

  .col-status {
    display: flex;
    align-items: center;
    justify-content: center; /* Center content in body */
  }

  .status-indicator-wrapper {
    display: flex;
    align-items: center;
    gap: var(--mv-space-2);
  }

  .status-dot {
    width: var(--mv-space-1-5);
    height: var(--mv-space-1-5);
    border-radius: var(--mv-radius-full);
    position: relative;
  }

  .status-pulse {
    position: absolute;
    top: calc(-1 * var(--mv-space-0-75));
    left: calc(-1 * var(--mv-space-0-75));
    width: var(--mv-space-3);
    height: var(--mv-space-3);
    border-radius: var(--mv-radius-full);
    border: var(--mv-border-width-sm) solid;
    opacity: 0;
    animation: pulse var(--mv-duration-slow) infinite;
  }

  .status-label {
    font-size: var(--mv-font-size-xs);
    font-weight: var(--mv-font-weight-semibold);
  }

  .col-activity {
    color: var(--mv-color-text-secondary);
    font-family: var(--mv-font-mono);
    font-size: var(--mv-font-size-xs);
    /* opacity: var(--mv-opacity-80); Removed to fix contrast on charts */
    padding-right: var(--mv-space-3);
    display: flex;
    align-items: center;
  }

  /* Metrics Grid */
  .metrics-grid {
    display: flex;
    align-items: center;
    gap: 16px;
    width: 100%;
  }

  .metric-item {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 120px;
  }

  .metric-label {
    font-size: 10px;
    color: var(--mv-color-text-muted);
    width: 24px;
  }

  .metric-bar-bg {
    flex: 1;
    height: 4px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 2px;
    overflow: hidden;
  }

  .metric-bar-fill {
    height: 100%;
    border-radius: 2px;
    transition: width 0.3s ease;
  }

  .metric-val {
    font-size: 10px;
    width: 30px;
    text-align: right;
  }

  .detail-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    opacity: 0.8;
  }

  .detail-text.extra {
    margin-left: 8px;
    font-size: 10px;
    opacity: 0.6;
  }

  @keyframes pulse {
    0% {
      transform: scale(var(--mv-scale-half));
      opacity: 0;
    }
    50% {
      opacity: var(--mv-opacity-60);
    }
    100% {
      transform: scale(var(--mv-scale-150));
      opacity: 0;
    }
  }
</style>
