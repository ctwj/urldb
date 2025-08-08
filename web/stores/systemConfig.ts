import { defineStore } from 'pinia'
import { useApiFetch } from '~/composables/useApiFetch'

export const useSystemConfigStore = defineStore('systemConfig', {
  state: () => ({
    config: null as any,
    initialized: false
  }),
  actions: {
    async initConfig(force = false) {
      if (this.initialized && !force) return
      try {
        // 使用公开的系统配置API，不需要管理员权限
        const response = await useApiFetch('/public/system-config')
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