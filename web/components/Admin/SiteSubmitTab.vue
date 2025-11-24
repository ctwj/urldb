<template>
  <div class="tab-content-container">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">站点提交（待开发）</h3>
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
  </div>
</template>

<script setup lang="ts">
import { useMessage } from 'naive-ui'
import { ref } from 'vue'

// Props
interface Props {
  lastSubmitTime: {
    baidu: string
    google: string
    bing: string
    sogou: string
    shenma: string
    so360: string
  }
}

const props = withDefaults(defineProps<Props>(), {
  lastSubmitTime: () => ({
    baidu: '',
    google: '',
    bing: '',
    sogou: '',
    shenma: '',
    so360: ''
  })
})

// Emits
const emit = defineEmits<{
  'update:last-submit-time': [engine: string, time: string]
}>()

// 获取消息组件
const message = useMessage()

// 提交到百度
const submitToBaidu = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'baidu', time)
  message.success('已提交到百度')
}

// 提交到谷歌
const submitToGoogle = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'google', time)
  message.success('已提交到谷歌')
}

// 提交到必应
const submitToBing = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'bing', time)
  message.success('已提交到必应')
}

// 提交到搜狗
const submitToSogou = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'sogou', time)
  message.success('已提交到搜狗')
}

// 提交到神马搜索
const submitToShenma = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'shenma', time)
  message.success('已提交到神马搜索')
}

// 提交到360搜索
const submitTo360 = () => {
  const time = new Date().toLocaleString('zh-CN')
  emit('update:last-submit-time', 'so360', time)
  message.success('已提交到360搜索')
}

// 批量提交
const submitToAll = () => {
  const time = new Date().toLocaleString('zh-CN')
  const engines = ['baidu', 'google', 'bing', 'sogou', 'shenma', 'so360']

  engines.forEach(engine => {
    emit('update:last-submit-time', engine, time)
  })

  message.success('已批量提交到所有搜索引擎')
}
</script>

<style scoped>
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>