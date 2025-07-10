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
            @click="showAddModal = true" 
            class="w-full sm:w-auto px-4 py-2 bg-green-600 hover:bg-green-700 rounded-md transition-colors text-center flex items-center justify-center gap-2"
          >
            <i class="fas fa-plus"></i> 批量添加
          </button>
        </nav>
      </div>

      <!-- 批量添加模态框 -->
      <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg shadow-xl p-6 max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-bold">批量添加待处理资源</h3>
            <button @click="closeModal" class="text-gray-500 hover:text-gray-800">
              <i class="fas fa-times"></i>
            </button>
          </div>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">输入格式说明：</label>
            <div class="bg-gray-50 p-3 rounded text-sm text-gray-600 mb-4">
              <p class="mb-2"><strong>格式1：</strong>标题和URL两行一组</p>
              <pre class="bg-white p-2 rounded border text-xs">
电影标题1
https://pan.baidu.com/s/123456
电影标题2
https://pan.baidu.com/s/789012</pre>
              <p class="mt-2 mb-2"><strong>格式2：</strong>只有URL，系统自动判断</p>
              <pre class="bg-white p-2 rounded border text-xs">
https://pan.baidu.com/s/123456
https://pan.baidu.com/s/789012
https://pan.baidu.com/s/345678</pre>
            </div>
          </div>
          
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">资源内容：</label>
            <textarea
              v-model="resourceText"
              rows="15"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="请输入资源内容，支持两种格式..."
            ></textarea>
          </div>
          
          <div class="flex justify-end gap-2">
            <button @click="closeModal" class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50">
              取消
            </button>
            <button @click="handleBatchAdd" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
              批量添加
            </button>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold text-gray-900">待处理资源管理</h2>
        <div class="flex gap-2">
          <button 
            @click="refreshData" 
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 flex items-center gap-2"
          >
            <i class="fas fa-refresh"></i> 刷新
          </button>
          <button 
            @click="clearAll" 
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 flex items-center gap-2"
          >
            <i class="fas fa-trash"></i> 清空全部
          </button>
        </div>
      </div>

      <!-- 资源列表 -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="bg-slate-800 text-white">
                <th class="px-4 py-3 text-left text-sm">ID</th>
                <th class="px-4 py-3 text-left text-sm">标题</th>
                <th class="px-4 py-3 text-left text-sm">URL</th>
                <th class="px-4 py-3 text-left text-sm">创建时间</th>
                <th class="px-4 py-3 text-left text-sm">IP地址</th>
                <th class="px-4 py-3 text-left text-sm">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="loading" class="text-center py-8">
                <td colspan="6" class="text-gray-500">
                  <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
                </td>
              </tr>
              <tr v-else-if="readyResources.length === 0" class="text-center py-8">
                <td colspan="6" class="text-gray-500">暂无待处理资源</td>
              </tr>
              <tr 
                v-for="resource in readyResources" 
                :key="resource.id"
                class="hover:bg-gray-50"
              >
                <td class="px-4 py-3 text-sm text-gray-900">{{ resource.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-900">
                  <span v-if="resource.title">{{ resource.title }}</span>
                  <span v-else class="text-gray-400 italic">未设置</span>
                </td>
                <td class="px-4 py-3 text-sm">
                  <a 
                    :href="resource.url" 
                    target="_blank" 
                    class="text-blue-600 hover:text-blue-800 hover:underline break-all"
                  >
                    {{ resource.url }}
                  </a>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ formatTime(resource.create_time) }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ resource.ip || '-' }}
                </td>
                <td class="px-4 py-3 text-sm">
                  <button 
                    @click="deleteResource(resource.id)"
                    class="text-red-600 hover:text-red-800"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="mt-4 text-sm text-gray-600">
        共 {{ readyResources.length }} 个待处理资源
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ReadyResource {
  id: number
  title?: string
  url: string
  create_time: string
  ip?: string
}

const readyResources = ref<ReadyResource[]>([])
const loading = ref(false)
const showAddModal = ref(false)
const resourceText = ref('')

// 获取待处理资源API
const { useReadyResourceApi } = await import('~/composables/useApi')
const readyResourceApi = useReadyResourceApi()

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await readyResourceApi.getReadyResources() as any
    readyResources.value = response.resources || []
  } catch (error) {
    console.error('获取待处理资源失败:', error)
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchData()
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  resourceText.value = ''
}

// 批量添加
const handleBatchAdd = async () => {
  if (!resourceText.value.trim()) {
    alert('请输入资源内容')
    return
  }

  try {
    const response = await readyResourceApi.createReadyResourcesFromText(resourceText.value) as any
    console.log('批量添加成功:', response)
    closeModal()
    fetchData()
    alert(`成功添加 ${response.count} 个资源`)
  } catch (error) {
    console.error('批量添加失败:', error)
    alert('批量添加失败，请检查输入格式')
  }
}

// 删除资源
const deleteResource = async (id: number) => {
  if (!confirm('确定要删除这个待处理资源吗？')) {
    return
  }

  try {
    await readyResourceApi.deleteReadyResource(id)
    fetchData()
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败')
  }
}

// 清空全部
const clearAll = async () => {
  if (!confirm('确定要清空所有待处理资源吗？此操作不可恢复！')) {
    return
  }

  try {
    const response = await readyResourceApi.clearReadyResources() as any
    console.log('清空成功:', response)
    fetchData()
    alert(`成功清空 ${response.deleted_count} 个资源`)
  } catch (error) {
    console.error('清空失败:', error)
    alert('清空失败')
  }
}

// 格式化时间
const formatTime = (timeString: string) => {
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN')
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 