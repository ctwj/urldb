<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">数据推送</h1>
        <p class="text-gray-600 dark:text-gray-400">数据推送管理</p>
      </div>
    </template>

    <!-- 内容区 -->
    <template #content>
      <div class="text-center py-12">
        <div class="text-gray-400 dark:text-gray-500 mb-4">
          <i class="fas fa-upload text-4xl"></i>
        </div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能开发中</h3>
        <p class="text-gray-500 dark:text-gray-400">数据推送功能正在开发中，敬请期待...</p>
      </div>
    </template>
  </AdminPageLayout>
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'

// 数据推送管理页面
definePageMeta({
  layout: 'admin',
  middleware: ['auth']
})
</script> 