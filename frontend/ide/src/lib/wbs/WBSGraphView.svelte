<script lang="ts">
  import { onMount } from "svelte";
  import WBSGraphNode from "./WBSGraphNode.svelte";
  import WBSHeader from "./WBSHeader.svelte";
  import {
    wbsTree,
    expandedNodes,
    overallProgress,
  } from "../../stores/wbsStore";
  import type { WBSNode as WBSNodeType } from "../../stores/wbsStore";
  import {
    GRAPH_NODE_WIDTH as NODE_WIDTH,
    GRAPH_NODE_HEIGHT as NODE_HEIGHT,
    HORIZONTAL_GAP,
    VERTICAL_GAP,
    GRAPH_PADDING as PADDING,
  } from "./utils";

  // ノード位置を計算
  interface PositionedNode {
    node: WBSNodeType;
    x: number;
    y: number;
    level: number;
  }

  // ツリーからグラフレイアウトを生成
  function calculateLayout(
    nodes: WBSNodeType[],
    level: number = 0,
    startY: number = PADDING
  ): PositionedNode[] {
    const result: PositionedNode[] = [];
    let currentY = startY;

    for (const node of nodes) {
      const x = PADDING + level * (NODE_WIDTH + HORIZONTAL_GAP);
      const y = currentY;

      result.push({ node, x, y, level });

      // 子ノードを再帰的にレイアウト
      if (node.children.length > 0 && $expandedNodes.has(node.id)) {
        const childNodes = calculateLayout(node.children, level + 1, currentY);
        result.push(...childNodes);
        // 子ノードの高さ分だけ次の位置を調整
        const childHeight = childNodes.length * (NODE_HEIGHT + VERTICAL_GAP);
        currentY += Math.max(NODE_HEIGHT + VERTICAL_GAP, childHeight);
      } else {
        currentY += NODE_HEIGHT + VERTICAL_GAP;
      }
    }

    return result;
  }

  // 接続線のパスを生成
  function getConnectionPath(from: PositionedNode, to: PositionedNode): string {
    const startX = from.x + NODE_WIDTH;
    const startY = from.y + NODE_HEIGHT / 2;
    const endX = to.x;
    const endY = to.y + NODE_HEIGHT / 2;

    // ベジェ曲線で滑らかな接続
    const controlOffset = HORIZONTAL_GAP / 2;
    return `M ${startX} ${startY} C ${startX + controlOffset} ${startY}, ${endX - controlOffset} ${endY}, ${endX} ${endY}`;
  }

  // 親子関係からエッジを生成
  function getEdges(
    nodes: PositionedNode[]
  ): Array<{ from: PositionedNode; to: PositionedNode }> {
    const edges: Array<{ from: PositionedNode; to: PositionedNode }> = [];
    const nodeMap = new Map(nodes.map((n) => [n.node.id, n]));

    for (const positioned of nodes) {
      if (
        positioned.node.children.length > 0 &&
        $expandedNodes.has(positioned.node.id)
      ) {
        for (const child of positioned.node.children) {
          const childPositioned = nodeMap.get(child.id);
          if (childPositioned) {
            edges.push({ from: positioned, to: childPositioned });
          }
        }
      }
    }

    return edges;
  }

  $: positionedNodes = calculateLayout($wbsTree);
  $: edges = getEdges(positionedNodes);
  $: canvasWidth = Math.max(
    800,
    ...positionedNodes.map((n) => n.x + NODE_WIDTH + PADDING)
  );
  $: canvasHeight = Math.max(
    400,
    ...positionedNodes.map((n) => n.y + NODE_HEIGHT + PADDING)
  );

  // ドラッグスクロール
  let container: HTMLDivElement;
  let isDragging = false;
  let startX = 0;
  let startY = 0;
  let scrollLeft = 0;
  let scrollTop = 0;

  function handleMouseDown(e: MouseEvent) {
    if (e.button !== 0) return;
    isDragging = true;
    startX = e.clientX;
    startY = e.clientY;
    scrollLeft = container.scrollLeft;
    scrollTop = container.scrollTop;
  }

  function handleMouseMove(e: MouseEvent) {
    if (!isDragging) return;
    e.preventDefault();
    const deltaX = e.clientX - startX;
    const deltaY = e.clientY - startY;
    container.scrollLeft = scrollLeft - deltaX;
    container.scrollTop = scrollTop - deltaY;
  }

  function handleMouseUp() {
    isDragging = false;
  }
</script>

<div class="wbs-graph-view">
  <!-- ヘッダー -->
  <WBSHeader />

  <!-- グラフキャンバス -->
  <div
    class="graph-container"
    class:dragging={isDragging}
    bind:this={container}
    on:mousedown={handleMouseDown}
    on:mousemove={handleMouseMove}
    on:mouseup={handleMouseUp}
    on:mouseleave={handleMouseUp}
    role="application"
    aria-label="WBS グラフ"
    tabindex="0"
  >
    {#if positionedNodes.length === 0}
      <div class="empty-state">
        <p>タスクがありません</p>
        <p class="empty-hint">チャットからタスクを生成してください</p>
      </div>
    {:else}
      <div
        class="canvas"
        style:width="{canvasWidth}px"
        style:height="{canvasHeight}px"
      >
        <!-- 接続線 (SVG) -->
        <svg
          class="connections-layer"
          width={canvasWidth}
          height={canvasHeight}
        >
          <defs>
            <marker
              id="arrowhead"
              markerWidth="10"
              markerHeight="7"
              refX="9"
              refY="3.5"
              orient="auto"
            >
              <polygon
                points="0 0, 10 3.5, 0 7"
                fill="var(--mv-color-border-default)"
              />
            </marker>
          </defs>
          {#each edges as edge}
            <path
              class="connection-path"
              d={getConnectionPath(edge.from, edge.to)}
              marker-end="url(#arrowhead)"
            />
          {/each}
        </svg>

        <!-- ノード -->
        <div class="nodes-layer">
          {#each positionedNodes as { node, x, y } (node.id)}
            <WBSGraphNode {node} {x} {y} />
          {/each}
        </div>
      </div>
    {/if}
  </div>

  <!-- 操作ヒント -->
  <div class="controls-hint">
    <span>ドラッグでスクロール</span>
    <span>ノードクリックで展開/折りたたみ</span>
  </div>
</div>

<style>
  .wbs-graph-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--mv-color-surface-node);
  }

  /* グラフコンテナ */
  .graph-container {
    flex: 1;
    overflow: auto;
    position: relative;
    cursor: grab;
    touch-action: none;
  }

  .graph-container.dragging {
    cursor: grabbing;
  }

  /* カスタムスクロールバー */
  .graph-container::-webkit-scrollbar {
    width: var(--mv-scrollbar-width);
    height: var(--mv-scrollbar-width);
  }

  .graph-container::-webkit-scrollbar-track {
    background: var(--mv-color-surface-node);
  }

  .graph-container::-webkit-scrollbar-thumb {
    background: var(--mv-color-border-default);
    border-radius: var(--mv-scrollbar-radius);
  }

  .graph-container::-webkit-scrollbar-thumb:hover {
    background: var(--mv-color-border-strong);
  }

  .graph-container::-webkit-scrollbar-corner {
    background: var(--mv-color-surface-node);
  }

  /* キャンバス */
  .canvas {
    position: relative;
    min-width: var(--mv-size-full);
    min-height: var(--mv-size-full);
  }

  /* 接続線 */
  .connections-layer {
    position: absolute;
    top: 0;
    left: 0;
    pointer-events: none;
  }

  .connection-path {
    fill: none;
    stroke: var(--mv-color-border-default);
    stroke-width: 2;
    filter: drop-shadow(
      0 0 2px var(--mv-color-border-subtle)
    ); /* subtle glow */
    transition: stroke var(--mv-transition-hover);
  }

  /* ノードレイヤー */
  .nodes-layer {
    position: absolute;
    top: 0;
    left: 0;
  }

  /* 空状態 */
  .empty-state {
    position: absolute;
    top: var(--mv-size-half);
    left: var(--mv-size-half);
    transform: translate(-50%, -50%);
    text-align: center;
    color: var(--mv-color-text-muted);
  }

  .empty-state p {
    margin: var(--mv-spacing-xxs) 0;
  }

  .empty-hint {
    font-size: var(--mv-font-size-sm);
  }

  /* 操作ヒント */
  .controls-hint {
    display: flex;
    gap: var(--mv-spacing-sm);
    padding: var(--mv-spacing-xs) var(--mv-spacing-md);
    background: var(--mv-color-surface-secondary);
    border-top: var(--mv-border-width-thin) solid var(--mv-color-border-subtle);
    font-size: var(--mv-font-size-xs);
    color: var(--mv-color-text-muted);
    flex-shrink: 0;
  }

  .controls-hint span {
    background: var(--mv-color-surface-primary);
    padding: var(--mv-spacing-xxs) var(--mv-spacing-xs);
    border-radius: var(--mv-radius-sm);
  }
</style>
