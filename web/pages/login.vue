<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div class="max-w-md w-full">
      <div class="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-full max-w-md text-gray-900 dark:text-gray-100">
        <div class="text-center">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">管理员登录</h1>
          <p class="mt-2 text-sm text-gray-600">请输入管理员账号密码</p>
            <!-- <div class="mt-3 p-3 bg-blue-50 rounded-lg">  
            </div> -->
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-100">用户名</label>
            <n-input 
              type="text" 
              id="username" 
              v-model:value="form.username"
              required 
              :class="{ 'border-red-500': errors.username }"
            />
            <p v-if="errors.username" class="mt-1 text-sm text-red-600">{{ errors.username }}</p>
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-100">密码</label>
            <n-input 
              type="password" 
              id="password" 
              v-model:value="form.password"
              required 
              :class="{ 'border-red-500': errors.password }"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
          </div>

          <button 
            type="submit" 
            :disabled="userStore.loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="userStore.loading" class="inline-flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登录中...
            </span>
            <span v-else>登录</span>
          </button>
        </form>
        
        <div class="pt-2 text-center">
          <NuxtLink to="/register" class="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors mr-4">
            <i class="fas fa-user-plus mr-1"></i> 注册账号
          </NuxtLink>
          <NuxtLink to="/" class="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors">
            <i class="fas fa-home mr-1"></i> 返回首页
          </NuxtLink>
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
  
  if (result.success) {
    notification.success({
      content: '登录成功',
      duration: 3000
    })
    await router.push('/admin')
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

// 设置页面标题
useHead({
  title: '管理员登录 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 