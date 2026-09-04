<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white tracking-tight">我的资源</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1.5">提交网盘资源自动检测有效性，有效资源经处理公开发布</p>
      </div>
      <div class="flex items-center gap-3 flex-shrink-0">
        <n-button @click="showBatchModal = true">
          <template #icon>
            <i class="fas fa-layer-group"></i>
          </template>
          批量提交
        </n-button>
        <n-button type="primary" @click="openSubmitModal">
          <template #icon>
            <i class="fas fa-plus"></i>
          </template>
          提交资源
        </n-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div
        v-for="card in statCards"
        :key="card.key"
        class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200/80 dark:border-gray-700/60 px-5 py-4 flex items-center gap-4 transition-colors duration-200 hover:border-blue-300 dark:hover:border-blue-700 cursor-default"
      >
        <div
          class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0"
          :class="card.iconBg"
        >
          <i :class="[card.icon, card.iconColor, 'text-base']"></i>
        </div>
        <div class="min-w-0">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ card.label }}</p>
          <p class="text-xl font-semibold text-gray-900 dark:text-white leading-tight mt-0.5 tabular-nums tracking-tight">
            {{ card.value }}
          </p>
        </div>
      </div>
    </div>

    <!-- 资源列表 -->
    <n-card :bordered="false" class="shadow-sm rounded-xl">
      <template #header>
        <div class="flex items-center gap-2.5">
          <span class="w-1 h-4 rounded-full bg-blue-500" aria-hidden="true"></span>
          <span class="font-semibold text-gray-900 dark:text-white">资源列表</span>
          <n-tag v-if="pagination.itemCount" :bordered="false" size="small" type="default" class="ml-1">
            {{ pagination.itemCount }}
          </n-tag>
        </div>
      </template>
      <template #header-extra>
        <n-space :size="12" align="center">
          <n-select
            v-model:value="filterStatus"
            :options="statusFilterOptions"
            placeholder="状态"
            clearable
            size="small"
            style="width: 120px"
            @update:value="handleFilterChange"
          />
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索标题"
            size="small"
            clearable
            style="width: 180px"
            @keyup.enter="handleFilterChange"
          />
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button size="small" quaternary @click="handleFilterChange" aria-label="搜索">
                <template #icon><i class="fas fa-search"></i></template>
              </n-button>
            </template>
            搜索
          </n-tooltip>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button size="small" quaternary @click="fetchResources" aria-label="刷新">
                <template #icon><i class="fas fa-refresh"></i></template>
              </n-button>
            </template>
            刷新
          </n-tooltip>
        </n-space>
      </template>

      <n-data-table
        :columns="columns"
        :data="resources"
        :pagination="pagination"
        :loading="loading"
        remote
        striped
      >
        <template #empty>
          <div class="py-10 flex flex-col items-center">
            <div class="w-14 h-14 rounded-2xl bg-gray-50 dark:bg-gray-700/50 flex items-center justify-center mb-4">
              <i class="fas fa-cloud-upload-alt text-gray-300 dark:text-gray-500 text-2xl"></i>
            </div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-300">
              {{ searchKeyword || filterStatus ? '没有匹配的资源' : '还没有提交过资源' }}
            </p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-1.5">
              {{ searchKeyword || filterStatus ? '调整筛选条件后重试' : '点击右上角「提交资源」，粘贴分享链接即可自动检测' }}
            </p>
          </div>
        </template>
      </n-data-table>
    </n-card>

    <!-- 提交资源弹窗 -->
    <n-modal
      v-model:show="showSubmitModal"
      preset="card"
      style="width: 600px"
      :bordered="false"
      :title="undefined"
    >
      <template #header>
        <div class="flex items-center gap-2.5">
          <span class="w-8 h-8 rounded-lg bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
            <i class="fas fa-cloud-upload-alt text-blue-500 dark:text-blue-400 text-sm"></i>
          </span>
          <div>
            <p class="font-semibold text-gray-900 dark:text-white leading-tight">提交资源</p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">提交后自动检测链接有效性</p>
          </div>
        </div>
      </template>

      <n-form
        ref="submitFormRef"
        :model="submitForm"
        :rules="submitRules"
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="资源标题" path="title">
          <n-input
            v-model:value="submitForm.title"
            placeholder="请输入资源标题"
            :maxlength="255"
            show-count
            @keyup.enter="handleSubmit"
          />
        </n-form-item>
        <n-form-item label="分享链接" path="url">
          <n-input
            v-model:value="submitForm.url"
            placeholder="粘贴网盘分享链接"
            @keyup.enter="handleSubmit"
          />
        </n-form-item>
        <n-form-item label="资源描述" path="description">
          <n-input
            v-model:value="submitForm.description"
            type="textarea"
            :rows="2"
            placeholder="可选，简单描述资源内容"
          />
        </n-form-item>
        <div class="flex items-center gap-2 px-3 py-2.5 rounded-lg bg-gray-50 dark:bg-gray-700/40 text-xs text-gray-500 dark:text-gray-400">
          <i class="fas fa-info-circle text-gray-400"></i>
          <span>支持百度 / 阿里 / 夸克 / 迅雷 / 天翼 / 123 / 115 / UC / 光鸭的分享链接</span>
        </div>
      </n-form>

      <template #footer>
        <div class="flex justify-end gap-3">
          <n-button @click="showSubmitModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            提交并检测
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 批量提交弹窗 -->
    <n-modal
      v-model:show="showBatchModal"
      preset="card"
      style="width: 600px"
      :bordered="false"
    >
      <template #header>
        <div class="flex items-center gap-2.5">
          <span class="w-8 h-8 rounded-lg bg-blue-50 dark:bg-blue-900/30 flex items-center justify-center">
            <i class="fas fa-layer-group text-blue-500 dark:text-blue-400 text-sm"></i>
          </span>
          <div>
            <p class="font-semibold text-gray-900 dark:text-white leading-tight">批量提交资源</p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">每行一条，逐条检测并返回结果</p>
          </div>
        </div>
      </template>

      <n-input
        v-model:value="batchText"
        type="textarea"
        :rows="8"
        placeholder="资源标题|https://pan.quark.cn/s/xxxxx&#10;https://pan.baidu.com/s/xxxxx"
      />
      <div class="flex items-center gap-2 px-3 py-2.5 mt-3 rounded-lg bg-gray-50 dark:bg-gray-700/40 text-xs text-gray-500 dark:text-gray-400">
        <i class="fas fa-info-circle text-gray-400"></i>
        <span>支持「标题|链接」格式，空行自动跳过，单次最多 50 条</span>
      </div>

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-xs text-gray-400">已识别 {{ parsedBatchLines.length }} 条</span>
          <div class="flex gap-3">
            <n-button @click="showBatchModal = false">取消</n-button>
            <n-button type="primary" :loading="batchSubmitting" :disabled="parsedBatchLines.length === 0" @click="handleBatchSubmit">
              提交 {{ parsedBatchLines.length }} 条
            </n-button>
          </div>
        </div>
      </template>
    </n-modal>

    <!-- 批量结果弹窗 -->
    <n-modal
      v-model:show="showBatchResultModal"
      preset="card"
      style="width: 560px"
      :bordered="false"
      title="批量提交结果"
    >
      <n-space vertical :size="12">
        <n-tag :type="batchResult.failCount === 0 ? 'success' : 'warning'" :bordered="false">
          成功 {{ batchResult.successCount }} 条，失败 {{ batchResult.failCount }} 条
        </n-tag>
        <div class="max-h-72 overflow-y-auto">
          <n-list :show-divider="true">
            <n-list-item v-for="(r, idx) in batchResult.results" :key="idx">
              <div class="flex items-center justify-between gap-3 py-1">
                <span class="text-sm text-gray-700 dark:text-gray-200 truncate">
                  {{ r.index + 1 }}. {{ r.title || '(无标题)' }}
                </span>
                <n-tag :type="r.ok ? 'success' : 'error'" size="small" :bordered="false" class="flex-shrink-0">
                  {{ r.ok ? `成功（${statusLabel(r.status)}）` : r.reason || '失败' }}
                </n-tag>
              </div>
            </n-list-item>
          </n-list>
        </div>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, computed } from 'vue'
import { NButton, NSpace, NTag, NTooltip, useDialog, useNotification } from 'naive-ui'
import { useApiFetch } from '~/composables/useApiFetch'

// 页面元数据
definePageMeta({
  layout: 'user',
  title: '我的资源'
})

// 在 setup 同步上下文获取 naive-ui 实例（async 处理器内调用 useNotification 会因脱离 setup 上下文而失败）
const notification = useNotification()
const dialog = useDialog()

// 常量：状态展示配置
const STATUS_META: Record<string, { label: string; type: 'default' | 'success' | 'error' | 'info' | 'warning' }> = {
  pending: { label: '未检测', type: 'warning' },
  valid: { label: '有效', type: 'info' },
  invalid: { label: '无效', type: 'error' },
  processing: { label: '处理中', type: 'info' },
  published: { label: '已公开', type: 'success' }
}

const statusLabel = (status: string) => STATUS_META[status]?.label || status

// 提交弹窗
const showSubmitModal = ref(false)
const submitFormRef = ref()
const submitting = ref(false)
const submitForm = ref({ title: '', description: '', url: '' })
const submitRules = {
  title: [{ required: true, message: '请输入资源标题', trigger: 'blur' }],
  url: [{ required: true, message: '请粘贴网盘分享链接', trigger: 'blur' }]
}

const openSubmitModal = () => {
  showSubmitModal.value = true
}

// 列表与筛选
const loading = ref(false)
const resources = ref<any[]>([])
const stats = ref<Record<string, number>>({})
const searchKeyword = ref('')
const filterStatus = ref<string | null>(null)

const statusFilterOptions = Object.entries(STATUS_META).map(([value, m]) => ({ label: m.label, value }))

// 统计卡配置
const statCards = computed(() => [
  {
    key: 'total', label: '总资源数', value: stats.value.total || 0,
    icon: 'fas fa-cloud', iconBg: 'bg-blue-50 dark:bg-blue-900/30', iconColor: 'text-blue-500 dark:text-blue-400'
  },
  {
    key: 'published', label: '已公开', value: stats.value.published || 0,
    icon: 'fas fa-check-circle', iconBg: 'bg-green-50 dark:bg-green-900/30', iconColor: 'text-green-500 dark:text-green-400'
  },
  {
    key: 'pending', label: '待处理',
    value: (stats.value.pending || 0) + (stats.value.processing || 0),
    icon: 'fas fa-hourglass-half', iconBg: 'bg-orange-50 dark:bg-orange-900/30', iconColor: 'text-orange-500 dark:text-orange-400'
  },
  {
    key: 'invalid', label: '无效', value: stats.value.invalid || 0,
    icon: 'fas fa-times-circle', iconBg: 'bg-red-50 dark:bg-red-900/30', iconColor: 'text-red-500 dark:text-red-400'
  }
])

const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.value.page = page
    fetchResources()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    fetchResources()
  }
})

// 批量提交
const showBatchModal = ref(false)
const showBatchResultModal = ref(false)
const batchText = ref('')
const batchSubmitting = ref(false)
const batchResult = ref<{ results: any[]; successCount: number; failCount: number }>({
  results: [],
  successCount: 0,
  failCount: 0
})

const parsedBatchLines = computed(() => {
  return batchText.value
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => line.length > 0)
    .map((line) => {
      const sep = line.lastIndexOf('|')
      if (sep > 0 && sep < line.length - 1) {
        return { title: line.slice(0, sep).trim(), url: line.slice(sep + 1).trim() }
      }
      return { title: `批量资源 ${new Date().toLocaleString('zh-CN')}`, url: line }
    })
})

// 表格列
const columns = [
  {
    title: '资源标题',
    key: 'title',
    render: (row: any) =>
      h('div', [
        h('div', { class: 'font-medium text-gray-900 dark:text-gray-100' }, row.title),
        row.description
          ? h('div', { class: 'text-xs text-gray-400 dark:text-gray-500 mt-0.5' }, row.description)
          : null
      ])
  },
  {
    title: '平台',
    key: 'pan',
    width: 110,
    render: (row: any) => row.pan?.remark || row.pan?.name || '未知'
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const meta = STATUS_META[row.status] || { label: row.status, type: 'default' as const }
      const tag = h(
        NTag,
        { type: meta.type, size: 'small', bordered: false, rounded: true },
        { default: () => meta.label }
      )
      if (row.status === 'invalid' && row.fail_reason) {
        return h(NTooltip, null, {
          trigger: () => tag,
          default: () => row.fail_reason
        })
      }
      return tag
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 170,
    render: (row: any) => h('span', { class: 'text-sm text-gray-500 dark:text-gray-400' }, formatDateTime(row.created_at))
  },
  {
    title: '操作',
    key: 'actions',
    width: 170,
    render: (row: any) => {
      const buttons: any[] = []
      if (row.status === 'pending' || row.status === 'invalid') {
        buttons.push(
          h(
            NButton,
            { size: 'small', type: 'info', secondary: true, onClick: () => handleRecheck(row) },
            { default: () => '重新检测' }
          )
        )
      }
      buttons.push(
        h(
          NButton,
          { size: 'small', type: 'error', secondary: true, onClick: () => handleDelete(row) },
          { default: () => '删除' }
        )
      )
      return h(NSpace, { size: 'small' }, { default: () => buttons })
    }
  }
]

// 格式化时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 获取资源列表
const fetchResources = async () => {
  loading.value = true
  try {
    // 注意：parseApiResponse 对 data.list 形状会退化为数组，这里直接解包 data
    const res = await useApiFetch('/user/resources', {
      params: {
        page: pagination.value.page,
        page_size: pagination.value.pageSize,
        status: filterStatus.value || '',
        keyword: searchKeyword.value || ''
      }
    }) as any
    const payload = res?.data || {}

    resources.value = payload.list || []
    pagination.value.itemCount = payload.total || 0
    if (payload.stats) {
      stats.value = payload.stats
    }
  } catch (error: any) {
    console.error('获取资源列表失败:', error)
    if (process.client) {
      notification.error({
        content: error?.message || '获取资源列表失败',
        duration: 3000
      })
    }
  } finally {
    loading.value = false
  }
}

// 筛选变化
const handleFilterChange = () => {
  pagination.value.page = 1
  fetchResources()
}

// 提交资源
const handleSubmit = async () => {
  try {
    await submitFormRef.value?.validate()
  } catch {
    return
  }
  submitting.value = true
  try {
    const res = await useApiFetch('/user/resources', {
      method: 'POST',
      body: {
        title: submitForm.value.title,
        description: submitForm.value.description,
        url: submitForm.value.url
      }
    }) as any
    const response = res?.data || {}

    notification[response?.status === 'invalid' ? 'warning' : 'success']({
      content: response?.message || '提交成功',
      duration: 4000
    })

    showSubmitModal.value = false
    submitForm.value = { title: '', description: '', url: '' }
    pagination.value.page = 1
    fetchResources()
  } catch (error: any) {
    notification.error({
      content: error?.data?.message || error?.message || '提交失败，请稍后重试',
      duration: 4000
    })
  } finally {
    submitting.value = false
  }
}

// 重新检测
const handleRecheck = async (row: any) => {
  try {
    const res = await useApiFetch(`/user/resources/${row.id}/recheck`, {
      method: 'POST'
    }) as any
    const response = res?.data || {}

    notification.info({
      content: `重新检测完成：${statusLabel(response?.status || '')}`,
      duration: 3000
    })
    fetchResources()
  } catch (error: any) {
    notification.error({
      content: error?.data?.message || '重新检测失败',
      duration: 3000
    })
  }
}

// 删除资源
const handleDelete = (row: any) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除资源「${row.title}」吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await useApiFetch(`/user/resources/${row.id}`, { method: 'DELETE' })
        notification.success({ content: '删除成功', duration: 3000 })
        fetchResources()
      } catch (error: any) {
        notification.error({
          content: error?.data?.message || '删除失败',
          duration: 3000
        })
      }
    }
  })
}

// 批量提交
const handleBatchSubmit = async () => {
  const items = parsedBatchLines.value
  if (items.length === 0) return
  if (items.length > 50) {
    notification.error({ content: '单次最多提交 50 条', duration: 3000 })
    return
  }

  batchSubmitting.value = true
  try {
    const res = await useApiFetch('/user/resources/batch', {
      method: 'POST',
      body: { items }
    }) as any
    const payload = res?.data || {}

    batchResult.value = {
      results: payload.results || [],
      successCount: payload.success_count || 0,
      failCount: payload.fail_count || 0
    }
    showBatchModal.value = false
    showBatchResultModal.value = true
    batchText.value = ''
    pagination.value.page = 1
    fetchResources()
  } catch (error: any) {
    notification.error({
      content: error?.data?.message || '批量提交失败',
      duration: 4000
    })
  } finally {
    batchSubmitting.value = false
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchResources()
})
</script>
