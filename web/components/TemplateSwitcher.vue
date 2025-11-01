<template>
  <div class="template-switcher">
    <div class="relative">
      <button
        @click="toggleDropdown"
        class="flex items-center space-x-2 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        aria-haspopup="true"
        :aria-expanded="dropdownOpen"
      >
        <i class="fas fa-paint-brush"></i>
        <span>模板: {{ currentTemplateName }}</span>
        <i class="fas fa-chevron-down transition-transform" :class="{ 'rotate-180': dropdownOpen }"></i>
      </button>

      <div
        v-show="dropdownOpen"
        class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-md shadow-lg z-50 border border-gray-200 dark:border-gray-700"
        role="menu"
      >
        <div class="py-1">
          <button
            v-for="template in availableTemplates"
            :key="template.id"
            @click="switchTemplate(template.id)"
            class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-between"
            :class="{ 'bg-blue-50 dark:bg-blue-900/30': currentTemplate === template.id }"
            role="menuitem"
          >
            <span>{{ template.name }}</span>
            <i
              v-if="currentTemplate === template.id"
              class="fas fa-check text-blue-600 dark:text-blue-400"
            ></i>
          </button>
        </div>

        <div class="border-t border-gray-200 dark:border-gray-700 py-2 px-4">
          <NuxtLink
            to="/admin/templates"
            class="text-sm text-blue-600 dark:text-blue-400 hover:underline flex items-center"
          >
            <i class="fas fa-cog mr-2"></i>
            模板管理
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTemplateManager } from '~/composables/templateManager'

const { currentTemplate, getAvailableTemplates, setCurrentTemplate } = useTemplateManager()

const dropdownOpen = ref(false)
const availableTemplates = computed(() => Object.values(getAvailableTemplates()))

const currentTemplateName = computed(() => {
  const template = getAvailableTemplates()[currentTemplate.value]
  return template ? template.name : '未知模板'
})

const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}

const switchTemplate = (templateId: string) => {
  setCurrentTemplate(templateId)
  dropdownOpen.value = false

  // 显示切换成功的提示
  const message = useMessage()
  const template = getAvailableTemplates()[templateId]
  if (template) {
    message.success(`已切换到 ${template.name} 模板，页面将刷新...`)
  }

  // 延迟刷新页面以应用新模板
  setTimeout(() => {
    window.location.reload()
  }, 1500)
}

// 点击外部关闭下拉菜单
const handleClickOutside = (event: MouseEvent) => {
  const element = document.querySelector('.template-switcher')
  if (element && !element.contains(event.target as Node)) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.template-switcher {
  position: relative;
}
</style>