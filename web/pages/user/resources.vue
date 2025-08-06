<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">我的资源</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">管理您的个人资源</p>
      </div>
      <n-button type="primary" @click="handleAddResource">
        <template #icon>
          <i class="fas fa-plus"></i>
        </template>
        添加资源
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <i class="fas fa-cloud text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总资源数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <i class="fas fa-check text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">有效资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.valid || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-exclamation-triangle text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">无效资源</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.invalid || 0 }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 筛选和搜索 -->
    <n-card title="筛选和搜索" :bordered="false">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索资源..."
          clearable
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-select
          v-model:value="filterStatus"
          :options="statusOptions"
          placeholder="状态筛选"
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

    <!-- 资源列表 -->
    <n-card title="资源列表" :bordered="false">
      <template #header-extra>
        <n-space>
          <n-tag type="info" size="small">{{ resources.length }} 个资源</n-tag>
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
      
      <div v-if="resources.length === 0" class="text-center py-12">
        <n-empty description="暂无资源">
          <template #icon>
            <i class="fas fa-cloud text-gray-400 text-4xl"></i>
          </template>
          <template #extra>
            <n-button type="primary" @click="handleAddResource">
              <template #icon>
                <i class="fas fa-plus"></i>
              </template>
              添加第一个资源
            </n-button>
          </template>
        </n-empty>
      </div>
      
      <div v-else>
        <n-data-table
          :columns="columns"
          :data="resources"
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
  title: '我的资源'
})

// 响应式数据
const stats = ref({
  total: 0,
  valid: 0,
  invalid: 0
})

const searchKeyword = ref('')
const filterStatus = ref(null)
const filterPlatform = ref(null)
const loading = ref(false)

// 筛选选项
const statusOptions = [
  { label: '全部', value: '' },
  { label: '有效', value: 'valid' },
  { label: '无效', value: 'invalid' }
]

const platformOptions = [
  { label: '全部', value: '' },
  { label: '百度网盘', value: 'baidu' },
  { label: '阿里云盘', value: 'alipan' },
  { label: '夸克网盘', value: 'quark' },
  { label: 'UC网盘', value: 'uc' }
]

// 资源数据
const resources = ref<any[]>([])

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.value.page = page
    fetchResources()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchResources()
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
    title: '状态',
    key: 'status',
    render: (row: any) => {
      return h('n-tag', {
        type: row.is_valid ? 'success' : 'error',
        size: 'small'
      }, { default: () => row.is_valid ? '有效' : '无效' })
    }
  },
  {
    title: '添加时间',
    key: 'created_at',
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDate(row.created_at))
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
            type: 'info',
            onClick: () => handleEdit(row)
          }, { default: () => '编辑' }),
          h('n-button', {
            size: 'small',
            type: 'error',
            onClick: () => handleDelete(row)
          }, { default: () => '删除' })
        ]
      })
    }
  }
]

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 处理添加资源
const handleAddResource = () => {
  // TODO: 实现添加资源功能
  const notification = useNotification()
  notification.info({
    content: '添加资源功能开发中...',
    duration: 3000
  })
}

// 处理搜索
const handleSearch = () => {
  pagination.value.page = 1
  fetchResources()
}

// 处理刷新
const handleRefresh = () => {
  fetchResources()
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
const handleView = (resource: any) => {
  if (resource.url) {
    window.open(resource.url, '_blank')
  }
}

// 处理编辑资源
const handleEdit = (resource: any) => {
  const notification = useNotification()
  notification.info({
    content: '编辑功能开发中...',
    duration: 3000
  })
}

// 处理删除资源
const handleDelete = (resource: any) => {
  const dialog = useDialog()
  dialog.warning({
    title: '确认删除',
    content: `确定要删除资源"${resource.title}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      const notification = useNotification()
      notification.success({
        content: '删除成功',
        duration: 3000
      })
    }
  })
}

// 获取资源数据
const fetchResources = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取用户资源
    // const response = await userApi.getUserResources({
    //   page: pagination.value.page,
    //   pageSize: pagination.value.pageSize,
    //   keyword: searchKeyword.value,
    //   status: filterStatus.value,
    //   platform: filterPlatform.value
    // })
    
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
    resources.value = []
    stats.value = {
      total: 0,
      valid: 0,
      invalid: 0
    }
  } catch (error) {
    console.error('获取资源失败:', error)
    const notification = useNotification()
    notification.error({
      content: '获取资源失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchResources()
})
</script> 