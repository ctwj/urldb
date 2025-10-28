<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">公开API访问日志</h1>
        <p class="text-gray-600 dark:text-gray-400">查看公开API的访问记录和统计信息</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="refreshData" :loading="loading">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
        <n-button type="warning" @click="clearOldLogs" :loading="clearing">
          <template #icon>
            <i class="fas fa-trash-alt"></i>
          </template>
          清理旧日志
        </n-button>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和筛选 -->
    <template #filter-bar>
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索接口路径或IP..."
            @keyup.enter="handleSearch"
            clearable
          >
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </n-input>

          <n-date-picker
            v-model:value="startDate"
            type="date"
            placeholder="开始日期"
            clearable
          />

          <n-date-picker
            v-model:value="endDate"
            type="date"
            placeholder="结束日期"
            clearable
          />

          <n-button type="primary" @click="handleSearch" class="w-20">
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
          <span class="text-lg font-semibold">访问日志列表</span>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            共 {{ total }} 条日志
          </div>
        </div>
        <!-- 统计卡片 -->
        <div class="flex space-x-6">
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-blue-600">{{ summary.total_requests }}</div>
            <div class="text-xs text-gray-500">总请求</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-green-600">{{ summary.today_requests }}</div>
            <div class="text-xs text-gray-500">今日请求</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-purple-600">{{ summary.week_requests }}</div>
            <div class="text-xs text-gray-500">本周请求</div>
          </div>
          <div class="text-center flex items-base">
            <div class="text-2xl font-bold text-red-600">{{ summary.error_requests }}</div>
            <div class="text-xs text-gray-500">错误请求</div>
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
      <div v-else-if="logs.length === 0" class="flex flex-col items-center justify-center py-12">
        <i class="fas fa-file-alt text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500 dark:text-gray-400">暂无访问日志</p>
      </div>

      <!-- 日志列表 -->
      <div v-else class="space-y-2 h-full overflow-y-auto">
        <div
          v-for="log in logs"
          :key="log.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-3 mb-2 w-full">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ log.id }}</span>
                <n-tag :type="getMethodTagType(log.method)" size="small">
                  {{ log.method }}
                </n-tag>
                <n-tag :type="getStatusTagType(log.response_status)" size="small">
                  {{ log.response_status }}
                </n-tag>
                <code class="text-sm bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                    {{ log.endpoint }}
                  </code>
                <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 float-right">
                  <code class="text-sm bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">
                    {{ log.ip }}
                  </code>
                  <span class="flex items-center">
                    <i class="fas fa-clock mr-1"></i>
                    {{ formatDate(log.created_at) }}
                  </span>
                  <span v-if="log.processing_time > 0"  class="flex items-center">
                    <i class="fas fa-tachometer-alt mr-1"></i>
                    {{ log.processing_time }}ms
                  </span>
                </div>
              </div>

              <div v-if="log.request_params" class="mt-2 text-xs text-gray-600 dark:text-gray-400 flex items-center justify-between">
                <div class="flex items-center flex-1 min-w-0">
                  <strong class="mr-2 flex-0 whitespace-nowrap">请求参数:</strong>
                  <span class="truncate">{{ log.request_params }}</span>
                </div>
                <div class="flex items-center space-x-1 ml-2">
                  <n-button size="tiny" @click="copyParams(log.request_params)">
                    <template #icon>
                      <i class="fas fa-copy"></i>
                    </template>
                  </n-button>
                  <n-button size="tiny" @click="viewParams(log.request_params)">
                    <template #icon>
                      <i class="fas fa-eye"></i>
                    </template>
                  </n-button>
                </div>
              </div>

              <div v-if="log.error_message" class="mt-2 text-xs text-red-600 dark:text-red-400">
                <strong>错误信息:</strong> {{ log.error_message }}
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
            v-model:page="currentPage"
            v-model:page-size="pageSize"
            :item-count="total"
            :page-sizes="[10, 20, 50, 100]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </template>
  </AdminPageLayout>

    <!-- 请求参数详情模态框 -->
  <n-modal v-model:show="showModal" preset="card" title="请求参数详情" style="min-width: 600px;">
    <n-code
      :code="selectedParams"
      language="json"
      :folding="true"
      :show-line-numbers="true"
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded max-h-96 overflow-auto"
    />
  </n-modal>

  <!-- 请求参数详情模态框 -->
  <n-modal v-model:show="showModal" preset="card" title="请求参数详情" style="min-width: 600px;">
    <n-code
      :code="selectedParams"
      language="json"
      :folding="true"
      :show-line-numbers="true"
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded max-h-96 overflow-auto"
    />
  </n-modal>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

import { useApiAccessLogApi } from '~/composables/useApi'

const notification = useNotification()
const dialog = useDialog()

interface ApiAccessLog {
  id: number
  ip: string
  user_agent: string
  endpoint: string
  method: string
  request_params: string
  response_status: number
  response_data: string
  process_count: number
  error_message: string
  processing_time: number
  created_at: string
}

// 获取API实例
const apiAccessLogApi = useApiAccessLogApi()

// 响应式数据
const loading = ref(false)
const clearing = ref(false)
const logs = ref<ApiAccessLog[]>([])
const summary = ref({
  total_requests: 0,
  today_requests: 0,
  week_requests: 0,
  month_requests: 0,
  error_requests: 0,
  unique_ips: 0
})

// 筛选和搜索
const searchQuery = ref('')
const startDate = ref<number | null>(null)
const endDate = ref<number | null>(null)

// 分页
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 模态框
const showModal = ref(false)
const selectedParams = ref('')


// 获取日志数据
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    // 添加日期筛选
    if (startDate.value) {
      const date = new Date(startDate.value)
      params.start_date = date.toISOString().split('T')[0]
    }
    if (endDate.value) {
      const date = new Date(endDate.value)
      params.end_date = date.toISOString().split('T')[0]
    }

    // 添加搜索条件
    if (searchQuery.value) {
      params.endpoint = searchQuery.value
      params.ip = searchQuery.value
    }

    const response = await apiAccessLogApi.getApiAccessLogs(params) as any
    logs.value = response.data || []
    total.value = response.total || 0

  } catch (error) {
    console.error('获取API访问日志失败:', error)
    notification.error({
      content: '获取API访问日志失败',
      duration: 3000
    })
    logs.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 获取统计汇总
const fetchSummary = async () => {
  try {
    const response = await apiAccessLogApi.getApiAccessLogSummary()
    summary.value = response as any
  } catch (error) {
    console.error('获取统计汇总失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
  fetchSummary()
}

// 方法标签类型
const getMethodTagType = (method: string): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' => {
  const methodColors: Record<string, 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'> = {
    GET: 'info',
    POST: 'success',
    PUT: 'warning',
    DELETE: 'error',
    PATCH: 'warning'
  }
  return methodColors[method] || 'default'
}

// 状态标签类型
const getStatusTagType = (status: number) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 400 && status < 500) return 'warning'
  if (status >= 500) return 'error'
  return 'default'
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 复制参数
const copyParams = (params: string) => {
  navigator.clipboard.writeText(params).then(() => {
    notification.success({
      content: '已复制到剪贴板',
      duration: 2000
    })
  }).catch(() => {
    notification.error({
      content: '复制失败',
      duration: 2000
    })
  })
}

// 查看参数
const viewParams = (params: string) => {
  try {
    selectedParams.value = JSON.stringify(JSON.parse(params), null, 2)
  } catch {
    selectedParams.value = params
  }
  showModal.value = true
}


// 清理旧日志
const clearOldLogs = async () => {
  dialog.warning({
    title: '清理旧日志',
    content: '确定要清理30天前的旧日志吗？此操作不可恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        clearing.value = true
        await apiAccessLogApi.clearApiAccessLogs(30)

        notification.success({
          content: '旧日志清理成功',
          duration: 3000
        })

        refreshData()
      } catch (error) {
        console.error('清理旧日志失败:', error)
        notification.error({
          content: '清理旧日志失败',
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
  await Promise.all([fetchData(), fetchSummary()])
})
</script>

<style scoped>
.logs-content {
  padding: 1rem;
  background-color: var(--color-white, #ffffff);
}

.dark .logs-content {
  background-color: var(--color-dark-bg, #1f2937);
}
</style>