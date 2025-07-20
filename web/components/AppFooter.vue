<template>
  <footer class="mt-auto py-6 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
    <div class="max-w-7xl mx-auto text-center text-gray-600 dark:text-gray-400 text-sm px-3 sm:px-5">
      <p class="mb-2">本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小内删除!!!</p>
      <p class="flex items-center justify-center gap-2">
        <span>{{ systemConfig?.copyright || '© 2025 老九网盘资源数据库 By 老九' }}</span>
        <span v-if="versionInfo.version" class="text-gray-400 dark:text-gray-500">| v{{ versionInfo.version }}</span>
      </p>
    </div>
  </footer>
</template>

<script setup lang="ts">
// 使用版本信息组合式函数
const { versionInfo, fetchVersionInfo } = useVersion()

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('systemConfig', 
  () => $fetch('/api/system-config')
)

const systemConfig = computed(() => (systemConfigData.value as any)?.data || { copyright: '© 2025 老九网盘资源数据库 By 老九' })

// 组件挂载时获取版本信息
onMounted(() => {
  fetchVersionInfo()
})
</script> 