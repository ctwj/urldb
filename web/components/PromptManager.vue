<template>
  <div class="h-full w-full overflow-hidden flex flex-col">

    <!-- 主内容区域 - 自适应剩余高度 -->
    <div class="flex-1 bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div class="h-full flex flex-col lg:flex-row">
        <!-- 左侧：系统提示词列表 -->
        <div class="w-full lg:w-1/3 border-r border-gray-200 dark:border-gray-700 flex flex-col">
          <div class="flex-1 overflow-y-auto p-3">
            <div v-if="!systemPrompts || systemPrompts.length === 0" class="text-gray-500 text-center py-8">
              暂无系统提示词
            </div>
            <div v-else class="space-y-2">
              <div v-for="prompt in (systemPrompts || [])" :key="prompt.id"
                   class="border rounded-lg p-3 bg-blue-50 dark:bg-blue-900/20 cursor-pointer hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-colors"
                   :class="{ 'ring-2 ring-blue-500': selectedPrompt?.id === prompt.id }"
                   @click="selectPrompt(prompt)">
                <div class="flex items-center justify-between mb-1">
                  <div class="flex items-center space-x-2">
                    <div class="w-2 h-2 rounded-full bg-blue-400"></div>
                    <div class="font-medium text-sm text-gray-900 dark:text-gray-100">{{ prompt.name }}</div>
                  </div>
                  <div class="flex items-center space-x-1">
                    <!-- <n-tag type="info" size="tiny">
                      系统内置
                    </n-tag>
                    <n-button size="tiny" @click.stop="selectPrompt(prompt)" type="primary" v-if="selectedPrompt?.id !== prompt.id">
                      <template #icon>
                        <i class="fas fa-edit text-xs"></i>
                      </template>
                    </n-button>
                    <n-tag v-else type="success" size="tiny">
                      <template #icon>
                        <i class="fas fa-check text-xs"></i>
                      </template>
                    </n-tag> -->
                  </div>
                </div>
                <div class="text-xs text-gray-600 dark:text-gray-400">
                  {{ prompt.type_description }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：提示词编辑器 -->
        <div class="w-full lg:w-2/3 flex flex-col">
          <div v-if="!selectedPrompt" class="flex-1 flex items-center justify-center text-gray-500">
            <div class="text-center p-8">
              <i class="fas fa-mouse-pointer text-4xl mb-4"></i>
              <p>请选择一个提示词进行编辑</p>
            </div>
          </div>
          <div v-else class="flex-1 flex flex-col">
            <!-- 提示词信息头部 -->
            <div class="px-4 py-3 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-700">
              <div class="flex items-center justify-between mb-2">
                <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                  {{ selectedPrompt.name }}
                </h3>
                <div class="flex items-center space-x-1">
                  <n-button @click="testCurrentPrompt" size="small" type="info" :disabled="!editingPrompt.user_content" :loading="testing">
                    <template #icon>
                      <i class="fas fa-play"></i>
                    </template>
                    测试
                  </n-button>
                  <n-button @click="savePrompt" type="primary" size="small" :loading="saving" :disabled="!hasPromptChanges">
                    <template #icon>
                      <i class="fas fa-save"></i>
                    </template>
                    保存
                  </n-button>
                </div>
              </div>
              <div class="flex flex-wrap gap-1">
                <n-tag v-for="variable in editingPrompt.variables" :key="variable" size="tiny" type="info">
                  {{ variable }}
                </n-tag>
                <span v-if="editingPrompt.variables.length === 0" class="text-xs text-gray-500">无变量</span>
              </div>
            </div>

            <!-- 提示词内容编辑器 - 自适应剩余高度 -->
            <div class="flex-1 flex flex-col">
              <!-- 系统提示词编辑器 -->
              <div class="flex-1 flex flex-col border-b border-gray-200 dark:border-gray-700">
                <div class="px-3 py-2 bg-gray-50 dark:bg-gray-700">
                  <div class="flex items-center space-x-2">
                    <div class="w-2 h-2 rounded-full bg-red-400"></div>
                    <span class="text-xs font-medium text-gray-700 dark:text-gray-300">系统提示词 (System Prompt)</span>
                    <n-tag size="tiny" type="info">可编辑</n-tag>
                  </div>
                  <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    定义AI角色和行为规则，可自定义修改
                  </div>
                </div>
                <div class="flex-1 p-3">
                  <n-input
                    v-model:value="editingPrompt.system_content"
                    placeholder="系统提示词内容..."
                    type="textarea"
                    class="font-mono text-sm h-full border-none"
                    style="height: 100%;"
                  />
                </div>
              </div>

              <!-- 用户提示词编辑器 -->
              <div class="flex-1 flex flex-col">
                <div class="px-3 py-2 bg-gray-50 dark:bg-gray-700">
                  <div class="flex items-center space-x-2">
                    <div class="w-2 h-2 rounded-full bg-blue-400"></div>
                    <span class="text-xs font-medium text-gray-700 dark:text-gray-300">用户提示词 (User Prompt)</span>
                    <n-tag size="tiny" type="info">可编辑</n-tag>
                  </div>
                  <div class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    包含变量模板，可自定义修改
                  </div>
                </div>
                <div class="flex-1 p-3">
                  <n-input
                    v-model:value="editingPrompt.user_content"
                    placeholder="请输入用户提示词内容，支持 {{variable}} 格式的变量"
                    type="textarea"
                    class="font-mono text-sm h-full border-none"
                    style="height: 100%;"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useMessage, useNotification } from 'naive-ui'
import { usePromptApi } from '~/composables/usePromptApi'

// Props
interface Props {
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  saving: false
})

// Emits
const emit = defineEmits<{
  'update:saving': [value: boolean]
}>()

// Composables
const message = useMessage()
const notification = useNotification()
const { getPrompts, updatePrompt, testPrompt } = usePromptApi()

// 提示词管理相关变量
const prompts = ref<any[]>([])
const selectedPromptId = ref<number | null>(null)
const testing = ref(false)
const editingPrompt = ref({
  id: 0,
  name: '',
  type: '',
  system_content: '',
  user_content: '',
  description: '',
  variables: [],
  is_active: true
})


// 计算属性
// 系统提示词计算属性（显示所有系统提示词）
const systemPrompts = computed(() => {
  if (!prompts.value) return []
  return prompts.value.filter(p => [
    'classification',
    'content_generation',
    'tool_system',
    'qa_template',
    'analysis_template'
  ].includes(p.type))
})

const selectedPrompt = computed(() => {
  if (!selectedPromptId.value) return null
  return prompts.value.find(p => p.id === selectedPromptId.value) || null
})

const hasPromptChanges = computed(() => {
  if (!selectedPrompt.value) return false
  return editingPrompt.value.system_content !== selectedPrompt.value.system_content ||
         editingPrompt.value.user_content !== selectedPrompt.value.user_content
})

// 方法
const loadPrompts = async () => {
  try {
    const response = await getPrompts()
    prompts.value = response || []
    if (prompts.value.length > 0 && !selectedPromptId.value) {
      selectedPromptId.value = prompts.value[0].id
    }
  } catch (error) {
    console.error('加载提示词失败:', error)
    notification.error({
      content: '加载提示词失败',
      duration: 3000
    })
  }
}

const selectPrompt = (prompt: any) => {
  selectedPromptId.value = prompt.id
  // 创建编辑副本，避免直接修改原数据
  editingPrompt.value = { ...prompt }
}

const savePrompt = async () => {
  if (!selectedPromptId.value) return

  try {
    emit('update:saving', true)

    // 构建更新数据，支持系统提示词和用户提示词
    const updateData = {}

    if (editingPrompt.value.system_content !== selectedPrompt.value.system_content) {
      updateData.system_content = editingPrompt.value.system_content
    }

    if (editingPrompt.value.user_content !== selectedPrompt.value.user_content) {
      updateData.user_content = editingPrompt.value.user_content
    }

    await updatePrompt(selectedPromptId.value, updateData)

    notification.success({
      content: '提示词保存成功',
      duration: 3000
    })

    // 重新加载提示词列表
    await loadPrompts()
  } catch (error) {
    console.error('保存提示词失败:', error)
    notification.error({
      content: '保存提示词失败',
      duration: 3000
    })
  } finally {
    emit('update:saving', false)
  }
}

const resetPrompt = () => {
  if (selectedPrompt.value) {
    editingPrompt.value = { ...selectedPrompt.value }
  }
}

const testCurrentPrompt = async () => {
  if (!editingPrompt.value.user_content) {
    notification.warning({
      content: '请先输入用户提示词内容',
      duration: 3000
    })
    return
  }

  try {
    testing.value = true

    // 根据提示词类型准备测试数据
    let testData = {}
    let testMessage = ''

    switch (editingPrompt.value.type) {
      case 'classification':
        testData = {
          title: '人工智能技术发展趋势分析',
          description: '探讨人工智能在各个领域的应用前景和挑战',
          author: 'AI研究专家',
          tags: ['人工智能', '技术趋势', '未来发展']
        }
        testMessage = '请对以下资源进行智能分类推荐'
        break
      case 'content_generation':
        testData = {
          title: '机器学习基础教程',
          description: '这是一篇关于机器学习基础知识的详细介绍',
          current_content: '机器学习是人工智能的一个重要分支'
        }
        testMessage = '请为以下资源生成优化的标题和描述'
        break
      case 'tool_system':
        testData = {
          query: '现在几点了？',
          context: '用户询问当前时间'
        }
        testMessage = '用户询问时间，请使用工具查询'
        break
      case 'qa_template':
        testData = {
          question: '什么是人工智能？',
          context: '科技知识问答'
        }
        testMessage = '请回答以下问题'
        break
      case 'analysis_template':
        testData = {
          text: '人工智能正在改变我们的生活方式，从智能手机到自动驾驶汽车，AI技术无处不在。',
          analysis_type: '情感分析'
        }
        testMessage = '请分析以下文本'
        break
      default:
        testData = {
          input: '这是一个测试输入'
        }
        testMessage = '测试提示词功能'
    }

    // 构建测试请求数据
    const testRequest = {
      prompt_type: editingPrompt.value.type,
      system_content: editingPrompt.value.system_content,
      user_content: editingPrompt.value.user_content,
      test_data: testData,
      test_message: testMessage
    }

    const response = await testPrompt(testRequest)

    notification.success({
      content: '提示词测试成功',
      duration: 3000
    })

    // 可以选择显示测试结果
    console.log('测试结果:', response)

  } catch (error) {
    console.error('测试提示词失败:', error)
    notification.error({
      content: '测试提示词失败，请检查配置',
      duration: 3000
    })
  } finally {
    testing.value = false
  }
}



// 监听选中变化，同步编辑数据
watch(selectedPrompt, (newPrompt) => {
  if (newPrompt) {
    editingPrompt.value = { ...newPrompt }
  }
}, { immediate: true })

// 初始化
onMounted(() => {
  loadPrompts()
})

// 暴露方法给父组件
defineExpose({
  loadPrompts
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