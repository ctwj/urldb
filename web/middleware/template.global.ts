export default defineNuxtRouteMiddleware((to, from) => {
  // 从cookie获取当前模板
  const cookie = useCookie('site_template', { default: () => 'default' })
  const template = cookie.value || 'default'

  // 可以在这里添加模板特定的逻辑
  // 例如：为不同模板加载不同的CSS或JS

  // 将当前模板添加到nuxtApp中，以便在组件中使用
  const nuxtApp = useNuxtApp()
  nuxtApp.provide('currentTemplate', template)

  return true
})