<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col">
    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
        <!-- 头部 -->
        <div class="header-container bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
          <h1 class="text-2xl sm:text-3xl font-bold mb-4">
            <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
              热播剧榜单
            </a>
          </h1>
          <p class="text-gray-300 max-w-2xl mx-auto">实时获取豆瓣热门电影和电视剧榜单</p>
          <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
            <NuxtLink to="/" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-home text-xs"></i> 首页
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
          </nav>
        </div>

        <!-- 筛选器 -->
        <div class="mb-6 flex flex-wrap gap-4">
          <button
            v-for="category in categories"
            :key="category.value"
            @click="selectedCategory = category.value"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              selectedCategory === category.value
                ? 'bg-blue-600 text-white'
                : 'bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border border-gray-300 dark:border-gray-600'
            ]"
          >
            {{ category.label }}
          </button>
        </div>


        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>

        <!-- 热播剧列表 -->
        <div v-else-if="filteredDramas.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          <a
            v-for="(drama, index) in filteredDramas"
            :key="drama.id"
            :data-item-id="drama.id"
            :data-item-index="index"
            :href="`/?search=${encodeURIComponent(drama.title)}`"
            class="group relative bg-white/10 dark:bg-gray-800/10 backdrop-blur-md rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 border border-white/20 dark:border-gray-700/50 hover:scale-105 cursor-pointer no-underline block"
          >
            <!-- 海报图片 -->
            <div v-if="drama.poster_url" class="relative overflow-hidden h-52">
                <!-- 主图片（SSR数据立即显示，分页数据延迟加载） -->
              <img
                v-if="shouldShowImage(index, drama.id)"
                :src="getPosterUrl(drama.poster_url)"
                :alt="drama.title"
                class="w-full h-full object-cover"
                @error="handleImageError"
              />
              <!-- 图片上的遮罩和信息（只在图片显示后显示） -->
              <div v-if="shouldShowImage(index, drama.id)" class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>

              <!-- 新剧标签 -->
              <div
                v-if="drama.is_new && shouldShowImage(index, drama.id)"
                class="absolute top-3 right-3 bg-gradient-to-r from-red-500 to-red-600 text-white px-3 py-1 rounded-full text-xs font-semibold shadow-lg z-10"
              >
                🔥 HOT
              </div>

              <!-- 评分显示 -->
              <div v-if="shouldShowImage(index, drama.id)" class="absolute bottom-3 left-3 right-3 flex items-center justify-between z-20">
                <div class="bg-black/60 backdrop-blur-md px-2 py-1 rounded-lg">
                  <span class="text-yellow-400 font-bold text-lg">{{ drama.rating }}</span>
                  <span class="text-white/80 text-sm ml-1">分</span>
                </div>
                <div class="flex gap-1">
                  <span class="bg-black/60 backdrop-blur-md text-white/90 text-xs px-2 py-1 rounded-lg">{{ drama.category }}</span>
                  <span v-if="drama.sub_type" class="bg-black/60 backdrop-blur-md text-white/90 text-xs px-2 py-1 rounded-lg">{{ drama.sub_type }}</span>
                </div>
              </div>
            </div>

            <!-- 剧集信息 -->
            <div class="p-5">
              <!-- 标题 -->
              <div class="mb-3">
                <h3 class="text-base font-bold text-gray-900 dark:text-white line-clamp-2 leading-tight">
                  {{ drama.title }}
                </h3>
              </div>

              <!-- 副标题 -->
              <div v-if="drama.card_subtitle" class="mb-3">
                <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-2 leading-relaxed">{{ drama.card_subtitle }}</p>
              </div>

              <!-- 年份、地区信息 -->
              <div class="flex items-center gap-2 mb-3 flex-wrap">
                <span v-if="drama.year" class="text-xs text-white/80 bg-black/40 backdrop-blur-sm px-2 py-1 rounded-md">
                  {{ drama.year }}
                </span>
                <span v-if="drama.region" class="text-xs text-white/80 bg-black/40 backdrop-blur-sm px-2 py-1 rounded-md">
                  {{ drama.region }}
                </span>
              </div>

              <!-- 类型标签 -->
              <div v-if="drama.genres" class="mb-3">
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="genre in drama.genres.split(',').slice(0, 3)"
                    :key="genre"
                    class="text-xs text-white/90 bg-gradient-to-r from-blue-500/80 to-purple-500/80 backdrop-blur-sm px-2 py-1 rounded-md"
                  >
                    {{ genre.trim() }}
                  </span>
                </div>
              </div>
            </div>
          </a>
        </div>

        <!-- 加载更多按钮 -->
        <div v-if="filteredDramas.length > 0 && !loading && hasMore" class="mt-8 mb-4 flex justify-center">
          <button
            @click="loadMoreDramas"
            :disabled="paginationLoading"
            class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 flex items-center gap-2"
          >
            <span v-if="paginationLoading" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></span>
            <span>{{ paginationLoading ? '加载中...' : '加载更多' }}</span>
          </button>
        </div>

        <div v-if="!hasMore && filteredDramas.length > 0" class="text-center py-6 text-gray-500">
          <p>已经是全部数据了</p>
        </div>
      </div>
    </div>

   <!-- 页脚 -->
   <AppFooter />
  </div>
</template>

<script setup>
// 设置页面布局
definePageMeta({
  layout: 'default'
})

// 设置页面SEO
const { initSystemConfig, setHotDramasSeo } = useGlobalSeo()

onBeforeMount(async () => {
  await initSystemConfig()
  setHotDramasSeo()
})

const hotDramaApi = useHotDramaApi()
const { data: hotDramsaResponse, error } = await hotDramaApi.getHotDramas({
  page: 1,
  page_size: 20
})

const { getPosterUrl } = hotDramaApi

// 设置响应式数据
const dramas = ref(hotDramsaResponse.value?.items || [])
const total = ref(hotDramsaResponse.value?.total || 0)
const loading = ref(false)
const paginationLoading = ref(false)
const hasMore = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const selectedCategory = ref('')
const ssrLoadLength = ref(hotDramsaResponse.value?.items?.length || 0) // SSR加载的数据长度
let observer = null
const visibleItems = ref(new Set())

// 处理错误
if (error.value) {
  // SSR错误已在服务器端处理
}

// 分类选项
const categories = ref([
  { label: '全部', value: '' },
  { label: '热门电影', value: '电影-热门' },
  { label: '热门电视剧', value: '电视剧-热门' },
  { label: '热门综艺', value: '综艺-热门' },
  { label: '豆瓣Top250', value: '电影-Top250' }
])

// 计算属性
const filteredDramas = computed(() => {
  if (!selectedCategory.value) {
    return dramas.value
  }
  // Handle old categories
  if (selectedCategory.value === '电影') {
    return dramas.value.filter(drama => drama.category === '电影')
  }
  if (selectedCategory.value === '电视剧') {
    return dramas.value.filter(drama => drama.category === '电视剧')
  }
  // Handle new combined categories
  const [category, subType] = selectedCategory.value.split('-')
  if (subType) {
    return dramas.value.filter(drama => drama.category === category && drama.sub_type === subType)
  }
  return dramas.value
})

// 检查图片是否应该显示（SSR数据立即显示，其他数据延迟加载）
const shouldShowImage = (dramaIndex, dramaId) => {
  if (dramaIndex < ssrLoadLength.value) {
    return true
  }
  return visibleItems.value.has(dramaId)
}

const movieCount = computed(() => {
  return dramas.value.filter(drama => drama.category === '电影').length
})

const tvCount = computed(() => {
  return dramas.value.filter(drama => drama.category === '电视剧').length
})

const averageRating = computed(() => {
  const validRatings = dramas.value.filter(drama => drama.rating > 0)
  if (validRatings.length === 0) return '0.0'
  const sum = validRatings.reduce((acc, drama) => acc + drama.rating, 0)
  return (sum / validRatings.length).toFixed(1)
})

// 获取热播剧列表（重置分页）
const fetchDramas = async () => {
  loading.value = true
  try {
    // 解析分类参数，分割为category和sub_type
    const params = {
      page: 1,
      page_size: pageSize.value
    }

    if (selectedCategory.value) {
      const [category, subType] = selectedCategory.value.split('-')
      params.category = category
      if (subType) {
        params.sub_type = subType
      }
    }

    // 使用客户端版本的API
    const response = await hotDramaApi.getHotDramasClient(params)

    if (response && response.items) {
      dramas.value = response.items
      total.value = response.total || 0
      currentPage.value = 1
      hasMore.value = response.items.length === pageSize.value
      ssrLoadLength.value = response.items.length
      visibleItems.value.clear()
      nextTick(() => {
        initIntersectionObserver()
      })
    } else {
      dramas.value = Array.isArray(response) ? response : []
      total.value = dramas.value.length
      hasMore.value = false
      ssrLoadLength.value = dramas.value.length
      visibleItems.value.clear()
    }
  } catch (error) {
    dramas.value = []
    total.value = 0
    hasMore.value = false
  } finally {
    loading.value = false
  }
}

// 加载更多数据（按钮方式）
const loadMoreDramas = async () => {
  if (paginationLoading.value || !hasMore.value) return

  paginationLoading.value = true
  try {
    const nextPage = currentPage.value + 1
    // 解析分类参数，分割为category和sub_type
    const params = {
      page: nextPage,
      page_size: pageSize.value
    }

    if (selectedCategory.value) {
      const [category, subType] = selectedCategory.value.split('-')
      params.category = category
      if (subType) {
        params.sub_type = subType
      }
    }

    const response = await hotDramaApi.getHotDramasClient(params)

    if (response && response.items && response.items.length > 0) {
      dramas.value = [...dramas.value, ...response.items]
      currentPage.value = nextPage
      hasMore.value = response.items.length === pageSize.value
      nextTick(() => {
        initIntersectionObserver()
      })
    } else {
      hasMore.value = false
    }
  } catch (error) {
    // 加载更多剧集失败
    hasMore.value = false
  } finally {
    paginationLoading.value = false
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 格式化数字
const formatNumber = (num) => {
  if (!num) return '0'
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toString()
}

const initIntersectionObserver = () => {
  if (observer) observer.disconnect()

  observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      const itemId = entry.target.getAttribute('data-item-id')
      const itemIndex = parseInt(entry.target.getAttribute('data-item-index'))

      if (!itemId || !itemIndex) return
      if (itemIndex >= ssrLoadLength.value && entry.isIntersecting && entry.intersectionRatio > 0.01) {
        visibleItems.value.add(Number(itemId))
        observer.unobserve(entry.target)
      }
    })
  }, {
    root: null,
    rootMargin: '200px 0px 200px 0px',
    threshold: [0.01, 0.1, 0.5]
  })

  // 只观察分页加载的数据
  nextTick(() => {
    const cards = document.querySelectorAll('[data-item-index]')
    cards.forEach((card) => {
      const itemIndex = parseInt(card.getAttribute('data-item-index'))
      if (itemIndex >= ssrLoadLength.value) {
        observer?.observe(card)
      }
    })
  })
}

// 处理图片加载错误 - 显示占位图
const handleImageError = (event) => {
  // 设置占位图片
  event.target.src = 'data:image/svg+xml;base64,' + btoa(`
    <svg width="400" height="208" xmlns="http://www.w3.org/2000/svg">
      <rect width="100%" height="100%" fill="#374151"/>
      <text x="50%" y="50%" font-family="Arial" font-size="14" fill="#9CA3AF" text-anchor="middle" dy=".35em">暂无封面</text>
    </svg>
  `)
  event.target.style.background = '#374151'
}

// 清理Intersection Observer
const cleanupObserver = () => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
}

watch(selectedCategory, () => {
  currentPage.value = 1
  hasMore.value = true
  fetchDramas()
})

onMounted(() => {
  if (dramas.value.length === 0) {
    fetchDramas()
  } else {
    nextTick(() => {
      initIntersectionObserver()
    })
  }
})

onBeforeUnmount(() => {
  cleanupObserver()
})

</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.header-container{
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}
</style> 