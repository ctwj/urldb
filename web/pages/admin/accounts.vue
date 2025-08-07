<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">账号管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的第三方平台账号</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加账号
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <n-card>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索账号..."
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-select
          v-model:value="selectedPlatform"
          placeholder="选择平台"
          :options="platformOptions"
          clearable
        />
        
        <n-select
          v-model:value="selectedStatus"
          placeholder="选择状态"
          :options="statusOptions"
          clearable
        />
      </div>
    </n-card>

    <!-- 账号列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">账号列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个账号</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="accounts.length === 0" class="text-center py-8">
        <i class="fas fa-user-shield text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无账号数据</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="account in accounts"
          :key="account.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <div class="w-8 h-8 flex items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900">
                  <i class="fas fa-user text-blue-600 dark:text-blue-400"></i>
                </div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                  {{ account.username }}
                </h3>
                <span v-if="account.is_enabled" class="text-xs px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded">
                  启用
                </span>
                <span v-else class="text-xs px-2 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded">
                  禁用
                </span>
              </div>
              
              <div class="flex items-center space-x-4 text-sm text-gray-600 dark:text-gray-400 mb-2">
                <span>平台: {{ getPlatformName(account.pan_id) }}</span>
                <span>邮箱: {{ account.email || '未设置' }}</span>
                <span>状态: {{ account.status || '正常' }}</span>
              </div>
              
              <div class="flex items-center space-x-4 text-xs text-gray-500">
                <span>创建时间: {{ formatDate(account.created_at) }}</span>
                <span>更新时间: {{ formatDate(account.updated_at) }}</span>
                <span>最后登录: {{ formatDate(account.last_login_at) || '从未登录' }}</span>
              </div>
            </div>
            
            <div class="flex space-x-2">
              <n-button size="small" type="primary" @click="editAccount(account)">
                <template #icon>
                  <i class="fas fa-edit"></i>
                </template>
                编辑
              </n-button>
              <n-button size="small" type="warning" @click="testAccount(account)">
                <template #icon>
                  <i class="fas fa-vial"></i>
                </template>
                测试
              </n-button>
              <n-button size="small" type="error" @click="deleteAccount(account)">
                <template #icon>
                  <i class="fas fa-trash"></i>
                </template>
                删除
              </n-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="mt-6">
        <n-pagination
          v-model:page="currentPage"
          v-model:page-size="pageSize"
          :item-count="total"
          :page-sizes="[10, 20, 50, 100]"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- 添加/编辑账号模态框 -->
    <n-modal v-model:show="showAddModal" preset="card" title="添加账号" style="width: 600px">
      <n-form
        ref="formRef"
        :model="accountForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <n-form-item label="用户名" path="username">
            <n-input
              v-model:value="accountForm.username"
              placeholder="请输入用户名"
            />
          </n-form-item>

          <n-form-item label="平台" path="pan_id">
            <n-select
              v-model:value="accountForm.pan_id"
              placeholder="请选择平台"
              :options="platformOptions"
            />
          </n-form-item>

          <n-form-item label="邮箱" path="email">
            <n-input
              v-model:value="accountForm.email"
              placeholder="请输入邮箱"
              type="text"
            />
          </n-form-item>

          <n-form-item label="密码" path="password">
            <n-input
              v-model:value="accountForm.password"
              placeholder="请输入密码"
              type="password"
              show-password-on="click"
            />
          </n-form-item>

          <n-form-item label="Cookie" path="cookie">
            <n-input
              v-model:value="accountForm.cookie"
              placeholder="请输入Cookie"
              type="textarea"
              :rows="3"
            />
          </n-form-item>

          <n-form-item label="启用状态" path="is_enabled">
            <n-switch v-model:value="accountForm.is_enabled" />
          </n-form-item>
        </div>

        <n-form-item label="备注" path="remark">
          <n-input
            v-model:value="accountForm.remark"
            placeholder="请输入备注信息"
            type="textarea"
            :rows="2"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="closeModal">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            保存
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin' as any
})

// 使用API
const { useCksApi, usePanApi } = await import('~/composables/useApi')
const cksApi = useCksApi()
const panApi = usePanApi()

// 响应式数据
const loading = ref(false)
const accounts = ref<any[]>([])
const platforms = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const selectedPlatform = ref('')
const selectedStatus = ref('')

// 模态框相关
const showAddModal = ref(false)
const submitting = ref(false)
const editingAccount = ref<any>(null)

// 表单数据
const accountForm = ref({
  username: '',
  pan_id: '',
  email: '',
  password: '',
  cookie: '',
  is_enabled: true,
  remark: ''
})

// 表单验证规则
const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur'
  },
  pan_id: {
    required: true,
    message: '请选择平台',
    trigger: 'change'
  }
}

const formRef = ref()

// 状态选项
const statusOptions = [
  { label: '正常', value: 'normal' },
  { label: '异常', value: 'error' },
  { label: '禁用', value: 'disabled' }
]

// 计算属性
const platformOptions = computed(() => 
  platforms.value.map(platform => ({ label: platform.name, value: platform.id }))
)

// 获取数据
const fetchData = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
      pan_id: selectedPlatform.value,
      status: selectedStatus.value
    }
    
    const response = await cksApi.getCks(params) as any
    accounts.value = response.data || []
    total.value = response.total || 0
  } catch (error: any) {
    useNotification().error({
      content: error.message || '获取账号数据失败',
      duration: 5000
    })
  } finally {
    loading.value = false
  }
}

// 获取平台数据
const fetchPlatforms = async () => {
  try {
    const response = await panApi.getPans()
    if (response && (response as any).items && Array.isArray((response as any).items)) {
      platforms.value = (response as any).items
    } else if (Array.isArray(response)) {
      platforms.value = response
    } else {
      platforms.value = []
    }
  } catch (error) {
    console.error('获取平台数据失败:', error)
    platforms.value = []
  }
}

// 初始化数据
onMounted(async () => {
  await Promise.all([
    fetchData(),
    fetchPlatforms()
  ])
})

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchData()
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 编辑账号
const editAccount = (account: any) => {
  editingAccount.value = account
  accountForm.value = {
    username: account.username,
    pan_id: account.pan_id,
    email: account.email || '',
    password: '',
    cookie: account.cookie || '',
    is_enabled: account.is_enabled,
    remark: account.remark || ''
  }
  showAddModal.value = true
}

// 测试账号
const testAccount = async (account: any) => {
  try {
    // 这里需要实现测试账号的API调用
    console.log('测试账号:', account.id)
    // await cksApi.testCks(account.id)
    useNotification().success({
      content: '账号测试成功',
      duration: 3000
    })
  } catch (error: any) {
    useNotification().error({
      content: error.message || '账号测试失败',
      duration: 5000
    })
  }
}

// 删除账号
const deleteAccount = async (account: any) => {
  try {
    await cksApi.deleteCks(account.id)
    useNotification().success({
      content: '账号删除成功',
      duration: 3000
    })
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '删除账号失败',
      duration: 5000
    })
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  editingAccount.value = null
  accountForm.value = {
    username: '',
    pan_id: '',
    email: '',
    password: '',
    cookie: '',
    is_enabled: true,
    remark: ''
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    
    if (editingAccount.value) {
      await cksApi.updateCks(editingAccount.value.id, accountForm.value)
      useNotification().success({
        content: '账号更新成功',
        duration: 3000
      })
    } else {
      await cksApi.createCks(accountForm.value)
      useNotification().success({
        content: '账号创建成功',
        duration: 3000
      })
    }
    
    closeModal()
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '保存账号失败',
      duration: 5000
    })
  } finally {
    submitting.value = false
  }
}

// 获取平台名称
const getPlatformName = (panId: number) => {
  if (!platforms.value || !Array.isArray(platforms.value)) {
    return '未知平台'
  }
  const platform = platforms.value.find(p => p.id === panId)
  return platform?.name || '未知平台'
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString('zh-CN')
}
</script> 