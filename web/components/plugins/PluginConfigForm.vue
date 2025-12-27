<template>
  <div class="plugin-config-form">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="120px"
      require-mark-placement="right-hanging"
    >
      <div v-if="!configFields || Object.keys(configFields).length === 0" class="no-config">
        <n-empty description="该插件暂无可配置项" />
      </div>

      <template v-for="(field, fieldName) in configFields" :key="fieldName">
        <!-- 字符串输入框 -->
        <n-form-item
          v-if="field.type === 'string'"
          :label="getFieldLabel(field)"
          :path="fieldName"
          :rule="getFieldRule(field)"
        >
          <n-input
            v-model:value="formData[fieldName]"
            :placeholder="field.description || `请输入${field.label}`"
            clearable
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }}
            </n-text>
          </template>
        </n-form-item>

        <!-- 文本域 -->
        <n-form-item
          v-else-if="field.type === 'text'"
          :label="getFieldLabel(field)"
          :path="fieldName"
          :rule="getFieldRule(field)"
        >
          <n-input
            v-model:value="formData[fieldName]"
            type="textarea"
            :placeholder="field.description || `请输入${field.label}`"
            :autosize="{ minRows: 3, maxRows: 6 }"
            clearable
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }}
            </n-text>
          </template>
        </n-form-item>

        <!-- 数字输入框 -->
        <n-form-item
          v-else-if="field.type === 'number'"
          :label="getFieldLabel(field)"
          :path="fieldName"
          :rule="getFieldRule(field)"
        >
          <n-input-number
            v-model:value="formData[fieldName]"
            :placeholder="field.description || `请输入${field.label}`"
            style="width: 100%"
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }}
            </n-text>
          </template>
        </n-form-item>

        <!-- 布尔开关 -->
        <n-form-item
          v-else-if="field.type === 'boolean'"
          :label="getFieldLabel(field)"
          :path="fieldName"
          :rule="getFieldRule(field)"
        >
          <n-switch
            v-model:value="formData[fieldName]"
            :checked-text="field.description || '启用'"
            :unchecked-text="'禁用'"
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }}
            </n-text>
          </template>
        </n-form-item>

        <!-- 选择框 -->
        <n-form-item
          v-else-if="field.type === 'select'"
          :label="getFieldLabel(field)"
          :path="fieldName"
          :rule="getFieldRule(field)"
        >
          <n-select
            v-model:value="formData[fieldName]"
            :placeholder="`请选择${field.label}`"
            :options="getFieldOptions(field)"
            clearable
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }}
            </n-text>
          </template>
        </n-form-item>

        <!-- 未知类型 -->
        <n-form-item
          v-else
          :label="getFieldLabel(field)"
          :path="fieldName"
        >
          <n-input
            v-model:value="formData[fieldName]"
            :placeholder="field.description || `请输入${field.label}`"
            disabled
          />
          <template #tip>
            <n-text depth="3" style="font-size: 12px">
              {{ field.description }} (不支持的字段类型: {{ field.type }})
            </n-text>
          </template>
        </n-form-item>
      </template>
    </n-form>

    <div class="form-actions">
      <n-space justify="end">
        <n-button @click="handleReset">重置</n-button>
        <n-button type="primary" @click="handleSave" :loading="saving">
          保存配置
        </n-button>
      </n-space>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, computed, h } from 'vue'
import { useToast } from '~/composables/useToast'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  },
  config: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['save', 'reset'])

const toast = useToast()
const formRef = ref(null)
const saving = ref(false)

// 配置字段定义
const configFields = computed(() => {
  return props.plugin?.config_fields || {}
})

// 表单数据
const formData = reactive({})

// 表单验证规则
const formRules = reactive({})

// 初始化表单数据
const initializeForm = () => {
  // 清空现有数据
  Object.keys(formData).forEach(key => delete formData[key])
  Object.keys(formRules).forEach(key => delete formRules[key])

  // 根据配置字段初始化
  Object.entries(configFields.value).forEach(([fieldName, field]) => {
    // 设置默认值或当前配置值
    const currentValue = props.config[fieldName]
    formData[fieldName] = currentValue !== undefined ? currentValue : field.default
  })
}

// 获取字段标签
const getFieldLabel = (field) => {
  // 直接返回标签，让 UI 框架自己处理必填标识
  return field.label
}

// 获取字段验证规则
const getFieldRule = (field) => {
  const rules = []

  if (field.required) {
    let message = `请输入${field.label}`
    let trigger = ['input', 'blur', 'change']

    // 根据字段类型调整错误消息和触发事件
    switch (field.type) {
      case 'boolean':
        // 布尔字段：只要值不是 undefined/null 就算通过
        rules.push({
          validator: (rule, value) => {
            if (value === undefined || value === null) {
              return new Error(`请选择${field.label}`)
            }
            return true
          },
          trigger: ['change']
        })
        break
      case 'select':
        message = `请选择${field.label}`
        trigger = ['change']
        rules.push({
          required: true,
          message: message,
          trigger: trigger
        })
        break
      case 'number':
        message = `请输入${field.label}`
        trigger = ['input', 'blur']
        rules.push({
          required: true,
          type: 'number',
          message: message,
          trigger: trigger
        })
        break
      default:
        message = `请输入${field.label}`
        trigger = ['input', 'blur']
        rules.push({
          required: true,
          message: message,
          trigger: trigger
        })
    }
  }

  return rules
}

// 获取选择框选项
const getFieldOptions = (field) => {
  if (!field.options || !Array.isArray(field.options)) {
    return []
  }

  return field.options.map(option => ({
    label: option,
    value: option
  }))
}

// 保存配置
const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true

    emit('save', { ...formData })

    toast.success('插件配置保存成功')
  } catch (error) {
    console.error('表单验证失败:', error)
    toast.error('请检查表单输入')
  } finally {
    saving.value = false
  }
}

// 重置配置
const handleReset = () => {
  initializeForm()
  emit('reset')
  toast.info('配置已重置')
}

// 监听配置变化
watch(() => props.config, initializeForm, { deep: true })
watch(() => props.plugin, initializeForm, { deep: true })

onMounted(() => {
  initializeForm()
})
</script>

<style scoped>
.plugin-config-form {
  padding: 16px;
}

.no-config {
  text-align: center;
  padding: 40px 0;
}

.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

/* 必填字段的红色星号标识 */
:deep(.n-form-item-label) {
  color: var(--n-label-text-color);
}

:deep(.n-form-item-label .required-asterisk) {
  color: #d03050;
  margin-left: 2px;
}
</style>