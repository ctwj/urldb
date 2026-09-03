// 用户导航配置
export interface NavigationItem {
  to: string
  icon: string
  label: string
  active: (route: any) => boolean
  permission?: string // 权限要求
  description?: string // 页面描述
}

// 用户导航菜单配置
export const userNavigationConfig: NavigationItem[] = [
  {
    to: '/user/resources',
    icon: 'fas fa-cloud',
    label: '我的资源',
    active: (route: any) => route.path.startsWith('/user/resources'),
    description: '管理您的个人资源'
  },
  {
    to: '/user/downloads',
    icon: 'fas fa-download',
    label: '下载历史',
    active: (route: any) => route.path.startsWith('/user/downloads'),
    description: '查看和管理下载历史记录'
  },
  {
    to: '/user/settings',
    icon: 'fas fa-cog',
    label: '设置',
    active: (route: any) => route.path.startsWith('/user/settings'),
    description: '账户设置和偏好'
  }
]

// 管理员额外导航项
export const adminNavigationConfig: NavigationItem[] = [
  {
    to: '/admin',
    icon: 'fas fa-user-shield',
    label: '管理后台',
    active: (route: any) => route.path.startsWith('/admin'),
    permission: 'admin',
    description: '系统管理功能'
  }
]

// 获取完整导航配置
export const getFullNavigationConfig = (userRole?: string): NavigationItem[] => {
  const config = [...userNavigationConfig]
  
  // 如果是管理员，添加管理功能
  if (userRole === 'admin') {
    config.push(...adminNavigationConfig)
  }
  
  return config
}

// 用户菜单项配置
export interface UserMenuItem {
  to?: string
  icon?: string
  label?: string
  type: 'link' | 'button' | 'divider'
  action?: () => void
  className?: string
}

export const getUserMenuItems = (handleLogout: () => void): UserMenuItem[] => [
  {
    to: '/user/settings',
    icon: 'fas fa-cog',
    label: '设置',
    type: 'link'
  },
  {
    type: 'divider'
  },
  {
    type: 'button',
    icon: 'fas fa-sign-out-alt',
    label: '退出登录',
    action: handleLogout,
    className: 'text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300'
  }
] 