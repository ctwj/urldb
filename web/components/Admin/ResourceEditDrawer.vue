<template>
  <n-drawer v-model:show="visible" :width="800" placement="right" :trap-focus="false" :block-scroll="true">
    <n-drawer-content :title="editingResource ? `编辑资源 - ${editingResource.title}` : '编辑资源'" closable>
      <n-form ref="editFormRef" :model="editForm" :rules="editRules" label-placement="left" label-width="100">
        <n-form-item label="标题" path="title">
          <n-input v-model:value="editForm.title" placeholder="请输入资源标题" />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="editForm.description"
            type="textarea"
            placeholder="请输入资源描述"
            :autosize="{ minRows: 3, maxRows: 6 }"
            class="w-full"
          />
        </n-form-item>

        <n-form-item label="资源链接" path="url">
          <n-input v-model:value="editForm.url" placeholder="请输入资源链接" />
        </n-form-item>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="editForm.category_id"
            :options="categoryOptions"
            placeholder="请选择分类"
            clearable
          />
        </n-form-item>

        <n-form-item label="平台" path="pan_id">
          <n-select
            v-model:value="editForm.pan_id"
            :options="platformOptions"
            placeholder="请选择平台"
            clearable
          />
        </n-form-item>

        <n-form-item label="标签" path="tag_ids">
          <n-select
            v-model:value="editForm.tag_ids"
            :options="tagOptions"
            :loading="tagLoading"
            :filterable="true"
            :remote="true"
            :clearable="true"
            placeholder="请选择标签，支持搜索"
            multiple
            @search="handleTagSearch"
            @scroll="handleTagScroll"
          />
        </n-form-item>

        <n-form-item label="作者" path="author">
          <n-input v-model:value="editForm.author" placeholder="请输入作者" />
        </n-form-item>

        <n-form-item label="文件大小" path="file_size">
          <n-input v-model:value="editForm.file_size" placeholder="如：2.5GB" />
        </n-form-item>

        <n-form-item label="封面图片" path="cover">
          <n-input v-model:value="editForm.cover" placeholder="请输入封面图片URL" />
        </n-form-item>

        <n-form-item label="转存链接" path="save_url">
          <n-input v-model:value="editForm.save_url" placeholder="请输入转存链接" />
        </n-form-item>

        <n-form-item label="是否有效" path="is_valid">
          <n-switch v-model:value="editForm.is_valid" />
        </n-form-item>

        <n-form-item label="是否公开" path="is_public">
          <n-switch v-model:value="editForm.is_public" />
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
import { ref, reactive, computed, watch } from 'vue'
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
  modelValue: boolean
  resource: Resource | null
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
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
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const submitting = ref(false)
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
const loadTagOptionsWithCurrentTags = async (currentTagIds: number[]) => {
  if (currentTagIds.length === 0) {
    // 如果没有当前标签，只需加载默认标签列表
    await handleTagSearch('')
    return
  }

  // 先加载默认标签列表
  const { options: defaultOptions } = await loadTagOptions('', 1, tagPagination.pageSize)

  // 检查当前标签是否都在默认选项中
  const existingTagIds = defaultOptions.map((option: any) => option.value)
  const missingTagIds = currentTagIds.filter(id => !existingTagIds.includes(id))

  let allOptions = [...defaultOptions]

  if (missingTagIds.length > 0) {
    // 如果有缺失的标签ID，我们需要获取这些标签的详细信息
    // 方法是：尝试通过ID批量获取标签信息
    try {
      // 为了获取缺失标签的详细信息，我们可以逐个查询
      const missingTagDetails = []
      for (const tagId of missingTagIds) {
        // 尝试获取单个标签的信息
        try {
          // 假设存在一个根据ID获取单个标签的API方法
          // 如果API没有专门的方法，我们可以通过搜索来获取
          const searchResponse = await tagApi.getTags({
            search: tagId.toString(),
            page: 1,
            page_size: 1
          })

          const items = searchResponse?.items || searchResponse?.data || []
          const foundTag = items.find((tag: any) => tag.id === tagId)

          if (foundTag) {
            missingTagDetails.push({
              label: foundTag.name + (foundTag.description ? ` (${foundTag.description})` : ''),
              value: foundTag.id
            })
          } else {
            // 如果通过搜索找不到，使用临时值
            missingTagDetails.push({
              label: `标签 ${tagId}`,
              value: tagId
            })
          }
        } catch (error) {
          console.error(`获取标签 ${tagId} 信息失败:`, error)
          // 出错时使用临时值
          missingTagDetails.push({
            label: `标签 ${tagId}`,
            value: tagId
          })
        }
      }

      allOptions = [...defaultOptions, ...missingTagDetails]
    } catch (error) {
      console.error('获取缺失标签详情失败:', error)
      // 出错时仍然添加临时选项
      const missingOptions = missingTagIds.map(id => ({
        label: `标签 ${id}`,
        value: id
      }))

      allOptions = [...defaultOptions, ...missingOptions]
    }
  }

  tagOptions.value = allOptions
}

// 处理标签搜索
const handleTagSearch = async (keyword: string) => {
  tagLoading.value = true
  try {
    const { options } = await loadTagOptions(keyword, 1, tagPagination.pageSize)
    tagOptions.value = options
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
  visible.value = false
}

// 监听资源变化，初始化表单
watch(() => props.resource, async (newResource) => {
  if (newResource) {
    // 确保平台和分类数据已加载，以便选项可用
    await Promise.all([
      loadCategories(),
      loadPlatforms()
    ])

    // 加载标签选项，确保包含当前资源的所有标签
    await loadTagOptionsWithCurrentTags(newResource.tag_ids || [])

    editForm.value = {
      title: newResource.title,
      description: newResource.description || '',
      url: newResource.url,
      category_id: newResource.category_id || null,
      pan_id: newResource.pan_id || null,
      tag_ids: newResource.tag_ids || [],
      author: newResource.author || '',
      file_size: newResource.file_size || '',
      cover: newResource.cover || '',
      save_url: newResource.save_url || '',
      is_valid: newResource.is_valid !== undefined ? newResource.is_valid : true,
      is_public: newResource.is_public !== undefined ? newResource.is_public : true
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
    await loadTagOptionsWithCurrentTags(props.resource.tag_ids || [])
  }
})
</script>