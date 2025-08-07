// 管理后台导航配置
export interface AdminNavigationItem {
  to: string
  icon: string
  label: string
  active: (route: any) => boolean
  permission?: string // 权限要求
  description?: string // 页面描述
}

// 管理后台导航菜单配置
export const adminNewNavigationItems = [
  {
    key: 'dashboard',
    label: '仪表盘',
    icon: 'fas fa-tachometer-alt',
    to: '/admin',
    active: (route: any) => route.path === '/admin'
  },
  {
    key: 'resources',
    label: '资源管理',
    icon: 'fas fa-database',
    to: '/admin/resources',
    active: (route: any) => route.path.startsWith('/admin/resources')
  },
  {
    key: 'ready-resources',
    label: '待处理资源',
    icon: 'fas fa-clock',
    to: '/admin/ready-resources',
    active: (route: any) => route.path.startsWith('/admin/ready-resources')
  },
  {
    key: 'categories',
    label: '分类管理',
    icon: 'fas fa-folder',
    to: '/admin/categories',
    active: (route: any) => route.path.startsWith('/admin/categories')
  },
  {
    key: 'tags',
    label: '标签管理',
    icon: 'fas fa-tags',
    to: '/admin/tags',
    active: (route: any) => route.path.startsWith('/admin/tags')
  },
  {
    key: 'platforms',
    label: '平台管理',
    icon: 'fas fa-cloud',
    to: '/admin/platforms',
    active: (route: any) => route.path.startsWith('/admin/platforms')
  },
  {
    key: 'accounts',
    label: '账号管理',
    icon: 'fas fa-user-shield',
    to: '/admin/accounts',
    active: (route: any) => route.path.startsWith('/admin/accounts')
  },
  {
    key: 'hot-dramas',
    label: '热播剧管理',
    icon: 'fas fa-film',
    to: '/admin/hot-dramas',
    active: (route: any) => route.path.startsWith('/admin/hot-dramas')
  },
  {
    key: 'users',
    label: '用户管理',
    icon: 'fas fa-users',
    to: '/admin/users',
    active: (route: any) => route.path.startsWith('/admin/users')
  },
  {
    key: 'search-stats',
    label: '搜索统计',
    icon: 'fas fa-chart-line',
    to: '/admin/search-stats',
    active: (route: any) => route.path.startsWith('/admin/search-stats')
  },
  {
    key: 'system-config',
    label: '系统配置',
    icon: 'fas fa-cog',
    to: '/admin/system-config',
    active: (route: any) => route.path.startsWith('/admin/system-config')
  }
]

// 获取完整导航配置
export const getAdminNewNavigationConfig = (): AdminNavigationItem[] => {
  return [...adminNewNavigationItems]
}

// 管理员菜单项配置
export interface AdminMenuItem {
  to?: string
  icon?: string
  label?: string
  type: 'link' | 'button' | 'divider'
  action?: () => void
  className?: string
}

export const adminNewMenuItems = [
  {
    key: 'profile',
    label: '个人资料',
    icon: 'fas fa-user',
    to: '/admin/profile'
  },
  {
    key: 'settings',
    label: '设置',
    icon: 'fas fa-cog',
    to: '/admin/settings'
  },
  {
    key: 'logout',
    label: '退出登录',
    icon: 'fas fa-sign-out-alt',
    action: 'logout'
  }
] 