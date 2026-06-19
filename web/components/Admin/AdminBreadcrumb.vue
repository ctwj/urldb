<template>
  <nav aria-label="Breadcrumb" class="flex items-center text-sm">
    <!-- 第一级：管理后台（首页） -->
    <NuxtLink
      to="/admin"
      class="text-gray-500 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
    >
      管理后台
    </NuxtLink>

    <!-- 第二级：分组标题（dashboard 无分组时跳过） -->
    <template v-if="currentGroupTitle">
      <i class="fas fa-chevron-right mx-2 text-xs text-gray-300 dark:text-gray-600"></i>
      <span class="text-gray-500 dark:text-gray-400">{{ currentGroupTitle }}</span>
    </template>

    <!-- 第三级：当前页面 -->
    <template v-if="currentItemLabel">
      <i class="fas fa-chevron-right mx-2 text-xs text-gray-300 dark:text-gray-600"></i>
      <span class="text-gray-900 dark:text-white font-medium">{{ currentItemLabel }}</span>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAdminNav } from '~/composables/useAdminNav'

const { groups, currentGroup, isItemActive } = useAdminNav()

// 当前分组标题
const currentGroupTitle = computed(() => {
  const g = currentGroup.value
  if (!g) return ''
  return groups.find((x) => x.key === g)?.title ?? ''
})

// 当前活跃 nav item 的 label
const currentItemLabel = computed(() => {
  for (const g of groups) {
    const hit = g.items.find((item) => isItemActive(item))
    if (hit) return hit.label
  }
  return ''
})
</script>
