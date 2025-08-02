<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-900 rounded-lg shadow-xl max-w-2xl w-full mx-4" style="height:600px;">
      <div class="p-6 h-full flex flex-col text-gray-900 dark:text-gray-100">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100">
            添加资源
          </h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Tab 切换 -->
        <div class="flex mb-6 border-b flex-shrink-0">
          <button
            v-for="tab in tabs"
            :key="tab.value"
            :class="['px-4 py-2 -mb-px border-b-2', mode === tab.value ? 'border-blue-500 text-blue-600 font-bold' : 'border-transparent text-gray-500']"
            @click="mode = tab.value"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- 内容区域 -->
        <div class="flex-1 overflow-y-auto">
          <!-- 批量添加 -->
          <div v-if="mode === 'batch'">
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
          </div>

          <!-- 单个添加 -->
          <div v-else-if="mode === 'single'" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">标题</label>
              <input v-model="form.title" class="input-field dark:bg-gray-900 dark:text-gray-100 dark:border-gray-700" placeholder="输入标题" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">描述</label>
              <textarea v-model="form.description" rows="3" class="input-field dark:bg-gray-900 dark:text-gray-100 dark:border-gray-700" placeholder="输入资源描述"></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">类型</label>
              <select v-model="form.file_type" class="input-field dark:bg-gray-900 dark:text-gray-100 dark:border-gray-700">
                <option value="">选择类型</option>
                <option value="pan">网盘</option>
                <option value="link">直链</option>
                <option value="other">其他</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">标签</label>
              <div class="flex flex-wrap gap-2 mb-2">
                <span v-for="tag in form.tags" :key="tag" class="bg-blue-100 text-blue-700 px-2 py-1 rounded text-xs flex items-center">
                  {{ tag }}
                  <button type="button" class="ml-1 text-xs" @click="removeTag(tag)">×</button>
                </span>
              </div>
              <input v-model="newTag" @keyup.enter.prevent="addTag" class="input-field dark:bg-gray-900 dark:text-gray-100 dark:border-gray-700" placeholder="输入标签后回车添加" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">链接（可多行，每行一个链接）</label>
              <textarea v-model="form.url" rows="3" class="input-field dark:bg-gray-900 dark:text-gray-100 dark:border-gray-700" placeholder="https://a.com&#10;https://b.com"></textarea>
            </div>
          </div>

          <!-- API说明 -->
          <div v-else class="space-y-4">
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
          </div>
        </div>

        <!-- 按钮区域 -->
        <div class="flex-shrink-0 pt-4 border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900/90 sticky bottom-0 left-0 w-full flex justify-end space-x-3 z-10 backdrop-blur">
          <template v-if="mode === 'batch'">
            <button type="button" @click="$emit('close')" class="btn-secondary">取消</button>
            <button type="button" @click="handleBatchSubmit" class="btn-primary" :disabled="loading">
              {{ loading ? '保存中...' : '批量添加' }}
            </button>
          </template>
          <template v-else-if="mode === 'single'">
            <button type="button" @click="$emit('close')" class="btn-secondary">取消</button>
            <button type="button" @click="handleSingleSubmit" class="btn-primary" :disabled="loading">
              {{ loading ? '保存中...' : '添加' }}
            </button>
          </template>
          <template v-else>
            <button type="button" @click="$emit('close')" class="btn-secondary">关闭</button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useResourceStore } from '~/stores/resource'
import { storeToRefs } from 'pinia'
import { useReadyResourceApi } from '~/composables/useApi'

const notification = useNotification()
const store = useResourceStore()
const { categories } = storeToRefs(store)

const props = defineProps<{ resource?: any }>()
const emit = defineEmits(['close', 'save'])

const loading = ref(false)
const newTag = ref('')

const tabs = [
  { label: '批量添加', value: 'batch' },
  { label: '单个添加', value: 'single' },
  { label: 'API说明', value: 'api' },
]
const mode = ref('batch')

// 批量添加
const batchInput = ref('')

// 单个添加表单
const form = ref({
  title: '',
  description: '',
  url: '', // 多行
  category_id: '',
  tags: [] as string[],
  file_path: '',
  file_type: '',
  file_size: 0,
  is_public: true,
})

const readyResourceApi = useReadyResourceApi()

onMounted(() => {
  if (props.resource) {
    form.value = {
      title: props.resource.title || '',
      description: props.resource.description || '',
      url: props.resource.url || '',
      category_id: props.resource.category_id || '',
      tags: [...(props.resource.tags || [])],
      file_path: props.resource.file_path || '',
      file_type: props.resource.file_type || '',
      file_size: props.resource.file_size || 0,
      is_public: props.resource.is_public !== false,
    }
  }
})

const addTag = () => {
  const tag = newTag.value.trim()
  if (tag && !form.value.tags.includes(tag)) {
    form.value.tags.push(tag)
    newTag.value = ''
  }
}

const removeTag = (tag: string) => {
  const index = form.value.tags.indexOf(tag)
  if (index > -1) {
    form.value.tags.splice(index, 1)
  }
}

// 批量添加提交
const handleBatchSubmit = async () => {
  loading.value = true
  try {
    if (!batchInput.value.trim()) throw new Error('请输入资源内容')
    const res: any = await readyResourceApi.createReadyResourcesFromText(batchInput.value)
    notification.success({
      content: `成功添加 ${res.count || 0} 个资源，资源已进入待处理列表，处理完成后会自动入库`
    })
    batchInput.value = ''
  } catch (e: any) {
    notification.error({
      content: e.message || '批量添加失败'
    })
  } finally {
    loading.value = false
  }
}

// 单个添加提交
const handleSingleSubmit = async () => {
  loading.value = true
  try {
    // 多行链接
    const urls = form.value.url.split(/\r?\n/).map(l => l.trim()).filter(Boolean)
    if (!urls.length) throw new Error('请输入至少一个链接')
    for (const url of urls) {
      await store.createResource({
        ...form.value,
        url,
        tags: [...form.value.tags],
      })
    }
    notification.success({
      content:  '资源已进入待处理列表，处理完成后会自动入库'
    })
    // 清空表单
    form.value.title = ''
    form.value.description = ''
    form.value.url = ''
    form.value.tags = []
    form.value.file_type = ''
  } catch (e: any) {
    notification.error({
      content:  e.message || '添加失败'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 