import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: ".",
  testMatch: "*.spec.ts",
  timeout: 120000,
  retries: 0,
  use: {
    baseURL: "http://localhost:5173",
    headless: true,
    viewport: { width: 1280, height: 720 },
    launchOptions: {
      executablePath: "/usr/bin/google-chrome",
    },
  },
  webServer: {
    command: "cd /home/sam/GitHub/herbst-mud/admin && npx vite --port 5173",
    port: 5173,
    reuseExistingServer: true,
  },
});
