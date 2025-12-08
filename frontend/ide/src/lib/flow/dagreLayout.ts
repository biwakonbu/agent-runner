import dagre from 'dagre';
import { Position, type Node, type Edge } from '@xyflow/svelte';
import type { Task } from '../../types';

const BASE_HEIGHT = 100; // Safe default height
const WIDTH_SMALL = 180;
const WIDTH_MEDIUM = 240;
const WIDTH_LARGE = 320;

function getNodeDimensions(node: Node) {
  const task = node.data?.task as Task | undefined;
  if (!task) return { width: WIDTH_MEDIUM, height: BASE_HEIGHT };

  const len = task.title.length;
  let width = WIDTH_LARGE;
  if (len <= 15) width = WIDTH_SMALL;
  else if (len <= 30) width = WIDTH_MEDIUM;

  // Rough height estimation
  // Base 80 + ~20 per extra line of text (max 3 lines) if loose
  // But let's keep height semi-constant or safe max for layout stability
  // Just adding a bit of buffer
  return { width, height: BASE_HEIGHT };
}

const LAYOUT_GRID_X = 200; // grid.cellWidth (160) + grid.gap (40)


export const getLayoutedElements = (
  nodes: Node[],
  edges: Edge[],
  direction = 'TB' // 'TB' (top to bottom) or 'LR' (left to right)
) => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  const isHorizontal = direction === 'LR';
  
  // Align rank separation to grid
  // ranksep: vertical gap between layers
  // nodesep: horizontal gap between nodes in same layer
  dagreGraph.setGraph({ 
    rankdir: direction,
    ranksep: 60, // Slightly larger than grid.gap (40) to ensure clear separation, total ~160? No, let's try to hit 140 multiple.
                 // Node height 100. 140 - 100 = 40. So ranksep 40 is ideal for tight packing.
    nodesep: 50  // Horizontal gap.
  });

  nodes.forEach((node) => {
    const { width, height } = getNodeDimensions(node);
    // Artificially inflate width to grid multiples for reliable spacing
    // But this might make it too sparse. Let's trust dagre with real sizes first.
    dagreGraph.setNode(node.id, { width, height });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    const { width, height } = getNodeDimensions(node);
    
    // Snap center to grid
    // We want the center (nodeWithPosition.x, y) to be close to a grid intersection
    // Grid pattern starts at 0,0.
    // X step = 200. Y step = 140.
    
    let snappedX = nodeWithPosition.x;
    let snappedY = nodeWithPosition.y;

    // Apply snapping only if meaningful
    snappedX = Math.round(snappedX / (LAYOUT_GRID_X / 2)) * (LAYOUT_GRID_X / 2); // Snap to half-grid for flexibility?
    // User asked for "Grid Alignment". Let's try strict snapping to 100 (half-cell) or 200.
    // If we snap to 100, we have more freedom.
    snappedX = Math.round(snappedX / 100) * 100;
    snappedY = Math.round(snappedY / 70) * 70; // 140 / 2 = 70.

    // Calculate top-left position based on snapped center
    return {
      ...node,
      targetPosition: isHorizontal ? Position.Left : Position.Top,
      sourcePosition: isHorizontal ? Position.Right : Position.Bottom,
      position: {
        x: snappedX - width / 2,
        y: snappedY - height / 2,
      },
    };
  });

  return { nodes: layoutedNodes, edges };
};

export function convertTasksToFlowData(tasks: Task[]) {
  const nodes: Node[] = [];
  const edges: Edge[] = [];

  tasks.forEach((task) => {
    // Node
    nodes.push({
      id: task.id,
      type: 'task', // Custom node type
      position: { x: 0, y: 0 }, // Initial position, will be calculated by dagre
      data: { task },
    });

    // Edges (Dependencies)
    task.dependencies?.forEach((depId: string) => {
      edges.push({
        id: `e${depId}-${task.id}`,
        source: depId,
        target: task.id,
        type: 'dependency', // Custom edge type
        animated: true,
      });
    });
  });

  return { nodes, edges };
}
