import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.E2E_BASE_URL ?? "http://localhost:3000";
const retries = Number(process.env.E2E_RETRIES ?? (process.env.CI ? "0" : "0"));
const maxFailures = Number(process.env.E2E_MAX_FAILURES ?? (process.env.CI ? "1" : "0"));
const testTimeout = Number(process.env.E2E_TEST_TIMEOUT_MS ?? "45000");
const globalTimeout = Number(process.env.E2E_GLOBAL_TIMEOUT_MS ?? "300000");
const resultDir = process.env.E2E_RESULT_DIR ?? "result";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries,
  maxFailures: process.env.CI ? maxFailures : undefined,
  workers: process.env.CI ? 1 : undefined,
  timeout: testTimeout,
  globalTimeout,
  expect: {
    timeout: 10_000,
  },
  reporter: process.env.CI
    ? [
        ["github"],
        ["html", { open: "never", outputFolder: `${resultDir}/playwright-report` }],
        ["junit", { outputFile: `${resultDir}/test-results/e2e-junit.xml` }],
        ["list"],
      ]
    : [["list"], ["html", { open: "never", outputFolder: `${resultDir}/playwright-report` }]],
  use: {
    baseURL,
    trace: process.env.CI ? "on-first-retry" : "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
    viewport: { width: 1440, height: 900 },
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
    },
    {
      name: "webkit",
      use: { ...devices["Desktop Safari"] },
    },
  ],
  outputDir: `${resultDir}/test-results`,
});
