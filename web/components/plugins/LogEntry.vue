<template>
  <div class="log-entry p-4 hover:bg-gray-50 transition-colors cursor-pointer" @click="$emit('viewDetails', log)">
    <div class="flex items-start space-x-3">
      <!-- 状态图标 -->
      <div class="flex-shrink-0 mt-1">
        <div
          class="w-2 h-2 rounded-full"
          :class="statusClass"
        ></div>
      </div>

      <!-- 主要内容 -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <!-- 钩子名称 -->
            <span class="font-medium text-gray-900">{{ log.hook_name }}</span>

            <!-- 执行状态 -->
            <span
              class="px-2 py-0.5 text-xs font-medium rounded-full"
              :class="log.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
            >
              {{ log.success ? '成功' : '失败' }}
            </span>

            <!-- 执行时间 -->
            <span class="text-sm text-gray-500">
              {{ log.execution_time }}ms
            </span>
          </div>

          <!-- 时间戳 -->
          <span class="text-sm text-gray-500">
            {{ formatTime(log.created_at) }}
          </span>
        </div>

        <!-- 日志消息 -->
        <div v-if="log.message" class="mt-1">
          <p class="text-sm text-gray-700 break-all">{{ log.message }}</p>
        </div>

        <!-- 错误信息 -->
        <div v-if="!log.success && log.error_message" class="mt-1">
          <p class="text-sm text-red-600 break-all">{{ log.error_message }}</p>
        </div>

        <!-- 详细信息（可展开） -->
        <div v-if="showDetails" class="mt-3 space-y-2">
          <div class="bg-gray-50 rounded-lg p-3 text-sm">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <span class="text-gray-600">日志ID:</span>
                <span class="ml-2 font-mono text-gray-900">{{ log.id }}</span>
              </div>
              <div>
                <span class="text-gray-600">插件名称:</span>
                <span class="ml-2 text-gray-900">{{ log.plugin_name }}</span>
              </div>
              <div>
                <span class="text-gray-600">执行时间:</span>
                <span class="ml-2 text-gray-900">{{ log.execution_time }}ms</span>
              </div>
              <div>
                <span class="text-gray-600">执行状态:</span>
                <span class="ml-2" :class="log.success ? 'text-green-600' : 'text-red-600'">
                  {{ log.success ? '成功' : '失败' }}
                </span>
              </div>
            </div>

            <!-- 日志消息 -->
            <div v-if="log.message" class="mt-3">
              <span class="text-gray-600">日志消息:</span>
              <pre class="mt-1 p-2 bg-blue-50 border border-blue-200 rounded text-blue-800 text-xs overflow-x-auto">{{ log.message }}</pre>
            </div>

            <!-- 错误信息 -->
            <div v-if="log.error_message" class="mt-3">
              <span class="text-gray-600">错误信息:</span>
              <pre class="mt-1 p-2 bg-red-50 border border-red-200 rounded text-red-800 text-xs overflow-x-auto">{{ log.error_message }}</pre>
            </div>

            <div class="mt-3">
              <span class="text-gray-600">创建时间:</span>
              <span class="ml-2 text-gray-900">{{ formatFullTime(log.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex items-center space-x-1">
        <button
          @click.stop="toggleDetails"
          class="p-1 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
          :title="showDetails ? '收起详情' : '展开详情'"
        >
          <Icon :name="showDetails ? 'chevron-up' : 'chevron-down'" class="w-4 h-4" />
        </button>
        <button
          @click.stop="copyLog"
          class="p-1 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
          title="复制日志"
        >
          <Icon name="copy" class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useToast } from '~/composables/useToast'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  log: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['viewDetails'])

const toast = useToast()
const showDetails = ref(false)

// 计算属性
const statusClass = computed(() => {
  if (!props.log.success) {
    return 'bg-red-500'
  }
  if (props.log.execution_time > 1000) {
    return 'bg-yellow-500'
  }
  return 'bg-green-500'
})

// 方法
const toggleDetails = () => {
  showDetails.value = !showDetails.value
}

const copyLog = async () => {
  try {
    const logText = `插件: ${props.log.plugin_name}
钩子: ${props.log.hook_name}
状态: ${props.log.success ? '成功' : '失败'}
执行时间: ${props.log.execution_time}ms
创建时间: ${formatFullTime(props.log.created_at)}
${props.log.message ? `日志消息: ${props.log.message}` : ''}
${props.log.error_message ? `错误信息: ${props.log.error_message}` : ''}`

    await navigator.clipboard.writeText(logText)
    toast.success('日志已复制到剪贴板')
  } catch (error) {
    toast.error('复制失败: ' + error.message)
  }
}

const formatTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date

  // 如果是今天，只显示时间
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  // 如果是昨天，显示"昨天"
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return '昨天 ' + date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  // 其他情况显示完整日期
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatFullTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
</script>

<style scoped>
.log-entry {
  border-left: 3px solid transparent;
  transition: all 0.2s ease;
}

.log-entry:hover {
  border-left-color: #3b82f6;
  background-color: #f9fafb;
}

/* 状态指示器动画 */
.w-2.h-2.rounded-full {
  transition: background-color 0.3s ease;
}

/* 代码块样式 */
pre {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 点击效果 */
.log-entry:hover {
  transform: translateX(2px);
}
</style>