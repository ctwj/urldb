<template>
  <AdminPageLayout>
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
          <i class="fas fa-flag text-red-500 mr-2"></i>
          举报管理
        </h1>
        <p class="text-gray-600 dark:text-gray-400">管理用户提交的资源举报信息</p>
      </div>
    </template>

    <!-- 过滤栏 - 搜索和操作 -->
    <template #filter-bar>
      <div class="flex justify-between items-center">
        <div class="flex gap-2">
          <!-- 空白区域用于按钮 -->
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input
              v-model:value="filters.resourceKey"
              @input="debounceSearch"
              type="text"
              placeholder="搜索资源Key..."
              clearable
            >
              <template #prefix>
                <i class="fas fa-search text-gray-400 text-sm"></i>
              </template>
            </n-input>
          </div>
          <n-select
            v-model:value="filters.status"
            :options="[
              { label: '全部状态', value: '' },
              { label: '待处理', value: 'pending' },
              { label: '已批准', value: 'approved' },
              { label: '已拒绝', value: 'rejected' }
            ]"
            placeholder="状态"
            clearable
            @update:value="fetchReports"
            style="width: 150px"
          />
          <n-button @click="resetFilters" type="tertiary">
            <template #icon>
              <i class="fas fa-redo"></i>
            </template>
            重置
          </n-button>
          <n-button @click="fetchReports" type="tertiary">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区 - 举报数据 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex h-full items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="reports.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无举报记录</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">目前没有用户提交的举报信息</div>
      </div>

      <!-- 数据表格 - 自适应高度 -->
      <div v-else class="flex flex-col h-full overflow-auto">
        <n-data-table
          :columns="columns"
          :data="reports"
          :pagination="false"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          :scroll-x="1200"
          class="h-full"
        />
      </div>
    </template>

    <!-- 内容区footer - 分页组件 -->
    <template #content-footer>
      <div class="p-4">
        <div class="flex justify-center">
          <n-pagination
            v-model:page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :item-count="pagination.total"
            :page-sizes="[50, 100, 200, 500]"
            show-size-picker
            @update:page="fetchReports"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </template>

  </AdminPageLayout>

  <!-- 查看举报详情模态框 -->
  <n-modal v-model:show="showDetailModal" :mask-closable="false" preset="card" :style="{ maxWidth: '600px', width: '90%' }" title="举报详情">
    <div v-if="selectedReport" class="space-y-4">
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">举报ID</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.id }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">资源Key</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.resource_key }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">举报原因</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ getReasonLabel(selectedReport.reason) }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">详细描述</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.description }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">联系方式</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.contact || '未提供' }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">提交时间</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDateTime(selectedReport.created_at) }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">IP地址</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.ip_address || '未知' }}</p>
      </div>
      <div v-if="selectedReport.note">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">处理备注</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedReport.note }}</p>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
// 设置页面标题和元信息
useHead({
  title: '举报管理 - 管理后台',
  meta: [
    { name: 'description', content: '管理用户提交的资源举报信息' }
  ]
})

// 设置页面布局和认证保护
definePageMeta({
  layout: 'admin',
  middleware: ['auth', 'admin']
})

import { h } from 'vue'
const message = useMessage()
const notification = useNotification()
const dialog = useDialog()

const { resourceApi } = useApi()
const loading = ref(false)
const reports = ref<any[]>([])
const showDetailModal = ref(false)
const selectedReport = ref<any>(null)

// 分页和筛选状态
const pagination = ref({
  page: 1,
  pageSize: 50,
  total: 0
})

const filters = ref({
  status: '',
  resourceKey: ''
})

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 30,
    render: (row: any) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '资源',
    key: 'resource_key',
    width: 200,
    render: (row: any) => {
      const resourceInfo = getResourceInfo(row);
      return h('div', { class: 'space-y-1' }, [
        // 第一行：标题（单行，省略号）
        h('div', {
          class: 'font-medium text-sm truncate max-w-[200px]',
          style: { maxWidth: '200px' },
          title: resourceInfo.title // 鼠标hover显示完整标题
        }, resourceInfo.title),
        // 第二行：详情（单行，省略号）
        h('div', {
          class: 'text-xs text-gray-500 dark:text-gray-400 truncate max-w-[200px]',
          style: { maxWidth: '200px' },
          title: resourceInfo.description // 鼠标hover显示完整描述
        }, resourceInfo.description),
        // 第三行：分类图片和链接数
        h('div', { class: 'flex items-center gap-1' }, [
          h('i', {
            class: `fas fa-${getCategoryIcon(resourceInfo.category)} text-blue-500 text-xs`,
            // 鼠标hover显示第一个资源的链接地址
            title: resourceInfo.resources.length > 0 ? `链接地址: ${resourceInfo.resources[0].save_url || resourceInfo.resources[0].url}` : `资源链接地址: ${row.resource_key}`
          }),
          h('span', { class: 'text-xs text-gray-400' }, `链接数: ${resourceInfo.resources.length}`)
        ])
      ])
    }
  },
  {
    title: '举报原因',
    key: 'reason',
    width: 100,
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        // 举报原因和描述提示
        h('div', {
          class: 'flex items-center gap-1 truncate max-w-[80px]',
          style: { maxWidth: '80px' }
        }, [
          h('span', null, getReasonLabel(row.reason)),
          // 添加描述提示图片
          h('i', {
            class: 'fas fa-info-circle text-blue-400 cursor-pointer text-xs ml-1',
            title: row.description  // 鼠标hover显示描述
          })
        ]),
        // 举报时间
        h('div', {
          class: 'text-xs text-gray-400 truncate max-w-[80px]',
          style: { maxWidth: '80px' }
        }, `举报时间: ${formatDateTime(row.created_at)}`),
        // 联系方式
        h('div', {
          class: 'text-xs text-gray-500 dark:text-gray-400 truncate max-w-[80px]',
          style: { maxWidth: '80px' }
        }, `联系方式: ${row.contact || '未提供'}`)
      ])
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 50,
    render: (row: any) => {
      const type = getStatusType(row.status)
      return h('n-tag', {
        type: type,
        size: 'small',
        bordered: false
      }, { default: () => getStatusLabel(row.status) })
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row: any) => {
      const buttons = [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors mr-1',
          onClick: () => viewReport(row)
        }, [
          h('i', { class: 'fas fa-eye mr-1 text-xs' }),
          '查看'
        ])
      ]

      if (row.status === 'pending') {
        buttons.push(
          h('button', {
            class: 'px-2 py-1 text-xs bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400 rounded transition-colors mr-1',
            onClick: () => updateReportStatus(row, 'approved')
          }, [
            h('i', { class: 'fas fa-check mr-1 text-xs' }),
            '批准'
          ]),
          h('button', {
            class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
            onClick: () => updateReportStatus(row, 'rejected')
          }, [
            h('i', { class: 'fas fa-times mr-1 text-xs' }),
            '拒绝'
          ])
        )
      }

      return h('div', { class: 'flex items-center gap-1' }, buttons)
    }
  }
]

// 搜索防抖
let searchTimeout: NodeJS.Timeout | null = null
const debounceSearch = () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  searchTimeout = setTimeout(() => {
    pagination.value.page = 1
    fetchReports()
  }, 300)
}

// 获取举报列表
const fetchReports = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }

    if (filters.value.status) params.status = filters.value.status
    if (filters.value.resourceKey) params.resource_key = filters.value.resourceKey

    // 使用原始API调用以获取完整的分页信息
    const rawResponse = await resourceApi.getReportsRaw(params)
    console.log(rawResponse)

    // 检查响应格式并处理
    if (rawResponse && rawResponse.data && rawResponse.data.list !== undefined) {
      // 如果后端返回了分页格式，使用正确的字段
      reports.value = rawResponse.data.list || []
      pagination.value.total = rawResponse.data.total || 0
    } else {
      // 如果是其他格式，尝试直接使用响应
      reports.value = rawResponse || []
      pagination.value.total = rawResponse.length || 0
    }
  } catch (error) {
    console.error('获取举报列表失败:', error)
    // 显示错误提示
    if (process.client) {
      notification.error({
        content: '获取举报列表失败',
        duration: 3000
      })
    }
  } finally {
    loading.value = false
  }
}

// 重置筛选条件
const resetFilters = () => {
  filters.value = {
    status: '',
    resourceKey: ''
  }
  pagination.value.page = 1
  fetchReports()
}

// 处理页面大小变化
const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  fetchReports()
}

// 查看举报详情
const viewReport = (report: any) => {
  selectedReport.value = report
  showDetailModal.value = true
}

// 更新举报状态
const updateReportStatus = async (report: any, status: string) => {
  try {
    // 获取处理备注（如果需要）
    let note = ''
    if (status === 'rejected') {
      note = await getRejectionNote()
      if (note === null) return // 用户取消操作
    }

    const response = await resourceApi.updateReport(report.id, {
      status,
      note
    })

    // 更新本地数据
    const index = reports.value.findIndex(r => r.id === report.id)
    if (index !== -1) {
      reports.value[index] = response
    }

    // 更新详情模态框中的数据
    if (selectedReport.value && selectedReport.value.id === report.id) {
      selectedReport.value = response
    }

    if (process.client) {
      notification.success({
        content: '状态更新成功',
        duration: 3000
      })
    }
  } catch (error) {
    console.error('更新举报状态失败:', error)
    if (process.client) {
      notification.error({
        content: '状态更新失败',
        duration: 3000
      })
    }
  }
}

// 获取拒绝原因输入
const getRejectionNote = (): Promise<string | null> => {
  return new Promise((resolve) => {
    // 使用naive-ui的dialog API
    const { dialog } = useDialog()

    let inputValue = ''

    dialog.warning({
      title: '输入拒绝原因',
      content: () => h(nInput, {
        value: inputValue,
        onUpdateValue: (value) => inputValue = value,
        placeholder: '请输入拒绝的原因...',
        type: 'textarea',
        rows: 4
      }),
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        if (!inputValue.trim()) {
          const { message } = useNotification()
          message.warning('请输入拒绝原因')
          return false // 不关闭对话框
        }
        resolve(inputValue)
      },
      onNegativeClick: () => {
        resolve(null)
      }
    })
  })
}

// 状态类型和标签
const getStatusType = (status: string) => {
  switch (status) {
    case 'pending': return 'warning'
    case 'approved': return 'success'
    case 'rejected': return 'error'
    default: return 'default'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'pending': return '待处理'
    case 'approved': return '已批准'
    case 'rejected': return '已拒绝'
    default: return status
  }
}

// 举报原因标签
const getReasonLabel = (reason: string) => {
  const reasonMap: Record<string, string> = {
    'link_invalid': '链接已失效',
    'download_failed': '资源无法下载',
    'content_mismatch': '资源内容不符',
    'malicious': '包含恶意软件',
    'copyright': '版权问题',
    'other': '其他问题'
  }
  return reasonMap[reason] || reason
}

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取分类图标
const getCategoryIcon = (category: string) => {
  if (!category) return 'folder';

  // 根据分类名称返回对应的图标
  const categoryMap: Record<string, string> = {
    '文档': 'file-alt',
    '文档资料': 'file-alt',
    '压缩包': 'file-archive',
    '图片': 'images',
    '视频': 'film',
    '音乐': 'music',
    '电子书': 'book',
    '软件': 'cogs',
    '应用': 'mobile-alt',
    '游戏': 'gamepad',
    '资料': 'folder',
    '其他': 'file',
    'folder': 'folder',
    'file': 'file'
  };

  return categoryMap[category] || 'folder';
}

// 获取资源信息显示
const getResourceInfo = (row: any) => {
  // 从后端返回的资源列表中获取信息
  const resources = row.resources || [];

  if (resources.length > 0) {
    // 如果有多个资源，可以选择第一个或合并信息
    const resource = resources[0];
    return {
      title: resource.title || `资源: ${row.resource_key}`,
      description: resource.description || `资源详情: ${row.resource_key}`,
      category: resource.category || 'folder',
      resources: resources // 返回所有资源用于显示链接数量等
    }
  } else {
    // 如果没有关联资源，使用默认值
    return {
      title: `资源: ${row.resource_key}`,
      description: `资源详情: ${row.resource_key}`,
      category: 'folder',
      resources: []
    }
  }
}

// 初始化数据
onMounted(() => {
  fetchReports()
})
</script>