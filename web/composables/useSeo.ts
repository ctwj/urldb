import { ref } from 'vue'
import { useRoute } from '#imports'

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
  const { getPublicSystemConfig } = usePublicSystemConfigApi()

  // 获取系统配置
  const fetchSystemConfig = async () => {
    try {
      const response = await getPublicSystemConfig() as any
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
    if (systemConfig.value && systemConfig.value.site_title) {
      return `${systemConfig.value.site_title} - ${pageTitle}`
    }
    return `${pageTitle} - 老九网盘资源数据库`
  }

  // 生成页面元数据
  const generateMeta = (customMeta?: Record<string, string>) => {
    const defaultMeta = {
      description: (systemConfig.value && systemConfig.value.site_description) || '专业的老九网盘资源数据库',
      keywords: (systemConfig.value && systemConfig.value.keywords) || '网盘,资源管理,文件分享',
      author: (systemConfig.value && systemConfig.value.author) || '系统管理员',
      copyright: (systemConfig.value && systemConfig.value.copyright) || '© 2024 老九网盘资源数据库'
    }

    return {
      ...defaultMeta,
      ...customMeta
    }
  }

  // 生成动态OG图片URL
  const generateOgImageUrl = (title: string, description?: string, theme: string = 'default') => {
    // 获取运行时配置
    const config = useRuntimeConfig()
    const ogApiUrl = config.public.ogApiUrl || '/api/og-image'

    // 构建URL参数
    const params = new URLSearchParams()
    params.set('title', title)

    if (description) {
      // 限制描述长度
      const trimmedDesc = description.length > 200 ? description.substring(0, 200) + '...' : description
      params.set('description', trimmedDesc)
    }

    params.set('site_name', (systemConfig.value && systemConfig.value.site_title) || '老九网盘资源数据库')
    params.set('theme', theme)
    params.set('width', '1200')
    params.set('height', '630')

    // 如果是相对路径，添加当前域名
    if (ogApiUrl.startsWith('/')) {
      if (process.client) {
        const origin = window.location.origin
        return `${origin}${ogApiUrl}?${params.toString()}`
      }
      // 服务端渲染时使用配置的API基础URL
      const apiBase = config.public.apiBase || 'http://localhost:8080'
      return `${apiBase}${ogApiUrl}?${params.toString()}`
    }

    return `${ogApiUrl}?${params.toString()}`
  }

  // 生成动态SEO元数据
  const generateDynamicSeo = (pageTitle: string, customMeta?: Record<string, string>, routeQuery?: Record<string, any>, useRawTitle: boolean = false) => {
    const title = useRawTitle ? pageTitle : generateTitle(pageTitle)
    const meta = generateMeta(customMeta)
    const route = routeQuery || useRoute()

    // 根据路由参数生成动态描述
    const searchKeyword = route.query?.search as string || ''
    const platformId = route.query?.platform as string || ''

    let dynamicDescription = meta.description
    if (searchKeyword && platformId) {
      dynamicDescription = `在${platformId}中搜索"${searchKeyword}"的相关资源。${meta.description}`
    } else if (searchKeyword) {
      dynamicDescription = `搜索"${searchKeyword}"的相关资源。${meta.description}`
    }

    // 动态关键词
    let dynamicKeywords = meta.keywords
    if (searchKeyword) {
      dynamicKeywords = `${searchKeyword},${meta.keywords}`
    }

    // 生成动态OG图片URL
    const theme = searchKeyword ? 'blue' : platformId ? 'green' : 'default'
    const ogImageUrl = generateOgImageUrl(title, dynamicDescription, theme)

    return {
      title,
      description: dynamicDescription,
      keywords: dynamicKeywords,
      ogTitle: title,
      ogDescription: dynamicDescription,
      ogType: 'website',
      ogImage: ogImageUrl,
      ogSiteName: (systemConfig.value && systemConfig.value.site_title) || '老九网盘资源数据库',
      twitterCard: 'summary_large_image',
      robots: 'index, follow'
    }
  }

  // 设置页面SEO - 使用Nuxt3最佳实践
  const setPageSeo = (pageTitle: string, customMeta?: Record<string, string>, routeQuery?: Record<string, any>) => {
    // 检测标题是否已包含站点名（以避免重复）
    const isTitleFormatted = systemConfig.value && pageTitle.includes(systemConfig.value.site_title || '');
    const seoData = generateDynamicSeo(pageTitle, customMeta, routeQuery, isTitleFormatted)

    useSeoMeta({
      title: seoData.title,
      description: seoData.description,
      keywords: seoData.keywords,
      ogTitle: seoData.ogTitle,
      ogDescription: seoData.ogDescription,
      ogType: seoData.ogType,
      ogImage: seoData.ogImage,
      ogSiteName: seoData.ogSiteName,
      twitterCard: seoData.twitterCard,
      robots: seoData.robots
    })
  }

  // 设置服务端SEO（适用于不需要在客户端更新的元数据）
  const setServerSeo = (pageTitle: string, customMeta?: Record<string, string>) => {
    if (import.meta.server) {
      const title = generateTitle(pageTitle)
      const meta = generateMeta(customMeta)

      useServerSeoMeta({
        title: title,
        description: meta.description,
        keywords: meta.keywords,
        robots: 'index, follow'
      })
    }
  }

  return {
    systemConfig,
    fetchSystemConfig,
    generateTitle,
    generateMeta,
    generateOgImageUrl,
    generateDynamicSeo,
    setPageSeo,
    setServerSeo
  }
} 