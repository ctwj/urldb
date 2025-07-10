<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 p-8">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-2xl font-bold mb-6">登录状态测试</h1>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
        <h2 class="text-lg font-semibold mb-4">当前状态</h2>
        <div class="space-y-2">
          <div class="flex justify-between">
            <span>是否已认证:</span>
            <span :class="userStore.isAuthenticated ? 'text-green-600' : 'text-red-600'">
              {{ userStore.isAuthenticated ? '是' : '否' }}
            </span>
          </div>
          <div class="flex justify-between">
            <span>用户名:</span>
            <span>{{ userStore.userInfo?.username || '未登录' }}</span>
          </div>
          <div class="flex justify-between">
            <span>角色:</span>
            <span>{{ userStore.userInfo?.role || '未登录' }}</span>
          </div>
          <div class="flex justify-between">
            <span>Token:</span>
            <span>{{ userStore.token ? '存在' : '不存在' }}</span>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
        <h2 class="text-lg font-semibold mb-4">localStorage 状态</h2>
        <div class="space-y-2">
          <div class="flex justify-between">
            <span>Token:</span>
            <span>{{ localStorageToken ? '存在' : '不存在' }}</span>
          </div>
          <div class="flex justify-between">
            <span>用户信息:</span>
            <span>{{ localStorageUser ? '存在' : '不存在' }}</span>
          </div>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
        <h2 class="text-lg font-semibold mb-4">操作</h2>
        <div class="space-y-4">
          <button 
            @click="initAuth"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            重新初始化认证状态
          </button>
          <button 
            @click="clearStorage"
            class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 ml-2"
          >
            清除localStorage
          </button>
          <button 
            @click="refreshPage"
            class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 ml-2"
          >
            刷新页面
          </button>
        </div>
      </div>
      
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold mb-4">调试信息</h2>
        <pre class="bg-gray-100 dark:bg-gray-700 p-4 rounded text-sm overflow-auto">{{ debugInfo }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
const userStore = useUserStore()

const localStorageToken = ref('')
const localStorageUser = ref('')
const debugInfo = ref('')

// 检查localStorage
const checkLocalStorage = () => {
  if (typeof window !== 'undefined') {
    localStorageToken.value = localStorage.getItem('token') ? '存在' : '不存在'
    localStorageUser.value = localStorage.getItem('user') ? '存在' : '不存在'
  }
}

// 初始化认证状态
const initAuth = () => {
  userStore.initAuth()
  checkLocalStorage()
  updateDebugInfo()
}

// 清除localStorage
const clearStorage = () => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    checkLocalStorage()
    updateDebugInfo()
  }
}

// 刷新页面
const refreshPage = () => {
  window.location.reload()
}

// 更新调试信息
const updateDebugInfo = () => {
  debugInfo.value = JSON.stringify({
    store: {
      isAuthenticated: userStore.isAuthenticated,
      user: userStore.userInfo,
      token: userStore.token ? '存在' : '不存在'
    },
    localStorage: {
      token: localStorageToken.value,
      user: localStorageUser.value
    }
  }, null, 2)
}

// 页面加载时检查状态
onMounted(() => {
  userStore.initAuth()
  checkLocalStorage()
  updateDebugInfo()
})
</script> 