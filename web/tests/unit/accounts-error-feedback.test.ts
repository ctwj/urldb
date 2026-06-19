import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

/**
 * accounts.vue 错误反馈路径契约测试
 *
 * 背景：feature 005-baidu-pan-account US2 要求后端 400 错误的 message 字段
 * （如 "该百度网盘账号已存在，请使用编辑功能更新凭证" 或
 * "无法获取用户信息，账号创建失败: xxx"）能正确传递到前端 UI。
 *
 * 实现链：
 *   cksApi.createCks(form.value)
 *     → useApiFetch('/cks', { method: 'POST' })
 *     → parseApiResponse 在 code !== 0 时 throw Error(message)
 *     → createCks 的 catch (error) 捕获
 *     → dialog.error({ content: '创建账号失败: ' + error.message })
 *
 * 因项目约定不挂载 Nuxt 页面（见 accounts-platform-entry.test.ts 头注释），
 * 本测试采用源码契约方式：读取 accounts.vue 源码，断言错误反馈链的每一环
 * 都存在且语义正确。这能在未来有人误删 try/catch 或改用 console.log 时
 * 立即失败。
 *
 * 关联：
 * - spec.md FR-009（重复账号拒绝）
 * - spec.md FR-010（失效凭证同步验证）
 * - contracts/cks-api.md（ErrorResponse 顶层 message 字段契约）
 */

const ACCOUNTS_VUE = resolve(__dirname, '../../pages/admin/accounts.vue')
const source = () => readFileSync(ACCOUNTS_VUE, 'utf-8')

/**
 * 从 createCks 开始切片到下一个顶层 `const`，得到完整函数体。
 * 因函数内有嵌套块（try/catch/finally/object literal），简单正则无法
 * 平衡花括号；用"下一个 const 声明"作为终止标志更稳。
 */
function createCksBody(): string {
  const src = source()
  const startIdx = src.indexOf('const createCks')
  expect(startIdx).toBeGreaterThan(-1)
  // 找下一个顶层 const（出现在行首、不在对象字面量内）
  const rest = src.slice(startIdx + 1)
  const nextConst = rest.search(/\nconst\s+\w/)
  return nextConst === -1 ? src.slice(startIdx) : src.slice(startIdx, startIdx + 1 + nextConst)
}

describe('FR-009 + FR-010: 后端 400 错误能正确传递到 UI', () => {
  it('createCks 必须通过 useApi 封装调用（禁止裸 fetch）', () => {
    // Principle II — Unified API Contract
    const body = createCksBody()
    expect(body).toMatch(/cksApi\.createCks\(/)
    expect(body).not.toMatch(/fetch\(['"`]\/api\/cks/)
  })

  it('createCks 必须有 try/catch，catch 提取 error.message', () => {
    const body = createCksBody()
    expect(body).toMatch(/\btry\s*\{/)
    expect(body).toMatch(/catch\s*\(\s*(\w+)\s*\)/)
    const errVar = body.match(/catch\s*\(\s*(\w+)\s*\)/)![1]
    expect(body).toMatch(new RegExp(`${errVar}\\.message`))
  })

  it('catch 必须调用 dialog.error 显示失败信息，content 拼接 error.message', () => {
    const body = createCksBody()
    expect(body).toMatch(/dialog\.error\s*\(/)
    expect(body).toMatch(/['"]创建账号失败[^'"]*['"]\s*\+\s*\(?(\w+\.message|\w+\.message\s*\|\|)/)
  })

  it('后端查重错误消息能通过 parseApiResponse 透传（useApi 封装存在）', () => {
    const src = source()
    expect(src).toMatch(/useApi|useCksApi/)
  })
})
