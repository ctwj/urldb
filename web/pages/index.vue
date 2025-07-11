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
            {{ systemConfig?.site_title || '网盘资源管理系统' }}
          </a>
        </h1>
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-4 right-4 top-0 absolute">
          <NuxtLink 
            to="/hot-dramas" 
            class="w-full sm:w-auto px-4 py-2 bg-purple-600 hover:bg-purple-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-film"></i> 热播剧
          </NuxtLink>
          <NuxtLink 
            to="/monitor" 
            class="w-full sm:w-auto px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-chart-line"></i> 系统监控
          </NuxtLink>
          <NuxtLink 
            v-if="authInitialized && !userStore.isAuthenticated"
            to="/login" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-sign-in-alt"></i> 登录
          </NuxtLink>
          <NuxtLink 
            v-if="authInitialized && userStore.isAuthenticated"
            to="/admin" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-user-shield"></i> 管理后台
          </NuxtLink>
        </nav>
      </div>

      <!-- 搜索区域 -->
      <div class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
        <div class="relative">
          <input 
            v-model="searchQuery"
            @keyup="debounceSearch"
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
            <span v-html="getPlatformIcon(platform.name)"></span> {{ platform.name }}
          </button>
        </div>
        
        <!-- 统计信息 -->
        <div class="flex justify-between mt-3 text-sm text-gray-600 dark:text-gray-300 px-2">
          <div class="flex items-center">
            <i class="fas fa-calendar-day text-pink-600 mr-1"></i>
            今日更新: <span class="font-medium text-pink-600 ml-1 count-up" :data-target="safeTodayUpdates">0</span>
          </div>
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 mr-1"></i>
            总资源数: <span class="font-medium text-blue-600 ml-1 count-up" :data-target="safeStats?.total_resources || 0">0</span>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg shadow">
        <table class="w-full">
          <thead>
            <tr class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100">
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">
                <div class="flex items-center">
                  <i class="fas fa-cloud mr-1 text-gray-300"></i> 文件名
                </div>
              </th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell">链接</th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">更新时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="safeLoading" class="text-center py-8">
              <td colspan="3" class="text-gray-500 dark:text-gray-400">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
            </tr>
            <tr v-else-if="safeResources.length === 0" class="text-center py-8">
              <td colspan="3" class="text-gray-500 dark:text-gray-400">暂无数据</td>
            </tr>
            <tr 
              v-for="(resource, index) in visibleResources" 
              :key="resource.id"
              :class="isUpdatedToday(resource.updated_at) ? 'hover:bg-pink-50 dark:hover:bg-pink-900 bg-pink-50/30 dark:bg-pink-900/30' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
              v-intersection="onIntersection"
              :data-index="index"
            >
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm">
                <div class="flex items-start">
                  <span class="mr-2 flex-shrink-0" v-html="getPlatformIcon(getPlatformName(resource.pan_id || 0))"></span>
                  <span class="break-words">{{ resource.title }}</span>
                </div>
                <div class="sm:hidden mt-1">
                  <button 
                    class="text-blue-600 hover:text-blue-800 text-xs flex items-center gap-1 show-link-btn" 
                    @click="toggleLink(resource)"
                  >
                    <i class="fas fa-eye"></i> 显示链接
                  </button>
                  <a 
                    v-if="resource.showLink" 
                    :href="resource.url" 
                    target="_blank" 
                    class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline text-xs break-all"
                  >
                    {{ resource.url }}
                  </a>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm hidden sm:table-cell">
                <button 
                  class="text-blue-600 hover:text-blue-800 flex items-center gap-1 show-link-btn" 
                  @click="toggleLink(resource)"
                >
                  <i class="fas fa-eye"></i> 显示链接
                </button>
                <a 
                  v-if="resource.showLink" 
                  :href="resource.url" 
                  target="_blank" 
                  class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline"
                >
                  {{ resource.url }}
                </a>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm text-gray-500" :title="resource.updated_at">
                {{ formatRelativeTime(resource.updated_at) }}
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

    <!-- 添加资源模态框 -->
    <!-- <ResourceModal
      v-if="showAddResourceModal"
      :resource="editingResource"
      @close="closeModal"
      @save="handleSaveResource"
    /> -->
    </div>

    <!-- 页脚 -->
    <footer class="mt-auto py-6 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="max-w-7xl mx-auto text-center text-gray-600 dark:text-gray-400 text-sm px-3 sm:px-5">
        <p class="mb-2">本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小内删除!!!</p>
        <p>{{ systemConfig?.copyright || '© 2025 网盘资源管理系统 By 小七' }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// 响应式数据
const searchQuery = ref('')
const selectedPlatform = ref('')
const authInitialized = ref(false) // 添加认证状态初始化标志
const pageLoading = ref(true) // 添加页面加载状态
const systemConfig = ref<SystemConfig | null>(null) // 添加系统配置状态

// 虚拟滚动相关
const visibleResources = ref<any[]>([])
const hasMoreData = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const isLoadingMore = ref(false)

// 延迟初始化store，避免SSR过程中的错误
let store: any = null
let userStore: any = null

// 本地状态管理，避免SSR过程中的store访问
const localResources = ref<any[]>([])
const localCategories = ref<any[]>([])
const localStats = ref<any>({ total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0 })
const localLoading = ref(false)

// 安全的store访问
const safeResources = computed(() => {
  try {
    if (process.client && store) {
      const storeRefs = storeToRefs(store)
      return (storeRefs as any).resources?.value || localResources.value
    }
    return localResources.value
  } catch (error) {
    console.error('获取resources时出错:', error)
    return localResources.value
  }
})

const safeCategories = computed(() => {
  try {
    if (process.client && store) {
      const storeRefs = storeToRefs(store)
      return (storeRefs as any).categories?.value || localCategories.value
    }
    return localCategories.value
  } catch (error) {
    console.error('获取categories时出错:', error)
    return localCategories.value
  }
})

const safeStats = computed(() => {
  try {
    if (process.client && store) {
      const storeRefs = storeToRefs(store)
      return (storeRefs as any).stats?.value || localStats.value
    }
    return localStats.value
  } catch (error) {
    console.error('获取stats时出错:', error)
    return localStats.value
  }
})

const safeLoading = computed(() => {
  try {
    if (process.client && store) {
      const storeRefs = storeToRefs(store)
      return (storeRefs as any).loading?.value || localLoading.value
    }
    return localLoading.value
  } catch (error) {
    console.error('获取loading时出错:', error)
    return localLoading.value
  }
})

// 动态SEO配置
const seoConfig = computed(() => ({
  title: systemConfig.value?.site_title || '网盘资源管理系统',
  meta: [
    { 
      name: 'description', 
      content: systemConfig.value?.site_description || '专业的网盘资源管理系统' 
    },
    { 
      name: 'keywords', 
      content: systemConfig.value?.keywords || '网盘,资源管理,文件分享' 
    },
    { 
      name: 'author', 
      content: systemConfig.value?.author || '系统管理员' 
    },
    { 
      name: 'copyright', 
      content: systemConfig.value?.copyright || '© 2024 网盘资源管理系统' 
    }
  ]
}))

// 页面元数据 - 使用watchEffect来避免组件卸载时的错误
watchEffect(() => {
  try {
    if (systemConfig.value && systemConfig.value.site_title) {
      useHead({
        title: systemConfig.value.site_title || '网盘资源管理系统',
        meta: [
          { 
            name: 'description', 
            content: systemConfig.value.site_description || '专业的网盘资源管理系统' 
          },
          { 
            name: 'keywords', 
            content: systemConfig.value.keywords || '网盘,资源管理,文件分享' 
          },
          { 
            name: 'author', 
            content: systemConfig.value.author || '系统管理员' 
          },
          { 
            name: 'copyright', 
            content: systemConfig.value.copyright || '© 2024 网盘资源管理系统' 
          }
        ]
      })
    } else {
      // 默认SEO配置
      useHead({
        title: '网盘资源管理系统',
        meta: [
          { name: 'description', content: '专业的网盘资源管理系统' },
          { name: 'keywords', content: '网盘,资源管理,文件分享' },
          { name: 'author', content: '系统管理员' },
          { name: 'copyright', content: '© 2024 网盘资源管理系统' }
        ]
      })
    }
  } catch (error) {
    console.error('设置页面元数据时出错:', error)
    // 使用默认配置作为后备
    useHead({
      title: '网盘资源管理系统',
      meta: [
        { name: 'description', content: '专业的网盘资源管理系统' },
        { name: 'keywords', content: '网盘,资源管理,文件分享' },
        { name: 'author', content: '系统管理员' },
        { name: 'copyright', content: '© 2024 网盘资源管理系统' }
      ]
    })
  }
})

// API
// const { getSystemConfig } = useSystemConfigApi()

// const showAddResourceModal = ref(false)
const editingResource = ref<any>(null)
const totalPages = ref(1)
interface Platform {
  id: number
  name: string
  key?: number
  ck?: string
  is_valid: boolean
  space: number
  left_space: number
  remark: string
  created_at: string
  updated_at: string
  icon?: string // 新增图标字段
}

interface ExtendedResource {
  id: number
  title: string
  description: string
  url: string
  pan_id?: number
  quark_url?: string
  file_size?: string
  category_id?: number
  category_name?: string
  view_count: number
  is_valid?: boolean
  is_public?: boolean
  created_at: string
  updated_at: string
  showLink?: boolean
}

interface SystemConfig {
  id: number
  site_title: string
  site_description: string
  keywords: string
  author: string
  copyright: string
  auto_process_ready_resources: boolean
  auto_process_interval: number
  page_size: number
  maintenance_mode: boolean
  created_at: string
  updated_at: string
}

const platforms = ref<Platform[]>([])
const todayUpdates = ref(0)

// 安全地计算今日更新数量
const safeTodayUpdates = computed(() => {
  try {
    const resources = safeResources.value
    if (!resources || !Array.isArray(resources)) {
      return 0
    }
    const today = new Date().toDateString()
    return resources.filter((resource: any) => {
      if (!resource || !resource.updated_at) return false
      try {
        return new Date(resource.updated_at).toDateString() === today
      } catch (dateError) {
        console.error('解析日期时出错:', dateError)
        return false
      }
    }).length
  } catch (error) {
    console.error('计算今日更新数量时出错:', error)
    return 0
  }
})

// 防抖搜索
let searchTimeout: NodeJS.Timeout
const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    handleSearch()
  }, 500)
}

// 获取系统配置
const fetchSystemConfig = async () => {
  try {
    // const response = await getSystemConfig() as any
    // console.log('首页系统配置响应:', response)
    // if (response && response.success && response.data) {
    //   systemConfig.value = response.data
    // } else if (response && response.data) {
    //   // 兼容非标准格式
    //   systemConfig.value = response.data
    // }
    console.log('系统配置功能暂时禁用')
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 获取数据
onMounted(async () => {
  console.log('首页 - onMounted 开始')
  
  try {
    // 初始化store
    store = useResourceStore()
    userStore = useUserStore()
    
    // 初始化用户状态
    userStore.initAuth()
    authInitialized.value = true // 设置认证状态初始化完成
    
    console.log('首页 - authInitialized:', authInitialized.value)
    console.log('首页 - isAuthenticated:', userStore.isAuthenticated)
    console.log('首页 - user:', userStore.userInfo)
    
    // 设置超时时间（5秒）
    const timeout = 5000
    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => reject(new Error('请求超时')), timeout)
    })
    
    // 使用Promise.race来添加超时机制，并优化请求顺序
    try {
      // 首先加载最重要的数据（资源列表）
      const resourcesPromise = store.fetchResources().then((data: any) => {
        localResources.value = data.resources || []
        return data
      }).catch((e: any) => {
        console.error('获取资源失败:', e)
        return { resources: [] }
      })
      
      // 等待资源数据加载完成或超时
      await Promise.race([resourcesPromise, timeoutPromise])
      
      // 然后并行加载其他数据
      const otherDataPromise = Promise.allSettled([
        store.fetchCategories().then((data: any) => {
          localCategories.value = data.categories || []
          return data
        }).catch((e: any) => {
          console.error('获取分类失败:', e)
          return { categories: [] }
        }),
        store.fetchStats().then((data: any) => {
          localStats.value = data || { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0 }
          return data
        }).catch((e: any) => {
          console.error('获取统计失败:', e)
          return { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0 }
        }),
        fetchPlatforms().catch((e: any) => console.error('获取平台失败:', e)),
        fetchSystemConfig().catch((e: any) => console.error('获取系统配置失败:', e)),
      ])
      
      // 检查哪些请求成功了
      otherDataPromise.then(results => {
        results.forEach((result, index) => {
          if (result.status === 'rejected') {
            console.error(`请求 ${index} 失败:`, result.reason)
          }
        })
      })
      
    } catch (timeoutError) {
      console.warn('部分数据加载超时，使用默认数据:', timeoutError)
      // 超时后使用默认数据，不阻塞页面显示
    }
    
    animateCounters()
  } catch (error) {
    console.error('页面数据加载失败:', error)
  } finally {
    // 无论成功还是失败，都要关闭加载状态
    pageLoading.value = false
    console.log('首页 - onMounted 完成')
  }
})

// 获取平台列表
const fetchPlatforms = async () => {
  try {
    const { usePanApi } = await import('~/composables/useApi')
    const panApi = usePanApi()
    const response = await panApi.getPans() as any
    // 后端直接返回数组，不需要 .pans
    platforms.value = Array.isArray(response) ? response : []
    console.log('获取到的平台数据:', platforms.value)
  } catch (error) {
    console.error('获取平台列表失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  try {
    if (!store || !process.client) {
      console.error('store未初始化或不在客户端')
      return
    }
    const platformId = selectedPlatform.value ? parseInt(selectedPlatform.value) : undefined
    store.searchResources(searchQuery.value, platformId).then((data: any) => {
      localResources.value = data.resources || []
    }).catch((error: any) => {
      console.error('搜索失败:', error)
    })
  } catch (error) {
    console.error('搜索处理时出错:', error)
  }
}

// 按平台筛选
const filterByPlatform = (platformId: string | number) => {
  selectedPlatform.value = platformId.toString()
  currentPage.value = 1
  handleSearch()
}

// 获取平台图标
const getPlatformIcon = (platformName: string) => {
  // 首先尝试从平台列表中查找对应的平台
  const platform = platforms.value.find((p: Platform) => p.name === platformName)
  if (platform && platform.icon) {
    return platform.icon
  }
  
  // 如果找不到对应的平台或没有图标，使用默认图标
  const defaultIcons: Record<string, string> = {
    'unknown': '<i class="fas fa-question-circle text-gray-400"></i>',
    'other': '<i class="fas fa-cloud text-gray-500"></i>',
    'magnet': '<i class="fas fa-magnet text-red-600"></i>',
    'uc': '<i class="fas fa-cloud-download-alt text-purple-600"></i>',
    '夸克网盘': '<i class="fas fa-cloud text-blue-600"></i>',
    '阿里云盘': '<i class="fas fa-cloud text-orange-600"></i>',
    '百度网盘': '<i class="fas fa-cloud text-blue-500"></i>',
    '天翼云盘': '<i class="fas fa-cloud text-red-500"></i>',
    'OneDrive': '<i class="fas fa-cloud text-blue-700"></i>',
    'Google Drive': '<i class="fas fa-cloud text-green-600"></i>'
  }
  
  return defaultIcons[platformName] || defaultIcons['unknown']
}

// 获取平台名称
const getPlatformName = (platformId: number) => {
  const platform = platforms.value.find((p: Platform) => p.id === platformId)
  return platform?.name || 'unknown'
}

// 切换链接显示
const toggleLink = (resource: any) => {
  resource.showLink = !resource.showLink
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
    return `${diffHour}小时前`
  } else if (diffDay < 7) {
    return `${diffDay}天前`
  } else if (diffWeek < 4) {
    return `${diffWeek}周前`
  } else if (diffMonth < 12) {
    return `${diffMonth}个月前`
  } else {
    return `${diffYear}年前`
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
  
  counters.forEach((counter: Element) => {
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
const goToPage = (page: number) => {
  currentPage.value = page
  handleSearch()
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// 编辑资源
// const editResource = (resource: any) => {
//   editingResource.value = resource
//   showAddResourceModal.value = true
// }

// 删除资源
const deleteResource = async (id: number) => {
  if (confirm('确定要删除这个资源吗？')) {
    try {
      await store.deleteResource(id)
    } catch (error) {
      console.error('删除失败:', error)
    }
  }
}

// 关闭模态框
// const closeModal = () => {
//   showAddResourceModal.value = false
//   editingResource.value = null
// }

// 保存资源
const handleSaveResource = async (resourceData: any) => {
  try {
    if (editingResource.value) {
      await store.updateResource(editingResource.value.id, resourceData)
    } else {
      await store.createResource(resourceData)
    }
    // closeModal() // 移除未定义的函数调用
  } catch (error) {
    console.error('保存失败:', error)
  }
}

// 虚拟滚动相关方法
const onIntersection = (entries: IntersectionObserverEntry[]) => {
  entries.forEach(entry => {
    if (entry.isIntersecting && hasMoreData.value && !isLoadingMore.value) {
      loadMore()
    }
  })
}

const loadMore = async () => {
  if (isLoadingMore.value || !hasMoreData.value) return
  
  isLoadingMore.value = true
  try {
    currentPage.value++
    const newResources = await fetchResources(currentPage.value, pageSize.value)
    if (newResources && newResources.length > 0) {
      visibleResources.value.push(...newResources)
    } else {
      hasMoreData.value = false
    }
  } catch (error) {
    console.error('加载更多失败:', error)
    currentPage.value-- // 回退页码
  } finally {
    isLoadingMore.value = false
  }
}

const fetchResources = async (page: number, size: number) => {
  try {
    const { useResourceApi } = await import('~/composables/useApi')
    const resourceApi = useResourceApi()
    const response = await resourceApi.getResources({
      page,
      page_size: size
    }) as any
    return response.resources || []
  } catch (error) {
    console.error('获取资源失败:', error)
    return []
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