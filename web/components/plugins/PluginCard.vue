<template>
  <div class="plugin-card p-6 hover:bg-gray-50 transition-colors">
    <div class="flex items-start justify-between">
      <!-- 插件基本信息 -->
      <div class="flex-1">
        <div class="flex items-center space-x-3">
          <h3 class="text-lg font-semibold text-gray-900">{{ plugin.name }}</h3>
          <span class="px-2 py-1 text-xs font-medium rounded-full" :class="statusClass">
            {{ statusText }}
          </span>
          <span class="px-2 py-1 text-xs font-medium bg-gray-100 text-gray-800 rounded-full">
            {{ plugin.category }}
          </span>
        </div>

        <p class="text-gray-600 mt-1">{{ plugin.description }}</p>

        <div class="flex items-center space-x-4 mt-2 text-sm text-gray-500">
          <span>版本 {{ plugin.version }}</span>
          <span>作者: {{ plugin.author }}</span>
          <span>大小: {{ formatFileSize(plugin.file_size) }}</span>
          <span>更新: {{ formatDate(plugin.last_updated) }}</span>
        </div>

        <!-- 执行统计 -->
        <div v-if="plugin.execution_stats" class="mt-3 flex items-center space-x-6 text-sm">
          <div class="flex items-center space-x-1">
            <Icon name="activity" class="w-4 h-4 text-gray-400" />
            <span class="text-gray-600">执行: {{ plugin.execution_stats.total_executions }}</span>
          </div>
          <div class="flex items-center space-x-1">
            <Icon name="check-circle" class="w-4 h-4 text-green-500" />
            <span class="text-green-600">成功率: {{ plugin.execution_stats.success_rate.toFixed(1) }}%</span>
          </div>
          <div class="flex items-center space-x-1">
            <Icon name="clock" class="w-4 h-4 text-gray-400" />
            <span class="text-gray-600">平均耗时: {{ plugin.execution_stats.average_time }}ms</span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex items-center space-x-2 ml-4">
        <!-- 启用/禁用按钮 -->
        <button
          @click="togglePlugin"
          :class="plugin.enabled ? 'bg-red-100 hover:bg-red-200 text-red-700' : 'bg-green-100 hover:bg-green-200 text-green-700'"
          class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors"
        >
          <Icon :name="plugin.enabled ? 'x-circle' : 'check-circle'" class="w-4 h-4" />
        </button>

        <!-- 配置按钮 -->
        <button
          @click="$emit('configure', plugin)"
          class="p-1.5 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors"
          title="配置插件"
        >
          <Icon name="settings" class="w-4 h-4" />
        </button>

        <!-- 日志按钮 -->
        <button
          @click="$emit('viewLogs', plugin)"
          class="p-1.5 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors"
          title="查看日志"
        >
          <Icon name="file-text" class="w-4 h-4" />
        </button>

        <!-- 详情按钮 -->
        <button
          @click="$emit('viewDetails', plugin)"
          class="p-1.5 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors"
          title="查看详情"
        >
          <Icon name="info" class="w-4 h-4" />
        </button>

        <!-- 更多操作菜单 -->
        <div class="relative">
          <button
            @click="showMenu = !showMenu"
            class="p-1.5 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-md transition-colors"
          >
            <Icon name="more-vertical" class="w-4 h-4" />
          </button>

          <!-- 下拉菜单 -->
          <div
            v-if="showMenu"
            v-click-outside="() => showMenu = false"
            class="absolute right-0 mt-1 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200"
          >
            <div class="py-1">
              <button
                @click="exportPlugin"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center space-x-2"
              >
                <Icon name="download" class="w-4 h-4" />
                <span>导出插件</span>
              </button>
              <button
                @click="duplicatePlugin"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center space-x-2"
              >
                <Icon name="copy" class="w-4 h-4" />
                <span>复制插件</span>
              </button>
              <hr class="my-1">
              <button
                @click="uninstallPlugin"
                class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center space-x-2"
              >
                <Icon name="trash" class="w-4 h-4" />
                <span>卸载插件</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 插件标签 -->
    <div v-if="plugin.tags && plugin.tags.length > 0" class="mt-3 flex flex-wrap gap-1">
      <span
        v-for="tag in plugin.tags"
        :key="tag"
        class="px-2 py-1 text-xs bg-blue-50 text-blue-700 rounded-md"
      >
        {{ tag }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useToast } from '~/composables/useToast'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['enable', 'disable', 'configure', 'viewLogs', 'viewDetails'])

const toast = useToast()
const showMenu = ref(false)

// 计算属性
const statusClass = computed(() => {
  if (!props.plugin.enabled) {
    return 'bg-gray-100 text-gray-800'
  }
  if (props.plugin.status === 'error') {
    return 'bg-red-100 text-red-800'
  }
  return 'bg-green-100 text-green-800'
})

const statusText = computed(() => {
  if (!props.plugin.enabled) {
    return '已禁用'
  }
  if (props.plugin.status === 'error') {
    return '错误'
  }
  return '运行中'
})

// 方法
const togglePlugin = () => {
  if (props.plugin.enabled) {
    emit('disable', props.plugin.name)
  } else {
    emit('enable', props.plugin.name)
  }
}

const exportPlugin = () => {
  // 导出插件逻辑
  toast.info('导出功能开发中...')
  showMenu.value = false
}

const duplicatePlugin = () => {
  // 复制插件逻辑
  toast.info('复制功能开发中...')
  showMenu.value = false
}

const uninstallPlugin = () => {
  // 卸载插件逻辑
  toast.info('卸载功能开发中...')
  showMenu.value = false
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// v-click-outside 指令
const vClickOutside = {
  mounted(el, binding) {
    el.__ClickOutsideHandler__ = event => {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value(event)
      }
    }
    document.addEventListener('click', el.__ClickOutsideHandler__)
  },
  unmounted(el) {
    document.removeEventListener('click', el.__ClickOutsideHandler__)
  }
}
</script>

<style scoped>
.plugin-card {
  border-left: 4px solid transparent;
  transition: all 0.2s ease;
}

.plugin-card:hover {
  border-left-color: #3b82f6;
}

.relative {
  position: relative;
}

/* 下拉菜单动画 */
.relative > div {
  transform-origin: top right;
  transition: all 0.1s ease;
}

/* 状态指示器动画 */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.status-indicator {
  animation: pulse 2s infinite;
}
</style>