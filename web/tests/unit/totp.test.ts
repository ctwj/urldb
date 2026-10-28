import { describe, expect, it } from 'vitest'
import { generateTotpCode, totpRemainingSeconds, TOTP_PERIOD_SECONDS } from '~/utils/totp'

// 测试向量来自《微信公众号动态口令对接说明.md》第四节，与后端 utils/totp_test.go 同源
describe('generateTotpCode', () => {
  const secret = 'moss-wxpass-8f4c2a1e9b7d'

  it('与对接文档测试向量一致', () => {
    expect(generateTotpCode(secret, 1700000000)).toBe('534579')
    expect(generateTotpCode(secret, 1700001800)).toBe('079502')
    expect(generateTotpCode(secret, 1735689600)).toBe('428950')
  })

  it('时间窗内口令一致，跨窗后变化', () => {
    // 1700000000 所在窗口为 [1699999200, 1700000999]，取窗内另一时刻对比
    expect(generateTotpCode(secret, 1700000000)).toBe(generateTotpCode(secret, 1700000800))
    expect(generateTotpCode(secret, 1700000000)).not.toBe(generateTotpCode(secret, 1700001800))
  })

  it('不足 6 位时补前导 0', () => {
    expect(generateTotpCode(secret, 1700001800)).toMatch(/^079502$/)
    expect(generateTotpCode(secret, 1700001800).length).toBe(6)
  })

  it('自定义周期生效', () => {
    // counter = 28333333 的 60 秒窗口为 [1699999980, 1700000039]
    expect(generateTotpCode(secret, 1699999980, 60)).toBe(generateTotpCode(secret, 1699999980 + 59, 60))
    expect(generateTotpCode(secret, 1699999980, 60)).not.toBe(generateTotpCode(secret, 1699999980, 1800))
  })
})

describe('totpRemainingSeconds', () => {
  it('窗口开头返回满周期，结尾返回 1', () => {
    expect(totpRemainingSeconds(1699999200)).toBe(TOTP_PERIOD_SECONDS)
    expect(totpRemainingSeconds(1699999200 + TOTP_PERIOD_SECONDS - 1)).toBe(1)
  })
})
