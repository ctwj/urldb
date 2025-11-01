export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  const { template_id } = body

  // 验证模板ID
  if (!template_id) {
    throw createError({
      statusCode: 400,
      message: '模板ID不能为空'
    })
  }

  // 验证模板是否存在
  const validTemplates = ['default', 'blog']
  if (!validTemplates.includes(template_id)) {
    throw createError({
      statusCode: 400,
      message: '无效的模板ID'
    })
  }

  // 设置cookie
  setCookie(event, 'site_template', template_id, {
    maxAge: 60 * 60 * 24 * 365, // 一年
    httpOnly: false,
    sameSite: true,
    path: '/'
  })

  return {
    success: true,
    message: '模板切换成功',
    data: {
      template_id: template_id
    }
  }
})