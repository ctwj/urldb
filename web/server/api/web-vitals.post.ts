/**
 * Web Vitals 上报端点
 *
 * 接收前端 sendBeacon 上报的 LCP/FID/CLS/TTFB/INP 数据。
 * 当前阶段仅日志输出（开发可见），后续可接入 PostgreSQL 持久化做趋势分析。
 *
 * 设计要点：
 * - POST + sendBeacon：浏览器在页面卸载时也能送达
 * - 200 快速响应：不阻塞业务
 * - 不要求鉴权（性能指标属匿名遥测；如需关联用户后续再扩展）
 */
export default defineEventHandler(async (event) => {
  try {
    const body = await readBody(event)
    if (!body || !Array.isArray(body.metrics)) {
      setResponseStatus(event, 400)
      return { code: 400, message: 'Invalid payload', data: null }
    }

    // 开发阶段：仅控制台输出便于调试
    // 生产阶段：此处可改为 INSERT 到 PostgreSQL web_vitals 表
    if (process.env.NODE_ENV !== 'production') {
      // eslint-disable-next-line no-console
      console.log('[web-vitals]', body.metrics)
    }

    return { code: 0, message: 'ok', data: { received: body.metrics.length } }
  } catch (err) {
    setResponseStatus(event, 500)
    return {
      code: 500,
      message: 'Internal error',
      data: null,
    }
  }
})
