export default defineNuxtRouteMiddleware(async (to, from) => {
  // 只在客户端执行认证检查
  if (!process.client) {
    return
  }

  const userStore = useUserStore()

  // 初始化用户状态
  userStore.initAuth()

  // 等待一小段时间确保认证状态初始化完成
  await new Promise(resolve => setTimeout(resolve, 100))

  // 检查认证状态
  if (!userStore.isAuthenticated) {
    console.log('admin middleware - 用户未认证，重定向到登录页面')
    return navigateTo('/login')
  }

  // 检查用户是否为管理员（通常通过用户角色或权限判断）
  // 这里可以根据具体实现来调整，例如检查 userStore.user?.is_admin 字段
  const isAdmin = userStore.user?.is_admin || userStore.user?.role === 'admin' || userStore.user?.username === 'admin'

  if (!isAdmin) {
    console.log('admin middleware - 用户不是管理员，重定向到首页')
    return navigateTo('/')
  }

  console.log('admin middleware - 用户已认证且为管理员，继续访问')
})