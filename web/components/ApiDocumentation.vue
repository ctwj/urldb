<template>
  <div class="space-y-4">
    <div class="text-gray-700 dark:text-gray-300 text-sm">
      <p>你可以通过API批量添加资源：</p>
      <pre class="bg-gray-100 dark:bg-gray-800 p-3 rounded text-xs overflow-x-auto mt-2">
POST /api/resources/batch
Content-Type: application/json
Body:
[
  { "title": "资源A", "url": "https://a.com", "file_type": "pan", ... },
  { "title": "资源B", "url": "https://b.com", ... }
]
      </pre>
      <p>参数说明：<br/>
        title: 标题<br/>
        url: 资源链接<br/>
        file_type: 类型（pan/link/other）<br/>
        tags: 标签数组（可选）<br/>
        description: 描述（可选）<br/>
        ... 其他字段参考文档
      </p>
    </div>
    
    <div class="flex justify-end space-x-3 pt-4">
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