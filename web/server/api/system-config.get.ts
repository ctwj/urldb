export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  
  try {
    // 在服务端调用后端 API
    const response = await $fetch('/system/config', {
      baseURL: String(process.server ? config.public.apiServer : config.public.apiBase),
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    return response
  } catch (error) {
    console.error('服务端获取系统配置失败:', error)
    // 返回默认配置而不是抛出错误
    return {
      site_title: '老九网盘资源数据库',
      site_description: '一个现代化的资源管理系统',
      keywords: '网盘资源,资源管理,数据库',
      author: '老九',
      copyright: '© 2025 老九网盘资源数据库'
    }
  }
}) 