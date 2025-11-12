<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 px-4 py-8 sm:px-6 sm:py-12">
    <div class="w-full max-w-sm sm:max-w-md">
      <div class="bg-white dark:bg-gray-800 p-6 sm:p-8 rounded-2xl shadow-2xl w-full text-gray-900 dark:text-gray-100 border border-gray-100 dark:border-gray-700">
        <!-- Logo/标题区域 -->
        <div class="text-center mb-6 sm:mb-6">
          <div class="mb-4 sm:mb-4">
            <i class="fas fa-user-circle text-4xl sm:text-5xl text-blue-500 dark:text-blue-400"></i>
          </div>
          <h1 class="text-xl sm:text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">欢迎回来</h1>
          <p class="text-sm sm:text-base text-gray-600 dark:text-gray-400">请输入您的登录信息</p>
        </div>

        <!-- 登录表单 -->
        <form @submit.prevent="handleLogin" class="space-y-4 sm:space-y-5">
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

          <button 
            type="submit" 
            :disabled="userStore.loading"
            class="w-full flex justify-center py-3 px-6 border border-transparent rounded-xl shadow-lg text-sm font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
          >
            <span v-if="userStore.loading" class="inline-flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登录中...
            </span>
            <span v-else class="inline-flex items-center">
              <i class="fas fa-sign-in-alt mr-2"></i>
              登录
            </span>
          </button>
        </form>
        
        <!-- 底部链接 -->
        <div class="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4 text-xs sm:text-sm">
            <NuxtLink to="/register" class="inline-flex items-center text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors font-medium">
              <i class="fas fa-user-plus mr-2"></i> 注册新账号
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
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()

const form = reactive({
  username: '',
  password: ''
})

const errors = reactive({
  username: '',
  password: ''
})


const validateForm = () => {
  errors.username = ''
  errors.password = ''
  
  console.log('validateForm - username:', form.username)
  console.log('validateForm - password:', form.password ? '***' : 'empty')
  
  if (!form.username || !form.username.trim()) {
    errors.username = '请输入用户名'
    return false
  }
  
  if (!form.password || !form.password.trim()) {
    errors.password = '请输入密码'
    return false
  }
  
  return true
}

const handleLogin = async () => {
  console.log('handleLogin - 开始登录，表单数据:', {
    username: form.username,
    password: form.password ? '***' : 'empty'
  })
  
  if (!validateForm()) {
    console.log('handleLogin - 表单验证失败')
    return
  }
  
  console.log('handleLogin - 表单验证通过，开始调用登录API')
  
  const result = await userStore.login({
    username: form.username,
    password: form.password
  })
  
  console.log('handleLogin - 登录结果:', result)
  
  if (result && result.success) {
    notification.success({
        content: '登录成功',
        duration: 3000
      })
    
    // 根据用户角色跳转到不同页面
    if (userStore.user?.role === 'admin') {
      await router.push('/admin')
    } else {
      await router.push('/user')
    }
  } else {
    // 根据错误类型提供更友好的提示
    let message = '登录失败'
    if (result.message) {
      if (result.message.includes('用户名或密码错误')) {
        message = '用户名或密码错误，请检查后重试'
      } else if (result.message.includes('账户已被禁用')) {
        message = '账户已被禁用，请联系管理员'
      } else if (result.message.includes('网络连接')) {
        message = '网络连接失败，请检查网络后重试'
      } else {
        message = result.message
      }
    }
    notification.error({
        content: message,
        duration: 3000
      })
  }
}

definePageMeta({
  layout: 'single',
  ssr: false
})

// 设置页面SEO
const { initSystemConfig, setLoginSeo } = useGlobalSeo()

onBeforeMount(async () => {
  await initSystemConfig()
  setLoginSeo()
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
}

/* 平板端优化 */
@media (min-width: 641px) and (max-width: 1024px) {
  .bg-white {
    margin: 1rem;
    padding: 2rem;
  }
}
</style> 