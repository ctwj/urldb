export default defineNuxtRouteMiddleware(async (to, from) => {
  // 从cookie获取当前模板
  const cookie = useCookie('site_template', { default: () => 'default' })
  const template = cookie.value || 'default'

  // 将当前模板信息添加到路由元数据
  to.meta.currentTemplate = template

  // 检查模板特定页面是否存在
  try {
    // 这里只是记录日志，实际的路由逻辑在组件中处理
    console.log(`当前模板: ${template}, 路径: ${to.path}`)
  } catch (error) {
    console.warn('检查模板页面时出现错误:', error)
  }

  return true
})