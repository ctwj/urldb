export default defineNuxtPlugin(() => {
  // 只在客户端执行
  if (process.client) {
    const userStore = useUserStore()
    
    // 在应用启动时初始化用户认证状态
    // 使用 nextTick 确保在 DOM 更新后执行
    nextTick(() => {
      userStore.initAuth()
    })
  }
}) 