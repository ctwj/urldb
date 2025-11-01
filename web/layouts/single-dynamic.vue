<template>
  <div class="single-layout">
    <n-notification-provider>
      <component :is="layoutComponent" />
    </n-notification-provider>
  </div>
</template>

<script setup>
import { useTemplateManager } from '~/composables/templateManager'

// 初始化模板管理器
const { currentTemplate, init } = useTemplateManager()
init()

// 对于single布局，我们仍然使用single布局，但可以根据模板做一些调整
const layoutComponent = computed(() => {
  // single布局保持简单，不根据模板变化
  return defineAsyncComponent(() => import('~/layouts/single.vue'))
})
</script>

<style>
body, html {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  background-color: #fff;
}

.single-layout {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>