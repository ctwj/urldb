<template>
  <div>
    <!-- 暗色模式切换按钮 -->
    <button
      class="fixed top-4 right-4 z-50 w-8 h-8 flex items-center justify-center rounded-full shadow-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 transition-all duration-200 hover:bg-blue-100 dark:hover:bg-blue-900 hover:scale-110 focus:outline-none"
      @click="toggleDarkMode"
      aria-label="切换明暗模式"
    >
      <span class="text-2xl transition-transform duration-300" :class="isDark ? 'rotate-0' : 'rotate-180'">
        <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="w-5 h-5">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12.79A9 9 0 1111.21 3a7 7 0 109.79 9.79z" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="w-5 h-5">
          <circle cx="12" cy="12" r="5" stroke="currentColor" stroke-width="2" />
          <path stroke-linecap="round" stroke-width="2" d="M12 1v2m0 18v2m11-11h-2M3 12H1m16.95 6.95l-1.41-1.41M6.46 6.46L5.05 5.05m12.02 0l-1.41 1.41M6.46 17.54l-1.41 1.41" />
        </svg>
      </span>
    </button>
 
    <n-notification-provider>
      <n-dialog-provider>
        <NuxtPage />
      </n-dialog-provider>
    </n-notification-provider>
  </div>
</template>

<script setup lang="ts">
import { lightTheme } from 'naive-ui'
import { ref, onMounted } from 'vue'

const theme = lightTheme
const isDark = ref(false)

const toggleDarkMode = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

onMounted(() => {
  // 初始化主题
  if (localStorage.getItem('theme') === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
})
</script> 