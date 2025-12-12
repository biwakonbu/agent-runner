/**
 * Storybook VRT（Visual Regression Testing）
 *
 * 全 Storybook ストーリーのスクリーンショットを比較し、
 * 意図しない視覚的変更を検知する。
 *
 * 使用方法:
 *   pnpm test:vrt         - スナップショット比較
 *   pnpm test:vrt:update  - スナップショット更新
 */
import { test, expect } from '@playwright/test';
import { getStories, getStoryUrl } from './stories';

// ページ読み込み後の待機時間（アニメーション完了用）
const RENDER_WAIT_TIME = 100;
const STORY_RENDER_TIMEOUT_TIME = 30_000;

// ストーリー一覧を取得して個別テストとして登録
const stories = await getStories();

for (const story of stories) {
  test(`${story.title} / ${story.name}`, async ({ page }) => {
    // ストーリーの iframe URL にアクセス
    const storyUrl = getStoryUrl(story.id);
    await page.goto(storyUrl);

    // Storybook のプレビューが "preparing" 状態のままスクショを撮ると、
    // ローディング画面（白背景 + スピナー）が撮れて VRT が不安定になる。
    await page.waitForFunction(() => {
      const classes = document.body.classList;
      return classes.contains('sb-show-main') || classes.contains('sb-show-errordisplay');
    }, null, { timeout: STORY_RENDER_TIMEOUT_TIME });
    await page.waitForTimeout(RENDER_WAIT_TIME);

    // スクリーンショット比較
    await expect(page).toHaveScreenshot(`${story.id}.png`, {
      fullPage: true,
    });
  });
}
