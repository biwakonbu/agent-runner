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
  insetShadow: string;
  textShadowMd: string;
  textShadowXs: string;
} {
  // 共通の inset shadow（CSS変数は直接使えないためハードコード）
  const insetShadowDark = 'inset 0 1px 2px rgba(0, 0, 0, 0.3)';
  const insetShadowLight = 'inset 0 1px 2px rgba(0, 0, 0, 0.2)';

  // text-shadow のテンプレート関数
  const makeTextShadowMd = (glow: string) => `0 0 8px ${glow}`;
  const makeTextShadowXs = (glow: string) => `0 0 1px ${glow}`;

  if (percent === 0) {
    const glow = 'rgba(92, 102, 119, 0.5)';
    return {
      fill: 'var(--mv-primitive-neutral-500)',
      glow,
      bg: 'rgba(92, 102, 119, 0.2)',
      insetShadow: insetShadowDark,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }
  if (percent < 30) {
    const glow = 'rgba(218, 126, 135, 0.4)';
    return {
      fill: 'var(--mv-primitive-aurora-red)',
      glow,
      bg: 'rgba(218, 126, 135, 0.2)',
      insetShadow: insetShadowDark,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }
  if (percent < 50) {
    const glow = 'rgba(224, 152, 128, 0.4)';
    return {
      fill: 'var(--mv-primitive-aurora-orange)',
      glow,
      bg: 'rgba(224, 152, 128, 0.2)',
      insetShadow: insetShadowDark,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }
  if (percent < 70) {
    const glow = 'rgba(245, 222, 170, 0.4)';
    return {
      fill: 'var(--mv-primitive-aurora-yellow)',
      glow,
      bg: 'rgba(245, 222, 170, 0.2)',
      insetShadow: insetShadowDark,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }
  if (percent < 90) {
    const glow = 'rgba(181, 232, 255, 0.5)';
    return {
      fill: 'var(--mv-progress-bar-fill)',
      glow,
      bg: 'rgba(181, 232, 255, 0.2)',
      insetShadow: insetShadowLight,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }
  if (percent < 100) {
    const glow = 'rgba(143, 191, 159, 0.5)';
    return {
      fill: 'var(--mv-primitive-pastel-green)',
      glow,
      bg: 'rgba(143, 191, 159, 0.2)',
      insetShadow: insetShadowLight,
      textShadowMd: makeTextShadowMd(glow),
      textShadowXs: makeTextShadowXs(glow),
    };
  }

  // 100%
  const glow = 'rgba(184, 222, 166, 0.6)';
  return {
    fill: 'var(--mv-primitive-aurora-green)',
    glow,
    bg: 'rgba(184, 222, 166, 0.2)',
    insetShadow: insetShadowLight,
    textShadowMd: makeTextShadowMd(glow),
    textShadowXs: makeTextShadowXs(glow),
  };
}
