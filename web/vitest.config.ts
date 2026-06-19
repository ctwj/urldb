import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'

// 前端测试基础设施 (Constitution Principle I: Test-First)
export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'happy-dom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['composables/**', 'components/Admin/**', 'utils/**'],
      exclude: [
        '**/*.test.ts',
        '**/*.spec.ts',
        'node_modules/**',
        '.nuxt/**',
        'output/**',
        'tests/**',
      ],
    },
    include: ['tests/**/*.{test,spec}.ts'],
  },
  resolve: {
    alias: {
      '~': fileURLToPath(new URL('./', import.meta.url)),
      '@': fileURLToPath(new URL('./', import.meta.url)),
    },
  },
})
