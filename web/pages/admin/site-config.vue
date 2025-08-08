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
      <n-form
        ref="formRef"
        :model="configForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 网站标题 -->
          <n-form-item label="网站标题" path="site_title">
            <n-input
              v-model:value="configForm.site_title"
              placeholder="请输入网站标题"
            />
          </n-form-item>

          <!-- 网站描述 -->
          <n-form-item label="网站描述" path="site_description">
            <n-input
              v-model:value="configForm.site_description"
              placeholder="请输入网站描述"
            />
          </n-form-item>

          <!-- 关键词 -->
          <n-form-item label="关键词" path="keywords">
            <n-input
              v-model:value="configForm.keywords"
              placeholder="请输入关键词，用逗号分隔"
            />
          </n-form-item>

          <!-- 版权信息 -->
          <n-form-item label="版权信息" path="copyright">
            <n-input
              v-model:value="configForm.copyright"
              placeholder="请输入版权信息"
            />
          </n-form-item>

          <!-- 维护模式 -->
          <n-form-item label="维护模式" path="maintenance_mode">
            <n-switch v-model:value="configForm.maintenance_mode" />
            <template #help>
              开启后网站将显示维护页面
            </template>
          </n-form-item>

          <!-- 开启注册 -->
          <n-form-item label="开启注册" path="enable_register">
            <n-switch v-model:value="configForm.enable_register" />
            <template #help>
              开启后用户才能注册新账号，关闭后注册页面将显示"当前系统已关闭注册功能"
            </template>
          </n-form-item>

          <!-- 违禁词 -->
          <n-form-item label="违禁词" path="forbidden_words" class="md:col-span-2">
            <n-input
              v-model:value="configForm.forbidden_words"
              placeholder="请输入违禁词，用逗号分隔"
              type="textarea"
              :rows="4"
            />
            <template #help>
              包含这些词汇的资源将被过滤
            </template>
          </n-form-item>
        </div>
      </n-form>
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

// 配置表单数据
const configForm = ref<{
  site_title: string
  site_description: string
  keywords: string
  copyright: string
  maintenance_mode: boolean
  enable_register: boolean
  forbidden_words: string
}>({
  site_title: '',
  site_description: '',
  keywords: '',
  copyright: '',
  maintenance_mode: false,
  enable_register: false, // 新增：开启注册开关
  forbidden_words: ''
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
        forbidden_words: response.forbidden_words || ''
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
      forbidden_words: configForm.value.forbidden_words
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