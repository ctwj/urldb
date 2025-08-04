<template>
  <div>
    <div class="flex justify-between mb-4">
      <div class="mb-4 flex-1 w-1">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">输入格式说明：</label>
        <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded text-sm text-gray-600 dark:text-gray-300 mb-4">
          <p class="mb-2"><strong>格式要求：</strong>标题和URL为一组，标题必填, 同一标题URL支持多行</p>
          <pre class="bg-white dark:bg-gray-800 p-2 rounded border text-xs">
电影1
https://pan.baidu.com/s/123456
https://pan.quark.com/s/123456
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
      <div class="mb-4 flex-1 w-1">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">资源内容：</label>
        <n-input v-model:value="batchInput" type="textarea"
          :autosize="{ minRows: 10, maxRows: 15 }"
          placeholder="请输入资源内容，格式：标题和URL为一组..." />
      </div>
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

  // 首行必须为标题
  if (/^https?:\/\//i.test(lines[0])) {
    // 你可以用 alert、ElMessage 或其它方式提示
    alert('首行必须为标题，不能为链接！')
    return
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

    let currentTitle = ''
    let currentUrls = []

    for (const line of lines) {
      // 判断是否为 url（以 http/https 开头）
      if (/^https?:\/\//i.test(line)) {
        currentUrls.push(line)
      } else {
        // 新标题，先保存上一个
        if (currentTitle && currentUrls.length) {
          resources.push({
            title: currentTitle,
            url: currentUrls.slice()
          })
        }
        currentTitle = line
        currentUrls = []
      }
    }
    // 处理最后一组
    if (currentTitle && currentUrls.length) {
      resources.push({
        title: currentTitle,
        url: currentUrls.slice()
      })
    }

    // 调用API添加资源
    const res: any = await readyResourceApi.batchCreateReadyResources({resources})
    console.log(res)
    emit('success', res.message)
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