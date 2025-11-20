<template>
  <ClientOnly>
    <!-- 自定义背景遮罩 -->
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-start justify-center pt-[20vh]"
      @click="handleBackdropClick"
    >
      <!-- 背景模糊遮罩 -->
      <div class="absolute inset-0 bg-black/20 backdrop-blur-sm"></div>

      <!-- 搜索弹窗 -->
      <div
        class="relative w-full max-w-2xl mx-4 transform transition-all duration-200 ease-out"
        :class="visible ? 'scale-100 opacity-100' : 'scale-95 opacity-0'"
        @click.stop
      >
        <!-- 搜索输入区域 -->
        <div class="relative bg-white dark:bg-gray-900 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          <!-- 顶部装饰条 -->
          <div class="h-1 bg-gradient-to-r from-green-500 via-emerald-500 to-teal-500"></div>

          <!-- 搜索输入框 -->
          <div class="relative px-6 py-5">
            <div class="relative flex items-center">
              <!-- 搜索图标 -->
              <div class="absolute left-4 flex items-center pointer-events-none">
                <div class="w-5 h-5 rounded-full bg-gradient-to-r from-green-500 to-emerald-500 flex items-center justify-center">
                  <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                  </svg>
                </div>
              </div>

              <!-- 输入框 -->
              <input
                ref="searchInput"
                v-model="searchQuery"
                type="text"
                placeholder="搜索资源..."
                class="w-full pl-12 pr-32 py-4 bg-transparent border-0 text-lg text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-0"
                @keyup.enter="handleSearch"
                @input="handleInputChange"
                @keydown.escape="handleClose"
              >

              <!-- 搜索按钮 -->
              <div class="absolute right-2 flex items-center gap-2">
                <button
                  v-if="searchQuery.trim()"
                  type="button"
                  @click="clearSearch"
                  class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
                <button
                  type="button"
                  @click="handleSearch"
                  :disabled="!searchQuery.trim()"
                  :loading="searching"
                  class="px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-500 text-white text-sm font-medium rounded-lg hover:from-green-600 hover:to-emerald-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
                >
                  <span v-if="!searching" class="flex items-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                    </svg>
                    搜索
                  </span>
                  <span v-else class="flex items-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    搜索中
                  </span>
                </button>
              </div>
            </div>
          </div>

          <!-- 搜索建议下拉 -->
          <div v-if="showSuggestions && suggestions.length > 0" class="border-t border-gray-200 dark:border-gray-700">
            <div class="max-h-60 overflow-y-auto">
              <div class="px-6 py-3">
                <div class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">搜索建议</div>
                <div class="space-y-1">
                  <button
                    v-for="(suggestion, index) in suggestions"
                    :key="index"
                    @click="selectSuggestion(suggestion)"
                    class="w-full flex items-center gap-3 px-3 py-2.5 text-left rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors group"
                  >
                    <div class="w-8 h-8 rounded-lg bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20 flex items-center justify-center group-hover:scale-110 transition-transform">
                      <svg class="w-4 h-4 text-green-500 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                      </svg>
                    </div>
                    <div class="flex-1">
                      <div class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ suggestion }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">点击搜索 "{{ suggestion }}"</div>
                    </div>
                    <svg class="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- 搜索历史 -->
          <div v-if="searchHistory.length > 0" class="border-t border-gray-200 dark:border-gray-700">
            <div class="px-6 py-4">
              <div class="flex items-center justify-between mb-3">
                <div class="flex items-center gap-2">
                  <div class="w-6 h-6 rounded-lg bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                    <svg class="w-3 h-3 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                  <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">最近搜索</span>
                </div>
                <button
                  @click="clearHistory"
                  class="text-xs text-gray-500 hover:text-red-500 dark:text-gray-400 dark:hover:text-red-400 transition-colors"
                >
                  清空
                </button>
              </div>
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="(item, index) in searchHistory.slice(0, 8)"
                  :key="index"
                  @click="selectHistory(item)"
                  class="inline-flex items-center gap-2 px-3 py-1.5 text-sm bg-gray-50 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-all duration-200 group"
                >
                  <svg class="w-3 h-3 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  {{ item }}
                </button>
              </div>
            </div>
          </div>

          <!-- 搜索提示 -->
          <div class="border-t border-gray-200 dark:border-gray-700 bg-gradient-to-r from-green-50 via-emerald-50 to-teal-50 dark:from-green-900/20 dark:via-emerald-900/20 dark:to-teal-900/20">
            <div class="px-6 py-4">
              <div class="flex items-center gap-3 mb-2">
                <div class="w-8 h-8 rounded-full bg-white dark:bg-gray-800 shadow-sm flex items-center justify-center">
                  <svg class="w-4 h-4 text-green-500 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">搜索技巧</span>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs text-gray-600 dark:text-gray-400">
                <div class="flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-green-400"></span>
                  <span>支持多关键词搜索，用空格分隔</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
                  <span>按 <kbd class="px-1.5 py-0.5 text-xs bg-white dark:bg-gray-700 rounded border border-gray-300 dark:border-gray-600 font-mono">Ctrl+K</kbd> 快速打开</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-teal-400"></span>
                  <span>按 <kbd class="px-1.5 py-0.5 text-xs bg-white dark:bg-gray-700 rounded border border-gray-300 dark:border-gray-600 font-mono">Esc</kbd> 关闭弹窗</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full bg-green-400"></span>
                  <span>搜索历史自动保存，方便下次使用</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </ClientOnly>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

// 组件状态 - 完全内部管理
const visible = ref(false)
const searchInput = ref<any>(null)
const searchQuery = ref('')
const searching = ref(false)
const showSuggestions = ref(false)
const searchHistory = ref<string[]>([])

// 路由器
const router = useRouter()

// 计算属性
const suggestions = computed(() => {
  if (!searchQuery.value.trim()) return []

  const query = searchQuery.value.toLowerCase().trim()

  return searchHistory.value
    .filter(item => item.toLowerCase().includes(query))
    .filter(item => item.toLowerCase() !== query)
    .slice(0, 5)
})

// 初始化搜索历史
const initSearchHistory = () => {
  if (process.client && typeof localStorage !== 'undefined') {
    const history = localStorage.getItem('searchHistory')
    if (history) {
      try {
        searchHistory.value = JSON.parse(history)
      } catch (e) {
        searchHistory.value = []
      }
    }
  }
}

// 保存搜索历史
const saveSearchHistory = () => {
  if (process.client && typeof localStorage !== 'undefined') {
    localStorage.setItem('searchHistory', JSON.stringify(searchHistory.value))
  }
}

// 处理输入变化
const handleInputChange = () => {
  showSuggestions.value = searchQuery.value.trim().length > 0
}

// 处理搜索
const handleSearch = () => {
  const query = searchQuery.value.trim()
  if (!query) return

  searching.value = true

  // 添加到搜索历史
  if (!searchHistory.value.includes(query)) {
    searchHistory.value.unshift(query)
    if (searchHistory.value.length > 10) {
      searchHistory.value = searchHistory.value.slice(0, 10)
    }
    saveSearchHistory()
  }

  // 关闭弹窗
  visible.value = false

  // 跳转到搜索页面
  nextTick(() => {
    router.push(`/?search=${encodeURIComponent(query)}`)
  })

  setTimeout(() => {
    searching.value = false
  }, 500)
}

// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  showSuggestions.value = false
  nextTick(() => {
    searchInput.value?.focus()
  })
}

// 选择搜索建议
const selectSuggestion = (suggestion: string) => {
  searchQuery.value = suggestion
  showSuggestions.value = false
  nextTick(() => {
    searchInput.value?.focus()
  })
}

// 选择历史记录
const selectHistory = (item: string) => {
  searchQuery.value = item
  handleSearch()
}

// 清空历史
const clearHistory = () => {
  searchHistory.value = []
  saveSearchHistory()
}

// 处理背景点击
const handleBackdropClick = () => {
  handleClose()
}

// 处理关闭
const handleClose = () => {
  visible.value = false
  searchQuery.value = ''
  showSuggestions.value = false
}

// 监听弹窗显示状态
watch(visible, (newValue) => {
  if (newValue && process.client) {
    nextTick(() => {
      searchInput.value?.focus()
      initSearchHistory()
    })
  } else {
    setTimeout(() => {
      searchQuery.value = ''
      showSuggestions.value = false
    }, 300)
  }
})

// 键盘事件监听
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    if (!visible.value) {
      visible.value = true
    }
  }

  if (e.key === 'Escape' && visible.value) {
    handleClose()
  }
}

// 组件挂载时添加键盘事件监听器
onMounted(() => {
  if (process.client && typeof document !== 'undefined') {
    document.addEventListener('keydown', handleKeydown)
  }
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  if (process.client && typeof document !== 'undefined') {
    document.removeEventListener('keydown', handleKeydown)
  }
})

// 暴露给父组件的方法
defineExpose({
  show: () => { visible.value = true },
  hide: () => { handleClose() },
  toggle: () => { visible.value = !visible.value }
})
</script>

<style scoped>
/* 自定义动画 */
.transform {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 滚动条样式 */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.3);
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.5);
}

/* 深色模式滚动条 */
.dark .overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(75, 85, 99, 0.3);
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(75, 85, 99, 0.5);
}

/* 键盘快捷键样式 */
kbd {
  box-shadow: 0 1px 0 1px rgba(0, 0, 0, 0.1), 0 1px 0 rgba(0, 0, 0, 0.1);
}

/* 按钮悬停效果 */
button {
  transition: all 0.15s ease-in-out;
}

/* 输入框聚焦效果 */
input:focus {
  box-shadow: none;
}

/* 渐变动画 */
@keyframes gradient {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.bg-gradient-to-r {
  background-size: 200% 200%;
  animation: gradient 3s ease infinite;
}
</style>