<template>
  <section class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
    <header
      v-if="title || $slots.header"
      class="px-5 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between gap-3"
    >
      <div class="min-w-0">
        <h3 class="text-base font-semibold text-gray-900 dark:text-white truncate">
          <slot name="header">{{ title }}</slot>
        </h3>
        <p
          v-if="description"
          class="mt-1 text-sm text-gray-500 dark:text-gray-400"
        >
          {{ description }}
        </p>
      </div>
      <div v-if="$slots.actions" class="shrink-0">
        <slot name="actions" />
      </div>
    </header>
    <div class="px-5 py-4">
      <slot />
    </div>
    <footer
      v-if="$slots.footer"
      class="px-5 py-3 bg-gray-50 dark:bg-gray-900/40 border-t border-gray-200 dark:border-gray-700 flex items-center justify-end gap-2"
    >
      <slot name="footer" />
    </footer>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
}>()
</script>
