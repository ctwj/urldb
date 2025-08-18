<template>
  <footer class="footer-container mt-auto py-6 border-t border-gray-700 bg-white dark:bg-gray-800">
    <div class="max-w-7xl mx-auto text-center text-gray-400  text-sm px-3 sm:px-5">
      <p class="mb-2">本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小内删除!!!</p>
      <p class="flex items-center justify-center gap-2">
        <span>{{ systemConfig?.copyright || '© 2025 老九网盘资源数据库 By 老九' }}</span>
        <span v-if="versionInfo && versionInfo.version" class="text-gray-400 dark:text-gray-500">| v  <n-a 
            href="https://github.com/ctwj/urldb"
            target="_blank"
            rel="noopener noreferrer"
            referrerpolicy="no-referrer"
            aria-label="在 GitHub 上查看版本信息"
            class="github-link"
          ><span>{{ versionInfo.version }}</span></n-a>
        </span>
      </p>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'
// 使用版本信息组合式函数
const { versionInfo, fetchVersionInfo } = useVersion()

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('footerSystemConfig', 
  () => useApiFetch('/system/config').then(parseApiResponse)
)

const systemConfig = computed(() => (systemConfigData.value as any) || { copyright: '© 2025 老九网盘资源数据库 By 老九' })

// 组件挂载时获取版本信息
onMounted(() => {
  fetchVersionInfo()
})
</script> 

<style scoped>
.footer-container{
  background: url(/assets/images/footer-banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}
</style>