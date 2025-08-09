<template>
  <div class="space-y-6">
    <!-- 说明信息 -->
    <n-alert type="info" show-icon>
      <template #icon>
        <i class="fas fa-info-circle"></i>
      </template>
      批量转存功能：支持批量输入资源URL进行转存操作。每行一个链接，系统将自动处理转存任务。
    </n-alert>

    <!-- 输入区域 -->
    <n-card title="批量转存配置">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- 左侧：资源输入 -->
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              资源信息 <span class="text-red-500">*</span>
            </label>
            <n-input
              v-model:value="resourceText"
              type="textarea"
              placeholder="请输入资源信息，每行格式：标题|链接地址&#10;例如：&#10;电影名称1|https://pan.quark.cn/s/xxx&#10;电影名称2|https://pan.baidu.com/s/xxx"
              :rows="12"
              show-count
              :maxlength="10000"
            />
            <p class="text-xs text-gray-500 mt-1">
              每行一个资源，格式：标题|链接地址（用竖线分隔）
            </p>
          </div>
        </div>

        <!-- 右侧：配置选项 -->
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              默认分类
            </label>
            <n-select
              v-model:value="selectedCategory"
              placeholder="选择分类"
              :options="categoryOptions"
              clearable
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              标签
            </label>
            <n-select
              v-model:value="selectedTags"
              placeholder="选择标签"
              :options="tagOptions"
              multiple
              clearable
            />
          </div>

          <!-- 操作按钮 -->
          <div class="space-y-3 pt-4">
            <n-button 
              type="primary" 
              block 
              size="large"
              :loading="processing"
              :disabled="!resourceText.trim() || processing"
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
              清空输入
            </n-button>
          </div>
        </div>
      </div>
    </n-card>

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
import { ref, computed, onMounted, h } from 'vue'
import { useCategoryApi, useTagApi, usePanApi } from '~/composables/useApi'

// 数据状态
const resourceText = ref('')
const processing = ref(false)
const results = ref([])

// 配置选项
const selectedCategory = ref(null)
const selectedTags = ref([])
const selectedPlatform = ref(null)
const autoValidate = ref(true)
const skipExisting = ref(true)
const autoTransfer = ref(false)

// 选项数据
const categoryOptions = ref([])
const tagOptions = ref([])
const platformOptions = ref([])

// API实例
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()

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
  return results.value.filter((r: any) => r.status === 'success').length
})

const failedCount = computed(() => {
  return results.value.filter((r: any) => r.status === 'failed').length
})

const processingCount = computed(() => {
  return results.value.filter((r: any) => r.status === 'processing').length
})

// 结果表格列
const resultColumns = [
  {
    title: '链接',
    key: 'url',
    width: 300,
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
        processing: { color: 'info', text: '处理中', icon: 'fas fa-spinner fa-spin' }
      }
      const status = statusMap[row.status] || statusMap.failed
      return h('n-tag', { type: status.color }, {
        icon: () => h('i', { class: status.icon }),
        default: () => status.text
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
      if (row.saveUrl) {
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
    const diskDomains = ['quark.cn', 'pan.baidu.com', 'aliyundrive.com']
    return diskDomains.some(domain => url.includes(domain))
  } catch {
    return false
  }
}

// 获取分类选项
const fetchCategories = async () => {
  try {
    const result = await categoryApi.getCategories() as any
    if (result && result.items) {
      categoryOptions.value = result.items.map((item: any) => ({
        label: item.name,
        value: item.id
      }))
    }
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取标签选项
const fetchTags = async () => {
  try {
    const result = await tagApi.getTags() as any
    if (result && result.items) {
      tagOptions.value = result.items.map((item: any) => ({
        label: item.name,
        value: item.id
      }))
    }
  } catch (error) {
    console.error('获取标签失败:', error)
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
    $message.warning('请输入资源链接')
    return
  }

  processing.value = true
  results.value = []

  try {
    const lines = resourceText.value.split('\n').filter(line => line.trim())
    const validLines = lines.filter(line => isValidUrl(line.trim()))

    if (validLines.length === 0) {
      $message.warning('没有找到有效的资源链接')
      return
    }

    // 初始化结果
    results.value = validLines.map(url => ({
      url: url.trim(),
      status: 'processing',
      message: '准备处理...',
      saveUrl: null
    }))

    // 这里应该调用实际的批量转存API
    // 由于只是UI展示，这里模拟处理过程
    for (let i = 0; i < results.value.length; i++) {
      const result = results.value[i]
      
      // 模拟处理延迟
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      // 模拟随机成功/失败
      const isSuccess = Math.random() > 0.3
      
      if (isSuccess) {
        result.status = 'success'
        result.message = '转存成功'
        result.saveUrl = `https://pan.quark.cn/s/mock${Date.now()}`
      } else {
        result.status = 'failed'
        result.message = '转存失败：网络错误'
      }

      // 触发响应式更新
      results.value = [...results.value]
    }

    $message.success(`批量转存完成，成功 ${successCount.value} 个，失败 ${failedCount.value} 个`)

  } catch (error) {
    console.error('批量转存失败:', error)
    $message.error('批量转存失败')
  } finally {
    processing.value = false
  }
}

// 清空输入
const clearInput = () => {
  resourceText.value = ''
  results.value = []
}

// 初始化
onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchPlatforms()
})
</script>