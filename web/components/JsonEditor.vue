<template>
  <div class="json-editor-container">
    <!-- 工具栏 -->
    <div class="editor-toolbar flex items-center justify-between mb-2 p-2 bg-gray-100 dark:bg-gray-700 rounded-t-lg">
      <div class="flex items-center space-x-2">
        <div class="text-sm text-gray-600 dark:text-gray-300">
          <span v-if="cursorPosition">行 {{ cursorPosition.line }}, 列 {{ cursorPosition.column }}</span>
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400">
          {{ totalLines }} 行
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <n-button size="tiny" @click="toggleWordWrap" :type="wordWrap ? 'primary' : 'default'">
          <template #icon>
            <i :class="wordWrap ? 'fas fa-text-width' : 'fas fa-align-left'"></i>
          </template>
          自动换行
        </n-button>
        <n-button size="tiny" @click="findAndReplace">
          <template #icon>
            <i class="fas fa-search-replace"></i>
          </template>
          查找替换
        </n-button>
        <n-button size="tiny" @click="copyToClipboard">
          <template #icon>
            <i class="fas fa-copy"></i>
          </template>
          复制
        </n-button>
      </div>
    </div>

    <!-- 编辑器主体 -->
    <div class="editor-main flex rounded-b-lg border border-t-0 border-gray-300 dark:border-gray-600 overflow-hidden">
      <!-- 行号 -->
      <div class="line-numbers bg-gray-50 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-sm font-mono p-2 text-right select-none border-r border-gray-300 dark:border-gray-600">
        <div
          v-for="line in totalLines"
          :key="line"
          class="line-number h-5 leading-5"
          :class="{ 'bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-300': line === currentLine }"
        >
          {{ line }}
        </div>
      </div>

      <!-- 代码编辑区 -->
      <div class="editor-content flex-1">
        <textarea
          ref="textareaRef"
          v-model="content"
          @input="handleInput"
          @keydown="handleKeydown"
          @scroll="handleScroll"
          @click="updateCursorPosition"
          @keyup="updateCursorPosition"
          class="editor-textarea w-full h-full p-3 font-mono text-sm bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 resize-none outline-none border-none rounded-b-lg"
          :class="{ 'word-wrap': wordWrap }"
          :style="{ minHeight: '400px' }"
          spellcheck="false"
          placeholder="输入 JSON 配置内容..."
        ></textarea>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="errorMessage" class="error-message mt-2 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
      <div class="flex items-start space-x-2">
        <i class="fas fa-exclamation-triangle text-red-500 mt-0.5"></i>
        <div>
          <div class="text-sm font-medium text-red-800 dark:text-red-200">JSON 语法错误</div>
          <div class="text-xs text-red-600 dark:text-red-300 mt-1">{{ errorMessage }}</div>
        </div>
      </div>
    </div>

    <!-- 查找替换弹窗 -->
    <n-modal v-model:show="showFindReplace" preset="dialog" title="查找和替换">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">查找</label>
          <n-input
            v-model:value="findText"
            placeholder="输入要查找的文本..."
            @keydown.enter="findNext"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">替换</label>
          <n-input
            v-model:value="replaceText"
            placeholder="输入替换文本..."
            @keydown.enter="replaceNext"
          />
        </div>
        <div class="flex space-x-2">
          <n-button @click="findNext">查找下一个</n-button>
          <n-button @click="replaceNext">替换下一个</n-button>
          <n-button @click="replaceAll">全部替换</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

interface Props {
  modelValue: string
  placeholder?: string
  minHeight?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'validate', isValid: boolean, error?: string): void
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '输入 JSON 配置内容...',
  minHeight: '400px'
})

const emit = defineEmits<Emits>()

// 响应式数据
const content = ref(props.modelValue)
const wordWrap = ref(false)
const showFindReplace = ref(false)
const findText = ref('')
const replaceText = ref('')
const errorMessage = ref('')
const cursorPosition = ref({ line: 1, column: 1 })
const currentLine = ref(1)
const totalLines = computed(() => content.value.split('\n').length)

// DOM 引用
const textareaRef = ref<HTMLTextAreaElement>()
const highlightRef = ref<HTMLElement>()

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== content.value) {
    content.value = newValue
    validateJSON()
  }
})

// 监听内容变化
watch(content, (newValue) => {
  emit('update:modelValue', newValue)
  emit('change', newValue)
  validateJSON()
  updateCursorPosition()
})

// JSON 语法高亮
const highlightedContent = computed(() => {
  if (!content.value) return ''

  try {
    // 简单的 JSON 语法高亮
    let highlighted = content.value
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')

    // 高亮字符串
    highlighted = highlighted.replace(/"([^"\\]|\\.)*"/g, '<span class="string">"$1"</span>')

    // 高亮数字
    highlighted = highlighted.replace(/\b(\d+\.?\d*)\b/g, '<span class="number">$1</span>')

    // 高亮布尔值和 null
    highlighted = highlighted.replace(/\b(true|false|null)\b/g, '<span class="boolean">$1</span>')

    // 高亮键名
    highlighted = highlighted.replace(/"([^"\\]+)":\s*/g, '<span class="key">"$1"</span>: ')

    // 高亮括号
    highlighted = highlighted.replace(/([{}[\]])/g, '<span class="bracket">$1</span>')

    // 高亮逗号和冒号
    highlighted = highlighted.replace(/([,])/g, '<span class="comma">$1</span>')

    return highlighted
  } catch (e) {
    return content.value
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
  }
})

// JSON 验证
const validateJSON = () => {
  if (!content.value.trim()) {
    errorMessage.value = ''
    emit('validate', true)
    return
  }

  try {
    JSON.parse(content.value)
    errorMessage.value = ''
    emit('validate', true)
  } catch (error: any) {
    errorMessage.value = error.message || '无效的 JSON 格式'
    emit('validate', false, errorMessage.value)
  }
}

// 更新光标位置
const updateCursorPosition = () => {
  if (!textareaRef.value) return

  const textarea = textareaRef.value
  const text = textarea.value
  const cursorIndex = textarea.selectionStart

  // 计算行和列
  const lines = text.substring(0, cursorIndex).split('\n')
  currentLine.value = lines.length
  cursorPosition.value = {
    line: lines.length,
    column: lines[lines.length - 1].length + 1
  }
}

// 处理输入
const handleInput = () => {
  nextTick(() => {
    syncScroll()
  })
}

// 处理键盘事件
const handleKeydown = (event: KeyboardEvent) => {
  // Tab 键插入空格
  if (event.key === 'Tab') {
    event.preventDefault()
    const textarea = event.target as HTMLTextAreaElement
    const start = textarea.selectionStart
    const end = textarea.selectionEnd

    const newValue = content.value.substring(0, start) + '  ' + content.value.substring(end)
    content.value = newValue

    nextTick(() => {
      textarea.selectionStart = textarea.selectionEnd = start + 2
      updateCursorPosition()
    })
  }

  // Ctrl+S 保存
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault()
    emit('save')
  }
}

// 同步滚动
const handleScroll = () => {
  syncScroll()
}

const syncScroll = () => {
  if (!textareaRef.value || !highlightRef.value) return

  const textarea = textareaRef.value
  const highlight = highlightRef.value

  highlight.scrollTop = textarea.scrollTop
  highlight.scrollLeft = textarea.scrollLeft
}

// 切换自动换行
const toggleWordWrap = () => {
  wordWrap.value = !wordWrap.value
}

// 查找替换功能
const findAndReplace = () => {
  showFindReplace.value = true
}

const findNext = () => {
  if (!findText.value || !textareaRef.value) return

  const textarea = textareaRef.value
  const text = textarea.value
  const currentIndex = textarea.selectionEnd
  const index = text.indexOf(findText.value, currentIndex)

  if (index !== -1) {
    textarea.selectionStart = index
    textarea.selectionEnd = index + findText.value.length
    textarea.focus()
  } else {
    // 从头开始查找
    const firstIndex = text.indexOf(findText.value)
    if (firstIndex !== -1) {
      textarea.selectionStart = firstIndex
      textarea.selectionEnd = firstIndex + findText.value.length
      textarea.focus()
    }
  }
}

const replaceNext = () => {
  if (!findText.value || !textareaRef.value) return

  const textarea = textareaRef.value
  const text = textarea.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd

  if (text.substring(start, end) === findText.value) {
    const newText = text.substring(0, start) + replaceText.value + text.substring(end)
    content.value = newText

    nextTick(() => {
      textarea.selectionStart = start
      textarea.selectionEnd = start + replaceText.value.length
      findNext()
    })
  } else {
    findNext()
  }
}

const replaceAll = () => {
  if (!findText.value) return

  const newText = content.value.replace(new RegExp(findText.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g'), replaceText.value)
  content.value = newText
}

// 复制到剪贴板
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(content.value)
    // 这里可以添加成功提示
  } catch (error) {
    // 降级方案
    if (textareaRef.value) {
      textareaRef.value.select()
      document.execCommand('copy')
    }
  }
}

// 格式化 JSON
const formatJSON = () => {
  try {
    const parsed = JSON.parse(content.value)
    content.value = JSON.stringify(parsed, null, 2)
  } catch (error) {
    // JSON 格式不正确，不进行格式化
  }
}

// 压缩 JSON
const minifyJSON = () => {
  try {
    const parsed = JSON.parse(content.value)
    content.value = JSON.stringify(parsed)
  } catch (error) {
    // JSON 格式不正确，不进行压缩
  }
}

// 暴露方法给父组件
defineExpose({
  formatJSON,
  minifyJSON,
  validateJSON,
  copyToClipboard,
  findAndReplace
})
</script>

<style scoped>
.json-editor-container {
  @apply w-full;
}

.editor-textarea {
  position: relative;
  z-index: 2;
  background: transparent;
  resize: none;
  caret-color: text-color;
}

.editor-highlight {
  z-index: 1;
  color: transparent;
}

.word-wrap {
  white-space: pre-wrap !important;
  word-break: break-all !important;
}

/* 语法高亮样式 */
:deep(.string) {
  @apply text-green-600 dark:text-green-400;
}

:deep(.number) {
  @apply text-blue-600 dark:text-blue-400;
}

:deep(.boolean) {
  @apply text-purple-600 dark:text-purple-400;
}

:deep(.key) {
  @apply text-red-600 dark:text-red-400 font-medium;
}

:deep(.bracket) {
  @apply text-orange-600 dark:text-orange-400 font-bold;
}

:deep(.comma) {
  @apply text-gray-600 dark:text-gray-400;
}

.line-number {
  transition: background-color 0.15s ease-in-out;
}

/* 滚动条样式 */
.editor-textarea::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.editor-textarea::-webkit-scrollbar-track {
  @apply bg-gray-100 dark:bg-gray-800;
}

.editor-textarea::-webkit-scrollbar-thumb {
  @apply bg-gray-300 dark:bg-gray-600 rounded;
}

.editor-textarea::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400 dark:bg-gray-500;
}
</style>