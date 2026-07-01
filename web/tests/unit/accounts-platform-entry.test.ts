import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

/**
 * accounts.vue 平台账号入口契约测试
 *
 * 背景：feature 005-baidu-pan-account 修复"新增平台账号时无法选择百度网盘"。
 * 根因是 web/pages/admin/accounts.vue 曾有一处硬编码白名单 panEnables=['quark','xunlei']
 * 把 baidu 等平台过滤掉了。
 *
 * 本测试用源码契约（source-inspection）固化架构规则——不挂载 Nuxt 页面
 * （项目约定 web/tests/ 只覆盖纯工具与展示型组件；pages 由手动 smoke 覆盖）。
 * 这些断言会在有人重新引入硬编码白名单时立即失败，从而保护 FR-001 与 FR-008。
 *
 * 关联：
 * - spec.md FR-001（平台选择列表必须包含百度网盘）
 * - spec.md FR-008（MUST NOT 引入独立启用开关；展示行为与其他平台一致）
 * - spec.md FR-002（选中百度网盘后呈现 Cookie 字段）
 */

const ACCOUNTS_VUE = resolve(__dirname, '../../pages/admin/accounts.vue')
const source = () => readFileSync(ACCOUNTS_VUE, 'utf-8')

describe('FR-001 + FR-008: 平台选择列表不得有硬编码白名单', () => {
  it('不得定义 panEnables 白名单常量', () => {
    // panEnables 是历史根因；任何形式的重引入都直接违反 FR-008
    expect(source()).not.toMatch(/panEnables/)
  })

  it('不得用本地 Enabled 数组过滤 /api/pans 返回的平台列表', () => {
    // 防止换个变量名（如 enabledPlatforms / allowedPans）重引入同一 bug
    expect(source()).not.toMatch(/\.filter\([^)]*[Ee]nable[^)]*\)/)
    expect(source()).not.toMatch(/\.filter\([^)]*[Aa]llow[^)]*\)/)
  })

  it('新增账号表单的平台下拉必须直接渲染 /api/pans 全量平台', () => {
    // 关键修复点：modal 内 n-select 的 :options 由 platforms.map(...) 而非 .filter(...).map(...) 生成
    const src = source()
    // 找到表单内 v-model 绑定 form.pan_id 的下拉（Naive UI 使用 v-model:value）
    expect(src).toMatch(/v-model:value="form\.pan_id"/)
    // 该下拉的 :options 必须是 platforms.map(...)（不含 filter 前置）
    const modalSelectMatch = src.match(/<n-select[^>]*v-model:value="form\.pan_id"[^>]*:options="([^"]+)"/)
    expect(modalSelectMatch).not.toBeNull()
    const optionsExpr = modalSelectMatch![1]
    expect(optionsExpr).not.toMatch(/\.filter\(/)
    expect(optionsExpr).toMatch(/platforms\.map/)
  })
})

describe('FR-002: 选中百度网盘后呈现 Cookie 字段', () => {
  it('必须存在 isBaidu 响应式变量', () => {
    expect(source()).toMatch(/const\s+isBaidu\s*=\s*ref\(false\)/)
  })

  it('watch(form.pan_id) 必须在 pan.name === "baidu" 时设置 isBaidu=true', () => {
    const src = source()
    // watch 必须存在并切换 isBaidu
    expect(src).toMatch(/watch\(\(\)\s*=>\s*form\.value\.pan_id/)
    expect(src).toMatch(/pan\.name\s*===?\s*['"]baidu['"]/)
    expect(src).toMatch(/isBaidu\.value\s*=\s*true/)
  })

  it('Cookie textarea 的 v-if 必须同时覆盖 quark 和 baidu', () => {
    // FR-002：百度网盘的 Cookie 字段语义与夸克一致
    const src = source()
    // 找到 Cookie textarea 所在 div 的 v-if 表达式
    const cookieBlockMatch = src.match(/<div\s+v-if="([^"]+)"[^>]*>\s*<label[^>]*>\s*Cookie/)
    expect(cookieBlockMatch).not.toBeNull()
    const vifExpr = cookieBlockMatch![1]
    // 必须同时引用 isQuark 和 isBaidu（顺序不拘，允许 || 或 OR 风格）
    expect(vifExpr).toMatch(/isQuark/)
    expect(vifExpr).toMatch(/isBaidu/)
    // 必须是"或"逻辑，而非"与"
    expect(vifExpr).toMatch(/\|\||\bOR\b/)
  })
})

describe('阿里云盘: 选中后呈现 refresh_token 字段（feature 008）', () => {
  // 背景：阿里云盘用 refresh_token 授权（非 Cookie），曾因 accounts.vue 的 watch/表单
  // 缺少 alipan 分支，导致选中阿里云盘后无任何输入项。此处固化修复。

  it('必须存在 isAlipan 响应式变量', () => {
    expect(source()).toMatch(/const\s+isAlipan\s*=\s*ref\(false\)/)
  })

  it('watch(form.pan_id) 必须在 pan.name 为 alipan/aliyun 时设置 isAlipan=true', () => {
    const src = source()
    expect(src).toMatch(/pan\.name\s*===?\s*['"]alipan['"]/)
    expect(src).toMatch(/isAlipan\.value\s*=\s*true/)
  })

  it('必须存在阿里云盘 refresh_token 输入区（v-if="isAlipan"）', () => {
    const src = source()
    expect(src).toMatch(/<div\s+v-if="isAlipan"/)
    expect(src).toMatch(/Refresh Token/)
  })
})
