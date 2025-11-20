import { defineStore } from 'pinia'
import { useApiFetch } from '~/composables/useApiFetch'
import { parseApiResponse } from '~/composables/useApi'

// 缓存配置
const CACHE_KEY = 'system-config-cache'
const CACHE_TIMESTAMP_KEY = 'system-config-cache-timestamp'
const CACHE_DURATION = 30 * 60 * 1000 // 30分钟缓存

// 安全的客户端检查函数
const isClient = () => {
  return typeof process !== 'undefined' && process.client
}

interface CacheData {
  config: any
  timestamp: number
  version?: string
}

export const useSystemConfigStore = defineStore('systemConfig', {
  state: () => ({
    config: null as any,
    initialized: false,
    lastFetchTime: 0 as number,
    isLoading: false as boolean,
    error: null as string | null
  }),

  getters: {
    // 检查缓存是否有效
    isCacheValid(): boolean {
      if (!isClient()) return false

      try {
        const cacheData = localStorage.getItem(CACHE_KEY)
        const cacheTimestamp = localStorage.getItem(CACHE_TIMESTAMP_KEY)

        if (!cacheData || !cacheTimestamp) return false

        const timestamp = parseInt(cacheTimestamp)
        const now = Date.now()

        // 检查缓存是否过期
        const isValid = (now - timestamp) < CACHE_DURATION

        console.log(`[SystemConfig] 缓存检查: ${isValid ? '有效' : '已过期'}, 剩余时间: ${Math.max(0, CACHE_DURATION - (now - timestamp)) / 1000 / 60}分钟`)

        return isValid
      } catch (error) {
        console.error('[SystemConfig] 缓存检查失败:', error)
        return false
      }
    },

    // 获取缓存的数据
    cachedConfig(): any {
      if (!isClient() || !this.isCacheValid) return null

      try {
        const cacheData = localStorage.getItem(CACHE_KEY)
        if (cacheData) {
          const parsed = JSON.parse(cacheData) as CacheData
          console.log('[SystemConfig] 使用缓存数据')
          return parsed.config
        }
      } catch (error) {
        console.error('[SystemConfig] 读取缓存失败:', error)
        this.clearCache()
      }

      return null
    },

    // 获取缓存剩余时间（秒）
    cacheTimeRemaining(): number {
      if (!isClient()) return 0

      try {
        const cacheTimestamp = localStorage.getItem(CACHE_TIMESTAMP_KEY)
        if (!cacheTimestamp) return 0

        const timestamp = parseInt(cacheTimestamp)
        const now = Date.now()
        const remaining = Math.max(0, CACHE_DURATION - (now - timestamp))

        return Math.floor(remaining / 1000)
      } catch (error) {
        return 0
      }
    }
  },

  actions: {
    // 清除缓存
    clearCache() {
      if (!isClient()) return

      try {
        localStorage.removeItem(CACHE_KEY)
        localStorage.removeItem(CACHE_TIMESTAMP_KEY)
        console.log('[SystemConfig] 缓存已清除')
      } catch (error) {
        console.error('[SystemConfig] 清除缓存失败:', error)
      }
    },

    // 保存到缓存
    saveToCache(config: any) {
      if (!isClient()) return

      try {
        const cacheData: CacheData = {
          config,
          timestamp: Date.now(),
          version: '1.0'
        }

        localStorage.setItem(CACHE_KEY, JSON.stringify(cacheData))
        localStorage.setItem(CACHE_TIMESTAMP_KEY, cacheData.timestamp.toString())

        console.log('[SystemConfig] 配置已缓存，有效期30分钟')
      } catch (error) {
        console.error('[SystemConfig] 保存缓存失败:', error)
      }
    },

    // 从缓存加载
    loadFromCache(): boolean {
      if (!isClient()) return false

      const cachedConfig = this.cachedConfig
      if (cachedConfig) {
        this.config = cachedConfig
        this.initialized = true

        // 从缓存时间戳设置 lastFetchTime
        const cacheTimestamp = localStorage.getItem(CACHE_TIMESTAMP_KEY)
        if (cacheTimestamp) {
          this.lastFetchTime = parseInt(cacheTimestamp)
        } else {
          this.lastFetchTime = Date.now()
        }

        console.log('[SystemConfig] 从缓存加载配置成功')
        return true
      }

      return false
    },

    // 初始化配置（带缓存支持）
    async initConfig(force = false, useAdminApi = false) {
      // 如果已经初始化且不强制刷新，直接返回
      if (this.initialized && !force) {
        console.log('[SystemConfig] 配置已初始化，直接返回')
        return
      }

      // 如果不强制刷新，先尝试从缓存加载
      if (!force && this.loadFromCache()) {
        return
      }

      // 防止重复请求
      if (this.isLoading) {
        console.log('[SystemConfig] 正在加载中，等待完成...')
        return
      }

      this.isLoading = true
      this.error = null

      try {
        // console.log(`[SystemConfig] 开始获取配置 (force: ${force}, useAdminApi: ${useAdminApi})`)

        // 根据上下文选择API：管理员页面使用管理员API，其他页面使用公开API
        const apiUrl = useAdminApi ? '/system/config' : '/public/system-config'
        const response = await useApiFetch(apiUrl)

        // 使用parseApiResponse正确解析API响应
        const data = parseApiResponse(response)

        this.config = data
        this.initialized = true
        this.lastFetchTime = Date.now()
        this.isLoading = false

        // 保存到缓存（仅在客户端）
        this.saveToCache(data)

        // console.log('[SystemConfig] 配置获取并缓存成功')
        // console.log('[SystemConfig] 自动处理状态:', data.auto_process_ready_resources)
        // console.log('[SystemConfig] 自动转存状态:', data.auto_transfer_enabled)

      } catch (error) {
        this.isLoading = false
        this.error = error instanceof Error ? error.message : '获取配置失败'
        // console.error('[SystemConfig] 获取系统配置失败:', error)

        // 如果网络请求失败，尝试使用过期的缓存作为降级方案
        if (!force) {
          try {
            const expiredCache = localStorage.getItem(CACHE_KEY)
            if (expiredCache) {
              const parsed = JSON.parse(expiredCache) as CacheData
              this.config = parsed.config
              this.initialized = true
              console.log('[SystemConfig] 网络请求失败，使用过期缓存作为降级方案')
              return
            }
          } catch (cacheError) {
            console.error('[SystemConfig] 降级缓存方案也失败:', cacheError)
          }
        }

        this.config = null
        this.initialized = false
      }
    },

    // 强制刷新配置
    async refreshConfig(useAdminApi = false) {
      console.log('[SystemConfig] 强制刷新配置')
      this.clearCache()
      await this.initConfig(true, useAdminApi)
    },

    // 检查并自动刷新缓存（如果即将过期）
    async checkAndRefreshCache(useAdminApi = false) {
      if (!isClient()) return

      const timeRemaining = this.cacheTimeRemaining

      // 如果缓存剩余时间少于5分钟，自动刷新
      if (timeRemaining > 0 && timeRemaining < 5 * 60) {
        console.log(`[SystemConfig] 缓存即将过期（剩余${timeRemaining}秒），自动刷新`)
        await this.refreshConfig(useAdminApi)
      }
    },

    // 手动设置配置（用于管理员更新配置后）
    setConfig(newConfig: any) {
      this.config = newConfig
      this.initialized = true
      this.lastFetchTime = Date.now()

      // 更新缓存
      this.saveToCache(newConfig)

      console.log('[SystemConfig] 配置已手动更新并缓存')
    },

    // 获取配置状态信息
    getStatus() {
      return {
        initialized: this.initialized,
        isLoading: this.isLoading,
        error: this.error,
        lastFetchTime: this.lastFetchTime,
        cacheTimeRemaining: this.cacheTimeRemaining,
        isCacheValid: this.isCacheValid
      }
    }
  }
})