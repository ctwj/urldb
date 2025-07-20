<template>

  <!-- 操作按钮 -->
  <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <button @click="showAddModal = true"
          class="px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-white text-sm flex items-center gap-2">
          <i class="fas fa-plus"></i> 添加分类
        </button>
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <input v-model="searchQuery" @keyup="debounceSearch" type="text"
              class="w-64 px-3 py-2 rounded-md border border-gray-300 dark:border-gray-700 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200 dark:bg-gray-900 dark:text-gray-100 dark:placeholder-gray-500 transition-all text-sm"
              placeholder="搜索分类名称..." />
            <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
              <i class="fas fa-search text-gray-400 text-sm"></i>
            </div>
          </div>
          <button @click="refreshData"
            class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 flex items-center gap-2">
            <i class="fas fa-refresh"></i> 刷新
          </button>
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
                <button @click="showAddModal = true"
                  class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors text-sm flex items-center gap-2">
                  <i class="fas fa-plus"></i> 添加分类
                </button>
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
                <button @click="editCategory(category)"
                  class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
                  title="编辑分类">
                  <i class="fas fa-edit"></i>
                </button>
                <button @click="deleteCategory(category.id)"
                  class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors"
                  title="删除分类">
                  <i class="fas fa-trash"></i>
                </button>
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
          <button @click="closeModal"
            class="text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">分类名称：</label>
            <input v-model="formData.name" type="text" required
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100 dark:placeholder-gray-500"
              placeholder="请输入分类名称" />
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">描述：</label>
            <textarea v-model="formData.description" rows="3"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-100 dark:placeholder-gray-500"
              placeholder="请输入分类描述（可选）"></textarea>
          </div>

          <div class="flex justify-end gap-3">
            <button type="button" @click="closeModal"
              class="px-4 py-2 text-gray-700 dark:text-gray-300 bg-gray-200 dark:bg-gray-600 rounded-md hover:bg-gray-300 dark:hover:bg-gray-500 transition-colors">
              取消
            </button>
            <button type="submit" :disabled="submitting"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              {{ submitting ? '提交中...' : (editingCategory ? '更新' : '添加') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

const router = useRouter()
const userStore = useUserStore()
const config = useRuntimeConfig()

// 页面状态
const pageLoading = ref(true)
const loading = ref(false)
const categories = ref([])

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
const editingCategory = ref(null)

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

    const response = await $fetch('/categories', {
      baseURL: config.public.apiBase,
      params
    })


    // 解析响应
    if (response && typeof response === 'object' && 'code' in response && response.code === 200) {
      categories.value = response.data.items || []
      totalCount.value = response.data.total || 0
      totalPages.value = Math.ceil(totalCount.value / pageSize.value)
    } else {
      categories.value = response.items || []
      totalCount.value = response.total || 0
      totalPages.value = Math.ceil(totalCount.value / pageSize.value)
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
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
  editingCategory.value = category
  formData.value = {
    name: category.name,
    description: category.description || ''
  }
  showAddModal.value = true
}

// 删除分类
const deleteCategory = async (categoryId: number) => {
  if (!confirm(`确定要删除分类吗？`)) {
    return
  }

  try {
    await $fetch(`/categories/${categoryId}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders()
    })
    await fetchCategories()
  } catch (error) {
    console.error('删除分类失败:', error)
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true

    if (editingCategory.value) {
      await $fetch(`/categories/${editingCategory.value.id}`, {
        baseURL: config.public.apiBase,
        method: 'PUT',
        body: formData.value,
        headers: getAuthHeaders()
      })
    } else {
      await $fetch('/categories', {
        baseURL: config.public.apiBase,
        method: 'POST',
        body: formData.value,
        headers: getAuthHeaders()
      })
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