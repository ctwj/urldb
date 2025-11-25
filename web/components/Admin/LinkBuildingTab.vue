<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">外链建设（待开发）</h3>
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
          <n-button type="primary" @click="$emit('add-new-link')">
            <template #icon>
              <i class="fas fa-plus"></i>
            </template>
            添加外链
          </n-button>
        </div>

        <n-data-table
          :columns="linkColumns"
          :data="linkList"
          :pagination="pagination"
          :loading="loading"
          :bordered="false"
          striped
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h } from 'vue'

// Props
interface Props {
  linkStats: {
    total: number
    valid: number
    pending: number
    invalid: number
  }
  linkList: Array<{
    id: number
    url: string
    title: string
    status: string
    domain: string
    created_at: string
  }>
  loading: boolean
  pagination: any
}

const props = withDefaults(defineProps<Props>(), {
  linkStats: () => ({
    total: 0,
    valid: 0,
    pending: 0,
    invalid: 0
  }),
  linkList: () => [],
  loading: false,
  pagination: () => ({})
})

// Emits
const emit = defineEmits<{
  'add-new-link': []
  'edit-link': [row: any]
  'delete-link': [row: any]
  'load-link-list': [page: number]
}>()

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
          onClick: () => emit('edit-link', row)
        }, '编辑'),
        h('button', {
          class: 'text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300',
          onClick: () => emit('delete-link', row)
        }, '删除')
      ])
    }
  }
]
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>