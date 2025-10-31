<template>
  <div class="tab-content-container">
    <div class="space-y-6">
      <!-- 基础配置 -->
      <n-card title="基础配置" class="mb-6">
        <n-form :model="configForm" label-placement="left" label-width="120px">
          <n-form-item label="AppID">
            <n-input v-model:value="configForm.app_id" placeholder="请输入微信公众号AppID" />
          </n-form-item>
          <n-form-item label="AppSecret">
            <n-input v-model:value="configForm.app_secret" type="password" placeholder="请输入AppSecret" show-password-on="click" />
          </n-form-item>
          <n-form-item label="Token">
            <n-input v-model:value="configForm.token" placeholder="请输入Token（用于消息验证）" />
          </n-form-item>
          <n-form-item label="EncodingAESKey">
            <n-input v-model:value="configForm.encoding_aes_key" type="password" placeholder="请输入EncodingAESKey（可选，用于消息加密）" show-password-on="click" />
          </n-form-item>
        </n-form>
      </n-card>

      <!-- 功能配置 -->
      <n-card title="功能配置" class="mb-6">
        <n-form :model="configForm" label-placement="left" label-width="120px">
          <n-form-item label="启用机器人">
            <n-switch v-model:value="configForm.enabled" />
          </n-form-item>
          <n-form-item label="自动回复">
            <n-switch v-model:value="configForm.auto_reply_enabled" />
          </n-form-item>
          <n-form-item label="欢迎消息">
            <n-input v-model:value="configForm.welcome_message" type="textarea" :rows="3" placeholder="新用户关注时发送的欢迎消息" />
          </n-form-item>
          <n-form-item label="搜索结果限制">
            <n-input-number v-model:value="configForm.search_limit" :min="1" :max="10" placeholder="搜索结果返回数量" />
          </n-form-item>
        </n-form>
      </n-card>

      <!-- 微信公众号验证文件上传 -->
      <n-card title="微信公众号验证文件" class="mb-6">
        <div class="space-y-4">
          <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <h4 class="font-medium text-blue-800 dark:text-blue-200 mb-2">验证文件上传说明</h4>
            <p class="text-sm text-gray-700 dark:text-gray-300 mb-2">
              微信公众号需要上传一个TXT验证文件到网站根目录。请按照以下步骤操作：
            </p>
            <ol class="text-sm text-gray-700 dark:text-gray-300 list-decimal list-inside space-y-1">
              <li>点击下方"选择文件"按钮，选择微信提供的TXT验证文件</li>
              <li>点击"上传验证文件"按钮上传文件</li>
              <li>上传成功后，文件将可通过网站根目录直接访问</li>
              <li>在微信公众平台完成域名验证</li>
            </ol>
          </div>

          <div class="flex items-center space-x-4">
            <n-upload
              ref="uploadRef"
              :show-file-list="false"
              :accept="'.txt'"
              :max="1"
              :custom-request="handleUpload"
              @before-upload="beforeUpload"
              @change="handleFileChange"
            >
              <n-button type="primary">
                <template #icon>
                  <i class="fas fa-file-upload"></i>
                </template>
                选择TXT文件
              </n-button>
            </n-upload>
            <n-button
              type="success"
              @click="triggerUpload"
              :disabled="!selectedFile"
              :loading="uploading"
            >
              <template #icon>
                <i class="fas fa-cloud-upload-alt"></i>
              </template>
              上传验证文件
            </n-button>
          </div>

          <div v-if="uploadResult" class="p-3 rounded-md bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300">
            <div class="flex items-center">
              <i class="fas fa-check-circle mr-2"></i>
              <span>文件上传成功！</span>
            </div>
            <p class="text-xs mt-1">文件名: {{ uploadResult.file_name }}</p>
            <p class="text-xs">访问地址: {{ getFullUrl(uploadResult.access_url) }}</p>
          </div>
        </div>
      </n-card>

      <!-- 微信公众号平台配置说明 -->
      <n-card title="微信公众号平台配置" class="mb-6">
        <div class="space-y-4">
          <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <h4 class="font-medium text-blue-800 dark:text-blue-200 mb-2">服务器配置</h4>
            <p class="text-sm text-gray-700 dark:text-gray-300 mb-2">
              在微信公众平台后台的<strong>开发 > 基本配置 > 服务器配置</strong>中设置：
            </p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 block">URL</label>
                <div class="flex items-center space-x-2">
                  <n-input :value="serverUrl" readonly class="flex-1" />
                  <n-button size="small" @click="copyToClipboard(serverUrl)" type="primary">
                    <template #icon>
                      <i class="fas fa-copy"></i>
                    </template>
                    复制
                  </n-button>
                </div>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 block">Token</label>
                <div class="flex items-center space-x-2">
                  <n-input :value="configForm.token" readonly class="flex-1" />
                  <n-button size="small" @click="copyToClipboard(configForm.token)" type="primary">
                    <template #icon>
                      <i class="fas fa-copy"></i>
                    </template>
                    复制
                  </n-button>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
            <h4 class="font-medium text-green-800 dark:text-green-200 mb-2">消息加解密配置</h4>
            <p class="text-sm text-gray-700 dark:text-gray-300">
              如果需要启用消息加密，请在微信公众平台选择<strong>安全模式</strong>，并填写上面的EncodingAESKey。
            </p>
          </div>

          <div class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
            <h4 class="font-medium text-yellow-800 dark:text-yellow-200 mb-2">注意事项</h4>
            <ul class="text-sm text-gray-700 dark:text-gray-300 list-disc list-inside space-y-1">
              <li>服务器必须支持HTTPS（微信要求）</li>
              <li>域名必须已备案</li>
              <li>首次配置时，微信会发送GET请求验证服务器</li>
              <li>配置完成后记得点击"启用"按钮</li>
            </ul>
          </div>
        </div>
      </n-card>

      <!-- 操作按钮 -->
      <div class="flex justify-end space-x-4">
        <n-button @click="resetForm">重置</n-button>
        <n-button type="primary" @click="saveConfig" :loading="loading">保存配置</n-button>
      </div>

      <!-- 状态信息 -->
      <n-card title="运行状态" class="mt-6">
        <div class="space-y-4">
          <div class="flex items-center space-x-4">
            <n-tag :type="botStatus.overall_status ? 'success' : 'default'">
              {{ botStatus.status_text || '未知状态' }}
            </n-tag>
            <n-tag v-if="botStatus.config" :type="botStatus.config.enabled ? 'success' : 'default'">
              配置状态: {{ botStatus.config.enabled ? '已启用' : '已禁用' }}
            </n-tag>
            <n-tag v-if="botStatus.config" :type="botStatus.config.app_id_configured ? 'success' : 'warning'">
              AppID: {{ botStatus.config.app_id_configured ? '已配置' : '未配置' }}
            </n-tag>
          </div>

          <div v-if="!botStatus.overall_status && botStatus.config && botStatus.config.enabled" class="bg-orange-50 dark:bg-orange-900/20 rounded-lg p-3">
            <p class="text-sm text-orange-800 dark:text-orange-200">
              <i class="fas fa-exclamation-triangle mr-2"></i>
              机器人已启用但未运行，请检查配置是否正确或查看系统日志。
            </p>
          </div>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useNotification } from 'naive-ui'
import { useWechatApi } from '~/composables/useApi'

// 定义配置表单类型
interface WechatBotConfigForm {
  enabled: boolean
  app_id: string
  app_secret: string
  token: string
  encoding_aes_key: string
  welcome_message: string
  auto_reply_enabled: boolean
  search_limit: number
}

const notification = useNotification()
const loading = ref(false)
const wechatApi = useWechatApi()
const botStatus = ref({
  overall_status: false,
  status_text: '',
  config: null as any,
  runtime: null as any
})
// 验证文件上传相关
const uploadRef = ref()
const selectedFile = ref<File | null>(null)
const uploading = ref(false)
const uploadResult = ref<any>(null)

// 配置表单 - 直接使用 reactive
const configForm = reactive<WechatBotConfigForm>({
  enabled: false,
  app_id: '',
  app_secret: '',
  token: '',
  encoding_aes_key: '',
  welcome_message: '欢迎关注老九网盘资源库！发送关键词即可搜索资源。',
  auto_reply_enabled: true,
  search_limit: 5
})

// 计算服务器URL
const serverUrl = computed(() => {
  if (process.client) {
    return `${window.location.origin}/api/wechat/callback`
  }
  return 'https://yourdomain.com/api/wechat/callback'
})

// 获取机器人配置
const fetchBotConfig = async () => {
  try {
    const response = await wechatApi.getBotConfig()

    if (response) {
      // 直接更新 configForm
      configForm.enabled = response.enabled || false
      configForm.app_id = response.app_id || ''
      configForm.app_secret = response.app_secret || '' // 现在所有字段都不敏感
      configForm.token = response.token || ''
      configForm.encoding_aes_key = response.encoding_aes_key || ''
      configForm.welcome_message = response.welcome_message || '欢迎关注老九网盘资源库！发送关键词即可搜索资源。'
      configForm.auto_reply_enabled = response.auto_reply_enabled || true
      configForm.search_limit = response.search_limit || 5
    }
  } catch (error) {
    console.error('获取微信机器人配置失败:', error)
    notification.error({
      content: '获取配置失败',
      duration: 3000
    })
  }
}

// 获取机器人状态
const fetchBotStatus = async () => {
  try {
    const response = await wechatApi.getBotStatus()

    if (response) {
      botStatus.value = response
    }
  } catch (error) {
    console.error('获取微信机器人状态失败:', error)
    notification.error({
      content: '获取状态失败',
      duration: 3000
    })
  }
}

// 保存配置
const saveConfig = async () => {
  loading.value = true
  try {
    // 直接保存所有字段，不检测变更
    const payload = {
      enabled: configForm.enabled,
      app_id: configForm.app_id,
      app_secret: configForm.app_secret,
      token: configForm.token,
      encoding_aes_key: configForm.encoding_aes_key,
      welcome_message: configForm.welcome_message,
      auto_reply_enabled: configForm.auto_reply_enabled,
      search_limit: configForm.search_limit
    }

    const response = await wechatApi.updateBotConfig(payload)

    if (response.success) {
      notification.success({
        content: '配置保存成功',
        duration: 3000
      })
      // 重新获取状态和配置
      await fetchBotConfig()
      await fetchBotStatus()
    } else {
      throw new Error(response.message || '保存失败')
    }
  } catch (error: any) {
    console.error('保存微信机器人配置失败:', error)
    notification.error({
      content: error.message || '保存配置失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  // 重新获取原始配置
  fetchBotConfig()
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    notification.success({
      content: '已复制到剪贴板',
      duration: 2000
    })
  } catch (error) {
    console.error('复制失败:', error)
    notification.error({
      content: '复制失败',
      duration: 2000
    })
  }
}

// 验证文件上传相关函数
const beforeUpload = (options: { file: any, fileList: any[] }) => {
  // 从 options 中提取文件
  const file = options?.file?.file || options?.file

  // 检查文件对象是否有效
  if (!file || !file.name) {
    notification.error({
      content: '文件选择失败，请重试',
      duration: 2000
    })
    return false
  }

  // 验证文件类型 - 使用多重检查确保是TXT文件
  const isValid = file.type === 'text/plain' ||
                  file.name.toLowerCase().endsWith('.txt') ||
                  file.type === 'text/plain;charset=utf-8'

  if (!isValid) {
    notification.error({
      content: '请上传TXT文件',
      duration: 2000
    })
    selectedFile.value = null // 清空之前的选择
    return false // 阻止上传无效文件
  }

  // 保存选中的文件并更新状态
  selectedFile.value = file

  return false // 阻止自动上传
}

const handleUpload = ({ file, onSuccess, onError }: any) => {
  // 这个函数不会被调用，因为我们阻止了自动上传
}

// 文件选择变化时的处理函数
const handleFileChange = (options: { file: any, fileList: any[] }) => {
  // 从 change 事件中提取文件信息
  const file = options?.file?.file || options?.file

  if (file && file.name) {
    // 更新选中的文件
    selectedFile.value = file
  }
}

const triggerUpload = async () => {
  if (!selectedFile.value) {
    notification.warning({
      content: '请选择要上传的TXT文件',
      duration: 2000
    })
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)

    const response = await wechatApi.uploadVerifyFile(formData)

    if (response.success) {
      uploadResult.value = response
      notification.success({
        content: '验证文件上传成功',
        duration: 3000
      })
      // 清空选择的文件
      selectedFile.value = null
      if (uploadRef.value) {
        uploadRef.value.clear()
      }
    } else {
      throw new Error(response.message || '上传失败')
    }
  } catch (error: any) {
    console.error('上传验证文件失败:', error)
    notification.error({
      content: error.message || '上传验证文件失败',
      duration: 3000
    })
  } finally {
    uploading.value = false
  }
}

const getFullUrl = (path: string) => {
  if (process.client) {
    return `${window.location.origin}${path}`
  }
  return path
}

// 页面加载时获取配置和状态
onMounted(async () => {
  await fetchBotConfig()
  await fetchBotStatus()
  // 定期刷新状态
  const interval = setInterval(fetchBotStatus, 30000)
  onUnmounted(() => clearInterval(interval))
})
</script>

<style scoped>
/* 微信公众号机器人标签样式 */
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>