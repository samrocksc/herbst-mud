import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: ".",
  testMatch: "*.spec.ts",
  timeout: 120000,
  retries: 0,
  use: {
    baseURL: "http://localhost:5174",
    headless: true,
    viewport: { width: 1280, height: 720 },
    launchOptions: {
      executablePath: "/usr/bin/google-chrome",
    },
  },
  webServer: {
    command: "cd /home/sam/GitHub/herbst-mud/web-client && npx vite --port 5174",
    port: 5174,
    reuseExistingServer: true,
  },
});
