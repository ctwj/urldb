<template>
  <div class="min-h-screen bg-gray-50 text-gray-800 p-3 sm:p-5">
    <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <div class="bg-slate-800 text-white rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4">
          <NuxtLink to="/" class="text-white hover:text-gray-200 no-underline">网盘资源管理系统</NuxtLink>
        </h1>
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-4">
          <NuxtLink 
            to="/" 
            class="w-full sm:w-auto px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-home"></i> 返回首页
          </NuxtLink>
          <button 
            @click="showAddResourceModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-plus"></i> 添加资源
          </button>
        </nav>
      </div>

      <!-- 管理功能区域 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <!-- 资源管理 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-blue-100 rounded-lg">
              <i class="fas fa-cloud text-blue-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">资源管理</h3>
              <p class="text-sm text-gray-600">管理所有资源</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToResourceManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">查看所有资源</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddResourceModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">添加新资源</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 平台管理 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-green-100 rounded-lg">
              <i class="fas fa-server text-green-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">平台管理</h3>
              <p class="text-sm text-gray-600">管理网盘平台</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToPlatformManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">管理平台</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddPlatformModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">添加平台</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 分类管理 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-purple-100 rounded-lg">
              <i class="fas fa-folder text-purple-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">分类管理</h3>
              <p class="text-sm text-gray-600">管理资源分类</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToCategoryManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">管理分类</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddCategoryModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">添加分类</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 标签管理 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-orange-100 rounded-lg">
              <i class="fas fa-tags text-orange-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">标签管理</h3>
              <p class="text-sm text-gray-600">管理资源标签</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToTagManagement" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">管理标签</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="showAddTagModal = true" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">添加标签</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-red-100 rounded-lg">
              <i class="fas fa-chart-bar text-red-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">统计信息</h3>
              <p class="text-sm text-gray-600">系统统计数据</p>
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">总资源数</span>
              <span class="text-lg font-semibold text-gray-900">{{ stats?.total_resources || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">总浏览量</span>
              <span class="text-lg font-semibold text-gray-900">{{ stats?.total_views || 0 }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">分类数量</span>
              <span class="text-lg font-semibold text-gray-900">{{ stats?.total_categories || 0 }}</span>
            </div>
          </div>
        </div>

        <!-- 待处理资源 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-yellow-100 rounded-lg">
              <i class="fas fa-clock text-yellow-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">待处理资源</h3>
              <p class="text-sm text-gray-600">批量添加和管理</p>
            </div>
          </div>
          <div class="space-y-2">
            <NuxtLink to="/ready-resources" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors block">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">管理待处理资源</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </NuxtLink>
            <button @click="goToBatchAdd" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">批量添加资源</span>
                <i class="fas fa-plus text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>

        <!-- 系统设置 -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center mb-4">
            <div class="p-3 bg-gray-100 rounded-lg">
              <i class="fas fa-cog text-gray-600 text-xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">系统设置</h3>
              <p class="text-sm text-gray-600">系统配置</p>
            </div>
          </div>
          <div class="space-y-2">
            <button @click="goToSystemSettings" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">系统配置</span>
                <i class="fas fa-chevron-right text-gray-400"></i>
              </div>
            </button>
            <button @click="goToBackupRestore" class="w-full text-left p-3 rounded-lg hover:bg-gray-50 transition-colors">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">备份恢复</span>
                <i class="fas fa-database text-gray-400"></i>
              </div>
            </button>
          </div>
        </div>
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
const { stats } = storeToRefs(store)

const showAddResourceModal = ref(false)
const showAddPlatformModal = ref(false)
const showAddCategoryModal = ref(false)
const showAddTagModal = ref(false)
const editingResource = ref(null)

// 获取统计数据
onMounted(async () => {
  await store.fetchStats()
})

// 导航到资源管理
const goToResourceManagement = () => {
  // 这里可以导航到资源管理页面
  console.log('导航到资源管理')
}

// 导航到平台管理
const goToPlatformManagement = () => {
  // 这里可以导航到平台管理页面
  console.log('导航到平台管理')
}

// 导航到分类管理
const goToCategoryManagement = () => {
  // 这里可以导航到分类管理页面
  console.log('导航到分类管理')
}

// 导航到标签管理
const goToTagManagement = () => {
  // 这里可以导航到标签管理页面
  console.log('导航到标签管理')
}

// 导航到系统设置
const goToSystemSettings = () => {
  // 这里可以导航到系统设置页面
  console.log('导航到系统设置')
}

// 导航到备份恢复
const goToBackupRestore = () => {
  // 这里可以导航到备份恢复页面
  console.log('导航到备份恢复')
}

// 导航到批量添加
const goToBatchAdd = () => {
  // 导航到待处理资源页面
  navigateTo('/ready-resources')
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
/* 可以添加自定义样式 */
</style> 