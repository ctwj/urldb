<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">搜索统计</h1>
        <p class="text-gray-600 dark:text-gray-400">查看搜索量统计和热门关键词分析</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="refreshData" type="primary">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <n-card>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400">
            <i class="fas fa-search text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日搜索</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.todaySearches || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-green-100 dark:bg-green-900 text-green-600 dark:text-green-400">
            <i class="fas fa-chart-line text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本周搜索</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.weekSearches || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-purple-100 dark:bg-purple-900 text-purple-600 dark:text-purple-400">
            <i class="fas fa-calendar text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本月搜索</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.monthSearches || 0 }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 搜索趋势图表 -->
    <n-card>
      <template #header>
        <span class="text-xl font-semibold text-gray-900 dark:text-white">搜索趋势</span>
      </template>
      <div class="h-64">
        <canvas ref="trendChart"></canvas>
      </div>
    </n-card>

    <!-- 热门关键词 -->
    <n-card>
      <template #header>
        <span class="text-xl font-semibold text-gray-900 dark:text-white">热门关键词</span>
      </template>
      <div class="space-y-4">
        <div v-for="keyword in stats.hotKeywords" :key="keyword.keyword" 
             class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
          <div class="flex items-center">
            <span class="inline-flex items-center justify-center w-8 h-8 bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400 rounded-full text-sm font-medium mr-3">
              {{ keyword.rank }}
            </span>
            <span class="text-gray-900 dark:text-white font-medium">{{ keyword.keyword }}</span>
          </div>
          <div class="flex items-center">
            <span class="text-gray-600 dark:text-gray-400 mr-2">{{ keyword.count }}次</span>
            <div class="w-24 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
              <div class="bg-blue-600 h-2 rounded-full" 
                   :style="{ width: getPercentage(keyword.count) + '%' }"></div>
            </div>
          </div>
        </div>
        <div v-if="!stats.hotKeywords || stats.hotKeywords.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
          暂无热门关键词数据
        </div>
      </div>
    </n-card>

    <!-- 搜索记录 -->
    <n-card>
      <template #header>
        <span class="text-xl font-semibold text-gray-900 dark:text-white">搜索记录</span>
      </template>
      <n-data-table
        :columns="columns"
        :data="searchList"
        :pagination="pagination"
        :loading="loading"
        :bordered="false"
        striped
      />
      <div v-if="searchList.length === 0 && !loading" class="text-center py-8 text-gray-500 dark:text-gray-400">
        暂无搜索记录
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: ['auth']
})

import { ref, onMounted, computed, nextTick } from 'vue'
import Chart from 'chart.js/auto'
import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'

// 响应式数据
const stats = ref<{
  todaySearches: number
  weekSearches: number
  monthSearches: number
  hotKeywords: Array<{
    keyword: string
    count: number
    rank: number
  }>
  searchTrend: {
    days: string[]
    values: number[]
  }
}>({
  todaySearches: 0,
  weekSearches: 0,
  monthSearches: 0,
  hotKeywords: [],
  searchTrend: {
    days: [],
    values: []
  }
})

const searchList = ref<Array<{
  id: number
  keyword: string
  count: number
  date: string
  created_at: string
}>>([])
const loading = ref(false)
const trendChart = ref<HTMLCanvasElement | null>(null)
let chart: any = null

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    pagination.value.page = page
    loadSearchRecords()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    loadSearchRecords()
  }
})

// 表格列配置
const columns = [
  {
    title: '关键词',
    key: 'keyword',
    width: 200
  },
  {
    title: '搜索次数',
    key: 'count',
    width: 120
  },
  {
    title: '日期',
    key: 'date',
    width: 150,
    render: (row: any) => {
      return row.date ? new Date(row.date).toLocaleDateString() : ''
    }
  }
]

// 获取百分比
const getPercentage = (count: number) => {
  if (!stats.value.hotKeywords || stats.value.hotKeywords.length === 0) return 0
  const maxCount = Math.max(...stats.value.hotKeywords.map((k: any) => k.count))
  return Math.round((count / maxCount) * 100)
}

// 加载搜索统计
const loadSearchStats = async () => {
  try {
    loading.value = true
    
    // 1. 汇总卡片
    const summary = await useApiFetch('/search-stats/summary').then(parseApiResponse) as any
    stats.value.todaySearches = summary?.today || 0
    stats.value.weekSearches = summary?.week || 0
    stats.value.monthSearches = summary?.month || 0
    
    // 2. 热门关键词
    const hotKeywords = await useApiFetch('/search-stats/hot-keywords').then(parseApiResponse) as any[]
    stats.value.hotKeywords = hotKeywords || []
    
    // 3. 趋势
    const trend = await useApiFetch('/search-stats/trend').then(parseApiResponse) as any[]
    stats.value.searchTrend.days = (trend || []).map((item: any) => item.date ? new Date(item.date).toLocaleDateString() : '')
    stats.value.searchTrend.values = (trend || []).map((item: any) => item.total_searches)
    
    // 4. 更新图表
    await nextTick()
    updateChart()
  } catch (error) {
    console.error('加载搜索统计失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载搜索记录
const loadSearchRecords = async () => {
  try {
    loading.value = true
    const response = await useApiFetch('/search-stats').then(parseApiResponse) as any
    searchList.value = response?.data || []
  } catch (error) {
    console.error('加载搜索记录失败:', error)
  } finally {
    loading.value = false
  }
}

// 更新图表
const updateChart = () => {
  if (chart) {
    chart.destroy()
  }

  if (!trendChart.value) return

  const ctx = trendChart.value.getContext('2d')
  if (!ctx) return

  chart = new Chart(ctx as any, {
    type: 'line',
    data: {
      labels: stats.value.searchTrend.days,
      datasets: [{
        label: '搜索量',
        data: stats.value.searchTrend.values,
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
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

// 刷新数据
const refreshData = () => {
  loadSearchStats()
  loadSearchRecords()
}

// 初始化
onMounted(() => {
  loadSearchStats()
  loadSearchRecords()
})
</script> 