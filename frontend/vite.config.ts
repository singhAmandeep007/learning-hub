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
  const validProducts = (env.VITE_VALID_PRODUCTS || process.env.VALID_PRODUCTS || "").trim();

  if (!validProducts) {
    throw new Error("VITE_VALID_PRODUCTS (or VALID_PRODUCTS) is required");
  }

  return {
    plugins: [react()],
    define: {
      "import.meta.env.VITE_VALID_PRODUCTS": JSON.stringify(validProducts),
    },
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
    },
    build: {
      outDir: "dist",
      sourcemap: false,
      rolldownOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes("node_modules")) {
              return;
            }

            if (id.includes("react") || id.includes("react-dom") || id.includes("scheduler")) {
              return "react-vendor";
            }

            if (id.includes("@tanstack")) {
              return "query-vendor";
            }

            if (id.includes("@tiptap") || id.includes("prosemirror") || id.includes("@floating-ui")) {
              return "editor-vendor";
            }

            return "vendor";
          },
        },
      },
    },
  };
});
