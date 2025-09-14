<template>
  <div class="space-y-6 h-full">

    <!-- 输入区域 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- 左侧：资源输入 -->
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            资源内容 <span class="text-red-500">*</span>
          </label>
          <n-input
            v-model:value="resourceText"
            type="textarea"
            placeholder="请输入资源内容，格式：标题和URL为一组..."
            :autosize="{ minRows: 10, maxRows: 15 }"
            show-count
            :maxlength="100000"
          />
        </div>
      </div>

      <!-- 右侧：配置选项 -->
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            默认分类
          </label>
          <CategorySelector
            v-model="selectedCategory"
            placeholder="选择分类"
            clearable
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            标签
          </label>
          <TagSelector
            v-model="selectedTags"
            placeholder="选择标签"
            multiple
            clearable
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            网盘账号 <span class="text-red-500">*</span>
          </label>
          <n-select
            v-model:value="selectedAccounts"
            :options="accountOptions"
            placeholder="选择网盘账号"
            multiple
            filterable
            :loading="accountsLoading"
            @update:value="handleAccountChange"
          >
            <template #option="scope">
              <div class="flex items-center justify-between w-full" v-if="scope && scope.option">
                <div class="flex items-center space-x-2">
                  <span class="text-sm">{{ scope.option.label || '未知账号' }}</span>
                  <n-tag v-if="scope.option.is_valid" type="success" size="small">有效</n-tag>
                  <n-tag v-else type="error" size="small">无效</n-tag>
                </div>
                <div class="text-xs text-gray-500">
                  {{ formatSpace(scope.option.left_space || 0) }}
                </div>
              </div>
            </template>
          </n-select>
          <div class="text-xs text-gray-500 mt-1">
            请选择要使用的网盘账号，系统将使用选中的账号进行转存操作
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="space-y-3 pt-4">
          <n-button 
            type="primary" 
            block 
            size="large"
            :loading="processing"
            :disabled="!resourceText.trim() || !selectedAccounts.length || processing"
            @click="handleBatchTransfer"
          >
            <template #icon>
              <i class="fas fa-upload"></i>
            </template>
            开始批量转存
          </n-button>
          
          <n-button 
            block 
            @click="clearInput"
            :disabled="processing"
          >
            <template #icon>
              <i class="fas fa-trash"></i>
            </template>
            清空内容
          </n-button>
        </div>
      </div>
    </div>

    <!-- 处理结果 -->
    <n-card v-if="results.length > 0" title="转存结果">
      <div class="space-y-4">
        <!-- 结果统计 -->
        <div class="grid grid-cols-4 gap-4">
          <div class="text-center p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
            <div class="text-xl font-bold text-blue-600">{{ results.length }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">总处理数</div>
          </div>
          <div class="text-center p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
            <div class="text-xl font-bold text-green-600">{{ successCount }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">成功</div>
          </div>
          <div class="text-center p-3 bg-red-50 dark:bg-red-900/20 rounded-lg">
            <div class="text-xl font-bold text-red-600">{{ failedCount }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">失败</div>
          </div>
          <div class="text-center p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
            <div class="text-xl font-bold text-yellow-600">{{ processingCount }}</div>
            <div class="text-sm text-gray-600 dark:text-gray-400">处理中</div>
          </div>
        </div>

        <!-- 结果列表 -->
        <n-data-table
          :columns="resultColumns"
          :data="results"
          :pagination="false"
          max-height="300"
          size="small"
        />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, h } from 'vue'
import { usePanApi, useTaskApi, useCksApi } from '~/composables/useApi'
import { useMessage } from 'naive-ui'

// 数据状态
const resourceText = ref('')
const processing = ref(false)
const results = ref<any[]>([])

// 任务状态
const currentTaskId = ref<number | null>(null)
const taskStatus = ref<any>(null)
const taskStats = ref({
  total: 0,
  pending: 0,
  processing: 0,
  completed: 0,
  failed: 0
})
const statusCheckInterval = ref<NodeJS.Timeout | null>(null)

// 配置选项
const selectedCategory = ref(null)
const selectedTags = ref([])
const selectedPlatform = ref(null)
const autoValidate = ref(true)
const skipExisting = ref(true)
const autoTransfer = ref(false)
const selectedAccounts = ref<number[]>([])

// 选项数据
const platformOptions = ref<any[]>([])
const accountOptions = ref<any[]>([])
const accountsLoading = ref(false)

// API实例
const panApi = usePanApi()
const taskApi = useTaskApi()
const cksApi = useCksApi()
const message = useMessage()

// 计算属性
const totalLines = computed(() => {
  return resourceText.value ? resourceText.value.split('\n').filter(line => line.trim()).length : 0
})

const validUrls = computed(() => {
  if (!resourceText.value) return 0
  const lines = resourceText.value.split('\n').filter(line => line.trim())
  return lines.filter(line => isValidUrl(line.trim())).length
})

const invalidUrls = computed(() => {
  return totalLines.value - validUrls.value
})

const successCount = computed(() => {
  return results.value ? results.value.filter((r: any) => r && r.status === 'success').length : 0
})

const failedCount = computed(() => {
  return results.value ? results.value.filter((r: any) => r && r.status === 'failed').length : 0
})

const processingCount = computed(() => {
  return results.value ? results.value.filter((r: any) => r && r.status === 'processing').length : 0
})

// 结果表格列
const resultColumns = [
  {
    title: '标题',
    key: 'title',
    width: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '链接',
    key: 'url',
    width: 250,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const statusMap = {
        success: { color: 'success', text: '成功', icon: 'fas fa-check' },
        failed: { color: 'error', text: '失败', icon: 'fas fa-times' },
        processing: { color: 'info', text: '处理中', icon: 'fas fa-spinner fa-spin' },
        pending: { color: 'warning', text: '等待中', icon: 'fas fa-clock' }
      }
      const status = statusMap[row.status as keyof typeof statusMap] || statusMap.failed
      const safeStatus = status || statusMap.failed
      return h('n-tag', { type: safeStatus.color }, {
        icon: () => h('i', { class: safeStatus.icon }),
        default: () => safeStatus.text
      })
    }
  },
  {
    title: '消息',
    key: 'message',
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '转存链接',
    key: 'saveUrl',
    width: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      if (row && row.saveUrl) {
        return h('a', {
          href: row.saveUrl,
          target: '_blank',
          class: 'text-blue-500 hover:text-blue-700'
        }, '查看')
      }
      return '-'
    }
  }
]

// URL验证
const isValidUrl = (url: string) => {
  try {
    new URL(url)
    // 简单检查是否包含常见网盘域名
    const diskDomains = ['quark.cn', 'pan.baidu.com', 'aliyundrive.com', 'pan.xunlei.com']
    return diskDomains.some(domain => url.includes(domain))
  } catch {
    return false
  }
}

// 获取平台选项
const fetchPlatforms = async () => {
  try {
    const result = await panApi.getPans() as any
    if (result && Array.isArray(result)) {
      platformOptions.value = result.map((item: any) => ({
        label: item.remark || item.name,
        value: item.id
      }))
    }
  } catch (error) {
    console.error('获取平台失败:', error)
  }
}

// 处理批量转存
const handleBatchTransfer = async () => {
  if (!resourceText.value.trim()) {
    message.warning('请输入资源内容')
    return
  }

  if (!selectedAccounts.value || selectedAccounts.value.length === 0) {
    message.warning('请选择至少一个网盘账号')
    return
  }

  processing.value = true
  results.value = []

  try {
    // 第一步：拆解资源信息，按照一行标题，一行链接的形式
    const resourceList = parseResourceText(resourceText.value)
    
    if (resourceList.length === 0) {
      message.warning('没有找到有效的资源信息，请按照格式要求输入：标题和URL为一组，标题必填')
      return
    }

    // 第二步：生成任务标题和数据
    const taskTitle = `批量转存任务_${new Date().toLocaleString('zh-CN')}`
    const taskData = {
      title: taskTitle,
      description: `批量转存 ${resourceList.length} 个资源，使用 ${selectedAccounts.value.length} 个账号`,
      resources: resourceList.map(item => {
        const resource: any = {
          title: item.title,
          url: item.url
        }
        if (selectedCategory.value) {
          resource.category_id = selectedCategory.value
        }
        if (selectedTags.value && selectedTags.value.length > 0) {
          resource.tags = selectedTags.value
        }
        return resource
      }),
      // 添加选择的账号信息
      selected_accounts: selectedAccounts.value
    }

    console.log('创建任务数据:', taskData)

    // 第三步：创建任务
    const taskResponse = await taskApi.createBatchTransferTask(taskData) as any
    console.log('任务创建响应:', taskResponse)

    if (!taskResponse || !taskResponse.task_id) {
      throw new Error('创建任务失败：响应数据无效')
    }

    currentTaskId.value = taskResponse.task_id
    
    // 第四步：启动任务
    await taskApi.startTask(currentTaskId.value!)
    
    // 第五步：开始实时监控任务状态
    startTaskMonitoring()
    
    message.success('任务已创建并启动，开始处理...')

  } catch (error: any) {
    console.error('创建任务失败:', error)
    message.error('创建任务失败: ' + (error.message || '未知错误'))
    processing.value = false
  } finally {
    processing.value = false
  }
}

// 解析资源文本，按照 标题\n链接 的格式（支持同一标题多个URL）
const parseResourceText = (text: string) => {
  const lines = text.split('\n').filter((line: string) => line.trim())
  const resourceList = []
  
  let currentTitle = ''
  let currentUrls = []
  
  for (const line of lines) {
    // 判断是否为 url（以 http/https 开头）
    if (/^https?:\/\//i.test(line)) {
      currentUrls.push(line.trim())
    } else {
      // 新标题，先保存上一个
      if (currentTitle && currentUrls.length > 0) {
        // 为每个URL创建一个资源项
        for (const url of currentUrls) {
          if (isValidUrl(url)) {
            resourceList.push({
              title: currentTitle,
              url: url,
              category_id: selectedCategory.value || 0,
              tags: selectedTags.value || []
            })
          }
        }
      }
      currentTitle = line.trim()
      currentUrls = []
    }
  }
  
  // 处理最后一组
  if (currentTitle && currentUrls.length > 0) {
    for (const url of currentUrls) {
      if (isValidUrl(url)) {
        resourceList.push({
          title: currentTitle,
          url: url,
          category_id: selectedCategory.value || 0,
          tags: selectedTags.value || []
        })
      }
    }
  }
  
  return resourceList
}

// 开始任务监控
const startTaskMonitoring = () => {
  if (statusCheckInterval.value) {
    clearInterval(statusCheckInterval.value)
  }
  
  statusCheckInterval.value = setInterval(async () => {
    try {
      const status = await taskApi.getTaskStatus(currentTaskId.value!) as any
      console.log('任务状态更新:', status)
      
      taskStatus.value = status
      taskStats.value = status.stats || {
        total: 0,
        pending: 0,
        processing: 0,
        completed: 0,
        failed: 0
      }
      
      // 更新结果显示
      updateResultsDisplay()
      
      // 如果任务完成，停止监控
      if (status.status === 'completed' || status.status === 'failed' || status.status === 'partial_success') {
        stopTaskMonitoring()
        processing.value = false
        
        const { completed, failed } = taskStats.value
        message.success(`批量转存完成！成功: ${completed}, 失败: ${failed}`)
      }
      
    } catch (error) {
      console.error('获取任务状态失败:', error)
      // 如果连续失败，停止监控
      stopTaskMonitoring()
      processing.value = false
    }
  }, 2000) // 每2秒检查一次
}

// 停止任务监控
const stopTaskMonitoring = () => {
  if (statusCheckInterval.value) {
    clearInterval(statusCheckInterval.value)
    statusCheckInterval.value = null
  }
}

// 更新结果显示
const updateResultsDisplay = () => {
  if (!taskStatus.value) return
  
  // 如果还没有结果，初始化
  if (results.value.length === 0) {
    const resourceList = parseResourceText(resourceText.value)
    results.value = resourceList.map(item => ({
      title: item.title,
      url: item.url,
      status: 'pending',
      message: '等待处理...',
      saveUrl: null
    }))
  }
  
  // 更新整体进度显示
  const { pending, processing, completed, failed } = taskStats.value
  const processed = completed + failed
  
  // 简单的状态更新逻辑 - 这里可以根据需要获取详细的任务项状态
  for (let i = 0; i < results.value.length; i++) {
    const result = results.value[i]
    
    if (i < completed) {
      // 已完成的项目
      result.status = 'success'
      result.message = '转存成功'
    } else if (i < completed + failed) {
      // 失败的项目
      result.status = 'failed'
      result.message = '转存失败'
    } else if (i < processed + processing) {
      // 正在处理的项目
      result.status = 'processing'
      result.message = '正在处理...'
    } else {
      // 等待处理的项目
      result.status = 'pending'
      result.message = '等待处理...'
    }
  }
}

// 获取网盘账号选项
const getAccountOptions = async () => {
  accountsLoading.value = true
  try {
    const response = await cksApi.getCks() as any
    const accounts = Array.isArray(response) ? response : []
    
    accountOptions.value = accounts.map((account: any) => {
      if (!account) return null
      return {
        label: `${account.username || '未知用户'} (${account.pan?.name || '未知平台'})`,
        value: account.id,
        is_valid: account.is_valid || false,
        left_space: account.left_space || 0,
        username: account.username || '未知用户',
        pan_name: account.pan?.name || '未知平台'
      }
    }).filter(option => option !== null) as any[]
  } catch (error) {
    console.error('获取网盘账号选项失败:', error)
    message.error('获取网盘账号失败')
  } finally {
    accountsLoading.value = false
  }
}

// 处理账号选择变化
const handleAccountChange = (value: number[]) => {
  selectedAccounts.value = value
  console.log('选择的账号:', value)
}

// 格式化空间大小
const formatSpace = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 清空输入
const clearInput = () => {
  resourceText.value = ''
  results.value = []
  selectedAccounts.value = []
}

// 初始化
onMounted(() => {
  fetchPlatforms()
  getAccountOptions()
})

// 组件销毁时清理定时器
onBeforeUnmount(() => {
  stopTaskMonitoring()
})
</script>