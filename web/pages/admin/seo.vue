<template>
  <div class="p-6">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">SEO管理</h1>
      <p class="text-gray-600 dark:text-gray-400 mt-2">搜索引擎优化管理</p>
    </div>

    <!-- Tab导航 -->
    <n-tabs v-model:value="activeTab" type="line" animated>
      <n-tab-pane name="site-submit" tab="站点提交">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="mb-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">站点提交</h3>
            <p class="text-gray-600 dark:text-gray-400">向各大搜索引擎提交站点信息</p>
          </div>

          <!-- 搜索引擎列表 -->
          <div class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- 百度 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
                      <i class="fas fa-search text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">百度</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">baidu.com</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitToBaidu">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.baidu || '未提交' }}
                </div>
              </div>

              <!-- 谷歌 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-red-500 rounded flex items-center justify-center">
                      <i class="fas fa-globe text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">谷歌</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">google.com</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitToGoogle">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.google || '未提交' }}
                </div>
              </div>

              <!-- 必应 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-green-500 rounded flex items-center justify-center">
                      <i class="fas fa-search text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">必应</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">bing.com</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitToBing">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.bing || '未提交' }}
                </div>
              </div>

              <!-- 搜狗 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-orange-500 rounded flex items-center justify-center">
                      <i class="fas fa-search text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">搜狗</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">sogou.com</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitToSogou">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.sogou || '未提交' }}
                </div>
              </div>

              <!-- 神马搜索 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-purple-500 rounded flex items-center justify-center">
                      <i class="fas fa-mobile-alt text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">神马搜索</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">sm.cn</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitToShenma">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.shenma || '未提交' }}
                </div>
              </div>

              <!-- 360搜索 -->
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <div class="flex items-center space-x-3">
                    <div class="w-8 h-8 bg-green-600 rounded flex items-center justify-center">
                      <i class="fas fa-shield-alt text-white text-sm"></i>
                    </div>
                    <div>
                      <h4 class="font-medium text-gray-900 dark:text-white">360搜索</h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400">so.com</p>
                    </div>
                  </div>
                  <n-button size="small" type="primary" @click="submitTo360">
                    <template #icon>
                      <i class="fas fa-paper-plane"></i>
                    </template>
                    提交
                  </n-button>
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  最后提交时间：{{ lastSubmitTime.so360 || '未提交' }}
                </div>
              </div>
            </div>

            <!-- 批量提交 -->
            <div class="mt-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="font-medium text-blue-900 dark:text-blue-100">批量提交</h4>
                  <p class="text-sm text-blue-700 dark:text-blue-300 mt-1">
                    一键提交到所有支持的搜索引擎
                  </p>
                </div>
                <n-button type="primary" @click="submitToAll">
                  <template #icon>
                    <i class="fas fa-rocket"></i>
                  </template>
                  批量提交
                </n-button>
              </div>
            </div>
          </div>
        </div>
      </n-tab-pane>

      <n-tab-pane name="link-building" tab="外链建设">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="mb-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">外链建设</h3>
            <p class="text-gray-600 dark:text-gray-400">管理和监控外部链接建设情况</p>
          </div>

          <!-- 外链统计 -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
              <div class="flex items-center">
                <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                  <i class="fas fa-link text-blue-600 dark:text-blue-400"></i>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-gray-600 dark:text-gray-400">总外链数</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ linkStats.total }}</p>
                </div>
              </div>
            </div>

            <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
              <div class="flex items-center">
                <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                  <i class="fas fa-check text-green-600 dark:text-green-400"></i>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-gray-600 dark:text-gray-400">有效外链</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ linkStats.valid }}</p>
                </div>
              </div>
            </div>

            <div class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
              <div class="flex items-center">
                <div class="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                  <i class="fas fa-clock text-yellow-600 dark:text-yellow-400"></i>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-gray-600 dark:text-gray-400">待审核</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ linkStats.pending }}</p>
                </div>
              </div>
            </div>

            <div class="bg-red-50 dark:bg-red-900/20 rounded-lg p-4">
              <div class="flex items-center">
                <div class="p-2 bg-red-100 dark:bg-red-900 rounded-lg">
                  <i class="fas fa-times text-red-600 dark:text-red-400"></i>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-gray-600 dark:text-gray-400">失效外链</p>
                  <p class="text-xl font-bold text-gray-900 dark:text-white">{{ linkStats.invalid }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 外链列表 -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h4 class="text-lg font-medium text-gray-900 dark:text-white">外链列表</h4>
              <n-button type="primary" @click="addNewLink">
                <template #icon>
                  <i class="fas fa-plus"></i>
                </template>
                添加外链
              </n-button>
            </div>

            <n-data-table
              :columns="linkColumns"
              :data="linkList"
              :pagination="linkPagination"
              :loading="linkLoading"
              :bordered="false"
              striped
            />
          </div>
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
// SEO管理页面
definePageMeta({
  layout: 'admin'
})

import { useMessage } from 'naive-ui'

// 获取消息组件
const message = useMessage()

// 当前激活的Tab
const activeTab = ref('site-submit')

// 最后提交时间
const lastSubmitTime = ref({
  baidu: '',
  google: '',
  bing: '',
  sogou: '',
  shenma: '',
  so360: ''
})

// 外链统计
const linkStats = ref({
  total: 156,
  valid: 142,
  pending: 8,
  invalid: 6
})

// 外链列表
const linkList = ref([
  {
    id: 1,
    url: 'https://example1.com',
    title: '示例外链1',
    status: 'valid',
    domain: 'example1.com',
    created_at: '2024-01-15'
  },
  {
    id: 2,
    url: 'https://example2.com',
    title: '示例外链2',
    status: 'pending',
    domain: 'example2.com',
    created_at: '2024-01-16'
  },
  {
    id: 3,
    url: 'https://example3.com',
    title: '示例外链3',
    status: 'invalid',
    domain: 'example3.com',
    created_at: '2024-01-17'
  }
])

const linkLoading = ref(false)

// 分页配置
const linkPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    linkPagination.value.page = page
    loadLinkList()
  },
  onUpdatePageSize: (pageSize: number) => {
    linkPagination.value.pageSize = pageSize
    linkPagination.value.page = 1
    loadLinkList()
  }
})

// 表格列配置
const linkColumns = [
  {
    title: 'URL',
    key: 'url',
    width: 300,
    render: (row: any) => {
      return h('a', {
        href: row.url,
        target: '_blank',
        class: 'text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300'
      }, row.url)
    }
  },
  {
    title: '标题',
    key: 'title',
    width: 200
  },
  {
    title: '域名',
    key: 'domain',
    width: 150
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const statusMap = {
        valid: { text: '有效', class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
        pending: { text: '待审核', class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' },
        invalid: { text: '失效', class: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' }
      }
      const status = statusMap[row.status as keyof typeof statusMap]
      return h('span', {
        class: `px-2 py-1 text-xs font-medium rounded ${status.class}`
      }, status.text)
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 120
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row: any) => {
      return h('div', { class: 'space-x-2' }, [
        h('button', {
          class: 'text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300',
          onClick: () => editLink(row)
        }, '编辑'),
        h('button', {
          class: 'text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300',
          onClick: () => deleteLink(row)
        }, '删除')
      ])
    }
  }
]

// 提交到百度
const submitToBaidu = () => {
  // 模拟提交
  lastSubmitTime.value.baidu = new Date().toLocaleString('zh-CN')
  message.success('已提交到百度')
}

// 提交到谷歌
const submitToGoogle = () => {
  // 模拟提交
  lastSubmitTime.value.google = new Date().toLocaleString('zh-CN')
  message.success('已提交到谷歌')
}

// 提交到必应
const submitToBing = () => {
  // 模拟提交
  lastSubmitTime.value.bing = new Date().toLocaleString('zh-CN')
  message.success('已提交到必应')
}

// 提交到搜狗
const submitToSogou = () => {
  // 模拟提交
  lastSubmitTime.value.sogou = new Date().toLocaleString('zh-CN')
  message.success('已提交到搜狗')
}

// 提交到神马搜索
const submitToShenma = () => {
  // 模拟提交
  lastSubmitTime.value.shenma = new Date().toLocaleString('zh-CN')
  message.success('已提交到神马搜索')
}

// 提交到360搜索
const submitTo360 = () => {
  // 模拟提交
  lastSubmitTime.value.so360 = new Date().toLocaleString('zh-CN')
  message.success('已提交到360搜索')
}

// 批量提交
const submitToAll = () => {
  // 模拟批量提交
  lastSubmitTime.value.baidu = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.google = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.bing = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.sogou = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.shenma = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.so360 = new Date().toLocaleString('zh-CN')
  message.success('已批量提交到所有搜索引擎')
}

// 加载外链列表
const loadLinkList = () => {
  // 模拟加载数据
  linkLoading.value = true
  setTimeout(() => {
    linkLoading.value = false
  }, 1000)
}

// 添加新外链
const addNewLink = () => {
  message.info('添加外链功能开发中')
}

// 编辑外链
const editLink = (row: any) => {
  message.info(`编辑外链: ${row.title}`)
}

// 删除外链
const deleteLink = (row: any) => {
  message.warning(`删除外链: ${row.title}`)
}

// 初始化
onMounted(() => {
  loadLinkList()
})
</script> 