import type { Meta, StoryObj } from '@storybook/svelte';
import ToastContainer from './ToastContainer.svelte';
import { toasts } from '../../stores/toastStore';

const meta = {
  title: 'Components/ToastContainer',
  component: ToastContainer,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
} satisfies Meta<ToastContainer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  play: async () => {
    // Reset toasts
    toasts.reset();
    
    // Add demo toasts
    toasts.add('This is an info toast', 'info');
    toasts.add('This is a success toast', 'success');
    toasts.add('This is a warning toast', 'warning');
    toasts.add('This is an error toast', 'error');
  }
};

export const WithAction: Story = {
    play: async () => {
        toasts.reset();
        toasts.add('Toast with action', 'info', 3000, {
                label: 'Undo',
                onClick: () => console.log('Undo clicked')
        });
    }
}
