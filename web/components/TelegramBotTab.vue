<template>
  <div class="tab-content-container h-full flex flex-col overflow-hidden">
    <div class="space-y-8 h-1 flex-1 overflow-auto">
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
                <span>{{  apiKeyValidationResult.valid ? '' : 'Fail' }}</span>
                <span v-if="apiKeyValidationResult.valid && apiKeyValidationResult.bot_info" class="mt-2 text-xs">
                  机器人：@{{ apiKeyValidationResult.bot_info.username }} ({{ apiKeyValidationResult.bot_info.first_name }})
                </span>
              </div>
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

          <div class="flex items-center space-x-2">
            <n-button
              @click="refreshChannels"
              :loading="refreshingChannels"
            >
              <template #icon>
                <i class="fas fa-sync-alt"></i>
              </template>
              刷新
            </n-button>
            <n-button
              @click="testBotConnection"
              :loading="testingConnection"
            >
              <template #icon>
                <i class="fas fa-robot"></i>
              </template>
              测试连接
            </n-button>
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
    <div class="flex justify-end p-2">
       <n-button @click="showLogDrawer = true">
         <template #icon>
           <i class="fas fa-list-alt"></i>
         </template>
         日志
       </n-button>
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

  <!-- Telegram 日志抽屉 -->
  <n-drawer
    v-model:show="showLogDrawer"
    title="Telegram 机器人日志"
    width="80%"
    placement="right"
  >
    <n-drawer-content>
      <div class="space-y-4">
        <!-- 日志控制栏 -->
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-600 dark:text-gray-400">时间范围:</span>
            <n-select
              v-model:value="logHours"
              :options="[
                { label: '1小时', value: 1 },
                { label: '6小时', value: 6 },
                { label: '24小时', value: 24 },
                { label: '72小时', value: 72 }
              ]"
              size="small"
              style="width: 100px"
              @update:value="refreshLogs"
            />
          </div>
          <div class="flex items-center space-x-2">
            <n-button size="small" @click="refreshLogs" :loading="loadingLogs">
              <template #icon>
                <i class="fas fa-sync-alt"></i>
              </template>
              刷新
            </n-button>
          </div>
        </div>

        <!-- 日志列表 -->
        <div class="space-y-2 max-h-96 overflow-y-auto">
          <div v-if="telegramLogs.length === 0 && !loadingLogs" class="text-center py-8">
            <i class="fas fa-list-alt text-4xl text-gray-400 mb-4"></i>
            <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">暂无日志</h3>
            <p class="text-gray-500 dark:text-gray-400">机器人运行日志将显示在这里</p>
          </div>

          <div v-else-if="loadingLogs" class="text-center py-8">
            <n-spin size="large" />
            <p class="text-gray-500 dark:text-gray-400 mt-4">加载日志中...</p>
          </div>

          <div v-for="log in telegramLogs" :key="`${log.timestamp}-${Math.random()}`"
               class="flex items-start space-x-3 p-3 rounded-lg border"
               :class="getLogItemClass(log.level)">
            <div class="flex-shrink-0">
              <i :class="getLogIcon(log.level)" class="text-lg"></i>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-2 mb-1">
                <span class="text-xs font-medium" :class="getLogLevelClass(log.level)">
                  {{ log.level.toUpperCase() }}
                </span>
                <span class="text-xs text-gray-400">{{ formatTimestamp(log.timestamp) }}</span>
                <n-tag v-if="log.category" size="small" :type="getCategoryTagType(log.category)">
                  {{ getCategoryLabel(log.category) }}
                </n-tag>
              </div>
              <p class="text-sm text-gray-900 dark:text-white break-words">{{ log.message }}</p>
            </div>
          </div>
        </div>

        <!-- 日志统计 -->
        <div class="flex justify-between items-center text-sm text-gray-600 dark:text-gray-400">
          <span>显示 {{ telegramLogs.length }} 条日志</span>
          <span v-if="telegramLogs.length > 0">
            加载于 {{ formatTimestamp(new Date().toISOString()) }}
          </span>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useNotification, useDialog } from 'naive-ui'
import { useTelegramApi } from '~/composables/useApi'

// Telegram 相关数据和状态
const telegramBotConfig = ref<any>({
  bot_enabled: false,
  bot_api_key: '',
  auto_reply_enabled: true,
  auto_reply_template: '您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。',
  auto_delete_enabled: false,
  auto_delete_interval: 60,
})

const telegramChannels = ref<any[]>([])
const validatingApiKey = ref(false)
const savingBotConfig = ref(false)
const apiKeyValidationResult = ref<any>(null)
const hasBotConfigChanges = ref(false)
const showRegisterChannelDialog = ref(false)
const showLogDrawer = ref(false)
const refreshingChannels = ref(false)
const testingConnection = ref(false)
const telegramLogs = ref<any[]>([])
const loadingLogs = ref(false)
const logHours = ref(24)

// 使用统一的Telegram API
const telegramApi = useTelegramApi()

// 获取 Telegram 配置
const fetchTelegramConfig = async () => {
  try {
    const data = await telegramApi.getBotConfig() as any
    if (data) {
      telegramBotConfig.value = { ...data }
    }
  } catch (error) {
    console.error('获取 Telegram 配置失败:', error)
  }
}

// 获取频道列表
const fetchTelegramChannels = async () => {
  try {
    const data = await telegramApi.getChannels() as any[]
    if (data) {
      telegramChannels.value = data
    }
  } catch (error: any) {
    console.error('获取频道列表失败:', error)
    // 如果是表不存在的错误，给出更友好的提示
    if (error?.message?.includes('telegram_channels') ||
        error?.message?.includes('does not exist') ||
        error?.message?.includes('relation') && error?.message?.includes('does not exist')) {
      notification.error({
        content: '频道列表表不存在，请重启服务器以创建表',
        duration: 5000
      })
    } else {
      notification.error({
        content: '获取频道列表失败，请稍后重试',
        duration: 3000
      })
    }
    // 清空列表以显示空状态
    telegramChannels.value = []
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
    const data = await telegramApi.validateApiKey({
      api_key: telegramBotConfig.value.bot_api_key
    }) as any

    console.log('API Key 校验结果:', data)
    if (data) {
      apiKeyValidationResult.value = data
      if (data.valid) {
        notification.success({
          content: 'API Key 校验成功',
          duration: 2000
        })
      } else {
        notification.error({
          content: data.error,
          duration: 3000
        })
      }
    }
  } catch (error: any) {
    apiKeyValidationResult.value = {
      valid: false,
      error: error?.message || '校验失败'
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

  // 先校验key 是否有效
  try {
    if (telegramBotConfig.value.bot_enabled) {
      const data = await telegramApi.validateApiKey({
        api_key: telegramBotConfig.value.bot_api_key
      }) as any

      console.log('API Key 校验结果:', data)
      if (data) {
        apiKeyValidationResult.value = data
        if (!data.valid) {
          notification.error({
            content: data.error,
            duration: 3000
          })
          return
        }
      }
    }

  } catch (error: any) {
    apiKeyValidationResult.value = {
      valid: false,
      error: error?.message || '校验失败'
    }
    notification.error({
      content: 'API Key 校验失败',
      duration: 2000
    })
    return
  }

  try {
    const configRequest: any = {}
    if (hasBotConfigChanges.value) {
      const config = telegramBotConfig.value as any
      configRequest.bot_enabled = config.bot_enabled
      configRequest.bot_api_key = config.bot_api_key
      configRequest.auto_reply_enabled = config.auto_reply_enabled
      configRequest.auto_reply_template = config.auto_reply_template
      configRequest.auto_delete_enabled = config.auto_delete_enabled
      configRequest.auto_delete_interval = config.auto_delete_interval
    }

    await telegramApi.updateBotConfig(configRequest)

    notification.success({
      content: '配置保存成功',
      duration: 2000
    })
    hasBotConfigChanges.value = false
    // 重新获取配置以确保同步
    await fetchTelegramConfig()
  } catch (error: any) {
    notification.error({
      content: error?.message || '配置保存失败',
      duration: 3000
    })
  } finally {
    savingBotConfig.value = false
  }
}

// 编辑频道
const editChannel = (channel: any) => {
  // TODO: 实现编辑频道功能
  console.log('编辑频道:', channel)
}

// 注销频道（带确认）
const unregisterChannel = async (channel: any) => {
  try {
    // 使用 Naïve UI 的确认对话框
    dialog.create({
      title: '确认注销频道',
      content: `确定要注销频道 "${channel.chat_name}" 吗？\n\n此操作将停止向该频道推送内容，并会向频道发送注销通知。`,
      positiveText: '确定注销',
      negativeText: '取消',
      type: 'warning',
      onPositiveClick: async () => {
        await performUnregisterChannel(channel)
      },
      onNegativeClick: () => {
        console.log('用户取消了注销操作')
      }
    })
  } catch (error) {
    console.error('创建确认对话框失败:', error)
  }
}

const performUnregisterChannel = async (channel: any) => {
  try {
    await telegramApi.deleteChannel(channel.id)

    notification.success({
      content: `频道 "${channel.chat_name}" 已成功注销`,
      duration: 3000
    })

    // 重新获取频道列表，更新UI
    await fetchTelegramChannels()

    // 尝试向频道发送通知（可选）
    await sendChannelNotification(channel)

  } catch (error: any) {
    console.error('注销频道失败:', error)

    // 提供更详细的错误信息
    let errorMessage = '取消注册失败'
    if (error?.message?.includes('telegram_channels') ||
        error?.message?.includes('does not exist')) {
      errorMessage = '频道表不存在，请重启服务器创建表'
    } else if (error?.message) {
      errorMessage = `注销失败: ${error.message}`
    }

    notification.error({
      content: errorMessage,
      duration: 4000
    })

    // 如果删除失败，仍然尝试刷新列表以确保UI同步
    try {
      await fetchTelegramChannels()
    } catch (refreshError) {
      console.warn('刷新频道列表失败:', refreshError)
    }
  }
}

// 向频道发送注销通知
const sendChannelNotification = async (channel: any) => {
  try {
    const message = `📢 **频道注销通知**\n\n频道 **${channel.chat_name}** 已从机器人推送系统中移除。\n\n❌ 停止推送：此频道将不会再收到自动推送内容\n\n💡 如需继续接收推送，请联系管理员重新注册此频道。`

    await telegramApi.testBotMessage({
      chat_id: channel.chat_id,
      text: message
    })

    notification.success({
      content: `已向频道 "${channel.chat_name}" 发送注销通知`,
      duration: 3000
    })

    console.log(`已向频道 ${channel.chat_name} 发送注销通知`)
  } catch (error: any) {
    console.warn(`向频道 ${channel.chat_name} 发送通知失败:`, error)
    notification.warning({
      content: `向频道 "${channel.chat_name}" 发送通知失败，但频道已从系统中移除`,
      duration: 4000
    })
    // 不抛出错误，因为主操作（删除频道）已经成功
  }
}

// 刷新频道列表
const refreshChannels = async () => {
  refreshingChannels.value = true
  try {
    await fetchTelegramChannels()
    notification.success({
      content: '频道列表已刷新',
      duration: 2000
    })
  } catch (error) {
    notification.error({
      content: '刷新频道列表失败',
      duration: 2000
    })
  } finally {
    refreshingChannels.value = false
  }
}

// 测试机器人连接
const testBotConnection = async () => {
  testingConnection.value = true
  try {
    const data = await telegramApi.getBotStatus() as any
    if (data && data.service_running) {
      notification.success({
        content: `机器人连接正常！用户名：@${data.bot_username}`,
        duration: 3000
      })
    } else {
      notification.warning({
        content: '机器人服务未运行或未配置',
        duration: 3000
      })
    }
  } catch (error: any) {
    notification.error({
      content: '测试连接失败：' + (error?.message || '请检查配置'),
      duration: 3000
    })
  } finally {
    testingConnection.value = false
  }
}

// 获取 Telegram 日志
const fetchTelegramLogs = async () => {
  loadingLogs.value = true
  try {
    const data = await telegramApi.getLogs({
      hours: logHours.value,
      limit: 100
    }) as any
    if (data && data.logs) {
      telegramLogs.value = data.logs
    }
  } catch (error: any) {
    console.error('获取 Telegram 日志失败:', error)
    notification.error({
      content: '获取日志失败：' + (error?.message || '请稍后重试'),
      duration: 3000
    })
  } finally {
    loadingLogs.value = false
  }
}

// 刷新日志
const refreshLogs = async () => {
  await fetchTelegramLogs()
  notification.success({
    content: '日志已刷新',
    duration: 2000
  })
}

// 获取日志图标
const getLogIcon = (level: string) => {
  switch (level.toLowerCase()) {
    case 'info': return 'fas fa-info-circle text-blue-500'
    case 'warn': return 'fas fa-exclamation-triangle text-yellow-500'
    case 'error': return 'fas fa-times-circle text-red-500'
    case 'fatal': return 'fas fa-skull-crossbones text-red-700'
    default: return 'fas fa-circle text-gray-400'
  }
}

// 格式化时间戳
const formatTimestamp = (timestamp: string) => {
  return new Date(timestamp).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 获取日志项样式类
const getLogItemClass = (level: string) => {
  switch (level.toLowerCase()) {
    case 'error': return 'border-red-200 bg-red-50 dark:bg-red-900/10'
    case 'warn': return 'border-yellow-200 bg-yellow-50 dark:bg-yellow-900/10'
    case 'info': return 'border-blue-200 bg-blue-50 dark:bg-blue-900/10'
    default: return 'border-gray-200 bg-gray-50 dark:bg-gray-900/10'
  }
}

// 获取日志级别样式类
const getLogLevelClass = (level: string) => {
  switch (level.toLowerCase()) {
    case 'error': return 'text-red-600 dark:text-red-400'
    case 'warn': return 'text-yellow-600 dark:text-yellow-400'
    case 'info': return 'text-blue-600 dark:text-blue-400'
    default: return 'text-gray-600 dark:text-gray-400'
  }
}

// 获取分类标签类型
const getCategoryTagType = (category: string): "success" | "error" | "warning" | "default" | "primary" | "info" => {
  switch (category?.toLowerCase()) {
    case 'push': return 'success'
    case 'message': return 'info'
    case 'channel': return 'warning'
    case 'service': return 'default'
    default: return 'default'
  }
}

// 获取类别标签文字
const getCategoryLabel = (category: string): string => {
  switch (category?.toLowerCase()) {
    case 'push': return '推送'
    case 'message': return '消息'
    case 'channel': return '频道'
    case 'service': return '服务'
    default: return category || '通用'
  }
}

const notification = useNotification()
const dialog = useDialog()

// 页面加载时获取配置
onMounted(async () => {
  await fetchTelegramConfig()
  await fetchTelegramChannels()
  console.log('Telegram 机器人标签已加载')
})

// 监听日志抽屉打开状态，打开时获取日志
watch(showLogDrawer, async (newValue) => {
  if (newValue && telegramBotConfig.value.bot_enabled) {
    await refreshLogs()
  }
})
</script>

<style scoped>
/* Telegram 机器人标签样式 */
</style>