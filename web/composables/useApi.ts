// 统一响应解析函数
export const parseApiResponse = <T>(response: any): T => {
  console.log('parseApiResponse - 原始响应:', response)
  
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
      return response.data
    } else {
      throw new Error(response.message || '请求失败')
    }
  }
  
  // 兼容旧格式，直接返回响应
  return response
}

// 使用 $fetch 替代 axios，更好地处理 SSR
export const useResourceApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getResources = async (params?: any) => {
    const response = await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getResource = async (id: number) => {
    const response = await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createResource = async (data: any) => {
    const response = await $fetch('/resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updateResource = async (id: number, data: any) => {
    const response = await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteResource = async (id: number) => {
    const response = await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const searchResources = async (params: any) => {
    const response = await $fetch('/search', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getResourcesByPan = async (panId: number, params?: any) => {
    const response = await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params: { ...params, pan_id: panId },
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getResources,
    getResource,
    createResource,
    updateResource,
    deleteResource,
    searchResources,
    getResourcesByPan,
  }
}

// 认证相关API
export const useAuthApi = () => {
  const config = useRuntimeConfig()
  
  const login = async (data: any) => {
    const response = await $fetch('/auth/login', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
    return parseApiResponse(response)
  }

  const register = async (data: any) => {
    const response = await $fetch('/auth/register', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
    return parseApiResponse(response)
  }

  const getProfile = async () => {
    const token = localStorage.getItem('token')
    const response = await $fetch('/auth/profile', {
      baseURL: config.public.apiBase,
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    return parseApiResponse(response)
  }

  return {
    login,
    register,
    getProfile,
  }
}

// 分类相关API
export const useCategoryApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getCategories = async () => {
    const response = await $fetch('/categories', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createCategory = async (data: any) => {
    const response = await $fetch('/categories', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updateCategory = async (id: number, data: any) => {
    const response = await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteCategory = async (id: number) => {
    const response = await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  }
}

// 平台相关API
export const usePanApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getPans = async () => {
    const response = await $fetch('/pans', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getPan = async (id: number) => {
    const response = await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createPan = async (data: any) => {
    const response = await $fetch('/pans', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updatePan = async (id: number, data: any) => {
    const response = await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deletePan = async (id: number) => {
    const response = await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getPans,
    getPan,
    createPan,
    updatePan,
    deletePan,
  }
}

// Cookie相关API
export const useCksApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getCks = async (params?: any) => {
    const response = await $fetch('/cks', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getCksByID = async (id: number) => {
    const response = await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createCks = async (data: any) => {
    const response = await $fetch('/cks', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updateCks = async (id: number, data: any) => {
    const response = await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteCks = async (id: number) => {
    const response = await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getCks,
    getCksByID,
    createCks,
    updateCks,
    deleteCks,
  }
}

// 标签相关API
export const useTagApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getTags = async () => {
    const response = await $fetch('/tags', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getTagsByCategory = async (categoryId: number, params?: any) => {
    const response = await $fetch(`/categories/${categoryId}/tags`, {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getTag = async (id: number) => {
    const response = await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createTag = async (data: any) => {
    const response = await $fetch('/tags', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updateTag = async (id: number, data: any) => {
    const response = await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteTag = async (id: number) => {
    const response = await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const getResourceTags = async (resourceId: number) => {
    const response = await $fetch(`/resources/${resourceId}/tags`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getTags,
    getTagsByCategory,
    getTag,
    createTag,
    updateTag,
    deleteTag,
    getResourceTags,
  }
}

// 待处理资源相关API
export const useReadyResourceApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getReadyResources = async (params?: any) => {
    const response = await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createReadyResource = async (data: any) => {
    const response = await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const batchCreateReadyResources = async (data: any) => {
    const response = await $fetch('/ready-resources/batch', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const createReadyResourcesFromText = async (text: string) => {
    const formData = new FormData()
    formData.append('text', text)
    
    const response = await $fetch('/ready-resources/text', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: formData,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteReadyResource = async (id: number) => {
    const response = await $fetch(`/ready-resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const clearReadyResources = async () => {
    const response = await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getReadyResources,
    createReadyResource,
    batchCreateReadyResources,
    createReadyResourcesFromText,
    deleteReadyResource,
    clearReadyResources,
  }
}

// 统计相关API
export const useStatsApi = () => {
  const config = useRuntimeConfig()
  
  const getStats = async () => {
    const response = await $fetch('/stats', {
      baseURL: config.public.apiBase,
    })
    return parseApiResponse(response)
  }

  return {
    getStats,
  }
}

// 系统配置相关API
export const useSystemConfigApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getSystemConfig = async () => {
    const response = await $fetch('/system/config', {
      baseURL: config.public.apiBase,
      // GET接口不需要认证头
    })
    return parseApiResponse(response)
  }

  const updateSystemConfig = async (data: any) => {
    const response = await $fetch('/system/config', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getSystemConfig,
    updateSystemConfig,
  }
}

// 热播剧相关API
export const useHotDramaApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getHotDramas = async (params?: any) => {
    const response = await $fetch('/hot-dramas', {
      baseURL: config.public.apiBase,
      params,
    })
    return parseApiResponse(response)
  }

  const createHotDrama = async (data: any) => {
    const response = await $fetch('/hot-dramas', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const updateHotDrama = async (id: number, data: any) => {
    const response = await $fetch(`/hot-dramas/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const deleteHotDrama = async (id: number) => {
    const response = await $fetch(`/hot-dramas/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  const fetchHotDramas = async () => {
    const response = await $fetch('/hot-dramas/fetch', {
      baseURL: config.public.apiBase,
      method: 'POST',
      headers: getAuthHeaders() as Record<string, string>
    })
    return parseApiResponse(response)
  }

  return {
    getHotDramas,
    createHotDrama,
    updateHotDrama,
    deleteHotDrama,
    fetchHotDramas,
  }
}

// 监控相关API
export const useMonitorApi = () => {
  const config = useRuntimeConfig()
  
  const getPerformanceStats = async () => {
    const response = await $fetch('/performance', {
      baseURL: config.public.apiBase,
    })
    return parseApiResponse(response)
  }

  const getSystemInfo = async () => {
    const response = await $fetch('/system/info', {
      baseURL: config.public.apiBase,
    })
    return parseApiResponse(response)
  }

  const getBasicStats = async () => {
    const response = await $fetch('/stats', {
      baseURL: config.public.apiBase,
    })
    return parseApiResponse(response)
  }

  return {
    getPerformanceStats,
    getSystemInfo,
    getBasicStats,
  }
} 