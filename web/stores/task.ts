import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useTaskApi } from '~/composables/useApi'

export interface TaskStats {
  total: number
  running: number
  pending: number
  completed: number
  failed: number
  paused: number
}

export interface TaskInfo {
  id: number
  title: string
  type: string
  status: string
  total_items: number
  processed_items?: number
  success_items?: number
  failed_items?: number
  created_at: string
  updated_at: string
  is_running?: boolean  // 任务是否在TaskManager中运行
}

export const useTaskStore = defineStore('task', () => {
  const taskApi = useTaskApi()
  
  // 任务统计信息
  const taskStats = ref<TaskStats>({
    total: 0,
    running: 0,
    pending: 0,
    completed: 0,
    failed: 0,
    paused: 0
  })
  
  // 正在运行的任务列表
  const runningTasks = ref<TaskInfo[]>([])
  
  // 未完成的任务列表（pending + running + paused）
  const incompleteTasks = ref<TaskInfo[]>([])
  
  // 更新定时器
  let updateInterval: NodeJS.Timeout | null = null
  
  // 计算属性：是否有活跃任务
  const hasActiveTasks = computed(() => {
    return taskStats.value.running > 0 || taskStats.value.pending > 0 || taskStats.value.paused > 0
  })
  
  // 计算属性：活跃任务总数
  const activeTaskCount = computed(() => {
    return taskStats.value.running + taskStats.value.pending + taskStats.value.paused
  })
  
  // 计算属性：正在运行的任务数
  const runningTaskCount = computed(() => {
    return taskStats.value.running
  })
  
  // 获取任务统计信息
  const fetchTaskStats = async () => {
    try {
      const response = await taskApi.getTasks() as any
      console.log('原始任务API响应:', response)
      
      // 处理API响应格式
      let tasks: TaskInfo[] = []
      if (response && response.items && Array.isArray(response.items)) {
        tasks = response.items
      } else if (Array.isArray(response)) {
        tasks = response
      }
      
      console.log('解析后的任务列表:', tasks)
      
      if (tasks && tasks.length >= 0) {
        // 重置统计
        const stats: TaskStats = {
          total: tasks.length,
          running: 0,
          pending: 0,
          completed: 0,
          failed: 0,
          paused: 0
        }
        
        const running: TaskInfo[] = []
        const incomplete: TaskInfo[] = []
        
        // 统计各种状态的任务
        tasks.forEach((task: TaskInfo) => {
          console.log('处理任务:', task.id, '状态:', task.status, '是否运行中:', task.is_running)
          
          // 如果任务标记为运行中，优先使用running状态
          let currentStatus = task.status
          if (task.is_running) {
            currentStatus = 'running'
          }
          
          switch (currentStatus) {
            case 'running':
              stats.running++
              running.push(task)
              incomplete.push(task)
              break
            case 'pending':
              stats.pending++
              incomplete.push(task)
              break
            case 'completed':
              stats.completed++
              break
            case 'failed':
              stats.failed++
              break
            case 'paused':
              stats.paused++
              incomplete.push(task)
              break
          }
        })
        
        // 更新状态
        taskStats.value = stats
        runningTasks.value = running
        incompleteTasks.value = incomplete
        
        console.log('任务统计更新:', stats)
        console.log('运行中的任务:', running)
        console.log('未完成的任务:', incomplete)
      }
    } catch (error) {
      console.error('获取任务统计失败:', error)
    }
  }
  
  // 开始定时更新
  const startAutoUpdate = () => {
    if (updateInterval) {
      clearInterval(updateInterval)
    }
    
    // 立即执行一次
    fetchTaskStats()
    
    // 每5秒更新一次
    updateInterval = setInterval(() => {
      fetchTaskStats()
    }, 5000)
    
    console.log('任务状态自动更新已启动')
  }
  
  // 停止定时更新
  const stopAutoUpdate = () => {
    if (updateInterval) {
      clearInterval(updateInterval)
      updateInterval = null
      console.log('任务状态自动更新已停止')
    }
  }
  
  // 获取特定任务的详细状态
  const getTaskStatus = async (taskId: number) => {
    try {
      const status = await taskApi.getTaskStatus(taskId)
      return status
    } catch (error) {
      console.error('获取任务状态失败:', error)
      return null
    }
  }
  
  // 启动任务
  const startTask = async (taskId: number) => {
    try {
      await taskApi.startTask(taskId)
      // 立即更新状态
      await fetchTaskStats()
      return true
    } catch (error) {
      console.error('启动任务失败:', error)
      return false
    }
  }
  
  // 停止任务
  const stopTask = async (taskId: number) => {
    try {
      await taskApi.stopTask(taskId)
      // 立即更新状态
      await fetchTaskStats()
      return true
    } catch (error) {
      console.error('停止任务失败:', error)
      return false
    }
  }
  
  // 删除任务
  const deleteTask = async (taskId: number) => {
    try {
      await taskApi.deleteTask(taskId)
      // 立即更新状态
      await fetchTaskStats()
      return true
    } catch (error) {
      console.error('删除任务失败:', error)
      return false
    }
  }
  
  return {
    // 状态
    taskStats,
    runningTasks,
    incompleteTasks,
    
    // 计算属性
    hasActiveTasks,
    activeTaskCount,
    runningTaskCount,
    
    // 方法
    fetchTaskStats,
    startAutoUpdate,
    stopAutoUpdate,
    getTaskStatus,
    startTask,
    stopTask,
    deleteTask
  }
})