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
          <!-- 空白区域用于按钮 -->
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
          :scroll-x="1200"
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
  <n-modal v-model:show="showConfigModal" :mask-closable="false" preset="card" :style="{ maxWidth: '700px', width: '90%' }" title="插件配置">
    <div v-if="configPlugin" class="space-y-4">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ configPlugin.display_name || configPlugin.name }} 配置</h3>
        <div class="flex items-center space-x-2">
          <n-button @click="resetConfig" size="small" type="warning">重置</n-button>
          <n-button @click="saveConfig" :loading="saving" type="primary">保存</n-button>
        </div>
      </div>

      <div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4">
        <n-code :code="JSON.stringify(pluginConfig, null, 2)" language="json" show-line-numbers />
      </div>

      <div class="text-sm text-gray-600 dark:text-gray-400">
        <i class="fas fa-info-circle mr-1"></i>
        请确保配置格式正确，保存后将立即生效
      </div>
    </div>
  </n-modal>

  <!-- 插件日志模态框 -->
  <n-modal v-model:show="showLogsModal" :mask-closable="false" preset="card" :style="{ maxWidth: '800px', width: '90%' }" title="插件日志">
    <div v-if="logsPlugin" class="space-y-4">
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

      <div v-if="loadingLogs" class="flex items-center justify-center py-8">
        <n-spin size="medium" />
      </div>

      <div v-else-if="pluginLogs.length === 0" class="text-center py-8">
        <n-empty description="暂无日志" />
      </div>

      <div v-else class="space-y-2 max-h-96 overflow-y-auto">
        <div
          v-for="log in pluginLogs"
          :key="log.id"
          class="p-3 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700"
        >
          <div class="flex items-center justify-between mb-1">
            <span class="text-sm font-medium text-gray-900 dark:text-white">
              {{ log.level }}
            </span>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              {{ formatTime(log.created_at) }}
            </span>
          </div>
          <p class="text-sm text-gray-700 dark:text-gray-300">{{ log.message }}</p>
        </div>
      </div>
    </div>
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
const message = useMessage()
const notification = useNotification()

const loading = ref(false)
const plugins = ref<any[]>([])
const showDetailModal = ref(false)
const showConfigModal = ref(false)
const showLogsModal = ref(false)
const selectedPlugin = ref<any>(null)
const configPlugin = ref<any>(null)
const logsPlugin = ref<any>(null)
const pluginConfig = ref({})
const pluginLogs = ref([])
const loadingLogs = ref(false)
const saving = ref(false)

// 分页和筛选状态
const filters = ref({
  status: '',
  search: ''
})

// 表格列定义
const columns = [
  {
    title: '插件信息',
    key: 'name',
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        // 第一行：名称和版本
        h('div', { class: 'flex items-center gap-2' }, [
          h('i', { class: 'fas fa-plug text-blue-500 text-sm' }),
          h('span', { class: 'font-medium text-sm' }, row.display_name || row.name),
          h('span', { class: 'text-xs text-gray-400' }, `v${row.version}`)
        ]),
        // 第二行：描述
        h('div', {
          class: 'text-xs text-gray-500 dark:text-gray-400 line-clamp-2 max-w-[330px]'
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
    title: '状态',
    key: 'enabled',
    width: 80,
    render: (row: any) => {
      return h('n-tag', {
        type: row.enabled ? 'success' : 'error',
        size: 'small',
        bordered: false
      }, { default: () => row.enabled ? '已启用' : '已禁用' })
    }
  },
  {
    title: '定时任务',
    key: 'tasks',
    width: 100,
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
    width: 120,
    render: (row: any) => {
      // 第一排按钮
      const firstRowButtons = [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors mr-1',
          onClick: () => viewPluginDetails(row)
        }, [
          h('i', { class: 'fas fa-info-circle mr-1 text-xs' }),
          '详情'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-yellow-100 hover:bg-yellow-200 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400 rounded transition-colors',
          onClick: () => configurePlugin(row)
        }, [
          h('i', { class: 'fas fa-cog mr-1 text-xs' }),
          '配置'
        ])
      ]

      // 第二排按钮
      const secondRowButtons = [
        h('button', {
          class: 'px-2 py-1 text-xs bg-purple-100 hover:bg-purple-200 text-purple-700 dark:bg-purple-900/20 dark:text-purple-400 rounded transition-colors mr-1',
          onClick: () => viewPluginLogs(row)
        }, [
          h('i', { class: 'fas fa-file-alt mr-1 text-xs' }),
          '日志'
        ]),
        h('button', {
          class: `px-2 py-1 text-xs ${row.enabled ? 'bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400' : 'bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400'} rounded transition-colors`,
          onClick: () => togglePlugin(row)
        }, [
          h('i', {
            class: `fas ${row.enabled ? 'fa-stop' : 'fa-play'} mr-1 text-xs`
          }),
          row.enabled ? '禁用' : '启用'
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
    if (response.success) {
      pluginConfig.value = response.data.config || {}
    }
  } catch (error) {
    console.error('加载插件配置失败:', error)
    pluginConfig.value = {}
  }
  showConfigModal.value = true
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
      pluginLogs.value = response.data || []
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

// 保存配置
const saveConfig = async () => {
  if (!configPlugin.value) return

  try {
    saving.value = true
    // 移除.plugin后缀，因为后端API期望不带后缀的名称
    const pluginName = configPlugin.value.name.replace('.plugin', '')
    const response = await $fetch(`/api/plugins/${pluginName}/config`, {
      method: 'PUT',
      body: {
        config: pluginConfig.value
      }
    })
    if (response.success) {
      if (process.client) {
        notification.success({
          content: '配置已保存',
          duration: 3000
        })
      }
      showConfigModal.value = false
      await fetchPlugins()
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    if (process.client) {
      notification.error({
        content: '保存配置失败',
        duration: 3000
      })
    }
  } finally {
    saving.value = false
  }
}

const resetConfig = () => {
  pluginConfig.value = {}
}

const formatTime = (timestamp: string) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString()
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