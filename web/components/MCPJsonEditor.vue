<template>
  <div class="mcp-json-editor">
    <!-- 工具栏 -->
    <div class="editor-toolbar flex items-center justify-between mb-2">
      <div class="flex items-center gap-2">
        <n-button size="small" @click="formatJSON">
          <template #icon>
            <n-icon><Code /></n-icon>
          </template>
          格式化
        </n-button>
        <n-button size="small" @click="compactJSON">
          <template #icon>
            <n-icon><Minimize /></n-icon>
          </template>
          压缩
        </n-button>
        <n-button size="small" @click="copyToClipboard">
          <template #icon>
            <n-icon><Copy /></n-icon>
          </template>
          复制
        </n-button>
        <n-button size="small" @click="validateConfig">
          <template #icon>
            <n-icon><Check /></n-icon>
          </template>
          验证
        </n-button>
      </div>

      <div class="flex items-center gap-2">
        <n-text depth="3" class="text-xs">
          {{ cursorPosition }}
        </n-text>
        <n-switch v-model:value="wordWrap" size="small">
          <template #checked>换行</template>
          <template #unchecked>不换行</template>
        </n-switch>
      </div>
    </div>

    <!-- 编辑器容器 -->
    <div class="editor-container relative">
      <codemirror
        v-model="content"
        :style="{ height: height }"
        :extensions="extensions"
        :autofocus="autofocus"
        :indent-with-tab="true"
        :tab-size="2"
        @ready="handleReady"
        @change="handleChange"
        @focus="handleFocus"
        @blur="handleBlur"
        @update="handleUpdate"
      />

      <!-- 验证状态指示器 -->
      <div v-if="validationState.show" class="validation-status absolute top-2 right-2">
        <n-tag
          :type="validationState.valid ? 'success' : 'error'"
          size="small"
          :bordered="false"
        >
          {{ validationState.message }}
        </n-tag>
      </div>
    </div>

    <!-- 错误和警告信息 -->
    <div v-if="errors.length > 0 || warnings.length > 0" class="mt-2">
      <n-alert v-if="errors.length > 0" type="error" title="错误" class="mb-2">
        <ul class="list-disc list-inside text-sm">
          <li v-for="error in errors" :key="error">{{ error }}</li>
        </ul>
      </n-alert>

      <n-alert v-if="warnings.length > 0" type="warning" title="警告">
        <ul class="list-disc list-inside text-sm">
          <li v-for="warning in warnings" :key="warning">{{ warning }}</li>
        </ul>
      </n-alert>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { json } from '@codemirror/lang-json'
import { linter } from '@codemirror/lint'
import { search } from '@codemirror/search'
import { keymap, EditorView } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { NButton, NIcon, NText, NSwitch, NTag, NAlert } from 'naive-ui'
import { Code, Minimize, Copy, Check } from '@vicons/tabler'
import { validateMCPConfig } from '~/schemas/mcp-schema'

interface Props {
  modelValue: string
  height?: string
  autofocus?: boolean
  placeholder?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'validate', valid: boolean, errors?: string[], warnings?: string[]): void
  (e: 'change', value: string, view: any): void
  (e: 'ready', view: any): void
  (e: 'focus', view: any): void
  (e: 'blur', view: any): void
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  autofocus: false,
  placeholder: '请输入MCP配置JSON...'
})

const emit = defineEmits<Emits>()

// 响应式数据
const content = ref(props.modelValue)
const wordWrap = ref(true)
const cursorPosition = ref('行 1, 列 1')
const view = ref<any>(null)

// 验证状态
const validationState = ref({
  show: false,
  valid: false,
  message: ''
})

const errors = ref<string[]>([])
const warnings = ref<string[]>([])

// MCP配置验证器
const mcpLinter = linter((view) => {
  const content = view.state.doc.toString()
  if (!content.trim()) return []

  const diagnostics = []

  // JSON语法验证
  try {
    const config = JSON.parse(content)
    const validation = validateMCPConfig(config)

    if (!validation.valid) {
      validation.errors.forEach((error, index) => {
        diagnostics.push({
          from: 0,
          to: content.length,
          severity: 'error',
          message: error,
          source: 'mcp-validator'
        })
      })
    }

    if (validation.warnings.length > 0) {
      validation.warnings.forEach((warning, index) => {
        diagnostics.push({
          from: 0,
          to: content.length,
          severity: 'warning',
          message: warning,
          source: 'mcp-validator'
        })
      })
    }

    // 更新验证状态
    errors.value = validation.errors
    warnings.value = validation.warnings
    emit('validate', validation.valid, validation.errors, validation.warnings)

  } catch (error: any) {
    diagnostics.push({
      from: 0,
      to: content.length,
      severity: 'error',
      message: `JSON语法错误: ${error.message}`,
      source: 'json-parser'
    })

    errors.value = [error.message]
    warnings.value = []
    emit('validate', false, [error.message], [])
  }

  return diagnostics
})

// CodeMirror扩展
const extensions = computed(() => [
  json(),
  mcpLinter,
  search(),
  EditorView.theme({
    '&': {
      fontSize: '14px',
      height: '100%'
    },
    '.cm-content': {
      padding: '12px',
      minHeight: '100%'
    },
    '.cm-editor': {
      borderRadius: '6px',
      border: '1px solid var(--n-border-color)',
      height: '100%',
      overflow: 'hidden'
    },
    '.cm-editor.cm-focused': {
      borderColor: 'var(--n-primary-color)',
      boxShadow: '0 0 0 2px rgba(24, 160, 88, 0.2)'
    },
    '.cm-scroller': {
      overflow: 'auto',
      fontFamily: 'inherit',
      height: '100%'
    },
    '.cm-gutters': {
      backgroundColor: 'var(--n-color-embedded)',
      borderRight: '1px solid var(--n-border-color)'
    },
    '.cm-lineNumbers .cm-gutterElement': {
      color: 'var(--n-text-color-3)'
    },
    '.cm-activeLineGutter': {
      backgroundColor: 'var(--n-color-embedded-modal)'
    },
    '.cm-cursor': {
      borderLeftColor: 'var(--n-primary-color)'
    },
    '.cm-selectionBackground': {
      backgroundColor: 'var(--n-primary-color-hover)'
    }
  })
])

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== content.value) {
    content.value = newValue
  }
})

watch(content, (newValue) => {
  emit('update:modelValue', newValue)
})

// 方法
const formatJSON = () => {
  try {
    const parsed = JSON.parse(content.value)
    content.value = JSON.stringify(parsed, null, 2)
    showValidationMessage('格式化成功', true)
  } catch (error: any) {
    showValidationMessage('JSON格式错误，无法格式化', false)
  }
}

const compactJSON = () => {
  try {
    const parsed = JSON.parse(content.value)
    content.value = JSON.stringify(parsed)
    showValidationMessage('压缩成功', true)
  } catch (error: any) {
    showValidationMessage('JSON格式错误，无法压缩', false)
  }
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(content.value)
    showValidationMessage('已复制到剪贴板', true)
  } catch (error) {
    showValidationMessage('复制失败', false)
  }
}

const validateConfig = () => {
  try {
    const config = JSON.parse(content.value)
    const validation = validateMCPConfig(config)
    showValidationMessage(
      validation.valid ? '配置验证通过' : '配置验证失败',
      validation.valid
    )
  } catch (error: any) {
    showValidationMessage('JSON格式错误', false)
  }
}

const showValidationMessage = (message: string, valid: boolean) => {
  validationState.value = {
    show: true,
    valid,
    message
  }

  setTimeout(() => {
    validationState.value.show = false
  }, 3000)
}

// 事件处理器
const handleReady = (payload: any) => {
  view.value = payload.view
  emit('ready', payload.view)

  // 更新光标位置
  updateCursorPosition()
}

const handleChange = (value: string, viewUpdate: any) => {
  emit('change', value, viewUpdate)
  updateCursorPosition()
}

const handleFocus = (payload: any) => {
  emit('focus', payload.view)
}

const handleBlur = (payload: any) => {
  emit('blur', payload.view)
}

const handleUpdate = (viewUpdate: any) => {
  updateCursorPosition()
}

const updateCursorPosition = () => {
  if (view.value) {
    const pos = view.value.state.selection.main.head
    const line = view.value.state.doc.lineAt(pos)
    cursorPosition.value = `行 ${line.number}, 列 ${pos - line.from + 1}`
  }
}

// 生命周期
onMounted(() => {
  // 初始验证
  if (content.value.trim()) {
    validateConfig()
  }
})
</script>

<style scoped>
.mcp-json-editor {
  height: 100%;
  overflow: hidden;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.editor-toolbar {
  flex-shrink: 0;
}

.editor-container {
  flex: 1;
  position: relative;
  border-radius: 6px;
  overflow: hidden;
  min-height: 0;
  max-height: 500px;
}

.validation-status {
  pointer-events: none;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  :deep(.cm-editor) {
    background-color: var(--n-color-embedded);
    color: var(--n-text-color);
  }
}
</style>