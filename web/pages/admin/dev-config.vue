<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">开发配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理API和开发相关配置</p>
      </div>
      <n-button type="primary" @click="saveConfig" :loading="saving">
        <template #icon>
          <i class="fas fa-save"></i>
        </template>
        保存配置
      </n-button>
    </div>

    <!-- 配置表单 -->
    <n-card>
      <div class="space-y-6">
        <!-- API Token -->
        <div>
          <n-form-item label="公开API访问令牌" path="api_token">
            <div class="flex gap-2">
              <n-input
                v-model:value="configForm.api_token"
                type="password"
                placeholder="输入API Token，用于公开API访问认证"
                show-password-on="click"
              />
              <n-button
                v-if="!configForm.api_token"
                type="primary"
                @click="generateApiToken"
              >
                生成
              </n-button>
              <template v-else>
                <n-button
                  type="primary"
                  @click="copyApiToken"
                >
                  复制
                </n-button>
                <n-button
                  type="warning"
                  @click="regenerateApiToken"
                >
                  重新生成
                </n-button>
              </template>
            </div>
            <template #help>
              API Token用于公开API的访问认证，请妥善保管
            </template>
          </n-form-item>
        </div>

        <!-- API文档链接 -->
        <div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <h3 class="text-lg font-medium text-blue-900 dark:text-blue-100 mb-2">
            API文档
          </h3>
          <p class="text-sm text-blue-700 dark:text-blue-300 mb-3">
            查看完整的API文档和使用说明
          </p>
          <n-button type="primary" @click="openApiDocs">
            <template #icon>
              <i class="fas fa-book"></i>
            </template>
            查看API文档
          </n-button>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

const notification = useNotification()
const saving = ref(false)

// 配置表单数据
const configForm = ref({
  api_token: ''
})

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig()
    
    if (response) {
      configForm.value = {
        api_token: response.api_token || ''
      }
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
    notification.error({
      content: '获取系统配置失败',
      duration: 3000
    })
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    saving.value = true
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    await systemConfigApi.updateSystemConfig({
      api_token: configForm.value.api_token
    })
    
    notification.success({
      content: '开发配置保存成功',
      duration: 3000
    })
  } catch (error) {
    console.error('保存开发配置失败:', error)
    notification.error({
      content: '保存开发配置失败',
      duration: 3000
    })
  } finally {
    saving.value = false
  }
}

// 生成API Token
const generateApiToken = async () => {
  try {
    const token = Math.random().toString(36).substring(2) + Date.now().toString(36)
    configForm.value.api_token = token
    
    notification.success({
      content: 'API Token生成成功',
      duration: 3000
    })
  } catch (error) {
    console.error('生成API Token失败:', error)
    notification.error({
      content: '生成API Token失败',
      duration: 3000
    })
  }
}

// 复制API Token
const copyApiToken = async () => {
  try {
    await navigator.clipboard.writeText(configForm.value.api_token)
    notification.success({
      content: 'API Token已复制到剪贴板',
      duration: 3000
    })
  } catch (error) {
    console.error('复制API Token失败:', error)
    notification.error({
      content: '复制API Token失败',
      duration: 3000
    })
  }
}

// 重新生成API Token
const regenerateApiToken = async () => {
  try {
    const token = Math.random().toString(36).substring(2) + Date.now().toString(36)
    configForm.value.api_token = token
    
    notification.success({
      content: 'API Token重新生成成功',
      duration: 3000
    })
  } catch (error) {
    console.error('重新生成API Token失败:', error)
    notification.error({
      content: '重新生成API Token失败',
      duration: 3000
    })
  }
}

// 打开API文档
const openApiDocs = () => {
  window.open('/api-docs', '_blank')
}

// 打开API测试工具
const openApiTest = () => {
  window.open('/api-test', '_blank')
}

// 导出配置
const exportConfig = async () => {
  try {
    const configData = {
      api_token: configForm.value.api_token,
      export_time: new Date().toISOString()
    }
    
    const blob = new Blob([JSON.stringify(configData, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `dev-config-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    notification.success({
      content: '配置导出成功',
      duration: 3000
    })
  } catch (error) {
    console.error('导出配置失败:', error)
    notification.error({
      content: '导出配置失败',
      duration: 3000
    })
  }
}

// 导入配置
const importConfig = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
      try {
        const text = await file.text()
        const configData = JSON.parse(text)
        
        if (configData.api_token) {
          configForm.value.api_token = configData.api_token
          notification.success({
            content: '配置导入成功',
            duration: 3000
          })
        } else {
          notification.error({
            content: '配置文件格式错误',
            duration: 3000
          })
        }
      } catch (error) {
        console.error('导入配置失败:', error)
        notification.error({
          content: '导入配置失败',
          duration: 3000
        })
      }
    }
  }
  input.click()
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
})


</script>

<style scoped>
/* 自定义样式 */
</style> 