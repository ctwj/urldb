<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">资源管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的所有资源</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="navigateTo('/admin/add-resource')">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加资源
        </n-button>
        <n-button @click="openBatchModal" type="info">
          <template #icon>
            <i class="fas fa-list"></i>
          </template>
          批量操作
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和筛选 -->
    <template #filter-bar>
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索资源..."
            @keyup.enter="handleSearch"
            clearable
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
            v-model:value="selectedPlatform"
            placeholder="选择平台"
            :options="platformOptions"
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

    <!-- 内容区header - 资源列表头部 -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <span class="text-lg font-semibold">资源列表</span>
          <div class="flex items-center space-x-2">
            <n-checkbox
              :checked="isAllSelected"
              @update:checked="toggleSelectAll"
              :indeterminate="isIndeterminate"
            />
            <span class="text-sm text-gray-500 dark:text-gray-400">全选</span>
          </div>
        </div>
        <span class="text-sm text-gray-500 dark:text-gray-400">共 {{ total }} 个资源，已选择 {{ selectedResources.length }} 个</span>
      </div>
    </template>

    <!-- 内容区content - 资源列表 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="resources.length === 0" class="flex flex-col items-center justify-center py-12">
        <i class="fas fa-inbox text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500 dark:text-gray-400">暂无资源数据</p>
      </div>

      <!-- 虚拟列表容器 -->
      <div v-else class="flex-1 h-full overflow-hidden">
        <n-virtual-list
          :items="resources"
          :item-size="120"
          :key-field="'id'"
          class="h-full"
        >
          <template #default="{ item: resource }">
            <div class="border m-2 border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors mb-2">
              <!-- 头部：标题和操作按钮 -->
              <div class="flex items-center justify-between p-3 border-b border-gray-200 dark:border-gray-700">
                <!-- 左侧：标题和元信息 -->
                <div class="flex items-center space-x-2 flex-1 min-w-0">
                  <n-checkbox
                    :value="resource.id"
                    :checked="selectedResources.includes(resource.id)"
                    @update:checked="(checked) => toggleResourceSelection(resource.id, checked)"
                  />
                  <span class="text-xs text-gray-500 dark:text-gray-400 font-mono flex-shrink-0">#{{ resource.id }}</span>

                  <!-- 转存标记 - 头部显示 -->
                  <div v-if="resource.save_url" class="flex items-center space-x-1 flex-shrink-0 bg-green-50 dark:bg-green-900/30 px-2 py-1 rounded-md">
                    <i class="fas fa-save text-green-600 dark:text-green-400 text-xs"></i>
                    <span class="text-xs font-medium text-green-700 dark:text-green-300">已转存</span>
                  </div>

                  <span v-if="resource.pan_id" class="text-xs px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded flex-shrink-0">
                    {{ getPlatformName(resource.pan_id) }}
                  </span>

                  <h3 class="text-base font-medium text-gray-900 dark:text-white line-clamp-1 flex-1 min-w-0">
                    {{ resource.title }}
                  </h3>

                  <span v-if="resource.category_id" class="text-xs px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded flex-shrink-0">
                    {{ getCategoryName(resource.category_id) }}
                  </span>
                </div>

                <!-- 右侧：操作按钮 -->
                <div class="flex items-center space-x-1 ml-2 flex-shrink-0">
                  <n-button size="tiny" type="info" @click="openSmartManage(resource)">
                    <template #icon>
                      <i class="fas fa-robot text-xs"></i>
                    </template>
                  </n-button>
                  <n-button size="tiny" type="primary" @click="editResource(resource)">
                    <template #icon>
                      <i class="fas fa-edit text-xs"></i>
                    </template>
                  </n-button>
                  <n-button size="tiny" type="error" @click="deleteResource(resource)">
                    <template #icon>
                      <i class="fas fa-trash text-xs"></i>
                    </template>
                  </n-button>
                </div>
              </div>

              <!-- 主体：图片和内容 -->
              <div class="flex p-3">
                <!-- 左侧：预览图片 -->
                <div class="flex-shrink-0 mr-3">
                  <div class="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-md overflow-hidden flex items-center justify-center">
                    <img
                      v-if="resource.cover"
                      :src="resource.cover"
                      :alt="resource.title"
                      class="w-full h-full object-cover"
                      @error="handleImageError"
                    />
                    <div v-else class="text-gray-400 dark:text-gray-500 text-center">
                      <i class="fas fa-image text-lg"></i>
                    </div>
                  </div>
                </div>

                <!-- 右侧：描述和其他信息 -->
                <div class="flex-1 min-w-0">
                  <!-- 描述 -->
                  <p v-if="resource.description" class="text-gray-600 dark:text-gray-400 text-sm mb-2 line-clamp-2">
                    {{ resource.description }}
                  </p>

                  <!-- 元信息行 -->
                  <div class="flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400 mb-2">
                    <!-- 转存标记 -->
                    <div v-if="resource.save_url" class="flex items-center space-x-1 flex-shrink-0">
                      <i class="fas fa-save text-green-600 dark:text-green-400"></i>
                      <span class="text-green-600 dark:text-green-400 font-medium">已转存</span>
                    </div>
                    <!-- 访问次数 -->
                    <div class="flex items-center space-x-1 flex-shrink-0">
                      <i class="fas fa-eye text-blue-600 dark:text-blue-400"></i>
                      <span class="text-blue-600 dark:text-blue-400 font-medium">{{ resource.view_count || 0 }}</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <i class="fas fa-user flex-shrink-0"></i>
                      <span class="truncate">{{ resource.author || '未知作者' }}</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <i class="fas fa-file flex-shrink-0"></i>
                      <span class="truncate">{{ resource.file_size || '未知大小' }}</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <i class="fas fa-link flex-shrink-0"></i>
                      <span class="truncate max-w-[150px]">{{ resource.url }}</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <i class="fas fa-clock flex-shrink-0"></i>
                      <span>{{ formatDate(resource.updated_at) }}</span>
                    </div>
                  </div>

                  <!-- 标签 -->
                  <div v-if="resource.tags && resource.tags.length > 0" class="flex flex-wrap gap-1">
                    <span
                      v-for="tag in resource.tags.slice(0, 6)"
                      :key="tag.id"
                      class="text-xs px-1.5 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                    >
                      {{ tag.name }}
                    </span>
                    <span
                      v-if="resource.tags.length > 6"
                      class="text-xs px-1.5 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 rounded"
                    >
                      +{{ resource.tags.length - 6 }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </n-virtual-list>
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
            :page-sizes="[100, 200, 500, 1000]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </template>
  </AdminPageLayout>

  <!-- 模态框 - 在AdminPageLayout外部 -->
  <!-- 批量操作模态框 -->
  <n-modal v-model:show="showBatchModal" preset="card" title="批量操作" style="width: 600px">
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <span class="font-medium">已选择 {{ selectedResources.length }} 个资源</span>
          <p class="text-sm text-gray-500 mt-1">
            {{ isAllSelected ? '已全选当前页面' : isIndeterminate ? '部分选中' : '未选择' }}
          </p>
        </div>
        <n-button size="small" @click="clearSelection">清空选择</n-button>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <n-button type="error" @click="batchDelete" :disabled="selectedResources.length === 0">
          <template #icon>
            <i class="fas fa-trash"></i>
          </template>
          批量删除
        </n-button>
        <n-button type="warning" @click="batchUpdate" :disabled="selectedResources.length === 0">
          <template #icon>
            <i class="fas fa-edit"></i>
          </template>
          批量更新
        </n-button>
      </div>
    </div>
  </n-modal>

  <!-- 资源编辑抽屉 -->
  <ResourceEditDrawer
    v-model:show="showEditDrawer"
    :resource="editingResource"
    @updated="handleResourceUpdated"
  />

  <!-- 智能管理抽屉 -->
  <ResourceSmartManageDrawer
    v-model:show="showSmartManageDrawer"
    :resource="smartManageResource"
    @updated="handleResourceUpdated"
  />
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})
interface Resource {
  id: number
  title: string
  description?: string
  url: string
  category_id?: number
  pan_id?: number
  tag_ids?: number[]
  tags?: Array<{ id: number; name: string; description?: string }>
  author?: string
  file_size?: string
  view_count?: number
  cover?: string
  save_url?: string
  is_valid: boolean
  is_public: boolean
  created_at: string
  updated_at: string
}

// 使用computed延迟获取notification和dialog实例，避免SSR问题
const notification = computed(() => {
  if (process.client) {
    return useNotification()
  }
  return null
})

const dialog = computed(() => {
  if (process.client) {
    return useDialog()
  }
  return null
})

const resources = ref<Resource[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(200)
const searchQuery = ref('')
const selectedCategory = ref(null)
const selectedPlatform = ref(null)
const selectedResources = ref<number[]>([])
const showBatchModal = ref(false)
const showEditDrawer = ref(false)
const editingResource = ref<Resource | null>(null)

// 获取资源API
import { useResourceApi, useCategoryApi, useTagApi, usePanApi, useAIApi } from '~/composables/useApi'
import { useMessage } from 'naive-ui'
import ResourceSmartManageDrawer from '~/components/Admin/ResourceSmartManageDrawer.vue'
import ResourceEditDrawer from '~/components/Admin/ResourceEditDrawer.vue'

// 用户状态管理
const userStore = useUserStore()
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()
const aiApi = useAIApi()
const message = useMessage()

// 智能管理抽屉状态
const showSmartManageDrawer = ref(false)
const smartManageResource = ref<Resource | null>(null)

// 获取分类数据
const { data: categoriesData } = await useAsyncData('resourceCategories', () => categoryApi.getCategories())

// 标签搜索和加载相关状态
const tagSearchKeyword = ref('')
const tagLoading = ref(false)
const tagOptions = ref([])
const tagPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取平台数据
const { data: platformsData } = await useAsyncData('resourcePlatforms', () => panApi.getPans())

// 分类选项
const categoryOptions = computed(() => {
  const data = categoriesData.value as any
  const categories = data?.items || data || []

  // 确保categories是数组
  if (!Array.isArray(categories)) {
    console.warn('categoryOptions: categories不是数组', categories)
    return []
  }

  return categories.map((cat: any) => ({
    label: cat.name,
    value: cat.id
  }))
})

// 平台选项
const platformOptions = computed(() => {
  const data = platformsData.value as any
  const platforms = data?.items || data || []

  // 确保platforms是数组
  if (!Array.isArray(platforms)) {
    console.warn('platformOptions: platforms不是数组', platforms)
    return []
  }

  return platforms.map((pan: any) => ({
    label: pan.name,
    value: pan.id
  }))
})

// 选择所有资源
const isAllSelected = computed(() => {
  return resources.value.length > 0 && selectedResources.value.length === resources.value.length
})

// 选择状态（部分选中）
const isIndeterminate = computed(() => {
  return selectedResources.value.length > 0 && !isAllSelected.value
})

// 切换选择所有资源
const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedResources.value = []
  } else {
    selectedResources.value = resources.value.map(resource => resource.id)
  }
}

// 切换单个资源选择状态
const toggleResourceSelection = (resourceId: number, checked: boolean) => {
  if (checked) {
    if (!selectedResources.value.includes(resourceId)) {
      selectedResources.value.push(resourceId)
    }
  } else {
    const index = selectedResources.value.indexOf(resourceId)
    if (index > -1) {
      selectedResources.value.splice(index, 1)
    }
  }
}

// 清空选择
const clearSelection = () => {
  selectedResources.value = []
}

// 打开批量操作模态框
const openBatchModal = () => {
  // 如果没有选择任何资源，自动全选当前页面
  if (selectedResources.value.length === 0 && resources.value.length > 0) {
    selectedResources.value = resources.value.map(resource => resource.id)
    notification.value?.info({
      content: '已自动全选当前页面资源',
      duration: 2000
    })
  }
  showBatchModal.value = true
}

// 编辑资源
const editResource = (resource: Resource) => {
  editingResource.value = resource
  showEditDrawer.value = true
}

// 打开智能管理抽屉
const openSmartManage = (resource: Resource) => {
  smartManageResource.value = resource
  showSmartManageDrawer.value = true
}

// 删除资源
const deleteResource = async (resource: Resource) => {
  console.log('删除资源被点击:', resource.title, 'ID:', resource.id)

  // 使用原生确认对话框
  if (confirm(`确定要删除资源"${resource.title}"吗？`)) {
    try {
      console.log('开始删除资源:', resource.id)
      await resourceApi.deleteResource(resource.id)
      console.log('删除成功')

      notification.value?.success({
        content: '资源删除成功',
        duration: 2000
      })

      // 从本地列表中移除
      const index = resources.value.findIndex(r => r.id === resource.id)
      if (index > -1) {
        resources.value.splice(index, 1)
        total.value--
      }
    } catch (error) {
      console.error('删除资源失败:', error)
      notification.value?.error({
        content: '资源删除失败',
        duration: 2000
      })
    }
  }
}

// 批量删除资源
const batchDelete = async () => {
  if (selectedResources.value.length === 0) {
    message.warning('请先选择要删除的资源')
    return
  }

  const confirmResult = confirm(`确定要删除选中的 ${selectedResources.value.length} 个资源吗？`)
  if (!confirmResult) return

  try {
    await resourceApi.batchDeleteResources(selectedResources.value)

    notification.value?.success({
      content: `成功删除 ${selectedResources.value.length} 个资源`,
      duration: 2000
    })

    // 重新加载数据
    await fetchData()
    selectedResources.value = []
    showBatchModal.value = false
  } catch (error) {
    console.error('批量删除失败:', error)
    notification.value?.error({
      content: '批量删除失败',
      duration: 2000
    })
  }
}

// 获取资源数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await resourceApi.getResources({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
      category_id: selectedCategory.value,
      pan_id: selectedPlatform.value
    })

    if (response) {
      // 根据API响应格式，response已经是解析后的数据（数组或分页对象）
      if (Array.isArray(response)) {
        // 如果是数组，直接赋值（适用于不需要分页的列表）
        resources.value = response
        total.value = response.length
      } else if (response && typeof response === 'object') {
        // 如果是分页对象，处理items和total
        resources.value = response.items || response.resources || response.data || []
        total.value = response.total || response.count || response.total_count || 0
      } else {
        // 如果是其他格式，初始化为空数组
        resources.value = []
        total.value = 0
      }
    }
  } catch (error) {
    console.error('获取资源失败:', error)
    notification.value?.error({
      content: '获取资源失败',
      duration: 2000
    })
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = async () => {
  currentPage.value = 1
  await fetchData()
}

// 处理分页变化
const handlePageChange = async (page: number) => {
  currentPage.value = page
  await fetchData()
}

// 处理每页大小变化
const handlePageSizeChange = async (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  await fetchData()
}

// 刷新数据
const refreshData = async () => {
  await fetchData()
}

// 刷新标签选项
const loadTagOptions = async (keyword: string, page: number, pageSize: number) => {
  try {
    const response = await tagApi.getTags({
      search: keyword,
      page: page,
      page_size: pageSize
    })

    if (response) {
      const items = response.items || response.data || []
      const total = response.total || response.count || 0

      tagPagination.total = total

      const options = items.map((tag: any) => ({
        label: tag.name + (tag.description ? ` (${tag.description})` : ''),
        value: tag.id
      }))

      return { options, total }
    }
  } catch (error) {
    console.error('加载标签失败:', error)
  }
  return { options: [], total: 0 }
}

// 处理标签搜索
const handleTagSearch = async (keyword: string) => {
  tagLoading.value = true
  try {
    const { options } = await loadTagOptions(keyword, 1, tagPagination.pageSize)
    tagOptions.value = options
  } catch (error) {
    console.error('搜索标签失败:', error)
  } finally {
    tagLoading.value = false
  }
}

// 处理标签滚动加载
const handleTagScroll = async () => {
  if (tagOptions.value.length >= tagPagination.total) {
    return // 已加载全部数据
  }

  const nextPage = Math.floor(tagOptions.value.length / tagPagination.pageSize) + 1
  if (nextPage > 1) {
    tagLoading.value = true
    try {
      const { options } = await loadTagOptions(tagSearchKeyword.value, nextPage, tagPagination.pageSize)
      tagOptions.value = [...tagOptions.value, ...options]
    } catch (error) {
      console.error('加载更多标签失败:', error)
    } finally {
      tagLoading.value = false
    }
  }
}

// 获取平台名称
const getPlatformName = (platformId?: number) => {
  if (!platformId || !platformOptions.value) return null
  const platform = platformOptions.value.find(opt => opt.value === platformId)
  return platform?.label
}

// 获取分类名称
const getCategoryName = (categoryId?: number) => {
  if (!categoryId || !categoryOptions.value) return null
  const category = categoryOptions.value.find(opt => opt.value === categoryId)
  return category?.label
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知时间'
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateString
  }
}

// 处理图片加载错误
const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
  const parent = img.parentElement
  if (parent) {
    parent.innerHTML = '<div class="text-gray-400 dark:text-gray-500 text-center p-2"><i class="fas fa-image text-2xl"></i></div>'
  }
}




// 批量更新处理
const batchUpdate = async () => {
  // 这里可以实现批量更新逻辑
  notification.value?.info({
    content: '批量更新功能开发中',
    duration: 2000
  })
}

onMounted(async () => {
  // 初始化用户认证状态
  const userStore = useUserStore()
  userStore.initAuth()

  // 初始化加载第一页标签
  const { options } = await loadTagOptions('', 1, tagPagination.pageSize)
  tagOptions.value = options

  fetchData()
})

// 处理资源更新事件
const handleResourceUpdated = (updatedResource: Resource) => {
  // 更新资源列表中的对应资源
  const index = resources.value.findIndex(r => r.id === updatedResource.id)
  if (index !== -1) {
    resources.value[index] = updatedResource
  }

  message.success('资源已更新')
}


</script>

<style scoped>
/* 自定义样式 */
.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}

.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
</style>