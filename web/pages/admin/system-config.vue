<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">系统配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统配置和设置</p>
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
      <n-tabs type="line" animated>
        <!-- 站点配置 -->
        <n-tab-pane name="站点配置" tab="站点配置">
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
        </n-tab-pane>

        <!-- 功能配置 -->
        <n-tab-pane name="功能配置" tab="功能配置">
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
        </n-tab-pane>

        <!-- API配置 -->
        <n-tab-pane name="API配置" tab="API配置">
          <div class="space-y-6">
            <!-- API Token -->
            <div>
              <n-form-item label="公开API访问令牌" path="api_token">
                <div class="flex gap-2">
                  <n-input
                    v-model:value="configForm.api_token"
                    type="password"
                    placeholder="输入API Token，用于公开API访问认证"
                    show-password-on="click"
                  />
                  <n-button
                    v-if="!configForm.api_token"
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
                <template #help>
                  用于公开API的访问认证，建议使用随机字符串
                </template>
              </n-form-item>
            </div>

            <!-- API使用说明 -->
            <n-card>
              <template #header>
                <span class="text-lg font-semibold">API使用说明</span>
              </template>
              <div class="space-y-2 text-sm">
                <p><strong>批量添加资源:</strong> POST /api/public/resources/batch-add</p>
                <p><strong>资源搜索:</strong> GET /api/public/resources/search</p>
                <p><strong>热门剧:</strong> GET /api/public/hot-dramas</p>
              </div>
            </n-card>
          </div>
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- 系统状态 -->
    <n-card>
      <template #header>
        <span class="text-lg font-semibold">系统状态</span>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
          <div class="flex items-center">
            <i class="fas fa-server text-green-600 text-xl mr-3"></i>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">系统状态</p>
              <p class="text-lg font-semibold text-green-600">正常运行</p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 text-xl mr-3"></i>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">数据库</p>
              <p class="text-lg font-semibold text-blue-600">连接正常</p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
          <div class="flex items-center">
            <i class="fas fa-clock text-yellow-600 text-xl mr-3"></i>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">自动处理</p>
              <p class="text-lg font-semibold text-yellow-600">
                {{ configForm.auto_process_enabled ? '已启用' : '已禁用' }}
              </p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg">
          <div class="flex items-center">
            <i class="fas fa-exchange-alt text-purple-600 text-xl mr-3"></i>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">自动转存</p>
              <p class="text-lg font-semibold text-purple-600">
                {{ configForm.auto_transfer_enabled ? '已启用' : '已禁用' }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin' as any
})

// 获取运行时配置
const config = useRuntimeConfig()

// 使用API
const { useSystemConfigApi } = await import('~/composables/useApi')
const systemConfigApi = useSystemConfigApi()

// 响应式数据
const saving = ref(false)
const formRef = ref()

// 配置表单
const configForm = ref({
  site_title: '',
  site_description: '',
  keywords: '',
  copyright: '',
  maintenance_mode: false,
  api_token: '',
  auto_process_enabled: false,
  auto_process_interval: '30',
  auto_transfer_enabled: false,
  auto_transfer_limit_days: '30',
  auto_transfer_min_space: '500',
  forbidden_words: '',
  hot_drama_auto_fetch: false
})

// 表单验证规则
const rules = {
  site_title: {
    required: true,
    message: '请输入网站标题',
    trigger: 'blur'
  }
}

// 获取配置
const { data: configData, refresh: refreshConfig } = await useAsyncData(
  'systemConfig',
  () => systemConfigApi.getSystemConfig()
)

// 监听配置数据变化
watch(configData, (newData) => {
  if (newData && (newData as any)?.data) {
    configForm.value = {
      ...configForm.value,
      ...(newData as any).data
    }
  }
}, { immediate: true })

// 保存配置
const saveConfig = async () => {
  try {
    saving.value = true
    await systemConfigApi.updateSystemConfig(configForm.value)
    useNotification().success({
      content: '配置保存成功',
      duration: 3000
    })
    await refreshConfig()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '保存配置失败',
      duration: 5000
    })
  } finally {
    saving.value = false
  }
}

// 生成API Token
const generateApiToken = () => {
  const token = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
  configForm.value.api_token = token
  useNotification().success({
    content: 'API Token已生成',
    duration: 3000
  })
}

// 复制API Token
const copyApiToken = async () => {
  try {
    await navigator.clipboard.writeText(configForm.value.api_token)
    useNotification().success({
      content: 'API Token已复制到剪贴板',
      duration: 3000
    })
  } catch (error) {
    useNotification().error({
      content: '复制失败',
      duration: 3000
    })
  }
}
</script>

<style scoped>
/* 确保Font Awesome图标正确显示 */
.fas {
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
}
</style> 