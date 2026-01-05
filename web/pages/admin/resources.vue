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
          :item-size="100"
          class="h-full"
        >
          <template #default="{ item: resource }">
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors mb-4">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center space-x-2 mb-2">
                    <n-checkbox
                      :value="resource.id"
                      :checked="selectedResources.includes(resource.id)"
                      @update:checked="(checked) => toggleResourceSelection(resource.id, checked)"
                    />
                    <span class="text-sm text-gray-500 dark:text-gray-400">{{ resource.id }}</span>

                    <span v-if="resource.pan_id" class="text-xs px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded flex-shrink-0">
                      {{ getPlatformName(resource.pan_id) }}
                    </span>
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white flex-1 line-clamp-1">
                      {{ resource.title }}
                    </h3>
                    <span v-if="resource.category_id" class="text-xs px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded flex-shrink-0">
                      {{ getCategoryName(resource.category_id) }}
                    </span>

                  </div>

                  <p v-if="resource.description" class="text-gray-600 dark:text-gray-400 mb-2 line-clamp-2">
                    {{ resource.description }}
                  </p>

                  <div class="flex items-center space-x-4 text-sm text-gray-500 dark:text-gray-400">
                    <span>
                      <i class="fas fa-link mr-1"></i>
                      {{ resource.url }}
                    </span>
                    <span v-if="resource.author">
                      <i class="fas fa-user mr-1"></i>
                      {{ resource.author }}
                    </span>
                    <span v-if="resource.file_size">
                      <i class="fas fa-file mr-1"></i>
                      {{ resource.file_size }}
                    </span>
                    <span>
                      <i class="fas fa-eye mr-1"></i>
                      {{ resource.view_count || 0 }}
                    </span>
                    <span>
                      <i class="fas fa-clock mr-1"></i>
                      {{ resource.updated_at }}
                    </span>
                  </div>

                  <div v-if="resource.tags && resource.tags.length > 0" class="mt-2">
                    <div class="flex flex-wrap gap-1">
                      <span
                        v-for="tag in resource.tags"
                        :key="tag.id"
                        class="text-xs px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="flex items-center space-x-2 ml-4">
                  <n-button size="small" type="primary" @click="editResource(resource)">
                    <template #icon>
                      <i class="fas fa-edit"></i>
                    </template>
                    编辑
                  </n-button>
                  <n-button size="small" type="error" @click="deleteResource(resource)">
                    <template #icon>
                      <i class="fas fa-trash"></i>
                    </template>
                    删除
                  </n-button>
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

  <!-- 编辑资源模态框 -->
  <n-modal v-model:show="showEditModal" preset="card" title="编辑资源" style="width: 700px; max-height: 80vh">
    <n-scrollbar style="max-height: 60vh">
      <n-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-placement="left"
        label-width="80px"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="标题" path="title">
          <n-input v-model:value="editForm.title" placeholder="请输入资源标题" />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="editForm.description"
            type="textarea"
            placeholder="请输入资源描述"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="URL" path="url">
          <n-input v-model:value="editForm.url" placeholder="请输入资源链接" />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="editForm.category_id"
            :options="categoryOptions"
            placeholder="请选择分类"
            clearable
          />
        </n-form-item>

        <n-form-item label="平台" path="pan_id">
          <n-select
            v-model:value="editForm.pan_id"
            :options="platformOptions"
            placeholder="请选择平台"
            clearable
          />
        </n-form-item>

        <n-form-item label="标签" path="tag_ids">
          <n-select
            v-model:value="editForm.tag_ids"
            :options="tagOptions"
            :loading="tagLoading"
            :filterable="true"
            :remote="true"
            :clearable="true"
            placeholder="请选择标签，支持搜索"
            multiple
            @search="handleTagSearch"
            @scroll="handleTagScroll"
          />
        </n-form-item>

        <n-form-item label="作者" path="author">
          <n-input v-model:value="editForm.author" placeholder="请输入作者" />
        </n-form-item>

        <n-form-item label="文件大小" path="file_size">
          <n-input v-model:value="editForm.file_size" placeholder="如：2.5GB" />
        </n-form-item>

        <n-form-item label="封面图片" path="cover">
          <n-input v-model:value="editForm.cover" placeholder="请输入封面图片URL" />
        </n-form-item>

        <n-form-item label="转存链接" path="save_url">
          <n-input v-model:value="editForm.save_url" placeholder="请输入转存链接" />
        </n-form-item>

        <n-form-item label="是否有效" path="is_valid">
          <n-switch v-model:value="editForm.is_valid" />
        </n-form-item>

        <n-form-item label="是否公开" path="is_public">
          <n-switch v-model:value="editForm.is_public" />
        </n-form-item>
      </n-form>
    </n-scrollbar>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <n-button @click="showEditModal = false">取消</n-button>
        <n-button type="primary" @click="handleEditSubmit" :loading="editing">
          保存
        </n-button>
      </div>
    </template>
  </n-modal>
 

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
const showEditModal = ref(false)
const editing = ref(false)
const editingResource = ref<Resource | null>(null)
const editFormRef = ref()

// 编辑表单
const editForm = ref({
  title: '',
  description: '',
  url: '',
  category_id: null as number | null,
  pan_id: null as number | null,
  tag_ids: [] as number[],
  author: '',
  file_size: '',
  cover: '',
  save_url: '',
  is_valid: true,
  is_public: true
})

// 编辑验证规则
const editRules = {
  title: {
    required: true,
    message: '请输入资源标题',
    trigger: 'blur'
  },
  url: {
    required: true,
    message: '请输入资源链接',
    trigger: 'blur'
  }
}

// 获取资源API
import { useResourceApi, useCategoryApi, useTagApi, usePanApi } from '~/composables/useApi'
import { useMessage } from 'naive-ui'

// 用户状态管理
const userStore = useUserStore()
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()
const message = useMessage()

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
  const platforms = data?.data || data || []
  
  // 确保platforms是数组
  if (!Array.isArray(platforms)) {
    console.warn('platformOptions: platforms不是数组', platforms)
    return []
  }
  
  return platforms.map((platform: any) => ({
    label: platform.remark || platform.name,
    value: platform.id
  }))
})

// 获取分类名称
const getCategoryName = (categoryId: number) => {
  const category = (categoriesData.value as any)?.data?.find((cat: any) => cat.id === categoryId)
  return category?.name || '未知分类'
}

// 获取平台名称
const getPlatformName = (platformId: number) => {
  // console.log('platformId', platformId, platformsData.value)
  const platform = (platformsData.value as any)?.find((plat: any) => plat.id === platformId)
  return platform?.remark || platform?.name || '未知平台'
}

// 加载标签选项（支持搜索和分页）
const loadTagOptions = async (search = '', page = 1, pageSize = 20) => {
  tagLoading.value = true
  try {
    const params = {
      page,
      page_size: pageSize,
      search
    }

    const response = await tagApi.getTags(params)
    const data = response?.items || response?.data || []
    const total = response?.total || 0

    // 确保tags是数组
    if (!Array.isArray(data)) {
      console.warn('loadTagOptions: tags不是数组', data)
      return { options: [], total: 0 }
    }

    const options = data.map((tag: any) => ({
      label: tag.name,
      value: tag.id
    }))

    tagPagination.total = total

    return { options, total }
  } catch (error) {
    console.error('加载标签选项失败:', error)
    message.error('加载标签失败')
    return { options: [], total: 0 }
  } finally {
    tagLoading.value = false
  }
}

// 标签搜索处理函数
const handleTagSearch = async (query = '') => {
  tagSearchKeyword.value = query
  const { options } = await loadTagOptions(query, 1, tagPagination.pageSize)
  tagOptions.value = options
}

// 标签滚动加载更多函数
const handleTagScroll = async (e: any) => {
  const { scrollTop, scrollHeight, clientHeight } = e

  // 检查是否滚动到底部
  if (scrollTop + clientHeight >= scrollHeight - 10 &&
      tagOptions.value.length < tagPagination.total &&
      !tagLoading.value) {

    tagPagination.page++
    const { options } = await loadTagOptions(
      tagSearchKeyword.value,
      tagPagination.page,
      tagPagination.pageSize
    )

    // 合并选项，避免重复
    const existingIds = new Set(tagOptions.value.map((opt: any) => opt.value))
    const newOptions = options.filter((opt: any) => !existingIds.has(opt.value))
    tagOptions.value = [...tagOptions.value, ...newOptions]
  }
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    }
    
    // 添加分类筛选
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
      // console.log('添加分类筛选:', selectedCategory.value)
    }
    
    // 添加平台筛选
    if (selectedPlatform.value) {
      params.pan_id = selectedPlatform.value
      // console.log('添加平台筛选:', selectedPlatform.value)
    }
    
    const response = await resourceApi.getResources(params) as any
    
    if (response && response.data) {
      // 处理嵌套的data结构：{data: {data: [...], total: ...}}
      if (response.data.data && Array.isArray(response.data.data)) {
        resources.value = response.data.data
        total.value = response.data.total || 0
      } else {
        // 处理直接的data结构：{data: [...], total: ...}
        resources.value = response.data
        total.value = response.total || 0
      }
      // 清空选择（因为数据已更新）
      selectedResources.value = []
    } else {
      resources.value = []
      total.value = 0
      selectedResources.value = []
    }
  } catch (error) {
    console.error('获取资源失败:', error)
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
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
  // 清空选择
  selectedResources.value = []
  // 重新获取数据
  fetchData()
}

// 切换资源选择
const toggleResourceSelection = (resourceId: number, checked: boolean) => {
  if (checked) {
    selectedResources.value.push(resourceId)
  } else {
    const index = selectedResources.value.indexOf(resourceId)
    if (index > -1) {
      selectedResources.value.splice(index, 1)
    }
  }
}

// 全选状态计算
const isAllSelected = computed(() => {
  return resources.value.length > 0 && selectedResources.value.length === resources.value.length
})

// 部分选中状态计算
const isIndeterminate = computed(() => {
  return selectedResources.value.length > 0 && selectedResources.value.length < resources.value.length
})

// 切换全选
const toggleSelectAll = (checked: boolean) => {
  if (checked) {
    // 全选：添加所有当前页面的资源ID
    selectedResources.value = resources.value.map(resource => resource.id)
  } else {
    // 取消全选：清空选择
    selectedResources.value = []
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
const editResource = async (resource: Resource) => {
  editingResource.value = resource

  // 从资源的tags数组中提取tag_ids，确保tags是数组
  let tagIds: number[] = []
  if (resource.tags && Array.isArray(resource.tags)) {
    tagIds = resource.tags.map(tag => tag.id)
  } else if (resource.tag_ids && Array.isArray(resource.tag_ids)) {
    // 如果tags不存在但tag_ids存在，直接使用tag_ids
    tagIds = resource.tag_ids
  }

  // 如果存在已选择的标签ID，确保这些标签在选项中显示
  if (tagIds && tagIds.length > 0) {
    // 获取已选择的标签详情，以便显示标签名称
    const selectedTags = []
    for (const tagId of tagIds) {
      // 检查标签是否已在当前选项中
      const existingTag = tagOptions.value.find((opt: any) => opt.value === tagId)
      if (existingTag) {
        selectedTags.push(existingTag)
      } else {
        // 如果标签不在当前选项中，单独获取其信息
        try {
          const tagDetail = await tagApi.getTag(tagId)
          if (tagDetail) {
            selectedTags.push({
              label: tagDetail.name,
              value: tagDetail.id
            })
          }
        } catch (error) {
          console.error('获取标签详情失败:', error)
          // 作为回退，添加一个临时标签
          selectedTags.push({
            label: `标签${tagId}`,
            value: tagId
          })
        }
      }
    }

    // 确保这些标签出现在选项中
    const newOptions = [...tagOptions.value]
    for (const selectedTag of selectedTags) {
      const exists = newOptions.some((opt: any) => opt.value === selectedTag.value)
      if (!exists) {
        newOptions.push(selectedTag)
      }
    }
    tagOptions.value = newOptions
  }

  editForm.value = {
    title: resource.title,
    description: resource.description || '',
    url: resource.url,
    category_id: resource.category_id || null,
    pan_id: resource.pan_id || null,
    tag_ids: tagIds,
    author: resource.author || '',
    file_size: resource.file_size || '',
    cover: resource.cover || '',
    save_url: resource.save_url || '',
    is_valid: resource.is_valid !== undefined ? resource.is_valid : true,
    is_public: resource.is_public !== undefined ? resource.is_public : true
  }
  showEditModal.value = true
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
        content: '删除成功',
        duration: 3000
      })

      // 从当前列表中移除
      const index = resources.value.findIndex(r => r.id === resource.id)
      if (index > -1) {
        resources.value.splice(index, 1)
        total.value = Math.max(0, total.value - 1)
      }
      // 从选择列表中移除
      const selectedIndex = selectedResources.value.indexOf(resource.id)
      if (selectedIndex > -1) {
        selectedResources.value.splice(selectedIndex, 1)
      }
    } catch (error) {
      console.error('删除失败:', error)
      notification.value?.error({
        content: '删除失败: ' + (error as Error).message,
        duration: 3000
      })
    }
  }
}

// 批量删除
const batchDelete = async () => {
  if (selectedResources.value.length === 0) {
    notification.value?.warning({
      content: '请先选择要删除的资源',
      duration: 3000
    })
    return
  }

  const deleteCount = selectedResources.value.length

  dialog.value?.warning({
    title: '警告',
    content: `确定要删除选中的 ${deleteCount} 个资源吗？此操作不可恢复！`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        // 调用批量删除API
        await resourceApi.batchDeleteResources(selectedResources.value)
        notification.value?.success({
          content: `成功删除 ${deleteCount} 个资源`,
          duration: 3000
        })
        // 清空选择
        selectedResources.value = []
        showBatchModal.value = false
        // 重新获取数据
        await fetchData()
      } catch (error) {
        console.error('批量删除失败:', error)
        notification.value?.error({
          content: '批量删除失败: ' + (error as Error).message,
          duration: 3000
        })
      }
    }
  })
}

// 批量更新
const batchUpdate = () => {
  if (selectedResources.value.length === 0) {
    notification.value?.warning({
      content: '请先选择要更新的资源',
      duration: 3000
    })
    return
  }
  
  // 这里可以实现批量更新功能
  console.log('批量更新:', selectedResources.value)
  notification.value?.info({
    content: '批量更新功能开发中',
    duration: 3000
  })
}

// 提交编辑
const handleEditSubmit = async () => {
  try {
    editing.value = true
    await editFormRef.value?.validate()

    await resourceApi.updateResource(editingResource.value!.id, editForm.value)

    notification.value?.success({
      content: '更新成功',
      duration: 3000
    })

    // 关闭模态框
    showEditModal.value = false
    editingResource.value = null
    
    // 重新获取数据以确保tags等关联数据正确显示
    await fetchData()

  } catch (error) {
    console.error('更新失败:', error)
    notification.value?.error({
      content: '更新失败',
      duration: 3000
    })
  } finally {
    editing.value = false
  }
}

// 页面加载时获取数据
onMounted(async () => {
  // 初始化用户认证状态
  const userStore = useUserStore()
  userStore.initAuth()

  // 初始化加载第一页标签
  const { options } = await loadTagOptions('', 1, tagPagination.pageSize)
  tagOptions.value = options

  fetchData()
})


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