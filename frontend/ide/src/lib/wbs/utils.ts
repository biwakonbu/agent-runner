/**
 * WBS関連の定数・ユーティリティ
 */

// グラフ描画定数
export const GRAPH_NODE_WIDTH = 200;
export const GRAPH_NODE_HEIGHT = 60;
export const HORIZONTAL_GAP = 80;
export const VERTICAL_GAP = 40;
export const GRAPH_PADDING = 40;

/**
 * Returns the CSS variable for the progress bar color based on percentage.
 * 0% -> Gray
 * 1-29% -> Red
 * 30-49% -> Orange
 * 50-69% -> Yellow
 * 70-89% -> Cyan/Blue
 * 90-99% -> Light Green
 * 100% -> Intense Green
 */
export function getProgressColor(percent: number): {
  fill: string;
  glow: string;
  bg: string;
} {
  if (percent === 0) {
    return {
      fill: 'var(--mv-primitive-neutral-500)', // #5c6677
      glow: 'rgba(92, 102, 119, 0.5)',
      bg: 'rgba(92, 102, 119, 0.2)',
    };
  }
  if (percent < 30) {
    return {
      fill: 'var(--mv-primitive-aurora-red)', // #da7e87 -> 218, 126, 135
      glow: 'rgba(218, 126, 135, 0.4)',
      bg: 'rgba(218, 126, 135, 0.2)',
    };
  }
  if (percent < 50) {
    return {
      fill: 'var(--mv-primitive-aurora-orange)', // #e09880 -> 224, 152, 128
      glow: 'rgba(224, 152, 128, 0.4)',
      bg: 'rgba(224, 152, 128, 0.2)',
    };
  }
  if (percent < 70) {
    return {
      fill: 'var(--mv-primitive-aurora-yellow)', // #f5deaa -> 245, 222, 170
      glow: 'rgba(245, 222, 170, 0.4)',
      bg: 'rgba(245, 222, 170, 0.2)',
    };
  }
  if (percent < 90) {
    return {
      fill: 'var(--mv-progress-bar-fill)', // #b5e8ff -> 181, 232, 255
      glow: 'rgba(181, 232, 255, 0.5)',
      bg: 'rgba(181, 232, 255, 0.2)',
    };
  }
  if (percent < 100) {
    return {
      fill: 'var(--mv-primitive-pastel-green)', // #8fbf9f -> 143, 191, 159
      glow: 'rgba(143, 191, 159, 0.5)',
      bg: 'rgba(143, 191, 159, 0.2)',
    };
  }
  // 100%
  return {
    fill: 'var(--mv-primitive-aurora-green)', // #b8dea6 -> 184, 222, 166
    glow: 'rgba(184, 222, 166, 0.6)',
    bg: 'rgba(184, 222, 166, 0.2)',
  };
}
