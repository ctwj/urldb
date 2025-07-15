import { defineStore } from 'pinia'

export interface Resource {
  id: number
  title: string
  description: string
  url: string
  file_path: string
  file_size: number
  file_type: string
  category_id?: number
  category_name: string
  tags: string[]
  download_count: number
  view_count: number
  is_public: boolean
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface Stats {
  total_resources: number
  total_categories: number
  total_downloads: number
  total_views: number
}

export const useResourceStore = defineStore('resource', {
  state: () => ({
    resources: [] as Resource[],
    categories: [] as Category[],
    stats: null as Stats | null,
    loading: false,
    currentPage: 1,
    totalPages: 1,
    searchQuery: '',
    selectedCategory: null as number | null,
  }),

  getters: {
    getResourceById: (state) => (id: number) => {
      return state.resources.find(resource => resource.id === id)
    },
    
    getCategoryById: (state) => (id: number) => {
      return state.categories.find(category => category.id === id)
    },
  },

  actions: {
    async fetchResources(params?: any) {
      this.loading = true
      try {
        const { getResources } = useResourceApi()
        // 确保有默认参数
        const defaultParams = {
          page: 1,
          page_size: 100,
          ...params
        }
        console.log('fetchResources - 请求参数:', defaultParams)
        const data = await getResources(defaultParams) as any
        console.log('fetchResources - 返回数据:', data)
        
        // 添加更详细的错误检查
        if (!data) {
          console.error('fetchResources - 数据为空')
          this.resources = []
          return
        }
        
        if (!data.resources) {
          console.error('fetchResources - 缺少resources字段:', data)
          this.resources = []
          return
        }
        
        this.resources = data.resources || []
        this.currentPage = data.page || 1
        this.totalPages = Math.ceil((data.total || 0) / (data.page_size || 100))
        console.log('fetchResources - 设置成功:', {
          resourcesCount: this.resources.length,
          currentPage: this.currentPage,
          totalPages: this.totalPages
        })
      } catch (error) {
        console.error('获取资源失败:', error)
        this.resources = []
      } finally {
        this.loading = false
      }
    },

    async fetchCategories() {
      try {
        const { getCategories } = useCategoryApi()
        this.categories = await getCategories() as any
      } catch (error) {
        console.error('获取分类失败:', error)
      }
    },

    async fetchStats() {
      try {
        const { getStats } = useStatsApi()
        this.stats = await getStats() as any
      } catch (error) {
        console.error('获取统计失败:', error)
      }
    },

    async searchResources(query: string, categoryId?: number) {
      this.loading = true
      try {
        const { searchResources } = useResourceApi()
        const params = { q: query, category_id: categoryId }
        const data = await searchResources(params) as any
        this.resources = data.resources || []
        this.searchQuery = query
        this.selectedCategory = categoryId || null
      } catch (error) {
        console.error('搜索资源失败:', error)
      } finally {
        this.loading = false
      }
    },

    async createResource(resourceData: any) {
      try {
        const { createResource } = useResourceApi()
        await createResource(resourceData)
        await this.fetchResources()
      } catch (error) {
        console.error('创建资源失败:', error)
        throw error
      }
    },

    async updateResource(id: number, resourceData: any) {
      try {
        const { updateResource } = useResourceApi()
        await updateResource(id, resourceData)
        await this.fetchResources()
      } catch (error) {
        console.error('更新资源失败:', error)
        throw error
      }
    },

    async deleteResource(id: number) {
      try {
        const { deleteResource } = useResourceApi()
        await deleteResource(id)
        await this.fetchResources()
      } catch (error) {
        console.error('删除资源失败:', error)
        throw error
      }
    },

    async createCategory(categoryData: any) {
      try {
        const { createCategory } = useCategoryApi()
        await createCategory(categoryData)
        await this.fetchCategories()
      } catch (error) {
        console.error('创建分类失败:', error)
        throw error
      }
    },

    async updateCategory(id: number, categoryData: any) {
      try {
        const { updateCategory } = useCategoryApi()
        await updateCategory(id, categoryData)
        await this.fetchCategories()
      } catch (error) {
        console.error('更新分类失败:', error)
        throw error
      }
    },

    async deleteCategory(id: number) {
      try {
        const { deleteCategory } = useCategoryApi()
        await deleteCategory(id)
        await this.fetchCategories()
      } catch (error) {
        console.error('删除分类失败:', error)
        throw error
      }
    },

    // 直接设置资源列表，用于搜索结果显示
    setResources(resources: Resource[]) {
      this.resources = resources
      console.log('setResources - 设置资源:', resources.length)
    },

    clearSearch() {
      this.searchQuery = ''
      this.selectedCategory = null
      this.fetchResources()
    },
  },
}) 