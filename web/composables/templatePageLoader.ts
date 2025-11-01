import { useTemplateManager } from '~/composables/templateManager'

interface TemplatePageInfo {
  path: string
  template: string
  exists: boolean
  component: any
}

export const useTemplatePageLoader = () => {
  const { currentTemplate } = useTemplateManager()

  /**
   * 加载模板页面组件
   * @param pagePath 页面路径
   * @returns 页面组件Promise
   */
  const loadTemplatePage = async (pagePath: string) => {
    // 标准化页面路径
    const normalizedPath = pagePath === '/' ? '/index' : pagePath

    try {
      // 首先尝试加载模板特定页面
      const templateComponent = await loadTemplateSpecificPage(normalizedPath)
      if (templateComponent) {
        return templateComponent
      }

      // 如果模板特定页面不存在，加载默认页面
      return await loadDefaultPage(normalizedPath)
    } catch (error) {
      console.error('加载模板页面失败:', error)
      // 最后回退到默认页面
      return await loadDefaultPage(normalizedPath)
    }
  }

  /**
   * 加载模板特定页面
   * @param pagePath 页面路径
   * @returns 页面组件或null
   */
  const loadTemplateSpecificPage = async (pagePath: string) => {
    try {
      const component = await import(`~/templates/${currentTemplate.value}/pages${pagePath}.vue`)
      return component.default
    } catch (error) {
      // 模板特定页面不存在
      return null
    }
  }

  /**
   * 加载默认页面
   * @param pagePath 页面路径
   * @returns 页面组件
   */
  const loadDefaultPage = async (pagePath: string) => {
    try {
      const component = await import(`~/pages${pagePath}.vue`)
      return component.default
    } catch (error) {
      throw new Error(`无法加载页面: ${pagePath}`)
    }
  }

  /**
   * 检查模板页面是否存在
   * @param pagePath 页面路径
   * @returns Promise<boolean>
   */
  const templatePageExists = async (pagePath: string): Promise<boolean> => {
    try {
      await import(`~/templates/${currentTemplate.value}/pages${pagePath}.vue`)
      return true
    } catch (error) {
      return false
    }
  }

  /**
   * 获取当前模板的页面信息
   * @param pagePath 页面路径
   * @returns TemplatePageInfo
   */
  const getTemplatePageInfo = async (pagePath: string): Promise<TemplatePageInfo> => {
    const normalizedPath = pagePath === '/' ? '/index' : pagePath
    const exists = await templatePageExists(normalizedPath)

    return {
      path: normalizedPath,
      template: currentTemplate.value,
      exists,
      component: exists ? await loadTemplateSpecificPage(normalizedPath) : await loadDefaultPage(normalizedPath)
    }
  }

  return {
    loadTemplatePage,
    templatePageExists,
    getTemplatePageInfo
  }
}