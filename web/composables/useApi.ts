// 使用 $fetch 替代 axios，更好地处理 SSR
export const useResourceApi = () => {
  const config = useRuntimeConfig()
  
  const getAuthHeaders = () => {
    const userStore = useUserStore()
    return userStore.authHeaders
  }
  
  const getResources = async (params?: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getResource = async (id: number) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createResource = async (data: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const updateResource = async (id: number, data: any) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deleteResource = async (id: number) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const searchResources = async (params: any) => {
    return await $fetch('/search', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getResourcesByPan = async (panId: number, params?: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params: { ...params, pan_id: panId },
      headers: getAuthHeaders() as Record<string, string>
    })
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
    return await $fetch('/auth/login', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const register = async (data: any) => {
    return await $fetch('/auth/register', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const getProfile = async () => {
    const token = localStorage.getItem('token')
    return await $fetch('/auth/profile', {
      baseURL: config.public.apiBase,
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
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
    return await $fetch('/categories', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createCategory = async (data: any) => {
    return await $fetch('/categories', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const updateCategory = async (id: number, data: any) => {
    return await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deleteCategory = async (id: number) => {
    return await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
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
    return await $fetch('/pans', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getPan = async (id: number) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createPan = async (data: any) => {
    return await $fetch('/pans', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const updatePan = async (id: number, data: any) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deletePan = async (id: number) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
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
    return await $fetch('/cks', {
      baseURL: config.public.apiBase,
      params,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getCksByID = async (id: number) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createCks = async (data: any) => {
    return await $fetch('/cks', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const updateCks = async (id: number, data: any) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deleteCks = async (id: number) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
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
    return await $fetch('/tags', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getTag = async (id: number) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createTag = async (data: any) => {
    return await $fetch('/tags', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const updateTag = async (id: number, data: any) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deleteTag = async (id: number) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const getResourceTags = async (resourceId: number) => {
    return await $fetch(`/resources/${resourceId}/tags`, {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  return {
    getTags,
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
  
  const getReadyResources = async () => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createReadyResource = async (data: any) => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const batchCreateReadyResources = async (data: any) => {
    return await $fetch('/ready-resources/batch', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const createReadyResourcesFromText = async (text: string) => {
    const formData = new FormData()
    formData.append('text', text)
    
    return await $fetch('/ready-resources/text', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: formData,
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const deleteReadyResource = async (id: number) => {
    return await $fetch(`/ready-resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
  }

  const clearReadyResources = async () => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'DELETE',
      headers: getAuthHeaders() as Record<string, string>
    })
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
    return await $fetch('/stats', {
      baseURL: config.public.apiBase,
    })
  }

  return {
    getStats,
  }
} 