<template>
  <Teleport to="body">
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-5xl w-full mx-4 max-h-[90vh] overflow-hidden">
        <!-- 模态框头部 -->
        <div class="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-900">插件日志</h2>
            <p class="text-sm text-gray-600 mt-1">{{ plugin.name }} - 执行日志记录</p>
          </div>
          <button
            @click="$emit('close')"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <Icon name="x" class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <!-- 模态框内容 -->
        <div class="flex flex-col h-[calc(90vh-140px)]">
          <!-- 过滤和控制栏 -->
          <div class="border-b border-gray-200 px-6 py-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <!-- 日志级别过滤 -->
                <select
                  v-model="levelFilter"
                  class="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">全部级别</option>
                  <option value="info">信息</option>
                  <option value="warn">警告</option>
                  <option value="error">错误</option>
                </select>

                <!-- 时间范围过滤 -->
                <select
                  v-model="timeRange"
                  class="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="1h">最近1小时</option>
                  <option value="24h">最近24小时</option>
                  <option value="7d">最近7天</option>
                  <option value="30d">最近30天</option>
                </select>

                <!-- 搜索框 -->
                <div class="relative">
                  <input
                    v-model="searchQuery"
                    type="text"
                    placeholder="搜索日志..."
                    class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  <Icon name="search" class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" />
                </div>
              </div>

              <div class="flex items-center space-x-2">
                <!-- 自动刷新 -->
                <button
                  @click="toggleAutoRefresh"
                  :class="autoRefresh ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'"
                  class="px-3 py-2 rounded-lg hover:bg-opacity-80 transition-colors flex items-center"
                >
                  <Icon name="refresh-cw" class="w-4 h-4 mr-2" :class="{ 'animate-spin': autoRefresh }" />
                  {{ autoRefresh ? '自动刷新' : '手动刷新' }}
                </button>

                <!-- 清空日志 -->
                <button
                  @click="clearLogs"
                  class="px-3 py-2 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition-colors flex items-center"
                >
                  <Icon name="trash" class="w-4 h-4 mr-2" />
                  清空
                </button>

                <!-- 导出日志 -->
                <button
                  @click="exportLogs"
                  class="px-3 py-2 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors flex items-center"
                >
                  <Icon name="download" class="w-4 h-4 mr-2" />
                  导出
                </button>
              </div>
            </div>
          </div>

          <!-- 日志列表 -->
          <div class="flex-1 overflow-hidden">
            <div v-if="loading" class="flex items-center justify-center h-full">
              <div class="text-center">
                <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p class="mt-2 text-gray-600">加载日志...</p>
              </div>
            </div>

            <div v-else-if="filteredLogs.length === 0" class="flex items-center justify-center h-full">
              <div class="text-center">
                <Icon name="file-text" class="w-12 h-12 text-gray-400 mx-auto" />
                <p class="mt-2 text-gray-600">没有找到日志记录</p>
              </div>
            </div>

            <div v-else class="h-full overflow-y-auto">
              <div class="divide-y divide-gray-200">
                <LogEntry
                  v-for="log in paginatedLogs"
                  :key="log.id"
                  :log="log"
                  @view-details="showLogDetails"
                />
              </div>

              <!-- 分页 -->
              <div v-if="totalPages > 1" class="p-4 border-t border-gray-200 bg-gray-50">
                <div class="flex items-center justify-between">
                  <div class="text-sm text-gray-700">
                    显示 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, filteredLogs.length) }} 条，共 {{ filteredLogs.length }} 条日志
                  </div>
                  <div class="flex items-center space-x-2">
                    <button
                      @click="currentPage--"
                      :disabled="currentPage === 1"
                      class="px-3 py-1 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      上一页
                    </button>
                    <span class="px-3 py-1 text-sm text-gray-700">
                      {{ currentPage }} / {{ totalPages }}
                    </span>
                    <button
                      @click="currentPage++"
                      :disabled="currentPage === totalPages"
                      class="px-3 py-1 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      下一页
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 模态框底部 -->
        <div class="border-t border-gray-200 px-6 py-4 flex justify-between items-center">
          <div class="text-sm text-gray-600">
            <span v-if="logs.length > 0">
              共 {{ logs.length }} 条日志记录
            </span>
            <span v-else>
              暂无日志记录
            </span>
          </div>
          <button
            @click="$emit('close')"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- 日志详情模态框 -->
    <LogDetailModal
      v-if="selectedLog"
      :log="selectedLog"
      @close="selectedLog = null"
    />
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useToast } from '~/composables/useToast'
import LogEntry from './LogEntry.vue'
import LogDetailModal from './LogDetailModal.vue'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

const toast = useToast()

// 响应式数据
const loading = ref(false)
const logs = ref([])
const levelFilter = ref('')
const timeRange = ref('24h')
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(50)
const autoRefresh = ref(false)
const refreshInterval = ref(null)
const selectedLog = ref(null)

// 计算属性
const filteredLogs = computed(() => {
  let filtered = logs.value

  // 级别过滤
  if (levelFilter.value) {
    filtered = filtered.filter(log => log.level === levelFilter.value)
  }

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log =>
      log.hook_name.toLowerCase().includes(query) ||
      (log.message && log.message.toLowerCase().includes(query)) ||
      (log.error_message && log.error_message.toLowerCase().includes(query))
    )
  }

  return filtered
})

const totalPages = computed(() => {
  return Math.ceil(filteredLogs.value.length / pageSize.value)
})

const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredLogs.value.slice(start, end)
})

// 方法
const loadLogs = async () => {
  try {
    loading.value = true
    const response = await $fetch(`/api/plugins/${props.plugin.name}/logs?page=${currentPage.value}&limit=${pageSize.value}`)
    if (response.success) {
      logs.value = response.data.logs || []
    }
  } catch (error) {
    toast.error('加载日志失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value

  if (autoRefresh.value) {
    refreshInterval.value = setInterval(() => {
      loadLogs()
    }, 5000) // 每5秒刷新一次
    toast.info('已启用自动刷新')
  } else {
    if (refreshInterval.value) {
      clearInterval(refreshInterval.value)
      refreshInterval.value = null
    }
    toast.info('已关闭自动刷新')
  }
}

const clearLogs = async () => {
  if (!confirm('确定要清空所有日志吗？此操作不可恢复。')) {
    return
  }

  try {
    // 这里应该调用清空日志的API
    toast.success('日志清空功能开发中...')
  } catch (error) {
    toast.error('清空日志失败: ' + error.message)
  }
}

const exportLogs = () => {
  if (filteredLogs.value.length === 0) {
    toast.warning('没有日志可导出')
    return
  }

  try {
    const logText = filteredLogs.value.map(log => {
      let message = ''
      if (log.message) {
        message = ` - ${log.message}`
      }
      if (log.error_message) {
        message += (message ? ' | ' : ' - ') + `ERROR: ${log.error_message}`
      }
      return `[${log.created_at}] ${log.hook_name}: ${log.success ? 'SUCCESS' : 'ERROR'} (${log.execution_time}ms)${message}`
    }).join('\n')

    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${props.plugin.name}_logs_${new Date().toISOString().slice(0, 10)}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    toast.success('日志已导出')
  } catch (error) {
    toast.error('导出日志失败: ' + error.message)
  }
}

const showLogDetails = (log) => {
  selectedLog.value = log
}

// 生命周期
onMounted(() => {
  loadLogs()
})

onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})

// 监听过滤器变化
watch([levelFilter, timeRange, searchQuery], () => {
  currentPage.value = 1
})

// 监听分页变化
watch(currentPage, () => {
  loadLogs()
})
</script>

<style scoped>
/* 模态框动画 */
.fixed > div {
  animation: modalFadeIn 0.3s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 自动刷新动画 */
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>