<template>
  <div class="space-y-4">
    <div class="text-gray-700 dark:text-gray-300 text-sm">
      <p class="mb-4">你可以通过API批量添加资源到待处理列表：</p>
      
      <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg mb-4">
        <h4 class="font-semibold mb-2">单个资源添加：</h4>
        <pre class="bg-white dark:bg-gray-900 p-3 rounded text-xs overflow-x-auto">
POST /api/ready-resources
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN

{
  "title": "资源标题",
  "url": "https://pan.baidu.com/s/123456",
  "category": "电影",
  "tags": "动作,科幻,2024",
  "img": "https://example.com/cover.jpg",
  "source": "手动添加",
  "extra": "{\"size\": \"2GB\", \"quality\": \"1080p\"}"
}
        </pre>
      </div>

      <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg mb-4">
        <h4 class="font-semibold mb-2">批量资源添加：</h4>
        <pre class="bg-white dark:bg-gray-900 p-3 rounded text-xs overflow-x-auto">
POST /api/ready-resources/batch
Content-Type: application/json
Authorization: Bearer YOUR_TOKEN

{
  "resources": [
    {
      "title": "资源A",
      "url": "https://pan.baidu.com/s/123456",
      "category": "电影",
      "tags": "动作,科幻",
      "img": "https://example.com/cover1.jpg",
      "source": "API导入",
      "extra": "{\"size\": \"2GB\"}"
    },
    {
      "title": "资源B", 
      "url": "https://pan.baidu.com/s/789012",
      "category": "电视剧",
      "tags": "悬疑,犯罪",
      "source": "API导入"
    }
  ]
}
        </pre>
      </div>

      <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg mb-4">
        <h4 class="font-semibold mb-2">从文本批量添加：</h4>
        <pre class="bg-white dark:bg-gray-900 p-3 rounded text-xs overflow-x-auto">
POST /api/ready-resources/text
Content-Type: multipart/form-data
Authorization: Bearer YOUR_TOKEN

Form Data:
text: |
  电影标题1
  https://pan.baidu.com/s/123456
  电影标题2
  https://pan.baidu.com/s/789012
        </pre>
      </div>

      <div class="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg">
        <h4 class="font-semibold mb-2 text-blue-800 dark:text-blue-200">字段说明：</h4>
        <ul class="space-y-1 text-xs">
          <li><strong>title</strong>: 资源标题（可选，留空则自动从URL提取）</li>
          <li><strong>url</strong>: 资源链接（必填，支持百度网盘、阿里云盘等）</li>
          <li><strong>category</strong>: 资源分类（可选，如：电影、电视剧、动漫等）</li>
          <li><strong>tags</strong>: 资源标签（可选，多个标签用逗号分隔）</li>
          <li><strong>img</strong>: 封面图片链接（可选）</li>
          <li><strong>source</strong>: 数据来源（可选，如：手动添加、API导入、爬虫等）</li>
          <li><strong>extra</strong>: 额外数据（可选，JSON格式字符串）</li>
        </ul>
      </div>
    </div>
    
    <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
      <button type="button" @click="$emit('cancel')" class="btn-secondary">关闭</button>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits(['cancel'])
</script>

<style scoped>
.btn-secondary {
  @apply px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-md transition-colors;
}
</style> 