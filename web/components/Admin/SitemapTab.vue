<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <!-- Sitemap配置 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Sitemap配置</h3>
        <p class="text-gray-600 dark:text-gray-400">管理网站的Sitemap生成和配置</p>
      </div>

      <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">自动生成Sitemap</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              开启后系统将定期自动生成Sitemap文件
            </p>
          </div>
          <n-switch
            v-model:value="sitemapConfig.autoGenerate"
            @update:value="updateSitemapConfig"
            :loading="configLoading"
            size="large"
          >
            <template #checked>已开启</template>
            <template #unchecked>已关闭</template>
          </n-switch>
        </div>

        <!-- 配置详情 -->
        <div class="border-t border-gray-200 dark:border-gray-600 pt-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">站点URL</label>
              <n-input
                :value="systemConfig?.site_url || '站点URL未配置'"
                :disabled="true"
                placeholder="请先在站点配置中设置站点URL"
              >
                <template #prefix>
                  <i class="fas fa-globe text-gray-400"></i>
                </template>
              </n-input>
            </div>
            <div class="flex flex-col justify-end">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">最后生成时间</label>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                {{ sitemapConfig.lastGenerate || '尚未生成' }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Sitemap统计 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
          <div class="flex items-center">
            <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
              <i class="fas fa-link text-blue-600 dark:text-blue-400"></i>
            </div>
            <div class="ml-3">
              <p class="text-sm text-gray-600 dark:text-gray-400">资源总数</p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.total_resources || 0 }}</p>
            </div>
          </div>
        </div>

        <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
          <div class="flex items-center">
            <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
              <i class="fas fa-sitemap text-green-600 dark:text-green-400"></i>
            </div>
            <div class="ml-3">
              <p class="text-sm text-gray-600 dark:text-gray-400">页面数量</p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.total_pages || 0 }}</p>
            </div>
          </div>
        </div>

        <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
          <div class="flex items-center">
            <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
              <i class="fas fa-history text-purple-600 dark:text-purple-400"></i>
            </div>
            <div class="ml-3">
              <p class="text-sm text-gray-600 dark:text-gray-400">最后更新</p>
              <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.last_generate || 'N/A' }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex flex-wrap gap-3 mb-6">
        <n-button
          type="primary"
          @click="generateSitemap"
          :loading="isGenerating"
          size="large"
        >
          <template #icon>
            <i class="fas fa-cog"></i>
          </template>
          生成Sitemap
        </n-button>

        <n-button
          type="success"
          @click="viewSitemap"
          size="large"
        >
          <template #icon>
            <i class="fas fa-external-link-alt"></i>
          </template>
          查看Sitemap
        </n-button>

        <n-button
          type="info"
          @click="$emit('refresh-status')"
          size="large"
        >
          <template #icon>
            <i class="fas fa-sync-alt"></i>
          </template>
          刷新状态
        </n-button>
      </div>

      <!-- 生成状态 -->
      <div v-if="generateStatus" class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
        <div class="flex items-center">
          <i class="fas fa-info-circle text-blue-500 dark:text-blue-400 mr-2"></i>
          <span class="text-blue-700 dark:text-blue-300">{{ generateStatus }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { useApi } from '~/composables/useApi'
import { ref } from 'vue'

// Props
interface Props {
  systemConfig?: any
  sitemapConfig: any
  sitemapStats: any
  configLoading: boolean
  isGenerating: boolean
  generateStatus: string
}

const props = withDefaults(defineProps<Props>(), {
  systemConfig: null,
  sitemapConfig: () => ({}),
  sitemapStats: () => ({}),
  configLoading: false,
  isGenerating: false,
  generateStatus: ''
})

// Emits
const emit = defineEmits<{
  'update:sitemap-config': [value: boolean]
  'refresh-status': []
}>()

// 获取消息组件
const message = useMessage()

// 更新Sitemap配置
const updateSitemapConfig = async (value: boolean) => {
  try {
    const api = useApi()
    await api.sitemapApi.updateSitemapConfig({
      autoGenerate: value,
      lastGenerate: props.sitemapConfig.lastGenerate,
      lastUpdate: new Date().toISOString()
    })
    message.success(value ? '自动生成功能已开启' : '自动生成功能已关闭')
  } catch (error) {
    message.error('更新配置失败')
    // 恢复之前的值
    props.sitemapConfig.autoGenerate = !value
  }
}

// 生成Sitemap
const generateSitemap = async () => {
  // 使用已经加载的系统配置
  const siteUrl = props.systemConfig?.site_url || ''
  if (!siteUrl) {
    message.warning('请先在站点配置中设置站点URL，然后再生成Sitemap')
    return
  }

  try {
    const api = useApi()
    const response = await api.sitemapApi.generateSitemap({ site_url: siteUrl })

    if (response) {
      message.success(`Sitemap生成任务已启动，使用站点URL: ${siteUrl}`)
      // 更新统计信息
      emit('refresh-status')
    }
  } catch (error: any) {
    message.error('Sitemap生成失败')
  }
}

// 查看Sitemap
const viewSitemap = () => {
  window.open('/sitemap.xml', '_blank')
}
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>