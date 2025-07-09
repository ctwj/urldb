<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold text-gray-900">
            {{ resource ? '编辑资源' : '添加资源' }}
          </h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- 标题 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">标题 *</label>
            <input
              v-model="form.title"
              type="text"
              required
              class="input-field"
              placeholder="输入资源标题"
            />
          </div>

          <!-- 描述 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">描述</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="input-field"
              placeholder="输入资源描述"
            ></textarea>
          </div>

          <!-- URL -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">URL</label>
            <input
              v-model="form.url"
              type="url"
              class="input-field"
              placeholder="https://example.com"
            />
          </div>

          <!-- 分类 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">分类</label>
            <select v-model="form.category_id" class="input-field">
              <option value="">选择分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>

          <!-- 标签 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">标签</label>
            <div class="flex flex-wrap gap-2 mb-2">
              <span
                v-for="tag in form.tags"
                :key="tag"
                class="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full flex items-center"
              >
                {{ tag }}
                <button
                  @click="removeTag(tag)"
                  type="button"
                  class="ml-2 text-blue-600 hover:text-blue-800"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </span>
            </div>
            <div class="flex gap-2">
              <input
                v-model="newTag"
                @keyup.enter="addTag"
                type="text"
                class="input-field flex-1"
                placeholder="输入标签后按回车"
              />
              <button @click="addTag" type="button" class="btn-secondary">
                添加
              </button>
            </div>
          </div>

          <!-- 文件信息 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">文件路径</label>
              <input
                v-model="form.file_path"
                type="text"
                class="input-field"
                placeholder="/path/to/file"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">文件类型</label>
              <input
                v-model="form.file_type"
                type="text"
                class="input-field"
                placeholder="pdf, doc, mp4..."
              />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">文件大小 (字节)</label>
              <input
                v-model.number="form.file_size"
                type="number"
                class="input-field"
                placeholder="1024"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">是否公开</label>
              <div class="flex items-center mt-2">
                <input
                  v-model="form.is_public"
                  type="checkbox"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label class="ml-2 text-sm text-gray-700">公开显示</label>
              </div>
            </div>
          </div>

          <!-- 按钮 -->
          <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
            <button
              type="button"
              @click="$emit('close')"
              class="btn-secondary"
            >
              取消
            </button>
            <button
              type="submit"
              class="btn-primary"
              :disabled="loading"
            >
              {{ loading ? '保存中...' : (resource ? '更新' : '创建') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const store = useResourceStore()
const { categories } = storeToRefs(store)

const props = defineProps<{
  resource?: any
}>()

const emit = defineEmits(['close', 'save'])

const loading = ref(false)
const newTag = ref('')

const form = ref({
  title: '',
  description: '',
  url: '',
  category_id: '',
  tags: [] as string[],
  file_path: '',
  file_type: '',
  file_size: 0,
  is_public: true,
})

// 初始化表单
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

// 添加标签
const addTag = () => {
  const tag = newTag.value.trim()
  if (tag && !form.value.tags.includes(tag)) {
    form.value.tags.push(tag)
    newTag.value = ''
  }
}

// 移除标签
const removeTag = (tag: string) => {
  const index = form.value.tags.indexOf(tag)
  if (index > -1) {
    form.value.tags.splice(index, 1)
  }
}

// 提交表单
const handleSubmit = async () => {
  loading.value = true
  try {
    const data = {
      ...form.value,
      category_id: form.value.category_id ? parseInt(form.value.category_id) : null,
    }
    emit('save', data)
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    loading.value = false
  }
}
</script> 