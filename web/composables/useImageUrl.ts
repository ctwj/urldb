export const useImageUrl = () => {
  const getImageUrl = (url: string) => {
    if (!url) return ''
    
    // 如果已经是完整URL，直接返回
    if (url.startsWith('http://') || url.startsWith('https://')) {
      return url
    }
    
    // 如果是相对路径，在开发环境中添加后端地址
    if (process.env.NODE_ENV === 'development') {
      const fullUrl = `http://localhost:8080${url}`
      // console.log('useImageUrl - 开发环境:', { original: url, processed: fullUrl })
      return fullUrl
    }
    
    // 生产环境中直接返回相对路径（通过Nginx代理）
    // console.log('useImageUrl - 生产环境:', { original: url, processed: url })
    return url
  }
  
  return {
    getImageUrl
  }
} 