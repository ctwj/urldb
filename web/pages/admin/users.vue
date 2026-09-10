<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和操作按钮 -->
    <template #page-header>
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
    </template>

    <!-- 通知区域 -->
    <template #notice-section>
      <n-alert title="用户管理功能，可以创建、编辑、删除用户，以及修改用户密码" type="info" />
    </template>

    <!-- 内容区header -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <span class="text-lg font-semibold">用户列表</span>
        <span class="text-sm text-gray-500">共 {{ total }} 个用户</span>
      </div>
    </template>

    <!-- 内容区 - 用户列表 -->
    <template #content>

        <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <AdminEmptyState
        v-else-if="users.length === 0"
        icon="fas fa-users"
        title="暂无用户"
        description='你可以点击上方"添加用户"按钮创建新用户'
      >
        <template #action>
          <n-button @click="showCreateModal = true" type="primary">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加用户
          </n-button>
        </template>
      </AdminEmptyState>

      <div v-else class="h-full">
        <n-data-table
          :columns="columns"
          :data="users"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          :scroll-x="1200"
          @update:page="handlePageChange"
        />
      </div>
    </template>

  
    <!-- 内容区footer - 分页组件 -->
    <template #content-footer>
      <div class="p-4">
        <div class="flex justify-center">
          <n-pagination
            v-model:page="currentPage"
            v-model:page-size="pageSize"
            :item-count="total"
            :page-sizes="[100, 200, 500, 1000]"
            show-size-picker
            @update:page="fetchData"
            @update:page-size="(size) => { pageSize = size; currentPage = 1; fetchData() }"
          />
        </div>
      </div>
    </template>
  </AdminPageLayout>

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

    <!-- 用户上传资源列表模态框（管理员） -->
    <n-modal v-model:show="showResourcesModal" preset="card" :title="`用户「${resourcesUser?.username}」的上传资源`" style="width: 960px">
      <div class="mb-3 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <n-select
            v-model:value="resourceStatus"
            :options="resourceStatusOptions"
            placeholder="状态筛选"
            clearable
            style="width: 150px"
            @update:value="() => { resourcePage = 1; fetchUserResources() }"
          />
          <span class="text-sm text-gray-500">共 {{ resourceTotal }} 条</span>
        </div>
        <n-button size="small" @click="fetchUserResources">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>

      <n-data-table
        :columns="resourceColumns"
        :data="resourceList"
        :loading="resourcesLoading"
        :bordered="false"
        size="small"
        :row-key="(row: any) => row.id"
      />

      <div class="mt-3 flex justify-center">
        <n-pagination
          v-model:page="resourcePage"
          :item-count="resourceTotal"
          :page-size="resourcePageSize"
          @update:page="fetchUserResources"
        />
      </div>
    </n-modal>
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'

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
  upload_disabled: boolean
  resource_count?: number
  last_login?: string
  created_at: string
  updated_at: string
}

// 用户上传资源（管理弹窗）
interface UserResourceItem {
  id: number
  title: string
  url: string
  status: string
  fail_reason?: string
  created_at: string
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
    ellipsis: { tooltip: true },
    width: 140,
    render: (row: User) => {
      return h('span', { title: row.username }, row.username)
    }
  },
  {
    title: '邮箱',
    key: 'email',
    ellipsis: { tooltip: true },
    width: 200,
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
      return h('div', { class: 'flex flex-col gap-1 items-start' }, [
        h('span', { class: statusClass }, row.is_active ? '激活' : '禁用'),
        row.upload_disabled
          ? h('span', { class: 'px-2 py-0.5 text-xs font-medium rounded-full bg-orange-100 text-orange-800 dark:bg-orange-900/20 dark:text-orange-400' }, '禁传')
          : null
      ])
    }
  },
  {
    title: '上传资源',
    key: 'resource_count',
    width: 100,
    render: (row: User) => {
      return h('button', {
        class: 'px-2 py-1 text-xs bg-cyan-100 hover:bg-cyan-200 text-cyan-700 dark:bg-cyan-900/20 dark:text-cyan-400 rounded transition-colors',
        onClick: () => openResourcesModal(row),
        title: `查看 ${row.username} 上传的资源列表`
      }, [
        h('i', { class: 'fas fa-cloud-upload-alt mr-1' }),
        String(row.resource_count ?? 0)
      ])
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
    width: 260,
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
          class: row.upload_disabled
            ? 'px-2 py-1 text-xs bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400 rounded transition-colors'
            : 'px-2 py-1 text-xs bg-orange-100 hover:bg-orange-200 text-orange-700 dark:bg-orange-900/20 dark:text-orange-400 rounded transition-colors',
          onClick: () => toggleUploadDisabled(row),
          title: row.upload_disabled ? '恢复该用户上传资源的权限' : '禁止该用户上传资源'
        }, [
          h('i', { class: 'fas fa-ban mr-1' }),
          row.upload_disabled ? '解禁上传' : '禁止上传'
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



// ===== 用户上传资源管理 =====

// 上传资源状态选项（对齐后端五态）
const resourceStatusOptions = [
  { label: '未检测', value: 'pending' },
  { label: '有效', value: 'valid' },
  { label: '无效', value: 'invalid' },
  { label: '处理中', value: 'processing' },
  { label: '已公开', value: 'published' }
]

const statusLabel = (status: string) => resourceStatusOptions.find(o => o.value === status)?.label || status
const statusBadgeClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400',
    valid: 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400',
    invalid: 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400',
    processing: 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400',
    published: 'bg-purple-100 text-purple-800 dark:bg-purple-900/20 dark:text-purple-400'
  }
  return `px-2 py-1 text-xs font-medium rounded-full ${map[status] || map.pending}`
}

const showResourcesModal = ref(false)
const resourcesUser = ref<User | null>(null)
const resourceList = ref<UserResourceItem[]>([])
const resourceTotal = ref(0)
const resourcePage = ref(1)
const resourcePageSize = 10
const resourceStatus = ref<string | null>(null)
const resourcesLoading = ref(false)

// 打开用户上传资源弹窗
const openResourcesModal = (user: User) => {
  resourcesUser.value = user
  resourcePage.value = 1
  resourceStatus.value = null
  showResourcesModal.value = true
  fetchUserResources()
}

// 拉取该用户的上传资源列表
// 注意：getUserResources 返回原始响应（data.list 形状会被 parseApiResponse 退化为数组丢失 total），这里手动解包 data
const fetchUserResources = async () => {
  if (!resourcesUser.value) return
  resourcesLoading.value = true
  try {
    const params: any = { page: resourcePage.value, page_size: resourcePageSize }
    if (resourceStatus.value) params.status = resourceStatus.value
    const res = await userApi.getUserResources(resourcesUser.value.id, params) as any
    const payload = res?.data || {}
    resourceList.value = payload.list || []
    resourceTotal.value = payload.total || 0
  } catch (error) {
    resourceList.value = []
    resourceTotal.value = 0
    notification.error({ content: '获取资源列表失败', duration: 3000 })
  } finally {
    resourcesLoading.value = false
  }
}

// 管理员删除用户资源
const deleteUserResource = (row: UserResourceItem) => {
  dialog.warning({
    title: '警告',
    content: `确定删除资源「${row.title}」吗？若该资源已发布为公共资源，公共池不受影响。`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await userApi.deleteUserResource(row.id)
        notification.success({ content: '删除成功', duration: 3000 })
        fetchUserResources()
        fetchData() // 同步刷新列表中的上传数
      } catch (error) {
        notification.error({ content: '删除失败', duration: 3000 })
      }
    }
  })
}

// 禁止/恢复上传权限
const toggleUploadDisabled = async (row: User) => {
  try {
    await userApi.updateUserUploadStatus(row.id, !row.upload_disabled)
    notification.success({
      content: row.upload_disabled ? '已恢复该用户上传权限' : '已禁止该用户上传资源',
      duration: 3000
    })
    fetchData()
  } catch (error) {
    notification.error({ content: '操作失败', duration: 3000 })
  }
}

// 上传资源弹窗表格列
const resourceColumns = [
  { title: 'ID', key: 'id', width: 60, render: (row: UserResourceItem) => h('span', { class: 'font-medium' }, row.id) },
  {
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true },
    render: (row: UserResourceItem) => h('span', { title: row.title, class: 'cursor-help' }, row.title)
  },
  {
    title: '链接',
    key: 'url',
    width: 180,
    ellipsis: { tooltip: true },
    render: (row: UserResourceItem) => h('span', { title: row.url, class: 'text-gray-500 cursor-help' }, row.url)
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row: UserResourceItem) => h('span', { class: statusBadgeClass(row.status) }, statusLabel(row.status))
  },
  {
    title: '失败原因',
    key: 'fail_reason',
    width: 140,
    ellipsis: { tooltip: true },
    render: (row: UserResourceItem) => h('span', { title: row.fail_reason || '', class: 'text-gray-500 cursor-help' }, row.fail_reason || '-')
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 160,
    render: (row: UserResourceItem) => h('span', { class: 'text-gray-500' }, formatDate(row.created_at))
  },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row: UserResourceItem) => {
      return h('button', {
        class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
        onClick: () => deleteUserResource(row),
        title: '删除该上传记录'
      }, [
        h('i', { class: 'fas fa-trash mr-1' }),
        '删除'
      ])
    }
  }
]

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

.config-content {
  padding: 1rem;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}
</style>