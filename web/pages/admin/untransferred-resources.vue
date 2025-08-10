<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">未转存列表</h1>
        <p class="text-gray-600 dark:text-gray-400">显示夸克网盘中尚未转存的资源</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="refreshData" type="info">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
        <n-button @click="batchTransfer" type="primary" :disabled="selectedResources.length === 0">
          <template #icon>
            <i class="fas fa-cloud-upload-alt"></i>
          </template>
          批量转存 ({{ selectedResources.length }})
        </n-button>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <i class="fas fa-clock text-orange-500 text-2xl"></i>
          </div>
          <div class="ml-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">待转存总数</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ total }}</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <i class="fas fa-check-circle text-green-500 text-2xl"></i>
          </div>
          <div class="ml-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">已选择</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ selectedResources.length }}</div>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <i class="fas fa-quark text-blue-500 text-2xl"></i>
          </div>
          <div class="ml-4">
            <div class="text-sm font-medium text-gray-500 dark:text-gray-400">夸克网盘</div>
            <div class="text-2xl font-bold text-gray-900 dark:text-white">{{ total }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <n-card>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索资源标题..."
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-select
          v-model:value="selectedCategory"
          placeholder="选择分类"
          :options="categoryOptions"
          clearable
        />
        
        <n-select
          v-model:value="sortBy"
          placeholder="排序方式"
          :options="sortOptions"
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
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <span class="text-lg font-semibold">未转存资源列表</span>
            <div class="flex items-center space-x-2">
              <n-checkbox 
                :checked="isAllSelected"
                @update:checked="toggleSelectAll"
                :indeterminate="isIndeterminate"
              />
              <span class="text-sm text-gray-500">全选</span>
            </div>
          </div>
          <span class="text-sm text-gray-500">共 {{ total }} 个资源，已选择 {{ selectedResources.length }} 个</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="resources.length === 0" class="text-center py-8">
        <i class="fas fa-check-circle text-4xl text-green-400 mb-4"></i>
        <p class="text-gray-500">暂无未转存的资源</p>
      </div>

      <div v-else>
        <n-data-table
          :columns="columns"
          :data="resources"
          :pagination="paginationConfig"
          :bordered="false"
          size="small"
          :scroll-x="800"
          @update:checked-row-keys="handleSelectionChange"
        />
      </div>
    </n-card>

    <!-- 批量转存确认对话框 -->
    <n-modal v-model:show="showBatchTransferModal" preset="dialog" title="确认批量转存">
      <div class="space-y-4">
        <p>确定要将选中的 <strong>{{ selectedResources.length }}</strong> 个资源进行批量转存吗？</p>
        <div class="bg-yellow-50 dark:bg-yellow-900/20 p-3 rounded border border-yellow-200 dark:border-yellow-800">
          <div class="flex items-start space-x-2">
            <i class="fas fa-exclamation-triangle text-yellow-500 mt-0.5"></i>
            <div class="text-sm text-yellow-800 dark:text-yellow-200">
              <p>• 转存过程可能需要较长时间</p>
              <p>• 请确保夸克网盘账号有足够的存储空间</p>
              <p>• 转存完成后可在"已转存列表"中查看结果</p>
            </div>
          </div>
        </div>
      </div>
      <template #action>
        <div class="flex space-x-2">
          <n-button @click="showBatchTransferModal = false">取消</n-button>
          <n-button type="primary" @click="confirmBatchTransfer" :loading="transferring">
            {{ transferring ? '转存中...' : '确认转存' }}
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

import { ref, computed, onMounted } from 'vue'
import { useResourceApi, useCategoryApi, useTaskApi } from '~/composables/useApi'

// API实例
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const taskApi = useTaskApi()

// 数据状态
const resources = ref([])
const categories = ref([])
const loading = ref(false)
const transferring = ref(false)
const total = ref(0)
const selectedResourceIds = ref([])

// 搜索和筛选
const searchQuery = ref('')
const selectedCategory = ref(null)
const sortBy = ref('created_at')

// 分页
const currentPage = ref(1)
const pageSize = ref(20)

// 模态框
const showBatchTransferModal = ref(false)

// 排序选项
const sortOptions = [
  { label: '创建时间 (最新)', value: 'created_at' },
  { label: '创建时间 (最早)', value: 'created_at_asc' },
  { label: '更新时间 (最新)', value: 'updated_at' },
  { label: '标题 (A-Z)', value: 'title' },
  { label: '标题 (Z-A)', value: 'title_desc' }
]

// 分类选项
const categoryOptions = computed(() => {
  return categories.value.map(cat => ({
    label: cat.name,
    value: cat.id
  }))
})

// 分页配置
const paginationConfig = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchResources()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchResources()
  }
}))

// 选择状态
const selectedResources = computed(() => {
  return resources.value.filter(r => selectedResourceIds.value.includes(r.id))
})

const isAllSelected = computed(() => {
  return resources.value.length > 0 && selectedResourceIds.value.length === resources.value.length
})

const isIndeterminate = computed(() => {
  return selectedResourceIds.value.length > 0 && selectedResourceIds.value.length < resources.value.length
})

// 表格列定义
const columns = [
  {
    type: 'selection',
    width: 50
  },
  {
    title: 'ID',
    key: 'id',
    width: 80,
    minWidth: 80
  },
  {
    title: '标题',
    key: 'title',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '分类',
    key: 'category',
    width: 120,
    render: (row: any) => {
      if (!row.category) return '-'
      return h('n-tag', { type: 'info', size: 'small' }, { default: () => row.category.name })
    }
  },
  {
    title: '原始链接',
    key: 'url',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      return h('a', {
        href: row.url,
        target: '_blank',
        class: 'text-blue-600 hover:text-blue-800 text-xs break-all'
      }, row.url)
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: (row: any) => {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 160,
    render: (row: any) => {
      return new Date(row.updated_at).toLocaleString('zh-CN')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row: any) => {
      return h('div', { class: 'flex space-x-1' }, [
        h('n-button', {
          size: 'small',
          type: 'primary',
          onClick: () => singleTransfer(row)
        }, { default: () => '转存' }),
        h('n-button', {
          size: 'small',
          type: 'default',
          onClick: () => viewResource(row.id)
        }, { default: () => '查看' })
      ])
    }
  }
]

// 获取未转存资源
const fetchResources = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      pan_name: 'quark', // 只获取夸克网盘的资源
      has_save_url: false, // 只获取没有转存链接的资源
      search: searchQuery.value,
      category_id: selectedCategory.value,
      sort_by: sortBy.value
    }
    
    const response = await resourceApi.getResources(params)
    resources.value = response.resources || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取未转存资源失败:', error)
    notification.error({
      content: '获取未转存资源失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchResources()
}

// 刷新数据
const refreshData = () => {
  fetchResources()
}

// 选择处理
const handleSelectionChange = (keys: any[]) => {
  selectedResourceIds.value = keys
}

// 全选/取消全选
const toggleSelectAll = (checked: boolean) => {
  if (checked) {
    selectedResourceIds.value = resources.value.map(r => r.id)
  } else {
    selectedResourceIds.value = []
  }
}

// 单个转存
const singleTransfer = async (resource: any) => {
  try {
    const taskData = {
      title: `转存资源: ${resource.title}`,
      description: `转存单个资源: ${resource.title}`,
      resources: [{
        title: resource.title,
        url: resource.url,
        category_id: resource.category_id || 0
      }]
    }
    
    const response = await taskApi.createBatchTransferTask(taskData)
    notification.success({
      content: '转存任务已创建',
      duration: 3000
    })
    
    // 跳转到任务详情页
    navigateTo(`/admin/tasks/${response.id}`)
  } catch (error) {
    console.error('创建转存任务失败:', error)
    notification.error({
      content: '创建转存任务失败',
      duration: 3000
    })
  }
}

// 批量转存
const batchTransfer = () => {
  if (selectedResources.value.length === 0) {
    notification.warning({
      content: '请先选择要转存的资源',
      duration: 3000
    })
    return
  }
  showBatchTransferModal.value = true
}

// 确认批量转存
const confirmBatchTransfer = async () => {
  transferring.value = true
  try {
    const taskData = {
      title: `批量转存 ${selectedResources.value.length} 个资源`,
      description: `批量转存 ${selectedResources.value.length} 个夸克网盘资源`,
      resources: selectedResources.value.map(r => ({
        title: r.title,
        url: r.url,
        category_id: r.category_id || 0
      }))
    }
    
    const response = await taskApi.createBatchTransferTask(taskData)
    notification.success({
      content: `批量转存任务已创建，共 ${selectedResources.value.length} 个资源`,
      duration: 3000
    })
    
    // 跳转到任务详情页
    navigateTo(`/admin/tasks/${response.id}`)
  } catch (error) {
    console.error('创建批量转存任务失败:', error)
    notification.error({
      content: '创建批量转存任务失败',
      duration: 3000
    })
  } finally {
    transferring.value = false
    showBatchTransferModal.value = false
  }
}

// 查看资源详情
const viewResource = (id: number) => {
  navigateTo(`/admin/resources/${id}`)
}

// 页面加载
onMounted(async () => {
  await Promise.all([
    fetchCategories(),
    fetchResources()
  ])
})
</script>

<style scoped>
/* 自定义样式 */
</style>
