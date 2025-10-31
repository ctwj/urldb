<template>
  <div class="tab-content-container h-full flex flex-col overflow-hidden">
    <div class="space-y-8 h-1 flex-1 overflow-auto">
      <!-- 机器人基本配置 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">机器人配置</h3>
        </div>

        <div class="space-y-4">
          <!-- 机器人启用开关 -->
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300">启用 Telegram 机器人</label>
              <p class="text-xs text-gray-500 dark:text-gray-400">开启后机器人将开始工作</p>
            </div>
            <div class="flex items-center space-x-3">
              <div v-if="botStatus" class="flex items-center space-x-2">
                <n-tag
                  :type="botStatus.overall_status ? 'success' : (telegramBotConfig.bot_enabled ? 'warning' : 'default')"
                  size="small"
                  class="min-w-16 text-center"
                >
                  {{ botStatus.status_text }}
                </n-tag>
                <!-- 当机器人已启用但未运行时，显示启动按钮 -->
                <n-button
                  v-if="telegramBotConfig.bot_enabled && !botStatus.overall_status"
                  size="small"
                  type="primary"
                  @click="startBotService"
                  :loading="startingBot"
                >
                  <template #icon>
                    <i class="fas fa-play"></i>
                  </template>
                  启动
                </n-button>
                <n-button
                  size="small"
                  @click="refreshBotStatus"
                  :loading="statusRefreshing"
                  circle
                >
                  <template #icon>
                    <i class="fas fa-sync-alt"></i>
                  </template>
                </n-button>
                <n-switch
                  v-model:value="telegramBotConfig.bot_enabled"
                  @update:value="handleBotConfigChange"
                />
              </div>
            </div>
          </div>

          <!-- 启用代理 -->
          <div class="flex items-center justify-between">
            <div>
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300">启用代理</label>
              <p class="text-xs text-gray-500 dark:text-gray-400">通过代理服务器连接 Telegram API</p>
            </div>
            <n-switch
              v-model:value="telegramBotConfig.proxy_enabled"
              @update:value="handleBotConfigChange"
            />
          </div>

          <!-- 代理设置 -->
          <div v-if="telegramBotConfig.proxy_enabled" class="space-y-4">
            <!-- 代理类型 -->
            <div>
              <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">代理类型</label>
              <n-select
                v-model:value="telegramBotConfig.proxy_type"
                :options="[
                  { label: 'HTTP', value: 'http' },
                  { label: 'HTTPS', value: 'https' },
                  { label: 'SOCKS5', value: 'socks5' }
                ]"
                placeholder="选择代理类型"
                @update:value="handleBotConfigChange"
              />
            </div>

            <!-- 代理主机和端口 -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">代理主机</label>
                <n-input
                  v-model:value="telegramBotConfig.proxy_host"
                  placeholder="例如: 127.0.0.1 或 proxy.example.com"
                  @input="handleBotConfigChange"
                />
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">代理端口</label>
                <n-input-number
                  v-model:value="telegramBotConfig.proxy_port"
                  :min="1"
                  :max="65535"
                  placeholder="例如: 8080"
                  @update:value="handleBotConfigChange"
                />
              </div>
            </div>

            <!-- 代理认证 -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">用户名 (可选)</label>
                <n-input
                  v-model:value="telegramBotConfig.proxy_username"
                  placeholder="代理用户名"
                  @input="handleBotConfigChange"
                />
              </div>
              <div>
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">密码 (可选)</label>
                <n-input
                  v-model:value="telegramBotConfig.proxy_password"
                  type="password"
                  placeholder="代理密码"
                  @input="handleBotConfigChange"
                />
              </div>
            </div>

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
              :disabled="telegramBotConfig.bot_enabled"
              @update:value="handleBotConfigChange"
            />
          </div>

          <!-- 回复模板 -->
          <div v-if="telegramBotConfig.auto_reply_enabled || telegramBotConfig.bot_enabled">
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
            <div v-if="channel.push_enabled" class="flex flex-wrap gap-6 mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
              <div class="flex-1 min-w-0">
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300">推送频率</label>
                <p class="text-sm text-gray-600 dark:text-gray-400">{{ channel.push_frequency }} 分钟</p>
              </div>
              <div class="flex-1 min-w-0">
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300">资源策略</label>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ getResourceStrategyLabel(channel.resource_strategy) }}
                </p>
              </div>
              <div class="flex-1 min-w-0">
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300">时间限制</label>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ getTimeLimitLabel(channel.time_limit) }}
                </p>
              </div>
              <div class="flex-1 min-w-0">
                <label class="text-sm font-medium text-gray-700 dark:text-gray-300">推送时间段</label>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ channel.push_start_time && channel.push_end_time ? `${channel.push_start_time}-${channel.push_end_time}` : '全天' }}
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
    <div class="flex justify-end p-2 gap-2">
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
        @click="debugBotConnection"
      >
        <template #icon>
          <i class="fas fa-bug"></i>
        </template>
        调试
      </n-button>
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
    :bordered="false"
    :segmented="false"
    :style="{ width: '800px' }"
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

  <!-- 编辑频道对话框 -->
  <n-modal
    v-model:show="showEditChannelDialog"
    preset="card"
    :title="`编辑频道 - ${editingChannel?.chat_name || ''}`"
    size="large"
    :bordered="false"
    :segmented="false"
  >
    <div v-if="editingChannel" class="space-y-6">
      <!-- 频道基本信息 -->
      <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
        <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-2">频道信息</h4>
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-600 dark:text-gray-400">频道名称:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ editingChannel.chat_name }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">频道ID:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ editingChannel.chat_id }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">类型:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ editingChannel.chat_type === 'channel' ? '频道' : '群组' }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">状态:</span>
            <n-tag :type="editingChannel.is_active ? 'success' : 'warning'" size="small" class="ml-2">
              {{ editingChannel.is_active ? '活跃' : '非活跃' }}
            </n-tag>
          </div>
        </div>
      </div>

      <!-- 推送设置 -->
      <div class="space-y-4">
        <h4 class="text-base font-medium text-gray-900 dark:text-white">推送设置</h4>

        <!-- 启用推送 -->
        <div class="flex items-center justify-between">
          <div>
            <label class="text-sm font-medium text-gray-700 dark:text-gray-300">启用推送</label>
            <p class="text-xs text-gray-500 dark:text-gray-400">是否向此频道推送内容</p>
          </div>
          <n-switch
            v-model:value="editingChannel.push_enabled"
          />
        </div>

        <!-- 推送频率 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">推送频率</label>
          <n-select
            v-model:value="editingChannel.push_frequency"
            :options="[
              { label: '每1分钟', value: 1 },
              { label: '每2分钟', value: 2 },
              { label: '每3分钟', value: 3 },
              { label: '每5分钟', value: 5 },
              { label: '每10分钟', value: 10 },
              { label: '每15分钟', value: 15 },
              { label: '每20分钟', value: 20 },
              { label: '每30分钟', value: 30 },
              { label: '每45分钟', value: 45 },
              { label: '每60分钟', value: 60 }
            ]"
            placeholder="选择推送频率"
          />
        </div>

        <!-- 资源策略 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">资源策略</label>
          <n-select
            v-model:value="editingChannel.resource_strategy"
            :options="[
              { label: '纯随机', value: 'random' },
              { label: '最新优先', value: 'latest' },
              { label: '已转存优先', value: 'transferred' }
            ]"
            placeholder="选择资源策略"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
            纯随机：完全随机推送资源；最新优先：优先推送最新资源；已转存优先：优先推送已转存的资源
          </p>
        </div>

        <!-- 时间限制 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">时间限制</label>
          <n-select
            v-model:value="editingChannel.time_limit"
            :options="[
              { label: '无限制', value: 'none' },
              { label: '一周内', value: 'week' },
              { label: '一月内', value: 'month' }
            ]"
            placeholder="选择时间限制"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
            无限制：推送所有时间段的资源；一周内：仅推送最近一周的资源；一月内：仅推送最近一个月的资源
          </p>
        </div>

        <!-- 推送时间段 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">推送时间段</label>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-xs text-gray-600 dark:text-gray-400">开始时间</label>
              <n-time-picker
                v-model:value="editingChannel.push_start_time"
                format="HH:mm"
                placeholder="选择开始时间"
                clearable
                :value-format="'HH:mm'"
                :actions="['clear', 'confirm']"
              />
            </div>
            <div>
              <label class="text-xs text-gray-600 dark:text-gray-400">结束时间</label>
              <n-time-picker
                v-model:value="editingChannel.push_end_time"
                format="HH:mm"
                placeholder="选择结束时间"
                clearable
                :value-format="'HH:mm'"
                :actions="['clear', 'confirm']"
              />
            </div>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
            留空表示全天推送，不设置时间限制
          </p>
        </div>

        <!-- 内容分类 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">内容分类</label>
          <n-input
            v-model:value="editingChannel.content_categories"
            placeholder="输入内容分类，多个用逗号分隔 (如: 电影,电视剧,动漫)"
            type="textarea"
            :rows="2"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">留空表示推送所有分类的内容</p>
        </div>

        <!-- 标签过滤 -->
        <div v-if="editingChannel.push_enabled">
          <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 block">标签过滤</label>
          <n-input
            v-model:value="editingChannel.content_tags"
            placeholder="输入标签关键词，多个用逗号分隔 (如: 高清,1080p,蓝光)"
            type="textarea"
            :rows="2"
          />
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">留空表示推送所有标签的内容</p>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-600">
        <n-button @click="showEditChannelDialog = false">
          取消
        </n-button>
        <n-button
          type="primary"
          :loading="savingChannel"
          @click="saveChannelSettings"
        >
          保存设置
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
      <div class="space-y-4 h-full overflow-y-auto flex flex-col">
        <!-- 日志控制栏 -->
        <div class="flex-0 flex items-center justify-between">
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
        <div class="h-1 flex-1 space-y-2 overflow-y-auto">
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
        <div class="flex-0 flex justify-between items-center text-sm text-gray-600 dark:text-gray-400">
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
  proxy_enabled: false,
  proxy_type: 'http',
  proxy_host: '',
  proxy_port: 8080,
  proxy_username: '',
  proxy_password: '',
})

const telegramChannels = ref<any[]>([])
const validatingApiKey = ref(false)
const savingBotConfig = ref(false)
const apiKeyValidationResult = ref<any>(null)
const hasBotConfigChanges = ref(false)
const showRegisterChannelDialog = ref(false)
const showEditChannelDialog = ref(false)
const showLogDrawer = ref(false)
const refreshingChannels = ref(false)
const testingConnection = ref(false)
const telegramLogs = ref<any[]>([])
const loadingLogs = ref(false)
const logHours = ref(24)
const editingChannel = ref<any>(null)
const savingChannel = ref(false)

// 机器人状态相关变量
const botStatus = ref<any>(null)
const statusRefreshing = ref(false)
const startingBot = ref(false)

// 使用统一的Telegram API
const telegramApi = useTelegramApi()

// 获取 Telegram 配置
const fetchTelegramConfig = async () => {
  try {
    const data = await telegramApi.getBotConfig() as any
    if (data) {
      // 确保当机器人启用时，自动回复始终为true
      if (data.bot_enabled) {
        data.auto_reply_enabled = true
      }
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
    if (data !== undefined && data !== null) {
      telegramChannels.value = Array.isArray(data) ? data : []
    } else {
      telegramChannels.value = []
    }
    console.log('频道列表已更新:', telegramChannels.value.length, '个频道')
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
  // 当机器人启用时，自动回复必须为true
  if (telegramBotConfig.value.bot_enabled) {
    telegramBotConfig.value.auto_reply_enabled = true
  }
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
    // 构建校验请求，包含代理配置
    const validateRequest: any = {
      api_key: telegramBotConfig.value.bot_api_key
    }

    // 如果启用了代理，包含代理配置
    if (telegramBotConfig.value.proxy_enabled) {
      validateRequest.proxy_enabled = telegramBotConfig.value.proxy_enabled
      validateRequest.proxy_type = telegramBotConfig.value.proxy_type
      validateRequest.proxy_host = telegramBotConfig.value.proxy_host
      validateRequest.proxy_port = telegramBotConfig.value.proxy_port
      validateRequest.proxy_username = telegramBotConfig.value.proxy_username
      validateRequest.proxy_password = telegramBotConfig.value.proxy_password
    }

    const data = await telegramApi.validateApiKey(validateRequest) as any

    console.log('API Key 校验结果:', data)
    if (data) {
      apiKeyValidationResult.value = data
      if (data.valid) {
        // 显示校验成功的提示，如果使用了代理则特别说明
        let successMessage = 'API Key 校验成功'
        if (telegramBotConfig.value.proxy_enabled) {
          successMessage += ' (通过代理)'
        }
        notification.success({
          content: successMessage,
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
      const validateRequest: any = {
        api_key: telegramBotConfig.value.bot_api_key
      }
      if (telegramBotConfig.value.proxy_enabled) {
        validateRequest.proxy_enabled = telegramBotConfig.value.proxy_enabled
        validateRequest.proxy_type = telegramBotConfig.value.proxy_type
        validateRequest.proxy_host = telegramBotConfig.value.proxy_host
        validateRequest.proxy_port = telegramBotConfig.value.proxy_port
        validateRequest.proxy_username = telegramBotConfig.value.proxy_username
        validateRequest.proxy_password = telegramBotConfig.value.proxy_password
      }
      const data = await telegramApi.validateApiKey(validateRequest) as any
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
    savingBotConfig.value = false
    return
  }

  try {
    const configRequest: any = {}
    if (hasBotConfigChanges.value) {
      const config = telegramBotConfig.value as any
      configRequest.bot_enabled = config.bot_enabled
      configRequest.bot_api_key = config.bot_api_key
      // 当机器人启用时，自动回复必须为true
      configRequest.auto_reply_enabled = config.bot_enabled ? true : config.auto_reply_enabled
      configRequest.auto_reply_template = config.auto_reply_template
      configRequest.auto_delete_enabled = config.auto_delete_enabled
      configRequest.auto_delete_interval = config.auto_delete_interval
      configRequest.proxy_enabled = config.proxy_enabled
      configRequest.proxy_type = config.proxy_type
      configRequest.proxy_host = config.proxy_host
      configRequest.proxy_port = config.proxy_port
      configRequest.proxy_username = config.proxy_username
      configRequest.proxy_password = config.proxy_password
    }

    await telegramApi.updateBotConfig(configRequest)

    notification.success({
      content: '配置保存成功，机器人服务已重新加载配置',
      duration: 3000
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
  // 复制频道数据并处理时间字段
  const channelCopy = { ...channel }

  // 处理时间字段，确保时间选择器可以正确显示
  try {
    console.log('处理编辑频道时间字段:')
    console.log('原始开始时间:', channelCopy.push_start_time)
    console.log('原始结束时间:', channelCopy.push_end_time)

    // 处理开始时间
    if (channelCopy.push_start_time) {
      if (isValidTimeString(channelCopy.push_start_time)) {
        // 数据库中是 "HH:mm" 格式的时间字符串
        console.log('开始时间是有效格式，保持原样:', channelCopy.push_start_time)
      } else {
        console.log('开始时间格式无效，设为null')
        channelCopy.push_start_time = null
      }
    } else {
      console.log('开始时间为空，设为null')
      channelCopy.push_start_time = null
    }

    // 处理结束时间
    if (channelCopy.push_end_time) {
      if (isValidTimeString(channelCopy.push_end_time)) {
        // 数据库中是 "HH:mm" 格式的时间字符串
        console.log('结束时间是有效格式，保持原样:', channelCopy.push_end_time)
      } else {
        console.log('结束时间格式无效，设为null')
        channelCopy.push_end_time = null
      }
    } else {
      console.log('结束时间为空，设为null')
      channelCopy.push_end_time = null
    }

    console.log('处理后时间字段:', {
      push_start_time: channelCopy.push_start_time,
      push_end_time: channelCopy.push_end_time
    })

    // 尝试转换为时间戳格式（毫秒），因为时间选择器可能期望这种格式
    if (channelCopy.push_start_time) {
      const timeStr = channelCopy.push_start_time // 格式如 "08:30"
      const parts = timeStr.split(':')
      if (parts.length === 2) {
        const hours = parseInt(parts[0], 10)
        const minutes = parseInt(parts[1], 10)
        // 创建今天的日期，然后设置小时和分钟
        const today = new Date()
        today.setHours(hours, minutes, 0, 0)
        const timestamp = today.getTime()
        console.log('转换开始时间戳:', timestamp)
        channelCopy.push_start_time = timestamp
      }
    }

    if (channelCopy.push_end_time) {
      const timeStr = channelCopy.push_end_time // 格式如 "11:30"
      const parts = timeStr.split(':')
      if (parts.length === 2) {
        const hours = parseInt(parts[0], 10)
        const minutes = parseInt(parts[1], 10)
        // 创建今天的日期，然后设置小时和分钟
        const today = new Date()
        today.setHours(hours, minutes, 0, 0)
        const timestamp = today.getTime()
        console.log('转换结束时间戳:', timestamp)
        channelCopy.push_end_time = timestamp
      }
    }

    console.log('最终时间字段格式:', {
      push_start_time: channelCopy.push_start_time,
      push_end_time: channelCopy.push_end_time
    })
  } catch (error) {
    console.warn('处理频道时间字段时出错:', error)
    channelCopy.push_start_time = null
    channelCopy.push_end_time = null
  }

  editingChannel.value = channelCopy
  showEditChannelDialog.value = true
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

    // 添加短暂延迟确保数据库事务完成
    await new Promise(resolve => setTimeout(resolve, 500))

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
    if (data && data.overall_status) {
      notification.success({
        content: `机器人连接正常！用户名：@${data.runtime?.username || '未知'}`,
        duration: 3000
      })
    } else {
      let warningMessage = '机器人服务未运行或未配置'
      if (data?.config?.enabled) {
        warningMessage = '机器人已启用但未运行，请检查 API Key 配置'
      } else if (!data?.config?.api_key_configured) {
        warningMessage = 'API Key 未配置，请先配置有效的 API Key'
      }
      notification.warning({
        content: warningMessage,
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

// 获取资源策略标签
const getResourceStrategyLabel = (strategy: string): string => {
  switch (strategy) {
    case 'random': return '纯随机'
    case 'latest': return '最新优先'
    case 'transferred': return '已转存优先'
    default: return '纯随机'
  }
}

// 获取时间限制标签
const getTimeLimitLabel = (timeLimit: string): string => {
  switch (timeLimit) {
    case 'none': return '无限制'
    case 'week': return '一周内'
    case 'month': return '一月内'
    default: return '无限制'
  }
}

const notification = useNotification()
const dialog = useDialog()

// 检查时间字符串是否有效
const isValidTimeString = (timeStr: string): boolean => {
  if (!timeStr || typeof timeStr !== 'string') {
    return false
  }

  // 检查 HH:mm 格式
  const timeRegex = /^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/
  if (!timeRegex.test(timeStr)) {
    return false
  }

  return true
}

// 规范化时间字符串
const normalizeTimeString = (timeStr: string): string => {
  if (!timeStr || typeof timeStr !== 'string') {
    return timeStr
  }

  // 确保 HH:mm 格式，补齐前导零
  const parts = timeStr.split(':')
  if (parts.length === 2) {
    const hours = parts[0].padStart(2, '0')
    const minutes = parts[1].padStart(2, '0')
    return `${hours}:${minutes}`
  }

  return timeStr
}

// 保存频道设置
const saveChannelSettings = async () => {
  if (!editingChannel.value) return

  savingChannel.value = true
  try {
    // 处理时间字段，确保保存为字符串格式
    const updateData = {
      chat_id: editingChannel.value.chat_id,
      chat_name: editingChannel.value.chat_name,
      chat_type: editingChannel.value.chat_type,
      push_enabled: editingChannel.value.push_enabled,
      push_frequency: editingChannel.value.push_frequency,
      push_start_time: formatTimeForSave(editingChannel.value.push_start_time),
      push_end_time: formatTimeForSave(editingChannel.value.push_end_time),
      content_categories: editingChannel.value.content_categories,
      content_tags: editingChannel.value.content_tags,
      is_active: editingChannel.value.is_active,
      resource_strategy: editingChannel.value.resource_strategy,
      time_limit: editingChannel.value.time_limit
    }

    console.log('准备提交频道更新数据:', updateData)
    console.log('频道ID:', editingChannel.value.id)
    console.log('推送开始时间原始值:', editingChannel.value.push_start_time)
    console.log('推送结束时间原始值:', editingChannel.value.push_end_time)
    console.log('格式化后推送开始时间:', formatTimeForSave(editingChannel.value.push_start_time))
    console.log('格式化后推送结束时间:', formatTimeForSave(editingChannel.value.push_end_time))
    console.log('资源策略:', editingChannel.value.resource_strategy)
    console.log('时间限制:', editingChannel.value.time_limit)

    await telegramApi.updateChannel(editingChannel.value.id, updateData)
    console.log('频道更新提交完成')

    notification.success({
      content: `频道 "${editingChannel.value.chat_name}" 设置已更新`,
      duration: 3000
    })

    // 关闭对话框
    showEditChannelDialog.value = false

    // 刷新频道列表
    await fetchTelegramChannels()

  } catch (error: any) {
    notification.error({
      content: `保存频道设置失败: ${error?.message || '请稍后重试'}`,
      duration: 3000
    })
  } finally {
    savingChannel.value = false
  }
}

// 格式化时间字段以便保存
const formatTimeForSave = (timeValue: any): string | null => {
  console.log('formatTimeForSave 输入值:', timeValue, '类型:', typeof timeValue)

  if (!timeValue) {
    console.log('formatTimeForSave: 空值，返回 null')
    return null
  }

  // 如果已经是字符串格式，直接返回
  if (typeof timeValue === 'string') {
    console.log('formatTimeForSave: 字符串格式，直接返回:', timeValue)
    return timeValue
  }

  // 如果是数组（Naive UI Time Picker 可能返回这种格式）
  if (Array.isArray(timeValue)) {
    console.log('formatTimeForSave: 数组格式，处理数组:', timeValue)
    if (timeValue.length >= 2) {
      const hours = timeValue[0].toString().padStart(2, '0')
      const minutes = timeValue[1].toString().padStart(2, '0')
      const result = `${hours}:${minutes}`
      console.log('formatTimeForSave: 数组转换为:', result)
      return result
    }
  }

  // 如果是 Date 对象，格式化为 HH:mm
  if (timeValue instanceof Date) {
    const hours = timeValue.getHours().toString().padStart(2, '0')
    const minutes = timeValue.getMinutes().toString().padStart(2, '0')
    const result = `${hours}:${minutes}`
    console.log('formatTimeForSave: Date 对象转换为:', result)
    return result
  }

  // 如果是有 hour 和 minute 属性的对象
  if (timeValue && typeof timeValue === 'object' && 'hour' in timeValue && 'minute' in timeValue) {
    const hours = timeValue.hour.toString().padStart(2, '0')
    const minutes = timeValue.minute.toString().padStart(2, '0')
    const result = `${hours}:${minutes}`
    console.log('formatTimeForSave: 对象格式转换为:', result)
    return result
  }

  // 如果是时间戳（毫秒）
  if (typeof timeValue === 'number' && timeValue > 0) {
    console.log('formatTimeForSave: 时间戳格式，转换为日期')
    const date = new Date(timeValue)
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    const result = `${hours}:${minutes}`
    console.log('formatTimeForSave: 时间戳转换为:', result)
    return result
  }

  console.log('formatTimeForSave: 无法识别的格式，返回 null')
  return null
}

// 启动机器人服务
const startBotService = async () => {
  startingBot.value = true
  try {
    // 重新保存配置以启动机器人
    await saveBotConfig()

    // 等待一秒后刷新状态
    await new Promise(resolve => setTimeout(resolve, 1000))
    await refreshBotStatus()

    notification.success({
      content: '机器人服务启动中，请稍后刷新状态查看',
      duration: 3000
    })
  } catch (error: any) {
    notification.error({
      content: '启动机器人服务失败：' + (error?.message || '请稍后重试'),
      duration: 3000
    })
  } finally {
    startingBot.value = false
  }
}

// 刷新机器人状态
const refreshBotStatus = async () => {
  statusRefreshing.value = true
  try {
    const data = await telegramApi.getBotStatus() as any
    botStatus.value = data
    notification.success({
      content: '机器人状态已刷新',
      duration: 2000
    })
  } catch (error: any) {
    notification.error({
      content: '刷新状态失败：' + (error?.message || '请稍后重试'),
      duration: 3000
    })
  } finally {
    statusRefreshing.value = false
  }
}

// 调试机器人连接
const debugBotConnection = async () => {
  try {
    const data = await telegramApi.getBotStatus() as any

    let message = `🔍 **Telegram 机器人调试信息**\n\n`
    message += `🤖 机器人状态: ${data.runtime?.is_running ? '✅ 运行中' : '❌ 未运行'}\n`
    message += `👤 用户名: @${data.runtime?.username || '未知'}\n`
    message += `⚡ 工作模式: 长轮询\n\n`

    message += `📋 **故障排查步骤:**\n`
    message += `1. 检查服务器控制台是否有 [TELEGRAM] 日志\n`
    message += `2. 确认机器人已添加到群组并设为管理员\n`
    message += `3. 验证 API Key 配置是否正确\n`
    message += `4. 确认自动回复功能已启用\n`
    message += `5. 重启服务器重新加载配置\n\n`

    message += `🔧 **预期日志输出:**\n`
    message += `• [TELEGRAM:SERVICE] Telegram Bot (@用户名) 已启动\n`
    message += `• [TELEGRAM:MESSAGE] 收到消息: ChatID=xxx, Text='/register'\n`
    message += `• [TELEGRAM:MESSAGE] 处理 /register 命令 from ChatID=xxx\n`
    message += `• [TELEGRAM:MESSAGE:SUCCESS] 消息发送成功\n\n`

    message += `💡 **如果没有日志输出:**\n`
    message += `• 服务器可能未正确启动机器人服务\n`
    message += `• API Key 可能有误\n`
    message += `• 数据库配置可能有问题`

    notification.info({
      title: '🤖 机器人连接调试',
      content: message,
      duration: 15000,
      keepAliveOnHover: true
    })
  } catch (error: any) {
    notification.error({
      title: '🔧 调试失败',
      content: `无法获取机器人状态: ${error?.message || '网络错误或服务未运行'}`,
      duration: 5000
    })
  }
}

// 页面加载时获取配置
onMounted(async () => {
  await fetchTelegramConfig()
  await fetchTelegramChannels()
  await refreshBotStatus() // 初始化机器人状态
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