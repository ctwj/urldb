<template>
  <div class="space-y-4">
    <!-- 搜索和筛选 -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索未转存资源..."
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

      <n-select
        v-model:value="selectedStatus"
        placeholder="资源状态"
        :options="statusOptions"
        clearable
      />
      
      <n-button type="primary" @click="handleSearch">
        <template #icon>
          <i class="fas fa-search"></i>
        </template>
        搜索
      </n-button>
    </div>

    <!-- 批量操作 -->
    <n-card>
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <n-checkbox 
              :checked="isAllSelected"
              @update:checked="toggleSelectAll"
              :indeterminate="isIndeterminate"
            />
            <span class="text-sm text-gray-600 dark:text-gray-400">全选</span>
          </div>
          <span class="text-sm text-gray-500">
            共 {{ total }} 个资源，已选择 {{ selectedResources.length }} 个
          </span>
        </div>
        
        <div class="flex space-x-2">
          <n-button 
            type="primary"
            :disabled="selectedResources.length === 0"
            :loading="batchTransferring"
            @click="handleBatchTransfer"
          >
            <template #icon>
              <i class="fas fa-exchange-alt"></i>
            </template>
            批量转存 ({{ selectedResources.length }})
          </n-button>
          
          <n-button @click="refreshData">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- 资源列表 -->
    <n-card>
      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="resources.length === 0" class="text-center py-8">
        <i class="fas fa-inbox text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无未转存的夸克资源</p>
      </div>

      <div v-else>
        <!-- 虚拟列表 -->
        <n-virtual-list
          :items="resources"
          :item-size="120"
          style="max-height: 500px"
          container-style="height: 500px;"
        >
          <template #default="{ item }">
            <div class="p-4 border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
              <div class="flex items-start space-x-4">
                <!-- 选择框 -->
                <div class="pt-2">
                  <n-checkbox 
                    :checked="selectedResources.includes(item.id)"
                    @update:checked="(checked) => toggleResourceSelection(item.id, checked)"
                  />
                </div>

                <!-- 资源信息 -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between">
                    <div class="flex-1 min-w-0">
                      <!-- 标题和状态 -->
                      <div class="flex items-center space-x-2 mb-2">
                        <h3 class="text-lg font-medium text-gray-900 dark:text-white line-clamp-1">
                          {{ item.title || '未命名资源' }}
                        </h3>
                        <n-tag :type="getStatusType(item)" size="small">
                          {{ getStatusText(item) }}
                        </n-tag>
                      </div>

                      <!-- 描述 -->
                      <p class="text-gray-600 dark:text-gray-400 text-sm line-clamp-2 mb-2">
                        {{ item.description || '暂无描述' }}
                      </p>

                      <!-- 元信息 -->
                      <div class="flex items-center space-x-4 text-sm text-gray-500">
                        <span class="flex items-center">
                          <i class="fas fa-folder mr-1"></i>
                          {{ item.category_name || '未分类' }}
                        </span>
                        <span class="flex items-center">
                          <i class="fas fa-cloud mr-1"></i>
                          夸克网盘
                        </span>
                        <span class="flex items-center">
                          <i class="fas fa-eye mr-1"></i>
                          {{ item.view_count || 0 }} 次浏览
                        </span>
                        <span class="flex items-center">
                          <i class="fas fa-calendar mr-1"></i>
                          {{ formatDate(item.created_at) }}
                        </span>
                      </div>

                      <!-- 原始链接 -->
                      <div class="mt-2">
                        <div class="flex items-center space-x-2">
                          <span class="text-xs text-gray-400">原始链接:</span>
                          <a 
                            :href="item.url" 
                            target="_blank"
                            class="text-xs text-blue-500 hover:text-blue-700 truncate max-w-xs"
                          >
                            {{ item.url }}
                          </a>
                        </div>
                      </div>
                    </div>

                    <!-- 操作按钮 -->
                    <div class="flex flex-col space-y-2 ml-4">
                      <n-button 
                        size="small" 
                        type="primary"
                        :loading="item.transferring"
                        @click="handleSingleTransfer(item)"
                      >
                        <template #icon>
                          <i class="fas fa-exchange-alt"></i>
                        </template>
                        {{ item.transferring ? '转存中' : '立即转存' }}
                      </n-button>
                      
                      <n-button 
                        size="small" 
                        @click="viewResource(item)"
                      >
                        <template #icon>
                          <i class="fas fa-eye"></i>
                        </template>
                        查看详情
                      </n-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </n-virtual-list>

        <!-- 分页 -->
        <div class="mt-4 flex justify-center">
          <n-pagination
            v-model:page="currentPage"
            v-model:page-size="pageSize"
            :item-count="total"
            :page-sizes="[10000, 20000, 50000, 100000]"
            show-size-picker
            show-quick-jumper
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </n-card>

    <!-- 转存结果模态框 -->
    <n-modal v-model:show="showTransferResult" preset="card" title="转存结果" style="width: 600px">
      <div v-if="transferResults.length > 0" class="space-y-4">
        <div class="grid grid-cols-3 gap-4">
          <div class="text-center p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
            <div class="text-xl font-bold text-green-600">{{ transferSuccessCount }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">成功</div>
          </div>
          <div class="text-center p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
            <div class="text-xl font-bold text-red-600">{{ transferFailedCount }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">失败</div>
          </div>
          <div class="text-center p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
            <div class="text-xl font-bold text-blue-600">{{ transferResults.length }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">总计</div>
          </div>
        </div>

        <div class="max-h-300 overflow-y-auto">
          <div v-for="result in transferResults" :key="result.id" class="p-3 border rounded mb-2">
            <div class="flex items-center justify-between">
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium truncate">{{ result.title }}</div>
                <div class="text-xs text-gray-500 truncate">{{ result.url }}</div>
              </div>
              <n-tag :type="result.success ? 'success' : 'error'" size="small">
                {{ result.success ? '成功' : '失败' }}
              </n-tag>
            </div>
            <div v-if="result.message" class="text-xs text-gray-600 mt-1">
              {{ result.message }}
            </div>
          </div>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useResourceApi, useCategoryApi, useTagApi } from '~/composables/useApi'

// 数据状态
const loading = ref(false)
const resources = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10000)

// 搜索条件
const searchQuery = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)
const selectedStatus = ref(null)

// 选择状态
const selectedResources = ref([])

// 批量操作状态
const batchTransferring = ref(false)
const showTransferResult = ref(false)
const transferResults = ref([])

// 选项数据
const categoryOptions = ref([])
const tagOptions = ref([])
const statusOptions = [
  { label: '有效', value: 'valid' },
  { label: '无效', value: 'invalid' },
  { label: '待验证', value: 'pending' }
]

// API实例
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()

// 计算属性
const isAllSelected = computed(() => {
  return resources.value.length > 0 && selectedResources.value.length === resources.value.length
})

const isIndeterminate = computed(() => {
  return selectedResources.value.length > 0 && selectedResources.value.length < resources.value.length
})

const transferSuccessCount = computed(() => {
  return transferResults.value.filter(r => r.success).length
})

const transferFailedCount = computed(() => {
  return transferResults.value.filter(r => !r.success).length
})

// 获取未转存资源（夸克网盘且无save_url）
const fetchUntransferredResources = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      no_save_url: true, // 筛选没有转存链接的资源
      pan_name: 'quark' // 仅夸克网盘资源
    }

    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }

    const result = await resourceApi.getResources(params) as any
    console.log('未转存资源结果:', result)

    if (result && result.data) {
      resources.value = result.data
      total.value = result.total || 0
    } else if (Array.isArray(result)) {
      resources.value = result
      total.value = result.length
    }

    // 清空选择
    selectedResources.value = []
  } catch (error) {
    console.error('获取未转存资源失败:', error)
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 获取分类选项
const fetchCategories = async () => {
  try {
    const result = await categoryApi.getCategories() as any
    if (result && result.items) {
      categoryOptions.value = result.items.map((item: any) => ({
        label: item.name,
        value: item.id
      }))
    }
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取标签选项
const fetchTags = async () => {
  try {
    const result = await tagApi.getTags() as any
    if (result && result.items) {
      tagOptions.value = result.items.map((item: any) => ({
        label: item.name,
        value: item.id
      }))
    }
  } catch (error) {
    console.error('获取标签失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchUntransferredResources()
}

// 刷新数据
const refreshData = () => {
  fetchUntransferredResources()
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchUntransferredResources()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchUntransferredResources()
}

// 选择处理
const toggleSelectAll = (checked: boolean) => {
  if (checked) {
    selectedResources.value = resources.value.map(r => r.id)
  } else {
    selectedResources.value = []
  }
}

const toggleResourceSelection = (id: number, checked: boolean) => {
  if (checked) {
    if (!selectedResources.value.includes(id)) {
      selectedResources.value.push(id)
    }
  } else {
    const index = selectedResources.value.indexOf(id)
    if (index > -1) {
      selectedResources.value.splice(index, 1)
    }
  }
}

// 单个转存
const handleSingleTransfer = async (resource: any) => {
  resource.transferring = true
  
  try {
    // 这里应该调用实际的转存API
    // 由于只是UI展示，模拟转存过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // 模拟随机成功/失败
    const isSuccess = Math.random() > 0.3
    
    if (isSuccess) {
      $message.success(`${resource.title} 转存成功`)
      // 刷新列表
      refreshData()
    } else {
      $message.error(`${resource.title} 转存失败`)
    }
  } catch (error) {
    console.error('转存失败:', error)
    $message.error('转存失败')
  } finally {
    resource.transferring = false
  }
}

// 批量转存
const handleBatchTransfer = async () => {
  if (selectedResources.value.length === 0) {
    $message.warning('请选择要转存的资源')
    return
  }

  batchTransferring.value = true
  transferResults.value = []

  try {
    const selectedItems = resources.value.filter(r => selectedResources.value.includes(r.id))
    
    // 这里应该调用实际的批量转存API
    // 由于只是UI展示，模拟批量转存过程
    for (const item of selectedItems) {
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      const isSuccess = Math.random() > 0.3
      transferResults.value.push({
        id: item.id,
        title: item.title,
        url: item.url,
        success: isSuccess,
        message: isSuccess ? '转存成功' : '转存失败：网络错误'
      })
    }

    showTransferResult.value = true
    
    // 刷新列表
    refreshData()
    
  } catch (error) {
    console.error('批量转存失败:', error)
    $message.error('批量转存失败')
  } finally {
    batchTransferring.value = false
  }
}

// 查看资源详情
const viewResource = (resource: any) => {
  console.log('查看资源详情:', resource)
  // 这里可以打开资源详情模态框
}

// 获取状态类型
const getStatusType = (resource: any) => {
  if (resource.is_valid === false) return 'error'
  if (resource.is_valid === true) return 'success'
  return 'warning'
}

// 获取状态文本
const getStatusText = (resource: any) => {
  if (resource.is_valid === false) return '无效'
  if (resource.is_valid === true) return '有效'
  return '待验证'
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

// 初始化
onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchUntransferredResources()
})
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>