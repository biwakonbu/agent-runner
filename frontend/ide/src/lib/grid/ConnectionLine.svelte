<script lang="ts">
  import { grid, gridToCanvas } from "../../design-system";
  import type { TaskEdge } from "../../stores/taskStore";
  import { taskNodes } from "../../stores";

  // Props
  export let edge: TaskEdge;

  // ノード位置のマップを取得
  $: nodePositions = new Map(
    $taskNodes.map((n) => [n.task.id, { col: n.col, row: n.row }])
  );

  // 始点と終点の座標を計算
  $: fromPos = nodePositions.get(edge.from);
  $: toPos = nodePositions.get(edge.to);

  // SVGパスを計算
  $: pathData = calculatePath(fromPos, toPos);

  function calculatePath(
    from: { col: number; row: number } | undefined,
    to: { col: number; row: number } | undefined
  ): string {
    if (!from || !to) return "";

    const fromCanvas = gridToCanvas(from.col, from.row);
    const toCanvas = gridToCanvas(to.col, to.row);

    // ノードの中心から接続点を計算
    // 始点: ノードの右端中央
    const startX = fromCanvas.x + grid.cellWidth;
    const startY = fromCanvas.y + grid.cellHeight / 2;

    // 終点: ノードの左端中央
    const endX = toCanvas.x;
    const endY = toCanvas.y + grid.cellHeight / 2;

    // Smooth Bezier Curve Calculation
    // 水平距離に基づいて制御点を調整
    const dist = Math.abs(endX - startX);

    // 単純な直線に近い場合は制御点を短く、遠い場合は長く
    const controlOffset = Math.max(dist * 0.5, 80);

    // 水平方向が順方向（左から右）かつ十分な距離がある場合
    if (endX > startX + 40) {
      return `M ${startX} ${startY} C ${startX + controlOffset} ${startY}, ${endX - controlOffset} ${endY}, ${endX} ${endY}`;
    } else {
      // 逆方向（右から左）または近すぎる場合の「迂回パス」
      // よりスムーズなS字カーブを描くために制御点を調整

      const loopOffset = 80; // 迂回する幅
      const verticalGap = Math.abs(endY - startY);
      const isBelow = endY > startY;

      // 垂直方向の中間点
      const midY = startY + (endY - startY) / 2;

      // 迂回パスのための垂直オフセット
      // ノードと重ならないように、少し大きめに
      const verticalBypass = isBelow
        ? Math.max(verticalGap, grid.cellHeight + 40)
        : -Math.max(verticalGap, grid.cellHeight + 40);

      // 中間点X (バックトラックの中間)
      const midX = (startX + endX) / 2;

      // 3次ベジェ曲線をつなぎ合わせてスムーズなループを作る
      // 1. Start -> Right Out
      // 2. Right -> Vertical Up/Down
      // 3. Vertical -> Left In

      return `M ${startX} ${startY}
              C ${startX + loopOffset} ${startY},
                ${startX + loopOffset} ${startY + verticalBypass / 2},
                ${midX} ${startY + verticalBypass / 2}
              S ${endX - loopOffset} ${endY},
                ${endX} ${endY}`;
    }
  }

  // 線のスタイルクラス
  $: lineClass = edge.satisfied ? "satisfied" : "unsatisfied";
  $: strokeColor = edge.satisfied
    ? "var(--mv-color-status-succeeded-border)"
    : "var(--mv-color-status-blocked-border)";
</script>

{#if pathData}
  <g class="connection-line {lineClass}">
    <!-- グロー効果用のぼかし背景 -->
    <path
      class="path-glow"
      d={pathData}
      fill="none"
      stroke={strokeColor}
      stroke-width="4"
    />

    <!-- 背景パス（ヒット判定拡大用） -->
    <path
      class="path-hit"
      d={pathData}
      fill="none"
      stroke="transparent"
      stroke-width="20"
    />

    <!-- メインパス -->
    <path
      class="path-main"
      d={pathData}
      fill="none"
      stroke={strokeColor}
      marker-end="url(#arrowhead-{edge.satisfied
        ? 'satisfied'
        : 'unsatisfied'})"
    />

    <!-- data flow animation particle (for unsatisfied/active dependencies) -->
    {#if !edge.satisfied}
      <circle r="3" fill="var(--mv-primitive-aurora-purple)">
        <animateMotion
          dur="2s"
          repeatCount="indefinite"
          path={pathData}
          keyPoints="0;1"
          keyTimes="0;1"
          calcMode="linear"
        />
      </circle>
    {/if}
  </g>
{/if}

<style>
  .connection-line {
    pointer-events: stroke;
    transition: opacity var(--mv-duration-normal);
  }

  .path-glow {
    opacity: 0.4;
    filter: blur(4px);
    transition: stroke var(--mv-duration-normal);
  }

  .path-main {
    stroke-width: 2;
    stroke-linecap: round;
    transition:
      stroke var(--mv-duration-normal),
      stroke-width var(--mv-duration-fast);
  }

  /* 満たされた依存 */
  .satisfied .path-main {
    stroke-opacity: 0.6;
  }

  .satisfied .path-glow {
    opacity: 0.2;
    stroke: var(--mv-color-status-succeeded-text);
  }

  /* 未満の依存（アクティブ） */
  .unsatisfied .path-main {
    stroke-opacity: 0.8;
  }

  .unsatisfied .path-glow {
    opacity: 0.5;
    stroke: var(--mv-color-status-blocked-text);
  }

  /* ホバー時 */
  .connection-line:hover .path-main {
    stroke-width: 3;
    stroke-opacity: 1;
  }

  .connection-line:hover .path-glow {
    opacity: 0.8;
    stroke-width: 6;
    filter: blur(6px);
  }

  .satisfied:hover .path-main {
    stroke: var(--mv-color-status-succeeded-text);
  }

  .unsatisfied:hover .path-main {
    stroke: var(--mv-color-status-blocked-text);
  }
</style>
