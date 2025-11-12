import { useSeo } from './useSeo'

export const useGlobalSeo = () => {
  const { systemConfig, fetchSystemConfig, setPageSeo, setServerSeo } = useSeo()

  // 初始化系统配置
  const initSystemConfig = async () => {
    if (!systemConfig.value) {
      await fetchSystemConfig()
    }
  }

  // 为首页设置SEO
  const setHomeSeo = (customMeta?: Record<string, string>) => {
    setPageSeo('首页', {
      description: systemConfig.value?.site_description || '老九网盘资源数据库 - 专业的网盘资源管理系统',
      keywords: systemConfig.value?.keywords || '网盘资源,资源管理,数据库,文件分享',
      ...customMeta
    })
  }

  // 为登录页设置SEO
  const setLoginSeo = (customMeta?: Record<string, string>) => {
    setPageSeo('用户登录', {
      description: systemConfig.value?.site_description ? `${systemConfig.value.site_description} - 用户登录页面` : '老九网盘资源数据库登录页面',
      keywords: `${systemConfig.value?.keywords || '网盘资源,登录'},用户登录,账号登录`,
      ...customMeta
    })
  }

  // 为注册页设置SEO
  const setRegisterSeo = (customMeta?: Record<string, string>) => {
    setPageSeo('用户注册', {
      description: systemConfig.value?.site_description ? `${systemConfig.value.site_description} - 用户注册页面` : '老九网盘资源数据库注册页面',
      keywords: `${systemConfig.value?.keywords || '网盘资源,注册'},用户注册,账号注册,免费注册`,
      ...customMeta
    })
  }

  // 为热门剧页面设置SEO
  const setHotDramasSeo = (customMeta?: Record<string, string>) => {
    setPageSeo('热播剧榜单', {
      description: systemConfig.value?.site_description ? `${systemConfig.value.site_description} - 实时获取豆瓣热门电影和电视剧榜单` : '实时获取豆瓣热门电影和电视剧榜单，包括热门电影、热门电视剧、热门综艺和豆瓣Top250等分类',
      keywords: `${systemConfig.value?.keywords || '网盘资源'},热播剧,热门电影,热门电视剧,豆瓣榜单,Top250,影视推荐,电影榜单`,
      ...customMeta
    })
  }

  // 为API文档页面设置SEO
  const setApiDocsSeo = (customMeta?: Record<string, string>) => {
    setPageSeo('API文档', {
      description: systemConfig.value?.site_description ? `${systemConfig.value.site_description} - 公开API接口文档` : '老九网盘资源数据库的公开API接口文档，支持资源添加、搜索和热门剧获取等功能',
      keywords: `${systemConfig.value?.keywords || '网盘资源'},API,接口文档,资源搜索,批量添加,API接口,开发者`,
      ...customMeta
    })
  }

  return {
    initSystemConfig,
    systemConfig,
    setHomeSeo,
    setLoginSeo,
    setRegisterSeo,
    setHotDramasSeo,
    setApiDocsSeo
  }
}