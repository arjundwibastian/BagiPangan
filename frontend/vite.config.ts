import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ command }) => ({
  plugins: [react()],
  // Development always uses the local proxy, even with an old .env.local.
  define: command === "serve" ? { "import.meta.env.VITE_API_BASE_URL": '""' } : undefined,
  server: {
    port: 3000,
    proxy: {
      "^/api/v1/users/[^/]+/food-requests": {
        target: "http://127.0.0.1:18083",
        changeOrigin: true,
      },
      "/api/v1/food-listings": {
        target: "http://127.0.0.1:18082",
        changeOrigin: true,
      },
      "/api/v1/food-requests": {
        target: "http://127.0.0.1:18083",
        changeOrigin: true,
      },
      "/api/v1/claims": {
        target: "http://127.0.0.1:18084",
        changeOrigin: true,
      },
      "/api/v1/users": {
        target: "http://127.0.0.1:18081",
        changeOrigin: true,
      },
      "/api/v1/auth": {
        target: "http://127.0.0.1:18081",
        changeOrigin: true,
      },
    },
  },
}));
