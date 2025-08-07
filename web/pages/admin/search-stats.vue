<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">搜索统计</h1>
        <p class="text-gray-600 dark:text-gray-400">查看系统搜索统计数据</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <n-card>
        <div class="flex items-center">
          <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-search text-blue-600 dark:text-blue-400"></i>
          </div>
          <div class="ml-3">
            <p class="text-sm text-gray-600 dark:text-gray-400">总资源数</p>
            <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ stats.total_resources || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-users text-green-600 dark:text-green-400"></i>
          </div>
          <div class="ml-3">
            <p class="text-sm text-gray-600 dark:text-gray-400">总用户数</p>
            <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ stats.total_users || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
            <i class="fas fa-chart-line text-yellow-600 dark:text-yellow-400"></i>
          </div>
          <div class="ml-3">
            <p class="text-sm text-gray-600 dark:text-gray-400">总浏览量</p>
            <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ stats.total_views || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-calendar text-purple-600 dark:text-purple-400"></i>
          </div>
          <div class="ml-3">
            <p class="text-sm text-gray-600 dark:text-gray-400">今日更新</p>
            <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ stats.today_updates || 0 }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 系统状态 -->
    <n-card>
      <template #header>
        <span class="text-lg font-semibold">系统状态</span>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">分类统计</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">总分类数: {{ stats.total_categories || 0 }}</p>
        </div>

        <div class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">标签统计</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">总标签数: {{ stats.total_tags || 0 }}</p>
        </div>

        <div class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">平台统计</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">总平台数: {{ stats.total_platforms || 0 }}</p>
        </div>

        <div class="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">热播剧统计</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">总热播剧数: {{ stats.total_hot_dramas || 0 }}</p>
        </div>
      </div>
    </n-card>

    <!-- 最近活动 -->
    <n-card>
      <template #header>
        <span class="text-lg font-semibold">最近活动</span>
      </template>

      <div class="space-y-4">
        <div class="flex items-center space-x-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
            <i class="fas fa-plus text-blue-600 dark:text-blue-400 text-sm"></i>
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-900 dark:text-white">新资源添加</p>
            <p class="text-xs text-gray-600 dark:text-gray-400">系统正常运行中</p>
          </div>
        </div>

        <div class="flex items-center space-x-4 p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
          <div class="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-full flex items-center justify-center">
            <i class="fas fa-check text-green-600 dark:text-green-400 text-sm"></i>
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-900 dark:text-white">自动处理</p>
            <p class="text-xs text-gray-600 dark:text-gray-400">待处理资源自动处理中</p>
          </div>
        </div>

        <div class="flex items-center space-x-4 p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
          <div class="w-8 h-8 bg-yellow-100 dark:bg-yellow-900 rounded-full flex items-center justify-center">
            <i class="fas fa-sync text-yellow-600 dark:text-yellow-400 text-sm"></i>
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-900 dark:text-white">数据同步</p>
            <p class="text-xs text-gray-600 dark:text-gray-400">系统数据同步正常</p>
          </div>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin' as any
})

// 使用API
const { useStatsApi } = await import('~/composables/useApi')
const statsApi = useStatsApi()

// 响应式数据
const stats = ref<any>({})

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await statsApi.getStats() as any
    stats.value = response.data || {}
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 初始化数据
onMounted(() => {
  fetchStats()
})

// 刷新数据
const refreshData = () => {
  fetchStats()
}
</script> 