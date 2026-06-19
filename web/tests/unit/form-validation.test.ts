/**
 * 表单即时校验纯函数测试（T033）
 *
 * 覆盖 validateField 核心契约：
 * - 必填校验
 * - 类型校验（url/email/number/string）
 * - 边界校验（min/max for number, minLength/maxLength for string）
 * - 正则 pattern
 * - 自定义 validator（返回 true/string）
 * - 多规则短路求值（第一失败即返回）
 * - 自定义 message 覆盖默认文案
 *
 * 配置页（site-config/feature-config/dev-config/seo）的 FormSection 即时校验复用此函数。
 */
import { describe, it, expect } from 'vitest'
import { validateField, validateForm } from '~/utils/formValidation'

describe('validateField 必填校验', () => {
  it('空字符串触发必填错误', () => {
    expect(validateField('', { required: true })).toBe('此字段为必填')
    expect(validateField('   ', { required: true })).toBe('此字段为必填')
  })

  it('null / undefined / 空数组 触发必填错误', () => {
    expect(validateField(null, { required: true })).toBe('此字段为必填')
    expect(validateField(undefined, { required: true })).toBe('此字段为必填')
    expect(validateField([], { required: true })).toBe('此字段为必填')
  })

  it('有值时必填通过', () => {
    expect(validateField('hello', { required: true })).toBeNull()
    expect(validateField(0, { required: true })).toBeNull()
    expect(validateField(false, { required: true })).toBeNull()
  })

  it('非必填且为空时通过', () => {
    expect(validateField('', { required: false })).toBeNull()
    expect(validateField(null, {})).toBeNull()
  })

  it('支持自定义必填消息', () => {
    expect(validateField('', { required: true, message: '站点名称不能为空' }))
      .toBe('站点名称不能为空')
  })
})

describe('validateField 类型校验', () => {
  it('type=url 合法 URL 通过', () => {
    expect(validateField('https://example.com', { type: 'url' })).toBeNull()
    expect(validateField('http://localhost:3000/path?q=1', { type: 'url' })).toBeNull()
  })

  it('type=url 非法 URL 返回错误', () => {
    expect(validateField('not-a-url', { type: 'url' })).toContain('URL')
    expect(validateField('ftp://example.com', { type: 'url' })).toContain('URL')
  })

  it('type=url 空值跳过类型校验（必填由 required 控制）', () => {
    expect(validateField('', { type: 'url' })).toBeNull()
  })

  it('type=email 合法邮箱通过', () => {
    expect(validateField('user@example.com', { type: 'email' })).toBeNull()
  })

  it('type=email 非法邮箱返回错误', () => {
    expect(validateField('not-an-email', { type: 'email' })).toContain('邮箱')
    expect(validateField('a@', { type: 'email' })).toContain('邮箱')
  })

  it('type=number 合法数字通过', () => {
    expect(validateField(42, { type: 'number' })).toBeNull()
    expect(validateField('42', { type: 'number' })).toBeNull()
    expect(validateField(0, { type: 'number' })).toBeNull()
  })

  it('type=number 非数字返回错误', () => {
    expect(validateField('abc', { type: 'number' })).toContain('数字')
    expect(validateField(NaN, { type: 'number' })).toContain('数字')
  })
})

describe('validateField 边界校验', () => {
  it('number min 边界', () => {
    expect(validateField(5, { type: 'number', min: 5 })).toBeNull()
    expect(validateField(4, { type: 'number', min: 5 })).toContain('大于等于')
  })

  it('number max 边界', () => {
    expect(validateField(5, { type: 'number', max: 5 })).toBeNull()
    expect(validateField(6, { type: 'number', max: 5 })).toContain('小于等于')
  })

  it('string minLength 边界', () => {
    expect(validateField('abc', { minLength: 3 })).toBeNull()
    expect(validateField('ab', { minLength: 3 })).toContain('至少')
  })

  it('string maxLength 边界', () => {
    expect(validateField('abc', { maxLength: 3 })).toBeNull()
    expect(validateField('abcd', { maxLength: 3 })).toContain('不能超过')
  })
})

describe('validateField 正则与自定义', () => {
  it('pattern 匹配通过', () => {
    expect(validateField('abc123', { pattern: /^[a-z0-9]+$/ })).toBeNull()
  })

  it('pattern 不匹配返回错误', () => {
    expect(validateField('ABC!', { pattern: /^[a-z0-9]+$/ })).toContain('格式')
  })

  it('validator 返回 true 视为通过', () => {
    expect(validateField('hello', { validator: (v) => v === 'hello' })).toBeNull()
  })

  it('validator 返回 false 使用默认错误', () => {
    expect(validateField('world', { validator: (v) => v === 'hello' })).toContain('无效')
  })

  it('validator 返回字符串作为自定义错误', () => {
    expect(validateField('x', { validator: (v) => v.length >= 3 ? true : '长度不足' }))
      .toBe('长度不足')
  })
})

describe('validateField 多规则短路', () => {
  it('数组规则：第一个失败即返回', () => {
    const result = validateField('', [
      { required: true },
      { type: 'url' },
    ])
    expect(result).toBe('此字段为必填')
  })

  it('数组规则：前序通过后序继续', () => {
    const result = validateField('not-a-url', [
      { required: true },
      { type: 'url' },
    ])
    expect(result).toContain('URL')
  })

  it('数组规则：全部通过返回 null', () => {
    const result = validateField('https://example.com', [
      { required: true },
      { type: 'url' },
      { pattern: /^https:/, message: '必须 HTTPS' },
    ])
    expect(result).toBeNull()
  })
})

describe('validateForm 表单级聚合', () => {
  it('所有字段通过返回空对象', () => {
    const schema = {
      name: [{ required: true }],
      url: [{ required: true }, { type: 'url' }],
    }
    const values = { name: '站点', url: 'https://example.com' }
    expect(validateForm(values, schema)).toEqual({})
  })

  it('返回字段名到错误消息的映射', () => {
    const schema = {
      name: [{ required: true, message: '名称必填' }],
      url: [{ type: 'url' }],
    }
    const values = { name: '', url: 'bad' }
    const errors = validateForm(values, schema)
    expect(errors.name).toBe('名称必填')
    expect(errors.url).toContain('URL')
  })

  it('未在 schema 中的字段被忽略', () => {
    const errors = validateForm(
      { extra: 'whatever' },
      { name: [{ required: true }] },
    )
    expect(errors.name).toBe('此字段为必填')
    expect(errors.extra).toBeUndefined()
  })

  it('schema 中字段存在但值为 undefined', () => {
    const errors = validateForm(
      { name: undefined },
      { name: [{ required: true }] },
    )
    expect(errors.name).toBe('此字段为必填')
  })

  it('空 schema 返回空对象', () => {
    expect(validateForm({ name: 'x' }, {})).toEqual({})
  })
})
