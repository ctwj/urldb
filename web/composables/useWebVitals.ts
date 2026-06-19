/**
 * Web Vitals 性能埋点
 *
 * 采集 LCP / CLS / TTFB / INP 四项核心指标（SC-001 性能基线），
 * 注：FID 已被 web-vitals v4+ 移除，由 INP（Interaction to Next Paint）取代
 * 通过 navigator.sendBeacon 上报到 `/api/web-vitals` 端点。
 *
 * 设计：
 * - 仅客户端采集（SSR 阶段不执行）
 * - 单次会话最多上报一次（避免重复）
 * - sendBeacon 在页面隐藏/卸载时也能可靠发送
 * - 失败静默，绝不阻塞业务逻辑
 *
 * 后端 `/api/web-vitals`（Nuxt server route）目前仅日志，后续可接入 PG 持久化。
 */
import { onLCP, onCLS, onTTFB, onINP, type Metric } from 'web-vitals'

const ENDPOINT = '/api/web-vitals'
const SESSION_FLAG = 'urldb:vitals:sent'

export interface WebVitalsPayload {
  name: Metric['name']
  value: number
  rating: Metric['rating']
  id: string
  path: string
  ts: number
}

export function useWebVitals() {
  // 仅客户端执行
  if (!import.meta.client) return

  // 单次会话避免重复上报
  if (sessionStorage.getItem(SESSION_FLAG)) return

  const collected: WebVitalsPayload[] = []
  const currentPath = window.location.pathname

  const record = (metric: Metric) => {
    collected.push({
      name: metric.name,
      value: Number(metric.value.toFixed(2)),
      rating: metric.rating,
      id: metric.id,
      path: currentPath,
      ts: Date.now(),
    })
    // 每个指标采集后立即尝试上报（LCP/CLS 等可能在页面卸载时才触发最终值）
    flush()
  }

  const flush = () => {
    if (collected.length === 0) return
    // 使用 sendBeacon 保证页面卸载时也能送达
    const payload = JSON.stringify({ metrics: collected.slice() })
    if (navigator.sendBeacon) {
      try {
        navigator.sendBeacon(ENDPOINT, payload)
        sessionStorage.setItem(SESSION_FLAG, '1')
        collected.length = 0
      } catch {
        // 静默失败
      }
    }
  }

  try {
    onLCP(record)
    onCLS(record)
    onTTFB(record)
    onINP(record)
  } catch {
    // 静默失败：web-vitals 内部异常不应影响业务
  }

  // 页面隐藏时强制 flush（visibilitychange 是最可靠的"用户即将离开"信号）
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'hidden') flush()
  })

  // 兜底：beforeunload 也尝试 flush
  window.addEventListener('pagehide', flush)
}
