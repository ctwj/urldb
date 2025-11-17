<template>
  <n-modal
    :show="visible"
    @update:show="handleClose"
    :mask-closable="true"
    preset="card"
    title="举报资源失效"
    class="max-w-md w-full"
    :style="{ maxWidth: '90vw' }"
  >
    <template #header-extra>
      <n-button size="small" circle @click="handleClose">
        <template #icon>
          <i class="fas fa-times"></i>
        </template>
      </n-button>
    </template>

    <div class="space-y-4">
      <div class="text-gray-600 dark:text-gray-300 text-sm">
        请选择举报原因，我们会尽快核实处理：
      </div>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="举报原因" path="reason">
          <n-select
            v-model:value="formData.reason"
            :options="reasonOptions"
            placeholder="请选择举报原因"
            :loading="loading"
          />
        </n-form-item>

        <n-form-item label="详细描述" path="description">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="请详细描述问题，帮助我们更好地处理..."
            :autosize="{ minRows: 3, maxRows: 6 }"
            maxlength="500"
            show-count
          />
        </n-form-item>

        <n-form-item label="联系方式（选填）" path="contact">
          <n-input
            v-model:value="formData.contact"
            placeholder="邮箱或手机号，便于我们反馈处理结果"
          />
        </n-form-item>
      </n-form>

      <div class="text-xs text-gray-500 dark:text-gray-400">
        <i class="fas fa-info-circle mr-1"></i>
        我们承诺保护您的隐私，举报信息仅用于核实处理
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <n-button @click="handleClose" :disabled="submitting">
          取消
        </n-button>
        <n-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          提交举报
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { useMessage } from 'naive-ui'

interface Props {
  visible: boolean
  resourceKey: string
}

interface Emits {
  (e: 'close'): void
  (e: 'submitted'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const message = useMessage()

// 表单数据
const formData = ref({
  reason: '',
  description: '',
  contact: ''
})

// 表单引用
const formRef = ref()

// 状态
const loading = ref(false)
const submitting = ref(false)

// 举报原因选项
const reasonOptions = [
  { label: '链接已失效', value: 'link_invalid' },
  { label: '资源无法下载', value: 'download_failed' },
  { label: '资源内容不符', value: 'content_mismatch' },
  { label: '包含恶意软件', value: 'malicious' },
  { label: '版权问题', value: 'copyright' },
  { label: '其他问题', value: 'other' }
]

// 表单验证规则
const rules = {
  reason: {
    required: true,
    message: '请选择举报原因',
    trigger: ['blur', 'change']
  },
  description: {
    required: true,
    message: '请详细描述问题',
    trigger: 'blur'
  }
}

// 关闭模态框
const handleClose = () => {
  if (!submitting.value) {
    emit('close')
    resetForm()
  }
}

// 重置表单
const resetForm = () => {
  formData.value = {
    reason: '',
    description: '',
    contact: ''
  }
  formRef.value?.restoreValidation()
}

// 提交举报
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    submitting.value = true

    // 这里可以调用实际的举报API
    // const reportData = {
    //   resource_key: props.resourceKey,
    //   reason: formData.value.reason,
    //   description: formData.value.description,
    //   contact: formData.value.contact,
    //   user_agent: navigator.userAgent,
    //   ip_address: await getClientIP()
    // }
    // await reportApi.submitReport(reportData)

    // 模拟提交过程
    await new Promise(resolve => setTimeout(resolve, 1500))

    message.success('举报提交成功，我们会尽快核实处理')
    emit('submitted')
  } catch (error) {
    console.error('提交举报失败:', error)
    message.error('提交失败，请重试')
  } finally {
    submitting.value = false
  }
}
</script>