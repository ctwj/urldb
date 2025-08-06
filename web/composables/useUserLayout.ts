import { useUserStore } from '~/stores/user'
import { getFullNavigationConfig, getUserMenuItems as getMenuItems } from '~/config/userNavigation'

export const useUserLayout = () => {
  const userStore = useUserStore()
  const router = useRouter()

  // 检查用户认证状态
  const checkAuth = () => {
    if (!userStore.isAuthenticated) {
      router.push('/login')
      return false
    }
    return true
  }

  // 检查用户权限
  const checkPermission = (requiredRole: string = 'user') => {
    if (!userStore.isAuthenticated) {
      router.push('/login')
      return false
    }

    if (requiredRole === 'admin' && userStore.user?.role !== 'admin') {
      router.push('/user')
      return false
    }

    return true
  }

  // 获取用户信息
  const getUserInfo = () => {
    return {
      username: userStore.user?.username || '用户',
      email: userStore.user?.email || '',
      role: userStore.user?.role || 'user',
      isAdmin: userStore.user?.role === 'admin'
    }
  }

  // 处理退出登录
  const handleLogout = () => {
    userStore.logout()
    router.push('/login')
  }

  // 获取导航菜单项
  const getNavigationItems = () => {
    return getFullNavigationConfig(userStore.user?.role)
  }

  // 获取用户菜单项
  const getUserMenuItems = () => {
    return getMenuItems(handleLogout)
  }

  return {
    checkAuth,
    checkPermission,
    getUserInfo,
    handleLogout,
    getNavigationItems,
    getUserMenuItems
  }
} 