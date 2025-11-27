<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <!-- Bing索引配置 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Bing索引配置</h3>
        <p class="text-gray-600 dark:text-gray-400">配置Bing网站地图提交相关设置</p>
      </div>

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
    </div>
  </div>
</template>

<script setup lang="ts">
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
}>()

// 更新Bing索引配置
const updateBingIndexConfig = async (value: boolean) => {
  emit('update:bing-index-config', value)
}
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>