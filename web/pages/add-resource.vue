<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 页面头部 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-4xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button 
              @click="$router.back()" 
              class="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100"
            >
              <i class="fas fa-arrow-left text-xl"></i>
            </button>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">添加资源</h1>
          </div>
          <div class="flex items-center space-x-2">
            <NuxtLink 
              to="/admin" 
              class="px-4 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded-md transition-colors text-sm"
            >
              <i class="fas fa-cog mr-1"></i>管理后台
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg">
        <!-- Tab 切换 -->
        <div class="border-b border-gray-200 dark:border-gray-700">
          <div class="flex">
            <button
              v-for="tab in tabs"
              :key="tab.value"
              :class="[
                'px-6 py-4 text-sm font-medium border-b-2 transition-colors',
                mode === tab.value 
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400' 
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
              ]"
              @click="mode = tab.value"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="p-6">
          <!-- 批量添加 -->
          <BatchAddResource 
            v-if="mode === 'batch'"
            @success="handleSuccess"
            @error="handleError"
            @cancel="handleCancel"
          />

          <!-- 单个添加 -->
          <SingleAddResource 
            v-else-if="mode === 'single'"
            @success="handleSuccess"
            @error="handleError"
            @cancel="handleCancel"
          />

          <!-- API说明 -->
          <ApiDocumentation 
            v-else
            @cancel="handleCancel"
          />
        </div>
      </div>
    </div>

    <!-- 成功/失败提示 -->
    <SuccessToast v-if="showSuccess" :message="successMsg" @close="showSuccess = false" />
    <ErrorToast v-if="showError" :message="errorMsg" @close="showError = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BatchAddResource from '~/components/BatchAddResource.vue'
import SingleAddResource from '~/components/SingleAddResource.vue'
import ApiDocumentation from '~/components/ApiDocumentation.vue'
import SuccessToast from '~/components/SuccessToast.vue'
import ErrorToast from '~/components/ErrorToast.vue'

const router = useRouter()
const showSuccess = ref(false)
const successMsg = ref('')
const showError = ref(false)
const errorMsg = ref('')

const tabs = [
  { label: '批量添加', value: 'batch' },
  { label: '单个添加', value: 'single' },
  { label: 'API说明', value: 'api' },
]
const mode = ref('batch')

// 检查用户权限
onMounted(() => {
  const userStore = useUserStore()
  if (!userStore.isAuthenticated) {
    router.push('/login')
    return
  }
})

// 事件处理
const handleSuccess = (message: string) => {
  successMsg.value = message
  showSuccess.value = true
}

const handleError = (message: string) => {
  errorMsg.value = message
  showError.value = true
}

const handleCancel = () => {
  router.back()
}

// 设置页面标题
useHead({
  title: '添加资源 - 网盘资源数据库'
})
</script>

<style scoped>
/* 自定义样式 */
</style> 