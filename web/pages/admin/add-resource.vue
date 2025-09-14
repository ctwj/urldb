<template>
  <AdminPageLayout :is-sub-page="true">
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">添加资源</h1>
        <p class="text-gray-600 dark:text-gray-400">添加新的资源到系统</p>
      </div>
      <div></div>
    </template>

    <!-- 内容区 -->
    <template #content>
      <div class="resource-content h-full">
        <n-tabs
          v-model:value="activeTab"
          type="line"
          animated
          class="mb-6"
        >
          <n-tab-pane name="batch" tab="批量添加">
            <div class="tab-content-container">
              <div class="space-y-8">
                <AdminBatchAddResource
                  @success="handleSuccess"
                  @error="handleError"
                  @cancel="handleCancel"
                />
              </div>
            </div>
          </n-tab-pane>

        <n-tab-pane name="singal" tab="单个添加">
          <div class="tab-content-container">
            <AdminSingleAddResource
              @success="handleSuccess"
              @error="handleError"
              @cancel="handleCancel"
            />
          </div>
        </n-tab-pane>

        </n-tabs>
      </div>
    </template>
  </AdminPageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AdminPageLayout from '~/components/AdminPageLayout.vue'

// 设置页面布局
definePageMeta({
  layout: 'admin'
})

// 根据 Nuxt 3 组件规则，位于 components/Admin/ 的组件会自动以 Admin 前缀导入
const activeTab = ref('batch')
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


</script>

<style scoped>
/* 添加资源页面样式 */

.resource-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .resource-content {
  background-color: var(--color-dark-bg, #1f2937);
}
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>