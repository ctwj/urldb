<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="container mx-auto px-4 py-8">
      <div class="max-w-4xl mx-auto">
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
        <VersionInfo />

        <!-- 版本历史 -->
        <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            <i class="fas fa-history mr-2 text-green-500"></i>
            版本历史
          </h3>
          
          <div class="space-y-4">
            <div v-for="(version, index) in versionHistory" :key="index" 
                 class="border-l-4 border-blue-500 pl-4 py-2">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="font-medium text-gray-900 dark:text-white">
                    v{{ version.version }}
                  </h4>
                  <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                    {{ version.date }}
                  </p>
                </div>
                <span class="px-2 py-1 text-xs rounded-full" 
                      :class="getVersionTypeClass(version.type)">
                  {{ version.type }}
                </span>
              </div>
              <ul class="mt-2 space-y-1">
                <li v-for="(change, changeIndex) in version.changes" :key="changeIndex"
                    class="text-sm text-gray-600 dark:text-gray-400 flex items-start">
                  <span class="mr-2 mt-1" :class="getChangeTypeClass(change.type)">
                    {{ getChangeTypeIcon(change.type) }}
                  </span>
                  {{ change.description }}
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- 构建信息 -->
        <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            <i class="fas fa-cogs mr-2 text-purple-500"></i>
            构建信息
          </h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded">
              <span class="text-sm text-gray-600 dark:text-gray-400">构建环境</span>
              <p class="font-mono text-gray-900 dark:text-white">Go 1.23.0</p>
            </div>
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded">
              <span class="text-sm text-gray-600 dark:text-gray-400">前端框架</span>
              <p class="font-mono text-gray-900 dark:text-white">Nuxt.js 3.8.0</p>
            </div>
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded">
              <span class="text-sm text-gray-600 dark:text-gray-400">数据库</span>
              <p class="font-mono text-gray-900 dark:text-white">PostgreSQL 15+</p>
            </div>
            <div class="p-3 bg-gray-50 dark:bg-gray-700 rounded">
              <span class="text-sm text-gray-600 dark:text-gray-400">部署方式</span>
              <p class="font-mono text-gray-900 dark:text-white">Docker</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin'
})

// 页面元数据
useHead({
  title: '版本信息 - 网盘资源数据库',
  meta: [
    { name: 'description', content: '查看系统版本信息和更新状态' }
  ]
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