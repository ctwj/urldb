<script setup lang="ts">
// 用户中心首页已移除：进入 /user 直接跳转「我的资源」
definePageMeta({
  layout: 'user'
})

onMounted(() => {
  navigateTo('/user/resources', { replace: true })
})
</script>

<template>
  <div class="flex items-center justify-center py-24 text-gray-400">
    <n-spin size="small" />
  </div>
</template>
