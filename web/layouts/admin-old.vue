<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在初始化管理后台</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 管理页面头部 -->
    <div class="p-3 sm:p-5">
      <AdminHeader :title="pageTitle" />
    </div>
    
    <!-- 主要内容区域 -->
    <div class="p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
        <ClientOnly>
          <n-notification-provider>
            <n-dialog-provider>
              <!-- 页面内容插槽 -->
              <slot />
            </n-dialog-provider>
          </n-notification-provider>
        </ClientOnly>
      </div>
    </div>

    <!-- 页脚 -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import { useSystemConfigStore } from '~/stores/systemConfig'
import { useUserLayout } from '~/composables/useUserLayout'

// 使用用户布局组合式函数
const { checkAuth, checkPermission } = useUserLayout()

// 页面加载状态
const pageLoading = ref(false)

// 页面标题
const route = useRoute()
const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/admin-old': '管理后台',
    '/admin-old/users': '用户管理',
    '/admin-old/categories': '分类管理',
    '/admin-old/tags': '标签管理',
    '/admin-old/tasks': '任务管理',
    '/admin-old/system-config': '系统配置',
    '/admin-old/resources': '资源管理',
    '/admin-old/cks': '平台账号管理',
    '/admin-old/ready-resources': '待处理资源',
    '/admin-old/search-stats': '搜索统计',
    '/admin-old/hot-dramas': '热播剧管理',
    '/admin-old/monitor': '系统监控',
    '/admin-old/add-resource': '添加资源',
    '/admin-old/api-docs': 'API文档',
    '/admin-old/version': '版本信息'
  }
  return titles[route.path] || '管理后台'
})

// 监听路由变化，显示加载状态
watch(() => route.path, () => {
  pageLoading.value = true
  setTimeout(() => {
    pageLoading.value = false
  }, 300)
})

const systemConfigStore = useSystemConfigStore()
onMounted(() => {
  // 检查用户认证和权限
  if (!checkAuth()) {
    return
  }
  
  // 检查是否为管理员
  if (!checkPermission('admin')) {
    return
  }
  
  systemConfigStore.initConfig()
  pageLoading.value = true
  setTimeout(() => {
    pageLoading.value = false
  }, 300)
})
</script>

<style scoped>
/* 管理后台专用样式 */
</style> 