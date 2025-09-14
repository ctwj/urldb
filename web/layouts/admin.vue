<template>
  <div class="min-h-screen max-h-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
    <!-- 设置通用title -->
    <Head>
      <title>管理后台 - 老九网盘资源数据库</title>
    </Head>
    <!-- 顶部导航栏 -->
    <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="flex items-center justify-between px-6 py-4">
        <!-- 左侧：Logo和标题 -->
        <div class="flex items-center">
          <NuxtLink to="/admin" class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <i class="fas fa-shield-alt text-white text-sm"></i>
            </div>
            <div class="flex items-center space-x-2">
              <h1 class="text-xl font-bold text-gray-900 dark:text-white">管理后台</h1>
              <!-- 版本信息 -->
              <NuxtLink 
                to="/admin/version" 
                class="text-xs text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
              >
                v{{ versionInfo.version }}
              </NuxtLink>
            </div>
          </NuxtLink>
        </div>



        <!-- 右侧：状态信息和用户菜单 -->
        <div class="flex items-center space-x-4">
          <!-- 自动处理状态 -->
          <div class="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-lg px-3 py-2">
            <div class="w-2 h-2 rounded-full animate-pulse" :class="{ 
              'bg-red-400': !isAutoProcessEnabled,
              'bg-green-400': isAutoProcessEnabled
            }"></div>
            <span class="text-xs text-gray-700 dark:text-gray-300 font-medium">
              自动处理已<span>{{ isAutoProcessEnabled ? '开启' : '关闭' }}</span>
            </span>
          </div>
          
          <!-- 自动转存状态 -->
          <div class="flex items-center gap-2 bg-gray-100 dark:bg-gray-700 rounded-lg px-3 py-2">
            <div class="w-2 h-2 rounded-full animate-pulse" :class="{ 
              'bg-red-400': !isAutoTransferEnabled,
              'bg-green-400': isAutoTransferEnabled
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
          
          <NuxtLink to="/" class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
            <i class="fas fa-home text-lg"></i>
          </NuxtLink>
          
          <!-- 用户信息和下拉菜单 -->
          <div class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center space-x-2 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
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
                <!-- 链接菜单项 -->
                <NuxtLink
                  v-if="item.type === 'link' && item.to"
                  :to="item.to"
                  class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </NuxtLink>
                
                <!-- 按钮菜单项 -->
                <button
                  v-else-if="item.type === 'button'"
                  @click="item.action"
                  class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
                  :class="item.className || 'text-gray-700 dark:text-gray-300'"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </button>
                
                <!-- 分割线 -->
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
      <!-- 侧边栏 -->
      <aside class="w-64 bg-white dark:bg-gray-800 shadow-sm border-r border-gray-200 dark:border-gray-700 h-full overflow-y-auto">
        <nav class="mt-8">
          <div class="px-4 space-y-6">
            <!-- 仪表盘 -->
            <div>
              <h3 class="px-4 mb-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                仪表盘
              </h3>
              <div class="space-y-1">
                <NuxtLink
                  v-for="item in dashboardItems"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center px-4 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
                  :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active(useRoute()) }"
                >
                  <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>

            <!-- 数据管理 -->
            <div>
              <button
                @click="toggleGroup('dataManagement')"
                class="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              >
                <span>数据管理</span>
                <i 
                  class="fas fa-chevron-down text-xs transition-transform duration-200"
                  :class="{ 'rotate-180': expandedGroups.dataManagement }"
                ></i>
              </button>
              <div 
                v-show="expandedGroups.dataManagement"
                class="space-y-1 mt-2"
              >
                <NuxtLink
                  v-for="item in dataManagementItems"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center px-8 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
                  :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active(useRoute()) }"
                >
                  <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>

            <!-- 系统配置 -->
            <div>
              <button
                @click="toggleGroup('systemConfig')"
                class="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              >
                <span>系统配置</span>
                <i 
                  class="fas fa-chevron-down text-xs transition-transform duration-200"
                  :class="{ 'rotate-180': expandedGroups.systemConfig }"
                ></i>
              </button>
              <div 
                v-show="expandedGroups.systemConfig"
                class="space-y-1 mt-2"
              >
                <NuxtLink
                  v-for="item in systemConfigItems"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center px-8 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
                  :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active(useRoute()) }"
                >
                  <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>

            <!-- 运营管理 -->
            <div>
              <button
                @click="toggleGroup('operation')"
                class="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              >
                <span>运营管理</span>
                <i 
                  class="fas fa-chevron-down text-xs transition-transform duration-200"
                  :class="{ 'rotate-180': expandedGroups.operation }"
                ></i>
              </button>
              <div 
                v-show="expandedGroups.operation"
                class="space-y-1 mt-2"
              >
                <NuxtLink
                  v-for="item in operationItems"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center px-8 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
                  :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active(useRoute()) }"
                >
                  <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>

            <!-- 统计分析 -->
            <div>
              <button
                @click="toggleGroup('statistics')"
                class="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              >
                <span>统计分析</span>
                <i 
                  class="fas fa-chevron-down text-xs transition-transform duration-200"
                  :class="{ 'rotate-180': expandedGroups.statistics }"
                ></i>
              </button>
              <div 
                v-show="expandedGroups.statistics"
                class="space-y-1 mt-2"
              >
                <NuxtLink
                  v-for="item in statisticsItems"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center px-8 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
                  :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active(useRoute()) }"
                >
                  <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>
          </div>
        </nav>
      </aside>
      
      <!-- 主内容区域 -->
      <main class="flex-1 p-4 h-full overflow-y-auto">
        <ClientOnly>
          <n-message-provider>
            <n-notification-provider>
              <n-dialog-provider>
                <!-- 页面内容插槽 -->
                <slot />
              </n-dialog-provider>
            </n-notification-provider>
          </n-message-provider>
        </ClientOnly>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '~/stores/user'
import { useSystemConfigStore } from '~/stores/systemConfig'
import { useTaskStore } from '~/stores/task'

// 用户状态管理
const userStore = useUserStore()
const router = useRouter()

// 系统配置store
const systemConfigStore = useSystemConfigStore()

// 任务状态管理
const taskStore = useTaskStore()

systemConfigStore.initConfig(false, true).catch(console.error)

// 版本信息
const versionInfo = ref({
  version: '1.1.0'
})

// 获取版本信息
const fetchVersionInfo = async () => {
  try {
    const response = await $fetch('/api/version') as any
    if (response.success) {
      versionInfo.value = response.data
    }
  } catch (error) {
    console.error('获取版本信息失败:', error)
  }
}

// 初始化版本信息和任务状态管理
onMounted(() => {
  fetchVersionInfo()

  // 启动任务状态自动更新
  taskStore.startAutoUpdate()

  // 确保在客户端配置被正确载入（防止SSR水合问题）
  setTimeout(async () => {
    try {
      await systemConfigStore.initConfig(true, true) // 强制刷新，防止SSR水合问题
    } catch (error) {
      console.error('Admin layout: onMounted 配置刷新失败', error)
    }
  }, 100) // 延迟100ms，确保组件渲染完成
})

// 组件销毁时清理任务状态管理
onBeforeUnmount(() => {
  // 停止任务状态自动更新
  taskStore.stopAutoUpdate()
  console.log('Admin layout: 任务状态自动更新已停止')
})

// 系统配置
const systemConfig = computed(() => {
  const config = systemConfigStore.config || {}
  return config
})

// 自动处理状态（确保布尔值）
const isAutoProcessEnabled = computed(() => {
  const value = systemConfig.value?.auto_process_ready_resources
  return value === true || value === 'true' || value === '1'
})

// 自动转存状态（确保布尔值）
const isAutoTransferEnabled = computed(() => {
  const value = systemConfig.value?.auto_transfer_enabled
  return value === true || value === 'true' || value === '1'
})

// 用户菜单状态
const showUserMenu = ref(false)

// 展开/折叠状态管理
const expandedGroups = ref({
  dataManagement: false,
  systemConfig: false,
  operation: false,
  statistics: false
})

// 切换分组展开/折叠状态
const toggleGroup = (groupName: string) => {
  expandedGroups.value[groupName as keyof typeof expandedGroups.value] = !expandedGroups.value[groupName as keyof typeof expandedGroups.value]
}

// 处理退出登录
const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 管理员菜单项
const userMenuItems = computed(() => [
  {
    to: '/admin/tasks',
    icon: 'fas fa-tasks',
    label: '任务列表',
    type: 'link'
  },
  {
    to: '/admin/accounts',
    icon: 'fas fa-user-shield',
    label: '账号管理',
    type: 'link'
  },
  {
    to: '/admin/system-config',
    icon: 'fas fa-cog',
    label: '系统配置',
    type: 'link'
  },
  {
    to: '/admin/version',
    icon: 'fas fa-code-branch',
    label: '版本信息',
    type: 'link'
  },
  {
    type: 'divider'
  },
  {
    type: 'button',
    icon: 'fas fa-sign-out-alt',
    label: '退出登录',
    action: handleLogout,
    className: 'text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300'
  }
])

// 仪表盘菜单项
const dashboardItems = ref([
  {
    to: '/admin',
    label: '仪表盘',
    icon: 'fas fa-tachometer-alt',
    active: (route: any) => route.path === '/admin'
  }
])

// 数据管理菜单项
const dataManagementItems = ref([
  {
    to: '/admin/resources',
    label: '资源管理',
    icon: 'fas fa-database',
    active: (route: any) => route.path.startsWith('/admin/resources')
  },
  {
    to: '/admin/ready-resources',
    label: '待处理资源',
    icon: 'fas fa-clock',
    active: (route: any) => route.path.startsWith('/admin/ready-resources')
  },

  {
    to: '/admin/tags',
    label: '标签管理',
    icon: 'fas fa-tags',
    active: (route: any) => route.path.startsWith('/admin/tags')
  },
  {
    to: '/admin/categories',
    label: '分类管理',
    icon: 'fas fa-folder',
    active: (route: any) => route.path.startsWith('/admin/categories')
  },
  {
    to: '/admin/accounts',
    label: '平台账号',
    icon: 'fas fa-user-shield',
    active: (route: any) => route.path.startsWith('/admin/accounts')
  },
  {
    to: '/admin/files',
    label: '文件管理',
    icon: 'fas fa-file-upload',
    active: (route: any) => route.path.startsWith('/admin/files')
  }
])

// 系统配置菜单项
const systemConfigItems = ref([
  {
    to: '/admin/site-config',
    label: '站点配置',
    icon: 'fas fa-globe',
    active: (route: any) => route.path.startsWith('/admin/site-config')
  },
  {
    to: '/admin/feature-config',
    label: '功能配置',
    icon: 'fas fa-sliders-h',
    active: (route: any) => route.path.startsWith('/admin/feature-config')
  },
  {
    to: '/admin/dev-config',
    label: '开发配置',
    icon: 'fas fa-code',
    active: (route: any) => route.path.startsWith('/admin/dev-config')
  },
  {
    to: '/admin/users',
    label: '用户管理',
    icon: 'fas fa-users',
    active: (route: any) => route.path.startsWith('/admin/users')
  }
])

// 运营管理菜单项
const operationItems = ref([
  {
    to: '/admin/data-transfer',
    label: '数据转存管理',
    icon: 'fas fa-exchange-alt',
    active: (route: any) => route.path.startsWith('/admin/data-transfer')
  },
  {
    to: '/admin/data-push',
    label: '数据推送',
    icon: 'fas fa-upload',
    active: (route: any) => route.path.startsWith('/admin/data-push')
  },
  {
    to: '/admin/bot',
    label: '机器人',
    icon: 'fas fa-robot',
    active: (route: any) => route.path.startsWith('/admin/bot')
  },
  {
    to: '/admin/seo',
    label: 'SEO',
    icon: 'fas fa-search',
    active: (route: any) => route.path.startsWith('/admin/seo')
  }
])

// 统计分析菜单项
const statisticsItems = ref([
  {
    to: '/admin/search-stats',
    label: '搜索统计',
    icon: 'fas fa-chart-line',
    active: (route: any) => route.path.startsWith('/admin/search-stats')
  },
  {
    to: '/admin/third-party-stats',
    label: '三方统计',
    icon: 'fas fa-chart-bar',
    active: (route: any) => route.path.startsWith('/admin/third-party-stats')
  }
])

// 自动展开当前页面所在的分组
const autoExpandCurrentGroup = () => {
  const currentPath = useRoute().path
  
  // 检查当前页面属于哪个分组并展开
  if (currentPath.startsWith('/admin/resources') || currentPath.startsWith('/admin/ready-resources') || currentPath.startsWith('/admin/tags') || currentPath.startsWith('/admin/categories') || currentPath.startsWith('/admin/accounts') || currentPath.startsWith('/admin/files')) {
    expandedGroups.value.dataManagement = true
  } else if (currentPath.startsWith('/admin/site-config') || currentPath.startsWith('/admin/feature-config') || currentPath.startsWith('/admin/dev-config') || currentPath.startsWith('/admin/users') || currentPath.startsWith('/admin/version')) {
    expandedGroups.value.systemConfig = true
  } else if (currentPath.startsWith('/admin/data-transfer') || currentPath.startsWith('/admin/seo') || currentPath.startsWith('/admin/data-push') || currentPath.startsWith('/admin/bot')) {
    expandedGroups.value.operation = true
  } else if (currentPath.startsWith('/admin/search-stats') || currentPath.startsWith('/admin/third-party-stats')) {
    expandedGroups.value.statistics = true
  }
}

// 监听路由变化，自动展开对应分组
watch(() => useRoute().path, (newPath) => {
  // 重置所有分组状态
  expandedGroups.value = {
    dataManagement: false,
    systemConfig: false,
    operation: false,
    statistics: false
  }
  
  // 根据新路径展开对应分组
  if (newPath.startsWith('/admin/resources') || newPath.startsWith('/admin/ready-resources') || newPath.startsWith('/admin/tags') || newPath.startsWith('/admin/categories') || newPath.startsWith('/admin/accounts') || newPath.startsWith('/admin/files')) {
    expandedGroups.value.dataManagement = true
  } else if (newPath.startsWith('/admin/site-config') || newPath.startsWith('/admin/feature-config') || newPath.startsWith('/admin/dev-config') || newPath.startsWith('/admin/users') || newPath.startsWith('/admin/version')) {
    expandedGroups.value.systemConfig = true
  } else if (newPath.startsWith('/admin/data-transfer') || newPath.startsWith('/admin/seo') || newPath.startsWith('/admin/data-push') || newPath.startsWith('/admin/bot')) {
    expandedGroups.value.operation = true
  } else if (newPath.startsWith('/admin/search-stats') || newPath.startsWith('/admin/third-party-stats')) {
    expandedGroups.value.statistics = true
  }
}, { immediate: true })

// 点击外部关闭菜单
onMounted(() => {
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('.relative')) {
      showUserMenu.value = false
    }
  })
})

// 导航到任务列表页面
const navigateToTasks = () => {
  router.push('/admin/tasks')
}
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
.main-content {
  height: calc(100vh - 85px);
}
</style> 