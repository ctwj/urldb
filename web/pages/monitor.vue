<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col">
    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
        <!-- 头部 -->
        <div class="header-container bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
          <h1 class="text-2xl sm:text-3xl font-bold mb-4">
            <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
              系统性能监控
            </a>
          </h1>
          <p class="text-gray-300 max-w-2xl mx-auto">实时监控系统运行状态和性能指标</p>
          <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
            <NuxtLink to="/" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-home text-xs"></i> 首页
              </n-button>
            </NuxtLink>
            <NuxtLink to="/hot-dramas" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-film text-xs"></i> 热播剧
              </n-button>
            </NuxtLink>
            <NuxtLink to="/api-docs" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-book text-xs"></i> API文档
              </n-button>
            </NuxtLink>
          </nav>
        </div>

        <!-- 刷新按钮 -->
        <div class="mb-6 flex justify-between items-center">
          <div class="flex items-center space-x-4">
            <button
              @click="refreshData"
              :disabled="loading"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
            >
              <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
              <span>{{ loading ? '刷新中...' : '刷新数据' }}</span>
            </button>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              最后更新: {{ lastUpdateTime }}
            </div>
          </div>
          <div class="flex items-center space-x-2">
            <label class="text-sm text-gray-600 dark:text-gray-400">自动刷新:</label>
            <n-checkbox
              v-model:checked="autoRefresh"
            />
            <span class="text-sm text-gray-500 dark:text-gray-400">{{ autoRefreshInterval }}秒</span>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>

        <!-- 监控数据 -->
        <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
          <!-- 系统信息卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-server mr-2 text-blue-600"></i>
              系统信息
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">运行时间:</span>
                <span class="font-medium">{{ systemInfo.uptime || 'N/A' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">启动时间:</span>
                <span class="font-medium">{{ systemInfo.start_time || 'N/A' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">版本:</span>
                <span class="font-medium">{{ systemInfo.version || 'N/A' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">运行模式:</span>
                <span class="font-medium">{{ systemInfo.environment?.gin_mode || 'N/A' }}</span>
              </div>
            </div>
          </div>

          <!-- 内存使用卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-memory mr-2 text-green-600"></i>
              内存使用
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">当前分配:</span>
                <span class="font-medium">{{ formatBytes(performanceStats.memory?.alloc) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">总分配:</span>
                <span class="font-medium">{{ formatBytes(performanceStats.memory?.total_alloc) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">系统内存:</span>
                <span class="font-medium">{{ formatBytes(performanceStats.memory?.sys) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">堆内存:</span>
                <span class="font-medium">{{ formatBytes(performanceStats.memory?.heap_alloc) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">GC次数:</span>
                <span class="font-medium">{{ performanceStats.memory?.num_gc || 0 }}</span>
              </div>
            </div>
          </div>

          <!-- 数据库连接卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-database mr-2 text-purple-600"></i>
              数据库连接
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">最大连接数:</span>
                <span class="font-medium">{{ performanceStats.database?.max_open_connections || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">当前连接数:</span>
                <span class="font-medium">{{ performanceStats.database?.open_connections || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">使用中:</span>
                <span class="font-medium">{{ performanceStats.database?.in_use || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">空闲:</span>
                <span class="font-medium">{{ performanceStats.database?.idle || 0 }}</span>
              </div>
            </div>
          </div>

          <!-- 系统资源卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-microchip mr-2 text-orange-600"></i>
              系统资源
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">CPU核心数:</span>
                <span class="font-medium">{{ performanceStats.system?.cpu_count || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Go版本:</span>
                <span class="font-medium">{{ performanceStats.system?.go_version || 'N/A' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">协程数:</span>
                <span class="font-medium">{{ performanceStats.goroutines || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">时间戳:</span>
                <span class="font-medium">{{ formatTimestamp(performanceStats.timestamp) }}</span>
              </div>
            </div>
          </div>

          <!-- 基础统计卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-chart-bar mr-2 text-red-600"></i>
              基础统计
            </h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">资源总数:</span>
                <span class="font-medium">{{ basicStats.total_resources || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">分类总数:</span>
                <span class="font-medium">{{ basicStats.total_categories || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">标签总数:</span>
                <span class="font-medium">{{ basicStats.total_tags || 0 }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">总浏览量:</span>
                <span class="font-medium">{{ basicStats.total_views || 0 }}</span>
              </div>
            </div>
          </div>

          <!-- 性能图表卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center">
              <i class="fas fa-chart-line mr-2 text-indigo-600"></i>
              性能趋势
            </h3>
            <div class="space-y-3">
              <div class="text-center py-4">
                <div class="text-2xl font-bold text-indigo-600">
                  {{ formatBytes(performanceStats.memory?.alloc) }}
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-400">当前内存使用</div>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div 
                  class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
                  :style="{ width: memoryUsagePercentage + '%' }"
                ></div>
              </div>
              <div class="text-center text-sm text-gray-500 dark:text-gray-400">
                内存使用率: {{ memoryUsagePercentage.toFixed(1) }}%
              </div>
            </div>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="mt-6 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          <div class="flex items-center">
            <i class="fas fa-exclamation-triangle mr-2"></i>
            <span>{{ error }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'default'
})

import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useMonitorApi } from '~/composables/useApi'
const monitorApi = useMonitorApi()

// 响应式数据
const loading = ref(false)
const error = ref('')
const lastUpdateTime = ref('')
const autoRefresh = ref(false)
const autoRefreshInterval = ref(30)

// 监控数据
const systemInfo = ref<any>({})
const performanceStats = ref<any>({})
const basicStats = ref<any>({})

// 计算内存使用率
const memoryUsagePercentage = computed(() => {
  const memory = performanceStats.value.memory
  if (!memory || !memory.sys) return 0
  return (memory.alloc / memory.sys) * 100
})

// 格式化字节数
const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化时间戳
const formatTimestamp = (timestamp: number) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const response = await monitorApi.getSystemInfo()
    systemInfo.value = response
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
}

// 获取性能统计
const fetchPerformanceStats = async () => {
  try {
    const response = await monitorApi.getPerformanceStats()
    performanceStats.value = response
    console.log('性能统计数据:', response)
    console.log('数据库连接信息:', (response as any).database)
  } catch (error) {
    console.error('获取性能统计失败:', error)
  }
}

// 获取基础统计
const fetchBasicStats = async () => {
  try {
    const response = await monitorApi.getBasicStats()
    basicStats.value = response
  } catch (error) {
    console.error('获取基础统计失败:', error)
  }
}

// 刷新所有数据
const refreshData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    await Promise.all([
      fetchSystemInfo(),
      fetchPerformanceStats(),
      fetchBasicStats()
    ])
    lastUpdateTime.value = new Date().toLocaleString('zh-CN')
  } catch (err: any) {
    error.value = err.message || '获取监控数据失败'
  } finally {
    loading.value = false
  }
}

// 自动刷新定时器
let autoRefreshTimer: NodeJS.Timeout | null = null

// 监听自动刷新设置
const startAutoRefresh = () => {
  if (autoRefresh.value) {
    autoRefreshTimer = setInterval(refreshData, autoRefreshInterval.value * 1000)
  }
}

const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

// 监听自动刷新变化
watch(autoRefresh, (newValue) => {
  if (newValue) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
})

// 页面加载时获取数据
onMounted(() => {
  refreshData()
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

// 页面卸载时清理定时器
onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
/* 可以添加自定义样式 */
.header-container{
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}
</style> 