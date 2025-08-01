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
      include: ['vueuc', 'date-fns'],
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
      title: '老九网盘资源数据库',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '老九网盘资源管理数据庫，现代化的网盘资源数据库，支持多网盘自动化转存分享，支持百度网盘，阿里云盘，夸克网盘， 天翼云盘，迅雷云盘，123云盘，115网盘，UC网盘' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      // 开发环境：直接访问后端，生产环境：通过 Nginx 反代
      apiBase: process.env.NODE_ENV === 'production' ? '/api' : 'http://localhost:8080/api',
      // 服务端：开发环境直接访问，生产环境容器内访问
      apiServer: process.env.NODE_ENV === 'production' ? 'http://backend:8080/api' : 'http://localhost:8080/api'
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