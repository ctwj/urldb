<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题和按钮 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">账号管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理平台账号信息</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="showCreateModal = true" type="primary">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加账号
        </n-button>
        <n-button @click="goToExpansionManagement" type="warning">
          <template #icon>
            <i class="fas fa-expand"></i>
          </template>
          账号扩容
        </n-button>
        <n-button @click="refreshData" type="info">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和筛选 -->
    <template #filter-bar>
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex flex-col md:flex-row gap-4">
          <n-input v-model:value="searchQuery" placeholder="搜索账号..." clearable class="flex-1">
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </n-input>

          <n-select v-model:value="platform" placeholder="选择平台" :options="platformOptions" clearable
            @update:value="onPlatformChange" class="w-full md:w-48" />

          <n-button type="primary" @click="handleSearch" class="w-full md:w-auto md:min-w-[100px]">
            <template #icon>
              <i class="fas fa-search"></i>
            </template>
            搜索
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区header - 账号列表头部 -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <span class="text-lg font-semibold text-gray-900 dark:text-white">账号列表</span>
        <div class="text-sm text-gray-500 dark:text-gray-400">
          共 {{ filteredCksList.length }} 个账号
        </div>
      </div>
    </template>

    <!-- 内容区content - 账号列表表格 -->
    <template #content>
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large">
          <template #description>
            <span class="text-gray-500">加载中...</span>
          </template>
        </n-spin>
      </div>

      <AdminErrorState
        v-else-if="errorMessage"
        icon="fas fa-exclamation-triangle"
        :message="errorMessage"
        :on-retry="refreshData"
      />

      <AdminEmptyState
        v-else-if="filteredCksList.length === 0"
        icon="fas fa-user-circle"
        title="暂无账号"
      >
        <template #action>
          <n-button @click="showCreateModal = true" type="primary">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加账号
          </n-button>
        </template>
      </AdminEmptyState>

      <!-- 账号列表和分页 -->
      <div v-else class="flex flex-col flex-1 h-full overflow-y-auto">
        <div
          v-for="item in filteredCksList"
          :key="item.id"
          class="border-b border-gray-200 dark:border-gray-700 p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-center justify-between">
            <!-- 左侧信息 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-4">
                <!-- ID -->
                <div class="w-16 text-sm font-medium text-gray-900 dark:text-gray-100">
                  #{{ item.id }}
                </div>

                <!-- 平台 -->
                <div class="flex items-center space-x-2">
                  <span v-html="getPlatformIcon(item.pan?.name || '')" class="text-lg"></span>
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {{ item.pan?.name || '未知平台' }}
                  </span>
                </div>

                <!-- 用户名 -->
                <div class="flex-1 min-w-0">
                  <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1"
                    :title="item.username || '未知用户'">
                    {{ item.username || '未知用户' }}
                  </h3>
                </div>
              </div>

              <!-- 状态和容量信息 -->
              <div class="mt-2 flex items-center space-x-4">
                <n-tag :type="item.is_valid ? 'success' : 'error'" size="small">
                  {{ item.is_valid ? '有效' : '无效' }}
                </n-tag>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  总空间: {{ formatFileSize(item.space) }}
                </span>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  已使用: {{ formatFileSize(Math.max(0, item.used_space || (item.space - item.left_space))) }}
                </span>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  剩余: {{ formatFileSize(Math.max(0, item.left_space)) }}
                </span>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  已转存: {{ item.transferred_count || 0 }}
                </span>
              </div>

              <!-- 备注 -->
              <div v-if="item.remark" class="mt-1">
                <span class="text-xs text-gray-600 dark:text-gray-400 line-clamp-1" :title="item.remark">
                  备注: {{ item.remark }}
                </span>
              </div>
            </div>

            <!-- 右侧操作按钮 -->
            <div class="flex items-center space-x-2 ml-4">
              <n-button size="small" :type="item.is_valid ? 'warning' : 'success'" @click="toggleStatus(item)"
                :title="item.is_valid ? '禁用账号' : '启用账号'" text>
                {{ item.is_valid ? '禁用' : '启用' }}
              </n-button>
              <n-button size="small" type="info" @click="refreshCapacity(item.id)" title="刷新容量" text>
                刷新容量
              </n-button>
              <n-button size="small" type="primary" @click="editCks(item)" title="编辑账号" text>
                编辑
              </n-button>
              <n-button size="small" type="error" @click="deleteCks(item.id)" title="删除账号" text>
                删除
              </n-button>
              <n-button size="small" type="warning" @click="deleteRelatedResources(item.id)" title="删除关联资源" text>
                删除关联
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 内容区footer - 分页组件 -->
    <template #content-footer>
      <div class="p-4">
        <div class="flex justify-center">
          <n-pagination v-model:page="currentPage" v-model:page-size="itemsPerPage" :item-count="filteredCksList.length"
            :page-sizes="[10, 20, 50, 100]" show-size-picker @update:page="goToPage"
            @update:page-size="(size) => { itemsPerPage = size; currentPage = 1; }" />
        </div>
      </div>
    </template>
  </AdminPageLayout>

  <!-- 创建/编辑账号模态框 -->
  <n-modal :show="showCreateModal || showEditModal" preset="card" title="账号管理" style="width: 500px"
    @update:show="(show) => { if (!show) closeModal() }">
    <template #header>
      <div class="flex items-center space-x-2">
        <i class="fas fa-user-circle text-lg"></i>
        <span>{{ showEditModal ? '编辑账号' : '添加账号' }}</span>
      </div>
    </template>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          平台类型 <span class="text-red-500">*</span>
        </label>
        <n-select v-model:value="form.pan_id" placeholder="请选择平台"
          :options="platforms.map(pan => ({ label: pan.remark, value: pan.id }))"
          :disabled="showEditModal" required />
        <p v-if="showEditModal" class="mt-1 text-xs text-gray-500">编辑时不允许修改平台类型</p>
      </div>

      <div v-if="showEditModal && editingCks?.username">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">用户名</label>
        <n-input :value="editingCks.username" disabled readonly />
      </div>

      <div v-if="isQuark || isBaidu">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Cookie <span class="text-red-500">*</span>
        </label>
        <n-input v-model:value="form.ck" type="textarea" placeholder="请输入Cookie内容，系统将自动识别容量" :rows="4" required />
        <n-alert v-if="isBaidu" type="warning" class="mt-2" :show-icon="true">
          请先登录百度网盘网页版，在根目录手动创建名为 <code class="px-1 bg-gray-100 dark:bg-gray-700 rounded">urldb</code> 的目录。系统转存的文件会保存到该目录。
        </n-alert>
      </div>

      <div v-if="isXunlei">
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">登录方式</label>
            <n-radio-group v-model:value="xunleiForm.loginMode">
              <n-radio value="password">账号密码</n-radio>
              <n-radio value="refresh_token">Refresh Token（推荐）</n-radio>
            </n-radio-group>
          </div>
          <template v-if="xunleiForm.loginMode === 'refresh_token'">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Refresh Token <span class="text-red-500">*</span>
              </label>
              <n-input v-model:value="xunleiForm.refreshToken" type="textarea" :rows="3" placeholder="从手机迅雷APP抓包获取（xluser-ssl.xunlei.com/v1/auth/token 响应的 refresh_token）" required />
              <n-alert type="info" class="mt-2" :show-icon="true">
                账号密码登录会触发迅雷短信验证（浏览器无法完成）。推荐用 refresh_token：手机迅雷APP 抓包获取后填入，永久免验证、自动续期。
              </n-alert>
            </div>
          </template>
          <template v-else>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                手机号 <span class="text-red-500">*</span>
              </label>
              <n-input v-model:value="xunleiForm.username" placeholder="请输入手机号（不需要+86前缀）" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                密码 <span class="text-red-500">*</span>
              </label>
              <n-input v-model:value="xunleiForm.password" type="password" placeholder="请输入密码" show-password-on="click" required />
            </div>
          </template>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">备注</label>
        <n-input v-model:value="form.remark" placeholder="可选，备注信息" />
      </div>

      <div v-if="showEditModal">
        <n-checkbox v-model:checked="form.is_valid">
          账号有效
        </n-checkbox>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <n-button type="tertiary" @click="closeModal">
          <template #icon>
            <i class="fas fa-times"></i>
          </template>
          取消
        </n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">
          <template #icon>
            <i class="fas fa-check"></i>
          </template>
          {{ showEditModal ? '更新' : '创建' }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  ssr: false
})

const isQuark = ref(false)
const isXunlei = ref(false)
const isBaidu = ref(false)

const notification = useNotification()
const router = useRouter()
const userStore = useUserStore()

const cksList = ref([])
const errorMessage = ref('')
const platforms = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingCks = ref(null)
const form = ref({
  pan_id: '',
  ck: '',
  is_valid: true,
  remark: ''
})

// 迅雷专用表单数据
const xunleiForm = ref({
  username: '',
  password: '',
  loginMode: 'password', // 'password' | 'refresh_token'
  refreshToken: ''
})

watch(() => form.value.pan_id, (newVal) => {
  isQuark.value = false
  isXunlei.value = false
  isBaidu.value = false
  const list = platforms.value.filter(it => it.id === newVal)
  if (!list || list.length === 0) {
    return
  }
  const pan = list[0]
  if (pan.name === 'quark') {
    isQuark.value = true
  } else if (pan.name === 'xunlei') {
    isXunlei.value = true
  } else if (pan.name === 'baidu') {
    isBaidu.value = true
  }
})

// 搜索和分页逻辑
const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(10)
const totalPages = ref(1)
const loading = ref(true)
const pageLoading = ref(true)
const submitting = ref(false)
const platform = ref(null)
const dialog = useDialog()

import { useCksApi, usePanApi } from '~/composables/useApi'
const cksApi = useCksApi()
const panApi = usePanApi()

const { data: pansData } = await useAsyncData('pans', () => panApi.getPans())
const pans = computed(() => {
  // 统一接口格式后直接为数组
  return Array.isArray(pansData.value) ? pansData.value : (pansData.value?.list || [])
})
const platformOptions = computed(() => {
  const options = [
    { label: '全部平台', value: null }
  ]

  pans.value.forEach(pan => {
    options.push({
      label: pan.remark || pan.name || `平台${pan.id}`,
      value: pan.id
    })
  })

  return options
})

// 检查认证
const checkAuth = () => {
  userStore.initAuth()
  if (!userStore.isAuthenticated) {
    router.push('/login')
    return
  }
}

// 获取账号列表
const fetchCks = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await cksApi.getCks()
    cksList.value = Array.isArray(response) ? response : []
  } catch (error) {
    errorMessage.value = '加载数据失败，请检查网络或后端服务'
    cksList.value = []
  } finally {
    loading.value = false
    pageLoading.value = false
  }
}

// 获取平台列表
const fetchPlatforms = async () => {
  try {
    const response = await panApi.getPans()
    platforms.value = Array.isArray(response) ? response : []
  } catch (error) {
  }
}

// 创建账号
const createCks = async () => {
  submitting.value = true
  try {
    await cksApi.createCks(form.value)
    await fetchCks()
    closeModal()
  } catch (error) {
    dialog.error({
      title: '错误',
      content: '创建账号失败: ' + (error.message || '未知错误'),
      positiveText: '确定'
    })
  } finally {
    submitting.value = false
  }
}

// 更新账号
const updateCks = async () => {
  submitting.value = true
  try {
    await cksApi.updateCks(editingCks.value.id, form.value)
    await fetchCks()
    closeModal()
  } catch (error) {
    notification.error({
      title: '失败',
      content: '更新账号失败: ' + (error.message || '未知错误'),
      duration: 3000
    })
  } finally {
    submitting.value = false
  }
}

// 删除账号
const deleteCks = async (id) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除这个账号吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await cksApi.deleteCks(id)
        await fetchCks()
      } catch (error) {
        notification.error({
          title: '失败',
          content: '删除账号失败: ' + (error.message || '未知错误'),
          duration: 3000
        })
      }
    }
  })
}

// 删除关联资源
const deleteRelatedResources = async (id) => {
  dialog.warning({
    title: '警告',
    content: '确定要删除与此账号关联的所有资源吗？这将清空这些资源的转存信息，变为未转存状态。',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        // 调用API删除关联资源
        await cksApi.deleteRelatedResources(id)
        await fetchCks()
        notification.success({
          title: '成功',
          content: '关联资源已删除！',
          duration: 3000
        })
      } catch (error) {
        notification.error({
          title: '失败',
          content: '删除关联资源失败: ' + (error.message || '未知错误'),
          duration: 3000
        })
      }
    }
  })
}

// 刷新容量
const refreshCapacity = async (id) => {
  dialog.warning({
    title: '警告',
    content: '确定要刷新此账号的容量信息吗？',
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await cksApi.refreshCapacity(id)
        await fetchCks()
        notification.success({
          title: '成功',
          content: '容量信息已刷新！',
          duration: 3000
        })
      } catch (error) {
        notification.error({
          title: '失败',
          content: '刷新容量失败: ' + (error.message || '未知错误'),
          duration: 3000
        })
      }
    }
  })
}

// 切换账号状态
const toggleStatus = async (cks) => {
  const newStatus = !cks.is_valid
  dialog.warning({
    title: '警告',
    content: `确定要${cks.is_valid ? '禁用' : '启用'}此账号吗？`,
    positiveText: '确定',
    negativeText: '取消',
    draggable: true,
    onPositiveClick: async () => {
      try {
        await cksApi.updateCks(cks.id, { is_valid: newStatus })
        await fetchCks()
        notification.success({
          title: '成功',
          content: `账号已${newStatus ? '启用' : '禁用'}！`,
          duration: 3000
        })
      } catch (error) {
        notification.error({
          title: '失败',
          content: `切换账号状态失败: ${error.message || '未知错误'}`,
          duration: 3000
        })
      }
    }
  })
}

// 编辑账号
const editCks = (cks) => {
  editingCks.value = cks
  form.value = {
    pan_id: cks.pan_id,
    ck: cks.ck,
    is_valid: cks.is_valid,
    remark: cks.remark || ''
  }

  // 如果是迅雷账号，解析ck字段来设置表单（支持账号密码 / refresh_token 两种模式回显）
  if (cks.pan?.name === 'xunlei') {
    try {
      const parsed = JSON.parse(cks.ck)
      if (parsed.refresh_token && !parsed.username) {
        // refresh_token 模式
        xunleiForm.value = {
          username: '',
          password: '',
          loginMode: 'refresh_token',
          refreshToken: parsed.refresh_token
        }
      } else {
        // 账号密码模式
        xunleiForm.value = {
          username: parsed.username || '',
          password: parsed.password || '',
          loginMode: 'password',
          refreshToken: ''
        }
      }
    } catch (e) {
      // 解析失败：ck 非 JSON（刷新后系统存的是纯 refresh_token 字符串），按 refresh_token 模式回显
      xunleiForm.value = {
        username: '',
        password: '',
        loginMode: 'refresh_token',
        refreshToken: cks.ck || ''
      }
    }
  }

  showEditModal.value = true
}

// 关闭模态框
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingCks.value = null
  form.value = {
    pan_id: '',
    ck: '',
    is_valid: true,
    remark: ''
  }
  // 重置迅雷表单
  xunleiForm.value = {
    username: '',
    password: '',
    loginMode: 'password',
    refreshToken: ''
  }
}

// 提交表单
const handleSubmit = async () => {
  // 如果是迅雷账号，按登录方式构造 CK JSON
  if (isXunlei.value) {
    if (xunleiForm.value.loginMode === 'refresh_token') {
      let token = (xunleiForm.value.refreshToken || '').trim()
      if (!token) {
        notification.error({
          title: '失败',
          content: '请填写 refresh_token',
          duration: 3000
        })
        return
      }
      // 容错：若填入的是整个 token 响应 JSON，提取其中的 refresh_token 字段
      if (token.startsWith('{')) {
        try {
          const parsed = JSON.parse(token)
          if (parsed && parsed.refresh_token) token = parsed.refresh_token
        } catch (e) { /* 非法 JSON，按原样作为 token 处理 */ }
      }
      form.value.ck = JSON.stringify({
        refresh_token: token
      })
    } else {
      if (!xunleiForm.value.username || !xunleiForm.value.password) {
        notification.error({
          title: '失败',
          content: '请填写完整的账号和密码',
          duration: 3000
        })
        return
      }
      form.value.ck = JSON.stringify({
        username: xunleiForm.value.username,
        password: xunleiForm.value.password
      })
    }
  }

  if (showEditModal.value) {
    await updateCks()
  } else {
    await createCks()
  }
}

// 获取平台图标
const getPlatformIcon = (platformName) => {
  const defaultIcons = {
    'unknown': '<i class="fas fa-question-circle text-gray-400"></i>',
    'other': '<i class="fas fa-cloud text-gray-500"></i>',
    'magnet': '<i class="fas fa-magnet text-red-600"></i>',
    'uc': '<i class="fas fa-cloud-download-alt text-purple-600"></i>',
    '夸克网盘': '<i class="fas fa-cloud text-blue-600"></i>',
    '阿里云盘': '<i class="fas fa-cloud text-orange-600"></i>',
    '百度网盘': '<i class="fas fa-cloud text-blue-500"></i>',
    '天翼云盘': '<i class="fas fa-cloud text-red-500"></i>',
    'OneDrive': '<i class="fas fa-cloud text-blue-700"></i>',
    'Google Drive': '<i class="fas fa-cloud text-green-600"></i>'
  }

  return defaultIcons[platformName] || defaultIcons['unknown']
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes <= 0) return '0 B'

  const tb = bytes / (1024 * 1024 * 1024 * 1024)
  if (tb >= 1) {
    return tb.toFixed(2) + ' TB'
  }

  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1) {
    return gb.toFixed(2) + ' GB'
  }

  const mb = bytes / (1024 * 1024)
  if (mb >= 1) {
    return mb.toFixed(2) + ' MB'
  }

  const kb = bytes / 1024
  if (kb >= 1) {
    return kb.toFixed(2) + ' KB'
  }

  return bytes + ' B'
}

// 过滤和分页计算
const filteredCksList = computed(() => {
  let filtered = cksList.value
  // 平台过滤
  if (platform.value !== null && platform.value !== undefined) {
    filtered = filtered.filter(cks => cks.pan_id === platform.value)
  }

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(cks =>
      cks.pan?.name?.toLowerCase().includes(query) ||
      cks.remark?.toLowerCase().includes(query)
    )
  }

  totalPages.value = Math.ceil(filtered.length / itemsPerPage.value)
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filtered.slice(start, end)
})

// 防抖搜索
let searchTimeout = null
const debounceSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
  }, 500)
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 平台变化处理
const onPlatformChange = () => {
  currentPage.value = 1
}

// 刷新数据
const refreshData = () => {
  currentPage.value = 1
  // 保持当前的过滤条件，只刷新数据
  fetchCks()
  fetchPlatforms()
}

// 分页跳转
const goToPage = (page) => {
  currentPage.value = page
}

// 跳转到扩容管理页面
const goToExpansionManagement = () => {
  router.push('/admin/accounts-expansion')
}

// 页面加载
onMounted(async () => {
  try {
    checkAuth()
    await Promise.all([
      fetchCks(),
      fetchPlatforms()
    ])
  } catch (error) {
  }
})
</script>

<style scoped>
/* 自定义样式 */
.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}

.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}
</style>