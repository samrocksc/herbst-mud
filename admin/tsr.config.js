import { defineConfig } from '@tanstack/router-cli'

export default defineConfig({
  routesDirectory: './src/routes',
  generatedRouteTree: './src/routeTree.gen.ts',
  routeFileIgnorePrefix: '-',
  quoteStyle: 'single',
})