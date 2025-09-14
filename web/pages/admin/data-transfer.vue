<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">数据转存管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理资源转存任务和状态</p>
      </div>
    </template>

    <!-- 内容区 - 配置表单 -->
    <template #content>
      <div class="config-content h-full">
        <!-- 顶部Tabs -->
        <n-tabs
          v-model:value="activeTab"
          type="line"
          animated
          class="mb-6"
        >
          <!-- 手动批量转存 -->
          <n-tab-pane name="manual" tab="手动批量转存">
            <div class="tab-content-container">
              <AdminManualBatchTransfer />
            </div>
          </n-tab-pane>

          <!-- 已转存列表 -->
          <n-tab-pane name="transferred" tab="已转存列表">
            <div class="tab-content-container">
              <AdminTransferredList ref="transferredListRef" />
            </div>
          </n-tab-pane>

          <!-- 未转存列表 -->
          <n-tab-pane name="untransferred" tab="未转存列表">
            <div class="tab-content-container">
              <AdminUntransferredList ref="untransferredListRef" />
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

// 页面配置
definePageMeta({
  layout: 'admin',
  middleware: ['auth']
})

// 活动标签
const activeTab = ref('manual')

// 组件引用
const transferredListRef = ref(null)
const untransferredListRef = ref(null)
</script>

<style scoped>
/* 数据转存管理页面样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}

/* tab内容容器 - 个别内容滚动 */
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>