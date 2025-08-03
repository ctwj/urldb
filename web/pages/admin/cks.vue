<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-3 sm:p-5">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在加载账号数据</p>
          </div>
        </div>
      </div>
    </div>

    <div>
      <n-alert class="mb-4" title="平台账号管理当前只支持夸克" type="warning" />

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <div class="flex gap-2">
          <n-button 
            @click="showCreateModal = true"
            type="success"
          >
            <i class="fas fa-plus"></i> 添加账号
          </n-button>
        </div>
        <div class="flex gap-2">
          <div class="relative w-40">
            <n-select v-model:value="platform" :options="platformOptions" @update:value="onPlatformChange" />
          </div>
          <n-button 
            @click="refreshData"
            type="tertiary"
          >
            <i class="fas fa-refresh"></i> 刷新
          </n-button>
        </div>
      </div>

      <!-- 账号列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full min-w-full">
            <thead>
              <tr class="bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100">
                <th class="px-4 py-3 text-left text-sm font-medium">ID</th>
                <th class="px-4 py-3 text-left text-sm font-medium">平台</th>
                <th class="px-4 py-3 text-left text-sm font-medium">用户名</th>
                <th class="px-4 py-3 text-left text-sm font-medium">状态</th>
                <th class="px-4 py-3 text-left text-sm font-medium">总空间</th>
                <th class="px-4 py-3 text-left text-sm font-medium">已使用</th>
                <th class="px-4 py-3 text-left text-sm font-medium">剩余空间</th>
                <th class="px-4 py-3 text-left text-sm font-medium">备注</th>
                <th class="px-4 py-3 text-left text-sm font-medium">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-if="loading" class="text-center py-8">
                <td colspan="9" class="text-gray-500 dark:text-gray-400">
                  <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
                </td>
              </tr>
              <tr v-else-if="filteredCksList.length === 0" class="text-center py-8">
                <td colspan="9" class="text-gray-500 dark:text-gray-400">
                  <div class="flex flex-col items-center justify-center py-12">
                    <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 48 48">
                      <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
                      <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
                    </svg>
                    <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无账号</div>
                    <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">你可以点击上方"添加账号"按钮创建新账号</div>
                    <n-button 
                      @click="showCreateModal = true" 
                      type="primary"
                    >
                      <i class="fas fa-plus"></i> 添加账号
                    </n-button>
                  </div>
                </td>
              </tr>
              <tr 
                v-for="cks in filteredCksList" 
                :key="cks.id"
                class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100 font-medium">{{ cks.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  <div class="flex items-center">
                    <span v-html="getPlatformIcon(cks.pan?.name || '')" class="mr-2"></span>
                    {{ cks.pan?.name || '未知平台' }}
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  <span v-if="cks.username" :title="cks.username">{{ cks.username }}</span>
                  <span v-else class="text-gray-400 dark:text-gray-500 italic">未知用户</span>
                </td>
                <td class="px-4 py-3 text-sm">
                  <span :class="cks.is_valid ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-200'" 
                        class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ cks.is_valid ? '有效' : '无效' }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  {{ formatFileSize(cks.space) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  {{ formatFileSize(Math.max(0, cks.used_space || (cks.space - cks.left_space))) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  {{ formatFileSize(Math.max(0, cks.left_space)) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                  <span v-if="cks.remark" :title="cks.remark">{{ cks.remark }}</span>
                  <span v-else class="text-gray-400 dark:text-gray-500 italic">无备注</span>
                </td>
                <td class="px-4 py-3 text-sm">
                  <div class="flex items-center gap-2">
                    <button 
                      @click="toggleStatus(cks)"
                      :class="cks.is_valid ? 'text-orange-600 hover:text-orange-800 dark:text-orange-400 dark:hover:text-orange-300' : 'text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300'"
                      class="transition-colors"
                      :title="cks.is_valid ? '禁用账号' : '启用账号'"
                    >
                      <i :class="cks.is_valid ? 'fas fa-ban' : 'fas fa-check'"></i>
                    </button>
                    <button 
                      @click="refreshCapacity(cks.id)"
                      class="text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300 transition-colors"
                      title="刷新容量"
                    >
                      <i class="fas fa-sync-alt"></i>
                    </button>
                    <button 
                      @click="editCks(cks)"
                      class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
                      title="编辑账号"
                    >
                      <i class="fas fa-edit"></i>
                    </button>
                    <button 
                      @click="deleteCks(cks.id)"
                      class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 transition-colors"
                      title="删除账号"
                    >
                      <i class="fas fa-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex flex-wrap justify-center gap-1 sm:gap-2 mt-6">
        <button 
          v-if="currentPage > 1"
          @click="goToPage(currentPage - 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          <i class="fas fa-chevron-left mr-1"></i> 上一页
        </button>
        
        <button 
          @click="goToPage(1)"
          :class="currentPage === 1 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          1
        </button>
        
        <button 
          v-if="totalPages > 1"
          @click="goToPage(2)"
          :class="currentPage === 2 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          2
        </button>
        
        <span v-if="currentPage > 2" class="px-2 py-1 sm:px-3 sm:py-2 text-gray-500 text-sm">...</span>
        
        <button 
          v-if="currentPage !== 1 && currentPage !== 2 && currentPage > 2"
          class="bg-slate-800 text-white px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          {{ currentPage }}
        </button>
        
        <button 
          v-if="currentPage < totalPages"
          @click="goToPage(currentPage + 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          下一页 <i class="fas fa-chevron-right ml-1"></i>
        </button>
      </div>

      <!-- 统计信息 -->
      <div v-if="totalPages <= 1" class="mt-4 text-center">
        <div class="inline-flex items-center bg-white dark:bg-gray-800 rounded-lg shadow px-6 py-3">
          <div class="text-sm text-gray-600 dark:text-gray-400">
            共 <span class="font-semibold text-gray-900 dark:text-gray-100">{{ filteredCksList.length }}</span> 个账号
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑账号模态框 -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white dark:bg-gray-800">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4">
            {{ showEditModal ? '编辑账号' : '添加账号' }}
          </h3>
          <form @submit.prevent="handleSubmit">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">平台类型 <span class="text-red-500">*</span></label>
                <select 
                  v-model="form.pan_id" 
                  required
                  :disabled="showEditModal"
                  :class="showEditModal ? 'bg-gray-100 dark:bg-gray-600 cursor-not-allowed' : ''"
                  class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                >
                  <option value="">请选择平台</option>
                  <option v-for="pan in platforms.filter(pan => pan.name === 'quark')" :key="pan.id" :value="pan.id">
                    {{ pan.remark }}
                  </option>
                </select>
                <p v-if="showEditModal" class="mt-1 text-xs text-gray-500 dark:text-gray-400">编辑时不允许修改平台类型</p>
              </div>
              <div v-if="showEditModal && editingCks?.username">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">用户名</label>
                <div class="mt-1 px-3 py-2 bg-gray-100 dark:bg-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100">
                  {{ editingCks.username }}
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Cookie <span class="text-red-500">*</span></label>
                <n-input 
                  v-model:value="form.ck" 
                  required
                  type="textarea"
                  placeholder="请输入Cookie内容，系统将自动识别容量"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">备注</label>
                <n-input 
                  v-model:value="form.remark" 
                  type="text" 
                  placeholder="可选，备注信息"
                />
              </div>
              <div v-if="showEditModal">
                <label class="flex items-center">
                  <n-checkbox 
                    v-model:checked="form.is_valid"
                  />
                  <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">账号有效</span>
                </label>
              </div>
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <n-button 
                type="tertiary" 
                @click="closeModal"
              >
                取消
              </n-button>
              <n-button 
                type="primary" 
                :disabled="submitting"
                @click="handleSubmit"
              >
                {{ submitting ? '处理中...' : (showEditModal ? '更新' : '创建') }}
              </n-button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  ssr: false
})

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
    alert('更新账号失败: ' + (error.message || '未知错误'))
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
        alert('删除账号失败: ' + (error.message || '未知错误'))
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
        alert('容量信息已刷新！')
      } catch (error) {
        console.error('刷新容量失败:', error)
        alert('刷新容量失败: ' + (error.message || '未知错误'))
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
        alert(`账号已${newStatus ? '启用' : '禁用'}！`)
      } catch (error) {
        console.error('切换账号状态失败:', error)
        alert(`切换账号状态失败: ${error.message || '未知错误'}`)
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