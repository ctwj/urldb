import { ref } from 'vue'
import type { Ref } from 'vue'

export interface ConfigChangeDetectionOptions {
  // 是否启用自动检测
  autoDetect?: boolean
  // 是否在控制台输出调试信息
  debug?: boolean
  // 自定义比较函数
  customCompare?: (key: string, currentValue: any, originalValue: any) => boolean
  // 配置项映射（前端字段名 -> 后端字段名）
  fieldMapping?: Record<string, string>
}

export interface ConfigSubmitOptions {
  // 是否只提交改动的字段
  onlyChanged?: boolean
  // 是否包含所有配置项（用于后端识别）
  includeAllFields?: boolean
  // 自定义提交数据转换
  transformSubmitData?: (data: any) => any
}

export const useConfigChangeDetection = <T extends Record<string, any>>(
  options: ConfigChangeDetectionOptions = {}
) => {
  const { autoDetect = true, debug = false, customCompare, fieldMapping = {} } = options
  
  // 原始配置数据
  const originalConfig = ref<T>({} as T)
  
  // 当前配置数据
  const currentConfig = ref<T>({} as T)
  
  // 是否已初始化
  const isInitialized = ref(false)

  /**
   * 设置原始配置数据
   */
  const setOriginalConfig = (config: T) => {
    originalConfig.value = { ...config }
    currentConfig.value = { ...config }
    isInitialized.value = true
    
    if (debug) {
      console.log('useConfigChangeDetection - 设置原始配置:', config)
    }
  }

  /**
   * 更新当前配置数据
   */
  const updateCurrentConfig = (config: Partial<T>) => {
    currentConfig.value = { ...currentConfig.value, ...config }
    
    if (debug) {
      console.log('useConfigChangeDetection - 更新当前配置:', config)
    }
  }

  /**
   * 检测配置改动
   */
  const getChangedConfig = (): Partial<T> => {
    if (!isInitialized.value) {
      if (debug) {
        console.warn('useConfigChangeDetection - 配置未初始化')
      }
      return {}
    }

    const changedConfig: Partial<T> = {}

    // 遍历所有配置项
    for (const key in currentConfig.value) {
      const currentValue = currentConfig.value[key]
      const originalValue = originalConfig.value[key]

      // 使用自定义比较函数或默认比较
      const hasChanged = customCompare
        ? customCompare(key, currentValue, originalValue)
        : currentValue !== originalValue

      if (hasChanged) {
        changedConfig[key as keyof T] = currentValue
      }
    }

    if (debug) {
      console.log('useConfigChangeDetection - 检测到的改动:', changedConfig)
    }

    return changedConfig
  }

  /**
   * 检查是否有改动
   */
  const hasChanges = (): boolean => {
    const changedConfig = getChangedConfig()
    return Object.keys(changedConfig).length > 0
  }

  /**
   * 获取改动的字段列表
   */
  const getChangedFields = (): string[] => {
    const changedConfig = getChangedConfig()
    return Object.keys(changedConfig)
  }

  /**
   * 获取改动的详细信息
   */
  const getChangedDetails = (): Array<{
    key: string
    originalValue: any
    currentValue: any
  }> => {
    if (!isInitialized.value) {
      return []
    }

    const details: Array<{
      key: string
      originalValue: any
      currentValue: any
    }> = []
    
    for (const key in currentConfig.value) {
      const currentValue = currentConfig.value[key]
      const originalValue = originalConfig.value[key]
      
      const hasChanged = customCompare 
        ? customCompare(key, currentValue, originalValue)
        : currentValue !== originalValue
      
      if (hasChanged) {
        details.push({
          key,
          originalValue,
          currentValue
        })
      }
    }
    
    return details
  }

  /**
   * 重置为原始配置
   */
  const resetToOriginal = () => {
    currentConfig.value = { ...originalConfig.value }
    
    if (debug) {
      console.log('useConfigChangeDetection - 重置为原始配置')
    }
  }

  /**
   * 更新原始配置（通常在保存成功后调用）
   */
  const updateOriginalConfig = () => {
    originalConfig.value = { ...currentConfig.value }
    
    if (debug) {
      console.log('useConfigChangeDetection - 更新原始配置')
    }
  }

  /**
   * 获取配置快照
   */
  const getSnapshot = () => {
    return {
      original: { ...originalConfig.value },
      current: { ...currentConfig.value },
      changed: getChangedConfig(),
      hasChanges: hasChanges()
    }
  }

  /**
   * 准备提交数据
   */
  const prepareSubmitData = (submitOptions: ConfigSubmitOptions = {}): any => {
    const { onlyChanged = true, includeAllFields = true, transformSubmitData } = submitOptions
    
    let submitData: any = {}
    
    if (onlyChanged) {
      // 只提交改动的字段
      submitData = getChangedConfig()
    } else {
      // 提交所有字段
      submitData = { ...currentConfig.value }
    }
    
    // 应用字段映射
    if (Object.keys(fieldMapping).length > 0) {
      const mappedData: any = {}
      for (const [frontendKey, backendKey] of Object.entries(fieldMapping)) {
        if (submitData[frontendKey] !== undefined) {
          mappedData[backendKey] = submitData[frontendKey]
        }
      }
      submitData = mappedData
    }
    
    // 如果包含所有字段，添加未改动的字段（值为undefined，让后端知道这些字段存在但未改动）
    if (includeAllFields && onlyChanged) {
      for (const key in originalConfig.value) {
        if (submitData[key] === undefined) {
          submitData[key] = undefined
        }
      }
    }
    
    // 应用自定义转换
    if (transformSubmitData) {
      submitData = transformSubmitData(submitData)
    }
    
    if (debug) {
      console.log('useConfigChangeDetection - 准备提交数据:', submitData)
    }
    
    return submitData
  }

  /**
   * 通用配置保存函数
   */
  const saveConfig = async (
    apiFunction: (data: any) => Promise<any>,
    submitOptions: ConfigSubmitOptions = {},
    onSuccess?: () => void,
    onError?: (error: any) => void
  ) => {
    try {
      // 检测是否有改动
      if (!hasChanges()) {
        if (debug) {
          console.log('useConfigChangeDetection - 没有检测到改动，跳过保存')
        }
        return { success: true, message: '没有检测到任何改动' }
      }
      
      // 准备提交数据
      const submitData = prepareSubmitData(submitOptions)
      
      if (debug) {
        console.log('useConfigChangeDetection - 提交数据:', submitData)
      }
      
      // 调用API
      const response = await apiFunction(submitData)
      
      // 更新原始配置
      updateOriginalConfig()
      
      if (debug) {
        console.log('useConfigChangeDetection - 保存成功')
      }
      
      // 调用成功回调
      if (onSuccess) {
        onSuccess()
      }
      
      return { success: true, response }
    } catch (error) {
      if (debug) {
        console.error('useConfigChangeDetection - 保存失败:', error)
      }
      
      // 调用错误回调
      if (onError) {
        onError(error)
      }
      
      throw error
    }
  }

  return {
    // 响应式数据
    originalConfig: originalConfig as Ref<T>,
    currentConfig: currentConfig as Ref<T>,
    isInitialized,
    
    // 方法
    setOriginalConfig,
    updateCurrentConfig,
    getChangedConfig,
    hasChanges,
    getChangedFields,
    getChangedDetails,
    resetToOriginal,
    updateOriginalConfig,
    getSnapshot,
    prepareSubmitData,
    saveConfig
  }
} 