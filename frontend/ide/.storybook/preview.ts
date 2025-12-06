import type { Preview } from '@storybook/svelte';

// デザインシステムのCSS変数をインポート
import '../src/design-system/variables.css';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i
      }
    },
    backgrounds: {
      default: 'multiverse-app',
      values: [
        { name: 'multiverse-app', value: '#16181e' }, // --mv-primitive-deep-0
        { name: 'dark', value: '#1a1a1a' },
        { name: 'light', value: '#ffffff' }
      ]
    }
  }
};

export default preview;
