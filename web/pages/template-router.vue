<template>
  <div>
    <TemplatePageWrapper :page-path="route.path" />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from '#app'
import { useTemplateManager } from '~/composables/templateManager'

definePageMeta({
  layout: 'dynamic'
})

const route = useRoute()
const { currentTemplate } = useTemplateManager()

// SEO设置
useSeoMeta({
  title: '页面加载中 - 老九网盘资源数据库',
  description: '正在加载页面内容...',
})

// 监听模板变化并更新页面标题
watch(currentTemplate, () => {
  useHead({
    title: `页面加载中 - ${currentTemplate.value}模板`
  })
}, { immediate: true })
</script>