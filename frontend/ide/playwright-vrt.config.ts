import { defineConfig, devices } from '@playwright/test';

/**
 * Storybook VRT（Visual Regression Testing）専用設定
 *
 * 使用方法:
 *   pnpm test:vrt         - スナップショット比較
 *   pnpm test:vrt:update  - スナップショット更新
 */
export default defineConfig({
  testDir: './tests/vrt',
  fullyParallel: true, // 並列実行で高速化
  forbidOnly: !!process.env.CI,
  retries: 0,
  workers: 4, // ワーカー数増加
  reporter: [['html', { outputFolder: 'playwright-vrt-report' }]],

  use: {
    baseURL: 'http://localhost:6006',
    trace: 'off', // トレース無効化で高速化
  },

  expect: {
    toHaveScreenshot: {
      maxDiffPixels: 250, // 許容差分ピクセル数（フレーキーテスト対策）
      threshold: 0.2, // 許容差分率
      animations: 'disabled', // アニメーション無効化
    },
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  webServer: {
    command: 'pnpm run storybook --ci',
    port: 6006,
    reuseExistingServer: true,
    timeout: 120 * 1000,
  },
});
