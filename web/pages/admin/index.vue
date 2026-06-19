<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">仪表盘</h1>
      <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">系统运行状态概览</p>
    </div>

    <!-- 加载态 -->
    <div v-if="pending" class="flex justify-center py-12">
      <n-spin size="large" />
    </div>

    <!-- 错误态 -->
    <ErrorState
      v-else-if="error"
      icon="fas fa-exclamation-circle"
      :message="error.message || '加载数据失败'"
      :on-retry="refresh"
    />

    <!-- 数据展示 -->
    <template v-else-if="summary">
      <!-- 指标卡片（含环比昨日） -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard
          label="今日新增资源"
          icon="fas fa-database"
          color="blue"
          :today="summary.resources.today"
          :yesterday="summary.resources.yesterday"
          :total="summary.resources.total"
          to="/admin/resources"
        />
        <StatCard
          label="今日浏览量"
          icon="fas fa-eye"
          color="orange"
          :today="summary.views.today"
          :yesterday="summary.views.yesterday"
          to="/admin/search-stats"
        />
        <StatCard
          label="今日搜索量"
          icon="fas fa-search"
          color="green"
          :today="summary.searches.today"
          :yesterday="summary.searches.yesterday"
          to="/admin/search-stats"
        />
      </div>

      <!-- 待办聚合区 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center mb-4">
          <div class="p-2 bg-amber-100 dark:bg-amber-900/40 rounded-lg">
            <i class="fas fa-tasks text-amber-600 dark:text-amber-400"></i>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white ml-3 dark:text-gray-900">待办事项</h3>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <NuxtLink to="/admin/ready-resources" class="flex items-center justify-between p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-600 transition-colors cursor-pointer">
            <span class="text-sm text-gray-600 dark:text-gray-400">待处理资源</span>
            <span class="text-xl font-bold" :class="summary.todos.ready_resources > 0 ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-400'">{{ summary.todos.ready_resources }}</span>
          </NuxtLink>
          <NuxtLink to="/admin/tasks" class="flex items-center justify-between p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-red-300 dark:hover:border-red-600 transition-colors cursor-pointer">
            <span class="text-sm text-gray-600 dark:text-gray-400">失败任务</span>
            <span class="text-xl font-bold" :class="summary.todos.failed_tasks > 0 ? 'text-red-600 dark:text-red-400' : 'text-gray-400 dark:text-gray-400'">{{ summary.todos.failed_tasks }}</span>
          </NuxtLink>
          <NuxtLink to="/admin/reports" class="flex items-center justify-between p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-amber-300 dark:hover:border-amber-600 transition-colors cursor-pointer">
            <span class="text-sm text-gray-600 dark:text-gray-400">待审核举报</span>
            <span class="text-xl font-bold" :class="summary.todos.pending_reports > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400 dark:text-gray-400'">{{ summary.todos.pending_reports }}</span>
          </NuxtLink>
        </div>
      </div>

      <!-- 趋势图表 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white dark:text-gray-900">访问量趋势</h3>
            <div class="p-2 bg-orange-100 dark:bg-orange-900/40 rounded-full">
              <i class="fas fa-chart-line text-orange-600 dark:text-orange-400 text-sm"></i>
            </div>
          </div>
          <div v-if="hasViewsData" class="h-40">
            <canvas ref="viewsChart"></canvas>
          </div>
          <EmptyState
            v-else
            icon="fas fa-chart-line"
            title="暂无访问数据"
          />
        </div>

        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white dark:text-gray-900">搜索量趋势</h3>
            <div class="p-2 bg-green-100 dark:bg-green-900/40 rounded-full">
              <i class="fas fa-chart-line text-green-600 dark:text-green-400 text-sm"></i>
            </div>
          </div>
          <div v-if="hasSearchesData" class="h-40">
            <canvas ref="searchesChart"></canvas>
          </div>
          <EmptyState
            v-else
            icon="fas fa-chart-line"
            title="暂无搜索数据"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin' })

import { useStatsApi, type StatsSummary } from '~/composables/useApi'
import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'
import { chartDefaults } from '~/composables/useTheme'
import Chart from 'chart.js/auto'

const statsApi = useStatsApi()

// 聚合统计（单次请求 GET /api/stats/summary，含环比与待办）
const { data: summary, pending, error, refresh } = await useAsyncData<StatsSummary>(
  'adminSummary',
  () => statsApi.getSummary()
)

// 趋势数据
const weeklyViews = ref<Array<{ label: string; value: number }>>([])
const weeklySearches = ref<Array<{ label: string; value: number }>>([])

const hasViewsData = computed(() => weeklyViews.value.some((d) => d.value > 0))
const hasSearchesData = computed(() => weeklySearches.value.some((d) => d.value > 0))

const fetchTrendData = async () => {
  try {
    const [viewsRes, searchesRes] = await Promise.all([
      useApiFetch('/stats/views-trend').then(parseApiResponse),
      useApiFetch('/stats/searches-trend').then(parseApiResponse),
    ])
    weeklyViews.value = Array.isArray(viewsRes)
      ? viewsRes.map((item: any) => ({
          label: item.date ? new Date(item.date).toLocaleDateString('zh-CN', { weekday: 'short' }) : '',
          value: Number(item.views) || 0,
        }))
      : []
    weeklySearches.value = Array.isArray(searchesRes)
      ? searchesRes.map((item: any) => ({
          label: item.date ? new Date(item.date).toLocaleDateString('zh-CN', { weekday: 'short' }) : '',
          value: Number(item.searches) || 0,
        }))
      : []
  } catch {
    weeklyViews.value = []
    weeklySearches.value = []
  }
}

// 图表实例
const viewsChart = ref<HTMLCanvasElement | null>(null)
const searchesChart = ref<HTMLCanvasElement | null>(null)
let viewsChartInstance: Chart | null = null
let searchesChartInstance: Chart | null = null

function chartOptions() {
  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: false } },
    scales: {
      y: { beginAtZero: true, grid: { color: chartDefaults.gridColor }, ticks: { color: chartDefaults.textColor } },
      x: { grid: { color: chartDefaults.gridColor }, ticks: { color: chartDefaults.textColor } },
    },
  } as any
}

const initCharts = () => {
  if (hasViewsData.value && viewsChart.value) {
    viewsChartInstance?.destroy()
    const ctx = viewsChart.value.getContext('2d')
    if (ctx) {
      viewsChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
          labels: weeklyViews.value.map((d) => d.label),
          datasets: [{
            label: '访问量',
            data: weeklyViews.value.map((d) => d.value),
            borderColor: '#f97316',
            backgroundColor: 'rgba(249,115,22,0.1)',
            tension: 0.4,
            fill: true,
          }],
        },
        options: chartOptions(),
      })
    }
  }
  if (hasSearchesData.value && searchesChart.value) {
    searchesChartInstance?.destroy()
    const ctx = searchesChart.value.getContext('2d')
    if (ctx) {
      searchesChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
          labels: weeklySearches.value.map((d) => d.label),
          datasets: [{
            label: '搜索量',
            data: weeklySearches.value.map((d) => d.value),
            borderColor: '#22c55e',
            backgroundColor: 'rgba(34,197,94,0.1)',
            tension: 0.4,
            fill: true,
          }],
        },
        options: chartOptions(),
      })
    }
  }
}

watch([weeklyViews, weeklySearches], () => nextTick(initCharts))

onMounted(async () => {
  await fetchTrendData()
  nextTick(initCharts)
})

onBeforeUnmount(() => {
  viewsChartInstance?.destroy()
  searchesChartInstance?.destroy()
})
</script>

<style scoped>
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style>
