export default defineNuxtPlugin(() => {
  // 全局错误处理
  const handleApiError = (error: any) => {
    console.error('API错误:', error)
    
    let message = '操作失败'
    
    // 根据错误类型提供不同的提示
    if (error?.message) {
      if (error.message.includes('需要管理员权限')) {
        message = '需要管理员权限，请使用管理员账号登录'
      } else if (error.message.includes('登录已过期')) {
        message = '登录已过期，请重新登录'
      } else if (error.message.includes('网络连接')) {
        message = '网络连接失败，请检查网络后重试'
      } else {
        message = error.message
      }
    }
    
    // 显示错误提示
    try {
      const notification = useNotification()
      notification.error({
        content: message,
        duration: 5000
      })
    } catch (e) {
      // 如果 notification provider 不可用，使用 console.error
      console.error('通知显示失败:', e)
      console.error('错误信息:', message)
    }
  }

  return {
    provide: {
      handleApiError
    }
  }
}) 