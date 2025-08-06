<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载系统配置</p>
          </div>
        </div>
      </div>
    </div>

    <div class="">
      <div class="max-w-7xl mx-auto">
        <!-- 配置表单 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <form @submit.prevent="saveConfig" class="space-y-6">

            <n-tabs type="line" animated>
              <n-tab-pane name="站点配置" tab="站点配置">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 网站标题 -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    网站标题 *
                  </label>
                  <n-input 
                    v-model:value="config.siteTitle" 
                    type="text" 
                    required
                    placeholder="老九网盘资源数据库"
                  />
                </div>

                <!-- 网站描述 -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    网站描述
                  </label>
                  <n-input 
                  v-model:value="config.siteDescription" 
                    type="text" 
                    placeholder="专业的老九网盘资源数据库"
                  />
                </div>

                <!-- 关键词 -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    关键词 (用逗号分隔)
                  </label>
                  <n-input 
                  v-model:value="config.keywords" 
                    type="text" 
                    placeholder="网盘,资源管理,文件分享"
                  />
                </div>

              
                <!-- 版权信息 -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    版权信息
                  </label>
                  <n-input 
                  v-model:value="config.copyright" 
                    type="text" 
                    placeholder="© 2024 老九网盘资源数据库"
                  />
                </div>

                <!-- 禁止词 -->
                <div class="md:col-span-2">
                  <div class="flex items-center justify-between mb-2">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                      违禁词
                    </label>
                    <div class="flex gap-2">
                      <n-button 
                        type="default" 
                        size="small" 
                        @click="openForbiddenWordsSource"
                      >
                        开源违禁词
                      </n-button>
                    </div>
                  </div>
                  <n-input
                  v-model:value="config.forbiddenWords"
                    type="textarea"
                    placeholder=""
                    :autosize="{ minRows: 4, maxRows: 8 }"
                  />
                </div>

                <!-- <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    每页显示数量
                  </label>
                  <select 
                    v-model.number="config.pageSize" 
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                  >
                    <option value="20">20 条</option>
                    <option value="50">50 条</option>
                    <option value="100">100 条</option>
                    <option value="200">200 条</option>
                  </select>
                </div> -->

                <!-- 系统维护模式 -->
                <div class="md:col-span-2 flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                  <div class="flex-1">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                      维护模式
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      开启后，普通用户无法访问系统
                    </p>
                  </div>
                  <div class="ml-4">
                    <label class="relative inline-flex items-center cursor-pointer">
                      <n-switch v-model:value="config.maintenanceMode" />
                    </label>
                  </div>
                </div>
              </div>
              </n-tab-pane>
              <n-tab-pane name="功能配置" tab="功能配置">
                <div class="space-y-4">
                  <div class="flex flex-col gap-1">
                    <!-- 待处理资源自动处理 -->
                    <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                    <div class="flex-1">
                      <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                        待处理资源自动处理
                      </h3>
                      <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                        开启后，系统将自动处理待处理的资源，无需手动操作
                      </p>
                    </div>
                    <div class="ml-4">
                      <label class="relative inline-flex items-center cursor-pointer">
                        <n-switch v-model:value="config.autoProcessReadyResources" />
                      </label>
                    </div>
                  </div>
                  <div v-if="config.autoProcessReadyResources" class="ml-6">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      自动处理间隔 (分钟)
                    </label>
                    <n-input 
                      v-model:value="config.autoProcessInterval" 
                      type="number" 
                      placeholder="30"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      建议设置 5-60 分钟，避免过于频繁的处理
                    </p>
                  </div>

                </div>
                
                

                <!-- 自动转存 -->
                <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                  <div class="flex-1">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                      自动转存
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      开启后，系统将自动转存资源到其他网盘平台
                    </p>
                  </div>
                  <div class="ml-4">
                    <label class="relative inline-flex items-center cursor-pointer">
                      <n-switch v-model:value="config.autoTransferEnabled" />
                    </label>
                  </div>
                </div>

                <!-- 自动转存配置（仅在开启时显示） -->
                <div v-if="config.autoTransferEnabled" class="ml-6 space-y-4">
                  <!-- 自动转存限制天数 -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      自动转存限制（n天内资源）
                    </label>
                    <n-input 
                      v-model:value="config.autoTransferLimitDays" 
                      type="number" 
                      placeholder="30"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      只转存指定天数内的资源，0表示不限制时间
                    </p>
                  </div>

                  <!-- 最小存储空间 -->
                  <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      最小存储空间（GB）
                    </label>
                    <n-input 
                      v-model:value="config.autoTransferMinSpace" 
                      type="number" 
                      placeholder="500"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      当网盘剩余空间小于此值时，停止自动转存（100-1024GB）
                    </p>
                  </div>
                </div>

                <!-- 自动拉取热播剧 -->
                <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                  <div class="flex-1">
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                      自动拉取热播剧
                    </h3>
                    <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      开启后，系统将自动从豆瓣获取热播剧信息
                    </p>
                  </div>
                  <div class="ml-4">
                    <label class="relative inline-flex items-center cursor-pointer">
                      <n-switch v-model:value="config.autoFetchHotDramaEnabled" />
                    </label>
                  </div>
                </div>

                <!-- 自动处理间隔 -->
                
              </div>
              </n-tab-pane>
              <n-tab-pane name="API配置" tab="API配置">
                <div class="space-y-4">
                <!-- API Token -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    公开API访问令牌
                  </label>
                  <div class="flex gap-2">
                    <n-input 
                      v-model:value="config.apiToken" 
                      type="password" 
                      placeholder="输入API Token，用于公开API访问认证"
                      :show-password-on="'click'"
                    />
                    <n-button 
                      v-if="!config.apiToken"
                      type="primary"
                      @click="generateApiToken"
                    >
                      生成
                    </n-button>
                    <template v-else>
                      <n-button 
                        type="primary"
                        @click="copyApiToken"
                      >
                        复制
                      </n-button>
                      <n-button 
                        type="default"
                        @click="generateApiToken"
                      >
                        重新生成
                      </n-button>
                    </template>
                  </div>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    用于公开API的访问认证，建议使用随机字符串
                  </p>
                </div>

                <!-- API使用说明 -->
                <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                  <h3 class="text-sm font-medium text-blue-800 dark:text-blue-200 mb-2">
                    <i class="fas fa-info-circle mr-1"></i>
                    API使用说明
                  </h3>
                  <div class="text-xs text-blue-700 dark:text-blue-300 space-y-1">
                    <p>• 批量添加资源: POST /api/public/resources/batch-add</p>
                    <p>• 资源搜索: GET /api/public/resources/search</p>
                    <p>• 热门剧: GET /api/public/hot-dramas</p>
                  </div>
                </div>
              </div>
              </n-tab-pane>
            </n-tabs>

            <!-- 保存按钮 -->
            <div class="flex justify-end space-x-4 pt-6">
              <n-button 
                type="tertiary"
                @click="resetForm"
              >
                重置
              </n-button>
              <n-button 
                type="primary"
                :disabled="saving"
                @click="saveConfig"
              >
                <i v-if="saving" class="fas fa-spinner fa-spin mr-2"></i>
                {{ saving ? '保存中...' : '保存配置' }}
              </n-button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

import { ref, onMounted } from 'vue'
import { useSystemConfigApi } from '~/composables/useApi'
import { useSystemConfigStore } from '~/stores/systemConfig'

// 权限检查已在 admin 布局中处理

const systemConfigStore = useSystemConfigStore()

// API
const systemConfigApi = useSystemConfigApi()
const notification = useNotification()

// 响应式数据
const loading = ref(false)
const loadingForbiddenWords = ref(false)
const config = ref({
  // SEO 配置
  siteTitle: '老九网盘资源数据库',
  siteDescription: '专业的老九网盘资源数据库',
  keywords: '网盘,资源管理,文件分享',
  author: '系统管理员',
  copyright: '© 2024 老九网盘资源数据库',
  
  // 自动处理配置
  autoProcessReadyResources: false,
  autoProcessInterval: 30,
  autoTransferEnabled: false, // 新增
  autoTransferLimitDays: 30, // 新增：自动转存限制天数
  autoTransferMinSpace: 500, // 新增：最小存储空间（GB）
  autoFetchHotDramaEnabled: false, // 新增
  
  // 其他配置
  pageSize: 100,
  maintenanceMode: false,
  apiToken: '' // 新增
})

// 系统配置状态（用于SEO）
const systemConfig = ref({
  site_title: '老九网盘资源数据库',
  site_description: '系统配置管理页面',
  keywords: '系统配置,管理',
  author: '系统管理员'
})
const originalConfig = ref(null)

// 页面元数据 - 移到变量声明之后
useHead({
  title: () => `${systemConfig.value.site_title} - 系统配置`,
  meta: [
    { 
      name: 'description', 
      content: () => systemConfig.value.site_description
    },
    { 
      name: 'keywords', 
      content: () => systemConfig.value.keywords
    },
    { 
      name: 'author', 
      content: () => systemConfig.value.author
    }
  ]
})

// 加载配置
const loadConfig = async () => {
  try {
    loading.value = true
    const response = await systemConfigApi.getSystemConfig()
    console.log('系统配置响应:', response)
    
    // 使用新的统一响应格式，直接使用response
    if (response) {
      const newConfig =  {
        siteTitle: response.site_title || '老九网盘资源数据库',
        siteDescription: response.site_description || '专业的老九网盘资源数据库',
        keywords: response.keywords || '网盘,资源管理,文件分享',
        author: response.author || '系统管理员',
        copyright: response.copyright || '© 2024 老九网盘资源数据库',
        autoProcessReadyResources: response.auto_process_ready_resources || false,
        autoProcessInterval: String(response.auto_process_interval || 30),
        autoTransferEnabled: response.auto_transfer_enabled || false, // 新增
        autoTransferLimitDays: String(response.auto_transfer_limit_days || 30), // 新增：自动转存限制天数
        autoTransferMinSpace: String(response.auto_transfer_min_space || 500), // 新增：最小存储空间（GB）
        autoFetchHotDramaEnabled: response.auto_fetch_hot_drama_enabled || false, // 新增
        forbiddenWords: formatForbiddenWordsForDisplay(response.forbidden_words || ''),
        pageSize: String(response.page_size || 100),
        maintenanceMode: response.maintenance_mode || false,
        apiToken: response.api_token || '' // 加载API Token
      }
      config.value = newConfig
      originalConfig.value = JSON.parse(JSON.stringify(newConfig)) // 深拷贝保存原始数据
      systemConfig.value = response // 更新系统配置状态
    }
  } catch (error) {
    console.error('加载配置失败:', error)
    // 显示错误提示
  } finally {
    loading.value = false
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    loading.value = true

    const changes = {}
    const currentConfig = config.value
    const original = originalConfig.value
    
    // 检查每个字段是否有变化
    if (currentConfig.siteTitle !== original.siteTitle) {
      changes.site_title = currentConfig.siteTitle
    }
    if (currentConfig.siteDescription !== original.siteDescription) {
      changes.site_description = currentConfig.siteDescription
    }
    if (currentConfig.keywords !== original.keywords) {
      changes.keywords = currentConfig.keywords
    }
    if (currentConfig.author !== original.author) {
      changes.author = currentConfig.author
    }
    if (currentConfig.copyright !== original.copyright) {
      changes.copyright = currentConfig.copyright
    }
    if (currentConfig.autoProcessReadyResources !== original.autoProcessReadyResources) {
      changes.auto_process_ready_resources = currentConfig.autoProcessReadyResources
    }
    if (currentConfig.autoProcessInterval !== original.autoProcessInterval) {
      changes.auto_process_interval = parseInt(currentConfig.autoProcessInterval) || 0
    }
    if (currentConfig.autoTransferEnabled !== original.autoTransferEnabled) {
      changes.auto_transfer_enabled = currentConfig.autoTransferEnabled
    }
    if (currentConfig.autoTransferLimitDays !== original.autoTransferLimitDays) {
      changes.auto_transfer_limit_days = parseInt(currentConfig.autoTransferLimitDays) || 0
    }
    if (currentConfig.autoTransferMinSpace !== original.autoTransferMinSpace) {
      changes.auto_transfer_min_space = parseInt(currentConfig.autoTransferMinSpace) || 0
    }
    if (currentConfig.autoFetchHotDramaEnabled !== original.autoFetchHotDramaEnabled) {
      changes.auto_fetch_hot_drama_enabled = currentConfig.autoFetchHotDramaEnabled
    }
    if (currentConfig.forbiddenWords !== original.forbiddenWords) {
      changes.forbidden_words = formatForbiddenWordsForSave(currentConfig.forbiddenWords)
    }
    if (currentConfig.pageSize !== original.pageSize) {
      changes.page_size = parseInt(currentConfig.pageSize) || 0
    }
    if (currentConfig.maintenanceMode !== original.maintenanceMode) {
      changes.maintenance_mode = currentConfig.maintenanceMode
    }
    if (currentConfig.apiToken !== original.apiToken) {
      changes.api_token = currentConfig.apiToken
    }
    
    console.log('检测到的变化:', changes)
    if (Object.keys(changes).length === 0) {
      notification.warning({
        content: '没有需要保存的配置',
        duration: 3000
      })
      return
    }
    const response = await systemConfigApi.updateSystemConfig(changes)
    // 使用新的统一响应格式，直接检查response是否存在
    if (response) {
      notification.success({
        content: '配置保存成功！',
        duration: 3000
      })
      await loadConfig()
      // 自动更新 systemConfig store（强制刷新）
      await systemConfigStore.initConfig(true)
    } else {
      notification.error({
        content: '保存配置失败：未知错误',
        duration: 3000
      })
    }
  } catch (error) {
    notification.error({
      content: '保存配置失败：' + (error.message || '未知错误'),
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  notification.confirm({
      title: '确定要重置所有配置吗？',
      content: '重置后，所有配置将恢复为修改前配置',
      duration: 3000,
      onOk: () => {
      loadConfig()
    }
  })
}

// 生成API Token
const generateApiToken = () => {
  const newToken = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
  config.value.apiToken = newToken;
  notification.success({
    content: '新API Token已生成: ' + newToken,
    duration: 3000
  })
};

// 复制API Token
const copyApiToken = async () => {
  try {
    await navigator.clipboard.writeText(config.value.apiToken);
    notification.success({
      content: 'API Token已复制到剪贴板',
      duration: 3000
    });
  } catch (err) {
    // 降级方案：使用传统的复制方法
    const textArea = document.createElement('textarea');
    textArea.value = config.value.apiToken;
    document.body.appendChild(textArea);
    textArea.select();
    try {
      document.execCommand('copy');
      notification.success({
        content: 'API Token已复制到剪贴板',
        duration: 3000
      });
    } catch (fallbackErr) {
      notification.error({
        content: '复制失败，请手动复制',
        duration: 3000
      });
    }
    document.body.removeChild(textArea);
  }
};


// 打开违禁词源文件
const openForbiddenWordsSource = () => {
  const url = 'https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/db/forbidden.txt'
  window.open(url, '_blank', 'noopener,noreferrer')
}

// 格式化违禁词用于显示（逗号分隔转为多行）
const formatForbiddenWordsForDisplay = (forbiddenWords) => {
  if (!forbiddenWords) return ''
  
  // 按逗号分割，过滤空字符串，然后按行显示
  return forbiddenWords.split(',')
    .map(word => word.trim())
    .filter(word => word.length > 0)
    .join('\n')
}

// 格式化违禁词用于保存（多行转为逗号分隔）
const formatForbiddenWordsForSave = (forbiddenWords) => {
  if (!forbiddenWords) return ''
  
  // 按行分割，过滤空行，然后用逗号连接
  return forbiddenWords.split('\n')
    .map(line => line.trim())
    .filter(line => line.length > 0)
    .join(',')
}

// 页面加载时获取配置
onMounted(() => {
  loadConfig()
})
</script> 