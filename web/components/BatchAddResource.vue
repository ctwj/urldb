<template>
  <div>
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">输入格式说明：</label>
      <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded text-sm text-gray-600 dark:text-gray-300 mb-4">
        <p class="mb-2"><strong>格式要求：</strong>标题和URL两行为一组，标题为必填项</p>
        <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
电影标题1
https://pan.baidu.com/s/123456
电影标题2
https://pan.baidu.com/s/789012
电视剧标题3
https://pan.quark.cn/s/345678</pre>
        <p class="mt-2 text-xs text-red-600 dark:text-red-400">
          <i class="fas fa-exclamation-triangle mr-1"></i>
          注意：标题为必填项，不能为空
        </p>
      </div>
    </div>
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">资源内容：</label>
      <textarea
        v-model="batchInput"
        rows="15"
        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-900 dark:text-gray-100"
        placeholder="请输入资源内容，格式：标题和URL两行为一组..."
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

// 验证输入格式
const validateInput = () => {
  if (!batchInput.value.trim()) {
    throw new Error('请输入资源内容')
  }
  
  const lines = batchInput.value.split(/\r?\n/).map(line => line.trim()).filter(Boolean)
  
  if (lines.length === 0) {
    throw new Error('请输入有效的资源内容')
  }
  
  // 检查是否为偶数行（标题+URL为一组）
  if (lines.length % 2 !== 0) {
    throw new Error('资源格式错误：标题和URL必须成对出现，请检查是否缺少标题或URL')
  }
  
  // 检查每组的标题是否为空
  for (let i = 0; i < lines.length; i += 2) {
    const title = lines[i]
    const url = lines[i + 1]
    
    if (!title) {
      throw new Error(`第${i + 1}行标题不能为空`)
    }
    
    if (!url) {
      throw new Error(`第${i + 2}行URL不能为空`)
    }
    
    // 验证URL格式
    try {
      new URL(url)
    } catch {
      throw new Error(`第${i + 2}行URL格式无效: ${url}`)
    }
  }
}

// 批量添加提交
const handleSubmit = async () => {
  loading.value = true
  try {
    validateInput()
    
    // 解析输入内容
    const lines = batchInput.value.split(/\r?\n/).map(line => line.trim()).filter(Boolean)
    const resources = []
    
    for (let i = 0; i < lines.length; i += 2) {
      const title = lines[i]
      const url = lines[i + 1]
      
      resources.push({
        title: title,
        url: url,
        source: '批量添加'
      })
    }
    
    // 调用API添加资源
    const res: any = await readyResourceApi.batchCreateReadyResources(resources)
    emit('success', `成功添加 ${res.count || resources.length} 个资源，资源已进入待处理列表，处理完成后会自动入库`)
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