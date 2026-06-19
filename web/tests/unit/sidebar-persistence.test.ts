import { describe, it, expect } from 'vitest'
import {
  serializeSidebarState,
  deserializeSidebarState,
  expandForGroup,
  type SidebarState,
} from '~/utils/sidebarState'

describe('serializeSidebarState', () => {
  it('正常状态序列化为 JSON 字符串', () => {
    const state: SidebarState = { dataManagement: true, systemConfig: false }
    const raw = serializeSidebarState(state)
    expect(JSON.parse(raw)).toEqual(state)
  })

  it('空对象序列化为 {}', () => {
    expect(serializeSidebarState({})).toBe('{}')
  })

  it('剥离未知键（仅保留 4 个合法分组）', () => {
    // dashboard 不参与持久化；非法键应被过滤
    const raw = serializeSidebarState({
      dataManagement: true,
      dashboard: true,
      invalidGroup: true,
    } as any)
    const parsed = JSON.parse(raw)
    expect(parsed).toEqual({ dataManagement: true })
    expect(parsed.dashboard).toBeUndefined()
    expect(parsed.invalidGroup).toBeUndefined()
  })

  it('所有 false 值仍保留（用户显式折叠）', () => {
    const state: SidebarState = { dataManagement: false, systemConfig: false }
    const parsed = JSON.parse(serializeSidebarState(state))
    expect(parsed).toEqual(state)
  })
})

describe('deserializeSidebarState', () => {
  it('正常 JSON 反序列化', () => {
    const state = deserializeSidebarState('{"dataManagement":true,"systemConfig":false}')
    expect(state).toEqual({ dataManagement: true, systemConfig: false })
  })

  it('空字符串返回空对象', () => {
    expect(deserializeSidebarState('')).toEqual({})
  })

  it('null 返回空对象', () => {
    expect(deserializeSidebarState(null)).toEqual({})
  })

  it('undefined 返回空对象', () => {
    expect(deserializeSidebarState(undefined as any)).toEqual({})
  })

  it('损坏 JSON 返回空对象（不抛异常）', () => {
    expect(deserializeSidebarState('{not json')).toEqual({})
    expect(deserializeSidebarState('}}}}')).toEqual({})
  })

  it('非对象 JSON（数组/字符串/数字）返回空对象', () => {
    expect(deserializeSidebarState('[1,2,3]')).toEqual({})
    expect(deserializeSidebarState('"hello"')).toEqual({})
    expect(deserializeSidebarState('123')).toEqual({})
  })

  it('剥离非法键与 dashboard', () => {
    const raw = JSON.stringify({
      dataManagement: true,
      dashboard: true,
      invalidGroup: true,
    })
    const state = deserializeSidebarState(raw)
    expect(state).toEqual({ dataManagement: true })
  })

  it('非布尔值被强制转换为布尔', () => {
    const raw = JSON.stringify({ dataManagement: 1, systemConfig: 0 })
    const state = deserializeSidebarState(raw)
    expect(state.dataManagement).toBe(true)
    expect(state.systemConfig).toBe(false)
  })
})

describe('expandForGroup', () => {
  it('指定分组时该分组在结果中为 true', () => {
    const result = expandForGroup({}, 'dataManagement')
    expect(result.dataManagement).toBe(true)
  })

  it('保留其他已展开的分组', () => {
    const prev: SidebarState = { dataManagement: true, systemConfig: true }
    const result = expandForGroup(prev, 'statistics')
    expect(result.dataManagement).toBe(true)
    expect(result.systemConfig).toBe(true)
    expect(result.statistics).toBe(true)
  })

  it('dashboard 不参与持久化（传入 dashboard 时返回原状态）', () => {
    const prev: SidebarState = { dataManagement: true }
    // dashboard 不会作为合法 group key，函数应直接返回 prev（无副作用）
    const result = expandForGroup(prev, null)
    expect(result).toEqual(prev)
    expect(result.dashboard).toBeUndefined()
  })

  it('不修改传入的原对象', () => {
    const prev: SidebarState = { dataManagement: true }
    const result = expandForGroup(prev, 'systemConfig')
    // 原对象不应被修改
    expect(prev).toEqual({ dataManagement: true })
    expect(result).not.toBe(prev)
    expect(result.systemConfig).toBe(true)
  })
})
