<template>
  <n-select
    v-model:value="selectedValue"
    :placeholder="placeholder"
    :options="tagOptions"
    :loading="loading"
    :multiple="multiple"
    :clearable="clearable"
    :filterable="true"
    :remote="true"
    :clear-filter-after-select="false"
    @search="handleSearch"
    @focus="loadInitialTags"
    @update:value="handleUpdate"
  />
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useTagApi } from '~/composables/useApi'

// Props定义
interface Props {
  modelValue?: number | number[] | null
  placeholder?: string
  multiple?: boolean
  clearable?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '选择标签',
  multiple: false,
  clearable: true,
  disabled: false
})

// Emits定义
const emit = defineEmits<{
  'update:modelValue': [value: number | number[] | null]
}>()

// 定义选项类型
interface TagOption {
  label: string
  value: number
  disabled: boolean
}

// 内部状态
const selectedValue = ref(props.modelValue)
const tagOptions = ref<TagOption[]>([])
const loading = ref(false)
const searchCache = ref(new Map<string, TagOption[]>()) // 搜索缓存

// API实例
const tagApi = useTagApi()

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    selectedValue.value = newValue
  }
)

// 监听内部值变化并向外发射
const handleUpdate = (value: number | number[] | null) => {
  selectedValue.value = value
  emit('update:modelValue', value)
}

// 加载标签数据
const loadTags = async (query: string = '') => {
  // 检查缓存
  if (searchCache.value.has(query)) {
    const cachedOptions = searchCache.value.get(query)
    if (cachedOptions) {
      tagOptions.value = cachedOptions
      return
    }
  }

  loading.value = true
  try {
    const result = await tagApi.getTags({
      search: query,
      page: 1,
      page_size: 50 // 限制返回数量，避免数据过多
    }) as any
    
    const options: TagOption[] = []
    if (result && result.items) {
      options.push(...result.items.map((item: any) => ({
        label: item.name,
        value: item.id,
        disabled: false
      })))
    }
    
    // 缓存结果
    searchCache.value.set(query, options)
    tagOptions.value = options
  } catch (error) {
    console.error('获取标签失败:', error)
    tagOptions.value = []
  } finally {
    loading.value = false
  }
}

// 初始加载标签
const loadInitialTags = async () => {
  if (tagOptions.value.length === 0) {
    await loadTags('')
  }
}

// 搜索处理
const handleSearch = async (query: string) => {
  await loadTags(query)
}

// 组件挂载时预加载一些标签（可选）
onMounted(() => {
  // 可以选择在挂载时就加载标签，或者等用户聚焦时再加载
  // loadInitialTags()
})
</script>