export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const query = getQuery(event)
  
  try {
    // 在服务端调用后端 API
    const response = await $fetch('/resources', {
      baseURL: config.public.apiBase,
      query,
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    return response
  } catch (error) {
    console.error('服务端获取资源失败:', error)
    throw createError({
      statusCode: 500,
      statusMessage: '获取资源失败'
    })
  }
}) 