<template>
  <div class="flex flex-col items-center justify-center py-16 px-4 text-center">
    <div class="w-16 h-16 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center mb-4">
      <i :class="[icon, 'text-3xl text-gray-400 dark:text-gray-500']"></i>
    </div>
    <h3 class="text-base font-medium text-gray-900 dark:text-white mb-1">{{ title }}</h3>
    <p v-if="description" class="text-sm text-gray-500 dark:text-gray-400 mb-4 max-w-sm">{{ description }}</p>
    <NuxtLink
      v-if="ctaTo && ctaText"
      :to="ctaTo"
      class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-700 hover:bg-blue-800 dark:bg-blue-600 dark:hover:bg-blue-500 rounded-lg transition-colors cursor-pointer"
    >
      {{ ctaText }}
    </NuxtLink>
    <slot name="action" />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  icon: string
  title: string
  description?: string
  ctaText?: string
  ctaTo?: string
}>()
</script>
