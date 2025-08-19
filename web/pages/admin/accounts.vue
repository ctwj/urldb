<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
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
        <n-button @click="refreshData" type="info">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <n-card>
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
    </n-card>

    <!-- 账号列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">账号列表</span>
          <div class="text-sm text-gray-500">
            共 {{ filteredCksList.length }} 个账号
          </div>
        </div>
      </template>

      <!-- 加载状态 -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <n-spin size="large">
          <template #description>
            <span class="text-gray-500">加载中...</span>
          </template>
        </n-spin>
      </div>

      <!-- 空状态 -->
      <div v-else-if="filteredCksList.length === 0" class="flex flex-col items-center justify-center py-12">
        <n-empty description="暂无账号">
          <template #icon>
            <i class="fas fa-user-circle text-4xl text-gray-400"></i>
          </template>
          <template #extra>
            <n-button @click="showCreateModal = true" type="primary">
              <template #icon>
                <i class="fas fa-plus"></i>
              </template>
              添加账号
            </n-button>
          </template>
        </n-empty>
      </div>

      <!-- 账号列表 -->
      <div v-else>
        <n-virtual-list :items="filteredCksList" :item-size="100" style="max-height: 500px">
          <template #default="{ item }">
            <div
              class="border-b border-gray-200 dark:border-gray-700 p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
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
                    <span class="text-xs text-gray-500">
                      总空间: {{ formatFileSize(item.space) }}
                    </span>
                    <span class="text-xs text-gray-500">
                      已使用: {{ formatFileSize(Math.max(0, item.used_space || (item.space - item.left_space))) }}
                    </span>
                    <span class="text-xs text-gray-500">
                      剩余: {{ formatFileSize(Math.max(0, item.left_space)) }}
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
                    :title="item.is_valid ? '禁用账号' : '启用账号'">
                    <template #icon>
                      <i :class="item.is_valid ? 'fas fa-ban' : 'fas fa-check'"></i>
                    </template>
                  </n-button>
                  <n-button size="small" type="info" @click="refreshCapacity(item.id)" title="刷新容量">
                    <template #icon>
                      <i class="fas fa-sync-alt"></i>
                    </template>
                  </n-button>
                  <n-button size="small" type="primary" @click="editCks(item)" title="编辑账号">
                    <template #icon>
                      <i class="fas fa-edit"></i>
                    </template>
                  </n-button>
                  <n-button size="small" type="error" @click="deleteCks(item.id)" title="删除账号">
                    <template #icon>
                      <i class="fas fa-trash"></i>
                    </template>
                  </n-button>
                </div>
              </div>
            </div>
          </template>
        </n-virtual-list>
      </div>
    </n-card>

    <!-- 分页 -->
    <div class="flex justify-center">
      <n-pagination v-model:page="currentPage" v-model:page-size="itemsPerPage" :item-count="filteredCksList.length"
        :page-sizes="[10, 20, 50, 100]" show-size-picker @update:page="goToPage"
        @update:page-size="(size) => { itemsPerPage = size; currentPage = 1; }" />
    </div>
  </div>

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
          :options="platforms.filter(pan => pan.name === 'quark').map(pan => ({ label: pan.remark, value: pan.id }))"
          :disabled="showEditModal" required />
        <p v-if="showEditModal" class="mt-1 text-xs text-gray-500">编辑时不允许修改平台类型</p>
      </div>

      <div v-if="showEditModal && editingCks?.username">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">用户名</label>
        <n-input :value="editingCks.username" disabled readonly />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Cookie <span class="text-red-500">*</span>
        </label>
        <n-input v-model:value="form.ck" type="textarea" placeholder="请输入Cookie内容，系统将自动识别容量" :rows="4" required />
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
          取消
        </n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">
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

const notification = useNotification()
const router = useRouter()
const userStore = useUserStore()

const cksList = ref([])
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
  try {
    console.log('开始获取账号列表...')
    const response = await cksApi.getCks()
    cksList.value = Array.isArray(response) ? response : []
    console.log('获取账号列表成功，数据:', cksList.value)
  } catch (error) {
    console.error('获取账号列表失败:', error)
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
    console.error('获取平台列表失败:', error)
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
    console.error('更新账号失败:', error)
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
        console.error('删除账号失败:', error)
        notification.error({
          title: '失败',
          content: '删除账号失败: ' + (error.message || '未知错误'),
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
        console.error('刷新容量失败:', error)
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
        console.log('切换状态 - 账号ID:', cks.id, '当前状态:', cks.is_valid, '新状态:', newStatus)
        await cksApi.updateCks(cks.id, { is_valid: newStatus })
        console.log('状态更新成功，正在刷新数据...')
        await fetchCks()
        console.log('数据刷新完成')
        notification.success({
          title: '成功',
          content: `账号已${newStatus ? '启用' : '禁用'}！`,
          duration: 3000
        })
      } catch (error) {
        console.error('切换账号状态失败:', error)
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
}

// 提交表单
const handleSubmit = async () => {
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
  console.log('原始账号数量:', filtered.length)

  // 平台过滤
  if (platform.value !== null && platform.value !== undefined) {
    filtered = filtered.filter(cks => cks.pan_id === platform.value)
    console.log('平台过滤后数量:', filtered.length, '平台ID:', platform.value)
  }

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(cks =>
      cks.pan?.name?.toLowerCase().includes(query) ||
      cks.remark?.toLowerCase().includes(query)
    )
    console.log('搜索过滤后数量:', filtered.length, '搜索词:', searchQuery.value)
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
  console.log('执行搜索，搜索词:', searchQuery.value)
  console.log('当前过滤后的账号数量:', filteredCksList.value.length)
}

// 平台变化处理
const onPlatformChange = () => {
  currentPage.value = 1
  console.log('平台过滤条件变化:', platform.value)
  console.log('当前过滤后的账号数量:', filteredCksList.value.length)
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

// 页面加载
onMounted(async () => {
  try {
    checkAuth()
    await Promise.all([
      fetchCks(),
      fetchPlatforms()
    ])
  } catch (error) {
    console.error('页面初始化失败:', error)
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