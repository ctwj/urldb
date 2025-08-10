<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">用户管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的用户账户</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加用户
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 提示信息 -->
    <n-alert title="用户管理功能，可以创建、编辑、删除用户，以及修改用户密码" type="info" />

    <!-- 用户列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">用户列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个用户</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="users.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无用户</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加用户"按钮创建新用户</div>
        <n-button @click="showCreateModal = true" type="primary">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加用户
        </n-button>
      </div>

      <div v-else>
        <n-data-table
          :columns="columns"
          :data="users"
          :pagination="pagination"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          @update:page="handlePageChange"
        />
      </div>
    </n-card>

    <!-- 创建/编辑用户模态框 -->
    <n-modal v-model:show="showModal" preset="card" :title="showEditModal ? '编辑用户' : '创建用户'" style="width: 500px">
      <div v-if="showEditModal && editingUser?.username === 'admin'" class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
        <p class="text-sm text-yellow-800">
          <i class="fas fa-exclamation-triangle mr-2"></i>
          管理员用户信息不可修改，只能通过修改密码功能来更新密码。
        </p>
      </div>
      <div v-if="showEditModal && editingUser?.username !== 'admin'" class="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
        <p class="text-sm text-blue-800">
          <i class="fas fa-info-circle mr-2"></i>
          编辑模式：用户名和邮箱不可修改，只能修改角色和激活状态。
        </p>
      </div>
      
      <n-form
        ref="formRef"
        :model="userForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="用户名" path="username">
          <n-input
            v-model:value="userForm.username"
            placeholder="请输入用户名"
            :disabled="showEditModal"
          />
        </n-form-item>

        <n-form-item label="邮箱" path="email">
          <n-input
            v-model:value="userForm.email"
            placeholder="请输入邮箱"
            :disabled="showEditModal"
          />
        </n-form-item>

        <n-form-item v-if="!showEditModal" label="密码" path="password">
          <n-input
            v-model:value="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
          />
        </n-form-item>

        <n-form-item label="角色" path="role">
          <n-select
            v-model:value="userForm.role"
            :options="roleOptions"
            placeholder="请选择角色"
          />
        </n-form-item>

        <n-form-item label="状态" path="is_active">
          <n-switch v-model:value="userForm.is_active" />
          <span class="ml-2 text-sm text-gray-500">{{ userForm.is_active ? '激活' : '禁用' }}</span>
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="closeModal">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ showEditModal ? '更新' : '创建' }}
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 修改密码模态框 -->
    <n-modal v-model:show="showChangePasswordModal" preset="card" title="修改密码" style="width: 400px">
      <n-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="新密码" path="new_password">
          <n-input
            v-model:value="passwordForm.new_password"
            type="password"
            placeholder="请输入新密码"
            show-password-on="click"
          />
        </n-form-item>

        <n-form-item label="确认密码" path="confirm_password">
          <n-input
            v-model:value="passwordForm.confirm_password"
            type="password"
            placeholder="请再次输入新密码"
            show-password-on="click"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="showChangePasswordModal = false">取消</n-button>
          <n-button type="primary" @click="handleChangePassword" :loading="changingPassword">
            修改密码
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

interface User {
  id: number
  username: string
  email: string
  role: string
  is_active: boolean
  last_login?: string
  created_at: string
  updated_at: string
}

const notification = useNotification()
const dialog = useDialog()
const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showChangePasswordModal = ref(false)
const editingUser = ref<User | null>(null)
const changingPasswordUser = ref<User | null>(null)
const submitting = ref(false)
const changingPassword = ref(false)
const formRef = ref()
const passwordFormRef = ref()

// 用户表单
const userForm = ref({
  username: '',
  email: '',
  password: '',
  role: 'user',
  is_active: true
})

// 密码表单
const passwordForm = ref({
  new_password: '',
  confirm_password: ''
})

// 角色选项
const roleOptions = [
  { label: '用户', value: 'user' },
  { label: '管理员', value: 'admin' }
]

// 表单验证规则
const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur'
  },
  email: {
    required: true,
    message: '请输入邮箱',
    trigger: 'blur',
    pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur',
    min: 6
  },
  role: {
    required: true,
    message: '请选择角色',
    trigger: 'change'
  }
}

// 密码验证规则
const passwordRules = {
  new_password: {
    required: true,
    message: '请输入新密码',
    trigger: 'blur',
    min: 6
  },
  confirm_password: {
    required: true,
    message: '请确认密码',
    trigger: 'blur',
    validator: (rule: any, value: string) => {
      if (value !== passwordForm.value.new_password) {
        return new Error('两次输入的密码不一致')
      }
      return true
    }
  }
}

// 获取用户API
import { useUserApi } from '~/composables/useApi'
import { h } from 'vue'
const userApi = useUserApi()

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render: (row: User) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '用户名',
    key: 'username',
    render: (row: User) => {
      return h('span', { title: row.username }, row.username)
    }
  },
  {
    title: '邮箱',
    key: 'email',
    render: (row: User) => {
      return h('span', { title: row.email }, row.email)
    }
  },
  {
    title: '角色',
    key: 'role',
    width: 100,
    render: (row: User) => {
      const roleClass = row.role === 'admin' 
        ? 'px-2 py-1 text-xs font-medium rounded-full bg-purple-100 text-purple-800 dark:bg-purple-900/20 dark:text-purple-400'
        : 'px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400'
      return h('span', { class: roleClass }, row.role)
    }
  },
  {
    title: '状态',
    key: 'is_active',
    width: 100,
    render: (row: User) => {
      const statusClass = row.is_active
        ? 'px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
        : 'px-2 py-1 text-xs font-medium rounded-full bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
      return h('span', { class: statusClass }, row.is_active ? '激活' : '禁用')
    }
  },
  {
    title: '最后登录',
    key: 'last_login',
    width: 180,
    render: (row: User) => {
      return h('span', { class: 'text-gray-500' }, row.last_login ? formatDate(row.last_login) : '从未登录')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row: User) => {
      return h('div', { class: 'flex items-center gap-2' }, [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors',
          onClick: () => editUser(row),
          title: row.username === 'admin' ? '管理员用户信息不可修改' : '编辑用户'
        }, [
          h('i', { class: 'fas fa-edit mr-1' }),
          row.username === 'admin' ? '编辑(只读)' : '编辑'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-yellow-100 hover:bg-yellow-200 text-yellow-700 dark:bg-yellow-900/20 dark:text-yellow-400 rounded transition-colors',
          onClick: () => showChangePasswordModalFunc(row)
        }, [
          h('i', { class: 'fas fa-key mr-1' }),
          '修改密码'
        ]),
        h('button', {
          class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
          onClick: () => deleteUser(row.id),
          disabled: row.username === 'admin'
        }, [
          h('i', { class: 'fas fa-trash mr-1' }),
          '删除'
        ])
      ])
    }
  }
]

// 分页配置
const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    currentPage.value = page
    fetchData()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchData()
  }
}))

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await userApi.getUsers({
      page: currentPage.value,
      page_size: pageSize.value
    }) as any
    
    if (response && response.data) {
      users.value = response.data
      total.value = response.total || 0
    } else if (Array.isArray(response)) {
      users.value = response
      total.value = response.length
    } else {
      users.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取用户失败:', error)
    users.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 编辑用户
const editUser = (user: User) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    email: user.email,
    password: '',
    role: user.role,
    is_active: user.is_active
  }
  showEditModal.value = true
}

// 删除用户
const deleteUser = async (userId: number) => {
  const user = users.value.find(u => u.id === userId)
  if (user?.username === 'admin') {
    notification.error({
      content: '不能删除管理员用户',
      duration: 3000
    })
    return
  }

  dialog.warning({
    title: '警告',
    content: `确定要删除用户"${user?.username}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await userApi.deleteUser(userId)
        notification.success({
          content: '删除成功',
          duration: 3000
        })
        fetchData()
      } catch (error) {
        console.error('删除失败:', error)
        notification.error({
          content: '删除失败',
          duration: 3000
        })
      }
    }
  })
}

// 显示修改密码模态框
const showChangePasswordModalFunc = (user: User) => {
  changingPasswordUser.value = user
  passwordForm.value = {
    new_password: '',
    confirm_password: ''
  }
  showChangePasswordModal.value = true
}

// 关闭模态框
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingUser.value = null
  userForm.value = {
    username: '',
    email: '',
    password: '',
    role: 'user',
    is_active: true
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    await formRef.value?.validate()
    
    if (showEditModal.value) {
      await userApi.updateUser(editingUser.value!.id, userForm.value)
      notification.success({
        content: '更新成功',
        duration: 3000
      })
    } else {
      await userApi.createUser(userForm.value)
      notification.success({
        content: '创建成功',
        duration: 3000
      })
    }
    
    closeModal()
    fetchData()
  } catch (error) {
    console.error('提交失败:', error)
    notification.error({
      content: '操作失败',
      duration: 3000
    })
  } finally {
    submitting.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  try {
    changingPassword.value = true
    await passwordFormRef.value?.validate()
    
    await userApi.changePassword(changingPasswordUser.value!.id, passwordForm.value.new_password)
    
    notification.success({
      content: '密码修改成功',
      duration: 3000
    })
    
    showChangePasswordModal.value = false
    changingPasswordUser.value = null
    passwordForm.value = {
      new_password: '',
      confirm_password: ''
    }
  } catch (error) {
    console.error('修改密码失败:', error)
    notification.error({
      content: '修改密码失败',
      duration: 3000
    })
  } finally {
    changingPassword.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})



// 计算属性
const showModal = computed({
  get: () => showCreateModal.value || showEditModal.value,
  set: (value: boolean) => {
    if (!value) {
      showCreateModal.value = false
      showEditModal.value = false
    }
  }
})
</script>

<style scoped>
/* 自定义样式 */
</style> 