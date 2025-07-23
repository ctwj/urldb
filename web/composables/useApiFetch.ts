import { useRuntimeConfig } from '#app'
import { useUserStore } from '~/stores/user'

export function useApiFetch<T = any>(
  url: string,
  options: any = {}
): Promise<T> {
  const config = useRuntimeConfig()
  const userStore = useUserStore()
  const baseURL = process.server
    ? String(config.public.apiServer)
    : String(config.public.apiBase)

  // 自动带上 token
  const headers = {
    ...(options.headers || {}),
    ...(userStore.authHeaders || {})
  }

  return $fetch<T>(url, {
    baseURL,
    ...options,
    headers,
    onResponse({ response }) {
      // 统一处理 code/message
      if (response._data && response._data.code && response._data.code !== 200) {
        throw new Error(response._data.message || '请求失败')
      }
    },
    onResponseError({ error }: { error: any }) {
      // 检查是否为"无效的令牌"错误
      if (error?.data?.error === '无效的令牌') {
        // 清除用户状态
        userStore.logout()
        // 跳转到登录页面
        if (process.client) {
          window.location.href = '/login'
        }
        throw new Error('登录已过期，请重新登录')
      }
      
      // 统一错误提示
      // 你可以用 naive-ui 的 useMessage() 这里弹窗
      // useMessage().error(error.message)
      throw error
    }
  })
} 