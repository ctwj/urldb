<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">系统日志</h1>
        <p class="text-gray-600 dark:text-gray-400">查看系统运行日志和错误信息</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="refreshSystemLogs" :loading="loading">
          <template #icon>
            <i class="fas fa-sync-alt"></i>
          </template>
          刷新
        </n-button>
        <n-button type="warning" @click="clearSystemLogs" :loading="clearing">
          <template #icon>
            <i class="fas fa-trash-alt"></i>
          </template>
          清理日志
        </n-button>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和筛选 -->
    <template #filter-bar>
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
          <n-select
            v-model:value="systemLogLevel"
            :options="logLevelOptions"
            placeholder="选择日志级别"
            clearable
          />
          <n-input
            v-model:value="systemLogSearch"
            placeholder="搜索日志内容..."
            @keyup.enter="handleSystemLogSearch"
            clearable
          >
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </n-input>
          <n-date-picker
            v-model:value="systemStartDate"
            type="date"
            placeholder="开始日期"
            clearable
          />
          <n-date-picker
            v-model:value="systemEndDate"
            type="date"
            placeholder="结束日期"
            clearable
          />
          <n-button type="primary" @click="handleSystemLogSearch">
            <template #icon>
              <i class="fas fa-search"></i>
            </template>
            搜索
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区header - 统计信息 -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-lg font-semibold">系统日志列表</span>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            共 {{ systemTotal }} 条日志
          </div>
        </div>
        <!-- 统计卡片 -->
        <div class="flex space-x-6" v-if="systemLogSummary">
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-blue-600">{{ systemLogSummary.total }}</div>
            <div class="text-xs text-gray-500">总日志</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-gray-500">{{ systemLogSummary.debug }}</div>
            <div class="text-xs text-gray-500">调试</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-green-600">{{ systemLogSummary.info }}</div>
            <div class="text-xs text-gray-500">信息</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-yellow-600">{{ systemLogSummary.warn }}</div>
            <div class="text-xs text-gray-500">警告</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-red-600">{{ systemLogSummary.error }}</div>
            <div class="text-xs text-gray-500">错误</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-purple-600">{{ systemLogSummary.fatal }}</div>
            <div class="text-xs text-gray-500">致命</div>
          </div>
        </div>
      </div>
    </template>

    <!-- 内容区content - 日志列表 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="systemLogs.length === 0" class="flex flex-col items-center justify-center py-12">
        <i class="fas fa-file-alt text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500 dark:text-gray-400">暂无系统日志</p>
      </div>

      <!-- 日志列表 -->
      <div v-else class="space-y-2 h-full overflow-y-auto">
        <div
          v-for="(log, index) in systemLogs"
          :key="index"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          :class="getLogItemClass(log.level)"
        >
          <div class="flex items-start">
            <div class="flex-shrink-0 w-3 h-3 rounded-full mt-1.5 mr-3" :class="getLogLevelColor(log.level)"></div>
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <n-tag :type="getLogLevelTagType(log.level)" size="small">
                  {{ log.level }}
                </n-tag>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ formatLogTime(log.timestamp) }}
                </span>
                <span v-if="log.file" class="text-xs text-gray-500 dark:text-gray-400">
                  {{ log.file }}:{{ log.line }}
                </span>
              </div>
              <div class="text-sm text-gray-800 dark:text-gray-200 font-mono break-words">
                {{ log.message }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 内容区footer - 分页组件 -->
    <template #content-footer>
      <div class="p-4">
        <div class="flex justify-center">
          <n-pagination
            v-model:page="systemCurrentPage"
            v-model:page-size="systemPageSize"
            :item-count="systemTotal"
            :page-sizes="[20, 50, 100]"
            show-size-picker
            @update:page="handleSystemLogPageChange"
            @update:page-size="handleSystemLogPageSizeChange"
          />
        </div>
      </div>
    </template>
  </AdminPageLayout>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

import { useSystemLogApi } from '~/composables/useApi'

const notification = useNotification()
const dialog = useDialog()

// 获取API实例
const systemLogApi = useSystemLogApi()

// 响应式数据
const loading = ref(false)
const clearing = ref(false)
const systemLogs = ref<any[]>([])
const systemLogSummary = ref<any>(null)

// 过滤和搜索
const systemLogLevel = ref<string | null>(null)
const systemLogSearch = ref('')
const systemStartDate = ref<number | null>(null)
const systemEndDate = ref<number | null>(null)

// 分页
const systemCurrentPage = ref(1)
const systemPageSize = ref(50)
const systemTotal = ref(0)

// 日志级别选项
const logLevelOptions = [
  { label: 'DEBUG', value: 'debug' },
  { label: 'INFO', value: 'info' },
  { label: 'WARN', value: 'warn' },
  { label: 'ERROR', value: 'error' },
  { label: 'FATAL', value: 'fatal' }
]

// 获取系统日志数据
const fetchSystemLogs = async () => {
  loading.value = true
  try {
    const params: any = {
      page: systemCurrentPage.value,
      page_size: systemPageSize.value
    }

    // 添加级别筛选
    if (systemLogLevel.value) {
      params.level = systemLogLevel.value
    }

    // 添加日期筛选
    if (systemStartDate.value) {
      const date = new Date(systemStartDate.value)
      params.start_date = date.toISOString().split('T')[0]
    }
    if (systemEndDate.value) {
      const date = new Date(systemEndDate.value)
      params.end_date = date.toISOString().split('T')[0]
    }

    // 添加搜索条件
    if (systemLogSearch.value) {
      params.search = systemLogSearch.value
    }

    const response = await systemLogApi.getSystemLogs(params) as any
    systemLogs.value = response.data || []
    systemTotal.value = response.total || 0
  } catch (error) {
    console.error('获取系统日志失败:', error)
    notification.error({
      content: '获取系统日志失败',
      duration: 3000
    })
    systemLogs.value = []
    systemTotal.value = 0
  } finally {
    loading.value = false
  }
}

// 获取系统日志统计
const fetchSystemLogSummary = async () => {
  try {
    const response = await systemLogApi.getSystemLogSummary()
    systemLogSummary.value = response.summary || null
  } catch (error) {
    console.error('获取系统日志统计失败:', error)
  }
}

// 刷新系统日志
const refreshSystemLogs = () => {
  fetchSystemLogs()
  fetchSystemLogSummary()
}

// 系统日志搜索处理
const handleSystemLogSearch = () => {
  systemCurrentPage.value = 1
  fetchSystemLogs()
}

// 系统日志分页处理
const handleSystemLogPageChange = (page: number) => {
  systemCurrentPage.value = page
  fetchSystemLogs()
}

const handleSystemLogPageSizeChange = (size: number) => {
  systemPageSize.value = size
  systemCurrentPage.value = 1
  fetchSystemLogs()
}

// 获取日志级别标签类型
const getLogLevelTagType = (level: string): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' => {
  const levelMap: Record<string, 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'> = {
    'DEBUG': 'default',
    'INFO': 'info',
    'WARN': 'warning',
    'ERROR': 'error',
    'FATAL': 'error'
  }
  return levelMap[level?.toUpperCase()] || 'default'
}

// 获取日志级别颜色
const getLogLevelColor = (level: string): string => {
  const colorMap: Record<string, string> = {
    'DEBUG': 'bg-gray-400',
    'INFO': 'bg-blue-500',
    'WARN': 'bg-yellow-500',
    'ERROR': 'bg-red-500',
    'FATAL': 'bg-purple-500'
  }
  return colorMap[level?.toUpperCase()] || 'bg-gray-400'
}

// 获取日志项类名
const getLogItemClass = (level: string): string => {
  const classMap: Record<string, string> = {
    'DEBUG': 'bg-gray-50 dark:bg-gray-800',
    'INFO': 'bg-blue-50 dark:bg-blue-900/20',
    'WARN': 'bg-yellow-50 dark:bg-yellow-900/20',
    'ERROR': 'bg-red-50 dark:bg-red-900/20',
    'FATAL': 'bg-purple-50 dark:bg-purple-900/20'
  }
  return classMap[level?.toUpperCase()] || ''
}

// 格式化日志时间
const formatLogTime = (timestamp: string) => {
  if (!timestamp) return '-'
  try {
    const date = new Date(timestamp)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    console.error('时间格式化错误:', error)
    return timestamp
  }
}


// 清理系统日志
const clearSystemLogs = async () => {
  dialog.warning({
    title: '清理系统日志',
    content: '确定要清理30天前的系统日志吗？此操作不可恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        clearing.value = true
        await systemLogApi.clearSystemLogs(30)

        notification.success({
          content: '系统日志清理成功',
          duration: 3000
        })

        refreshSystemLogs()
      } catch (error) {
        console.error('清理系统日志失败:', error)
        notification.error({
          content: '清理系统日志失败',
          duration: 3000
        })
      } finally {
        clearing.value = false
      }
    }
  })
}

// 页面加载时获取数据
onMounted(async () => {
  await Promise.all([fetchSystemLogs(), fetchSystemLogSummary()])
})
</script>

<style scoped>
/* 日志条目悬停效果 */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
}

.dark .hover\:bg-gray-800:hover {
  background-color: #1f2937;
}
</style>