import { useApiFetch } from './useApiFetch'
import { useUserStore } from '~/stores/user'

// 统一响应解析函数
export const parseApiResponse = <T>(response: any): T => {
  log('parseApiResponse - 原始响应:', response)
  
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
      // 特殊处理资源接口返回的data.list格式，转换为resources格式
      if (response.data && response.data.list && Array.isArray(response.data.list)) {
        return {
          resources: response.data.list,
          total: response.data.total,
          page: response.data.page,
          page_size: response.data.limit
        } as T
      }
      // 特殊处理失败资源接口，返回完整的data结构
      if (response.data && response.data.data && Array.isArray(response.data.data) && response.data.total !== undefined) {
        return response.data
      }
      // 特殊处理登录接口，直接返回data部分（包含token和user）
      if (response.data && response.data.token && response.data.user) {
        console.log('parseApiResponse - 登录接口处理，返回data:', response.data)
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
  return { getResources, getResource, createResource, updateResource, deleteResource, searchResources, getResourcesByPan, incrementViewCount, batchDeleteResources, getResourceLink }
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
  return { getCks, getCksByID, createCks, updateCks, deleteCks, refreshCapacity }
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

export const useSystemConfigApi = () => {
  const getSystemConfig = () => useApiFetch('/system/config').then(parseApiResponse)
  const updateSystemConfig = (data: any) => useApiFetch('/system/config', { method: 'POST', body: data }).then(parseApiResponse)
  const toggleAutoProcess = (enabled: boolean) => useApiFetch('/system/config/toggle-auto-process', { method: 'POST', body: { auto_process_ready_resources: enabled } }).then(parseApiResponse)
  return { getSystemConfig, updateSystemConfig, toggleAutoProcess }
}

export const useHotDramaApi = () => {
  const getHotDramas = (params?: any) => useApiFetch('/hot-dramas', { params }).then(parseApiResponse)
  const createHotDrama = (data: any) => useApiFetch('/hot-dramas', { method: 'POST', body: data }).then(parseApiResponse)
  const updateHotDrama = (id: number, data: any) => useApiFetch(`/hot-dramas/${id}`, { method: 'PUT', body: data }).then(parseApiResponse)
  const deleteHotDrama = (id: number) => useApiFetch(`/hot-dramas/${id}`, { method: 'DELETE' }).then(parseApiResponse)
  const fetchHotDramas = () => useApiFetch('/hot-dramas/fetch', { method: 'POST' }).then(parseApiResponse)
  return { getHotDramas, createHotDrama, updateHotDrama, deleteHotDrama, fetchHotDramas }
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

// 日志函数：只在开发环境打印
function log(...args: any[]) {
  if (process.env.NODE_ENV !== 'production') {
    console.log(...args)
  }
} 