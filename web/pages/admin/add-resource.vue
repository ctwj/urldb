<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">添加资源</h1>
        <p class="text-gray-600 dark:text-gray-400">添加新的资源到系统</p>
      </div>
      <n-button @click="navigateTo('/admin/resources')">
        <template #icon>
          <i class="fas fa-arrow-left"></i>
        </template>
        返回资源管理
      </n-button>
    </div>

    <!-- 主要内容 -->
    <n-card>
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
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

import { ref } from 'vue'
import BatchAddResource from '~/components/BatchAddResource.vue'
import SingleAddResource from '~/components/SingleAddResource.vue'

const tabs = [
  { label: '批量添加', value: 'batch' },
  { label: '单个添加', value: 'single' },
]
const mode = ref('batch')
const notification = useNotification()

// 事件处理
const handleSuccess = (message: string) => {
  notification.success({
    content: message,
    duration: 3000
  })
}

const handleError = (message: string) => {
  notification.error({
    content: message,
    duration: 3000
  })
}

const handleCancel = () => {
  navigateTo('/admin/resources')
}

// 设置页面标题
useHead({
  title: '添加资源 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 自定义样式 */
</style> 