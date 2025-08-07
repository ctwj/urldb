<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-bold text-gray-900 dark:text-white">布局测试页面</h1>
    <p class="text-gray-600 dark:text-gray-400">如果您能看到这个页面，说明布局已经生效了。</p>
    
    <n-card>
      <template #header>
        <span class="text-lg font-semibold">测试信息</span>
      </template>
      
      <div class="space-y-4">
        <div class="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
          <h3 class="text-lg font-medium text-green-800 dark:text-green-200 mb-2">✅ 布局生效</h3>
          <p class="text-sm text-green-700 dark:text-green-300">
            如果您能看到这个绿色卡片，说明新管理后台的布局已经正确加载。
          </p>
        </div>
        
        <div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <h3 class="text-lg font-medium text-blue-800 dark:text-blue-200 mb-2">🔧 组件状态</h3>
          <p class="text-sm text-blue-700 dark:text-blue-300">
            顶部导航栏和侧边栏应该已经显示。如果没有显示，请检查控制台错误。
          </p>
        </div>
        
        <div class="p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
          <h3 class="text-lg font-medium text-yellow-800 dark:text-yellow-200 mb-2">⚠️ 注意事项</h3>
          <p class="text-sm text-yellow-700 dark:text-yellow-300">
            确保您已经登录并且具有管理员权限。
          </p>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin' as any
})
</script> 