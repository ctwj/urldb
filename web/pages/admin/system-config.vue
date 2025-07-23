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
              <n-tab-pane name="SEO 配置" tab="SEO 配置">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 网站标题 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    网站标题 *
                  </label>
                  <input 
                    v-model="config.siteTitle" 
                    type="text" 
                    required
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="老九网盘资源数据库"
                  />
                </div>

                <!-- 网站描述 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    网站描述
                  </label>
                  <input 
                    v-model="config.siteDescription" 
                    type="text" 
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="专业的老九网盘资源数据库"
                  />
                </div>

                <!-- 关键词 -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    关键词 (用逗号分隔)
                  </label>
                  <input 
                    v-model="config.keywords" 
                    type="text" 
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="网盘,资源管理,文件分享"
                  />
                </div>

                <!-- 作者 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    作者
                  </label>
                  <input 
                    v-model="config.author" 
                    type="text" 
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="系统管理员"
                  />
                </div>

                <!-- 版权信息 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    版权信息
                  </label>
                  <input 
                    v-model="config.copyright" 
                    type="text" 
                    class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="© 2024 老九网盘资源数据库"
                  />
                </div>
              </div>
              </n-tab-pane>
              <n-tab-pane name="自动处理配置" tab="自动处理配置">
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
                        <input 
                          v-model="config.autoProcessReadyResources" 
                          type="checkbox" 
                          class="sr-only peer"
                        />
                        <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                      </label>
                    </div>
                  </div>
                  <div v-if="config.autoProcessReadyResources" class="ml-6">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      自动处理间隔 (分钟)
                    </label>
                    <input 
                      v-model.number="config.autoProcessInterval" 
                      type="number" 
                      min="1"
                      max="1440"
                      class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
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
                      <input 
                        v-model="config.autoTransferEnabled" 
                        type="checkbox" 
                        class="sr-only peer"
                      />
                      <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
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
                    <input 
                      v-model.number="config.autoTransferLimitDays" 
                      type="number" 
                      min="0"
                      max="365"
                      class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
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
                    <input 
                      v-model.number="config.autoTransferMinSpace" 
                      type="number" 
                      min="100"
                      max="1024"
                      class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
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
                      <input 
                        v-model="config.autoFetchHotDramaEnabled" 
                        type="checkbox" 
                        class="sr-only peer"
                      />
                      <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
                    </label>
                  </div>
                </div>

                <!-- 自动处理间隔 -->
                
              </div>
              </n-tab-pane>
              <n-tab-pane name="其他配置" tab="其他配置">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 每页显示数量 -->
                <div>
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
                </div>

                <!-- 系统维护模式 -->
                <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
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
                      <input 
                        v-model="config.maintenanceMode" 
                        type="checkbox" 
                        class="sr-only peer"
                      />
                      <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-red-300 dark:peer-focus:ring-red-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-red-600"></div>
                    </label>
                  </div>
                </div>
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
                    <input 
                      v-model="config.apiToken" 
                      type="text" 
                      class="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                      placeholder="输入API Token，用于公开API访问认证"
                    />
                    <button 
                      type="button"
                      @click="generateApiToken"
                      class="px-4 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700 transition-colors"
                    >
                      生成
                    </button>
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
                    <p>• 单个添加资源: POST /api/public/resources/add</p>
                    <p>• 批量添加资源: POST /api/public/resources/batch-add</p>
                    <p>• 资源搜索: GET /api/public/resources/search</p>
                    <p>• 热门剧: GET /api/public/hot-dramas</p>
                    <p>• 认证方式: 在请求头中添加 X-API-Token 或在查询参数中添加 api_token</p>
                    <p>• Swagger文档: <a href="/swagger/index.html" target="_blank" class="underline">查看完整API文档</a></p>
                  </div>
                </div>
              </div>
              </n-tab-pane>
            </n-tabs>

            <!-- 保存按钮 -->
            <div class="flex justify-end space-x-4 pt-6">
              <button 
                type="button"
                @click="resetForm"
                class="px-6 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                重置
              </button>
              <button 
                type="submit"
                :disabled="saving"
                class="px-6 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <i v-if="saving" class="fas fa-spinner fa-spin mr-2"></i>
                {{ saving ? '保存中...' : '保存配置' }}
              </button>
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
  layout: 'admin'
})

import { ref, onMounted } from 'vue'
import { useSystemConfigApi } from '~/composables/useApi'
import { useSystemConfigStore } from '~/stores/systemConfig'
const systemConfigStore = useSystemConfigStore()

// API
const systemConfigApi = useSystemConfigApi()

// 响应式数据
const loading = ref(false)
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
const systemConfig = ref(null)

// 页面元数据 - 移到变量声明之后
useHead({
  title: () => systemConfig.value?.site_title ? `${systemConfig.value.site_title} - 系统配置` : '系统配置 - 老九网盘资源数据库',
  meta: [
    { 
      name: 'description', 
      content: () => systemConfig.value?.site_description || '系统配置管理页面' 
    },
    { 
      name: 'keywords', 
      content: () => systemConfig.value?.keywords || '系统配置,管理' 
    },
    { 
      name: 'author', 
      content: () => systemConfig.value?.author || '系统管理员' 
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
      config.value = {
        siteTitle: response.site_title || '老九网盘资源数据库',
        siteDescription: response.site_description || '专业的老九网盘资源数据库',
        keywords: response.keywords || '网盘,资源管理,文件分享',
        author: response.author || '系统管理员',
        copyright: response.copyright || '© 2024 老九网盘资源数据库',
        autoProcessReadyResources: response.auto_process_ready_resources || false,
        autoProcessInterval: response.auto_process_interval || 30,
        autoTransferEnabled: response.auto_transfer_enabled || false, // 新增
        autoTransferLimitDays: response.auto_transfer_limit_days || 30, // 新增：自动转存限制天数
        autoTransferMinSpace: response.auto_transfer_min_space || 500, // 新增：最小存储空间（GB）
        autoFetchHotDramaEnabled: response.auto_fetch_hot_drama_enabled || false, // 新增
        pageSize: response.page_size || 100,
        maintenanceMode: response.maintenance_mode || false,
        apiToken: response.api_token || '' // 加载API Token
      }
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
    
    const requestData = {
      site_title: config.value.siteTitle,
      site_description: config.value.siteDescription,
      keywords: config.value.keywords,
      author: config.value.author,
      copyright: config.value.copyright,
      auto_process_ready_resources: config.value.autoProcessReadyResources,
      auto_process_interval: config.value.autoProcessInterval,
      auto_transfer_enabled: config.value.autoTransferEnabled, // 新增
      auto_transfer_limit_days: config.value.autoTransferLimitDays, // 新增：自动转存限制天数
      auto_transfer_min_space: config.value.autoTransferMinSpace, // 新增：最小存储空间（GB）
      auto_fetch_hot_drama_enabled: config.value.autoFetchHotDramaEnabled, // 新增
      page_size: config.value.pageSize,
      maintenance_mode: config.value.maintenanceMode,
      api_token: config.value.apiToken // 保存API Token
    }
    
    const response = await systemConfigApi.updateSystemConfig(requestData)
    // 使用新的统一响应格式，直接检查response是否存在
    if (response) {
      alert('配置保存成功！')
      await loadConfig()
      // 自动更新 systemConfig store（强制刷新）
      await systemConfigStore.initConfig(true)
    } else {
      alert('保存配置失败：未知错误')
    }
  } catch (error) {
    console.error('保存配置失败:', error)
    alert('保存配置失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  if (confirm('确定要重置所有配置吗？')) {
    loadConfig()
  }
}

// 生成API Token
const generateApiToken = () => {
  const newToken = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
  config.value.apiToken = newToken;
  alert('新API Token已生成: ' + newToken);
};

// 页面加载时获取配置
onMounted(() => {
  loadConfig()
})
</script> 