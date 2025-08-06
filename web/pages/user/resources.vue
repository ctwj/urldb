<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">我的资源</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">管理您的个人资源</p>
      </div>
      <n-button type="primary" @click="handleAddResource">
        <i class="fas fa-plus mr-2"></i>
        添加资源
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总资源数</p>
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
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">有效资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.valid || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-exclamation-triangle text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">无效资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.invalid || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 资源列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">资源列表</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">您添加的所有资源</p>
      </div>
      
      <div class="p-6">
        <div class="text-center py-12">
          <i class="fas fa-cloud text-4xl text-gray-400 mb-4"></i>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无资源</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">您还没有添加任何资源</p>
          <n-button type="primary" @click="handleAddResource">
            <i class="fas fa-plus mr-2"></i>
            添加第一个资源
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
  title: '我的资源'
})

// 响应式数据
const stats = ref({
  total: 0,
  valid: 0,
  invalid: 0
})

// 处理添加资源
const handleAddResource = () => {
  // TODO: 实现添加资源功能
  console.log('添加资源')
}

// 页面加载时获取数据
onMounted(() => {
  // TODO: 获取用户资源数据
  console.log('加载用户资源数据')
})
</script> 