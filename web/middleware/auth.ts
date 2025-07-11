export default defineNuxtRouteMiddleware((to, from) => {
  const userStore = useUserStore()
  
  // 初始化用户状态
  userStore.initAuth()
  
  // 如果用户未登录，重定向到登录页面
  if (!userStore.isAuthenticated) {
    return navigateTo('/login')
  }
}) 