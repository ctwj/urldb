<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
        <i class="fas fa-plug mr-3 text-blue-500"></i>
        插件管理
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        管理系统插件，配置插件参数，监控插件运行状态
      </p>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 dark:bg-blue-900/20 rounded-lg">
            <i class="fas fa-plug text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600 dark:text-gray-400">总插件数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total_plugins }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 dark:bg-green-900/20 rounded-lg">
            <i class="fas fa-check-circle text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600 dark:text-gray-400">启用插件</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.enabled_plugins }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center">
          <div class="p-3 bg-red-100 dark:bg-red-900/20 rounded-lg">
            <i class="fas fa-times-circle text-red-600 dark:text-red-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600 dark:text-gray-400">禁用插件</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.disabled_plugins }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
        <div class="flex items-center">
          <div class="p-3 bg-purple-100 dark:bg-purple-900/20 rounded-lg">
            <i class="fas fa-chart-line text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600 dark:text-gray-400">成功率</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.success_rate.toFixed(1) }}%</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 操作栏 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex flex-col md:flex-row gap-4">
        <!-- 搜索框 -->
        <div class="flex-1">
          <n-input v-model:value="searchQuery" placeholder="搜索插件..." clearable>
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </n-input>
        </div>

        <!-- 状态过滤 -->
        <n-select v-model:value="statusFilter" placeholder="全部状态" :options="statusOptions" clearable class="w-full md:w-40" />

        <!-- 分类过滤 -->
        <n-select v-model:value="categoryFilter" placeholder="全部分类" :options="categoryOptions" clearable class="w-full md:w-40" />

        <!-- 操作按钮 -->
        <div class="flex items-center space-x-2">
          <n-button @click="refreshPlugins" :loading="loading" type="info">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
          <n-button @click="openPluginMarket" type="primary">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            插件市场
          </n-button>
        </div>
      </div>
    </div>

    <!-- 插件列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large">
          <template #description>
            <span class="text-gray-500 dark:text-gray-400">加载中...</span>
          </template>
        </n-spin>
      </div>

      <div v-else-if="filteredPlugins.length === 0" class="flex flex-col items-center justify-center py-12">
        <n-empty description="暂无插件">
          <template #icon>
            <i class="fas fa-plug text-4xl text-gray-400"></i>
          </template>
          <template #extra>
            <n-button @click="openPluginMarket" type="primary">
              <template #icon>
                <i class="fas fa-plus"></i>
              </template>
              添加插件
            </n-button>
          </template>
        </n-empty>
      </div>

      <!-- 插件列表和分页 -->
      <div v-else class="flex flex-col">
        <div class="divide-y divide-gray-200 dark:divide-gray-700">
          <div
            v-for="plugin in paginatedPlugins"
            :key="plugin.id"
            class="p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            <div class="flex items-center justify-between">
              <!-- 左侧信息 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-4">
                  <!-- 插件图标 -->
                  <div class="flex-shrink-0">
                    <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
                      <i class="fas fa-plug text-blue-600 dark:text-blue-400 text-lg"></i>
                    </div>
                  </div>

                  <!-- 插件信息 -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center space-x-3">
                      <h3 class="text-lg font-medium text-gray-900 dark:text-white truncate">
                        {{ plugin.name }}
                      </h3>
                      <n-tag :type="plugin.enabled ? 'success' : 'error'" size="small">
                        {{ plugin.enabled ? '运行中' : '已禁用' }}
                      </n-tag>
                      <span class="text-sm text-gray-500 dark:text-gray-400">
                        v{{ plugin.version }}
                      </span>
                    </div>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mt-1 line-clamp-1">
                      {{ plugin.description }}
                    </p>
                    <div class="flex items-center space-x-4 mt-2 text-xs text-gray-500 dark:text-gray-400">
                      <span><i class="fas fa-user mr-1"></i>{{ plugin.author }}</span>
                      <span><i class="fas fa-code-branch mr-1"></i>{{ plugin.category || 'utility' }}</span>
                      <span v-if="plugin.executions">
                        <i class="fas fa-play mr-1"></i>{{ plugin.executions }}次执行
                      </span>
                      <span v-if="plugin.avg_duration">
                        <i class="fas fa-clock mr-1"></i>{{ plugin.avg_duration }}ms
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 右侧操作按钮 -->
              <div class="flex items-center space-x-2 ml-4">
                <n-button size="small" @click="viewPluginDetails(plugin)" type="info" text>
                  <template #icon>
                    <i class="fas fa-info-circle"></i>
                  </template>
                  详情
                </n-button>
                <n-button size="small" @click="configurePlugin(plugin)" type="warning" text>
                  <template #icon>
                    <i class="fas fa-cog"></i>
                  </template>
                  配置
                </n-button>
                <n-button size="small" @click="viewPluginLogs(plugin)" type="tertiary" text>
                  <template #icon>
                    <i class="fas fa-file-alt"></i>
                  </template>
                  日志
                </n-button>
                <n-button
                  size="small"
                  :type="plugin.enabled ? 'error' : 'success'"
                  @click="togglePlugin(plugin)"
                  text
                >
                  <template #icon>
                    <i :class="plugin.enabled ? 'fas fa-stop' : 'fas fa-play'"></i>
                  </template>
                  {{ plugin.enabled ? '禁用' : '启用' }}
                </n-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="totalPages > 1" class="p-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700 dark:text-gray-400">
              显示 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, filteredPlugins.length) }} 个，共 {{ filteredPlugins.length }} 个插件
            </div>
            <n-pagination
              v-model:page="currentPage"
              v-model:page-size="pageSize"
              :item-count="filteredPlugins.length"
              :page-sizes="[10, 20, 50, 100]"
              show-size-picker
              @update:page="currentPage = $event"
              @update:page-size="pageSize = $event; currentPage = 1"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 插件详情模态框 -->
    <n-modal v-model:show="showDetailModal" preset="card" title="插件详情" style="width: 600px">
      <div v-if="selectedPlugin" class="space-y-4">
        <div class="flex items-center space-x-4">
          <div class="w-16 h-16 bg-blue-100 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
            <i class="fas fa-plug text-blue-600 dark:text-blue-400 text-2xl"></i>
          </div>
          <div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">{{ selectedPlugin.name }}</h3>
            <p class="text-gray-600 dark:text-gray-400">{{ selectedPlugin.description }}</p>
            <div class="flex items-center space-x-2 mt-1">
              <n-tag :type="selectedPlugin.enabled ? 'success' : 'error'" size="small">
                {{ selectedPlugin.enabled ? '运行中' : '已禁用' }}
              </n-tag>
              <span class="text-sm text-gray-500">v{{ selectedPlugin.version }}</span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">作者</label>
            <p class="text-gray-900 dark:text-white">{{ selectedPlugin.author }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">分类</label>
            <p class="text-gray-900 dark:text-white">{{ selectedPlugin.category || 'utility' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">执行次数</label>
            <p class="text-gray-900 dark:text-white">{{ selectedPlugin.executions || 0 }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">平均耗时</label>
            <p class="text-gray-900 dark:text-white">{{ selectedPlugin.avg_duration || 0 }}ms</p>
          </div>
        </div>

        <div v-if="selectedPlugin.permissions">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">权限</label>
          <div class="flex flex-wrap gap-2">
            <n-tag v-for="permission in selectedPlugin.permissions" :key="permission" size="small">
              {{ permission }}
            </n-tag>
          </div>
        </div>

        <div v-if="selectedPlugin.dependencies">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">依赖</label>
          <div class="flex flex-wrap gap-2">
            <n-tag v-for="dep in selectedPlugin.dependencies" :key="dep" size="small" type="warning">
              {{ dep }}
            </n-tag>
          </div>
        </div>
      </div>
    </n-modal>

    <!-- 插件配置模态框 -->
    <n-modal v-model:show="showConfigModal" preset="card" title="插件配置" style="width: 700px">
      <div v-if="configPlugin" class="space-y-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ configPlugin.name }} 配置</h3>
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
    <n-modal v-model:show="showLogsModal" preset="card" title="插件日志" style="width: 800px">
      <div v-if="logsPlugin" class="space-y-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ logsPlugin.name }} 日志</h3>
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
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  ssr: false
})

const notification = useNotification()

// 响应式数据
const loading = ref(false)
const plugins = ref([])
const stats = ref({
  total_plugins: 0,
  enabled_plugins: 0,
  disabled_plugins: 0,
  success_rate: 0
})

// 过滤和搜索
const searchQuery = ref('')
const statusFilter = ref('')
const categoryFilter = ref('')

// 分页
const currentPage = ref(1)
const pageSize = ref(10)

// 模态框状态
const showDetailModal = ref(false)
const showConfigModal = ref(false)
const showLogsModal = ref(false)
const selectedPlugin = ref(null)
const configPlugin = ref(null)
const logsPlugin = ref(null)
const pluginConfig = ref({})
const pluginLogs = ref([])
const loadingLogs = ref(false)
const saving = ref(false)

// 过滤选项
const statusOptions = [
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' }
]

const categoryOptions = [
  { label: '工具', value: 'utility' },
  { label: '分析', value: 'analytics' },
  { label: '通知', value: 'notification' }
]

// 计算属性
const filteredPlugins = computed(() => {
  let filtered = plugins.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(plugin =>
      plugin.name.toLowerCase().includes(query) ||
      plugin.description.toLowerCase().includes(query) ||
      plugin.author.toLowerCase().includes(query)
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter(plugin => {
      if (statusFilter.value === 'enabled') return plugin.enabled
      if (statusFilter.value === 'disabled') return !plugin.enabled
      return true
    })
  }

  // 分类过滤
  if (categoryFilter.value) {
    filtered = filtered.filter(plugin => plugin.category === categoryFilter.value)
  }

  return filtered
})

const totalPages = computed(() => {
  return Math.ceil(filteredPlugins.value.length / pageSize.value)
})

const paginatedPlugins = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredPlugins.value.slice(start, end)
})

// 方法
const loadPlugins = async () => {
  try {
    loading.value = true
    const response = await $fetch('/api/plugins')
    if (response.success) {
      plugins.value = response.data
    }
  } catch (error) {
    console.error('加载插件列表失败:', error)
    notification.error({
      title: '失败',
      content: '加载插件列表失败: ' + (error.message || '未知错误'),
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await $fetch('/api/plugins/stats')
    if (response.success) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

const refreshPlugins = async () => {
  await Promise.all([loadPlugins(), loadStats()])
}

const togglePlugin = async (plugin) => {
  try {
    // 移除.pb后缀，因为后端API期望不带后缀的名称
    const pluginName = plugin.name.replace('.pb', '')
    const action = plugin.enabled ? 'disable' : 'enable'
    const response = await $fetch(`/api/plugins/${pluginName}/${action}`, {
      method: 'POST'
    })
    if (response.success) {
      notification.success({
        title: '成功',
        content: `插件已${plugin.enabled ? '禁用' : '启用'}`,
        duration: 3000
      })
      await loadPlugins()
      await loadStats()
    }
  } catch (error) {
    console.error(`${plugin.enabled ? '禁用' : '启用'}插件失败:`, error)
    notification.error({
      title: '失败',
      content: `${plugin.enabled ? '禁用' : '启用'}插件失败: ` + (error.message || '未知错误'),
      duration: 3000
    })
  }
}

const viewPluginDetails = (plugin) => {
  selectedPlugin.value = plugin
  showDetailModal.value = true
}

const configurePlugin = async (plugin) => {
  configPlugin.value = plugin
  // 加载插件配置
  try {
    // 移除.pb后缀，因为后端API期望不带后缀的名称
    const pluginName = plugin.name.replace('.pb', '')
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

const viewPluginLogs = async (plugin) => {
  logsPlugin.value = plugin
  await loadPluginLogs()
  showLogsModal.value = true
}

const loadPluginLogs = async () => {
  if (!logsPlugin.value) return

  try {
    loadingLogs.value = true
    // 移除.pb后缀，因为后端API期望不带后缀的名称
    const pluginName = logsPlugin.value.name.replace('.pb', '')
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

const saveConfig = async () => {
  if (!configPlugin.value) return

  try {
    saving.value = true
    // 移除.pb后缀，因为后端API期望不带后缀的名称
    const pluginName = configPlugin.value.name.replace('.pb', '')
    const response = await $fetch(`/api/plugins/${pluginName}/config`, {
      method: 'PUT',
      body: {
        config: pluginConfig.value
      }
    })
    if (response.success) {
      notification.success({
        title: '成功',
        content: '配置已保存',
        duration: 3000
      })
      showConfigModal.value = false
      await loadPlugins()
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    notification.error({
      title: '失败',
      content: '保存配置失败: ' + (error.message || '未知错误'),
      duration: 3000
    })
  } finally {
    saving.value = false
  }
}

const resetConfig = () => {
  pluginConfig.value = {}
}

const openPluginMarket = () => {
  notification.info({
    title: '提示',
    content: '插件市场功能开发中...',
    duration: 3000
  })
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString()
}

// 生命周期
onMounted(() => {
  refreshPlugins()
})

// 监听过滤器变化，重置分页
watch([searchQuery, statusFilter, categoryFilter], () => {
  currentPage.value = 1
})
</script>

<style scoped>
.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}
</style>