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
      <div class="space-y-6">
        <!-- 自动处理 -->
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
            <n-switch v-model:value="configForm.auto_process_enabled" />
          </div>
        </div>

        <!-- 自动处理间隔 -->
        <div v-if="configForm.auto_process_enabled" class="ml-6">
          <n-form-item label="自动处理间隔 (分钟)" path="auto_process_interval">
            <n-input
              v-model:value="configForm.auto_process_interval"
              type="text"
              placeholder="30"
            />
            <template #help>
              建议设置 5-60 分钟，避免过于频繁的处理
            </template>
          </n-form-item>
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
            <n-switch v-model:value="configForm.auto_transfer_enabled" />
          </div>
        </div>

        <!-- 自动转存配置 -->
        <div v-if="configForm.auto_transfer_enabled" class="ml-6 space-y-4">
          <n-form-item label="自动转存限制（n天内资源）" path="auto_transfer_limit_days">
            <n-input
              v-model:value="configForm.auto_transfer_limit_days"
              type="text"
              placeholder="30"
            />
            <template #help>
              只转存指定天数内的资源，0表示不限制时间
            </template>
          </n-form-item>

          <n-form-item label="最小存储空间（GB）" path="auto_transfer_min_space">
            <n-input
              v-model:value="configForm.auto_transfer_min_space"
              type="text"
              placeholder="500"
            />
            <template #help>
              当网盘剩余空间小于此值时，停止自动转存（100-1024GB）
            </template>
          </n-form-item>
        </div>

        <!-- 热播剧自动获取 -->
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
            <n-switch v-model:value="configForm.hot_drama_auto_fetch" />
          </div>
        </div>
      </div>
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

// 配置表单数据
const configForm = ref({
  auto_process_enabled: false,
  auto_process_interval: 30,
  auto_transfer_enabled: false,
  auto_transfer_limit_days: 30,
  auto_transfer_min_space: 500,
  hot_drama_auto_fetch: false
})

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig()
    
    if (response) {
      configForm.value = {
        auto_process_enabled: response.auto_process_ready_resources || false,
        auto_process_interval: response.auto_process_interval || 30,
        auto_transfer_enabled: response.auto_transfer_enabled || false,
        auto_transfer_limit_days: response.auto_transfer_limit_days || 30,
        auto_transfer_min_space: response.auto_transfer_min_space || 500,
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
      auto_process_interval: configForm.value.auto_process_interval,
      auto_transfer_enabled: configForm.value.auto_transfer_enabled,
      auto_transfer_limit_days: configForm.value.auto_transfer_limit_days,
      auto_transfer_min_space: configForm.value.auto_transfer_min_space,
      auto_fetch_hot_drama_enabled: configForm.value.hot_drama_auto_fetch
    })
    
    notification.success({
      content: '功能配置保存成功',
      duration: 3000
    })
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

// 设置页面标题
useHead({
  title: '功能配置 - 老九网盘资源数据库'
})
</script>

<style scoped>
/* 自定义样式 */
</style> 