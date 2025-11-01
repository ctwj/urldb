export default defineEventHandler(async (event) => {
  // 从系统配置获取可用模板列表
  const templates = [
    {
      id: 'default',
      name: '默认模板',
      description: '系统默认模板',
      enabled: true
    },
    {
      id: 'blog',
      name: '博客模板',
      description: '博客风格模板',
      enabled: true
    }
  ]

  return {
    success: true,
    data: templates
  }
})