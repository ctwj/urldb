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
          :scroll-x="1020"
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
        <div class="mt-1 space-y-2">
          <div
            v-for="(file, index) in getProofFiles(selectedClaim.proof_files)"
            :key="index"
            class="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors cursor-pointer"
            @click="downloadFile(file)"
          >
            <div class="flex items-center space-x-2">
              <i class="fas fa-file-download text-blue-500"></i>
              <span class="text-sm text-gray-900 dark:text-gray-100">{{ getFileName(file) }}</span>
            </div>
            <i class="fas fa-download text-gray-400 hover:text-blue-500 transition-colors"></i>
          </div>
        </div>
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
    width: 60,
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        h('div', { class: 'font-medium text-sm' }, row.id),
        h('div', {
          class: 'text-xs text-gray-400',
          title: `IP: ${row.ip_address || '未知'}`
        }, row.ip_address ? `IP: ${row.ip_address.slice(0, 8)}...` : 'IP:未知')
      ])
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
    title: '申述人信息',
    key: 'claimant_info',
    width: 180,
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        // 第一行：姓名和身份
        h('div', { class: 'font-medium text-sm' }, [
          h('i', { class: 'fas fa-user text-green-500 mr-1 text-xs' }),
          row.claimant_name || '未知'
        ]),
        h('div', {
          class: 'text-xs text-blue-600 dark:text-blue-400 truncate max-w-[180px]',
          title: getIdentityLabel(row.identity)
        }, getIdentityLabel(row.identity)),
        // 第二行：联系方式
        h('div', {
          class: 'text-xs text-gray-500 dark:text-gray-400 truncate max-w-[180px]',
          title: row.contact_info
        }, [
          h('i', { class: 'fas fa-phone text-purple-500 mr-1' }),
          row.contact_info || '未提供'
        ]),
        // 第三行：证明类型
        h('div', {
          class: 'text-xs text-orange-600 dark:text-orange-400 truncate max-w-[180px]',
          title: getProofTypeLabel(row.proof_type)
        }, [
          h('i', { class: 'fas fa-certificate text-orange-500 mr-1' }),
          getProofTypeLabel(row.proof_type)
        ])
      ])
    }
  },
  {
    title: '申述详情',
    key: 'claim_details',
    width: 280,
    render: (row: any) => {
      return h('div', { class: 'space-y-1' }, [
        // 第一行：申述理由和提交时间
        h('div', { class: 'space-y-1' }, [
          h('div', { class: 'text-xs text-gray-500 dark:text-gray-400' }, '申述理由:'),
          h('div', {
            class: 'text-sm text-gray-700 dark:text-gray-300 line-clamp-2 max-h-10',
            title: row.reason
          }, row.reason || '无'),
          h('div', { class: 'text-xs text-gray-400' }, [
            h('i', { class: 'fas fa-clock mr-1' }),
            `提交时间: ${formatDateTime(row.created_at)}`
          ])
        ]),
        // 第二行：证明文件
        row.proof_files ?
          h('div', { class: 'space-y-1' }, [
            h('div', { class: 'text-xs text-gray-500 dark:text-gray-400' }, '证明文件:'),
            ...getProofFiles(row.proof_files).slice(0, 2).map((file, index) =>
              h('div', {
                class: 'text-xs text-blue-600 dark:text-blue-400 truncate max-w-[280px] cursor-pointer hover:text-blue-500 hover:underline',
                title: `点击下载: ${file}`,
                onClick: () => downloadFile(file)
              }, [
                h('i', { class: 'fas fa-download text-blue-500 mr-1' }),
                getFileName(file)
              ])
            ),
            getProofFiles(row.proof_files).length > 2 ?
              h('div', { class: 'text-xs text-gray-400' }, `还有 ${getProofFiles(row.proof_files).length - 2} 个文件...`) : null
          ]) :
          h('div', { class: 'text-xs text-gray-400' }, '无证明文件'),
        // 第三行：处理备注（如果有）
        row.note ?
          h('div', { class: 'space-y-1' }, [
            h('div', { class: 'text-xs text-gray-500 dark:text-gray-400' }, '处理备注:'),
            h('div', {
              class: 'text-xs text-yellow-600 dark:text-yellow-400 truncate max-w-[280px]',
              title: row.note
            }, [
              h('i', { class: 'fas fa-sticky-note text-yellow-500 mr-1' }),
              row.note.length > 30 ? `${row.note.slice(0, 30)}...` : row.note
            ])
          ]) : null
      ].filter(Boolean))
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const type = getStatusType(row.status)
      return h('div', { class: 'space-y-1' }, [
        h('n-tag', {
          type: type,
          size: 'small',
          bordered: false
        }, { default: () => getStatusLabel(row.status) }),
        // 显示处理时间（如果已处理）
        (row.status !== 'pending' && row.updated_at) ?
          h('div', {
            class: 'text-xs text-gray-400',
            title: `处理时间: ${formatDateTime(row.updated_at)}`
          }, `更新: ${new Date(row.updated_at).toLocaleDateString()}`) : null
      ].filter(Boolean))
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    render: (row: any) => {
      const buttons = [
        h('button', {
          class: 'px-2 py-1 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 dark:bg-blue-900/20 dark:text-blue-400 rounded transition-colors mb-1 w-full',
          onClick: () => viewClaim(row)
        }, [
          h('i', { class: 'fas fa-eye mr-1 text-xs' }),
          '查看详情'
        ])
      ]

      if (row.status === 'pending') {
        buttons.push(
          h('button', {
            class: 'px-2 py-1 text-xs bg-green-100 hover:bg-green-200 text-green-700 dark:bg-green-900/20 dark:text-green-400 rounded transition-colors mb-1 w-full',
            onClick: () => updateClaimStatus(row, 'approved')
          }, [
            h('i', { class: 'fas fa-check mr-1 text-xs' }),
            '批准'
          ]),
          h('button', {
            class: 'px-2 py-1 text-xs bg-red-100 hover:bg-red-200 text-red-700 dark:bg-red-900/20 dark:text-red-400 rounded transition-colors w-full',
            onClick: () => updateClaimStatus(row, 'rejected')
          }, [
            h('i', { class: 'fas fa-times mr-1 text-xs' }),
            '拒绝'
          ])
        )
      }

      return h('div', { class: 'flex flex-col gap-1' }, buttons)
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
    console.log(response)

    // 检查响应格式并处理
    if (response && response.data && response.data.list !== undefined) {
      // 如果后端返回了分页格式，使用正确的字段
      claims.value = response.data.list || []
      pagination.value.total = response.data.total || 0
    } else {
      // 如果是其他格式，尝试直接使用响应
      claims.value = response || []
      pagination.value.total = response.length || 0
    }
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

// 解析证明文件字符串
const getProofFiles = (proofFiles: string) => {
  if (!proofFiles) return []

  console.log('原始证明文件数据:', proofFiles)

  try {
    // 尝试解析为JSON格式
    const parsed = JSON.parse(proofFiles)
    console.log('JSON解析结果:', parsed)

    if (Array.isArray(parsed)) {
      // 处理对象数组格式：[{id: "xxx", name: "文件名.pdf", status: "pending"}]
      const fileObjects = parsed.filter(item => item && typeof item === 'object')
      if (fileObjects.length > 0) {
        // 返回原始对象，包含完整信息
        console.log('解析出文件对象数组:', fileObjects)
        return fileObjects
      }

      // 如果不是对象数组，尝试作为字符串数组处理
      const files = parsed.filter(file => file && typeof file === 'string' && file.trim()).map(file => file.trim())
      if (files.length > 0) {
        console.log('解析出的文件字符串数组:', files)
        return files
      }
    } else if (typeof parsed === 'object' && parsed.url) {
      console.log('解析出的单个文件:', parsed.url)
      return [parsed.url]
    } else if (typeof parsed === 'object' && parsed.files) {
      // 处理 {files: ["url1", "url2"]} 格式
      if (Array.isArray(parsed.files)) {
        const files = parsed.files.filter(file => file && file.trim()).map(file => file.trim())
        console.log('解析出的files数组:', files)
        return files
      }
    }
  } catch (e) {
    console.log('JSON解析失败，尝试分隔符解析:', e.message)
    // 如果不是JSON格式，按分隔符解析
    // 假设文件URL以逗号、分号或换行符分隔
    const files = proofFiles.split(/[,;\n\r]+/).filter(file => file.trim()).map(file => file.trim())
    console.log('分隔符解析结果:', files)
    return files
  }

  console.log('未解析出任何文件')
  return []
}

// 获取文件名
const getFileName = (fileInfo: any) => {
  if (!fileInfo) return '未知文件'

  // 如果是对象，优先使用name字段
  if (typeof fileInfo === 'object') {
    return fileInfo.name || fileInfo.id || '未知文件'
  }

  // 如果是字符串，从URL中提取文件名
  const fileName = fileInfo.split('/').pop() || fileInfo.split('\\').pop() || fileInfo

  // 如果URL太长，截断显示
  return fileName.length > 50 ? fileName.substring(0, 47) + '...' : fileName
}

// 下载文件
const downloadFile = async (fileInfo: any) => {
  console.log('尝试下载文件:', fileInfo)

  if (!fileInfo) {
    console.error('文件信息为空')
    if (process.client) {
      notification.warning({
        content: '文件信息无效',
        duration: 3000
      })
    }
    return
  }

  try {
    let downloadUrl = ''
    let fileName = ''

    // 处理文件对象格式：{id: "xxx", name: "文件名.pdf", status: "pending"}
    if (typeof fileInfo === 'object' && fileInfo.id) {
      fileName = fileInfo.name || fileInfo.id
      // 构建下载API URL，假设有 /api/files/{id} 端点
      downloadUrl = `/api/files/${fileInfo.id}`
      console.log('文件对象下载:', { id: fileInfo.id, name: fileName, url: downloadUrl })
    }
    // 处理字符串格式（直接是URL）
    else if (typeof fileInfo === 'string') {
      downloadUrl = fileInfo
      fileName = getFileName(fileInfo)

      // 检查是否是文件名（不包含http://或https://或/开头）
      if (!fileInfo.match(/^https?:\/\//) && !fileInfo.startsWith('/')) {
        console.log('检测到纯文件名，需要通过API下载:', fileName)

        if (process.client) {
          notification.info({
            content: `文件 "${fileName}" 需要通过API下载，功能开发中...`,
            duration: 3000
          })
        }
        return
      }

      // 处理相对路径URL
      if (fileInfo.startsWith('/uploads/')) {
        downloadUrl = `${window.location.origin}${fileInfo}`
        console.log('处理本地文件URL:', downloadUrl)
      }
    }

    if (!downloadUrl) {
      console.error('无法确定下载URL')
      if (process.client) {
        notification.warning({
          content: '无法确定下载地址',
          duration: 3000
        })
      }
      return
    }

    // 创建下载链接
    const link = document.createElement('a')
    link.href = downloadUrl
    link.target = '_blank' // 在新标签页打开，避免跨域问题

    // 设置下载文件名
    link.download = fileName.includes('.') ? fileName : fileName + '.file'

    console.log('下载参数:', {
      originalInfo: fileInfo,
      downloadUrl: downloadUrl,
      fileName: fileName
    })

    // 添加到页面并触发点击
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    if (process.client) {
      notification.success({
        content: `开始下载: ${fileName}`,
        duration: 2000
      })
    }
  } catch (error) {
    console.error('下载文件失败:', error)
    if (process.client) {
      notification.error({
        content: `下载失败: ${error.message}`,
        duration: 3000
      })
    }
  }
}

// 初始化数据
onMounted(() => {
  fetchClaims()
})
</script>