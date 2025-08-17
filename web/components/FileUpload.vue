<template>
  <div class="file-upload-container">
    <n-upload
    multiple
    directory-dnd
    :custom-request="customRequest"
    :on-before-upload="handleBeforeUpload"
    :on-finish="handleUploadFinish"
    :on-error="handleUploadError"
    :on-remove="handleFileRemove"
    :max="5"
    ref="uploadRef"
  >
    <n-upload-dragger>
      <div style="margin-bottom: 12px">
        <i class="fas fa-cloud-upload-alt text-4xl text-gray-400"></i>
      </div>
      <n-text style="font-size: 16px">
        点击或者拖动文件到该区域来上传
      </n-text>
      <n-p depth="3" style="margin: 8px 0 0 0">
        请不要上传敏感数据，比如你的银行卡号和密码，信用卡号有效期和安全码
      </n-p>
    </n-upload-dragger>
  </n-upload>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { useFileApi } from '~/composables/useFileApi'

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

interface UploadFileInfo {
  id: string
  name: string
  status: 'pending' | 'uploading' | 'finished' | 'error' | 'removed'
  url?: string
  file?: File
}

const message = useMessage()
const fileApi = useFileApi()

// 响应式数据
const uploadRef = ref()
const fileList = ref<FileItem[]>([])
const isPublic = ref(true) // 默认公开
const maxFiles = ref(10)
const maxFileSize = ref(5 * 1024 * 1024) // 5MB
const acceptTypes = ref('image/*')

// 添加状态标记：用于跟踪已上传的文件
const uploadedFiles = ref<Map<string, boolean>>(new Map()) // 文件哈希 -> 是否已上传
const uploadingFiles = ref<Set<string>>(new Set()) // 正在上传的文件哈希

// 生成文件哈希值（基于文件名、大小和修改时间）
const generateFileHash = (file: File): string => {
  return `${file.name}_${file.size}_${file.lastModified}`
}

// 检查文件是否已经上传过
const isFileAlreadyUploaded = (file: File): boolean => {
  const fileHash = generateFileHash(file)
  return uploadedFiles.value.has(fileHash)
}

// 检查文件是否正在上传
const isFileUploading = (file: File): boolean => {
  const fileHash = generateFileHash(file)
  return uploadingFiles.value.has(fileHash)
}

// 标记文件为已上传
const markFileAsUploaded = (file: File) => {
  const fileHash = generateFileHash(file)
  uploadedFiles.value.set(fileHash, true)
  uploadingFiles.value.delete(fileHash)
}

// 标记文件为正在上传
const markFileAsUploading = (file: File) => {
  const fileHash = generateFileHash(file)
  uploadingFiles.value.add(fileHash)
}

// 标记文件上传失败
const markFileAsFailed = (file: File) => {
  const fileHash = generateFileHash(file)
  uploadingFiles.value.delete(fileHash)
}

// 自定义上传请求
const customRequest = async (options: any) => {
  const { file, onProgress, onSuccess, onError } = options
  
  // 检查文件是否已经上传过
  if (isFileAlreadyUploaded(file.file)) {
    message.warning(`${file.name} 已经上传过了，跳过重复上传`)
    if (onSuccess) {
      onSuccess({ message: '文件已存在，跳过上传' })
    }
    return
  }
  
  // 检查文件是否正在上传
  if (isFileUploading(file.file)) {
    message.warning(`${file.name} 正在上传中，请稍候`)
    return
  }
  
  // 标记文件为正在上传
  markFileAsUploading(file.file)
  
  console.log('开始上传文件:', file.name, file.file)
  
  try {
    // 创建FormData
    const formData = new FormData()
    formData.append('file', file.file)
    formData.append('is_public', isPublic.value.toString())
    
    // 调用统一的API接口
    const response = await fileApi.uploadFile(formData)
    
    console.log('文件上传成功:', file.name, response)
    
    // 标记文件为已上传
    markFileAsUploaded(file.file)
    
    if (onSuccess) {
      onSuccess(response)
    }
  } catch (error) {
    console.error('文件上传失败:', file.name, error)
    // 标记文件上传失败
    markFileAsFailed(file.file)
    if (onError) {
      onError(error)
    }
  }
}


// 默认文件列表（从props传入）
const defaultFileList = ref<UploadFileInfo[]>([])

// 方法
const handleBeforeUpload = (data: { file: Required<UploadFileInfo> }) => {
  const { file } = data
  
  // 检查文件是否已经上传过
  if (file.file && isFileAlreadyUploaded(file.file)) {
    //message.warning(`${file.name} 已经上传过了，请勿重复上传`)
    return false
  }
  
  // 检查文件是否正在上传
  if (file.file && isFileUploading(file.file)) {
    message.warning(`${file.name} 正在上传中，请稍候`)
    return false
  }
  
  // 检查文件大小
  if (file.file && file.file.size > maxFileSize.value) {
    message.error(`文件大小不能超过 ${formatFileSize(maxFileSize.value)}`)
    return false
  }
  
  // 检查文件类型
  if (file.file) {
    const fileName = file.file.name.toLowerCase()
    const acceptedTypes = acceptTypes.value.split(',')
    const isAccepted = acceptedTypes.some(type => {
      if (type === 'image/*') {
        return file.file!.type.startsWith('image/')
      }
      if (type.startsWith('.')) {
        return fileName.endsWith(type)
      }
      return file.file!.type === type
    })
    
    if (!isAccepted) {
      message.error('只支持图片格式文件')
      return false
    }
  }
  
  return true
}

const handleUploadFinish = (data: { file: Required<UploadFileInfo> }) => {
  const { file } = data
  
  if (file.status === 'finished') {
    message.success(`${file.name} 上传成功`)
    // 确保文件被标记为已上传
    if (file.file) {
      markFileAsUploaded(file.file)
    }
  }
}

const handleUploadError = (data: { file: Required<UploadFileInfo> }) => {
  const { file } = data
  message.error(`${file.name} 上传失败`)
  // 标记文件上传失败
  if (file.file) {
    markFileAsFailed(file.file)
  }
}

const handleFileRemove = (data: { file: Required<UploadFileInfo> }) => {
  const { file } = data
  message.info(`已移除 ${file.name}`)
  // 从上传状态中移除文件
  if (file.file) {
    const fileHash = generateFileHash(file.file)
    uploadingFiles.value.delete(fileHash)
  }
}

const loadFileList = async () => {
  try {
    const response = await fileApi.getFileList({
      page: 1,
      page_size: 50
    })
    fileList.value = response.data.files || []
  } catch (error) {
    console.error('加载文件列表失败:', error)
    message.error('加载文件列表失败')
  }
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

const getFileIconClass = (fileType: string) => {
  const iconMap: Record<string, string> = {
    'image': 'fas fa-image text-blue-500',
    'document': 'fas fa-file-alt text-green-500',
    'video': 'fas fa-video text-red-500',
    'audio': 'fas fa-music text-purple-500',
    'archive': 'fas fa-archive text-orange-500',
    'other': 'fas fa-file text-gray-500'
  }
  return iconMap[fileType] || iconMap.other
}

const getFileDescription = (file: FileItem) => {
  return `${formatFileSize(file.file_size)} | ${file.file_type} | ${file.created_at}`
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 生命周期
// 不在组件挂载时加载文件列表，由父组件管理

// 重置上传组件状态
const resetUpload = () => {
  if (uploadRef.value) {
    uploadRef.value.clear()
  }
  // 清空上传状态
  uploadedFiles.value.clear()
  uploadingFiles.value.clear()
}

// 清空已上传文件状态（用于重新开始上传）
const clearUploadedFiles = () => {
  uploadedFiles.value.clear()
  uploadingFiles.value.clear()
}

// 暴露方法给父组件
defineExpose({
  loadFileList,
  fileList,
  resetUpload,
  clearUploadedFiles,
  uploadedFiles,
  uploadingFiles
})
</script>

<style scoped>
</style> 