<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- 顶部导航栏 -->
    <UserHeader />
    
    <!-- 侧边栏和主内容区域 -->
    <div class="flex">
      <!-- 侧边栏 -->
      <UserSidebar />
      
      <!-- 主内容区域 -->
      <main class="flex-1 p-8">
        <ClientOnly>
          <n-notification-provider>
            <n-dialog-provider>
              <!-- 页面内容插槽 -->
              <slot />
            </n-dialog-provider>
          </n-notification-provider>
        </ClientOnly>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserLayout } from '~/composables/useUserLayout'

// 使用用户布局组合式函数
const { checkAuth } = useUserLayout()

// 页面加载时检查认证状态
onMounted(() => {
  checkAuth()
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 