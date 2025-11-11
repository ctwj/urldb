<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col">
    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
        <!-- 头部 -->
        <div class="header-container bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
          <h1 class="text-2xl sm:text-3xl font-bold mb-4">
            <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
              老九网盘资源数据库 - API文档
            </a>
          </h1>
          <p class="text-gray-300 max-w-2xl mx-auto">公开API接口文档，支持资源添加、搜索和热门剧获取等功能</p>
          <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
            <NuxtLink to="/" class="hidden sm:flex">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-home text-xs"></i> 首页
              </n-button>
            </NuxtLink>
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
          </nav>
        </div>

        <!-- 认证说明 -->
        <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6 mb-8">
          <h2 class="text-xl font-semibold text-blue-800 dark:text-blue-200 mb-4 flex items-center">
            <i class="fas fa-key mr-2"></i>
            API认证说明
          </h2>
          <div class="space-y-3 text-blue-700 dark:text-blue-300">
            <p><strong>认证方式：</strong>所有API都需要提供API Token进行认证</p>
            <p><strong>请求头方式：</strong><code class="bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded">X-API-Token: your_token</code></p>
            <p><strong>查询参数方式：</strong><code class="bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded">?api_token=your_token</code></p>
            <p><strong>获取Token：</strong>请联系管理员在系统配置中设置API Token</p>
          </div>
        </div>

        <!-- API接口列表 -->
        <div class="space-y-8">
          <!-- 批量添加资源 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
            <div class="bg-purple-600 text-white px-6 py-4">
              <h3 class="text-xl font-semibold flex items-center">
                <i class="fas fa-layer-group mr-2"></i>
                批量添加资源
              </h3>
              <p class="text-purple-100 mt-1">批量添加多个资源到待处理列表，每个资源可包含多个链接（url为数组），标题和url为必填项</p>
            </div>
            <div class="p-6">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求信息</h4>
                  <div class="space-y-2 text-sm">
                    <p><strong>方法：</strong><span class="bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-200 px-2 py-1 rounded">POST</span></p>
                    <p><strong>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">/api/public/resources/batch-add</code></p>
                    <p><strong>认证：</strong><span class="text-red-600 dark:text-red-400">必需</span>（X-API-Token）</p>
                  </div>
                </div>
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求参数</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">title 和 url 是必填项，其他字段均为选填</p>
                  <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                    <pre class="text-sm overflow-x-auto"><code>{
  "resources": [
    {
      "title": "资源1",
      "description": "描述1",
      "url": ["链接1", "链接2"],
      "category": "分类",
      "tags": "标签1,标签2",
      "img": "图片链接",
      "source": "数据来源",
      "extra": "额外信息"
    },
    {
      "title": "资源2",
      "url": ["链接3"],
      "description": "描述2"
    }
  ]
}</code></pre>
                  </div>
                </div>
              </div>
              <div class="mt-6">
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">响应示例</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "success": true,
  "message": "操作成功",
  "data": {
    "created_count": 2,
    "created_ids": [123, 124]
  },
  "code": 200
}</code></pre>
                </div>
              </div>
            </div>
          </div>

          <!-- 资源搜索 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
            <div class="bg-blue-600 text-white px-6 py-4">
              <h3 class="text-xl font-semibold flex items-center">
                <i class="fas fa-search mr-2"></i>
                资源搜索
              </h3>
              <p class="text-blue-100 mt-1">搜索资源，支持关键词、标签、分类过滤，自动过滤包含违禁词的资源</p>
            </div>
            <div class="p-6">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求信息</h4>
                  <div class="space-y-2 text-sm">
                    <p><strong>方法：</strong><span class="bg-blue-100 dark:bg-blue-800 text-blue-800 dark:text-blue-200 px-2 py-1 rounded">GET</span></p>
                    <p><strong>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">/api/public/resources/search</code></p>
                    <p><strong>认证：</strong><span class="text-red-600 dark:text-red-400">必需</span></p>
                  </div>
                </div>
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">查询参数</h4>
                  <div class="space-y-2 text-sm">
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">keyword</code> - 搜索关键词</p>
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">tag</code> - 标签过滤</p>
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">category</code> - 分类过滤</p>
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">page</code> - 页码（默认1）</p>
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">page_size</code> - 每页数量（默认20，最大100）</p>
                  </div>
                </div>
              </div>
              <div class="mt-6">
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">响应示例</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "success": true,
  "message": "操作成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "资源标题",
        "url": "资源链接",
        "description": "资源描述",
        "view_count": 100,
        "created_at": "2024-12-19 10:00:00",
        "updated_at": "2024-12-19 10:00:00"
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 20
  },
  "code": 200
}</code></pre>
                </div>
                <div class="mt-4 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
                  <h5 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2 flex items-center">
                    <i class="fas fa-exclamation-triangle mr-2"></i>
                    违禁词过滤说明
                  </h5>
                  <p class="text-sm text-yellow-700 dark:text-yellow-300 mb-2">当搜索结果包含违禁词时，响应会包含额外的过滤信息：</p>
                  <pre class="text-xs bg-yellow-100 dark:bg-yellow-800 rounded p-2"><code>{
  "success": true,
  "message": "操作成功",
  "data": {
    "list": [...],
    "total": 45,
    "page": 1,
    "limit": 20,
    "forbidden_words_filtered": true,
    "filtered_forbidden_words": ["违禁词1", "违禁词2"],
    "original_total": 50,
    "filtered_count": 5
  },
  "code": 200
}</code></pre>
                </div>
              </div>
            </div>
          </div>

          <!-- 热门剧 -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
            <div class="bg-orange-600 text-white px-6 py-4">
              <h3 class="text-xl font-semibold flex items-center">
                <i class="fas fa-film mr-2"></i>
                热门剧列表
              </h3>
              <p class="text-orange-100 mt-1">获取热门剧列表，支持分页</p>
            </div>
            <div class="p-6">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求信息</h4>
                  <div class="space-y-2 text-sm">
                    <p><strong>方法：</strong><span class="bg-orange-100 dark:bg-orange-800 text-orange-800 dark:text-orange-200 px-2 py-1 rounded">GET</span></p>
                    <p><strong>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">/api/public/hot-dramas</code></p>
                    <p><strong>认证：</strong><span class="text-red-600 dark:text-red-400">必需</span></p>
                  </div>
                </div>
                <div>
                  <h4 class="font-semibold text-gray-900 dark:text-white mb-3">查询参数</h4>
                  <div class="space-y-2 text-sm">
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">page</code> - 页码（默认1）</p>
                    <p><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">page_size</code> - 每页数量（默认20，最大100）</p>
                  </div>
                </div>
              </div>
              <div class="mt-6">
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">响应示例</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "success": true,
  "message": "操作成功",
  "data": {
    "hot_dramas": [
      {
        "id": 1,
        "title": "剧名",
        "description": "剧集描述",
        "img": "封面图片",
        "url": "详情链接",
        "rating": 8.5,
        "year": "2024",
        "region": "中国大陆",
        "genres": "剧情,悬疑",
        "category": "电视剧",
        "created_at": "2024-12-19 10:00:00",
        "updated_at": "2024-12-19 10:00:00"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 20
  },
  "code": 200
}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 错误码说明 -->
        <div class="mt-12 bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
          <div class="bg-red-600 text-white px-6 py-4">
            <h3 class="text-xl font-semibold flex items-center">
              <i class="fas fa-exclamation-triangle mr-2"></i>
              错误码说明
            </h3>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">HTTP状态码</h4>
                <div class="space-y-2 text-sm">
                  <p><span class="bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-200 px-2 py-1 rounded">200</span> - 请求成功</p>
                  <p><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded">400</span> - 请求参数错误</p>
                  <p><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded">401</span> - 认证失败（Token无效或缺失）</p>
                  <p><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded">500</span> - 服务器内部错误</p>
                  <p><span class="bg-yellow-100 dark:bg-yellow-800 text-yellow-800 dark:text-yellow-200 px-2 py-1 rounded">503</span> - 系统维护中或API Token未配置</p>
                </div>
              </div>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">响应格式</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "success": true/false,
  "message": "响应消息",
  "data": {}, // 响应数据
  "code": 200 // 状态码
}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 使用示例 -->
        <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
          <div class="bg-indigo-600 text-white px-6 py-4">
            <h3 class="text-xl font-semibold flex items-center">
              <i class="fas fa-code mr-2"></i>
              使用示例
            </h3>
          </div>
          <div class="p-6">
            <h4 class="font-semibold text-gray-900 dark:text-white mb-3">cURL示例</h4>
            <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
              <pre class="text-sm overflow-x-auto"><code># 设置API Token
API_TOKEN="your_api_token_here"

# 批量添加资源
curl -X POST "http://localhost:8080/api/public/resources/batch-add" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: $API_TOKEN" \
  -d '{
    "resources": [
      { "title": "测试资源1", "url": ["https://example.com/resource1"], "description": "描述1" },
      { "title": "测试资源2", "url": ["https://example.com/resource2", "https://example.com/resource3"], "description": "描述2" }
    ]
  }'

# 搜索资源
curl -X GET "http://localhost:8080/api/public/resources/search?keyword=测试" \
  -H "X-API-Token: $API_TOKEN"

# 获取热门剧
curl -X GET "http://localhost:8080/api/public/hot-dramas?page=1&page_size=5" \
  -H "X-API-Token: $API_TOKEN"</code></pre>
            </div>
            <h4 class="font-semibold text-gray-900 dark:text-white mb-3">JavaScript fetch 示例</h4>
            <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
              <pre class="text-sm overflow-x-auto"><code>// 资源搜索
fetch('/api/public/resources/search?keyword=测试', { 
  headers: { 'X-API-Token': 'your_token' } 
})
  .then(res => res.json())
  .then(res => {
    if (res.success) {
      const list = res.data.list // 资源列表
      const total = res.data.total
      console.log('搜索结果:', list)
    } else {
      console.error('搜索失败:', res.message)
    }
  })

// 批量添加资源
fetch('/api/public/resources/batch-add', {
  method: 'POST',
  headers: { 
    'Content-Type': 'application/json', 
    'X-API-Token': 'your_token' 
  },
  body: JSON.stringify({
    resources: [
      { title: '测试资源1', url: ['https://example.com/resource1'], description: '描述1' },
      { title: '测试资源2', url: ['https://example.com/resource2'], description: '描述2' }
    ]
  })
})
  .then(res => res.json())
  .then(res => {
    if (res.success) {
      console.log('添加成功，ID:', res.data.created_ids)
    } else {
      console.error('添加失败:', res.message)
    }
  })

// 获取热门剧
fetch('/api/public/hot-dramas?page=1&page_size=10', {
  headers: { 'X-API-Token': 'your_token' }
})
  .then(res => res.json())
  .then(res => {
    if (res.success) {
      const dramas = res.data.hot_dramas
      console.log('热门剧:', dramas)
    } else {
      console.error('获取失败:', res.message)
    }
  })</code></pre>
            </div>
            <h4 class="font-semibold text-gray-900 dark:text-white mb-3">Python requests 示例</h4>
            <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
              <pre class="text-sm overflow-x-auto"><code>import requests

API_TOKEN = 'your_api_token_here'
BASE_URL = 'http://localhost:8080/api'

headers = {
    'X-API-Token': API_TOKEN,
    'Content-Type': 'application/json'
}

# 搜索资源
def search_resources(keyword, page=1, page_size=20):
    params = {
        'keyword': keyword,
        'page': page,
        'page_size': page_size
    }
    response = requests.get(
        f'{BASE_URL}/public/resources/search',
        headers={'X-API-Token': API_TOKEN},
        params=params
    )
    return response.json()

# 批量添加资源
def batch_add_resources(resources):
    data = {'resources': resources}
    response = requests.post(
        f'{BASE_URL}/public/resources/batch-add',
        headers=headers,
        json=data
    )
    return response.json()

# 获取热门剧
def get_hot_dramas(page=1, page_size=20):
    params = {
        'page': page,
        'page_size': page_size
    }
    response = requests.get(
        f'{BASE_URL}/public/hot-dramas',
        headers={'X-API-Token': API_TOKEN},
        params=params
    )
    return response.json()

# 使用示例
if __name__ == '__main__':
    # 搜索资源
    result = search_resources('测试')
    if result['success']:
        print('搜索结果:', result['data']['list'])
    
    # 批量添加资源
    resources = [
        {'title': '测试资源1', 'url': ['https://example.com/resource1']},
        {'title': '测试资源2', 'url': ['https://example.com/resource2']}
    ]
    result = batch_add_resources(resources)
    if result['success']:
        print('添加成功，ID:', result['data']['created_ids'])
    
    # 获取热门剧
    result = get_hot_dramas()
    if result['success']:
        print('热门剧:', result['data']['hot_dramas'])</code></pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <AppFooter />
  </div>
</template>

<script setup>
// 设置页面布局
definePageMeta({
  layout: 'default'
})

// 页面元数据 - 使用Nuxt3 SEO最佳实践
useSeoMeta({
  title: 'API文档 - 老九网盘资源数据库',
  description: '老九网盘资源数据库的公开API接口文档，支持资源添加、搜索和热门剧获取等功能',
  keywords: 'API,接口文档,网盘资源管理,资源搜索,批量添加',
  ogTitle: 'API文档 - 老九网盘资源数据库',
  ogDescription: '老九网盘资源数据库的公开API接口文档，支持资源添加、搜索和热门剧获取等功能',
  ogType: 'website',
  ogImage: '/assets/images/logo.webp',
  twitterCard: 'summary_large_image',
  robots: 'index, follow'
})

useHead({
  htmlAttrs: {
    lang: 'zh-CN'
  },
  link: [
    {
      rel: 'canonical',
      href: 'https://yourdomain.com/api-docs'
    }
  ],
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        "@context": "https://schema.org",
        "@type": "WebPage",
        "name": "API文档 - 老九网盘资源数据库",
        "description": "老九网盘资源数据库的公开API接口文档，支持资源添加、搜索和热门剧获取等功能",
        "url": "https://yourdomain.com/api-docs"
      })
    }
  ]
})
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}

code {
  font-family: 'Courier New', Courier, monospace;
}
.header-container {
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}
</style> 