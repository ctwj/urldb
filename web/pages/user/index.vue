<template>
  <div class="space-y-8">
    <!-- 欢迎区域 -->
    <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg p-8 text-white">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold mb-2">
            欢迎回来，{{ userStore.user?.username || '用户' }}！
          </h1>
          <p class="text-blue-100 text-lg">
            这里是您的个人中心，您可以管理您的资源、收藏和历史记录。
          </p>
        </div>
        <div class="hidden md:block">
          <div class="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center">
            <i class="fas fa-user text-2xl"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- 我的资源 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">我的资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.resources || 0 }}</p>
          </div>
        </div>
        <NuxtLink
          to="/user/resources"
          class="mt-4 inline-flex items-center text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300"
        >
          查看详情
          <i class="fas fa-arrow-right ml-1"></i>
        </NuxtLink>
      </div>

      <!-- 收藏夹 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-heart text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">收藏夹</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.favorites || 0 }}</p>
          </div>
        </div>
        <NuxtLink
          to="/user/favorites"
          class="mt-4 inline-flex items-center text-sm text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300"
        >
          查看详情
          <i class="fas fa-arrow-right ml-1"></i>
        </NuxtLink>
      </div>

      <!-- 浏览历史 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-history text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">浏览历史</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.history || 0 }}</p>
          </div>
        </div>
        <NuxtLink
          to="/user/history"
          class="mt-4 inline-flex items-center text-sm text-green-600 dark:text-green-400 hover:text-green-800 dark:hover:text-green-300"
        >
          查看详情
          <i class="fas fa-arrow-right ml-1"></i>
        </NuxtLink>
      </div>

      <!-- 最近活动 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-clock text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">最近活动</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.recent || 0 }}</p>
          </div>
        </div>
        <NuxtLink
          to="/user/activity"
          class="mt-4 inline-flex items-center text-sm text-purple-600 dark:text-purple-400 hover:text-purple-800 dark:hover:text-purple-300"
        >
          查看详情
          <i class="fas fa-arrow-right ml-1"></i>
        </NuxtLink>
      </div>
    </div>

    <!-- 快速操作 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 最近资源 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
        <div class="p-6 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">最近资源</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">您最近浏览或收藏的资源</p>
        </div>
        <div class="p-6">
          <div v-if="recentResources.length === 0" class="text-center py-8">
            <i class="fas fa-cloud text-gray-400 text-3xl mb-4"></i>
            <p class="text-gray-500 dark:text-gray-400">暂无最近资源</p>
            <NuxtLink
              to="/"
              class="mt-4 inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              去发现资源
            </NuxtLink>
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="resource in recentResources"
              :key="resource.id"
              class="flex items-center space-x-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
            >
              <div class="flex-1">
                <h4 class="font-medium text-gray-900 dark:text-white">{{ resource.title }}</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">{{ resource.description }}</p>
              </div>
              <button
                @click="viewResource(resource)"
                class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300"
              >
                <i class="fas fa-external-link-alt"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 快速操作 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
        <div class="p-6 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">快速操作</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">常用功能快捷入口</p>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-2 gap-4">
            <NuxtLink
              to="/user/profile"
              class="flex flex-col items-center p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-colors"
            >
              <i class="fas fa-user-edit text-blue-600 dark:text-blue-400 text-2xl mb-2"></i>
              <span class="text-sm font-medium text-gray-900 dark:text-white">个人资料</span>
            </NuxtLink>

            <NuxtLink
              to="/user/settings"
              class="flex flex-col items-center p-4 bg-green-50 dark:bg-green-900/20 rounded-lg hover:bg-green-100 dark:hover:bg-green-900/30 transition-colors"
            >
              <i class="fas fa-cog text-green-600 dark:text-green-400 text-2xl mb-2"></i>
              <span class="text-sm font-medium text-gray-900 dark:text-white">设置</span>
            </NuxtLink>

            <NuxtLink
              to="/"
              class="flex flex-col items-center p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-colors"
            >
              <i class="fas fa-search text-purple-600 dark:text-purple-400 text-2xl mb-2"></i>
              <span class="text-sm font-medium text-gray-900 dark:text-white">搜索资源</span>
            </NuxtLink>

            <button
              @click="exportData"
              class="flex flex-col items-center p-4 bg-orange-50 dark:bg-orange-900/20 rounded-lg hover:bg-orange-100 dark:hover:bg-orange-900/30 transition-colors"
            >
              <i class="fas fa-download text-orange-600 dark:text-orange-400 text-2xl mb-2"></i>
              <span class="text-sm font-medium text-gray-900 dark:text-white">导出数据</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">账户信息</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400">您的账户基本信息</p>
      </div>
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h4 class="font-medium text-gray-900 dark:text-white mb-4">基本信息</h4>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">用户名</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ userStore.user?.username }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">邮箱</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ userStore.user?.email || '未设置' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">注册时间</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ formatDate(userStore.user?.created_at) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">最后登录</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ formatDate(userStore.user?.last_login_at) }}</span>
              </div>
            </div>
          </div>
          <div>
            <h4 class="font-medium text-gray-900 dark:text-white mb-4">账户状态</h4>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">账户状态</span>
                <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                  <i class="fas fa-check-circle mr-1"></i>
                  正常
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">用户角色</span>
                <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                  普通用户
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">账户类型</span>
                <span class="font-medium text-gray-900 dark:text-white">免费账户</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'user',
  ssr: false
})

// 用户状态管理
const userStore = useUserStore()

// 响应式数据
const userStats = ref({
  resources: 0,
  favorites: 0,
  history: 0,
  recent: 0
})

const recentResources = ref([])

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 查看资源
const viewResource = (resource: any) => {
  // 这里可以跳转到资源详情页或打开链接
  if (resource.url) {
    window.open(resource.url, '_blank')
  }
}

// 导出数据
const exportData = () => {
  // 这里可以实现数据导出功能
  alert('导出功能开发中...')
}

// 获取用户统计数据
const fetchUserStats = async () => {
  try {
    // 这里可以调用API获取用户统计数据
    // const response = await userApi.getUserStats()
    // userStats.value = response.data
    
    // 模拟数据
    userStats.value = {
      resources: 12,
      favorites: 8,
      history: 25,
      recent: 3
    }
  } catch (error) {
    console.error('获取用户统计数据失败:', error)
  }
}

// 获取最近资源
const fetchRecentResources = async () => {
  try {
    // 这里可以调用API获取最近资源
    // const response = await userApi.getRecentResources()
    // recentResources.value = response.data
    
    // 模拟数据
    recentResources.value = [
      {
        id: 1,
        title: '示例资源1',
        description: '这是一个示例资源描述',
        url: 'https://example.com'
      },
      {
        id: 2,
        title: '示例资源2',
        description: '这是另一个示例资源描述',
        url: 'https://example.com'
      }
    ]
  } catch (error) {
    console.error('获取最近资源失败:', error)
  }
}

// 页面加载时获取数据
onMounted(async () => {
  await fetchUserStats()
  await fetchRecentResources()
})

// 设置页面标题
useHead({
  title: '用户中心 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 