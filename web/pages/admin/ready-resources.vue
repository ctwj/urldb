<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">待处理资源</h1>
        <p class="text-gray-600 dark:text-gray-400">管理待处理的资源</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="navigateTo('/admin/failed-resources')" type="error">
          <template #icon>
            <i class="fas fa-exclamation-triangle"></i>
          </template>
          错误资源
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 自动处理配置状态 -->
    <n-card>
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
          <n-button 
            @click="refreshConfig"
            :disabled="updatingConfig"
            size="small"
            type="tertiary"
            title="刷新配置"
          >
            <template #icon>
              <i class="fas fa-sync-alt"></i>
            </template>
          </n-button>
          <n-button 
            @click="toggleAutoProcess"
            :disabled="updatingConfig"
            :type="systemConfig?.auto_process_ready_resources ? 'error' : 'success'"
            size="small"
          >
            <template #icon>
              <i v-if="updatingConfig" class="fas fa-spinner fa-spin"></i>
              <i v-else :class="systemConfig?.auto_process_ready_resources ? 'fas fa-pause' : 'fas fa-play'"></i>
            </template>
            {{ systemConfig?.auto_process_ready_resources ? '关闭' : '开启' }}
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- 资源列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">待处理资源列表</span>
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-500">共 {{ totalCount }} 个待处理资源</span>
            <n-button 
              @click="clearAll" 
              type="error"
              size="small"
            >
              <template #icon>
                <i class="fas fa-trash"></i>
              </template>
              清空全部
            </n-button>
          </div>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="readyResources.length === 0" class="text-center py-8">
        <i class="fas fa-inbox text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无待处理资源</p>
        <p class="text-sm text-gray-400 mt-2">你可以点击上方"添加资源"按钮快速导入资源</p>
        <div class="mt-4">
          <n-button type="primary" @click="navigateTo('/admin/add-resource')">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加资源
          </n-button>
        </div>
      </div>

      <div v-else>
        <n-data-table
          :columns="columns"
          :data="readyResources"
          :pagination="pagination"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          @update:page="handlePageChange"
        />
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

interface ReadyResource {
  id: number
  title?: string
  url: string
  create_time: string
  ip?: string
}

const notification = useNotification()
const readyResources = ref<ReadyResource[]>([])
const loading = ref(false)

// 分页相关状态
const currentPage = ref(1)
const pageSize = ref(100)
const totalCount = ref(0)
const totalPages = ref(0)

// 获取待处理资源API
import { useReadyResourceApi, useSystemConfigApi } from '~/composables/useApi'
import { useSystemConfigStore } from '~/stores/systemConfig'
import { h } from 'vue'
const readyResourceApi = useReadyResourceApi()
const systemConfigApi = useSystemConfigApi()
const systemConfigStore = useSystemConfigStore()

// 获取系统配置
const systemConfig = ref<any>(null)
const updatingConfig = ref(false)
const dialog = useDialog()

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render: (row: ReadyResource) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '标题',
    key: 'title',
    render: (row: ReadyResource) => {
      if (row.title) {
        return h('span', { title: row.title }, escapeHtml(row.title))
      } else {
        return h('span', { class: 'text-gray-400 italic' }, '未设置')
      }
    }
  },
  {
    title: 'URL',
    key: 'url',
    render: (row: ReadyResource) => {
      return h('a', {
        href: checkUrlSafety(row.url),
        target: '_blank',
        rel: 'noopener noreferrer',
        class: 'text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 hover:underline break-all',
        title: row.url
      }, escapeHtml(row.url))
    }
  },
  {
    title: '创建时间',
    key: 'create_time',
    width: 180,
    render: (row: ReadyResource) => {
      return h('span', { class: 'text-gray-500' }, formatTime(row.create_time))
    }
  },
  {
    title: 'IP地址',
    key: 'ip',
    width: 120,
    render: (row: ReadyResource) => {
      return h('span', { class: 'text-gray-500' }, escapeHtml(row.ip || '-'))
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row: ReadyResource) => {
      return h('div', [
        h('button', {
          class: 'text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors',
          onClick: () => deleteResource(row.id),
          title: '删除此资源'
        }, [
          h('i', { class: 'fas fa-trash' })
        ])
      ])
    }
  }
]

// 分页配置
const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: totalCount.value,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchData()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchData()
  }
}))

const fetchSystemConfig = async () => {
  try {
    const response = await systemConfigApi.getSystemConfig()
    systemConfig.value = response
    systemConfigStore.setConfig(response)
    console.log('ready-resources页面系统配置:', response)
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
    
    if (response && response.data) {
      readyResources.value = response.data
      totalCount.value = response.total || 0
      totalPages.value = Math.ceil((response.total || 0) / pageSize.value)
    } else if (Array.isArray(response)) {
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

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 刷新配置
const refreshConfig = () => {
  fetchSystemConfig()
}

// 删除资源
const deleteResource = async (id: number) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除这个待处理资源吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await readyResourceApi.deleteReadyResource(id)
        if (readyResources.value.length === 1 && currentPage.value > 1) {
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

// 清空全部
const clearAll = async () => {
  dialog.warning({
    title: '警告',
    content: '确定要清空所有待处理资源吗？此操作不可恢复！',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        const response = await readyResourceApi.clearReadyResources() as any
        console.log('清空成功:', response)
        currentPage.value = 1
        fetchData()
        notification.success({
          content: `成功清空 ${response.data.deleted_count} 个资源`,
          duration: 3000
        })
      } catch (error) {
        console.error('清空失败:', error)
        notification.error({
          content: '清空失败',
          duration: 3000
        })
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

// 切换自动处理配置
const toggleAutoProcess = async () => {
  if (updatingConfig.value) {
    return
  }
  updatingConfig.value = true
  try {
    const newValue = !systemConfig.value?.auto_process_ready_resources
    console.log('切换自动处理配置:', newValue)
    
    const response = await systemConfigApi.toggleAutoProcess(newValue)
    console.log('切换响应:', response)
    
    systemConfig.value = response
    systemConfigStore.setConfig(response)
    
    notification.success({
      content: `自动处理配置已${newValue ? '开启' : '关闭'}`,
      duration: 3000
    })
  } catch (error: any) {
    notification.error({
      content: `切换自动处理配置失败`,
      duration: 3000
    })
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
  }
})

// 设置页面标题
useHead({
  title: '待处理资源管理 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 自定义样式 */
</style> 