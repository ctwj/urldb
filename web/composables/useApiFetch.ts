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
    onResponseError({ error }) {
      // 统一错误提示
      // 你可以用 naive-ui 的 useMessage() 这里弹窗
      // useMessage().error(error.message)
      throw error
    }
  })
} 