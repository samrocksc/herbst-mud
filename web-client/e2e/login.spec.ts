import { test, expect } from "@playwright/test";

test.describe("Player web client login", () => {
  test("logs in as sma, selects Ooze Surfers and smack, reaches game", async ({ page }) => {
    await page.goto("/");
    await expect(page.getByRole("heading", { name: "HERBST MUD" })).toBeVisible();

    // Login
    await page.fill("input[placeholder='username']", "sma");
    await page.fill("input[type='password']", "sma");
    await page.click("button[type='submit']:has-text('ENTER THE WORLD')");

    // World selection screen
    await expect(page.getByText("Available Worlds")).toBeVisible({ timeout: 15000 });
    await page.getByText("Ooze Surfers", { exact: true }).click();

    // Character selection screen
    await expect(page.getByText("Select Character")).toBeVisible({ timeout: 15000 });
    await page.getByText("smack", { exact: true }).first().click();

    // Game screen — room view should render (smack is Level 24 in Ooze Surfers)
    await expect(page.getByText("A Dank Ass Closet")).toBeVisible({ timeout: 15000 });
  });
});
