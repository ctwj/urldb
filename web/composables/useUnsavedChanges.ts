/**
 * 未保存变更检测
 *
 * 提供：
 * - `hasChanges` 响应式布尔：当前表单是否有未保存变更
 * - `markDirty()` / `markClean()` 手动标记
 * - `beforeunload` 事件监听：浏览器关闭/刷新前提示
 * - 路由切换守卫：onBeforeRouteLeave 提示
 *
 * 用法：
 *   const { hasChanges, markDirty, markClean } = useUnsavedChanges()
 *   watch(form, markDirty, { deep: true })
 *   // 保存成功后
 *   markClean()
 *
 * 决策依据：spec Clarifications —— 多标签页采用 last-write-wins，不引入冲突检测，
 * 因此本 composable 仅在「同标签页内的路由切换/刷新」时提示，不做跨标签页冲突感知。
 */
import { onBeforeRouteLeave } from 'vue-router'

export interface UnsavedChangesOptions {
  /** 提示文案 */
  message?: string
  /** 是否启用路由守卫（默认 true） */
  routeGuard?: boolean
  /** 是否启用 beforeunload（默认 true） */
  beforeUnload?: boolean
}

export function useUnsavedChanges(options: UnsavedChangesOptions = {}) {
  const {
    message = '当前页面有未保存的更改，确定离开吗？',
    routeGuard = true,
    beforeUnload = true,
  } = options

  const hasChanges = ref(false)

  const markDirty = () => {
    hasChanges.value = true
  }

  const markClean = () => {
    hasChanges.value = false
  }

  // 浏览器原生 beforeunload（关闭/刷新标签页）
  const onBeforeUnload = (e: BeforeUnloadEvent) => {
    if (!hasChanges.value) return
    e.preventDefault()
    e.returnValue = message
    return message
  }

  if (import.meta.client) {
    onMounted(() => {
      if (beforeUnload) window.addEventListener('beforeunload', onBeforeUnload)
    })
    onBeforeUnmount(() => {
      if (beforeUnload) window.removeEventListener('beforeunload', onBeforeUnload)
    })
  }

  // 路由切换守卫（仅 SPA 路由切换）
  if (routeGuard) {
    onBeforeRouteLeave(() => {
      if (!hasChanges.value) return true
      // 浏览器 confirm 是 Nuxt/Vue 路由守卫的标准做法（同步返回）
      if (typeof window !== 'undefined' && window.confirm(message)) {
        return true
      }
      return false
    })
  }

  return {
    hasChanges: readonly(hasChanges),
    markDirty,
    markClean,
  }
}
