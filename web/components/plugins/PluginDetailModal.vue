<template>
  <Teleport to="body">
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-hidden">
        <!-- 模态框头部 -->
        <div class="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <h2 class="text-xl font-semibold text-gray-900">{{ plugin.name }}</h2>
            <span class="px-2 py-1 text-xs font-medium rounded-full" :class="statusClass">
              {{ statusText }}
            </span>
          </div>
          <button
            @click="$emit('close')"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <Icon name="x" class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <!-- 模态框内容 -->
        <div class="p-6 overflow-y-auto max-h-[calc(90vh-80px)]">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- 左侧：基本信息 -->
            <div class="lg:col-span-2 space-y-6">
              <!-- 基本信息 -->
              <div class="bg-gray-50 rounded-lg p-4">
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="info" class="w-5 h-5 mr-2 text-blue-600" />
                  基本信息
                </h3>
                <dl class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <dt class="text-gray-600">插件名称</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.name }}</dd>
                  </div>
                  <div>
                    <dt class="text-gray-600">版本</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.version }}</dd>
                  </div>
                  <div>
                    <dt class="text-gray-600">作者</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.author || '未知' }}</dd>
                  </div>
                  <div>
                    <dt class="text-gray-600">许可证</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.license || '未指定' }}</dd>
                  </div>
                  <div>
                    <dt class="text-gray-600">分类</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.category }}</dd>
                  </div>
                  <div>
                    <dt class="text-gray-600">文件大小</dt>
                    <dd class="font-medium text-gray-900">{{ formatFileSize(plugin.file_size) }}</dd>
                  </div>
                  <div class="col-span-2">
                    <dt class="text-gray-600">描述</dt>
                    <dd class="font-medium text-gray-900">{{ plugin.description }}</dd>
                  </div>
                </dl>
              </div>

              <!-- 描述 -->
              <div>
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="file-text" class="w-5 h-5 mr-2 text-blue-600" />
                  详细描述
                </h3>
                <p class="text-gray-700 leading-relaxed">
                  {{ plugin.description || '暂无详细描述' }}
                </p>
              </div>

              <!-- 权限和依赖 -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-blue-50 rounded-lg p-4">
                  <h4 class="font-medium text-blue-900 mb-2 flex items-center">
                    <Icon name="shield" class="w-4 h-4 mr-2" />
                    所需权限
                  </h4>
                  <div v-if="plugin.permissions && plugin.permissions.length > 0" class="space-y-1">
                    <div
                      v-for="permission in plugin.permissions"
                      :key="permission"
                      class="text-sm text-blue-800 bg-blue-100 px-2 py-1 rounded"
                    >
                      {{ permission }}
                    </div>
                  </div>
                  <p v-else class="text-sm text-blue-700">无特殊权限要求</p>
                </div>

                <div class="bg-orange-50 rounded-lg p-4">
                  <h4 class="font-medium text-orange-900 mb-2 flex items-center">
                    <Icon name="package" class="w-4 h-4 mr-2" />
                    依赖插件
                  </h4>
                  <div v-if="plugin.dependencies && plugin.dependencies.length > 0" class="space-y-1">
                    <div
                      v-for="dependency in plugin.dependencies"
                      :key="dependency"
                      class="text-sm text-orange-800 bg-orange-100 px-2 py-1 rounded"
                    >
                      {{ dependency }}
                    </div>
                  </div>
                  <p v-else class="text-sm text-orange-700">无依赖</p>
                </div>
              </div>

              <!-- 钩子列表 -->
              <div v-if="plugin.hooks && plugin.hooks.length > 0">
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="link" class="w-5 h-5 mr-2 text-blue-600" />
                  注册的钩子
                </h3>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="hook in plugin.hooks"
                    :key="hook"
                    class="px-3 py-1 bg-purple-100 text-purple-800 text-sm rounded-full"
                  >
                    {{ hook }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 右侧：状态和统计 -->
            <div class="space-y-6">
              <!-- 运行状态 -->
              <div class="bg-white border border-gray-200 rounded-lg p-4">
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="activity" class="w-5 h-5 mr-2 text-green-600" />
                  运行状态
                </h3>
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">状态</span>
                    <span class="px-2 py-1 text-xs font-medium rounded-full" :class="statusClass">
                      {{ statusText }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">启用状态</span>
                    <span class="px-2 py-1 text-xs font-medium rounded-full" :class="plugin.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'">
                      {{ plugin.enabled ? '已启用' : '已禁用' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">最后更新</span>
                    <span class="text-sm text-gray-900">{{ formatDate(plugin.last_updated) }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">文件路径</span>
                    <span class="text-sm text-gray-900 font-mono text-xs">{{ plugin.file_path }}</span>
                  </div>
                </div>
              </div>

              <!-- 执行统计 -->
              <div v-if="plugin.execution_stats" class="bg-white border border-gray-200 rounded-lg p-4">
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="bar-chart" class="w-5 h-5 mr-2 text-blue-600" />
                  执行统计
                </h3>
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">总执行次数</span>
                    <span class="text-lg font-bold text-gray-900">{{ plugin.execution_stats.total_executions.toLocaleString() }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">成功率</span>
                    <span class="text-lg font-bold" :class="plugin.execution_stats.success_rate >= 95 ? 'text-green-600' : plugin.execution_stats.success_rate >= 80 ? 'text-yellow-600' : 'text-red-600'">
                      {{ plugin.execution_stats.success_rate.toFixed(1) }}%
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">平均耗时</span>
                    <span class="text-lg font-bold text-gray-900">{{ plugin.execution_stats.average_time }}ms</span>
                  </div>
                  <div v-if="plugin.execution_stats.last_execution" class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">最后执行</span>
                    <span class="text-sm text-gray-900">{{ formatDate(plugin.execution_stats.last_execution) }}</span>
                  </div>
                </div>

                <!-- 成功率进度条 -->
                <div class="mt-4">
                  <div class="flex items-center justify-between text-sm mb-1">
                    <span class="text-gray-600">成功率</span>
                    <span class="font-medium">{{ plugin.execution_stats.success_rate.toFixed(1) }}%</span>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2">
                    <div
                      class="h-2 rounded-full transition-all duration-300"
                      :class="plugin.execution_stats.success_rate >= 95 ? 'bg-green-500' : plugin.execution_stats.success_rate >= 80 ? 'bg-yellow-500' : 'bg-red-500'"
                      :style="{ width: plugin.execution_stats.success_rate + '%' }"
                    ></div>
                  </div>
                </div>
              </div>

              <!-- 快速操作 -->
              <div class="bg-white border border-gray-200 rounded-lg p-4">
                <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                  <Icon name="zap" class="w-5 h-5 mr-2 text-purple-600" />
                  快速操作
                </h3>
                <div class="space-y-2">
                  <button
                    @click="togglePlugin"
                    :class="plugin.enabled ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
                    class="w-full px-4 py-2 text-white rounded-lg transition-colors flex items-center justify-center"
                  >
                    <Icon :name="plugin.enabled ? 'x-circle' : 'check-circle'" class="w-4 h-4 mr-2" />
                    {{ plugin.enabled ? '禁用插件' : '启用插件' }}
                  </button>
                  <button
                    @click="viewLogs"
                    class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center justify-center"
                  >
                    <Icon name="file-text" class="w-4 h-4 mr-2" />
                    查看日志
                  </button>
                  <button
                    @click="configure"
                    class="w-full px-4 py-2 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors flex items-center justify-center"
                  >
                    <Icon name="settings" class="w-4 h-4 mr-2" />
                    配置插件
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 模态框底部 -->
        <div class="border-t border-gray-200 px-6 py-4 flex justify-end space-x-3">
          <button
            @click="$emit('close')"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            关闭
          </button>
          <button
            @click="configure"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            配置插件
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue'
import { useToast } from '~/composables/useToast'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  plugin: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'enable', 'disable', 'configure', 'viewLogs'])

const toast = useToast()

// 计算属性
const statusClass = computed(() => {
  if (!props.plugin.enabled) {
    return 'bg-gray-100 text-gray-800'
  }
  if (props.plugin.status === 'error') {
    return 'bg-red-100 text-red-800'
  }
  return 'bg-green-100 text-green-800'
})

const statusText = computed(() => {
  if (!props.plugin.enabled) {
    return '已禁用'
  }
  if (props.plugin.status === 'error') {
    return '错误'
  }
  return '运行中'
})

// 方法
const togglePlugin = () => {
  if (props.plugin.enabled) {
    emit('disable', props.plugin.name)
  } else {
    emit('enable', props.plugin.name)
  }
}

const configure = () => {
  emit('configure', props.plugin)
}

const viewLogs = () => {
  emit('viewLogs', props.plugin)
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
/* 模态框动画 */
.fixed > div {
  animation: modalFadeIn 0.3s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 进度条动画 */
.h-2.rounded-full {
  transition: width 0.6s ease;
}
</style>