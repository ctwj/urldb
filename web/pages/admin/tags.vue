<template>
  <AdminPageLayout>
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">标签管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的资源标签</p>
      </div>
      <div class="flex space-x-3">
      </div>
    </template>

    <!-- 提示信息区域 -->
    <template #notice-section>
      <n-alert title="提交的数据中，如果包含标签，数据添加成功，会自动添加标签" type="info" />
    </template>

    <!-- 过滤栏 - 搜索和操作 -->
    <template #filter-bar>
      <div class="flex justify-between items-center">
        <div class="flex gap-2">
          <n-button @click="showAddModal = true" type="success">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加标签
          </n-button>
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input
              v-model:value="searchQuery"
              @input="debounceSearch"
              type="text"
              placeholder="搜索标签名称..."
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
    </template>

    <!-- 内容区header - 标签列表标题 -->
    <!-- <template #content-header>
      <div class="flex items-center justify-between">
        <span class="text-lg font-semibold">标签列表</span>
        <span class="text-sm text-gray-500">共 {{ total }} 个标签</span>
      </div>
    </template> -->

    <!-- 内容区 - 标签数据 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex h-full items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="tags.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无标签</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加标签"按钮创建新标签</div>
        <n-button @click="showAddModal = true" type="primary">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加标签
        </n-button>
      </div>
      <!-- 数据表格 - 自适应高度 -->
      <div v-else class="flex flex-col h-full overflow-auto">
        <n-data-table
            :columns="columns"
            :data="tags"
            :pagination="false"
            :bordered="false"
            :single-line="false"
            :loading="loading"
            :scroll-x="800"
            class="h-full"
          />
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
            @update:page="fetchData"
            @update:page-size="(size) => { pageSize = size; currentPage = 1; fetchData() }"
          />
        </div>
      </div>
    </template>

 </AdminPageLayout>
     <!-- 添加/编辑标签模态框 -->
    <n-modal v-model:show="showAddModal" preset="card" :title="editingTag ? '编辑标签' : '添加标签'" style="width: 500px">
      <n-form
        ref="formRef"
        :model="tagForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="标签名称" path="name">
          <n-input
            v-model:value="tagForm.name"
            placeholder="请输入标签名称"
          />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="tagForm.category_id"
            :options="categoryOptions"
            placeholder="请选择分类（可选）"
            clearable
          />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="tagForm.description"
            type="textarea"
            placeholder="请输入标签描述（可选）"
            :rows="3"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="showAddModal = false">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ editingTag ? '更新' : '添加' }}
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

interface Tag {
  id: number
  name: string
  description?: string
  category_id?: number
  category_name?: string
  resource_count?: number
  created_at: string
  updated_at: string
}

const notification = useNotification()
const dialog = useDialog()
const tags = ref<Tag[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const showAddModal = ref(false)
const editingTag = ref<Tag | null>(null)
const submitting = ref(false)
const formRef = ref()

// 标签表单
const tagForm = ref({
  name: '',
  description: '',
  category_id: null as number | null
})

// 表单验证规则
const rules = {
  name: {
    required: true,
    message: '请输入标签名称',
    trigger: 'blur'
  }
}

// 获取标签API
import { useTagApi, useCategoryApi } from '~/composables/useApi'
import { h } from 'vue'
const tagApi = useTagApi()
const categoryApi = useCategoryApi()

// 获取分类数据
const { data: categoriesData } = await useAsyncData('tagCategories', () => categoryApi.getCategories())

// 分类选项
const categoryOptions = computed(() => {
  const data = categoriesData.value as any
  const categories = data?.data || data || []
  return categories.items.map((cat: any) => ({
    label: cat.name,
    value: cat.id
  }))
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render: (row: Tag) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '标签名称',
    key: 'name',
    render: (row: Tag) => {
      return h('span', { title: row.name }, row.name)
    }
  },
  {
    title: '分类',
    key: 'category_name',
    width: 120,
    render: (row: Tag) => {
      if (row.category_name) {
        return h('span', {
          class: 'px-2 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-300 rounded-full text-xs'
        }, row.category_name)
      } else {
        return h('span', { class: 'text-gray-400 italic text-xs' }, '无分类')
      }
    }
  },
  {
    title: '描述',
    key: 'description',
    render: (row: Tag) => {
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
    render: (row: Tag) => {
      return h('span', {
        class: 'px-2 py-1 bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-300 rounded-full text-xs'
      }, row.resource_count || 0)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: Tag) => {
      return h('div', { class: 'flex items-center gap-2' }, [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors',
          onClick: () => editTag(row)
        }, [
          h('i', { class: 'fas fa-edit mr-1' }),
          '编辑'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
          onClick: () => deleteTag(row)
        }, [
          h('i', { class: 'fas fa-trash mr-1' }),
          '删除'
        ])
      ])
    }
  }
]

// 分页配置已经被移到模板中处理

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await tagApi.getTags({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    }) as any
    
    if (response && response.items) {
      tags.value = response.items
      total.value = response.total || 0
    } else if (Array.isArray(response)) {
      tags.value = response
      total.value = response.length
    } else {
      tags.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取标签失败:', error)
    tags.value = []
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

// 编辑标签
const editTag = (tag: Tag) => {
  editingTag.value = tag
  tagForm.value = {
    name: tag.name,
    description: tag.description || '',
    category_id: tag.category_id || null
  }
  showAddModal.value = true
}

// 删除标签
const deleteTag = async (tag: Tag) => {
  dialog.warning({
    title: '警告',
    content: `确定要删除标签"${tag.name}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await tagApi.deleteTag(tag.id)
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
    
    if (editingTag.value) {
      await tagApi.updateTag(editingTag.value.id, tagForm.value)
      notification.success({
        content: '更新成功',
        duration: 3000
      })
    } else {
      await tagApi.createTag(tagForm.value)
      notification.success({
        content: '添加成功',
        duration: 3000
      })
    }
    
    showAddModal.value = false
    editingTag.value = null
    tagForm.value = { name: '', description: '', category_id: null }
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