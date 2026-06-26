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

      <div v-if="isQuark || isBaidu || isUC">
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
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">客户端类型</label>
            <n-radio-group v-model:value="xunleiForm.clientType">
              <n-radio value="android">安卓（迅雷下载管家 APP）— refresh_token</n-radio>
              <n-radio value="browser" disabled>浏览器（迅雷浏览器 APP）— 暂未启用</n-radio>
            </n-radio-group>
            <p class="text-xs text-gray-500 mt-1">目前仅支持安卓方案：用手机迅雷下载管家 APP 抓包获取 refresh_token 填入，永久免验证、自动续期。</p>
          </div>

          <!-- 安卓（下载管家）：refresh_token -->
          <div v-if="xunleiForm.clientType === 'android'">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Refresh Token <span class="text-red-500">*</span>
            </label>
            <n-input v-model:value="xunleiForm.refreshToken" type="textarea" :rows="3" placeholder="从手机迅雷下载管家 APP 抓包获取（xluser-ssl.xunlei.com/v1/auth/token 响应的 refresh_token）" required />
            <n-alert type="info" class="mt-2" :show-icon="true">
              用手机迅雷下载管家 APP 抓包获取 refresh_token 填入，永久免验证、自动续期。
            </n-alert>
          </div>

          <!-- 浏览器（迅雷浏览器APP）：账号密码 + creditkey 闭环 -->
          <div v-else>
            <div class="space-y-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  账号（手机号） <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="xunleiForm.username" placeholder="迅雷账号手机号（不含 +86 前缀）" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  密码 <span class="text-red-500">*</span>
                </label>
                <n-input v-model:value="xunleiForm.password" type="password" show-password-on="click" :placeholder="showEditModal ? '迅雷账号密码（留空表示不修改）' : '迅雷账号密码'" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Creditkey <span class="text-xs text-gray-400 font-normal">（可选，触发验证后填写）</span>
                </label>
                <n-input v-model:value="xunleiForm.creditkey" placeholder="触发安全验证后系统会自动填入；也可手动粘贴" />
              </div>
            </div>

            <!-- review 安全验证提示区 -->
            <n-alert v-if="xunleiForm.reviewUrl" type="warning" class="mt-3" :show-icon="true" closable @close="xunleiForm.reviewUrl = ''">
              <div class="font-medium">本次登录需要安全验证</div>
              <div class="mt-1">1. 打开
                <a :href="xunleiForm.reviewUrl" target="_blank" class="text-blue-600 underline">验证链接</a>
                完成短信验证（或在迅雷官方端登录该账号一次）。</div>
              <div>2. 系统已自动填入 creditkey，确认后点击「创建」重新提交即可。</div>
            </n-alert>

            <n-alert type="info" class="mt-2" :show-icon="true">
              迅雷浏览器 APP 用账号密码登录。新设备首次登录大概率触发短信验证，按提示完成验证并填入 creditkey 后重新提交即可绕过。
            </n-alert>
          </div>
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
const isUC = ref(false)

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
  clientType: 'android', // 'android'（下载管家，refresh_token）| 'browser'（迅雷浏览器APP，账号密码）
  refreshToken: '',
  username: '',
  password: '',
  creditkey: '',
  reviewUrl: ''
})

watch(() => form.value.pan_id, (newVal) => {
  isQuark.value = false
  isXunlei.value = false
  isBaidu.value = false
  isUC.value = false
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
  } else if (pan.name === 'uc') {
    isUC.value = true
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
    const result = await cksApi.createCks(form.value)
    // 迅雷浏览器账号密码登录触发 review：作为业务状态（HTTP 200）返回，走 creditkey 闭环（不关闭弹窗）
    if (result && result.need_review) {
      xunleiForm.value.creditkey = result.creditkey || xunleiForm.value.creditkey
      xunleiForm.value.reviewUrl = result.review_url || ''
      dialog.warning({
        title: '需要安全验证',
        content: '请在验证链接完成短信验证（或在迅雷官方端登录该账号一次），creditkey 已自动填入，确认后点击「创建」重新提交。',
        positiveText: '打开验证链接',
        onPositiveClick: () => {
          if (xunleiForm.value.reviewUrl) {
            window.open(xunleiForm.value.reviewUrl, '_blank')
          }
        }
      })
      return
    }
    await fetchCks()
    closeModal()
  } catch (error) {
    const respData = error && error.data
    dialog.error({
      title: '错误',
      content: '创建账号失败: ' + (respData?.message || error.message || '未知错误'),
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

  // 如果是迅雷账号，按客户端类型回显
  if (cks.pan?.name === 'xunlei') {
    let clientType = 'android'
    let refreshToken = ''
    let username = ''
    let password = ''
    let creditkey = ''
    try {
      const parsed = JSON.parse(cks.ck)
      clientType = parsed.client_type || 'android'
      refreshToken = parsed.refresh_token || ''
      username = parsed.username || ''
      creditkey = parsed.creditkey || ''
      // 密码留空表示不修改（编辑时不强制重填）
    } catch (e) {
      // ck 非 JSON（刷新后系统存纯 refresh_token 字符串），默认 android
      refreshToken = cks.ck || ''
    }
    xunleiForm.value = {
      clientType,
      refreshToken,
      username,
      password,
      creditkey,
      reviewUrl: ''
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
    clientType: 'android',
    refreshToken: '',
    username: '',
    password: '',
    creditkey: '',
    reviewUrl: ''
  }
}

// 提交表单
const handleSubmit = async () => {
  // 如果是迅雷账号，按客户端类型构造 CK JSON
  if (isXunlei.value) {
    if (xunleiForm.value.clientType === 'android') {
      let token = (xunleiForm.value.refreshToken || '').trim()
      if (!token) {
        notification.error({ title: '失败', content: '请填写 refresh_token', duration: 3000 })
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
        refresh_token: token,
        client_type: 'android'
      })
    } else {
      // 浏览器（迅雷浏览器APP）：账号密码
      const username = (xunleiForm.value.username || '').trim()
      const password = xunleiForm.value.password || ''
      if (!username) {
        notification.error({ title: '失败', content: '请填写账号（手机号）', duration: 3000 })
        return
      }
      if (showEditModal.value && !password) {
        // 编辑时密码留空：保留原 ck 不变（仅更新备注等字段，不重新登录）
      } else {
        if (!password) {
          notification.error({ title: '失败', content: '请填写密码', duration: 3000 })
          return
        }
        form.value.ck = JSON.stringify({
          username,
          password,
          creditkey: (xunleiForm.value.creditkey || '').trim(),
          client_type: 'browser'
        })
      }
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