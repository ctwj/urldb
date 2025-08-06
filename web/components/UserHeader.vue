<template>
  <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- 左侧 Logo 和标题 -->
        <div class="flex items-center">
          <NuxtLink to="/user" class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <i class="fas fa-user text-white text-sm"></i>
            </div>
            <span class="text-xl font-bold text-gray-900 dark:text-white">用户中心</span>
          </NuxtLink>
        </div>

        <!-- 右侧用户信息和操作 -->
        <div class="flex items-center space-x-4">
          <!-- 用户信息 -->
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
              <i class="fas fa-user text-blue-600 dark:text-blue-400 text-sm"></i>
            </div>
            <div class="hidden md:block">
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ userStore.user?.username || '用户' }}
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">普通用户</p>
            </div>
          </div>

          <!-- 下拉菜单 -->
          <div class="relative">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center space-x-2 text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors"
            >
              <i class="fas fa-chevron-down text-xs"></i>
            </button>

            <!-- 下拉菜单内容 -->
            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 z-50 border border-gray-200 dark:border-gray-700"
            >
              <template v-for="item in userMenuItems" :key="item.label || item.type">
                <!-- 链接菜单项 -->
                <NuxtLink
                  v-if="item.type === 'link' && item.to"
                  :to="item.to"
                  class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </NuxtLink>
                
                <!-- 按钮菜单项 -->
                <button
                  v-else-if="item.type === 'button'"
                  @click="item.action"
                  class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
                  :class="item.className || 'text-gray-700 dark:text-gray-300'"
                >
                  <i :class="item.icon + ' mr-2'"></i>
                  {{ item.label }}
                </button>
                
                <!-- 分割线 -->
                <div
                  v-else-if="item.type === 'divider'"
                  class="border-t border-gray-200 dark:border-gray-700 my-1"
                ></div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useUserLayout } from '~/composables/useUserLayout'

// 用户状态管理
const userStore = useUserStore()

// 使用用户布局组合式函数
const { getUserMenuItems } = useUserLayout()

// 用户菜单状态
const showUserMenu = ref(false)

// 获取用户菜单项
const userMenuItems = computed(() => getUserMenuItems())

// 点击外部关闭菜单
onMounted(() => {
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('.relative')) {
      showUserMenu.value = false
    }
  })
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 