<template>
  <AdminPageLayout>
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
          <i class="fas fa-balance-scale text-blue-500 mr-2"></i>
          版权申述管理
        </h1>
        <p class="text-gray-600 dark:text-gray-400">管理用户提交的版权申述信息</p>
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
            @update:value="fetchClaims"
            style="width: 150px"
          />
          <n-button @click="resetFilters" type="tertiary">
            <template #icon>
              <i class="fas fa-redo"></i>
            </template>
            重置
          </n-button>
          <n-button @click="fetchClaims" type="tertiary">
            <template #icon>
              <i class="fas fa-refresh"></i>
            </template>
            刷新
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区 - 版权申述数据 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex h-full items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <!-- 空状态 -->
      <div v-else-if="claims.length === 0" class="text-center py-8">
        <svg class="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <circle cx="24" cy="24" r="20" stroke-width="3" stroke-dasharray="6 6" />
          <path d="M16 24h16M24 16v16" stroke-width="3" stroke-linecap="round" />
        </svg>
        <div class="text-lg font-semibold text-gray-400 dark:text-gray-500 mb-2">暂无版权申述记录</div>
        <div class="text-sm text-gray-400 dark:text-gray-600 mb-4">目前没有用户提交的版权申述信息</div>
      </div>

      <!-- 数据表格 - 自适应高度 -->
      <div v-else class="flex flex-col h-full overflow-auto">
        <n-data-table
          :columns="columns"
          :data="claims"
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
            @update:page="fetchClaims"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </template>

  </AdminPageLayout>

  <!-- 查看申述详情模态框 -->
  <n-modal v-model:show="showDetailModal" :mask-closable="false" preset="card" :style="{ maxWidth: '600px', width: '90%' }" title="版权申述详情">
    <div v-if="selectedClaim" class="space-y-4">
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">申述ID</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.id }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">资源Key</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.resource_key }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">申述人身份</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ getIdentityLabel(selectedClaim.identity) }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">证明类型</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ getProofTypeLabel(selectedClaim.proof_type) }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">申述理由</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.reason }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">联系方式</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.contact_info }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">申述人姓名</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.claimant_name }}</p>
      </div>
      <div v-if="selectedClaim.proof_files">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">证明文件</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100 break-all">{{ selectedClaim.proof_files }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">提交时间</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDateTime(selectedClaim.created_at) }}</p>
      </div>
      <div>
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">IP地址</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.ip_address || '未知' }}</p>
      </div>
      <div v-if="selectedClaim.note">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">处理备注</h3>
        <p class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ selectedClaim.note }}</p>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
// 设置页面标题和元信息
useHead({
  title: '版权申述管理 - 管理后台',
  meta: [
    { name: 'description', content: '管理用户提交的版权申述信息' }
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
const claims = ref<any[]>([])
const showDetailModal = ref(false)
const selectedClaim = ref<any>(null)

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
    width: 80,
    render: (row: any) => {
      return h('span', { class: 'font-medium' }, row.id)
    }
  },
  {
    title: '资源Key',
    key: 'resource_key',
    width: 180,
    render: (row: any) => {
      return h('n-tag', {
        type: 'info',
        size: 'small',
        class: 'truncate max-w-xs'
      }, { default: () => row.resource_key })
    }
  },
  {
    title: '申述人身份',
    key: 'identity',
    width: 120,
    render: (row: any) => {
      return h('span', null, getIdentityLabel(row.identity))
    }
  },
  {
    title: '证明类型',
    key: 'proof_type',
    width: 140,
    render: (row: any) => {
      return h('span', null, getProofTypeLabel(row.proof_type))
    }
  },
  {
    title: '申述人姓名',
    key: 'claimant_name',
    width: 120,
    render: (row: any) => {
      return h('span', null, row.claimant_name)
    }
  },
  {
    title: '联系方式',
    key: 'contact_info',
    width: 150,
    render: (row: any) => {
      return h('span', null, row.contact_info)
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
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
    title: '提交时间',
    key: 'created_at',
    width: 180,
    render: (row: any) => {
      return h('span', null, formatDateTime(row.created_at))
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
          onClick: () => viewClaim(row)
        }, [
          h('i', { class: 'fas fa-eye mr-1 text-xs' }),
          '查看'
        ])
      ]

      if (row.status === 'pending') {
        buttons.push(
          h('button', {
            class: 'px-2 py-1 text-xs bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400 rounded transition-colors mr-1',
            onClick: () => updateClaimStatus(row, 'approved')
          }, [
            h('i', { class: 'fas fa-check mr-1 text-xs' }),
            '批准'
          ]),
          h('button', {
            class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors',
            onClick: () => updateClaimStatus(row, 'rejected')
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
    fetchClaims()
  }, 300)
}

// 获取版权申述列表
const fetchClaims = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }

    if (filters.value.status) params.status = filters.value.status
    if (filters.value.resourceKey) params.resource_key = filters.value.resourceKey

    const response = await resourceApi.getCopyrightClaims(params)
    claims.value = response.items || []
    pagination.value.total = response.total || 0
  } catch (error) {
    console.error('获取版权申述列表失败:', error)
    // 显示错误提示
    if (process.client) {
      notification.error({
        content: '获取版权申述列表失败',
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
  fetchClaims()
}

// 处理页面大小变化
const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  fetchClaims()
}

// 查看申述详情
const viewClaim = (claim: any) => {
  selectedClaim.value = claim
  showDetailModal.value = true
}

// 更新申述状态
const updateClaimStatus = async (claim: any, status: string) => {
  try {
    // 获取处理备注（如果需要）
    let note = ''
    if (status === 'rejected') {
      note = await getRejectionNote()
      if (note === null) return // 用户取消操作
    }

    const response = await resourceApi.updateCopyrightClaim(claim.id, {
      status,
      note
    })

    // 更新本地数据
    const index = claims.value.findIndex(c => c.id === claim.id)
    if (index !== -1) {
      claims.value[index] = response
    }

    // 更新详情模态框中的数据
    if (selectedClaim.value && selectedClaim.value.id === claim.id) {
      selectedClaim.value = response
    }

    if (process.client) {
      notification.success({
        content: '状态更新成功',
        duration: 3000
      })
    }
  } catch (error) {
    console.error('更新版权申述状态失败:', error)
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

// 申述人身份标签
const getIdentityLabel = (identity: string) => {
  const identityMap: Record<string, string> = {
    'copyright_owner': '版权所有者',
    'authorized_agent': '授权代表',
    'law_firm': '律师事务所',
    'other': '其他'
  }
  return identityMap[identity] || identity
}

// 证明类型标签
const getProofTypeLabel = (proofType: string) => {
  const proofTypeMap: Record<string, string> = {
    'copyright_certificate': '版权登记证书',
    'first_publish_proof': '作品首发证明',
    'authorization_letter': '授权委托书',
    'identity_document': '身份证明文件',
    'other_proof': '其他证明材料'
  }
  return proofTypeMap[proofType] || proofType
}

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 初始化数据
onMounted(() => {
  fetchClaims()
})
</script>