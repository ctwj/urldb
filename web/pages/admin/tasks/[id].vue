<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题和返回按钮 -->
    <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
      <div class="flex items-center space-x-4">
        <n-button
          quaternary
          size="small"
          @click="navigateTo('/admin/tasks')"
        >
          <template #icon>
            <i class="fas fa-arrow-left"></i>
          </template>
          <span class="hidden sm:inline">返回任务列表</span>
          <span class="sm:hidden">返回</span>
        </n-button>
        <div>
          <h1 class="text-xl md:text-2xl font-bold text-gray-900 dark:text-white">任务详情</h1>
          <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">查看任务的详细信息和执行状态</p>
        </div>
      </div>
      
      <!-- 操作按钮 -->
      <div class="flex items-center space-x-2 flex-wrap" v-if="task">
        <n-button
          v-if="task.status === 'pending'"
          type="primary"
          size="small"
          @click="startTask"
          :loading="actionLoading"
        >
          启动任务
        </n-button>
        
        <n-button
          v-if="task.status === 'running'"
          type="warning"
          size="small"
          @click="pauseTask"
          :loading="actionLoading"
        >
          暂停任务
        </n-button>
        
        <n-button
          v-if="task.status === 'paused'"
          type="primary"
          size="small"
          @click="resumeTask"
          :loading="actionLoading"
        >
          继续任务
        </n-button>
        
        <n-button
          v-if="task.status === 'failed'"
          type="info"
          size="small"
          @click="retryTask"
          :loading="actionLoading"
        >
          重试任务
        </n-button>
        
        <n-button
          v-if="['completed', 'failed'].includes(task.status)"
          type="error"
          size="small"
          @click="deleteTask"
          :loading="actionLoading"
        >
          删除任务
        </n-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center py-8">
      <n-spin size="large" />
    </div>

    <!-- 任务详情 -->
    <div v-else-if="task" class="space-y-6">
      <!-- 基本信息卡片 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 md:p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">基本信息</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">任务ID</label>
            <p class="text-sm text-gray-900 dark:text-white mt-1">{{ task.id }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">任务标题</label>
            <p class="text-sm text-gray-900 dark:text-white mt-1 break-words">{{ task.title }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">任务类型</label>
            <div class="mt-1">
              <n-tag :type="getTaskTypeColor(task.task_type)" size="small">
                {{ getTaskTypeText(task.task_type) }}
              </n-tag>
            </div>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">任务状态</label>
            <div class="mt-1">
              <n-tag :type="getTaskStatusColor(task.status)" size="small">
                {{ getTaskStatusText(task.status) }}
              </n-tag>
            </div>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">创建时间</label>
            <p class="text-sm text-gray-900 dark:text-white mt-1">{{ formatDate(task.created_at) }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">更新时间</label>
            <p class="text-sm text-gray-900 dark:text-white mt-1">{{ formatDate(task.updated_at) }}</p>
          </div>
        </div>
        
        <div v-if="task.description" class="mt-4">
          <label class="text-sm font-medium text-gray-500 dark:text-gray-400">任务描述</label>
          <p class="text-sm text-gray-900 dark:text-white mt-1 break-words">{{ task.description }}</p>
        </div>
      </div>

      <!-- 进度信息卡片 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4 md:p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">进度信息</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">总项目数</label>
            <p class="text-xl md:text-2xl font-bold text-gray-900 dark:text-white mt-1">{{ task.total_items || 0 }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">已处理</label>
            <p class="text-xl md:text-2xl font-bold text-blue-600 dark:text-blue-400 mt-1">{{ task.processed_items || 0 }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">成功</label>
            <p class="text-xl md:text-2xl font-bold text-green-600 dark:text-green-400 mt-1">{{ task.success_items || 0 }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-500 dark:text-gray-400">失败</label>
            <p class="text-xl md:text-2xl font-bold text-red-600 dark:text-red-400 mt-1">{{ task.failed_items || 0 }}</p>
          </div>
        </div>
        
        <!-- 进度条 -->
        <div class="mt-4" v-if="task.total_items > 0">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">总体进度</span>
            <span class="text-sm text-gray-500 dark:text-gray-400">
              {{ Math.round((task.processed_items / task.total_items) * 100) }}%
            </span>
          </div>
          <n-progress
            type="line"
            :percentage="Math.round((task.processed_items / task.total_items) * 100)"
            :height="8"
            :show-indicator="false"
          />
        </div>
      </div>

      <!-- 任务项列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
        <div class="p-4 md:p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">任务项列表</h2>
        </div>
        
        <div class="p-4 md:p-6 overflow-x-auto">
          <n-data-table
            :columns="taskItemColumns"
            :data="taskItems"
            :loading="itemsLoading"
            :pagination="itemsPaginationConfig"
            size="small"
            :scroll-x="600"
          />
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="text-center py-8">
      <n-empty description="任务不存在或已被删除" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTaskStore } from '~/stores/task'
import { useMessage, useDialog } from 'naive-ui'

// 路由和状态管理
const route = useRoute()
const router = useRouter()
const taskStore = useTaskStore()
const message = useMessage()
const dialog = useDialog()

// 数据状态
const task = ref(null)
const taskItems = ref([])
const loading = ref(false)
const itemsLoading = ref(false)
const actionLoading = ref(false)

// 分页配置
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const itemsPaginationConfig = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchTaskItems()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchTaskItems()
  }
}))

// 任务项表格列定义
const taskItemColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    minWidth: 80
  },
  {
    title: '输入数据',
    key: 'input',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      if (!row.input) return h('span', { class: 'text-sm text-gray-500' }, '无输入数据')
      return h('span', { class: 'text-sm' }, row.input.title || row.input.url || '无标题')
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    minWidth: 100,
    render: (row: any) => {
      const statusMap: Record<string, { text: string; color: string }> = {
        pending: { text: '待处理', color: 'warning' },
        processing: { text: '处理中', color: 'info' },
        completed: { text: '已完成', color: 'success' },
        failed: { text: '失败', color: 'error' }
      }
      const status = statusMap[row.status] || { text: row.status, color: 'default' }
      return h('n-tag', { type: status.color, size: 'small' }, { default: () => status.text })
    }
  },
  {
    title: '输出数据',
    key: 'output',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      if (!row.output) return h('span', { class: 'text-sm text-gray-500' }, '无输出')
      return h('span', { class: 'text-sm' }, row.output.error || row.output.save_url || '处理完成')
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    minWidth: 160,
    render: (row: any) => {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
  }
]

// 获取任务详情
const fetchTask = async () => {
  loading.value = true
  try {
    const { useTaskApi } = await import('~/composables/useApi')
    const taskApi = useTaskApi()
    
    const response = await taskApi.getTaskStatus(parseInt(route.params.id as string)) as any
    task.value = response
  } catch (error) {
    console.error('获取任务详情失败:', error)
    message.error('获取任务详情失败')
  } finally {
    loading.value = false
  }
}

// 获取任务项列表
const fetchTaskItems = async () => {
  itemsLoading.value = true
  try {
    const { useTaskApi } = await import('~/composables/useApi')
    const taskApi = useTaskApi()
    
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    const response = await taskApi.getTaskItems(parseInt(route.params.id as string), params) as any
    
    if (response && response.items) {
      taskItems.value = response.items
      total.value = response.total || 0
    }
  } catch (error) {
    console.error('获取任务项列表失败:', error)
    message.error('获取任务项列表失败')
  } finally {
    itemsLoading.value = false
  }
}

// 任务操作
const startTask = async () => {
  actionLoading.value = true
  try {
    const success = await taskStore.startTask(task.value.id)
    if (success) {
      message.success('任务启动成功')
      await fetchTask()
    } else {
      message.error('任务启动失败')
    }
  } catch (error) {
    console.error('启动任务失败:', error)
    message.error('启动任务失败')
  } finally {
    actionLoading.value = false
  }
}

const pauseTask = async () => {
  actionLoading.value = true
  try {
    const success = await taskStore.pauseTask(task.value.id)
    if (success) {
      message.success('任务暂停成功')
      await fetchTask()
    } else {
      message.error('任务暂停失败')
    }
  } catch (error) {
    console.error('暂停任务失败:', error)
    message.error('暂停任务失败')
  } finally {
    actionLoading.value = false
  }
}

const resumeTask = async () => {
  actionLoading.value = true
  try {
    const success = await taskStore.startTask(task.value.id)
    if (success) {
      message.success('任务继续成功')
      await fetchTask()
    } else {
      message.error('任务继续失败')
    }
  } catch (error) {
    console.error('继续任务失败:', error)
    message.error('继续任务失败')
  } finally {
    actionLoading.value = false
  }
}

const retryTask = async () => {
  actionLoading.value = true
  try {
    const success = await taskStore.startTask(task.value.id)
    if (success) {
      message.success('任务重试成功')
      await fetchTask()
    } else {
      message.error('任务重试失败')
    }
  } catch (error) {
    console.error('重试任务失败:', error)
    message.error('重试任务失败')
  } finally {
    actionLoading.value = false
  }
}

const deleteTask = async () => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个任务吗？此操作不可恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      actionLoading.value = true
      try {
        const success = await taskStore.deleteTask(task.value.id)
        if (success) {
          message.success('任务删除成功')
          router.push('/admin/tasks')
        } else {
          message.error('任务删除失败')
        }
      } catch (error) {
        console.error('删除任务失败:', error)
        message.error('删除任务失败')
      } finally {
        actionLoading.value = false
      }
    }
  })
}

// 工具函数
const getTaskTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    transfer: '转存任务'
  }
  return typeMap[type] || type
}

const getTaskTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    transfer: 'blue'
  }
  return colorMap[type] || 'default'
}

const getTaskStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    paused: '暂停'
  }
  return statusMap[status] || status
}

const getTaskStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    pending: 'warning',
    running: 'info',
    completed: 'success',
    failed: 'error',
    paused: 'default'
  }
  return colorMap[status] || 'default'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

// 页面加载
onMounted(async () => {
  await fetchTask()
  await fetchTaskItems()
})

// 设置页面meta
definePageMeta({
  layout: 'admin'
})
</script>
