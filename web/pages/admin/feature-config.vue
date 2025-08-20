<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">功能配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统功能开关和参数设置</p>
      </div>
      <n-button type="primary" @click="saveConfig" :loading="saving">
        <template #icon>
          <i class="fas fa-save"></i>
        </template>
        保存配置
      </n-button>
    </div>

    <!-- 配置表单 -->
    <n-card>
      <!-- 顶部Tabs -->
      <n-tabs
        v-model:value="activeTab"
        type="line"
        animated
        class="mb-6"
      >
        <n-tab-pane name="resource" tab="资源处理">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-8">
              <!-- 自动处理配置组 -->
              <div class="space-y-4">
                <div class="flex items-center space-x-2 mb-4">
                  <div class="w-1 h-6 bg-blue-500 rounded-full"></div>
                  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">自动处理配置</h3>
                </div>
                
                <!-- 自动处理 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">待处理资源自动处理</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">开启后，系统将自动处理待处理的资源，无需手动操作</span>
                  </div>
                  <n-switch v-model:value="configForm.auto_process_enabled" />
                </div>

                <!-- 自动处理间隔 -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">自动处理间隔 (分钟)</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">建议设置 5-60 分钟，避免过于频繁的处理</span>
                  </div>
                  <n-input
                    v-model:value="configForm.auto_process_interval"
                    type="text"
                    placeholder="30"
                    :disabled="!configForm.auto_process_enabled"
                  />
                </div>
              </div>

              <!-- Meilisearch搜索优化配置组 -->
              <div class="space-y-4">
                <div class="flex items-center space-x-2 mb-4">
                  <div class="w-1 h-6 bg-green-500 rounded-full"></div>
                  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">搜索优化配置</h3>
                </div>
                
                <!-- 启用Meilisearch -->
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">启用Meilisearch搜索优化</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">开启后，系统将使用Meilisearch提供更快的搜索体验</span>
                  </div>
                  <n-switch v-model:value="configForm.meilisearch_enabled" />
                </div>

                <!-- Meilisearch服务器配置 -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4" :class="{ 'opacity-50': !configForm.meilisearch_enabled }">
                  <!-- 服务器地址 -->
                  <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300">服务器地址</label>
                    <n-input
                      v-model:value="configForm.meilisearch_host"
                      placeholder="localhost"
                      :disabled="!configForm.meilisearch_enabled"
                    />
                  </div>

                  <!-- 端口 -->
                  <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300">端口</label>
                    <n-input
                      v-model:value="configForm.meilisearch_port"
                      placeholder="7700"
                      :disabled="!configForm.meilisearch_enabled"
                    />
                  </div>

                  <!-- 主密钥 -->
                  <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300">主密钥 (可选)</label>
                    <n-input
                      v-model:value="configForm.meilisearch_master_key"
                      placeholder="留空表示无认证"
                      type="password"
                      show-password-on="click"
                      :disabled="!configForm.meilisearch_enabled"
                    />
                  </div>

                  <!-- 索引名称 -->
                  <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-700 dark:text-gray-300">索引名称</label>
                    <n-input
                      v-model:value="configForm.meilisearch_index_name"
                      placeholder="resources"
                      :disabled="!configForm.meilisearch_enabled"
                    />
                  </div>
                </div>

                <!-- 操作按钮组 -->
                <div class="flex items-center space-x-3">
                  <n-button 
                    type="info" 
                    size="small"
                    :disabled="!configForm.meilisearch_enabled"
                    @click="testMeilisearchConnection"
                    :loading="testingConnection"
                  >
                    <template #icon>
                      <i class="fas fa-plug"></i>
                    </template>
                    测试连接
                  </n-button>
                  
                  <n-button 
                    type="primary" 
                    size="small"
                    @click="navigateTo('/admin/meilisearch-management')"
                  >
                    <template #icon>
                      <i class="fas fa-cogs"></i>
                    </template>
                    搜索优化管理
                  </n-button>

                  <!-- 健康状态和未同步数量显示 -->
                  <div v-if="meilisearchStatus" class="flex items-center space-x-4 ml-4">
                    <div class="flex items-center space-x-2">
                      <div class="w-2 h-2 rounded-full" :class="meilisearchStatus.healthy ? 'bg-green-500' : 'bg-red-500'"></div>
                      <span class="text-sm text-gray-600 dark:text-gray-400">健康状态: {{ meilisearchStatus.healthy ? '正常' : '异常' }}</span>
                    </div>
                    <div class="flex items-center space-x-2">
                      <i class="fas fa-sync-alt text-purple-500"></i>
                      <span class="text-sm text-gray-600 dark:text-gray-400">未同步: {{ unsyncedCount || 0 }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="transfer" tab="转存配置">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 自动转存 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">自动转存</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后，访问夸克资源，将自动转存，并提供转存后分享链接</span>
                </div>
                <n-switch v-model:value="configForm.auto_transfer_enabled" />
              </div>



              <!-- 最小存储空间 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">最小存储空间（GB）</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">当网盘剩余空间小于此值时，停止自动转存（100-1024GB）</span>
                </div>
                <n-input
                  v-model:value="configForm.auto_transfer_min_space"
                  type="text"
                  placeholder="500"
                  :disabled="!configForm.auto_transfer_enabled"
                />
              </div>

              <!-- 广告关键词 -->
              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">广告关键词</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">设置广告关键词，转存时，如果文件名包含广告关键词，则文件被删除</span>
                  </div>
                  <a 
                    href="https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/db/ad.txt" 
                    target="_blank" 
                    class="text-xs text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 underline"
                  >
                    开源广告关键词
                  </a>
                </div>
                <n-input
                  v-model:value="configForm.ad_keywords"
                  type="text"
                  placeholder="电影,电视剧,综艺"
                  :disabled="!configForm.auto_transfer_enabled"
                />
              </div>

              <!-- 自动插入广告 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">自动插入广告</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">在分享链接中的广告内容，会在转存时自动插入到转存文件夹</span>
                </div>
                <n-input
                  v-model:value="configForm.auto_insert_ad"
                  type="textarea"
                  placeholder="请输入广告内容..."
                  :rows="3"
                  :disabled="!configForm.auto_transfer_enabled"
                />
              </div>
            </div>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="drama" tab="热播剧">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 热播剧自动获取 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">自动拉取热播剧</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后，系统将自动从豆瓣获取热播剧信息</span>
                </div>
                <n-switch v-model:value="configForm.hot_drama_auto_fetch" />
              </div>
            </div>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNotification } from 'naive-ui'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'

// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

// 配置表单数据类型
interface FeatureConfigForm {
  auto_process_enabled: boolean
  auto_process_interval: string
  auto_transfer_enabled: boolean
  auto_transfer_min_space: string
  ad_keywords: string
  auto_insert_ad: string
  hot_drama_auto_fetch: boolean
  meilisearch_enabled: boolean
  meilisearch_host: string
  meilisearch_port: string
  meilisearch_master_key: string
  meilisearch_index_name: string
}

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  updateOriginalConfig,
  saveConfig: saveConfigWithDetection
} = useConfigChangeDetection<FeatureConfigForm>({
  debug: true,
  // 字段映射：前端字段名 -> 后端字段名
  fieldMapping: {
    auto_process_enabled: 'auto_process_ready_resources',
    auto_process_interval: 'auto_process_interval',
    auto_transfer_enabled: 'auto_transfer_enabled',
    auto_transfer_min_space: 'auto_transfer_min_space',
    ad_keywords: 'ad_keywords',
    auto_insert_ad: 'auto_insert_ad',
    hot_drama_auto_fetch: 'auto_fetch_hot_drama_enabled',
    meilisearch_enabled: 'meilisearch_enabled',
    meilisearch_host: 'meilisearch_host',
    meilisearch_port: 'meilisearch_port',
    meilisearch_master_key: 'meilisearch_master_key',
    meilisearch_index_name: 'meilisearch_index_name'
  }
})

const notification = useNotification()
const saving = ref(false)
const activeTab = ref('resource')
const testingConnection = ref(false)

// Meilisearch状态
const meilisearchStatus = ref<any>(null)
const unsyncedCount = ref(0)

// 配置表单数据
const configForm = ref<FeatureConfigForm>({
  auto_process_enabled: false,
  auto_process_interval: '30',
  auto_transfer_enabled: false,
  auto_transfer_min_space: '500',
  ad_keywords: '',
  auto_insert_ad: '',
  hot_drama_auto_fetch: false,
  meilisearch_enabled: false,
  meilisearch_host: '',
  meilisearch_port: '',
  meilisearch_master_key: '',
  meilisearch_index_name: ''
})

// 表单验证规则
const rules = {} as any

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig() as any
    
    if (response) {
      const configData = {
        auto_process_enabled: response.auto_process_ready_resources || false,
        auto_process_interval: String(response.auto_process_interval || 30),
        auto_transfer_enabled: response.auto_transfer_enabled || false,
        auto_transfer_min_space: String(response.auto_transfer_min_space || 500),
        ad_keywords: response.ad_keywords || '',
        auto_insert_ad: response.auto_insert_ad || '',
        hot_drama_auto_fetch: response.auto_fetch_hot_drama_enabled || false,
        meilisearch_enabled: response.meilisearch_enabled || false,
        meilisearch_host: response.meilisearch_host || '',
        meilisearch_port: String(response.meilisearch_port || 7700),
        meilisearch_master_key: response.meilisearch_master_key || '',
        meilisearch_index_name: response.meilisearch_index_name || 'resources'
      }
      
      configForm.value = { ...configData }
      setOriginalConfig(configData)
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
    notification.error({
      content: '获取系统配置失败',
      duration: 3000
    })
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    saving.value = true
    
    // 更新当前配置数据
    updateCurrentConfig({
      auto_process_enabled: configForm.value.auto_process_enabled,
      auto_process_interval: configForm.value.auto_process_interval,
      auto_transfer_enabled: configForm.value.auto_transfer_enabled,
      auto_transfer_min_space: configForm.value.auto_transfer_min_space,
      ad_keywords: configForm.value.ad_keywords,
      auto_insert_ad: configForm.value.auto_insert_ad,
      hot_drama_auto_fetch: configForm.value.hot_drama_auto_fetch,
      meilisearch_enabled: configForm.value.meilisearch_enabled,
      meilisearch_host: configForm.value.meilisearch_host,
      meilisearch_port: configForm.value.meilisearch_port,
      meilisearch_master_key: configForm.value.meilisearch_master_key,
      meilisearch_index_name: configForm.value.meilisearch_index_name
    })
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    // 使用通用保存函数
    const result = await saveConfigWithDetection(
      systemConfigApi.updateSystemConfig,
      {
        onlyChanged: true,
        includeAllFields: true,
        // 自定义数据转换
        transformSubmitData: (data) => {
          // 转换字符串为数字
          if (data.auto_process_interval !== undefined) {
            data.auto_process_interval = parseInt(data.auto_process_interval) || 30
          }
          if (data.auto_transfer_min_space !== undefined) {
            data.auto_transfer_min_space = parseInt(data.auto_transfer_min_space) || 500
          }
          return data
        }
      },
      // 成功回调
      async () => {
        notification.success({
          content: '功能配置保存成功',
          duration: 3000
        })
        
        // 刷新系统配置状态，确保顶部导航同步更新
        const { useSystemConfigStore } = await import('~/stores/systemConfig')
        const systemConfigStore = useSystemConfigStore()
        await systemConfigStore.initConfig(true, true)
      },
      // 错误回调
      (error) => {
        console.error('保存功能配置失败:', error)
        notification.error({
          content: '保存功能配置失败',
          duration: 3000
        })
      }
    )
    
    // 如果没有改动，显示提示
    if (result && result.message === '没有检测到任何改动') {
      notification.info({
        content: '没有检测到任何改动',
        duration: 3000
      })
    }
  } finally {
    saving.value = false
  }
}

// 测试Meilisearch连接
const testMeilisearchConnection = async () => {
  testingConnection.value = true
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    await meilisearchApi.testConnection({
      host: configForm.value.meilisearch_host,
      port: parseInt(configForm.value.meilisearch_port, 10),
      masterKey: configForm.value.meilisearch_master_key,
      indexName: configForm.value.meilisearch_index_name || 'resources'
    })
    notification.success({
      content: 'Meilisearch连接测试成功！',
      duration: 3000
    })
  } catch (error: any) {
    console.error('Meilisearch连接测试失败:', error)
    notification.error({
      content: `Meilisearch连接测试失败: ${error?.message || error}`,
      duration: 5000
    })
  } finally {
    testingConnection.value = false
  }
}

// 获取Meilisearch状态
const fetchMeilisearchStatus = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const status = await meilisearchApi.getStatus()
    meilisearchStatus.value = status
  } catch (error: any) {
    console.error('获取Meilisearch状态失败:', error)
    notification.error({
      content: `获取Meilisearch状态失败: ${error?.message || error}`,
      duration: 5000
    })
  }
}

// 获取未同步文档数量
const fetchUnsyncedCount = async () => {
  try {
    const { useMeilisearchApi } = await import('~/composables/useApi')
    const meilisearchApi = useMeilisearchApi()
    const response = await meilisearchApi.getUnsyncedCount() as any
    unsyncedCount.value = response?.count || 0
  } catch (error: any) {
    console.error('获取未同步文档数量失败:', error)
    notification.error({
      content: `获取未同步文档数量失败: ${error?.message || error}`,
      duration: 5000
    })
  }
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
  fetchMeilisearchStatus()
  fetchUnsyncedCount()
})


</script>

<style scoped>
/* 自定义样式 */
</style> 