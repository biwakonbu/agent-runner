import { test, expect } from '@playwright/test';

test.describe('Task Integration & UI Updates', () => {
  test.beforeEach(async ({ page }) => {
    // Clear mock storage
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    // Setup console log forwarding for debugging
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should open workspace and display main UI', async ({ page }) => {
    // 1. Open Workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // 2. Verify main UI elements are visible
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
    await expect(page.getByTitle('Graph View')).toBeVisible();
    await expect(page.getByTitle('WBS View')).toBeVisible();

    // 3. Verify chat window is visible (placeholderでtextareaを検索)
    await expect(page.locator('textarea[placeholder="Ask Multiverse..."]')).toBeVisible();
  });

  test('should display tasks with dependencies in Graph view', async ({ page }) => {
    // Setup tasks with dependencies
    await page.addInitScript(() => {
      window.localStorage.setItem('mock_tasks', JSON.stringify([
        {
          id: '1',
          title: 'Task A',
          status: 'SUCCEEDED',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        },
        {
          id: '2',
          title: 'Task B',
          status: 'PENDING',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          dependencies: ['1']
        }
      ]));
    });

    // Reload to pick up mock data
    await page.reload();
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // Wait for tasks to load
    await expect(page.getByText('Task A')).toBeVisible();
    await expect(page.getByText('Task B')).toBeVisible();
  });

  test('should switch to WBS view and display tasks', async ({ page }) => {
    // Setup tasks
    await page.addInitScript(() => {
      window.localStorage.setItem('mock_tasks', JSON.stringify([
        {
          id: '1',
          title: 'WBS Task',
          status: 'PENDING',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]));
    });

    await page.reload();
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

    // Switch to WBS View (title属性でボタンを検索)
    await page.getByTitle('WBS View').click();

    // Verify WBS view is displayed (WBS List Viewクラスを確認)
    await expect(page.locator('.wbs-list-view')).toBeVisible();

    // Verify WBS tree exists and has at least one node
    await expect(page.locator('.wbs-list-view .wbs-node')).toHaveCount(1, { timeout: 10000 });
  });

  test('should show execution state in toolbar', async ({ page }) => {
    // Open workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // Verify initial execution state
    await expect(page.getByText('IDLE')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should toggle backlog panel', async ({ page }) => {
    // Open workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

    // Toggle backlog
    await page.getByRole('button', { name: 'Toggle Backlog' }).click();
    await expect(page.locator('.backlog-panel h3')).toBeVisible();

    // Close backlog
    await page.getByRole('button', { name: 'Toggle Backlog' }).click();
    await expect(page.locator('.backlog-panel')).not.toBeVisible();
  });

  test('should display progress percentage', async ({ page }) => {
    // Setup tasks with different statuses for progress calculation
    await page.addInitScript(() => {
      window.localStorage.setItem('mock_tasks', JSON.stringify([
        {
          id: '1',
          title: 'Done Task',
          status: 'SUCCEEDED',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        },
        {
          id: '2',
          title: 'Pending Task',
          status: 'PENDING',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]));
    });

    await page.reload();
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // Verify progress is displayed (50% = 1/2 tasks completed)
    await expect(page.getByText('50%')).toBeVisible();
  });
});
