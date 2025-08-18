<template>
  <div class="p-6 space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">三方统计</h1>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">配置第三方统计代码，监控网站访问数据</p>
      </div>
    </div>

    <!-- 统计代码配置 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
      <div class="p-4 md:p-6 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">统计代码配置</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">配置第三方统计代码，代码将自动加载到网站首页</p>
      </div>
      
      <div class="p-4 md:p-6">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              统计代码
            </label>
            <n-input
              v-model:value="statsCode"
              type="textarea"
              placeholder="请输入第三方统计代码（如百度统计、Google Analytics等）"
              :autosize="{ minRows: 6, maxRows: 12 }"
              class="font-mono text-sm"
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
              支持HTML、JavaScript代码，代码将自动插入到网站首页的 &lt;/body&gt; 标签之前
            </p>
          </div>
          
          <div class="flex items-center justify-end space-x-3">
            <n-button
              @click="resetCode"
              size="small"
            >
              重置
            </n-button>
            <n-button
              @click="saveCode"
              type="primary"
              size="small"
              :loading="saving"
            >
              保存配置
            </n-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 常用三方统计列表 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm">
      <div class="p-4 md:p-6 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">常用三方统计</h2>
        <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">常用的第三方统计服务及其配置说明</p>
      </div>
      
      <div class="p-4 md:p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <!-- 微软统计 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
                <i class="fas fa-chart-bar text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">微软统计</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">免费网站分析服务</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              微软统计（Microsoft Clarity）是微软提供的免费网站分析工具，提供热力图、用户行为录制、性能分析等功能。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://clarity.microsoft.com" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://clarity.microsoft.com
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">免费、热力图、用户录制、性能分析</span>
              </div>
            </div>
          </div>

          <!-- Umami -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-green-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-leaf text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">Umami</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">开源隐私友好统计</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              Umami 是一个开源的、注重隐私的网站分析工具，可以自托管，不收集个人数据，符合GDPR要求。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://umami.is" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://umami.is
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">开源、隐私友好、可自托管、轻量级</span>
              </div>
            </div>
          </div>

          <!-- 百度统计 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-chart-line text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">百度统计</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">免费网站统计服务</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              百度统计是百度推出的一款免费的专业网站流量分析工具，提供网站流量分析、用户行为分析等功能。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://tongji.baidu.com" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://tongji.baidu.com
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">免费、中文界面、数据准确</span>
              </div>
            </div>
          </div>

          <!-- Google Analytics -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-red-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-chart-bar text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">Google Analytics</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">全球领先的网站分析工具</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              Google Analytics 是 Google 提供的免费网站分析服务，功能强大，数据详细，支持多维度分析。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://analytics.google.com" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://analytics.google.com
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">功能强大、数据详细、国际化</span>
              </div>
            </div>
          </div>

          <!-- 51.la统计 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-purple-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-chart-area text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">51.la统计</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">轻量级网站统计</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              51.la 统计是一款轻量级的网站统计工具，代码简洁，加载速度快，适合小型网站使用。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://www.51.la" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://www.51.la
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">轻量级、加载快、简单易用</span>
              </div>
            </div>
          </div>

          <!-- 腾讯云分析 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
                <i class="fas fa-cloud text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">腾讯云分析</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">免费网站分析服务</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              腾讯云分析是腾讯云提供的免费网站分析服务，与腾讯云生态深度集成，提供全面的数据分析能力。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://cloud.tencent.com/product/ta" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://cloud.tencent.com/product/ta
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">免费、云生态集成、数据安全</span>
              </div>
            </div>
          </div>

          <!-- 阿里云日志服务 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-orange-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-file-alt text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">阿里云日志服务</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">免费日志分析服务</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              阿里云日志服务提供免费的日志收集、存储、查询和分析功能，适合需要深度数据分析的网站。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://www.aliyun.com/product/sls" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://www.aliyun.com/product/sls
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">免费额度、实时分析、高可用</span>
              </div>
            </div>
          </div>

          <!-- 神策数据 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-brain text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">神策数据</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">免费用户行为分析</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              神策数据提供免费的用户行为分析服务，支持精细化运营分析，数据准确度高。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">官网：</span>
                <a href="https://www.sensorsdata.cn" target="_blank" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                  https://www.sensorsdata.cn
                </a>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">免费版、数据准确、行为分析</span>
              </div>
            </div>
          </div>

          <!-- 自定义统计 -->
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-8 h-8 bg-gray-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-cog text-white text-sm"></i>
              </div>
              <div>
                <h3 class="font-medium text-gray-900 dark:text-white">自定义统计</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">自建统计系统</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              如果您有自己的统计系统或使用其他第三方统计服务，可以直接将统计代码粘贴到上方配置框中。
            </p>
            <div class="space-y-2">
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">支持：</span>
                <span class="text-gray-600 dark:text-gray-400">HTML、JavaScript、CSS</span>
              </div>
              <div class="text-xs">
                <span class="font-medium text-gray-700 dark:text-gray-300">特点：</span>
                <span class="text-gray-600 dark:text-gray-400">灵活配置、完全控制</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'

// 页面配置
definePageMeta({
  layout: 'admin'
})

// 定义配置表单类型
interface ThirdPartyStatsForm {
  third_party_stats_code: string
}

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  updateOriginalConfig,
  saveConfig: saveConfigWithDetection
} = useConfigChangeDetection<ThirdPartyStatsForm>({
  debug: true,
  fieldMapping: {
    third_party_stats_code: 'third_party_stats_code'
  }
})

// 状态管理
const message = useMessage()
const statsCode = ref('')
const saving = ref(false)

// 获取当前配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig()
    
    if (response) {
      const configData = {
        third_party_stats_code: (response as any).third_party_stats_code || ''
      }
      
      statsCode.value = configData.third_party_stats_code
      setOriginalConfig(configData)
    }
  } catch (error) {
    console.error('获取配置失败:', error)
    message.error('获取配置失败')
  }
}

// 保存配置
const saveCode = async () => {
  try {
    saving.value = true
    
    // 更新当前配置
    updateCurrentConfig({
      third_party_stats_code: statsCode.value
    })
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    // 使用通用保存函数
    const result = await saveConfigWithDetection(
      systemConfigApi.updateSystemConfig,
      {
        onlyChanged: true,
        includeAllFields: true
      },
      // 成功回调
      () => {
        message.success('配置保存成功')
      },
      // 错误回调
      (error) => {
        console.error('保存配置失败:', error)
        message.error('保存配置失败')
      }
    )
    
    // 如果没有改动，显示提示
    if (result && result.message === '没有检测到任何改动') {
      message.info('没有检测到任何改动')
    }
  } finally {
    saving.value = false
  }
}

// 重置代码
const resetCode = () => {
  statsCode.value = ''
  message.info('已重置统计代码')
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
})
</script> 