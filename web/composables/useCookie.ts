export const useCookie = () => {
    const get = (name: string): string | null => {
      if (process.server) {
        // 服务端处理
        const { req } = useRequestEvent()
        const cookies = parse(req.headers.cookie || '')
        return cookies[name] || null
      } else {
        // 客户端处理
        const cookieString = document.cookie
        const cookies = cookieString.split('; ')
        for (const cookie of cookies) {
          const [cookieName, cookieValue] = cookie.split('=')
          if (cookieName === name) {
            return decodeURIComponent(cookieValue)
          }
        }
        return null
      }
    }
  
    const set = (name: string, value: string, options?: any) => {
      if (process.server) {
        // 服务端设置 cookie
        const { res } = useRequestEvent()
        res.setHeader('Set-Cookie', `${name}=${encodeURIComponent(value)}; ${serializeOptions(options)}`)
      } else {
        // 客户端设置 cookie
        document.cookie = `${name}=${encodeURIComponent(value)}; ${serializeOptions(options)}`
      }
    }
  
    const remove = (name: string) => {
      set(name, '', { maxAge: -1 })
    }
  
    // 序列化 cookie 选项
    const serializeOptions = (options: any = {}): string => {
      const {
        maxAge,
        expires,
        path = '/',
        domain,
        secure,
        httpOnly,
        sameSite
      } = options
  
      let result = ''
  
      if (maxAge !== undefined) result += `Max-Age=${maxAge}; `
      if (expires instanceof Date) result += `Expires=${expires.toUTCString()}; `
      if (path) result += `Path=${path}; `
      if (domain) result += `Domain=${domain}; `
      if (secure) result += 'Secure; '
      if (httpOnly) result += 'HttpOnly; '
      if (sameSite) result += `SameSite=${sameSite}; `
  
      return result
    }
  
    // 解析 cookie 字符串
    const parse = (cookieString: string): Record<string, string> => {
      const cookies: Record<string, string> = {}
      cookieString.split(';').forEach(cookie => {
        const [name, value] = cookie.split('=')
        if (name && value) {
          cookies[name.trim()] = decodeURIComponent(value.trim())
        }
      })
      return cookies
    }
  
    return {
      get,
      set,
      remove
    }
  }