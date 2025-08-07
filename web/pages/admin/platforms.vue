<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">平台管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的网盘平台</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加平台
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 搜索 -->
    <n-card>
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索平台..."
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <i class="fas fa-search"></i>
        </template>
      </n-input>
    </n-card>

    <!-- 平台列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">平台列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个平台</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="platforms.length === 0" class="text-center py-8">
        <i class="fas fa-cloud text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无平台数据</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="platform in platforms"
          :key="platform.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <div class="w-8 h-8 flex items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900">
                  <span v-html="platform.icon" class="text-blue-600 dark:text-blue-400"></span>
                </div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                  {{ platform.name }}
                </h3>
                <span v-if="platform.is_enabled" class="text-xs px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded">
                  启用
                </span>
                <span v-else class="text-xs px-2 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded">
                  禁用
                </span>
              </div>
              
              <p v-if="platform.description" class="text-sm text-gray-600 dark:text-gray-400 mb-2">
                {{ platform.description }}
              </p>
              
              <div class="flex items-center space-x-4 text-xs text-gray-500">
                <span>资源数量: {{ platform.resource_count || 0 }}</span>
                <span>创建时间: {{ formatDate(platform.created_at) }}</span>
                <span>更新时间: {{ formatDate(platform.updated_at) }}</span>
              </div>
            </div>
            
            <div class="flex space-x-2">
              <n-button size="small" type="primary" @click="editPlatform(platform)">
                <template #icon>
                  <i class="fas fa-edit"></i>
                </template>
                编辑
              </n-button>
              <n-button size="small" type="error" @click="deletePlatform(platform)">
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

    <!-- 添加/编辑平台模态框 -->
    <n-modal v-model:show="showAddModal" preset="card" title="添加平台" style="width: 500px">
      <n-form
        ref="formRef"
        :model="platformForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="平台名称" path="name">
          <n-input
            v-model:value="platformForm.name"
            placeholder="请输入平台名称"
          />
        </n-form-item>

        <n-form-item label="平台图标" path="icon">
          <n-input
            v-model:value="platformForm.icon"
            placeholder="请输入平台图标HTML"
          />
          <template #help>
            支持HTML格式的图标，如：&lt;i class="fas fa-cloud"&gt;&lt;/i&gt;
          </template>
        </n-form-item>

        <n-form-item label="平台描述" path="description">
          <n-input
            v-model:value="platformForm.description"
            placeholder="请输入平台描述"
            type="textarea"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="启用状态" path="is_enabled">
          <n-switch v-model:value="platformForm.is_enabled" />
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
const { usePanApi } = await import('~/composables/useApi')
const panApi = usePanApi()

// 响应式数据
const loading = ref(false)
const platforms = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')

// 模态框相关
const showAddModal = ref(false)
const submitting = ref(false)
const editingPlatform = ref<any>(null)

// 表单数据
const platformForm = ref({
  name: '',
  icon: '',
  description: '',
  is_enabled: true
})

// 表单验证规则
const rules = {
  name: {
    required: true,
    message: '请输入平台名称',
    trigger: 'blur'
  },
  icon: {
    required: true,
    message: '请输入平台图标',
    trigger: 'blur'
  }
}

const formRef = ref()

// 获取数据
const fetchData = async () => {
  try {
    loading.value = true
    
    const response = await panApi.getPans() as any
    platforms.value = response.data || []
    total.value = platforms.value.length
  } catch (error: any) {
    useNotification().error({
      content: error.message || '获取平台数据失败',
      duration: 5000
    })
  } finally {
    loading.value = false
  }
}

// 初始化数据
onMounted(() => {
  fetchData()
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

// 编辑平台
const editPlatform = (platform: any) => {
  editingPlatform.value = platform
  platformForm.value = {
    name: platform.name,
    icon: platform.icon,
    description: platform.description || '',
    is_enabled: platform.is_enabled
  }
  showAddModal.value = true
}

// 删除平台
const deletePlatform = async (platform: any) => {
  try {
    await panApi.deletePan(platform.id)
    useNotification().success({
      content: '平台删除成功',
      duration: 3000
    })
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '删除平台失败',
      duration: 5000
    })
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  editingPlatform.value = null
  platformForm.value = {
    name: '',
    icon: '',
    description: '',
    is_enabled: true
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    
    if (editingPlatform.value) {
      await panApi.updatePan(editingPlatform.value.id, platformForm.value)
      useNotification().success({
        content: '平台更新成功',
        duration: 3000
      })
    } else {
      await panApi.createPan(platformForm.value)
      useNotification().success({
        content: '平台创建成功',
        duration: 3000
      })
    }
    
    closeModal()
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '保存平台失败',
      duration: 5000
    })
  } finally {
    submitting.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script> 