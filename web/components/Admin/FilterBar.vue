<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-center">
      <!-- 搜索框（可选） -->
      <n-input
        v-if="config.search"
        :value="getValue(config.search.key)"
        :placeholder="config.search.placeholder"
        clearable
        @update:value="(v) => setValue(config.search!.key, v)"
        @keyup.enter="emit('search')"
      >
        <template #prefix>
          <i class="fas fa-search"></i>
        </template>
      </n-input>

      <!-- 下拉筛选 -->
      <n-select
        v-for="sel in config.selects"
        :key="sel.key"
        :value="getValue(sel.key)"
        :placeholder="sel.placeholder"
        :options="sel.options"
        clearable
        @update:value="(v) => setValue(sel.key, v)"
      />

      <!-- 操作按钮 -->
      <div class="flex gap-2">
        <n-button type="primary" @click="emit('search')">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
        <n-button @click="handleReset">
          <template #icon>
            <i class="fas fa-undo"></i>
          </template>
          重置
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * 统一筛选栏组件
 *
 * 通过 config 声明式配置搜索框与下拉筛选，
 * 父组件用 v-model 绑定整个 values 对象，FilterBar 通过 key 读写对应字段。
 * emits: search（触发查询）、reset（清空并触发查询）
 */

export interface FilterConfig {
  /** 搜索框配置（可选） */
  search?: { placeholder: string; key: string }
  /** 下拉筛选配置数组 */
  selects: Array<{
    key: string
    placeholder: string
    options: Array<{ label: string; value: string | number }>
  }>
}

const props = defineProps<{
  config: FilterConfig
  modelValue: Record<string, any>
}>()

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
  search: []
  reset: []
}>()

const getValue = (key: string) => props.modelValue[key] ?? null

const setValue = (key: string, value: any) => {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

const handleReset = () => {
  const cleared: Record<string, any> = {}
  if (props.config.search) cleared[props.config.search.key] = ''
  for (const sel of props.config.selects) cleared[sel.key] = null
  emit('update:modelValue', { ...props.modelValue, ...cleared })
  emit('reset')
}
</script>
