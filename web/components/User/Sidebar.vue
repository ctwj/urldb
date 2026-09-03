<template>
  <aside class="w-64 bg-white dark:bg-gray-800 shadow-sm border-r border-gray-200 dark:border-gray-700 min-h-screen flex flex-col">
    <!-- 品牌区 -->
    <div class="px-6 pt-6 pb-4">
      <NuxtLink to="/" class="flex items-center gap-3 cursor-pointer group">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-sm group-hover:shadow-blue-200 dark:group-hover:shadow-blue-900/50 transition-shadow duration-200">
          <i class="fas fa-cloud text-white text-sm"></i>
        </div>
        <div>
          <p class="text-sm font-semibold text-gray-900 dark:text-white leading-tight">用户中心</p>
          <p class="text-xs text-gray-400 dark:text-gray-500 leading-tight mt-0.5">老九网盘资源库</p>
        </div>
      </NuxtLink>
    </div>

    <nav class="flex-1 px-3">
      <div class="space-y-1">
        <!-- 导航菜单 -->
        <NuxtLink
          v-for="item in navigationItems"
          :key="item.to"
          :to="item.to"
          class="relative flex items-center px-3 py-2.5 rounded-lg text-sm transition-colors duration-200 cursor-pointer group"
          :class="item.active($route)
            ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 font-medium'
            : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700/50 hover:text-gray-900 dark:hover:text-gray-200'"
        >
          <!-- active 左侧指示条 -->
          <span
            v-if="item.active($route)"
            class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-5 rounded-r-full bg-blue-500"
            aria-hidden="true"
          ></span>
          <i
            :class="[
              item.icon,
              'w-5 text-center mr-3 text-base transition-transform duration-200',
              item.active($route) ? 'text-blue-500 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300'
            ]"
          ></i>
          <span>{{ item.label }}</span>
        </NuxtLink>

        <div class="border-t border-gray-100 dark:border-gray-700 my-3"></div>

        <NuxtLink
          to="/"
          class="flex items-center px-3 py-2.5 rounded-lg text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700/50 hover:text-gray-900 dark:hover:text-gray-200 rounded-lg transition-colors duration-200 cursor-pointer group"
        >
          <i class="fas fa-arrow-left w-5 text-center mr-3 text-base text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300"></i>
          <span>返回网站首页</span>
        </NuxtLink>
      </div>
    </nav>

    <!-- 底部用户信息卡 -->
    <div class="p-3 border-t border-gray-100 dark:border-gray-700">
      <div class="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors duration-200">
        <div class="w-9 h-9 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center flex-shrink-0">
          <i class="fas fa-user text-white text-xs"></i>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white truncate">{{ userInfo.username }}</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">{{ userInfo.isAdmin ? '管理员' : '普通用户' }}</p>
        </div>
        <button
          type="button"
          class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-200 cursor-pointer"
          aria-label="退出登录"
          title="退出登录"
          @click="handleLogout"
        >
          <i class="fas fa-sign-out-alt text-sm"></i>
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useUserLayout } from '~/composables/useUserLayout'

// 使用用户布局组合式函数
const { getNavigationItems, getUserInfo, handleLogout } = useUserLayout()

// 获取导航菜单项
const navigationItems = computed(() => getNavigationItems())

// 当前用户信息
const userInfo = computed(() => getUserInfo())
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style>
