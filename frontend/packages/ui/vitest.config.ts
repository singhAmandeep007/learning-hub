import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    watch: false,
    testTimeout: 10000,
    hookTimeout: 10000,
    teardownTimeout: 2000,
    environment: "jsdom",
    setupFiles: ["./src/test/setup.ts"],
    include: ["src/**/*.{test,spec}.{ts,tsx}"],
    exclude: ["src/**/*.stories.*"],
  },
});
