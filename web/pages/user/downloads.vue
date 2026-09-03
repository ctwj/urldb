<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">下载历史</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">记录您获取过的资源链接</p>
      </div>
      <n-button type="primary" :disabled="historyRecords.length === 0" @click="handleClearHistory">
        <template #icon>
          <i class="fas fa-trash"></i>
        </template>
        清空历史
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6">
      <n-card :bordered="false" class="shadow-sm hover:shadow-md transition-shadow duration-200 cursor-default">
        <div class="flex items-center gap-4">
          <div class="w-11 h-11 rounded-xl bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center flex-shrink-0">
            <i class="fas fa-download text-blue-500 dark:text-blue-400 text-lg"></i>
          </div>
          <div class="min-w-0">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">下载资源总数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white leading-tight mt-0.5">{{ stats.total }}</p>
          </div>
        </div>
      </n-card>

      <n-card :bordered="false" class="shadow-sm hover:shadow-md transition-shadow duration-200 cursor-default">
        <div class="flex items-center gap-4">
          <div class="w-11 h-11 rounded-xl bg-green-50 dark:bg-green-900/30 flex items-center justify-center flex-shrink-0">
            <i class="fas fa-calendar-day text-green-500 dark:text-green-400 text-lg"></i>
          </div>
          <div class="min-w-0">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">今日下载</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white leading-tight mt-0.5">{{ stats.today }}</p>
          </div>
        </div>
      </n-card>

      <n-card :bordered="false" class="shadow-sm hover:shadow-md transition-shadow duration-200 cursor-default">
        <div class="flex items-center gap-4">
          <div class="w-11 h-11 rounded-xl bg-amber-50 dark:bg-amber-900/30 flex items-center justify-center flex-shrink-0">
            <i class="fas fa-calendar-week text-amber-500 dark:text-amber-400 text-lg"></i>
          </div>
          <div class="min-w-0">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">最近7天</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white leading-tight mt-0.5">{{ stats.thisWeek }}</p>
          </div>
        </div>
      </n-card>

      <n-card :bordered="false" class="shadow-sm hover:shadow-md transition-shadow duration-200 cursor-default">
        <div class="flex items-center gap-4">
          <div class="w-11 h-11 rounded-xl bg-purple-50 dark:bg-purple-900/30 flex items-center justify-center flex-shrink-0">
            <i class="fas fa-calendar-days text-purple-500 dark:text-purple-400 text-lg"></i>
          </div>
          <div class="min-w-0">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">最近30天</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white leading-tight mt-0.5">{{ stats.thisMonth }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 历史记录列表 -->
    <n-card :bordered="false" class="shadow-sm">
      <n-data-table
        :columns="columns"
        :data="historyRecords"
        :loading="loading"
        :pagination="pagination"
        :bordered="false"
      >
        <template #empty>
          <n-empty description="还没有下载过资源">
            <template #icon>
              <i class="fas fa-download text-gray-300 dark:text-gray-600 text-4xl"></i>
            </template>
            <template #extra>
              <n-button type="primary" size="small" @click="navigateTo('/')">
                去发现资源
              </n-button>
            </template>
          </n-empty>
        </template>
      </n-data-table>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { parseApiResponse } from '~/composables/useApi'
import { useApiFetch } from '~/composables/useApiFetch'

// 在 setup 同步上下文获取 naive-ui 实例（async 处理器内调用 useNotification 会因脱离 setup 上下文而失败）
const notification = useNotification()
const dialog = useDialog()

// 页面元数据
definePageMeta({
  layout: 'user',
  title: '下载历史'
})

// 统计数据
const stats = ref({
  total: 0,
  today: 0,
  thisWeek: 0,
  thisMonth: 0
})

const loading = ref(false)
const historyRecords = ref<any[]>([])

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  itemCount: 0,
  onChange: (page: number) => {
    pagination.value.page = page
    fetchHistory()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchHistory()
  }
})

// 表格列配置
const columns = [
  {
    title: '资源名称',
    key: 'title',
    render: (row: any) => {
      return h('div', [
        h('div', { class: 'font-medium' }, row.title || '（资源已不存在）'),
        h('div', { class: 'text-sm text-gray-500 line-clamp-1' }, row.description || '')
      ])
    }
  },
  {
    title: '平台',
    key: 'platform',
    width: 120,
    render: (row: any) => {
      if (!row.platform) return h('span', { class: 'text-sm text-gray-400' }, '未知')
      return h('n-tag', { type: 'info', size: 'small' }, { default: () => row.platform })
    }
  },
  {
    title: '下载次数',
    key: 'download_count',
    width: 100,
    render: (row: any) => {
      return h('span', { class: 'font-medium' }, String(row.download_count ?? 1))
    }
  },
  {
    title: '首次下载',
    key: 'first_download',
    width: 170,
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDateTime(row.first_download))
    }
  },
  {
    title: '最近下载',
    key: 'last_download',
    width: 170,
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDateTime(row.last_download))
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: any) => {
      return h('n-space', { size: 'small' }, {
        default: () => [
          h('n-button', {
            size: 'small',
            type: 'primary',
            disabled: !row.resource_key,
            onClick: () => handleView(row)
          }, { default: () => '查看' }),
          h('n-button', {
            size: 'small',
            type: 'warning',
            onClick: () => handleRemoveRecord(row)
          }, { default: () => '删除' })
        ]
      })
    }
  }
]

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取下载历史
const fetchHistory = async () => {
  loading.value = true
  try {
    // 注意：parseApiResponse 对 data.list 形状会退化为数组，这里直接解包 data
    const res = await useApiFetch('/user/download-history', {
      params: {
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      }
    }) as any
    const payload = res?.data || {}

    historyRecords.value = payload.list || []
    pagination.value.itemCount = payload.total || 0
  } catch (error: any) {
    console.error('获取下载历史失败:', error)
    if (process.client) {
      notification.error({
        content: error?.message || '获取下载历史失败',
        duration: 3000
      })
    }
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await useApiFetch('/user/download-history/stats').then(parseApiResponse) as any
    if (response) {
      stats.value = {
        total: response.total || 0,
        today: response.today || 0,
        thisWeek: response.thisWeek || 0,
        thisMonth: response.thisMonth || 0
      }
    }
  } catch (error) {
    console.error('获取下载统计失败:', error)
  }
}

// 查看资源详情
const handleView = (row: any) => {
  navigateTo(`/r/${row.resource_key}`)
}

// 删除单条记录
const handleRemoveRecord = (row: any) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除「${row.title || '该资源'}」的下载记录吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await useApiFetch(`/user/download-history/${row.id}`, { method: 'DELETE' })
        if (process.client) {
          notification.success({ content: '记录已删除', duration: 3000 })
        }
        await Promise.all([fetchHistory(), fetchStats()])
      } catch (error: any) {
        if (process.client) {
          notification.error({ content: error?.message || '删除失败', duration: 3000 })
        }
      }
    }
  })
}

// 清空全部历史
const handleClearHistory = () => {
  dialog.warning({
    title: '确认清空历史',
    content: '确定要清空所有下载历史记录吗？此操作不可撤销。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await useApiFetch('/user/download-history', { method: 'DELETE' })
        if (process.client) {
          notification.success({ content: '历史记录已清空', duration: 3000 })
        }
        pagination.value.page = 1
        await Promise.all([fetchHistory(), fetchStats()])
      } catch (error: any) {
        if (process.client) {
          notification.error({ content: error?.message || '清空失败', duration: 3000 })
        }
      }
    }
  })
}

// 页面加载时获取数据
onMounted(() => {
  fetchHistory()
  fetchStats()
})
</script>
