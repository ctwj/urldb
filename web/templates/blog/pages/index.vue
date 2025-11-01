<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-12">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 博客风格的首页 -->
      <header class="text-center mb-16">
        <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">欢迎来到老九博客</h1>
        <p class="text-xl text-gray-600 dark:text-gray-300">分享最新的网盘资源和影视资讯</p>
      </header>

      <!-- 特色内容 -->
      <section class="mb-16">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-8 pb-2 border-b border-gray-200 dark:border-gray-700">最新资源</h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- 资源卡片 -->
          <div
            v-for="resource in featuredResources"
            :key="resource.id"
            class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
          >
            <div class="p-6">
              <div class="flex items-center mb-4">
                <span class="text-blue-600 dark:text-blue-400 text-xl mr-2">
                  <i class="fas fa-cloud"></i>
                </span>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ resource.title }}</h3>
              </div>
              <p class="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">
                {{ resource.description || '暂无描述' }}
              </p>
              <div class="flex justify-between items-center">
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ formatDate(resource.updated_at) }}
                </span>
                <button
                  @click="showResourceLink(resource)"
                  class="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 transition-colors"
                >
                  获取链接
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 热门分类 -->
      <section class="mb-16">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-8 pb-2 border-b border-gray-200 dark:border-gray-700">热门分类</h2>

        <div class="flex flex-wrap gap-3">
          <a
            v-for="category in categories"
            :key="category.id"
            :href="`/category/${category.id}`"
            class="px-4 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-full text-gray-700 dark:text-gray-300 hover:bg-blue-50 dark:hover:bg-blue-900/30 hover:border-blue-300 dark:hover:border-blue-700 transition-colors"
          >
            {{ category.name }}
          </a>
        </div>
      </section>

      <!-- 统计信息 -->
      <section>
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-8 pb-2 border-b border-gray-200 dark:border-gray-700">站点统计</h2>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 text-center">
            <div class="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">{{ stats.total_resources }}</div>
            <div class="text-gray-600 dark:text-gray-400">总资源数</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 text-center">
            <div class="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">{{ stats.today_resources }}</div>
            <div class="text-gray-600 dark:text-gray-400">今日新增</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 text-center">
            <div class="text-3xl font-bold text-purple-600 dark:text-purple-400 mb-2">{{ stats.total_categories }}</div>
            <div class="text-gray-600 dark:text-gray-400">资源分类</div>
          </div>
        </div>
      </section>
    </div>

    <!-- 资源链接模态框 -->
    <QrCodeModal
      v-if="showLinkModal"
      :visible="showLinkModal"
      :url="selectedResource?.url"
      :save_url="selectedResource?.save_url"
      :loading="selectedResource?.loading"
      :linkType="selectedResource?.linkType"
      :platform="selectedResource?.platform"
      :message="selectedResource?.message"
      :error="selectedResource?.error"
      :forbidden="selectedResource?.forbidden"
      :forbidden_words="selectedResource?.forbidden_words"
      @close="showLinkModal = false"
    />
  </div>
</template>

<script setup lang="ts">
// 博客模板首页特定的SEO设置
useSeoMeta({
  title: '老九博客 - 分享最新的网盘资源和影视资讯',
  description: '老九博客专注于分享高质量的网盘资源，包括电影、电视剧、文档等各类资源',
  ogTitle: '老九博客 - 分享最新的网盘资源和影视资讯',
  ogDescription: '老九博客专注于分享高质量的网盘资源，包括电影、电视剧、文档等各类资源',
  ogImage: '/og-image.jpg',
  ogUrl: 'https://pan.l9.lc',
  twitterCard: 'summary_large_image',
  twitterTitle: '老九博客 - 分享最新的网盘资源和影视资讯',
  twitterDescription: '老九博客专注于分享高质量的网盘资源，包括电影、电视剧、文档等各类资源',
  twitterImage: '/og-image.jpg'
})

// 响应式数据
const showLinkModal = ref(false)
const selectedResource = ref<any>(null)

// 模拟数据 - 在实际应用中这些数据应该从API获取
const featuredResources = ref([
  {
    id: 1,
    title: '最新热门电影合集',
    description: '包含2025年最新上映的热门电影，高清无水印',
    updated_at: '2025-11-01T10:30:00Z',
    pan_id: 1
  },
  {
    id: 2,
    title: '经典电视剧回顾',
    description: '经典电视剧系列合集，包含多个季度',
    updated_at: '2025-10-30T15:45:00Z',
    pan_id: 2
  },
  {
    id: 3,
    title: '学习资料大礼包',
    description: '各类学习资料，包括编程、设计、语言学习等',
    updated_at: '2025-10-29T09:15:00Z',
    pan_id: 3
  }
])

const categories = ref([
  { id: 1, name: '电影' },
  { id: 2, name: '电视剧' },
  { id: 3, name: '纪录片' },
  { id: 4, name: '动漫' },
  { id: 5, name: '学习资料' },
  { id: 6, name: '软件工具' }
])

const stats = ref({
  total_resources: 12840,
  today_resources: 42,
  total_categories: 15
})

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 显示资源链接
const showResourceLink = (resource: any) => {
  selectedResource.value = { ...resource, loading: true }
  showLinkModal.value = true

  // 模拟API调用
  setTimeout(() => {
    selectedResource.value = {
      ...resource,
      loading: false,
      url: 'https://example.com/resource-link',
      save_url: 'https://example.com/save-link',
      linkType: 'direct',
      platform: 'baidu'
    }
  }, 1000)
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