<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">浏览历史</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">查看浏览历史记录</p>
      </div>
      <div class="flex space-x-2">
        <n-button type="info" @click="handleRefresh">
          <i class="fas fa-sync-alt mr-2"></i>
          刷新
        </n-button>
        <n-button type="warning" @click="handleClearHistory">
          <i class="fas fa-trash mr-2"></i>
          清空历史
        </n-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-history text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总记录数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-clock text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日浏览</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.today || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-calendar text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本周浏览</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.thisWeek || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 历史记录列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">浏览记录</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">您的浏览历史记录</p>
      </div>
      
      <div class="p-6">
        <div class="text-center py-12">
          <i class="fas fa-history text-4xl text-gray-400 mb-4"></i>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无浏览记录</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">您还没有浏览过任何资源</p>
          <n-button type="primary" @click="handleGoToHome">
            <i class="fas fa-home mr-2"></i>
            去首页浏览资源
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 页面元数据
definePageMeta({
  layout: 'user',
  title: '浏览历史'
})

// 响应式数据
const stats = ref({
  total: 0,
  today: 0,
  thisWeek: 0
})

// 处理刷新
const handleRefresh = () => {
  // TODO: 实现刷新功能
  console.log('刷新浏览历史')
}

// 处理清空历史
const handleClearHistory = () => {
  // TODO: 实现清空历史功能
  console.log('清空浏览历史')
}

// 处理跳转到首页
const handleGoToHome = () => {
  navigateTo('/')
}

// 页面加载时获取数据
onMounted(() => {
  // TODO: 获取用户浏览历史数据
  console.log('加载用户浏览历史数据')
})
</script> 