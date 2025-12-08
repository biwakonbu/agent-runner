import type { Meta, StoryObj } from '@storybook/svelte';
import MockMainView from './MockMainView.svelte';

import { windowStore } from '../../../stores/windowStore';

const meta = {
  title: 'Features/Chat/FloatingChatWindow',
  component: MockMainView,
  tags: ['autodocs'],
  argTypes: {},
  parameters: {
    layout: 'fullscreen',
    docs: {
      description: {
        component: 'フローティングチャットウィンドウ。ドラッグ可能なウィンドウ内にチャットUIを表示します。',
      },
    },
  },
  decorators: [
    (Story) => {
        windowStore.update((s: any) => ({ 
             ...s, 
             chat: { ...s.chat, isOpen: true, position: { x: 50, y: 50 } } 
        }));
        return Story();
    }
  ]
} as Meta<typeof MockMainView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  parameters: {
    docs: {
      description: {
        story: 'デフォルト状態のフローティングチャットウィンドウ。',
      },
    },
  },
};
