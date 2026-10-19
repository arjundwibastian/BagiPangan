import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      "^/api/v1/users/[^/]+/food-requests": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
      "/api/v1/food-listings": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
      "/api/v1/food-requests": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
      "/api/v1/claims": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
      "/api/v1/users": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
      "/api/v1/auth": {
        target: "https://app.bagipangan.my.id",
        changeOrigin: true,
      },
    },
  },
});
