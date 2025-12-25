<template>
  <Teleport to="body">
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-3xl w-full mx-4 max-h-[90vh] overflow-hidden">
        <!-- 模态框头部 -->
        <div class="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-900">配置插件</h2>
            <p class="text-sm text-gray-600 mt-1">{{ plugin.name }} v{{ plugin.version }}</p>
          </div>
          <button
            @click="$emit('close')"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <Icon name="x" class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <!-- 模态框内容 -->
        <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
          <!-- 配置编辑器 -->
          <div class="space-y-6">
            <!-- 配置说明 -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div class="flex items-start">
                <Icon name="info" class="w-5 h-5 text-blue-600 mt-0.5 mr-3 flex-shrink-0" />
                <div class="text-sm text-blue-800">
                  <p class="font-medium mb-1">配置说明</p>
                  <p>修改插件配置后点击保存即可生效。配置格式为JSON，请确保语法正确。</p>
                </div>
              </div>
            </div>

            <!-- JSON编辑器 -->
            <div>
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">
                  插件配置
                </label>
                <div class="flex items-center space-x-2">
                  <!-- 格式化按钮 -->
                  <button
                    @click="formatConfig"
                    class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors flex items-center"
                  >
                    <Icon name="code" class="w-4 h-4 mr-1" />
                    格式化
                  </button>
                  <!-- 验证按钮 -->
                  <button
                    @click="validateConfig"
                    class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors flex items-center"
                  >
                    <Icon name="check" class="w-4 h-4 mr-1" />
                    验证
                  </button>
                  <!-- 重置按钮 -->
                  <button
                    @click="resetConfig"
                    class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors flex items-center"
                  >
                    <Icon name="refresh-cw" class="w-4 h-4 mr-1" />
                    重置
                  </button>
                </div>
              </div>

              <div class="relative">
                <textarea
                  v-model="configText"
                  @input="onConfigChange"
                  class="w-full h-96 p-4 font-mono text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="请输入JSON格式的配置..."
                  spellcheck="false"
                ></textarea>

                <!-- 行号显示 -->
                <div class="absolute left-0 top-0 w-12 h-full bg-gray-50 border-r border-gray-300 rounded-l-lg p-4 text-xs text-gray-500 font-mono leading-6 pointer-events-none">
                  <div v-for="n in lineCount" :key="n" class="text-right">{{ n }}</div>
                </div>
              </div>

              <!-- 错误提示 -->
              <div v-if="configError" class="mt-2 p-3 bg-red-50 border border-red-200 rounded-lg">
                <div class="flex items-start">
                  <Icon name="alert-circle" class="w-5 h-5 text-red-600 mt-0.5 mr-2 flex-shrink-0" />
                  <div class="text-sm text-red-800">
                    <p class="font-medium">配置错误</p>
                    <p>{{ configError }}</p>
                  </div>
                </div>
              </div>

              <!-- 验证成功提示 -->
              <div v-if="isValid && !configError" class="mt-2 p-3 bg-green-50 border border-green-200 rounded-lg">
                <div class="flex items-center">
                  <Icon name="check-circle" class="w-5 h-5 text-green-600 mr-2" />
                  <span class="text-sm text-green-800">配置格式正确</span>
                </div>
              </div>
            </div>

            <!-- 配置预览 -->
            <div v-if="parsedConfig">
              <h3 class="text-lg font-medium text-gray-900 mb-3">配置预览</h3>
              <div class="bg-gray-50 rounded-lg p-4">
                <pre class="text-sm text-gray-800 whitespace-pre-wrap">{{ JSON.stringify(parsedConfig, null, 2) }}</pre>
              </div>
            </div>

            <!-- 配置模式切换 -->
            <div>
              <h3 class="text-lg font-medium text-gray-900 mb-3">快速配置</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- 常用配置项 -->
                <div v-for="(preset, name) in configPresets" :key="name" class="border border-gray-200 rounded-lg p-4">
                  <h4 class="font-medium text-gray-900 mb-2">{{ name }}</h4>
                  <p class="text-sm text-gray-600 mb-3">{{ preset.description }}</p>
                  <button
                    @click="applyPreset(preset.config)"
                    class="w-full px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
                  >
                    应用预设
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 模态框底部 -->
        <div class="border-t border-gray-200 px-6 py-4 flex justify-between items-center">
          <div class="text-sm text-gray-600">
            <span v-if="hasChanges" class="text-orange-600 font-medium">
              <Icon name="alert-circle" class="w-4 h-4 inline mr-1" />
              有未保存的更改
            </span>
            <span v-else>
              配置未修改
            </span>
          </div>
          <div class="flex space-x-3">
            <button
              @click="$emit('close')"
              class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              取消
            </button>
            <button
              @click="saveConfig"
              :disabled="!isValid || configError || saving"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
            >
              <Icon v-if="saving" name="loader" class="w-4 h-4 mr-2 animate-spin" />
              {{ saving ? '保存中...' : '保存配置' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useToast } from '~/composables/useToast'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'saved'])

const toast = useToast()

// 响应式数据
const configText = ref('')
const originalConfig = ref('')
const parsedConfig = ref(null)
const configError = ref('')
const isValid = ref(false)
const saving = ref(false)

// 计算属性
const lineCount = computed(() => {
  return configText.value.split('\n').length
})

const hasChanges = computed(() => {
  return configText.value !== originalConfig.value
})

// 配置预设
const configPresets = ref({
  '开发模式': {
    description: '启用调试和详细日志',
    config: {
      enabled: true,
      debug: true,
      log_level: 'debug',
      verbose_logging: true
    }
  },
  '生产模式': {
    description: '优化性能，减少日志',
    config: {
      enabled: true,
      debug: false,
      log_level: 'error',
      verbose_logging: false,
      performance_mode: true
    }
  },
  '最小配置': {
    description: '仅启用基本功能',
    config: {
      enabled: true
    }
  },
  '全功能': {
    description: '启用所有功能',
    config: {
      enabled: true,
      debug: false,
      log_level: 'info',
      verbose_logging: true,
      performance_mode: false,
      advanced_features: true,
      monitoring: true,
      auto_update: false
    }
  }
})

// 方法
const loadConfig = async () => {
  try {
    const response = await $fetch(`/api/plugins/${props.plugin.name}`)
    if (response.success && response.data.plugin.config) {
      const config = response.data.plugin.config
      configText.value = JSON.stringify(config, null, 2)
      originalConfig.value = configText.value
      parsedConfig.value = config
      validateConfig()
    } else {
      // 使用默认配置
      const defaultConfig = {
        enabled: true,
        log_level: 'info'
      }
      configText.value = JSON.stringify(defaultConfig, null, 2)
      originalConfig.value = configText.value
      parsedConfig.value = defaultConfig
      isValid.value = true
    }
  } catch (error) {
    toast.error('加载配置失败: ' + error.message)
  }
}

const onConfigChange = () => {
  validateConfig()
}

const validateConfig = () => {
  try {
    const config = JSON.parse(configText.value)
    parsedConfig.value = config
    configError.value = ''
    isValid.value = true
  } catch (error) {
    configError.value = error.message
    parsedConfig.value = null
    isValid.value = false
  }
}

const formatConfig = () => {
  try {
    const config = JSON.parse(configText.value)
    configText.value = JSON.stringify(config, null, 2)
    validateConfig()
    toast.success('配置已格式化')
  } catch (error) {
    toast.error('格式化失败: ' + error.message)
  }
}

const resetConfig = () => {
  configText.value = originalConfig.value
  validateConfig()
  toast.info('配置已重置')
}

const applyPreset = (preset) => {
  configText.value = JSON.stringify(preset, null, 2)
  validateConfig()
  toast.success('预设配置已应用')
}

const saveConfig = async () => {
  if (!isValid.value || configError.value) {
    toast.error('请先修正配置错误')
    return
  }

  try {
    saving.value = true
    const config = JSON.parse(configText.value)

    const response = await $fetch(`/api/plugins/${props.plugin.name}/config`, {
      method: 'PUT',
      body: {
        config: config
      }
    })

    if (response.success) {
      originalConfig.value = configText.value
      emit('saved')
      toast.success('配置保存成功')
    } else {
      throw new Error(response.error || '保存失败')
    }
  } catch (error) {
    toast.error('保存配置失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

// 生命周期
onMounted(() => {
  loadConfig()
})

// 监听插件变化
watch(() => props.plugin.name, () => {
  loadConfig()
})
</script>

<style scoped>
/* JSON编辑器样式 */
textarea {
  padding-left: 5rem; /* 为行号留出空间 */
  tab-size: 2;
}

/* 行号样式 */
.absolute > div {
  line-height: 1.5; /* 与textarea的行高匹配 */
}

/* 保存按钮动画 */
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* 模态框动画 */
.fixed > div {
  animation: modalFadeIn 0.3s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>