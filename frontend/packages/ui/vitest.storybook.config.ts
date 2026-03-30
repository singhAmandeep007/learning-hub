import path from "node:path";
import { fileURLToPath } from "node:url";

import { storybookTest } from "@storybook/addon-vitest/vitest-plugin";
import { playwright } from "@vitest/browser-playwright";
import { defineConfig } from "vitest/config";

const dirname = typeof __dirname !== "undefined" ? __dirname : path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [storybookTest({ configDir: path.join(dirname, ".storybook") })],
  test: {
    watch: false,
    testTimeout: 20000,
    hookTimeout: 20000,
    browser: {
      enabled: true,
      provider: playwright({ launchOptions: { headless: true } }),
      instances: [{ browser: "chromium" }],
    },
  },
});
