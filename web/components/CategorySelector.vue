<template>
  <n-select
    v-model:value="selectedValue"
    :placeholder="placeholder"
    :options="categoryOptions"
    :loading="loading"
    :clearable="clearable"
    :filterable="true"
    :disabled="disabled"
    @update:value="handleUpdate"
  />
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useCategoryApi } from '~/composables/useApi'

// Props定义
interface Props {
  modelValue?: number | null
  placeholder?: string
  clearable?: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '选择分类',
  clearable: true,
  disabled: false
})

// Emits定义
const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

// 定义选项类型
interface CategoryOption {
  label: string
  value: number
  disabled: boolean
}

// 内部状态
const selectedValue = ref(props.modelValue)
const categoryOptions = ref<CategoryOption[]>([])
const loading = ref(false)

// API实例
const categoryApi = useCategoryApi()

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    selectedValue.value = newValue
  }
)

// 监听内部值变化并向外发射
const handleUpdate = (value: number | null) => {
  selectedValue.value = value
  emit('update:modelValue', value)
}

// 加载分类数据
const loadCategories = async () => {
  // 如果已经加载过，直接返回
  if (categoryOptions.value.length > 0) {
    return
  }

  loading.value = true
  try {
    const result = await categoryApi.getCategories() as any
    
    const options: CategoryOption[] = []
    if (result && result.items) {
      options.push(...result.items.map((item: any) => ({
        label: item.name,
        value: item.id,
        disabled: false
      })))
    } else if (Array.isArray(result)) {
      options.push(...result.map((item: any) => ({
        label: item.name,
        value: item.id,
        disabled: false
      })))
    }
    
    categoryOptions.value = options
  } catch (error) {
    console.error('获取分类失败:', error)
    categoryOptions.value = []
  } finally {
    loading.value = false
  }
}

// 组件挂载时立即加载分类
onMounted(() => {
  loadCategories()
})
</script>