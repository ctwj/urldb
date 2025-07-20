  <template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在初始化系统</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <div class="bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4">
          <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
            {{ systemConfig?.site_title || '网盘资源数据库' }}
          </a>
        </h1>
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
          <NuxtLink to="/hot-dramas" class="hidden sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-film text-xs"></i> 热播剧
            </n-button>
          </NuxtLink>
          <NuxtLink to="/monitor" class="hidden sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-chart-line text-xs"></i> 系统监控
            </n-button>
          </NuxtLink>
          <NuxtLink to="/api-docs" class="hidden sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-book text-xs"></i> API文档
            </n-button>
          </NuxtLink>
          <NuxtLink v-if="authInitialized && !userStore.isAuthenticated" to="/login" class="sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-sign-in-alt text-xs"></i> 登录
            </n-button>
          </NuxtLink>
          <NuxtLink v-if="authInitialized && userStore.isAuthenticated" to="/admin" class="hidden sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-user-shield text-xs"></i> 管理后台
            </n-button>
          </NuxtLink>
        </nav>
      </div>

      <!-- 搜索区域 -->
      <div class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
        <div class="relative">
          <input 
            v-model="searchQuery"
            @keyup="handleSearch"
            type="text" 
            class="w-full px-4 py-3 rounded-full border-2 border-gray-300 dark:border-gray-700 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200 dark:bg-gray-900 dark:text-gray-100 dark:placeholder-gray-500 transition-all"
            placeholder="输入文件名或链接进行搜索..."
          />
          <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
            <i class="fas fa-search text-gray-400"></i>
          </div>
        </div>
        
        <!-- 平台类型筛选 -->
        <div class="mt-3 flex flex-wrap gap-2" id="platformFilters">
          <button 
            class="px-2 py-1 text-xs rounded-full bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 active-filter" 
            @click="filterByPlatform('')"
          >
            全部
          </button>
          <button 
            v-for="platform in platforms" 
            :key="platform.id"
            class="px-2 py-1 text-xs rounded-full bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
            @click="filterByPlatform(platform.id)"
          >
            <span v-html="platform.icon"></span> {{ platform.name }}
          </button>
        </div>
        
        <!-- 统计信息 -->
        <div class="flex justify-between mt-3 text-sm text-gray-600 dark:text-gray-300 px-2">
          <div class="flex items-center">
            <i class="fas fa-calendar-day text-pink-600 mr-1"></i>
            今日更新: <span class="font-medium text-pink-600 ml-1 count-up" :data-target="safeStats?.today_updates || 0">0</span>
          </div>
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 mr-1"></i>
            总资源数: <span class="font-medium text-blue-600 ml-1 count-up" :data-target="safeStats?.total_resources || 0">0</span>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg shadow">
        <table class="w-full min-w-full table-fixed">
          <thead>
            <tr class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100">
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm w-1/2 sm:w-4/6">
                <div class="flex items-center">
                  <i class="fas fa-cloud mr-1 text-gray-300"></i> 文件名
                </div>
              </th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-1/6">链接</th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-1/6">更新时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="safeLoading" class="text-center py-8">
              <td colspan="1" class="text-gray-500 dark:text-gray-400 sm:hidden">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
              <td colspan="3" class="text-gray-500 dark:text-gray-400 hidden sm:table-cell">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
            </tr>
            <tr v-else-if="safeResources.length === 0" class="text-center py-8">
              <td colspan="1" class="text-gray-500 dark:text-gray-400 sm:hidden">暂无数据</td>
              <td colspan="3" class="text-gray-500 dark:text-gray-400 hidden sm:table-cell">暂无数据</td>
            </tr>
            <tr 
              v-for="(resource, index) in safeResources" 
              :key="resource.id"
              :class="isUpdatedToday(resource.updated_at) ? 'hover:bg-pink-50 dark:hover:bg-pink-900 bg-pink-50/30 dark:bg-pink-900/30' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
              :data-index="index"
            >
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm w-1/2 sm:w-2/5">
                <div class="flex items-start">
                  <span class="mr-2 flex-shrink-0" v-html="getPlatformIcon(resource.pan_id || 0)"></span>
                  <span class="break-words">{{ resource.title }}</span>
                </div>
                <div class="sm:hidden mt-1 space-y-1">
                  <!-- 移动端显示更新时间 -->
                  <div class="text-xs text-gray-500" :title="resource.updated_at">
                    <span v-html="formatRelativeTime(resource.updated_at)"></span>
                  </div>
                  <!-- 移动端显示链接按钮 -->
                  <button 
                    class="text-blue-600 hover:text-blue-800 text-xs flex items-center gap-1 show-link-btn" 
                    @click="toggleLink(resource)"
                  >
                    <i class="fas fa-eye"></i> 显示链接
                  </button>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm hidden sm:table-cell w-1/5">
                <button 
                  class="text-blue-600 hover:text-blue-800 flex items-center gap-1 show-link-btn" 
                  @click="toggleLink(resource)"
                >
                  <i class="fas fa-eye"></i> 显示链接
                </button>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm text-gray-500 hidden sm:table-cell w-2/5" :title="resource.updated_at">
                <span v-html="formatRelativeTime(resource.updated_at)"></span>
              </td>
            </tr>
          </tbody>
        </table>
        
        <!-- 加载更多按钮 -->
        <div v-if="hasMoreData && !safeLoading" class="text-center py-4">
          <button 
            @click="loadMore"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
          >
            加载更多
          </button>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex flex-wrap justify-center gap-1 sm:gap-2 my-4 sm:my-8 px-2">
        <button 
          v-if="currentPage > 1"
          @click="goToPage(currentPage - 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          <i class="fas fa-chevron-left mr-1"></i> 上一页
        </button>
        
        <button 
          @click="goToPage(1)"
          :class="currentPage === 1 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          1
        </button>
        
        <button 
          v-if="totalPages > 1"
          @click="goToPage(2)"
          :class="currentPage === 2 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          2
        </button>
        
        <span v-if="currentPage > 2" class="px-2 py-1 sm:px-3 sm:py-2 text-gray-500 text-sm">...</span>
        
        <button 
          v-if="currentPage !== 1 && currentPage !== 2 && currentPage > 2"
          class="bg-slate-800 text-white px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          {{ currentPage }}
        </button>
        
        <button 
          v-if="currentPage < totalPages"
          @click="goToPage(currentPage + 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          下一页 <i class="fas fa-chevron-right ml-1"></i>
        </button>
      </div>
    </div>

    </div>

    <!-- 二维码模态框 -->
    <QrCodeModal 
      :visible="showLinkModal" 
      :url="selectedResource?.url" 
      @close="showLinkModal = false" 
    />

    <!-- 页脚 -->
    <footer class="mt-auto py-6 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="max-w-7xl mx-auto text-center text-gray-600 dark:text-gray-400 text-sm px-3 sm:px-5">
        <p class="mb-2">本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小内删除!!!</p>
        <p>{{ systemConfig?.copyright || '© 2025 网盘资源数据库 By 老九' }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// 页面元数据
useHead({
  title: '网盘资源数据库 - 首页',
  meta: [
    { name: 'description', content: '网盘资源数据库 - 一个现代化的资源管理系统' },
    { name: 'keywords', content: '网盘资源,资源管理,数据库' }
  ]
})

// 获取运行时配置
const config = useRuntimeConfig()

// 获取路由参数
const route = useRoute()
const router = useRouter()

// 响应式数据
const searchQuery = ref(route.query.q as string || '')
const currentPage = ref(parseInt(route.query.page as string) || 1)
const pageSize = ref(200)
const selectedPlatform = ref(route.query.platform as string || '')
const showLinkModal = ref(false)
const selectedResource = ref<any>(null)
const authInitialized = ref(false)
const isLoadingMore = ref(false)
const hasMoreData = ref(true)
const pageLoading = ref(false)

console.log(pageSize.value, currentPage.value)

// 使用 useAsyncData 获取资源数据
const { data: resourcesData, pending, refresh } = await useAsyncData(
  () => `resources-${currentPage.value}-${searchQuery.value}-${selectedPlatform.value}`,
  () => $fetch('/api/resources', {
    params: {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
      pan_id: selectedPlatform.value
    }
  })
)

// 获取统计数据
const { data: statsData } = await useAsyncData('stats', 
  () => $fetch('/api/stats')
)

// 获取平台数据
const { data: platformsData } = await useAsyncData('platforms', 
  () => $fetch('/api/pans')
)

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('systemConfig', 
  () => $fetch('/api/system-config')
)

const sysConfig = (systemConfigData.value as any)?.data as any
const panList = (platformsData.value as any)?.data?.list as any[]
const resourceList = (resourcesData.value as any)?.data?.resources as any[]
const total = (resourcesData.value as any)?.data?.total as number

// 从 SSR 数据中获取值
const safeResources = computed(() => (resourcesData.value as any)?.data?.resources || [])
const safeStats = computed(() => (statsData.value as any)?.data || { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0, today_updates: 0 })
const platforms = computed(() => panList || [])
const systemConfig = computed(() => sysConfig || { site_title: '网盘资源数据库' })
const safeLoading = computed(() => pending.value)

// 计算属性
const totalPages = computed(() => {
  const total = (resourcesData.value as any)?.data?.total || 0
  return Math.ceil(total / pageSize.value)
})

// 用户状态管理
const userStore = useUserStore()

// 初始化认证状态
onMounted(() => {
  userStore.initAuth()
  authInitialized.value = true
  animateCounters()
})

// 搜索处理
const handleSearch = async () => {
  currentPage.value = 1
  
  // 更新URL参数
  const query = { ...route.query }
  if (searchQuery.value.trim()) {
    query.q = searchQuery.value.trim()
  } else {
    delete query.q
  }
  if (selectedPlatform.value) {
    query.platform = selectedPlatform.value
  } else {
    delete query.platform
  }
  delete query.page // 重置页码
  
  // 更新URL（不刷新页面）
  await router.push({ query })
  
  // 刷新数据
  await refresh()
}

// 平台筛选
const filterByPlatform = async (platformId: string) => {
  selectedPlatform.value = platformId
  currentPage.value = 1
  
  // 更新URL参数
  const query = { ...route.query }
  if (platformId) {
    query.platform = platformId
  } else {
    delete query.platform
  }
  delete query.page // 重置页码
  
  // 更新URL（不刷新页面）
  await router.push({ query })
  
  // 刷新数据
  await refresh()
}

// 获取平台名称
const getPlatformIcon = (panId: string) => {
  const platform = platforms.value.find(p => p.id === panId)
  return platform?.icon || '未知平台'
}

// 切换链接显示
const toggleLink = (resource: any) => {
  selectedResource.value = resource
  showLinkModal.value = true
}

// 复制到剪贴板
const copyToClipboard = async (text: any) => {
  try {
    await navigator.clipboard.writeText(text)
    const button = document.querySelector('.show-link-btn')
    if (button) {
      const originalText = button.innerHTML
      button.innerHTML = '<i class="fas fa-check"></i> 已复制'
      button.classList.add('bg-green-600')
      setTimeout(() => {
        button.innerHTML = originalText
        button.classList.remove('bg-green-600')
      }, 2000)
    }
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 跳转到链接
const openLink = (url: string) => {
  window.open(url, '_blank')
}

// 格式化相对时间
const formatRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  const diffWeek = Math.floor(diffDay / 7)
  const diffMonth = Math.floor(diffDay / 30)
  const diffYear = Math.floor(diffDay / 365)
  
  const isToday = date.toDateString() === now.toDateString()
  
  if (isToday) {
    if (diffMin < 1) {
      return '<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>刚刚更新</span>'
    } else if (diffHour < 1) {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffMin}分钟前</span>`
    } else {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffHour}小时前</span>`
    }
  } else if (diffDay < 1) {
    return `<span class="text-gray-600">${diffHour}小时前</span>`
  } else if (diffDay < 7) {
    return `<span class="text-gray-600">${diffDay}天前</span>`
  } else if (diffWeek < 4) {
    return `<span class="text-gray-600">${diffWeek}周前</span>`
  } else if (diffMonth < 12) {
    return `<span class="text-gray-600">${diffMonth}个月前</span>`
  } else {
    return `<span class="text-gray-600">${diffYear}年前</span>`
  }
}

// 检查是否为今天更新
const isUpdatedToday = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  return date.toDateString() === now.toDateString()
}

// 数字动画效果
const animateCounters = () => {
  const counters = document.querySelectorAll('.count-up')
  const speed = 200
  
  counters.forEach((counter) => {
    const target = parseInt(counter.getAttribute('data-target') || '0')
    const increment = Math.ceil(target / speed)
    let count = 0
    
    const updateCount = () => {
      if (count < target) {
        count += increment
        if (count > target) count = target
        counter.textContent = count.toString()
        setTimeout(updateCount, 1)
      } else {
        counter.textContent = target.toString()
      }
    }
    
    updateCount()
  })
}

// 页面跳转
const goToPage = async (page: number) => {
  currentPage.value = page
  
  // 更新URL参数
  const query = { ...route.query }
  query.page = page.toString()
  await router.push({ query })
  
  // 刷新数据
  await refresh()
  
  // 滚动到顶部
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}



const loadMore = async () => {
  if (isLoadingMore.value || !hasMoreData.value) return
  
  isLoadingMore.value = true
  try {
    currentPage.value++
    
    // 使用 refresh 获取更多数据
    await refresh()
    
    // 检查是否还有更多数据
    const currentTotal = (resourcesData.value as any)?.data?.total || 0
    const currentLoaded = safeResources.value.length
    
    if (currentLoaded >= currentTotal) {
      hasMoreData.value = false
    }
  } catch (error) {
    console.error('加载更多失败:', error)
    currentPage.value-- // 回退页码
  } finally {
    isLoadingMore.value = false
  }
}
</script>

<style scoped>
.active-filter {
  @apply bg-slate-800 text-white;
}

.count-up {
  transition: all 0.3s ease;
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}
</style> 