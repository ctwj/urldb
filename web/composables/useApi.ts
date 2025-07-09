import axios from 'axios'

const config = useRuntimeConfig()
const api = axios.create({
  baseURL: config.public.apiBase,
  timeout: 10000,
})

// 资源相关API
export const useResourceApi = () => {
  const getResources = async (params?: any) => {
    const response = await api.get('/resources', { params })
    return response.data
  }

  const getResource = async (id: number) => {
    const response = await api.get(`/resources/${id}`)
    return response.data
  }

  const createResource = async (data: any) => {
    const response = await api.post('/resources', data)
    return response.data
  }

  const updateResource = async (id: number, data: any) => {
    const response = await api.put(`/resources/${id}`, data)
    return response.data
  }

  const deleteResource = async (id: number) => {
    const response = await api.delete(`/resources/${id}`)
    return response.data
  }

  const searchResources = async (params: any) => {
    const response = await api.get('/search', { params })
    return response.data
  }

  return {
    getResources,
    getResource,
    createResource,
    updateResource,
    deleteResource,
    searchResources,
  }
}

// 分类相关API
export const useCategoryApi = () => {
  const getCategories = async () => {
    const response = await api.get('/categories')
    return response.data
  }

  const createCategory = async (data: any) => {
    const response = await api.post('/categories', data)
    return response.data
  }

  const updateCategory = async (id: number, data: any) => {
    const response = await api.put(`/categories/${id}`, data)
    return response.data
  }

  const deleteCategory = async (id: number) => {
    const response = await api.delete(`/categories/${id}`)
    return response.data
  }

  return {
    getCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  }
}

// 统计相关API
export const useStatsApi = () => {
  const getStats = async () => {
    const response = await api.get('/stats')
    return response.data
  }

  return {
    getStats,
  }
} 