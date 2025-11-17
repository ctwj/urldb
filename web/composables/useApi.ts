import { useApiFetch } from './useApiFetch'
import { useUserStore } from '~/stores/user'

// 统一响应解析函数
export const parseApiResponse = <T>(response: any): T => {
  // log('parseApiResponse - 原始响应:', response)
  
  // 检查是否是新的统一响应格式
  if (response && typeof response === 'object' && 'code' in response && 'data' in response) {
    if (response.code === 200) {
      // 特殊处理pan接口返回的data.list格式
      if (response.data && response.data.list && Array.isArray(response.data.list)) {
        return response.data.list
      }
      return response.data
    } else {
      throw new Error(response.message || '请求失败')
    }
  }
  
  // 检查是否是包含items字段的响应格式（如分类接口）
  if (response && typeof response === 'object' && 'items' in response) {
    return response
  }
  
  // 检查是否是包含success字段的响应格式（如登录接口）
  if (response && typeof response === 'object' && 'success' in response && 'data' in response) {
    if (response.success) {
      // 特殊处理登录接口，直接返回data部分（包含token和user）
      if (response.data && response.data.token && response.data.user) {
        // console.log('parseApiResponse - 登录接口处理，返回data:', response.data)
        return response.data
      }
      // 特殊处理删除操作响应，直接返回data部分
      if (response.data && response.data.affected_rows !== undefined) {
        return response.data
      }
      return response.data
    } else {
      throw new Error(response.message || '请求失败')
    }
  }
  
  // 兼容旧格式，直接返回响应
  return response
}

export const useResourceApi = () => {
  const getResources = (params?: any) => useApiFetch('/resources', { params }).then(parseApiResponse)
  const getResource = (id: number) => useApiFetch(`/resources/${id}`).then(parseApiResponse)
  const getResourcesByKey = (key: string) => useApiFetch(`/resources/key/${key}`).then(parseApiResponse)
  const createResource = (data: any) => useApiFetch('/resources', { method: 'POST', body: data }).then(parseApiResponse)
  const updateResource = (id: number, data: any) => useApiFetch(`/resources/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteResource = (id: number) => useApiFetch(`/resources/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const searchResources = (params: any) => useApiFetch('/search', { params }).then(parseApiResponse)
  const getResourcesByPan = (panId: number, params?: any) => useApiFetch('/resources', { params: { ...params, pan_id: panId } }).then(parseApiResponse)
  // 新增：统一的资源访问次数上报（注意：getResourceLink 已包含访问统计，通常不需要单独调用此方法）
  const incrementViewCount = (id: number) => useApiFetch(`/resources/${id}/view`, { method: 'POST' })
  // 新增：批量删除资源
  const batchDeleteResources = (ids: number[]) => useApiFetch('/resources/batch', { method: 'DELETE', body: { ids } }).then(parseApiResponse)
  // 新增：获取资源链接（智能转存）
  const getResourceLink = (id: number) => useApiFetch(`/resources/${id}/link`).then(parseApiResponse)
  return { getResources, getResource, getResourcesByKey, createResource, updateResource, deleteResource, searchResources, getResourcesByPan, incrementViewCount, batchDeleteResources, getResourceLink }
}

export const useAuthApi = () => {
  const login = (data: any) => useApiFetch('/auth/login', { method: 'POST', body: data }).then(parseApiResponse)
  const register = (data: any) => useApiFetch('/auth/register', { method: 'POST', body: data }).then(parseApiResponse)
  const getProfile = () => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : ''
    return useApiFetch('/auth/profile', { headers: token ? { Authorization: `Bearer ${token}` } : {} }).then(parseApiResponse)
  }
  return { login, register, getProfile }
}

export const useCategoryApi = () => {
  const getCategories = (params?: any) => useApiFetch('/categories', { params }).then(parseApiResponse)
  const createCategory = (data: any) => useApiFetch('/categories', { method: 'POST', body: data }).then(parseApiResponse)
  const updateCategory = (id: number, data: any) => useApiFetch(`/categories/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteCategory = (id: number) => useApiFetch(`/categories/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  return { getCategories, createCategory, updateCategory, deleteCategory }
}

export const usePanApi = () => {
  const getPans = () => useApiFetch('/pans').then(parseApiResponse)
  const getPan = (id: number) => useApiFetch(`/pans/${id}`).then(parseApiResponse)
  const createPan = (data: any) => useApiFetch('/pans', { method: 'POST', body: data }).then(parseApiResponse)
  const updatePan = (id: number, data: any) => useApiFetch(`/pans/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deletePan = (id: number) => useApiFetch(`/pans/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  return { getPans, getPan, createPan, updatePan, deletePan }
}

export const useCksApi = () => {
  const getCks = (params?: any) => useApiFetch('/cks', { params }).then(parseApiResponse)
  const getCksByID = (id: number) => useApiFetch(`/cks/${id}`).then(parseApiResponse)
  const createCks = (data: any) => useApiFetch('/cks', { method: 'POST', body: data }).then(parseApiResponse)
  const updateCks = (id: number, data: any) => useApiFetch(`/cks/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteCks = (id: number) => useApiFetch(`/cks/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const refreshCapacity = (id: number) => useApiFetch(`/cks/${id}/refresh-capacity`, { method: 'POST' }).then(parseApiResponse)
  const deleteRelatedResources = (id: number) => useApiFetch(`/cks/${id}/delete-related-resources`, { method: 'POST' }).then(parseApiResponse)
  return { getCks, getCksByID, createCks, updateCks, deleteCks, refreshCapacity, deleteRelatedResources }
}

export const useTagApi = () => {
  const getTags = (params?: any) => useApiFetch('/tags', { params }).then(parseApiResponse)
  const getTagsByCategory = (categoryId: number, params?: any) => useApiFetch(`/categories/${categoryId}/tags`, { params }).then(parseApiResponse)
  const getTag = (id: number) => useApiFetch(`/tags/${id}`).then(parseApiResponse)
  const createTag = (data: any) => useApiFetch('/tags', { method: 'POST', body: data }).then(parseApiResponse)
  const updateTag = (id: number, data: any) => useApiFetch(`/tags/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteTag = (id: number) => useApiFetch(`/tags/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const getResourceTags = (resourceId: number) => useApiFetch(`/resources/${resourceId}/tags`).then(parseApiResponse)
  return { getTags, getTagsByCategory, getTag, createTag, updateTag, deleteTag, getResourceTags }
}

export const useReadyResourceApi = () => {
  const getReadyResources = (params?: any) => useApiFetch('/ready-resources', { params }).then(parseApiResponse)
  const getFailedResources = (params?: any) => useApiFetch('/ready-resources/errors', { params }).then(parseApiResponse)
  const createReadyResource = (data: any) => useApiFetch('/ready-resources', { method: 'POST', body: data }).then(parseApiResponse)
  const batchCreateReadyResources = (data: any) => useApiFetch('/ready-resources/batch', { method: 'POST', body: data }).then(parseApiResponse)
  const createReadyResourcesFromText = (text: string) => {
    const formData = new FormData()
    formData.append('text', text)
    return useApiFetch('/ready-resources/text', { method: 'POST', body: formData }).then(parseApiResponse)
  }
  const deleteReadyResource = (id: number) => useApiFetch(`/ready-resources/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const clearReadyResources = () => useApiFetch('/ready-resources', { method: 'DELETE' }).then(parseApiResponse)
  const clearErrorMsg = (id: number) => useApiFetch(`/ready-resources/${id}/clear-error`, { method: 'POST' }).then(parseApiResponse)
  const retryFailedResources = () => useApiFetch('/ready-resources/retry-failed', { method: 'POST' }).then(parseApiResponse)
  const batchRestoreToReadyPool = (ids: number[]) => useApiFetch('/ready-resources/batch-restore', { method: 'POST', body: { ids } }).then(parseApiResponse)
  const batchRestoreToReadyPoolByQuery = (queryParams: any) => useApiFetch('/ready-resources/batch-restore-by-query', { method: 'POST', body: queryParams }).then(parseApiResponse)
  const clearAllErrorsByQuery = (queryParams: any) => useApiFetch('/ready-resources/clear-all-errors-by-query', { method: 'POST', body: queryParams }).then(parseApiResponse)
  return { 
    getReadyResources, 
    getFailedResources, 
    createReadyResource, 
    batchCreateReadyResources, 
    createReadyResourcesFromText, 
    deleteReadyResource, 
    clearReadyResources,
    clearErrorMsg,
    retryFailedResources,
    batchRestoreToReadyPool,
    batchRestoreToReadyPoolByQuery,
    clearAllErrorsByQuery
  }
}

export const useStatsApi = () => {
  const getStats = () => useApiFetch('/stats').then(parseApiResponse)
  return { getStats }
}

export const useSearchStatsApi = () => {
  const getSearchStats = (params?: any) => useApiFetch('/search-stats', { params }).then(parseApiResponse)
  const getHotKeywords = (params?: any) => useApiFetch('/search-stats/hot-keywords', { params }).then(parseApiResponse)
  const getDailyStats = (params?: any) => useApiFetch('/search-stats/daily', { params }).then(parseApiResponse)
  const getSearchTrend = (params?: any) => useApiFetch('/search-stats/trend', { params }).then(parseApiResponse)
  const getKeywordTrend = (keyword: string, params?: any) => useApiFetch(`/search-stats/keyword/${keyword}/trend`, { params }).then(parseApiResponse)
  const getSearchStatsSummary = () => useApiFetch('/search-stats/summary').then(parseApiResponse)
  const recordSearch = (data: { keyword: string }) => useApiFetch('/search-stats/record', { method: 'POST', body: data }).then(parseApiResponse)
  return { 
    getSearchStats, 
    getHotKeywords, 
    getDailyStats, 
    getSearchTrend, 
    getKeywordTrend, 
    getSearchStatsSummary,
    recordSearch
  }
}

export const useSystemConfigApi = () => {
  const getSystemConfig = () => useApiFetch('/system/config').then(parseApiResponse)
  const updateSystemConfig = (data: any) => useApiFetch('/system/config', { method: 'POST', body: data }).then(parseApiResponse)
  const getConfigStatus = () => useApiFetch('/system/config/status').then(parseApiResponse)
  const toggleAutoProcess = (enabled: boolean) => useApiFetch('/system/config/toggle-auto-process', { method: 'POST', body: { auto_process_ready_resources: enabled } }).then(parseApiResponse)
  return { getSystemConfig, updateSystemConfig, getConfigStatus, toggleAutoProcess }
}
export const useHotDramaApi = () => {
  // 为SSR优化版本，使用Nuxt3的useApiFetch
  const getHotDramas = (params?: any) => {
    return useAsyncData('hot-dramas', () => $fetch('/api/hot-dramas', {
      params,
      baseURL: process.server ? 'http://localhost:8080' : undefined,
      headers: {
        'Content-Type': 'application/json',
      }
    }).then(parseApiResponse), {
      server: true, // 启用服务端渲染
      default: () => ({
        items: [],
        total: 0
      })
    })
  }

  // 客户端交互的版本
  const getHotDramasClient = (params?: any) => useApiFetch('/hot-dramas', { params }).then(parseApiResponse)

  // 其他方法保持不变，但添加SSR支持
  const createHotDrama = (data: any) => useApiFetch('/hot-dramas', { method: 'POST', body: data }).then(parseApiResponse)
  const updateHotDrama = (id: number, data: any) => useApiFetch(`/hot-dramas/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteHotDrama = (id: number) => useApiFetch(`/hot-dramas/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const fetchHotDramas = () => useApiFetch('/hot-dramas/fetch', { method: 'POST' }).then(parseApiResponse)


  const getPosterUrl = (posterUrl: string): string => {
    if (!posterUrl) return ''
    return `/api/hot-dramas/poster?url=${encodeURIComponent(posterUrl)}`
  }

  return {
    getHotDramas,
    getHotDramasClient,
    createHotDrama,
    updateHotDrama,
    deleteHotDrama,
    fetchHotDramas,
    getPosterUrl
  }
}

export const useMonitorApi = () => {
  const getPerformanceStats = () => useApiFetch('/performance').then(parseApiResponse)
  const getSystemInfo = () => useApiFetch('/system/info').then(parseApiResponse)
  const getBasicStats = () => useApiFetch('/stats').then(parseApiResponse)
  return { getPerformanceStats, getSystemInfo, getBasicStats }
}

export const useUserApi = () => {
  const getUsers = (params?: any) => useApiFetch('/users', { params }).then(parseApiResponse)
  const getUser = (id: number) => useApiFetch(`/users/${id}`).then(parseApiResponse)
  const createUser = (data: any) => useApiFetch('/users', { method: 'POST', body: data }).then(parseApiResponse)
  const updateUser = (id: number, data: any) => useApiFetch(`/users/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteUser = (id: number) => useApiFetch(`/users/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const changePassword = (id: number, newPassword: string) => useApiFetch(`/users/${id}/password`, { method: 'PUT', body: { new_password: newPassword } }).then(parseApiResponse)
  return { getUsers, getUser, createUser, updateUser, deleteUser, changePassword }
} 

// 公开获取系统配置API
export const usePublicSystemConfigApi = () => {
  const getPublicSystemConfig = () => useApiFetch('/public/system-config').then(res => res)
  return { getPublicSystemConfig }
} 

// 任务管理API
export const useTaskApi = () => {
  const createBatchTransferTask = (data: any) => useApiFetch('/tasks/transfer', { method: 'POST', body: data }).then(parseApiResponse)
  const createExpansionTask = (data: any) => useApiFetch('/tasks/expansion', { method: 'POST', body: data }).then(parseApiResponse)
  const getExpansionAccounts = () => useApiFetch('/tasks/expansion/accounts').then(parseApiResponse)
  const getTasks = (params?: any) => useApiFetch('/tasks', { params }).then(parseApiResponse)
  const getTaskStatus = (id: number) => useApiFetch(`/tasks/${id}`).then(parseApiResponse)
  const startTask = (id: number) => useApiFetch(`/tasks/${id}/start`, { method: 'POST' }).then(parseApiResponse)
  const stopTask = (id: number) => useApiFetch(`/tasks/${id}/stop`, { method: 'POST' }).then(parseApiResponse)
  const pauseTask = (id: number) => useApiFetch(`/tasks/${id}/pause`, { method: 'POST' }).then(parseApiResponse)
  const deleteTask = (id: number) => useApiFetch(`/tasks/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const getTaskItems = (id: number, params?: any) => useApiFetch(`/tasks/${id}/items`, { params }).then(parseApiResponse)
  const getExpansionOutput = (accountId: number) => useApiFetch(`/tasks/expansion/accounts/${accountId}/output`).then(parseApiResponse)
  return { createBatchTransferTask, createExpansionTask, getExpansionAccounts, getTasks, getTaskStatus, startTask, stopTask, pauseTask, deleteTask, getTaskItems, getExpansionOutput }
}

// 日志函数：只在开发环境打印
function log(...args: any[]) {
  if (process.env.NODE_ENV !== 'production') {
    console.log(...args)
  }
}

// Telegram Bot管理API
export const useTelegramApi = () => {
  const getBotConfig = () => useApiFetch('/telegram/bot-config').then(parseApiResponse)
  const updateBotConfig = (data: any) => useApiFetch('/telegram/bot-config', { method: 'PUT', body: data }).then(parseApiResponse)
  const validateApiKey = (data: any) => useApiFetch('/telegram/validate-api-key', { method: 'POST', body: data }).then(parseApiResponse)
  const getBotStatus = () => useApiFetch('/telegram/bot-status').then(parseApiResponse)
  const debugBotConnection = () => useApiFetch('/telegram/debug-connection').then(parseApiResponse)
  const reloadBotConfig = () => useApiFetch('/telegram/reload-config', { method: 'POST' }).then(parseApiResponse)
  const testBotMessage = (data: any) => useApiFetch('/telegram/test-message', { method: 'POST', body: data }).then(parseApiResponse)
  const getChannels = () => useApiFetch('/telegram/channels').then(parseApiResponse)
  const createChannel = (data: any) => useApiFetch('/telegram/channels', { method: 'POST', body: data }).then(parseApiResponse)
  const updateChannel = (id: number, data: any) => useApiFetch(`/telegram/channels/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteChannel = (id: number) => useApiFetch(`/telegram/channels/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const getLogs = (params?: any) => useApiFetch('/telegram/logs', { params }).then(parseApiResponse)
  const getLogStats = (params?: any) => useApiFetch('/telegram/logs/stats', { params }).then(parseApiResponse)
  const clearLogs = (params?: any) => useApiFetch('/telegram/logs/clear', { method: 'POST', body: params }).then(parseApiResponse)
  return {
    getBotConfig,
    updateBotConfig,
    validateApiKey,
    getBotStatus,
    debugBotConnection,
    reloadBotConfig,
    testBotMessage,
    getChannels,
    createChannel,
    updateChannel,
    deleteChannel,
    getLogs,
    getLogStats,
    clearLogs
  }
}

// Meilisearch管理API
export const useMeilisearchApi = () => {
  const getStatus = () => useApiFetch('/meilisearch/status').then(parseApiResponse)
  const getUnsyncedCount = () => useApiFetch('/meilisearch/unsynced-count').then(parseApiResponse)
  const getUnsyncedResources = (params?: any) => useApiFetch('/meilisearch/unsynced', { params }).then(parseApiResponse)
  const getSyncedResources = (params?: any) => useApiFetch('/meilisearch/synced', { params }).then(parseApiResponse)
  const getAllResources = (params?: any) => useApiFetch('/meilisearch/resources', { params }).then(parseApiResponse)
  const testConnection = (data: any) => useApiFetch('/meilisearch/test-connection', { method: 'POST', body: data }).then(parseApiResponse)
  const syncAllResources = () => useApiFetch('/meilisearch/sync-all', { method: 'POST' }).then(parseApiResponse)
  const stopSync = () => useApiFetch('/meilisearch/stop-sync', { method: 'POST' }).then(parseApiResponse)
  const clearIndex = () => useApiFetch('/meilisearch/clear-index', { method: 'POST' }).then(parseApiResponse)
  const updateIndexSettings = () => useApiFetch('/meilisearch/update-settings', { method: 'POST' }).then(parseApiResponse)
  const getSyncProgress = () => useApiFetch('/meilisearch/sync-progress').then(parseApiResponse)
  const debugGetAllDocuments = () => useApiFetch('/meilisearch/debug/documents').then(parseApiResponse)
  return {
    getStatus,
    getUnsyncedCount,
    getUnsyncedResources,
    getSyncedResources,
    getAllResources,
    testConnection,
    syncAllResources,
    stopSync,
    clearIndex,
    updateIndexSettings,
    getSyncProgress,
    debugGetAllDocuments
  }
}

// API访问日志管理API
export const useApiAccessLogApi = () => {
  const getApiAccessLogs = (params?: any) => useApiFetch('/api/api-access-logs', { params }).then(parseApiResponse)
  const getApiAccessLogSummary = () => useApiFetch('/api/api-access-logs/summary').then(parseApiResponse)
  const getApiAccessLogStats = () => useApiFetch('/api/api-access-logs/stats').then(parseApiResponse)
  const clearApiAccessLogs = (days: number) => useApiFetch('/api/api-access-logs', { method: 'DELETE', body: { days } }).then(parseApiResponse)
  return {
    getApiAccessLogs,
    getApiAccessLogSummary,
    getApiAccessLogStats,
    clearApiAccessLogs
  }
}

// 系统日志管理API
export const useSystemLogApi = () => {
  const getSystemLogs = (params?: any) => useApiFetch('/api/system-logs', { params }).then(parseApiResponse)
  const getSystemLogFiles = () => useApiFetch('/api/system-logs/files').then(parseApiResponse)
  const getSystemLogSummary = () => useApiFetch('/api/system-logs/summary').then(parseApiResponse)
  const clearSystemLogs = (days: number) => useApiFetch('/api/system-logs', { method: 'DELETE', body: { days } }).then(parseApiResponse)
  return {
    getSystemLogs,
    getSystemLogFiles,
    getSystemLogSummary,
    clearSystemLogs
  }
}

// 微信机器人管理API
export const useWechatApi = () => {
  const getBotConfig = () => useApiFetch('/wechat/bot-config').then(parseApiResponse)
  const updateBotConfig = (data: any) => useApiFetch('/wechat/bot-config', { method: 'PUT', body: data }).then(parseApiResponse)
  const getBotStatus = () => useApiFetch('/wechat/bot-status').then(parseApiResponse)
  const uploadVerifyFile = (formData: FormData) => useApiFetch('/wechat/verify-file', { method: 'POST', body: formData }).then(parseApiResponse)
  return {
    getBotConfig,
    updateBotConfig,
    getBotStatus,
    uploadVerifyFile
  }
}