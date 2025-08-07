import { useUserLayout } from '~/composables/useUserLayout'
import { adminNewNavigationItems, adminNewMenuItems } from '~/config/adminNewNavigation'

export const useAdminNewLayout = () => {
  // 直接复用 useUserLayout
  const userLayout = useUserLayout()

  // 管理后台专用的认证检查 - 要求管理员权限
  const checkAdminAuth = () => {
    return userLayout.checkPermission('admin')
  }

  // 管理后台专用的用户信息
  const getAdminInfo = () => {
    const userInfo = userLayout.getUserInfo()
    return {
      ...userInfo,
      username: userInfo.username || '管理员',
      role: userInfo.role || 'admin',
      isAdmin: true
    }
  }

  // 管理后台专用的导航菜单
  const getAdminNavigationItems = () => {
    return adminNewNavigationItems
  }

  // 管理后台专用的菜单项
  const getAdminMenuItems = () => {
    return adminNewMenuItems
  }

  return {
    // 管理后台专用方法
    checkAuth: checkAdminAuth,
    getAdminInfo,
    getNavigationItems: getAdminNavigationItems,
    getAdminMenuItems,
    
    // 复用 useUserLayout 的所有方法
    checkPermission: userLayout.checkPermission,
    getUserInfo: userLayout.getUserInfo,
    handleLogout: userLayout.handleLogout,
    getUserMenuItems: userLayout.getUserMenuItems
  }
} 