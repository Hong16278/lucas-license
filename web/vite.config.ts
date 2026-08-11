import path from "node:path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(import.meta.dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api": "http://localhost:9000",
      "/pay": "http://localhost:9000",
      "/health": "http://localhost:9000",
    },
  },
  build: {
    // 管理后台是单入口应用；当前生产包 gzip 后约 308 KB，可接受。
    chunkSizeWarningLimit: 1200,
  },
})
