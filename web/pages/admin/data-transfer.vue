<template>
  <div class="max-w-7xl mx-auto space-y-6">
    <!-- 页面标题 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">数据转存管理</h1>
      <p class="text-gray-600 dark:text-gray-400">管理资源转存任务和状态</p>
    </div>

    <!-- 主要内容 -->
    <n-card>
      <n-tabs v-model:value="activeTab" type="line" animated>
        <!-- 手动批量转存 -->
        <n-tab-pane name="manual" tab="手动批量转存">
          <AdminManualBatchTransfer />
        </n-tab-pane>

        <!-- 已转存列表 -->
        <n-tab-pane name="transferred" tab="已转存列表">
          <AdminTransferredList ref="transferredListRef" />
        </n-tab-pane>

        <!-- 未转存列表 -->
        <n-tab-pane name="untransferred" tab="未转存列表">
          <AdminUntransferredList ref="untransferredListRef" />
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

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