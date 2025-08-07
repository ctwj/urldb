<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载分类数据</p>
          </div>
        </div>
      </div>
    </div>

    <div class="p-6">
      <n-alert class="mb-4" title="分类用于对资源进行分类管理，可以关联多个标签" type="info" />

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <n-button @click="showAddModal = true" type="success">
            <i class="fas fa-plus"></i> 添加分类
          </n-button>
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input v-model:value="searchQuery" @input="debounceSearch" type="text"
              placeholder="搜索分类名称..." />
            <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
              <i class="fas fa-search text-gray-400 text-sm"></i>
            </div>
          </div>
          <n-button @click="refreshData" type="tertiary">
            <i class="fas fa-refresh"></i> 刷新
          </n-button>
        </div>
      </div>

  <!-- 分类列表 -->
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full min-w-full">
        <thead>
          <tr class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100">
            <th class="px-4 py-3 text-left text-sm font-medium">ID</th>
            <th class="px-4 py-3 text-left text-sm font-medium">分类名称</th>
            <th class="px-4 py-3 text-left text-sm font-medium">描述</th>
            <th class="px-4 py-3 text-left text-sm font-medium">资源数量</th>
            <th class="px-4 py-3 text-left text-sm font-medium">关联标签</th>
            <th class="px-4 py-3 text-left text-sm font-medium">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-if="loading" class="text-center py-8">
            <td colspan="6" class="text-gray-500 dark:text-gray-400">
              <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
            </td>
          </tr>
          <tr v-else-if="categories.length === 0" class="text-center py-8">
            <td colspan="6" class="text-gray-500 dark:text-gray-400">
              <div class="flex flex-col items-center justify-center py-12">
                <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor"
                  viewBox="0 0 48 48">
                  <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
                  <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
                </svg>
                <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无分类</div>
                <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加分类"按钮创建新分类</div>
                <n-button @click="showAddModal = true" type="primary">
                  <i class="fas fa-plus"></i> 添加分类
                </n-button>
              </div>
            </td>
          </tr>
          <tr v-for="category in categories" :key="category.id"
            class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 font-medium">{{ category.id }}</td>
            <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
              <span :title="category.name">{{ category.name }}</span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
              <span v-if="category.description" :title="category.description">{{ category.description }}</span>
              <span v-else class="text-gray-400 dark:text-gray-500 italic">无描述</span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
              <span
                class="px-2 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-300 rounded-full text-xs">
                {{ category.resource_count || 0 }}
              </span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
              <span v-if="category.tag_names && category.tag_names.length > 0" class="text-gray-800 dark:text-gray-200">
                {{ category.tag_names.join(', ') }}
              </span>
              <span v-else class="text-gray-400 dark:text-gray-500 italic text-xs">无标签</span>
            </td>
            <td class="px-4 py-3 text-sm">
              <div class="flex items-center gap-2">
                <n-button @click="editCategory(category)" type="info" size="small">
                  <i class="fas fa-edit"></i>
                </n-button>
                <n-button @click="deleteCategory(category.id)" type="error" size="small">
                  <i class="fas fa-trash"></i>
                </n-button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <!-- 分页 -->
  <div v-if="totalPages > 1" class="flex flex-wrap justify-center gap-1 sm:gap-2 mt-6">
    <button v-if="currentPage > 1" @click="goToPage(currentPage - 1)"
      class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center">
      <i class="fas fa-chevron-left mr-1"></i> 上一页
    </button>

    <button @click="goToPage(1)"
      :class="currentPage === 1 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
      class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm">
      1
    </button>

    <button v-if="totalPages > 1" @click="goToPage(2)"
      :class="currentPage === 2 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
      class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm">
      2
    </button>

    <span v-if="currentPage > 2" class="px-2 py-1 sm:px-3 sm:py-2 text-gray-500 text-sm">...</span>

    <button v-if="currentPage !== 1 && currentPage !== 2 && currentPage > 2"
      class="bg-slate-800 text-white px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm">
      {{ currentPage }}
    </button>

    <button v-if="currentPage < totalPages" @click="goToPage(currentPage + 1)"
      class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center">
      下一页 <i class="fas fa-chevron-right ml-1"></i>
    </button>
  </div>

  <!-- 统计信息 -->
  <div v-if="totalPages <= 1" class="mt-4 text-center">
    <div class="inline-flex items-center bg-white dark:bg-gray-800 rounded-lg shadow px-6 py-3">
      <div class="text-sm text-gray-600 dark:text-gray-400">
        共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个分类
      </div>
    </div>
  </div>

  <!-- 添加/编辑分类模态框 -->
  <div v-if="showAddModal" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
      <div class="p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
            {{ editingCategory ? '编辑分类' : '添加分类' }}
          </h3>
          <n-button @click="closeModal" type="tertiary" size="small">
            <i class="fas fa-times"></i>
          </n-button>
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">分类名称：</label>
            <n-input v-model:value="formData.name" type="text" required
              placeholder="请输入分类名称" />
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">描述：</label>
            <n-input v-model:value="formData.description" type="textarea"
              placeholder="请输入分类描述（可选）" />
          </div>

          <div class="flex justify-end gap-3">
            <n-button type="tertiary" @click="closeModal">
              取消
            </n-button>
            <n-button type="primary" :disabled="submitting" @click="handleSubmit">
              {{ submitting ? '提交中...' : (editingCategory ? '更新' : '添加') }}
            </n-button>
          </div>
        </form>
      </div>
    </div>
  </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin-old',
  ssr: false
})

const router = useRouter()
const userStore = useUserStore()
const config = useRuntimeConfig()
import { useCategoryApi } from '~/composables/useApi'
const categoryApi = useCategoryApi()

// 页面状态
const pageLoading = ref(true)
const loading = ref(false)
const categories = ref<any[]>([])

// 分页状态
const currentPage = ref(1)
const pageSize = ref(20)
const totalCount = ref(0)
const totalPages = ref(0)

// 搜索状态
const searchQuery = ref('')
let searchTimeout: NodeJS.Timeout | null = null

// 模态框状态
const showAddModal = ref(false)
const submitting = ref(false)
const editingCategory = ref<any>(null)
const dialog = useDialog()

// 表单数据
const formData = ref({
  name: '',
  description: ''
})

// 获取认证头
const getAuthHeaders = () => {
  return userStore.authHeaders
}

// 页面元数据
useHead({
  title: '分类管理 - 老九网盘资源数据库',
  meta: [
    { name: 'description', content: '管理网盘资源分类' },
    { name: 'keywords', content: '分类管理,资源管理' }
  ]
})

// 检查认证状态
const checkAuth = () => {
  userStore.initAuth()
  if (!userStore.isAuthenticated) {
    router.push('/')
    return
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    }
    console.log('获取分类列表参数:', params)
    const response = await categoryApi.getCategories(params)
    console.log('分类接口响应:', response)
    console.log('响应类型:', typeof response)
    console.log('响应是否为数组:', Array.isArray(response))
    
    // 适配后端API响应格式
    if (response && (response as any).items && Array.isArray((response as any).items)) {
      console.log('使用 items 格式:', (response as any).items)
      categories.value = (response as any).items
      totalCount.value = (response as any).total || 0
      totalPages.value = Math.ceil(totalCount.value / pageSize.value)
    } else if (Array.isArray(response)) {
      console.log('使用数组格式:', response)
      // 兼容旧格式
      categories.value = response
      totalCount.value = response.length
      totalPages.value = 1
    } else {
      console.log('使用默认格式:', response)
      categories.value = []
      totalCount.value = 0
      totalPages.value = 1
    }
    console.log('最终分类数据:', categories.value)
    console.log('分类数据长度:', categories.value.length)
  } catch (error) {
    console.error('获取分类列表失败:', error)
    categories.value = []
    totalCount.value = 0
    totalPages.value = 1
  } finally {
    loading.value = false
  }
}

// 搜索防抖
const debounceSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchCategories()
  }, 300)
}

// 刷新数据
const refreshData = () => {
  fetchCategories()
}

// 分页跳转
const goToPage = (page: number) => {
  currentPage.value = page
  fetchCategories()
}

// 编辑分类
const editCategory = (category: any) => {
  console.log('编辑分类:', category)
  editingCategory.value = category
  formData.value = {
    name: category.name,
    description: category.description || ''
  }
  console.log('设置表单数据:', formData.value)
  showAddModal.value = true
}

// 删除分类
const deleteCategory = async (categoryId: number) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除分类吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await categoryApi.deleteCategory(categoryId)
        await fetchCategories()
      } catch (error) {
        console.error('删除分类失败:', error)
      }
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    let response: any
    if (editingCategory.value) {
      response = await categoryApi.updateCategory(editingCategory.value.id, formData.value)
    } else {
      response = await categoryApi.createCategory(formData.value)
    }
    console.log('分类操作响应:', response)
    
    // 检查是否是恢复操作
    if (response && response.message && response.message.includes('恢复成功')) {
      console.log('检测到分类恢复操作，延迟刷新数据')
      console.log('恢复的分类信息:', response.category)
      closeModal()
      // 延迟一点时间再刷新，确保数据库状态已更新
      setTimeout(async () => {
        console.log('开始刷新分类数据...')
        await fetchCategories()
        console.log('分类数据刷新完成')
      }, 500)
      return
    }
    
    closeModal()
    await fetchCategories()
  } catch (error) {
    console.error('提交分类失败:', error)
  } finally {
    submitting.value = false
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  editingCategory.value = null
  formData.value = {
    name: '',
    description: ''
  }
}

// 格式化时间
const formatTime = (timestamp: string) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 退出登录
const handleLogout = () => {
  userStore.logout()
  navigateTo('/login')
}

// 页面加载
onMounted(async () => {
  try {
    checkAuth()
    await fetchCategories()

    // 检查URL参数，如果action=add则自动打开新增弹窗
    const route = useRoute()
    if (route.query.action === 'add') {
      showAddModal.value = true
    }
  } catch (error) {
    console.error('分类管理页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
/* 自定义样式 */
</style>