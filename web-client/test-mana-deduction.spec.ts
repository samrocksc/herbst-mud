import { test, expect } from "@playwright/test";

/**
 * Test: Casting "Trash Bolt" correctly deducts mana from the player's UI.
 *
 * Flow:
 * 1. Login as sma/sma
 * 2. Select Ooze Surfers world
 * 3. Select "smack" character (has Trash Bolt in slot 1, mana_cost=5)
 * 4. Enter game, wait for room screen
 * 5. Click Gizmo -> Attack -> Confirm to start combat
 * 6. Read initial mana from header
 * 7. Press "1" to queue Trash Bolt
 * 8. Wait for combat tick (~1.5s)
 * 9. Read mana after casting
 * 10. Assert mana decreased by exactly 5
 */

test("Trash Bolt deducts 5 mana from player", async ({ page }) => {
  // ── 1. Login ──────────────────────────────────────────────────────────────
  await page.goto("http://localhost:5174");
  await page.waitForSelector('text=Herbst MUD', { timeout: 10000 });

  // Input component renders <label> without htmlFor, so use placeholder-based locators
  await page.locator('input[placeholder="username"]').fill("sma");
  await page.locator('input[type="password"]').fill("sma");
  await page.getByRole("button", { name: "ENTER THE WORLD" }).click();

  // ── 2. Select Ooze Surfers world ──────────────────────────────────────────
  await page.waitForSelector("text=Available Worlds", { timeout: 10000 });
  await page.locator("li").filter({ hasText: "Ooze Surfers" }).click();

  // ── 3. Select "smack" character ───────────────────────────────────────────
  await page.waitForSelector("text=Select Character", { timeout: 10000 });
  await page.locator("li").filter({ hasText: "smack" }).first().click();

  // ── 4. Wait for game to load and WebSocket to connect ────────────────────
  // The header shows connection status. Wait for "online" indicator.
  await page.waitForSelector("text=online", { timeout: 15000 });

  // Wait for room screen to appear with Gizmo
  await page.waitForSelector("text=Gizmo", { timeout: 15000 });

  // ── 5. Start combat with Gizmo ────────────────────────────────────────────
  // Click on Gizmo's name to expand the action menu
  const gizmoButton = page.locator("button", { hasText: "Gizmo" }).first();
  await gizmoButton.click();

  // Click "Attack" button
  const attackButton = page.locator("button", { hasText: "Attack" }).first();
  await attackButton.click();

  // Click "Confirm" button
  const confirmButton = page.locator("button", { hasText: "Confirm" }).first();
  await confirmButton.click();

  // ── 6. Wait for combat to start ──────────────────────────────────────────
  await page.waitForSelector("text=Combat", { timeout: 10000 });
  await page.waitForSelector("text=Round", { timeout: 10000 });

  // ── 7. Read initial mana from header ─────────────────────────────────────
  // The header shows: "MANA {value}/{max}"
  const manaHeader = page.locator("header").locator("text=MANA");
  await manaHeader.waitFor({ state: "visible", timeout: 5000 });

  const manaText = await manaHeader.textContent();
  console.log(`[TEST] Initial mana text: "${manaText}"`);

  // Parse "MANA 50/50" -> { current: 50, max: 50 }
  const manaMatch = manaText?.match(/MANA\s+(\d+)\/(\d+)/);
  expect(manaMatch, `Could not parse mana from "${manaText}"`).not.toBeNull();

  const initialMana = parseInt(manaMatch![1], 10);
  const maxMana = parseInt(manaMatch![2], 10);
  console.log(`[TEST] Initial mana: ${initialMana}/${maxMana}`);

  expect(initialMana).toBeGreaterThanOrEqual(5);
  expect(maxMana).toBeGreaterThan(0);

  // ── 8. Queue Trash Bolt by pressing "1" ──────────────────────────────────
  // Trash Bolt is in slot 1 for the "smack" character
  await page.keyboard.press("1");

  // Wait for the queued action indicator
  await page.waitForSelector("text=Queued", { timeout: 5000 });

  // ── 9. Wait for combat tick to resolve the action ────────────────────────
  // Combat ticks every 1.5s. Wait 3s to ensure the tick fires and resolves.
  await page.waitForTimeout(3500);

  // ── 10. Read mana after casting ──────────────────────────────────────────
  const manaTextAfter = await manaHeader.textContent();
  console.log(`[TEST] Mana after casting: "${manaTextAfter}"`);

  const manaMatchAfter = manaTextAfter?.match(/MANA\s+(\d+)\/(\d+)/);
  expect(manaMatchAfter, `Could not parse mana from "${manaTextAfter}"`).not.toBeNull();

  const manaAfter = parseInt(manaMatchAfter![1], 10);
  console.log(`[TEST] Mana after: ${manaAfter}/${maxMana}`);

  // ── 11. Assert mana decreased by exactly 5 ───────────────────────────────
  const expectedMana = initialMana - 5;
  expect(manaAfter).toBe(expectedMana);
});
