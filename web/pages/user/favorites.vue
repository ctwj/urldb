<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">收藏夹</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">查看和管理收藏的资源</p>
      </div>
      <n-button type="info" @click="handleRefresh">
        <i class="fas fa-sync-alt mr-2"></i>
        刷新
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-heart text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">收藏总数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-check text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">有效收藏</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.valid || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-clock text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本月新增</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.thisMonth || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 收藏列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">收藏列表</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">您收藏的所有资源</p>
      </div>
      
      <div class="p-6">
        <div class="text-center py-12">
          <i class="fas fa-heart text-4xl text-gray-400 mb-4"></i>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无收藏</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">您还没有收藏任何资源</p>
          <n-button type="primary" @click="handleGoToHome">
            <i class="fas fa-home mr-2"></i>
            去首页发现资源
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
  title: '收藏夹'
})

// 响应式数据
const stats = ref({
  total: 0,
  valid: 0,
  thisMonth: 0
})

// 处理刷新
const handleRefresh = () => {
  // TODO: 实现刷新功能
  console.log('刷新收藏数据')
}

// 处理跳转到首页
const handleGoToHome = () => {
  navigateTo('/')
}

// 页面加载时获取数据
onMounted(() => {
  // TODO: 获取用户收藏数据
  console.log('加载用户收藏数据')
})
</script> 