<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div class="max-w-md w-full">
      <div class="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-full max-w-md text-gray-900 dark:text-gray-100">
        <div class="text-center">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">用户注册</h1>
          <p class="mt-2 text-sm text-gray-600">创建新的用户账户</p>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
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
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-100">邮箱</label>
            <n-input 
              type="email" 
              id="email" 
              v-model:value="form.email"
              required 
              :class="{ 'border-red-500': errors.email }"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
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

          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 dark:text-gray-100">确认密码</label>
            <n-input 
              type="password" 
              id="confirmPassword" 
              v-model:value="form.confirmPassword"
              required 
              :class="{ 'border-red-500': errors.confirmPassword }"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
          </div>

          <button 
            type="submit" 
            :disabled="userStore.loading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 dark:bg-blue-700 dark:hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span v-if="userStore.loading" class="inline-flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              注册中...
            </span>
            <span v-else>注册</span>
          </button>
        </form>
        
        <div class="pt-2 text-center">
          <NuxtLink to="/login" class="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors mr-4">
            <i class="fas fa-sign-in-alt mr-1"></i> 已有账号？登录
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
</style> 