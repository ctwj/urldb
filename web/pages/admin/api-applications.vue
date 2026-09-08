<template>
  <AdminPageLayout>
    <!-- 页面头部 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white flex items-center">
          <i class="fas fa-user-check text-blue-500 mr-2"></i>
          API 审核
        </h1>
        <p class="text-gray-600 dark:text-gray-400">审核用户的 API 开通申请，管理已开通用户的访问权限</p>
      </div>
      <div class="flex space-x-3">
        <n-button @click="refreshAll">
          <template #icon>
            <i class="fas fa-refresh"></i>
          </template>
          刷新
        </n-button>
      </div>
    </template>

    <!-- 通知区域 -->
    <template #notice-section>
      <n-alert title="同意申请后将自动为用户签发 API 调用凭证；已同意的申请可对其用户执行「停用权限」（立即失效）与「启用权限」（恢复调用），并可调整凭证有效期" type="info" />
    </template>

    <!-- 过滤栏 - 状态筛选与搜索 -->
    <template #filter-bar>
      <div class="flex justify-between items-center gap-3 flex-wrap">
        <div class="flex gap-2">
          <n-select
            v-model:value="activeTab"
            :options="statusOptions"
            placeholder="状态"
            style="width: 150px"
            @update:value="handleStatusChange"
          />
        </div>
        <div class="flex gap-2">
          <div class="relative">
            <n-input
              v-model:value="keyword"
              @input="debouncedSearch"
              @clear="debouncedSearch"
              type="text"
              placeholder="搜索用户名 / 邮箱 / 用途..."
              clearable
            >
              <template #prefix>
                <i class="fas fa-search text-gray-400 text-sm dark:text-gray-400"></i>
              </template>
            </n-input>
          </div>
          <n-button @click="resetFilters" type="tertiary">
            <template #icon>
              <i class="fas fa-redo"></i>
            </template>
            重置
          </n-button>
        </div>
      </div>
    </template>

    <!-- 内容区 -->
    <template #content>
      <!-- 加载状态 -->
      <div v-if="loading" class="flex h-full items-center justify-center py-8">
        <n-spin size="large" />
      </div>

      <!-- 错误状态 -->
      <AdminErrorState
        v-else-if="errorMessage"
        icon="fas fa-exclamation-triangle"
        :message="errorMessage"
        :on-retry="fetchList"
      />

      <!-- 空状态 -->
      <AdminEmptyState
        v-else-if="list.length === 0"
        icon="fas fa-user-check"
        title="暂无申请记录"
        :description="activeTab === 'pending' ? '当前没有待审核的 API 开通申请' : '当前状态下没有申请记录'"
      />

      <!-- 数据表格 - 自适应高度 -->
      <div v-else class="flex flex-col h-full overflow-auto">
        <n-data-table
          :columns="columns"
          :data="list"
          :pagination="false"
          :bordered="false"
          :single-line="false"
          :loading="loading"
          :scroll-x="1300"
          class="h-full"
        />
      </div>
    </template>

    <!-- 内容区footer - 分页组件 -->
    <template #content-footer>
      <div class="p-4">
        <div class="flex justify-center">
          <n-pagination
            v-model:page="page"
            v-model:page-size="pageSize"
            :item-count="total"
            :page-sizes="[20, 50, 100, 200]"
            show-size-picker
            @update:page="fetchList"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>
    </template>
  </AdminPageLayout>

  <!-- 拒绝理由 -->
  <n-modal v-model:show="showRejectModal" :mask-closable="false" preset="dialog" type="warning" title="拒绝申请"
    positive-text="确认拒绝" negative-text="取消"
    :loading="acting"
    @positive-click="confirmReject">
    <div class="space-y-3">
      <p class="text-sm text-gray-500 dark:text-gray-400">拒绝理由为可选项，将展示给申请人：</p>
      <n-input v-model:value="rejectReason" type="textarea" placeholder="例如：用途描述不清，请补充说明" :rows="3" />
    </div>
  </n-modal>

  <!-- 调整有效期 -->
  <n-modal v-model:show="showExpiryModal" :mask-closable="false" preset="dialog" type="info" title="调整 API 有效期"
    positive-text="保存" negative-text="取消"
    :loading="savingExpiry"
    @positive-click="confirmExpiry">
    <div class="space-y-4">
      <p class="text-sm text-gray-500 dark:text-gray-400">
        用户「{{ expiryTarget?.username }}」当前的到期时间：
        <span class="font-medium text-gray-700 dark:text-gray-200">{{ expiryTarget?.expires_at ? formatDateTime(expiryTarget.expires_at) : '永久有效' }}</span>
      </p>
      <div class="flex items-center gap-3">
        <n-switch v-model:value="expiryPermanent" size="small" />
        <span class="text-sm text-gray-700 dark:text-gray-200">永久有效</span>
      </div>
      <div v-if="!expiryPermanent">
        <p class="text-xs text-gray-400 dark:text-gray-500 mb-2">选择新的到期时间（设为过去的时间将立即过期）</p>
        <n-date-picker v-model:value="expiryTs" type="datetime" clearable placeholder="选择到期时间" style="width: 100%" />
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import { h, computed, onMounted, ref } from 'vue'
import { NButton, NSpace, NTag, NTooltip, useDialog, useNotification } from 'naive-ui'
import { useApiFetch } from '~/composables/useApiFetch'
import { useTimeFormat } from '~/composables/useTimeFormat'
import AdminPageLayout from '~/components/AdminPageLayout.vue'

definePageMeta({
  layout: 'admin',
  ssr: false
})

const notification = useNotification()
const dialog = useDialog()
const { formatDateTime } = useTimeFormat()

const loading = ref(false)
const acting = ref(false)
const errorMessage = ref('')
const list = ref<any[]>([])
const activeTab = ref('pending')
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const stats = ref<Record<string, number | null>>({ pending: null, approved: null, rejected: null })

const showRejectModal = ref(false)
const rejectReason = ref('')
const rejectTarget = ref<any>(null)

// 调整有效期
const showExpiryModal = ref(false)
const expiryTarget = ref<any>(null)
const expiryPermanent = ref(false)
const expiryTs = ref<number | null>(null)
const savingExpiry = ref(false)

const statusOptions = computed(() => [
  { label: `待审核${stats.value.pending != null ? ` (${stats.value.pending})` : ''}`, value: 'pending' },
  { label: `已同意${stats.value.approved != null ? ` (${stats.value.approved})` : ''}`, value: 'approved' },
  { label: `已拒绝${stats.value.rejected != null ? ` (${stats.value.rejected})` : ''}`, value: 'rejected' }
])

const STATUS_META: Record<string, { label: string; type: 'default' | 'success' | 'error' | 'warning'; icon: string }> = {
  pending: { label: '待审核', type: 'warning', icon: 'fas fa-hourglass-half' },
  approved: { label: '已同意', type: 'success', icon: 'fas fa-circle-check' },
  rejected: { label: '已拒绝', type: 'error', icon: 'fas fa-circle-xmark' }
}

const notify = (type: 'success' | 'error' | 'warning' | 'info', title: string, content: string) => {
  notification[type]({ title, content, duration: 4000 })
}

const fetchStats = async () => {
  try {
    const res: any = await useApiFetch('/admin/api-applications/stats')
    stats.value = {
      pending: res?.data?.pending ?? 0,
      approved: res?.data?.approved ?? 0,
      rejected: res?.data?.rejected ?? 0
    }
  } catch {
    // 统计失败不阻塞列表
  }
}

const fetchList = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res: any = await useApiFetch('/admin/api-applications', {
      params: { status: activeTab.value, keyword: keyword.value.trim(), page: page.value, page_size: pageSize.value }
    })
    list.value = res?.data?.list || []
    total.value = res?.data?.total || 0
  } catch (e: any) {
    errorMessage.value = e?.message || '获取申请列表失败，请稍后重试'
    list.value = []
  } finally {
    loading.value = false
  }
}

const refreshAll = () => {
  fetchStats()
  fetchList()
}

const handleStatusChange = () => {
  page.value = 1
  fetchList()
}

const resetFilters = () => {
  keyword.value = ''
  activeTab.value = 'pending'
  page.value = 1
  fetchList()
}

const handlePageSizeChange = () => {
  page.value = 1
  fetchList()
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debouncedSearch = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    fetchList()
  }, 350)
}

const approve = (row: any) => {
  dialog.warning({
    title: '同意申请',
    content: `确定同意用户「${row.username}」的 API 开通申请吗？同意后将自动签发调用凭证。`,
    positiveText: '确定同意',
    negativeText: '取消',
    onPositiveClick: async () => {
      acting.value = true
      try {
        await useApiFetch(`/admin/api-applications/${row.id}/approve`, { method: 'POST' })
        notify('success', '已同意', `用户「${row.username}」的 API 权限已开通`)
        refreshAll()
      } catch (e: any) {
        notify('error', '操作失败', e?.message || '请稍后重试')
      } finally {
        acting.value = false
      }
    }
  })
}

const openReject = (row: any) => {
  rejectTarget.value = row
  rejectReason.value = ''
  showRejectModal.value = true
}

const confirmReject = async () => {
  if (!rejectTarget.value) return
  acting.value = true
  try {
    await useApiFetch(`/admin/api-applications/${rejectTarget.value.id}/reject`, {
      method: 'POST',
      body: { reason: rejectReason.value }
    })
    notify('success', '已拒绝', '申请已拒绝，用户可查看理由后重新申请')
    showRejectModal.value = false
    refreshAll()
  } catch (e: any) {
    notify('error', '操作失败', e?.message || '请稍后重试')
  } finally {
    acting.value = false
  }
}

const disableUser = (row: any) => {
  dialog.warning({
    title: '停用 API 权限',
    content: `确定停用用户「${row.username}」的 API 权限吗？停用后其凭证立即失效。`,
    positiveText: '确定停用',
    negativeText: '取消',
    onPositiveClick: async () => {
      acting.value = true
      try {
        await useApiFetch(`/admin/api-credentials/${row.user_id}/disable`, { method: 'POST' })
        notify('success', '已停用', `用户「${row.username}」的 API 权限已停用`)
        await fetchList()
      } catch (e: any) {
        notify('error', '操作失败', e?.message || '该用户可能没有有效的 API 权限')
      } finally {
        acting.value = false
      }
    }
  })
}

// 启用 API 权限（恢复被停用的凭证，B6）
const enableUser = (row: any) => {
  dialog.warning({
    title: '启用 API 权限',
    content: `确定恢复用户「${row.username}」的 API 权限吗？启用后其凭证立即恢复调用能力${row.credential_expired ? '（该凭证已过期，启用后请调整有效期）' : ''}。`,
    positiveText: '确定启用',
    negativeText: '取消',
    onPositiveClick: async () => {
      acting.value = true
      try {
        await useApiFetch(`/admin/api-credentials/${row.user_id}/enable`, { method: 'POST' })
        notify('success', '已启用', `用户「${row.username}」的 API 权限已恢复`)
        await fetchList()
      } catch (e: any) {
        notify('error', '操作失败', e?.message || '该用户可能没有已停用的 API 凭证')
      } finally {
        acting.value = false
      }
    }
  })
}

// 调整凭证有效期（approved 行，B5）
const openExpiry = (row: any) => {
  expiryTarget.value = row
  if (row.expires_at) {
    expiryTs.value = new Date(row.expires_at).getTime()
    expiryPermanent.value = false
  } else {
    expiryTs.value = null
    expiryPermanent.value = true
  }
  showExpiryModal.value = true
}

const confirmExpiry = async () => {
  if (!expiryTarget.value) return false
  savingExpiry.value = true
  try {
    await useApiFetch(`/admin/api-credentials/${expiryTarget.value.user_id}/expiry`, {
      method: 'POST',
      body: {
        expires_at: !expiryPermanent.value && expiryTs.value
          ? new Date(expiryTs.value).toISOString()
          : null
      }
    })
    notify('success', '已调整', `用户「${expiryTarget.value.username}」的凭证有效期已更新`)
    showExpiryModal.value = false
    await fetchList()
  } catch (e: any) {
    notify('error', '操作失败', e?.message || '请稍后重试')
    return false
  } finally {
    savingExpiry.value = false
  }
}

const renderStatus = (row: any) => {
  const meta = STATUS_META[row.status] || { label: row.status, type: 'default' as const, icon: 'fas fa-tag' }
  return h(NTag, { type: meta.type, size: 'small', round: true }, {
    icon: () => h('i', { class: meta.icon }),
    default: () => meta.label
  })
}

const renderPurpose = (row: any) => {
  const text = row.purpose || '-'
  return h(NTooltip, { trigger: 'hover', style: 'max-width: 360px' }, {
    trigger: () => h('span', { class: 'text-gray-700 dark:text-gray-300' }, text.length > 40 ? text.slice(0, 40) + '…' : text),
    default: () => h('span', { class: 'whitespace-pre-wrap break-all' }, text)
  })
}

const columns = [
  {
    title: '申请人',
    key: 'username',
    width: 130,
    render: (row: any) => h('div', { class: 'flex items-center gap-2' }, [
      h('div', { class: 'w-7 h-7 rounded-full bg-blue-100 dark:bg-blue-900/40 flex items-center justify-center flex-shrink-0' }, [
        h('i', { class: 'fas fa-user text-blue-600 dark:text-blue-400 text-xs' })
      ]),
      h('span', { class: 'font-medium text-gray-800 dark:text-gray-200 truncate' }, row.username || '-')
    ])
  },
  { title: '邮箱', key: 'email', width: 190, ellipsis: { tooltip: true } },
  { title: '用途说明', key: 'purpose', minWidth: 200, render: renderPurpose },
  { title: '状态', key: 'status', width: 100, render: renderStatus },
  {
    title: '提交时间',
    key: 'created_at',
    width: 165,
    render: (row: any) => h('span', { class: 'text-gray-500 dark:text-gray-400 tabular-nums' }, formatDateTime(row.created_at))
  },
  {
    title: '到期时间',
    key: 'expires_at',
    width: 180,
    render: (row: any) => {
      if (row.status !== 'approved' || !row.credential_status) return h('span', { class: 'text-gray-400' }, '-')
      const cells: any[] = [
        h('span', { class: 'tabular-nums text-gray-600 dark:text-gray-300' },
          row.expires_at ? formatDateTime(row.expires_at) : '永久有效')
      ]
      if (row.credential_status === 'disabled') {
        cells.push(h(NTag, { type: 'error', size: 'tiny', class: 'ml-1.5' }, { default: () => '已停用' }))
      } else if (row.credential_expired) {
        cells.push(h(NTag, { type: 'error', size: 'tiny', class: 'ml-1.5' }, { default: () => '已过期' }))
      }
      return h('div', { class: 'flex items-center' }, cells)
    }
  },
  {
    title: '审核信息',
    key: 'review',
    width: 200,
    render: (row: any) => {
      if (row.status === 'pending') return h('span', { class: 'text-gray-400' }, '-')
      const parts: string[] = []
      if (row.reviewer_name) parts.push(row.reviewer_name)
      if (row.reviewed_at) parts.push(formatDateTime(row.reviewed_at))
      if (row.status === 'rejected' && row.reject_reason) {
        return h(NTooltip, { trigger: 'hover', style: 'max-width: 320px' }, {
          trigger: () => h('span', { class: 'text-gray-500 dark:text-gray-400' },
            (parts.join(' · ') || '-') + '（理由）'),
          default: () => h('span', { class: 'whitespace-pre-wrap' }, `拒绝理由：${row.reject_reason}`)
        })
      }
      return h('span', { class: 'text-gray-500 dark:text-gray-400' }, parts.join(' · ') || '-')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 260,
    render: (row: any) => {
      const buttons: any[] = []
      if (row.status === 'pending') {
        buttons.push(
          h(NButton, {
            size: 'small', type: 'primary', onClick: () => approve(row),
            'aria-label': `同意 ${row.username} 的申请`
          }, { icon: () => h('i', { class: 'fas fa-check' }), default: () => '同意' }),
          h(NButton, {
            size: 'small', type: 'error', secondary: true, onClick: () => openReject(row),
            'aria-label': `拒绝 ${row.username} 的申请`
          }, { icon: () => h('i', { class: 'fas fa-xmark' }), default: () => '拒绝' })
        )
      }
      if (row.status === 'approved') {
        if (row.credential_status === 'disabled') {
          // 已停用：展示恢复入口（若同时已过期，提示启用后需调整有效期）
          buttons.push(
            h(NButton, {
              size: 'small', type: 'success', secondary: true, onClick: () => enableUser(row),
              'aria-label': `启用 ${row.username} 的 API 权限`
            }, { icon: () => h('i', { class: 'fas fa-circle-play' }), default: () => '启用权限' })
          )
        } else {
          buttons.push(
            h(NButton, {
              size: 'small', type: 'info', secondary: true,
              disabled: row.credential_status === 'disabled',
              onClick: () => openExpiry(row),
              'aria-label': `调整 ${row.username} 的凭证有效期`
            }, { icon: () => h('i', { class: 'fas fa-calendar-days' }), default: () => '调整有效期' }),
            h(NButton, {
              size: 'small', type: 'warning', secondary: true, onClick: () => disableUser(row),
              'aria-label': `停用 ${row.username} 的 API 权限`
            }, { icon: () => h('i', { class: 'fas fa-ban' }), default: () => '停用权限' })
          )
        }
      }
      if (buttons.length === 0) return h('span', { class: 'text-gray-400' }, '-')
      return h(NSpace, { size: 'small', wrap: false }, { default: () => buttons })
    }
  }
]

onMounted(refreshAll)
</script>
