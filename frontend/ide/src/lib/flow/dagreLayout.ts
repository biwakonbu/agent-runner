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

export const getLayoutedElements = (
  nodes: Node[],
  edges: Edge[],
  direction = 'TB' // 'TB' (top to bottom) or 'LR' (left to right)
) => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  const isHorizontal = direction === 'LR';
  dagreGraph.setGraph({ rankdir: direction });

  nodes.forEach((node) => {
    const { width, height } = getNodeDimensions(node);
    dagreGraph.setNode(node.id, { width, height });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  const layoutedNodes = nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    const { width, height } = getNodeDimensions(node);
    
    // We are shifting the dagre node position (anchor=center center) to the top left
    // so it matches the React Flow node anchor point (top left).
    return {
      ...node,
      targetPosition: isHorizontal ? Position.Left : Position.Top,
      sourcePosition: isHorizontal ? Position.Right : Position.Bottom,
      position: {
        x: nodeWithPosition.x - width / 2,
        y: nodeWithPosition.y - height / 2,
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
