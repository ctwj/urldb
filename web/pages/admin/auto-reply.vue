<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">自动回复</h1>
        <p class="text-gray-600 dark:text-gray-400">管理各平台的自动回复配置</p>
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
        <n-tab-pane name="qq" tab="QQ机器人">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- QQ机器人配置占位符 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">QQ机器人开关</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启QQ机器人自动回复功能</span>
                </div>
                <n-switch v-model:value="configForm.qq_bot_enabled" />
              </div>

              <!-- 占位符内容 -->
              <div class="p-8 text-center text-gray-500 dark:text-gray-400">
                <i class="fas fa-cog text-4xl mb-4"></i>
                <p class="text-lg font-medium mb-2">QQ机器人配置</p>
                <p class="text-sm">QQ机器人自动回复功能配置区域</p>
                <p class="text-xs mt-2">具体配置项待开发...</p>
              </div>
            </div>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="wechat" tab="微信公众号">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 微信公众号配置占位符 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">微信公众号开关</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启微信公众号自动回复功能</span>
                </div>
                <n-switch v-model:value="configForm.wechat_mp_enabled" />
              </div>

              <!-- 占位符内容 -->
              <div class="p-8 text-center text-gray-500 dark:text-gray-400">
                <i class="fas fa-comment-dots text-4xl mb-4"></i>
                <p class="text-lg font-medium mb-2">微信公众号配置</p>
                <p class="text-sm">微信公众号自动回复功能配置区域</p>
                <p class="text-xs mt-2">具体配置项待开发...</p>
              </div>
            </div>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="wechat_open" tab="微信对话开放平台">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 微信对话开放平台配置占位符 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">微信对话开放平台开关</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启微信对话开放平台自动回复功能</span>
                </div>
                <n-switch v-model:value="configForm.wechat_open_enabled" />
              </div>

              <!-- 占位符内容 -->
              <div class="p-8 text-center text-gray-500 dark:text-gray-400">
                <i class="fas fa-comments text-4xl mb-4"></i>
                <p class="text-lg font-medium mb-2">微信对话开放平台配置</p>
                <p class="text-sm">微信对话开放平台自动回复功能配置区域</p>
                <p class="text-xs mt-2">具体配置项待开发...</p>
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
const activeTab = ref('qq')

// 配置表单数据
const configForm = ref({
  qq_bot_enabled: false,
  wechat_mp_enabled: false,
  wechat_open_enabled: false
})

// 表单验证规则
const rules = {
  // 暂时为空，后续添加验证规则
}

// 获取配置
const fetchConfig = async () => {
  try {
    // 暂时使用模拟数据
    configForm.value = {
      qq_bot_enabled: false,
      wechat_mp_enabled: false,
      wechat_open_enabled: false
    }
  } catch (error) {
    console.error('获取自动回复配置失败:', error)
    notification.error({
      content: '获取自动回复配置失败',
      duration: 3000
    })
  }
}

// 保存配置
const saveConfig = async () => {
  try {
    saving.value = true
    
    // 暂时使用模拟保存
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    notification.success({
      content: '自动回复配置保存成功',
      duration: 3000
    })
  } catch (error) {
    console.error('保存自动回复配置失败:', error)
    notification.error({
      content: '保存自动回复配置失败',
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