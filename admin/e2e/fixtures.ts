import { test as base } from "@playwright/test";

export const test = base.extend({
  adminPage: async ({ page }, use) => {
    await page.goto("/login");
    await page.fill("#username", "sma");
    await page.fill("#password", "sma");
    await page.click("button[type='submit']");
    await page.waitForURL("**/dashboard");
    await use(page);
  },
});
