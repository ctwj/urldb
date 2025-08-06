<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">个人资料</h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">编辑您的个人信息</p>
    </div>

    <!-- 个人资料表单 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">基本信息</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">更新您的个人资料信息</p>
      </div>
      
      <div class="p-6">
        <n-form
          ref="formRef"
          :model="profileForm"
          :rules="formRules"
          label-placement="left"
          label-width="auto"
          require-mark-placement="right-hanging"
        >
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 用户名 -->
            <n-form-item label="用户名" path="username">
              <n-input
                v-model:value="profileForm.username"
                placeholder="请输入用户名"
                :disabled="true"
              />
              <template #feedback>
                <span class="text-xs text-gray-500">用户名不可修改</span>
              </template>
            </n-form-item>

            <!-- 邮箱 -->
            <n-form-item label="邮箱" path="email">
              <n-input
                v-model:value="profileForm.email"
                placeholder="请输入邮箱地址"
              />
            </n-form-item>

            <!-- 昵称 -->
            <n-form-item label="昵称" path="nickname">
              <n-input
                v-model:value="profileForm.nickname"
                placeholder="请输入昵称"
              />
            </n-form-item>

            <!-- 手机号 -->
            <n-form-item label="手机号" path="phone">
              <n-input
                v-model:value="profileForm.phone"
                placeholder="请输入手机号"
              />
            </n-form-item>

            <!-- 性别 -->
            <n-form-item label="性别" path="gender">
              <n-select
                v-model:value="profileForm.gender"
                :options="genderOptions"
                placeholder="请选择性别"
              />
            </n-form-item>

            <!-- 生日 -->
            <n-form-item label="生日" path="birthday">
              <n-date-picker
                v-model:value="profileForm.birthday"
                type="date"
                placeholder="请选择生日"
              />
            </n-form-item>
          </div>

          <!-- 个人简介 -->
          <n-form-item label="个人简介" path="bio">
            <n-input
              v-model:value="profileForm.bio"
              type="textarea"
              placeholder="请输入个人简介"
              :rows="4"
            />
          </n-form-item>

          <!-- 操作按钮 -->
          <div class="flex justify-end space-x-4 pt-6">
            <n-button @click="handleReset">
              重置
            </n-button>
            <n-button type="primary" @click="handleSave" :loading="saving">
              保存修改
            </n-button>
          </div>
        </n-form>
      </div>
    </div>

    <!-- 账户信息 -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">账户信息</h3>
        <p class="text-gray-600 dark:text-gray-400 mt-1">您的账户基本信息</p>
      </div>
      
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-4">
            <div class="flex justify-between items-center py-3 border-b border-gray-200 dark:border-gray-700">
              <span class="text-gray-600 dark:text-gray-400">用户ID</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ accountInfo.userId }}</span>
            </div>
            <div class="flex justify-between items-center py-3 border-b border-gray-200 dark:border-gray-700">
              <span class="text-gray-600 dark:text-gray-400">注册时间</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ accountInfo.registerTime }}</span>
            </div>
            <div class="flex justify-between items-center py-3 border-b border-gray-200 dark:border-gray-700">
              <span class="text-gray-600 dark:text-gray-400">最后登录</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ accountInfo.lastLogin }}</span>
            </div>
            <div class="flex justify-between items-center py-3 border-b border-gray-200 dark:border-gray-700">
              <span class="text-gray-600 dark:text-gray-400">账户状态</span>
              <n-tag :type="accountInfo.status === 'active' ? 'success' : 'warning'">
                {{ accountInfo.status === 'active' ? '正常' : '待激活' }}
              </n-tag>
            </div>
            <div class="flex justify-between items-center py-3">
              <span class="text-gray-600 dark:text-gray-400">用户角色</span>
              <n-tag type="info">
                {{ accountInfo.role === 'admin' ? '管理员' : '普通用户' }}
              </n-tag>
            </div>
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
  title: '个人资料'
})

// 表单引用
const formRef = ref()

// 保存状态
const saving = ref(false)

// 表单数据
const profileForm = ref({
  username: '',
  email: '',
  nickname: '',
  phone: '',
  gender: null,
  birthday: null,
  bio: ''
})

// 性别选项
const genderOptions = [
  { label: '男', value: 'male' },
  { label: '女', value: 'female' },
  { label: '其他', value: 'other' }
]

// 表单验证规则
const formRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  nickname: [
    { max: 20, message: '昵称不能超过20个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ]
}

// 账户信息
const accountInfo = ref({
  userId: '',
  registerTime: '',
  lastLogin: '',
  status: 'active',
  role: 'user'
})

// 处理保存
const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true
    
    // TODO: 调用API保存个人资料
    console.log('保存个人资料:', profileForm.value)
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 显示成功提示
    if (process.client) {
      const notification = useNotification()
      notification.success({
        content: '个人资料保存成功',
        duration: 3000
      })
    }
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

// 处理重置
const handleReset = () => {
  // TODO: 重置表单数据
  console.log('重置表单')
}

// 页面加载时获取数据
onMounted(() => {
  // TODO: 获取用户个人资料数据
  console.log('加载用户个人资料数据')
  
  // 模拟数据
  profileForm.value = {
    username: 'testuser',
    email: 'test@example.com',
    nickname: '测试用户',
    phone: '13800138000',
    gender: null,
    birthday: null,
    bio: '这是一个测试用户的个人简介'
  }
  
  accountInfo.value = {
    userId: 'U001',
    registerTime: '2024-01-01 10:00:00',
    lastLogin: '2024-01-15 14:30:00',
    status: 'active',
    role: 'user'
  }
})
</script> 