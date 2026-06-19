<template>
  <div class="flex flex-col items-center justify-center py-16 px-4 text-center">
    <div class="w-16 h-16 rounded-full bg-red-50 dark:bg-red-900/20 flex items-center justify-center mb-4">
      <i :class="[icon, 'text-3xl text-red-500 dark:text-red-400']"></i>
    </div>
    <p class="text-sm text-gray-600 dark:text-gray-400 mb-4 max-w-sm">{{ message }}</p>
    <button
      v-if="onRetry"
      @click="onRetry"
      class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-700 hover:bg-blue-800 dark:bg-blue-600 dark:hover:bg-blue-500 rounded-lg transition-colors cursor-pointer"
    >
      <i class="fas fa-redo mr-2"></i>重试
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  icon: string
  message: string
  onRetry?: () => void
}>()
</script>
