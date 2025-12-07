import { test, expect } from '@playwright/test';

test.describe('Task Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Clear mock storage before each test
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    // Go to the app
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should display welcome screen initially', async ({ page }) => {
    // Welcome 画面が表示されることを確認
    await expect(page.locator('body')).toBeVisible();
    // "Workspaceを開く" ボタンが表示されている
    await expect(page.getByRole('button', { name: 'Workspaceを開く' })).toBeVisible();
  });

  test('should display recent workspaces list', async ({ page }) => {
    // 最近使ったワークスペースの表示を確認
    await expect(page.getByText('最近使用したワークスペース')).toBeVisible();
    // モックワークスペースが表示されていること
    await expect(page.getByText('My Project')).toBeVisible();
    await expect(page.getByText('Another Project')).toBeVisible();
  });

  test('should open workspace and show toolbar', async ({ page }) => {
    // 1. Select Workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // 2. Wait for Toolbar to appear (indicates app loaded backend data)
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();

    // 3. Verify view mode buttons are visible
    await expect(page.getByRole('button', { name: 'Graph View' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'WBS View' })).toBeVisible();
  });
});

test.describe('Workspace Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should show welcome screen on initial load', async ({ page }) => {
    // Welcome 画面が表示されることを確認
    await expect(page.locator('body')).toBeVisible();
    // "Workspaceを開く" ボタンが表示されている
    await expect(page.getByRole('button', { name: 'Workspaceを開く' })).toBeVisible();
  });

  test('should open recent workspace', async ({ page }) => {
    // 最近使ったワークスペースをクリック
    await page.getByRole('button', { name: /My Project/ }).click();

    // ワークスペースが開かれ、Toolbarが表示される
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });
});

test.describe('Execution Controls', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');

    // Open workspace first
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should show execution controls', async ({ page }) => {
    // 初期状態はIDLE
    await expect(page.getByText('IDLE')).toBeVisible();

    // Start ボタンが表示されていること
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should start execution', async ({ page }) => {
    // Start ボタンをクリック
    await page.getByRole('button', { name: 'Start' }).click();

    // 実行状態が変わることを期待（モックでは即座にRUNNINGになる）
    // Note: モックの実装によっては状態変化を待つ必要があるかもしれない
  });
});

test.describe('View Mode', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');

    // Open workspace first
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should switch between Graph and WBS view', async ({ page }) => {
    // 初期状態はGraph View - title属性でボタンを検索
    const graphButton = page.getByTitle('Graph View');
    const wbsButton = page.getByTitle('WBS View');

    await expect(graphButton).toBeVisible();
    await expect(wbsButton).toBeVisible();

    // 初期状態でGraph ViewのGridCanvasがDOMに存在することを確認
    // (canvas-containerは絶対配置で高さ0になる場合があるため、要素の存在を確認)
    await expect(page.locator('.canvas-container')).toHaveCount(1);

    // WBS Viewに切り替え
    await wbsButton.click();

    // WBS要素が表示されることを確認（WBS Listビューのクラスを確認）
    await expect(page.locator('.wbs-list-view')).toBeVisible();
    // Graph ViewのGridCanvasはDOMから消えている
    await expect(page.locator('.canvas-container')).toHaveCount(0);

    // Graph Viewに戻す
    await graphButton.click();

    // Graph ViewのGridCanvasがDOMに再追加されることを確認
    await expect(page.locator('.canvas-container')).toHaveCount(1);
    // WBS Listビューは消えている
    await expect(page.locator('.wbs-list-view')).toHaveCount(0);
  });
});

test.describe('Chat Window', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');

    // Open workspace first
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should display chat window', async ({ page }) => {
    // チャットウィンドウが表示されていること（placeholderでtextareaを検索）
    await expect(page.locator('textarea[placeholder="Ask Multiverse..."]')).toBeVisible();

    // タブが表示されていること（exact: trueで完全一致を要求）
    await expect(page.getByRole('button', { name: 'General', exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Log', exact: true })).toBeVisible();
  });

  test('should close and reopen chat window', async ({ page }) => {
    // Close ボタンをクリック
    await page.getByRole('button', { name: 'Close' }).click();

    // チャットウィンドウが非表示になる
    await expect(page.locator('textarea[placeholder="Ask Multiverse..."]')).not.toBeVisible();

    // FABボタンでチャットを再表示
    await page.getByRole('button', { name: 'Open Chat' }).click();

    // チャットウィンドウが再表示される
    await expect(page.locator('textarea[placeholder="Ask Multiverse..."]')).toBeVisible();
  });
});

test.describe('Backlog Panel', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      window.localStorage.removeItem('mock_tasks');
      window.localStorage.removeItem('mock_workspaces');
      window.localStorage.removeItem('mock_backlog');
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');

    // Open workspace first
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should toggle backlog panel', async ({ page }) => {
    // Toggle Backlog ボタンをクリック
    await page.getByRole('button', { name: 'Toggle Backlog' }).click();

    // バックログパネルが表示される（h3ヘッダーを確認）
    await expect(page.locator('.backlog-panel h3')).toBeVisible();

    // もう一度クリックで閉じる
    await page.getByRole('button', { name: 'Toggle Backlog' }).click();

    // バックログパネルが非表示になる
    await expect(page.locator('.backlog-panel')).not.toBeVisible();
  });
});

test.describe('Task Status Display', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      // 既存タスクをセットアップ
      window.localStorage.setItem('mock_tasks', JSON.stringify([
        {
          id: 'task-pending',
          title: 'Pending Task',
          status: 'PENDING',
          poolId: 'default',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        },
        {
          id: 'task-running',
          title: 'Running Task',
          status: 'RUNNING',
          poolId: 'codegen',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]));
    });
    page.on('console', msg => console.log('PAGE LOG:', msg.text()));
    await page.goto('/');
  });

  test('should display tasks in grid', async ({ page }) => {
    // 1. Select Workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();

    // 2. タスクが表示されるのを待つ
    await expect(page.getByText('Pending Task')).toBeVisible();
    await expect(page.getByText('Running Task')).toBeVisible();
  });
});
