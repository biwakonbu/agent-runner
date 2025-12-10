import { writable } from 'svelte/store';

export type LayoutDirection = 'TB' | 'LR';

function createLayoutStore() {
  const { subscribe, set, update } = writable<LayoutDirection>('LR'); // Default LR for Time-Flow like WBS

  return {
    subscribe,
    setDirection: (dir: LayoutDirection) => set(dir),
    toggle: () => update(d => d === 'TB' ? 'LR' : 'TB')
  };
}

export const layoutDirection = createLayoutStore();
