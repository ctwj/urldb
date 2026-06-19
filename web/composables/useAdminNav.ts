import { computed } from 'vue'
import { useRoute } from '#imports'

/**
 * 后台导航数据源（侧边栏与命令面板共用）
 *
 * 从原 layouts/admin.vue 抽出的 5 个分组菜单数据，
 * 作为侧边栏渲染与 Cmd/Ctrl+K 命令面板索引的单一数据源。
 */

export interface NavItem {
  to: string
  label: string
  icon: string
  active: (route: { path: string }) => boolean
}

export interface NavGroup {
  key: 'dashboard' | 'dataManagement' | 'systemConfig' | 'operation' | 'statistics'
  title: string
  items: NavItem[]
}

// 命令面板索引项（由 NavGroup 扁平化得到）
export interface CommandPaletteItem {
  id: string
  title: string
  group: string
  icon: string
  to: string
  keywords: string[]
}

const dashboardItems: NavItem[] = [
  {
    to: '/admin',
    label: '仪表盘',
    icon: 'fas fa-tachometer-alt',
    active: (r) => r.path === '/admin',
  },
]

const dataManagementItems: NavItem[] = [
  { to: '/admin/resources', label: '资源管理', icon: 'fas fa-database', active: (r) => r.path.startsWith('/admin/resources') },
  { to: '/admin/ready-resources', label: '待处理资源', icon: 'fas fa-clock', active: (r) => r.path.startsWith('/admin/ready-resources') },
  { to: '/admin/tags', label: '标签管理', icon: 'fas fa-tags', active: (r) => r.path.startsWith('/admin/tags') },
  { to: '/admin/categories', label: '分类管理', icon: 'fas fa-folder', active: (r) => r.path.startsWith('/admin/categories') },
  { to: '/admin/accounts', label: '平台账号', icon: 'fas fa-user-shield', active: (r) => r.path.startsWith('/admin/accounts') },
  { to: '/admin/files', label: '文件管理', icon: 'fas fa-file-upload', active: (r) => r.path.startsWith('/admin/files') },
  { to: '/admin/reports', label: '举报管理', icon: 'fas fa-flag', active: (r) => r.path.startsWith('/admin/reports') },
  { to: '/admin/copyright-claims', label: '版权申述', icon: 'fas fa-balance-scale', active: (r) => r.path.startsWith('/admin/copyright-claims') },
]

const systemConfigItems: NavItem[] = [
  { to: '/admin/site-config', label: '站点配置', icon: 'fas fa-globe', active: (r) => r.path.startsWith('/admin/site-config') },
  { to: '/admin/feature-config', label: '功能配置', icon: 'fas fa-sliders-h', active: (r) => r.path.startsWith('/admin/feature-config') },
  { to: '/admin/dev-config', label: '开发配置', icon: 'fas fa-code', active: (r) => r.path.startsWith('/admin/dev-config') },
  { to: '/admin/plugins', label: '插件管理', icon: 'fas fa-plug', active: (r) => r.path.startsWith('/admin/plugins') },
  { to: '/admin/users', label: '用户管理', icon: 'fas fa-users', active: (r) => r.path.startsWith('/admin/users') },
]

const operationItems: NavItem[] = [
  { to: '/admin/data-transfer', label: '数据转存管理', icon: 'fas fa-exchange-alt', active: (r) => r.path.startsWith('/admin/data-transfer') },
  { to: '/admin/data-push', label: '数据推送', icon: 'fas fa-upload', active: (r) => r.path.startsWith('/admin/data-push') },
  { to: '/admin/bot', label: '机器人', icon: 'fas fa-robot', active: (r) => r.path.startsWith('/admin/bot') },
  { to: '/admin/seo', label: 'SEO', icon: 'fas fa-search', active: (r) => r.path.startsWith('/admin/seo') },
]

const statisticsItems: NavItem[] = [
  { to: '/admin/search-stats', label: '搜索统计', icon: 'fas fa-chart-line', active: (r) => r.path.startsWith('/admin/search-stats') },
  { to: '/admin/third-party-stats', label: '三方统计', icon: 'fas fa-chart-bar', active: (r) => r.path.startsWith('/admin/third-party-stats') },
]

const groups: NavGroup[] = [
  { key: 'dashboard', title: '仪表盘', items: dashboardItems },
  { key: 'dataManagement', title: '数据管理', items: dataManagementItems },
  { key: 'systemConfig', title: '系统配置', items: systemConfigItems },
  { key: 'operation', title: '运营管理', items: operationItems },
  { key: 'statistics', title: '统计分析', items: statisticsItems },
]

// 路径 → 分组 key 映射，用于 currentGroup 计算
const pathToGroup: Record<string, NavGroup['key']> = {}
for (const g of groups) {
  if (g.key === 'dashboard') continue
  for (const item of g.items) {
    // 取 item.to 的 /admin/xxx 段作为前缀
    const seg = item.to.replace('/admin/', '/admin/')
    pathToGroup[seg] = g.key
  }
}

export function useAdminNav() {
  const route = useRoute()

  const currentGroup = computed<NavGroup['key'] | null>(() => {
    const path = route.path
    for (const g of groups) {
      if (g.key === 'dashboard') continue
      if (g.items.some((item) => item.active({ path }))) {
        return g.key
      }
    }
    return null
  })

  const isItemActive = (item: NavItem): boolean => item.active({ path: route.path })

  return { groups, currentGroup, isItemActive }
}

// 将 NavGroup[] 扁平化为命令面板索引（含 keywords）
export function buildCommandPaletteItems(navGroups: NavGroup[]): CommandPaletteItem[] {
  const items: CommandPaletteItem[] = []
  for (const g of navGroups) {
    for (const item of g.items) {
      items.push({
        id: item.to,
        title: item.label,
        group: g.title,
        icon: item.icon,
        to: item.to,
        keywords: [item.label, g.title],
      })
    }
  }
  return items
}
