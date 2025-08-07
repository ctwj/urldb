<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-3 sm:p-5">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-red-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载失败资源列表</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">失败资源列表</h1>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">显示处理失败的资源，包含错误信息</p>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <button 
            @click="retryAllFailed" 
            :disabled="!errorFilter.trim() || isProcessing"
            :class="[
              'w-full sm:w-auto px-4 py-2 rounded-md transition-colors text-center flex items-center justify-center gap-2',
              errorFilter.trim() && !isProcessing 
                ? 'bg-green-600 hover:bg-green-700' 
                : 'bg-gray-400 text-gray-200 cursor-not-allowed'
            ]"
          >
            <i v-if="isProcessing" class="fas fa-spinner fa-spin"></i>
            <i v-else class="fas fa-redo"></i> 
            {{ isProcessing ? '处理中...' : '重新放入待处理池' }}
          </button>
          <button 
            @click="clearAllErrors" 
            :disabled="!errorFilter.trim() || isProcessing"
            :class="[
              'w-full sm:w-auto px-4 py-2 rounded-md transition-colors text-center flex items-center justify-center gap-2',
              errorFilter.trim() && !isProcessing 
                ? 'bg-yellow-600 hover:bg-yellow-700' 
                : 'bg-gray-400 text-gray-200 cursor-not-allowed'
            ]"
          >
            <i v-if="isProcessing" class="fas fa-spinner fa-spin"></i>
            <i v-else class="fas fa-trash"></i> 
            {{ isProcessing ? '处理中...' : '删除失败资源' }}
          </button>
        </div>
        <div class="flex gap-2 items-center">
          <!-- 错误信息过滤 -->
          <div class="flex items-center gap-2">
            <n-input
              v-model:value="errorFilter"
              type="text"
              placeholder="过滤错误信息..."
              class="w-48"
              clearable
              @input="onErrorFilterChange"
            />
            <button 
              v-if="errorFilter"
              @click="clearErrorFilter"
              class="px-2 py-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
              title="清除过滤条件"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
          <button 
            @click="refreshData" 
            class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 flex items-center gap-2"
          >
            <i class="fas fa-refresh"></i> 刷新
          </button>
        </div>
      </div>



      <!-- 失败资源列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-red-800 dark:bg-red-900 text-white dark:text-gray-100 sticky top-0 z-10">
              <tr>
                <th class="px-4 py-3 text-left text-sm font-medium">ID</th>
                <th class="px-4 py-3 text-left text-sm font-medium">状态</th>
                <th class="px-4 py-3 text-left text-sm font-medium">标题</th>
                <th class="px-4 py-3 text-left text-sm font-medium">URL</th>
                <th class="px-4 py-3 text-left text-sm font-medium">错误信息</th>
                <th class="px-4 py-3 text-left text-sm font-medium">创建时间</th>
                <th class="px-4 py-3 text-left text-sm font-medium">IP地址</th>
                <th class="px-4 py-3 text-left text-sm font-medium">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700 max-h-96 overflow-y-auto">
              <tr v-if="loading" class="text-center py-8">
                <td colspan="8" class="text-gray-500 dark:text-gray-400">
                  <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
                </td>
              </tr>
              <tr v-else-if="failedResources.length === 0">
                <td colspan="8">
                  <div class="flex flex-col items-center justify-center py-12">
                    <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 48 48">
                      <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
                      <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
                    </svg>
                    <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无失败资源</div>
                    <div class="text-sm text-gray-400 dark:text-gray-600">所有资源处理成功</div>
                  </div>
                </td>
              </tr>
              <tr 
                v-for="resource in failedResources" 
                :key="resource.id"
                :class="[
                  'hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors',
                  resource.is_deleted ? 'bg-gray-100 dark:bg-gray-700' : ''
                ]"
              >
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 font-medium">{{ resource.id }}</td>
                <td class="px-4 py-3 text-sm">
                  <span 
                    v-if="resource.is_deleted" 
                    class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                    title="已删除"
                  >
                    <i class="fas fa-trash mr-1"></i>已删除
                  </span>
                  <span 
                    v-else 
                    class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
                    title="正常"
                  >
                    <i class="fas fa-check mr-1"></i>正常
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  <span v-if="resource.title && resource.title !== null" :title="resource.title">{{ escapeHtml(resource.title) }}</span>
                  <span v-else class="text-gray-400 dark:text-gray-500 italic">未设置</span>
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
                <td class="px-4 py-3 text-sm">
                  <div class="max-w-xs">
                    <span 
                      class="text-red-600 dark:text-red-400 text-xs bg-red-50 dark:bg-red-900/20 px-2 py-1 rounded"
                      :title="resource.error_msg"
                    >
                      {{ truncateError(resource.error_msg) }}
                    </span>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ formatTime(resource.create_time) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ escapeHtml(resource.ip || '-') }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <div class="flex gap-2">
                    <button 
                      @click="retryResource(resource.id)"
                      class="text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300 transition-colors"
                      title="重试此资源"
                    >
                      <i class="fas fa-redo"></i>
                    </button>
                    <button 
                      @click="clearError(resource.id)"
                      class="text-yellow-600 hover:text-yellow-800 dark:text-yellow-400 dark:hover:text-yellow-300 transition-colors"
                      title="清除错误信息"
                    >
                      <i class="fas fa-broom"></i>
                    </button>
                    <button 
                      @click="deleteResource(resource.id)"
                      class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors"
                      title="删除此资源"
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
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个失败资源
            <span v-if="errorFilter" class="ml-2 text-blue-600 dark:text-blue-400">
              (已过滤)
            </span>
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
                  ? 'bg-red-600 text-white shadow-md' 
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
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个失败资源
            <span v-if="errorFilter" class="ml-2 text-blue-600 dark:text-blue-400">
              (已过滤)
            </span>
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

interface FailedResource {
  id: number
  title?: string | null
  url: string
  error_msg: string
  create_time: string
  ip?: string | null
  deleted_at?: string | null
  is_deleted: boolean
}

const notification = useNotification()
const failedResources = ref<FailedResource[]>([])
const loading = ref(false)
const pageLoading = ref(true)

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(100)
const totalCount = ref(0)
const totalPages = ref(0)


const dialog = useDialog()

// 过滤相关状态
const errorFilter = ref('')

// 获取失败资源API
import { useReadyResourceApi } from '~/composables/useApi'
const readyResourceApi = useReadyResourceApi()

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    // 如果有过滤条件，添加到查询参数中
    if (errorFilter.value.trim()) {
      params.error_filter = errorFilter.value.trim()
    }
    
    console.log('fetchData - 开始获取失败资源，参数:', params)
    
    const response = await readyResourceApi.getFailedResources(params) as any
    
    console.log('fetchData - 原始响应:', response)
    
    if (response && response.data && Array.isArray(response.data)) {
      console.log('fetchData - 使用response.data格式（数组）')
      failedResources.value = response.data
      totalCount.value = response.total || 0
      totalPages.value = Math.ceil((response.total || 0) / pageSize.value)
    } else {
      console.log('fetchData - 使用空数据格式')
      failedResources.value = []
      totalCount.value = 0
      totalPages.value = 1
    }
    
    console.log('fetchData - 处理后的数据:', {
      failedResourcesCount: failedResources.value.length,
      totalCount: totalCount.value,
      totalPages: totalPages.value
    })
    
    // 打印第一个资源的数据结构（如果存在）
    if (failedResources.value.length > 0) {
      console.log('fetchData - 第一个资源的数据结构:', failedResources.value[0])
    }
  } catch (error) {
    console.error('获取失败资源失败:', error)
    failedResources.value = []
    totalCount.value = 0
    totalPages.value = 1
  } finally {
    loading.value = false
  }
}

// 跳转到指定页面
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    fetchData()
  }
}

// 计算可见的页码
const visiblePages = computed(() => {
  const pages: (number | string)[] = []
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

// 防抖函数
const debounce = (func: Function, delay: number) => {
  let timeoutId: NodeJS.Timeout
  return (...args: any[]) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func.apply(null, args), delay)
  }
}

// 错误过滤输入变化处理（防抖）
const onErrorFilterChange = debounce(() => {
  currentPage.value = 1 // 重置到第一页
  fetchData()
}, 300)

// 清除错误过滤
const clearErrorFilter = () => {
  errorFilter.value = ''
  currentPage.value = 1 // 重置到第一页
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 重试单个资源
const retryResource = async (id: number) => {
  dialog.warning({
    title: '警告',
    content: '确定要重试这个资源吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await readyResourceApi.clearErrorMsg(id)
        notification.success({
          title: '成功',
          content: '错误信息已清除，资源将在下次调度时重新处理',
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('重试失败:', error)
        notification.error({
          title: '失败',
          content: '重试失败',
          duration: 3000
        })
      }
    }
  })
}

// 清除单个资源错误
const clearError = async (id: number) => {
  dialog.warning({
    title: '警告',
    content: '确定要清除这个资源的错误信息吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await readyResourceApi.clearErrorMsg(id)
        notification.success({
          title: '成功',
          content: '错误信息已清除',
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('清除错误失败:', error)
        notification.error({
          title: '失败',
          content: '清除错误失败',
          duration: 3000
        })
      }
    }
  })
}

// 删除资源
const deleteResource = async (id: number) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除这个失败资源吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await readyResourceApi.deleteReadyResource(id)
        if (failedResources.value.length === 1 && currentPage.value > 1) {
          currentPage.value--
        }
        fetchData()
      } catch (error) {
        console.error('删除失败:', error)
        notification.error({
          title: '失败',
          content: '删除失败',
          duration: 3000
        })
      }
    }
  })
}

// 处理状态
const isProcessing = ref(false)

// 重新放入待处理池
const retryAllFailed = async () => {
  if (totalCount.value === 0) {
    notification.error({
      title: '失败',
      content: '没有可处理的资源',
      duration: 3000
    })
    return
  }
  
  // 检查是否有过滤条件
  if (!errorFilter.value.trim()) {
    notification.error({
      title: '失败',
      content: '请先设置过滤条件，以避免处理所有失败资源',
      duration: 3000
    })
    return
  }
  
  // 构建查询条件
  const queryParams: any = {}
  
  // 如果有过滤条件，添加到查询参数中
  if (errorFilter.value.trim()) {
    queryParams.error_filter = errorFilter.value.trim()
  }
  
  const count = totalCount.value
  
  dialog.warning({
    title: '确认操作',
    content: `确定要将 ${count} 个资源重新放入待处理池吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      if (isProcessing.value) return // 防止重复点击
      
      isProcessing.value = true
      
      try {
        const response = await readyResourceApi.batchRestoreToReadyPoolByQuery(queryParams) as any
        notification.success({
          title: '成功',
          content: `操作完成：\n总数量：${response.total_count}\n成功处理：${response.success_count}\n失败：${response.failed_count}`,
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('重新放入待处理池失败:', error)
        notification.error({
          title: '失败',
          content: '操作失败',
          duration: 3000
        })
      } finally {
        isProcessing.value = false
      }
    }
  })
}

// 清除所有错误
const clearAllErrors = async () => {
  // 检查是否有过滤条件
  if (!errorFilter.value.trim()) {
    notification.error({
      title: '失败',
      content: '请先设置过滤条件，以避免删除所有失败资源',
      duration: 3000
    })
    return
  }
  
  // 构建查询条件
  const queryParams: any = {}
  
  // 如果有过滤条件，添加到查询参数中
  if (errorFilter.value.trim()) {
    queryParams.error_filter = errorFilter.value.trim()
  }
  
  const count = totalCount.value
  
  dialog.warning({
    title: '警告',
    content: `确定要删除 ${count} 个失败资源吗？此操作将永久删除这些资源，不可恢复！`,
    positiveText: '确定删除',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      if (isProcessing.value) return // 防止重复点击
      
      isProcessing.value = true
      
      try {
        console.log('开始调用删除API，参数:', queryParams)
        const response = await readyResourceApi.clearAllErrorsByQuery(queryParams) as any
        // console.log('删除API响应:', response)
        notification.success({
          title: '成功',
          content: `操作完成：\n删除失败资源：${response.affected_rows} 个资源`,
          duration: 3000
        })
        fetchData()
      } catch (error: any) {
        console.error('删除失败资源失败:', error)
        console.error('错误详情:', {
          message: error?.message,
          stack: error?.stack,
          response: error?.response
        })
        notification.error({
          title: '失败',
          content: '删除失败',
          duration: 3000
        })
      } finally {
        isProcessing.value = false
      }
    }
  })
}

// 格式化时间
const formatTime = (timeString: string) => {
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

// 转义HTML防止XSS
const escapeHtml = (text: string) => {
  if (!text) return text
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

// 验证URL安全性
const checkUrlSafety = (url: string) => {
  if (!url) return '#'
  try {
    const urlObj = new URL(url)
    if (urlObj.protocol !== 'http:' && urlObj.protocol !== 'https:') {
      return '#'
    }
    return url
  } catch {
    return '#'
  }
}

// 截断错误信息
const truncateError = (errorMsg: string) => {
  if (!errorMsg) return ''
  return errorMsg.length > 50 ? errorMsg.substring(0, 50) + '...' : errorMsg
}



// 页面加载时获取数据
onMounted(async () => {
  try {
    await fetchData()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})

// 设置页面标题
useHead({
  title: '失败资源列表 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 表格滚动样式 */
.overflow-x-auto {
  max-height: 500px;
  overflow-y: auto;
}

/* 表格头部固定 */
thead {
  position: sticky;
  top: 0;
  z-index: 10;
}

/* 分页按钮悬停效果 */
.pagination-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* 当前页码按钮效果 */
.current-page {
  box-shadow: 0 4px 6px -1px rgba(220, 38, 38, 0.3), 0 2px 4px -1px rgba(220, 38, 38, 0.2);
}

/* 表格行悬停效果 */
tbody tr:hover {
  background-color: rgba(220, 38, 38, 0.05);
}

/* 暗黑模式下的表格行悬停 */
.dark tbody tr:hover {
  background-color: rgba(220, 38, 38, 0.1);
}
</style> 