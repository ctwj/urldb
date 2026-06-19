<template>
  <nav aria-label="Sidebar navigation" class="mt-8">
    <div class="px-4 space-y-6">
      <div v-for="group in groups" :key="group.key">
        <!-- dashboard：固定展示无折叠 -->
        <template v-if="group.key === 'dashboard'">
          <h3 class="px-4 mb-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
            {{ group.title }}
          </h3>
          <div class="space-y-1">
            <NuxtLink
              v-for="item in group.items"
              :key="item.to"
              :to="item.to"
              class="flex items-center px-4 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
              :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': isItemActive(item) }"
              @click="emit('navigate')"
            >
              <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
              <span>{{ item.label }}</span>
            </NuxtLink>
          </div>
        </template>

        <!-- 其他 4 组：可折叠 -->
        <template v-else>
          <button
            type="button"
            @click="toggle(group.key)"
            class="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider hover:text-gray-700 dark:hover:text-gray-300 transition-colors cursor-pointer"
          >
            <span>{{ group.title }}</span>
            <i
              class="fas fa-chevron-down text-xs transition-transform duration-200"
              :class="{ 'rotate-180': isExpanded(group.key) }"
            ></i>
          </button>
          <div v-show="isExpanded(group.key)" class="space-y-1 mt-2">
            <NuxtLink
              v-for="item in group.items"
              :key="item.to"
              :to="item.to"
              class="flex items-center px-8 py-3 text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 rounded-lg transition-colors"
              :class="{ 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400': isItemActive(item) }"
              @click="emit('navigate')"
            >
              <i :class="item.icon + ' w-5 h-5 mr-3'"></i>
              <span>{{ item.label }}</span>
            </NuxtLink>
          </div>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useStorage } from '@vueuse/core'
import { useAdminNav } from '~/composables/useAdminNav'
import {
  expandForGroup,
  type SidebarState,
  type SidebarGroupKey,
} from '~/utils/sidebarState'

const STORAGE_KEY = 'urldb:admin:sidebar:expanded-groups'

const emit = defineEmits<{
  navigate: []
}>()

const { groups, currentGroup, isItemActive } = useAdminNav()

// localStorage 持久化（@vueuse/core 自动 JSON 序列化，SSR 时回退内存）
const expanded = useStorage<SidebarState>(STORAGE_KEY, {})

const isExpanded = (key: SidebarGroupKey): boolean => Boolean(expanded.value[key])

const toggle = (key: SidebarGroupKey) => {
  expanded.value = { ...expanded.value, [key]: !isExpanded(key) }
}

// 进入路由时自动展开当前分组（覆盖存储值，保留其他偏好）
// currentGroup 为 null（dashboard）时返回 prev 拷贝，不引入副作用
watch(
  currentGroup,
  (g) => {
    expanded.value = expandForGroup(expanded.value, g)
  },
  { immediate: true },
)
</script>
