<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">任务管理</h1>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">查看和管理系统中的所有任务</p>
      </div>
    </div>

    <!-- 任务列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
      <div class="p-4 md:p-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">任务列表</h2>
          
          <!-- 筛选条件 -->
          <div class="flex flex-col sm:flex-row items-start sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
            <div class="flex items-center space-x-2">
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">任务状态：</label>
              <n-select
                v-model:value="statusFilter"
                :options="statusOptions"
                placeholder="全部状态"
                style="width: 120px"
                size="small"
                clearable
                @update:value="onStatusFilterChange"
              />
            </div>
            
            <div class="flex items-center space-x-2">
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300 whitespace-nowrap">任务类型：</label>
              <n-select
                v-model:value="typeFilter"
                :options="typeOptions"
                placeholder="全部类型"
                style="width: 100px"
                size="small"
                clearable
                @update:value="onTypeFilterChange"
              />
            </div>
            
            <n-button
              type="primary"
              size="small"
              @click="refreshTasks"
            >
              <template #icon>
                <i class="fas fa-refresh"></i>
              </template>
              刷新
            </n-button>
          </div>
        </div>
      </div>
      
      <div class="p-4 md:p-6">
        <div class="overflow-x-auto">
          <n-data-table
            :columns="taskColumns"
            :data="tasks"
            :loading="loading"
            :pagination="paginationConfig"
            :row-class-name="getRowClassName"
            size="small"
            :scroll-x="800"
            class="min-w-full"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import { NButton } from 'naive-ui'
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

// 筛选条件
const statusFilter = ref(null)
const typeFilter = ref(null)

// 状态选项
const statusOptions = [
  { label: '待处理', value: 'pending' },
  { label: '运行中', value: 'running' },
  { label: '已完成', value: 'completed' },
  { label: '失败', value: 'failed' },
  { label: '暂停', value: 'paused' }
]

// 类型选项
const typeOptions = [
  { label: '转存任务', value: 'transfer' }
]

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
    width: 60,
    minWidth: 60,
    maxWidth: 60,
    sorter: true
  },
  {
    title: '任务标题',
    key: 'title',
    minWidth: 180,
    ellipsis: {
      tooltip: true
    },
    render: (row: any) => {
      return h('a', {
        href: `/admin/tasks/${row.id}`,
        class: 'text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 cursor-pointer hover:underline',
        onClick: (e: Event) => {
          e.preventDefault()
          navigateTo(`/admin/tasks/${row.id}`)
        }
      }, row.title)
    }
  },
  {
    title: '类型',
    key: 'task_type',
    width: 80,
    minWidth: 80,
    maxWidth: 80,
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
    width: 90,
    minWidth: 90,
    maxWidth: 90,
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
    width: 110,
    minWidth: 110,
    maxWidth: 110,
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
          style: { width: '50px' }
        })
      ])
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 140,
    minWidth: 140,
    maxWidth: 140,
    render: (row: any) => {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    minWidth: 120,
    maxWidth: 120,
    render: (row: any) => {
      return renderActions(row)
    }
  }
]

// 获取任务状态
const getTaskStatus = (row: any) => {
  // 优先使用数据库中的状态，因为它是最权威的
  // 只有当数据库状态是running且TaskManager也确认在运行时，才显示为running
  if (row.status === 'running' && row.is_running) {
    return 'running'
  }
  // 其他所有情况都返回数据库状态
  return row.status
}

// 渲染操作按钮
const renderActions = (row: any) => {
  const currentStatus = getTaskStatus(row)
  const buttons = []
  
  // 待处理状态：启动按钮
  if (currentStatus === 'pending') {
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'primary',
        quaternary: true,
        onClick: () => startTask(row.id)
      }, { default: () => '启动' })
    )
  }
  
  // 运行中状态：暂停和停止按钮
  if (currentStatus === 'running') {
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'warning',
        quaternary: true,
        onClick: () => pauseTask(row.id)
      }, { default: () => '暂停' })
    )
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'error',
        quaternary: true,
        onClick: () => stopTask(row.id)
      }, { default: () => '停止' })
    )
  }
  
  // 暂停状态：继续和停止按钮
  if (currentStatus === 'paused') {
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'primary',
        quaternary: true,
        onClick: () => resumeTask(row.id)
      }, { default: () => '继续' })
    )
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'error',
        quaternary: true,
        onClick: () => stopTask(row.id)
      }, { default: () => '停止' })
    )
  }
  
  // 失败状态：重试按钮
  if (currentStatus === 'failed') {
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'info',
        quaternary: true,
        onClick: () => retryTask(row.id)
      }, { default: () => '重试' })
    )
  }
  
  // 已完成或失败状态：删除按钮
  if (currentStatus === 'completed' || currentStatus === 'failed') {
    buttons.push(
      h('n-button', {
        size: 'small',
        type: 'error',
        quaternary: true,
        onClick: () => deleteTask(row.id)
      }, { default: () => '删除' })
    )
  }
  
  return h('div', { class: 'task-actions flex flex-wrap gap-1 justify-center' }, buttons)
}

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
    
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    // 添加筛选条件
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    if (typeFilter.value) {
      params.task_type = typeFilter.value
    }
    
    const response = await taskApi.getTasks(params) as any
    
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

// 筛选条件变化处理
const onStatusFilterChange = () => {
  currentPage.value = 1
  fetchTasks()
}

const onTypeFilterChange = () => {
  currentPage.value = 1
  fetchTasks()
}

// 刷新任务列表
const refreshTasks = async () => {
  // 强制刷新任务状态和列表
  await taskStore.fetchTaskStats()
  await fetchTasks()
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

// 暂停任务
const pauseTask = async (taskId: number) => {
  try {
    const success = await taskStore.pauseTask(taskId)
    if (success) {
      message.success('任务暂停成功')
      await fetchTasks()
    } else {
      message.error('任务暂停失败')
    }
  } catch (error) {
    console.error('暂停任务失败:', error)
    message.error('暂停任务失败')
  }
}

// 继续任务（恢复暂停的任务）
const resumeTask = async (taskId: number) => {
  try {
    const success = await taskStore.startTask(taskId)
    if (success) {
      message.success('任务继续成功')
      await fetchTasks()
    } else {
      message.error('任务继续失败')
    }
  } catch (error) {
    console.error('继续任务失败:', error)
    message.error('继续任务失败')
  }
}

// 重试失败的任务
const retryTask = async (taskId: number) => {
  try {
    const success = await taskStore.startTask(taskId)
    if (success) {
      message.success('任务重试成功')
      await fetchTasks()
    } else {
      message.error('任务重试失败')
    }
  } catch (error) {
    console.error('重试任务失败:', error)
    message.error('重试任务失败')
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

/* 表格自适应优化 */
:deep(.n-data-table) {
  min-width: 100%;
}

:deep(.n-data-table-wrapper) {
  overflow-x: auto;
}

/* 任务操作按钮 hover 效果 */
:deep(.task-actions .n-button) {
  transition: all 0.2s ease;
}

:deep(.task-actions .n-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 不同类型按钮的 hover 阴影颜色 */
:deep(.task-actions .n-button--primary-type:hover) {
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

:deep(.task-actions .n-button--warning-type:hover) {
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

:deep(.task-actions .n-button--error-type:hover) {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

:deep(.task-actions .n-button--info-type:hover) {
  box-shadow: 0 4px 12px rgba(14, 165, 233, 0.3);
}

:deep(.task-actions .n-button:active) {
  transform: translateY(0);
}

/* 响应式优化 */
@media (max-width: 768px) {
  :deep(.n-data-table) {
    font-size: 12px;
  }
  
  :deep(.n-data-table .n-button) {
    font-size: 11px;
    padding: 2px 6px;
  }
  
  :deep(.n-data-table .n-tag) {
    font-size: 11px;
  }
}
</style>