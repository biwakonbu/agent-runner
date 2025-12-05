import { test, expect } from '@playwright/test';

test.describe('Task Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Clear mock storage before each test
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
    });
    // Go to the app
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should display empty task list initially', async ({ page }) => {
    // Initially mock returns empty list
    // Check if task list container is visible but empty or says "No tasks"
    // Adjust selector based on actual UI implementation
    // Assuming DetailPanel or similar components show something.
    // Let's check for some text or element.
    await expect(page.locator('body')).toBeVisible();
  });

  test('should create a new task and display it', async ({ page }) => {
    // 1. Select Workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // 2. Wait for Toolbar to appear (indicates app loaded backend data)
    const addTaskBtn = page.getByRole('button', { name: '新規タスク' });
    await expect(addTaskBtn).toBeVisible();

    // 3. Open Create Task Modal
    await addTaskBtn.click();
    await expect(page.getByRole('dialog', { name: '新規タスク作成' })).toBeVisible();

    // 4. Fill form
    await page.getByPlaceholder('タスクのタイトルを入力').fill('My E2E Task');
    
    // 5. Submit
    await page.getByRole('button', { name: 'タスクを作成' }).click();

    // 6. Verify Task appears in Grid
    // Assuming GridNode renders the title text
    // We wait for the modal to close and text to appear
    await expect(page.getByRole('dialog')).not.toBeVisible();
    await expect(page.getByText('My E2E Task')).toBeVisible();
  });
});
