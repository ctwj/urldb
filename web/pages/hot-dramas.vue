<template>
  <div class="min-h-screen bg-gray-50">
    <div class="container mx-auto px-4 py-8">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">热播剧榜单</h1>
        <p class="text-gray-600">实时获取豆瓣热门电影和电视剧榜单</p>
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
              : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
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
      <div v-else-if="dramas.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div
          v-for="drama in filteredDramas"
          :key="drama.id"
          class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
        >
          <!-- 剧集信息 -->
          <div class="p-6">
            <div class="flex items-start justify-between mb-3">
              <h3 class="text-lg font-semibold text-gray-900 line-clamp-2">
                {{ drama.title }}
              </h3>
              <div class="flex items-center ml-2">
                <span class="text-yellow-500 text-sm font-medium">{{ drama.rating }}</span>
                <span class="text-gray-400 text-xs ml-1">分</span>
              </div>
            </div>

            <!-- 年份和分类 -->
            <div class="flex items-center gap-2 mb-3">
              <span v-if="drama.year" class="text-sm text-gray-500">{{ drama.year }}</span>
              <span class="text-sm text-blue-600 bg-blue-100 px-2 py-1 rounded">
                {{ drama.category }}
              </span>
              <span v-if="drama.sub_type" class="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">
                {{ drama.sub_type }}
              </span>
            </div>

            <!-- 导演 -->
            <div v-if="drama.directors" class="mb-2">
              <span class="text-xs text-gray-500">导演：</span>
              <span class="text-sm text-gray-700">{{ drama.directors }}</span>
            </div>

            <!-- 演员 -->
            <div v-if="drama.actors" class="mb-3">
              <span class="text-xs text-gray-500">主演：</span>
              <span class="text-sm text-gray-700 line-clamp-2">{{ drama.actors }}</span>
            </div>

            <!-- 数据来源 -->
            <div class="flex items-center justify-between text-xs text-gray-400">
              <span>来源：{{ drama.source }}</span>
              <span>{{ formatDate(drama.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">暂无热播剧数据</h3>
        <p class="text-gray-500">请稍后再试或联系管理员</p>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="mt-8 flex justify-center">
        <nav class="flex items-center space-x-2">
          <button
            @click="changePage(currentPage - 1)"
            :disabled="currentPage <= 1"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            上一页
          </button>
          <span class="px-3 py-2 text-sm text-gray-700">
            第 {{ currentPage }} 页，共 {{ Math.ceil(total / pageSize) }} 页
          </span>
          <button
            @click="changePage(currentPage + 1)"
            :disabled="currentPage >= Math.ceil(total / pageSize)"
            class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一页
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'

// 响应式数据
const loading = ref(false)
const dramas = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
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

// 获取热播剧列表
const fetchDramas = async () => {
  loading.value = true
  try {
    const { useHotDramaApi } = await import('~/composables/useApi')
    const hotDramaApi = useHotDramaApi()
    
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    const response = await hotDramaApi.getHotDramas(params)
    
    // 使用新的统一响应格式
    if (response && response.items) {
      dramas.value = response.items
      total.value = response.total || 0
    } else {
      // 兼容旧格式
      dramas.value = Array.isArray(response) ? response : []
      total.value = dramas.value.length
    }
  } catch (error) {
    console.error('获取热播剧列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 切换页面
const changePage = (page) => {
  if (page >= 1 && page <= Math.ceil(total.value / pageSize.value)) {
    currentPage.value = page
    fetchDramas()
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 监听分类变化
watch(selectedCategory, () => {
  currentPage.value = 1
  fetchDramas()
})

// 页面加载时获取数据
onMounted(() => {
  fetchDramas()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 