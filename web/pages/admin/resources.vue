<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载资源数据</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto">

      <!-- 搜索和筛选区域 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <!-- 搜索框 -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">搜索资源</label>
            <div class="relative">
              <input 
                v-model="searchQuery"
                @keyup.enter="handleSearch"
                type="text" 
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100 dark:placeholder-gray-500"
                placeholder="输入文件名或链接进行搜索..."
              />
              <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
                <i class="fas fa-search text-gray-400"></i>
              </div>
            </div>
          </div>
          
          <!-- 平台筛选 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">平台筛选</label>
            <select 
              v-model="selectedPlatform"
              @change="handleSearch"
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
            >
              <option value="">全部平台</option>
              <option v-for="platform in platforms" :key="platform.id" :value="platform.id">
                {{ platform.name }}
              </option>
            </select>
          </div>
          
          <!-- 分类筛选 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">分类筛选</label>
            <select 
              v-model="selectedCategory"
              @change="handleSearch"
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
            >
              <option value="">全部分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
        </div>
        
        <!-- 搜索按钮 -->
        <div class="mt-4 flex justify-between items-center">
          <div class="flex gap-2">
            <button 
              @click="handleSearch" 
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center gap-2"
            >
              <i class="fas fa-search"></i> 搜索
            </button>
            <button 
              @click="clearFilters" 
              class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 flex items-center gap-2"
            >
              <i class="fas fa-times"></i> 清除筛选
            </button>
          </div>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            共找到 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个资源
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <button 
            @click="showBatchModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-list"></i> 批量操作
          </button>
        </div>
        <div class="flex gap-2">
          <button 
            @click="refreshData" 
            class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 flex items-center gap-2"
          >
            <i class="fas fa-refresh"></i> 刷新
          </button>
          <button 
            @click="exportData" 
            class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 flex items-center gap-2"
          >
            <i class="fas fa-download"></i> 导出
          </button>
        </div>
      </div>

      <!-- 批量操作模态框 -->
      <div v-if="showBatchModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 max-w-2xl w-full mx-4 text-gray-900 dark:text-gray-100">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-bold">批量操作</h3>
            <button @click="closeBatchModal" class="text-gray-500 hover:text-gray-800">
              <i class="fas fa-times"></i>
            </button>
          </div>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">选择操作：</label>
              <select 
                v-model="batchAction"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
              >
                <option value="">请选择操作</option>
                <option value="delete">批量删除</option>
                <option value="update_category">批量更新分类</option>
                <option value="update_tags">批量更新标签</option>
              </select>
            </div>
            
            <div v-if="batchAction === 'update_category'">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">选择分类：</label>
              <select 
                v-model="batchCategory"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
              >
                <option value="">请选择分类</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">
                  {{ category.name }}
                </option>
              </select>
            </div>
            
            <div v-if="batchAction === 'update_tags'">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">选择标签：</label>
              <div class="space-y-2">
                <div v-for="tag in tags" :key="tag.id" class="flex items-center">
                  <input 
                    type="checkbox" 
                    :value="tag.id" 
                    v-model="batchTags"
                    class="mr-2"
                  />
                  <span class="text-sm">{{ tag.name }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <div class="flex justify-end gap-2 mt-6">
            <button @click="closeBatchModal" class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50">
              取消
            </button>
            <button @click="handleBatchAction" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              执行操作
            </button>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 sticky top-0 z-10">
              <tr>
                <th class="px-4 py-3 text-left text-sm font-medium">
                  <input 
                    type="checkbox" 
                    v-model="selectAll"
                    @change="toggleSelectAll"
                    class="mr-2"
                  />
                  ID
                </th>
                <th class="px-4 py-3 text-left text-sm font-medium">标题</th>
                <th class="px-4 py-3 text-left text-sm font-medium">平台</th>
                <th class="px-4 py-3 text-left text-sm font-medium">分类</th>
                <th class="px-4 py-3 text-left text-sm font-medium">链接</th>
                <th class="px-4 py-3 text-left text-sm font-medium">浏览量</th>
                <th class="px-4 py-3 text-left text-sm font-medium">更新时间</th>
                <th class="px-4 py-3 text-left text-sm font-medium">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700 max-h-96 overflow-y-auto">
              <tr v-if="loading" class="text-center py-8">
                <td colspan="8" class="text-gray-500 dark:text-gray-400">
                  <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
                </td>
              </tr>
              <tr v-else-if="resources.length === 0">
                <td colspan="8">
                  <div class="flex flex-col items-center justify-center py-12">
                    <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 48 48">
                      <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
                      <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
                    </svg>
                    <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无资源数据</div>
                    <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加资源"按钮快速导入资源</div>
                    <div class="flex gap-2">
                      <NuxtLink 
                        to="/add-resource" 
                        class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md transition-colors text-sm flex items-center gap-2"
                      >
                        <i class="fas fa-plus"></i> 添加资源
                      </NuxtLink>
                      <button 
                        @click="showBatchModal = true" 
                        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors text-sm flex items-center gap-2"
                      >
                        <i class="fas fa-list"></i> 批量操作
                      </button>
                    </div>
                  </div>
                </td>
              </tr>
              <tr 
                v-for="resource in resources" 
                :key="resource.id"
                class="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 font-medium">
                  <input 
                    type="checkbox" 
                    :value="resource.id"
                    v-model="selectedResources"
                    class="mr-2"
                  />
                  {{ resource.id }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  <span :title="resource.title">{{ escapeHtml(resource.title) }}</span>
                </td>
                <td class="px-4 py-3 text-sm">
                  <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                    {{ getPlatformName(resource.pan_id) }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ getCategoryName(resource.category_id) || '-' }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <a 
                    :href="checkUrlSafety(resource.url)" 
                    target="_blank" 
                    rel="noopener noreferrer"
                    class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline break-all"
                    :title="resource.url"
                  >
                    {{ escapeHtml(resource.url) }}
                  </a>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ resource.view_count || 0 }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ formatTime(resource.updated_at) }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <div class="flex gap-2">
                    <button 
                      @click="editResource(resource)"
                      class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
                      title="编辑资源"
                    >
                      <i class="fas fa-edit"></i>
                    </button>
                    <button 
                      @click="deleteResource(resource.id)"
                      class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors"
                      title="删除资源"
                    >
                      <i class="fas fa-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 分页组件 -->
      <div v-if="totalPages > 1" class="mt-6 flex justify-center">
        <div class="flex items-center space-x-4 bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4">
          <!-- 总资源数 -->
          <div class="text-sm text-gray-600 dark:text-gray-400">
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个资源
          </div>
          
          <div class="w-px h-6 bg-gray-300 dark:bg-gray-600"></div>
          
          <!-- 上一页 -->
          <button 
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage <= 1"
            class="px-4 py-2 text-sm font-medium text-gray-500 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center gap-2"
          >
            <i class="fas fa-chevron-left"></i>
            <span>上一页</span>
          </button>
          
          <!-- 页码 -->
          <template v-for="page in visiblePages" :key="page">
            <button 
              v-if="typeof page === 'number'"
              @click="goToPage(page)"
              :class="[
                'px-4 py-2 text-sm font-medium rounded-md transition-all duration-200 min-w-[40px]',
                page === currentPage 
                  ? 'bg-blue-600 text-white shadow-md' 
                  : 'text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 hover:bg-gray-200 dark:hover:bg-gray-600'
              ]"
            >
              {{ page }}
            </button>
            <span v-else class="px-3 py-2 text-sm text-gray-500">...</span>
          </template>
          
          <!-- 下一页 -->
          <button 
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage >= totalPages"
            class="px-4 py-2 text-sm font-medium text-gray-500 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center gap-2"
          >
            <span>下一页</span>
            <i class="fas fa-chevron-right"></i>
          </button>
        </div>
      </div>

      <!-- 统计信息 -->
      <div v-if="totalPages <= 1" class="mt-4 text-center">
        <div class="inline-flex items-center bg-white dark:bg-gray-800 rounded-lg shadow px-6 py-3">
          <div class="text-sm text-gray-600 dark:text-gray-400">
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个资源
          </div>
        </div>
      </div>
    </div>
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
  url: string
  pan_id?: number
  category_id?: number
  view_count: number
  created_at: string
  updated_at: string
}

interface Platform {
  id: number
  name: string
}

interface Category {
  id: number
  name: string
}

interface Tag {
  id: number
  name: string
}

const resources = ref<Resource[]>([])
const platforms = ref<Platform[]>([])
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const loading = ref(false)
const pageLoading = ref(true)

// 搜索和筛选状态
const searchQuery = ref('')
const selectedPlatform = ref('')
const selectedCategory = ref('')

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(50)
const totalCount = ref(0)
const totalPages = ref(0)

// 批量操作状态
const showBatchModal = ref(false)
const batchAction = ref('')
const batchCategory = ref('')
const batchTags = ref<number[]>([])
const selectedResources = ref<number[]>([])
const selectAll = ref(false)

// API
import { useResourceApi, usePanApi, useCategoryApi, useTagApi } from '~/composables/useApi'

const resourceApi = useResourceApi()
const panApi = usePanApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    
    if (selectedPlatform.value) {
      params.pan_id = selectedPlatform.value
    }
    
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }
    
    const response = await resourceApi.getResources(params) as any
    
    if (response && response.resources) {
      resources.value = response.resources
      totalCount.value = response.total || 0
      totalPages.value = Math.ceil((response.total || 0) / pageSize.value)
    } else if (Array.isArray(response)) {
      resources.value = response
      totalCount.value = response.length
      totalPages.value = 1
    } else {
      resources.value = []
      totalCount.value = 0
      totalPages.value = 1
    }
  } catch (error) {
    console.error('获取资源失败:', error)
    resources.value = []
    totalCount.value = 0
    totalPages.value = 1
  } finally {
    loading.value = false
  }
}

// 获取平台列表
const fetchPlatforms = async () => {
  try {
    const response = await panApi.getPans()
    platforms.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('获取平台列表失败:', error)
    platforms.value = []
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    categories.value = []
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await tagApi.getTags()
    tags.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('获取标签列表失败:', error)
    tags.value = []
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 清除筛选
const clearFilters = () => {
  searchQuery.value = ''
  selectedPlatform.value = ''
  selectedCategory.value = ''
  currentPage.value = 1
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 导出数据
const exportData = () => {
  // 实现导出功能
  console.log('导出数据功能待实现')
}

// 分页处理
const goToPage = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 计算可见页码
const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  
  if (totalPages.value <= maxVisible) {
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i)
    }
  } else {
    if (currentPage.value <= 3) {
      for (let i = 1; i <= 4; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages.value)
    } else if (currentPage.value >= totalPages.value - 2) {
      pages.push(1)
      pages.push('...')
      for (let i = totalPages.value - 3; i <= totalPages.value; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = currentPage.value - 1; i <= currentPage.value + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages.value)
    }
  }
  
  return pages
})

// 全选/取消全选
const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedResources.value = resources.value.map(r => r.id)
  } else {
    selectedResources.value = []
  }
}

// 批量操作
const handleBatchAction = async () => {
  if (selectedResources.value.length === 0) {
    alert('请选择要操作的资源')
    return
  }
  
  if (!batchAction.value) {
    alert('请选择操作类型')
    return
  }
  
  try {
    switch (batchAction.value) {
      case 'delete':
        if (confirm(`确定要删除选中的 ${selectedResources.value.length} 个资源吗？`)) {
          await resourceApi.batchDeleteResources(selectedResources.value)
          alert('批量删除成功')
        }
        break
      case 'update_category':
        if (!batchCategory.value) {
          alert('请选择分类')
          return
        }
        await Promise.all(selectedResources.value.map(id => 
          resourceApi.updateResource(id, { category_id: batchCategory.value })
        ))
        alert('批量更新分类成功')
        break
      case 'update_tags':
        await Promise.all(selectedResources.value.map(id => 
          resourceApi.updateResource(id, { tag_ids: batchTags.value })
        ))
        alert('批量更新标签成功')
        break
    }
    
    closeBatchModal()
    fetchData()
  } catch (error) {
    console.error('批量操作失败:', error)
    alert('批量操作失败')
  }
}

// 关闭批量操作模态框
const closeBatchModal = () => {
  showBatchModal.value = false
  batchAction.value = ''
  batchCategory.value = ''
  batchTags.value = []
  selectedResources.value = []
  selectAll.value = false
}

// 编辑资源
const editResource = (resource: Resource) => {
  // 跳转到编辑页面或打开编辑模态框
  console.log('编辑资源:', resource)
}

// 删除资源
const deleteResource = async (id: number) => {
  if (confirm('确定要删除这个资源吗？')) {
    try {
      await resourceApi.deleteResource(id)
      alert('删除成功')
      fetchData()
    } catch (error) {
      console.error('删除失败:', error)
      alert('删除失败')
    }
  }
}

// 工具函数
const escapeHtml = (text: string) => {
  if (!text) return ''
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

const checkUrlSafety = (url: string) => {
  if (!url) return '#'
  // 检查URL安全性，这里可以添加更多检查逻辑
  return url
}

const formatTime = (timeString: string) => {
  if (!timeString) return '-'
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

const getPlatformName = (panId?: number) => {
  if (!panId) return '未知'
  const platform = platforms.value.find(p => p.id === panId)
  return platform?.name || '未知'
}

const getCategoryName = (categoryId?: number) => {
  if (!categoryId) return null
  const category = categories.value.find(c => c.id === categoryId)
  return category?.name || null
}

// 页面初始化
onMounted(async () => {
  try {
    await Promise.all([
      fetchData(),
      fetchPlatforms(),
      fetchCategories(),
      fetchTags()
    ])
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 