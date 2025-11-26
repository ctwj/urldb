import type { UseApiFetchOptions } from 'nuxt/app'

// Bing API 相关类型定义
export interface BingIndexConfig {
  enabled: boolean
  submitInterval: number
  batchSize: number
  retryCount: number
}


export interface UpdateBingConfigRequest {
  enabled: boolean
}

// Bing API Hook
export const useBingApi = () => {
  const { $fetch } = useNuxtApp()

  // 获取Bing配置
  const getConfig = async (): Promise<{ success: boolean; data: BingIndexConfig }> => {
    const options: UseApiFetchOptions = {
      method: 'GET'
    }

    return $fetch<{ success: boolean; data: BingIndexConfig }>('/api/bing/config', options)
  }

  // 更新Bing配置
  const updateConfig = async (data: UpdateBingConfigRequest): Promise<{ success: boolean; message: string }> => {
    const options: UseApiFetchOptions = {
      method: 'POST',
      body: data
    }

    return $fetch<{ success: boolean; message: string }>('/api/bing/config', options)
  }

  return {
    getConfig,
    updateConfig
  }
}