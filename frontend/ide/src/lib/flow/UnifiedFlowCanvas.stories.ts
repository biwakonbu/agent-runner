import type { Meta, StoryObj } from '@storybook/svelte';
import UnifiedFlowCanvas from './UnifiedFlowCanvas.svelte';
import type { Task, TaskStatus, PhaseName } from '../../types';
import { viewMode } from '../../stores/wbsStore';

const meta = {
  title: 'Flow/UnifiedFlowCanvas',
  component: UnifiedFlowCanvas,
  tags: ['autodocs'],
  argTypes: {
    taskList: { control: 'object' },
  },
} satisfies Meta<UnifiedFlowCanvas>;

export default meta;
type Story = StoryObj<typeof meta>;

// --- Mock Data Generators ---

const validPhases: PhaseName[] = ['概念設計', '実装設計', '実装', '検証'];

function generateShowcaseTasks(count: number): Task[] {
  const tasks: Task[] = [];
  for (let i = 1; i <= count; i++) {
    const phaseIndex = Math.floor((i - 1) / (count / 4)) % 4; // Distribute across phases
    const phase = validPhases[phaseIndex];
    
    // Status variation
    let status: TaskStatus = 'PENDING';
    const mod = i % 10;
    if (mod === 0) status = 'FAILED';
    else if (mod === 1) status = 'BLOCKED';
    else if (mod === 2 || mod === 3) status = 'RUNNING';
    else if (mod === 4 || mod === 5) status = 'COMPLETED';
    else if (mod === 6) status = 'READY';
    else if (mod === 7) status = 'SUCCEEDED';
    else if (mod === 8) status = 'RETRY_WAIT';
    else status = 'PENDING';

    // Dependencies: link to previous task deterministically (no random for VRT stability)
    const dependencies: string[] = [];
    if (i > 1 && i % 3 === 0) {
      dependencies.push(String(i - 1));
    }
    if (i > 5 && i % 7 === 0) {
      dependencies.push(String(i - 5)); // Longer range dep
    }

    // Varying title lengths
    let title = `Task ${i} (${phase})`;
    if (i % 12 === 0) {
      title = `Task ${i} is a very long task title to test text wrappping and clamping behaviors within the task node UI component to ensure it handles overflow correctly.`;
    } else if (i % 7 === 0) {
      title = `Task ${i} - Medium length title description`;
    }

    tasks.push({
      id: String(i),
      title,
      status,
      phaseName: phase,
      dependencies,
      poolId: 'p1',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });
  }
  return tasks;
}

// --- Stories ---

export const Default: Story = {
  args: {
    taskList: generateShowcaseTasks(5),
  },
  play: async () => {
    viewMode.setGraph();
  },
};

export const Showcase: Story = {
  args: {
    taskList: generateShowcaseTasks(40),
  },
  parameters: {
     docs: {
         description: {
             story: '約40個のノードを配置し、様々なステータスやタイトルの長さを確認できるショーケース。'
         }
     }
  },
  play: async () => {
    viewMode.setGraph();
  },
};

export const Empty: Story = {
  args: {
    taskList: [],
  },
  play: async () => {
    viewMode.setGraph();
  },
};

export const WBSView: Story = {
  args: {
    taskList: generateShowcaseTasks(10),
  },
  play: async () => {
    viewMode.setWBS();
  },
};
