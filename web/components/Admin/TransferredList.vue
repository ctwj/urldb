<template>
  <div class="space-y-4">
    <!-- 搜索和筛选 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索已转存资源..."
        @keyup.enter="handleSearch"
        clearable
      >
        <template #prefix>
          <i class="fas fa-search"></i>
        </template>
      </n-input>
      
      <CategorySelector
        v-model="selectedCategory"
        placeholder="选择分类"
        clearable
      />
      
      <TagSelector
        v-model="selectedTag"
        placeholder="选择标签"
        clearable
      />
      
      <n-button type="primary" @click="handleSearch">
        <template #icon>
          <i class="fas fa-search"></i>
        </template>
        搜索
      </n-button>
    </div>

    <!-- 数据表格 -->
    <n-data-table
      :columns="columns"
      :data="resources"
      :loading="loading"
      :pagination="pagination"
      :remote="true"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
      :row-key="(row: any) => row.id"
      max-height="500"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useResourceApi } from '~/composables/useApi'

// 数据状态
const loading = ref(false)
const resources = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索条件
const searchQuery = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)

// API实例
const resourceApi = useResourceApi()

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  pageSizes: [10, 20, 50, 100],
  showSizePicker: true,
  showQuickJumper: true,
  prefix: ({ itemCount }: any) => `共 ${itemCount} 条`
})

// 表格列配置
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 60,
    fixed: 'left'
  },
  {
    title: '标题',
    key: 'title',
    width: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '分类',
    key: 'category_name',
    width: 80
  },
  {
    title: '平台',
    key: 'platform_name',
    width: 80,
    render: (row: any) => {
      const platform = platformOptions.value.find((p: any) => p.value === row.pan_id)
      return platform?.label || '未知'
    }
  },
  {
    title: '转存链接',
    key: 'save_url',
    width: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      return h('a', {
        href: row.save_url,
        target: '_blank',
        class: 'text-green-500 hover:text-green-700'
      }, row.save_url.length > 30 ? row.save_url.substring(0, 30) + '...' : row.save_url)
    }
  },
  {
    title: '转存时间',
    key: 'updated_at',
    width: 130,
    render: (row: any) => {
      return new Date(row.updated_at).toLocaleDateString()
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    fixed: 'right',
    render: (row: any) => {
      return [
        h('n-button', {
          size: 'small',
          type: 'primary',
          onClick: () => viewResource(row)
        }, '查看'),
        h('n-button', {
          size: 'small',
          type: 'info',
          style: { marginLeft: '8px' },
          onClick: () => copyLink(row.save_url)
        }, '复制')
      ]
    }
  }
]

// 平台选项
const platformOptions = ref([])

// 获取已转存资源
const fetchTransferredResources = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      has_save_url: true // 筛选有转存链接的资源
    }

    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }

    const result = await resourceApi.getResources(params) as any
    console.log('已转存资源结果:', result)

    if (result && result.resources) {
      resources.value = result.resources
      total.value = result.total || 0
      pagination.itemCount = result.total || 0
    } else if (Array.isArray(result)) {
      resources.value = result
      total.value = result.length
      pagination.itemCount = result.length
    }
  } catch (error) {
    console.error('获取已转存资源失败:', error)
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  pagination.page = 1
  fetchTransferredResources()
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  pagination.page = page
  fetchTransferredResources()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  pagination.pageSize = size
  currentPage.value = 1
  pagination.page = 1
  fetchTransferredResources()
}

// 查看资源
const viewResource = (resource: any) => {
  // 这里可以打开资源详情模态框
  console.log('查看资源:', resource)
}

// 复制链接
const copyLink = async (url: string) => {
  try {
    await navigator.clipboard.writeText(url)
    $message.success('链接已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    $message.error('复制失败')
  }
}

// 初始化
onMounted(() => {
  fetchTransferredResources()
})
</script>