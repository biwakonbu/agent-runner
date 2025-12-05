import type { Meta, StoryObj } from '@storybook/svelte';
import FloatingChatWindow from './FloatingChatWindow.svelte';
import MockMainView from './MockMainView.svelte';

const meta = {
  title: 'Features/Chat/FloatingChatWindow',
  component: FloatingChatWindow,
  tags: ['autodocs'],
  argTypes: {},
  parameters: {
      layout: 'centered',
  }
} satisfies Meta<FloatingChatWindow>;

export default meta;
type Story = StoryObj<typeof meta>;

const now = new Date();
const timeMinus2 = new Date(now.getTime() - 2 * 60000).toISOString();
const timeMinus1 = new Date(now.getTime() - 1 * 60000).toISOString();
const timeNow = now.toISOString();

const defaultMessages: Array<{
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
}> = [
  { id: '1', role: 'system', content: 'Agent connected to Multiverse.', timestamp: timeMinus2 },
  { id: '2', role: 'assistant', content: 'こんにちは！タスクを開始します。どのようなサポートが必要ですか？', timestamp: timeMinus1 },
  { id: '3', role: 'user', content: 'フロントエンドのチャットUIを作成してほしい。', timestamp: timeNow },
  { id: '4', role: 'assistant', content: '承知しました。MMO風のフローティングウィンドウを作成しましょう。', timestamp: timeNow },
];

export const Standalone: Story = {
  args: {
    initialPosition: { x: 0, y: 0 }, // Relative to story container
    messages: defaultMessages
  },
  parameters: {
      layout: 'centered', // Show in center to focus on component
  }
};

export const InLayout: Story = {
    render: (args) => ({
        Component: MockMainView,
        props: {
            chatMessages: args.messages
        }
    }),
    args: {
        messages: defaultMessages
    },
    parameters: {
        layout: 'fullscreen'
    }
}
