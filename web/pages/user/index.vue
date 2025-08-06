<template>
  <div class="space-y-8">
    <!-- 欢迎区域 -->
    <n-card class="bg-gradient-to-r from-blue-500 to-purple-600 text-white border-0">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold mb-2 text-white">
            欢迎回来，{{ userStore.user?.username || '用户' }}！
          </h1>
          <p class="text-blue-100 text-lg">
            这里是您的个人中心，您可以管理您的资源、收藏和历史记录。
          </p>
        </div>
        <div class="hidden md:block">
          <div class="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center">
            <i class="fas fa-user text-2xl text-white"></i>
          </div>
        </div>
      </div>
    </n-card>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- 我的资源 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">我的资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.resources || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="primary" @click="navigateTo('/user/resources')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- 收藏夹 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-heart text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">收藏夹</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.favorites || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="error" @click="navigateTo('/user/favorites')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- 浏览历史 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-history text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">浏览历史</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.history || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="success" @click="navigateTo('/user/history')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- 最近活动 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-clock text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">最近活动</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ userStats.recent || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="info" @click="navigateTo('/user/activity')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>
    </div>

    <!-- 快速操作 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 最近资源 -->
      <n-card title="最近资源" :bordered="false">
        <template #header-extra>
          <n-tag type="info" size="small">最新</n-tag>
        </template>
        
        <div v-if="recentResources.length === 0" class="text-center py-8">
          <n-empty description="暂无最近资源">
            <template #icon>
              <i class="fas fa-cloud text-gray-400 text-3xl"></i>
            </template>
            <template #extra>
              <n-button type="primary" @click="navigateTo('/')">
                去发现资源
              </n-button>
            </template>
          </n-empty>
        </div>
        
        <div v-else class="space-y-4">
          <n-list>
            <n-list-item v-for="resource in recentResources" :key="resource.id">
              <template #prefix>
                <n-avatar round size="small">
                  <i class="fas fa-file-alt"></i>
                </n-avatar>
              </template>
              
              <n-thing :title="resource.title" :description="resource.description">
                <template #avatar>
                  <n-avatar round size="small">
                    <i class="fas fa-file-alt"></i>
                  </n-avatar>
                </template>
              </n-thing>
              
              <template #suffix>
                <n-button text type="primary" @click="viewResource(resource)">
                  <template #icon>
                    <i class="fas fa-external-link-alt"></i>
                  </template>
                </n-button>
              </template>
            </n-list-item>
          </n-list>
        </div>
      </n-card>

      <!-- 快速操作 -->
      <n-card title="快速操作" :bordered="false">
        <template #header-extra>
          <n-tag type="success" size="small">快捷</n-tag>
        </template>
        
        <div class="grid grid-cols-2 gap-4">
          <n-button 
            quaternary 
            type="primary" 
            size="large" 
            class="h-24"
            @click="navigateTo('/user/profile')"
          >
            <template #icon>
              <i class="fas fa-user-edit text-2xl"></i>
            </template>
            <div class="flex flex-col items-center">
              <span class="text-sm font-medium">个人资料</span>
            </div>
          </n-button>

          <n-button 
            quaternary 
            type="success" 
            size="large" 
            class="h-24"
            @click="navigateTo('/user/settings')"
          >
            <template #icon>
              <i class="fas fa-cog text-2xl"></i>
            </template>
            <div class="flex flex-col items-center">
              <span class="text-sm font-medium">设置</span>
            </div>
          </n-button>

          <n-button 
            quaternary 
            type="info" 
            size="large" 
            class="h-24"
            @click="navigateTo('/')"
          >
            <template #icon>
              <i class="fas fa-search text-2xl"></i>
            </template>
            <div class="flex flex-col items-center">
              <span class="text-sm font-medium">搜索资源</span>
            </div>
          </n-button>

          <n-button 
            quaternary 
            type="warning" 
            size="large" 
            class="h-24"
            @click="exportData"
          >
            <template #icon>
              <i class="fas fa-download text-2xl"></i>
            </template>
            <div class="flex flex-col items-center">
              <span class="text-sm font-medium">导出数据</span>
            </div>
          </n-button>
        </div>
      </n-card>
    </div>

    <!-- 系统信息 -->
    <n-card title="账户信息" :bordered="false">
      <template #header-extra>
        <n-tag type="primary" size="small">账户</n-tag>
      </template>
      
      <n-descriptions :column="2" bordered>
        <n-descriptions-item label="用户名">
          <n-text>{{ userStore.user?.username }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="邮箱">
          <n-text>{{ userStore.user?.email || '未设置' }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="注册时间">
          <n-text>{{ formatDate(userStore.user?.created_at || '') }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="最后登录">
          <n-text>{{ formatDate(userStore.user?.last_login_at || '') }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="账户状态">
          <n-tag type="success" size="small">
            <template #icon>
              <i class="fas fa-check-circle"></i>
            </template>
            正常
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="用户角色">
          <n-tag type="primary" size="small">
            {{ userStore.user?.role === 'admin' ? '管理员' : '普通用户' }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="账户类型">
          <n-text>免费账户</n-text>
        </n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'user',
  ssr: false
})

// 用户状态管理
const userStore = useUserStore()

// 响应式数据
const userStats = ref({
  resources: 0,
  favorites: 0,
  history: 0,
  recent: 0
})

const recentResources = ref<any[]>([])

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 查看资源
const viewResource = (resource: any) => {
  // 这里可以跳转到资源详情页或打开链接
  if (resource.url) {
    window.open(resource.url, '_blank')
  }
}

// 导出数据
const exportData = () => {
  // 这里可以实现数据导出功能
  const notification = useNotification()
  notification.info({
    content: '导出功能开发中...',
    duration: 3000
  })
}

// 获取用户统计数据
const fetchUserStats = async () => {
  try {
    // 这里可以调用API获取用户统计数据
    // const response = await userApi.getUserStats()
    // userStats.value = response.data
    
    // 模拟数据
    userStats.value = {
      resources: 12,
      favorites: 8,
      history: 25,
      recent: 3
    }
  } catch (error) {
    console.error('获取用户统计数据失败:', error)
  }
}

// 获取最近资源
const fetchRecentResources = async () => {
  try {
    // 这里可以调用API获取最近资源
    // const response = await userApi.getRecentResources()
    // recentResources.value = response.data
    
    // 模拟数据
    recentResources.value = [
      {
        id: 1,
        title: '示例资源1',
        description: '这是一个示例资源描述',
        url: 'https://example.com'
      },
      {
        id: 2,
        title: '示例资源2',
        description: '这是另一个示例资源描述',
        url: 'https://example.com'
      }
    ]
  } catch (error) {
    console.error('获取最近资源失败:', error)
  }
}

// 页面加载时获取数据
onMounted(async () => {
  await fetchUserStats()
  await fetchRecentResources()
})

// 设置页面标题
useHead({
  title: '用户中心 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 