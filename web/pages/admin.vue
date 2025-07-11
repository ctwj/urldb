<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-3 sm:p-5">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在初始化管理后台</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <div class="bg-slate-800 text-white rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center">
        <div class="flex justify-between items-center mb-4">
          <h1 class="text-2xl sm:text-3xl font-bold">
            <NuxtLink to="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
              {{ systemConfig?.site_title || '网盘资源管理系统' }}
            </NuxtLink>
          </h1>
                      <div class="flex items-center gap-4">
              <div class="text-sm">
                <span>欢迎，{{ userStore.userInfo?.username }}</span>
                <span class="ml-2 px-2 py-1 bg-blue-600 rounded text-xs">{{ userStore.userInfo?.role }}</span>
              </div>
              <button 
                @click="handleLogout" 
                class="px-3 py-1 bg-red-600 hover:bg-red-700 rounded text-sm transition-colors"
              >
                退出登录
              </button>
            </div>
        </div>
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-4">
          <NuxtLink 
            to="/" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-home"></i> 返回首页
          </NuxtLink>
          <button 
            @click="showAddResourceModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-plus"></i> 添加资源
          </button>
        </nav>
      </div>

      <!-- 管理功能区域 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <!-- 资源管理 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-blue-100 rounded-lg">
              <i class="fas fa-cloud text-blue-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">资源管理</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">管理所有资源</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToResourceManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">查看所有资源</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddResourceModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加新资源</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 平台管理 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-green-100 rounded-lg">
              <i class="fas fa-server text-green-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">平台管理</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">管理网盘平台</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToPlatformManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理平台</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddPlatformModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加平台</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 分类管理 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-purple-100 rounded-lg">
              <i class="fas fa-folder text-purple-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">分类管理</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">管理资源分类</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToCategoryManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理分类</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddCategoryModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加分类</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 标签管理 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-orange-100 rounded-lg">
              <i class="fas fa-tags text-orange-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">标签管理</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">管理资源标签</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToTagManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理标签</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddTagModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加标签</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-red-100 rounded-lg">
              <i class="fas fa-chart-bar text-red-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">统计信息</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">系统统计数据</p>
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">总资源数</span>
              <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ stats?.total_resources || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">总浏览量</span>
              <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ stats?.total_views || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">分类数量</span>
              <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ stats?.total_categories || 0 }}</span>
            </div>
          </div>
        </div>

        <!-- 待处理资源 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-yellow-100 rounded-lg">
              <i class="fas fa-clock text-yellow-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">待处理资源</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">批量添加和管理</p>
            </div>
          </div>
          <div class="space-y-2">
            <NuxtLink to="/ready-resources" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理待处理资源</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </NuxtLink>
            <button @click="goToBatchAdd" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">批量添加资源</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 搜索统计 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-indigo-100 dark:bg-indigo-900 rounded-lg">
              <i class="fas fa-search text-indigo-600 dark:text-indigo-300 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">搜索统计</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">搜索量分析和热门关键词</p>
            </div>
          </div>
          <div class="space-y-2">
            <NuxtLink to="/search-stats" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">查看搜索统计</span>
                <i class="fas fa-chart-line text-gray-400 dark:text-gray-300"></i>
              </div>
            </NuxtLink>
            <button @click="goToHotKeywords" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">热门关键词</span>
                <i class="fas fa-fire text-gray-400 dark:text-gray-300"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 系统设置 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-gray-100 dark:bg-gray-700 rounded-lg">
              <i class="fas fa-cog text-gray-600 dark:text-gray-300 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">系统设置</h3>
              <p class="text-sm text-gray-600 dark:text-gray-400">系统配置</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToSystemSettings" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">系统配置</span>
                <i class="fas fa-chevron-right text-gray-400 dark:text-gray-300"></i>
              </div>
            </button>
            <button @click="goToUserManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700 dark:text-gray-200">用户管理</span>
                <i class="fas fa-users text-gray-400 dark:text-gray-300"></i>
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- 模态框组件 -->
      <ResourceModal v-if="showAddResourceModal" @close="showAddResourceModal = false" />
    </div>
  </div>
</template>

<script setup>
import ResourceModal from '~/components/ResourceModal.vue'

definePageMeta({
  middleware: 'auth'
})

// API
const { getSystemConfig } = useSystemConfigApi()

const router = useRouter()
const userStore = useUserStore()
const { $api } = useNuxtApp()

const user = ref(null)
const stats = ref(null)
const showAddResourceModal = ref(false)
const pageLoading = ref(true) // 添加页面加载状态
const systemConfig = ref(null) // 添加系统配置状态

// 页面元数据 - 移到变量声明之后
useHead({
  title: () => systemConfig.value?.site_title ? `${systemConfig.value.site_title} - 管理后台` : '管理后台 - 网盘资源管理系统',
  meta: [
    { 
      name: 'description', 
      content: () => systemConfig.value?.site_description || '网盘资源管理系统管理后台' 
    },
    { 
      name: 'keywords', 
      content: () => systemConfig.value?.keywords || '网盘,资源管理,管理后台' 
    },
    { 
      name: 'author', 
      content: () => systemConfig.value?.author || '系统管理员' 
    }
  ]
})

// 获取系统配置
const fetchSystemConfig = async () => {
  try {
    const response = await getSystemConfig()
    console.log('admin系统配置响应:', response)
    // 使用新的统一响应格式，直接使用response
    if (response) {
      systemConfig.value = response
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 检查认证状态
const checkAuth = () => {
  console.log('admin - checkAuth 开始')
  userStore.initAuth()
  
  console.log('admin - isAuthenticated:', userStore.isAuthenticated)
  console.log('admin - user:', userStore.userInfo)
  
  if (!userStore.isAuthenticated) {
    console.log('admin - 用户未认证，重定向到首页')
    router.push('/')
    return
  }
  
  console.log('admin - 用户已认证，继续')
}

// 获取统计信息
const fetchStats = async () => {
  try {
    const response = await $api.get('/stats')
    stats.value = response.data
  } catch (error) {
    console.error('获取统计信息失败:', error)
  }
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  router.push('/')
}

// 页面跳转方法
const goToResourceManagement = () => {
  // 实现资源管理页面跳转
}

const goToPlatformManagement = () => {
  // 实现平台管理页面跳转
}

const goToCategoryManagement = () => {
  // 实现分类管理页面跳转
}

const goToTagManagement = () => {
  // 实现标签管理页面跳转
}

const goToBatchAdd = () => {
  router.push('/ready-resources')
}

const goToSystemSettings = () => {
  router.push('/system-config')
}

const goToUserManagement = () => {
  router.push('/users')
}

const goToHotKeywords = () => {
  router.push('/search-stats')
}

// 页面加载时检查认证
onMounted(async () => {
  try {
    checkAuth()
    await Promise.all([
      fetchStats(),
      fetchSystemConfig()
    ])
  } catch (error) {
    console.error('admin页面初始化失败:', error)
  } finally {
    // 所有数据加载完成后，关闭加载状态
    pageLoading.value = false
  }
})
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 