<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <!-- Bing索引配置 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Bing索引配置</h3>
        <p class="text-gray-600 dark:text-gray-400">配置Bing网站地图提交相关设置</p>
      </div>

      <!-- 功能开关 -->
      <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">Bing索引功能</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              开启后系统将在生成sitemap后自动提交到Bing搜索引擎
            </p>
          </div>
          <n-switch
            :value="bingIndexConfig.enabled"
            @update:value="updateBingIndexConfig"
            :loading="configLoading"
            size="large"
          >
            <template #checked>已开启</template>
            <template #unchecked>已关闭</template>
          </n-switch>
        </div>
      </div>

      <!-- API配置 -->
      <div class="space-y-6">
        <!-- API密钥 -->
        <div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <div class="mb-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">
              <i class="fas fa-key mr-2 text-blue-600"></i>
              Bing Webmaster API密钥
            </h4>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              配置Bing Webmaster API密钥以获得更好的提交效果。如果未配置，将使用传统的ping API。
            </p>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-2">
              <i class="fas fa-info-circle mr-1"></i>
              网站URL将使用系统配置中的网站URL
            </p>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                API密钥
              </label>
              <n-input
                v-model:value="localConfig.apiKey"
                type="password"
                show-password-on="click"
                placeholder="请输入Bing Webmaster API密钥"
                :disabled="!bingIndexConfig.enabled || configLoading"
              />
            </div>

            <div class="text-xs text-gray-500 dark:text-gray-400">
              <p>• 在 <a href="https://www.bing.com/webmasters/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Bing Webmaster Tools</a> 中获取API密钥</p>
              <p>• API密钥将自动提交sitemap到Bing索引</p>
              <p>• 如果不配置API密钥，系统将使用备用的ping API</p>
              <p>• 网站URL来源于系统配置中的网站URL设置</p>
            </div>
          </div>
        </div>

        <!-- 保存按钮 -->
        <div class="flex justify-end">
          <n-button
            type="primary"
            :loading="configLoading"
            :disabled="!bingIndexConfig.enabled"
            @click="saveConfig"
          >
            <template #icon>
              <i class="fas fa-save"></i>
            </template>
            保存配置
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

// Props
interface Props {
  bingIndexConfig: any
  configLoading: boolean
}

const props = withDefaults(defineProps<Props>(), {
  bingIndexConfig: () => ({}),
  configLoading: false
})

// Emits
const emit = defineEmits<{
  'update:bing-index-config': [value: boolean]
  'save-config': [config: { enabled: boolean; apiKey: string }]
}>()

// 本地配置状态
const localConfig = ref({
  apiKey: ''
})

// 监听props变化，更新本地配置
watch(() => props.bingIndexConfig, (newConfig) => {
  if (newConfig) {
    localConfig.value.apiKey = newConfig.apiKey || ''
  }
}, { immediate: true })

// 更新Bing索引配置
const updateBingIndexConfig = async (value: boolean) => {
  emit('update:bing-index-config', value)
}

// 保存配置
const saveConfig = () => {
  const config = {
    enabled: props.bingIndexConfig.enabled,
    apiKey: localConfig.value.apiKey
  }
  emit('save-config', config)
}
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>