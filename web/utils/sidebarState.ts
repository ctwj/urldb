/**
 * 侧边栏分组展开状态持久化纯函数
 *
 * 设计要点：
 *   - 仅 4 个可折叠分组参与持久化：dataManagement / systemConfig / operation / statistics
 *   - dashboard 永远单独展示，不参与持久化
 *   - 进入路由时通过 expandForGroup 自动展开当前分组（覆盖存储值，但不抹除其他用户偏好）
 *   - 损坏数据容错（返回空对象，不抛异常）
 */

export type SidebarGroupKey =
  | 'dataManagement'
  | 'systemConfig'
  | 'operation'
  | 'statistics'

export type SidebarState = Partial<Record<SidebarGroupKey, boolean>>

const VALID_KEYS: readonly SidebarGroupKey[] = [
  'dataManagement',
  'systemConfig',
  'operation',
  'statistics',
]

function sanitize(state: Record<string, unknown>): SidebarState {
  const clean: SidebarState = {}
  for (const key of VALID_KEYS) {
    if (key in state) {
      clean[key] = Boolean(state[key])
    }
  }
  return clean
}

export function serializeSidebarState(state: SidebarState): string {
  return JSON.stringify(sanitize(state as Record<string, unknown>))
}

export function deserializeSidebarState(raw: string | null | undefined): SidebarState {
  if (!raw) return {}
  let parsed: unknown
  try {
    parsed = JSON.parse(raw)
  } catch {
    return {}
  }
  // 仅接受纯对象字面量（排除数组、字符串、数字、null）
  if (
    typeof parsed !== 'object' ||
    parsed === null ||
    Array.isArray(parsed)
  ) {
    return {}
  }
  return sanitize(parsed as Record<string, unknown>)
}

/**
 * 进入路由时自动展开当前分组
 *
 * - groupKey 为合法可持久化分组时：在保留 prev 其他偏好的基础上，将该分组置为 true
 * - groupKey 为 null（如 dashboard 或未匹配）：直接返回 prev 拷贝，不引入副作用
 * - 不修改传入的原对象
 */
export function expandForGroup(
  prev: SidebarState,
  groupKey: SidebarGroupKey | null,
): SidebarState {
  const next: SidebarState = { ...prev }
  if (groupKey && (VALID_KEYS as readonly string[]).includes(groupKey)) {
    next[groupKey] = true
  }
  return next
}
