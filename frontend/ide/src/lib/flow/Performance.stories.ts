import type { Meta, StoryObj } from "@storybook/svelte";
import UnifiedFlowCanvas from "./UnifiedFlowCanvas.svelte";
import type { Task } from "../../types";

const meta = {
  title: "Flow/Performance",
  component: UnifiedFlowCanvas,
  tags: ["skip-vrt"],
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta<UnifiedFlowCanvas>;

export default meta;
type Story = StoryObj<UnifiedFlowCanvas>;

// Generate 2000 nodes
const generateManyTasks = (count: number): Task[] => {
  const tasks: Task[] = [];
  for (let i = 0; i < count; i++) {
    const id = `task-${i}`;
    const parentId = i > 0 ? `task-${Math.floor((i - 1) / 2)}` : undefined; // Binary tree structure
    tasks.push({
      id,
      title: `Task ${i} - Performance Test Node`,
      status: i % 2 === 0 ? "COMPLETED" : "PENDING",
      description: `Description for task ${i}`,
      poolId: "default",
      phaseName: "実装",
      wbsLevel: 1,
      dependencies: parentId ? [parentId] : [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });
  }
  return tasks;
};

export const TwoThousandNodes: Story = {
  args: {
    taskList: generateManyTasks(2000),
  },
};
