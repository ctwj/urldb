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
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path
        },
        '/uploads': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path
        }
      }
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
      // 客户端API地址：开发环境通过代理，生产环境通过Nginx
      apiBase: '/api',
      // 服务端API地址：通过环境变量配置，支持不同部署方式
      apiServer: process.env.NUXT_PUBLIC_API_SERVER || (process.env.NODE_ENV === 'production' ? 'http://backend:8080/api' : '/api')
    }
  },
  build: {
    transpile: ['naive-ui', 'vueuc', '@css-render/vue3-ssr', '@juggle/resize-observer']
  },
  ssr: true,
  nitro: {
    logLevel: 'info',
    preset: 'node-server',
    storage: {
      redis: {
        driver: 'memory',
        max: 1000
      }
    }
  }
})