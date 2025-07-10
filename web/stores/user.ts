import { defineStore } from 'pinia'

interface User {
  id: number
  username: string
  email: string
  role: string
  created_at: string
}

interface LoginForm {
  username: string
  password: string
}

interface RegisterForm {
  username: string
  email: string
  password: string
}

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null as User | null,
    token: null as string | null,
    isAuthenticated: false,
    loading: false
  }),

  getters: {
    // 检查是否为管理员
    isAdmin: (state) => state.user?.role === 'admin',
    
    // 获取用户信息
    userInfo: (state) => state.user,
    
    // 获取认证头
    authHeaders: (state) => {
      return state.token ? { Authorization: `Bearer ${state.token}` } : {}
    }
  },

  actions: {
    // 初始化用户状态（从localStorage恢复）
    initAuth() {
      const token = localStorage.getItem('token')
      const userStr = localStorage.getItem('user')
      
      if (token && userStr) {
        try {
          this.token = token
          this.user = JSON.parse(userStr)
          this.isAuthenticated = true
        } catch (error) {
          console.error('解析用户信息失败:', error)
          this.logout()
        }
      }
    },

      // 登录
  async login(credentials: LoginForm) {
    this.loading = true
    try {
      const authApi = useAuthApi()
      const response = await authApi.login(credentials)
      
      if (response.token) {
        this.token = response.token
        this.user = response.user
        this.isAuthenticated = true
        
        // 保存到localStorage
        localStorage.setItem('token', response.token)
        localStorage.setItem('user', JSON.stringify(response.user))
        
        return { success: true }
      } else {
        return { success: false, message: '登录失败，服务器未返回有效令牌' }
      }
    } catch (error: any) {
      console.error('登录错误:', error)
      // 处理HTTP错误响应
      if (error.data && error.data.error) {
        return { 
          success: false, 
          message: error.data.error 
        }
      }
      return { 
        success: false, 
        message: error.message || '登录失败，请检查网络连接' 
      }
    } finally {
      this.loading = false
    }
  },

      // 注册
  async register(userData: RegisterForm) {
    this.loading = true
    try {
      const authApi = useAuthApi()
      await authApi.register(userData)
      return { success: true }
    } catch (error: any) {
      console.error('注册错误:', error)
      // 处理HTTP错误响应
      if (error.data && error.data.error) {
        return { 
          success: false, 
          message: error.data.error 
        }
      }
      return { 
        success: false, 
        message: error.message || '注册失败，请检查网络连接' 
      }
    } finally {
      this.loading = false
    }
  },

    // 登出
    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      
      // 清除localStorage
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },

    // 获取用户资料
    async fetchProfile() {
      try {
        const authApi = useAuthApi()
        const user = await authApi.getProfile()
        this.user = user
        return { success: true }
      } catch (error: any) {
        console.error('获取用户资料失败:', error)
        // 如果获取失败，可能是token过期，清除登录状态
        this.logout()
        return { success: false, message: error.message }
      }
    }
  }
}) 