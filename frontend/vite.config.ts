import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // By default, the dev command runs in 'development' mode and the build command runs in 'production' mode.
  const env = loadEnv(mode, process.cwd());
  const isDevMode = mode === "development";

  const port = env.VITE_PORT ? parseInt(env.VITE_PORT, 10) : 3000;
  const apiBaseURL = env.VITE_API_BASE_URL || "/api/v1";
  const proxyApiHost = env.VITE_PROXY_API_HOST || "http://localhost:8000";

  return {
    server: {
      port: port,
      strictPort: true,
      host: "0.0.0.0",
      ...(isDevMode && {
        proxy: {
          [apiBaseURL]: {
            target: proxyApiHost,
            changeOrigin: true,
            secure: false,
          },
        },
      }),
      plugins: [react()],
      build: {
        outDir: "dist",
        sourcemap: false,
        minify: "esbuild",
        rollupOptions: {
          output: {
            manualChunks: {
              vendor: ["react", "react-dom"],
            },
          },
        },
      },
    },
  };
});
