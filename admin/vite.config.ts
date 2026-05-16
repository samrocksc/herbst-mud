import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    host: true,
    allowedHosts: [
      "the-sewer.taild22ae7.ts.net",
      ".ts.net"
    ],
    cors: {
      origin: "*",
      credentials: false
    },
    proxy: {
      // All API routes — frontend calls /api/*, backend serves /api/*
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      // Unprotected backend routes that are fetched bare (not under /api).
      // Caution: Vite proxies prefix-match. Any path that is also an SPA
      // route (or has SPA sub-routes) must use bypass to distinguish.
      "/rooms": "http://localhost:8080",
      "/npcs": {
        target: "http://localhost:8080",
        changeOrigin: true,
        bypass(req) {
          // Only proxy exact /npcs (dashboard stats); /npcs/* are SPA routes
          if (req.url !== "/npcs") return "/index.html";
        },
      },
      "/characters": "http://localhost:8080",
      "/equipment": "http://localhost:8080",
      "/skills": {
        target: "http://localhost:8080",
        changeOrigin: true,
        bypass(req) {
          // /skills is also a sidebar SPA page.
          // API data fetch → no text/html Accept; SPA page load → text/html.
          const accept = req.headers.accept ?? "";
          if (accept.includes("text/html")) return "/index.html";
        },
      },
      // Auth endpoint
      "/users": "http://localhost:8080",
      // Health endpoint
      "/healthz": "http://localhost:8080",
      // Admin routes
      "/admin": "http://localhost:8080",
    }
  }
});
