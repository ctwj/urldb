/**
 * 仪表盘指标卡片相关纯函数
 *
 * computeComparison: 计算环比昨日对比
 * formatNumber: 数值格式化
 */

export interface StatCardComparison {
  /** 变化绝对值（today - yesterday） */
  delta: number
  /** 变化百分比（yesterday 为 0 时返回 null，前端显示"新增"） */
  percent: number | null
  /** 方向标识 */
  direction: 'up' | 'down' | 'flat' | 'new'
}

/**
 * 计算今日 vs 昨日的环比对比
 *
 * 规则（见 data-model.md StatCardComparison）:
 * - yesterday=0 且 today>0 → direction='new', percent=null
 * - yesterday=0 且 today=0 → direction='flat', percent=null
 * - delta>0 → 'up', delta<0 → 'down', delta=0 → 'flat'
 * - percent 四舍五入到整数
 */
export function computeComparison(today: number, yesterday: number): StatCardComparison {
  const delta = today - yesterday

  if (yesterday === 0) {
    return {
      delta,
      percent: null,
      direction: today > 0 ? 'new' : 'flat',
    }
  }

  const percent = Math.round((Math.abs(delta) / yesterday) * 100)

  let direction: 'up' | 'down' | 'flat'
  if (delta > 0) direction = 'up'
  else if (delta < 0) direction = 'down'
  else direction = 'flat'

  return { delta, percent, direction }
}

/**
 * 数值格式化
 * - integer: 千分位（如 12,345）
 * - compact: ≥10000 时缩写为"万"（如 1.2 万）
 */
export function formatNumber(value: number, format: 'integer' | 'compact' = 'integer'): string {
  if (format === 'compact' && value >= 10000) {
    return (value / 10000).toFixed(1) + ' 万'
  }
  return value.toLocaleString('zh-CN')
}
