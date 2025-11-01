<template>
  <div>
    <component :is="adminLayoutComponent" />
  </div>
</template>

<script setup lang="ts">
import { useTemplateManager } from '~/composables/templateManager'

// 初始化模板管理器
const { currentTemplate, init } = useTemplateManager()
init()

// 根据当前模板动态导入admin布局组件
const adminLayoutComponent = computed(() => {
  switch(currentTemplate.value) {
    case 'blog':
      // 如果博客模板有特殊的admin布局，可以在这里定义
      return defineAsyncComponent(() => import('~/layouts/admin.vue'))
    case 'default':
    default:
      return defineAsyncComponent(() => import('~/layouts/admin.vue'))
  }
})
</script>