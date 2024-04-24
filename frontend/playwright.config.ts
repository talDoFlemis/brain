import { defineConfig, devices } from "@playwright/test";

// Set webServer.url and use.baseURL with the location of the WebServer respecting the correct set port
const baseURL = "http://localhost:3000";

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: baseURL,
    trace: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        launchOptions: {
          executablePath: process.env
            .PLAYWRIGHT_LAUNCH_OPTIONS_EXECUTABLE_PATH as string,
        },
      },
    },
    ...(process.env.CI
      ? [
          {
            name: "firefox",
            use: { ...devices["Desktop Firefox"] },
          },
          {
            name: "webkit",
            use: { ...devices["Desktop Safari"] },
          },
        ]
      : []),
  ],

  webServer: {
    command: "task dev",
    url: baseURL,
    reuseExistingServer: true,
  },
});
