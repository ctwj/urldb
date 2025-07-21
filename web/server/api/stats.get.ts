export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  
  try {
    // 在服务端调用后端 API
    const response = await $fetch('/stats', {
      baseURL: String(process.server ? config.public.apiServer : config.public.apiBase),
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    return response
  } catch (error) {
    console.error('服务端获取统计数据失败:', error)
    throw createError({
      statusCode: 500,
      statusMessage: '获取统计数据失败'
    })
  }
}) 