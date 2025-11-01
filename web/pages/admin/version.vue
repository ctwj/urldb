<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
        <i class="fas fa-code-branch mr-3 text-blue-500"></i>
        版本信息
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        查看系统版本信息和更新状态
      </p>
    </div>

    <!-- 版本信息组件 -->
    <div class="max-w-4xl mx-auto">
      <VersionInfo />
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin-dynamic',
  middleware: ['auth']
})



interface VersionChange {
  type: 'feature' | 'fix' | 'improvement' | 'breaking'
  description: string
}

interface VersionHistory {
  version: string
  date: string
  type: 'major' | 'minor' | 'patch'
  changes: VersionChange[]
}

const versionHistory: VersionHistory[] = [
  {
    version: '1.1.0',
    date: '2024-01-20',
    type: 'minor',
    changes: [
      { type: 'feature', description: '🎉 新增数据转存管理功能' },
      { type: 'feature', description: '🔍 新增SEO管理功能' },
      { type: 'feature', description: '📤 新增数据推送功能' },
      { type: 'feature', description: '🤖 新增机器人管理功能' },
      { type: 'feature', description: '📊 新增三方统计功能' },
      { type: 'improvement', description: '🎨 优化搜索统计页面界面' },
      { type: 'improvement', description: '📱 改进响应式设计' },
      { type: 'fix', description: '🐛 修复部分UI显示问题' }
    ]
  },
  {
    version: '1.0.0',
    date: '2024-01-15',
    type: 'major',
    changes: [
      { type: 'feature', description: '🎉 首次发布' },
      { type: 'feature', description: '📁 多平台网盘支持' },
      { type: 'feature', description: '🔍 智能搜索功能' },
      { type: 'feature', description: '📊 数据统计和分析' },
      { type: 'feature', description: '🏷️ 标签系统' },
      { type: 'feature', description: '👥 用户权限管理' },
      { type: 'feature', description: '📦 批量资源管理' },
      { type: 'feature', description: '🔄 自动处理功能' },
      { type: 'feature', description: '📈 热播剧管理' },
      { type: 'feature', description: '⚙️ 系统配置管理' },
      { type: 'feature', description: '🔐 JWT认证系统' },
      { type: 'feature', description: '📱 响应式设计' },
      { type: 'feature', description: '🌙 深色模式支持' },
      { type: 'feature', description: '🎨 现代化UI界面' }
    ]
  }
]

// 获取版本类型样式
const getVersionTypeClass = (type: string) => {
  switch (type) {
    case 'major':
      return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
    case 'minor':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
    case 'patch':
      return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
    default:
      return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
  }
}

// 获取变更类型样式
const getChangeTypeClass = (type: string) => {
  switch (type) {
    case 'feature':
      return 'text-green-600 dark:text-green-400'
    case 'fix':
      return 'text-red-600 dark:text-red-400'
    case 'improvement':
      return 'text-blue-600 dark:text-blue-400'
    case 'breaking':
      return 'text-orange-600 dark:text-orange-400'
    default:
      return 'text-gray-600 dark:text-gray-400'
  }
}

// 获取变更类型图标
const getChangeTypeIcon = (type: string) => {
  switch (type) {
    case 'feature':
      return '✨'
    case 'fix':
      return '🐛'
    case 'improvement':
      return '🔧'
    case 'breaking':
      return '💥'
    default:
      return '��'
  }
}
</script> 