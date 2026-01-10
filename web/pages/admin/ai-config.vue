<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和保存按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">AI 配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理AI服务配置和设置</p>
      </div>
      <div class="flex space-x-2">
        <n-button @click="testAIConnectionClick" :loading="testingConnection">
          <template #icon>
            <i class="fas fa-plug"></i>
          </template>
          测试连接
        </n-button>
        <n-button type="primary" @click="saveConfig" :loading="saving">
          <template #icon>
            <i class="fas fa-save"></i>
          </template>
          保存配置
        </n-button>
      </div>
    </template>

    <!-- 内容区 - 配置表单 -->
    <template #content>
      <div class="config-content h-full">
        <!-- 顶部Tabs -->
        <n-tabs
          v-model:value="activeTab"
          type="line"
          animated
        >
          <n-tab-pane name="openai" tab="OpenAI 配置">
            <div class="tab-content-container">
              <n-form
                ref="formRef"
                :model="configForm"
                :rules="rules"
                label-placement="left"
                label-width="auto"
                require-mark-placement="right-hanging"
              >
                <div class="space-y-6">
                  <!-- API Key -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">API Key</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">AI服务的API密钥</span>
                    </div>
                    <div class="flex space-x-2">
                      <n-input
                        v-model:value="configForm.api_key"
                        type="password"
                        placeholder="请输入API Key"
                        show-password-on="click"
                        class="flex-1"
                      />
                      <n-button
                        @click="testAIConnectionClick"
                        :loading="testingConnection"
                        type="primary"
                        ghost
                        size="medium"
                      >
                        <template #icon>
                          <i class="fas fa-plug"></i>
                        </template>
                        验证
                      </n-button>
                    </div>
                    <div class="flex items-center space-x-2">
                      <div v-if="configForm.api_key_configured && configForm.api_key" class="text-sm text-green-600 dark:text-green-400">
                        <i class="fas fa-check-circle mr-1"></i> API Key 已配置
                      </div>
                      <div v-if="configForm.api_key && !configForm.api_key_configured" class="text-sm text-yellow-600 dark:text-yellow-400">
                        <i class="fas fa-exclamation-triangle mr-1"></i> 请验证新的 API Key
                      </div>
                      <div v-if="!configForm.api_key" class="text-sm text-gray-500 dark:text-gray-400">
                        <i class="fas fa-info-circle mr-1"></i> 请输入 API Key
                      </div>
                    </div>
                  </div>

                  <!-- API URL -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">API URL</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">AI服务的API地址，例如：https://api.openai.com/v1</span>
                    </div>
                    <n-input
                      v-model:value="configForm.api_url"
                      placeholder="请输入API URL，如：https://api.openai.com/v1"
                    />
                  </div>

                  <!-- 模型 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">模型</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">用于AI处理的模型名称，例如：gpt-3.5-turbo</span>
                    </div>
                    <n-input
                      v-model:value="configForm.model"
                      placeholder="请输入模型名称，如：gpt-3.5-turbo"
                    />
                  </div>

                  <!-- 最大令牌数 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">最大令牌数</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">单次请求的最大令牌数</span>
                    </div>
                    <n-input-number
                      v-model:value="configForm.max_tokens"
                      :min="1"
                      :max="4096"
                      placeholder="请输入最大令牌数"
                      class="w-full"
                    />
                  </div>

                  <!-- 温度 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">温度</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">控制AI输出的随机性，范围0.0-2.0</span>
                    </div>
                    <n-slider
                      v-model:value="configForm.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                      class="w-full"
                    />
                    <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400">
                      <span>确定性</span>
                      <span class="font-medium">{{ configForm.temperature.toFixed(1) }}</span>
                      <span>创造性</span>
                    </div>
                  </div>

                  
                  <!-- 超时时间 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">超时时间</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">API请求的超时时间（秒）</span>
                    </div>
                    <n-input-number
                      v-model:value="configForm.timeout"
                      :min="1"
                      :max="300"
                      placeholder="请输入超时时间（秒）"
                      class="w-full"
                    />
                  </div>

                  <!-- 重试次数 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">重试次数</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">API请求失败时的重试次数</span>
                    </div>
                    <n-input-number
                      v-model:value="configForm.retry_count"
                      :min="0"
                      :max="10"
                      placeholder="请输入重试次数"
                      class="w-full"
                    />
                  </div>
                </div>
              </n-form>
            </div>
          </n-tab-pane>

          <n-tab-pane name="mcp" tab="MCP 配置">
            <div class="tab-content-container h-[600px]">
              <div class="flex h-full gap-6">
                <!-- 左侧：MCP 服务状态 -->
                <div class="w-1/3 space-y-4 flex flex-col">
                  <div class="flex items-center space-x-2 flex-shrink-0">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">MCP 服务状态</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">当前运行的MCP服务</span>
                  </div>
                  <div class="border rounded-lg p-4 flex-1 overflow-y-auto bg-white dark:bg-gray-800">
                    <div v-if="mcpServices.length === 0" class="text-gray-500 text-center py-8">
                      暂无MCP服务
                    </div>
                    <div v-else class="space-y-3">
                      <div v-for="service in mcpServices" :key="service.name" class="border rounded-lg p-4 bg-gray-50 dark:bg-gray-700">
                        <div class="flex items-center justify-between mb-3">
                          <div class="flex items-center space-x-2">
                            <div class="w-3 h-3 rounded-full" :class="{
                              'bg-green-400': service.status === 'running',
                              'bg-red-400': service.status === 'stopped',
                              'bg-yellow-400': service.status === 'starting' || service.status === 'stopping'
                            }"></div>
                            <div class="font-medium text-gray-900 dark:text-gray-100">{{ service.name }}</div>
                          </div>
                        </div>
                        <div class="text-sm text-gray-600 dark:text-gray-400 mb-3">
                          状态: {{ service.status }}
                        </div>
                        <div class="flex space-x-2">
                          <n-button size="small" type="primary" @click="startMCPService(service.name)" :disabled="service.status === 'running'">
                            <template #icon>
                              <i class="fas fa-play"></i>
                            </template>
                            启动
                          </n-button>
                          <n-button size="small" @click="stopMCPService(service.name)" :disabled="service.status === 'stopped'">
                            <template #icon>
                              <i class="fas fa-stop"></i>
                            </template>
                            停止
                          </n-button>
                          <n-button size="small" @click="restartMCPService(service.name)">
                            <template #icon>
                              <i class="fas fa-redo"></i>
                            </template>
                            重启
                          </n-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 右侧：MCP 配置编辑器 -->
                <div class="w-2/3 space-y-4 flex flex-col">
                  <div class="flex items-center justify-between flex-shrink-0">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">MCP 配置文件</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">MCP服务配置文件编辑器</span>
                    </div>
                    <div class="flex space-x-2">
                      <!-- 配置模板下拉菜单 -->
                      <n-dropdown trigger="click" :options="templateOptions" @select="applyTemplate">
                        <n-button size="small">
                          <template #icon>
                            <i class="fas fa-file-alt"></i>
                          </template>
                          配置模板
                        </n-button>
                      </n-dropdown>
                      <n-button @click="loadMCPConfig" :loading="loadingMCPConfig" size="small">
                        <template #icon>
                          <i class="fas fa-sync"></i>
                        </template>
                        重新加载
                      </n-button>
                    </div>
                  </div>

                  <!-- JSON 编辑器组件 -->
                  <div class="flex-1 border rounded-lg overflow-hidden min-h-0">
                    <JsonEditorSimple
                      v-model="mcpConfigContent"
                      @validate="onEditorValidate"
                      @change="onEditorChange"
                      :min-height="'400px'"
                    />
                  </div>

                  <div class="flex items-center justify-between flex-shrink-0">
                    <div class="text-sm text-gray-600 dark:text-gray-400">
                      <span v-if="mcpConfigValid" class="text-green-600 dark:text-green-400">
                        <i class="fas fa-check-circle mr-1"></i> JSON格式正确
                      </span>
                      <span v-else class="text-red-600 dark:text-red-400">
                        <i class="fas fa-exclamation-circle mr-1"></i> JSON格式错误
                      </span>
                    </div>
                    <div class="flex space-x-2">
                      <n-button type="primary" @click="saveMCPConfig" :loading="savingMCPConfig" :disabled="!mcpConfigValid">
                        <template #icon>
                          <i class="fas fa-save"></i>
                        </template>
                        保存MCP配置
                      </n-button>
                      <n-button @click="testMCPConfig">
                        <template #icon>
                          <i class="fas fa-check-circle"></i>
                        </template>
                        验证配置
                      </n-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </n-tab-pane>
        </n-tabs>
      </div>
    </template>
  </AdminPageLayout>

  <!-- AI 测试弹窗 -->
  <n-modal v-model:show="showTestModal" :mask-closable="false" preset="card" :style="{ width: '900px', maxWidth: '90vw' }">
    <template #header>
      <div class="flex items-center space-x-2">
        <i class="fas fa-plug text-blue-600"></i>
        <span>AI 连接测试详情</span>
      </div>
    </template>

    <div class="space-y-6 max-h-[70vh] overflow-y-auto">
      <!-- 发送的问题 -->
      <div>
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-blue-500 rounded-full animate-pulse"></div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">发送给 AI 的问题</h3>
          </div>
          <n-button size="small" @click="copyToClipboard(testData.sent.prompt)">
            <template #icon>
              <i class="fas fa-copy"></i>
            </template>
            复制问题
          </n-button>
        </div>
        <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <p class="text-base text-blue-800 dark:text-blue-200 leading-relaxed">{{ testData.sent.prompt }}</p>
        </div>
        <div class="mt-2">
          <details class="cursor-pointer">
            <summary class="text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300">
              <i class="fas fa-cog mr-1"></i>查看配置参数
            </summary>
            <div class="mt-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
              <pre class="text-xs font-mono text-gray-600 dark:text-gray-300 whitespace-pre-wrap break-all">{{ JSON.stringify(testData.sent.config, null, 2) }}</pre>
            </div>
          </details>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="testData.loading" class="flex flex-col items-center justify-center py-12 space-y-4">
        <n-spin size="large" />
        <div class="text-center space-y-2">
          <span class="text-gray-600 dark:text-gray-400 font-medium">正在测试 AI 连接...</span>
          <div class="text-xs text-gray-500 dark:text-gray-500">请稍候，正在发送测试请求</div>
        </div>
      </div>

      <!-- AI 的回答 -->
      <div v-else-if="testData.received && testData.received.response" class="mt-4">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-blue-500 rounded-full"></div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">AI 的详细回答</h3>
              <div class="text-xs text-gray-500 dark:text-gray-400">获取模型的详细信息、版本和能力说明</div>
            </div>
          </div>
          <div class="flex space-x-2">
            <n-button size="small" @click="copyToClipboard(testData.received.response)">
              <template #icon>
                <i class="fas fa-copy"></i>
              </template>
              复制回答
            </n-button>
            <n-button size="small" @click="formatAIResponse">
              <template #icon>
                <i class="fas fa-align-left"></i>
              </template>
              格式化
            </n-button>
          </div>
        </div>
          <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <div class="text-base text-gray-800 dark:text-gray-200 leading-relaxed whitespace-pre-wrap">
            {{ formattedAIResponse }}
          </div>
        </div>
        <div v-if="testData.received && testData.received.response !== formattedAIResponse" class="mt-2 text-xs text-gray-500 dark:text-gray-400">
          <i class="fas fa-info-circle mr-1"></i>
          回答已格式化显示，点击"复制回答"获取原始内容
        </div>
        <div class="mt-4 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <i class="fas fa-check-circle text-green-500 text-xl"></i>
            <span class="text-green-600 dark:text-green-400 font-medium">测试成功！AI 连接正常</span>
          </div>
          <n-tag type="success" size="small">
            响应时间: {{ new Date().toLocaleTimeString() }}
          </n-tag>
        </div>
      </div>

      <!-- 错误信息 -->
      <div v-else-if="!testData.loading && testData.error">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-red-500 rounded-full"></div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">错误信息</h3>
          </div>
          <n-button size="small" @click="copyToClipboard(testData.error)">
            <template #icon>
              <i class="fas fa-copy"></i>
            </template>
            复制
          </n-button>
        </div>
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p class="text-sm font-mono text-red-800 dark:text-red-200">{{ testData.error }}</p>
        </div>
        <div class="mt-4 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <i class="fas fa-exclamation-triangle text-red-500 text-xl"></i>
            <span class="text-red-600 dark:text-red-400 font-medium">测试失败，请检查配置</span>
          </div>
          <n-tag type="error" size="small">
            错误时间: {{ new Date().toLocaleTimeString() }}
          </n-tag>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-between items-center">
        <div class="text-sm text-gray-500 dark:text-gray-400">
          <i class="fas fa-info-circle mr-1"></i>
          测试将使用当前表单中的配置参数，不会影响已保存的配置
        </div>
        <div class="flex space-x-2">
          <n-button @click="showTestModal = false">
            关闭
          </n-button>
          <n-button v-if="!testData.loading && testData.received" type="primary" @click="showTestModal = false">
            确认
          </n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  ssr: false
})
import { ref, reactive, onMounted, computed } from 'vue'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'
import { useAIApi, useMCPApi } from '~/composables/useApi'
import JsonEditorSimple from '~/components/JsonEditorSimple.vue'

// MCP 配置模板
const templateOptions = [
  {
    label: '基础模板 - 空配置',
    key: 'basic',
    icon: 'i-carbon-document-blank'
  },
  {
    label: 'Web搜索模板 - Bing搜索',
    key: 'web-search',
    icon: 'i-carbon-search'
  },
  {
    label: '文件系统模板 - 本地文件',
    key: 'filesystem',
    icon: 'i-carbon-folder'
  },
  {
    label: '完整模板 - 所有服务',
    key: 'full',
    icon: 'i-carbon-document-add'
  }
]
const activeTab = ref('openai')
const saving = ref(false)
const testingConnection = ref(false)
const loadingMCPConfig = ref(false)
const savingMCPConfig = ref(false)
const formRef = ref()
const notification = useNotification()

// 测试弹窗相关
const showTestModal = ref(false)
const testData = ref({
  sent: {} as any,
  received: null as any,
  loading: false,
  error: null as string | null
})

// 格式化后的 AI 回答
const formattedAIResponse = computed(() => {
  if (!testData.value.received?.response) return ''

  const response = testData.value.received.response

  // 移除开头的换行符
  let formatted = response.replace(/^\n+/, '')

  // 移除结尾的换行符
  formatted = formatted.replace(/\n+$/, '')

  // 保留完整的回答，不再截断
  return formatted
})

// 格式化 AI 回答（切换显示模式）
const formatAIResponse = () => {
  // 这里可以添加更多格式化逻辑，目前使用计算属性自动格式化
  notification.info({
    content: 'AI 回答已自动格式化显示',
    duration: 2000
  })
}

// AI配置表单数据
interface AIConfigForm {
  api_key?: string
  api_key_configured: boolean
  api_url: string
  model: string
  max_tokens: number
  temperature: number
  timeout: number
  retry_count: number
}

// 配置表单数据
const configForm = ref<AIConfigForm>({
  api_key_configured: false,
  api_url: 'https://api.openai.com/v1',
  model: 'gpt-3.5-turbo',
  max_tokens: 1000,
  temperature: 0.7,
  timeout: 30,
  retry_count: 3
})

// MCP配置
const mcpConfigContent = ref('')
const mcpServices = ref<any[]>([])

// 表单验证规则
const rules = {
  api_url: {
    required: false,
    trigger: ['blur', 'input'],
    message: '请输入API URL'
  },
  model: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入模型名称'
  }
}

// 获取AI配置API
const { getAIConfig, updateAIConfig, testAIConnection } = useAIApi()
const { getMCPConfig, updateMCPConfig, listMCPServices, startMCPService: startMCP, stopMCPService: stopMCP, restartMCPService: restartMCP } = useMCPApi()

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  getChangedDetails,
  updateOriginalConfig
} = useConfigChangeDetection<AIConfigForm>({
  debug: true,
  // 自定义比较函数
  customCompare: (key: string, currentValue: any, originalValue: any) => {
    return currentValue !== originalValue
  }
})

// 获取AI配置
const loadAIConfig = async () => {
  try {
    const response = await getAIConfig()
    if (response) {
      // 更新表单数据，使用默认值避免undefined
      const configData: AIConfigForm = {
        api_key: response.api_key || undefined,
        api_key_configured: response.ai_api_key_configured || false,
        api_url: response.ai_api_url || 'https://api.openai.com/v1',
        model: response.ai_model || 'gpt-3.5-turbo',
        max_tokens: response.ai_max_tokens ? parseInt(response.ai_max_tokens.toString()) : 1000,
        temperature: response.ai_temperature ? parseFloat(response.ai_temperature.toString()) : 0.7,
        timeout: response.ai_timeout ? parseInt(response.ai_timeout.toString()) : 30,
        retry_count: response.ai_retry_count ? parseInt(response.ai_retry_count.toString()) : 3
      }

      configForm.value = { ...configData }
      setOriginalConfig(configData)
    }
  } catch (error) {
    console.error('获取AI配置失败:', error)
    notification.error({
      content: '获取AI配置失败',
      duration: 3000
    })
  }
}

// 获取MCP配置
const loadMCPConfig = async () => {
  try {
    loadingMCPConfig.value = true
    console.log('正在加载MCP配置...')
    const response = await getMCPConfig()
    console.log('MCP配置响应:', response)

    if (response && response.config) {
      mcpConfigContent.value = response.config
      console.log('已加载MCP配置:', response.config)
      notification.success({
        content: 'MCP配置加载成功',
        duration: 2000
      })
    } else {
      console.log('API未返回配置，使用默认配置')
      // 如果API没有返回配置，则使用默认配置
      mcpConfigContent.value = '{\n' +
  '  "mcpServers": {\n' +
  '    "web-search": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-web-search"],\n' +
  '      "env": {\n' +
  '        "BING_API_KEY": "${BING_API_KEY}"\n' +
  '      },\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": true\n' +
  '    },\n' +
  '    "filesystem": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": false\n' +
  '    }\n' +
  '  }\n' +
  '}'
      notification.info({
        content: '已加载默认MCP配置',
        duration: 2000
      })
    }
  } catch (error) {
    console.error('获取MCP配置失败:', error)
    notification.error({
      content: '获取MCP配置失败',
      duration: 3000
    })
    // 如果API失败，也提供一个默认配置
    mcpConfigContent.value = '{\n' +
  '  "mcpServers": {\n' +
  '    "web-search": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-web-search"],\n' +
  '      "env": {\n' +
  '        "BING_API_KEY": "${BING_API_KEY}"\n' +
  '      },\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": true\n' +
  '    },\n' +
  '    "filesystem": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": false\n' +
  '    }\n' +
  '  }\n' +
  '}'
  } finally {
    loadingMCPConfig.value = false
  }
}

// 获取MCP服务状态
const loadMCPStatus = async () => {
  try {
    const services = await listMCPServices()
    mcpServices.value = services || []
  } catch (error) {
    console.error('获取MCP服务状态失败:', error)
    // 不显示错误，因为MCP可能未配置
  }
}

// 保存AI配置
const saveConfig = async () => {
  try {
    await formRef.value?.validate()

    saving.value = true

    // 更新当前配置数据
    updateCurrentConfig(configForm.value)

    // 准备要发送的数据
    const updateData: any = {
      ai_api_url: configForm.value.api_url,
      ai_model: configForm.value.model,
      ai_max_tokens: configForm.value.max_tokens,
      ai_temperature: configForm.value.temperature,
      ai_timeout: configForm.value.timeout,
      ai_retry_count: configForm.value.retry_count
    }

    // 只有当用户输入了新API Key时才更新
    if (configForm.value.api_key && configForm.value.api_key.trim() !== '') {
      updateData.api_key = configForm.value.api_key
    }

    // 调用API更新配置
    await updateAIConfig(updateData)

    notification.success({
      content: 'AI配置保存成功',
      duration: 3000
    })

    // 更新原始配置数据
    updateOriginalConfig(configForm.value)

    // 重新加载配置以确保数据同步
    await loadAIConfig()
  } catch (error) {
    console.error('保存AI配置失败:', error)
    notification.error({
      content: '保存AI配置失败',
      duration: 3000
    })
  } finally {
    saving.value = false
  }
}

// 保存MCP配置
const saveMCPConfig = async () => {
  try {
    savingMCPConfig.value = true

    // 验证JSON格式
    try {
      JSON.parse(mcpConfigContent.value)
    } catch (e) {
      notification.error({
        content: 'MCP配置JSON格式错误',
        duration: 3000
      })
      return
    }

    // 调用API保存MCP配置
    await updateMCPConfig({ config: mcpConfigContent.value })

    notification.success({
      content: 'MCP配置保存成功（请手动重启服务以应用更改）',
      duration: 5000
    })

  } catch (error) {
    console.error('保存MCP配置失败:', error)
    notification.error({
      content: '保存MCP配置失败',
      duration: 3000
    })
  } finally {
    savingMCPConfig.value = false
  }
}

// 配置模板定义
const configTemplates = {
  basic: `{
  "mcpServers": {

  }
}`,
  'web-search': '{\n' +
  '  "mcpServers": {\n' +
  '    "web-search": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-web-search"],\n' +
  '      "env": {\n' +
  '        "BING_API_KEY": "${BING_API_KEY}"\n' +
  '      },\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": true\n' +
  '    }\n' +
  '  }\n' +
  '}',
  filesystem: `{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],
      "transport": "stdio",
      "enabled": true,
      "auto_start": false
    }
  }
}`,
  full: '{\n' +
  '  "mcpServers": {\n' +
  '    "web-search": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-web-search"],\n' +
  '      "env": {\n' +
  '        "BING_API_KEY": "${BING_API_KEY}"\n' +
  '      },\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": true\n' +
  '    },\n' +
  '    "filesystem": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": true,\n' +
  '      "auto_start": false\n' +
  '    },\n' +
  '    "git": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-git", "/path/to/repo"],\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": false,\n' +
  '      "auto_start": false\n' +
  '    },\n' +
  '    "postgres": {\n' +
  '      "command": "npx",\n' +
  '      "args": ["@modelcontextprotocol/server-postgres"],\n' +
  '      "env": {\n' +
  '        "POSTGRES_CONNECTION_STRING": "postgresql://user:password@localhost:5432/dbname"\n' +
  '      },\n' +
  '      "transport": "stdio",\n' +
  '      "enabled": false,\n' +
  '      "auto_start": false\n' +
  '    }\n' +
  '  }\n' +
  '}'
}

// 应用配置模板
const applyTemplate = (templateKey: string) => {
  const template = configTemplates[templateKey as keyof typeof configTemplates]
  if (template) {
    mcpConfigContent.value = template
    notification.success({
      content: `已应用${templateOptions.find(t => t.key === templateKey)?.label}`,
      duration: 3000
    })
  }
}

// 编辑器验证事件处理
const onEditorValidate = (isValid: boolean, error?: string) => {
  mcpConfigValid.value = isValid
  if (error) {
    console.log('JSON验证错误:', error)
  }
}

// 编辑器内容变化事件处理
const onEditorChange = (value: string) => {
  // 内容变化时的处理逻辑
  mcpConfigContent.value = value
}

// MCP配置内容变化时验证
const onMCPConfigChange = () => {
  try {
    JSON.parse(mcpConfigContent.value)
    mcpConfigValid.value = true
  } catch (e) {
    mcpConfigValid.value = false
  }
}

// 测试MCP配置
const testMCPConfig = () => {
  try {
    const parsed = JSON.parse(mcpConfigContent.value)
    console.log('MCP配置验证成功:', parsed)

    // 检查基本结构
    if (!parsed.mcpServers) {
      notification.warning({
        content: 'MCP配置缺少 mcpServers 字段',
        duration: 3000
      })
      return
    }

    const serverCount = Object.keys(parsed.mcpServers).length
    notification.success({
      content: `MCP配置格式正确，包含 ${serverCount} 个服务`,
      duration: 3000
    })
  } catch (e: any) {
    console.error('MCP配置验证失败:', e)
    notification.error({
      content: `MCP配置JSON格式错误: ${e.message}`,
      duration: 5000
    })
  }
}

// MCP服务控制方法
const startMCPService = async (name: string) => {
  try {
    await startMCP(name)
    notification.success({
      content: `MCP服务 ${name} 启动成功`,
      duration: 3000
    })
    await loadMCPStatus()
  } catch (error) {
    console.error(`启动MCP服务 ${name} 失败:`, error)
    notification.error({
      content: `启动MCP服务 ${name} 失败`,
      duration: 3000
    })
  }
}

const stopMCPService = async (name: string) => {
  try {
    await stopMCP(name)
    notification.success({
      content: `MCP服务 ${name} 停止成功`,
      duration: 3000
    })
    await loadMCPStatus()
  } catch (error) {
    console.error(`停止MCP服务 ${name} 失败:`, error)
    notification.error({
      content: `停止MCP服务 ${name} 失败`,
      duration: 3000
    })
  }
}

const restartMCPService = async (name: string) => {
  try {
    await restartMCP(name)
    notification.success({
      content: `MCP服务 ${name} 重启成功`,
      duration: 3000
    })
    await loadMCPStatus()
  } catch (error) {
    console.error(`重启MCP服务 ${name} 失败:`, error)
    notification.error({
      content: `重启MCP服务 ${name} 失败`,
      duration: 3000
    })
  }
}

// 页面加载时获取配置
onMounted(async () => {
  await loadAIConfig()
  await loadMCPConfig()
  await loadMCPStatus()
})

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    notification.success({
      content: '已复制到剪贴板',
      duration: 2000
    })
  } catch (error) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    notification.success({
      content: '已复制到剪贴板',
      duration: 2000
    })
  }
}

// 测试AI连接
const testAIConnectionClick = async () => {
  if (!configForm.value.api_key || !configForm.value.api_url) {
    notification.warning({
      content: '请先填写 API Key 和 API URL',
      duration: 3000
    })
    return
  }

  // 构建测试配置，使用当前表单的参数
  const testConfig = {
    api_key: configForm.value.api_key,
    ai_api_url: configForm.value.api_url,
    ai_model: configForm.value.model,
    ai_max_tokens: configForm.value.max_tokens,
    ai_temperature: configForm.value.temperature,
    ai_timeout: configForm.value.timeout,
    ai_retry_count: configForm.value.retry_count
  }

  // 设置测试数据
  testData.value.sent = {
    prompt: "你是什么AI模型？请详细介绍你的名称、版本和能力。",
    config: testConfig
  }
  testData.value.received = null
  testData.value.loading = true
  testData.value.error = null
  showTestModal.value = true

  try {
    const response = await testAIConnection(testConfig)
    testData.value.received = response

    // 如果测试成功，更新 API Key 配置状态
    configForm.value.api_key_configured = true
  } catch (error: any) {
    console.error('AI连接测试失败:', error)
    testData.value.error = error?.response?.data?.message || error?.message || 'AI连接测试失败'

    // 如果测试失败，标记 API Key 未配置
    configForm.value.api_key_configured = false
  } finally {
    testData.value.loading = false
  }
}
</script>

<style scoped>
/* 自定义样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}

/* tab内容容器 - 个别内容滚动 */
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>