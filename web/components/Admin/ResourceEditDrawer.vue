<template>
  <n-drawer v-model:show="visible" :width="800" placement="right" :trap-focus="false" :block-scroll="true">
    <n-drawer-content :title="props.resource ? `编辑资源 - ${props.resource.title}` : '编辑资源'" closable>
      <!-- 初始化加载状态 -->
      <div v-if="initializing" class="flex justify-center items-center h-64">
        <n-spin size="large" />
      </div>

      <!-- 表单内容 -->
      <n-form v-else ref="editFormRef" :model="editForm" :rules="editRules" label-placement="left" label-width="100">
        <n-form-item label="标题" path="title">
          <n-input
            v-model:value="editForm.title"
            placeholder="请输入资源标题"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="editForm.description"
            type="textarea"
            placeholder="请输入资源描述"
            :autosize="{ minRows: 3, maxRows: 6 }"
            class="w-full"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="资源链接" path="url">
          <n-input
            v-model:value="editForm.url"
            placeholder="请输入资源链接"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="editForm.category_id"
            :options="categoryOptions"
            placeholder="请选择分类"
            clearable
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="平台" path="pan_id">
          <n-select
            v-model:value="editForm.pan_id"
            :options="platformOptions"
            placeholder="请选择平台"
            clearable
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="标签" path="tag_ids">
          <n-select
            key="tag-select"
            v-model:value="editForm.tag_ids"
            :options="tagOptions"
            :loading="tagLoading"
            :filterable="true"
            :remote="true"
            :clearable="true"
            :fallback-to-options="false"
            placeholder="请选择标签，支持搜索"
            multiple
            :disabled="initializing"
            @search="handleTagSearch"
            @scroll="handleTagScroll"
          />
          <div v-if="tagLoading" class="text-sm text-gray-500 mt-1">
            正在加载标签...
          </div>
        </n-form-item>

        <n-form-item label="作者" path="author">
          <n-input
            v-model:value="editForm.author"
            placeholder="请输入作者"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="文件大小" path="file_size">
          <n-input
            v-model:value="editForm.file_size"
            placeholder="如：2.5GB"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="封面图片" path="cover">
          <n-input
            v-model:value="editForm.cover"
            placeholder="请输入封面图片URL"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="转存链接" path="save_url">
          <n-input
            v-model:value="editForm.save_url"
            placeholder="请输入转存链接"
            :disabled="initializing"
          />
        </n-form-item>

        <n-form-item label="是否有效" path="is_valid">
          <n-switch v-model:value="editForm.is_valid" :disabled="initializing" />
        </n-form-item>

        <n-form-item label="是否公开" path="is_public">
          <n-switch v-model:value="editForm.is_public" :disabled="initializing" />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="handleClose">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            保存
          </n-button>
        </div>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { useResourceApi, useCategoryApi, useTagApi, usePanApi } from '~/composables/useApi'
import { useMessage, useNotification } from 'naive-ui'

interface Resource {
  id: number
  title: string
  description?: string
  url: string
  category_id?: number
  pan_id?: number
  tag_ids?: number[]
  tags?: Array<{ id: number; name: string; description?: string }>
  author?: string
  file_size?: string
  view_count?: number
  cover?: string
  save_url?: string
  is_valid: boolean
  is_public: boolean
  created_at: string
  updated_at: string
}

interface Category {
  id: number
  name: string
  description?: string
}

interface Platform {
  id: number
  name: string
  description?: string
}

// Props
interface Props {
  show: boolean
  resource: Resource | null
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'update:show': [value: boolean]
  'updated': [resource: Resource]
}>()

// 消息提示
const message = useMessage()
const notification = useNotification()

// API
const resourceApi = useResourceApi()
const categoryApi = useCategoryApi()
const tagApi = useTagApi()
const panApi = usePanApi()

// 状态管理
const visible = computed({
  get: () => props.show,
  set: (value) => emit('update:show', value)
})

const submitting = ref(false)
const initializing = ref(false)
const editFormRef = ref()

// 编辑表单
const editForm = ref({
  title: '',
  description: '',
  url: '',
  category_id: null as number | null,
  pan_id: null as number | null,
  tag_ids: [] as number[],
  author: '',
  file_size: '',
  cover: '',
  save_url: '',
  is_valid: true,
  is_public: true
})

// 表单验证规则
const editRules = {
  title: {
    required: true,
    message: '请输入资源标题',
    trigger: 'blur'
  },
  url: {
    required: true,
    message: '请输入资源链接',
    trigger: 'blur'
  }
}

// 分类数据
const categories = ref<Category[]>([])
const platforms = ref<Platform[]>([])

// 标签搜索和加载相关状态
const tagLoading = ref(false)
const tagOptions = ref([])
const tagPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算属性
const categoryOptions = computed(() => {
  return categories.value.map(cat => ({
    label: cat.name,
    value: cat.id
  }))
})

const platformOptions = computed(() => {
  return platforms.value.map(pan => ({
    label: pan.name,
    value: pan.id
  }))
})

// 加载分类数据
const loadCategories = async () => {
  try {
    const response = await categoryApi.getCategories()
    categories.value = response?.items || response || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 加载平台数据
const loadPlatforms = async () => {
  try {
    const response = await panApi.getPans()
    platforms.value = response?.items || response || []
  } catch (error) {
    console.error('加载平台失败:', error)
  }
}

// 刷新标签选项
const loadTagOptions = async (keyword: string, page: number, pageSize: number) => {
  try {
    const response = await tagApi.getTags({
      search: keyword,
      page: page,
      page_size: pageSize
    })

    if (response) {
      const items = response.items || response.data || []
      const total = response.total || response.count || 0

      tagPagination.total = total

      const options = items.map((tag: any) => ({
        label: tag.name + (tag.description ? ` (${tag.description})` : ''),
        value: tag.id
      }))

      return { options, total }
    }
  } catch (error) {
    console.error('加载标签失败:', error)
  }
  return { options: [], total: 0 }
}

// 加载标签选项，确保包含当前资源的所有标签
const loadTagOptionsWithCurrentTags = async (currentTagIds: number[], currentTags?: Array<{ id: number; name: string; description?: string }>) => {
  tagLoading.value = true
  try {
    let allOptions: any[] = []

    // 专注于当前资源的标签显示
    if (currentTagIds.length > 0 && currentTags && currentTags.length > 0) {
      // 将当前资源的标签转换为选项格式
      const currentTagOptions = currentTags.map(tag => ({
        label: tag.name + (tag.description ? ` (${tag.description})` : ''),
        value: tag.id
      }))
      allOptions = [...currentTagOptions]
    } else if (currentTagIds.length > 0) {
      // 如果没有标签详细信息，创建临时选项
      const tempOptions = currentTagIds.map(id => ({
        label: `标签 ${id}`,
        value: id
      }))
      allOptions = [...tempOptions]
    }

    // 只在有标签时才加载少量补充选项，避免性能问题
    if (allOptions.length > 0) {
      try {
        // 只加载20个常用标签作为补充
        const { options: defaultOptions } = await loadTagOptions('', 1, 20)

        // 过滤掉已经存在的标签（避免重复）
        const existingTagIds = allOptions.map(option => option.value)
        const additionalOptions = defaultOptions.filter(option => !existingTagIds.includes(option.value))

        // 将补充标签添加到列表后面
        allOptions = [...allOptions, ...additionalOptions]
      } catch (error) {
        console.error('加载补充标签失败:', error)
        // 即使加载补充标签失败，也要确保当前标签显示
      }
    }

    // 🔧 修复：强制设置响应式数据
    tagOptions.value = []
    await nextTick() // 等待清空生效
    tagOptions.value = allOptions
    await nextTick() // 等待设置生效
  } catch (error) {
    console.error('加载标签选项失败:', error)
    // 出错时至少显示当前标签
    if (currentTagIds.length > 0) {
      const errorOptions = currentTagIds.map(id => ({
        label: `标签 ${id}`,
        value: id
      }))
      tagOptions.value = errorOptions
    } else {
      tagOptions.value = []
    }
  } finally {
    tagLoading.value = false
  }
}

// 处理标签搜索
const handleTagSearch = async (keyword: string) => {
  // 如果是空搜索且已有标签选项，不执行搜索（避免覆盖当前标签）
  if (!keyword.trim() && tagOptions.value.length > 0) {
    return
  }

  tagLoading.value = true
  try {
    const { options } = await loadTagOptions(keyword, 1, tagPagination.pageSize)

    // 保存当前已选中的标签
    const currentlySelected = editForm.value.tag_ids || []

    const currentSelectedOptions = tagOptions.value.filter(option =>
      currentlySelected.includes(option.value)
    )

    // 合并搜索结果，确保已选中标签始终在前面
    const existingIds = new Set(options.map(opt => opt.value))
    const missingSelectedOptions = currentSelectedOptions.filter(option =>
      !existingIds.has(option.value)
    )

    tagOptions.value = [...missingSelectedOptions, ...options]
  } catch (error) {
    console.error('搜索标签失败:', error)
  } finally {
    tagLoading.value = false
  }
}

// 处理标签滚动加载
const handleTagScroll = async () => {
  if (tagOptions.value.length >= tagPagination.total) {
    return // 已加载全部数据
  }

  const nextPage = Math.floor(tagOptions.value.length / tagPagination.pageSize) + 1
  if (nextPage > 1) {
    tagLoading.value = true
    try {
      const { options } = await loadTagOptions('', nextPage, tagPagination.pageSize)
      tagOptions.value = [...tagOptions.value, ...options]
    } catch (error) {
      console.error('加载更多标签失败:', error)
    } finally {
      tagLoading.value = false
    }
  }
}

// 提交编辑
const handleSubmit = async () => {
  if (!props.resource) return

  try {
    submitting.value = true

    const formData = {
      ...editForm.value,
      tag_ids: editForm.value.tag_ids || []
    }

    await resourceApi.updateResource(props.resource.id, formData)

    notification.success({
      content: '资源更新成功',
      duration: 2000
    })

    // 更新资源数据
    const updatedResource = { ...props.resource, ...editForm.value }
    emit('updated', updatedResource)

    visible.value = false
  } catch (error) {
    console.error('更新资源失败:', error)
    notification.error({
      content: '资源更新失败',
      duration: 2000
    })
  } finally {
    submitting.value = false
  }
}

// 关闭抽屉
const handleClose = () => {
  // 重置表单状态
  editForm.value = {
    title: '',
    description: '',
    url: '',
    category_id: null,
    pan_id: null,
    tag_ids: [],
    author: '',
    file_size: '',
    cover: '',
    save_url: '',
    is_valid: true,
    is_public: true
  }

  // 重置加载状态
  initializing.value = false
  submitting.value = false
  tagLoading.value = false

  // 重置标签选项
  tagOptions.value = []

  // 关闭抽屉
  visible.value = false
}

// 监听资源变化，初始化表单
watch(() => props.resource, async (newResource) => {
  if (newResource) {
    initializing.value = true
    try {
      // 确保平台和分类数据已加载，以便选项可用
      await Promise.all([
        loadCategories(),
        loadPlatforms()
      ])

      // 加载标签选项，确保包含当前资源的所有标签
      // 🔧 修复：从 tags 数组中提取 tag_ids
      const extractedTagIds = newResource.tags && Array.isArray(newResource.tags)
        ? newResource.tags.map((tag: any) => tag.id)
        : (newResource.tag_ids || [])

      await loadTagOptionsWithCurrentTags(extractedTagIds, newResource.tags || [])

      // 🔧 修复：使用提取的tag_ids设置表单
      const extractedTagIdsForForm = newResource.tags && Array.isArray(newResource.tags)
        ? newResource.tags.map((tag: any) => tag.id)
        : []

      editForm.value = {
        title: newResource.title,
        description: newResource.description || '',
        url: newResource.url,
        category_id: newResource.category_id || null,
        pan_id: newResource.pan_id || null,
        tag_ids: extractedTagIdsForForm,
        author: newResource.author || '',
        file_size: newResource.file_size || '',
        cover: newResource.cover || '',
        save_url: newResource.save_url || '',
        is_valid: newResource.is_valid !== undefined ? newResource.is_valid : true,
        is_public: newResource.is_public !== undefined ? newResource.is_public : true
      }

      // 🔧 强制刷新n-select组件
      await nextTick()

    } catch (error) {
      console.error('初始化表单失败:', error)
      message.error('初始化表单失败，请重试')
    } finally {
      initializing.value = false
    }
  }
})

// 监听抽屉打开，加载必要数据
watch(visible, async (isOpen) => {
  if (isOpen && props.resource && (categories.value.length === 0 || platforms.value.length === 0)) {
    // 如果平台或分类数据为空，则加载数据（例如首次打开时）
    await Promise.all([
      loadCategories(),
      loadPlatforms()
    ])
    // 加载标签选项
    // 🔧 修复：从 tags 数组中提取 tag_ids
    const extractedTagIds2 = props.resource.tags && Array.isArray(props.resource.tags)
      ? props.resource.tags.map((tag: any) => tag.id)
      : (props.resource.tag_ids || [])
    await loadTagOptionsWithCurrentTags(extractedTagIds2, props.resource.tags || [])
  }
})
</script>