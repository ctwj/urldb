<template>
  <div class="space-y-6 max-w-4xl">
    <!-- 页面标题 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">设置</h1>
      <p class="text-gray-500 dark:text-gray-400 mt-1">账户设置和安全偏好</p>
    </div>

    <!-- 密码修改 -->
    <n-card :bordered="false" class="shadow-sm">
      <template #header>
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
            <i class="fas fa-shield-alt text-blue-500 dark:text-blue-400 text-sm"></i>
          </div>
          <div>
            <p class="font-semibold text-gray-900 dark:text-white leading-tight">密码修改</p>
            <p class="text-xs text-gray-400 dark:text-gray-500">建议定期更换密码，保障账户安全</p>
          </div>
        </div>
      </template>

      <n-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 当前密码 -->
          <n-form-item label="当前密码" path="currentPassword">
            <n-input
              v-model:value="passwordForm.currentPassword"
              type="password"
              placeholder="请输入当前密码"
              show-password-on="click"
            />
          </n-form-item>

          <!-- 新密码 -->
          <n-form-item label="新密码" path="newPassword">
            <n-input
              v-model:value="passwordForm.newPassword"
              type="password"
              placeholder="至少 6 位，建议字母 + 数字组合"
              show-password-on="click"
            />
          </n-form-item>

          <!-- 确认新密码 -->
          <n-form-item label="确认新密码" path="confirmPassword">
            <n-input
              v-model:value="passwordForm.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              show-password-on="click"
            />
          </n-form-item>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-end gap-3 pt-6 border-t border-gray-100 dark:border-gray-700 mt-2">
          <n-button @click="handleResetPassword">
            重置
          </n-button>
          <n-button type="primary" @click="handleChangePassword" :loading="changingPassword">
            <template #icon>
              <i class="fas fa-check"></i>
            </template>
            确认修改
          </n-button>
        </div>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { useApiFetch } from '~/composables/useApiFetch'

// 页面元数据
definePageMeta({
  layout: 'user',
  title: '设置'
})

// 表单引用
const passwordFormRef = ref()

// 加载状态
const changingPassword = ref(false)

// 密码表单
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码验证规则
const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string) => {
        if (value !== passwordForm.value.newPassword) {
          return new Error('两次输入的密码不一致')
        }
        return true
      },
      trigger: 'blur'
    }
  ]
}

// 处理修改密码
const handleChangePassword = async () => {
  try {
    await passwordFormRef.value?.validate()
    changingPassword.value = true

    // 调用后端接口修改密码
    await useApiFetch('/auth/password', {
      method: 'PUT',
      body: {
        current_password: passwordForm.value.currentPassword,
        new_password: passwordForm.value.newPassword
      }
    })

    // 显示成功提示
    if (process.client) {
      const notification = useNotification()
      notification.success({
        content: '密码修改成功',
        duration: 3000
      })
    }

    // 重置表单
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error: any) {
    console.error('修改密码失败:', error)
    if (process.client) {
      const notification = useNotification()
      notification.error({
        content: error?.data?.message || error?.message || '密码修改失败，请检查当前密码是否正确',
        duration: 4000
      })
    }
  } finally {
    changingPassword.value = false
  }
}

// 处理重置密码表单
const handleResetPassword = () => {
  passwordForm.value = {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}
</script> 