export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080/api'
  
  try {
    const response = await $fetch(`${apiBase}/version`)
    return response
  } catch (error: any) {
    throw createError({
      statusCode: error.statusCode || 500,
      statusMessage: error.statusMessage || '获取版本信息失败'
    })
  }
}) 