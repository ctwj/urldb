<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 p-3 sm:p-5">
    <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6">
        <div class="flex justify-between items-center">
          <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">第三方平台账号管理</h1>
          <div class="flex gap-2">
            <NuxtLink 
              to="/admin" 
              class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700"
            >
              返回管理
            </NuxtLink>
            <button 
              @click="showCreateModal = true" 
              class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              添加账号
            </button>
          </div>
        </div>
      </div>

      <!-- 账号列表 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">账号列表</h2>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead class="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">平台</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">索引</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">状态</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">总空间</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">剩余空间</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">备注</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="cks in cksList" :key="cks.id">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">{{ cks.id }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                  <div class="flex items-center">
                    <span v-html="getPlatformIcon(cks.pan?.name || '')" class="mr-2"></span>
                    {{ cks.pan?.name || '未知平台' }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">{{ cks.idx || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="cks.is_valid ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-200'" 
                        class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ cks.is_valid ? '有效' : '无效' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                  {{ formatFileSize(cks.space) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                  {{ formatFileSize(cks.left_space) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">{{ cks.remark || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button @click="editCks(cks)" class="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 mr-3">编辑</button>
                  <button @click="deleteCks(cks.id)" class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
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
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">平台 <span class="text-red-500">*</span></label>
                  <select 
                    v-model="form.pan_id" 
                    required
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                  >
                    <option value="">请选择平台</option>
                    <option v-for="pan in platforms" :key="pan.id" :value="pan.id">
                      {{ pan.name }}
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Cookie <span class="text-red-500">*</span></label>
                  <textarea 
                    v-model="form.ck" 
                    required
                    rows="4"
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                    placeholder="请输入Cookie内容"
                  ></textarea>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">索引</label>
                  <input 
                    v-model.number="form.idx" 
                    type="number" 
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                    placeholder="可选，账号索引"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">总空间 (GB)</label>
                  <input 
                    v-model.number="form.space" 
                    type="number" 
                    step="0.01"
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                    placeholder="可选，总空间大小"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">剩余空间 (GB)</label>
                  <input 
                    v-model.number="form.left_space" 
                    type="number" 
                    step="0.01"
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                    placeholder="可选，剩余空间大小"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">备注</label>
                  <input 
                    v-model="form.remark" 
                    type="text" 
                    class="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-gray-100"
                    placeholder="可选，备注信息"
                  />
                </div>
                <div>
                  <label class="flex items-center">
                    <input 
                      v-model="form.is_valid" 
                      type="checkbox" 
                      class="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 dark:bg-gray-700 dark:border-gray-600"
                    />
                    <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">账号有效</span>
                  </label>
                </div>
              </div>
              <div class="mt-6 flex justify-end space-x-3">
                <button 
                  type="button" 
                  @click="closeModal"
                  class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
                >
                  取消
                </button>
                <button 
                  type="submit" 
                  class="px-4 py-2 bg-indigo-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-indigo-700"
                >
                  {{ showEditModal ? '更新' : '创建' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
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
  idx: null,
  ck: '',
  is_valid: true,
  space: null,
  left_space: null,
  remark: ''
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
  try {
    const { useCksApi } = await import('~/composables/useApi')
    const cksApi = useCksApi()
    const response = await cksApi.getCks()
    cksList.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('获取账号列表失败:', error)
  }
}

// 获取平台列表
const fetchPlatforms = async () => {
  try {
    const { usePanApi } = await import('~/composables/useApi')
    const panApi = usePanApi()
    const response = await panApi.getPans()
    platforms.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('获取平台列表失败:', error)
  }
}

// 创建账号
const createCks = async () => {
  try {
    const { useCksApi } = await import('~/composables/useApi')
    const cksApi = useCksApi()
    await cksApi.createCks(form.value)
    await fetchCks()
    closeModal()
  } catch (error) {
    console.error('创建账号失败:', error)
    alert('创建账号失败: ' + (error.message || '未知错误'))
  }
}

// 更新账号
const updateCks = async () => {
  try {
    const { useCksApi } = await import('~/composables/useApi')
    const cksApi = useCksApi()
    await cksApi.updateCks(editingCks.value.id, form.value)
    await fetchCks()
    closeModal()
  } catch (error) {
    console.error('更新账号失败:', error)
    alert('更新账号失败: ' + (error.message || '未知错误'))
  }
}

// 删除账号
const deleteCks = async (id) => {
  if (!confirm('确定要删除这个账号吗？')) return
  
  try {
    const { useCksApi } = await import('~/composables/useApi')
    const cksApi = useCksApi()
    await cksApi.deleteCks(id)
    await fetchCks()
  } catch (error) {
    console.error('删除账号失败:', error)
    alert('删除账号失败: ' + (error.message || '未知错误'))
  }
}

// 编辑账号
const editCks = (cks) => {
  editingCks.value = cks
  form.value = {
    pan_id: cks.pan_id,
    idx: cks.idx,
    ck: cks.ck,
    is_valid: cks.is_valid,
    space: cks.space,
    left_space: cks.left_space,
    remark: cks.remark
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
    idx: null,
    ck: '',
    is_valid: true,
    space: null,
    left_space: null,
    remark: ''
  }
}

// 提交表单
const handleSubmit = () => {
  if (showEditModal.value) {
    updateCks()
  } else {
    createCks()
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
  if (!bytes || bytes === 0) return '-'
  
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

onMounted(() => {
  checkAuth()
  fetchCks()
  fetchPlatforms()
})
</script> 