<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">搜索优化管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理 Meilisearch 搜索服务状态和数据同步</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="refreshStatus" :loading="refreshing" :disabled="syncProgress.is_running">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新状态
        </n-button>
        <n-button @click="navigateTo('/admin/feature-config')" type="info" :disabled="syncProgress.is_running">
          <template #icon>
            <i class="fas fa-cog"></i>
          </template>
          配置设置
        </n-button>
      </div>
    </div>

    <!-- 状态卡片 -->
    <n-card class="mb-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
        <!-- 启用状态 -->
        <div class="flex items-center space-x-2 p-2 bg-gray-50 dark:bg-gray-800 rounded">
          <i class="fas fa-power-off text-sm" :class="status.enabled ? 'text-green-500' : 'text-red-500'"></i>
          <div>
            <p class="text-xs text-gray-500 dark:text-gray-400">启用状态</p>
            <p class="text-sm font-medium" :class="status.enabled ? 'text-green-600' : 'text-red-600'">
              {{ status.enabled ? '已启用' : '未启用' }}
            </p>
          </div>
        </div>

        <!-- 健康状态 -->
        <div class="flex items-center space-x-2 p-2 bg-gray-50 dark:bg-gray-800 rounded">
          <i class="fas fa-heartbeat text-sm" :class="status.healthy ? 'text-green-500' : 'text-red-500'"></i>
          <div>
            <p class="text-xs text-gray-500 dark:text-gray-400">健康状态</p>
            <p class="text-sm font-medium" :class="status.healthy ? 'text-green-600' : 'text-red-600'">
              {{ status.healthy ? '正常' : '异常' }}
            </p>
          </div>
        </div>

        <!-- 文档数量 -->
        <div class="flex items-center space-x-2 p-2 bg-gray-50 dark:bg-gray-800 rounded">
          <i class="fas fa-database text-sm text-blue-500"></i>
          <div>
            <p class="text-xs text-gray-500 dark:text-gray-400">索引文档</p>
            <p class="text-sm font-medium text-blue-600">{{ status.documentCount || 0 }}</p>
          </div>
        </div>

        <!-- 最后检查时间 -->
        <div class="flex items-center space-x-2 p-2 bg-gray-50 dark:bg-gray-800 rounded">
          <i class="fas fa-clock text-sm text-purple-500"></i>
          <div>
            <p class="text-xs text-gray-500 dark:text-gray-400">最后检查</p>
            <p class="text-xs font-medium text-purple-600">{{ formatTime(status.lastCheck) }}</p>
          </div>
        </div>
      </div>

      <!-- 错误信息 -->
      <div v-if="status.lastError" class="mt-3 p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded">
        <div class="flex items-start space-x-2">
          <i class="fas fa-exclamation-triangle text-red-500 mt-0.5 text-sm"></i>
          <div>
            <p class="text-xs font-medium text-red-800 dark:text-red-200">错误信息</p>
            <p class="text-xs text-red-700 dark:text-red-300">{{ status.lastError }}</p>
          </div>
        </div>
      </div>
    </n-card>

    <!-- 数据同步管理 -->
    <n-card  class="mb-6">
      <div class="space-y-4">
        <!-- 标题、过滤条件和操作按钮 -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h4 class="text-lg font-semibold">资源列表</h4>
            <div class="flex space-x-3">
              <n-button 
                type="primary" 
                @click="syncAllResources" 
                :loading="syncing"
                :disabled="unsyncedCount === 0 || syncProgress.is_running"
                size="small"
              >
                <template #icon>
                  <i class="fas fa-upload"></i>
                </template>
                同步所有资源
              </n-button>
              <!-- 停止同步按钮已隐藏 -->
              <n-button 
                type="error" 
                @click="clearIndex"
                :loading="clearing"
                :disabled="syncProgress.is_running"
                size="small"
              >
                <template #icon>
                  <i class="fas fa-trash"></i>
                </template>
                清空索引
              </n-button>
              <!-- <n-button 
                type="info" 
                @click="updateIndexSettings"
                :loading="updatingSettings"
                :disabled="syncProgress.is_running"
                size="small"
              >
                <template #icon>
                  <i class="fas fa-cogs"></i>
                </template>
                更新索引设置
              </n-button> -->
            </div>
          </div>
          
          <!-- 过滤条件 -->
          <div class="flex items-center space-x-4">
            <div class="flex items-center space-x-2">
              <span class="text-sm text-gray-600 dark:text-gray-400">同步状态:</span>
              <n-select
                v-model:value="syncFilter"
                :options="syncFilterOptions"
                size="small"
                style="width: 120px"
                :disabled="syncProgress.is_running"
                @update:value="onSyncFilterChange"
              />
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm text-gray-600 dark:text-gray-400">总计: {{ totalCount }} 个</span>
            </div>
          </div>
        </div>

        <!-- 同步进度显示 -->
        <div v-if="syncProgress.is_running" class="mt-4 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <div class="flex items-center justify-between mb-2">
            <h5 class="text-sm font-medium text-blue-800 dark:text-blue-200">同步进度</h5>
            <span class="text-xs text-blue-600 dark:text-blue-300">
              批次 {{ syncProgress.current_batch }}/{{ syncProgress.total_batches }}
            </span>
          </div>
          
          <!-- 进度条 -->
          <div class="w-full bg-blue-200 dark:bg-blue-800 rounded-full h-2 mb-2">
            <div 
              class="bg-blue-600 h-2 rounded-full transition-all duration-300"
              :style="{ width: progressPercentage + '%' }"
            ></div>
          </div>
          
          <!-- 进度信息 -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-xs">
            <div>
              <span class="text-blue-600 dark:text-blue-300">已同步:</span>
              <span class="font-medium">{{ syncProgress.synced_count }}/{{ syncProgress.total_count }}</span>
            </div>
            <div>
              <span class="text-blue-600 dark:text-blue-300">进度:</span>
              <span class="font-medium">{{ progressPercentage.toFixed(1) }}%</span>
            </div>
            <div>
              <span class="text-blue-600 dark:text-blue-300">预估剩余:</span>
              <span class="font-medium">{{ syncProgress.estimated_time || '计算中...' }}</span>
            </div>
            <div>
              <span class="text-blue-600 dark:text-blue-300">开始时间:</span>
              <span class="font-medium">{{ formatTime(syncProgress.start_time) }}</span>
            </div>
          </div>
          
          <!-- 错误信息 -->
          <div v-if="syncProgress.error_message" class="mt-2 p-2 bg-red-100 dark:bg-red-900/20 rounded text-xs text-red-700 dark:text-red-300">
            <i class="fas fa-exclamation-triangle mr-1"></i>
            {{ syncProgress.error_message }}
          </div>
        </div>

        <!-- 资源列表 -->
        <div v-if="resources.length > 0">
          <n-data-table
            :columns="columns"
            :data="resources"
            :pagination="pagination"
            :max-height="400"
            virtual-scroll
            :loading="loadingResources"
          />
        </div>
        <div v-else-if="!loadingResources" class="text-center py-8 text-gray-500">
          暂无资源数据
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNotification, useDialog } from 'naive-ui'

// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

const notification = useNotification()
const dialog = useDialog()

// 状态数据
const status = ref({
  enabled: false,
  healthy: false,
  documentCount: 0,
  lastCheck: null as Date | null,
  lastError: '',
  errorCount: 0
})

const systemConfig = ref({
  meilisearch_host: '',
  meilisearch_port: '',
  meilisearch_master_key: '',
  meilisearch_index_name: ''
})

// 定义资源类型
interface Resource {
  id: number
  title: string
  category?: {
    name: string
  }
  synced_to_meilisearch: boolean
  synced_at?: string
  created_at: string
}

// 同步状态过滤选项
const syncFilterOptions = [
  { label: '全部', value: 'all' },
  { label: '已同步', value: 'synced' },
  { label: '未同步', value: 'unsynced' }
]

const syncFilter = ref('unsynced') // 默认显示未同步
const totalCount = ref(0)
const resources = ref<Resource[]>([])
const unsyncedCount = ref(0)

// 加载状态
const refreshing = ref(false)
const syncing = ref(false)
const clearing = ref(false)
const updatingSettings = ref(false)
const loadingResources = ref(false)
const stopping = ref(false)

// 同步进度
const syncProgress = ref({
  is_running: false,
  total_count: 0,
  processed_count: 0,
  synced_count: 0,
  failed_count: 0,
  start_time: null as Date | null,
  estimated_time: '',
  current_batch: 0,
  total_batches: 0,
  error_message: ''
})

// 计算进度百分比
const progressPercentage = computed(() => {
  if (syncProgress.value.total_count === 0) return 0
  return (syncProgress.value.synced_count / syncProgress.value.total_count) * 100
})

// 分页配置
const pagination = ref({
  page: 1,
  pageSize: 1000,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [500, 1000, 2000],
  onChange: (page: number) => {
    pagination.value.page = page
    fetchResources()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchResources()
  }
})

// 表格列配置
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '标题',
    key: 'title',
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '分类',
    key: 'category',
    width: 120,
    render: (row: Resource) => {
      return row.category?.name || '-'
    }
  },
  {
    title: '同步状态',
    key: 'synced_to_meilisearch',
    width: 100,
    render: (row: Resource) => {
      return row.synced_to_meilisearch ? '已同步' : '未同步'
    }
  },
  {
    title: '同步时间',
    key: 'synced_at',
    width: 180,
    render: (row: Resource) => {
      return row.synced_at ? formatTime(row.synced_at) : '-'
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render: (row: Resource) => {
      return formatTime(row.created_at)
    }
  }
]

// 格式化时间
const formatTime = (time: Date | string | null) => {
  if (!time) return '未知'
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

// 获取状态
const fetchStatus = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const response = await meilisearchApi.getStatus() as any
    
    if (response) {
      status.value = {
        enabled: response.enabled || false,
        healthy: response.healthy || false,
        documentCount: response.document_count || response.documentCount || 0,
        lastCheck: response.last_check ? new Date(response.last_check) : response.lastCheck ? new Date(response.lastCheck) : null,
        lastError: response.last_error || response.lastError || '',
        errorCount: response.error_count || response.errorCount || 0
      }
    }
  } catch (error: any) {
    console.error('获取状态失败:', error)
    notification.error({
      content: `获取状态失败: ${error?.message || error}`,
      duration: 3000
    })
  }
}

// 获取系统配置
const fetchSystemConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig() as any
    
    if (response) {
      systemConfig.value = {
        meilisearch_host: response.meilisearch_host || '',
        meilisearch_port: response.meilisearch_port || '',
        meilisearch_master_key: response.meilisearch_master_key || '',
        meilisearch_index_name: response.meilisearch_index_name || ''
      }
    }
  } catch (error: any) {
    console.error('获取系统配置失败:', error)
  }
}

// 获取未同步数量
const fetchUnsyncedCount = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const response = await meilisearchApi.getUnsyncedCount() as any
    
    if (response) {
      unsyncedCount.value = response.count || 0
    }
  } catch (error: any) {
    console.error('获取未同步数量失败:', error)
  }
}

// 刷新状态
const refreshStatus = async () => {
  refreshing.value = true
  try {
    await Promise.all([
      fetchStatus(),
      fetchSystemConfig(),
      fetchUnsyncedCount()
    ])
    notification.success({
      content: '状态刷新成功',
      duration: 2000
    })
  } catch (error: any) {
    console.error('刷新状态失败:', error)
  } finally {
    refreshing.value = false
  }
}

// 同步所有资源
const syncAllResources = async () => {
  syncing.value = true
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    await meilisearchApi.syncAllResources()
    
    notification.success({
      content: '同步已开始，请查看进度',
      duration: 3000
    })
    
    // 开始轮询进度
    startProgressPolling()
  } catch (error: any) {
    console.error('同步资源失败:', error)
    notification.error({
      content: `同步资源失败: ${error?.message || error}`,
      duration: 5000
    })
  } finally {
    syncing.value = false
  }
}

// 停止同步
const stopSync = async () => {
  stopping.value = true
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    await meilisearchApi.stopSync()
    
    notification.success({
      content: '同步已停止',
      duration: 3000
    })
    
    // 立即更新进度状态为已停止
    syncProgress.value.is_running = false
    syncProgress.value.error_message = '同步已停止'
    
    // 停止轮询
    stopProgressPolling()
    
    // 刷新状态
    await refreshStatus()
  } catch (error: any) {
    console.error('停止同步失败:', error)
    notification.error({
      content: `停止同步失败: ${error?.message || error}`,
      duration: 5000
    })
  } finally {
    stopping.value = false
  }
}

// 进度轮询
let progressInterval: NodeJS.Timeout | null = null

const startProgressPolling = () => {
  // 清除之前的轮询
  if (progressInterval) {
    clearInterval(progressInterval)
  }
  
  // 立即获取一次进度
  fetchSyncProgress()
  
  // 每2秒轮询一次
  progressInterval = setInterval(() => {
    fetchSyncProgress()
  }, 2000)
}

const stopProgressPolling = () => {
  if (progressInterval) {
    clearInterval(progressInterval)
    progressInterval = null
  }
}

const fetchSyncProgress = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const progress = await meilisearchApi.getSyncProgress() as any
    
    if (progress) {
      syncProgress.value = {
        is_running: progress.is_running || false,
        total_count: progress.total_count || 0,
        processed_count: progress.processed_count || 0,
        synced_count: progress.synced_count || 0,
        failed_count: progress.failed_count || 0,
        start_time: progress.start_time ? new Date(progress.start_time) : null,
        estimated_time: progress.estimated_time || '',
        current_batch: progress.current_batch || 0,
        total_batches: progress.total_batches || 0,
        error_message: progress.error_message || ''
      }
      
      // 如果同步完成或出错，停止轮询
      if (!progress.is_running) {
        stopProgressPolling()
        
        // 只有在有同步进度时才显示完成消息
        if (progress.synced_count > 0 || progress.error_message) {
          if (progress.error_message) {
            notification.error({
              content: `同步失败: ${progress.error_message}`,
              duration: 5000
            })
          } else {
            notification.success({
              content: `同步完成，共同步 ${progress.synced_count} 个资源`,
              duration: 3000
            })
          }
        }
        
        // 刷新状态和表格
        await Promise.all([
          refreshStatus(),
          fetchResources()
        ])
      }
    }
  } catch (error: any) {
    console.error('获取同步进度失败:', error)
  }
}

// 静默获取同步进度，不显示任何提示
const fetchSyncProgressSilent = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const progress = await meilisearchApi.getSyncProgress() as any
    
    if (progress) {
      syncProgress.value = {
        is_running: progress.is_running || false,
        total_count: progress.total_count || 0,
        processed_count: progress.processed_count || 0,
        synced_count: progress.synced_count || 0,
        failed_count: progress.failed_count || 0,
        start_time: progress.start_time ? new Date(progress.start_time) : null,
        estimated_time: progress.estimated_time || '',
        current_batch: progress.current_batch || 0,
        total_batches: progress.total_batches || 0,
        error_message: progress.error_message || ''
      }
      
      // 如果同步完成或出错，停止轮询
      if (!progress.is_running) {
        stopProgressPolling()
        
        // 静默刷新状态和表格，不显示任何提示
        await Promise.all([
          refreshStatus(),
          fetchResources()
        ])
      }
    }
  } catch (error: any) {
    console.error('获取同步进度失败:', error)
  }
}

// 同步状态过滤变化处理
const onSyncFilterChange = () => {
  pagination.value.page = 1
  fetchResources()
}

// 获取资源列表
const fetchResources = async () => {
  loadingResources.value = true
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    
    let response: any
    if (syncFilter.value === 'unsynced') {
      // 获取未同步资源
      response = await meilisearchApi.getUnsyncedResources({
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      })
    } else if (syncFilter.value === 'synced') {
      // 获取已同步资源
      response = await meilisearchApi.getSyncedResources({
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      })
    } else {
      // 获取所有资源
      response = await meilisearchApi.getAllResources({
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      })
    }
    
    if (response && response.resources) {
      resources.value = response.resources
      totalCount.value = response.total || 0
      // 更新分页信息
      if (response.total !== undefined) {
        pagination.value.itemCount = response.total
      }
    }
  } catch (error: any) {
    console.error('获取资源失败:', error)
    notification.error({
      content: `获取资源失败: ${error?.message || error}`,
      duration: 3000
    })
  } finally {
    loadingResources.value = false
  }
}

// 获取未同步资源（保留兼容性）
const fetchUnsyncedResources = async () => {
  syncFilter.value = 'unsynced'
  await fetchResources()
}

// 清空索引
const clearIndex = async () => {
  try {
    await new Promise((resolve, reject) => {
      dialog.error({
        title: '确认清空索引',
        content: '此操作将清空所有 Meilisearch 索引数据，确定要继续吗？',
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: resolve,
        onNegativeClick: reject
      })
    })
    
    clearing.value = true
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    await meilisearchApi.clearIndex()
    
    notification.success({
      content: '索引清空成功',
      duration: 3000
    })
    
    // 刷新状态
    await refreshStatus()
  } catch (error: any) {
    if (error) {
      console.error('清空索引失败:', error)
      notification.error({
        content: `清空索引失败: ${error?.message || error}`,
        duration: 5000
      })
    }
  } finally {
    clearing.value = false
  }
}

// 更新索引设置
const updateIndexSettings = async () => {
  updatingSettings.value = true
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    await meilisearchApi.updateIndexSettings()
    
    notification.success({
      content: '索引设置已更新',
      duration: 3000
    })
    
    // 刷新状态
    await refreshStatus()
  } catch (error: any) {
    console.error('更新索引设置失败:', error)
    notification.error({
      content: `更新索引设置失败: ${error?.message || error}`,
      duration: 5000
    })
  } finally {
    updatingSettings.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  refreshStatus()
  fetchResources()
  // 静默检查同步进度，不显示任何提示
  fetchSyncProgressSilent().then(() => {
    // 如果检测到有同步在进行，开始轮询
    if (syncProgress.value.is_running) {
      startProgressPolling()
    }
  })
})

// 页面卸载时清理轮询
onUnmounted(() => {
  stopProgressPolling()
})
</script>

<style scoped>
/* 自定义样式 */
</style> 