// 使用 $fetch 替代 axios，更好地处理 SSR
export const useResourceApi = () => {
  const config = useRuntimeConfig()
  
  const getResources = async (params?: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params
    })
  }

  const getResource = async (id: number) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
    })
  }

  const createResource = async (data: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const updateResource = async (id: number, data: any) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data
    })
  }

  const deleteResource = async (id: number) => {
    return await $fetch(`/resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
    })
  }

  const searchResources = async (params: any) => {
    return await $fetch('/search', {
      baseURL: config.public.apiBase,
      params
    })
  }

  const getResourcesByPan = async (panId: number, params?: any) => {
    return await $fetch('/resources', {
      baseURL: config.public.apiBase,
      params: { ...params, pan_id: panId }
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

// 分类相关API
export const useCategoryApi = () => {
  const config = useRuntimeConfig()
  
  const getCategories = async () => {
    return await $fetch('/categories', {
      baseURL: config.public.apiBase,
    })
  }

  const createCategory = async (data: any) => {
    return await $fetch('/categories', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const updateCategory = async (id: number, data: any) => {
    return await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data
    })
  }

  const deleteCategory = async (id: number) => {
    return await $fetch(`/categories/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
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
  
  const getPans = async () => {
    return await $fetch('/pans', {
      baseURL: config.public.apiBase,
    })
  }

  const getPan = async (id: number) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
    })
  }

  const createPan = async (data: any) => {
    return await $fetch('/pans', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const updatePan = async (id: number, data: any) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data
    })
  }

  const deletePan = async (id: number) => {
    return await $fetch(`/pans/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
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
  
  const getCks = async (params?: any) => {
    return await $fetch('/cks', {
      baseURL: config.public.apiBase,
      params
    })
  }

  const getCksByID = async (id: number) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
    })
  }

  const createCks = async (data: any) => {
    return await $fetch('/cks', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const updateCks = async (id: number, data: any) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data
    })
  }

  const deleteCks = async (id: number) => {
    return await $fetch(`/cks/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
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
  
  const getTags = async () => {
    return await $fetch('/tags', {
      baseURL: config.public.apiBase,
    })
  }

  const getTag = async (id: number) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
    })
  }

  const createTag = async (data: any) => {
    return await $fetch('/tags', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const updateTag = async (id: number, data: any) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'PUT',
      body: data
    })
  }

  const deleteTag = async (id: number) => {
    return await $fetch(`/tags/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
    })
  }

  const getResourceTags = async (resourceId: number) => {
    return await $fetch(`/resources/${resourceId}/tags`, {
      baseURL: config.public.apiBase,
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

// 统计相关API
// 待处理资源相关API
export const useReadyResourceApi = () => {
  const config = useRuntimeConfig()
  
  const getReadyResources = async () => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
    })
  }

  const createReadyResource = async (data: any) => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const batchCreateReadyResources = async (data: any) => {
    return await $fetch('/ready-resources/batch', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: data
    })
  }

  const createReadyResourcesFromText = async (text: string) => {
    return await $fetch('/ready-resources/text', {
      baseURL: config.public.apiBase,
      method: 'POST',
      body: { text }
    })
  }

  const deleteReadyResource = async (id: number) => {
    return await $fetch(`/ready-resources/${id}`, {
      baseURL: config.public.apiBase,
      method: 'DELETE'
    })
  }

  const clearReadyResources = async () => {
    return await $fetch('/ready-resources', {
      baseURL: config.public.apiBase,
      method: 'DELETE'
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