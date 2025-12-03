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
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">站点URL</label>
            <n-input
              :value="systemConfig?.website_url || '站点URL未配置'"
              :disabled="true"
              placeholder="请先在站点配置中设置站点URL"
            >
              <template #prefix>
                <i class="fas fa-globe text-gray-400"></i>
              </template>
            </n-input>
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
              <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatLastGenerate(sitemapStats.last_generate) }}</p>
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

      <!-- Sitemap文件列表 -->
      <div v-if="sitemapStats.files && sitemapStats.files.length > 0" class="mb-4">
        <h4 class="text-md font-medium text-gray-900 dark:text-white mb-3">
          <i class="fas fa-folder-open mr-2"></i>
          Sitemap 文件列表 ({{ sitemapStats.file_count || sitemapStats.files.length }} 个文件)
        </h4>
        <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-600">
            <thead class="bg-gray-100 dark:bg-gray-700">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">文件名</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">大小</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">最后修改</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
              <tr v-for="file in sitemapStats.files" :key="file.name" class="hover:bg-gray-100 dark:hover:bg-gray-600">
                <td class="px-4 py-2 text-sm text-gray-900 dark:text-white">
                  <i class="fas fa-file-code text-green-500 mr-2"></i>
                  {{ file.name }}
                </td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ file.size_format }}</td>
                <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ file.mod_time }}</td>
                <td class="px-4 py-2 text-sm">
                  <a 
                    :href="'/' + file.name" 
                    target="_blank" 
                    class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                  >
                    <i class="fas fa-external-link-alt mr-1"></i>查看
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 无文件提示 -->
      <div v-else class="mb-4 p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800">
        <div class="flex items-center">
          <i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-2"></i>
          <span class="text-yellow-700 dark:text-yellow-300">尚未生成 Sitemap 文件，请点击"生成Sitemap"按钮</span>
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
    emit('update:sitemap-config', value)
  } catch (error) {
    message.error('更新配置失败')
  }
}

// 生成Sitemap
const generateSitemap = async () => {
  // 使用已经加载的系统配置
  const siteUrl = props.systemConfig?.website_url || ''
  if (!siteUrl) {
    message.warning('请先在站点配置中设置站点URL，然后再生成Sitemap')
    return
  }

  try {
    const api = useApi()
    const response = await api.sitemapApi.generateSitemap({ website_url: siteUrl })

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

// 格式化最后生成时间
const formatLastGenerate = (time: string | null | undefined): string => {
  if (!time || time === '' || time.startsWith('0001-01-01')) {
    return '尚未生成'
  }
  return time
}
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>