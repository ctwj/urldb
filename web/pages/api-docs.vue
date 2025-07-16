<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100">
    <!-- 头部 -->
    <div class="bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center">
          <h1 class="text-3xl sm:text-4xl font-bold mb-4">
            网盘资源管理系统 - API文档
          </h1>
          <p class="text-lg text-gray-300 max-w-2xl mx-auto">
            公开API接口文档，支持资源添加、搜索和热门剧获取等功能
          </p>
          <div class="mt-6 flex flex-col sm:flex-row gap-4 justify-center">
            <NuxtLink 
              to="/" 
              class="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors text-center"
            >
              <i class="fas fa-home mr-2"></i>返回首页
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
        <!-- 单个添加资源 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
          <div class="bg-green-600 text-white px-6 py-4">
            <h3 class="text-xl font-semibold flex items-center">
              <i class="fas fa-plus-circle mr-2"></i>
              单个添加资源
            </h3>
            <p class="text-green-100 mt-1">添加单个资源到待处理列表</p>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求信息</h4>
                <div class="space-y-2 text-sm">
                  <p><strong>方法：</strong><span class="bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-200 px-2 py-1 rounded">POST</span></p>
                  <p><strong>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">/api/public/resources/add</code></p>
                  <p><strong>认证：</strong><span class="text-red-600 dark:text-red-400">必需</span></p>
                </div>
              </div>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求参数</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "title": "资源标题",
  "description": "资源描述",
  "url": "资源链接",
  "category": "分类名称",
  "tags": "标签1,标签2",
  "img": "封面图片链接",
  "source": "数据来源",
  "extra": "额外信息"
}</code></pre>
                </div>
              </div>
            </div>
            <div class="mt-6">
              <h4 class="font-semibold text-gray-900 dark:text-white mb-3">响应示例</h4>
              <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                <pre class="text-sm overflow-x-auto"><code>{
  "success": true,
  "message": "资源添加成功，已进入待处理列表",
  "data": {
    "id": 123
  },
  "code": 200
}</code></pre>
              </div>
            </div>
          </div>
        </div>

        <!-- 批量添加资源 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
          <div class="bg-purple-600 text-white px-6 py-4">
            <h3 class="text-xl font-semibold flex items-center">
              <i class="fas fa-layer-group mr-2"></i>
              批量添加资源
            </h3>
            <p class="text-purple-100 mt-1">批量添加多个资源到待处理列表</p>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求信息</h4>
                <div class="space-y-2 text-sm">
                  <p><strong>方法：</strong><span class="bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-200 px-2 py-1 rounded">POST</span></p>
                  <p><strong>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded">/api/public/resources/batch-add</code></p>
                  <p><strong>认证：</strong><span class="text-red-600 dark:text-red-400">必需</span></p>
                </div>
              </div>
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white mb-3">请求参数</h4>
                <div class="bg-gray-50 dark:bg-gray-700 rounded p-4">
                  <pre class="text-sm overflow-x-auto"><code>{
  "resources": [
    {
      "title": "资源1",
      "url": "链接1",
      "description": "描述1"
    },
    {
      "title": "资源2", 
      "url": "链接2",
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
  "message": "批量添加成功，共添加 2 个资源",
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
            <p class="text-blue-100 mt-1">搜索资源，支持关键词、标签、分类过滤</p>
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
  "message": "搜索成功",
  "data": {
    "resources": [
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
    "page_size": 20
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
  "message": "获取热门剧成功",
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

# 单个添加资源
curl -X POST "http://localhost:8080/api/public/resources/add" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: $API_TOKEN" \
  -d '{
    "title": "测试资源",
    "url": "https://example.com/resource",
    "description": "测试描述"
  }'

# 搜索资源
curl -X GET "http://localhost:8080/api/public/resources/search?keyword=测试" \
  -H "X-API-Token: $API_TOKEN"

# 获取热门剧
curl -X GET "http://localhost:8080/api/public/hot-dramas?page=1&page_size=5" \
  -H "X-API-Token: $API_TOKEN"</code></pre>
          </div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <footer class="bg-gray-800 dark:bg-gray-900 text-gray-300 py-8 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <p>&copy; 2024 网盘资源管理系统. 保留所有权利.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
// 页面元数据
useHead({
  title: 'API文档 - 网盘资源管理系统',
  meta: [
    { name: 'description', content: '网盘资源管理系统的公开API接口文档' },
    { name: 'keywords', content: 'API,接口文档,网盘资源管理' }
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
</style> 