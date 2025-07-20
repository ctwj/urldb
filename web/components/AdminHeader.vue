<template>
  <div class="bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 text-center relative">
    <!-- 页面标题和面包屑 -->
    <div class="mb-4">
      <h1 class="text-2xl sm:text-3xl font-bold mb-2">
        <NuxtLink to="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
          {{ systemConfig?.site_title || '老九网盘资源数据库' }}
        </NuxtLink>
      </h1>
      <!-- 面包屑导航 -->
      <div v-if="currentPageTitle && currentPageTitle !== '管理后台'" class="absolute left-4 bottom-4 flex items-center justify-start text-sm text-white/80">
        <NuxtLink to="/admin" class="hover:text-white transition-colors">
          <i class="fas fa-home mr-1"></i>管理后台
        </NuxtLink>
        <i class="fas fa-angle-right mx-2 text-white/60"></i>
        <span class="text-white">
          <i :class="currentPageIcon + ' mr-1'"></i>
          {{ currentPageTitle }}
        </span>
      </div>
      <!-- 页面描述 -->
      <!-- <div v-if="currentPageDescription && currentPageTitle !== '管理后台'" class="text-xs text-white/60 mt-1">
        {{ currentPageDescription }}
      </div> -->
    </div>

    <div class="absolute left-4 top-4 flex items-center gap-2">
      <NuxtLink to="/" class="sm:flex">
          <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
            <i class="fas fa-home text-xs"></i> 前端首页
          </n-button>
        </NuxtLink>
    </div>
    
    <!-- 右上角用户信息和操作按钮 -->
    <div class="absolute right-4 top-4 flex items-center gap-2">
      <!-- 用户信息 -->
      <div v-if="userStore.isAuthenticated" class="hidden sm:flex items-center gap-2">
        <span class="text-sm text-white/80">欢迎，{{ userStore.user?.username || '管理员' }}</span>
        <span class="px-2 py-1 bg-blue-600/80 rounded text-xs text-white">{{ userStore.user?.role || 'admin' }}</span>
      </div>
      
      <!-- 操作按钮 -->
      <div class="flex gap-1">
        <button 
          v-if="userStore.isAuthenticated"
          @click="logout" 
          class="sm:flex"
        >
          <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
            <i class="fas fa-sign-out-alt text-xs"></i> 退出
          </n-button>
        </button>
      </div>
    </div>
    
    <!-- 移动端用户信息 -->
    <div v-if="userStore.isAuthenticated" class="sm:hidden mt-4 text-sm text-white/80">
      <span>欢迎，{{ userStore.user?.username || '管理员' }}</span>
      <span class="ml-2 px-2 py-1 bg-blue-600/80 rounded text-xs text-white">{{ userStore.user?.role || 'admin' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  title: '管理后台'
})

// 用户状态管理
const userStore = useUserStore()
const router = useRouter()

// 页面配置
const route = useRoute()
const pageConfig = computed(() => {
  const configs: Record<string, { title: string; icon: string; description?: string }> = {
    '/admin': { title: '管理后台', icon: 'fas fa-tachometer-alt', description: '系统管理总览' },
    '/admin/users': { title: '用户管理', icon: 'fas fa-users', description: '管理系统用户' },
    '/admin/categories': { title: '分类管理', icon: 'fas fa-folder', description: '管理资源分类' },
    '/admin/tags': { title: '标签管理', icon: 'fas fa-tags', description: '管理资源标签' },
    '/admin/system-config': { title: '系统配置', icon: 'fas fa-cog', description: '系统参数设置' },
    '/admin/resources': { title: '资源管理', icon: 'fas fa-database', description: '管理网盘资源' },
    '/admin/cks': { title: '平台账号管理', icon: 'fas fa-key', description: '管理第三方平台账号' },
    '/admin/ready-resources': { title: '待处理资源', icon: 'fas fa-clock', description: '批量处理资源' },
    '/admin/search-stats': { title: '搜索统计', icon: 'fas fa-chart-bar', description: '搜索数据分析' },
    '/admin/hot-dramas': { title: '热播剧管理', icon: 'fas fa-film', description: '管理热门剧集' },
    '/monitor': { title: '系统监控', icon: 'fas fa-desktop', description: '系统性能监控' },
    '/add-resource': { title: '添加资源', icon: 'fas fa-plus', description: '添加新资源' },
    '/api-docs': { title: 'API文档', icon: 'fas fa-book', description: '接口文档说明' },
    '/admin/version': { title: '版本信息', icon: 'fas fa-code-branch', description: '系统版本详情' }
  }
  return configs[route.path] || { title: props.title, icon: 'fas fa-cog', description: '管理页面' }
})

const currentPageTitle = computed(() => pageConfig.value.title)
const currentPageIcon = computed(() => pageConfig.value.icon)
const currentPageDescription = computed(() => pageConfig.value.description)

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('systemConfig', 
  () => $fetch('/api/system-config')
)

const systemConfig = computed(() => (systemConfigData.value as any)?.data || { site_title: '老九网盘资源数据库' })

// 退出登录
const logout = async () => {
  await userStore.logout()
  await router.push('/login')
}
</script>

<style scoped>
/* 确保样式与首页完全一致 */
</style> 