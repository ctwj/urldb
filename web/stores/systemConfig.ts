import { defineStore } from 'pinia'
import { useApiFetch } from '~/composables/useApiFetch'

export const useSystemConfigStore = defineStore('systemConfig', {
  state: () => ({
    config: null as any,
    initialized: false
  }),
  actions: {
    async initConfig(force = false, useAdminApi = false) {
      if (this.initialized && !force) return
      try {
        // 根据上下文选择API：管理员页面使用管理员API，其他页面使用公开API
        const apiUrl = useAdminApi ? '/system/config' : '/public/system-config'
        const response = await useApiFetch(apiUrl)
        console.log('Store API响应:', response) // 调试信息
        
        // 正确处理API响应结构
        const data = response.data || response
        console.log('Store 处理后的数据:', data) // 调试信息
        
        this.config = data
        this.initialized = true
      } catch (e) {
        console.error('Store 获取系统配置失败:', e) // 调试信息
        // 可根据需要处理错误
        this.config = null
        this.initialized = false
      }
    },
    setConfig(newConfig: any) {
      this.config = newConfig
      this.initialized = true
    }
  }
}) 