import { useApiFetch } from './useApiFetch'

export const useFileApi = () => {
  const getFileList = (params?: any) => useApiFetch('/files', { params }).then(parseApiResponse)
  const uploadFile = (data: FormData) => useApiFetch('/files/upload', { 
    method: 'POST', 
    body: data,
    headers: {
      // 不设置Content-Type，让浏览器自动设置multipart/form-data
    }
  }).then(parseApiResponse)
  const deleteFiles = (ids: number[]) => useApiFetch('/files', { 
    method: 'DELETE', 
    body: { ids } 
  }).then(parseApiResponse)
  const updateFile = (data: any) => useApiFetch('/files', { 
    method: 'PUT', 
    body: data 
  }).then(parseApiResponse)

  return {
    getFileList,
    uploadFile,
    deleteFiles,
    updateFile
  }
}

// 解析API响应
function parseApiResponse(response: any) {
  if (response.success) {
    return response
  } else {
    throw new Error(response.message || '请求失败')
  }
} 