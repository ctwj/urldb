<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">机器人管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理各平台的机器人配置和自动回复功能</p>
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
          class="mb-6"
        >
          <n-tab-pane name="qq" tab="QQ机器人">
            <div class="tab-content-container">
              <div class="space-y-8">
            <!-- 步骤1：Astrobot 安装指南 -->
            <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-6">
              <div class="flex items-center mb-4">
                <div class="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center mr-3">
                  <span class="text-sm font-bold">1</span>
                </div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">安装 Astrobot</h3>
              </div>
              <div class="space-y-4">
                <div class="flex items-start space-x-3">
                  <i class="fas fa-github text-gray-600 dark:text-gray-400 mt-1"></i>
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white mb-1">开源地址</p>
                    <a 
                      href="https://github.com/Astrian/astrobot" 
                      target="_blank" 
                      class="text-blue-600 dark:text-blue-400 hover:underline text-sm"
                    >
                      https://github.com/Astrian/astrobot
                    </a>
                  </div>
                </div>
                <div class="flex items-start space-x-3">
                  <i class="fas fa-book text-gray-600 dark:text-gray-400 mt-1"></i>
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white mb-1">安装教程</p>
                    <a 
                      href="https://github.com/Astrian/astrobot/wiki" 
                      target="_blank" 
                      class="text-blue-600 dark:text-blue-400 hover:underline text-sm"
                    >
                      https://github.com/Astrian/astrobot/wiki
                    </a>
                  </div>
                </div>
              </div>
            </div>

            <!-- 步骤2：插件安装 -->
            <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-6">
              <div class="flex items-center mb-4">
                <div class="w-8 h-8 bg-green-600 text-white rounded-full flex items-center justify-center mr-3">
                  <span class="text-sm font-bold">2</span>
                </div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">安装插件</h3>
              </div>
              <div class="space-y-4">
                <div class="flex items-start space-x-3">
                  <i class="fas fa-puzzle-piece text-gray-600 dark:text-gray-400 mt-1"></i>
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white mb-1">插件地址</p>
                    <a 
                      href="https://github.com/ctwj/astrbot_plugin_urldb" 
                      target="_blank" 
                      class="text-green-600 dark:text-green-400 hover:underline text-sm"
                    >
                      https://github.com/ctwj/astrbot_plugin_urldb
                    </a>
                  </div>
                </div>
                <div class="bg-gray-100 dark:bg-gray-800 rounded p-3">
                  <p class="text-sm text-gray-700 dark:text-gray-300">
                    <strong>插件特性：</strong><br>
                    • 支持@机器人搜索功能<br>
                    • 可配置API域名和密钥<br>
                    • 自动格式化搜索结果<br>
                    • 支持超时时间配置<br><br>
                    <strong>安装步骤：</strong><br>
                    1. Astrbot 插件市场， 搜索 urldb 安装<br>
                    2. 根据下面的配置信息配置插件
                  </p>
                </div>
              </div>
            </div>

            <!-- 步骤3：配置信息 -->
            <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-6">
              <div class="flex items-center mb-4">
                <div class="w-8 h-8 bg-purple-600 text-white rounded-full flex items-center justify-center mr-3">
                  <span class="text-sm font-bold">3</span>
                </div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">配置信息</h3>
              </div>
              <div class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">网站域名</label>
                    <div class="flex items-center space-x-2">
                      <n-input 
                        :value="siteDomain" 
                        readonly 
                        class="flex-1"
                      />
                      <n-button 
                        size="small" 
                        @click="copyToClipboard(siteDomain)"
                        type="primary"
                      >
                        <template #icon>
                          <i class="fas fa-copy"></i>
                        </template>
                        复制
                      </n-button>
                    </div>
                  </div>
                  <div>
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">API Token</label>
                    <div class="flex items-center space-x-2">
                      <n-input 
                        :value="apiToken" 
                        readonly 
                        type="password"
                        show-password-on="click"
                        class="flex-1"
                      />
                      <n-button 
                        size="small" 
                        @click="copyToClipboard(apiToken)"
                        type="primary"
                      >
                        <template #icon>
                          <i class="fas fa-copy"></i>
                        </template>
                        复制
                      </n-button>
                    </div>
                  </div>
                </div>
                <div class="bg-gray-100 dark:bg-gray-800 rounded p-3">
                  <p class="text-sm text-gray-700 dark:text-gray-300">
                    <strong>配置说明：</strong><br>
                    将上述信息配置到 Astrobot 的插件配置文件中，插件将自动连接到本系统进行资源搜索。
                  </p>
                </div>
              </div>
            </div>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="wechat" tab="微信公众号">
          <div class="tab-content-container">
            <div class="text-center py-12">
              <i class="fas fa-lock text-4xl text-gray-400 mb-4"></i>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能暂未开放</h3>
              <p class="text-gray-500 dark:text-gray-400">微信公众号机器人功能正在开发中，敬请期待</p>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="telegram" tab="Telegram机器人">
          <div class="tab-content-container">
            <div class="text-center py-12">
              <i class="fas fa-lock text-4xl text-gray-400 mb-4"></i>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能暂未开放</h3>
              <p class="text-gray-500 dark:text-gray-400">Telegram机器人功能正在开发中，敬请期待</p>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="wechat_open" tab="微信开放平台">
          <div class="tab-content-container">
            <div class="text-center py-12">
              <i class="fas fa-lock text-4xl text-gray-400 mb-4"></i>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能暂未开放</h3>
              <p class="text-gray-500 dark:text-gray-400">微信开放平台机器人功能正在开发中，敬请期待</p>
            </div>
          </div>
        </n-tab-pane>
        </n-tabs>
      </div>
    </template>
  </AdminPageLayout>
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'

// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

// 定义配置表单类型
interface BotConfigForm {
  api_token: string
}

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  updateOriginalConfig,
  saveConfig: saveConfigWithDetection
} = useConfigChangeDetection<BotConfigForm>({
  debug: true,
  fieldMapping: {
    api_token: 'api_token'
  }
})

const notification = useNotification()
const activeTab = ref('qq')

// 获取网站域名和API Token
const siteDomain = computed(() => {
  if (process.client) {
    return window.location.origin
  }
  return 'https://yourdomain.com'
})

const apiToken = ref('')

// 获取API Token
const fetchApiToken = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig()
    
    if (response) {
      const configData = {
        api_token: (response as any).api_token || ''
      }
      
      apiToken.value = configData.api_token || '未配置API Token'
      setOriginalConfig(configData)
    } else {
      apiToken.value = '未配置API Token'
    }
  } catch (error) {
    console.error('获取API Token失败:', error)
    apiToken.value = '获取失败'
  }
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

// 页面加载时获取配置
onMounted(() => {
  fetchApiToken()
  console.log('机器人管理页面已加载')
})
</script>

<style scoped>
/* 机器人管理页面样式 */

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