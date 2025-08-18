<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">收藏夹</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">管理您收藏的资源</p>
      </div>
      <n-button type="primary" @click="handleSyncFavorites">
        <template #icon>
          <i class="fas fa-sync-alt"></i>
        </template>
        同步收藏
      </n-button>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <i class="fas fa-heart text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总收藏数</p>
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
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">有效收藏</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.valid || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
            <i class="fas fa-clock text-yellow-600 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">最近收藏</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.recent || 0 }}</p>
          </div>
        </div>
      </n-card>

      <n-card>
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
            <i class="fas fa-folder text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">收藏夹</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.folders || 0 }}</p>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 筛选和搜索 -->
    <n-card title="筛选和搜索" :bordered="false">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索收藏..."
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
        
        <n-button type="primary" @click="handleSearch" class="w-20">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
      </div>
    </n-card>

    <!-- 收藏列表 -->
    <n-card title="收藏列表" :bordered="false">
      <template #header-extra>
        <n-space>
          <n-tag type="info" size="small">{{ favorites.length }} 个收藏</n-tag>
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
      
      <div v-if="favorites.length === 0" class="text-center py-12">
        <n-empty description="暂无收藏">
          <template #icon>
            <i class="fas fa-heart text-gray-400 text-4xl"></i>
          </template>
          <template #extra>
            <n-button type="primary" @click="navigateTo('/')">
              <template #icon>
                <i class="fas fa-search"></i>
              </template>
              去发现资源
            </n-button>
          </template>
        </n-empty>
      </div>
      
      <div v-else>
        <n-data-table
          :columns="columns"
          :data="favorites"
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
  title: '收藏夹'
})

// 响应式数据
const stats = ref({
  total: 0,
  valid: 0,
  recent: 0,
  folders: 0
})

const searchKeyword = ref('')
const filterStatus = ref('')
const filterPlatform = ref('')
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

// 收藏数据
const favorites = ref<any[]>([])

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.value.page = page
    fetchFavorites()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchFavorites()
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
    title: '收藏时间',
    key: 'favorited_at',
    render: (row: any) => {
      return h('span', { class: 'text-sm text-gray-500' }, formatDate(row.favorited_at))
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
            onClick: () => handleRemoveFavorite(row)
          }, { default: () => '取消收藏' })
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

// 处理同步收藏
const handleSyncFavorites = () => {
  const notification = useNotification()
  notification.info({
    content: '同步收藏功能开发中...',
    duration: 3000
  })
}

// 处理搜索
const handleSearch = () => {
  pagination.value.page = 1
  fetchFavorites()
}

// 处理刷新
const handleRefresh = () => {
  fetchFavorites()
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
const handleView = (favorite: any) => {
  if (favorite.url) {
    window.open(favorite.url, '_blank')
  }
}

// 处理取消收藏
const handleRemoveFavorite = (favorite: any) => {
  const dialog = useDialog()
  dialog.warning({
    title: '确认取消收藏',
    content: `确定要取消收藏"${favorite.title}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      const notification = useNotification()
      notification.success({
        content: '取消收藏成功',
        duration: 3000
      })
      fetchFavorites()
    }
  })
}

// 获取收藏数据
const fetchFavorites = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取用户收藏
    // const response = await userApi.getUserFavorites({
    //   page: pagination.value.page,
    //   pageSize: pagination.value.pageSize,
    //   keyword: searchKeyword.value,
    //   status: filterStatus.value,
    //   platform: filterPlatform.value
    // })
    
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
    favorites.value = []
    stats.value = {
      total: 0,
      valid: 0,
      recent: 0,
      folders: 0
    }
  } catch (error) {
    console.error('获取收藏失败:', error)
    const notification = useNotification()
    notification.error({
      content: '获取收藏失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchFavorites()
})
</script> 