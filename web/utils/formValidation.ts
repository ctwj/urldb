/**
 * 表单即时校验纯函数
 *
 * 设计原则：
 * - 纯函数，无副作用，便于单元测试
 * - 短路求值：数组规则遇到首个失败即返回
 * - 空值（null/undefined/''/[]）跳过非必填校验
 * - 自定义 message 优先于默认文案
 *
 * 配置页 FormSection 的字段失焦/输入即时校验复用此函数。
 */

export type FieldType = 'url' | 'email' | 'number' | 'string'

export interface FieldRule {
  /** 字段类型，触发对应格式校验 */
  type?: FieldType
  /** 是否必填 */
  required?: boolean
  /** type=number 时最小值 */
  min?: number
  /** type=number 时最大值 */
  max?: number
  /** 字符串最小长度 */
  minLength?: number
  /** 字符串最大长度 */
  maxLength?: number
  /** 正则匹配 */
  pattern?: RegExp
  /** 自定义错误消息（覆盖默认文案） */
  message?: string
  /** 自定义校验函数：返回 true 视为通过，返回 string 视为错误消息，返回 false 使用 message 或默认文案 */
  validator?: (value: unknown) => boolean | string
}

export type FieldRules = FieldRule | FieldRule[]

export type FormSchema = Record<string, FieldRules>

export type FormErrors = Record<string, string>

/** 判断是否为"空值"（必填校验的依据） */
export function isEmpty(value: unknown): boolean {
  if (value === null || value === undefined) return true
  if (typeof value === 'string') return value.trim() === ''
  if (Array.isArray(value)) return value.length === 0
  return false
}

const URL_PATTERN = /^https?:\/\/[^\s/$.?#].[^\s]*$/i
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

/** 单条规则校验：返回错误消息或 null */
function checkRule(value: unknown, rule: FieldRule): string | null {
  // 1. required
  if (rule.required && isEmpty(value)) {
    return rule.message ?? '此字段为必填'
  }

  // 空值跳过后续校验（除非 required 已触发）
  if (isEmpty(value)) return null

  // 2. type
  if (rule.type === 'url') {
    const v = String(value)
    if (!URL_PATTERN.test(v)) {
      return rule.message ?? '请输入合法的 URL（以 http:// 或 https:// 开头）'
    }
  } else if (rule.type === 'email') {
    const v = String(value)
    if (!EMAIL_PATTERN.test(v)) {
      return rule.message ?? '请输入合法的邮箱地址'
    }
  } else if (rule.type === 'number') {
    const num = Number(value)
    if (!Number.isFinite(num)) {
      return rule.message ?? '请输入有效数字'
    }
    // 3. min / max
    if (rule.min !== undefined && num < rule.min) {
      return rule.message ?? `数值必须大于等于 ${rule.min}`
    }
    if (rule.max !== undefined && num > rule.max) {
      return rule.message ?? `数值必须小于等于 ${rule.max}`
    }
  }

  // 4. minLength / maxLength（字符串）
  if (typeof value === 'string' || typeof value === 'number') {
    const str = String(value)
    if (rule.minLength !== undefined && str.length < rule.minLength) {
      return rule.message ?? `至少 ${rule.minLength} 个字符`
    }
    if (rule.maxLength !== undefined && str.length > rule.maxLength) {
      return rule.message ?? `不能超过 ${rule.maxLength} 个字符`
    }
    // 5. pattern
    if (rule.pattern && !rule.pattern.test(str)) {
      return rule.message ?? '格式不正确'
    }
  }

  // 6. 自定义 validator
  if (rule.validator) {
    const result = rule.validator(value)
    if (result === false) return rule.message ?? '值无效'
    if (typeof result === 'string') return result
  }

  return null
}

/**
 * 校验单个字段
 *
 * @param value 字段当前值
 * @param rules 单条规则或规则数组（数组按顺序短路求值）
 * @returns 错误消息字符串；null 表示通过
 */
export function validateField(value: unknown, rules: FieldRules): string | null {
  const ruleList = Array.isArray(rules) ? rules : [rules]
  for (const rule of ruleList) {
    const err = checkRule(value, rule)
    if (err) return err
  }
  return null
}

/**
 * 校验整个表单
 *
 * @param values 字段值映射
 * @param schema 字段到规则的映射
 * @returns 字段到错误消息的映射（仅包含失败字段）
 */
export function validateForm(values: Record<string, unknown>, schema: FormSchema): FormErrors {
  const errors: FormErrors = {}
  for (const fieldName of Object.keys(schema)) {
    const err = validateField(values[fieldName], schema[fieldName])
    if (err) errors[fieldName] = err
  }
  return errors
}
