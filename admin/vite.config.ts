import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    host: true,
    allowedHosts: [
      'the-sewer.taild22ae7.ts.net',
      '.ts.net'
    ],
    cors: {
      origin: '*',
      credentials: false
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      },
      '/users': 'http://localhost:8080',
      '/rooms': 'http://localhost:8080',
      '/npcs': 'http://localhost:8080',
      '/equipment': 'http://localhost:8080',
      '/skills': 'http://localhost:8080',
      '/talents': 'http://localhost:8080',
      '/characters': 'http://localhost:8080',
      '/worlds': 'http://localhost:8080',
      '/game-configs': 'http://localhost:8080',
      '/healthz': 'http://localhost:8080',
      '/admin': 'http://localhost:8080',
    }
  }
})
