import { ref } from 'vue'
import type { GlobalThemeOverrides } from 'naive-ui'

/**
 * 主题管理 composable
 *
 * 提供 Naive UI 主题覆盖、Chart.js 默认配色与明暗模式切换。
 * 切换 mode 时同步 <html> 的 dark class（Tailwind darkMode:'class' 依赖）。
 * Chart.js 实例需在使用方监听 mode 变化时销毁重建（参考 data-model.md ThemeConfig）。
 */

export type ThemeMode = 'light' | 'dark'

// 模块级单例：跨组件共享同一主题状态
const mode = ref<ThemeMode>('light')

// Naive UI 全局主题覆盖（slate + blue-700 主色系）
const naiveOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#1d4ed8',
    primaryColorHover: '#2563eb',
    primaryColorPressed: '#1e40af',
    primaryColorSuppl: '#2563eb',
    borderRadius: '6px',
    borderRadiusSmall: '4px',
  },
}

// Chart.js 默认配色（getter 响应 mode 变化，使用方需在 mode 切换后重建实例）
export const chartDefaults = {
  fontFamily: 'Inter, system-ui, sans-serif',
  get gridColor() {
    return mode.value === 'dark' ? 'rgba(255,255,255,0.08)' : 'rgba(15,23,42,0.08)'
  },
  get primaryColor() {
    return mode.value === 'dark' ? '#60a5fa' : '#1d4ed8'
  },
  get textColor() {
    return mode.value === 'dark' ? '#cbd5e1' : '#475569'
  },
}

function applyMode(m: ThemeMode) {
  if (typeof document === 'undefined') return
  document.documentElement.classList.toggle('dark', m === 'dark')
}

export function useTheme() {
  const toggle = () => {
    mode.value = mode.value === 'light' ? 'dark' : 'light'
    applyMode(mode.value)
  }

  const setMode = (m: ThemeMode) => {
    mode.value = m
    applyMode(m)
  }

  return { mode, naiveOverrides, chartDefaults, toggle, setMode }
}
