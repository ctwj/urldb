<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">网盘规则管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理网盘识别规则，支持动态添加和更新</p>
      </div>
      <div class="flex space-x-3">
        <n-button type="primary" @click="showAddModal = true">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          添加规则
        </n-button>
        <n-button @click="refreshData">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
        <n-button type="info" @click="handleRefreshCache">
          <template #icon>
            <i class="fas fa-database"></i>
          </template>
          刷新缓存
        </n-button>
      </div>
    </div>

    <n-card>
      <div class="flex gap-4">
        <n-select
          v-model:value="selectedPan"
          placeholder="选择网盘"
          class="w-40"
        >
          <n-option value="">全部网盘</n-option>
          <n-option v-for="pan in pans" :key="pan.id" :value="pan.id">
            {{ pan.name }}
          </n-option>
        </n-select>
        
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索规则名称..."
          @keyup.enter="handleSearch"
          class="flex-1"
          clearable
        >
          <template #prefix>
            <i class="fas fa-search"></i>
          </template>
        </n-input>
        
        <n-button type="primary" @click="handleSearch" class="w-20">
          <template #icon>
            <i class="fas fa-search"></i>
          </template>
          搜索
        </n-button>
      </div>
    </n-card>

    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">规则列表</span>
          <span class="text-sm text-gray-500">共 {{ total }} 条规则</span>
        </div>
      </template>

      <div v-if="loading" class="flex items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="rules.length === 0" class="text-center py-8">
        <i class="fas fa-file-alt text-4xl text-gray-400 mb-4"></i>
        <p class="text-gray-500">暂无规则数据</p>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="rule in rules"
          :key="rule.id"
          class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                  {{ rule.name }}
                </h3>
                <span v-if="rule.pan" class="text-xs px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 rounded">
                  {{ rule.pan.name }}
                </span>
                <span v-if="rule.enabled" class="text-xs px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded">
                  启用
                </span>
                <span v-else class="text-xs px-2 py-1 bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200 rounded">
                  禁用
                </span>
              </div>
              
              <div class="mb-2">
                <span class="text-xs text-gray-500">域名：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ rule.domains }}</span>
              </div>
              
              <div class="mb-2">
                <span class="text-xs text-gray-500">URL模式：</span>
                <span class="text-sm text-gray-700 dark:text-gray-300 font-mono">{{ rule.url_patterns || '未设置' }}</span>
              </div>
              
              <div class="flex items-center space-x-4 text-xs text-gray-500">
                <span>优先级: {{ rule.priority }}</span>
                <span>创建时间: {{ formatDate(rule.created_at) }}</span>
                <span>更新时间: {{ formatDate(rule.updated_at) }}</span>
              </div>
              
              <p v-if="rule.remark" class="text-sm text-gray-600 dark:text-gray-400 mt-2">
                {{ rule.remark }}
              </p>
            </div>
            
            <div class="flex space-x-2">
              <n-button size="small" type="primary" @click="editRule(rule)">
                <template #icon>
                  <i class="fas fa-edit"></i>
                </template>
                编辑
              </n-button>
              <n-button size="small" type="error" @click="deleteRule(rule)">
                <template #icon>
                  <i class="fas fa-trash"></i>
                </template>
                删除
              </n-button>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6">
        <n-pagination
          v-model:page="currentPage"
          v-model:page-size="pageSize"
          :item-count="total"
          :page-sizes="[10, 20, 50, 100]"
          show-size-picker
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <n-modal v-model:show="showAddModal" preset="card" :title="editingRule ? '编辑规则' : '添加规则'" style="width: 600px">
      <n-form
        ref="formRef"
        :model="ruleForm"
        :rules="rulesForm"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="网盘" path="pan_id">
          <n-select
            v-model:value="ruleForm.pan_id"
            placeholder="请选择网盘"
            :options="panOptions"
          />
        </n-form-item>

        <n-form-item label="规则名称" path="name">
          <n-input
            v-model:value="ruleForm.name"
            placeholder="请输入规则名称"
          />
        </n-form-item>

        <n-form-item label="域名列表" path="domains">
          <n-input
            v-model:value="ruleForm.domains"
            placeholder="多个域名用逗号分隔，如: pan.quark.cn,www.aliyundrive.com"
            type="textarea"
            :rows="2"
          />
          <template #help>
            支持多个域名，使用英文逗号分隔
          </template>
        </n-form-item>

        <n-form-item label="URL正则模式" path="url_patterns">
          <n-input
            v-model:value="ruleForm.url_patterns"
            placeholder="正则表达式，如: https?://pan\.quark\.cn/s/([a-zA-Z0-9]+)"
            type="textarea"
            :rows="3"
          />
          <template #help>
            正则表达式，分组捕获分享ID。支持多个模式用逗号分隔
          </template>
        </n-form-item>

        <n-form-item label="优先级" path="priority">
          <n-input-number
            v-model:value="ruleForm.priority"
            :min="1"
            :max="100"
            placeholder="优先级"
          />
          <template #help>
            数字越小优先级越高，默认1
          </template>
        </n-form-item>

        <n-form-item label="启用状态" path="enabled">
          <n-switch v-model:value="ruleForm.enabled" />
        </n-form-item>

        <n-form-item label="备注" path="remark">
          <n-input
            v-model:value="ruleForm.remark"
            placeholder="请输入备注说明"
            type="textarea"
            :rows="2"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end space-x-3">
          <n-button @click="closeModal">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            保存
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin' as any
})

const { usePanRuleApi } = await import('~/composables/useApi')
const { usePanApi } = await import('~/composables/useApi')
const panRuleApi = usePanRuleApi()
const panApi = usePanApi()

const loading = ref(false)
const rules = ref<any[]>([])
const pans = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const selectedPan = ref('')

const showAddModal = ref(false)
const submitting = ref(false)
const editingRule = ref<any>(null)

const ruleForm = ref({
  pan_id: 0,
  name: '',
  domains: '',
  url_patterns: '',
  priority: 1,
  enabled: true,
  remark: ''
})

const rulesForm = {
  pan_id: {
    required: true,
    message: '请选择网盘',
    trigger: 'blur'
  },
  name: {
    required: true,
    message: '请输入规则名称',
    trigger: 'blur'
  },
  domains: {
    required: true,
    message: '请输入域名列表',
    trigger: 'blur'
  }
}

const formRef = ref()

const panOptions = computed(() => {
  return pans.value.map(pan => ({
    label: pan.name,
    value: pan.id
  }))
})

const fetchPans = async () => {
  try {
    const response = await panApi.getPans() as any
    pans.value = response || []
  } catch (error: any) {
    useNotification().error({
      content: error.message || '获取网盘列表失败',
      duration: 5000
    })
  }
}

const fetchRules = async () => {
  try {
    loading.value = true
    
    const response = await panRuleApi.getPanRules() as any
    rules.value = response || []
    total.value = rules.value.length
  } catch (error: any) {
    useNotification().error({
      content: error.message || '获取规则数据失败',
      duration: 5000
    })
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchPans()
  await fetchRules()
})

const handleSearch = () => {
  currentPage.value = 1
  fetchRules()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchRules()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchRules()
}

const refreshData = () => {
  fetchRules()
}

const editRule = (rule: any) => {
  editingRule.value = rule
  ruleForm.value = {
    pan_id: rule.pan_id,
    name: rule.name,
    domains: rule.domains,
    url_patterns: rule.url_patterns || '',
    priority: rule.priority || 1,
    enabled: rule.enabled,
    remark: rule.remark || ''
  }
  showAddModal.value = true
}

const deleteRule = async (rule: any) => {
  try {
    await panRuleApi.deletePanRule(rule.id)
    useNotification().success({
      content: '规则删除成功',
      duration: 3000
    })
    await fetchRules()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '删除规则失败',
      duration: 5000
    })
  }
}

const closeModal = () => {
  showAddModal.value = false
  editingRule.value = null
  ruleForm.value = {
    pan_id: 0,
    name: '',
    domains: '',
    url_patterns: '',
    priority: 1,
    enabled: true,
    remark: ''
  }
}

const handleSubmit = async () => {
  try {
    submitting.value = true
    
    if (editingRule.value) {
      await panRuleApi.updatePanRule(editingRule.value.id, ruleForm.value)
      useNotification().success({
        content: '规则更新成功',
        duration: 3000
      })
    } else {
      await panRuleApi.createPanRule(ruleForm.value)
      useNotification().success({
        content: '规则创建成功',
        duration: 3000
      })
    }
    
    closeModal()
    await fetchRules()
  } catch (error: any) {
    useNotification().error({
      content: error.message || '保存规则失败',
      duration: 5000
    })
  } finally {
    submitting.value = false
  }
}

const handleRefreshCache = async () => {
  try {
    await panRuleApi.refreshPanRules()
    useNotification().success({
      content: '规则缓存刷新成功',
      duration: 3000
    })
  } catch (error: any) {
    useNotification().error({
      content: error.message || '刷新缓存失败',
      duration: 5000
    })
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>