  <template>
  <div v-if="!systemConfig.maintenance_mode" class="min-h-screen bg-gray-50 dark:bg-slate-900 text-gray-800 dark:text-slate-100 flex flex-col">
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
      <div class="header-container bg-slate-800 dark:bg-slate-800 text-white dark:text-slate-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4 flex items-center justify-center gap-3">
          <img 
            v-if="systemConfig?.site_logo" 
            :src="getImageUrl(systemConfig.site_logo)" 
            :alt="systemConfig?.site_title || 'Logo'"
            class="h-8 w-auto object-contain"
            @error="handleLogoError"
          />
          <img 
            v-else
            src="/assets/images/logo.webp" 
            alt="Logo" 
            class="h-8 w-auto object-contain"
          />
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
          <ClientOnly>
            <NuxtLink v-if="authInitialized && !userStore.isAuthenticated" to="/login" class="sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-sign-in-alt text-xs"></i> 登录
              </n-button>
            </NuxtLink>
            <NuxtLink v-if="authInitialized && userStore.isAuthenticated && userStore.user?.role === 'admin'" to="/admin" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-user-shield text-xs"></i> 管理后台
              </n-button>
            </NuxtLink>
            <NuxtLink v-if="authInitialized && userStore.isAuthenticated && userStore.user?.role !== 'admin'" to="/user" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-user text-xs"></i> 用户中心
              </n-button>
            </NuxtLink>
          </ClientOnly>
        </nav>
      </div>

      <!-- 公告信息 -->
      <div class="w-full max-w-3xl mx-auto mb-2 px-2 sm:px-0">
        <Announcement />
      </div>

      <!-- 搜索区域 -->
      <div class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
        <ClientOnly>
          <div class="relative">
            <n-input round placeholder="搜索" v-model:value="searchQuery" @blur="handleSearch" @keyup.enter="handleSearch" clearable>
                <template #prefix>
                <i class="fas fa-search text-gray-400"></i>
                </template>
              </n-input>
          </div>
        </ClientOnly>

        <!-- 平台类型筛选 -->
        <div class="mt-3 flex flex-wrap gap-2" id="platformFilters">
          <a 
            :href="`/?search=${$route.query.search || ''}&platform=`"
            class="px-2 py-1 text-xs rounded-full bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 hover:bg-slate-700 dark:hover:bg-gray-600 transition-colors"
            :class="{ 'active-filter': !selectedPlatform }"
          >
            全部
          </a>
          <a 
            v-for="platform in platforms" 
            :key="platform.id"
            :href="`/?search=${$route.query.search || ''}&platform=${platform.id}`"
            class="px-2 py-1 text-xs rounded-full bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
            :class="{ 'active-filter': selectedPlatform === platform.id }"
          >
            <span v-html="platform.icon"></span> {{ platform.name }}
          </a>
        </div>
        
        <!-- 统计信息 -->
        <div class="flex justify-between mt-3 text-sm text-gray-600 dark:text-gray-300 px-2">
          <div class="flex items-center">
            <i class="fas fa-calendar-day text-pink-600 mr-1"></i>
            今日资源: <span class="font-medium text-pink-600 ml-1 count-up" :data-target="safeStats?.today_resources || 0">0</span>
          </div>
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 mr-1"></i>
            总资源数: <span class="font-medium text-blue-600 ml-1 count-up" :data-target="safeStats?.total_resources || 0">0</span>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="overflow-x-auto bg-white dark:bg-slate-800 rounded-lg shadow-lg shadow-slate-900/10 dark:shadow-slate-900/50">
        <table class="w-full min-w-full">
          <thead>
            <tr class="bg-slate-800 dark:bg-slate-700 text-white dark:text-slate-100">
              <th class="text-left text-xs sm:text-sm w-20 pl-2 sm:pl-3">
                <div class="flex items-center">
                  <i class="fas fa-image mr-1 text-gray-300 dark:text-slate-300"></i> 封面
                </div>
              </th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">
                <div class="flex items-center">
                  <i class="fas fa-cloud mr-1 text-gray-300 dark:text-slate-300"></i> 文件名
                </div>
              </th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-24">链接</th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-32">更新时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-slate-700">
            <tr v-if="safeLoading" class="text-center py-8">
              <td colspan="1" class="text-gray-500 dark:text-gray-400 sm:hidden">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
              <td colspan="4" class="text-gray-500 dark:text-gray-400 hidden sm:table-cell">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
            </tr>
            <tr v-else-if="safeResources.length === 0" class="text-center py-12">
              <td colspan="4" class="text-gray-500 dark:text-slate-500">
                <div class="flex flex-col items-center justify-center space-y-4">
                  <img 
                    src="/assets/svg/empty.svg" 
                    alt="暂无数据" 
                    class="!w-64 !h-64 sm:w-64 sm:h-64 opacity-60 dark:opacity-40"
                  />
                  <div class="text-center">
                    <p class="text-lg font-medium text-gray-600 dark:text-gray-400 mb-2">
                      {{ searchQuery ? '没有找到相关资源' : '暂无资源数据' }}
                    </p>
                    <p class="text-sm text-gray-500 dark:text-gray-500">
                      {{ searchQuery ? '请尝试其他关键词或清除搜索条件' : '资源正在整理中，请稍后再来查看' }}
                    </p>
                  </div>
                </div>
              </td>
            </tr>
            <tr
              v-for="(resource, index) in safeResources"
              :key="resource.id"
              :class="isUpdatedToday(resource.updated_at) ? 'hover:bg-pink-50 dark:hover:bg-pink-500/10 bg-pink-50/30 dark:bg-pink-500/5' : 'hover:bg-gray-50 dark:hover:bg-slate-700/50'"
              :data-index="index"
            >
              <td class="text-xs sm:text-sm w-20 pl-2 sm:pl-3">
                <div class="flex justify-center">
                  <ClientOnly>
                    <n-image
                      :src="getResourceImageUrl(resource)"
                      :alt="resource.title || '资源图片'"
                      width="80"
                      class="rounded object-cover border border-gray-200 dark:border-slate-600 h-auto"
                      lazy
                      @error="handleResourceImageError"
                    />
                    <template #placeholder>
                      <div class="w-[80px] h-[80px] rounded bg-gray-200 dark:bg-slate-700 flex items-center justify-center">
                        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                      </div>
                    </template>
                    <template #fallback>
                      <div class="w-[80px] h-[80px] rounded bg-gray-200 dark:bg-slate-700 animate-pulse"></div>
                    </template>
                  </ClientOnly>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm">
                <div class="flex items-start">
                  <span class="mr-2 flex-shrink-0" v-html="getPlatformIcon(resource.pan_id || 0)"></span>
                  <div class="flex-1 min-w-0">
                    <div class="break-words font-medium" v-html="resource.title_highlight || resource.title"></div>
                    <!-- 显示标签 -->
                    <div v-if="resource.tags && resource.tags.length > 0" class="mt-1 flex flex-wrap gap-1">
                      <template v-for="(tag, index) in resource.tags" :key="tag.id">
                        <span
                          class="resource-tag inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-500/20 text-blue-800 dark:text-blue-100 border dark:border-blue-400/30"
                          :title="tag.name"
                        >
                          <i class="fas fa-tag mr-1 dark:text-blue-200"></i>
                          <span>{{ tag.name || '未知标签' }}</span>
                        </span>
                      </template>
                    </div>
                    <!-- 显示描述 -->
                    <div v-if="resource.description_highlight || resource.description" class="text-xs text-gray-600 dark:text-slate-400 mt-1 break-words line-clamp-2" v-html="resource.description_highlight || resource.description">
                    </div>
                  </div>
                </div>
                <div class="sm:hidden mt-2 space-y-2">
                  <!-- 移动端时间和链接按钮一行显示 -->
                  <div class="flex items-center gap-2">
                    <div class="flex-1 min-w-0">
                      <div class="text-xs text-gray-500 dark:text-slate-400 truncate" :title="resource.updated_at">
                        <span v-html="formatRelativeTime(resource.updated_at)"></span>
                      </div>
                    </div>
                    <div class="flex-1 flex justify-end">
                      <button
                        class="mobile-link-btn flex items-center gap-1 text-xs"
                        @click="toggleLink(resource)"
                      >
                        <i class="fas fa-eye"></i> 显示链接
                      </button>
                    </div>
                  </div>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm hidden sm:table-cell w-32">
                <button 
                  class="text-blue-600 hover:text-blue-800 flex items-center gap-1 show-link-btn" 
                  @click="toggleLink(resource)"
                >
                  <i class="fas fa-eye"></i> 显示链接
                </button>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm text-gray-500 hidden sm:table-cell w-32" :title="resource.updated_at">
                <span v-html="formatRelativeTime(resource.updated_at)"></span>
              </td>
            </tr>
          </tbody>
        </table>
        

      </div>

    </div>

    </div>

    <!-- 二维码模态框 -->
    <QrCodeModal 
      :visible="showLinkModal" 
      :url="selectedResource?.url" 
      :save_url="selectedResource?.save_url"
      :loading="selectedResource?.loading"
      :linkType="selectedResource?.linkType"
      :platform="selectedResource?.platform"
      :message="selectedResource?.message"
      :error="selectedResource?.error"
      :forbidden="selectedResource?.forbidden"
      :forbidden_words="selectedResource?.forbidden_words"
      @close="showLinkModal = false" 
    />

    <!-- 页脚 -->
    <AppFooter />

    <!-- 悬浮按钮组件 -->
    <FloatButtons />
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

import { useResourceApi, useStatsApi, usePanApi, useSystemConfigApi, usePublicSystemConfigApi, useSearchStatsApi } from '~/composables/useApi'

const resourceApi = useResourceApi()
const statsApi = useStatsApi()
const panApi = usePanApi()
const publicSystemConfigApi = usePublicSystemConfigApi()

// 获取路由参数
const route = useRoute()
const router = useRouter()

// 响应式数据
const showLinkModal = ref(false)
const selectedResource = ref<any>(null)
const pageLoading = ref(false)

// 使用 ClientOnly 包装器来处理认证状态
const authInitialized = ref(false)

// 用户状态管理
const userStore = useUserStore()

// 图片URL处理
const { getImageUrl } = useImageUrl()

// Logo错误处理
const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/assets/images/logo.webp'
}

// 获取资源图片URL，如果没有则返回随机默认封面
const getResourceImageUrl = (resource: any) => {
  // console.log('Resource data:', resource)
  // 如果资源有图片，使用资源图片（优先检查image_url，其次检查cover）
  if (resource.image_url) {
    return getImageUrl(resource.image_url)
  }

  if (resource.cover) {
    return getImageUrl(resource.cover)
  }

  // 否则随机选择默认封面图片 (cover1.webp 到 cover8.webp)
  const randomNum = Math.floor(Math.random() * 8) + 1
  return `/assets/images/cover${randomNum}.webp`
}

// 处理资源图片加载错误
const handleResourceImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  // 重新设置一个随机的默认封面图片
  const randomNum = Math.floor(Math.random() * 8) + 1
  img.src = `/assets/images/cover${randomNum}.webp`
}

// 使用 useAsyncData 获取资源数据
const { data: resourcesData, pending, refresh } = await useAsyncData(
  () => `resources-1-${route.query.search || ''}-${route.query.platform || ''}`,
  async () => {
    // 如果有搜索关键词，使用带搜索参数的资源接口（后端会优先使用Meilisearch）
    if (route.query.search) {
      return await resourceApi.getResources({
        page: 1,
        page_size: 200,
        search: route.query.search as string,
        pan_id: route.query.platform as string || ''
      })
    } else {
      // 没有搜索关键词时，使用普通资源接口获取最新数据
      return await resourceApi.getResources({
        page: 1,
        page_size: 200,
        pan_id: route.query.platform as string || ''
      })
    }
  }
)

// 获取统计数据
const { data: statsData, error: statsError } = await useAsyncData('stats', () => statsApi.getStats())

// 获取平台数据
const { data: platformsData, error: platformsError } = await useAsyncData('platforms', () => panApi.getPans())

// 获取系统配置
const { data: systemConfigData, error: systemConfigError } = await useAsyncData('systemConfig', () => publicSystemConfigApi.getPublicSystemConfig())

// 错误处理
const notification = ref()

// 监听错误
watch(statsError, (error) => {
  if (error && process.client) {
    console.error('获取统计数据失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取统计数据失败',
      duration: 5000
    })
  }
})

watch(platformsError, (error) => {
  if (error && process.client) {
    console.error('获取平台数据失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取平台数据失败',
      duration: 5000
    })
  }
})

watch(systemConfigError, (error) => {
  if (error && process.client) {
    console.error('获取系统配置失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取系统配置失败',
      duration: 5000
    })
  }
})

// 从 SSR 数据中获取值
const safeResources = computed(() => {
  const data = resourcesData.value as any
  // console.log('原始API数据结构:', JSON.stringify(data, null, 2))

  // 处理嵌套的data结构：{data: {data: [...], total: ...}}
  if (data?.data?.data && Array.isArray(data.data.data)) {
    const resources = data.data.data
    console.log('第一层嵌套资源:', resources)
    return resources
  }
  // 处理直接的data结构：{data: [...], total: ...}
  if (data?.data && Array.isArray(data.data)) {
    const resources = data.data
    // console.log('第二层嵌套资源:', resources)
    return resources
  }
  // 处理直接的数组结构
  if (Array.isArray(data)) {
    // console.log('直接数组结构:', data)
    return data
  }

  // console.log('未匹配到任何数据结构')
  return []
})
const safeStats = computed(() => (statsData.value as any) || { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0, today_resources: 0 })
const platforms = computed(() => (platformsData.value as any) || [])
const systemConfig = computed(() => (systemConfigData.value as any)?.data || { site_title: '老九网盘资源数据库' })
const safeLoading = computed(() => pending.value)


// 从路由参数获取当前状态
const searchQuery = ref(route.query.search as string || '')
const selectedPlatform = computed(() => route.query.platform as string || '')

// 记录搜索统计的函数
const recordSearchStats = (keyword: string) => {
  if (!keyword || keyword.trim().length === 0) {
    // console.log('搜索关键词为空，跳过统计记录')
    return
  }
  
  const trimmedKeyword = keyword.trim()
  // console.log('记录搜索统计:', trimmedKeyword)
  
  // 延迟执行，确保页面完全加载
  setTimeout(() => {
    const searchStatsApi = useSearchStatsApi()
    searchStatsApi.recordSearch({ keyword: trimmedKeyword }).catch(err => {
      console.error('记录搜索统计失败:', err)
    })
  }, 0)
}

const handleSearch = () => {
  const params = new URLSearchParams()
  if (searchQuery.value) params.set('search', searchQuery.value)
  if (selectedPlatform.value) params.set('platform', selectedPlatform.value)
  window.location.href = `/?${params.toString()}`
}

// 初始化认证状态
onMounted(() => {
  // 初始化认证状态
  authInitialized.value = true
  
  animateCounters()
  
  // 页面挂载完成时，如果有搜索关键词，记录搜索统计
  if (process.client && route.query.search) {
    const searchKeyword = route.query.search as string
    recordSearchStats(searchKeyword)
  } else {
    console.log('无搜索参数，跳过统计记录')
  }
})



// 获取平台名称
const getPlatformIcon = (panId: string | number) => {
  const platform = (platforms.value as any).find((p: any) => p.id == panId)
  return platform?.icon || '未知平台'
}

// 注意：链接访问统计已整合到 getResourceLink API 中

// 切换链接显示
const toggleLink = async (resource: any) => {
  // 如果包含违禁词，直接显示禁止访问，不发送请求
  if (resource.has_forbidden_words) {
    selectedResource.value = {
      ...resource,
      forbidden: true,
      error: '该资源包含违禁内容，无法访问',
      forbidden_words: resource.forbidden_words || []
    }
    showLinkModal.value = true
    return
  }

  // 显示加载状态
  selectedResource.value = { ...resource, loading: true }
  showLinkModal.value = true
  
  try {
    // 调用新的获取链接API（同时统计访问次数）
    const linkData = await resourceApi.getResourceLink(resource.id) as any
    console.log('获取到的链接数据:', linkData)
    
    // 更新资源信息，包含新的链接信息
    selectedResource.value = {
      ...resource,
      url: linkData.url,
      save_url: linkData.type === 'transferred' ? linkData.url : resource.save_url,
      loading: false,
      linkType: linkData.type,
      platform: linkData.platform,
      message: linkData.message
    }
  } catch (error: any) {
    console.error('获取资源链接失败:', error)
    
    // 其他错误
    selectedResource.value = {
      ...resource,
      loading: false,
      error: '检测有效性失败，请自行验证'
    }
  }
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

  // 处理今天更新的情况
  if (isToday) {
    if (diffMin < 1) {
      return '<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>刚刚更新</span>'
    } else if (diffHour < 1) {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffMin}分钟前</span>`
    } else {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffHour}小时前</span>`
    }
  }

  // 处理昨天更新的情况 - 显示具体时间
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  const isYesterday = date.toDateString() === yesterday.toDateString()

  if (isYesterday) {
    if (diffHour < 24) {
      // 昨天但不足24小时
      if (diffHour < 1) {
        return `<span class="text-gray-600">${diffMin}分钟前</span>`
      } else {
        return `<span class="text-gray-600">${diffHour}小时前</span>`
      }
    } else {
      // 超过24小时但仍然是昨天
      return `<span class="text-gray-600">${diffDay}天前</span>`
    }
  }

  // 处理其他情况
  if (diffDay < 7) {
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

/* 文本截断样式 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  word-wrap: break-word;
  word-break: break-word;
}

/* 表格单元格内容溢出控制 */
table td {
  overflow: hidden;
  word-wrap: break-word;
  word-break: break-word;
}

/* 确保flex容器不会溢出 */
.min-w-0 {
  min-width: 0;
}

/* 标签样式优化 */
.resource-tag {
  transition: all 0.2s ease;
}

.resource-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

/* 移动端按钮专用样式 */
.mobile-link-btn {
  border: 1px solid transparent;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 500;
  font-size: 11px;
  line-height: 1.2;
  transition: all 0.3s ease;
  min-height: 28px;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

.mobile-link-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s;
}

.mobile-link-btn:hover::before {
  left: 100%;
}

.mobile-link-btn:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  transform: translateY(-1px);
}

.mobile-link-btn:active {
  background: linear-gradient(135deg, #1d4ed8 0%, #172554 100%);
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.4);
}

.mobile-link-btn:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.6);
}
</style> 