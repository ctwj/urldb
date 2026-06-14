import{_ as l}from"./D8vZWwr2.js";import{_ as i}from"./Dd_fbjQY.js";import{u as p}from"./DG3JoL_b.js";import{_ as c,G as g,c as v,a,b as e,M as q,w as o,o as x,d}from"./DmHPR5lg.js";import{B as m}from"./DO8alW5h.js";import"./CpgoGwED.js";import"./g0YHQayI.js";import"./DFUnLHsf.js";import"./QpI9WcJO.js";import"./CoaUF789.js";const h={class:"min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-gray-100 flex flex-col"},y={class:"flex-1 p-3 sm:p-5"},f={class:"max-w-7xl mx-auto"},k={class:"header-container bg-slate-800 dark:bg-gray-800 text-white dark:text-gray-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative"},_={class:"mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute"},w={__name:"api-docs",setup(T){const{initSystemConfig:r,setApiDocsSeo:n}=p();return g(async()=>{await r(),n()}),(A,t)=>{const s=m,b=l,u=i;return x(),v("div",h,[a("div",y,[a("div",f,[a("div",k,[t[3]||(t[3]=a("h1",{class:"text-2xl sm:text-3xl font-bold mb-4"},[a("a",{href:"/",class:"text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline"}," 老九网盘资源数据库 - API文档 ")],-1)),t[4]||(t[4]=a("p",{class:"text-gray-300 max-w-2xl mx-auto"},"公开API接口文档，支持资源添加、搜索和热门剧获取等功能",-1)),a("nav",_,[e(b,{to:"/",class:"hidden sm:flex"},{default:o(()=>[e(s,{size:"tiny",type:"tertiary",round:"",ghost:"",class:"!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white"},{default:o(()=>t[0]||(t[0]=[a("i",{class:"fas fa-home text-xs"},null,-1),d(" 首页 ",-1)])),_:1,__:[0]})]),_:1}),e(b,{to:"/hot-dramas",class:"hidden sm:flex"},{default:o(()=>[e(s,{size:"tiny",type:"tertiary",round:"",ghost:"",class:"!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white"},{default:o(()=>t[1]||(t[1]=[a("i",{class:"fas fa-film text-xs"},null,-1),d(" 热播剧 ",-1)])),_:1,__:[1]})]),_:1}),e(b,{to:"/monitor",class:"hidden sm:flex"},{default:o(()=>[e(s,{size:"tiny",type:"tertiary",round:"",ghost:"",class:"!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white"},{default:o(()=>t[2]||(t[2]=[a("i",{class:"fas fa-chart-line text-xs"},null,-1),d(" 系统监控 ",-1)])),_:1,__:[2]})]),_:1})])]),t[5]||(t[5]=q(`<div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6 mb-8" data-v-bbb55402><h2 class="text-xl font-semibold text-blue-800 dark:text-blue-200 mb-4 flex items-center" data-v-bbb55402><i class="fas fa-key mr-2" data-v-bbb55402></i> API认证说明 </h2><div class="space-y-3 text-blue-700 dark:text-blue-300" data-v-bbb55402><p data-v-bbb55402><strong data-v-bbb55402>认证方式：</strong>所有API都需要提供API Token进行认证</p><p data-v-bbb55402><strong data-v-bbb55402>请求头方式：</strong><code class="bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded" data-v-bbb55402>X-API-Token: your_token</code></p><p data-v-bbb55402><strong data-v-bbb55402>查询参数方式：</strong><code class="bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded" data-v-bbb55402>?api_token=your_token</code></p><p data-v-bbb55402><strong data-v-bbb55402>获取Token：</strong>请联系管理员在系统配置中设置API Token</p></div></div><div class="space-y-8" data-v-bbb55402><div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden" data-v-bbb55402><div class="bg-purple-600 text-white px-6 py-4" data-v-bbb55402><h3 class="text-xl font-semibold flex items-center" data-v-bbb55402><i class="fas fa-layer-group mr-2" data-v-bbb55402></i> 批量添加资源 </h3><p class="text-purple-100 mt-1" data-v-bbb55402>批量添加多个资源到待处理列表，每个资源可包含多个链接（url为数组），标题和url为必填项</p></div><div class="p-6" data-v-bbb55402><div class="grid grid-cols-1 lg:grid-cols-2 gap-6" data-v-bbb55402><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>请求信息</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><strong data-v-bbb55402>方法：</strong><span class="bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-200 px-2 py-1 rounded" data-v-bbb55402>POST</span></p><p data-v-bbb55402><strong data-v-bbb55402>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>/api/public/resources/batch-add</code></p><p data-v-bbb55402><strong data-v-bbb55402>认证：</strong><span class="text-red-600 dark:text-red-400" data-v-bbb55402>必需</span>（X-API-Token）</p></div></div><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>请求参数</h4><p class="text-xs text-gray-500 dark:text-gray-400 mb-2" data-v-bbb55402>title 和 url 是必填项，其他字段均为选填</p><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>{
  &quot;resources&quot;: [
    {
      &quot;title&quot;: &quot;资源1&quot;,
      &quot;description&quot;: &quot;描述1&quot;,
      &quot;url&quot;: [&quot;链接1&quot;, &quot;链接2&quot;],
      &quot;category&quot;: &quot;分类&quot;,
      &quot;tags&quot;: &quot;标签1,标签2&quot;,
      &quot;img&quot;: &quot;图片链接&quot;,
      &quot;source&quot;: &quot;数据来源&quot;,
      &quot;extra&quot;: &quot;额外信息&quot;
    },
    {
      &quot;title&quot;: &quot;资源2&quot;,
      &quot;url&quot;: [&quot;链接3&quot;],
      &quot;description&quot;: &quot;描述2&quot;
    }
  ]
}</code></pre></div></div></div><div class="mt-6" data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>响应示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>{
  &quot;success&quot;: true,
  &quot;message&quot;: &quot;操作成功&quot;,
  &quot;data&quot;: {
    &quot;created_count&quot;: 2,
    &quot;created_ids&quot;: [123, 124]
  },
  &quot;code&quot;: 200
}</code></pre></div></div></div></div><div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden" data-v-bbb55402><div class="bg-blue-600 text-white px-6 py-4" data-v-bbb55402><h3 class="text-xl font-semibold flex items-center" data-v-bbb55402><i class="fas fa-search mr-2" data-v-bbb55402></i> 资源搜索 </h3><p class="text-blue-100 mt-1" data-v-bbb55402>搜索资源，支持关键词、标签、分类过滤，自动过滤包含违禁词的资源</p></div><div class="p-6" data-v-bbb55402><div class="grid grid-cols-1 lg:grid-cols-2 gap-6" data-v-bbb55402><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>请求信息</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><strong data-v-bbb55402>方法：</strong><span class="bg-blue-100 dark:bg-blue-800 text-blue-800 dark:text-blue-200 px-2 py-1 rounded" data-v-bbb55402>GET</span></p><p data-v-bbb55402><strong data-v-bbb55402>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>/api/public/resources/search</code></p><p data-v-bbb55402><strong data-v-bbb55402>认证：</strong><span class="text-red-600 dark:text-red-400" data-v-bbb55402>必需</span></p></div></div><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>查询参数</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>keyword</code> - 搜索关键词</p><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>tag</code> - 标签过滤</p><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>category</code> - 分类过滤</p><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>page</code> - 页码（默认1）</p><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>page_size</code> - 每页数量（默认20，最大100）</p></div></div></div><div class="mt-6" data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>响应示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>{
  &quot;success&quot;: true,
  &quot;message&quot;: &quot;操作成功&quot;,
  &quot;data&quot;: {
    &quot;list&quot;: [
      {
        &quot;id&quot;: 1,
        &quot;title&quot;: &quot;资源标题&quot;,
        &quot;url&quot;: &quot;资源链接&quot;,
        &quot;description&quot;: &quot;资源描述&quot;,
        &quot;view_count&quot;: 100,
        &quot;created_at&quot;: &quot;2024-12-19 10:00:00&quot;,
        &quot;updated_at&quot;: &quot;2024-12-19 10:00:00&quot;
      }
    ],
    &quot;total&quot;: 50,
    &quot;page&quot;: 1,
    &quot;limit&quot;: 20
  },
  &quot;code&quot;: 200
}</code></pre></div><div class="mt-4 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg" data-v-bbb55402><h5 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2 flex items-center" data-v-bbb55402><i class="fas fa-exclamation-triangle mr-2" data-v-bbb55402></i> 违禁词过滤说明 </h5><p class="text-sm text-yellow-700 dark:text-yellow-300 mb-2" data-v-bbb55402>当搜索结果包含违禁词时，响应会包含额外的过滤信息：</p><pre class="text-xs bg-yellow-100 dark:bg-yellow-800 rounded p-2" data-v-bbb55402><code data-v-bbb55402>{
  &quot;success&quot;: true,
  &quot;message&quot;: &quot;操作成功&quot;,
  &quot;data&quot;: {
    &quot;list&quot;: [...],
    &quot;total&quot;: 45,
    &quot;page&quot;: 1,
    &quot;limit&quot;: 20,
    &quot;forbidden_words_filtered&quot;: true,
    &quot;filtered_forbidden_words&quot;: [&quot;违禁词1&quot;, &quot;违禁词2&quot;],
    &quot;original_total&quot;: 50,
    &quot;filtered_count&quot;: 5
  },
  &quot;code&quot;: 200
}</code></pre></div></div></div></div><div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden" data-v-bbb55402><div class="bg-orange-600 text-white px-6 py-4" data-v-bbb55402><h3 class="text-xl font-semibold flex items-center" data-v-bbb55402><i class="fas fa-film mr-2" data-v-bbb55402></i> 热门剧列表 </h3><p class="text-orange-100 mt-1" data-v-bbb55402>获取热门剧列表，支持分页</p></div><div class="p-6" data-v-bbb55402><div class="grid grid-cols-1 lg:grid-cols-2 gap-6" data-v-bbb55402><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>请求信息</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><strong data-v-bbb55402>方法：</strong><span class="bg-orange-100 dark:bg-orange-800 text-orange-800 dark:text-orange-200 px-2 py-1 rounded" data-v-bbb55402>GET</span></p><p data-v-bbb55402><strong data-v-bbb55402>路径：</strong><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>/api/public/hot-dramas</code></p><p data-v-bbb55402><strong data-v-bbb55402>认证：</strong><span class="text-red-600 dark:text-red-400" data-v-bbb55402>必需</span></p></div></div><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>查询参数</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>page</code> - 页码（默认1）</p><p data-v-bbb55402><code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded" data-v-bbb55402>page_size</code> - 每页数量（默认20，最大100）</p></div></div></div><div class="mt-6" data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>响应示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>{
  &quot;success&quot;: true,
  &quot;message&quot;: &quot;操作成功&quot;,
  &quot;data&quot;: {
    &quot;hot_dramas&quot;: [
      {
        &quot;id&quot;: 1,
        &quot;title&quot;: &quot;剧名&quot;,
        &quot;description&quot;: &quot;剧集描述&quot;,
        &quot;img&quot;: &quot;封面图片&quot;,
        &quot;url&quot;: &quot;详情链接&quot;,
        &quot;rating&quot;: 8.5,
        &quot;year&quot;: &quot;2024&quot;,
        &quot;region&quot;: &quot;中国大陆&quot;,
        &quot;genres&quot;: &quot;剧情,悬疑&quot;,
        &quot;category&quot;: &quot;电视剧&quot;,
        &quot;created_at&quot;: &quot;2024-12-19 10:00:00&quot;,
        &quot;updated_at&quot;: &quot;2024-12-19 10:00:00&quot;
      }
    ],
    &quot;total&quot;: 20,
    &quot;page&quot;: 1,
    &quot;page_size&quot;: 20
  },
  &quot;code&quot;: 200
}</code></pre></div></div></div></div></div><div class="mt-12 bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden" data-v-bbb55402><div class="bg-red-600 text-white px-6 py-4" data-v-bbb55402><h3 class="text-xl font-semibold flex items-center" data-v-bbb55402><i class="fas fa-exclamation-triangle mr-2" data-v-bbb55402></i> 错误码说明 </h3></div><div class="p-6" data-v-bbb55402><div class="grid grid-cols-1 md:grid-cols-2 gap-6" data-v-bbb55402><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>HTTP状态码</h4><div class="space-y-2 text-sm" data-v-bbb55402><p data-v-bbb55402><span class="bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-200 px-2 py-1 rounded" data-v-bbb55402>200</span> - 请求成功</p><p data-v-bbb55402><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded" data-v-bbb55402>400</span> - 请求参数错误</p><p data-v-bbb55402><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded" data-v-bbb55402>401</span> - 认证失败（Token无效或缺失）</p><p data-v-bbb55402><span class="bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-200 px-2 py-1 rounded" data-v-bbb55402>500</span> - 服务器内部错误</p><p data-v-bbb55402><span class="bg-yellow-100 dark:bg-yellow-800 text-yellow-800 dark:text-yellow-200 px-2 py-1 rounded" data-v-bbb55402>503</span> - 系统维护中或API Token未配置</p></div></div><div data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>响应格式</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>{
  &quot;success&quot;: true/false,
  &quot;message&quot;: &quot;响应消息&quot;,
  &quot;data&quot;: {}, // 响应数据
  &quot;code&quot;: 200 // 状态码
}</code></pre></div></div></div></div></div><div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden" data-v-bbb55402><div class="bg-indigo-600 text-white px-6 py-4" data-v-bbb55402><h3 class="text-xl font-semibold flex items-center" data-v-bbb55402><i class="fas fa-code mr-2" data-v-bbb55402></i> 使用示例 </h3></div><div class="p-6" data-v-bbb55402><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>cURL示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402># 设置API Token
API_TOKEN=&quot;your_api_token_here&quot;

# 批量添加资源
curl -X POST &quot;http://localhost:8080/api/public/resources/batch-add&quot; \\
  -H &quot;Content-Type: application/json&quot; \\
  -H &quot;X-API-Token: $API_TOKEN&quot; \\
  -d &#39;{
    &quot;resources&quot;: [
      { &quot;title&quot;: &quot;测试资源1&quot;, &quot;url&quot;: [&quot;https://example.com/resource1&quot;], &quot;description&quot;: &quot;描述1&quot; },
      { &quot;title&quot;: &quot;测试资源2&quot;, &quot;url&quot;: [&quot;https://example.com/resource2&quot;, &quot;https://example.com/resource3&quot;], &quot;description&quot;: &quot;描述2&quot; }
    ]
  }&#39;

# 搜索资源
curl -X GET &quot;http://localhost:8080/api/public/resources/search?keyword=测试&quot; \\
  -H &quot;X-API-Token: $API_TOKEN&quot;

# 获取热门剧
curl -X GET &quot;http://localhost:8080/api/public/hot-dramas?page=1&amp;page_size=5&quot; \\
  -H &quot;X-API-Token: $API_TOKEN&quot;</code></pre></div><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>JavaScript fetch 示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>// 资源搜索
fetch(&#39;/api/public/resources/search?keyword=测试&#39;, { 
  headers: { &#39;X-API-Token&#39;: &#39;your_token&#39; } 
})
  .then(res =&gt; res.json())
  .then(res =&gt; {
    if (res.success) {
      const list = res.data.list // 资源列表
      const total = res.data.total
      console.log(&#39;搜索结果:&#39;, list)
    } else {
      console.error(&#39;搜索失败:&#39;, res.message)
    }
  })

// 批量添加资源
fetch(&#39;/api/public/resources/batch-add&#39;, {
  method: &#39;POST&#39;,
  headers: { 
    &#39;Content-Type&#39;: &#39;application/json&#39;, 
    &#39;X-API-Token&#39;: &#39;your_token&#39; 
  },
  body: JSON.stringify({
    resources: [
      { title: &#39;测试资源1&#39;, url: [&#39;https://example.com/resource1&#39;], description: &#39;描述1&#39; },
      { title: &#39;测试资源2&#39;, url: [&#39;https://example.com/resource2&#39;], description: &#39;描述2&#39; }
    ]
  })
})
  .then(res =&gt; res.json())
  .then(res =&gt; {
    if (res.success) {
      console.log(&#39;添加成功，ID:&#39;, res.data.created_ids)
    } else {
      console.error(&#39;添加失败:&#39;, res.message)
    }
  })

// 获取热门剧
fetch(&#39;/api/public/hot-dramas?page=1&amp;page_size=10&#39;, {
  headers: { &#39;X-API-Token&#39;: &#39;your_token&#39; }
})
  .then(res =&gt; res.json())
  .then(res =&gt; {
    if (res.success) {
      const dramas = res.data.hot_dramas
      console.log(&#39;热门剧:&#39;, dramas)
    } else {
      console.error(&#39;获取失败:&#39;, res.message)
    }
  })</code></pre></div><h4 class="font-semibold text-gray-900 dark:text-white mb-3" data-v-bbb55402>Python requests 示例</h4><div class="bg-gray-50 dark:bg-gray-700 rounded p-4" data-v-bbb55402><pre class="text-sm overflow-x-auto" data-v-bbb55402><code data-v-bbb55402>import requests

API_TOKEN = &#39;your_api_token_here&#39;
BASE_URL = &#39;http://localhost:8080/api&#39;

headers = {
    &#39;X-API-Token&#39;: API_TOKEN,
    &#39;Content-Type&#39;: &#39;application/json&#39;
}

# 搜索资源
def search_resources(keyword, page=1, page_size=20):
    params = {
        &#39;keyword&#39;: keyword,
        &#39;page&#39;: page,
        &#39;page_size&#39;: page_size
    }
    response = requests.get(
        f&#39;{BASE_URL}/public/resources/search&#39;,
        headers={&#39;X-API-Token&#39;: API_TOKEN},
        params=params
    )
    return response.json()

# 批量添加资源
def batch_add_resources(resources):
    data = {&#39;resources&#39;: resources}
    response = requests.post(
        f&#39;{BASE_URL}/public/resources/batch-add&#39;,
        headers=headers,
        json=data
    )
    return response.json()

# 获取热门剧
def get_hot_dramas(page=1, page_size=20):
    params = {
        &#39;page&#39;: page,
        &#39;page_size&#39;: page_size
    }
    response = requests.get(
        f&#39;{BASE_URL}/public/hot-dramas&#39;,
        headers={&#39;X-API-Token&#39;: API_TOKEN},
        params=params
    )
    return response.json()

# 使用示例
if __name__ == &#39;__main__&#39;:
    # 搜索资源
    result = search_resources(&#39;测试&#39;)
    if result[&#39;success&#39;]:
        print(&#39;搜索结果:&#39;, result[&#39;data&#39;][&#39;list&#39;])
    
    # 批量添加资源
    resources = [
        {&#39;title&#39;: &#39;测试资源1&#39;, &#39;url&#39;: [&#39;https://example.com/resource1&#39;]},
        {&#39;title&#39;: &#39;测试资源2&#39;, &#39;url&#39;: [&#39;https://example.com/resource2&#39;]}
    ]
    result = batch_add_resources(resources)
    if result[&#39;success&#39;]:
        print(&#39;添加成功，ID:&#39;, result[&#39;data&#39;][&#39;created_ids&#39;])
    
    # 获取热门剧
    result = get_hot_dramas()
    if result[&#39;success&#39;]:
        print(&#39;热门剧:&#39;, result[&#39;data&#39;][&#39;hot_dramas&#39;])</code></pre></div></div></div>`,4))])]),e(u)])}}},K=c(w,[["__scopeId","data-v-bbb55402"]]);export{K as default};
