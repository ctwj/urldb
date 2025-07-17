export default defineNuxtRouteMiddleware(async (to, from) => {
  // 只在客户端执行认证检查
  if (!process.client) {
    // auth middleware - 服务器端渲染，跳过认证检查
    return
  }
  
  const userStore = useUserStore()
  
  // 初始化用户状态
  userStore.initAuth()
  
  // 等待一小段时间确保认证状态初始化完成
  await new Promise(resolve => setTimeout(resolve, 100))
  
  // 检查认证状态
  if (!userStore.isAuthenticated) {
    // auth middleware - 用户未认证，重定向到登录页面
    // auth middleware - token: exists/not found
    // auth middleware - user: exists/not found
    return navigateTo('/login')
  }
  
      // auth middleware - 用户已认证，继续访问
    // auth middleware - 用户信息: username
}) 