// 管理后台导航配置
export interface AdminNavigationItem {
  to: string
  icon: string
  label: string
  active: (route: any) => boolean
  permission?: string // 权限要求
  description?: string // 页面描述
  group?: string // 分组
}

// 管理后台导航菜单配置
export const adminNewNavigationItems = [
  {
    key: 'dashboard',
    label: '仪表盘',
    icon: 'fas fa-tachometer-alt',
    to: '/admin',
    active: (route: any) => route.path === '/admin',
    group: 'dashboard'
  },
  // 运营管理分组
  {
    key: 'resources',
    label: '资源管理',
    icon: 'fas fa-database',
    to: '/admin/resources',
    active: (route: any) => route.path.startsWith('/admin/resources'),
    group: 'data'
  },
  {
    key: 'ready-resources',
    label: '待处理资源',
    icon: 'fas fa-clock',
    to: '/admin/ready-resources',
    active: (route: any) => route.path.startsWith('/admin/ready-resources'),
    group: 'data'
  },
  {
    key: 'categories',
    label: '分类管理',
    icon: 'fas fa-folder',
    to: '/admin/categories',
    active: (route: any) => route.path.startsWith('/admin/categories'),
    group: 'data'
  },
  {
    key: 'tags',
    label: '标签管理',
    icon: 'fas fa-tags',
    to: '/admin/tags',
    active: (route: any) => route.path.startsWith('/admin/tags'),
    group: 'data'
  },
  {
    key: 'platforms',
    label: '平台管理',
    icon: 'fas fa-cloud',
    to: '/admin/platforms',
    active: (route: any) => route.path.startsWith('/admin/platforms'),
    group: 'operation'
  },
  {
    key: 'accounts',
    label: '账号管理',
    icon: 'fas fa-user-shield',
    to: '/admin/accounts',
    active: (route: any) => route.path.startsWith('/admin/accounts'),
    group: 'operation'
  },
  {
    key: 'data-transfer',
    label: '数据转存管理',
    icon: 'fas fa-exchange-alt',
    to: '/admin/data-transfer',
    active: (route: any) => route.path.startsWith('/admin/data-transfer'),
    group: 'operation'
  },
  {
    key: 'tasks',
    label: '任务管理',
    icon: 'fas fa-tasks',
    to: '/admin/tasks',
    active: (route: any) => route.path.startsWith('/admin/tasks'),
    group: 'operation'
  },
  {
    key: 'seo',
    label: 'SEO',
    icon: 'fas fa-search',
    to: '/admin/seo',
    active: (route: any) => route.path.startsWith('/admin/seo'),
    group: 'operation'
  },
  {
    key: 'data-push',
    label: '数据推送',
    icon: 'fas fa-upload',
    to: '/admin/data-push',
    active: (route: any) => route.path.startsWith('/admin/data-push'),
    group: 'operation'
  },
  {
    key: 'files',
    label: '文件管理',
    icon: 'fas fa-file-upload',
    to: '/admin/files',
    active: (route: any) => route.path.startsWith('/admin/files'),
    group: 'data'
  },
  {
    key: 'bot',
    label: '机器人',
    icon: 'fas fa-robot',
    to: '/admin/bot',
    active: (route: any) => route.path.startsWith('/admin/bot'),
    group: 'operation'
  },
  // 统计分析分组
  {
    key: 'search-stats',
    label: '搜索统计',
    icon: 'fas fa-chart-line',
    to: '/admin/search-stats',
    active: (route: any) => route.path.startsWith('/admin/search-stats'),
    group: 'statistics'
  },
  {
    key: 'third-party-stats',
    label: '三方统计',
    icon: 'fas fa-chart-bar',
    to: '/admin/third-party-stats',
    active: (route: any) => route.path.startsWith('/admin/third-party-stats'),
    group: 'statistics'
  },
  // 系统管理分组
  {
    key: 'users',
    label: '用户管理',
    icon: 'fas fa-users',
    to: '/admin/users',
    active: (route: any) => route.path.startsWith('/admin/users'),
    group: 'system'
  },
  {
    key: 'system-config',
    label: '系统配置',
    icon: 'fas fa-cog',
    to: '/admin/site-config',
    active: (route: any) => route.path.startsWith('/admin/site-config'),
    group: 'system'
  },
  {
    key: 'version',
    label: '版本信息',
    icon: 'fas fa-code-branch',
    to: '/admin/version',
    active: (route: any) => route.path.startsWith('/admin/version'),
    group: 'system'
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