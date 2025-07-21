export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const apiBase = String(process.server ? config.public.apiServer : config.public.apiBase)
  
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