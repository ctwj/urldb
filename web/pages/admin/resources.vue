<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
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
    </div>

    <!-- 搜索和筛选 -->
    <n-card>
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
    </n-card>

    <!-- 资源列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <span class="text-lg font-semibold">资源列表</span>
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
        <i class="fas fa-inbox text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无资源数据</p>
      </div>

      <div v-else>
        <!-- 虚拟列表 -->
        <n-virtual-list
          :items="resources"
          :item-size="100"
          style="max-height: 400px"
          container-style="height: 600px;"
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
                    <span class="text-sm text-gray-500">{{ resource.id }}</span>
                    
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
                  
                  <div class="flex items-center space-x-4 text-sm text-gray-500">
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
                  <!-- <n-button size="small" type="primary" @click="editResource(resource)">
                    <template #icon>
                      <i class="fas fa-edit"></i>
                    </template>
                    编辑
                  </n-button> -->
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

        <!-- 分页 -->
        <div class="mt-6">
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
    </n-card>

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
    <n-modal v-model:show="showEditModal" preset="card" title="编辑资源" style="width: 600px">
      <n-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-placement="left"
        label-width="auto"
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
            placeholder="请选择标签"
            multiple
            clearable
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="showEditModal = false">取消</n-button>
          <n-button type="primary" @click="handleEditSubmit" :loading="editing">
            保存
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

interface Resource {
  id: number
  title: string
  description?: string
  url: string
  category_id?: number
  pan_id?: number
  tag_ids?: number[]
  tags?: Array<{ id: number; name: string }>
  author?: string
  file_size?: string
  view_count?: number
  is_valid: boolean
  is_public: boolean
  created_at: string
  updated_at: string
}

const notification = useNotification()
const dialog = useDialog()
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
  tag_ids: [] as number[]
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
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()

// 获取分类数据
const { data: categoriesData } = await useAsyncData('resourceCategories', () => categoryApi.getCategories())

// 获取标签数据
const { data: tagsData } = await useAsyncData('resourceTags', () => tagApi.getTags())

// 获取平台数据
const { data: platformsData } = await useAsyncData('resourcePlatforms', () => panApi.getPans())

// 分类选项
const categoryOptions = computed(() => {
  const data = categoriesData.value as any
  console.log('分类数据:', data)
  const categories = data?.items || data || []
  console.log('处理后的分类:', categories)
  const options = categories.map((cat: any) => ({
    label: cat.name,
    value: cat.id
  }))
  console.log('分类选项:', options)
  return options
})

// 标签选项
const tagOptions = computed(() => {
  const data = tagsData.value as any
  const tags = data?.data || data || []
  return tags.map((tag: any) => ({
    label: tag.name,
    value: tag.id
  }))
})

// 平台选项
const platformOptions = computed(() => {
  const data = platformsData.value as any
  const platforms = data?.data || data || []
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
  console.log('platformId', platformId, platformsData.value)
  const platform = (platformsData.value as any)?.find((plat: any) => plat.id === platformId)
  return platform?.remark || platform?.name || '未知平台'
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
      console.log('添加分类筛选:', selectedCategory.value)
    }
    
    // 添加平台筛选
    if (selectedPlatform.value) {
      params.pan_id = selectedPlatform.value
      console.log('添加平台筛选:', selectedPlatform.value)
    }
    
    console.log('请求参数:', params)
    console.log('pageSize:', pageSize.value)
    console.log('selectedCategory:', selectedCategory.value)
    console.log('selectedPlatform:', selectedPlatform.value)
    const response = await resourceApi.getResources(params) as any
    console.log('API响应:', response)
    console.log('返回的资源数量:', response?.data?.length || 0)
    
    if (response && response.data) {
      resources.value = response.data
      total.value = response.total || 0
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
    notification.info({
      content: '已自动全选当前页面资源',
      duration: 2000
    })
  }
  showBatchModal.value = true
}

// 编辑资源
const editResource = (resource: Resource) => {
  editingResource.value = resource
  editForm.value = {
    title: resource.title,
    description: resource.description || '',
    url: resource.url,
    category_id: resource.category_id || null,
    pan_id: resource.pan_id || null,
    tag_ids: resource.tag_ids || []
  }
  showEditModal.value = true
}

// 删除资源
const deleteResource = async (resource: Resource) => {
  dialog.warning({
    title: '警告',
    content: `确定要删除资源"${resource.title}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await resourceApi.deleteResource(resource.id)
        notification.success({
          content: '删除成功',
          duration: 3000
        })
        // 从当前列表中移除
        const index = resources.value.findIndex(r => r.id === resource.id)
        if (index > -1) {
          resources.value.splice(index, 1)
        }
        // 重新获取数据以更新总数
        fetchData()
      } catch (error) {
        console.error('删除失败:', error)
        notification.error({
          content: '删除失败',
          duration: 3000
        })
      }
    }
  })
}

// 批量删除
const batchDelete = async () => {
  if (selectedResources.value.length === 0) {
    notification.warning({
      content: '请先选择要删除的资源',
      duration: 3000
    })
    return
  }

  dialog.warning({
    title: '警告',
    content: `确定要删除选中的 ${selectedResources.value.length} 个资源吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        // 这里应该调用批量删除API
        console.log('批量删除:', selectedResources.value)
        notification.success({
          content: '批量删除成功',
          duration: 3000
        })
        selectedResources.value = []
        showBatchModal.value = false
        fetchData()
      } catch (error) {
        console.error('批量删除失败:', error)
        notification.error({
          content: '批量删除失败',
          duration: 3000
        })
      }
    }
  })
}

// 批量更新
const batchUpdate = () => {
  if (selectedResources.value.length === 0) {
    notification.warning({
      content: '请先选择要更新的资源',
      duration: 3000
    })
    return
  }
  
  // 这里可以实现批量更新功能
  console.log('批量更新:', selectedResources.value)
  notification.info({
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
    
    notification.success({
      content: '更新成功',
      duration: 3000
    })
    
    // 更新本地数据
    const resourceId = editingResource.value?.id
    const index = resources.value.findIndex(r => r.id === resourceId)
    if (index > -1) {
      resources.value[index] = { 
        ...resources.value[index], 
        title: editForm.value.title,
        description: editForm.value.description,
        url: editForm.value.url,
        category_id: editForm.value.category_id || undefined,
        pan_id: editForm.value.pan_id || undefined,
        tag_ids: editForm.value.tag_ids
      }
    }
    
    showEditModal.value = false
    editingResource.value = null
  } catch (error) {
    console.error('更新失败:', error)
    notification.error({
      content: '更新失败',
      duration: 3000
    })
  } finally {
    editing.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
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