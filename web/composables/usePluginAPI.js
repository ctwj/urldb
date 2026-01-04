/**
 * 插件管理API工具
 * 提供插件相关的API调用方法
 */

export const usePluginAPI = () => {
  // 获取插件列表
  const getPlugins = async (params = {}) => {
    const query = new URLSearchParams(params).toString()
    const url = query ? `/api/plugins?${query}` : '/api/plugins'

    try {
      const response = await $fetch(url)
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '获取插件列表失败')
      }
    } catch (error) {
      console.error('获取插件列表失败:', error)
      throw error
    }
  }

  // 获取插件详情
  const getPlugin = async (pluginName) => {
    try {
      const response = await $fetch(`/api/plugins/${pluginName}`)
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '获取插件详情失败')
      }
    } catch (error) {
      console.error('获取插件详情失败:', error)
      throw error
    }
  }

  // 获取插件统计
  const getPluginStats = async () => {
    try {
      const response = await $fetch('/api/plugins/stats')
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '获取插件统计失败')
      }
    } catch (error) {
      console.error('获取插件统计失败:', error)
      throw error
    }
  }

  // 启用插件
  const enablePlugin = async (pluginName) => {
    try {
      const response = await $fetch(`/api/plugins/${pluginName}/enable`, {
        method: 'POST'
      })
      if (response.success) {
        return response
      } else {
        throw new Error(response.error || '启用插件失败')
      }
    } catch (error) {
      console.error('启用插件失败:', error)
      throw error
    }
  }

  // 禁用插件
  const disablePlugin = async (pluginName) => {
    try {
      const response = await $fetch(`/api/plugins/${pluginName}/disable`, {
        method: 'POST'
      })
      if (response.success) {
        return response
      } else {
        throw new Error(response.error || '禁用插件失败')
      }
    } catch (error) {
      console.error('禁用插件失败:', error)
      throw error
    }
  }

  // 更新插件配置
  const updatePluginConfig = async (pluginName, config) => {
    try {
      const response = await $fetch(`/api/plugins/${pluginName}/config`, {
        method: 'PUT',
        body: {
          config
        }
      })
      if (response.success) {
        return response
      } else {
        throw new Error(response.error || '更新插件配置失败')
      }
    } catch (error) {
      console.error('更新插件配置失败:', error)
      throw error
    }
  }

  // 获取插件日志
  const getPluginLogs = async (pluginName, params = {}) => {
    const query = new URLSearchParams(params).toString()
    const url = query ? `/api/plugins/${pluginName}/logs?${query}` : `/api/plugins/${pluginName}/logs`

    try {
      const response = await $fetch(url)
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '获取插件日志失败')
      }
    } catch (error) {
      console.error('获取插件日志失败:', error)
      throw error
    }
  }

  // 获取插件市场
  const getPluginMarket = async () => {
    try {
      const response = await $fetch('/api/plugins/market')
      if (response.success) {
        return response.data
      } else {
        throw new Error(response.error || '获取插件市场失败')
      }
    } catch (error) {
      console.error('获取插件市场失败:', error)
      throw error
    }
  }

  // 安装插件
  const installPlugin = async (pluginData) => {
    try {
      const response = await $fetch('/api/plugins/install', {
        method: 'POST',
        body: pluginData
      })
      if (response.success) {
        return response
      } else {
        throw new Error(response.error || '安装插件失败')
      }
    } catch (error) {
      console.error('安装插件失败:', error)
      throw error
    }
  }

  // 卸载插件
  const uninstallPlugin = async (pluginName) => {
    try {
      const response = await $fetch(`/api/plugins/${pluginName}`, {
        method: 'DELETE'
      })
      if (response.success) {
        return response
      } else {
        throw new Error(response.error || '卸载插件失败')
      }
    } catch (error) {
      console.error('卸载插件失败:', error)
      throw error
    }
  }

  // 验证插件配置
  const validatePluginConfig = async (pluginName, config) => {
    try {
      // 本地JSON验证
      JSON.parse(JSON.stringify(config))

      // 可以添加服务器端验证
      const response = await $fetch(`/api/plugins/${pluginName}/config/validate`, {
        method: 'POST',
        body: {
          config
        }
      })

      return response.success !== false
    } catch (error) {
      console.error('验证插件配置失败:', error)
      return false
    }
  }

  // 批量操作插件
  const batchOperation = async (pluginNames, operation) => {
    const results = []

    for (const pluginName of pluginNames) {
      try {
        let result
        switch (operation) {
          case 'enable':
            result = await enablePlugin(pluginName)
            break
          case 'disable':
            result = await disablePlugin(pluginName)
            break
          default:
            throw new Error(`不支持的操作: ${operation}`)
        }
        results.push({ pluginName, success: true, result })
      } catch (error) {
        results.push({ pluginName, success: false, error: error.message })
      }
    }

    return results
  }

  // 导出插件配置
  const exportPluginConfig = async (pluginName) => {
    try {
      const plugin = await getPlugin(pluginName)
      const configData = {
        plugin_name: plugin.plugin.name,
        version: plugin.plugin.version,
        config: plugin.plugin.config,
        export_time: new Date().toISOString(),
        exported_by: 'URLDB Plugin Manager'
      }

      const blob = new Blob([JSON.stringify(configData, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${pluginName}_config_${new Date().toISOString().slice(0, 10)}.json`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)

      return true
    } catch (error) {
      console.error('导出插件配置失败:', error)
      throw error
    }
  }

  // 导入插件配置
  const importPluginConfig = async (pluginName, configFile) => {
    try {
      const text = await configFile.text()
      const configData = JSON.parse(text)

      if (configData.plugin_name !== pluginName) {
        throw new Error('配置文件与插件名称不匹配')
      }

      return await updatePluginConfig(pluginName, configData.config)
    } catch (error) {
      console.error('导入插件配置失败:', error)
      throw error
    }
  }

  return {
    getPlugins,
    getPlugin,
    getPluginStats,
    enablePlugin,
    disablePlugin,
    updatePluginConfig,
    getPluginLogs,
    getPluginMarket,
    installPlugin,
    uninstallPlugin,
    validatePluginConfig,
    batchOperation,
    exportPluginConfig,
    importPluginConfig
  }
}