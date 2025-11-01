import { useTemplateManager } from '~/composables/templateManager'

export const useTemplatePageResolver = () => {
  const { currentTemplate } = useTemplateManager()

  /**
   * 解析模板页面路径
   * @param pagePath 页面路径 (例如: '/index', '/hot-dramas')
   * @returns 模板页面的实际路径
   */
  const resolveTemplatePage = (pagePath: string): string => {
    // 标准化页面路径
    const normalizedPath = pagePath.startsWith('/') ? pagePath : `/${pagePath}`

    // 构建模板页面路径
    const templatePagePath = `/templates/${currentTemplate.value}/pages${normalizedPath}.vue`

    return templatePagePath
  }

  /**
   * 检查模板页面是否存在
   * @param pagePath 页面路径
   * @returns Promise<boolean>
   */
  const templatePageExists = async (pagePath: string): Promise<boolean> => {
    try {
      const templatePagePath = resolveTemplatePage(pagePath)
      // 尝试导入组件来检查是否存在
      await import(`~/templates/${currentTemplate.value}/pages${pagePath}.vue`)
      return true
    } catch (error) {
      return false
    }
  }

  /**
   * 获取页面组件
   * @param pagePath 页面路径
   * @returns 页面组件或null
   */
  const getPageComponent = async (pagePath: string) => {
    try {
      // 首先检查模板特定页面是否存在
      const templatePageExists = await checkTemplatePageExists(pagePath)

      if (templatePageExists) {
        // 如果模板特定页面存在，返回它
        return defineAsyncComponent(() => import(`~/templates/${currentTemplate.value}/pages${pagePath}.vue`))
      } else {
        // 否则返回默认页面
        const defaultPagePath = pagePath === '/index' ? '/' : pagePath
        return defineAsyncComponent(() => import(`~/pages${defaultPagePath}.vue`))
      }
    } catch (error) {
      console.error('获取页面组件失败:', error)
      // 返回默认页面作为后备
      const defaultPagePath = pagePath === '/index' ? '/' : pagePath
      return defineAsyncComponent(() => import(`~/pages${defaultPagePath}.vue`))
    }
  }

  /**
   * 检查模板页面是否存在 (内部方法)
   * @param pagePath 页面路径
   * @returns Promise<boolean>
   */
  const checkTemplatePageExists = async (pagePath: string): Promise<boolean> => {
    try {
      await import(`~/templates/${currentTemplate.value}/pages${pagePath}.vue`)
      return true
    } catch (error) {
      return false
    }
  }

  return {
    resolveTemplatePage,
    templatePageExists,
    getPageComponent
  }
}