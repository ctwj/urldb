import { describe, it, expect } from 'vitest'
import { computeComparison, formatNumber } from '~/utils/statCard'

describe('computeComparison', () => {
  describe('yesterday 为 0', () => {
    it('today > 0 时返回 new 方向、null 百分比', () => {
      const result = computeComparison(5, 0)
      expect(result.direction).toBe('new')
      expect(result.percent).toBeNull()
      expect(result.delta).toBe(5)
    })

    it('today = 0 时返回 flat 方向、null 百分比', () => {
      const result = computeComparison(0, 0)
      expect(result.direction).toBe('flat')
      expect(result.percent).toBeNull()
      expect(result.delta).toBe(0)
    })
  })

  describe('上升', () => {
    it('today > yesterday 时返回 up 方向', () => {
      const result = computeComparison(15, 10)
      expect(result.direction).toBe('up')
      expect(result.delta).toBe(5)
      expect(result.percent).toBe(50)
    })

    it('百分比四舍五入到整数', () => {
      // 1/3 ≈ 33.33% → 33
      expect(computeComparison(4, 3).percent).toBe(33)
      // 2/3 ≈ 66.67% → 67
      expect(computeComparison(5, 3).percent).toBe(67)
    })
  })

  describe('下降', () => {
    it('today < yesterday 时返回 down 方向', () => {
      const result = computeComparison(5, 10)
      expect(result.direction).toBe('down')
      expect(result.delta).toBe(-5)
      expect(result.percent).toBe(50)
    })
  })

  describe('持平', () => {
    it('today = yesterday（非 0）时返回 flat、percent=0', () => {
      const result = computeComparison(10, 10)
      expect(result.direction).toBe('flat')
      expect(result.delta).toBe(0)
      expect(result.percent).toBe(0)
    })
  })

  describe('边界值', () => {
    it('100% 增长', () => {
      const result = computeComparison(20, 10)
      expect(result.percent).toBe(100)
      expect(result.direction).toBe('up')
    })

    it('100% 下降', () => {
      const result = computeComparison(0, 10)
      expect(result.percent).toBe(100)
      expect(result.direction).toBe('down')
    })
  })
})

describe('formatNumber', () => {
  it('integer 格式使用千分位', () => {
    expect(formatNumber(12345, 'integer')).toBe('12,345')
    expect(formatNumber(0, 'integer')).toBe('0')
  })

  it('compact 格式 < 10000 时保持千分位', () => {
    expect(formatNumber(9999, 'compact')).toBe('9,999')
  })

  it('compact 格式 ≥ 10000 时缩写为万', () => {
    expect(formatNumber(12345, 'compact')).toBe('1.2 万')
    expect(formatNumber(10000, 'compact')).toBe('1.0 万')
  })

  it('默认 format 为 integer', () => {
    expect(formatNumber(1234)).toBe('1,234')
  })
})
