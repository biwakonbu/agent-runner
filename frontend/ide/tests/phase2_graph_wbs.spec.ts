import { test, expect } from '@playwright/test';

test.describe('Phase 2: Graph and WBS', () => {
  test.beforeEach(async ({ page }) => {
    // Mock data with dependencies
    await page.addInitScript(() => {
        const tasks = [
            {
                id: 't1',
                title: 'Parent Task',
                status: 'SUCCEEDED',
                poolId: 'default',
                dependencies: [],
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString()
            },
            {
                id: 't2',
                title: 'Child Task',
                status: 'PENDING',
                poolId: 'default',
                dependencies: ['t1'],
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString()
            }
        ];
        window.localStorage.setItem('mock_tasks', JSON.stringify(tasks));
        window.localStorage.removeItem('mock_workspaces');
    });

    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should display tasks in Graph View', async ({ page }) => {
     // 1. Enter Workspace
     await page.getByRole('button', { name: 'Workspaceを開く' }).click();

     // 2. Wait for Toolbar to appear
     await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

     // 3. Verify tasks are visible in Graph View (default)
     await expect(page.getByText('Parent Task')).toBeVisible();
     await expect(page.getByText('Child Task')).toBeVisible();

     // 4. Verify we are in Graph View (canvas-containerがDOMに存在)
     await expect(page.locator('.canvas-container')).toHaveCount(1);
  });

  test('should display hierarchy in WBS View', async ({ page }) => {
      // 1. Enter Workspace
      await page.getByRole('button', { name: 'Workspaceを開く' }).click();

      // 2. Wait for Toolbar
      await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

      // 3. Switch to WBS View (title属性でボタンを検索)
      await page.getByTitle('WBS View').click();

      // 4. Verify WBS View is displayed (WBS List Viewクラスを確認)
      await expect(page.locator('.wbs-list-view')).toBeVisible();

      // 5. Verify WBS tree structure exists (wbs-treeが存在する)
      await expect(page.locator('.wbs-tree')).toBeVisible();

      // 6. WBS View内にノードが存在することを確認
      // (フェーズ/マイルストーンでグループ化されるため、最低1つのノードが存在)
      const nodeCount = await page.locator('.wbs-list-view .wbs-node').count();
      expect(nodeCount).toBeGreaterThanOrEqual(1);
  });

  test('should switch between Graph and WBS views', async ({ page }) => {
      // 1. Enter Workspace
      await page.getByRole('button', { name: 'Workspaceを開く' }).click();

      // 2. Wait for Toolbar
      await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

      // 3. Default is Graph View - verify canvas-container is in DOM
      await expect(page.locator('.canvas-container')).toHaveCount(1);

      // 4. Switch to WBS View (title属性でボタンを検索)
      await page.getByTitle('WBS View').click();
      await expect(page.locator('.wbs-list-view')).toBeVisible();
      // Graph View は消える
      await expect(page.locator('.canvas-container')).toHaveCount(0);

      // 5. Switch back to Graph View
      await page.getByTitle('Graph View').click();
      await expect(page.locator('.canvas-container')).toHaveCount(1);
      // WBS View は消える
      await expect(page.locator('.wbs-list-view')).toHaveCount(0);
  });
});
