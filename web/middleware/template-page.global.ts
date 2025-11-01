export default defineNuxtRouteMiddleware(async (to, from) => {
  // 从cookie获取当前模板
  const cookie = useCookie('site_template', { default: () => 'default' })
  const template = cookie.value || 'default'

  // 构建模板页面路径
  const templatePagePath = `/templates/${template}/pages${to.path}.vue`

  // 检查模板特定页面是否存在
  try {
    // 尝试导入模板特定页面
    await import(`../templates/${template}/pages${to.path}.vue`)

    // 如果存在，可以在这里设置一些路由元数据
    to.meta.template = template
    to.meta.usesTemplatePage = true
  } catch (error) {
    // 模板特定页面不存在，使用默认页面
    to.meta.template = template
    to.meta.usesTemplatePage = false
  }

  return true
})