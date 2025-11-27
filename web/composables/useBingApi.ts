import { useApiFetch } from './useApiFetch'

// Bing API 相关类型定义
export interface BingIndexConfig {
  enabled: boolean
  submitInterval?: number
  batchSize?: number
  retryCount?: number
}

export interface UpdateBingConfigRequest {
  enabled: boolean
}

// Bing API Hook
export const useBingApi = () => {
  // 获取Bing配置
  const getConfig = async (): Promise<{ success: boolean; data: BingIndexConfig }> => {
    return useApiFetch('/bing/config', { method: 'GET' })
  }

  // 更新Bing配置
  const updateConfig = async (data: UpdateBingConfigRequest): Promise<{ success: boolean; message: string }> => {
    return useApiFetch('/bing/config', { 
      method: 'POST', 
      body: data 
    })
  }

  return {
    getConfig,
    updateConfig
  }
}