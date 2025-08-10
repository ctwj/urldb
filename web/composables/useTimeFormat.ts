// 统一的时间格式化工具函数
export const useTimeFormat = () => {
  // 格式化日期时间（标准格式）
  const formatDateTime = (dateString: string | Date) => {
    if (!dateString) return '-'
    const date = dateString instanceof Date ? dateString : new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  // 格式化日期（仅日期）
  const formatDate = (dateString: string | Date) => {
    if (!dateString) return '-'
    const date = dateString instanceof Date ? dateString : new Date(dateString)
    return date.toLocaleDateString('zh-CN')
  }

  // 格式化时间（仅时间）
  const formatTime = (dateString: string | Date) => {
    if (!dateString) return '-'
    const date = dateString instanceof Date ? dateString : new Date(dateString)
    return date.toLocaleTimeString('zh-CN')
  }

  // 格式化相对时间
  const formatRelativeTime = (dateString: string | Date) => {
    if (!dateString) return '-'
    const date = dateString instanceof Date ? dateString : new Date(dateString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffSec = Math.floor(diffMs / 1000)
    const diffMin = Math.floor(diffSec / 60)
    const diffHour = Math.floor(diffMin / 60)
    const diffDay = Math.floor(diffHour / 24)
    const diffWeek = Math.floor(diffDay / 7)
    const diffMonth = Math.floor(diffDay / 30)
    const diffYear = Math.floor(diffDay / 365)
    
    const isToday = date.toDateString() === now.toDateString()
    
    if (isToday) {
      if (diffMin < 1) {
        return '刚刚'
      } else if (diffHour < 1) {
        return `${diffMin}分钟前`
      } else {
        return `${diffHour}小时前`
      }
    } else if (diffDay < 1) {
      return `${diffHour}小时前`
    } else if (diffDay < 7) {
      return `${diffDay}天前`
    } else if (diffWeek < 4) {
      return `${diffWeek}周前`
    } else if (diffMonth < 12) {
      return `${diffMonth}个月前`
    } else {
      return `${diffYear}年前`
    }
  }

  // 获取当前时间字符串
  const getCurrentTimeString = () => {
    return new Date().toLocaleString('zh-CN')
  }

  // 检查是否为今天
  const isToday = (dateString: string | Date) => {
    if (!dateString) return false
    const date = dateString instanceof Date ? dateString : new Date(dateString)
    const now = new Date()
    return date.toDateString() === now.toDateString()
  }

  return {
    formatDateTime,
    formatDate,
    formatTime,
    formatRelativeTime,
    getCurrentTimeString,
    isToday
  }
}
