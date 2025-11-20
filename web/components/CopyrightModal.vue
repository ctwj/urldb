<template>
  <n-modal
    :show="visible"
    @update:show="handleClose"
    :mask-closable="true"
    preset="card"
    title="版权申述"
    class="max-w-lg w-full"
    :style="{ maxWidth: '95vw' }"
  >

    <div class="space-y-4">
      <div class="bg-blue-50 dark:bg-blue-500/10 border border-blue-200 dark:border-blue-400/30 rounded-lg p-4">
        <div class="flex items-start gap-2">
          <i class="fas fa-info-circle text-blue-600 dark:text-blue-400 mt-0.5"></i>
          <div class="text-sm text-blue-800 dark:text-blue-200">
            <p class="font-medium mb-1">版权申述说明：</p>
            <ul class="space-y-1 text-xs">
              <li>• 请确保您是版权所有者或授权代表</li>
              <li>• 提供真实准确的版权证明材料</li>
              <li>• 虚假申述可能承担法律责任</li>
              <li>• 我们会在收到申述后及时处理</li>
            </ul>
          </div>
        </div>
      </div>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="申述人身份" path="identity">
          <n-select
            v-model:value="formData.identity"
            :options="identityOptions"
            placeholder="请选择您的身份"
            :loading="loading"
          />
        </n-form-item>

        <n-form-item label="权利证明" path="proof_type">
          <n-select
            v-model:value="formData.proof_type"
            :options="proofOptions"
            placeholder="请选择权利证明类型"
          />
        </n-form-item>

        <n-form-item label="版权证明文件" path="proof_files">
          <n-upload
            v-model:file-list="formData.proof_files"
            :max="5"
            :default-upload="false"
            accept=".pdf,.jpg,.jpeg,.png,.gif"
            @change="handleFileChange"
          >
            <n-upload-dragger>
              <div class="text-center">
                <i class="fas fa-cloud-upload-alt text-4xl text-gray-400 mb-2"></i>
                <p class="text-gray-600 dark:text-gray-300">
                  点击或拖拽上传版权证明文件
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  支持 PDF、JPG、PNG 格式，最多5个文件
                </p>
              </div>
            </n-upload-dragger>
          </n-upload>
        </n-form-item>

        <n-form-item label="申述理由" path="reason">
          <n-input
            v-model:value="formData.reason"
            type="textarea"
            placeholder="请详细说明版权申述理由，包括具体的侵权情况..."
            :autosize="{ minRows: 4, maxRows: 8 }"
            maxlength="1000"
            show-count
          />
        </n-form-item>

        <n-form-item label="联系信息" path="contact_info">
          <n-input
            v-model:value="formData.contact_info"
            placeholder="请提供有效的联系方式（邮箱/电话），以便我们与您联系"
          />
        </n-form-item>

        <n-form-item label="申述人姓名" path="claimant_name">
          <n-input
            v-model:value="formData.claimant_name"
            placeholder="请填写申述人真实姓名或公司名称"
          />
        </n-form-item>

        <n-form-item>
          <n-checkbox v-model:checked="formData.agreement">
            我确认以上信息真实有效，并承担相应的法律责任
          </n-checkbox>
        </n-form-item>
      </n-form>

      <div class="text-xs text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-500/10 rounded p-3">
        <i class="fas fa-exclamation-triangle mr-1"></i>
        请谨慎提交版权申述，虚假申述可能承担法律责任
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
          :disabled="!formData.agreement"
          @click="handleSubmit"
        >
          提交申述
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { useResourceApi } from '~/composables/useApi'

const resourceApi = useResourceApi()

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
  identity: '',
  proof_type: '',
  proof_files: [],
  reason: '',
  contact_info: '',
  claimant_name: '',
  agreement: false
})

// 表单引用
const formRef = ref()

// 状态
const loading = ref(false)
const submitting = ref(false)

// 身份选项
const identityOptions = [
  { label: '版权所有者', value: 'copyright_owner' },
  { label: '授权代表', value: 'authorized_agent' },
  { label: '律师事务所', value: 'law_firm' },
  { label: '其他', value: 'other' }
]

// 证明类型选项
const proofOptions = [
  { label: '版权登记证书', value: 'copyright_certificate' },
  { label: '作品首发证明', value: 'first_publish_proof' },
  { label: '授权委托书', value: 'authorization_letter' },
  { label: '身份证明文件', value: 'identity_document' },
  { label: '其他证明材料', value: 'other_proof' }
]

// 表单验证规则
const rules = {
  identity: {
    required: true,
    message: '请选择申述人身份',
    trigger: ['blur', 'change']
  },
  proof_type: {
    required: true,
    message: '请选择权利证明类型',
    trigger: ['blur', 'change']
  },
  reason: {
    required: true,
    message: '请详细说明申述理由',
    trigger: 'blur'
  },
  contact_info: {
    required: true,
    message: '请提供联系信息',
    trigger: 'blur'
  },
  claimant_name: {
    required: true,
    message: '请填写申述人姓名',
    trigger: 'blur'
  }
}

// 处理文件变化
const handleFileChange = (options: any) => {
  console.log('文件变化:', options)
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
    identity: '',
    proof_type: '',
    proof_files: [],
    reason: '',
    contact_info: '',
    claimant_name: '',
    agreement: false
  }
  formRef.value?.restoreValidation()
}

// 提交申述
const handleSubmit = async () => {
  try {
    await formRef.value?.validate()

    if (!formData.value.agreement) {
      message.warning('请确认申述信息真实有效并承担相应法律责任')
      return
    }

    submitting.value = true

    // 构建证明文件数组（从文件列表转换为字符串）
    const proofFilesArray = formData.value.proof_files.map((file: any) => ({
      id: file.id,
      name: file.name,
      status: file.status,
      percentage: file.percentage
    }))

    // 调用实际的版权申述API
    const copyrightData = {
      resource_key: props.resourceKey,
      identity: formData.value.identity,
      proof_type: formData.value.proof_type,
      reason: formData.value.reason,
      contact_info: formData.value.contact_info,
      claimant_name: formData.value.claimant_name,
      proof_files: JSON.stringify(proofFilesArray), // 将文件信息转换为JSON字符串
      user_agent: navigator.userAgent,
      ip_address: '' // 服务端获取IP
    }

    const result = await resourceApi.submitCopyrightClaim(copyrightData)
    console.log('版权申述提交结果:', result)

    message.success('版权申述提交成功，我们会在24小时内处理并回复')
    emit('submitted') // 发送提交事件
  } catch (error: any) {
    console.error('提交版权申述失败:', error)
    let errorMessage = '提交失败，请重试'
    if (error && typeof error === 'object' && error.data) {
      errorMessage = error.data.message || errorMessage
    } else if (error && typeof error === 'object' && error.message) {
      errorMessage = error.message
    }
    message.error(errorMessage)
  } finally {
    submitting.value = false
  }
}
</script>