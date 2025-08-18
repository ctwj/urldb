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

    <!-- 调试信息 -->
    <div class="text-sm text-gray-500 mb-2">
      数据数量: {{ resources.length }}, 总数: {{ total }}, 加载状态: {{ loading }}
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
      virtual-scroll
      max-height="500"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useResourceApi, usePanApi } from '~/composables/useApi'
import { useMessage } from 'naive-ui'

// 消息提示
const $message = useMessage()

// 数据状态
const loading = ref(false)
const resources = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10000)

// 搜索条件
const searchQuery = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)

// API实例
const resourceApi = useResourceApi()
const panApi = usePanApi()

// 获取平台数据
const { data: platformsData } = await useAsyncData('transferredPlatforms', () => panApi.getPans())

// 平台选项
const platformOptions = computed(() => {
  const data = platformsData.value as any
  const platforms = data?.data || data || []
  return platforms.map((platform: any) => ({
    label: platform.remark || platform.name,
    value: platform.id
  }))
})

// 获取平台名称
const getPlatformName = (platformId: number) => {
  const platform = (platformsData.value as any)?.data?.find((plat: any) => plat.id === platformId)
  return platform?.remark || platform?.name || '未知平台'
}

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10000,
  itemCount: 0,
  pageSizes: [10000, 20000, 50000, 100000],
  showSizePicker: true,
  showQuickJumper: true,
  prefix: ({ itemCount }: any) => `共 ${itemCount} 条`
})

// 表格列配置
const columns: any[] = [
  {
    title: 'ID',
    key: 'id',
    width: 60,
    fixed: 'left' as const
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
  }
]

// 获取已转存资源
const fetchTransferredResources = async () => {
  loading.value = true
  try {
    const params: any = {
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

    console.log('请求参数:', params)
    const result = await resourceApi.getResources(params) as any
    console.log('已转存资源结果:', result)
    console.log('结果类型:', typeof result)
    console.log('结果结构:', Object.keys(result || {}))

    if (result && result.data) {
      console.log('使用 resources 格式，数量:', result.data.length)
      resources.value = result.data
      total.value = result.total || 0
      pagination.itemCount = result.total || 0
    } else if (Array.isArray(result)) {
      console.log('使用数组格式，数量:', result.length)
      resources.value = result
      total.value = result.length
      pagination.itemCount = result.length
    } else {
      console.log('未知格式，设置空数组')
      resources.value = []
      total.value = 0
      pagination.itemCount = 0
    }
    
    console.log('最终 resources.value:', resources.value)
    console.log('最终 total.value:', total.value)
    
    // 检查是否有资源没有 save_url
    const resourcesWithoutSaveUrl = resources.value.filter((r: any) => !r.save_url || r.save_url.trim() === '')
    if (resourcesWithoutSaveUrl.length > 0) {
      console.warn('发现没有 save_url 的资源:', resourcesWithoutSaveUrl.map((r: any) => ({ id: r.id, title: r.title, save_url: r.save_url })))
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

// 初始化
onMounted(() => {
  fetchTransferredResources()
})
</script>