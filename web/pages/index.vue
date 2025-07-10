<template>
  <div class="min-h-screen bg-gray-50 text-gray-800 p-3 sm:p-5">
    <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <div class="bg-slate-800 text-white rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4">
          <a href="/" class="text-white hover:text-gray-200 no-underline">网盘资源管理系统</a>
        </h1>
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-4">
          <button 
            @click="showAddResourceModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-plus"></i> 添加资源
          </button>
          <NuxtLink 
            to="/admin" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-user-shield"></i> 管理员入口
          </NuxtLink>
        </nav>
      </div>

      <!-- 搜索区域 -->
      <div class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
        <div class="relative">
          <input 
            v-model="searchQuery"
            @keyup="debounceSearch"
            type="text" 
            class="w-full px-4 py-3 rounded-full border-2 border-gray-300 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-200 transition-all"
            placeholder="输入文件名或链接进行搜索..."
          />
          <div class="absolute right-3 top-1/2 transform -translate-y-1/2">
            <i class="fas fa-search text-gray-400"></i>
          </div>
        </div>
        
        <!-- 平台类型筛选 -->
        <div class="mt-3 flex flex-wrap gap-2" id="platformFilters">
          <button 
            class="px-2 py-1 text-xs rounded-full bg-slate-800 text-white active-filter" 
            @click="filterByPlatform('')"
          >
            全部
          </button>
          <button 
            v-for="platform in platforms" 
            :key="platform.id"
            class="px-2 py-1 text-xs rounded-full bg-gray-200 text-gray-800 hover:bg-gray-300 transition-colors"
            @click="filterByPlatform(platform.id)"
          >
            {{ getPlatformIcon(platform.name) }} {{ platform.name }}
          </button>
        </div>
        
        <!-- 统计信息 -->
        <div class="flex justify-between mt-3 text-sm text-gray-600 px-2">
          <div class="flex items-center">
            <i class="fas fa-calendar-day text-pink-600 mr-1"></i>
            今日更新: <span class="font-medium text-pink-600 ml-1 count-up" :data-target="todayUpdates">0</span>
          </div>
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 mr-1"></i>
            总资源数: <span class="font-medium text-blue-600 ml-1 count-up" :data-target="stats?.total_resources || 0">0</span>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="overflow-x-auto bg-white rounded-lg shadow">
        <table class="w-full">
          <thead>
            <tr class="bg-slate-800 text-white">
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">
                <div class="flex items-center">
                  <i class="fas fa-cloud mr-1 text-gray-300"></i> 文件名
                </div>
              </th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell">链接</th>
              <th class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">更新时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="loading" class="text-center py-8">
              <td colspan="3" class="text-gray-500">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
            </tr>
            <tr v-else-if="resources.length === 0" class="text-center py-8">
              <td colspan="3" class="text-gray-500">暂无数据</td>
            </tr>
                         <tr 
               v-for="resource in (resources as unknown as ExtendedResource[])" 
               :key="resource.id"
               :class="isUpdatedToday(resource.updated_at) ? 'hover:bg-pink-50 bg-pink-50/30' : 'hover:bg-gray-50'"
             >
               <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm">
                 <div class="flex items-start">
                   <span class="mr-2 flex-shrink-0">{{ getPlatformIcon(getPlatformName(resource.pan_id || 0)) }}</span>
                   <span class="break-words">{{ resource.title }}</span>
                 </div>
                 <div class="sm:hidden mt-1">
                   <button 
                     class="text-blue-600 hover:text-blue-800 text-xs flex items-center gap-1 show-link-btn" 
                     @click="toggleLink(resource)"
                   >
                     <i class="fas fa-eye"></i> 显示链接
                   </button>
                   <a 
                     v-if="resource.showLink" 
                     :href="resource.url" 
                     target="_blank" 
                     class="text-blue-600 hover:text-blue-800 hover:underline text-xs break-all"
                   >
                     {{ resource.url }}
                   </a>
                 </div>
               </td>
               <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm hidden sm:table-cell">
                 <button 
                   class="text-blue-600 hover:text-blue-800 flex items-center gap-1 show-link-btn" 
                   @click="toggleLink(resource)"
                 >
                   <i class="fas fa-eye"></i> 显示链接
                 </button>
                 <a 
                   v-if="resource.showLink" 
                   :href="resource.url" 
                   target="_blank" 
                   class="text-blue-600 hover:text-blue-800 hover:underline"
                 >
                   {{ resource.url }}
                 </a>
               </td>
               <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm text-gray-500" :title="resource.updated_at">
                 {{ formatRelativeTime(resource.updated_at) }}
               </td>
             </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex flex-wrap justify-center gap-1 sm:gap-2 my-4 sm:my-8 px-2">
        <button 
          v-if="currentPage > 1"
          @click="goToPage(currentPage - 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          <i class="fas fa-chevron-left mr-1"></i> 上一页
        </button>
        
        <button 
          @click="goToPage(1)"
          :class="currentPage === 1 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          1
        </button>
        
        <button 
          v-if="totalPages > 1"
          @click="goToPage(2)"
          :class="currentPage === 2 ? 'bg-slate-800 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
          class="px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          2
        </button>
        
        <span v-if="currentPage > 2" class="px-2 py-1 sm:px-3 sm:py-2 text-gray-500 text-sm">...</span>
        
        <button 
          v-if="currentPage !== 1 && currentPage !== 2 && currentPage > 2"
          class="bg-slate-800 text-white px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm"
        >
          {{ currentPage }}
        </button>
        
        <button 
          v-if="currentPage < totalPages"
          @click="goToPage(currentPage + 1)"
          class="bg-white text-gray-700 hover:bg-gray-50 px-2 py-1 sm:px-4 sm:py-2 rounded border transition-colors text-sm flex items-center"
        >
          下一页 <i class="fas fa-chevron-right ml-1"></i>
        </button>
      </div>
    </div>

    <!-- 添加资源模态框 -->
    <ResourceModal
      v-if="showAddResourceModal"
      :resource="editingResource"
      @close="closeModal"
      @save="handleSaveResource"
    />

    <!-- 页脚 -->
    <footer class="mt-8 py-6 border-t border-gray-200">
      <div class="max-w-7xl mx-auto text-center text-gray-600 text-sm">
        <p class="mb-2">本站内容由网络爬虫自动抓取。本站不储存、复制、传播任何文件，仅作个人公益学习，请在获取后24小内删除!!!</p>
        <p>© 2025 网盘资源管理系统 By 小七</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const store = useResourceStore()
const { resources, categories, stats, loading } = storeToRefs(store)

const searchQuery = ref('')
const selectedPlatform = ref('')
const showAddResourceModal = ref(false)
const editingResource = ref(null)
const currentPage = ref(1)
const totalPages = ref(1)
interface Platform {
  id: number
  name: string
  key?: number
  ck?: string
  is_valid: boolean
  space: number
  left_space: number
  remark: string
  created_at: string
  updated_at: string
}

interface ExtendedResource {
  id: number
  title: string
  description: string
  url: string
  pan_id?: number
  quark_url?: string
  file_size?: string
  category_id?: number
  category_name?: string
  view_count: number
  is_valid?: boolean
  is_public?: boolean
  created_at: string
  updated_at: string
  showLink?: boolean
}

const platforms = ref<Platform[]>([])
const todayUpdates = ref(0)

// 防抖搜索
let searchTimeout: NodeJS.Timeout
const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    handleSearch()
  }, 500)
}

// 获取数据
onMounted(async () => {
  await Promise.all([
    store.fetchResources(),
    store.fetchCategories(),
    store.fetchStats(),
    fetchPlatforms(),
  ])
  animateCounters()
})

// 获取平台列表
const fetchPlatforms = async () => {
  try {
    const { usePanApi } = await import('~/composables/useApi')
    const panApi = usePanApi()
    const response = await panApi.getPans() as any
    platforms.value = response.pans || []
  } catch (error) {
    console.error('获取平台列表失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  const platformId = selectedPlatform.value ? parseInt(selectedPlatform.value) : undefined
  store.searchResources(searchQuery.value, platformId)
}

// 按平台筛选
const filterByPlatform = (platformId: string | number) => {
  selectedPlatform.value = platformId.toString()
  currentPage.value = 1
  handleSearch()
}

// 获取平台图标
const getPlatformIcon = (platformName: string) => {
  const icons: Record<string, string> = {
    '百度网盘': '<i class="fas fa-cloud text-blue-500"></i>',
    '阿里云盘': '<i class="fas fa-cloud text-orange-500"></i>',
    '夸克网盘': '<i class="fas fa-atom text-purple-500"></i>',
    '天翼云盘': '<i class="fas fa-cloud text-cyan-500"></i>',
    '迅雷云盘': '<i class="fas fa-bolt text-yellow-500"></i>',
    '微云': '<i class="fas fa-cloud text-green-500"></i>',
    '蓝奏云': '<i class="fas fa-cloud text-blue-400"></i>',
    '123云盘': '<i class="fas fa-cloud text-red-500"></i>',
    '腾讯微云': '<i class="fas fa-cloud text-green-500"></i>',
    'OneDrive': '<i class="fab fa-microsoft text-blue-600"></i>',
    'Google云盘': '<i class="fab fa-google-drive text-green-600"></i>',
    'Dropbox': '<i class="fab fa-dropbox text-blue-500"></i>',
    '城通网盘': '<i class="fas fa-folder text-yellow-600"></i>',
    '115网盘': '<i class="fas fa-cloud-upload-alt text-green-600"></i>',
    '磁力链接': '<i class="fas fa-magnet text-red-600"></i>',
    'UC网盘': '<i class="fas fa-cloud-download-alt text-purple-600"></i>',
    '天翼云': '<i class="fas fa-cloud text-cyan-500"></i>',
    'unknown': '<i class="fas fa-question-circle text-gray-400"></i>',
    '其他': '<i class="fas fa-cloud text-gray-500"></i>'
  }
  
  return icons[platformName] || icons['unknown']
}

// 获取平台名称
const getPlatformName = (platformId: number) => {
  const platform = platforms.value.find((p: Platform) => p.id === platformId)
  return platform?.name || 'unknown'
}

// 切换链接显示
const toggleLink = (resource: any) => {
  resource.showLink = !resource.showLink
}

// 格式化相对时间
const formatRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  const diffWeek = Math.floor(diffDay / 7)
  const diffMonth = Math.floor(diffDay / 30)
  const diffYear = Math.floor(diffDay / 365)
  
  const isToday = date.toDateString() === now.toDateString()
  
  if (isToday) {
    if (diffMin < 1) {
      return '<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>刚刚更新</span>'
    } else if (diffHour < 1) {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffMin}分钟前</span>`
    } else {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffHour}小时前</span>`
    }
  } else if (diffDay < 1) {
    return `${diffHour}小时前`
  } else if (diffDay < 7) {
    return `${diffDay}天前`
  } else if (diffWeek < 4) {
    return `${diffWeek}周前`
  } else if (diffMonth < 12) {
    return `${diffMonth}个月前`
  } else {
    return `${diffYear}年前`
  }
}

// 检查是否为今天更新
const isUpdatedToday = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  return date.toDateString() === now.toDateString()
}

// 数字动画效果
const animateCounters = () => {
  const counters = document.querySelectorAll('.count-up')
  const speed = 200
  
  counters.forEach((counter: Element) => {
    const target = parseInt(counter.getAttribute('data-target') || '0')
    const increment = Math.ceil(target / speed)
    let count = 0
    
    const updateCount = () => {
      if (count < target) {
        count += increment
        if (count > target) count = target
        counter.textContent = count.toString()
        setTimeout(updateCount, 1)
      } else {
        counter.textContent = target.toString()
      }
    }
    
    updateCount()
  })
}

// 页面跳转
const goToPage = (page: number) => {
  currentPage.value = page
  handleSearch()
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// 编辑资源
const editResource = (resource: any) => {
  editingResource.value = resource
  showAddResourceModal.value = true
}

// 删除资源
const deleteResource = async (id: number) => {
  if (confirm('确定要删除这个资源吗？')) {
    try {
      await store.deleteResource(id)
    } catch (error) {
      console.error('删除失败:', error)
    }
  }
}

// 关闭模态框
const closeModal = () => {
  showAddResourceModal.value = false
  editingResource.value = null
}

// 保存资源
const handleSaveResource = async (resourceData: any) => {
  try {
    if (editingResource.value) {
      await store.updateResource(editingResource.value.id, resourceData)
    } else {
      await store.createResource(resourceData)
    }
    closeModal()
  } catch (error) {
    console.error('保存失败:', error)
  }
}
</script>

<style scoped>
.active-filter {
  @apply bg-slate-800 text-white;
}

.count-up {
  transition: all 0.3s ease;
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
}
</style> 