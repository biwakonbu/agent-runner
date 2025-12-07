import { test, expect } from '@playwright/test';

test.describe('Resource Matrix HUD', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/');
    });

    test('should be visible in the DOM', async ({ page }) => {
        // Workspace selection might be needed first if not mocked
        // Assuming we are on the main screen or HUD is global logic
        // But HUD is inside {#if workspaceId} block in App.svelte
        
        // Select a workspace first
        const workspaceSelector = page.locator('.workspace-selector');
        if (await workspaceSelector.isVisible()) {
             // Click first workspace or mock selection
             await page.click('.project-card');
        }

        const hud = page.locator('.process-hud-container');
        await expect(hud).toBeVisible();
        await expect(hud).toContainText('IDLE');
    });

    test('should toggle expansion', async ({ page }) => {
        // Select workspace
        const workspaceSelector = page.locator('.workspace-selector');
        if (await workspaceSelector.isVisible()) {
             await page.click('.project-card');
        }

        const hud = page.locator('.process-hud-container');
        const header = hud.locator('.hud-header');
        
        // Initial: Not expanded
        await expect(hud).not.toHaveClass(/expanded/);

        // Click to expand
        await header.click();
        await expect(hud).toHaveClass(/expanded/);

        // Check content visibility
        const content = hud.locator('.hud-content');
        await expect(content).toBeVisible();

        // Click to collapse
        await header.click();
        await expect(hud).not.toHaveClass(/expanded/);
    });
});
