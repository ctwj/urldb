<template>
  <div class="space-y-6">
    <!-- 标题 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        标题 <span class="text-red-500">*</span>
      </label>
      <n-input 
        v-model="form.title" 
        placeholder="输入资源标题（必填）" 
        required
      />
    </div>

    <!-- 描述 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        描述 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <n-input 
        v-model="form.description" 
        type="textarea"
        placeholder="输入资源描述，如：剧情简介、文件大小、清晰度等"
      />
    </div>

    <!-- URL -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        URL <span class="text-red-500">*</span>
      </label>
      <n-input 
        v-model="form.url" 
        type="textarea"
        placeholder="请输入资源链接，支持多行，每行一个链接"
        required
      />
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
        支持百度网盘、阿里云盘、夸克网盘等链接，每行一个链接
      </p>
    </div>

    <!-- 分类 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        分类 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <n-input 
        v-model="form.category" 
        placeholder="如：电影、电视剧、动漫、音乐等" 
      />
    </div>

    <!-- 标签 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        标签 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <div class="flex flex-wrap gap-2 mb-2">
        <span 
          v-for="tag in form.tags" 
          :key="tag" 
          class="bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 px-2 py-1 rounded text-xs flex items-center"
        >
          {{ tag }}
          <button 
            type="button" 
            class="ml-1 text-xs hover:text-red-500" 
            @click="removeTag(tag)"
          >
            ×
          </button>
        </span>
      </div>
      <n-input 
        v-model="newTag" 
        @keyup.enter.prevent="addTag" 
        placeholder="输入标签后回车添加，多个标签用逗号分隔" 
      />
    </div>

    <!-- 封面图片 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        封面图片 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <n-input 
        v-model="form.img" 
        placeholder="封面图片链接" 
      />
    </div>

    <!-- 数据来源 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        数据来源 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <n-input 
        v-model="form.source" 
        placeholder="如：手动添加、API导入、爬虫等" 
      />
    </div>

    <!-- 额外数据 -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        额外数据 <span class="text-gray-400 text-xs">(可选)</span>
      </label>
      <n-input 
        v-model="form.extra" 
        type="textarea"
        placeholder="JSON格式的额外数据，如：{'size': '2GB', 'quality': '1080p'}"
      />
    </div>
    
    <!-- 按钮区域 -->
    <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
      <button type="button" @click="$emit('cancel')" class="btn-secondary">
        取消
      </button>
      <button type="button" @click="handleSubmit" class="btn-primary" :disabled="loading">
        {{ loading ? '保存中...' : '添加资源' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useReadyResourceApi } from '~/composables/useApi'

const emit = defineEmits(['success', 'error', 'cancel'])

const readyResourceApi = useReadyResourceApi()
const loading = ref(false)
const newTag = ref('')

// 根据 ready_resource 表字段定义表单
const form = ref({
  title: '',
  description: '',
  url: '',
  category: '',
  tags: [] as string[],
  img: '',
  source: '',
  extra: '',
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

// 验证表单
const validateForm = () => {
  if (!form.value.title.trim()) {
    throw new Error('标题不能为空')
  }
  if (!form.value.url.trim()) {
    throw new Error('请输入至少一个URL')
  }
  
  // 验证URL格式
  const urls = form.value.url.split(/\r?\n/).map(l => l.trim()).filter(Boolean)
  for (const url of urls) {
    try {
      new URL(url)
    } catch {
      throw new Error(`无效的URL格式: ${url}`)
    }
  }
}

// 清空表单
const clearForm = () => {
  form.value = {
    title: '',
    description: '',
    url: '',
    category: '',
    tags: [],
    img: '',
    source: '',
    extra: '',
  }
  newTag.value = ''
}

// 单个添加提交 - 更新为使用批量添加方法
const handleSubmit = async () => {
  loading.value = true
  try {
    validateForm()
    
    // 多行链接处理
    const urls = form.value.url.split(/\r?\n/).map(l => l.trim()).filter(Boolean)
    
    // 使用批量添加方法，将多个URL作为一个资源的多个链接
    const resourceData = {
      resources: [{
        title: form.value.title || undefined, // 后端期望 *string 类型
        description: form.value.description || '',
        url: urls, // 现在 url 是一个数组
        category: form.value.category || '',
        tags: form.value.tags.join(','), // 转换为逗号分隔的字符串
        img: form.value.img || '',
        source: form.value.source || '手动添加',
        extra: form.value.extra || '',
      }]
    }
    
    const response = await readyResourceApi.batchCreateReadyResources(resourceData)
    
    emit('success', `成功添加资源，包含 ${urls.length} 个链接`)
    clearForm()
  } catch (e: any) {
    emit('error', e.message || '添加失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.input-field {
  @apply w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors;
}

.btn-primary {
  @apply px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-md transition-colors;
}
</style> 