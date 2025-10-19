<template>
  <n-modal v-model:show="showModal" preset="card" :title="title" style="width: 90vw; max-width: 1200px; max-height: 80vh;">
    <div class="space-y-4">
      <!-- 搜索 -->
      <div class="flex gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索文件名..."
          @keyup.enter="handleSearch"
          class="flex-1"
          clearable
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>

        <n-button type="primary" @click="handleSearch" class="w-20">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
      </div>

      <!-- 文件列表 -->
      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="fileList.length === 0" class="text-center py-8">
        <i class="fas fa-file-upload text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无图片文件</p>
      </div>

      <div v-else class="file-grid">
        <div
          v-for="file in fileList"
          :key="file.id"
          class="file-item cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg p-3 transition-colors"
          :class="{ 'bg-blue-50 dark:bg-blue-900/20 border-2 border-blue-300 dark:border-blue-600': selectedFileId === file.id }"
          @click="selectFile(file)"
        >
          <div class="image-preview">
            <n-image
              :src="getImageUrl(file.access_url)"
              :alt="file.original_name"
              :lazy="false"
              object-fit="cover"
              class="preview-image rounded"
              @error="handleImageError"
              @load="handleImageLoad"
            />

            <div class="image-info mt-2">
              <div class="file-name text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                {{ file.original_name }}
              </div>
              <div class="file-size text-xs text-gray-500 dark:text-gray-400">
                {{ formatFileSize(file.file_size) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <n-pagination
          v-model:page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-count="Math.ceil(pagination.total / pagination.pageSize)"
          :page-sizes="pagination.pageSizes"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>

    <template #footer>
      <n-space justify="end">
        <n-button @click="closeModal">取消</n-button>
        <n-button
          type="primary"
          @click="confirmSelection"
          :disabled="!selectedFileId"
        >
          确认选择
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
// Props
const props = defineProps<{
  show: boolean
  title: string
}>()

// Emits
const emit = defineEmits<{
  'update:show': [value: boolean]
  'select': [file: any]
}>()

// 使用图片URL composable
const { getImageUrl } = useImageUrl()

// 响应式数据
const showModal = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

const loading = ref(false)
const fileList = ref<any[]>([])
const selectedFileId = ref<number | null>(null)
const searchKeyword = ref('')

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  pageSizes: [10, 20, 50, 100]
})

// 监听show变化，重新加载数据
watch(() => props.show, (newValue) => {
  if (newValue) {
    loadFileList()
  } else {
    // 重置状态
    selectedFileId.value = null
    searchKeyword.value = ''
  }
})

// 加载文件列表
const loadFileList = async () => {
  try {
    loading.value = true
    const { useFileApi } = await import('~/composables/useFileApi')
    const fileApi = useFileApi()

    const response = await fileApi.getFileList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      search: searchKeyword.value,
      fileType: 'image',
      status: 'active'
    }) as any

    if (response && response.data) {
      fileList.value = response.data.files || []
      pagination.value.total = response.data.total || 0
    }
  } catch (error) {
    console.error('获取文件列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.value.page = 1
  loadFileList()
}

// 分页处理
const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadFileList()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  loadFileList()
}

// 文件选择
const selectFile = (file: any) => {
  selectedFileId.value = file.id
}

// 确认选择
const confirmSelection = () => {
  if (selectedFileId.value) {
    const file = fileList.value.find(f => f.id === selectedFileId.value)
    if (file) {
      emit('select', file)
      closeModal()
    }
  }
}

// 关闭模态框
const closeModal = () => {
  showModal.value = false
}

// 文件大小格式化
const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(1) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

// 图片加载处理
const handleImageError = (event: any) => {
  console.error('图片加载失败:', event)
}

const handleImageLoad = (event: any) => {
  console.log('图片加载成功:', event)
}
</script>

<style scoped>
.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  max-height: 400px;
  overflow-y: auto;
}

.file-item {
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.file-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.preview-image {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}
</style>