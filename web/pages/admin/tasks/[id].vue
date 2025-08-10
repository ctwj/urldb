<template>
  <div class="p-4 space-y-4">
    <!-- 页面标题和返回按钮 -->
    <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-3 lg:space-y-0">
      <div class="flex items-center space-x-3">
        <n-button
          quaternary
          size="small"
          @click="navigateTo('/admin/tasks')"
        >
          <template #icon>
            <i class="fas fa-arrow-left"></i>
          </template>
          <span class="hidden sm:inline">返回</span>
        </n-button>
        <div>
          <h1 class="text-lg md:text-xl font-bold text-gray-900 dark:text-white">任务详情</h1>
          <p class="text-xs text-gray-600 dark:text-gray-400">查看任务的详细信息和执行状态</p>
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
    <div v-if="loading" class="flex justify-center py-4">
      <n-spin size="medium" />
    </div>

    <!-- 任务详情 -->
    <div v-else-if="task" class="space-y-4">
      <!-- 整合的任务信息卡片 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-4">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
          <!-- 左侧：基本信息 -->
          <div class="flex-1">
            <div class="flex items-center space-x-4 mb-3">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ task.title }}</h2>
              <n-tag :type="getTaskStatusColor(task.status)" size="small">
                {{ getTaskStatusText(task.status) }}
              </n-tag>
              <n-tag :type="getTaskTypeColor(task.task_type)" size="small">
                {{ getTaskTypeText(task.task_type) }}
              </n-tag>
            </div>
            
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span class="text-gray-500 dark:text-gray-400">ID:</span>
                <span class="ml-1 text-gray-900 dark:text-white">{{ task.id }}</span>
              </div>
              <div>
                <span class="text-gray-500 dark:text-gray-400">总项目:</span>
                <span class="ml-1 text-gray-900 dark:text-white">{{ task.total_items || 0 }}</span>
              </div>
              <div>
                <span class="text-gray-500 dark:text-gray-400">已处理:</span>
                <span class="ml-1 text-blue-600 dark:text-blue-400 font-medium">{{ task.processed_items || 0 }}</span>
              </div>
              <div>
                <span class="text-gray-500 dark:text-gray-400">成功率:</span>
                <span class="ml-1 text-green-600 dark:text-green-400 font-medium">
                  {{ task.total_items > 0 ? Math.round((task.success_items / task.total_items) * 100) : 0 }}%
                </span>
              </div>
            </div>
            
            <!-- 进度条 -->
            <div class="mt-3" v-if="task.total_items > 0">
              <div class="flex items-center justify-between mb-1">
                <span class="text-xs text-gray-500 dark:text-gray-400">进度</span>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ Math.round((task.processed_items / task.total_items) * 100) }}%
                </span>
              </div>
              <n-progress
                type="line"
                :percentage="Math.round((task.processed_items / task.total_items) * 100)"
                :height="6"
                :show-indicator="false"
              />
            </div>
          </div>
          
          <!-- 右侧：统计信息 -->
          <div class="flex items-center space-x-6">
            <div class="text-center">
              <div class="text-lg font-bold text-green-600 dark:text-green-400">{{ task.success_items || 0 }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">成功</div>
            </div>
            <div class="text-center">
              <div class="text-lg font-bold text-red-600 dark:text-red-400">{{ task.failed_items || 0 }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">失败</div>
            </div>
            <div class="text-center">
              <div class="text-lg font-bold text-gray-600 dark:text-gray-400">{{ task.total_items - task.processed_items || 0 }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">待处理</div>
            </div>
          </div>
        </div>
        
        <!-- 任务描述（如果有） -->
        <div v-if="task.description" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ task.description }}</p>
        </div>
      </div>

      <!-- 任务项列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">任务项列表</h2>
          <span class="text-sm text-gray-500 dark:text-gray-400">共 {{ total }} 项</span>
        </div>
        
        <div class="overflow-x-auto">
          <n-data-table
            :columns="taskItemColumns"
            :data="taskItems"
            :loading="itemsLoading"
            :pagination="itemsPaginationConfig"
            size="small"
            :scroll-x="600"
            :bordered="false"
          />
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="text-center py-4">
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
    width: 60,
    minWidth: 60
  },
  {
    title: '输入数据',
    key: 'input',
    minWidth: 250,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      if (!row.input) return h('span', { class: 'text-sm text-gray-500' }, '无输入数据')
      const title = row.input.title || row.input.url || '无标题'
      return h('div', { class: 'text-sm' }, [
        h('div', { class: 'font-medium' }, title),
        h('div', { class: 'text-xs text-gray-500 mt-1' }, row.input.url || '')
      ])
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    minWidth: 80,
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
    title: '结果',
    key: 'output',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      if (!row.output) return h('span', { class: 'text-sm text-gray-500' }, '无输出')
      if (row.output.error) {
        return h('div', { class: 'text-sm' }, [
          h('div', { class: 'text-red-600 font-medium' }, '失败'),
          h('div', { class: 'text-xs text-gray-500 mt-1' }, row.output.error)
        ])
      }
      return h('div', { class: 'text-sm' }, [
        h('div', { class: 'text-green-600 font-medium' }, '成功'),
        h('div', { class: 'text-xs text-gray-500 mt-1' }, row.output.save_url || '处理完成')
      ])
    }
  },
  {
    title: '时间',
    key: 'created_at',
    width: 120,
    minWidth: 120,
    render: (row: any) => {
      return h('div', { class: 'text-sm' }, [
        h('div', new Date(row.created_at).toLocaleDateString('zh-CN')),
        h('div', { class: 'text-xs text-gray-500' }, new Date(row.created_at).toLocaleTimeString('zh-CN'))
      ])
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
