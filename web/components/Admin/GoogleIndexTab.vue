<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <!-- Google索引配置 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Google索引配置(功能测试中)</h3>
        <p class="text-gray-600 dark:text-gray-400">配置Google Search Console API和索引相关设置</p>
      </div>

      <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">Google索引功能</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              开启后系统将在生成sitemap后自动提交Sitemap到Google索引
            </p>
          </div>
          <n-switch
            :value="googleIndexConfig.enabled"
            @update:value="updateGoogleIndexConfig"
            :loading="configLoading"
            size="large"
          >
            <template #checked>已开启</template>
            <template #unchecked>已关闭</template>
          </n-switch>
        </div>

        <!-- 配置详情 -->
        <div class="border-t border-gray-200 dark:border-gray-600 pt-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">站点URL</label>
              <n-input
                :value="getSiteUrlDisplay"
                :disabled="true"
                placeholder="请先在站点配置中设置站点URL"
              >
                <template #prefix>
                  <i class="fas fa-globe text-gray-400"></i>
                </template>
              </n-input>
              <!-- 所有权验证按钮 -->
              <div class="mt-3">
                <n-button
                  type="info"
                  size="small"
                  ghost
                  @click="$emit('show-verification')"
                >
                  <template #icon>
                    <i class="fas fa-shield-alt"></i>
                  </template>
                  所有权验证
                </n-button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">凭据文件路径</label>
              <div class="flex flex-col space-y-2">
                <n-input
                  :value="credentialsFilePath"
                  placeholder="点击上传按钮选择文件"
                  :disabled="true"
                />
                <div class="flex space-x-2">
                  <!-- 申请凭据按钮 -->
                  <n-button
                    size="small"
                    type="warning"
                    ghost
                    @click="$emit('show-credentials-guide')"
                  >
                    <template #icon>
                      <i class="fas fa-question-circle"></i>
                    </template>
                    申请凭据
                  </n-button>
                  <!-- 上传按钮 -->
                  <n-button
                    size="small"
                    type="primary"
                    ghost
                    @click="$emit('select-credentials-file')"
                  >
                    <template #icon>
                      <i class="fas fa-upload"></i>
                    </template>
                    上传凭据
                  </n-button>
                  <!-- 验证按钮 -->
                  <n-button
                    size="small"
                    type="info"
                    ghost
                    @click="validateCredentials"
                    :loading="validatingCredentials"
                    :disabled="!credentialsFilePath"
                  >
                    <template #icon>
                      <i class="fas fa-check-circle"></i>
                    </template>
                    验证凭据
                  </n-button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 凭据状态 -->
        <div v-if="credentialsStatus" class="mt-4 p-3 rounded-lg border"
          :class="{
            'bg-green-50 border-green-200 text-green-700 dark:bg-green-900/20 dark:border-green-800 dark:text-green-300': credentialsStatus === 'valid',
            'bg-yellow-50 border-yellow-200 text-yellow-700 dark:bg-yellow-900/20 dark:border-yellow-800 dark:text-yellow-300': credentialsStatus === 'invalid',
            'bg-blue-50 border-blue-200 text-blue-700 dark:bg-blue-900/20 dark:border-blue-800 dark:text-blue-300': credentialsStatus === 'verifying'
          }"
        >
          <div class="flex items-center">
            <i
              :class="{
                'fas fa-check-circle text-green-500 dark:text-green-400': credentialsStatus === 'valid',
                'fas fa-exclamation-circle text-yellow-500 dark:text-yellow-400': credentialsStatus === 'invalid',
                'fas fa-spinner fa-spin text-blue-500 dark:text-blue-400': credentialsStatus === 'verifying'
              }"
              class="mr-2"
            ></i>
            <span>{{ credentialsStatusMessage }}</span>
          </div>
        </div>
      </div>

      
      <!-- 外部工具链接 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <i class="fas fa-chart-line text-purple-600 dark:text-purple-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">Google Search Console</p>
                <p class="text-sm font-medium text-gray-900 dark:text-white">查看详细分析数据</p>
              </div>
            </div>
            <a
              :href="getSearchConsoleUrl"
              target="_blank"
              class="px-3 py-1 bg-purple-600 text-white rounded-md hover:bg-purple-700 transition-colors text-sm"
            >
              打开控制台
            </a>
          </div>
        </div>

        <div class="bg-orange-50 dark:bg-orange-900/20 rounded-lg p-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="p-2 bg-orange-100 dark:bg-orange-900 rounded-lg">
                <i class="fas fa-chart-line text-orange-600 dark:text-orange-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">Google Analytics</p>
                <p class="text-sm font-medium text-gray-900 dark:text-white">网站流量分析仪表板</p>
              </div>
            </div>
            <a
              :href="getAnalyticsUrl"
              target="_blank"
              class="px-3 py-1 bg-orange-600 text-white rounded-md hover:bg-orange-700 transition-colors text-sm"
            >
              查看分析
            </a>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex flex-wrap gap-3 mb-6">
        <n-button
          type="primary"
          @click="$emit('manual-check-urls')"
          :loading="manualCheckLoading"
          size="large"
        >
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          手动检查URL
        </n-button>

        <n-button
          type="warning"
          @click="$emit('manual-submit-urls')"
          :loading="manualSubmitLoading"
          size="large"
        >
          <template #icon>
            <i class="fas fa-paper-plane"></i>
          </template>
          手动提交URL
        </n-button>

        <n-button
          type="success"
          @click="submitSitemap"
          :loading="submitSitemapLoading"
          size="large"
        >
          <template #icon>
            <i class="fas fa-upload"></i>
          </template>
          提交网站地图
        </n-button>

        <n-button
          type="error"
          @click="$emit('diagnose-permissions')"
          :loading="diagnoseLoading"
          size="large"
        >
          <template #icon>
            <i class="fas fa-stethoscope"></i>
          </template>
          权限诊断
        </n-button>

        <n-button
          type="info"
          @click="$emit('refresh-status')"
          size="large"
        >
          <template #icon>
            <i class="fas fa-sync-alt"></i>
          </template>
          刷新状态
        </n-button>
      </div>

      <!-- 任务列表 -->
      <div>
        <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4">索引任务列表</h4>
        <n-data-table
          :columns="taskColumns"
          :data="tasks"
          :pagination="pagination"
          :loading="tasksLoading"
          :bordered="false"
          striped
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { useApi } from '~/composables/useApi'
import { ref, computed, h, watch, onMounted } from 'vue'

// Props
interface Props {
  systemConfig?: any
  googleIndexConfig: any
  tasks: any[]
  credentialsStatus: string | null
  credentialsStatusMessage: string
  configLoading: boolean
  manualCheckLoading: boolean
  manualSubmitLoading: boolean
  submitSitemapLoading: boolean
  tasksLoading: boolean
  diagnoseLoading: boolean
  pagination: any
}

const props = withDefaults(defineProps<Props>(), {
  systemConfig: null,
  googleIndexConfig: () => ({}),
  tasks: () => [],
  credentialsStatus: null,
  credentialsStatusMessage: '',
  configLoading: false,
  manualCheckLoading: false,
  manualSubmitLoading: false,
  submitSitemapLoading: false,
  tasksLoading: false,
  diagnoseLoading: false,
  pagination: () => ({})
})

// Emits
const emit = defineEmits<{
  'update:google-index-config': [value: boolean]
  'show-verification': []
  'show-credentials-guide': []
  'select-credentials-file': []
  'manual-check-urls': []
  'manual-submit-urls': []
  'refresh-status': []
  'diagnose-permissions': []
}>()

// 获取消息组件
const message = useMessage()

// 本地状态
const validatingCredentials = ref(false)

// 计算属性，用于安全地访问凭据文件路径
const credentialsFilePath = computed(() => {
  const path = props.googleIndexConfig?.credentialsFile || ''
  console.log('Component computed credentialsFilePath:', path)
  return path
})


// 任务表格列
const taskColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '标题',
    key: 'name',
    width: 200
  },
  {
    title: '类型',
    key: 'type',
    width: 120,
    render: (row: any) => {
      const typeMap = {
        status_check: { text: '状态检查', class: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' },
        sitemap_submit: { text: '网站地图', class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
        url_indexing: { text: 'URL索引', class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' }
      }
      const type = typeMap[row.type as keyof typeof typeMap] || { text: row.type, class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' }
      return h('span', {
        class: `px-2 py-1 text-xs font-medium rounded ${type.class}`
      }, type.text)
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const statusMap = {
        pending: { text: '待处理', class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' },
        running: { text: '运行中', class: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' },
        completed: { text: '完成', class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
        failed: { text: '失败', class: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' }
      }
      const status = statusMap[row.status as keyof typeof statusMap] || { text: row.status, class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' }
      return h('span', {
        class: `px-2 py-1 text-xs font-medium rounded ${status.class}`
      }, status.text)
    }
  },
  {
    title: '总项目',
    key: 'totalItems',
    width: 100
  },
  {
    title: '成功/失败',
    key: 'progress',
    width: 120,
    render: (row: any) => {
      return h('span', `${row.successful_items} / ${row.failed_items}`)
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 150,
    render: (row: any) => {
      return row.created_at ? new Date(row.created_at).toLocaleString('zh-CN') : 'N/A'
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: any) => {
      return h('div', { class: 'space-x-2' }, [
        h('button', {
          class: 'text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 text-sm',
          onClick: () => emit('view-task-items', row.id)
        }, '详情'),
        h('button', {
          class: 'text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300 text-sm',
          disabled: row.status !== 'pending' && row.status !== 'running',
          onClick: () => emit('start-task', row.id)
        }, '启动')
      ].filter(btn => !btn.props?.disabled))
    }
  }
]

// 验证凭据
const validateCredentials = async () => {
  if (!credentialsFilePath.value) {
    message.warning('请先上传凭据文件')
    return
  }

  validatingCredentials.value = true

  try {
    const api = useApi()
    const response = await api.googleIndexApi.validateCredentials({})

    if (response?.valid) {
      message.success('凭据验证成功')
      emit('update:google-index-config')
    } else {
      message.error('凭据验证失败：' + (response?.message || '凭据无效或权限不足'))
    }
  } catch (error: any) {
    console.error('凭据验证失败:', error)
    message.error('凭据验证失败: ' + (error?.message || '网络错误'))
  } finally {
    validatingCredentials.value = false
  }
}

// 更新Google索引配置
const updateGoogleIndexConfig = async (value: boolean) => {
  emit('update:google-index-config', value)
}

// 提交网站地图
const submitSitemap = async () => {
  const siteUrl = getSiteUrl.value
  if (!siteUrl || siteUrl === 'https://example.com') {
    message.warning('请先在站点配置中设置正确的站点URL')
    return
  }

  const sitemapUrl = siteUrl.endsWith('/') ? siteUrl + 'sitemap.xml' : siteUrl + '/sitemap.xml'

  try {
    const api = useApi()
    const response = await api.googleIndexApi.createGoogleIndexTask({
      title: `网站地图提交任务 - ${new Date().toLocaleString('zh-CN')}`,
      type: 'sitemap_submit',
      description: `提交网站地图: ${sitemapUrl}`,
      SitemapURL: sitemapUrl
    })
    if (response) {
      message.success('网站地图提交任务已创建')
      emit('refresh-status')
    }
  } catch (error: any) {
    console.error('提交网站地图失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '提交网站地图失败'
    message.error('提交网站地图失败: ' + errorMsg)
  }
}

// 从 store 获取站点 URL（作为备用方案）
const siteUrlFromStore = ref('')

// 在组件挂载时从 store 获取配置
onMounted(async () => {
  try {
    const { useSystemConfigStore } = await import('~/stores/systemConfig')
    const systemConfigStore = useSystemConfigStore()
    await systemConfigStore.initConfig()
    siteUrlFromStore.value = systemConfigStore.config?.site_url || ''
  } catch (error) {
    console.error('从 store 获取站点 URL 失败:', error)
  }
})

// 获取站点 URL（优先使用 prop，如果为空则使用 store）
const getSiteUrl = computed(() => {
  return props.systemConfig?.site_url || siteUrlFromStore.value || ''
})

// 获取Google Search Console URL
const getSearchConsoleUrl = computed(() => {
  const siteUrl = getSiteUrl.value
  if (!siteUrl) {
    return 'https://search.google.com/search-console'
  }

  // 格式化URL用于Google Search Console
  const normalizedUrl = siteUrl.startsWith('http') ? siteUrl : `https://${siteUrl}`
  return `https://search.google.com/search-console/performance/search-analytics?resource_id=${encodeURIComponent(normalizedUrl)}`
})

// 获取Google Analytics URL
const getAnalyticsUrl = computed(() => {
  const siteUrl = getSiteUrl.value

  // 格式化URL用于Google Analytics
  const normalizedUrl = siteUrl.startsWith('http') ? siteUrl : `https://${siteUrl}`

  // 跳转到Google Analytics
  return 'https://analytics.google.com/'
})

// 获取站点URL显示文本
const getSiteUrlDisplay = computed(() => {
  const siteUrl = getSiteUrl.value
  if (!siteUrl) {
    return '站点URL未配置'
  }
  if (siteUrl === 'https://example.com') {
    return '请配置正确的站点URL'
  }
  return siteUrl
})
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>