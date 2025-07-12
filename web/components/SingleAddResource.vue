<template>
  <div class="space-y-6">
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
    
    <div class="flex justify-end space-x-3 pt-4">
      <button type="button" @click="$emit('cancel')" class="btn-secondary">取消</button>
      <button type="button" @click="handleSubmit" class="btn-primary" :disabled="loading">
        {{ loading ? '保存中...' : '添加' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useResourceStore } from '~/stores/resource'

const emit = defineEmits(['success', 'error', 'cancel'])

const store = useResourceStore()
const loading = ref(false)
const newTag = ref('')

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

// 单个添加提交
const handleSubmit = async () => {
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
    emit('success', '资源已进入待处理列表，处理完成后会自动入库')
    // 清空表单
    form.value.title = ''
    form.value.description = ''
    form.value.url = ''
    form.value.tags = []
    form.value.file_type = ''
  } catch (e: any) {
    emit('error', e.message || '添加失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.input-field {
  @apply w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500;
}

.btn-primary {
  @apply px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors disabled:opacity-50;
}

.btn-secondary {
  @apply px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-md transition-colors;
}
</style> 