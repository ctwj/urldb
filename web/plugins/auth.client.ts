export default defineNuxtPlugin(() => {
  const userStore = useUserStore()
  
  // 在客户端初始化时恢复用户状态
  userStore.initAuth()
}) 