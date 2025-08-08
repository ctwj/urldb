<template>
  <div class="p-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">数据转存管理</h1>
      <p class="text-gray-600 dark:text-gray-400 mt-2">管理数据转存任务和配置</p>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="text-center py-12">
        <div class="text-gray-400 dark:text-gray-500 mb-4">
          <i class="fas fa-exchange-alt text-4xl"></i>
        </div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能开发中</h3>
        <p class="text-gray-500 dark:text-gray-400">数据转存管理功能正在开发中，敬请期待...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 数据转存管理页面
definePageMeta({
  layout: 'admin',
  middleware: ['auth']
})
</script> 