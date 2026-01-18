/**
 * 提示词管理API工具
 * 提供提示词相关的API调用方法
 */

// 统一响应解析函数
const parseApiResponse = (response) => {
  // 检查是否是包含success字段的响应格式
  if (response && typeof response === 'object' && 'success' in response && 'data' in response) {
    if (response.success) {
      return response.data
    } else {
      throw new Error(response.message || '请求失败')
    }
  }
  return response
}

export const usePromptApi = () => {
  // 获取提示词列表
  const getPrompts = async () => {
    try {
      const response = await useApiFetch('/ai/prompts')
      return parseApiResponse(response)
    } catch (error) {
      console.error('获取提示词列表失败:', error)
      throw error
    }
  }

  // 更新提示词
  const updatePrompt = async (promptId, data) => {
    try {
      const response = await useApiFetch(`/ai/prompts/${promptId}`, {
        method: 'PUT',
        body: data
      })
      return parseApiResponse(response)
    } catch (error) {
      console.error('更新提示词失败:', error)
      throw error
    }
  }

  // 切换提示词状态
  const togglePrompt = async (promptId) => {
    try {
      const response = await useApiFetch(`/ai/prompts/${promptId}/toggle`, {
        method: 'POST'
      })
      return parseApiResponse(response)
    } catch (error) {
      console.error('切换提示词状态失败:', error)
      throw error
    }
  }

  // 测试提示词
  const testPrompt = async (data) => {
    try {
      const response = await useApiFetch('/ai/prompts/test', {
        method: 'POST',
        body: data
      })
      return parseApiResponse(response)
    } catch (error) {
      console.error('测试提示词失败:', error)
      throw error
    }
  }

  // 初始化默认提示词
  const initPrompts = async () => {
    try {
      const response = await useApiFetch('/ai/prompts/init', {
        method: 'POST'
      })
      return parseApiResponse(response)
    } catch (error) {
      console.error('初始化默认提示词失败:', error)
      throw error
    }
  }

  return {
    getPrompts,
    updatePrompt,
    togglePrompt,
    testPrompt,
    initPrompts
  }
}