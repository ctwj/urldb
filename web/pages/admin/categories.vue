<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">分类管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的资源分类</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加分类
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 提示信息 -->
    <n-alert title="分类用于对资源进行分类管理，可以关联多个标签" type="info" />

    <!-- 搜索和操作 -->
    <n-card>
      <div class="flex justify-between items-center">
        <div class="flex gap-2">
          <n-button @click="showAddModal = true" type="success">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加分类
          </n-button>
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input 
              v-model:value="searchQuery" 
              @input="debounceSearch" 
              type="text"
              placeholder="搜索分类名称..."
              clearable
            >
              <template #prefix>
                <i class="fas fa-search text-gray-400 text-sm"></i>
              </template>
            </n-input>
          </div>
          <n-button @click="refreshData" type="tertiary">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- 分类列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">分类列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个分类</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="categories.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无分类</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加分类"按钮创建新分类</div>
        <n-button @click="showAddModal = true" type="primary">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加分类
        </n-button>
      </div>

      <div v-else>
        <n-data-table
          :columns="columns"
          :data="categories"
          :pagination="pagination"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          @update:page="handlePageChange"
        />
      </div>
    </n-card>

    <!-- 添加/编辑分类模态框 -->
    <n-modal v-model:show="showAddModal" preset="card" :title="editingCategory ? '编辑分类' : '添加分类'" style="width: 500px">
      <n-form
        ref="formRef"
        :model="categoryForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="分类名称" path="name">
          <n-input
            v-model:value="categoryForm.name"
            placeholder="请输入分类名称"
          />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="categoryForm.description"
            type="textarea"
            placeholder="请输入分类描述（可选）"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="关联标签" path="tag_ids">
          <n-select
            v-model:value="categoryForm.tag_ids"
            :options="tagOptions"
            placeholder="请选择关联标签"
            multiple
            clearable
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="showAddModal = false">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingCategory ? '更新' : '添加' }}
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

interface Category {
  id: number
  name: string
  description?: string
  resource_count?: number
  tag_names?: string[]
  created_at: string
  updated_at: string
}

const notification = useNotification()
const dialog = useDialog()
const categories = ref<Category[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const showAddModal = ref(false)
const editingCategory = ref<Category | null>(null)
const submitting = ref(false)
const formRef = ref()

// 分类表单
const categoryForm = ref({
  name: '',
  description: '',
  tag_ids: []
})

// 表单验证规则
const rules = {
  name: {
    required: true,
    message: '请输入分类名称',
    trigger: 'blur'
  }
}

// 获取分类API
import { useCategoryApi, useTagApi } from '~/composables/useApi'
import { h } from 'vue'
const categoryApi = useCategoryApi()
const tagApi = useTagApi()

// 获取标签数据
const { data: tagsData } = await useAsyncData('categoryTags', () => tagApi.getTags())

// 标签选项
const tagOptions = computed(() => {
  const data = tagsData.value as any
  const tags = data?.items || data || []
  return tags.map((tag: any) => ({
    label: tag.name,
    value: tag.id
  }))
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render: (row: Category) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '分类名称',
    key: 'name',
    render: (row: Category) => {
      return h('span', { title: row.name }, row.name)
    }
  },
  {
    title: '描述',
    key: 'description',
    render: (row: Category) => {
      if (row.description) {
        return h('span', { title: row.description }, row.description)
      } else {
        return h('span', { class: 'text-gray-400 italic' }, '无描述')
      }
    }
  },
  {
    title: '资源数量',
    key: 'resource_count',
    width: 120,
    render: (row: Category) => {
      return h('span', {
        class: 'px-2 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-300 rounded-full text-xs'
      }, row.resource_count || 0)
    }
  },
  {
    title: '关联标签',
    key: 'tag_names',
    render: (row: Category) => {
      if (row.tag_names && row.tag_names.length > 0) {
        return h('span', { class: 'text-gray-800 dark:text-gray-200' }, row.tag_names.join(', '))
      } else {
        return h('span', { class: 'text-gray-400 dark:text-gray-500 italic text-xs' }, '无标签')
      }
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: Category) => {
      return h('div', { class: 'flex items-center gap-2' }, [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors',
          onClick: () => editCategory(row)
        }, [
          h('i', { class: 'fas fa-edit mr-1' }),
          '编辑'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
          onClick: () => deleteCategory(row)
        }, [
          h('i', { class: 'fas fa-trash mr-1' }),
          '删除'
        ])
      ])
    }
  }
]

// 分页配置
const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchData()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchData()
  }
}))

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await categoryApi.getCategories({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    }) as any
    
    if (response && response.items) {
      categories.value = response.items
      total.value = response.total || 0
    } else if (Array.isArray(response)) {
      categories.value = response
      total.value = response.length
    } else {
      categories.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取分类失败:', error)
    categories.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 搜索防抖
let searchTimeout: NodeJS.Timeout | null = null
const debounceSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchData()
  }, 300)
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 编辑分类
const editCategory = (category: Category) => {
  editingCategory.value = category
  categoryForm.value = {
    name: category.name,
    description: category.description || '',
    tag_ids: []
  }
  showAddModal.value = true
}

// 删除分类
const deleteCategory = async (category: Category) => {
  dialog.warning({
    title: '警告',
    content: `确定要删除分类"${category.name}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await categoryApi.deleteCategory(category.id)
        notification.success({
          content: '删除成功',
          duration: 3000
        })
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

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    await formRef.value?.validate()
    
    if (editingCategory.value) {
      await categoryApi.updateCategory(editingCategory.value.id, categoryForm.value)
      notification.success({
        content: '更新成功',
        duration: 3000
      })
    } else {
      await categoryApi.createCategory(categoryForm.value)
      notification.success({
        content: '添加成功',
        duration: 3000
      })
    }
    
    showAddModal.value = false
    editingCategory.value = null
    categoryForm.value = { name: '', description: '', tag_ids: [] }
    fetchData()
  } catch (error) {
    console.error('提交失败:', error)
    notification.error({
      content: '操作失败',
      duration: 3000
    })
  } finally {
    submitting.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})


</script>

<style scoped>
/* 自定义样式 */
</style> 