<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-3 sm:p-5">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载待处理资源</p>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto">

      <!-- 自动处理配置状态 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 mb-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="flex items-center space-x-2">
              <i class="fas fa-cog text-gray-600 dark:text-gray-400"></i>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">自动处理配置：</span>
            </div>
            <div class="flex items-center space-x-2">
              <div 
                :class="[
                  'w-3 h-3 rounded-full',
                  systemConfig?.auto_process_ready_resources 
                    ? 'bg-green-500 animate-pulse' 
                    : 'bg-red-500'
                ]"
              ></div>
              <span 
                :class="[
                  'text-sm font-medium',
                  systemConfig?.auto_process_ready_resources 
                    ? 'text-green-600 dark:text-green-400' 
                    : 'text-red-600 dark:text-red-400'
                ]"
              >
                {{ systemConfig?.auto_process_ready_resources ? '已开启' : '已关闭' }}
              </span>
            </div>
          </div>
          <div class="flex items-center space-x-3">
            <div class="text-xs text-gray-500 dark:text-gray-400">
              <i class="fas fa-info-circle mr-1"></i>
              {{ systemConfig?.auto_process_ready_resources 
                ? '系统会自动处理待处理资源并入库' 
                : '需要手动处理待处理资源' 
              }}
            </div>
            <button 
              @click="refreshConfig"
              :disabled="updatingConfig"
              class="px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 rounded-md transition-colors"
              title="刷新配置"
            >
              <i class="fas fa-sync-alt"></i>
            </button>
            <button 
              @click="toggleAutoProcess"
              :disabled="updatingConfig"
              :class="[
                'px-3 py-1 text-xs rounded-md transition-colors flex items-center gap-1',
                systemConfig?.auto_process_ready_resources
                  ? 'bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400'
                  : 'bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400'
              ]"
            >
              <i v-if="updatingConfig" class="fas fa-spinner fa-spin"></i>
              <i v-else :class="systemConfig?.auto_process_ready_resources ? 'fas fa-pause' : 'fas fa-play'"></i>
              {{ systemConfig?.auto_process_ready_resources ? '关闭' : '开启' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 批量添加模态框 -->
      <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl p-6 max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto text-gray-900 dark:text-gray-100">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-bold">批量添加待处理资源</h3>
            <button @click="closeModal" class="text-gray-500 hover:text-gray-800">
              <i class="fas fa-times"></i>
            </button>
          </div>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">输入格式说明：</label>
            <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded text-sm text-gray-600 dark:text-gray-300 mb-4">
              <p class="mb-2"><strong>格式1：</strong>标题和URL两行一组</p>
              <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
电影标题1
https://pan.baidu.com/s/123456
电影标题2
https://pan.baidu.com/s/789012</pre>
              <p class="mt-2 mb-2"><strong>格式2：</strong>只有URL，系统自动判断</p>
              <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
https://pan.baidu.com/s/123456
https://pan.baidu.com/s/789012
https://pan.baidu.com/s/345678</pre>
            </div>
          </div>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">资源内容：</label>
            <textarea
              v-model="resourceText"
              rows="15"
              class="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100 dark:placeholder-gray-500"
              placeholder="请输入资源内容，支持两种格式..."
            ></textarea>
          </div>
          
          <div class="flex justify-end gap-2">
            <button @click="closeModal" class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50">
              取消
            </button>
            <button @click="handleBatchAdd" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              批量添加
            </button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <NuxtLink 
            to="/add-resource" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-plus"></i> 添加资源
          </NuxtLink>
          <button 
            @click="showAddModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-list"></i> 批量添加
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
            @click="clearAll" 
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 flex items-center gap-2"
          >
            <i class="fas fa-trash"></i> 清空全部
          </button>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 sticky top-0 z-10">
              <tr>
                <th class="px-4 py-3 text-left text-sm font-medium">ID</th>
                <th class="px-4 py-3 text-left text-sm font-medium">标题</th>
                <th class="px-4 py-3 text-left text-sm font-medium">URL</th>
                <th class="px-4 py-3 text-left text-sm font-medium">创建时间</th>
                <th class="px-4 py-3 text-left text-sm font-medium">IP地址</th>
                <th class="px-4 py-3 text-left text-sm font-medium">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700 max-h-96 overflow-y-auto">
              <tr v-if="loading" class="text-center py-8">
                <td colspan="6" class="text-gray-500 dark:text-gray-400">
                  <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
                </td>
              </tr>
              <tr v-else-if="readyResources.length === 0">
                <td colspan="6">
                  <div class="flex flex-col items-center justify-center py-12">
                    <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 48 48">
                      <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
                      <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
                    </svg>
                    <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无待处理资源</div>
                    <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加资源"按钮快速导入资源</div>
                    <div class="flex gap-2">
                      <NuxtLink 
                        to="/add-resource" 
                        class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-md transition-colors text-sm flex items-center gap-2"
                      >
                        <i class="fas fa-plus"></i> 添加资源
                      </NuxtLink>
                      <button 
                        @click="showAddModal = true" 
                        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors text-sm flex items-center gap-2"
                      >
                        <i class="fas fa-list"></i> 批量添加
                      </button>
                    </div>
                  </div>
                </td>
              </tr>
              <tr 
                v-for="resource in readyResources" 
                :key="resource.id"
                class="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 font-medium">{{ resource.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  <span v-if="resource.title" :title="resource.title">{{ escapeHtml(resource.title) }}</span>
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
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ formatTime(resource.create_time) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ escapeHtml(resource.ip || '-') }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <button 
                    @click="deleteResource(resource.id)"
                    class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors"
                    title="删除此资源"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
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
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个待处理资源
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
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ totalCount }}</span> 个待处理资源
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

interface ReadyResource {
  id: number
  title?: string
  url: string
  create_time: string
  ip?: string
}

const readyResources = ref<ReadyResource[]>([])
const loading = ref(false)
const showAddModal = ref(false)
const resourceText = ref('')
const pageLoading = ref(true) // 添加页面加载状态

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(100)
const totalCount = ref(0)
const totalPages = ref(0)

// 获取待处理资源API
import { useReadyResourceApi, useSystemConfigApi } from '~/composables/useApi'
const readyResourceApi = useReadyResourceApi()
const systemConfigApi = useSystemConfigApi()

// 获取系统配置
const systemConfig = ref<any>(null)
const updatingConfig = ref(false) // 添加配置更新状态
const fetchSystemConfig = async () => {
  try {
    const response = await systemConfigApi.getSystemConfig()
    systemConfig.value = response
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await readyResourceApi.getReadyResources({
      page: currentPage.value,
      page_size: pageSize.value
    }) as any
    
    // 适配后端API响应格式
    if (response && response.data) {
      readyResources.value = response.data
      // 后端返回格式: {data: [...], page: 1, page_size: 100, total: 123}
      totalCount.value = response.total || 0
      totalPages.value = Math.ceil((response.total || 0) / pageSize.value)
    } else if (Array.isArray(response)) {
      // 如果直接返回数组
      readyResources.value = response
      totalCount.value = response.length
      totalPages.value = 1
    } else {
      readyResources.value = []
      totalCount.value = 0
      totalPages.value = 1
    }
  } catch (error) {
    console.error('获取待处理资源失败:', error)
    readyResources.value = []
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
    // 如果总页数不多，显示所有页码
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i)
    }
  } else {
    // 如果总页数很多，显示部分页码
    if (currentPage.value <= 3) {
      // 当前页在前几页
      for (let i = 1; i <= 4; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages.value)
    } else if (currentPage.value >= totalPages.value - 2) {
      // 当前页在后几页
      pages.push(1)
      pages.push('...')
      for (let i = totalPages.value - 3; i <= totalPages.value; i++) {
        pages.push(i)
      }
    } else {
      // 当前页在中间
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

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 刷新配置
const refreshConfig = () => {
  fetchSystemConfig()
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  resourceText.value = ''
}

// 批量添加
const handleBatchAdd = async () => {
  if (!resourceText.value.trim()) {
    alert('请输入资源内容')
    return
  }

  try {
    const response = await readyResourceApi.createReadyResourcesFromText(resourceText.value) as any
    console.log('批量添加成功:', response)
    closeModal()
    fetchData()
    alert(`成功添加 ${response.data.count} 个资源`)
  } catch (error) {
    console.error('批量添加失败:', error)
    alert('批量添加失败，请检查输入格式')
  }
}

// 删除资源
const deleteResource = async (id: number) => {
  if (!confirm('确定要删除这个待处理资源吗？')) {
    return
  }

  try {
    await readyResourceApi.deleteReadyResource(id)
    // 如果当前页没有数据了，回到上一页
    if (readyResources.value.length === 1 && currentPage.value > 1) {
      currentPage.value--
    }
    fetchData()
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败')
  }
}

// 清空全部
const clearAll = async () => {
  if (!confirm('确定要清空所有待处理资源吗？此操作不可恢复！')) {
    return
  }

  try {
    const response = await readyResourceApi.clearReadyResources() as any
    console.log('清空成功:', response)
    currentPage.value = 1 // 清空后回到第一页
    fetchData()
    alert(`成功清空 ${response.data.deleted_count} 个资源`)
  } catch (error) {
    console.error('清空失败:', error)
    alert('清空失败')
  }
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
    // 只允许http和https协议
    if (urlObj.protocol !== 'http:' && urlObj.protocol !== 'https:') {
      return '#'
    }
    return url
  } catch {
    return '#'
  }
}

// 切换自动处理配置
const toggleAutoProcess = async () => {
  if (updatingConfig.value) {
    return
  }
  updatingConfig.value = true
  try {
    const newValue = !systemConfig.value?.auto_process_ready_resources
    console.log('当前配置:', systemConfig.value)
    console.log('新值:', newValue)
    
    // 先获取当前配置，然后只更新需要的字段
    const currentConfig = await systemConfigApi.getSystemConfig() as any
    console.log('获取到的当前配置:', currentConfig)
    
    const updateData = {
      site_title: currentConfig.site_title || '',
      site_description: currentConfig.site_description || '',
      keywords: currentConfig.keywords || '',
      author: currentConfig.author || '',
      copyright: currentConfig.copyright || '',
      auto_process_ready_resources: newValue,
      auto_process_interval: currentConfig.auto_process_interval || 30,
      auto_transfer_enabled: currentConfig.auto_transfer_enabled || false,
      auto_fetch_hot_drama_enabled: currentConfig.auto_fetch_hot_drama_enabled || false,
      page_size: currentConfig.page_size || 100,
      maintenance_mode: currentConfig.maintenance_mode || false
    }
    
    console.log('更新数据:', updateData)
    const response = await systemConfigApi.updateSystemConfig(updateData)
    console.log('更新响应:', response)
    
    systemConfig.value = response
    alert(`自动处理配置已${newValue ? '开启' : '关闭'}`)
  } catch (error: any) {
    console.error('切换自动处理配置失败:', error)
    console.error('错误详情:', error.response || error)
    alert('切换自动处理配置失败')
  } finally {
    updatingConfig.value = false
  }
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    await fetchData()
    await fetchSystemConfig()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    // 数据加载完成后，关闭加载状态
    pageLoading.value = false
  }
})

// 设置页面标题
useHead({
  title: '待处理资源管理 - 老九网盘资源数据库'
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
  box-shadow: 0 4px 6px -1px rgba(59, 130, 246, 0.3), 0 2px 4px -1px rgba(59, 130, 246, 0.2);
}

/* 表格行悬停效果 */
tbody tr:hover {
  background-color: rgba(59, 130, 246, 0.05);
}

/* 暗黑模式下的表格行悬停 */
.dark tbody tr:hover {
  background-color: rgba(59, 130, 246, 0.1);
}

/* 统计信息卡片效果 */
.stats-card {
  backdrop-filter: blur(10px);
  background-color: rgba(255, 255, 255, 0.9);
}

.dark .stats-card {
  background-color: rgba(31, 41, 55, 0.9);
}
</style> 