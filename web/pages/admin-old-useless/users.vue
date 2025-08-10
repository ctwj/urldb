<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载用户数据</p>
          </div>
        </div>
      </div>
    </div>

    <div class="p-6">
      <n-alert class="mb-4" title="用户管理功能，可以创建、编辑、删除用户，以及修改用户密码" type="info" />

      <!-- 用户管理内容 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
    <div class="flex justify-between items-center">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">用户管理</h2>
      <div class="flex gap-2">
        <n-button 
          @click="showCreateModal = true" 
          type="primary"
        >
          添加用户
        </n-button>
      </div>
    </div>
  </div>

  <!-- 用户列表 -->
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
    <div class="px-6 py-4 border-b border-gray-200">
      <h2 class="text-lg font-semibold text-gray-900">用户列表</h2>
    </div>
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户名</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">邮箱</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">角色</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">最后登录</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="user in users" :key="user.id">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.id }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.username }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getRoleClass(user.role)" class="px-2 py-1 text-xs font-medium rounded-full">
                {{ user.role }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
                    class="px-2 py-1 text-xs font-medium rounded-full">
                {{ user.is_active ? '激活' : '禁用' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ user.last_login ? formatDate(user.last_login) : '从未登录' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <n-button @click="editUser(user)" type="info" size="small" class="mr-3" :title="user.username === 'admin' ? '管理员用户信息不可修改' : '编辑用户'">
                编辑{{ user.username === 'admin' ? ' (只读)' : '' }}
              </n-button>
              <n-button @click="showChangePasswordModal(user)" type="warning" size="small" class="mr-3">修改密码</n-button>
              <n-button @click="deleteUser(user.id)" type="error" size="small">删除</n-button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <!-- 创建/编辑用户模态框 -->
  <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <h3 class="text-lg font-medium text-gray-900 mb-4">
          {{ showEditModal ? '编辑用户' : '创建用户' }}
        </h3>
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
        <form @submit.prevent="handleSubmit">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">用户名</label>
              <n-input 
                v-model:value="form.username" 
                type="text" 
                required
                :disabled="showEditModal"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">邮箱</label>
              <n-input 
                v-model:value="form.email" 
                type="text" 
                required
                :disabled="showEditModal"
              />
            </div>
            <div v-if="showCreateModal">
              <label class="block text-sm font-medium text-gray-700">密码</label>
              <n-input 
                v-model:value="form.password" 
                type="password" 
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">角色</label>
              <select 
                v-model="form.role" 
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :disabled="showEditModal && editingUser?.username === 'admin'"
              >
                <option value="user">用户</option>
                <option value="admin">管理员</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">激活状态</label>
              <n-switch 
                v-model:value="form.is_active"
                size="medium"
                :disabled="showEditModal && editingUser?.username === 'admin'"
              />
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button 
              type="button" 
              @click="closeModal"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="px-4 py-2 bg-indigo-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="showEditModal && editingUser?.username === 'admin'"
            >
              {{ showEditModal ? '更新' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>

  <!-- 修改密码模态框 -->
  <div v-if="showPasswordModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <h3 class="text-lg font-medium text-gray-900 mb-4">
          修改用户密码
        </h3>
        <p class="text-sm text-gray-600 mb-4">
          正在为用户 <strong>{{ changingPasswordUser?.username }}</strong> 修改密码
        </p>
        <form @submit.prevent="handlePasswordChange">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">新密码</label>
              <n-input 
                v-model:value="passwordForm.newPassword" 
                type="password" 
                required
                minlength="6"
                placeholder="请输入新密码（至少6位）"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">确认新密码</label>
              <n-input 
                v-model:value="passwordForm.confirmPassword" 
                type="password" 
                required
                placeholder="请再次输入新密码"
              />
            </div>
          </div>
          <div class="mt-6 flex justify-end space-x-3">
            <button 
              type="button" 
              @click="closePasswordModal"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button 
              type="submit" 
              class="px-4 py-2 bg-yellow-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-yellow-700"
            >
              修改密码
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin-old',
  ssr: false
})

const router = useRouter()
const userStore = useUserStore()
const notification = useNotification()

const pageLoading = ref(true)
const users = ref<any[]>([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showPasswordModal = ref(false)
const editingUser = ref<any>(null)
const changingPasswordUser = ref<any>(null)
const dialog = useDialog()
const form = ref({
  username: '',
  email: '',
  password: '',
  role: 'user',
  is_active: true
})
const passwordForm = ref({
  newPassword: '',
  confirmPassword: ''
})

// 检查认证
const checkAuth = () => {
  userStore.initAuth()
  if (!userStore.isAuthenticated) {
    router.push('/login')
    return
  }
}

// 获取用户列表
const fetchUsers = async () => {
  try {
    const { useUserApi } = await import('~/composables/useApi')
    const userApi = useUserApi()
    const response = await userApi.getUsers()
    if (response && (response as any).items && Array.isArray((response as any).items)) {
      users.value = (response as any).items
    } else if (Array.isArray(response)) {
      users.value = response
    } else {
      users.value = []
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
}

// 创建用户
const createUser = async () => {
  try {
    const { useUserApi } = await import('~/composables/useApi')
    const userApi = useUserApi()
    await userApi.createUser(form.value)
    await fetchUsers()
    closeModal()
  } catch (error) {
    console.error('创建用户失败:', error)
  }
}

// 更新用户
const updateUser = async () => {
  try {
    const { useUserApi } = await import('~/composables/useApi')
    const userApi = useUserApi()
    await userApi.updateUser(editingUser.value.id, form.value)
    await fetchUsers()
    closeModal()
  } catch (error) {
    console.error('更新用户失败:', error)
  }
}

// 删除用户
const deleteUser = async (id: any) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除这个用户吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        const { useUserApi } = await import('~/composables/useApi')
        const userApi = useUserApi()
        await userApi.deleteUser(id)
        await fetchUsers()
      } catch (error) {
        console.error('删除用户失败:', error)
      }
    }
  })
}

// 显示修改密码模态框
const showChangePasswordModal = (user: any) => {
  changingPasswordUser.value = user
  passwordForm.value = {
    newPassword: '',
    confirmPassword: ''
  }
  showPasswordModal.value = true
}

// 关闭修改密码模态框
const closePasswordModal = () => {
  showPasswordModal.value = false
  changingPasswordUser.value = null
  passwordForm.value = {
    newPassword: '',
    confirmPassword: ''
  }
}

// 修改密码
const changePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    notification.error({
      title: '失败',
      content: '两次输入的密码不一致',
      duration: 3000
    })
    return
  }
  
  if (passwordForm.value.newPassword.length < 6) {
    notification.error({
      title: '失败',
      content: '密码长度至少6位',
      duration: 3000
    })
    return
  }
  
  try {
    const { useUserApi } = await import('~/composables/useApi')
    const userApi = useUserApi()
    await userApi.changePassword(changingPasswordUser.value.id, passwordForm.value.newPassword)
    notification.success({
      title: '成功',
      content: '密码修改成功',
      duration: 3000
    })
    closePasswordModal()
  } catch (error: any) {
    console.error('修改密码失败:', error)
    notification.error({
      title: '失败',
      content: '修改密码失败: ' + (error.message || '未知错误'),
      duration: 3000
    })
  }
}

// 处理密码修改表单提交
const handlePasswordChange = () => {
  changePassword()
}

// 编辑用户
const editUser = (user: any) => {
  editingUser.value = user
  form.value = {
    username: user.username,
    email: user.email,
    password: '',
    role: user.role,
    is_active: user.is_active
  }
  showEditModal.value = true
}

// 关闭模态框
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingUser.value = null
  form.value = {
    username: '',
    email: '',
    password: '',
    role: 'user',
    is_active: true
  }
}

// 提交表单
const handleSubmit = () => {
  if (showEditModal.value) {
    updateUser()
  } else {
    createUser()
  }
}

// 获取角色样式
const getRoleClass = (role: any) => {
  return role === 'admin' ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'
}

// 格式化日期
const formatDate = (dateString: any) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

onMounted(() => {
  checkAuth()
  fetchUsers()
})
</script> 