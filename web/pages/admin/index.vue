<template>
  <!-- 管理功能区域 -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
    <!-- 资源管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-blue-100 rounded-lg">
          <i class="fas fa-cloud text-blue-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">资源管理</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">管理所有资源</p>
        </div>
      </div>
      <div class="space-y-2">
        <NuxtLink to="/admin/resources" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">查看所有资源</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </NuxtLink>
        <NuxtLink to="/admin/add-resource" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">批量添加资源</span>
            <i class="fas fa-plus text-gray-400"></i>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- 平台管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-green-100 rounded-lg">
          <i class="fas fa-server text-green-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">平台管理</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">系统支持的网盘平台</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex flex-wrap gap-1 w-full text-left rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer">
          <div v-for="pan in pans" :key="pan.id"  class="h-6 px-1 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
            <span v-html="pan.icon"></span>&nbsp;{{ pan.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- 第三方平台账号管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-teal-100 rounded-lg">
          <i class="fas fa-key text-teal-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">平台账号管理</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">管理第三方平台账号</p>
        </div>
      </div>
      <div class="space-y-2">
        <NuxtLink to="/admin/cks" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理账号</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </NuxtLink>
        <NuxtLink to="/admin/cks" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加账号</span>
            <i class="fas fa-plus text-gray-400"></i>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- 分类管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-purple-100 rounded-lg">
          <i class="fas fa-folder text-purple-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">分类管理</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">管理资源分类</p>
        </div>
      </div>
      <div class="space-y-2">
        <button @click="goToCategoryManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理分类</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </button>
        <button @click="goToAddCategory" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加分类</span>
            <i class="fas fa-plus text-gray-400"></i>
          </div>
        </button>
      </div>
    </div>

    <!-- 标签管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-orange-100 rounded-lg">
          <i class="fas fa-tags text-orange-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">标签管理</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">管理资源标签</p>
        </div>
      </div>
      <div class="space-y-2">
        <button @click="goToTagManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理标签</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </button>
        <button @click="goToAddTag" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">添加标签</span>
            <i class="fas fa-plus text-gray-400"></i>
          </div>
        </button>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-red-100 rounded-lg">
          <i class="fas fa-chart-bar text-red-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">统计信息</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">系统统计数据</p>
        </div>
      </div>
      <div class="space-y-3">
        <NuxtLink to="/admin/search-stats" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">搜索统计</span>
            <i class="fas fa-chart-line text-gray-400"></i>
          </div>
        </NuxtLink>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600 dark:text-gray-400">总资源数</span>
          <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ stats?.total_resources || 0 }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600 dark:text-gray-400">总浏览量</span>
          <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ stats?.total_views || 0 }}</span>
        </div>
      </div>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-yellow-100 rounded-lg">
          <i class="fas fa-clock text-yellow-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">待处理资源</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">批量添加和管理</p>
        </div>
      </div>
      <div class="space-y-2">
        <NuxtLink to="/admin/ready-resources" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">管理待处理资源</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </NuxtLink>
        <NuxtLink to="/admin/ready-resources" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">批量处理</span>
            <i class="fas fa-tasks text-gray-400"></i>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- 系统配置 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-indigo-100 rounded-lg">
          <i class="fas fa-cog text-indigo-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">系统配置</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">系统参数设置</p>
        </div>
      </div>
      <div class="space-y-2">
        <NuxtLink to="/admin/users" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">用户管理</span>
            <i class="fas fa-users text-gray-400"></i>
          </div>
        </NuxtLink>
        <NuxtLink to="/admin/system-config" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">系统设置</span>
            <i class="fas fa-chevron-right text-gray-400"></i>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- 版本信息 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-green-100 rounded-lg">
          <i class="fas fa-code-branch text-green-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">版本信息</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">系统版本和文档</p>
        </div>
      </div>
      <div class="space-y-2">
        <NuxtLink to="/admin/version" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors block">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-200">版本信息</span>
            <i class="fas fa-code-branch text-gray-400"></i>
          </div>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

// 用户状态管理
const userStore = useUserStore()

// 统计数据
import { useStatsApi, usePanApi } from '~/composables/useApi'

const statsApi = useStatsApi()
const panApi = usePanApi()

const { data: statsData } = await useAsyncData('stats', () => statsApi.getStats())
const stats = computed(() => (statsData.value as any) || {})

// 平台数据
const { data: pansData } = await useAsyncData('pans', () => panApi.getPans())
const pans = computed(() => (pansData.value as any) || [])

// 分类管理相关
const goToCategoryManagement = () => {
  navigateTo('/admin/categories')
}

const goToAddCategory = () => {
  navigateTo('/admin/categories')
}

// 标签管理相关
const goToTagManagement = () => {
  navigateTo('/admin/tags')
}

const goToAddTag = () => {
  navigateTo('/admin/tags')
}

// 页面加载时检查用户权限
onMounted(() => {
  if (!userStore.isAuthenticated) {
    navigateTo('/login')
  }
})
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 