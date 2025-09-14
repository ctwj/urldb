<template>
  <div class="flex flex-col gap-2 h-full">
    <!-- 搜索和筛选 -->
    <div class="flex-0 grid grid-cols-1 md:grid-cols-4 gap-4">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索已转存资源..."
        @keyup.enter="handleSearch"
        clearable
      >
        <template #prefix>
          <i class="fas fa-search"></i>
        </template>
      </n-input>
      
      <CategorySelector
        v-model="selectedCategory"
        placeholder="选择分类"
        clearable
      />
      
      <TagSelector
        v-model="selectedTag"
        placeholder="选择标签"
        clearable
      />
      
      <n-button type="primary" @click="handleSearch">
        <template #icon>
          <i class="fas fa-search"></i>
        </template>
        搜索
      </n-button>
    </div>

    <!-- 调试信息 -->
    <div class="flex-0 text-sm text-gray-500">
      数据数量: {{ resources.length }}, 总数: {{ total }}, 加载状态: {{ loading }}
    </div>

    <!-- 资源列表 -->
    <div class="flex-1 h-1">
      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="resources.length === 0" class="text-center py-8">
        <i class="fas fa-inbox text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无已转存的资源</p>
      </div>

      <div v-else class="h-full">
        <!-- 虚拟列表 -->
        <n-virtual-list
          :items="resources"
          :item-size="120"
          class="h-full"
        >
          <template #default="{ item }">
            <div class="p-4 border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800">
              <div class="flex items-start space-x-4">
                <!-- 资源信息 -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-center space-x-2 mb-2">
                    <!-- ID -->
                    <span class="text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 px-2 py-1 rounded">
                      ID: {{ item.id }}
                    </span>

                    <!-- 标题 -->
                    <h3 class="text-lg font-medium text-gray-900 dark:text-white line-clamp-1 flex-1">
                      {{ item.title || '未命名资源' }}
                    </h3>
                  </div>

                  <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm mt-3">
                    <!-- 分类 -->
                    <div class="flex items-center">
                      <i class="fas fa-folder mr-1 text-gray-400"></i>
                      <span class="text-gray-600 dark:text-gray-400">分类:</span>
                      <span class="ml-2">{{ item.category_name || '未分类' }}</span>
                    </div>

                    <!-- 转存时间 -->
                    <div class="flex items-center">
                      <i class="fas fa-calendar mr-1 text-gray-400"></i>
                      <span class="text-gray-600 dark:text-gray-400">转存时间:</span>
                      <span class="ml-2">{{ formatDate(item.updated_at) }}</span>
                    </div>

                    <!-- 浏览数 -->
                    <div class="flex items-center">
                      <i class="fas fa-eye mr-1 text-gray-400"></i>
                      <span class="text-gray-600 dark:text-gray-400">浏览数:</span>
                      <span class="ml-2">{{ item.view_count || 0 }}</span>
                    </div>
                  </div>

                  <!-- 转存链接 -->
                  <div class="mt-3">
                    <div class="flex items-start space-x-2">
                      <span class="text-xs text-gray-400">转存链接:</span>
                      <a
                        v-if="item.save_url"
                        :href="item.save_url"
                        target="_blank"
                        class="text-xs text-green-500 hover:text-green-700 break-all"
                      >
                        {{ item.save_url.length > 60 ? item.save_url.substring(0, 60) + '...' : item.save_url }}
                      </a>
                      <span v-else class="text-xs text-gray-500">暂无转存链接</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </n-virtual-list>
      </div>
    </div>

    <!-- 分页 -->
    <div class="flex-0">
      <div class="flex justify-center">
        <n-pagination
          v-model:page="currentPage"
          v-model:page-size="pageSize"
          :item-count="total"
          :page-sizes="[10000, 20000, 50000, 100000]"
          show-size-picker
          show-quick-jumper
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useResourceApi, usePanApi } from '~/composables/useApi'
import { useMessage } from 'naive-ui'

// 消息提示
const $message = useMessage()

// 数据状态
const loading = ref(false)
const resources = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10000)

// 搜索条件
const searchQuery = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)

// API实例
const resourceApi = useResourceApi()
const panApi = usePanApi()

// 获取平台数据
const { data: platformsData } = await useAsyncData('transferredPlatforms', () => panApi.getPans())

// 获取平台名称
const getPlatformName = (platformId: number) => {
  const platform = (platformsData.value as any)?.data?.find((plat: any) => plat.id === platformId)
  return platform?.remark || platform?.name || '未知平台'
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知时间'
  return new Date(dateString).toLocaleDateString()
}

// 获取已转存资源
const fetchTransferredResources = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      has_save_url: true // 筛选有转存链接的资源
    }

    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }

    console.log('请求参数:', params)
    const result = await resourceApi.getResources(params) as any
    console.log('已转存资源结果:', result)
    console.log('结果类型:', typeof result)
    console.log('结果结构:', Object.keys(result || {}))

    if (result && result.data) {
      // 处理嵌套的data结构：{data: {data: [...], total: ...}}
      if (result.data.data && Array.isArray(result.data.data)) {
        console.log('使用嵌套data格式，数量:', result.data.data.length)
        resources.value = result.data.data
        total.value = result.data.total || 0
      } else {
        // 处理直接的data结构：{data: [...], total: ...}
        console.log('使用直接data格式，数量:', result.data.length)
        resources.value = result.data
        total.value = result.total || 0
      }
    } else if (Array.isArray(result)) {
      console.log('使用数组格式，数量:', result.length)
      resources.value = result
      total.value = result.length
    } else {
      console.log('未知格式，设置空数组')
      resources.value = []
      total.value = 0
    }
    
    console.log('最终 resources.value:', resources.value)
    console.log('最终 total.value:', total.value)
    
    // 检查是否有资源没有 save_url
    const resourcesWithoutSaveUrl = resources.value.filter((r: any) => !r.save_url || r.save_url.trim() === '')
    if (resourcesWithoutSaveUrl.length > 0) {
      console.warn('发现没有 save_url 的资源:', resourcesWithoutSaveUrl.map((r: any) => ({ id: r.id, title: r.title, save_url: r.save_url })))
    }
  } catch (error) {
    console.error('获取已转存资源失败:', error)
    resources.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchTransferredResources()
}

// 分页处理
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchTransferredResources()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchTransferredResources()
}

// 初始化
onMounted(() => {
  fetchTransferredResources()
})
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>