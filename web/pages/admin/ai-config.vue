<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和保存按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">AI 配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理AI服务配置和设置</p>
      </div>
      <n-button type="primary" @click="saveConfig" :loading="saving">
        <template #icon>
          <i class="fas fa-save"></i>
        </template>
        保存配置
      </n-button>
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
                    <n-input
                      v-model:value="configForm.api_key"
                      type="password"
                      placeholder="请输入API Key"
                      show-password-on="click"
                    />
                    <div v-if="configForm.api_key_configured" class="text-sm text-green-600 dark:text-green-400">
                      <i class="fas fa-check-circle mr-1"></i> API Key 已配置
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

                  <!-- 组织ID -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">组织ID</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">AI服务的组织ID（可选）</span>
                    </div>
                    <n-input
                      v-model:value="configForm.organization"
                      placeholder="请输入组织ID（可选）"
                    />
                  </div>

                  <!-- 代理设置 -->
                  <div class="space-y-2">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">代理设置</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">AI服务的代理地址（可选）</span>
                    </div>
                    <n-input
                      v-model:value="configForm.proxy"
                      placeholder="请输入代理地址（可选），如：http://proxy.example.com:8080"
                    />
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
            <div class="tab-content-container">
              <div class="space-y-6">
                <!-- MCP 配置编辑器 -->
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">MCP 配置</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">MCP服务配置文件编辑器</span>
                    </div>
                    <n-button @click="loadMCPConfig" :loading="loadingMCPConfig">
                      <template #icon>
                        <i class="fas fa-sync"></i>
                      </template>
                      重新加载
                    </n-button>
                  </div>
                  <div class="border rounded-lg overflow-hidden">
                    <textarea
                      v-model="mcpConfigContent"
                      class="w-full h-96 font-mono text-sm p-4 border-none outline-none resize-none"
                      placeholder="MCP配置JSON内容..."
                    ></textarea>
                  </div>
                  <div class="flex space-x-2">
                    <n-button type="primary" @click="saveMCPConfig" :loading="savingMCPConfig">
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

                <!-- MCP 服务状态 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">MCP 服务状态</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">当前运行的MCP服务</span>
                  </div>
                  <div class="border rounded-lg p-4">
                    <div v-if="mcpServices.length === 0" class="text-gray-500 text-center py-8">
                      暂无MCP服务
                    </div>
                    <div v-else class="space-y-3">
                      <div v-for="service in mcpServices" :key="service.name" class="flex items-center justify-between p-3 border rounded-lg">
                        <div>
                          <div class="font-medium">{{ service.name }}</div>
                          <div class="text-sm text-gray-600 dark:text-gray-400">{{ service.status }}</div>
                        </div>
                        <div class="flex space-x-2">
                          <n-button size="small" type="primary" @click="startMCPService(service.name)">
                            启动
                          </n-button>
                          <n-button size="small" @click="stopMCPService(service.name)">
                            停止
                          </n-button>
                          <n-button size="small" @click="restartMCPService(service.name)">
                            重启
                          </n-button>
                        </div>
                      </div>
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
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useToast } from '~/composables/useToast'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'
import { useAIApi, useMCPApi } from '~/composables/useApi'

const { notification } = useToast()
const activeTab = ref('openai')
const saving = ref(false)
const loadingMCPConfig = ref(false)
const savingMCPConfig = ref(false)
const formRef = ref()

// AI配置表单数据
interface AIConfigForm {
  api_key?: string
  api_key_configured: boolean
  api_url: string
  model: string
  max_tokens: number
  temperature: number
  organization?: string
  proxy?: string
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
        organization: response.ai_organization || undefined,
        proxy: response.ai_proxy || undefined,
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
    const response = await getMCPConfig()
    if (response && response.config) {
      mcpConfigContent.value = response.config
    } else {
      // 如果API没有返回配置，则使用默认配置
      mcpConfigContent.value = `{
  "mcpServers": {
    "web-search": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-web-search"],
      "env": {
        "BING_API_KEY": "\${BING_API_KEY}"
      },
      "transport": "stdio",
      "enabled": true,
      "auto_start": true
    },
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],
      "transport": "stdio",
      "enabled": true,
      "auto_start": false
    }
  }
}`
    }
  } catch (error) {
    console.error('获取MCP配置失败:', error)
    notification.error({
      content: '获取MCP配置失败',
      duration: 3000
    })
    // 如果API失败，也提供一个默认配置
    mcpConfigContent.value = `{
  "mcpServers": {
    "web-search": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-web-search"],
      "env": {
        "BING_API_KEY": "\${BING_API_KEY}"
      },
      "transport": "stdio",
      "enabled": true,
      "auto_start": true
    },
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],
      "transport": "stdio",
      "enabled": true,
      "auto_start": false
    }
  }
}`
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

    if (configForm.value.organization) {
      updateData.ai_organization = configForm.value.organization
    }

    if (configForm.value.proxy) {
      updateData.ai_proxy = configForm.value.proxy
    }

    // 调用API更新配置
    const result = await updateAIConfig(updateData)

    notification.success({
      content: 'AI配置保存成功',
      duration: 3000
    })

    // 更新原始配置数据
    updateOriginalConfig(configForm.value)
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

// 测试MCP配置
const testMCPConfig = () => {
  try {
    JSON.parse(mcpConfigContent.value)
    notification.success({
      content: 'MCP配置格式正确',
      duration: 3000
    })
  } catch (e) {
    notification.error({
      content: 'MCP配置JSON格式错误',
      duration: 3000
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

// 测试AI连接
const testAIConnectionClick = async () => {
  try {
    await testAIConnection()
    notification.success({
      content: 'AI连接测试成功',
      duration: 3000
    })
  } catch (error) {
    console.error('AI连接测试失败:', error)
    notification.error({
      content: 'AI连接测试失败',
      duration: 3000
    })
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