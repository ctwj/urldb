<template>
  <!-- AI聊天模态框 -->
  <n-modal v-model:show="visible" :mask-closable="false" preset="card" style="width: 900px; max-width: 95vw;">
    <template #header>
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
          <i class="fas fa-robot text-white text-lg"></i>
        </div>
        <div>
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">AI 智能对话</h2>
          <p class="text-sm text-gray-500 dark:text-gray-400">与AI助手进行对话，测试AI聊天和MCP工具调用</p>
        </div>
      </div>
    </template>

    <div class="chat-container" style="height: 600px;">
      <!-- 聊天消息区域 -->
      <div class="flex flex-col h-full">
        <!-- 消息列表 -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50 dark:bg-gray-900 rounded-lg" ref="chatContainer">
          <!-- 欢迎消息 -->
          <div v-if="messages.length === 0" class="text-center py-12">
            <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full mb-4">
              <i class="fas fa-comments text-white text-2xl"></i>
            </div>
            <h3 class="text-lg font-semibold text-gray-700 dark:text-gray-300 mb-2">开始与AI对话</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">输入您的问题，AI助手将为您提供帮助</p>
            <div class="mt-4 flex flex-wrap justify-center gap-2">
              <n-button size="small" @click="sendQuickMessage('你好，请介绍一下你自己')">
                <template #icon>
                  <i class="fas fa-hand-point-right"></i>
                </template>
                自我介绍
              </n-button>
              <n-button size="small" @click="sendQuickMessage('现在几点了？')">
                <template #icon>
                  <i class="fas fa-clock"></i>
                </template>
                查询时间
              </n-button>
              <n-button size="small" @click="sendQuickMessage('搜索最新AI技术资讯')">
                <template #icon>
                  <i class="fas fa-search"></i>
                </template>
                搜索资讯
              </n-button>
            </div>
          </div>

          <!-- 消息列表 -->
          <div v-for="(message, index) in messages" :key="index" class="flex" :class="message.role === 'user' ? 'justify-end' : 'justify-start'">
            <div class="flex max-w-[80%]">
              <!-- AI消息 -->
              <div v-if="message.role === 'assistant'" class="flex items-start space-x-3">
                <div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center flex-shrink-0">
                  <i class="fas fa-robot text-white text-sm"></i>
                </div>
                <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm border border-gray-200 dark:border-gray-700">
                  <div class="text-sm text-gray-800 dark:text-gray-200 leading-relaxed whitespace-pre-wrap">{{ message.content }}</div>
                  <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">
                    {{ formatTime(message.timestamp) }}
                  </div>
                </div>
              </div>

              <!-- 用户消息 -->
              <div v-else class="flex items-start space-x-3 flex-row-reverse">
                <div class="w-8 h-8 bg-gradient-to-br from-green-500 to-teal-600 rounded-full flex items-center justify-center flex-shrink-0">
                  <i class="fas fa-user text-white text-sm"></i>
                </div>
                <div class="bg-gradient-to-br from-blue-500 to-purple-600 text-white rounded-lg p-4 shadow-sm">
                  <div class="text-sm leading-relaxed whitespace-pre-wrap">{{ message.content }}</div>
                  <div class="mt-2 text-xs opacity-80">
                    {{ formatTime(message.timestamp) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 加载中状态 -->
          <div v-if="loading" class="flex justify-start">
            <div class="flex items-start space-x-3">
              <div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center flex-shrink-0">
                <i class="fas fa-robot text-white text-sm"></i>
              </div>
              <div class="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-sm border border-gray-200 dark:border-gray-700">
                <div class="flex items-center space-x-2">
                  <div class="flex space-x-1">
                    <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                    <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                    <div class="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
                  </div>
                  <span class="text-sm text-gray-500">AI正在思考...</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4 rounded-b-lg">
          <div class="flex space-x-2">
            <n-input
              v-model:value="input"
              placeholder="输入您的问题..."
              @keydown.enter="sendMessage"
              :disabled="loading"
              type="textarea"
              :autosize="{ minRows: 2, maxRows: 4 }"
              class="flex-1"
            />
            <n-button
              @click="sendMessage"
              :loading="loading"
              :disabled="!input.trim()"
              type="primary"
              size="medium"
            >
              <template #icon>
                <i class="fas fa-paper-plane"></i>
              </template>
              发送
            </n-button>
          </div>
          <div class="mt-2 flex items-center justify-between">
            <div class="text-xs text-gray-500 dark:text-gray-400">
              <i class="fas fa-info-circle mr-1"></i>
              AI将使用当前配置进行回复，支持MCP工具调用
            </div>
            <div class="flex space-x-2">
              <n-button size="small" @click="clearChat" :disabled="messages.length === 0">
                <template #icon>
                  <i class="fas fa-trash"></i>
                </template>
                清空对话
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </div>

      </n-modal>
</template>

<script setup lang="ts">
import { ref, nextTick, computed } from 'vue'
import { useNotification } from 'naive-ui'
import { useAIApi } from '~/composables/useApi'

// Props
interface Props {
  modelValue: boolean
  config?: {
    api_url?: string
    model?: string
    max_tokens?: number
    temperature?: number
    timeout?: number
    retry_count?: number
  }
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  config: () => ({})
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

// 响应式数据
const notification = useNotification()
const { generateTextWithTools } = useAIApi()
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const input = ref('')
const loading = ref(false)
const messages = ref<Array<{
  role: 'user' | 'assistant'
  content: string
  timestamp: number
}>>([])
const chatContainer = ref<HTMLElement | null>(null)

// 方法
const formatTime = (timestamp: number) => {
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const sendMessage = async () => {
  if (!input.value.trim()) return

  const userMessage = input.value.trim()

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: userMessage,
    timestamp: Date.now()
  })

  input.value = ''
  loading.value = true

  // 滚动到底部
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })

  try {
    // 使用AI的文本生成接口（支持工具调用）
    const response = await generateTextWithTools({
      prompt: userMessage
    })

    // 解析响应数据，获取result字段
    let aiResponse = '抱歉，我暂时无法回复您的问题。'
    if (response && response.result) {
      aiResponse = response.result
    } else if (typeof response === 'string') {
      aiResponse = response
    }

    // 添加AI回复
    messages.value.push({
      role: 'assistant',
      content: aiResponse,
      timestamp: Date.now()
    })
  } catch (error: any) {
    console.error('AI聊天失败:', error)

    // 添加错误消息
    messages.value.push({
      role: 'assistant',
      content: '抱歉，AI服务暂时不可用。请检查AI配置后重试。',
      timestamp: Date.now()
    })

    notification.error({
      content: 'AI聊天失败，请检查配置',
      duration: 3000
    })
  } finally {
    loading.value = false

    // 滚动到底部
    nextTick(() => {
      if (chatContainer.value) {
        chatContainer.value.scrollTop = chatContainer.value.scrollHeight
      }
    })
  }
}

const sendQuickMessage = (message: string) => {
  input.value = message
  sendMessage()
}

const clearChat = () => {
  messages.value = []
  notification.info({
    content: '对话已清空',
    duration: 2000
  })
}

const close = () => {
  visible.value = false
}
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 自定义滚动条样式 */
.chat-container .overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.chat-container .overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.chat-container .overflow-y-auto::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.3);
  border-radius: 3px;
}

.chat-container .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.5);
}

/* 暗色模式滚动条 */
.dark .chat-container .overflow-y-auto::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.3);
}

.dark .chat-container .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background-color: rgba(75, 85, 99, 0.5);
}
</style>