// 图片预加载和缓存工具

interface ImageCache {
  [key: string]: Promise<string>
}

class ImageLoader {
  private cache: ImageCache = {}

  /**
   * 预加载图片
   */
  async preloadImage(url: string): Promise<string> {
    if (this.cache[url]) {
      return this.cache[url]
    }

    const promise = new Promise<string>((resolve, reject) => {
      const img = new Image()
      img.onload = () => resolve(url)
      img.onerror = () => {
        // 如果加载失败，从缓存中移除
        delete this.cache[url]
        reject(new Error(`Failed to load image: ${url}`))
      }
      img.src = url
    })

    this.cache[url] = promise
    return promise
  }

  /**
   * 批量预加载图片
   */
  async preloadImages(urls: string[]): Promise<void> {
    const promises = urls.map(url => this.preloadImage(url).catch(() => {
      // 忽略单个图片加载失败
      console.warn(`Failed to preload image: ${url}`)
    }))
    await Promise.all(promises)
  }

  /**
   * 获取缓存中的图片
   */
  getCachedImage(url: string): Promise<string> | undefined {
    return this.cache[url]
  }

  /**
   * 清理缓存
   */
  clearCache(): void {
    this.cache = {}
  }

  /**
   * 获取缓存大小
   */
  getCacheSize(): number {
    return Object.keys(this.cache).length
  }
}

// 创建全局图片加载器实例
export const imageLoader = new ImageLoader()

// 预加载常用 Logo 图片
export const preloadCommonLogos = async () => {
  const commonLogoUrls = [
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23000',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%233B82F6',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FFF',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%238B5CF6',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%236B7280',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%2300D4FF',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23059669',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23DC2626',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%231E40AF',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%237ABE4A',
    'https://api.iconify.design/ion:logo-vercel.svg?color=%23000',
    'https://api.iconify.design/ion:logo-vercel.svg?color=%23FFF',
    'https://api.iconify.design/logos:supabase-icon.svg',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%237700ff',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FF6B6B',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23646CFF',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%2342D392',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23252f3f',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23cebe2c',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%232196b0',
    'https://api.iconify.design/ion:qr-code-outline.svg?color=%23000000',
    'https://api.iconify.design/simple-icons:qq.svg?color=%2371cdfc',
    'https://api.iconify.design/simple-icons:wechat.svg?color=%23000000'
  ]

  try {
    await imageLoader.preloadImages(commonLogoUrls)
    console.log(`Preloaded ${imageLoader.getCacheSize()} common logos`)
  } catch (error) {
    console.warn('Failed to preload some logos:', error)
  }
}