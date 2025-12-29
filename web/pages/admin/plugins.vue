<template>
  <AdminPageLayout>
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
          <i class="fas fa-plug text-blue-500 mr-2"></i>
          插件管理
        </h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统插件，配置插件参数，监控插件运行状态</p>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和操作 -->
    <template #filter-bar>
      <div class="flex justify-between items-center">
        <div class="flex gap-2">
          <n-button @click="showInstallModal = true" type="primary">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            安装插件
          </n-button>
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input
              v-model:value="filters.search"
              @input="debounceSearch"
              type="text"
              placeholder="搜索插件名称..."
              clearable
            >
              <template #prefix>
                <i class="fas fa-search text-gray-400 text-sm"></i>
              </template>
            </n-input>
          </div>
          <n-select
            v-model:value="filters.status"
            :options="[
              { label: '全部状态', value: '' },
              { label: '已启用', value: 'enabled' },
              { label: '已禁用', value: 'disabled' }
            ]"
            placeholder="状态"
            clearable
            @update:value="fetchPlugins"
            style="width: 150px"
          />
          <n-button @click="resetFilters" type="tertiary">
            <template #icon>
              <i class="fas fa-redo"></i>
            </template>
            重置
          </n-button>
          <n-button @click="fetchPlugins" type="tertiary">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区 - 插件数据 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex h-full items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="plugins.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无插件</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">目前没有已安装的插件</div>
      </div>

      <!-- 数据表格 - 自适应高度 -->
      <div v-else class="flex flex-col h-full overflow-auto">
        <n-data-table
          :columns="columns"
          :data="plugins"
          :pagination="false"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          class="h-full"
        />
      </div>
    </template>

  </AdminPageLayout>

  <!-- 插件详情模态框 -->
  <n-modal v-model:show="showDetailModal" :mask-closable="false" preset="card" :style="{ maxWidth: '600px', width: '90%' }" title="插件详情">
    <div v-if="selectedPlugin" class="space-y-4">
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">插件名称</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedPlugin.display_name || selectedPlugin.name }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">版本</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">v{{ selectedPlugin.version }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">描述</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedPlugin.description }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">作者</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedPlugin.author || '未知' }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">分类</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedPlugin.category || 'utility' }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">状态</h3>
        <p class="mt-1">
          <n-tag :type="selectedPlugin.enabled ? 'success' : 'error'" size="small">
            {{ selectedPlugin.enabled ? '已启用' : '已禁用' }}
          </n-tag>
        </p>
      </div>
      <div v-if="selectedPlugin.scheduled_tasks && selectedPlugin.scheduled_tasks.length > 0">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">定时任务</h3>
        <div class="mt-1 space-y-2">
          <div v-for="task in selectedPlugin.scheduled_tasks" :key="task.name" class="text-sm">
            <p class="text-gray-900 dark:text-gray-100">{{ task.name }} - {{ task.schedule }}</p>
            <p class="text-xs text-gray-500">{{ task.frequency?.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </n-modal>

  <!-- 插件配置模态框 -->
  <n-modal
    v-model:show="showConfigModal"
    :mask-closable="false"
    preset="card"
    :style="{ maxWidth: '800px', width: '90%' }"
    :title="`配置 ${configPlugin?.display_name || configPlugin?.name || ''} v${configPlugin?.version || ''}`"
  >
    <div v-if="configPlugin" class="space-y-4">
      <!-- 使用插件配置表单组件 -->
      <PluginConfigForm
        :plugin="configPlugin"
        :config="pluginConfig"
        @save="handleConfigSave"
        @reset="handleConfigReset"
        @cancel="closeConfigModal"
      />
    </div>
  </n-modal>

  <!-- 插件日志模态框 -->
  <n-modal v-model:show="showLogsModal" :mask-closable="false" preset="card" :style="{ maxWidth: '1200px', width: '95%', maxHeight: '90vh' }" title="插件日志">
    <div v-if="logsPlugin" class="space-y-4 h-full">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ logsPlugin.display_name || logsPlugin.name }} 日志</h3>
        <div class="flex items-center space-x-2">
          <n-button @click="refreshLogs" :loading="loadingLogs" size="small" type="info">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>

      <div v-if="loadingLogs" class="flex items-center justify-center py-16">
        <n-spin size="medium" />
      </div>

      <div v-else-if="pluginLogs.length === 0" class="text-center py-16">
        <n-empty description="暂无日志" />
      </div>

      <div v-else class="overflow-hidden" style="height: calc(90vh - 140px);">
        <n-data-table
          :columns="logColumns"
          :data="pluginLogs"
          :pagination="{ pageSize: 50 }"
          :scroll-x="800"
          :max-height="600"
          striped
          size="small"
        />
      </div>
    </div>
  </n-modal>

  <!-- 安装插件模态框 -->
  <n-modal v-model:show="showInstallModal" :mask-closable="false" preset="card" :style="{ maxWidth: '600px', width: '90%' }" title="安装插件">
    <div class="space-y-4">
      <!-- 安装方式选择 -->
      <n-tabs v-model:value="installType" default-value="url" justify-content="center">
        <n-tab-pane name="url" tab="从URL安装">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                插件URL
              </label>
              <n-input
                v-model:value="installUrl"
                type="text"
                placeholder="输入插件包URL (例如: https://example.com/plugin.zip)"
                clearable
              />
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                支持ZIP格式的插件包，包含package.json配置文件
              </p>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="file" tab="从文件安装">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                选择插件文件
              </label>
              <n-upload
                v-model:file-list="installFiles"
                :max="1"
                accept=".zip,.plugin.js"
                @update:file-list="handleFileChange"
              >
                <n-upload-dragger>
                  <div style="margin-bottom: 12px">
                    <i class="fas fa-cloud-upload-alt" style="font-size: 48px; color: #409eff"></i>
                  </div>
                  <n-text style="font-size: 16px">
                    点击或者拖动文件到该区域来上传
                  </n-text>
                  <n-p depth="3" style="margin: 8px 0 0 0">
                    支持 .zip 插件包和 .plugin.js 单文件插件
                  </n-p>
                </n-upload-dragger>
              </n-upload>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                ZIP格式支持多文件插件，单文件格式适合简单插件
              </p>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>

      <!-- 安装进度 -->
      <div v-if="installing" class="space-y-2">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-700 dark:text-gray-300">安装进度</span>
          <span class="text-sm text-gray-500">{{ installProgress }}%</span>
        </div>
        <n-progress :percentage="installProgress" :show-indicator="false" />
        <p class="text-xs text-gray-500">{{ installStatus }}</p>
      </div>

      <!-- 安装结果 -->
      <div v-if="installResult" class="p-4 rounded-lg" :class="installSuccess ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800' : 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'">
        <div class="flex items-center">
          <i :class="installSuccess ? 'fas fa-check-circle text-green-500' : 'fas fa-exclamation-circle text-red-500'" class="mr-2"></i>
          <span class="text-sm font-medium" :class="installSuccess ? 'text-green-700 dark:text-green-400' : 'text-red-700 dark:text-red-400'">
            {{ installResult }}
          </span>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end space-x-2">
        <n-button @click="closeInstallModal" :disabled="installing">
          取消
        </n-button>
        <n-button type="primary" @click="installPlugin" :loading="installing" :disabled="!canInstall">
          安装
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
// 设置页面标题和元信息
useHead({
  title: '插件管理 - 管理后台',
  meta: [
    { name: 'description', content: '管理系统插件，配置插件参数，监控插件运行状态' }
  ]
})

// 设置页面布局和认证保护
definePageMeta({
  layout: 'admin',
  middleware: ['auth', 'admin']
})

import { h } from 'vue'
import PluginConfigForm from '~/components/plugins/PluginConfigForm.vue'
const message = useMessage()
const notification = useNotification()

const loading = ref(false)
const plugins = ref<any[]>([])
const showDetailModal = ref(false)
const showConfigModal = ref(false)
const showLogsModal = ref(false)
const showInstallModal = ref(false)
const selectedPlugin = ref<any>(null)
const configPlugin = ref<any>(null)
const logsPlugin = ref<any>(null)
const pluginConfig = ref({})
const pluginLogs = ref([])
const loadingLogs = ref(false)

// 日志表格列定义
const logColumns = [
  {
    title: '状态',
    key: 'success',
    width: 80,
    render: (row) => {
      return h('span', {
        class: row.success ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
      }, row.success ? '成功' : '失败')
    }
  },
  {
    title: '钩子',
    key: 'hook_name',
    width: 120,
    render: (row) => {
      return h('span', { class: 'font-medium' }, row.hook_name)
    }
  },
  {
    title: '执行时间',
    key: 'execution_time',
    width: 100,
    render: (row) => {
      return h('span', { class: 'text-gray-600 dark:text-gray-400' }, `${row.execution_time}ms`)
    }
  },
  {
    title: '时间',
    key: 'created_at',
    width: 180,
    render: (row) => {
      return h('span', { class: 'text-gray-600 dark:text-gray-400' }, formatTime(row.created_at))
    }
  },
  {
    title: '日志信息',
    key: 'message',
    render: (row) => {
      if (row.error_message) {
        // 显示错误信息（红色）
        return h('span', { class: 'text-red-600 dark:text-red-400 text-sm' }, row.error_message)
      } else if (row.message) {
        // 显示普通日志消息（蓝色）
        return h('span', { class: 'text-blue-600 dark:text-blue-400 text-sm' }, row.message)
      } else {
        // 没有消息内容
        return h('span', { class: 'text-gray-700 dark:text-gray-300 text-sm' }, '执行成功')
      }
    }
  }
]
const saving = ref(false)

// 安装相关状态
const installType = ref('url')
const installUrl = ref('')
const installFiles = ref([])
const installing = ref(false)
const installProgress = ref(0)
const installStatus = ref('')
const installResult = ref('')
const installSuccess = ref(false)

// 分页和筛选状态
const filters = ref({
  status: '',
  search: ''
})

// 计算属性：是否可以安装
const canInstall = computed(() => {
  if (installing.value) return false
  if (installType.value === 'url') {
    return installUrl.value.trim() !== ''
  } else {
    return installFiles.value.length > 0
  }
})

// 表格列定义
const columns = [
  {
    title: '插件信息',
    key: 'name',
    width: 'auto',
    minWidth: 300,
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        // 第一行：名称和版本
        h('div', { class: 'flex items-center gap-2' }, [
          h('i', {
            class: `fas fa-plug text-sm ${row.enabled ? 'text-green-500' : 'text-red-500'}`
          }),
          h('span', { class: 'font-medium text-sm' }, row.display_name || row.name),
          h('span', { class: 'text-xs text-gray-400' }, `v${row.version}`)
        ]),
        // 第二行：描述
        h('div', {
          class: 'text-xs text-gray-500 dark:text-gray-400 line-clamp-2'
        }, row.description || '无描述'),
        // 第三行：作者和分类
        h('div', { class: 'flex items-center gap-2' }, [
          h('span', { class: 'text-xs text-gray-400' }, `作者: ${row.author || '未知'}`),
          h('span', { class: 'text-xs text-gray-400' }, `分类: ${row.category || 'utility'}`)
        ])
      ])
    }
  },
  {
    title: '定时任务',
    key: 'tasks',
    width: 200,
    render: (row: any) => {
      if (!row.scheduled_tasks || row.scheduled_tasks.length === 0) {
        return h('span', { class: 'text-xs text-gray-400' }, '无')
      }

      return h('div', { class: 'space-y-1' }, [
        h('div', { class: 'text-xs font-medium' }, `${row.scheduled_tasks.length}个任务`),
        ...row.scheduled_tasks.slice(0, 2).map(task =>
          h('div', {
            class: 'text-xs text-gray-500 truncate',
            title: `${task.name} - ${task.schedule}`
          }, `${task.name}`)
        )
      ])
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    render: (row: any) => {
      // 第一排按钮 - 3个按钮
      const firstRowButtons = [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors mr-1',
          onClick: () => viewPluginDetails(row)
        }, [
          h('i', { class: 'fas fa-info-circle mr-1 text-xs' }),
          '详情'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-yellow-100 hover:bg-yellow-200 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400 rounded transition-colors mr-1',
          onClick: () => configurePlugin(row)
        }, [
          h('i', { class: 'fas fa-cog mr-1 text-xs' }),
          '配置'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-purple-100 hover:bg-purple-200 text-purple-700 dark:bg-purple-900/20 dark:text-purple-400 rounded transition-colors',
          onClick: () => viewPluginLogs(row)
        }, [
          h('i', { class: 'fas fa-file-alt mr-1 text-xs' }),
          '日志'
        ])
      ]

      // 第二排按钮 - 2个按钮
      const secondRowButtons = [
        h('button', {
          class: `px-2 py-1 text-xs ${row.enabled ? 'bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400' : 'bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400'} rounded transition-colors mr-1`,
          onClick: () => togglePlugin(row)
        }, [
          h('i', {
            class: `fas ${row.enabled ? 'fa-stop' : 'fa-play'} mr-1 text-xs`
          }),
          row.enabled ? '禁用' : '启用'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
          onClick: () => uninstallPlugin(row)
        }, [
          h('i', { class: 'fas fa-trash mr-1 text-xs' }),
          '卸载'
        ])
      ]

      return h('div', { class: 'space-y-1' }, [
        h('div', { class: 'flex items-center gap-1' }, firstRowButtons),
        h('div', { class: 'flex items-center gap-1' }, secondRowButtons)
      ])
    }
  }
]

// 搜索防抖
let searchTimeout: NodeJS.Timeout | null = null
const debounceSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    fetchPlugins()
  }, 300)
}

// 获取插件列表
const fetchPlugins = async () => {
  loading.value = true
  try {
    const response = await $fetch('/api/plugins')
    if (response.success) {
      let filteredPlugins = response.data

      // 应用筛选条件
      if (filters.value.status) {
        filteredPlugins = filteredPlugins.filter((plugin: any) => {
          if (filters.value.status === 'enabled') return plugin.enabled
          if (filters.value.status === 'disabled') return !plugin.enabled
          return true
        })
      }

      if (filters.value.search) {
        const query = filters.value.search.toLowerCase()
        filteredPlugins = filteredPlugins.filter((plugin: any) =>
          plugin.name.toLowerCase().includes(query) ||
          (plugin.display_name && plugin.display_name.toLowerCase().includes(query)) ||
          (plugin.description && plugin.description.toLowerCase().includes(query))
        )
      }

      plugins.value = filteredPlugins
    }
  } catch (error) {
    console.error('获取插件列表失败:', error)
    if (process.client) {
      notification.error({
        content: '获取插件列表失败',
        duration: 3000
      })
    }
  } finally {
    loading.value = false
  }
}

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    status: '',
    search: ''
  }
  fetchPlugins()
}

// 查看插件详情
const viewPluginDetails = (plugin: any) => {
  selectedPlugin.value = plugin
  showDetailModal.value = true
}

// 配置插件
const configurePlugin = async (plugin: any) => {
  configPlugin.value = plugin
  try {
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = plugin.name.replace('.plugin', '')
    const response = await $fetch(`/api/plugins/${pluginName}`)
    if (response.success && response.data.plugin) {
      // 使用插件详细信息，包括配置字段定义
      configPlugin.value = response.data.plugin

      // 如果配置为空，根据config_fields生成默认配置
      let config = response.data.plugin.config || {}
      if (!config || Object.keys(config).length === 0) {
        config = generateDefaultConfig(response.data.plugin.config_fields || {})
      }
      pluginConfig.value = config
    } else {
      // 使用传入的插件基本信息
      pluginConfig.value = {}
    }
  } catch (error) {
    console.error('加载插件配置失败:', error)
    // 即使API失败，也使用传入的插件基本信息
    pluginConfig.value = {}
  }
  showConfigModal.value = true
}

// 生成默认配置
const generateDefaultConfig = (configFields: any) => {
  const defaultConfig: any = {}

  for (const [fieldName, fieldConfig] of Object.entries(configFields)) {
    const field = fieldConfig as any
    if (field.default !== undefined && field.default !== null) {
      defaultConfig[fieldName] = field.default
    } else {
      // 根据字段类型设置默认值
      switch (field.type) {
        case 'boolean':
          defaultConfig[fieldName] = false
          break
        case 'number':
          defaultConfig[fieldName] = 0
          break
        case 'string':
        case 'text':
          defaultConfig[fieldName] = ''
          break
        case 'select':
          defaultConfig[fieldName] = field.options && field.options.length > 0 ? field.options[0] : ''
          break
        default:
          defaultConfig[fieldName] = null
      }
    }
  }

  return defaultConfig
}

// 处理配置保存
const handleConfigSave = async (config: any) => {
  if (!configPlugin.value) return

  try {
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = configPlugin.value.name.replace('.plugin', '')
    const response = await $fetch(`/api/plugins/${pluginName}/config`, {
      method: 'PUT',
      body: {
        config: config
      }
    })
    if (response.success) {
      if (process.client) {
        notification.success({
          content: '配置已保存',
          duration: 3000
        })
      }
      // 更新当前配置值
      pluginConfig.value = config
      showConfigModal.value = false
      await fetchPlugins()
    } else {
      throw new Error(response.error || '保存失败')
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    if (process.client) {
      notification.error({
        content: '保存配置失败: ' + error.message,
        duration: 3000
      })
    }
  } finally {
    saving.value = false
  }
}

// 处理配置重置
const handleConfigReset = () => {
  if (configPlugin.value && configPlugin.value.config_fields) {
    pluginConfig.value = generateDefaultConfig(configPlugin.value.config_fields)
    if (process.client) {
      notification.info({
        content: '配置已重置为默认值',
        duration: 2000
      })
    }
  }
}

// 重置配置为默认值
const resetConfigToDefault = () => {
  handleConfigReset()
}

// 关闭配置模态框
const closeConfigModal = () => {
  showConfigModal.value = false
}

// 查看插件日志
const viewPluginLogs = async (plugin: any) => {
  logsPlugin.value = plugin
  await loadPluginLogs()
  showLogsModal.value = true
}

const loadPluginLogs = async () => {
  if (!logsPlugin.value) return

  try {
    loadingLogs.value = true
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = logsPlugin.value.name.replace('.plugin', '')
    const response = await $fetch(`/api/plugins/${pluginName}/logs?limit=50`)
    if (response.success) {
      pluginLogs.value = response.data?.logs || []
    }
  } catch (error) {
    console.error('加载插件日志失败:', error)
    pluginLogs.value = []
  } finally {
    loadingLogs.value = false
  }
}

const refreshLogs = async () => {
  await loadPluginLogs()
}

// 切换插件状态
const togglePlugin = async (plugin: any) => {
  try {
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = plugin.name.replace('.plugin', '')
    const action = plugin.enabled ? 'disable' : 'enable'
    const response = await $fetch(`/api/plugins/${pluginName}/${action}`, {
      method: 'POST'
    })
    if (response.success) {
      if (process.client) {
        notification.success({
          content: `插件已${plugin.enabled ? '禁用' : '启用'}`,
          duration: 3000
        })
      }
      await fetchPlugins()
    }
  } catch (error) {
    console.error(`${plugin.enabled ? '禁用' : '启用'}插件失败:`, error)
    if (process.client) {
      notification.error({
        content: `${plugin.enabled ? '禁用' : '启用'}插件失败`,
        duration: 3000
      })
    }
  }
}


const formatTime = (timestamp: string) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString()
}

// 安装插件相关方法
const handleFileChange = (fileList: any[]) => {
  installFiles.value = fileList
}

const closeInstallModal = () => {
  showInstallModal.value = false
  installUrl.value = ''
  installFiles.value = []
  installProgress.value = 0
  installStatus.value = ''
  installResult.value = ''
  installSuccess.value = false
  installType.value = 'url'
}

const installPlugin = async () => {
  if (!canInstall.value) return

  installing.value = true
  installProgress.value = 0
  installStatus.value = '准备安装...'
  installResult.value = ''

  try {
    let requestBody: any = {}

    if (installType.value === 'url') {
      requestBody.source = installUrl.value.trim()
      installStatus.value = '正在从URL下载插件包...'
    } else {
      // 文件上传
      const file = installFiles.value[0]?.file
      if (!file) {
        throw new Error('请选择要安装的文件')
      }

      // 创建FormData
      const formData = new FormData()
      formData.append('file', file)
      requestBody = formData
      installStatus.value = '正在上传插件文件...'
    }

    installProgress.value = 30

    const response = await $fetch('/api/plugins/install', {
      method: 'POST',
      headers: installType.value === 'url' ? {
        'Content-Type': 'application/json'
      } : {},
      body: requestBody
    })

    installProgress.value = 80
    installStatus.value = '正在完成安装...'

    if (response.success) {
      installProgress.value = 100
      installStatus.value = '安装完成'
      installResult.value = response.message || '插件安装成功'
      installSuccess.value = true

      if (process.client) {
        notification.success({
          content: '插件安装成功',
          duration: 3000
        })
      }

      // 刷新插件列表
      await fetchPlugins()

      // 延迟关闭模态框
      setTimeout(() => {
        closeInstallModal()
      }, 2000)
    } else {
      throw new Error(response.error || '安装失败')
    }
  } catch (error) {
    installProgress.value = 0
    installStatus.value = '安装失败'
    installResult.value = error.message || '安装过程中发生错误'
    installSuccess.value = false

    if (process.client) {
      notification.error({
        content: '插件安装失败: ' + error.message,
        duration: 3000
      })
    }
  } finally {
    installing.value = false
  }
}

// 卸载插件
const uninstallPlugin = async (plugin: any) => {
  if (!confirm(`确定要卸载插件 "${plugin.display_name || plugin.name}" 吗？此操作不可撤销。`)) {
    return
  }

  try {
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = plugin.name.replace('.plugin', '')

    const response = await $fetch(`/api/plugins/${pluginName}`, {
      method: 'DELETE'
    })

    if (response.success) {
      if (process.client) {
        notification.success({
          content: '插件卸载成功',
          duration: 3000
        })
      }

      // 刷新插件列表
      await fetchPlugins()
    } else {
      throw new Error(response.error || '卸载失败')
    }
  } catch (error) {
    console.error('卸载插件失败:', error)
    if (process.client) {
      notification.error({
        content: '卸载插件失败: ' + error.message,
        duration: 3000
      })
    }
  }
}

// 初始化数据
onMounted(() => {
  fetchPlugins()
})
</script>

<style scoped>
.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
</style>