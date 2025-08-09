<template>
  <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
    <div class="flex items-center justify-between px-6 py-4">
      <!-- 左侧：Logo和标题 -->
      <div class="flex items-center">
        <NuxtLink to="/admin" class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <i class="fas fa-shield-alt text-white text-sm"></i>
          </div>
          <div>
            <h1 class="text-xl font-bold text-gray-900 dark:text-white">管理后台</h1>
            <p class="text-xs text-gray-500 dark:text-gray-400">老九网盘资源数据库</p>
          </div>
        </NuxtLink>
      </div>

      <!-- 中间：状态信息 -->
      <div class="flex items-center space-x-4">
        <div class="flex items-center space-x-2">
          <div class="w-2 h-2 bg-green-500 rounded-full"></div>
          <span class="text-sm text-gray-600 dark:text-gray-300">系统正常</span>
        </div>
      </div>

      <!-- 右侧：用户菜单 -->
      <div class="flex items-center space-x-4">
        <NuxtLink to="/" class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
          <i class="fas fa-home text-lg"></i>
        </NuxtLink>
        <NuxtLink to="/admin-old" class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
          <i class="fas fa-arrow-left text-lg"></i>
        </NuxtLink>
        <div class="flex items-center space-x-2 cursor-pointer p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
          <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
            <i class="fas fa-user text-white text-sm"></i>
          </div>
          <div class="hidden md:block text-left">
            <p class="text-sm font-medium text-gray-900 dark:text-white">管理员</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">admin</p>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
// 简化的头部组件
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 