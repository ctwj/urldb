<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">设置</h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">账户设置和偏好</p>
    </div>

    <!-- 密码修改 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">密码修改</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">修改您的登录密码</p>
      </div>
      
      <div class="p-6">
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
                placeholder="请输入新密码"
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
          <div class="flex justify-end space-x-4 pt-6">
            <n-button @click="handleResetPassword">
              重置
            </n-button>
            <n-button type="primary" @click="handleChangePassword" :loading="changingPassword">
              修改密码
            </n-button>
          </div>
        </n-form>
      </div>
    </div>

    <!-- 通知设置 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">通知设置</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">管理您的通知偏好</p>
      </div>
      
      <div class="p-6">
        <div class="space-y-6">
          <!-- 邮件通知 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">邮件通知</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">接收重要通知的邮件</p>
            </div>
            <n-switch v-model:value="notificationSettings.email" />
          </div>

          <!-- 系统通知 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">系统通知</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">接收系统相关通知</p>
            </div>
            <n-switch v-model:value="notificationSettings.system" />
          </div>

          <!-- 资源更新通知 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">资源更新通知</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">接收资源更新相关通知</p>
            </div>
            <n-switch v-model:value="notificationSettings.resource" />
          </div>
        </div>

        <!-- 保存按钮 -->
        <div class="flex justify-end pt-6">
          <n-button type="primary" @click="handleSaveNotificationSettings" :loading="savingSettings">
            保存设置
          </n-button>
        </div>
      </div>
    </div>

    <!-- 隐私设置 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">隐私设置</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">管理您的隐私偏好</p>
      </div>
      
      <div class="p-6">
        <div class="space-y-6">
          <!-- 个人资料可见性 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">个人资料可见性</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">允许其他用户查看您的个人资料</p>
            </div>
            <n-switch v-model:value="privacySettings.profileVisible" />
          </div>

          <!-- 浏览历史记录 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">记录浏览历史</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">保存您的浏览历史记录</p>
            </div>
            <n-switch v-model:value="privacySettings.saveHistory" />
          </div>

          <!-- 数据收集 -->
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">数据收集</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">允许收集匿名使用数据以改善服务</p>
            </div>
            <n-switch v-model:value="privacySettings.dataCollection" />
          </div>
        </div>

        <!-- 保存按钮 -->
        <div class="flex justify-end pt-6">
          <n-button type="primary" @click="handleSavePrivacySettings" :loading="savingPrivacy">
            保存设置
          </n-button>
        </div>
      </div>
    </div>

    <!-- 危险操作 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-red-600 dark:text-red-400">危险操作</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">这些操作不可撤销，请谨慎操作</p>
      </div>
      
      <div class="p-6">
        <div class="space-y-4">
          <!-- 删除账户 -->
          <div class="flex items-center justify-between p-4 border border-red-200 dark:border-red-800 rounded-lg">
            <div>
              <h4 class="text-sm font-medium text-red-600 dark:text-red-400">删除账户</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">永久删除您的账户和所有数据</p>
            </div>
            <n-button type="error" @click="handleDeleteAccount">
              删除账户
            </n-button>
          </div>

          <!-- 导出数据 -->
          <div class="flex items-center justify-between p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div>
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">导出数据</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">导出您的个人数据</p>
            </div>
            <n-button type="info" @click="handleExportData">
              导出数据
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 页面元数据
definePageMeta({
  layout: 'user',
  title: '设置'
})

// 表单引用
const passwordFormRef = ref()

// 加载状态
const changingPassword = ref(false)
const savingSettings = ref(false)
const savingPrivacy = ref(false)

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

// 通知设置
const notificationSettings = ref({
  email: true,
  system: true,
  resource: false
})

// 隐私设置
const privacySettings = ref({
  profileVisible: true,
  saveHistory: true,
  dataCollection: false
})

// 处理修改密码
const handleChangePassword = async () => {
  try {
    await passwordFormRef.value?.validate()
    changingPassword.value = true
    
    // TODO: 调用API修改密码
    console.log('修改密码:', passwordForm.value)
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
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
  } catch (error) {
    console.error('修改密码失败:', error)
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

// 处理保存通知设置
const handleSaveNotificationSettings = async () => {
  savingSettings.value = true
  
  try {
    // TODO: 调用API保存通知设置
    console.log('保存通知设置:', notificationSettings.value)
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    if (process.client) {
      const notification = useNotification()
      notification.success({
        content: '通知设置保存成功',
        duration: 3000
      })
    }
  } catch (error) {
    console.error('保存通知设置失败:', error)
  } finally {
    savingSettings.value = false
  }
}

// 处理保存隐私设置
const handleSavePrivacySettings = async () => {
  savingPrivacy.value = true
  
  try {
    // TODO: 调用API保存隐私设置
    console.log('保存隐私设置:', privacySettings.value)
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    if (process.client) {
      const notification = useNotification()
      notification.success({
        content: '隐私设置保存成功',
        duration: 3000
      })
    }
  } catch (error) {
    console.error('保存隐私设置失败:', error)
  } finally {
    savingPrivacy.value = false
  }
}

// 处理删除账户
const handleDeleteAccount = () => {
  // TODO: 实现删除账户功能
  console.log('删除账户')
}

// 处理导出数据
const handleExportData = () => {
  // TODO: 实现导出数据功能
  console.log('导出数据')
}

// 页面加载时获取数据
onMounted(() => {
  // TODO: 获取用户设置数据
  console.log('加载用户设置数据')
})
</script> 