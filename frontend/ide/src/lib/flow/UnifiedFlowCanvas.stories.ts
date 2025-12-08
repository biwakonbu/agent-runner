import type { Meta, StoryObj } from '@storybook/svelte';
import UnifiedFlowCanvas from './UnifiedFlowCanvas.svelte';
import type { Task } from '../../types';
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

const mockTasks: Task[] = [
  {
    id: '1',
    title: 'Task 1',
    status: 'COMPLETED',
    phaseName: '概念設計',
    dependencies: [],
    poolId: 'p1',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: '2',
    title: 'Task 2',
    status: 'RUNNING',
    phaseName: '実装設計',
    dependencies: ['1'],
    poolId: 'p1',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: '3',
    title: 'Task 3',
    status: 'PENDING',
    phaseName: '実装',
    dependencies: ['2'],
    poolId: 'p1',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
];

export const Default: Story = {
  args: {
    taskList: mockTasks,
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

export const VariableWidth: Story = {
  args: {
    taskList: [
      {
        id: '1',
        title: 'Short',
        status: 'COMPLETED',
        phaseName: '概念設計',
        dependencies: [],
        poolId: 'p1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: '2',
        title: 'Medium Length Title Task',
        status: 'RUNNING',
        phaseName: '実装設計',
        dependencies: ['1'],
        poolId: 'p1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: '3',
        title: 'Very Long Title Task That Should Be Large Width And Clamp To Three Lines',
        status: 'PENDING',
        phaseName: '実装',
        dependencies: ['2'],
        poolId: 'p1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: '4',
        title: 'Another Short',
        status: 'PENDING',
        phaseName: '検証',
        dependencies: ['2', '3'],
        poolId: 'p2', // Fixed syntax error here
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: '4',
        title: 'Another Short',
        status: 'PENDING',
        phaseName: '検証',
        dependencies: ['2', '3'],
        poolId: 'p2',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      }
    ],
  },
  play: async () => {
    viewMode.setGraph();
  },
};

export const WBSView: Story = {
  args: {
    taskList: mockTasks,
  },
  play: async () => {
    viewMode.setWBS();
  },
};
