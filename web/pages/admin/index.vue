<template>
  <div class="space-y-8">
    <!-- 欢迎区域 -->
    <n-card class="bg-gradient-to-r from-blue-500 to-purple-600 text-white border-0 relative">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold mb-2 text-white">
            欢迎回来，管理员！
          </h1>
          <p class="text-blue-100 text-lg">
            这里是管理后台，您可以管理所有系统资源和配置。
          </p>
        </div>
        <div class="hidden md:block">
          <div class="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center">
            <i class="fas fa-shield-alt text-2xl text-white"></i>
          </div>
        </div>
      </div>
      

    </n-card>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <!-- 今日资源/总资源数 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日资源/总资源数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.today_resources || 0 }}/{{ stats.total_resources || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="primary" @click="navigateTo('/admin/resources')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- 今日浏览量 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-orange-100 dark:bg-orange-900 rounded-lg">
            <i class="fas fa-eye text-orange-600 dark:text-orange-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日浏览量</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.today_views || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="warning" @click="navigateTo('/admin/search-stats')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- 今日搜索量 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-search text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日搜索量</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.today_searches || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="success" @click="navigateTo('/admin/search-stats')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>

      <!-- Google索引统计 -->
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fab fa-google text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">已索引URL/总URL</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ googleIndexStats.indexedURLs || 0 }}/{{ googleIndexStats.totalURLs || 0 }}</p>
          </div>
        </div>
        <template #footer>
          <n-button text type="primary" @click="navigateTo('/admin/seo#google-index')">
            查看详情
            <template #icon>
              <i class="fas fa-arrow-right"></i>
            </template>
          </n-button>
        </template>
      </n-card>
    </div>

    <!-- 平台管理 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="flex items-center mb-4">
        <div class="p-3 bg-green-100 rounded-lg">
          <i class="fas fa-server text-green-600 text-xl"></i>
        </div>
        <div class="ml-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">平台列表</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400">系统支持的网盘平台</p>
        </div>
      </div>
      <div class="space-y-2">
        <div class="flex flex-wrap gap-1 w-full text-left rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer">
          <div v-for="pan in pans" :key="pan.id" class="h-6 px-1 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
            <span v-html="pan.icon"></span>&nbsp;{{ pan.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- 统计趋势图表区域 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- 访问量趋势图 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">访问量趋势</h3>
          <div class="p-2 bg-orange-100 dark:bg-orange-900 rounded-full">
            <i class="fas fa-chart-line text-orange-600 dark:text-orange-400 text-sm"></i>
          </div>
        </div>
        <div class="h-40">
          <canvas ref="viewsChart"></canvas>
        </div>
      </div>

      <!-- 搜索量趋势图 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">搜索量趋势</h3>
          <div class="p-2 bg-green-100 dark:bg-green-900 rounded-full">
            <i class="fas fa-chart-line text-green-600 dark:text-green-400 text-sm"></i>
          </div>
        </div>
        <div class="h-40">
          <canvas ref="searchesChart"></canvas>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

import { useStatsApi, usePanApi } from '~/composables/useApi'
import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'
import Chart from 'chart.js/auto'

// Google索引统计
const googleIndexStats = ref({
  totalURLs: 0,
  indexedURLs: 0,
  errorURLs: 0,
  totalTasks: 0,
  runningTasks: 0,
  completedTasks: 0,
  failedTasks: 0
})

// API
const statsApi = useStatsApi()
const panApi = usePanApi()

// 获取统计数据
const { data: statsData } = await useAsyncData('adminStats', () => statsApi.getStats())

// 获取平台数据
const { data: pansData } = await useAsyncData('adminPans', () => panApi.getPans())

// 统计数据
const stats = computed(() => {
  console.log('原始统计数据:', statsData.value)
  const result = (statsData.value as any) || { 
    total_resources: 0, 
    today_resources: 0,
    today_views: 0,
    today_searches: 0
  }
  console.log('处理后的统计数据:', result)
  return result
})

// 平台数据
const pans = computed(() => (pansData.value as any) || [])

// 趋势图数据
const weeklyViews = ref<any[]>([])
const weeklySearches = ref<any[]>([])

// 获取趋势数据
const fetchTrendData = async () => {
  console.log('开始获取趋势数据...')
  try {
    // 获取访问量趋势数据
    const viewsResponse = await useApiFetch('/stats/views-trend').then(parseApiResponse)
    if (viewsResponse && Array.isArray(viewsResponse)) {
      weeklyViews.value = viewsResponse.map((item: any) => ({
        label: item.date ? new Date(item.date).toLocaleDateString('zh-CN', { weekday: 'short' }) : '',
        value: item.views || 0
      }))
    } else {
      // 如果没有数据，使用模拟数据
      weeklyViews.value = [
        { label: '周一', value: 0 },
        { label: '周二', value: 0 },
        { label: '周三', value: 0 },
        { label: '周四', value: 0 },
        { label: '周五', value: 0 },
        { label: '周六', value: 0 },
        { label: '周日', value: 0 }
      ]
    }

    // 获取搜索量趋势数据
    const searchesResponse = await useApiFetch('/stats/searches-trend').then(parseApiResponse)
    console.log('搜索量趋势API响应:', searchesResponse)
    if (searchesResponse && Array.isArray(searchesResponse)) {
      weeklySearches.value = searchesResponse.map((item: any) => ({
        label: item.date ? new Date(item.date).toLocaleDateString('zh-CN', { weekday: 'short' }) : '',
        value: item.searches || 0
      }))
      console.log('处理后的搜索量数据:', weeklySearches.value)
    } else {
      // 如果没有数据，使用模拟数据
      weeklySearches.value = [
        { label: '周一', value: 45 },
        { label: '周二', value: 52 },
        { label: '周三', value: 67 },
        { label: '周四', value: 58 },
        { label: '周五', value: 73 },
        { label: '周六', value: 89 },
        { label: '周日', value: 81 }
      ]
      console.log('使用模拟搜索量数据:', weeklySearches.value)
    }
  } catch (error) {
    console.error('获取趋势数据失败:', error)
    // 使用默认模拟数据
    weeklyViews.value = [
      { label: '周一', value: 0 },
      { label: '周二', value: 0 },
      { label: '周三', value: 0 },
      { label: '周四', value: 0 },
      { label: '周五', value: 0 },
      { label: '周六', value: 0 },
      { label: '周日', value: 0 }
    ]
    weeklySearches.value = [
      { label: '周一', value: 0 },
      { label: '周二', value: 0 },
      { label: '周三', value: 0 },
      { label: '周四', value: 0 },
      { label: '周五', value: 0 },
      { label: '周六', value: 0 },
      { label: '周日', value: 0 }
    ]
  }
}

// 图表引用
const viewsChart = ref<HTMLCanvasElement | null>(null)
const searchesChart = ref<HTMLCanvasElement | null>(null)

// 图表实例
let viewsChartInstance: any = null
let searchesChartInstance: any = null

// 初始化图表
const initCharts = () => {
  console.log('开始初始化图表...')
  console.log('访问量数据:', weeklyViews.value)
  console.log('搜索量数据:', weeklySearches.value)
  
  // 访问量趋势图
  if (viewsChart.value) {
    console.log('访问量图表容器存在:', viewsChart.value)
    if (viewsChartInstance) {
      viewsChartInstance.destroy()
    }
    const ctx = viewsChart.value.getContext('2d')
    if (ctx) {
      console.log('访问量图表上下文获取成功')
      viewsChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
          labels: weeklyViews.value.map(day => day.label),
          datasets: [{
            label: '访问量',
            data: weeklyViews.value.map(day => day.value),
            borderColor: 'rgb(249, 115, 22)',
            backgroundColor: 'rgba(249, 115, 22, 0.1)',
            tension: 0.4,
            fill: true
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: false
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: 'rgba(0, 0, 0, 0.1)'
              }
            },
            x: {
              grid: {
                color: 'rgba(0, 0, 0, 0.1)'
              }
            }
          }
        }
      })
    }
  }

  // 搜索量趋势图
  if (searchesChart.value) {
    console.log('搜索量图表容器存在:', searchesChart.value)
    if (searchesChartInstance) {
      searchesChartInstance.destroy()
    }
    const ctx = searchesChart.value.getContext('2d')
    if (ctx) {
      console.log('搜索量图表上下文获取成功')
      console.log('初始化搜索量图表，数据:', weeklySearches.value)
      searchesChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
          labels: weeklySearches.value.map(day => day.label),
          datasets: [{
            label: '搜索量',
            data: weeklySearches.value.map(day => day.value),
            borderColor: 'rgb(34, 197, 94)',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            tension: 0.4,
            fill: true
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: false
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: 'rgba(0, 0, 0, 0.1)'
              }
            },
            x: {
              grid: {
                color: 'rgba(0, 0, 0, 0.1)'
              }
            }
          }
        }
      })
    }
  }
}



// 监听数据变化，更新图表
watch([weeklyViews, weeklySearches], () => {
  nextTick(() => {
    initCharts()
  })
})

// 加载Google索引统计
const loadGoogleIndexStats = async () => {
  try {
    const response = await useApiFetch('/google-index/status').then(parseApiResponse)
    if (response && response.data) {
      googleIndexStats.value = response.data
    }
  } catch (error) {
    console.error('加载Google索引统计失败:', error)
  }
}

// 组件挂载后初始化图表和数据
onMounted(async () => {
  console.log('组件挂载，开始初始化...')
  await fetchTrendData()
  await loadGoogleIndexStats() // 加载Google索引统计
  console.log('数据获取完成，准备初始化图表...')
  nextTick(() => {
    console.log('nextTick执行，开始初始化图表...')
    initCharts()
  })
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 