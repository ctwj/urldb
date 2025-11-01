export default defineEventHandler(async (event) => {
  const path = getRouterParam(event, 'path')
  const template = getQuery(event).template as string || 'default'

  // 验证模板名称
  const validTemplates = ['default', 'blog']
  if (!validTemplates.includes(template)) {
    throw createError({
      statusCode: 400,
      message: '无效的模板名称'
    })
  }

  try {
    // 尝试获取模板页面
    const pageModule = await import(`~/templates/${template}/pages/${path}.vue`)

    return {
      success: true,
      data: {
        template: template,
        path: path,
        hasTemplatePage: true,
        component: pageModule.default
      }
    }
  } catch (error) {
    // 模板页面不存在，返回默认页面
    try {
      const defaultPageModule = await import(`~/pages/${path}.vue`)

      return {
        success: true,
        data: {
          template: template,
          path: path,
          hasTemplatePage: false,
          component: defaultPageModule.default
        }
      }
    } catch (defaultError) {
      throw createError({
        statusCode: 404,
        message: '页面不存在'
      })
    }
  }
})