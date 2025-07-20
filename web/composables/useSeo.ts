import { ref, computed } from 'vue'

interface SystemConfig {
  id: number
  site_title: string
  site_description: string
  keywords: string
  author: string
  copyright: string
  auto_process_ready_resources: boolean
  auto_process_interval: number
  page_size: number
  maintenance_mode: boolean
  created_at: string
  updated_at: string
}

export const useSeo = () => {
  const systemConfig = ref<SystemConfig | null>(null)
  const { getSystemConfig } = useSystemConfigApi()

  // 获取系统配置
  const fetchSystemConfig = async () => {
    try {
      const response = await getSystemConfig() as any
      console.log('系统配置响应:', response)
      if (response && response.success && response.data) {
        systemConfig.value = response.data
      } else if (response && response.data) {
        // 兼容非标准格式
        systemConfig.value = response.data
      }
    } catch (error) {
      console.error('获取系统配置失败:', error)
    }
  }

  // 生成页面标题
  const generateTitle = (pageTitle: string) => {
    if (systemConfig.value?.site_title) {
      return `${systemConfig.value.site_title} - ${pageTitle}`
    }
    return `${pageTitle} - 老九网盘资源数据库`
  }

  // 生成页面元数据
  const generateMeta = (customMeta?: Record<string, string>) => {
    const defaultMeta = {
      description: systemConfig.value?.site_description || '专业的老九网盘资源数据库',
      keywords: systemConfig.value?.keywords || '网盘,资源管理,文件分享',
      author: systemConfig.value?.author || '系统管理员',
      copyright: systemConfig.value?.copyright || '© 2024 老九网盘资源数据库'
    }

    return {
      ...defaultMeta,
      ...customMeta
    }
  }

  // 设置页面SEO
  const setPageSeo = (pageTitle: string, customMeta?: Record<string, string>) => {
    const title = generateTitle(pageTitle)
    const meta = generateMeta(customMeta)

    useHead({
      title,
      meta: [
        { name: 'description', content: meta.description },
        { name: 'keywords', content: meta.keywords },
        { name: 'author', content: meta.author },
        { name: 'copyright', content: meta.copyright }
      ]
    })
  }

  return {
    systemConfig,
    fetchSystemConfig,
    generateTitle,
    generateMeta,
    setPageSeo
  }
} 