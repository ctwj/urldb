<template>
  <div class="version-info">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          <i class="fas fa-info-circle mr-2 text-blue-500"></i>
          版本信息
        </h3>
        <button 
          @click="refreshVersion"
          :disabled="loading"
          class="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
        >
          <i class="fas fa-sync-alt mr-1" :class="{ 'animate-spin': loading }"></i>
          刷新
        </button>
      </div>

      <div v-if="loading" class="text-center py-4">
        <i class="fas fa-spinner fa-spin text-blue-500 text-xl"></i>
        <p class="text-gray-600 dark:text-gray-400 mt-2">加载中...</p>
      </div>

      <div v-else-if="error" class="text-center py-4">
        <i class="fas fa-exclamation-triangle text-red-500 text-xl"></i>
        <p class="text-red-600 dark:text-red-400 mt-2">{{ error }}</p>
      </div>

      <div v-else class="space-y-4">
        <!-- 版本号 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">版本号</span>
          <span class="font-mono text-blue-600 dark:text-blue-400">v{{ versionInfo.version }}</span>
        </div>

        <!-- 构建时间 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">构建时间</span>
          <span class="text-gray-600 dark:text-gray-400">{{ formatTime(versionInfo.build_time) }}</span>
        </div>

        <!-- Git提交 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">Git提交</span>
          <span class="font-mono text-gray-600 dark:text-gray-400">{{ versionInfo.git_commit }}</span>
        </div>

        <!-- Git分支 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">Git分支</span>
          <span class="text-gray-600 dark:text-gray-400">{{ versionInfo.git_branch }}</span>
        </div>

        <!-- Go版本 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">Go版本</span>
          <span class="text-gray-600 dark:text-gray-400">{{ versionInfo.go_version }}</span>
        </div>

        <!-- 平台信息 -->
        <div class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded">
          <span class="text-gray-700 dark:text-gray-300">平台</span>
          <span class="text-gray-600 dark:text-gray-400">{{ versionInfo.platform }}/{{ versionInfo.arch }}</span>
        </div>

        <!-- 更新检查 -->
        <div class="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="font-medium text-blue-900 dark:text-blue-100">检查更新</h4>
              <p class="text-sm text-blue-700 dark:text-blue-300 mt-1">
                当前版本: v{{ versionInfo.version }}
              </p>
            </div>
            <button 
              @click="checkUpdate"
              :disabled="updateChecking"
              class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
            >
              <i class="fas fa-download mr-1" :class="{ 'animate-spin': updateChecking }"></i>
              检查更新
            </button>
          </div>
          
          <div v-if="updateInfo" class="mt-3 p-3 bg-white dark:bg-gray-800 rounded border">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-gray-400">最新版本</span>
              <span class="font-mono text-gray-900 dark:text-white">v{{ updateInfo.latest_version }}</span>
            </div>
            <div v-if="updateInfo.has_update" class="mt-2 p-2 bg-green-50 dark:bg-green-900/20 rounded border border-green-200 dark:border-green-800">
              <div class="flex items-center">
                <i class="fas fa-arrow-up text-green-600 dark:text-green-400 mr-2"></i>
                <span class="text-sm text-green-700 dark:text-green-300">有新版本可用</span>
              </div>
            </div>
            <div v-else class="mt-2 p-2 bg-gray-50 dark:bg-gray-700 rounded border">
              <div class="flex items-center">
                <i class="fas fa-check text-gray-600 dark:text-gray-400 mr-2"></i>
                <span class="text-sm text-gray-700 dark:text-gray-300">已是最新版本</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface VersionInfo {
  version: string
  build_time: string
  git_commit: string
  git_branch: string
  go_version: string
  node_version: string
  platform: string
  arch: string
}

interface UpdateInfo {
  current_version: string
  latest_version: string
  has_update: boolean
  update_available: boolean
}

const loading = ref(false)
const error = ref('')
const updateChecking = ref(false)
const versionInfo = ref<VersionInfo>({
  version: '',
  build_time: '',
  git_commit: '',
  git_branch: '',
  go_version: '',
  node_version: '',
  platform: '',
  arch: ''
})
const updateInfo = ref<UpdateInfo | null>(null)

// 获取版本信息
const fetchVersionInfo = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await $fetch('/api/version') as any
    if (response.success) {
      versionInfo.value = response.data
    } else {
      error.value = response.message || '获取版本信息失败'
    }
  } catch (err: any) {
    error.value = err.message || '网络错误'
  } finally {
    loading.value = false
  }
}

// 检查更新
const checkUpdate = async () => {
  updateChecking.value = true
  
  try {
    const response = await $fetch('/api/version/check-update') as any
    if (response.success) {
      updateInfo.value = response.data
    }
  } catch (err: any) {
    console.error('检查更新失败:', err)
  } finally {
    updateChecking.value = false
  }
}

// 刷新版本信息
const refreshVersion = () => {
  fetchVersionInfo()
  updateInfo.value = null
}

// 格式化时间
const formatTime = (timeStr: string) => {
  if (!timeStr) return 'N/A'
  try {
    const date = new Date(timeStr)
    return date.toLocaleString('zh-CN')
  } catch {
    return timeStr
  }
}

// 组件挂载时获取版本信息
onMounted(() => {
  fetchVersionInfo()
})
</script>

<style scoped>
.version-info {
  max-width: 600px;
  margin: 0 auto;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style> 