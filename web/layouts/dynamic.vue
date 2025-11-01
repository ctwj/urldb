<template>
  <div>
    <component :is="layoutComponent" />
  </div>
</template>

<script setup lang="ts">
import { useTemplateManager } from '~/composables/templateManager'
import { useRoute } from '#app'

// 初始化模板管理器
const { currentTemplate, init } = useTemplateManager()
init()

// 获取当前路由
const route = useRoute()

// 根据当前模板动态导入布局组件
const layoutComponent = computed(() => {
  switch(currentTemplate.value) {
    case 'blog':
      return defineAsyncComponent(() => import('~/layouts/blog/default.vue'))
    case 'default':
    default:
      return defineAsyncComponent(() => import('~/layouts/default/default.vue'))
  }
})
</script>