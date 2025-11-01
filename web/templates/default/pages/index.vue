<template>
  <div class="min-h-screen bg-gray-50 dark:bg-slate-900 text-gray-800 dark:text-slate-100 flex flex-col">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在初始化系统</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
        <!-- 头部 -->
        <div class="header-container bg-slate-800 dark:bg-slate-800 text-white dark:text-slate-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
          <h1 class="text-2xl sm:text-3xl font-bold mb-4 flex items-center justify-center gap-3">
            <img
              src="/assets/images/logo.webp"
              alt="Logo"
              class="h-8 w-auto object-contain"
            />
            <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
              老九网盘资源数据库 (模板版)
            </a>
          </h1>

          <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
            <NuxtLink to="/hot-dramas" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-film text-xs"></i> 热播剧
              </n-button>
            </NuxtLink>
            <NuxtLink to="/monitor" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-chart-line text-xs"></i> 系统监控
              </n-button>
            </NuxtLink>
            <NuxtLink to="/api-docs" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-book text-xs"></i> API文档
              </n-button>
            </NuxtLink>
          </nav>
        </div>

        <!-- 公告信息 -->
        <div class="w-full max-w-3xl mx-auto mb-2 px-2 sm:px-0">
          <div class="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-4 text-blue-800 dark:text-blue-200">
            <p class="text-center font-medium">这是默认模板的首页，展示了模板系统的工作原理</p>
          </div>
        </div>

        <!-- 搜索区域 -->
        <div class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
          <ClientOnly>
            <div class="relative">
              <n-input round placeholder="搜索" v-model:value="searchQuery" @blur="handleSearch" @keyup.enter="handleSearch" clearable>
                <template #prefix>
                  <i class="fas fa-search text-gray-400"></i>
                </template>
              </n-input>
            </div>
          </ClientOnly>

          <!-- 平台类型筛选 -->
          <div class="mt-3 flex flex-wrap gap-2" id="platformFilters">
            <a
              href="/"
              class="px-2 py-1 text-xs rounded-full bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 hover:bg-slate-700 dark:hover:bg-gray-600 transition-colors active-filter"
            >
              全部
            </a>
          </div>

          <!-- 统计信息 -->
          <div class="flex justify-between mt-3 text-sm text-gray-600 dark:text-gray-300 px-2">
            <div class="flex items-center">
              <i class="fas fa-calendar-day text-pink-600 mr-1"></i>
              今日资源: <span class="font-medium text-pink-600 ml-1">0</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-database text-blue-600 mr-1"></i>
              总资源数: <span class="font-medium text-blue-600 ml-1">0</span>
            </div>
          </div>
        </div>

        <!-- 模板说明 -->
        <div class="bg-white dark:bg-slate-800 rounded-lg shadow-lg shadow-slate-900/10 dark:shadow-slate-900/50 p-6 mb-8">
          <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">模板系统说明</h2>
          <div class="prose dark:prose-invert max-w-none">
            <p>这是默认模板的首页，展示了模板系统的工作原理：</p>
            <ul>
              <li>不同模板可以有不同的布局和页面内容</li>
              <li>模板切换可以通过后台管理页面进行</li>
              <li>每个模板可以有自己独特的设计风格</li>
              <li>页面内容可以根据模板进行定制</li>
            </ul>
            <p>当前使用的是 <strong>默认模板</strong>，对比查看博客模板可以看到不同的设计风格。</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-gray-600 dark:text-gray-400">
        <p>&copy; 2025 老九网盘资源数据库. 保留所有权利.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// 默认模板首页特定的SEO设置
useSeoMeta({
  title: '老九网盘资源数据库 (模板版) - 首页',
  description: '老九网盘资源数据库，现代化的网盘资源数据库，支持多网盘自动化转存分享',
  ogTitle: '老九网盘资源数据库 (模板版) - 首页',
  ogDescription: '老九网盘资源数据库，现代化的网盘资源数据库，支持多网盘自动化转存分享',
  ogImage: '/og-image.jpg',
  ogUrl: 'https://pan.l9.lc',
  twitterCard: 'summary_large_image',
  twitterTitle: '老九网盘资源数据库 (模板版) - 首页',
  twitterDescription: '老九网盘资源数据库，现代化的网盘资源数据库，支持多网盘自动化转存分享',
  twitterImage: '/og-image.jpg'
})

// 响应式数据
const pageLoading = ref(false)
const searchQuery = ref('')

// 搜索处理
const handleSearch = () => {
  if (searchQuery.value) {
    window.location.href = `/?search=${encodeURIComponent(searchQuery.value)}`
  }
}
</script>

<style scoped>
.header-container {
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom,
      rgba(0,0,0,0.1) 0%,
      rgba(0,0,0,0.25) 100%
  );
}

.active-filter {
  @apply bg-slate-800 text-white;
}
</style>