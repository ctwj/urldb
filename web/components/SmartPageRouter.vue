<template>
  <div class="smart-page-router">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600 dark:text-gray-400">加载页面中...</span>
    </div>
    <div v-else-if="error" class="text-center py-12">
      <div class="text-red-600 dark:text-red-400 mb-4">
        <i class="fas fa-exclamation-triangle text-2xl"></i>
      </div>
      <p class="text-gray-800 dark:text-gray-200">{{ error }}</p>
      <button
        @click="retryLoad"
        class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
      >
        重试
      </button>
    </div>
    <component
      :is="pageComponent"
      v-else-if="pageComponent"
      class="template-page-content"
    />
    <div v-else class="text-center py-12">
      <p class="text-gray-600 dark:text-gray-400">页面内容为空</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTemplatePageLoader } from '~/composables/templatePageLoader'
import { useTemplateManager } from '~/composables/templateManager'
import { useRoute } from '#app'

const props = defineProps({
  pagePath: {
    type: String,
    default: null
  }
})

const route = useRoute()
const { loadTemplatePage } = useTemplatePageLoader()
const { currentTemplate } = useTemplateManager()

// 响应式状态
const loading = ref(true)
const error = ref<string | null>(null)
const pageComponent = ref<any>(null)

// 计算实际页面路径
const actualPagePath = computed(() => {
  if (props.pagePath) {
    return props.pagePath
  }
  return route.path
})

// 加载页面
const loadPage = async () => {
  loading.value = true
  error.value = null

  try {
    pageComponent.value = await loadTemplatePage(actualPagePath.value)
  } catch (err: any) {
    error.value = err.message || '加载页面失败'
    console.error('加载页面失败:', err)
  } finally {
    loading.value = false
  }
}

// 重试加载
const retryLoad = () => {
  loadPage()
}

// 监听模板和路由变化
watch([currentTemplate, () => route.path], () => {
  loadPage()
}, { immediate: true })

// 组件挂载时加载页面
onMounted(() => {
  loadPage()
})
</script>

<style scoped>
.smart-page-router {
  width: 100%;
  min-height: 200px;
}

.template-page-content {
  width: 100%;
}
</style>