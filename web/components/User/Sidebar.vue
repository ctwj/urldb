<template>
  <aside class="w-64 bg-white dark:bg-gray-800 shadow-sm border-r border-gray-200 dark:border-gray-700 min-h-screen">
    <nav class="mt-8">
      <div class="px-4 space-y-2">
        <!-- 导航菜单 -->
        <NuxtLink
          v-for="item in navigationItems"
          :key="item.to"
          :to="item.to"
          class="flex items-center px-4 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
          :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': item.active($route) }"
        >
          <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
          <span>{{ item.label }}</span>
        </NuxtLink>

        <div class="border-t border-gray-200 dark:border-gray-700 my-4"></div>

        <NuxtLink
          to="/"
          class="flex items-center px-4 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
        >
          <i class="fas fa-arrow-left w-5 h-5 mr-3"></i>
          <span>返回首页</span>
        </NuxtLink>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { useUserLayout } from '~/composables/useUserLayout'

// 使用用户布局组合式函数
const { getNavigationItems } = useUserLayout()

// 获取导航菜单项
const navigationItems = computed(() => getNavigationItems())
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 