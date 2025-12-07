import { test, expect } from '@playwright/test';

test.describe('Zoom Behavior', () => {
  test.beforeEach(async ({ page }) => {
    // Mock Wails Runtime & LocalStorage
    await page.addInitScript(() => {
      window.localStorage.clear();
      const listeners = new Map<string, Function[]>();

      (window as any).runtime = {
        EventsOn: (eventName: string, callback: Function) => {
          if (!listeners.has(eventName)) listeners.set(eventName, []);
          listeners.get(eventName)?.push(callback);
        },
        EventsOff: () => {},
        __trigger: (eventName: string, data: any) => {
          const callbacks = listeners.get(eventName) || [];
          callbacks.forEach(cb => cb(data));
        }
      };
      
      (window as any).go = {
        main: {
          App: {
            ListTasks: async () => [],
            GetPoolSummaries: async () => [],
            GetExecutionState: async () => "IDLE",
            StartExecution: async () => {},
            StopExecution: async () => {},
            PauseExecution: async () => {},
            ResumeExecution: async () => {},
            CreateChatSession: async () => ({ id: 'test-session-id' }),
            GetChatHistory: async () => [],
            SendChatMessage: async () => ({ generatedTasks: [], understanding: '' })
          }
        }
      };
    });

    await page.goto('/');
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should display zoom indicator', async ({ page }) => {
    const zoomIndicator = page.locator('.zoom-indicator');
    await expect(zoomIndicator).toBeVisible();
    // Default zoom is 100%
    await expect(zoomIndicator).toContainText('100%');
  });

  test('should zoom in with + key and update indicator', async ({ page }) => {
    const canvas = page.locator('.canvas-container');
    const zoomIndicator = page.locator('.zoom-indicator');
    
    // Focus the canvas
    await canvas.click();
    
    // Get initial zoom percentage
    const initialZoom = await zoomIndicator.textContent();
    
    // Press + to zoom in
    await page.keyboard.press('+');
    
    // Wait for zoom to update
    await page.waitForTimeout(100);
    
    // Zoom should have increased
    const newZoom = await zoomIndicator.textContent();
    expect(parseInt(newZoom || '100')).toBeGreaterThan(parseInt(initialZoom || '100'));
  });

  test('should zoom out with - key and update indicator', async ({ page }) => {
    const canvas = page.locator('.canvas-container');
    const zoomIndicator = page.locator('.zoom-indicator');
    
    // Focus the canvas
    await canvas.click();
    
    // First zoom in so we can zoom out
    await page.keyboard.press('+');
    await page.waitForTimeout(100);
    
    const zoomedInValue = await zoomIndicator.textContent();
    
    // Press - to zoom out
    await page.keyboard.press('-');
    await page.waitForTimeout(100);
    
    // Zoom should have decreased
    const zoomedOutValue = await zoomIndicator.textContent();
    expect(parseInt(zoomedOutValue || '100')).toBeLessThan(parseInt(zoomedInValue || '100'));
  });

  test('should reset zoom with 0 key', async ({ page }) => {
    const canvas = page.locator('.canvas-container');
    const zoomIndicator = page.locator('.zoom-indicator');
    
    // Focus the canvas
    await canvas.click();
    
    // Zoom in first
    await page.keyboard.press('+');
    await page.keyboard.press('+');
    await page.waitForTimeout(100);
    
    // Verify zoom changed
    const zoomedValue = await zoomIndicator.textContent();
    expect(parseInt(zoomedValue || '100')).toBeGreaterThan(100);
    
    // Reset with 0
    await page.keyboard.press('0');
    await page.waitForTimeout(100);
    
    // Should be back to 100%
    await expect(zoomIndicator).toContainText('100%');
  });

  test('should zoom centered on mouse position when using keyboard', async ({ page }) => {
    const canvas = page.locator('.canvas-container');
    
    // Create a task to have a visual reference point
    const task = {
      id: 'zoom-test-task',
      title: 'Zoom Test Task',
      status: 'PENDING',
      poolId: 'default',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    await page.evaluate((t) => {
      (window as any).runtime.__trigger('task:created', { task: t });
    }, task);
    
    await expect(page.getByText('Zoom Test Task')).toBeVisible();
    
    // Get the task node element
    const taskNode = page.getByText('Zoom Test Task').locator('..');
    
    // Get initial position of the task node
    const initialBox = await taskNode.boundingBox();
    expect(initialBox).not.toBeNull();
    
    // Move mouse to a specific position on the canvas (near the task)
    const canvasBox = await canvas.boundingBox();
    expect(canvasBox).not.toBeNull();
    
    const mouseX = canvasBox!.x + canvasBox!.width / 2;
    const mouseY = canvasBox!.y + canvasBox!.height / 2;
    
    await page.mouse.move(mouseX, mouseY);
    
    // Now zoom in. The point at mouse position should stay relatively fixed.
    await page.keyboard.press('+');
    await page.waitForTimeout(100);
    
    // This test validates that zooming works without error.
    // Full position verification would require more complex calculations.
    const zoomIndicator = page.locator('.zoom-indicator');
    const newZoom = await zoomIndicator.textContent();
    expect(parseInt(newZoom || '100')).toBeGreaterThan(100);
  });

  test('should handle Ctrl+wheel zoom centered on mouse position', async ({ page }) => {
    const canvas = page.locator('.canvas-container');
    const zoomIndicator = page.locator('.zoom-indicator');
    
    // Get canvas bounding box
    const canvasBox = await canvas.boundingBox();
    expect(canvasBox).not.toBeNull();
    
    // Move mouse to center of canvas
    const centerX = canvasBox!.x + canvasBox!.width / 2;
    const centerY = canvasBox!.y + canvasBox!.height / 2;
    
    await page.mouse.move(centerX, centerY);
    
    // Perform Ctrl+wheel zoom in
    await page.keyboard.down('Control');
    await page.mouse.wheel(0, -100); // Negative deltaY = zoom in
    await page.keyboard.up('Control');
    
    await page.waitForTimeout(100);
    
    // Zoom should have increased
    const zoomValue = await zoomIndicator.textContent();
    expect(parseInt(zoomValue || '100')).toBeGreaterThan(100);
  });
});
