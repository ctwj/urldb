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

        <!-- 统计信息 -->
        <div class="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                <i class="fas fa-film text-blue-600 dark:text-blue-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-500 dark:text-gray-400">总数量</p>
                <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ total }}</p>
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
            <div class="flex items-center">
              <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                <i class="fas fa-video text-green-600 dark:text-green-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-500 dark:text-gray-400">电影</p>
                <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ movieCount }}</p>
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <i class="fas fa-tv text-purple-600 dark:text-purple-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-500 dark:text-gray-400">电视剧</p>
                <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ tvCount }}</p>
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
            <div class="flex items-center">
              <div class="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                <i class="fas fa-star text-yellow-600 dark:text-yellow-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-500 dark:text-gray-400">平均评分</p>
                <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ averageRating }}</p>
              </div>
            </div>
          </div>
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
          <div
            v-for="drama in filteredDramas"
            :key="drama.id"
            class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow border border-gray-200 dark:border-gray-700"
          >
            <!-- 剧集信息 -->
            <div class="p-6">
              <div class="flex items-start justify-between mb-3">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white line-clamp-2 flex-1">
                  {{ drama.title }}
                </h3>
                <div class="flex items-center ml-2 flex-shrink-0">
                  <span class="text-yellow-500 text-sm font-medium">{{ drama.rating }}</span>
                  <span class="text-gray-400 dark:text-gray-500 text-xs ml-1">分</span>
                </div>
              </div>

              <!-- 副标题 -->
              <div v-if="drama.card_subtitle" class="mb-3">
                <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-1">{{ drama.card_subtitle }}</p>
              </div>

              <!-- 年份、地区、类型 -->
              <div class="flex items-center gap-2 mb-3 flex-wrap">
                <span v-if="drama.year" class="text-sm text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                  {{ drama.year }}
                </span>
                <span v-if="drama.region" class="text-sm text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                  {{ drama.region }}
                </span>
                <span class="text-sm text-blue-600 dark:text-blue-400 bg-blue-100 dark:bg-blue-900 px-2 py-1 rounded">
                  {{ drama.category }}
                </span>
                <span v-if="drama.sub_type" class="text-sm text-purple-600 dark:text-purple-400 bg-purple-100 dark:bg-purple-900 px-2 py-1 rounded">
                  {{ drama.sub_type }}
                </span>
              </div>

              <!-- 类型标签 -->
              <div v-if="drama.genres" class="mb-3">
                <div class="flex flex-wrap gap-1">
                  <span 
                    v-for="genre in drama.genres.split(',')" 
                    :key="genre"
                    class="text-xs text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded"
                  >
                    {{ genre.trim() }}
                  </span>
                </div>
              </div>

              <!-- 导演 -->
              <div v-if="drama.directors" class="mb-2">
                <span class="text-xs text-gray-500 dark:text-gray-400">导演：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300 line-clamp-1">{{ drama.directors }}</span>
              </div>

              <!-- 演员 -->
              <div v-if="drama.actors" class="mb-3">
                <span class="text-xs text-gray-500 dark:text-gray-400">主演：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300 line-clamp-2">{{ drama.actors }}</span>
              </div>

              <!-- 集数信息 -->
              <div v-if="drama.episodes_info" class="mb-3">
                <span class="text-xs text-gray-500 dark:text-gray-400">集数：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ drama.episodes_info }}</span>
              </div>

              <!-- 评分人数 -->
              <div v-if="drama.rating_count" class="mb-3">
                <span class="text-xs text-gray-500 dark:text-gray-400">评分人数：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ formatNumber(drama.rating_count) }}</span>
              </div>

              <!-- 数据来源和时间 -->
              <div class="flex items-center justify-between text-xs text-gray-400 dark:text-gray-500 pt-3 border-t border-gray-200 dark:border-gray-600">
                <span>来源：{{ drama.source }}</span>
                <span>{{ formatDate(drama.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-12">
          <div class="flex flex-col items-center justify-center space-y-4">
            <img 
              src="/assets/svg/empty.svg" 
              alt="暂无热播剧数据" 
              class="!w-64 !h-64 sm:w-64 sm:h-64 opacity-60 dark:opacity-40"
            />
            <div class="text-center">
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无热播剧数据</h3>
              <p class="text-gray-500 dark:text-gray-400">请稍后再试或联系管理员</p>
            </div>
          </div>
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

import { ref, computed, onMounted, watch } from 'vue'
import { useHotDramaApi } from '~/composables/useApi'
const hotDramaApi = useHotDramaApi()

// 响应式数据
const loading = ref(false)
const dramas = ref([])
const total = ref(0)
const selectedCategory = ref('')

// 分类选项
const categories = ref([
  { label: '全部', value: '' },
  { label: '电影', value: '电影' },
  { label: '电视剧', value: '电视剧' }
])

// 计算属性
const filteredDramas = computed(() => {
  if (!selectedCategory.value) {
    return dramas.value
  }
  return dramas.value.filter(drama => drama.category === selectedCategory.value)
})

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

// 获取热播剧列表（获取所有数据）
const fetchDramas = async () => {
  loading.value = true
  try {
    const params = {
      page: 1,
      page_size: 1000
    }
    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }
    const response = await hotDramaApi.getHotDramas(params)
    if (response && response.items) {
      dramas.value = response.items
      total.value = response.total || 0
    } else {
      dramas.value = Array.isArray(response) ? response : []
      total.value = dramas.value.length
    }
  } catch (error) {
    console.error('获取热播剧列表失败:', error)
  } finally {
    loading.value = false
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

// 处理图片加载错误
const handleImageError = (event) => {
  console.log('图片加载失败:', event.target.src)
  event.target.style.display = 'none'
}

// 处理图片加载成功
const handleImageLoad = (event) => {
  console.log('图片加载成功:', event.target.src)
}

// 监听分类变化
watch(selectedCategory, () => {
  fetchDramas()
})

// 页面加载时获取数据
onMounted(() => {
  console.log('热播剧页面加载')
  fetchDramas()
})

// 监听数据变化
watch(dramas, (newDramas) => {
  console.log('dramas数据变化:', newDramas?.length)
  if (newDramas && newDramas.length > 0) {
    console.log('第一条数据:', newDramas[0])
    console.log('第一条数据的poster_url:', newDramas[0].poster_url)
  }
}, { deep: true })
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