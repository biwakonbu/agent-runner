import { test, expect } from '@playwright/test';

test.describe('E2E: CLI Session Validation', () => {
  test.beforeEach(async ({ page }) => {
    // 1. Mock Wails Runtime & Backend
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
            SelectWorkspace: async () => "mock-workspace-id",
            GetExecutionState: async () => "IDLE",
            StartExecution: async () => {}, // Default success
            StopExecution: async () => {},
            PauseExecution: async () => {},
            ResumeExecution: async () => {},
          }
        }
      };
    });

    await page.goto('/');
    
    // Open workspace
    await page.getByRole('button', { name: 'Workspaceを開く' }).click();
    await expect(page.getByRole('button', { name: 'Start' })).toBeVisible();
  });

  test('should show error notification when CLI session is missing', async ({ page }) => {
     // 1. Simulate CLI session missing error from backend
     const errorMessage = "failed to build agent tool plan: CLI session not found (please run 'codex login')";
     
     await page.evaluate((msg) => {
         (window as any).go.main.App.StartExecution = () => Promise.reject(msg);
     }, errorMessage);

      // 2. Trigger startExecution programmatically (since UI controls are removed)
      await page.evaluate(() => {
        // @ts-ignore
        if (window.startExecution) {
          // @ts-ignore
          window.startExecution();
        } else {
          console.error("window.startExecution is not defined");
          throw new Error("window.startExecution is not defined");
        }
      });

      // 3. Verify Error Toast appears
      const toast = page.locator('.toast.error');
      // Wait for toast to be attached first
      await toast.waitFor({ state: 'attached', timeout: 5000 });
      await expect(toast).toBeVisible();
     
     // Check for error message text
     await expect(toast).toContainText('CLI session not found');
     
     // 4. Verify Action Button exists and click it
     const retryButton = toast.locator('button.action-button', { hasText: 'Retry' });
     await expect(retryButton).toBeVisible();
     
     // Mock StartExecution again to verify retry works (this time succeed)
     await page.evaluate(() => {
         (window as any).go.main.App.StartExecution = () => Promise.resolve();
     });
     
     // Click retry
     await retryButton.click();
     
     // Toast should disappear (or new success toast appears, depending on implementation)
     // Since our impl clears old toast when duration passed or removed, but here we just re-run.
     // Let's expect success toast "Execution started"
     await expect(page.locator('.toast.success')).toContainText('Execution started');
  });
});
