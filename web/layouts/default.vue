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

// 使用 useCookie 来确保服务端和客户端状态一致
const themeCookie = useCookie('theme', { default: () => 'light' })

// 初始化主题状态
isDark.value = themeCookie.value === 'dark'

const toggleDarkMode = () => {
  isDark.value = !isDark.value
  const newTheme = isDark.value ? 'dark' : 'light'
  
  // 更新 cookie
  themeCookie.value = newTheme
  
  // 更新 DOM 类
  if (process.client) {
    if (isDark.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }
}



const injectRawScript = (rawScriptString: string) => {
  if (process.client) {
    // 创建一个临时容器来解析原始字符串
    const container = document.createElement('div');
    container.innerHTML = rawScriptString.trim();

    // 获取解析后的所有 script 元素
    const scripts = container.querySelectorAll('script');

    // 遍历并注入所有脚本
    scripts.forEach((script) => {
      // 创建新的 script 元素
      const newScript = document.createElement('script');

      // 复制所有属性（包括 data-*、async、defer 等）
      Array.from(script.attributes).forEach((attr) => {
        newScript.setAttribute(attr.name, attr.value);
      });

      // 复制内容（如果是内联脚本）
      if (script.innerHTML) {
        newScript.innerHTML = script.innerHTML;
      }

      // 插入到 DOM
      document.head.appendChild(newScript);
    });
  }
};

// 获取三方统计代码并直接加载
const fetchStatsCode = async () => {
  try {
    const { usePublicSystemConfigApi } = await import('~/composables/useApi')
    const publicSystemConfigApi = usePublicSystemConfigApi()
    const response = await publicSystemConfigApi.getPublicSystemConfig()
    
    if (response?.data && response.data.third_party_stats_code) {
      injectRawScript(response.data.third_party_stats_code);
      console.log('三方统计代码已加载')
    }
  } catch (error) {
    console.error('获取三方统计代码失败:', error)
  }
}



onMounted(async () => {
  // 初始化主题 - 使用 cookie 而不是 localStorage
  if (themeCookie.value === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
  
  // 获取三方统计代码并直接加载
  await fetchStatsCode()
})
</script> 