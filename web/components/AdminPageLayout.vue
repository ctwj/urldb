<template>
  <div class="h-full flex flex-col gap-3">
    <!-- 顶部标题和按钮区域 -->
    <div class="flex-0 w-full flex">
      <div v-if="isSubPage" class="flex-0 mr-4 flex items-center">
        <n-button @click="goBack" type="text" size="small">
          <template #icon>
            <i class="fas fa-arrow-left"></i>
          </template>
        </n-button>
      </div>
      <!-- 页面头部内容 -->
      <div class="flex-1 w-1 flex items-center justify-between">
        <slot name="page-header"></slot>
      </div>
    </div>

    <!-- 通知提示区域 -->
    <div v-if="hasNoticeSection" class="flex-shrink-0">
      <slot name="notice-section"></slot>
    </div>

    <!-- 过滤栏区域 -->
    <div v-if="hasFilterBar" class="flex-shrink-0">
      <slot name="filter-bar"></slot>
    </div>

    <!-- 内容区 - 自适应剩余高度 -->
    <div class="flex-1 bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col">
      <!-- 内容区header -->
      <div v-if="hasContentHeader" class="px-6 py-4 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-700 whitespace-nowrap">
        <slot name="content-header"></slot>
      </div>

      <!-- 内容区content - 自适应剩余高度 -->
      <div class="flex-1 h-1 content-slot overflow-hidden">
        <slot name="content"></slot>
      </div>

      <!-- 内容区footer -->
      <div v-if="hasContentFooter" class="flex-shrink-0 bg-gray-50 dark:bg-gray-700 border-t border-gray-200 dark:border-gray-700">
        <slot name="content-footer"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
const router = useRouter()
const $slots = useSlots()
const hasNoticeSection = computed(() => $slots['notice-section'] !== undefined)
const hasFilterBar = computed(() => $slots['filter-bar'] !== undefined)
const hasContentHeader = computed(() => $slots['content-header'] !== undefined)
const hasContentFooter = computed(() => $slots['content-footer'] !== undefined)

const goBack = () => {
  try {
    router.back()
  } catch (error) {
    navigateTo('/admin')
  }
}

defineProps({
  minHeight: {
    type: String,
    default: '400px'
  },
  isSubPage: {
    type: Boolean,
    default: false,
  }
})
</script>

<style scoped>
:deep(.content-slot) {
  min-height: 0;
}
</style>