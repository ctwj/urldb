<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">任务管理</h1>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">查看和管理系统中的所有任务</p>
      </div>
    </div>

    <!-- 任务统计卡片 -->
    <div class="grid grid-cols-6 gap-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">总任务数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ taskStore.taskStats.total }}</p>
          </div>
          <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-tasks text-blue-600 dark:text-blue-400"></i>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">运行中</p>
            <p class="text-2xl font-bold text-orange-600">{{ taskStore.taskStats.running }}</p>
          </div>
          <div class="w-8 h-8 bg-orange-100 dark:bg-orange-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-play text-orange-600 dark:text-orange-400"></i>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">待处理</p>
            <p class="text-2xl font-bold text-yellow-600">{{ taskStore.taskStats.pending }}</p>
          </div>
          <div class="w-8 h-8 bg-yellow-100 dark:bg-yellow-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-clock text-yellow-600 dark:text-yellow-400"></i>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">已完成</p>
            <p class="text-2xl font-bold text-green-600">{{ taskStore.taskStats.completed }}</p>
          </div>
          <div class="w-8 h-8 bg-green-100 dark:bg-green-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-check text-green-600 dark:text-green-400"></i>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">失败</p>
            <p class="text-2xl font-bold text-red-600">{{ taskStore.taskStats.failed }}</p>
          </div>
          <div class="w-8 h-8 bg-red-100 dark:bg-red-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-times text-red-600 dark:text-red-400"></i>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-400">暂停</p>
            <p class="text-2xl font-bold text-gray-600">{{ taskStore.taskStats.paused }}</p>
          </div>
          <div class="w-8 h-8 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center">
            <i class="fas fa-pause text-gray-600 dark:text-gray-400"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 任务列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">任务列表</h2>
      </div>
      
      <div class="p-6">
        <n-data-table
          :columns="taskColumns"
          :data="tasks"
          :loading="loading"
          :pagination="paginationConfig"
          :row-class-name="getRowClassName"
          size="small"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { useTaskStore } from '~/stores/task'
import { useMessage, useDialog } from 'naive-ui'

// 任务状态管理
const taskStore = useTaskStore()
const message = useMessage()
const dialog = useDialog()

// 数据状态
const tasks = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 分页配置
const paginationConfig = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchTasks()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchTasks()
  }
}))

// 表格列定义
const taskColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    sorter: true
  },
  {
    title: '任务标题',
    key: 'title',
    minWidth: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '类型',
    key: 'task_type',
    width: 100,
    render: (row: any) => {
      const typeMap: Record<string, { text: string; color: string }> = {
        transfer: { text: '转存', color: 'blue' }
      }
      const type = typeMap[row.task_type] || { text: row.task_type, color: 'gray' }
      return h('n-tag', { type: type.color, size: 'small' }, { default: () => type.text })
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    render: (row: any) => {
      const statusMap: Record<string, { text: string; color: string }> = {
        pending: { text: '待处理', color: 'warning' },
        running: { text: '运行中', color: 'info' },
        completed: { text: '已完成', color: 'success' },
        failed: { text: '失败', color: 'error' },
        paused: { text: '暂停', color: 'default' }
      }
      
      // 优先使用 is_running 状态
      let currentStatus = row.status
      if (row.is_running) {
        currentStatus = 'running'
      }
      
      const status = statusMap[currentStatus] || { text: currentStatus, color: 'default' }
      return h('n-tag', { type: status.color, size: 'small' }, { default: () => status.text })
    }
  },
  {
    title: '进度',
    key: 'progress',
    width: 120,
    render: (row: any) => {
      const total = row.total_items || 0
      const processed = (row.processed_items || 0)
      const percentage = total > 0 ? Math.round((processed / total) * 100) : 0
      
      return h('div', { class: 'flex items-center space-x-2' }, [
        h('span', { class: 'text-sm' }, `${processed}/${total}`),
        h('n-progress', { 
          type: 'line', 
          percentage, 
          height: 4,
          showIndicator: false,
          style: { width: '60px' }
        })
      ])
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render: (row: any) => {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render: (row: any) => {
      const buttons = []
      
      if (row.status === 'pending' || (row.status !== 'running' && !row.is_running)) {
        buttons.push(
          h('n-button', {
            size: 'small',
            type: 'primary',
            onClick: () => startTask(row.id)
          }, { default: () => '启动' })
        )
      }
      
      if (row.is_running) {
        buttons.push(
          h('n-button', {
            size: 'small',
            type: 'warning',
            onClick: () => stopTask(row.id)
          }, { default: () => '停止' })
        )
      }
      
      if (row.status === 'completed' || row.status === 'failed') {
        buttons.push(
          h('n-button', {
            size: 'small',
            type: 'error',
            onClick: () => deleteTask(row.id)
          }, { default: () => '删除' })
        )
      }
      
      return h('div', { class: 'flex space-x-2' }, buttons)
    }
  }
]

// 行样式
const getRowClassName = (row: any) => {
  if (row.is_running) {
    return 'bg-blue-50 dark:bg-blue-900/10'
  }
  return ''
}

// 获取任务列表
const fetchTasks = async () => {
  loading.value = true
  try {
    const { useTaskApi } = await import('~/composables/useApi')
    const taskApi = useTaskApi()
    
    const response = await taskApi.getTasks({
      page: currentPage.value,
      page_size: pageSize.value
    }) as any
    
    if (response && response.items) {
      tasks.value = response.items
      total.value = response.total || 0
    }
  } catch (error) {
    console.error('获取任务列表失败:', error)
    message.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

// 启动任务
const startTask = async (taskId: number) => {
  try {
    const success = await taskStore.startTask(taskId)
    if (success) {
      message.success('任务启动成功')
      await fetchTasks()
    } else {
      message.error('任务启动失败')
    }
  } catch (error) {
    console.error('启动任务失败:', error)
    message.error('启动任务失败')
  }
}

// 停止任务
const stopTask = async (taskId: number) => {
  try {
    const success = await taskStore.stopTask(taskId)
    if (success) {
      message.success('任务停止成功')
      await fetchTasks()
    } else {
      message.error('任务停止失败')
    }
  } catch (error) {
    console.error('停止任务失败:', error)
    message.error('停止任务失败')
  }
}

// 删除任务
const deleteTask = async (taskId: number) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个任务吗？此操作不可逆。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const success = await taskStore.deleteTask(taskId)
        if (success) {
          message.success('任务删除成功')
          await fetchTasks()
        } else {
          message.error('任务删除失败')
        }
      } catch (error) {
        console.error('删除任务失败:', error)
        message.error('删除任务失败')
      }
    }
  })
}

// 初始化
onMounted(() => {
  fetchTasks()
  // 确保任务状态管理已启动（因为页面可能直接访问，而不是通过layout）
  taskStore.startAutoUpdate()
})

// 设置页面meta
definePageMeta({
  layout: 'admin'
})
</script>

<style scoped>
:deep(.n-data-table-th) {
  background-color: var(--n-th-color);
}

:deep(.bg-blue-50) {
  background-color: rgb(239 246 255);
}

:deep(.dark .bg-blue-50) {
  background-color: rgb(30 58 138 / 0.1);
}
</style>