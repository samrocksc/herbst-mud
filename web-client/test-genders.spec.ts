import { test, expect } from '@playwright/test';

test('character creation shows genders', async ({ page }) => {
  // Navigate to web client
  await page.goto('http://100.67.206.65:5174');

  // Login
  await page.getByLabel('Email').fill('sma');
  await page.getByLabel('Password').fill('sma');
  await page.getByRole('button', { name: 'Login' }).click();

  // Wait for dashboard to load
  await page.waitForSelector('.world-card', { timeout: 5000 });

  // Click on Ooze Surfers world
  await page.locator('.world-card').filter({ hasText: 'Ooze Surfers' }).click();

  // Click Create New Character
  await page.getByRole('button', { name: 'Create New Character' }).click();

  // Wait for character creation form to load
  await page.waitForSelector('text=Create Character', { timeout: 5000 });

  // Check that the Gender dropdown is populated
  const genderSelect = page.getByLabel('Gender');
  await expect(genderSelect).toBeVisible();

  const options = await genderSelect.locator('option').all();
  expect(options.length).toBeGreaterThan(0);

  // Check if "it" gender is available (world 2 should have this gender)
  const optionValues = await Promise.all(options.map(o => o.getAttribute('value')));
  console.log('Available genders:', optionValues);

  // The "It" gender should be in the list for world 2
  expect(optionValues).toContain('It');
});
