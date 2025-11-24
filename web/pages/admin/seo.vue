<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">SEO管理</h1>
        <p class="text-gray-600 dark:text-gray-400">搜索引擎优化管理</p>
      </div>
    </template>

    <!-- 内容区 -->
    <template #content>
      <div class="config-content h-full">
        <!-- Tab导航 -->
        <n-tabs v-model:value="activeTab" type="line" animated>
          <!-- Google Index Tab -->
          <n-tab-pane name="google-index" tab="Google Index">
            <GoogleIndexTab
              :system-config="systemConfig"
              :google-index-config="googleIndexConfig"
              :google-index-stats="googleIndexStats"
              :tasks="googleIndexTasks"
              :credentials-status="credentialsStatus"
              :credentials-status-message="credentialsStatusMessage"
              :config-loading="configLoading"
              :manual-check-loading="manualCheckLoading"
              :submit-sitemap-loading="submitSitemapLoading"
              :tasks-loading="tasksLoading"
              :pagination="googleIndexPagination"
              @update:google-index-config="updateGoogleIndexConfig"
              @show-verification="showVerificationModal = true"
              @show-credentials-guide="showCredentialsGuide = true"
              @select-credentials-file="selectCredentialsFile"
              @manual-check-urls="manualCheckURLs"
              @refresh-status="refreshGoogleIndexStatus"
              @view-task-items="viewTaskItems"
              @start-task="startTask"
            />
          </n-tab-pane>

          <!-- Sitemap管理 Tab -->
          <n-tab-pane name="sitemap" tab="Sitemap管理">
            <SitemapTab
              :system-config="systemConfig"
              :sitemap-config="sitemapConfig"
              :sitemap-stats="sitemapStats"
              :config-loading="configLoading"
              :is-generating="isGenerating"
              :generate-status="generateStatus"
              @update:sitemap-config="updateSitemapConfig"
              @refresh-status="refreshSitemapStatus"
            />
          </n-tab-pane>

          <!-- 站点提交 Tab -->
          <n-tab-pane name="site-submit" tab="站点提交">
            <SiteSubmitTab
              :last-submit-time="lastSubmitTime"
              @update:last-submit-time="updateLastSubmitTime"
            />
          </n-tab-pane>

          <!-- 外链建设 Tab -->
          <n-tab-pane name="link-building" tab="外链建设">
            <LinkBuildingTab
              :link-stats="linkStats"
              :link-list="linkList"
              :loading="linkLoading"
              :pagination="linkPagination"
              @add-new-link="addNewLink"
              @edit-link="editLink"
              @delete-link="deleteLink"
              @load-link-list="loadLinkList"
            />
          </n-tab-pane>
        </n-tabs>
      </div>
    </template>
  </AdminPageLayout>

  <!-- URL检查模态框 -->
  <n-modal v-model:show="urlCheckModal.show" preset="card" title="手动检查URL" style="max-width: 600px;">
    <div class="space-y-4">
      <p class="text-gray-600 dark:text-gray-400">输入要检查索引状态的URL，每行一个</p>
      <n-input
        v-model:value="urlCheckModal.urls"
        type="textarea"
        :autosize="{ minRows: 4, maxRows: 8 }"
        placeholder="https://yoursite.com/page1&#10;https://yoursite.com/page2"
      />
      <div class="flex justify-end space-x-2">
        <n-button @click="urlCheckModal.show = false">取消</n-button>
        <n-button type="primary" @click="confirmManualCheckURLs" :loading="manualCheckLoading">确认</n-button>
      </div>
    </div>
  </n-modal>

  <!-- 所有权验证模态框 -->
  <n-modal v-model:show="showVerificationModal" preset="card" title="站点所有权验证" style="max-width: 600px;">
    <div class="space-y-6">
      <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="fas fa-info-circle text-blue-500 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800 dark:text-blue-200">DNS方式验证</h3>
            <div class="mt-2 text-sm text-blue-700 dark:text-blue-300">
              <p>推荐使用DNS方式验证站点所有权，这是最安全和可靠的方法：</p>
              <ol class="list-decimal list-inside mt-2 space-y-1">
                <li>登录您的域名注册商或DNS管理平台</li>
                <li>添加一条TXT记录</li>
                <li>在Google Search Console中输入您的验证字符串</li>
                <li>验证DNS TXT记录是否生效</li>
              </ol>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">注意事项</h3>
            <div class="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
              <ul class="list-disc list-inside space-y-1">
                <li>DNS验证比HTML标签更安全，不易被其他网站复制</li>
                <li>验证成功后，Google会自动检测您的站点所有权</li>
                <li>如果您的域名服务商不支持TXT记录，请联系客服寻求帮助</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end pt-2">
        <n-button type="primary" @click="showVerificationModal = false">
          确定
        </n-button>
      </div>
    </div>
  </n-modal>

  <!-- 申请凭据说明抽屉 -->
  <n-drawer v-model:show="showCredentialsGuide" :width="600" placement="right">
    <n-drawer-content title="如何申请Google Search Console API凭据" closable>
      <div class="space-y-6">
        <!-- 步骤1 -->
        <div class="border-l-4 border-blue-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-1 text-blue-500 mr-2"></i>创建Google Cloud项目
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            首先需要在Google Cloud Console中创建一个新项目或选择现有项目。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>访问 <a href="https://console.cloud.google.com/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Google Cloud Console</a></li>
            <li>点击顶部的项目选择器</li>
            <li>点击"新建项目"或选择现有项目</li>
            <li>输入项目名称，点击"创建"</li>
          </ol>
        </div>

        <!-- 步骤2 -->
        <div class="border-l-4 border-green-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-2 text-green-500 mr-2"></i>启用Search Console API
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            在项目中启用Google Search Console API。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>在导航菜单中选择"API和服务" > "库"</li>
            <li>搜索"Google Search Console API"</li>
            <li>点击搜索结果中的"Google Search Console API"</li>
            <li>点击"启用"按钮</li>
          </ol>
        </div>

        <!-- 步骤3 -->
        <div class="border-l-4 border-yellow-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-3 text-yellow-500 mr-2"></i>创建服务账号
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            创建服务账号并生成JSON密钥文件。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>在导航菜单中选择"API和服务" > "凭据"</li>
            <li>点击"创建凭据" > "服务账号"</li>
            <li>输入服务账号名称（如：google-index-api）</li>
            <li>点击"创建并继续"</li>
            <li>在角色选择中，选择"项目" > "编辑者"</li>
            <li>点击"继续"，然后点击"完成"</li>
          </ol>
        </div>

        <!-- 步骤4 -->
        <div class="border-l-4 border-purple-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-4 text-purple-500 mr-2"></i>生成JSON密钥
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            为服务账号生成JSON格式的密钥文件。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>在服务账号列表中找到刚创建的服务账号</li>
            <li>点击服务账号名称进入详情页面</li>
            <li>切换到"密钥"标签页</li>
            <li>点击"添加密钥" > "创建新密钥"</li>
            <li>选择"JSON"作为密钥类型</li>
            <li>点击"创建"，JSON文件将自动下载</li>
          </ol>
        </div>

        <!-- 步骤5 -->
        <div class="border-l-4 border-red-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-5 text-red-500 mr-2"></i>验证网站所有权
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            在Google Search Console中验证网站并添加服务账号权限。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>访问 <a href="https://search.google.com/search-console/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Google Search Console</a></li>
            <li>添加属性并验证网站所有权</li>
            <li>在设置中找到"用户和权限"</li>
            <li>点击"添加用户"，输入服务账号的邮箱地址</li>
            <li>授予"所有者"或"完整"权限</li>
          </ol>
        </div>

        <!-- 注意事项 -->
        <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
          <h5 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2">
            <i class="fas fa-exclamation-triangle mr-2"></i>重要注意事项
          </h5>
          <ul class="space-y-1 text-sm text-yellow-700 dark:text-yellow-300">
            <li>• 请妥善保管下载的JSON密钥文件，不要泄露给他人</li>
            <li>• 服务账号邮箱地址通常格式为：xxx@xxx.iam.gserviceaccount.com</li>
            <li>• API配额有限制，请合理使用避免超出限制</li>
            <li>• 确保网站已在Search Console中验证所有权</li>
          </ul>
        </div>

        <!-- 完成按钮 -->
        <div class="flex justify-end">
          <n-button type="primary" @click="showCredentialsGuide = false">
            我已了解
          </n-button>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>

  <!-- 隐藏的文件输入 -->
  <input
    type="file"
    ref="credentialsFileInput"
    accept=".json"
    @change="handleCredentialsFileSelect"
    style="display: none;"
  />
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'
import GoogleIndexTab from '~/components/Admin/GoogleIndexTab.vue'
import SitemapTab from '~/components/Admin/SitemapTab.vue'
import SiteSubmitTab from '~/components/Admin/SiteSubmitTab.vue'
import LinkBuildingTab from '~/components/Admin/LinkBuildingTab.vue'

// SEO管理页面
definePageMeta({
  layout: 'admin'
})

import { useMessage } from 'naive-ui'
import { useApi } from '~/composables/useApi'
import { ref, onMounted, watch } from 'vue'

// 获取消息组件
const message = useMessage()

// 当前激活的Tab - 默认显示 Google Index
const activeTab = ref('google-index')

// 获取系统配置
const systemConfig = ref<any>(null)


// Google索引配置
const googleIndexConfig = ref({
  enabled: false,
  siteURL: '',
  credentialsFile: '',
  checkInterval: 60,
  batchSize: 100,
  concurrency: 5
})

// 凭据验证相关
const credentialsStatus = ref<string | null>(null)
const credentialsStatusMessage = ref('')
const credentialsFileInput = ref<HTMLInputElement | null>(null)

// 申请凭据抽屉显示状态
const showCredentialsGuide = ref(false)

// 所有权验证相关
const showVerificationModal = ref(false)

// Google索引统计
const googleIndexStats = ref({
  totalURLs: 0,
  indexedURLs: 0,
  errorURLs: 0,
  totalTasks: 0,
  runningTasks: 0,
  completedTasks: 0,
  failedTasks: 0
})

// Google索引任务列表
const googleIndexTasks = ref([])
const tasksLoading = ref(false)

// 分页配置
const googleIndexPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  itemCount: 0,
  onChange: (page: number) => {
    googleIndexPagination.value.page = page
    loadGoogleIndexTasks()
  },
  onUpdatePageSize: (pageSize: number) => {
    googleIndexPagination.value.pageSize = pageSize
    googleIndexPagination.value.page = 1
    loadGoogleIndexTasks()
  }
})

// 模态框状态
const urlCheckModal = ref({
  show: false,
  urls: ''
})

// 加载状态
const configLoading = ref(false)
const manualCheckLoading = ref(false)
const submitSitemapLoading = ref(false)

// Sitemap管理相关
const sitemapConfig = ref({
  autoGenerate: false,
  lastGenerate: '',
  lastUpdate: ''
})

const sitemapStats = ref({
  total_resources: 0,
  total_pages: 0,
  last_generate: ''
})

const isGenerating = ref(false)
const generateStatus = ref('')

// 最后提交时间
const lastSubmitTime = ref({
  baidu: '',
  google: '',
  bing: '',
  sogou: '',
  shenma: '',
  so360: ''
})

// 外链统计
const linkStats = ref({
  total: 156,
  valid: 142,
  pending: 8,
  invalid: 6
})

// 外链列表
const linkList = ref([
  {
    id: 1,
    url: 'https://example1.com',
    title: '示例外链1',
    status: 'valid',
    domain: 'example1.com',
    created_at: '2024-01-15'
  },
  {
    id: 2,
    url: 'https://example2.com',
    title: '示例外链2',
    status: 'pending',
    domain: 'example2.com',
    created_at: '2024-01-16'
  },
  {
    id: 3,
    url: 'https://example3.com',
    title: '示例外链3',
    status: 'invalid',
    domain: 'example3.com',
    created_at: '2024-01-17'
  }
])

const linkLoading = ref(false)

// 分页配置
const linkPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    linkPagination.value.page = page
    loadLinkList()
  },
  onUpdatePageSize: (pageSize: number) => {
    linkPagination.value.pageSize = pageSize
    linkPagination.value.page = 1
    loadLinkList()
  }
})

// 加载系统配置
const loadSystemConfig = async () => {
  try {
    const { useSystemConfigStore } = await import('~/stores/systemConfig')
    const systemConfigStore = useSystemConfigStore()
    await systemConfigStore.initConfig(true, true)
    systemConfig.value = systemConfigStore.config
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 加载Google索引配置
const loadGoogleIndexConfig = async () => {
  try {
    console.log('开始加载 Google 索引配置...')
    const api = useApi()
    const configs = await api.googleIndexApi.getGoogleIndexConfig()
    console.log('获取到的配置:', configs)
    if (configs) {
      // 查找general配置
      const generalConfig = configs.find((c: any) => c.group === 'general')
      const authConfig = configs.find((c: any) => c.group === 'auth')
      console.log('找到的配置 - general:', generalConfig, 'auth:', authConfig)

      let newConfig = { ...googleIndexConfig.value }

      if (generalConfig) {
        const configData = JSON.parse(generalConfig.value)
        newConfig.enabled = configData.enabled || false
        newConfig.siteURL = configData.siteURL || ''
        newConfig.checkInterval = configData.checkInterval || 60
        newConfig.batchSize = configData.batchSize || 100
        newConfig.concurrency = configData.concurrency || 5
      }

      if (authConfig) {
        console.log('解析 auth 配置:', authConfig.value)
        const authData = JSON.parse(authConfig.value)
        console.log('解析后的 authData:', authData)
        newConfig.credentialsFile = authData.credentialsFile || authData.credentials_file || ''
        console.log('设置凭据文件路径:', newConfig.credentialsFile)
      }

      // 强制触发响应式更新
      googleIndexConfig.value = newConfig
      console.log('最终配置:', googleIndexConfig.value)
    }
  } catch (error) {
    console.error('获取Google索引配置失败:', error)
  }
}

// 选择凭据文件
const selectCredentialsFile = () => {
  if (credentialsFileInput.value) {
    credentialsFileInput.value.click()
  }
}

// 处理凭据文件选择
const handleCredentialsFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) {
    return
  }

  // 验证文件类型
  if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
    message.error('请上传JSON格式的凭据文件')
    return
  }

  // 验证文件大小 (2MB限制)
  if (file.size > 2 * 1024 * 1024) {
    message.error('文件大小不能超过2MB')
    return
  }

  // 上传文件
  try {
    const api = useApi()
    const response = await api.googleIndexApi.uploadCredentials(file)

    // 检查API是否成功（success字段为true）且包含有效的文件路径
    if (response?.success === true && response?.file_path) {
      console.log('上传成功，文件路径:', response.file_path)
      // 统一路径格式为 Unix 格式
      const normalizedPath = response.file_path.replace(/\\/g, '/')
      console.log('标准化路径:', normalizedPath)
      // 强制触发响应式更新
      googleIndexConfig.value = {
        ...googleIndexConfig.value,
        credentialsFile: normalizedPath
      }
      console.log('更新后的 googleIndexConfig:', googleIndexConfig.value)
      message.success(response.message || '凭据文件上传成功，请验证凭据')

      // 清空文件输入以允许重新选择相同文件
      if (credentialsFileInput.value) {
        credentialsFileInput.value.value = ''
      }

      // 上传成功后立即更新后端配置并重新加载配置
      try {
        const configData = {
          group: 'auth',
          key: 'credentials_file',
          value: JSON.stringify({
            credentials_file: googleIndexConfig.value.credentialsFile.replace(/\\/g, '/')
          })
        }
        console.log('更新后端配置，发送数据:', JSON.stringify(configData, null, 2))

        const updateResponse = await api.googleIndexApi.updateGoogleIndexGroupConfig(configData)
        console.log('后端配置更新响应:', updateResponse)

        // 等待一下确保后端处理完成
        await new Promise(resolve => setTimeout(resolve, 500))

        // 重新加载配置以确保UI状态与后端同步
        console.log('重新加载配置...')
        await loadGoogleIndexConfig()
        console.log('配置重新加载完成')
      } catch (configError) {
        console.error('更新配置失败:', configError)
        message.error('配置更新失败，但文件已上传')

        // 即使配置更新失败，也尝试刷新状态
        setTimeout(async () => {
          console.log('延迟重新加载配置...')
          await loadGoogleIndexConfig()
        }, 1000)
      }
    } else {
      // 如果API调用成功但返回的数据有问题，或者API调用失败
      message.error(response?.message || '上传响应格式错误')
    }
  } catch (error: any) {
    console.error('凭据文件上传失败:', error)
    message.error('凭据文件上传失败: ' + (error?.message || '未知错误'))
  }
}

// 更新Google索引配置
const updateGoogleIndexConfig = async () => {
  configLoading.value = true
  try {
    const api = useApi()

    // 更新general配置
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'general',
      key: 'general',
      value: JSON.stringify({
        enabled: googleIndexConfig.value.enabled,
        siteURL: systemConfig.value?.site_url || googleIndexConfig.value.siteURL,
        checkInterval: googleIndexConfig.value.checkInterval,
        batchSize: googleIndexConfig.value.batchSize,
        concurrency: googleIndexConfig.value.concurrency || 5
      })
    })

    message.success('Google索引配置已更新')
  } catch (error) {
    console.error('更新Google索引配置失败:', error)
    message.error('更新配置失败')
  } finally {
    configLoading.value = false
  }
}

// 刷新Google索引状态
const refreshGoogleIndexStatus = async () => {
  try {
    const api = useApi()

    // 加载统计信息
    const statsResponse = await api.googleIndexApi.getGoogleIndexStatus()
    if (statsResponse) {
      googleIndexStats.value = statsResponse
    }

    // 加载任务列表
    await loadGoogleIndexTasks()
  } catch (error) {
    console.error('刷新Google索引状态失败:', error)
    message.error('刷新状态失败')
  }
}

// 加载Google索引任务列表
const loadGoogleIndexTasks = async () => {
  tasksLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexTasks({
      page: googleIndexPagination.value.page,
      pageSize: googleIndexPagination.value.pageSize
    })
    if (response) {
      googleIndexTasks.value = response.tasks || []
      googleIndexPagination.value.itemCount = response.total || 0
    }
  } catch (error) {
    console.error('加载Google索引任务列表失败:', error)
    message.error('加载任务列表失败')
  } finally {
    tasksLoading.value = false
  }
}

// 手动检查URL
const manualCheckURLs = () => {
  urlCheckModal.value.show = true
  urlCheckModal.value.urls = ''
}

// 确认手动检查URL
const confirmManualCheckURLs = async () => {
  const urls = urlCheckModal.value.urls.split('\n').filter(url => url.trim() !== '')
  if (urls.length === 0) {
    message.warning('请至少输入一个URL')
    return
  }

  manualCheckLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.createGoogleIndexTask({
      title: `手动URL检查任务 - ${new Date().toLocaleString('zh-CN')}`,
      type: 'status_check',
      description: `手动检查 ${urls.length} 个URL的索引状态`,
      URLs: urls
    })
    if (response) {
      message.success('URL检查任务已创建')
      urlCheckModal.value.show = false
      await refreshGoogleIndexStatus()
    }
  } catch (error) {
    console.error('手动检查URL失败:', error)
    message.error('手动检查URL失败')
  } finally {
    manualCheckLoading.value = false
  }
}

// 查看任务详情
const viewTaskItems = async (taskId: number) => {
  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexTaskItems(taskId)
    if (response) {
      // 在新窗口中打开任务详情
      const items = response.items || []
      const content = `任务 ${taskId} 详情:\n\n` +
        items.map((item: any) =>
          `URL: ${item.URL}\n状态: ${item.status}\n索引状态: ${item.indexStatus}\n错误信息: ${item.errorMessage}\n---\n`
        ).join('')
      alert(content)
    }
  } catch (error) {
    console.error('获取任务详情失败:', error)
    message.error('获取任务详情失败')
  }
}

// 启动任务
const startTask = async (taskId: number) => {
  try {
    const api = useApi()
    const response = await api.googleIndexApi.startGoogleIndexTask(taskId)
    if (response) {
      message.success('任务已启动')
      await loadGoogleIndexTasks()
    }
  } catch (error) {
    console.error('启动任务失败:', error)
    message.error('启动任务失败')
  }
}

// 获取Sitemap配置
const loadSitemapConfig = async () => {
  try {
    const api = useApi()
    const response = await api.sitemapApi.getSitemapConfig()
    if (response) {
      sitemapConfig.value = response
    }
  } catch (error) {
    message.error('获取Sitemap配置失败')
  }
}

// 更新Sitemap配置
const updateSitemapConfig = async (value: boolean) => {
  configLoading.value = true
  try {
    const api = useApi()
    await api.sitemapApi.updateSitemapConfig({
      autoGenerate: value,
      lastGenerate: sitemapConfig.value.lastGenerate,
      lastUpdate: new Date().toISOString()
    })
    message.success(value ? '自动生成功能已开启' : '自动生成功能已关闭')
  } catch (error) {
    message.error('更新配置失败')
  } finally {
    configLoading.value = false
  }
}

// 刷新Sitemap状态
const refreshSitemapStatus = async () => {
  try {
    const api = useApi()
    const response = await api.sitemapApi.getSitemapStatus()
    if (response) {
      sitemapStats.value = response
      generateStatus.value = '状态已刷新'
    }
  } catch (error) {
    message.error('刷新状态失败')
  }
}

// 更新最后提交时间
const updateLastSubmitTime = (engine: string, time: string) => {
  lastSubmitTime.value[engine as keyof typeof lastSubmitTime.value] = time
}

// 加载外链列表
const loadLinkList = () => {
  linkLoading.value = true
  setTimeout(() => {
    linkLoading.value = false
  }, 1000)
}

// 添加新外链
const addNewLink = () => {
  message.info('添加外链功能开发中')
}

// 编辑外链
const editLink = (row: any) => {
  message.info(`编辑外链: ${row.title}`)
}

// 删除外链
const deleteLink = (row: any) => {
  message.warning(`删除外链: ${row.title}`)
}

// 初始化
onMounted(async () => {
  await loadSystemConfig()
  await loadGoogleIndexConfig()
  await refreshGoogleIndexStatus()
  await loadSitemapConfig()
  await refreshSitemapStatus()
  loadLinkList()
})
</script>

<style scoped>
/* SEO管理页面样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}
</style>