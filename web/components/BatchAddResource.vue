<template>
  <div>
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">输入格式说明：</label>
      <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded text-sm text-gray-600 dark:text-gray-300 mb-4">
        <p class="mb-2"><strong>格式1：</strong>标题和URL两行一组</p>
        <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
电影标题1
https://pan.baidu.com/s/123456
电影标题2
https://pan.baidu.com/s/789012</pre>
        <p class="mt-2 mb-2"><strong>格式2：</strong>只有URL，系统自动判断</p>
        <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
https://pan.baidu.com/s/123456
https://pan.baidu.com/s/789012
https://pan.baidu.com/s/345678</pre>
      </div>
    </div>
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">资源内容：</label>
      <textarea
        v-model="batchInput"
        rows="15"
        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
        placeholder="请输入资源内容，支持两种格式..."
      ></textarea>
    </div>
    
    <div class="flex justify-end space-x-3 pt-4">
      <button type="button" @click="$emit('cancel')" class="btn-secondary">取消</button>
      <button type="button" @click="handleSubmit" class="btn-primary" :disabled="loading">
        {{ loading ? '保存中...' : '批量添加' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useReadyResourceApi } from '~/composables/useApi'

const emit = defineEmits(['success', 'error', 'cancel'])

const loading = ref(false)
const batchInput = ref('')

const readyResourceApi = useReadyResourceApi()

// 批量添加提交
const handleSubmit = async () => {
  loading.value = true
  try {
    if (!batchInput.value.trim()) throw new Error('请输入资源内容')
    const res: any = await readyResourceApi.createReadyResourcesFromText(batchInput.value)
    emit('success', `成功添加 ${res.count || 0} 个资源，资源已进入待处理列表，处理完成后会自动入库`)
    batchInput.value = ''
  } catch (e: any) {
    emit('error', e.message || '批量添加失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.btn-primary {
  @apply px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors disabled:opacity-50;
}

.btn-secondary {
  @apply px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-md transition-colors;
}
</style> 