export interface TemplateConfig {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
}

export const useTemplateManager = () => {
  const currentTemplate = ref<string>('default')
  const templates = ref<Record<string, TemplateConfig>>({
    'default': {
      id: 'default',
      name: '默认模板',
      description: '系统默认模板',
      enabled: true
    },
    'blog': {
      id: 'blog',
      name: '博客模板',
      description: '博客风格模板',
      enabled: true
    }
  })

  // 获取当前模板配置
  const getCurrentTemplate = (): string => {
    // 首先检查查询参数
    const route = useRoute()
    if (route.query.template && templates.value[route.query.template as string]) {
      return route.query.template as string
    }

    // 然后检查cookie
    const cookie = useCookie('site_template', { default: () => 'default' })
    if (cookie.value && templates.value[cookie.value]) {
      return cookie.value
    }

    // 检查localStorage (客户端)
    if (process.client) {
      const storedTemplate = localStorage.getItem('site_template')
      if (storedTemplate && templates.value[storedTemplate]) {
        return storedTemplate
      }
    }

    // 返回默认模板
    return 'default'
  }

  // 设置当前模板
  const setCurrentTemplate = (templateId: string) => {
    if (templates.value[templateId]) {
      currentTemplate.value = templateId

      // 更新cookie
      const cookie = useCookie('site_template', {
        default: () => 'default',
        maxAge: 60 * 60 * 24 * 365, // 一年
        sameSite: true,
        path: '/'
      })
      cookie.value = templateId

      // 更新localStorage (客户端)
      if (process.client) {
        localStorage.setItem('site_template', templateId)
      }
    }
  }

  // 注册模板
  const registerTemplate = (templateId: string, config: TemplateConfig) => {
    templates.value[templateId] = config
  }

  // 获取可用模板列表
  const getAvailableTemplates = (): Record<string, TemplateConfig> => {
    return templates.value
  }

  // 获取模板配置
  const getTemplateConfig = (templateId: string): TemplateConfig | null => {
    return templates.value[templateId] || null
  }

  // 从API获取模板列表
  const fetchTemplates = async () => {
    try {
      const response = await $fetch('/api/templates')
      if (response.success && response.data) {
        const newTemplates: Record<string, TemplateConfig> = {}
        response.data.forEach((template: any) => {
          newTemplates[template.id] = {
            id: template.id,
            name: template.name,
            description: template.description,
            enabled: template.enabled
          }
        })
        templates.value = newTemplates
      }
    } catch (error) {
      console.error('获取模板列表失败:', error)
    }
  }

  // 初始化
  const init = () => {
    currentTemplate.value = getCurrentTemplate()
    // 在客户端获取最新的模板列表
    if (process.client) {
      fetchTemplates()
    }
  }

  return {
    currentTemplate: readonly(currentTemplate),
    setCurrentTemplate,
    registerTemplate,
    getAvailableTemplates,
    getCurrentTemplate,
    getTemplateConfig,
    fetchTemplates,
    init
  }
}