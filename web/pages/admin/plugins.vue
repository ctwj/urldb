<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">插件管理</h1>
        <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
          管理和监控系统插件
        </p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="refreshPlugins">
          <template #icon>
            <i class="fas fa-sync-alt"></i>
          </template>
          刷新
        </n-button>
        <n-button type="info" @click="validateDependencies">
          <template #icon>
            <i class="fas fa-check-circle"></i>
          </template>
          验证依赖
        </n-button>
      </div>
    </div>

    <!-- 插件列表 -->
    <n-card>
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">已安装插件</h2>
        <n-button type="success" @click="showInstallModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          安装插件
        </n-button>
      </div>

      <!-- 插件表格 -->
      <n-data-table
        :columns="columns"
        :data="plugins"
        :loading="loading"
        :pagination="pagination"
        remote
        @update:page="handlePageChange"
      />

      <!-- 空状态 -->
      <div v-if="!loading && plugins.length === 0" class="text-center py-12">
        <i class="fas fa-plug text-4xl text-gray-300 mb-4"></i>
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-1">暂无插件</h3>
        <p class="text-gray-500 dark:text-gray-400">当前系统中没有安装任何插件</p>
      </div>
    </n-card>

    <!-- 插件详情模态框 -->
    <n-modal v-model:show="showDetailModal" preset="card" style="width: 800px" title="插件详情">
      <template #header>
        <div class="flex items-center">
          <i class="fas fa-info-circle text-blue-500 mr-2"></i>
          <span>{{ selectedPlugin?.name }} 详情</span>
        </div>
      </template>

      <div v-if="selectedPlugin" class="space-y-6">
        <!-- 基本信息 -->
        <n-card size="small" title="基本信息">
          <n-descriptions label-placement="left" :column="2">
            <n-descriptions-item label="名称">{{ selectedPlugin.name }}</n-descriptions-item>
            <n-descriptions-item label="版本">{{ selectedPlugin.version }}</n-descriptions-item>
            <n-descriptions-item label="作者">{{ selectedPlugin.author }}</n-descriptions-item>
            <n-descriptions-item label="状态">
              <n-tag :type="getStatusType(selectedPlugin.status)">
                {{ selectedPlugin.status }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="描述" :span="2">{{ selectedPlugin.description }}</n-descriptions-item>
            <n-descriptions-item label="健康分数" v-if="selectedPlugin.health_score !== undefined">
              <n-progress
                type="line"
                :percentage="selectedPlugin.health_score"
                :color="getHealthColor(selectedPlugin.health_score)"
                :show-indicator="true"
              />
            </n-descriptions-item>
          </n-descriptions>
        </n-card>

        <!-- 依赖信息 -->
        <n-card size="small" title="依赖关系">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h4 class="font-medium mb-2">依赖项</h4>
              <div v-if="selectedPlugin.dependencies && selectedPlugin.dependencies.length > 0">
                <n-tag
                  v-for="dep in selectedPlugin.dependencies"
                  :key="dep"
                  class="mr-2 mb-2"
                  :type="dependencyStatus[dep] ? 'success' : 'error'"
                >
                  {{ dep }}
                  <template #icon>
                    <i :class="dependencyStatus[dep] ? 'fas fa-check' : 'fas fa-times'"></i>
                  </template>
                </n-tag>
              </div>
              <div v-else>
                <p class="text-gray-500">无依赖项</p>
              </div>
            </div>
            <div>
              <h4 class="font-medium mb-2">被依赖项</h4>
              <div v-if="dependents && dependents.length > 0">
                <n-tag
                  v-for="dep in dependents"
                  :key="dep"
                  class="mr-2 mb-2"
                >
                  {{ dep }}
                </n-tag>
              </div>
              <div v-else>
                <p class="text-gray-500">无被依赖项</p>
              </div>
            </div>
          </div>
        </n-card>

        <!-- 配置信息 -->
        <n-card size="small" title="配置信息">
          <div class="flex justify-between items-center mb-4">
            <h4 class="font-medium">当前配置</h4>
            <n-button size="small" @click="editConfig">
              <template #icon>
                <i class="fas fa-edit"></i>
              </template>
              编辑配置
            </n-button>
          </div>
          <n-code
            v-if="pluginConfig && Object.keys(pluginConfig).length > 0"
            :code="JSON.stringify(pluginConfig, null, 2)"
            language="json"
            show-line-numbers
          />
          <div v-else class="text-gray-500">
            暂无配置信息
          </div>
        </n-card>

        <!-- 统计信息 -->
        <n-card size="small" title="运行统计" v-if="selectedPlugin.total_executions > 0">
          <n-descriptions label-placement="left" :column="2">
            <n-descriptions-item label="总执行次数">{{ selectedPlugin.total_executions }}</n-descriptions-item>
            <n-descriptions-item label="总错误次数">{{ selectedPlugin.total_errors }}</n-descriptions-item>
            <n-descriptions-item label="总执行时间">{{ formatDuration(selectedPlugin.total_execution_time) }}</n-descriptions-item>
            <n-descriptions-item label="平均执行时间" v-if="selectedPlugin.total_executions > 0">
              {{ formatDuration(selectedPlugin.total_execution_time / selectedPlugin.total_executions) }}
            </n-descriptions-item>
            <n-descriptions-item label="重启次数">{{ selectedPlugin.restart_count }}</n-descriptions-item>
            <n-descriptions-item label="最后执行时间">{{ formatTime(selectedPlugin.last_execution_time) }}</n-descriptions-item>
          </n-descriptions>
        </n-card>
      </div>
    </n-modal>

    <!-- 插件配置编辑模态框 -->
    <n-modal v-model:show="showConfigModal" preset="card" style="width: 600px" title="编辑插件配置">
      <template #header>
        <div class="flex items-center">
          <i class="fas fa-cog text-blue-500 mr-2"></i>
          <span>编辑 {{ selectedPlugin?.name }} 配置</span>
        </div>
      </template>

      <div v-if="selectedPlugin" class="space-y-4">
        <n-form :model="configForm" label-placement="left" label-width="120">
          <n-form-item label="配置内容">
            <n-input
              v-model:value="configForm.configJson"
              type="textarea"
              :autosize="{ minRows: 10 }"
              placeholder="请输入JSON格式的配置"
            />
          </n-form-item>
        </n-form>

        <div class="flex justify-end space-x-3">
          <n-button @click="showConfigModal = false">取消</n-button>
          <n-button type="primary" @click="saveConfig" :loading="savingConfig">保存</n-button>
        </div>
      </div>
    </n-modal>

    <!-- 安装插件模态框 -->
    <n-modal v-model:show="showInstallModal" preset="card" style="width: 500px" title="安装插件">
      <template #header>
        <div class="flex items-center">
          <i class="fas fa-download text-blue-500 mr-2"></i>
          <span>安装插件</span>
        </div>
      </template>

      <div class="space-y-4">
        <n-alert type="info" title="提示">
          插件安装功能将在后续版本中实现
        </n-alert>

        <div class="flex justify-end">
          <n-button @click="showInstallModal = false">关闭</n-button>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'
import { useMessage, useDialog } from 'naive-ui'

// 获取消息和对话框实例
const $message = useMessage()
const $dialog = useDialog()

// 状态管理
const loading = ref(false)
const plugins = ref<any[]>([])
const selectedPlugin = ref<any>(null)
const pluginConfig = ref<any>({})
const dependencyStatus = ref<Record<string, boolean>>({})
const dependents = ref<string[]>([])
const showDetailModal = ref(false)
const showConfigModal = ref(false)
const showInstallModal = ref(false)
const savingConfig = ref(false)

// 配置表单
const configForm = ref({
  configJson: ''
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

// 表格列定义
const columns = [
  {
    title: '名称',
    key: 'name',
    render(row: any) {
      return h('div', [
        h('div', { class: 'font-medium' }, row.name),
        h('div', { class: 'text-sm text-gray-500' }, row.description)
      ])
    }
  },
  {
    title: '版本',
    key: 'version',
    width: 100
  },
  {
    title: '作者',
    key: 'author',
    width: 120
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row: any) {
      const type = getStatusType(row.status)
      return h(
        'n-tag',
        {
          type,
          size: 'small'
        },
        {
          default: () => row.status
        }
      )
    }
  },
  {
    title: '健康分数',
    key: 'health_score',
    width: 120,
    render(row: any) {
      if (row.health_score === undefined) return '-'
      return h(
        'n-progress',
        {
          type: 'line',
          percentage: row.health_score,
          showIndicator: false,
          height: 8,
          color: getHealthColor(row.health_score)
        }
      )
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render(row: any) {
      return h('div', { class: 'flex space-x-2' }, [
        h(
          'n-button',
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => showPluginDetail(row)
          },
          {
            default: () => '详情'
          }
        ),
        h(
          'n-dropdown',
          {
            options: getActionOptions(row),
            onSelect: (key: string) => handleAction(key, row),
            trigger: 'click'
          },
          {
            default: () => h(
              'n-button',
              {
                size: 'small',
                type: 'info',
                secondary: true
              },
              {
                default: () => '操作'
              }
            )
          }
        )
      ])
    }
  }
]

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'running':
      return 'success'
    case 'error':
      return 'error'
    case 'stopped':
      return 'warning'
    case 'initialized':
      return 'info'
    default:
      return 'default'
  }
}

// 获取健康分数颜色
const getHealthColor = (score: number) => {
  if (score >= 80) return '#10b981'
  if (score >= 60) return '#f59e0b'
  return '#ef4444'
}

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 格式化持续时间
const formatDuration = (duration: number) => {
  if (duration === undefined) return '-'
  // 如果是毫秒，转换为秒
  if (duration > 1000000000) {
    duration = duration / 1000000
  }

  if (duration < 1000) {
    return `${Math.round(duration)}ms`
  } else if (duration < 60000) {
    return `${(duration / 1000).toFixed(2)}s`
  } else {
    return `${(duration / 60000).toFixed(2)}m`
  }
}

// 获取操作选项
const getActionOptions = (row: any) => {
  const options = []

  if (row.status === 'registered' || row.status === 'stopped') {
    options.push({
      label: '初始化',
      key: 'initialize'
    })
  }

  if (row.status === 'initialized' || row.status === 'stopped') {
    options.push({
      label: '启动',
      key: 'start'
    })
  }

  if (row.status === 'running') {
    options.push({
      label: '停止',
      key: 'stop'
    })
  }

  options.push({
    label: '卸载',
    key: 'uninstall'
  })

  return options
}

// 处理分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchPlugins()
}

// 获取插件列表
const fetchPlugins = async () => {
  loading.value = true
  try {
    const response = await useApiFetch('/plugins').then(parseApiResponse)
    if (response && response.plugins) {
      plugins.value = response.plugins
      pagination.itemCount = response.plugins.length
    }
  } catch (error) {
    console.error('获取插件列表失败:', error)
    // 显示错误提示
  } finally {
    loading.value = false
  }
}

// 刷新插件列表
const refreshPlugins = () => {
  fetchPlugins()
}

// 显示插件详情
const showPluginDetail = async (plugin: any) => {
  selectedPlugin.value = plugin

  // 获取插件配置
  try {
    const configResponse = await useApiFetch(`/plugins/${plugin.name}/config`).then(parseApiResponse)
    if (configResponse && configResponse.config) {
      pluginConfig.value = configResponse.config
      configForm.value.configJson = JSON.stringify(configResponse.config, null, 2)
    } else {
      pluginConfig.value = {}
      configForm.value.configJson = '{}'
    }
  } catch (error) {
    console.error('获取插件配置失败:', error)
    pluginConfig.value = {}
    configForm.value.configJson = '{}'
  }

  // 获取依赖信息
  try {
    const depResponse = await useApiFetch(`/plugins/${plugin.name}/dependencies`).then(parseApiResponse)
    if (depResponse && depResponse.dependencies) {
      dependencyStatus.value = depResponse.dependencies
      dependents.value = depResponse.dependents || []
    }
  } catch (error) {
    console.error('获取依赖信息失败:', error)
    dependencyStatus.value = {}
    dependents.value = []
  }

  showDetailModal.value = true
}

// 编辑配置
const editConfig = () => {
  showDetailModal.value = false
  showConfigModal.value = true
}

// 保存配置
const saveConfig = async () => {
  if (!selectedPlugin.value) return

  savingConfig.value = true
  try {
    let configData
    try {
      configData = JSON.parse(configForm.value.configJson)
    } catch (e) {
      $message.error('配置格式不正确，请检查JSON格式')
      return
    }

    await useApiFetch(`/plugins/${selectedPlugin.value.name}/config`, {
      method: 'PUT',
      body: JSON.stringify(configData)
    }).then(parseApiResponse)

    $message.success('配置保存成功')
    showConfigModal.value = false

    // 刷新插件列表
    await fetchPlugins()
  } catch (error) {
    console.error('保存配置失败:', error)
    $message.error('保存配置失败: ' + (error as any).message)
  } finally {
    savingConfig.value = false
  }
}

// 处理插件操作
const handleAction = async (action: string, plugin: any) => {
  try {
    // 添加确认对话框（除初始化外）
    if (action === 'uninstall') {
      const confirmed = await new Promise((resolve) => {
        $dialog.warning({
          title: '确认卸载',
          content: `确定要卸载插件 "${plugin.name}" 吗？此操作不可逆！`,
          positiveText: '确定',
          negativeText: '取消',
          onPositiveClick: () => resolve(true),
          onNegativeClick: () => resolve(false)
        })
      })

      if (!confirmed) return
    }

    switch (action) {
      case 'initialize':
        await useApiFetch(`/plugins/${plugin.name}/initialize`, {
          method: 'POST',
          body: JSON.stringify({})
        }).then(parseApiResponse)
        $message.success('插件初始化成功')
        break

      case 'start':
        await useApiFetch(`/plugins/${plugin.name}/start`, {
          method: 'POST'
          // 空body
        }).then(parseApiResponse)
        $message.success('插件启动成功')
        break

      case 'stop':
        await useApiFetch(`/plugins/${plugin.name}/stop`, {
          method: 'POST'
          // 空body
        }).then(parseApiResponse)
        $message.success('插件停止成功')
        break

      case 'uninstall':
        const force = false
        await useApiFetch(`/plugins/${plugin.name}?force=${force}`, {
          method: 'DELETE'
        }).then(parseApiResponse)
        $message.success('插件卸载成功')
        break
    }

    // 刷新插件列表
    await fetchPlugins()
  } catch (error) {
    console.error(`${action}插件失败:`, error)
    $message.error(`${action}插件失败: ` + (error as any).message)
  }
}

// 验证依赖
const validateDependencies = async () => {
  try {
    const response = await useApiFetch('/plugins/validate-dependencies', {
      method: 'POST'
    }).then(parseApiResponse)

    if (response && response.valid) {
      $message.success('所有插件依赖验证通过')
    } else {
      $message.warning('插件依赖验证失败: ' + (response?.message || '未知错误'))
    }
  } catch (error) {
    console.error('验证依赖失败:', error)
    $message.error('验证依赖失败: ' + (error as any).message)
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchPlugins()
})
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style>