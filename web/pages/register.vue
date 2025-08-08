<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100 dark:from-gray-900 dark:to-gray-800 px-4 py-8 sm:px-6 sm:py-12">
    <div class="w-full max-w-sm sm:max-w-md">
      <div class="bg-white dark:bg-gray-800 p-6 sm:p-8 rounded-2xl shadow-2xl w-full text-gray-900 dark:text-gray-100 border border-gray-100 dark:border-gray-700">
        <!-- 加载状态 -->
        <div v-if="configLoading" class="text-center">
          <div class="mb-4 sm:mb-6">
            <i class="fas fa-spinner fa-spin text-4xl sm:text-6xl text-blue-500 dark:text-blue-400"></i>
          </div>
          <h1 class="text-xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">加载中...</h1>
          <p class="text-sm sm:text-lg text-gray-600 dark:text-gray-400">正在检查系统配置</p>
        </div>

        <!-- 注册功能关闭时的提示 -->
        <div v-else-if="!enableRegister" class="text-center">
          <div class="mb-4 sm:mb-6">
            <i class="fas fa-ban text-4xl sm:text-6xl text-red-500 dark:text-red-400"></i>
          </div>
          <h1 class="text-xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2">注册功能已关闭</h1>
          <p class="text-sm sm:text-lg text-gray-600 dark:text-gray-400 mb-6 sm:mb-8">当前系统已关闭注册功能</p>
          <div class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4">
            <NuxtLink to="/login" class="w-full sm:w-auto inline-flex items-center justify-center px-4 sm:px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium text-sm">
              <i class="fas fa-sign-in-alt mr-2"></i> 已有账号？登录
            </NuxtLink>
            <NuxtLink to="/" class="w-full sm:w-auto inline-flex items-center justify-center px-4 sm:px-6 py-3 bg-gray-600 hover:bg-gray-700 text-white rounded-lg transition-colors font-medium text-sm">
              <i class="fas fa-home mr-2"></i> 返回首页
            </NuxtLink>
          </div>
        </div>

        <!-- 注册表单 -->
        <div v-else>
          <!-- Logo/标题区域 -->
          <div class="text-center mb-6">
            <div class="mb-4">
              <i class="fas fa-user-plus text-4xl sm:text-5xl text-green-500 dark:text-green-400"></i>
            </div>
            <h1 class="text-xl sm:text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">创建新账户</h1>
            <p class="text-sm sm:text-base text-gray-600 dark:text-gray-400">加入我们的社区</p>
          </div>

          <form @submit.prevent="handleRegister" class="space-y-4 sm:space-y-5">
            <div class="space-y-2">
              <label for="username" class="block text-sm font-semibold text-gray-700 dark:text-gray-300">用户名</label>
              <n-input 
                type="text" 
                id="username" 
                v-model:value="form.username"
                required 
                placeholder="请输入用户名"
                :class="{ 'border-red-500': errors.username }"
              />
              <p v-if="errors.username" class="mt-1 text-xs text-red-600 dark:text-red-400">{{ errors.username }}</p>
            </div>

            <div class="space-y-2">
              <label for="email" class="block text-sm font-semibold text-gray-700 dark:text-gray-300">邮箱</label>
              <n-input 
                id="email" 
                v-model:value="form.email"
                required 
                placeholder="请输入邮箱地址"
                :class="{ 'border-red-500': errors.email }"
              />
              <p v-if="errors.email" class="mt-1 text-xs text-red-600 dark:text-red-400">{{ errors.email }}</p>
            </div>

            <div class="space-y-2">
              <label for="password" class="block text-sm font-semibold text-gray-700 dark:text-gray-300">密码</label>
              <n-input 
                type="password" 
                id="password" 
                v-model:value="form.password"
                required 
                placeholder="请输入密码"
                :class="{ 'border-red-500': errors.password }"
              />
              <p v-if="errors.password" class="mt-1 text-xs text-red-600 dark:text-red-400">{{ errors.password }}</p>
            </div>

            <div class="space-y-2">
              <label for="confirmPassword" class="block text-sm font-semibold text-gray-700 dark:text-gray-300">确认密码</label>
              <n-input 
                type="password" 
                id="confirmPassword" 
                v-model:value="form.confirmPassword"
                required 
                placeholder="请再次输入密码"
                :class="{ 'border-red-500': errors.confirmPassword }"
              />
              <p v-if="errors.confirmPassword" class="mt-1 text-xs text-red-600 dark:text-red-400">{{ errors.confirmPassword }}</p>
            </div>

            <button 
              type="submit" 
              :disabled="userStore.loading"
              class="w-full flex justify-center py-3 px-6 border border-transparent rounded-xl shadow-lg text-sm font-semibold text-white bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
            >
              <span v-if="userStore.loading" class="inline-flex items-center">
                <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                注册中...
              </span>
              <span v-else class="inline-flex items-center">
                <i class="fas fa-user-plus mr-2"></i>
                创建账户
              </span>
            </button>
          </form>
          
          <!-- 底部链接 -->
          <div class="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
            <div class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4 text-xs sm:text-sm">
              <NuxtLink to="/login" class="inline-flex items-center text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors font-medium">
                <i class="fas fa-sign-in-alt mr-2"></i> 已有账号？登录
              </NuxtLink>
              <span class="text-gray-400 dark:text-gray-500 hidden sm:inline">|</span>
              <NuxtLink to="/" class="inline-flex items-center text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-300 transition-colors font-medium">
                <i class="fas fa-home mr-2"></i> 返回首页
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()
const systemConfigStore = useSystemConfigStore()

// 注册功能开关
const enableRegister = ref(true)
const configLoading = ref(true) // 添加加载状态

// 获取系统配置
const fetchSystemConfig = async () => {
  try {
    configLoading.value = true // 开始加载
    console.log('开始获取系统配置...') // 调试信息
    
    // 使用systemConfig store
    await systemConfigStore.initConfig()
    
    console.log('系统配置响应:', systemConfigStore.config) // 调试信息
    console.log('系统配置响应类型:', typeof systemConfigStore.config) // 调试信息
    console.log('系统配置响应键:', Object.keys(systemConfigStore.config || {})) // 调试信息
    
    if (systemConfigStore.config) {
      // 检查enable_register字段
      const enableRegisterValue = systemConfigStore.config.enable_register
      console.log('enable_register值:', enableRegisterValue) // 调试信息
      console.log('enable_register类型:', typeof enableRegisterValue) // 调试信息
      
      // 如果enable_register为false，则关闭注册
      enableRegister.value = enableRegisterValue !== false
      console.log('最终enableRegister值:', enableRegister.value) // 调试信息
    } else {
      console.log('未获取到系统配置数据') // 调试信息
      // 如果获取失败，默认允许注册
      enableRegister.value = true
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
    // 如果获取失败，默认允许注册
    enableRegister.value = true
  } finally {
    configLoading.value = false // 结束加载
  }
}

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateForm = () => {
  errors.username = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''
  
  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    return false
  }
  
  if (form.username.length < 3) {
    errors.username = '用户名至少需要3个字符'
    return false
  }
  
  if (!form.email.trim()) {
    errors.email = '请输入邮箱'
    return false
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    errors.email = '请输入有效的邮箱地址'
    return false
  }
  
  if (!form.password.trim()) {
    errors.password = '请输入密码'
    return false
  }
  
  if (form.password.length < 6) {
    errors.password = '密码至少需要6个字符'
    return false
  }
  
  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    return false
  }
  
  return true
}

const handleRegister = async () => {
  if (!validateForm()) return
  
  const result = await userStore.register({
    username: form.username,
    email: form.email,
    password: form.password
  })
  
  if (result.success) {
    notification.success({
      content: '注册成功！请登录',
      duration: 3000
    })
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } else {
    // 根据错误类型提供更友好的提示
    let errorMessage = '注册失败'
    if (result.message) {
      if (result.message.includes('用户名已存在')) {
        errorMessage = '用户名已存在，请选择其他用户名'
      } else if (result.message.includes('邮箱已存在')) {
        errorMessage = '邮箱已被注册，请使用其他邮箱'
      } else if (result.message.includes('网络连接')) {
        errorMessage = '网络连接失败，请检查网络后重试'
      } else {
        errorMessage = result.message
      }
    }
    notification.error({
      content: errorMessage,
      duration: 3000
    })
  }
}

definePageMeta({
  layout: 'single',
  ssr: false
})

// 页面加载时获取系统配置
onMounted(() => {
  fetchSystemConfig()
})


// 设置页面标题
useHead({
  title: '用户注册 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}

/* 添加卡片悬停效果 */
.bg-white {
  transition: all 0.3s ease;
}

.bg-white:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

/* 输入框焦点效果 */
.n-input {
  transition: all 0.2s ease;
}

.n-input:focus-within {
  transform: scale(1.02);
}

/* 按钮悬停效果 */
button[type="submit"]:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

/* 链接悬停效果 */
a {
  transition: all 0.2s ease;
}

a:hover {
  transform: translateY(-1px);
}

/* 加载动画效果 */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.fa-spinner {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

/* 手机端优化 */
@media (max-width: 640px) {
  .bg-white {
    margin: 0.5rem;
    padding: 1.5rem;
  }
  
  /* 手机端按钮优化 */
  button[type="submit"] {
    min-height: 44px; /* 确保触摸目标足够大 */
  }
  
  /* 手机端输入框优化 */
  .n-input {
    min-height: 44px; /* 确保触摸目标足够大 */
  }
  
  /* 手机端链接优化 */
  a {
    min-height: 44px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  /* 手机端按钮组优化 */
  .flex.flex-col.sm\:flex-row {
    gap: 0.75rem;
  }
  
  /* 手机端按钮全宽 */
  .w-full.sm\:w-auto {
    width: 100%;
  }
}

/* 平板端优化 */
@media (min-width: 641px) and (max-width: 1024px) {
  .bg-white {
    margin: 1rem;
    padding: 2rem;
  }
}
</style> 