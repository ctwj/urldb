<template>
  <div class="min-h-screen max-h-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
    <Head>
      <title>管理后台 - 老九网盘资源数据库</title>
    </Head>

    <!-- 无障碍：跳到主内容 -->
    <a
      href="#main-content"
      class="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-[9999] focus:px-4 focus:py-2 focus:bg-blue-700 focus:text-white focus:rounded-lg focus:shadow-lg"
    >
      跳到主内容
    </a>

    <!-- 顶部导航栏 -->
    <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="flex items-center justify-between px-6 py-4">
        <!-- 左侧：Logo和标题 -->
        <div class="flex items-center">
          <!-- 窄屏：侧边栏抽屉触发按钮（≤767px 显示） -->
          <button
            type="button"
            class="md:hidden mr-2 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors cursor-pointer"
            aria-label="打开侧边栏菜单"
            @click="sidebarDrawerOpen = true"
          >
            <i class="fas fa-bars text-lg text-gray-700 dark:text-gray-300"></i>
          </button>
          <NuxtLink to="/admin" class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <i class="fas fa-shield-alt text-white text-sm"></i>
            </div>
            <div class="flex items-center space-x-2">
              <h1 class="text-xl font-bold text-gray-900 dark:text-white">管理后台</h1>
              <NuxtLink
                to="/admin/version"
                class="text-xs text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
              >
                v{{ versionInfo.version }}
              </NuxtLink>
            </div>
          </NuxtLink>
        </div>

        <!-- 右侧：状态信息与用户菜单 -->
        <div class="flex items-center space-x-4">
          <!-- ⌘K 快速搜索徽章（可发现性） -->
          <button
            type="button"
            @click="commandPaletteOpen = true"
            class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg text-xs text-gray-500 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors cursor-pointer"
            aria-label="打开命令面板"
          >
            <i class="fas fa-search"></i>
            <span>快速跳转</span>
            <kbd class="px-1.5 py-0.5 bg-white dark:bg-gray-800 rounded border border-gray-200 dark:border-gray-600 font-mono">⌘K</kbd>
          </button>

          <!-- 自动处理状态 -->
          <div class="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-lg px-3 py-2">
            <div class="w-2 h-2 rounded-full animate-pulse" :class="{
              'bg-red-400': !isAutoProcessEnabled,
              'bg-green-400': isAutoProcessEnabled,
            }"></div>
            <span class="text-xs text-gray-700 dark:text-gray-300 font-medium">
              自动处理已<span>{{ isAutoProcessEnabled ? '开启' : '关闭' }}</span>
            </span>
          </div>

          <!-- 自动转存状态 -->
          <div class="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-lg px-3 py-2">
            <div class="w-2 h-2 rounded-full animate-pulse" :class="{
              'bg-red-400': !isAutoTransferEnabled,
              'bg-green-400': isAutoTransferEnabled,
            }"></div>
            <span class="text-xs text-gray-700 dark:text-gray-300 font-medium">
              自动转存已<span>{{ isAutoTransferEnabled ? '开启' : '关闭' }}</span>
            </span>
          </div>

          <!-- 任务状态 -->
          <div
            v-if="taskStore.hasActiveTasks"
            @click="navigateToTasks"
            class="flex items-center gap-2 bg-orange-50 dark:bg-orange-900/20 rounded-lg px-3 py-2 cursor-pointer hover:bg-orange-100 dark:hover:bg-orange-900/30 transition-colors"
          >
            <div class="w-2 h-2 bg-orange-500 rounded-full animate-pulse"></div>
            <span class="text-xs text-orange-700 dark:text-orange-300 font-medium">
              <template v-if="taskStore.runningTaskCount > 0">
                {{ taskStore.runningTaskCount }}个任务运行中
              </template>
              <template v-else>
                {{ taskStore.activeTaskCount }}个任务待处理
              </template>
            </span>
          </div>

          <NuxtLink
            to="/"
            class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
            aria-label="返回前台首页"
          >
            <i class="fas fa-home text-lg" aria-hidden="true"></i>
          </NuxtLink>

          <!-- 用户信息和下拉菜单 -->
          <div ref="userMenuRef" class="relative">
            <button
              type="button"
              @click="showUserMenu = !showUserMenu"
              class="flex items-center space-x-2 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors cursor-pointer"
              :aria-expanded="showUserMenu"
              aria-haspopup="true"
            >
              <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                <i class="fas fa-user text-white text-sm"></i>
              </div>
              <div class="hidden md:block text-left">
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ userStore.user?.username || '管理员' }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">管理员</p>
              </div>
              <i class="fas fa-chevron-down text-xs text-gray-400"></i>
            </button>

            <!-- 下拉菜单内容 -->
            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 z-50 border border-gray-200 dark:border-gray-700"
            >
              <template v-for="item in userMenuItems" :key="item.label || item.type">
                <NuxtLink
                  v-if="item.type === 'link' && item.to"
                  :to="item.to"
                  class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </NuxtLink>

                <button
                  v-else-if="item.type === 'button'"
                  type="button"
                  @click="item.action"
                  class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer"
                  :class="item.className || 'text-gray-700 dark:text-gray-300'"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </button>

                <div
                  v-else-if="item.type === 'divider'"
                  class="border-t border-gray-200 dark:border-gray-700 my-1"
                ></div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- 侧边栏和主内容区域 -->
    <div class="flex main-content">
      <!-- 桌面侧边栏（≥768px 显示，抽出为独立组件，含 localStorage 持久化） -->
      <aside class="hidden md:block w-64 bg-white dark:bg-gray-800 shadow-sm border-r border-gray-200 dark:border-gray-700 h-full overflow-y-auto">
        <AdminSidebarNav />
      </aside>

      <!-- 窄屏侧边栏抽屉（<768px，由顶部 hamburger 触发） -->
      <ClientOnly>
        <n-drawer v-model:show="sidebarDrawerOpen" :width="280" placement="left">
          <n-drawer-content title="导航菜单" closable>
            <AdminSidebarNav @navigate="sidebarDrawerOpen = false" />
          </n-drawer-content>
        </n-drawer>
      </ClientOnly>

      <!-- 主内容区域 -->
      <main id="main-content" tabindex="-1" class="flex-1 p-4 h-full overflow-y-auto focus:outline-none">
        <ClientOnly>
          <n-config-provider :theme="naiveTheme" :theme-overrides="naiveOverrides">
            <n-message-provider>
              <n-notification-provider>
                <n-dialog-provider>
                  <!-- 统一面包屑 -->
                  <div class="mb-4">
                    <AdminBreadcrumb />
                  </div>
                  <!-- 页面内容插槽 -->
                  <slot />
                </n-dialog-provider>
              </n-notification-provider>
            </n-message-provider>
          </n-config-provider>
        </ClientOnly>
      </main>
    </div>

    <!-- Cmd/Ctrl+K 命令面板 -->
    <AdminCommandPalette v-model:open="commandPaletteOpen" />
  </div>
</template>

<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { darkTheme } from 'naive-ui'
import { useUserStore } from '~/stores/user'
import { useSystemConfigStore } from '~/stores/systemConfig'
import { useTaskStore } from '~/stores/task'
import { useTheme } from '~/composables/useTheme'
import { useWebVitals } from '~/composables/useWebVitals'

// 用户状态管理
const userStore = useUserStore()
const router = useRouter()

// 系统配置 store
const systemConfigStore = useSystemConfigStore()

// 任务状态管理
const taskStore = useTaskStore()

// 主题（Naive UI themeOverrides + dark mode 切换）
const { mode, naiveOverrides } = useTheme()
const naiveTheme = computed(() => (mode.value === 'dark' ? darkTheme : null))

systemConfigStore.initConfig(false, true).catch(() => {
  // 配置初始化失败已由 store 内部处理；调试 console.error 已移除
})

// 版本信息
const versionInfo = ref({ version: '1.1.0' })

const fetchVersionInfo = async () => {
  try {
    const response = (await $fetch('/api/version')) as any
    if (response.success) {
      versionInfo.value = response.data
    }
  } catch {
    // 版本信息获取失败不影响页面渲染
  }
}

// 初始化版本信息和任务状态管理
onMounted(() => {
  fetchVersionInfo()
  taskStore.startAutoUpdate()

  // Web Vitals 性能埋点（SC-001 基线，sendBeacon 上报到 /api/web-vitals）
  useWebVitals()

  // 确保在客户端配置被正确载入（防止 SSR 水合问题）
  setTimeout(async () => {
    try {
      await systemConfigStore.initConfig(true, true)
    } catch {
      // 配置刷新失败已由 store 内部处理
    }
  }, 100)
})

onBeforeUnmount(() => {
  taskStore.stopAutoUpdate()
})

// 系统配置
const systemConfig = computed(() => systemConfigStore.config || {})

const isAutoProcessEnabled = computed(() => {
  const value = systemConfig.value?.auto_process_ready_resources
  return value === true || value === 'true' || value === '1'
})

const isAutoTransferEnabled = computed(() => {
  const value = systemConfig.value?.auto_transfer_enabled
  return value === true || value === 'true' || value === '1'
})

// 用户菜单：使用 onClickOutside 替代脆弱的全局 click 监听
const userMenuRef = ref<HTMLElement | null>(null)
const showUserMenu = ref(false)
onClickOutside(userMenuRef, () => {
  showUserMenu.value = false
})

// 命令面板开关（与 CommandPalette 组件 v-model:open 双向绑定）
const commandPaletteOpen = ref(false)

// 窄屏侧边栏抽屉开关（≤767px 由 hamburger 按钮触发，路由切换自动关闭）
const sidebarDrawerOpen = ref(false)
const route = useRoute()
watch(
  () => route.path,
  () => {
    sidebarDrawerOpen.value = false
  },
)

// 处理退出登录
const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 管理员菜单项
const userMenuItems = computed(() => [
  { to: '/admin/tasks', icon: 'fas fa-tasks', label: '任务列表', type: 'link' },
  { to: '/admin/accounts', icon: 'fas fa-user-shield', label: '平台账号', type: 'link' },
  { to: '/admin/api-access-logs', icon: 'fas fa-history', label: 'API访问日志', type: 'link' },
  { to: '/admin/system-logs', icon: 'fas fa-file-alt', label: '系统日志', type: 'link' },
  { to: '/admin/version', icon: 'fas fa-code-branch', label: '版本信息', type: 'link' },
  { type: 'divider' },
  {
    type: 'button',
    icon: 'fas fa-sign-out-alt',
    label: '退出登录',
    action: handleLogout,
    className: 'text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300',
  },
])

// 导航到任务列表页面
const navigateToTasks = () => {
  router.push('/admin/tasks')
}
</script>

<style scoped>
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
.main-content {
  height: calc(100vh - 85px);
}
</style>
