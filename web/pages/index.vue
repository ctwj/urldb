  <template>
  <div v-if="!systemConfig.maintenance_mode" class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col">
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
      <div class="header-container bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4">
          <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
            {{ systemConfig?.site_title || '老九网盘资源数据库' }}
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
        <ClientOnly>
          <div class="relative">
            <n-input round placeholder="搜索" v-model="searchQuery" @blur="handleSearch" @keyup.enter="handleSearch">
              <template #suffix>
                <i class="fas fa-search text-gray-400"></i>
              </template>
            </n-input>
          </div>
        </ClientOnly>
        
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
    <AppFooter />
  </div>
  <div v-if="systemConfig.maintenance_mode" class="fixed inset-0 z-[1000000] flex items-center justify-center bg-gradient-to-br from-yellow-100/80 via-gray-900/90 to-yellow-200/80 backdrop-blur-sm">
    <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-2xl px-8 py-10 flex flex-col items-center max-w-xs w-full border border-yellow-200 dark:border-yellow-700">
      <i class="fas fa-tools text-yellow-500 text-5xl mb-6 animate-bounce-slow"></i>
      <h3 class="text-2xl font-extrabold text-yellow-600 dark:text-yellow-400 mb-2 tracking-wide drop-shadow">系统维护中</h3>
      <p class="text-base text-gray-600 dark:text-gray-300 mb-6 text-center leading-relaxed">
        我们正在进行系统升级和维护，预计很快恢复服务。<br>
        请稍后再试，感谢您的理解与支持！
      </p>
      <!-- 动态点点动画 -->
      <div class="flex space-x-1 mt-2">
        <span class="w-2 h-2 bg-yellow-400 rounded-full animate-blink"></span>
        <span class="w-2 h-2 bg-yellow-500 rounded-full animate-blink delay-200"></span>
        <span class="w-2 h-2 bg-yellow-600 rounded-full animate-blink delay-400"></span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 页面元数据
useHead({
  title: '老九网盘资源数据库 - 首页',
  meta: [
    { name: 'description', content: '老九网盘资源管理系统， 一个现代化的网盘资源数据库，支持多网盘自动化转存分享，支持百度网盘，阿里云盘，夸克网盘， 天翼云盘，迅雷云盘，123云盘，115网盘，UC网盘' },
    { name: 'keywords', content: '网盘资源,资源管理,数据库' }
  ]
})

// 获取运行时配置
const config = useRuntimeConfig()

import { useResourceApi, useStatsApi, usePanApi, useSystemConfigApi, usePublicSystemConfigApi } from '~/composables/useApi'

const resourceApi = useResourceApi()
const statsApi = useStatsApi()
const panApi = usePanApi()
const publicSystemConfigApi = usePublicSystemConfigApi()

// 获取路由参数
const route = useRoute()
const router = useRouter()

// 响应式数据
const initSerch = route.query.search || ''
const oldQuery = ref(initSerch)
const searchQuery = ref(oldQuery)
const currentPage = ref(parseInt(route.query.page as string) || 1)
const pageSize = ref(200)
const selectedPlatform = ref(route.query.platform as string || '')
const showLinkModal = ref(false)
const selectedResource = ref<any>(null)
const authInitialized = ref(true) // 在app.vue中已经初始化，这里直接设为true
const isLoadingMore = ref(false)
const hasMoreData = ref(true)
const pageLoading = ref(false)

// 用户状态管理
const userStore = useUserStore()

// 使用 useAsyncData 获取资源数据
const { data: resourcesData, pending, refresh } = await useAsyncData(
  () => `resources-${currentPage.value}-${searchQuery.value}-${selectedPlatform.value}`,
  () => resourceApi.getResources({
    page: currentPage.value,
    page_size: pageSize.value,
    search: searchQuery.value,
    pan_id: selectedPlatform.value
  })
)

// 获取统计数据
const { data: statsData } = await useAsyncData('stats', () => statsApi.getStats())

// 获取平台数据
const { data: platformsData } = await useAsyncData('platforms', () => panApi.getPans())

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('systemConfig', () => publicSystemConfigApi.getPublicSystemConfig())

// 从 SSR 数据中获取值
const safeResources = computed(() => (resourcesData.value as any)?.data || [])
const safeStats = computed(() => (statsData.value as any) || { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0, today_updates: 0 })
const platforms = computed(() => (platformsData.value as any) || [])
const systemConfig = computed(() => (systemConfigData.value as any).data || { site_title: '老九网盘资源数据库' })
const safeLoading = computed(() => pending.value)

// 计算属性
const totalPages = computed(() => {
  const total = (resourcesData.value as any)?.total || 0
  return Math.ceil(total / pageSize.value)
})

// 初始化认证状态
onMounted(() => {
  animateCounters()
})

// 搜索处理
const handleSearch = async (e?: any) => {
  if (e && e.target && typeof e.target.value === 'string') {
    searchQuery.value = e.target.value
  }
  if (oldQuery.value === searchQuery.value) {
    return
  }
  oldQuery.value = searchQuery.value
  currentPage.value = 1
  
  // 更新URL参数
  const query = { ...route.query }
  if ((searchQuery.value as string).trim()) {
    query.search = (searchQuery.value as string).trim()
  } else {
    delete query.search
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
  const platform = (platforms.value as any).find((p: any) => p.id === panId)
  return platform?.icon || '未知平台'
}

// 跳转到链接
const openLink = async (url: string, resourceId: number) => {
  try {
    await fetch(`/api/resources/${resourceId}/view`, { method: 'POST' })
  } catch (e) {}
  if (process.client) {
    window.open(url, '_blank')
  }
}

// 切换链接显示
const toggleLink = async (resource: any) => {
  try {
    await resourceApi.incrementViewCount(resource.id)
  } catch (e) {}
  selectedResource.value = resource
  showLinkModal.value = true
}

// 复制到剪贴板
const copyToClipboard = async (text: any) => {
  try {
    await navigator.clipboard.writeText(text)
    if (process.client) {
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
    }
  } catch (error) {
    console.error('复制失败:', error)
  }
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
  if (!process.client) return
  
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
  
  // 滚动到顶部（只在客户端执行）
  if (process.client) {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }
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
@keyframes bounce-slow {
  0%, 100% { transform: translateY(0);}
  50% { transform: translateY(-12px);}
}
.animate-bounce-slow {
  animation: bounce-slow 1.6s infinite;
}
@keyframes blink {
  0%, 80%, 100% { opacity: 0.2;}
  40% { opacity: 1;}
}
.animate-blink {
  animation: blink 1.4s infinite both;
}
.animate-blink.delay-200 { animation-delay: 0.2s; }
.animate-blink.delay-400 { animation-delay: 0.4s; }
.header-container{
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}
</style> 