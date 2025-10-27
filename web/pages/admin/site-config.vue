<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和保存按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">站点配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理网站基本信息和设置</p>
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
          <n-tab-pane name="basic" tab="基本信息">
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
                <!-- 网站标题 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站标题</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">网站的主要标识，显示在浏览器标签页和搜索结果中</span>
                  </div>
                  <n-input
                    v-model:value="configForm.site_title"
                    placeholder="请输入网站标题"
                  />
                </div>

                <!-- 网站描述 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站描述</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">网站的简要介绍，用于SEO和社交媒体分享</span>
                  </div>
                  <n-input
                    v-model:value="configForm.site_description"
                    placeholder="请输入网站描述"
                  />
                </div>

                <!-- 关键词 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">关键词</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">用于SEO优化，多个关键词用逗号分隔</span>
                  </div>
                  <n-input
                    v-model:value="configForm.keywords"
                    placeholder="请输入关键词，用逗号分隔"
                  />
                </div>

                <!-- 网站Logo -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站Logo</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">选择网站Logo图片，建议使用正方形图片</span>
                  </div>
                  <div class="flex items-center space-x-4">
                    <div v-if="configForm.site_logo" class="flex-shrink-0">
                      <n-image
                        :src="getImageUrl(configForm.site_logo)"
                        alt="网站Logo"
                        width="80"
                        height="80"
                        object-fit="cover"
                        class="rounded-lg border"
                      />
                    </div>
                    <div class="flex-1">
                      <n-button type="primary" @click="openLogoSelector">
                        <template #icon>
                          <i class="fas fa-image"></i>
                        </template>
                        {{ configForm.site_logo ? '更换Logo' : '选择Logo' }}
                      </n-button>
                      <n-button v-if="configForm.site_logo" @click="clearLogo" class="ml-2">
                        <template #icon>
                          <i class="fas fa-times"></i>
                        </template>
                        清除
                      </n-button>
                    </div>
                  </div>
                </div>

                <!-- 版权信息 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">版权信息</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">网站底部的版权声明信息</span>
                  </div>
                  <n-input
                    v-model:value="configForm.copyright"
                    placeholder="请输入版权信息"
                  />
                </div>

                <!-- 二维码样式 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">二维码样式</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">选择前台显示的二维码样式</span>
                  </div>
                  <div class="flex items-center space-x-4">
                    <div class="flex items-center space-x-3">
                      <n-button type="primary" @click="openQRStyleSelector">
                        <template #icon>
                          <i class="fas fa-qrcode"></i>
                        </template>
                        {{ configForm.qr_code_style ? '更换样式' : '选择样式' }}
                      </n-button>
                      <div v-if="configForm.qr_code_style" class="text-sm text-gray-600 dark:text-gray-400">
                        当前样式: <span class="font-semibold">{{ configForm.qr_code_style }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              </n-form>
            </div>
          </n-tab-pane>

          <n-tab-pane name="security" tab="安全设置">

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
                <!-- 维护模式 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">维护模式</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">开启后网站将显示维护页面，暂停用户访问</span>
                  </div>
                  <n-switch v-model:value="configForm.maintenance_mode" />
                </div>

                <!-- 违禁词 -->
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-2">
                      <label class="text-base font-semibold text-gray-800 dark:text-gray-200">违禁词</label>
                      <span class="text-xs text-gray-500 dark:text-gray-400">包含这些词汇的资源将被过滤，多个词汇用逗号分隔</span>
                    </div>
                    <a
                      href="https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/db/forbidden.txt"
                      target="_blank"
                      class="text-xs text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 underline"
                    >
                      开源违禁词
                    </a>
                  </div>
                  <n-input
                    v-model:value="configForm.forbidden_words"
                    placeholder="请输入违禁词，用逗号分隔"
                    type="textarea"
                    :rows="4"
                  />
                </div>

                <!-- 开启注册 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">开启注册</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">开启后用户才能注册新账号，关闭后注册页面将显示"当前系统已关闭注册功能"</span>
                  </div>
                  <n-switch v-model:value="configForm.enable_register" />
                </div>
              </div>
              </n-form>
            </div>
          </n-tab-pane>

          <n-tab-pane name="ui" tab="界面配置">
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
                 <!-- 公告配置组件 -->
                 <AnnouncementConfig
                   v-model="announcementConfig"
                   @update:modelValue="handleAnnouncementUpdate"
                 />

                 <!-- 浮动按钮配置组件 -->
                 <FloatButtonsConfig
                   v-model="floatButtonsConfig"
                   @update:modelValue="handleFloatButtonsUpdate"
                   @openWechatSelector="showWechatSelector = true"
                   @openTelegramSelector="showTelegramSelector = true"
                 />

                 <!-- 微信图片选择器 -->
                 <ImageSelectorModal
                   v-model:show="showWechatSelector"
                   title="选择微信搜一搜图片"
                   @select="handleWechatImageSelect"
                 />

                 <!-- Telegram图片选择器 -->
                 <ImageSelectorModal
                   v-model:show="showTelegramSelector"
                   title="选择Telegram二维码图片"
                   @select="handleTelegramImageSelect"
                 />
                </div>
              </n-form>
            </div>
          </n-tab-pane>
        </n-tabs>
      </div>
    </template>
</AdminPageLayout>
    <!-- ImageSelectorModal 组件 -->
    <ImageSelectorModal
      v-model:show="showLogoSelector"
      title="选择Logo图片"
      @select="handleLogoSelect"
    />

    <!-- QR Code Style Selector Modal -->
    <QRCodeStyleSelector
      v-model:show="showQRStyleSelector"
      :current-style="configForm.qr_code_style"
      @select="handleQRStyleSelect"
    />
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})


import { useImageUrl } from '~/composables/useImageUrl'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'
import AnnouncementConfig from '~/components/Admin/AnnouncementConfig.vue'
import FloatButtonsConfig from '~/components/Admin/FloatButtonsConfig.vue'
import ImageSelectorModal from '~/components/Admin/ImageSelectorModal.vue'
import QRCodeStyleSelector from '~/components/Admin/QRCodeStyleSelector.vue'

const notification = useNotification()
const { getImageUrl } = useImageUrl()
const formRef = ref()
const saving = ref(false)
const activeTab = ref('basic')

// Logo选择器相关数据
const showLogoSelector = ref(false)

// 微信和Telegram选择器相关数据
const showWechatSelector = ref(false)
const showTelegramSelector = ref(false)

// QR样式选择器相关数据
const showQRStyleSelector = ref(false)

// 公告类型接口
interface Announcement {
  content: string
  enabled: boolean
}

// 配置表单数据类型
interface SiteConfigForm {
  site_title: string
  site_description: string
  keywords: string
  copyright: string
  site_logo: string
  maintenance_mode: boolean
  enable_register: boolean
  forbidden_words: string
  enable_announcements: boolean
  announcements: Announcement[]
  enable_float_buttons: boolean
  wechat_search_image: string
  telegram_qr_image: string
  qr_code_style: string
}

// 公告配置子组件数据
const announcementConfig = computed({
  get: () => ({
    enable_announcements: configForm.value.enable_announcements,
    announcements: configForm.value.announcements
  }),
  set: (value: any) => {
    configForm.value.enable_announcements = value.enable_announcements
    configForm.value.announcements = value.announcements
  }
})

// 浮动按钮配置子组件数据
const floatButtonsConfig = computed({
  get: () => ({
    enable_float_buttons: configForm.value.enable_float_buttons,
    wechat_search_image: configForm.value.wechat_search_image,
    telegram_qr_image: configForm.value.telegram_qr_image
  }),
  set: (value: any) => {
    configForm.value.enable_float_buttons = value.enable_float_buttons
    configForm.value.wechat_search_image = value.wechat_search_image
    configForm.value.telegram_qr_image = value.telegram_qr_image
  }
})

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  getChangedDetails,
  updateOriginalConfig,
  saveConfig: saveConfigWithDetection
} = useConfigChangeDetection<SiteConfigForm>({
  debug: true,
  // 自定义比较函数，处理数组深层比较
  customCompare: (key: string, currentValue: any, originalValue: any) => {
    // 对于数组类型，使用JSON字符串比较
    if (Array.isArray(currentValue) && Array.isArray(originalValue)) {
      return JSON.stringify(currentValue) !== JSON.stringify(originalValue)
    }
    // 其他类型使用默认比较
    return currentValue !== originalValue
  },
  // 字段映射：前端字段名 -> 后端字段名
  fieldMapping: {
    site_title: 'site_title',
    site_description: 'site_description',
    keywords: 'keywords',
    copyright: 'copyright',
    site_logo: 'site_logo',
    maintenance_mode: 'maintenance_mode',
    enable_register: 'enable_register',
    forbidden_words: 'forbidden_words',
    enable_announcements: 'enable_announcements',
    announcements: 'announcements',
    enable_float_buttons: 'enable_float_buttons',
    wechat_search_image: 'wechat_search_image',
    telegram_qr_image: 'telegram_qr_image',
    qr_code_style: 'qr_code_style'
  }
})

// 公告类型选项（如果需要的话可以保留，但根据反馈暂时移除）

// 配置表单数据
const configForm = ref<SiteConfigForm>({
  site_title: '',
  site_description: '',
  keywords: '',
  copyright: '',
  site_logo: '',
  maintenance_mode: false,
  enable_register: false,
  forbidden_words: '',
  enable_announcements: false,
  announcements: [],
  enable_float_buttons: false,
  wechat_search_image: '',
  telegram_qr_image: '',
  qr_code_style: 'Plain'
})



// 表单验证规则
const rules = {
  site_title: {
    required: true,
    message: '请输入网站标题',
    trigger: 'blur'
  },
  site_description: {
    required: true,
    message: '请输入网站描述',
    trigger: 'blur'
  }
}

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig() as any
    
    if (response) {
      const configData = {
        site_title: response.site_title || '',
        site_description: response.site_description || '',
        keywords: response.keywords || '',
        copyright: response.copyright || '',
        site_logo: response.site_logo || '',
        maintenance_mode: response.maintenance_mode || false,
        enable_register: response.enable_register || false,
        forbidden_words: response.forbidden_words || '',
        enable_announcements: response.enable_announcements || false,
        announcements: response.announcements ? JSON.parse(response.announcements) : [],
        enable_float_buttons: response.enable_float_buttons || false,
        wechat_search_image: response.wechat_search_image || '',
        telegram_qr_image: response.telegram_qr_image || '',
        qr_code_style: response.qr_code_style || 'Plain'
      }

      // 设置表单数据和原始数据
      configForm.value = { ...configData }
      setOriginalConfig(configData)
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
    await formRef.value?.validate()
    
    saving.value = true
    
    // 更新当前配置数据
    updateCurrentConfig(configForm.value)
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    // 使用通用保存函数
    const result = await saveConfigWithDetection(
      systemConfigApi.updateSystemConfig,
      {
        onlyChanged: true,
        includeAllFields: true
      },
      // 成功回调
      async () => {
        notification.success({
          content: '站点配置保存成功',
          duration: 3000
        })
        
        // 刷新系统配置状态，确保顶部导航同步更新
        const { useSystemConfigStore } = await import('~/stores/systemConfig')
        const systemConfigStore = useSystemConfigStore()
        await systemConfigStore.initConfig(true, true)
      },
      // 错误回调
      (error) => {
        console.error('保存站点配置失败:', error)
        notification.error({
          content: '保存站点配置失败',
          duration: 3000
        })
      }
    )
    
    // 如果没有改动，显示提示
    if (result && result.message === '没有检测到任何改动') {
      notification.info({
        content: '没有检测到任何改动',
        duration: 3000
      })
    }
  } finally {
    saving.value = false
  }
}

// Logo选择器方法
const openLogoSelector = () => {
  showLogoSelector.value = true
}

const clearLogo = () => {
  configForm.value = {
    ...configForm.value,
    site_logo: ''
  }
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// 子组件更新处理方法
const handleAnnouncementUpdate = (newValue: any) => {
  // 直接更新整个表单
  configForm.value = {
    ...configForm.value,
    enable_announcements: newValue.enable_announcements,
    announcements: newValue.announcements
  }
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

const handleFloatButtonsUpdate = (newValue: any) => {
  configForm.value = {
    ...configForm.value,
    enable_float_buttons: newValue.enable_float_buttons,
    wechat_search_image: newValue.wechat_search_image,
    telegram_qr_image: newValue.telegram_qr_image
  }
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// Logo选择处理
const handleLogoSelect = (file: any) => {
  configForm.value = {
    ...configForm.value,
    site_logo: file.access_url
  }
  showLogoSelector.value = false
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// 微信图片选择处理
const handleWechatImageSelect = (file: any) => {
  configForm.value = {
    ...configForm.value,
    wechat_search_image: file.access_url
  }
  showWechatSelector.value = false
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// Telegram图片选择处理
const handleTelegramImageSelect = (file: any) => {
  configForm.value = {
    ...configForm.value,
    telegram_qr_image: file.access_url
  }
  showTelegramSelector.value = false
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// QR样式选择器方法
const openQRStyleSelector = () => {
  showQRStyleSelector.value = true
}

// QR样式选择处理
const handleQRStyleSelect = (preset: any) => {
  configForm.value = {
    ...configForm.value,
    qr_code_style: preset.name
  }
  showQRStyleSelector.value = false
  // 强制触发更新
  updateCurrentConfig({ ...configForm.value })
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
})


</script>

<style scoped>
/* 站点配置页面样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}

/* 配置标签容器 - 支持滚动 */
.config-tabs-container {
  height: calc(100vh - 200px);
  overflow-y: auto;
  padding: 0.5rem 0;
}

/* tab内容容器 - 个别内容滚动 */
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}

</style>