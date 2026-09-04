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
          <n-config-provider :theme="naiveTheme" :theme-overrides="naiveOverrides">
            <n-notification-provider>
              <n-dialog-provider>
                <!-- 页面内容插槽 -->
                <slot />
              </n-dialog-provider>
            </n-notification-provider>
          </n-config-provider>
        </ClientOnly>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { darkTheme } from 'naive-ui'
import { useUserLayout } from '~/composables/useUserLayout'
import { useTheme } from '~/composables/useTheme'

// 使用用户布局组合式函数
const { checkAuth } = useUserLayout()

// Naive UI 明暗主题（与 admin 布局一致，dark class 切换时组件库同步）
const { mode, naiveOverrides, setMode } = useTheme()
const naiveTheme = computed(() => (mode.value === 'dark' ? darkTheme : null))

// 页面加载时检查认证状态
onMounted(() => {
  // 同步系统明暗偏好，保证 Tailwind dark class 与 Naive UI 主题一致
  if (window.matchMedia?.('(prefers-color-scheme: dark)').matches) {
    setMode('dark')
  }
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
