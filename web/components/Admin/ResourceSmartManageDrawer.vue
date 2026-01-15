<template>
  <n-drawer v-model:show="visible" :width="1000" placement="right" :trap-focus="false" :block-scroll="true">
    <n-drawer-content :title="`AI优化 - ${resource?.title || '未命名资源'}`" closable>
      <!-- 顶部操作按钮区 -->
      <div class="flex space-x-3 mb-4">
        <n-button type="primary" @click="handleSmartProcess" :loading="processing.smartProcess" :disabled="processing.any">
          <template #icon>
            <i class="fas fa-magic"></i>
          </template>
          智能处理
        </n-button>
        <n-button @click="handleOptimizeTitle" :loading="processing.title" :disabled="processing.any">
          <template #icon>
            <i class="fas fa-heading"></i>
          </template>
          优化标题
        </n-button>
        <n-button @click="handleClassify" :loading="processing.classification" :disabled="processing.any">
          <template #icon>
            <i class="fas fa-tags"></i>
          </template>
          智能分类
        </n-button>
      </div>

      <!-- 左右分栏内容区 -->
      <div class="grid grid-cols-2 gap-4 h-full flex-grow overflow-hidden" style="height: calc(100% - 50px);">
        <!-- 左侧：当前内容预览 -->
        <div class="flex flex-col h-full overflow-y-auto">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">当前内容</h3>
          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 flex-grow overflow-auto">
            <div class="space-y-4">
              <!-- 标题 -->
              <div>
                <label class="text-sm font-medium text-gray-600 dark:text-gray-400">标题</label>
                <p class="mt-1 text-gray-900 dark:text-white whitespace-pre-wrap break-words">{{ resource?.title || '暂无标题' }}</p>
              </div>

              <!-- 描述 -->
              <div>
                <label class="text-sm font-medium text-gray-600 dark:text-gray-400">描述</label>
                <p class="mt-1 text-gray-900 dark:text-white whitespace-pre-wrap break-words">{{ resource?.description || '暂无描述' }}</p>
              </div>

              <!-- 分类 -->
              <div>
                <label class="text-sm font-medium text-gray-600 dark:text-gray-400">分类</label>
                <p class="mt-1 text-gray-900 dark:text-white">{{ getCategoryName(resource?.category_id) || '未分类' }}</p>
              </div>

              <!-- 标签 -->
              <div v-if="resource?.tags && resource.tags.length > 0">
                <label class="text-sm font-medium text-gray-600 dark:text-gray-400">标签</label>
                <div class="mt-1 flex flex-wrap gap-1">
                  <span
                    v-for="tag in resource.tags"
                    :key="tag.id"
                    class="text-xs px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded"
                  >
                    {{ tag.name }}
                  </span>
                </div>
              </div>

              <!-- 封面 -->
              <div v-if="resource?.cover">
                <label class="text-sm font-medium text-gray-600 dark:text-gray-400">封面图片</label>
                <div class="mt-1">
                  <img :src="resource.cover" alt="封面" class="max-w-full h-auto rounded" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：AI处理结果 -->
        <div class="flex flex-col h-full overflow-y-auto">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">AI处理结果</h3>
          <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 flex-grow overflow-auto">
            <div class="space-y-4">
              <!-- 处理中状态 -->
              <div v-if="processing.any" class="flex items-center justify-center h-32">
                <div class="text-center">
                  <n-spin size="medium" />
                  <p class="mt-2 text-gray-600 dark:text-gray-400">AI正在处理中...</p>
                </div>
              </div>

              <!-- 处理结果 -->
              <div v-else-if="hasResults" class="space-y-4">
                <!-- 优化后的标题 -->
                <div v-if="aiResult.title">
                  <div class="flex items-center justify-between mb-1">
                    <label class="text-sm font-medium text-gray-600 dark:text-gray-400">优化后的标题</label>
                    <n-button size="tiny" @click="applyTitle" :loading="applying.title">
                      <template #icon>
                        <i class="fas fa-check"></i>
                      </template>
                      应用标题
                    </n-button>
                  </div>
                  <p class="text-gray-900 dark:text-white bg-white dark:bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{{ aiResult.title }}</p>
                </div>

                <!-- 优化后的描述 -->
                <div v-if="aiResult.description">
                  <div class="flex items-center justify-between mb-1">
                    <label class="text-sm font-medium text-gray-600 dark:text-gray-400">优化后的描述</label>
                    <n-button size="tiny" @click="applyDescription" :loading="applying.description">
                      <template #icon>
                        <i class="fas fa-check"></i>
                      </template>
                      应用描述
                    </n-button>
                  </div>
                  <p class="text-gray-900 dark:text-white bg-white dark:bg-gray-800 p-2 rounded whitespace-pre-wrap break-words">{{ aiResult.description }}</p>
                </div>

                <!-- 建议分类 -->
                <div v-if="aiResult.category">
                  <div class="flex items-center justify-between mb-1">
                    <label class="text-sm font-medium text-gray-600 dark:text-gray-400">建议分类</label>
                    <n-button size="tiny" @click="applyCategory" :loading="applying.category">
                      <template #icon>
                        <i class="fas fa-check"></i>
                      </template>
                      应用分类
                    </n-button>
                  </div>
                  <p class="text-gray-900 dark:text-white bg-white dark:bg-gray-800 p-2 rounded">{{ aiResult.category.name }}</p>
                </div>
              </div>

              <!-- 空状态 -->
              <div v-else class="flex items-center justify-center h-32">
                <div class="text-center text-gray-500 dark:text-gray-400">
                  <i class="fas fa-robot text-4xl mb-2"></i>
                  <p>点击上方按钮开始AI处理</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部操作按钮 -->
      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="handleClose">关闭</n-button>
          <n-button
            type="primary"
            @click="handleApplyAll"
            :disabled="!hasResults"
            :loading="applying.all"
          >
            <template #icon>
              <i class="fas fa-check-double"></i>
            </template>
            应用所有更改
          </n-button>
        </div>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useAIApi, useCategoryApi } from '~/composables/useApi'
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
const aiApi = useAIApi()
const categoryApi = useCategoryApi()

// 状态管理
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const processing = reactive({
  smartProcess: false,
  title: false,
  classification: false,
  any: computed(() => processing.smartProcess || processing.title || processing.classification)
})

const applying = reactive({
  title: false,
  description: false,
  category: false,
  all: false
})

const aiResult = reactive({
  title: '',
  description: '',
  category: null as { id: number; name: string } | null
})

const categories = ref<Category[]>([])

// 计算属性
const hasResults = computed(() => {
  return aiResult.title || aiResult.description || aiResult.category
})

// 获取分类名称
const getCategoryName = (categoryId?: number) => {
  if (!categoryId || !categories.value) return null
  const category = categories.value.find(cat => cat.id === categoryId)
  return category?.name
}

// 智能处理（全部功能）
const handleSmartProcess = async () => {
  if (!props.resource) return

  processing.smartProcess = true
  clearResults()

  try {
    // 并行执行所有AI功能
    const promises = [
      optimizeTitle(),
      generateDescription(),
      classifyResource()
    ]

    await Promise.all(promises)

    message.success('智能处理完成')
  } catch (error) {
    console.error('智能处理失败:', error)
    message.error('智能处理失败')
  } finally {
    processing.smartProcess = false
  }
}

// 优化标题
const handleOptimizeTitle = async () => {
  await optimizeTitle()
}

const optimizeTitle = async () => {
  if (!props.resource) return

  processing.title = true

  try {
    const prompt = `根据以下资源信息生成一个简洁、吸引人的标题：
资源标题: ${props.resource.title}
资源描述: ${props.resource.description || ''}
资源URL: ${props.resource.url}

请生成一个更优化的标题，要求：
1. 简洁明了，突出重点
2. 吸引用户点击
3. 符合中文表达习惯
4. 长度控制在50字以内`

    const response = await aiApi.generateText({
      prompt: prompt,
      options: [
        { type: 'max_tokens', value: 200 },
        { type: 'temperature', value: 0.7 }
      ]
    })

    if (response?.result) {
      aiResult.title = response.result.trim()
      message.success('标题优化完成')
    }
  } catch (error) {
    console.error('优化标题失败:', error)
    message.error('优化标题失败')
  } finally {
    processing.title = false
  }
}

// 生成描述
const generateDescription = async () => {
  if (!props.resource) return

  try {
    const prompt = `根据以下资源信息生成一个详细、有吸引力的描述：
资源标题: ${props.resource.title}
资源URL: ${props.resource.url}
现有描述: ${props.resource.description || ''}

请生成一个更好的描述，要求：
1. 详细介绍资源内容和特点
2. 突出资源的价值
3. 语言生动，有吸引力
4. 长度控制在200-500字之间
5. 可以包含资源的具体特点、适用人群等信息`

    const response = await aiApi.generateText({
      prompt: prompt,
      options: [
        { type: 'max_tokens', value: 800 },
        { type: 'temperature', value: 0.7 }
      ]
    })

    if (response?.result) {
      aiResult.description = response.result.trim()
    }
  } catch (error) {
    console.error('生成描述失败:', error)
  }
}

// 智能分类
const handleClassify = async () => {
  await classifyResource()
}

const classifyResource = async () => {
  if (!props.resource) return

  // 确保分类数据已加载
  if (categories.value.length === 0) {
    try {
      const response = await categoryApi.getCategories()
      categories.value = response || []
      if (categories.value.length === 0) {
        message.warning('暂无可用分类，请先添加分类')
        return
      }
    } catch (error) {
      console.error('加载分类失败:', error)
      message.error('加载分类失败')
      return
    }
  }

  processing.classification = true

  try {
    const categoryList = categories.value.map(cat => `${cat.id}: ${cat.name}`).join('\n')

    const prompt = `根据以下资源信息，为其选择一个最合适的分类：

资源标题: ${props.resource.title}
资源描述: ${props.resource.description || ''}
资源URL: ${props.resource.url}

可用分类列表：
${categoryList}

请直接返回最适合的分类ID，不要返回其他内容。`

    const response = await aiApi.generateText({
      prompt: prompt,
      options: [
        { type: 'max_tokens', value: 10 },
        { type: 'temperature', value: 0.3 }
      ]
    })

    if (response?.result) {
      const result = response.result.trim()
      const categoryId = parseInt(result)

      const validCategory = categories.value.find(cat => cat.id === categoryId)
      if (validCategory) {
        aiResult.category = { id: validCategory.id, name: validCategory.name }
        message.success(`智能分类完成: ${validCategory.name}`)
      } else {
        const match = result.match(/\d+/)
        if (match) {
          const id = parseInt(match[0])
          const validCategory = categories.value.find(cat => cat.id === id)
          if (validCategory) {
            aiResult.category = { id: validCategory.id, name: validCategory.name }
            message.success(`智能分类完成: ${validCategory.name}`)
          }
        }
      }
    }
  } catch (error) {
    console.error('智能分类失败:', error)
    message.error('智能分类失败')
  } finally {
    processing.classification = false
  }
}

// 应用标题
const applyTitle = async () => {
  if (!aiResult.title || !props.resource) return

  applying.title = true

  try {
    const response = await aiApi.applyGeneratedContent({
      resource_id: props.resource.id,
      field: 'title',
      content: aiResult.title
    })

    if (response) {
      message.success('标题已应用')
      emit('updated', { ...props.resource, title: aiResult.title })
    }
  } catch (error) {
    console.error('应用标题失败:', error)
    message.error('应用标题失败')
  } finally {
    applying.title = false
  }
}

// 应用描述
const applyDescription = async () => {
  if (!aiResult.description || !props.resource) return

  applying.description = true

  try {
    const response = await aiApi.applyGeneratedContent({
      resource_id: props.resource.id,
      field: 'description',
      content: aiResult.description
    })

    if (response) {
      message.success('描述已应用')
      emit('updated', { ...props.resource, description: aiResult.description })
    }
  } catch (error) {
    console.error('应用描述失败:', error)
    message.error('应用描述失败')
  } finally {
    applying.description = false
  }
}

// 应用分类
const applyCategory = async () => {
  if (!aiResult.category || !props.resource) return

  applying.category = true

  try {
    const response = await aiApi.applyClassification({
      resource_id: props.resource.id,
      category_id: aiResult.category.id
    })

    if (response) {
      message.success('分类已应用')
      emit('updated', { ...props.resource, category_id: aiResult.category.id })
    }
  } catch (error) {
    console.error('应用分类失败:', error)
    message.error('应用分类失败')
  } finally {
    applying.category = false
  }
}

// 应用所有更改
const handleApplyAll = async () => {
  if (!props.resource) return

  applying.all = true

  try {
    const promises = []

    if (aiResult.title) {
      promises.push(
        aiApi.applyGeneratedContent({
          resource_id: props.resource.id,
          field: 'title',
          content: aiResult.title
        })
      )
    }

    if (aiResult.description) {
      promises.push(
        aiApi.applyGeneratedContent({
          resource_id: props.resource.id,
          field: 'description',
          content: aiResult.description
        })
      )
    }

    if (aiResult.category) {
      promises.push(
        aiApi.applyClassification({
          resource_id: props.resource.id,
          category_id: aiResult.category.id
        })
      )
    }

    const results = await Promise.all(promises)

    if (results.every(r => r)) {
      message.success('所有更改已应用')

      // 更新资源数据
      const updatedResource = { ...props.resource }
      if (aiResult.title) updatedResource.title = aiResult.title
      if (aiResult.description) updatedResource.description = aiResult.description
      if (aiResult.category) updatedResource.category_id = aiResult.category.id

      emit('updated', updatedResource)

      // 清空结果
      clearResults()
    } else {
      message.warning('部分更改应用失败')
    }
  } catch (error) {
    console.error('应用更改失败:', error)
    message.error('应用更改失败')
  } finally {
    applying.all = false
  }
}

// 清空结果
const clearResults = () => {
  aiResult.title = ''
  aiResult.description = ''
  aiResult.category = null
}

// 关闭抽屉
const handleClose = () => {
  visible.value = false
}

// 监听资源变化，加载分类数据
watch(() => props.resource, async (newResource) => {
  if (newResource && visible.value) {
    // 加载分类数据
    try {
      const response = await categoryApi.getCategories()
      categories.value = response || []
    } catch (error) {
      console.error('加载分类失败:', error)
    }

    // 清空之前的结果
    clearResults()
  }
})

// 监听抽屉打开
watch(visible, async (isOpen) => {
  if (isOpen && props.resource) {
    // 加载分类数据
    try {
      const response = await categoryApi.getCategories()
      categories.value = response || []
    } catch (error) {
      console.error('加载分类失败:', error)
    }
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>