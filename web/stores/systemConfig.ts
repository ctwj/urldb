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
        const data = await useApiFetch('/system/config').then((res: any) => res.data || res)
        this.config = data
        this.initialized = true
      } catch (e) {
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