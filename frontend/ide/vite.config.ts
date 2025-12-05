import { defineConfig } from 'vite'
import { svelte, vitePreprocess } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const isE2E = process.env.VITE_TEST_E2E === 'true';
  return {
    plugins: [svelte({
      preprocess: vitePreprocess(),
    })],
    resolve: {
      alias: isE2E ? [
        {
          find: /.*wailsjs\/go\/main\/App/,
          replacement: path.resolve(__dirname, './src/mocks/wails.js'),
        }
      ] : {}
    }
  }
})
