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
          <div
            v-for="drama in filteredDramas"
            :key="drama.id"
            :data-item-id="drama.id"
            class="group relative bg-white/10 dark:bg-gray-800/10 backdrop-blur-md rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 border border-white/20 dark:border-gray-700/50 hover:scale-105"
          >
            <!-- 海报图片 -->
            <div v-if="drama.poster_url" class="relative overflow-hidden">
              <!-- 品牌Logo占位符 -->
              <div
                v-if="!visibleItems.has(drama.id)"
                class="w-full h-52 bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-700 dark:to-gray-800 flex flex-col items-center justify-center relative overflow-hidden"
              >
                <!-- 装饰性背景图形 -->
                <div class="absolute inset-0 opacity-20">
                  <svg viewBox="0 0 200 100" class="w-full h-full">
                    <circle cx="30" cy="25" r="3" fill="currentColor"/>
                    <circle cx="80" cy="40" r="2" fill="currentColor"/>
                    <circle cx="150" cy="20" r="2" fill="currentColor"/>
                    <circle cx="120" cy="60" r="2" fill="currentColor"/>
                    <circle cx="50" cy="70" r="2" fill="currentColor"/>
                  </svg>
                </div>

                <!-- 主要品牌元素 -->
                <div class="flex flex-col items-center space-y-2 z-10">
                  <!-- 电影院图标 -->
                  <svg class="w-12 h-12 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M19 6c0-1.1-.9-2-2-2H7c-1.1 0-2 .9-2 2v1l.669.775C6.537 8.347 7.605 9.334 9.5 9.781c.015 0 .03.003.045.003s.03-.003.045-.003c1.895-.447 2.963-1.434 3.331-1.506L13 7V6h6v1l.669.775C20.537 8.347 21.605 9.334 23.5 9.781c.015 0 .03.003.045.003s.03-.003.045-.003c1.895-.447 2.963-1.434 3.331-1.506L5 7V6H1c0 1.1.9 2 2 2v1l.669.775C4.537 8.347 5.605 9.334 7.5 9.781c.015 0 .03.003.045.003s.03-.003.045-.003c1.895-.447 2.963-1.434 3.331-1.506L13 7V18H7c-1.1 0-2 .9-2 2s.9 2 2 2h10c1.1 0 2-.9 2-2s-.9-2-2-2h-6V6z"/>
                  </svg>

                  <!-- 品牌文字 -->
                  <div class="text-center">
                    <div class="text-lg font-bold text-blue-600 dark:text-blue-400">热播剧榜单</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400 animate-pulse">精彩剧集等你发现</div>
                  </div>

                  <!-- 装饰线条 -->
                  <div class="flex items-center space-x-2">
                    <div class="w-8 h-px bg-blue-300 dark:bg-blue-600"></div>
                    <div class="w-2 h-2 bg-blue-400 dark:bg-blue-500 rounded-full"></div>
                    <div class="w-8 h-px bg-blue-300 dark:bg-blue-600"></div>
                  </div>
                </div>
              </div>
              <!-- 主图片（只有在可视区域时才加载） -->
              <img
                v-if="visibleItems.has(drama.id)"
                :src="getPosterUrl(drama.poster_url)"
                :alt="drama.title"
                class="w-full h-52 object-cover transition-all duration-500 opacity-0"
                @load="$event.target.style.opacity = '1'"
                @error="handleImageError"
              />
              <!-- 图片上的遮罩和信息（始终显示） -->
              <div v-if="visibleItems.has(drama.id)" class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>

              <!-- 新剧标签 -->
              <div
                v-if="drama.is_new && visibleItems.has(drama.id)"
                class="absolute top-3 right-3 bg-gradient-to-r from-red-500 to-red-600 text-white px-3 py-1 rounded-full text-xs font-semibold shadow-lg z-10"
              >
                🔥 HOT
              </div>

              <!-- 评分显示 -->
              <div v-if="visibleItems.has(drama.id)" class="absolute bottom-3 left-3 right-3 flex items-center justify-between z-20">
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

              <!-- 数据来源和时间 -->
              <!-- <div class="flex items-center justify-between text-xs pt-4 border-t border-gray-100 dark:border-gray-700/50">
                <div class="flex items-center gap-2">
                  <span class="text-gray-500 dark:text-gray-400">{{ drama.source }}</span>
                  <div class="w-1 h-1 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
                  <span class="text-gray-400 dark:text-gray-500">{{ formatDate(drama.created_at) }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-green-600 dark:text-green-400 font-medium">{{ drama.episodes_info || '更新中' }}</span>
                  <div class="w-1 h-1 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
                  <a
                    v-if="drama.douban_uri"
                    :href="drama.douban_uri"
                    target="_blank"
                    class="bg-gradient-to-r from-blue-500 to-blue-600 text-white px-3 py-1 rounded-full text-xs font-medium hover:from-blue-600 hover:to-blue-700 transition-all duration-200"
                    @click.stop
                  >
                    View
                  </a>
                </div>
              </div> -->
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

import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useHotDramaApi } from '~/composables/useApi'
const hotDramaApi = useHotDramaApi()
const { getPosterUrl } = hotDramaApi

// 响应式数据
const loading = ref(false)
const dramas = ref([])
const total = ref(0)
const selectedCategory = ref('')
const visibleItems = ref(new Set()) // 存储当前可视区域的项目ID

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

// 处理图片加载错误 - 显示占位图
const handleImageError = (event) => {
  console.log('图片加载失败:', event.target.src)
  // 设置占位图片
  event.target.src = 'data:image/svg+xml;base64,' + btoa(`
    <svg width="400" height="208" xmlns="http://www.w3.org/2000/svg">
      <rect width="100%" height="100%" fill="#374151"/>
      <text x="50%" y="50%" font-family="Arial" font-size="14" fill="#9CA3AF" text-anchor="middle" dy=".35em">暂无封面</text>
    </svg>
  `)
  event.target.style.background = '#374151'
}

// 处理图片加载成功
const handleImageLoad = (event) => {
  console.log('图片加载成功:', event.target.src)
}

// 监听分类变化
watch(selectedCategory, () => {
  visibleItems.value.clear() // 清空可见项目集合
  fetchDramas()
})

// 页面加载时获取数据
onMounted(() => {
  console.log('热播剧页面加载')
  fetchDramas()
})

// Intersection Observer 用于懒加载图片
let observer = null
const initIntersectionObserver = () => {
  if (observer) observer.disconnect()

  observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      const itemId = entry.target.getAttribute('data-item-id')
      if (!itemId) return

      if (entry.isIntersecting) {
        // 元素进入视窗，添加到可见集合
        visibleItems.value.add(Number(itemId))
      } else {
        // 元素离开视窗，如果需要可以移除
        // visibleItems.value.delete(Number(itemId)) // 可选，如果想重复懒加载
      }
    })
  }, {
    rootMargin: '100px 0px 100px 0px', // 提前100px和延后100px
    threshold: 0.1
  })

  // 观察所有卡片
  nextTick(() => {
    const cards = document.querySelectorAll('[data-item-id]')
    cards.forEach(card => {
      observer?.observe(card)
    })
  })
}

const cleanupObserver = () => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
}

// 监听数据变化
watch(dramas, (newDramas) => {
  console.log('dramas数据变化:', newDramas?.length)
  if (newDramas && newDramas.length > 0) {
    console.log('第一条数据:', newDramas[0])
    console.log('第一条数据的poster_url:', newDramas[0].poster_url)

    visibleItems.value.clear()

    // 延迟一帧后初始化观察器
    nextTick(() => {
      initIntersectionObserver()
    })
  }
}, { deep: true })

// 页面加载时获取数据
onMounted(() => {
  console.log('热播剧页面加载')
  fetchDramas()
})

// 页面卸载时清理观察器
onUnmounted(() => {
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