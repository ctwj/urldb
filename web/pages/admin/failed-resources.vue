<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">失败资源列表</h1>
        <p class="text-gray-600 dark:text-gray-400">显示处理失败的资源，包含错误信息</p>
      </div>
      <div class="flex space-x-3">
                 <n-button 
           @click="retryAllFailed" 
           :disabled="selectedResources.length === 0 || isProcessing"
           :type="selectedResources.length > 0 && !isProcessing ? 'success' : 'default'"
           :loading="isProcessing"
         >
           <template #icon>
             <i v-if="isProcessing" class="fas fa-spinner fa-spin"></i>
             <i v-else class="fas fa-redo"></i>
           </template>
           {{ isProcessing ? '处理中...' : `重新放入待处理池 (${selectedResources.length})` }}
         </n-button>
         <n-button 
           @click="clearAllErrors" 
           :disabled="selectedResources.length === 0 || isProcessing"
           :type="selectedResources.length > 0 && !isProcessing ? 'warning' : 'default'"
           :loading="isProcessing"
         >
           <template #icon>
             <i v-if="isProcessing" class="fas fa-spinner fa-spin"></i>
             <i v-else class="fas fa-trash"></i>
           </template>
           {{ isProcessing ? '处理中...' : `删除失败资源 (${selectedResources.length})` }}
         </n-button>
        <n-button @click="refreshData" type="info">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

         <!-- 搜索和筛选 -->
     <n-card>
       <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
         
         <n-select
           v-model:value="errorFilter"
           placeholder="选择状态"
           :options="statusOptions"
           clearable
         />
         
         <n-button type="primary" @click="handleSearch">
           <template #icon>
             <i class="fas fa-search"></i>
           </template>
           搜索
         </n-button>
       </div>
     </n-card>

         <!-- 失败资源列表 -->
     <n-card>
       <template #header>
         <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
                <span class="text-lg font-semibold">失败资源列表</span>
                <n-checkbox
                    :checked="isAllSelected"
                    :indeterminate="isIndeterminate"
                    @update:checked="toggleSelectAll"
                >
                    全选
                </n-checkbox>
            </div>
           
           <div class="text-sm text-gray-500">
             <span class="text-sm text-gray-500">共 {{ totalCount }} 个资源，已选择 {{ selectedResources.length }} 个</span>
           </div>
         </div>
       </template>

       <!-- 加载状态 -->
       <div v-if="loading" class="flex items-center justify-center py-12">
         <n-spin size="large">
           <template #description>
             <span class="text-gray-500">加载中...</span>
           </template>
         </n-spin>
       </div>

       <!-- 虚拟列表 -->
       <n-virtual-list
         v-if="!loading"
         :items="failedResources"
         :item-size="80"
         :item-resizable="true"
         style="max-height: 400px"
         container-style="height: 600px;"
       >
         <template #default="{ item }">
           <div class="border-b border-gray-200 dark:border-gray-700 p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                        <div class="flex items-center justify-between">
             <!-- 左侧信息 -->
             <div class="flex-1 min-w-0">
               <div class="flex items-center space-x-4">
                 <!-- 复选框 -->
                 <n-checkbox
                   :checked="selectedResources.includes(item.id)"
                   @update:checked="(checked) => {
                     if (checked) {
                       selectedResources.push(item.id)
                     } else {
                       const index = selectedResources.indexOf(item.id)
                       if (index > -1) {
                         selectedResources.splice(index, 1)
                       }
                     }
                   }"
                 />
                 
                 <!-- ID -->
                 <div class="w-16 text-sm font-medium text-gray-900 dark:text-gray-100">
                   #{{ item.id }}
                 </div>
                   
                   <!-- 标题 -->
                   <div class="flex-1 min-w-0">
                     <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1" :title="item.title || '未设置'">
                       {{ item.title || '未设置' }}
                     </h3>
                   </div>
                 </div>
                 
                 <!-- 错误信息 -->
                 <div class="mt-2 flex items-center space-x-2">
                    <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-1 mt-1" :title="item.url">
                       <a 
                         :href="checkUrlSafety(item.url)" 
                         target="_blank" 
                         rel="noopener noreferrer"
                         class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline"
                       >
                         {{ item.url }}
                       </a>
                     </p>
                   <n-tag type="error" size="small" :title="item.error_msg">
                     {{ truncateError(item.error_msg) }}
                   </n-tag>
                 </div>
                 
                 <!-- 底部信息 -->
                 <div class="flex items-center space-x-4 mt-2 text-xs text-gray-500 dark:text-gray-400">
                   <span>创建时间: {{ formatTime(item.create_time) }}</span>
                   <span>IP: {{ item.ip || '-' }}</span>
                 </div>
               </div>
               
               <!-- 右侧操作按钮 -->
               <div class="flex items-center space-x-2 ml-4">
                 <n-button 
                   size="small" 
                   type="success"
                   @click="retryResource(item.id)"
                   title="重试此资源"
                 >
                   <template #icon>
                     <i class="fas fa-redo"></i>
                   </template>
                 </n-button>
                 <n-button 
                   size="small" 
                   type="warning"
                   @click="clearError(item.id)"
                   title="清除错误信息"
                 >
                   <template #icon>
                     <i class="fas fa-broom"></i>
                   </template>
                 </n-button>
                 <n-button 
                   size="small" 
                   type="error"
                   @click="deleteResource(item.id)"
                   title="删除此资源"
                 >
                   <template #icon>
                     <i class="fas fa-trash"></i>
                   </template>
                 </n-button>
               </div>
             </div>
           </div>
         </template>
       </n-virtual-list>

       <!-- 空状态 -->
       <div v-if="!loading && failedResources.length === 0" class="flex flex-col items-center justify-center py-12">
         <n-empty description="暂无失败资源">
           <template #icon>
             <i class="fas fa-check-circle text-4xl text-green-500"></i>
           </template>
           <template #extra>
             <span class="text-sm text-gray-500">所有资源处理成功</span>
           </template>
         </n-empty>
       </div>

       <!-- 分页 -->
       <div class="mt-6 flex justify-center">
         <n-pagination
           v-model:page="currentPage"
           v-model:page-size="pageSize"
           :item-count="totalCount"
           :page-sizes="[100, 200, 500, 1000]"
           show-size-picker
           @update:page="fetchData"
           @update:page-size="(size) => { pageSize = size; currentPage = 1; fetchData() }"
         />
       </div>
     </n-card>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
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
const dialog = useDialog()
const failedResources = ref<FailedResource[]>([])
const loading = ref(false)

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(200)
const totalCount = ref(0)
const totalPages = ref(0)

 // 过滤相关状态
 const errorFilter = ref('')
 const selectedStatus = ref<string | null>(null)
 
 // 处理状态
 const isProcessing = ref(false)
 
 // 选择相关状态
 const selectedResources = ref<number[]>([])
 
 // 状态选项
 const statusOptions = [
   { label: '好友已取消了分享', value: '好友已取消了分享' },
   { label: '用户封禁', value: '用户封禁' },
   { label: '分享地址已失效', value: '分享地址已失效' },
   { label: '链接无效: 链接状态检查失败', value: '链接无效: 链接状态检查失败' }
 ]

 // 获取失败资源API
 import { useReadyResourceApi } from '~/composables/useApi'
 const readyResourceApi = useReadyResourceApi()
 
 // 全选相关计算属性
 const isAllSelected = computed(() => {
   return failedResources.value.length > 0 && selectedResources.value.length === failedResources.value.length
 })
 
 const isIndeterminate = computed(() => {
   return selectedResources.value.length > 0 && selectedResources.value.length < failedResources.value.length
 })
 
 // 全选切换方法
 const toggleSelectAll = (checked: boolean) => {
   if (checked) {
     selectedResources.value = failedResources.value.map(resource => resource.id)
   } else {
     selectedResources.value = []
   }
 }



// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
         const params: any = {
       page: currentPage.value,
       page_size: pageSize.value
     }
     
     // 如果有错误信息过滤条件，添加到查询参数中
     if (errorFilter.value.trim()) {
       params.error_filter = errorFilter.value.trim()
     }
     
     // 如果有状态筛选条件，添加到查询参数中
     if (selectedStatus.value) {
       params.status = selectedStatus.value
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
     
     // 重置选择状态
     selectedResources.value = []
     
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



 // 搜索处理
 const handleSearch = () => {
   currentPage.value = 1 // 重置到第一页
   fetchData()
 }
 
 // 清除错误过滤
 const clearErrorFilter = () => {
   errorFilter.value = ''
   selectedStatus.value = null
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
          content: '错误信息已清除，资源将在下次调度时重新处理',
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('重试失败:', error)
        notification.error({
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
          content: '错误信息已清除',
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('清除错误失败:', error)
        notification.error({
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
        notification.success({
          content: '删除成功',
          duration: 3000
        })
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

 // 重新放入待处理池
 const retryAllFailed = async () => {
   if (selectedResources.value.length === 0) {
     notification.error({
       content: '请先选择要处理的资源',
       duration: 3000
     })
     return
   }
   
   const count = selectedResources.value.length
  
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
         // 使用选中的资源ID进行批量操作
         const response = await readyResourceApi.batchRestoreToReadyPool(selectedResources.value) as any
         notification.success({
           content: `操作完成：\n总数量：${count}\n成功处理：${response.success_count || count}\n失败：${response.failed_count || 0}`,
           duration: 3000
         })
         selectedResources.value = [] // 清空选择
         fetchData()
       } catch (error) {
         console.error('重新放入待处理池失败:', error)
         notification.error({
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
   if (selectedResources.value.length === 0) {
     notification.error({
       content: '请先选择要删除的资源',
       duration: 3000
     })
     return
   }
   
   const count = selectedResources.value.length
  
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
         console.log('开始调用删除API，选中的资源ID:', selectedResources.value)
         // 逐个删除选中的资源
         let successCount = 0
         for (const id of selectedResources.value) {
           try {
             await readyResourceApi.deleteReadyResource(id)
             successCount++
           } catch (error) {
             console.error(`删除资源 ${id} 失败:`, error)
           }
         }
         
         notification.success({
           content: `操作完成：\n删除失败资源：${successCount} 个资源`,
           duration: 3000
         })
         selectedResources.value = [] // 清空选择
         fetchData()
       } catch (error: any) {
         console.error('删除失败资源失败:', error)
         console.error('错误详情:', {
           message: error?.message,
           stack: error?.stack,
           response: error?.response
         })
         notification.error({
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
  }
})

// 设置页面标题
useHead({
  title: '失败资源列表 - 老九网盘资源数据库'
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