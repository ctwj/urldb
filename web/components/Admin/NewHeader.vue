<template>
  <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
    <div class="flex items-center justify-between px-6 py-4">
      <!-- 左侧：Logo和标题 -->
      <div class="flex items-center">
        <NuxtLink to="/admin" class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <i class="fas fa-shield-alt text-white text-sm"></i>
          </div>
          <div>
            <h1 class="text-xl font-bold text-gray-900 dark:text-white">管理后台</h1>
            <p class="text-xs text-gray-500 dark:text-gray-400">老九网盘资源数据库</p>
          </div>
        </NuxtLink>
      </div>

      <!-- 中间：状态信息 -->
      <div class="flex items-center space-x-6">
        <!-- 系统状态 -->
        <div class="flex items-center space-x-2">
          <div class="w-2 h-2 bg-green-500 rounded-full"></div>
          <span class="text-sm text-gray-600 dark:text-gray-300">系统正常</span>
        </div>
        
        <!-- 自动处理状态 -->
        <NuxtLink to="/admin/feature-config" class="flex items-center space-x-2 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 px-2 py-1 rounded transition-colors">
          <div :class="autoProcessEnabled ? 'w-2 h-2 bg-green-500 rounded-full' : 'w-2 h-2 bg-gray-400 rounded-full'"></div>
          <span class="text-sm text-gray-600 dark:text-gray-300">
            自动处理{{ autoProcessEnabled ? '已开启' : '已关闭' }}
          </span>
        </NuxtLink>
        
        <!-- 自动转存状态 -->
        <div class="flex items-center space-x-2">
          <div :class="autoTransferEnabled ? 'w-2 h-2 bg-blue-500 rounded-full' : 'w-2 h-2 bg-gray-400 rounded-full'"></div>
          <span class="text-sm text-gray-600 dark:text-gray-300">
            自动转存{{ autoTransferEnabled ? '已开启' : '已关闭' }}
          </span>
        </div>
        
        <!-- 任务状态 -->
        <div v-if="taskStore.hasActiveTasks" class="flex items-center space-x-2">
          <div class="w-2 h-2 bg-orange-500 rounded-full animate-pulse"></div>
          <span class="text-sm text-gray-600 dark:text-gray-300">
            <template v-if="taskStore.runningTaskCount > 0">
              {{ taskStore.runningTaskCount }}个任务运行中
            </template>
            <template v-else>
              {{ taskStore.activeTaskCount }}个任务待处理
            </template>
          </span>
        </div>
      </div>

      <!-- 右侧：用户菜单 -->
      <div class="flex items-center space-x-4">
        <NuxtLink to="/" class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
          <i class="fas fa-home text-lg"></i>
        </NuxtLink>
        <NuxtLink to="/admin-old" class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
          <i class="fas fa-arrow-left text-lg"></i>
        </NuxtLink>
        <div class="flex items-center space-x-2 cursor-pointer p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
          <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
            <i class="fas fa-user text-white text-sm"></i>
          </div>
          <div class="hidden md:block text-left">
            <p class="text-sm font-medium text-gray-900 dark:text-white">管理员</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">admin</p>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useTaskStore } from '~/stores/task'
import { useSystemConfigStore } from '~/stores/systemConfig'

// 任务状态管理
const taskStore = useTaskStore()

// 系统配置状态管理
const systemConfigStore = useSystemConfigStore()

// 自动处理和自动转存状态
const autoProcessEnabled = ref(false)
const autoTransferEnabled = ref(false)

// 获取系统配置状态
const fetchSystemStatus = async () => {
  try {
    await systemConfigStore.initConfig(false, true)
    
    // 从系统配置中获取自动处理和自动转存状态
    const config = systemConfigStore.config
    
    if (config) {
      // 检查自动处理状态
      autoProcessEnabled.value = config.auto_process_ready_resources === '1' || config.auto_process_ready_resources === true
      
      // 检查自动转存状态
      autoTransferEnabled.value = config.auto_transfer_enabled === '1' || config.auto_transfer_enabled === true
    }
    
  } catch (error) {
    console.error('获取系统状态失败:', error)
  }
}

// 组件挂载时启动
onMounted(() => {
  // 启动任务状态自动更新
  taskStore.startAutoUpdate()
  
  // 获取系统配置状态
  fetchSystemStatus()
  
  // 定期更新系统配置状态（每30秒）
  const configInterval = setInterval(fetchSystemStatus, 30000)
  
  // 保存定时器引用用于清理
  ;(globalThis as any).__configInterval = configInterval
})

// 组件销毁时清理
onBeforeUnmount(() => {
  // 停止任务状态自动更新
  taskStore.stopAutoUpdate()
  
  // 清理配置更新定时器
  if ((globalThis as any).__configInterval) {
    clearInterval((globalThis as any).__configInterval)
    delete (globalThis as any).__configInterval
  }
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 