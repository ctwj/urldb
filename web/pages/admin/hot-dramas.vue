<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">热播剧管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的热播剧信息</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加热播剧
        </n-button>
        <n-button type="info" @click="fetchFromDouban">
          <template #icon>
            <i class="fas fa-download"></i>
          </template>
          从豆瓣获取
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
        placeholder="搜索热播剧..."
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <i class="fas fa-search"></i>
        </template>
      </n-input>
    </n-card>

    <!-- 热播剧列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">热播剧列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个热播剧</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="hotDramas.length === 0" class="text-center py-8">
        <i class="fas fa-film text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无热播剧数据</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="drama in hotDramas"
          :key="drama.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <div class="w-12 h-16 bg-gray-200 dark:bg-gray-700 rounded overflow-hidden">
                  <img 
                    v-if="drama.cover" 
                    :src="drama.cover" 
                    :alt="drama.title"
                    class="w-full h-full object-cover"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center">
                    <i class="fas fa-image text-gray-400"></i>
                  </div>
                </div>
                <div class="flex-1">
                  <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                    {{ drama.title }}
                  </h3>
                  <p v-if="drama.original_title" class="text-sm text-gray-600 dark:text-gray-400">
                    {{ drama.original_title }}
                  </p>
                </div>
                <span v-if="drama.rating" class="text-xs px-2 py-1 bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200 rounded">
                  {{ drama.rating }}分
                </span>
              </div>
              
              <div class="flex items-center space-x-4 text-sm text-gray-600 dark:text-gray-400 mb-2">
                <span v-if="drama.year">年份: {{ drama.year }}</span>
                <span v-if="drama.country">国家: {{ drama.country }}</span>
                <span v-if="drama.genre">类型: {{ drama.genre }}</span>
                <span v-if="drama.episodes">集数: {{ drama.episodes }}</span>
              </div>
              
              <p v-if="drama.summary" class="text-sm text-gray-600 dark:text-gray-400 mb-2 line-clamp-2">
                {{ drama.summary }}
              </p>
              
              <div class="flex items-center space-x-4 text-xs text-gray-500">
                <span>豆瓣ID: {{ drama.douban_id || '无' }}</span>
                <span>创建时间: {{ formatDate(drama.created_at) }}</span>
                <span>更新时间: {{ formatDate(drama.updated_at) }}</span>
              </div>
            </div>
            
            <div class="flex space-x-2">
              <n-button size="small" type="primary" @click="editDrama(drama)">
                <template #icon>
                  <i class="fas fa-edit"></i>
                </template>
                编辑
              </n-button>
              <n-button size="small" type="error" @click="deleteDrama(drama)">
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

    <!-- 添加/编辑热播剧模态框 -->
    <n-modal v-model:show="showAddModal" preset="card" title="添加热播剧" style="width: 700px">
      <n-form
        ref="formRef"
        :model="dramaForm"
        :rules="rules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <n-form-item label="剧名" path="title">
            <n-input
              v-model:value="dramaForm.title"
              placeholder="请输入剧名"
            />
          </n-form-item>

          <n-form-item label="原名" path="original_title">
            <n-input
              v-model:value="dramaForm.original_title"
              placeholder="请输入原名"
            />
          </n-form-item>

          <n-form-item label="年份" path="year">
            <n-input
              v-model:value="dramaForm.year"
              placeholder="请输入年份"
              type="text"
            />
          </n-form-item>

          <n-form-item label="评分" path="rating">
            <n-input
              v-model:value="dramaForm.rating"
              placeholder="请输入评分"
              type="text"
            />
          </n-form-item>

          <n-form-item label="国家" path="country">
            <n-input
              v-model:value="dramaForm.country"
              placeholder="请输入国家"
            />
          </n-form-item>

          <n-form-item label="类型" path="genre">
            <n-input
              v-model:value="dramaForm.genre"
              placeholder="请输入类型"
            />
          </n-form-item>

          <n-form-item label="集数" path="episodes">
            <n-input
              v-model:value="dramaForm.episodes"
              placeholder="请输入集数"
              type="text"
            />
          </n-form-item>

          <n-form-item label="豆瓣ID" path="douban_id">
            <n-input
              v-model:value="dramaForm.douban_id"
              placeholder="请输入豆瓣ID"
              type="text"
            />
          </n-form-item>
        </div>

        <n-form-item label="封面" path="cover">
          <n-input
            v-model:value="dramaForm.cover"
            placeholder="请输入封面图片URL"
          />
        </n-form-item>

        <n-form-item label="简介" path="summary">
          <n-input
            v-model:value="dramaForm.summary"
            placeholder="请输入简介"
            type="textarea"
            :rows="4"
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
const { useHotDramaApi } = await import('~/composables/useApi')
const hotDramaApi = useHotDramaApi()

// 响应式数据
const loading = ref(false)
const hotDramas = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')

// 模态框相关
const showAddModal = ref(false)
const submitting = ref(false)
const editingDrama = ref<any>(null)

// 表单数据
const dramaForm = ref({
  title: '',
  original_title: '',
  year: '',
  rating: '',
  country: '',
  genre: '',
  episodes: '',
  douban_id: '',
  cover: '',
  summary: ''
})

// 表单验证规则
const rules = {
  title: {
    required: true,
    message: '请输入剧名',
    trigger: 'blur'
  }
}

const formRef = ref()

// 获取数据
const fetchData = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value
    }
    
    const response = await hotDramaApi.getHotDramas(params) as any
    hotDramas.value = response.data || []
    total.value = response.total || 0
  } catch (error: any) {
    useNotification().error({
      content: error.message || '获取热播剧数据失败',
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

// 从豆瓣获取数据
const fetchFromDouban = async () => {
  try {
    loading.value = true
    await hotDramaApi.fetchHotDramas()
    useNotification().success({
      content: '从豆瓣获取数据成功',
      duration: 3000
    })
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '从豆瓣获取数据失败',
      duration: 5000
    })
  } finally {
    loading.value = false
  }
}

// 编辑热播剧
const editDrama = (drama: any) => {
  editingDrama.value = drama
  dramaForm.value = {
    title: drama.title,
    original_title: drama.original_title || '',
    year: drama.year || '',
    rating: drama.rating || '',
    country: drama.country || '',
    genre: drama.genre || '',
    episodes: drama.episodes || '',
    douban_id: drama.douban_id || '',
    cover: drama.cover || '',
    summary: drama.summary || ''
  }
  showAddModal.value = true
}

// 删除热播剧
const deleteDrama = async (drama: any) => {
  try {
    await hotDramaApi.deleteHotDrama(drama.id)
    useNotification().success({
      content: '热播剧删除成功',
      duration: 3000
    })
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '删除热播剧失败',
      duration: 5000
    })
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  editingDrama.value = null
  dramaForm.value = {
    title: '',
    original_title: '',
    year: '',
    rating: '',
    country: '',
    genre: '',
    episodes: '',
    douban_id: '',
    cover: '',
    summary: ''
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    submitting.value = true
    
    if (editingDrama.value) {
      await hotDramaApi.updateHotDrama(editingDrama.value.id, dramaForm.value)
      useNotification().success({
        content: '热播剧更新成功',
        duration: 3000
      })
    } else {
      await hotDramaApi.createHotDrama(dramaForm.value)
      useNotification().success({
        content: '热播剧创建成功',
        duration: 3000
      })
    }
    
    closeModal()
    await fetchData()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '保存热播剧失败',
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

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 