<template>
  <div class="min-h-screen bg-gray-50 p-8">
    <div class="max-w-2xl mx-auto">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">Toast 组件测试</h1>
      
      <div class="bg-white rounded-lg shadow p-6 space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <button 
            @click="showErrorToast"
            class="bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700"
          >
            显示错误提示
          </button>
          
          <button 
            @click="showSuccessToast"
            class="bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700"
          >
            显示成功提示
          </button>
        </div>
        
        <div class="text-sm text-gray-600">
          <p>点击按钮测试不同的Toast提示效果</p>
        </div>
      </div>
    </div>
    
    <!-- 错误提示 -->
    <ErrorToast 
      v-if="showError" 
      :message="errorMessage" 
      @close="showError = false"
    />
    
    <!-- 成功提示 -->
    <SuccessToast 
      v-if="showSuccess" 
      :message="successMessage" 
      @close="showSuccess = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const showError = ref(false)
const showSuccess = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const showErrorToast = () => {
  errorMessage.value = '这是一个错误提示示例'
  showError.value = true
}

const showSuccessToast = () => {
  successMessage.value = '这是一个成功提示示例'
  showSuccess.value = true
}

// 设置页面标题
useHead({
  title: 'Toast 组件测试 - 网盘资源管理系统'
})
</script> 