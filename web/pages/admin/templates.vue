<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-6">
    <div class="max-w-7xl mx-auto">
      <div class="mb-6">
        <h1 class="text-2xl font-bold mb-2">模板管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理网站的模板配置</p>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-xl font-semibold mb-4">当前模板: {{ currentTemplate?.name || '未知' }}</h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div
            v-for="template in templates"
            :key="template.id"
            class="border rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <h3 class="font-medium text-lg mb-2">{{ template.name }}</h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm mb-4">{{ template.description }}</p>

            <div class="flex items-center justify-between">
              <span
                class="px-2 py-1 rounded text-xs"
                :class="template.enabled ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'"
              >
                {{ template.enabled ? '已启用' : '已禁用' }}
              </span>

              <button
                v-if="currentTemplateId !== template.id"
                @click="switchTemplate(template.id)"
                class="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:opacity-50"
                :disabled="switchingTemplate"
              >
                {{ switchingTemplate ? '切换中...' : '切换' }}
              </button>
              <span
                v-else
                class="px-3 py-1 bg-green-600 text-white rounded text-sm"
              >
                当前
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin-dynamic',
  middleware: ['auth']
})

const { data: templates, refresh } = await useAsyncData('templates', async () => {
  try {
    const response = await $fetch('/api/templates')
    return response.data
  } catch (error) {
    console.error('获取模板列表失败:', error)
    return []
  }
})

const { data: currentTemplateId } = await useAsyncData('current-template', async () => {
  const cookie = useCookie('site_template', { default: () => 'default' })
  return cookie.value
})

const currentTemplate = computed(() => {
  if (!templates.value || !currentTemplateId.value) return null
  return templates.value.find((t: any) => t.id === currentTemplateId.value)
})

const switchingTemplate = ref(false)

const switchTemplate = async (templateId: string) => {
  if (switchingTemplate.value) return

  switchingTemplate.value = true
  try {
    await $fetch('/api/templates/switch', {
      method: 'POST',
      body: { template_id: templateId }
    })

    // 刷新数据并提示用户
    await refresh()
    const message = useMessage()
    message.success('模板切换成功，页面将在2秒后刷新...')

    // 延迟刷新页面以应用新模板
    setTimeout(() => {
      window.location.reload()
    }, 2000)
  } catch (error: any) {
    console.error('切换模板失败:', error)
    const message = useMessage()
    message.error(error.message || '切换模板失败')
  } finally {
    switchingTemplate.value = false
  }
}
</script>