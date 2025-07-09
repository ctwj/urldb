<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 头部 -->
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center">
            <h1 class="text-2xl font-bold text-gray-900">资源管理系统</h1>
          </div>
          <div class="flex items-center space-x-4">
            <button @click="showAddResourceModal = true" class="btn-primary">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
              </svg>
              添加资源
            </button>
          </div>
        </div>
      </div>
    </header>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="card">
          <div class="flex items-center">
            <div class="p-2 bg-blue-100 rounded-lg">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">总资源</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats?.total_resources || 0 }}</p>
            </div>
          </div>
        </div>
        
        <div class="card">
          <div class="flex items-center">
            <div class="p-2 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">分类</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats?.total_categories || 0 }}</p>
            </div>
          </div>
        </div>
        
        <div class="card">
          <div class="flex items-center">
            <div class="p-2 bg-yellow-100 rounded-lg">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">下载</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats?.total_downloads || 0 }}</p>
            </div>
          </div>
        </div>
        
        <div class="card">
          <div class="flex items-center">
            <div class="p-2 bg-purple-100 rounded-lg">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">浏览</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats?.total_views || 0 }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 搜索和筛选 -->
      <div class="card mb-8">
        <div class="flex flex-col md:flex-row gap-4">
          <div class="flex-1">
            <div class="relative">
              <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
              <input
                v-model="searchQuery"
                @keyup.enter="handleSearch"
                type="text"
                placeholder="搜索资源..."
                class="input-field pl-10"
              />
            </div>
          </div>
          <div class="w-full md:w-48">
            <select v-model="selectedCategory" class="input-field">
              <option value="">全部分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          <div class="flex gap-2">
            <button @click="handleSearch" class="btn-primary">
              搜索
            </button>
            <button @click="clearSearch" class="btn-secondary">
              清除
            </button>
          </div>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div
          v-for="resource in resources"
          :key="resource.id"
          class="resource-card cursor-pointer"
          @click="viewResource(resource)"
        >
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1">
              <h3 class="font-semibold text-gray-900 truncate">{{ resource.title }}</h3>
              <p class="text-sm text-gray-600 mt-1 line-clamp-2">{{ resource.description }}</p>
            </div>
            <div class="flex items-center space-x-1 ml-2">
              <button
                @click.stop="editResource(resource)"
                class="p-1 text-gray-400 hover:text-blue-600"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                </svg>
              </button>
              <button
                @click.stop="deleteResource(resource.id)"
                class="p-1 text-gray-400 hover:text-red-600"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
          </div>
          
          <div class="flex items-center justify-between text-sm text-gray-500">
            <span>{{ resource.category_name || '未分类' }}</span>
            <div class="flex items-center space-x-3">
              <span class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
                {{ resource.view_count }}
              </span>
              <span class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                </svg>
                {{ resource.download_count }}
              </span>
            </div>
          </div>
          
          <div v-if="resource.tags && resource.tags.length > 0" class="mt-3">
            <div class="flex flex-wrap gap-1">
              <span
                v-for="tag in resource.tags.slice(0, 3)"
                :key="tag"
                class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded"
              >
                {{ tag }}
              </span>
              <span v-if="resource.tags.length > 3" class="px-2 py-1 text-gray-400 text-xs">
                +{{ resource.tags.length - 3 }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && resources.length === 0" class="text-center py-12">
        <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">暂无资源</h3>
        <p class="text-gray-600 mb-4">开始添加您的第一个资源吧</p>
        <button @click="showAddResourceModal = true" class="btn-primary">
          添加资源
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
  </div>
</template>

<script setup lang="ts">
const store = useResourceStore()
const { resources, categories, stats, loading } = storeToRefs(store)

const searchQuery = ref('')
const selectedCategory = ref('')
const showAddResourceModal = ref(false)
const editingResource = ref(null)

// 获取数据
onMounted(async () => {
  await Promise.all([
    store.fetchResources(),
    store.fetchCategories(),
    store.fetchStats(),
  ])
})

// 搜索处理
const handleSearch = () => {
  const categoryId = selectedCategory.value ? parseInt(selectedCategory.value) : undefined
  store.searchResources(searchQuery.value, categoryId)
}

// 清除搜索
const clearSearch = () => {
  searchQuery.value = ''
  selectedCategory.value = ''
  store.clearSearch()
}

// 查看资源
const viewResource = (resource: any) => {
  console.log('查看资源:', resource)
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
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style> 