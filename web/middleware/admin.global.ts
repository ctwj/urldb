export default defineNuxtRouteMiddleware((to) => {
  // 全局路由守卫：拦截所有 /admin 路由（导航前执行，页面不会渲染）
  if (!to.path.startsWith('/admin')) return

  if (!process.client) return

  const userStore = useUserStore()
  userStore.initAuth()

  if (!userStore.isAuthenticated) {
    return navigateTo('/login')
  }

  // 以角色为准，不做用户名特判
  if (userStore.user?.role !== 'admin') {
    return navigateTo('/')
  }
})
