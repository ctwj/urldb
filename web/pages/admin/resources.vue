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
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </template>

    <!-- 统一筛选栏 -->
    <template #filter-bar>
      <AdminFilterBar
        :config="filterConfig"
        v-model="filterValues"
        @search="handleSearch"
      />
    </template>

    <!-- 内容区header - 批量操作栏 + 全选 -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-4">
            <span class="text-lg font-semibold text-gray-900 dark:text-white">资源列表</span>
            <div class="flex items-center space-x-2">
              <n-checkbox
                :checked="isAllSelected"
                @update:checked="toggleSelectAll"
                :indeterminate="isIndeterminate"
              />
              <span class="text-sm text-gray-500 dark:text-gray-400">全选</span>
            </div>
          </div>
          <!-- 批量操作栏（选中时显示） -->
          <AdminBatchActionBar
            :actions="batchActions"
            :selected-ids="selectedResources"
            :total="total"
            @completed="handleBatchCompleted"
          />
        </div>
        <span class="text-sm text-gray-500 dark:text-gray-400">
          共 {{ total }} 个资源，已选择 {{ selectedResources.length }} 个
        </span>
      </div>
    </template>

    <!-- 内容区content - 资源列表 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large" />
      </div>

      <!-- 错误状态 -->
      <AdminErrorState
        v-else-if="error"
        icon="fas fa-exclamation-circle"
        :message="error?.message || '加载资源失败'"
        :on-retry="refreshData"
      />

      <!-- 空状态 -->
      <AdminEmptyState
        v-else-if="resources.length === 0"
        icon="fas fa-inbox"
        title="暂无资源数据"
        description="点击右上角「添加资源」创建第一个资源"
      />

      <!-- 资源列表容器（保留卡片式列表：信息密度高，适合资源管理场景） -->
      <div v-else class="h-full overflow-y-auto p-4">
        <div class="space-y-4">
          <div
            v-for="resource in resources"
            :key="resource.id"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
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
                  <!-- 转存清理状态标签（002-auto-cleanup-transfer） -->
                  <span
                    v-if="resource.cleaned_at"
                    class="text-xs px-2 py-1 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded flex-shrink-0"
                    :title="`清理时间：${resource.cleaned_at}`"
                  >
                    已清理
                  </span>
                  <span
                    v-else-if="resource.clean_error_msg"
                    class="text-xs px-2 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded flex-shrink-0"
                    :title="`清理失败：${resource.clean_error_msg}${resource.last_clean_error_at ? '（' + resource.last_clean_error_at + '）' : ''}`"
                  >
                    清理失败
                  </span>
                  <span
                    v-else-if="resource.transferred_at"
                    class="text-xs px-2 py-1 bg-indigo-100 dark:bg-indigo-900 text-indigo-800 dark:text-indigo-200 rounded flex-shrink-0"
                    :title="`转存时间：${resource.transferred_at}`"
                  >
                    已转存
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
                <n-button size="small" type="error" @click="confirmDelete(resource)">
                  <template #icon>
                    <i class="fas fa-trash"></i>
                  </template>
                  删除
                </n-button>
              </div>
            </div>
          </div>
        </div>
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
import { useResourceApi, useCategoryApi, useTagApi, usePanApi } from '~/composables/useApi'
import { useMessage, useDialog } from 'naive-ui'
import type { FilterConfig } from '~/components/Admin/FilterBar.vue'
import type { BatchAction } from '~/components/Admin/BatchActionBar.vue'

definePageMeta({ layout: 'admin' })

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
  // 002-auto-cleanup-transfer 自动清理相关字段
  transferred_at?: string | null
  cleaned_at?: string | null
  clean_error_msg?: string
  last_clean_error_at?: string | null
}

const userStore = useUserStore()
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()
const message = useMessage()
const dialog = useDialog()

// 列表状态
const resources = ref<Resource[]>([])
const loading = ref(false)
const error = ref<any>(null)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(200)
const selectedResources = ref<number[]>([])

// 统一筛选值（FilterBar v-model 绑定）
const filterValues = ref<Record<string, any>>({
  search: '',
  category: null,
  platform: null,
})

// 编辑相关
const showEditModal = ref(false)
const editing = ref(false)
const editingResource = ref<Resource | null>(null)
const editFormRef = ref()

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
  is_public: true,
})

const editRules = {
  title: { required: true, message: '请输入资源标题', trigger: 'blur' },
  url: { required: true, message: '请输入资源链接', trigger: 'blur' },
}

// 分类/平台数据加载
const { data: categoriesData } = await useAsyncData('resourceCategories', () => categoryApi.getCategories())
const { data: platformsData } = await useAsyncData('resourcePlatforms', () => panApi.getPans())

const categoryOptions = computed(() => {
  const data = categoriesData.value as any
  const categories = data?.items || data || []
  if (!Array.isArray(categories)) return []
  return categories.map((cat: any) => ({ label: cat.name, value: cat.id }))
})

const platformOptions = computed(() => {
  const data = platformsData.value as any
  const platforms = data?.data || data || []
  if (!Array.isArray(platforms)) return []
  return platforms.map((platform: any) => ({
    label: platform.remark || platform.name,
    value: platform.id,
  }))
})

// 筛选栏配置（声明式）
const filterConfig = computed<FilterConfig>(() => ({
  search: { placeholder: '搜索资源...', key: 'search' },
  selects: [
    { key: 'category', placeholder: '选择分类', options: categoryOptions.value },
    { key: 'platform', placeholder: '选择平台', options: platformOptions.value },
  ],
}))

// 批量操作配置（BatchAction.handler 内完成 API 调用，BatchActionBar 自动反馈）
const batchActions = computed<BatchAction[]>(() => [
  {
    key: 'delete',
    label: '批量删除',
    type: 'error',
    icon: 'fas fa-trash',
    confirm: {
      title: '批量删除确认',
      content: `确定要删除选中的 ${selectedResources.value.length} 个资源吗？此操作不可恢复！`,
    },
    handler: async (ids) => {
      await resourceApi.batchDeleteResources(ids as number[])
      return {
        success: true,
        affected: ids.length,
        message: `成功删除 ${ids.length} 个资源`,
      }
    },
  },
])

// 标签搜索（保留原有远程加载逻辑）
const tagSearchKeyword = ref('')
const tagLoading = ref(false)
const tagOptions = ref([])
const tagPagination = reactive({ page: 1, pageSize: 20, total: 0 })

const getCategoryName = (categoryId: number) => {
  const category = (categoriesData.value as any)?.data?.find((cat: any) => cat.id === categoryId)
  return category?.name || '未知分类'
}

const getPlatformName = (platformId: number) => {
  const platform = (platformsData.value as any)?.find((plat: any) => plat.id === platformId)
  return platform?.remark || platform?.name || '未知平台'
}

const loadTagOptions = async (search = '', page = 1, pageSize = 20) => {
  tagLoading.value = true
  try {
    const response = await tagApi.getTags({ page, page_size: pageSize, search })
    const data = response?.items || response?.data || []
    const totalCount = response?.total || 0
    if (!Array.isArray(data)) return { options: [], total: 0 }
    const options = data.map((tag: any) => ({ label: tag.name, value: tag.id }))
    tagPagination.total = totalCount
    return { options, total: totalCount }
  } catch (err) {
    message.error('加载标签失败')
    return { options: [], total: 0 }
  } finally {
    tagLoading.value = false
  }
}

const handleTagSearch = async (query = '') => {
  tagSearchKeyword.value = query
  const { options } = await loadTagOptions(query, 1, tagPagination.pageSize)
  tagOptions.value = options
}

const handleTagScroll = async (e: any) => {
  const { scrollTop, scrollHeight, clientHeight } = e
  if (
    scrollTop + clientHeight >= scrollHeight - 10 &&
    tagOptions.value.length < tagPagination.total &&
    !tagLoading.value
  ) {
    tagPagination.page++
    const { options } = await loadTagOptions(
      tagSearchKeyword.value,
      tagPagination.page,
      tagPagination.pageSize,
    )
    const existingIds = new Set(tagOptions.value.map((opt: any) => opt.value))
    const newOptions = options.filter((opt: any) => !existingIds.has(opt.value))
    tagOptions.value = [...tagOptions.value, ...newOptions]
  }
}

// 数据加载
const fetchData = async () => {
  loading.value = true
  error.value = null
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: filterValues.value.search,
    }
    if (filterValues.value.category) params.category_id = filterValues.value.category
    if (filterValues.value.platform) params.pan_id = filterValues.value.platform

    const response = (await resourceApi.getResources(params)) as any
    if (response && response.data) {
      if (response.data.data && Array.isArray(response.data.data)) {
        resources.value = response.data.data
        total.value = response.data.total || 0
      } else {
        resources.value = response.data
        total.value = response.total || 0
      }
      selectedResources.value = []
    } else {
      resources.value = []
      total.value = 0
      selectedResources.value = []
    }
  } catch (err) {
    error.value = err
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchData()
}

const refreshData = () => {
  selectedResources.value = []
  fetchData()
}

// 选择逻辑
const toggleResourceSelection = (resourceId: number, checked: boolean) => {
  if (checked) {
    selectedResources.value.push(resourceId)
  } else {
    const index = selectedResources.value.indexOf(resourceId)
    if (index > -1) selectedResources.value.splice(index, 1)
  }
}

const isAllSelected = computed(
  () => resources.value.length > 0 && selectedResources.value.length === resources.value.length,
)

const isIndeterminate = computed(
  () =>
    selectedResources.value.length > 0 &&
    selectedResources.value.length < resources.value.length,
)

const toggleSelectAll = (checked: boolean) => {
  selectedResources.value = checked ? resources.value.map((r) => r.id) : []
}

// 批量操作完成后刷新
const handleBatchCompleted = async () => {
  selectedResources.value = []
  await fetchData()
}

// 单项删除：用 useDialog 替代原生 confirm()
const confirmDelete = (resource: Resource) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除资源「${resource.title}」吗？此操作不可恢复。`,
    positiveText: '确定删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await resourceApi.deleteResource(resource.id)
        message.success('删除成功')
        const index = resources.value.findIndex((r) => r.id === resource.id)
        if (index > -1) {
          resources.value.splice(index, 1)
          total.value = Math.max(0, total.value - 1)
        }
        const selectedIndex = selectedResources.value.indexOf(resource.id)
        if (selectedIndex > -1) selectedResources.value.splice(selectedIndex, 1)
      } catch (err) {
        message.error('删除失败: ' + (err as Error).message)
      }
    },
  })
}

// 编辑资源
const editResource = async (resource: Resource) => {
  editingResource.value = resource

  let tagIds: number[] = []
  if (resource.tags && Array.isArray(resource.tags)) {
    tagIds = resource.tags.map((tag) => tag.id)
  } else if (resource.tag_ids && Array.isArray(resource.tag_ids)) {
    tagIds = resource.tag_ids
  }

  // 确保已选标签出现在选项中
  if (tagIds.length > 0) {
    const selectedTags = []
    for (const tagId of tagIds) {
      const existingTag = tagOptions.value.find((opt: any) => opt.value === tagId)
      if (existingTag) {
        selectedTags.push(existingTag)
      } else {
        try {
          const tagDetail = await tagApi.getTag(tagId)
          if (tagDetail) {
            selectedTags.push({ label: tagDetail.name, value: tagDetail.id })
          }
        } catch {
          selectedTags.push({ label: `标签${tagId}`, value: tagId })
        }
      }
    }
    const newOptions = [...tagOptions.value]
    for (const selectedTag of selectedTags) {
      const exists = newOptions.some((opt: any) => opt.value === selectedTag.value)
      if (!exists) newOptions.push(selectedTag)
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
    is_public: resource.is_public !== undefined ? resource.is_public : true,
  }
  showEditModal.value = true
}

const handleEditSubmit = async () => {
  try {
    editing.value = true
    await editFormRef.value?.validate()
    await resourceApi.updateResource(editingResource.value!.id, editForm.value)
    message.success('更新成功')
    showEditModal.value = false
    editingResource.value = null
    await fetchData()
  } catch (err) {
    // validate 失败或 API 失败：validate 错误不弹 message（表单已显示）
    if (err && typeof err === 'object' && 'message' in err) {
      message.error('更新失败')
    }
  } finally {
    editing.value = false
  }
}

onMounted(async () => {
  userStore.initAuth()
  const { options } = await loadTagOptions('', 1, tagPagination.pageSize)
  tagOptions.value = options
  fetchData()
})
</script>

<style scoped>
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
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style>
