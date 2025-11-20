<template>
  <div v-if="showCacheInfo && isClient" class="fixed bottom-4 right-4 z-50 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 p-4 max-w-sm">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">系统配置缓存状态</h3>
      <button
        @click="showCacheInfo = false"
        class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
      >
        <i class="fas fa-times"></i>
      </button>
    </div>

    <div class="space-y-2 text-xs">
      <!-- 初始化状态 -->
      <div class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">初始化状态:</span>
        <span :class="status.initialized ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'">
          {{ status.initialized ? '已初始化' : '未初始化' }}
        </span>
      </div>

      <!-- 加载状态 -->
      <div class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">加载状态:</span>
        <span :class="status.isLoading ? 'text-blue-600 dark:text-blue-400' : 'text-gray-600 dark:text-gray-400'">
          {{ status.isLoading ? '加载中...' : '空闲' }}
        </span>
      </div>

      <!-- 缓存状态 -->
      <div class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">缓存状态:</span>
        <span :class="status.isCacheValid ? 'text-green-600 dark:text-green-400' : 'text-orange-600 dark:text-orange-400'">
          {{ status.isCacheValid ? '有效' : '无效/过期' }}
        </span>
      </div>

      <!-- 缓存剩余时间 -->
      <div class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">缓存剩余:</span>
        <span class="text-gray-900 dark:text-gray-100">
          {{ formatTime(status.cacheTimeRemaining) }}
        </span>
      </div>

      <!-- 最后获取时间 -->
      <div class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">最后更新:</span>
        <span class="text-gray-900 dark:text-gray-100">
          {{ formatLastFetch(status.lastFetchTime) }}
        </span>
      </div>

      <!-- 错误信息 -->
      <div v-if="status.error" class="flex justify-between">
        <span class="text-gray-600 dark:text-gray-400">错误:</span>
        <span class="text-red-600 dark:text-red-400">
          {{ status.error }}
        </span>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="mt-4 flex gap-2">
      <button
        @click="refreshCache"
        :disabled="status.isLoading"
        class="flex-1 px-2 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <i class="fas fa-sync-alt mr-1"></i>
        刷新
      </button>
      <button
        @click="clearCache"
        class="flex-1 px-2 py-1 text-xs bg-red-500 text-white rounded hover:bg-red-600"
      >
        <i class="fas fa-trash mr-1"></i>
        清除
      </button>
    </div>
  </div>

  <!-- 浮动按钮（仅在开发环境和客户端显示） -->
  <button
    v-if="isDev && isClient"
    @click="showCacheInfo = !showCacheInfo"
    class="fixed bottom-4 right-4 z-40 w-12 h-12 bg-purple-500 text-white rounded-full shadow-lg hover:bg-purple-600 transition-colors flex items-center justify-center"
    title="系统配置缓存信息"
  >
    <i class="fas fa-database"></i>
  </button>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useSystemConfigStore } from '~/stores/systemConfig'

const systemConfigStore = useSystemConfigStore()
const showCacheInfo = ref(false)

// 检查是否为开发环境
const isDev = computed(() => {
  return process.env.NODE_ENV === 'development'
})

// 检查是否为客户端
const isClient = computed(() => {
  return process.client
})

// 获取状态信息 - 直接访问store的响应式状态以确保正确更新
const status = computed(() => ({
  initialized: systemConfigStore.initialized,
  isLoading: systemConfigStore.isLoading,
  error: systemConfigStore.error,
  lastFetchTime: systemConfigStore.lastFetchTime,
  cacheTimeRemaining: systemConfigStore.cacheTimeRemaining,
  isCacheValid: systemConfigStore.isCacheValid
}))

// 格式化时间显示
const formatTime = (seconds: number): string => {
  if (seconds <= 0) return '已过期'

  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60

  if (minutes > 0) {
    return `${minutes}分${remainingSeconds}秒`
  } else {
    return `${remainingSeconds}秒`
  }
}

// 格式化最后获取时间
const formatLastFetch = (timestamp: number): string => {
  if (!timestamp) return '从未'

  const now = Date.now()
  const diff = now - timestamp

  if (diff < 60 * 1000) {
    return '刚刚'
  } else if (diff < 60 * 60 * 1000) {
    const minutes = Math.floor(diff / (60 * 1000))
    return `${minutes}分钟前`
  } else if (diff < 24 * 60 * 60 * 1000) {
    const hours = Math.floor(diff / (60 * 60 * 1000))
    return `${hours}小时前`
  } else {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

// 刷新缓存
const refreshCache = async () => {
  try {
    await systemConfigStore.refreshConfig()
    console.log('[CacheInfo] 手动刷新缓存完成')
  } catch (error) {
    console.error('[CacheInfo] 刷新缓存失败:', error)
  }
}

// 清除缓存
const clearCache = () => {
  systemConfigStore.clearCache()
  console.log('[CacheInfo] 手动清除缓存完成')
}

// 键盘快捷键支持
const handleKeydown = (e: KeyboardEvent) => {
  // Ctrl+Shift+C 显示/隐藏缓存信息（仅在开发环境）
  if (isDev.value && e.ctrlKey && e.shiftKey && e.key === 'C') {
    e.preventDefault()
    showCacheInfo.value = !showCacheInfo.value
  }
}

onMounted(() => {
  if (isClient.value) {
    document.addEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  if (isClient.value) {
    document.removeEventListener('keydown', handleKeydown)
  }
})
</script>

<style scoped>
/* 添加一些动画效果 */
.transition-colors {
  transition: color 0.2s, background-color 0.2s;
}
</style>