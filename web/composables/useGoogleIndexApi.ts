import { useApiFetch } from './useApiFetch'
import { parseApiResponse } from './useApi'

// Google索引配置类型定义
export interface GoogleIndexConfig {
  id?: number
  group: string
  key: string
  value: string
  type?: string
}

// Google索引任务类型定义
export interface GoogleIndexTask {
  id: number
  title: string
  type: string
  status: string
  description: string
  totalItems: number
  processedItems: number
  successItems: number
  failedItems: number
  indexedURLs: number
  failedURLs: number
  errorMessage?: string
  configID?: number
  startedAt?: Date
  completedAt?: Date
  createdAt: Date
  updatedAt: Date
}

// Google索引任务项类型定义
export interface GoogleIndexTaskItem {
  id: number
  taskID: number
  URL: string
  status: string
  indexStatus: string
  errorMessage?: string
  inspectResult?: string
  mobileFriendly: boolean
  lastCrawled?: Date
  statusCode: number
  startedAt?: Date
  completedAt?: Date
  createdAt: Date
  updatedAt: Date
}

// URL状态类型定义
export interface GoogleIndexURLStatus {
  id: number
  URL: string
  indexStatus: string
  lastChecked: Date
  canonicalURL?: string
  lastCrawled?: Date
  changeFreq?: string
  priority?: number
  mobileFriendly: boolean
  robotsBlocked: boolean
  lastError?: string
  statusCode: number
  statusCodeText: string
  checkCount: number
  successCount: number
  failureCount: number
  createdAt: Date
  updatedAt: Date
}

// Google索引状态响应类型定义
export interface GoogleIndexStatusResponse {
  enabled: boolean
  siteURL: string
  lastCheckTime: Date
  totalURLs: number
  indexedURLs: number
  notIndexedURLs: number
  errorURLs: number
  lastSitemapSubmit: Date
  authValid: boolean
}

// Google索引任务列表响应类型定义
export interface GoogleIndexTaskListResponse {
  tasks: GoogleIndexTask[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// Google索引任务项分页响应类型定义
export interface GoogleIndexTaskItemPageResponse {
  items: GoogleIndexTaskItem[]
  total: number
  page: number
  size: number
}

// Google索引API封装
export const useGoogleIndexApi = () => {
  // 配置管理API
  const getGoogleIndexConfig = (params?: any) =>
    useApiFetch('/google-index/config', { params }).then(parseApiResponse<GoogleIndexConfig[]>)

  const getGoogleIndexConfigByKey = (key: string) =>
    useApiFetch(`/google-index/config/${key}`).then(parseApiResponse<GoogleIndexConfig>)

  const updateGoogleIndexConfig = (data: GoogleIndexConfig) =>
    useApiFetch('/google-index/config', { method: 'POST', body: data }).then(parseApiResponse<GoogleIndexConfig>)

  const deleteGoogleIndexConfig = (key: string) =>
    useApiFetch(`/google-index/config/${key}`, { method: 'DELETE' }).then(parseApiResponse<boolean>)

  // 任务管理API
  const getGoogleIndexTasks = (params?: any) =>
    useApiFetch('/google-index/tasks', { params }).then(parseApiResponse<GoogleIndexTaskListResponse>)

  const getGoogleIndexTask = (id: number) =>
    useApiFetch(`/google-index/tasks/${id}`).then(parseApiResponse<GoogleIndexTask>)

  const createGoogleIndexTask = (data: any) =>
    useApiFetch('/google-index/tasks', { method: 'POST', body: data }).then(parseApiResponse<GoogleIndexTask>)

  const startGoogleIndexTask = (id: number) =>
    useApiFetch(`/google-index/tasks/${id}/start`, { method: 'POST' }).then(parseApiResponse<boolean>)

  const stopGoogleIndexTask = (id: number) =>
    useApiFetch(`/google-index/tasks/${id}/stop`, { method: 'POST' }).then(parseApiResponse<boolean>)

  const deleteGoogleIndexTask = (id: number) =>
    useApiFetch(`/google-index/tasks/${id}`, { method: 'DELETE' }).then(parseApiResponse<boolean>)

  // 任务项管理API
  const getGoogleIndexTaskItems = (taskId: number, params?: any) =>
    useApiFetch(`/google-index/tasks/${taskId}/items`, { params }).then(parseApiResponse<GoogleIndexTaskItemPageResponse>)

  // URL状态管理API
  const getGoogleIndexURLStatus = (params?: any) =>
    useApiFetch('/google-index/urls/status', { params }).then(parseApiResponse<GoogleIndexURLStatus[]>)

  const getGoogleIndexURLStatusByURL = (url: string) =>
    useApiFetch(`/google-index/urls/status/${encodeURIComponent(url)}`).then(parseApiResponse<GoogleIndexURLStatus>)

  const checkGoogleIndexURLStatus = (data: { urls: string[] }) =>
    useApiFetch('/google-index/urls/check', { method: 'POST', body: data }).then(parseApiResponse<any>)

  const submitGoogleIndexURL = (data: { urls: string[] }) =>
    useApiFetch('/google-index/urls/submit', { method: 'POST', body: data }).then(parseApiResponse<any>)

  // 批量操作API
  const batchSubmitGoogleIndexURLs = (data: { urls: string[], operation: string }) =>
    useApiFetch('/google-index/batch/submit', { method: 'POST', body: data }).then(parseApiResponse<any>)

  const batchCheckGoogleIndexURLs = (data: { urls: string[], operation: string }) =>
    useApiFetch('/google-index/batch/check', { method: 'POST', body: data }).then(parseApiResponse<any>)

  // 网站地图提交API
  const submitGoogleIndexSitemap = (data: { sitemapURL: string }) =>
    useApiFetch('/google-index/sitemap/submit', { method: 'POST', body: data }).then(parseApiResponse<any>)

  // 状态查询API
  const getGoogleIndexStatus = () =>
    useApiFetch('/google-index/status').then(parseApiResponse<GoogleIndexStatusResponse>)

  // 验证凭据API
  const validateCredentials = (data: { credentialsFile: string }) =>
    useApiFetch('/google-index/validate-credentials', { method: 'POST', body: data }).then(parseApiResponse<any>)

  // 更新Google索引分组配置API
  const updateGoogleIndexGroupConfig = (data: GoogleIndexConfig) =>
    useApiFetch('/google-index/config/update', { method: 'POST', body: data }).then(parseApiResponse<GoogleIndexConfig>)

  // 上传凭据API
  const uploadCredentials = (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return useApiFetch('/google-index/upload-credentials', {
      method: 'POST',
      body: formData,
      headers: {
        // 注意：此处不应包含Authorization头，因为文件上传通常由use-upload组件处理
      }
    }).then(parseApiResponse<any>)
  }

  // 调度器控制API
  const startGoogleIndexScheduler = () =>
    useApiFetch('/google-index/scheduler/start', { method: 'POST' }).then(parseApiResponse<boolean>)

  const stopGoogleIndexScheduler = () =>
    useApiFetch('/google-index/scheduler/stop', { method: 'POST' }).then(parseApiResponse<boolean>)

  const getGoogleIndexSchedulerStatus = () =>
    useApiFetch('/google-index/scheduler/status').then(parseApiResponse<any>)

  return {
    // 配置管理
    getGoogleIndexConfig,
    getGoogleIndexConfigByKey,
    updateGoogleIndexConfig,
    updateGoogleIndexGroupConfig,
    deleteGoogleIndexConfig,

    // 凭据验证和上传
    validateCredentials,
    uploadCredentials,

    // 任务管理
    getGoogleIndexTasks,
    getGoogleIndexTask,
    createGoogleIndexTask,
    startGoogleIndexTask,
    stopGoogleIndexTask,
    deleteGoogleIndexTask,

    // 任务项管理
    getGoogleIndexTaskItems,

    // URL状态管理
    getGoogleIndexURLStatus,
    getGoogleIndexURLStatusByURL,
    checkGoogleIndexURLStatus,
    submitGoogleIndexURL,

    // 批量操作
    batchSubmitGoogleIndexURLs,
    batchCheckGoogleIndexURLs,

    // 网站地图提交
    submitGoogleIndexSitemap,

    // 状态查询
    getGoogleIndexStatus,

    // 调度器控制
    startGoogleIndexScheduler,
    stopGoogleIndexScheduler,
    getGoogleIndexSchedulerStatus
  }
}