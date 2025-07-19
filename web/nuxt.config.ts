import AutoImport from 'unplugin-auto-import/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  vite: {
    clearScreen: false,
    plugins: [
      AutoImport({
        imports: [
          {
            'naive-ui': [
              'useDialog',
              'useMessage',
              'useNotification',
              'useLoadingBar'
            ]
          }
        ]
      }),
      Components({
        resolvers: [NaiveUiResolver()]
      })
    ],
    optimizeDeps: {
      include: ['naive-ui', 'vueuc', 'date-fns'],
      exclude: ["oxc-parser"] // 强制使用 WASM 版本
    }
  },
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  css: [
    '~/assets/css/main.css',
    'vfonts/Lato.css',
    'vfonts/FiraCode.css',
    '@fortawesome/fontawesome-free/css/all.min.css', // 本地Font Awesome
  ],
  app: {
    head: {
      title: '开源网盘资源数据库',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '开源网盘资源数据库 - 一个现代化的资源管理系统' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.API_BASE || 'http://localhost:8080/api'
    }
  },
  build: {
    transpile: ['naive-ui', 'vueuc', '@css-render/vue3-ssr', '@juggle/resize-observer']
  },
  ssr: true,
  nitro: {
    logLevel: 'verbose',
    preset: 'node-server'
  }
})