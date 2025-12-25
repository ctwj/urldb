<template>
  <Teleport to="body">
    <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-hidden">
        <!-- 模态框头部 -->
        <div class="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-900">日志详情</h2>
            <p class="text-sm text-gray-600 mt-1">ID: {{ log.id }}</p>
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
          <div class="space-y-6">
            <!-- 基本信息 -->
            <div class="bg-gray-50 rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                <Icon name="info" class="w-5 h-5 mr-2 text-blue-600" />
                基本信息
              </h3>
              <dl class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <dt class="text-gray-600">插件名称</dt>
                  <dd class="font-medium text-gray-900">{{ log.plugin_name }}</dd>
                </div>
                <div>
                  <dt class="text-gray-600">钩子名称</dt>
                  <dd class="font-medium text-gray-900">{{ log.hook_name }}</dd>
                </div>
                <div>
                  <dt class="text-gray-600">执行状态</dt>
                  <dd class="font-medium">
                    <span
                      class="px-2 py-1 text-xs font-medium rounded-full"
                      :class="log.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                    >
                      {{ log.success ? '成功' : '失败' }}
                    </span>
                  </dd>
                </div>
                <div>
                  <dt class="text-gray-600">执行时间</dt>
                  <dd class="font-medium">
                    <span :class="getExecutionTimeClass(log.execution_time)">
                      {{ log.execution_time }}ms
                    </span>
                  </dd>
                </div>
                <div>
                  <dt class="text-gray-600">创建时间</dt>
                  <dd class="font-medium text-gray-900">{{ formatFullTime(log.created_at) }}</dd>
                </div>
                <div>
                  <dt class="text-gray-600">日志ID</dt>
                  <dd class="font-mono text-sm text-gray-900">{{ log.id }}</dd>
                </div>
              </dl>
            </div>

            <!-- 性能分析 -->
            <div v-if="!log.success || log.execution_time > 100" class="border rounded-lg p-4" :class="log.success ? 'border-yellow-200 bg-yellow-50' : 'border-red-200 bg-red-50'">
              <h3 class="text-lg font-medium mb-3 flex items-center" :class="log.success ? 'text-yellow-900' : 'text-red-900'">
                <Icon :name="log.success ? 'alert-triangle' : 'alert-circle'" class="w-5 h-5 mr-2" />
                {{ log.success ? '性能警告' : '错误分析' }}
              </h3>

              <div v-if="!log.success && log.error_message" class="space-y-2">
                <div class="bg-white rounded-lg p-3">
                  <h4 class="font-medium text-red-900 mb-2">错误信息</h4>
                  <pre class="text-sm text-red-800 whitespace-pre-wrap bg-red-50 p-3 rounded border border-red-200 font-mono">{{ log.error_message }}</pre>
                </div>

                <div class="text-sm text-red-700">
                  <p class="font-medium mb-1">建议解决方案：</p>
                  <ul class="list-disc list-inside space-y-1">
                    <li>检查插件代码是否有语法错误</li>
                    <li>验证插件依赖是否正确安装</li>
                    <li>确认数据库连接是否正常</li>
                    <li>查看插件配置是否有效</li>
                  </ul>
                </div>
              </div>

              <div v-else-if="log.execution_time > 1000">
                <div class="text-sm text-yellow-700">
                  <p class="font-medium mb-1">性能警告：</p>
                  <ul class="list-disc list-inside space-y-1">
                    <li>执行时间超过1秒，可能影响系统性能</li>
                    <li>建议优化插件代码或减少数据库查询</li>
                    <li>考虑使用缓存或异步处理</li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- 执行上下文 -->
            <div class="bg-white border border-gray-200 rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3 flex items-center">
                <Icon name="code" class="w-5 h-5 mr-2 text-blue-600" />
                执行上下文
              </h3>

              <div class="space-y-4">
                <!-- 系统信息 -->
                <div>
                  <h4 class="font-medium text-gray-900 mb-2">系统环境</h4>
                  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div class="bg-gray-50 p-2 rounded">
                      <span class="text-gray-600">用户代理:</span>
                      <div class="font-mono text-xs mt-1 truncate">{{ navigator.userAgent }}</div>
                    </div>
                    <div class="bg-gray-50 p-2 rounded">
                      <span class="text-gray-600">时区:</span>
                      <div class="font-mono text-xs mt-1">{{ Intl.DateTimeFormat().resolvedOptions().timeZone }}</div>
                    </div>
                    <div class="bg-gray-50 p-2 rounded">
                      <span class="text-gray-600">语言:</span>
                      <div class="font-mono text-xs mt-1">{{ navigator.language }}</div>
                    </div>
                    <div class="bg-gray-50 p-2 rounded">
                      <span class="text-gray-600">平台:</span>
                      <div class="font-mono text-xs mt-1">{{ navigator.platform }}</div>
                    </div>
                  </div>
                </div>

                <!-- 内存使用情况 -->
                <div>
                  <h4 class="font-medium text-gray-900 mb-2">内存使用</h4>
                  <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                    <div class="bg-blue-50 p-3 rounded">
                      <div class="text-blue-600 font-medium">已使用内存</div>
                      <div class="text-blue-900 font-bold text-lg">{{ getMemoryUsage().used }}</div>
                    </div>
                    <div class="bg-green-50 p-3 rounded">
                      <div class="text-green-600 font-medium">可用内存</div>
                      <div class="text-green-900 font-bold text-lg">{{ getMemoryUsage().available }}</div>
                    </div>
                    <div class="bg-purple-50 p-3 rounded">
                      <div class="text-purple-600 font-medium">内存使用率</div>
                      <div class="text-purple-900 font-bold text-lg">{{ getMemoryUsage().percentage }}%</div>
                    </div>
                  </div>
                </div>

                <!-- 相关日志 -->
                <div>
                  <h4 class="font-medium text-gray-900 mb-2">相关日志</h4>
                  <div class="text-sm text-gray-600">
                    <p>显示同一插件的其他执行记录</p>
                    <!-- 这里可以加载相关的日志记录 -->
                    <div class="mt-2 p-3 bg-gray-50 rounded text-center">
                      <Icon name="loader" class="w-4 h-4 inline animate-spin" />
                      <span class="ml-2">加载相关日志...</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 操作建议 -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h3 class="text-lg font-medium text-blue-900 mb-3 flex items-center">
                <Icon name="lightbulb" class="w-5 h-5 mr-2" />
                操作建议
              </h3>
              <div class="space-y-2 text-sm text-blue-800">
                <div v-if="!log.success">
                  <p class="font-medium">错误处理建议：</p>
                  <ul class="list-disc list-inside space-y-1 ml-4">
                    <li>检查插件代码中的异常处理逻辑</li>
                    <li>验证所有外部依赖的可用性</li>
                    <li>增加详细的错误日志记录</li>
                    <li>考虑实现重试机制</li>
                  </ul>
                </div>
                <div v-else-if="log.execution_time > 1000">
                  <p class="font-medium">性能优化建议：</p>
                  <ul class="list-disc list-inside space-y-1 ml-4">
                    <li>分析插件代码的性能瓶颈</li>
                    <li>优化数据库查询，添加适当的索引</li>
                    <li>考虑使用缓存减少重复计算</li>
                    <li>将耗时操作移到异步处理</li>
                  </ul>
                </div>
                <div v-else>
                  <p class="font-medium">正常维护建议：</p>
                  <ul class="list-disc list-inside space-y-1 ml-4">
                    <li>定期检查插件日志，及时发现潜在问题</li>
                    <li>监控插件的性能趋势</li>
                    <li>保持插件代码的更新和维护</li>
                    <li>备份重要的插件配置</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 模态框底部 -->
        <div class="border-t border-gray-200 px-6 py-4 flex justify-between items-center">
          <div class="flex items-center space-x-3">
            <button
              @click="copyLog"
              class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center"
            >
              <Icon name="copy" class="w-4 h-4 mr-2" />
              复制日志
            </button>
            <button
              @click="exportLog"
              class="px-4 py-2 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition-colors flex items-center"
            >
              <Icon name="download" class="w-4 h-4 mr-2" />
              导出详情
            </button>
          </div>
          <button
            @click="$emit('close')"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { useToast } from '~/composables/useToast'
import Icon from '~/components/Icon.vue'

const props = defineProps({
  log: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close'])

const toast = useToast()

// 方法
const getExecutionTimeClass = (time) => {
  if (time > 5000) return 'text-red-600 font-bold'
  if (time > 1000) return 'text-yellow-600 font-medium'
  return 'text-green-600'
}

const getMemoryUsage = () => {
  // 模拟内存使用情况
  return {
    used: '256 MB',
    available: '768 MB',
    percentage: 25
  }
}

const formatFullTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const copyLog = async () => {
  try {
    const logText = `日志详情
================
插件名称: ${props.log.plugin_name}
钩子名称: ${props.log.hook_name}
执行状态: ${props.log.success ? '成功' : '失败'}
执行时间: ${props.log.execution_time}ms
创建时间: ${formatFullTime(props.log.created_at)}
日志ID: ${props.log.id}
${props.log.error_message ? `错误信息: ${props.log.error_message}` : ''}

系统信息
================
用户代理: ${navigator.userAgent}
时区: ${Intl.DateTimeFormat().resolvedOptions().timeZone}
语言: ${navigator.language}
平台: ${navigator.platform}
内存使用: ${getMemoryUsage().used} / ${getMemoryUsage().percentage}%
`

    await navigator.clipboard.writeText(logText)
    toast.success('日志详情已复制到剪贴板')
  } catch (error) {
    toast.error('复制失败: ' + error.message)
  }
}

const exportLog = () => {
  try {
    const logData = {
      basic_info: {
        plugin_name: props.log.plugin_name,
        hook_name: props.log.hook_name,
        success: props.log.success,
        execution_time: props.log.execution_time,
        created_at: props.log.created_at,
        id: props.log.id
      },
      error_message: props.log.error_message,
      system_info: {
        user_agent: navigator.userAgent,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        language: navigator.language,
        platform: navigator.platform
      },
      memory_usage: getMemoryUsage(),
      export_time: new Date().toISOString()
    }

    const blob = new Blob([JSON.stringify(logData, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `log_detail_${props.log.id}_${new Date().toISOString().slice(0, 10)}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    toast.success('日志详情已导出')
  } catch (error) {
    toast.error('导出失败: ' + error.message)
  }
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

/* 代码块样式 */
pre {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 加载动画 */
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>