import { test } from "./fixtures";
import { expect } from "@playwright/test";

test.describe("Admin panel login", () => {
  test("logs in as sma and lands on the dashboard", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByRole("heading", { name: "Herbst MUD Admin" })).toBeVisible();
    await page.fill("#username", "sma");
    await page.fill("#password", "sma");
    await page.click("button[type='submit']");
    await page.waitForURL("**/dashboard");
    // Dashboard should render with its world selector present
    await expect(page.locator("select").first()).toBeVisible({ timeout: 15000 });
  });
});
