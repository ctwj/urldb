<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">文件管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理系统中的上传文件</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="openUploadModal">
          <template #icon>
            <i class="fas fa-upload"></i>
          </template>
          上传文件
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </div>

    <!-- 提示信息 -->
    <n-alert title="支持图片格式文件，最大文件大小5MB" type="info" />

    <!-- 搜索和筛选 -->
    <n-card>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索文件名..."
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-select
          v-model:value="filterType"
          placeholder="文件类型"
          :options="fileTypeOptions"
          clearable
        />
        
        <n-select
          v-model:value="filterStatus"
          placeholder="状态"
          :options="statusOptions"
          clearable
        />
        
        <n-button type="primary" @click="handleSearch">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
      </div>
    </n-card>

    <!-- 文件列表 -->
    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">文件列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 个文件</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="fileList.length === 0" class="text-center py-8">
        <i class="fas fa-file-upload text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无文件数据</p>
        <n-button @click="openUploadModal" type="primary" class="mt-4">
          <template #icon>
            <i class="fas fa-upload"></i>
          </template>
          上传文件
        </n-button>
      </div>

      <div v-else>
        <!-- 图片预览区域 -->
        <div class="image-preview-container">
          <n-image-group>
            <div class="image-grid">
              <div 
                v-for="file in fileList" 
                :key="file.id" 
                class="image-item"
                :class="{ 'is-image': isImageFile(file) }"
              >
                                 <!-- 图片文件显示预览 -->
                 <div v-if="isImageFile(file)" class="file-item cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg p-3 transition-colors border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600">
                   <div class="image-preview relative">
                     <n-image
                       :src="getImageUrl(file.access_url)"
                       :alt="file.original_name"
                       :lazy="false"
                       object-fit="cover"
                       class="preview-image rounded"
                       @error="handleImageError"
                       @load="handleImageLoad"
                     />
                     <div class="delete-button">
                       <n-button 
                         size="small" 
                         type="error" 
                         circle
                         @click="confirmDelete(file)"
                       >
                         <template #icon>
                           <i class="fas fa-trash"></i>
                         </template>
                       </n-button>
                     </div>
                   </div>
                   <div class="image-info mt-2">
                     <div class="file-name text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                       {{ file.original_name }}
                     </div>
                     <div class="file-size text-xs text-gray-500 dark:text-gray-400">
                       {{ formatFileSize(file.file_size) }}
                     </div>
                   </div>
                 </div>
                
                                 <!-- 非图片文件显示图标 -->
                 <div v-else class="file-item cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg p-3 transition-colors border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 relative">
                   <div class="file-icon">
                     <i :class="getFileIconClass(file.file_type)"></i>
                   </div>
                   <div class="file-info">
                     <div class="file-name text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                       {{ file.original_name }}
                     </div>
                     <div class="file-size text-xs text-gray-500 dark:text-gray-400">
                       {{ formatFileSize(file.file_size) }}
                     </div>
                   </div>
                   <div class="delete-button">
                     <n-button 
                       size="small" 
                       type="error" 
                       circle
                       @click="confirmDelete(file)"
                     >
                       <template #icon>
                         <i class="fas fa-trash"></i>
                       </template>
                     </n-button>
                   </div>
                 </div>
              </div>
            </div>
          </n-image-group>
          
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
      </div>
    </n-card>

    <!-- 上传模态框 -->
    <n-modal v-model:show="showUploadModal" preset="card" title="上传文件" style="width: 800px" @update:show="handleModalClose">
      <FileUpload ref="fileUploadRef" :key="uploadModalKey" />
      <template #footer>
        <n-space justify="end">
          <n-button @click="showUploadModal = false">取消</n-button>
          <n-button type="primary" @click="handleUploadSuccess">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 删除确认对话框 -->
    <n-modal v-model:show="showDeleteModal" preset="card" title="确认删除" style="width: 400px">
      <div class="text-center py-4">
        <i class="fas fa-exclamation-triangle text-yellow-500 text-4xl mb-4"></i>
        <p class="text-lg font-medium mb-2">确定要删除这个文件吗？</p>
        <p class="text-gray-600 mb-4">{{ fileToDelete?.original_name }}</p>
        <p class="text-sm text-gray-500">此操作不可撤销，文件将被永久删除。</p>
      </div>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showDeleteModal = false">取消</n-button>
          <n-button type="error" @click="handleConfirmDelete">确认删除</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useMessage } from 'naive-ui'
import { useFileApi } from '~/composables/useFileApi'
import { useImageUrl } from '~/composables/useImageUrl'


// 设置页面布局
definePageMeta({
  layout: 'admin'
})

interface FileItem {
  id: number
  original_name: string
  file_name: string
  file_path: string
  file_size: number
  file_type: string
  mime_type: string
  access_url: string
  user_id: number
  user: string
  status: string
  is_public: boolean
  is_deleted: boolean
  created_at: string
  updated_at: string
}

const message = useMessage()
const fileApi = useFileApi()
const { getImageUrl } = useImageUrl()

// 响应式数据
const loading = ref(false)
const fileList = ref<FileItem[]>([])
const searchKeyword = ref('')
const filterType = ref('')
const filterStatus = ref('')
const showUploadModal = ref(false)
const fileUploadRef = ref()
const uploadModalKey = ref(0)

// 删除确认相关
const showDeleteModal = ref(false)
const fileToDelete = ref<FileItem | null>(null)

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

// 总数
const total = computed(() => pagination.value.total)

// 选项
const fileTypeOptions = [
  { label: '全部', value: '' },
  { label: '图片', value: 'image' }
]

const statusOptions = [
  { label: '全部', value: '' },
  { label: '正常', value: 'active' },
  { label: '禁用', value: 'inactive' },
  { label: '已删除', value: 'deleted' }
]



// 方法
const loadFileList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize,
      search: searchKeyword.value,
      file_type: filterType.value,
      status: filterStatus.value
    }
    
    const response = await fileApi.getFileList(params)
    fileList.value = response.data.files || []
    pagination.value.total = response.data.total || 0
    
    console.log('文件列表加载完成:', {
      total: pagination.value.total,
      files: fileList.value.map(f => ({
        id: f.id,
        name: f.original_name,
        type: f.file_type,
        url: f.access_url,
        isImage: isImageFile(f)
      }))
    })
  } catch (error) {
    console.error('加载文件列表失败:', error)
    message.error('加载文件列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadFileList()
}

const resetFilters = () => {
  searchKeyword.value = ''
  filterType.value = ''
  filterStatus.value = ''
  pagination.value.page = 1
  loadFileList()
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadFileList()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  loadFileList()
}

const copyFileUrl = async (file: FileItem) => {
  try {
    await navigator.clipboard.writeText(file.access_url)
    message.success('文件链接已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    message.error('复制失败')
  }
}

const openFile = (file: FileItem) => {
  window.open(file.access_url, '_blank')
}

const toggleFilePublic = async (file: FileItem) => {
  try {
    await fileApi.updateFile({
      id: file.id,
      is_public: !file.is_public
    })
    message.success('文件状态更新成功')
    loadFileList()
  } catch (error) {
    console.error('更新文件状态失败:', error)
    message.error('更新文件状态失败')
  }
}

const confirmDelete = (file: FileItem) => {
  fileToDelete.value = file
  showDeleteModal.value = true
}

const handleConfirmDelete = async () => {
  if (!fileToDelete.value) return
  
  try {
    await fileApi.deleteFiles([fileToDelete.value.id])
    message.success('文件删除成功')
    showDeleteModal.value = false
    fileToDelete.value = null
    loadFileList()
  } catch (error) {
    console.error('删除文件失败:', error)
    message.error('删除文件失败')
  }
}

const deleteFile = async (file: FileItem) => {
  try {
    await fileApi.deleteFiles([file.id])
    message.success('文件删除成功')
    loadFileList()
  } catch (error) {
    console.error('删除文件失败:', error)
    message.error('删除文件失败')
  }
}

const refreshData = () => {
  loadFileList()
}

const handleUploadSuccess = () => {
  // 重置上传组件状态
  if (fileUploadRef.value && fileUploadRef.value.resetUpload) {
    fileUploadRef.value.resetUpload()
  }
  showUploadModal.value = false
  loadFileList()
  message.success('文件上传成功')
}

const openUploadModal = () => {
  uploadModalKey.value++ // 强制重新渲染组件
  showUploadModal.value = true
}

const handleModalClose = (show: boolean) => {
  if (!show) {
    // 模态框关闭时重置上传组件状态
    if (fileUploadRef.value && fileUploadRef.value.resetUpload) {
      fileUploadRef.value.resetUpload()
    }
  }
}

const getFileIconClass = (fileType: string) => {
  const iconMap: Record<string, string> = {
    'image': 'fas fa-image text-blue-500',
    'jpeg': 'fas fa-image text-blue-500',
    'jpg': 'fas fa-image text-blue-500',
    'png': 'fas fa-image text-green-500',
    'gif': 'fas fa-image text-purple-500',
    'webp': 'fas fa-image text-orange-500',
    'bmp': 'fas fa-image text-red-500',
    'svg': 'fas fa-image text-indigo-500'
  }
  return iconMap[fileType] || 'fas fa-image text-gray-500'
}

const getFileTypeLabel = (fileType: string) => {
  const labelMap: Record<string, string> = {
    'jpeg': 'JPEG',
    'jpg': 'JPEG',
    'png': 'PNG',
    'gif': 'GIF',
    'webp': 'WebP',
    'bmp': 'BMP',
    'svg': 'SVG'
  }
  return labelMap[fileType] || '图片'
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const isImageFile = (file: FileItem) => {
  // 后端返回的 file_type 是 "image"，所以直接检查这个值
  const isImageByType = file.file_type.toLowerCase() === 'image'
  
  // 检查文件名扩展名
  const imageExtensions = ['jpeg', 'jpg', 'png', 'gif', 'webp', 'bmp', 'svg']
  const fileNameLower = file.original_name.toLowerCase()
  const hasImageExtension = imageExtensions.some(ext => fileNameLower.endsWith(`.${ext}`))
  
  // 检查 MIME 类型
  const mimeTypeLower = (file.mime_type || '').toLowerCase()
  const isImageByMime = mimeTypeLower.startsWith('image/')
  
  // 综合判断
  const isImage = isImageByType || hasImageExtension || isImageByMime
  
  console.log('isImageFile 详细检查:', { 
    fileName: file.original_name, 
    fileType: file.file_type,
    mimeType: file.mime_type,
    isImageByType: isImageByType,
    hasImageExtension: hasImageExtension,
    isImageByMime: isImageByMime,
    finalResult: isImage,
    accessUrl: file.access_url,
    processedUrl: getImageUrl(file.access_url)
  })
  
  return isImage
}

const handleImageError = (event: any) => {
  console.error('图片加载失败:', event)
}

const handleImageLoad = (event: any) => {
  console.log('图片加载成功:', event)
}

// 生命周期
onMounted(() => {
  loadFileList()
})
</script>

<style scoped>
/* 文件管理页面样式 */


.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  max-height: 400px;
  overflow-y: auto;
}





.preview-image {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border: 1px solid #f3f4f6;
  border-radius: 4px;
}

.delete-button {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 10;
}

.image-preview:hover .delete-button,
.file-item:hover .delete-button {
  opacity: 1;
}

.delete-button .n-button {
  background: rgba(239, 68, 68, 0.9);
  backdrop-filter: blur(4px);
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  color: white;
  transition: all 0.3s ease;
}

.delete-button .n-button:hover {
  background: rgba(239, 68, 68, 1);
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

.delete-button .n-button i {
  font-size: 14px;
}





.file-name {
  font-weight: 500;
  font-size: 13px;
  color: #333;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-size {
  font-size: 11px;
  color: #666;
}



.file-icon {
  font-size: 48px;
  margin-bottom: 12px;
  color: #666;
}

.file-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

/* 滚动条样式 */
.image-preview-container::-webkit-scrollbar {
  width: 6px;
}

.image-preview-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.image-preview-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.image-preview-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style> 