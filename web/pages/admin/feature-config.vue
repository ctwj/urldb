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
            <div class="space-y-6">
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
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">广告关键词</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">设置广告关键词，转存时，如果文件名包含广告关键词，则文件被删除</span>
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
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

const notification = useNotification()
const saving = ref(false)
const activeTab = ref('resource')

// 配置表单数据
const configForm = ref({
  auto_process_enabled: false,
  auto_process_interval: '30',
  auto_transfer_enabled: false,
  auto_transfer_min_space: '500',
  ad_keywords: '',
  auto_insert_ad: '',
  hot_drama_auto_fetch: false
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
      configForm.value = {
        auto_process_enabled: response.auto_process_ready_resources || false,
        auto_process_interval: String(response.auto_process_interval || 30),
        auto_transfer_enabled: response.auto_transfer_enabled || false,
        auto_transfer_min_space: String(response.auto_transfer_min_space || 500),
        ad_keywords: response.ad_keywords || '',
        auto_insert_ad: response.auto_insert_ad || '',
        hot_drama_auto_fetch: response.auto_fetch_hot_drama_enabled || false
      }
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
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    await systemConfigApi.updateSystemConfig({
      auto_process_ready_resources: configForm.value.auto_process_enabled,
      auto_process_interval: parseInt(configForm.value.auto_process_interval) || 30,
      auto_transfer_enabled: configForm.value.auto_transfer_enabled,
      auto_transfer_min_space: parseInt(configForm.value.auto_transfer_min_space) || 500,
      ad_keywords: configForm.value.ad_keywords,
      auto_insert_ad: configForm.value.auto_insert_ad,
      auto_fetch_hot_drama_enabled: configForm.value.hot_drama_auto_fetch
    })
    
    notification.success({
      content: '功能配置保存成功',
      duration: 3000
    })
    
    // 刷新系统配置状态，确保顶部导航同步更新
    const { useSystemConfigStore } = await import('~/stores/systemConfig')
    const systemConfigStore = useSystemConfigStore()
    await systemConfigStore.initConfig(true, true) // 强制刷新，使用管理员API
  } catch (error) {
    console.error('保存功能配置失败:', error)
    notification.error({
      content: '保存功能配置失败',
      duration: 3000
    })
  } finally {
    saving.value = false
  }
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
})


</script>

<style scoped>
/* 自定义样式 */
</style> 