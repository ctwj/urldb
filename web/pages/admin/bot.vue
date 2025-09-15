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
            <div class="space-y-8">
              <!-- 机器人基本配置 -->
              <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <div class="flex items-center mb-6">
                  <div class="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center mr-3">
                    <span class="text-sm font-bold">1</span>
                  </div>
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-white">机器人配置</h3>
                </div>

                <div class="space-y-4">
                  <!-- 机器人启用开关 -->
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700 dark:text-gray-300">启用 Telegram 机器人</label>
                      <p class="text-xs text-gray-500 dark:text-gray-400">开启后机器人将开始工作</p>
                    </div>
                    <n-switch
                      v-model:value="telegramBotConfig.bot_enabled"
                      @update:value="handleBotConfigChange"
                    />
                  </div>

                  <!-- API Key 配置 -->
                  <div v-if="telegramBotConfig.bot_enabled" class="space-y-3">
                    <div>
                      <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">Bot API Key</label>
                      <div class="flex space-x-3">
                        <n-input
                          v-model:value="telegramBotConfig.bot_api_key"
                          placeholder="请输入 Telegram Bot API Key"
                          type="password"
                          show-password-on="click"
                          class="flex-1"
                          @input="handleBotConfigChange"
                        />
                        <n-button
                          type="primary"
                          :loading="validatingApiKey"
                          @click="validateApiKey"
                        >
                          校验
                        </n-button>
                      </div>
                      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        从 @BotFather 获取 API Key
                      </p>
                    </div>

                    <!-- 校验结果 -->
                    <div v-if="apiKeyValidationResult" class="p-3 rounded-md"
                         :class="apiKeyValidationResult.valid ? 'bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300' : 'bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-300'">
                      <div class="flex items-center">
                        <i :class="apiKeyValidationResult.valid ? 'fas fa-check-circle' : 'fas fa-times-circle'"
                           class="mr-2"></i>
                        <span>{{ apiKeyValidationResult.message }}</span>
                      </div>
                      <div v-if="apiKeyValidationResult.valid && apiKeyValidationResult.botInfo" class="mt-2 text-xs">
                        机器人：@{{ apiKeyValidationResult.botInfo.username }} ({{ apiKeyValidationResult.botInfo.first_name }})
                      </div>
                    </div>

                    <div class="flex justify-end">
                      <n-button
                        type="primary"
                        :loading="savingBotConfig"
                        :disabled="!hasBotConfigChanges"
                        @click="saveBotConfig"
                      >
                        保存配置
                      </n-button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 自动回复配置 -->
              <div v-if="telegramBotConfig.bot_enabled" class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <div class="flex items-center mb-6">
                  <div class="w-8 h-8 bg-green-600 text-white rounded-full flex items-center justify-center mr-3">
                    <span class="text-sm font-bold">2</span>
                  </div>
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-white">自动回复设置</h3>
                </div>

                <div class="space-y-4">
                  <!-- 自动回复开关 -->
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700 dark:text-gray-300">启用自动回复</label>
                      <p class="text-xs text-gray-500 dark:text-gray-400">收到消息时自动回复帮助信息</p>
                    </div>
                    <n-switch
                      v-model:value="telegramBotConfig.auto_reply_enabled"
                      @update:value="handleBotConfigChange"
                    />
                  </div>

                  <!-- 回复模板 -->
                  <div v-if="telegramBotConfig.auto_reply_enabled">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">回复模板</label>
                    <n-input
                      v-model:value="telegramBotConfig.auto_reply_template"
                      type="textarea"
                      placeholder="请输入自动回复内容"
                      :rows="3"
                      @input="handleBotConfigChange"
                    />
                  </div>

                  <!-- 自动删除开关 -->
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700 dark:text-gray-300">自动删除回复</label>
                      <p class="text-xs text-gray-500 dark:text-gray-400">定时删除机器人发送的回复消息</p>
                    </div>
                    <n-switch
                      v-model:value="telegramBotConfig.auto_delete_enabled"
                      @update:value="handleBotConfigChange"
                    />
                  </div>

                  <!-- 删除间隔 -->
                  <div v-if="telegramBotConfig.auto_delete_enabled">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">删除间隔（分钟）</label>
                    <n-input-number
                      v-model:value="telegramBotConfig.auto_delete_interval"
                      :min="1"
                      :max="1440"
                      @update:value="handleBotConfigChange"
                    />
                  </div>
                </div>
              </div>

              <!-- 频道和群组管理 -->
              <div v-if="telegramBotConfig.bot_enabled" class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <div class="flex items-center justify-between mb-6">
                  <div class="flex items-center">
                    <div class="w-8 h-8 bg-purple-600 text-white rounded-full flex items-center justify-center mr-3">
                      <span class="text-sm font-bold">3</span>
                    </div>
                    <div>
                      <h3 class="text-lg font-semibold text-gray-900 dark:text-white">频道和群组管理</h3>
                      <p class="text-sm text-gray-600 dark:text-gray-400">管理推送对象的频道和群组</p>
                    </div>
                  </div>

                  <n-button
                    type="primary"
                    @click="showRegisterChannelDialog = true"
                  >
                    <template #icon>
                      <i class="fas fa-plus"></i>
                    </template>
                    注册频道/群组
                  </n-button>
                </div>

                <!-- 频道列表 -->
                <div v-if="telegramChannels.length > 0" class="space-y-4">
                  <div v-for="channel in telegramChannels" :key="channel.id"
                       class="border border-gray-200 dark:border-gray-600 rounded-lg p-4">
                    <div class="flex items-center justify-between mb-4">
                      <div class="flex items-center space-x-3">
                        <i :class="channel.chat_type === 'channel' ? 'fab fa-telegram-plane' : 'fas fa-users'"
                           class="text-lg text-blue-600 dark:text-blue-400"></i>
                        <div>
                          <h4 class="font-medium text-gray-900 dark:text-white">{{ channel.chat_name }}</h4>
                          <p class="text-sm text-gray-500 dark:text-gray-400">
                            {{ channel.chat_type === 'channel' ? '频道' : '群组' }} • ID: {{ channel.chat_id }}
                          </p>
                        </div>
                      </div>

                      <div class="flex items-center space-x-2">
                        <n-tag :type="channel.is_active ? 'success' : 'warning'" size="small">
                          {{ channel.is_active ? '活跃' : '非活跃' }}
                        </n-tag>
                        <n-button size="small" @click="editChannel(channel)">
                          <template #icon>
                            <i class="fas fa-edit"></i>
                          </template>
                        </n-button>
                        <n-button size="small" type="error" @click="unregisterChannel(channel)">
                          <template #icon>
                            <i class="fas fa-trash"></i>
                          </template>
                        </n-button>
                      </div>
                    </div>

                    <!-- 推送配置 -->
                    <div v-if="channel.push_enabled" class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
                      <div>
                        <label class="text-sm font-medium text-gray-700 dark:text-gray-300">推送频率</label>
                        <p class="text-sm text-gray-600 dark:text-gray-400">{{ channel.push_frequency }} 小时</p>
                      </div>
                      <div>
                        <label class="text-sm font-medium text-gray-700 dark:text-gray-300">内容分类</label>
                        <p class="text-sm text-gray-600 dark:text-gray-400">
                          {{ channel.content_categories || '全部' }}
                        </p>
                      </div>
                      <div>
                        <label class="text-sm font-medium text-gray-700 dark:text-gray-300">标签</label>
                        <p class="text-sm text-gray-600 dark:text-gray-400">
                          {{ channel.content_tags || '全部' }}
                        </p>
                      </div>
                    </div>

                    <div v-else class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
                      <p class="text-sm text-gray-500 dark:text-gray-400">推送已禁用</p>
                    </div>
                  </div>
                </div>

                <!-- 空状态 -->
                <div v-else class="text-center py-8">
                  <i class="fab fa-telegram-plane text-4xl text-gray-400 mb-4"></i>
                  <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无频道或群组</h3>
                  <p class="text-gray-500 dark:text-gray-400 mb-4">点击上方按钮注册推送对象</p>
                  <n-button type="primary" @click="showRegisterChannelDialog = true">
                    立即注册
                  </n-button>
                </div>
              </div>
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

        <!-- 注册频道对话框 -->
        <n-modal
          v-model:show="showRegisterChannelDialog"
          preset="card"
          title="注册频道/群组"
          size="huge"
          :bordered="false"
          :segmented="false"
        >
          <div class="space-y-6">
            <div class="text-sm text-gray-600 dark:text-gray-400">
              将机器人添加到频道或群组，然后发送命令获取频道信息并注册为推送对象。
            </div>

            <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
              <div class="flex items-start space-x-3">
                <i class="fas fa-info-circle text-blue-600 dark:text-blue-400 mt-1"></i>
                <div>
                  <h4 class="text-sm font-medium text-blue-800 dark:text-blue-200 mb-2">注册步骤：</h4>
                  <ol class="text-sm text-blue-700 dark:text-blue-300 space-y-1 list-decimal list-inside">
                    <li>将 @{{ telegramBotConfig.bot_enabled ? '机器人用户名' : '机器人' }} 添加为频道管理员或群组成员</li>
                    <li>在频道/群组中发送 <code class="bg-blue-200 dark:bg-blue-800 px-1 rounded">/register</code> 命令</li>
                    <li>机器人将自动识别并注册该频道/群组</li>
                  </ol>
                </div>
              </div>
            </div>

            <div v-if="!telegramBotConfig.bot_enabled || !telegramBotConfig.bot_api_key" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
              <div class="flex items-start space-x-3">
                <i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400 mt-1"></i>
                <div>
                  <h4 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">配置未完成</h4>
                  <p class="text-sm text-yellow-700 dark:text-yellow-300">请先启用机器人并配置有效的 API Key。</p>
                </div>
              </div>
            </div>

            <div class="text-center py-4">
              <n-button
                type="primary"
                @click="showRegisterChannelDialog = false"
              >
                我知道了
              </n-button>
            </div>
          </div>
        </n-modal>
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

// Telegram 相关数据和状态
const telegramBotConfig = ref({
  bot_enabled: false,
  bot_api_key: '',
  auto_reply_enabled: true,
  auto_reply_template: '您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。',
  auto_delete_enabled: false,
  auto_delete_interval: 60,
})

const telegramChannels = ref([])
const validatingApiKey = ref(false)
const savingBotConfig = ref(false)
const apiKeyValidationResult = ref(null)
const hasBotConfigChanges = ref(false)
const showRegisterChannelDialog = ref(false)

// 获取 Telegram 配置
const fetchTelegramConfig = async () => {
  try {
    const response = await $fetch('/api/telegram/bot-config', {
      method: 'GET'
    })
    if (response && response.data) {
      telegramBotConfig.value = { ...response.data }
    }
  } catch (error) {
    console.error('获取 Telegram 配置失败:', error)
  }
}

// 获取频道列表
const fetchTelegramChannels = async () => {
  try {
    const response = await $fetch('/api/telegram/channels', {
      method: 'GET'
    })
    if (response && response.data) {
      telegramChannels.value = response.data
    }
  } catch (error) {
    console.error('获取频道列表失败:', error)
  }
}

// 处理机器人配置变更
const handleBotConfigChange = () => {
  hasBotConfigChanges.value = true
}

// 校验 API Key
const validateApiKey = async () => {
  if (!telegramBotConfig.value.bot_api_key) {
    notification.error({
      content: '请输入 API Key',
      duration: 2000
    })
    return
  }

  validatingApiKey.value = true
  try {
    const response = await $fetch('/api/telegram/validate-api-key', {
      method: 'POST',
      body: {
        api_key: telegramBotConfig.value.bot_api_key
      }
    })

    if (response && response.data) {
      apiKeyValidationResult.value = response.data
      if (response.data.valid) {
        notification.success({
          content: 'API Key 校验成功',
          duration: 2000
        })
      } else {
        notification.error({
          content: response.data.error,
          duration: 3000
        })
      }
    }
  } catch (error) {
    apiKeyValidationResult.value = {
      valid: false,
      error: error.data?.message || '校验失败'
    }
    notification.error({
      content: 'API Key 校验失败',
      duration: 2000
    })
  } finally {
    validatingApiKey.value = false
  }
}

// 保存机器人配置
const saveBotConfig = async () => {
  savingBotConfig.value = true
  try {
    const configRequest = {}
    if (hasBotConfigChanges.value) {
      configRequest.bot_enabled = telegramBotConfig.value.bot_enabled
      configRequest.bot_api_key = telegramBotConfig.value.bot_api_key
      configRequest.auto_reply_enabled = telegramBotConfig.value.auto_reply_enabled
      configRequest.auto_reply_template = telegramBotConfig.value.auto_reply_template
      configRequest.auto_delete_enabled = telegramBotConfig.value.auto_delete_enabled
      configRequest.auto_delete_interval = telegramBotConfig.value.auto_delete_interval
    }

    const response = await $fetch('/api/telegram/bot-config', {
      method: 'PUT',
      body: configRequest
    })

    if (response && response.success) {
      notification.success({
        content: '配置保存成功',
        duration: 2000
      })
      hasBotConfigChanges.value = false
      // 重新获取配置以确保同步
      await fetchTelegramConfig()
    }
  } catch (error) {
    notification.error({
      content: error.data?.message || '配置保存失败',
      duration: 3000
    })
  } finally {
    savingBotConfig.value = false
  }
}

// 编辑频道
const editChannel = (channel) => {
  // TODO: 实现编辑频道功能
  console.log('编辑频道:', channel)
}

// 注销频道
const unregisterChannel = async (channel) => {
  try {
    await notification.warning({
      content: `确定取消注册 "${channel.chat_name}" 吗？`,
      duration: 10000,
      action: () => performUnregisterChannel(channel)
    })
  } catch (error) {
    console.error('注销频道操作失败:', error)
  }
}

const performUnregisterChannel = async (channel) => {
  try {
    const response = await $fetch(`/api/telegram/channels/${channel.id}`, {
      method: 'DELETE'
    })

    if (response && response.success) {
      notification.success({
        content: '频道已取消注册',
        duration: 2000
      })
      // 重新获取频道列表
      await fetchTelegramChannels()
    }
  } catch (error) {
    notification.error({
      content: error.data?.message || '取消注册失败',
      duration: 3000
    })
  }
}

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
onMounted(async () => {
  fetchApiToken()
  await fetchTelegramConfig()
  await fetchTelegramChannels()
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