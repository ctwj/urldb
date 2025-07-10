<template>
  <div class="min-h-screen bg-gray-50 p-8">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">登录功能测试</h1>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- 测试登录 -->
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">测试登录</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">用户名</label>
              <input 
                v-model="loginForm.username" 
                type="text" 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="admin"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">密码</label>
              <input 
                v-model="loginForm.password" 
                type="password" 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="password"
              >
            </div>
            
            <button 
              @click="testLogin"
              class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700"
            >
              测试登录
            </button>
          </div>
          
          <div v-if="loginResult" class="mt-4 p-4 bg-gray-100 rounded">
            <h3 class="font-semibold">登录结果:</h3>
            <pre class="text-sm mt-2">{{ JSON.stringify(loginResult, null, 2) }}</pre>
          </div>
        </div>
        
        <!-- 测试注册 -->
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-semibold mb-4">测试注册</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">用户名</label>
              <input 
                v-model="registerForm.username" 
                type="text" 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="testuser"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">邮箱</label>
              <input 
                v-model="registerForm.email" 
                type="email" 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="test@example.com"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">密码</label>
              <input 
                v-model="registerForm.password" 
                type="password" 
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
                placeholder="123456"
              >
            </div>
            
            <button 
              @click="testRegister"
              class="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700"
            >
              测试注册
            </button>
          </div>
          
          <div v-if="registerResult" class="mt-4 p-4 bg-gray-100 rounded">
            <h3 class="font-semibold">注册结果:</h3>
            <pre class="text-sm mt-2">{{ JSON.stringify(registerResult, null, 2) }}</pre>
          </div>
        </div>
      </div>
      
      <!-- 默认账户信息 -->
      <div class="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h2 class="text-xl font-semibold text-blue-900 mb-4">默认账户信息</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <h3 class="font-semibold text-blue-800">管理员账户</h3>
            <p class="text-blue-700">用户名: admin</p>
            <p class="text-blue-700">密码: password</p>
            <p class="text-blue-700">角色: admin</p>
          </div>
          <div>
            <h3 class="font-semibold text-blue-800">测试说明</h3>
            <ul class="text-blue-700 text-sm space-y-1">
              <li>• 使用默认账户可以正常登录</li>
              <li>• 注册新用户后可以登录</li>
              <li>• 错误的用户名密码会返回401错误</li>
              <li>• 重复的用户名会返回400错误</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'

const loginForm = reactive({
  username: 'admin',
  password: 'password'
})

const registerForm = reactive({
  username: 'testuser',
  email: 'test@example.com',
  password: '123456'
})

const loginResult = ref(null)
const registerResult = ref(null)

const testLogin = async () => {
  try {
    const authApi = useAuthApi()
    const response = await authApi.login(loginForm)
    loginResult.value = { success: true, data: response }
  } catch (error: any) {
    loginResult.value = { 
      success: false, 
      error: error.data?.error || error.message || '登录失败' 
    }
  }
}

const testRegister = async () => {
  try {
    const authApi = useAuthApi()
    const response = await authApi.register(registerForm)
    registerResult.value = { success: true, data: response }
  } catch (error: any) {
    registerResult.value = { 
      success: false, 
      error: error.data?.error || error.message || '注册失败' 
    }
  }
}

// 设置页面标题
useHead({
  title: '登录功能测试 - 网盘资源管理系统'
})
</script> 