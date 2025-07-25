export default defineNuxtRouteMiddleware((to, from) => {
  // 只在客户端执行
  if (process.client) {
    const ua = navigator.userAgent
    // 检测是否为QQ或微信内置浏览器
    if (ua.includes('QQ/') || ua.includes('MicroMessenger')) {
      // 如果当前不在禁止访问页面，则跳转
      if (to.path !== '/forbidden') {
        return navigateTo('/forbidden')
      }
    }
  }
}) 