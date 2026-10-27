// 来源渠道常量（与后端 db/entity/source_constants.go 镜像）
// Feature: 009-statistics-enhancement
export const SourceWeb = 'web'
export const SourceWechat = 'wechat'
export const SourceTelegram = 'telegram'
export const SourceApi = 'api'

// 返回来源渠道的中文展示名；未知来源原样返回。
export const SourceDisplayName = (source: string): string => {
  switch (source) {
    case SourceWeb:
      return '网页'
    case SourceWechat:
      return '公众号'
    case SourceTelegram:
      return '电报'
    case SourceApi:
      return '用户API'
    default:
      return source
  }
}
