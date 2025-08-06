<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">浏览历史</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">查看您的浏览历史记录</p>
      </div>
      <n-button type="primary" @click="handleClearHistory">
        <template #icon>
          <i class="fas fa-trash"></i>
        </template>
        清空历史
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-history text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总记录数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-calendar-day text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日浏览</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.today || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
            <i class="fas fa-calendar-week text-yellow-600 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本周浏览</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.thisWeek || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-calendar-month text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">本月浏览</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.thisMonth || 0 }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 筛选和搜索 -->
    <n-card title="筛选和搜索" :bordered="false">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索历史记录..."
          clearable
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-date-picker
          v-model:value="dateRange"
          type="daterange"
          placeholder="选择日期范围"
          clearable
        />
        
        <n-select
          v-model:value="filterPlatform"
          :options="platformOptions"
          placeholder="平台筛选"
          clearable
        />
        
        <n-button type="primary" @click="handleSearch">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
      </div>
    </n-card>

    <!-- 历史记录列表 -->
    <n-card title="浏览历史" :bordered="false">
      <template #header-extra>
        <n-space>
          <n-tag type="info" size="small">{{ historyRecords.length }} 条记录</n-tag>
          <n-button-group>
            <n-button size="small" @click="handleRefresh">
              <template #icon>
                <i class="fas fa-refresh"></i>
              </template>
              刷新
            </n-button>
            <n-button size="small" @click="handleExport">
              <template #icon>
                <i class="fas fa-download"></i>
              </template>
              导出
            </n-button>
          </n-button-group>
        </n-space>
      </template>
      
      <div v-if="historyRecords.length === 0" class="text-center py-12">
        <n-empty description="暂无浏览历史">
          <template #icon>
            <i class="fas fa-history text-gray-400 text-4xl"></i>
          </template>
          <template #extra>
            <n-button type="primary" @click="navigateTo('/')">
              <template #icon>
                <i class="fas fa-search"></i>
              </template>
              去浏览资源
            </n-button>
          </template>
        </n-empty>
      </div>
      
      <div v-else>
        <n-data-table
          :columns="columns"
          :data="historyRecords"
          :pagination="pagination"
          :loading="loading"
          striped
        />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// 页面元数据
definePageMeta({
  layout: 'user',
  title: '浏览历史'
})

// 响应式数据
const stats = ref({
  total: 0,
  today: 0,
  thisWeek: 0,
  thisMonth: 0
})

const searchKeyword = ref('')
const dateRange = ref(null)
const filterPlatform = ref('')
const loading = ref(false)

// 筛选选项
const platformOptions = [
  { label: '全部', value: '' },
  { label: '百度网盘', value: 'baidu' },
  { label: '阿里云盘', value: 'alipan' },
  { label: '夸克网盘', value: 'quark' },
  { label: 'UC网盘', value: 'uc' }
]

// 历史记录数据
const historyRecords = ref<any[]>([])

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
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
        h('div', { class: 'font-medium' }, row.title),
        h('div', { class: 'text-sm text-gray-500' }, row.description)
      ])
    }
  },
  {
    title: '平台',
    key: 'platform',
    render: (row: any) => {
      const platformMap: Record<string, any> = {
        baidu: { label: '百度网盘', color: 'blue' },
        alipan: { label: '阿里云盘', color: 'green' },
        quark: { label: '夸克网盘', color: 'purple' },
        uc: { label: 'UC网盘', color: 'orange' }
      }
      const platform = platformMap[row.platform] || { label: '未知', color: 'default' }
      return h('n-tag', { type: platform.color, size: 'small' }, { default: () => platform.label })
    }
  },
  {
    title: '浏览时间',
    key: 'viewed_at',
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDateTime(row.viewed_at))
    }
  },
  {
    title: '停留时长',
    key: 'duration',
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDuration(row.duration))
    }
  },
  {
    title: '操作',
    key: 'actions',
    render: (row: any) => {
      return h('n-space', { size: 'small' }, {
        default: () => [
          h('n-button', {
            size: 'small',
            type: 'primary',
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

// 格式化时长
const formatDuration = (seconds: number) => {
  if (!seconds) return '未知'
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}分${remainingSeconds}秒`
}

// 处理清空历史
const handleClearHistory = () => {
  const dialog = useDialog()
  dialog.warning({
    title: '确认清空历史',
    content: '确定要清空所有浏览历史记录吗？此操作不可撤销。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      const notification = useNotification()
      notification.success({
        content: '历史记录已清空',
        duration: 3000
      })
      fetchHistory()
    }
  })
}

// 处理搜索
const handleSearch = () => {
  pagination.value.page = 1
  fetchHistory()
}

// 处理刷新
const handleRefresh = () => {
  fetchHistory()
}

// 处理导出
const handleExport = () => {
  const notification = useNotification()
  notification.info({
    content: '导出功能开发中...',
    duration: 3000
  })
}

// 处理查看资源
const handleView = (record: any) => {
  if (record.url) {
    window.open(record.url, '_blank')
  }
}

// 处理删除记录
const handleRemoveRecord = (record: any) => {
  const dialog = useDialog()
  dialog.warning({
    title: '确认删除',
    content: `确定要删除这条浏览记录吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      const notification = useNotification()
      notification.success({
        content: '记录删除成功',
        duration: 3000
      })
      fetchHistory()
    }
  })
}

// 获取历史记录数据
const fetchHistory = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取用户浏览历史
    // const response = await userApi.getUserHistory({
    //   page: pagination.value.page,
    //   pageSize: pagination.value.pageSize,
    //   keyword: searchKeyword.value,
    //   dateRange: dateRange.value,
    //   platform: filterPlatform.value
    // })
    
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
    historyRecords.value = []
    stats.value = {
      total: 0,
      today: 0,
      thisWeek: 0,
      thisMonth: 0
    }
  } catch (error) {
    console.error('获取历史记录失败:', error)
    const notification = useNotification()
    notification.error({
      content: '获取历史记录失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchHistory()
})
</script> 