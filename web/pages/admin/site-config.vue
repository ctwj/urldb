<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">站点配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理网站基本信息和设置</p>
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
        <n-tab-pane name="basic" tab="基本信息">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 网站标题 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站标题</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站的主要标识，显示在浏览器标签页和搜索结果中</span>
                </div>
                <n-input
                  v-model:value="configForm.site_title"
                  placeholder="请输入网站标题"
                />
              </div>

              <!-- 网站描述 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站描述</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站的简要介绍，用于SEO和社交媒体分享</span>
                </div>
                <n-input
                  v-model:value="configForm.site_description"
                  placeholder="请输入网站描述"
                />
              </div>

              <!-- 关键词 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">关键词</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">用于SEO优化，多个关键词用逗号分隔</span>
                </div>
                <n-input
                  v-model:value="configForm.keywords"
                  placeholder="请输入关键词，用逗号分隔"
                />
              </div>

              <!-- 版权信息 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">版权信息</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站底部的版权声明信息</span>
                </div>
                <n-input
                  v-model:value="configForm.copyright"
                  placeholder="请输入版权信息"
                />
              </div>
            </div>
          </n-form>
        </n-tab-pane>



        <n-tab-pane name="security" tab="安全设置">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 维护模式 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">维护模式</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后网站将显示维护页面，暂停用户访问</span>
                </div>
                <n-switch v-model:value="configForm.maintenance_mode" />
              </div>

              <!-- 违禁词 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">违禁词</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">包含这些词汇的资源将被过滤，多个词汇用逗号分隔</span>
                </div>
                <n-input
                  v-model:value="configForm.forbidden_words"
                  placeholder="请输入违禁词，用逗号分隔"
                  type="textarea"
                  :rows="4"
                />
              </div>

              <!-- 开启注册 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">开启注册</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后用户才能注册新账号，关闭后注册页面将显示"当前系统已关闭注册功能"</span>
                </div>
                <n-switch v-model:value="configForm.enable_register" />
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
const formRef = ref()
const saving = ref(false)
const activeTab = ref('basic')

// 配置表单数据
const configForm = ref<{
  site_title: string
  site_description: string
  keywords: string
  copyright: string
  maintenance_mode: boolean
  enable_register: boolean
  forbidden_words: string
  enable_sitemap: boolean
  sitemap_update_frequency: string
}>({
  site_title: '',
  site_description: '',
  keywords: '',
  copyright: '',
  maintenance_mode: false,
  enable_register: false, // 新增：开启注册开关
  forbidden_words: '',
  enable_sitemap: false,
  sitemap_update_frequency: 'daily'
})



// 表单验证规则
const rules = {
  site_title: {
    required: true,
    message: '请输入网站标题',
    trigger: 'blur'
  },
  site_description: {
    required: true,
    message: '请输入网站描述',
    trigger: 'blur'
  }
}

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig() as any
    
    if (response) {
      configForm.value = {
        site_title: response.site_title || '',
        site_description: response.site_description || '',
        keywords: response.keywords || '',
        copyright: response.copyright || '',
        maintenance_mode: response.maintenance_mode || false,
        enable_register: response.enable_register || false, // 新增：获取开启注册开关
        forbidden_words: response.forbidden_words || '',
        enable_sitemap: response.enable_sitemap || false,
        sitemap_update_frequency: response.sitemap_update_frequency || 'daily'
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
    await formRef.value?.validate()
    saving.value = true
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    await systemConfigApi.updateSystemConfig({
      site_title: configForm.value.site_title,
      site_description: configForm.value.site_description,
      keywords: configForm.value.keywords,
      copyright: configForm.value.copyright,
      maintenance_mode: configForm.value.maintenance_mode,
      enable_register: configForm.value.enable_register, // 新增：保存开启注册开关
      forbidden_words: configForm.value.forbidden_words,
      enable_sitemap: configForm.value.enable_sitemap,
      sitemap_update_frequency: configForm.value.sitemap_update_frequency
    })
    
    notification.success({
      content: '站点配置保存成功',
      duration: 3000
    })
  } catch (error) {
    console.error('保存站点配置失败:', error)
    notification.error({
      content: '保存站点配置失败',
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

// 设置页面标题
useHead({
  title: '站点配置 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 自定义样式 */
</style> 