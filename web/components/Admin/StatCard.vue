<template>
  <NuxtLink
    :to="to"
    class="block bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-5 hover:shadow-md hover:border-blue-300 dark:hover:border-blue-600 transition-all cursor-pointer group"
  >
    <div class="flex items-start justify-between">
      <div class="flex items-center">
        <div
          class="p-3 rounded-lg"
          :class="iconBgClass"
        >
          <i :class="[icon, 'text-xl', iconTextClass]"></i>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ label }}</p>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">
            {{ formatNumber(displayValue, format) }}
          </p>
          <p v-if="total !== undefined" class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
            总计 {{ formatNumber(total, format) }}
          </p>
        </div>
      </div>
      <!-- 环比徽章 -->
      <div class="flex flex-col items-end">
        <span
          v-if="comparison.direction !== 'new' && comparison.percent !== null"
          class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="badgeClass"
        >
          <i :class="['mr-0.5', arrowIcon]"></i>
          {{ comparison.percent }}%
        </span>
        <span
          v-else-if="comparison.direction === 'new'"
          class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300"
        >
          新增
        </span>
        <span class="text-xs text-gray-400 dark:text-gray-500 mt-1">vs 昨日</span>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { computeComparison, formatNumber } from '~/utils/statCard'

const props = defineProps<{
  label: string
  icon: string
  today: number
  yesterday: number
  total?: number
  format?: 'integer' | 'compact'
  to: string
  /** 图标配色主题 */
  color?: 'blue' | 'orange' | 'green'
}>()

const format = computed(() => props.format ?? 'integer')
const color = computed(() => props.color ?? 'blue')

const displayValue = computed(() => props.today)
const comparison = computed(() => computeComparison(props.today, props.yesterday))

const iconBgClass = computed(() => {
  const map = {
    blue: 'bg-blue-100 dark:bg-blue-900/40',
    orange: 'bg-orange-100 dark:bg-orange-900/40',
    green: 'bg-green-100 dark:bg-green-900/40',
  }
  return map[color.value]
})

const iconTextClass = computed(() => {
  const map = {
    blue: 'text-blue-600 dark:text-blue-400',
    orange: 'text-orange-600 dark:text-orange-400',
    green: 'text-green-600 dark:text-green-400',
  }
  return map[color.value]
})

const badgeClass = computed(() => {
  const dir = comparison.value.direction
  if (dir === 'up') return 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300'
  if (dir === 'down') return 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'
  return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
})

const arrowIcon = computed(() => {
  const dir = comparison.value.direction
  if (dir === 'up') return 'fas fa-arrow-up'
  if (dir === 'down') return 'fas fa-arrow-down'
  return 'fas fa-minus'
})
</script>
