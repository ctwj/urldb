<template>
  <AdminPageLayout :is-sub-page="true">
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">账号扩容管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理账号扩容任务和状态</p>
      </div>
    </template>

    <!-- 通知提示区域 - 扩容说明Alert -->
    <template #notice-section>
      <n-alert type="info" :show-icon="false">
        <div class="flex items-center space-x-2">
          <i class="fas fa-info-circle text-blue-500"></i>
          <span class="text-sm">
          <strong>20T扩容说明：</strong>建议 <span @click="showImageModal = true" style="color:red" class="cursor-pointer">蜂小推</span> quark 账号扩容。<span @click="drawActive = true" style="color:blue" class="cursor-pointer">【什么推荐蜂小推】</span><br>
            1. 20T扩容 只支持新号，等到蜂小推首次 6T 奖励 到账后进行扩容<br>
            2. 账号需要处于关闭状态， 开启状态可能会被用于，自动转存等任务，存咋影响<br>
            3. <strong><n-text style='font-size:16px' type="error">扩容完成后，并不直接获得容量</n-text>，账号将存储大量热门资源，<n-text  style='font-size:16px' type="error">需要手动推广</n-text></strong><br>
            4. 注意 推广获得20T容量，删除所有资源， 热门资源比较敏感，不建议，长期推广，仅用于扩容
          </span>
        </div>
      </n-alert>
    </template>

    <!-- 内容区header - 账号列表标题 -->
    <template #content-header>
      <div class="flex items-center justify-between">
        <span class="text-lg font-semibold">支持扩容的账号列表</span>
        <div class="text-sm text-gray-500 dark:text-gray-400">
          共 {{ expansionAccounts.length }} 个账号
        </div>
      </div>
    </template>

    <!-- 内容区content - 账号列表详情 -->
    <template #content>
      <div class="flex-1 h-full overflow-hidden">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <n-spin size="large">
            <template #description>
              <span class="text-gray-500">加载中...</span>
            </template>
          </n-spin>
        </div>

        <!-- 空状态 -->
        <div v-else-if="expansionAccounts.length === 0" class="flex flex-col items-center justify-center py-12">
          <n-empty description="暂无可扩容的账号，请先添加有效的 quark 账号">
            <template #icon>
              <i class="fas fa-user-circle text-4xl text-gray-400"></i>
            </template>
          </n-empty>
        </div>

        <!-- 账号列表内容的虚拟滚动列表 -->
        <div v-else class="h-full">
          <n-virtual-list class="h-full" :items="expansionAccounts" :item-size="80">
            <template #default="{ item }">
              <div class="border-b border-gray-200 dark:border-gray-700 p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                <div class="flex items-center justify-between">
                  <!-- 左侧信息 -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center space-x-4">
                      <!-- 平台图标 -->
                      <span v-html="getPlatformIcon(item.service_type === 'quark' ? '夸克网盘' : '其他')" class="text-lg"></span>

                      <!-- 账号信息 -->
                      <div class="flex-1 min-w-0">
                        <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1">
                          {{ item.name }}
                        </h3>
                        <p class="text-xs text-gray-500 dark:text-gray-400">
                          {{ item.service_type === 'quark' ? '夸克网盘' : '其他平台' }}
                        </p>
                      </div>

                      <!-- 扩容状态 -->
                      <div class="flex items-center space-x-2">
                        <n-tag v-if="item.expanded" type="success" size="small">
                          已扩容
                        </n-tag>
                        <n-tag v-else type="warning" size="small">
                          可扩容
                        </n-tag>
                      </div>
                    </div>

                    <!-- 创建时间和容量 -->
                    <div class="mt-2">
                      <span class="text-xs text-gray-600 dark:text-gray-400">
                        创建时间: {{ formatDate(item.created_at) }}
                      </span>
                      <span v-if="item.total_space" class="text-xs text-gray-600 dark:text-gray-400 ml-4">
                        容量: {{ formatCapacity(item.used_space, item.total_space) }}
                      </span>
                    </div>
                  </div>

                  <!-- 右侧操作按钮 -->
                  <div class="flex items-center space-x-2 ml-4">
                    <n-button
                      size="small"
                      type="primary"
                      :disabled="item.expanded"
                      :loading="expandingAccountId === item.id"
                      @click="handleExpansion(item)"
                    >
                      <template #icon>
                        <i class="fas fa-expand"></i>
                      </template>
                      {{ item.expanded ? '已扩容' : '扩容' }}
                    </n-button>

                    <!-- 提取按钮，仅对已扩容账号显示 -->
                    <n-button
                      v-if="item.expanded"
                      size="small"
                      type="success"
                      @click="handleExtract(item)"
                    >
                      <template #icon>
                        <i class="fas fa-download"></i>
                      </template>
                      提取
                    </n-button>
                  </div>
                </div>
              </div>
            </template>
          </n-virtual-list>
        </div>
      </div>
    </template>
  </AdminPageLayout>
      <!-- 抽屉和模态框 - 移动到AdminPageLayout外部 -->
    <n-drawer v-model:show="drawActive" :width="502" closable placement="right">
      <n-drawer-content title="扩容说明">
        <div class="space-y-6 p-4">
          <div class="mb-4">
            <p class="text-gray-700 text-large dark:text-gray-300 leading-relaxed">
              扩容是网盘公司提供给推广用户的<n-text type="success">特权</n-text>！需要先注册推广平台并<n-text type="success">达标</n-text>即可获得权益。
            </p>
          </div>

          <n-collapse arrow-placement="right">
            <n-collapse-item title="达标要求" name="0">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                <i class="fas fa-list-check text-blue-500 mr-2"></i>
                达标要求（以蜂小推为例）
              </h3>
              <span>首次账号累计7天转存 > 10 或 拉新 > 5</span>
            </n-collapse-item>
            <n-collapse-item title="注意事项" name="1">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                <i class="fas fa-exclamation-triangle text-orange-500 mr-2"></i>
                注意事项
              </h3>
              <span>每个人的转存，只有当日第一次转存，且通过手机转存，才算有效转存。</span>
            </n-collapse-item>
            <n-collapse-item title="扩容原理" name="2">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                <i class="fas fa-question-circle text-purple-500 mr-2"></i>
                扩容原理
              </h3>
              <span>大量转存热播资源，这样才能尽可能快的达标。</span>
            </n-collapse-item>
            <n-collapse-item title="为什么推荐蜂小推" name="3">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                <i class="fas fa-thumbs-up text-green-500 mr-2"></i>
                为什么推荐蜂小推
              </h3>
              <p class="text-gray-600 dark:text-gray-400 leading-relaxed">
                登记后，第二天，即会发送 <strong class="text-blue-600">6T 空间</strong>，满足大量存储资源的前提条件。
              </p>
            </n-collapse-item>
            <n-collapse-item title="蜂小推怎么注册" name="3">
              <p class="text-gray-600 dark:text-gray-400">
                请扫描下方二维码进行注册。
              </p>
              <div class="mt-3 p-4 bg-gray-100 dark:bg-gray-800 rounded-lg text-center">
                <n-qr-code :value="qrCode" />
              </div>
            </n-collapse-item>
          </n-collapse>
        </div>
      </n-drawer-content>
    </n-drawer>

    <!-- 图片模态框 -->
    <n-modal v-model:show="showImageModal" title="蜂小推" size="huge">
      <div class="text-center">
        <img src="/assets/images/fxt.jpg" alt="蜂小推" class="max-w-full max-h-screen object-contain rounded-lg shadow-lg" />
      </div>
    </n-modal>

    <!-- 数据源选择弹窗 -->
    <n-modal v-model:show="showDataSourceDialog" title="确认扩容操作" size="small" style="width: 600px; max-width: 90vw;">
      <n-card title="数据源选择" size="small">
        <div class="space-y-4">
          <p class="text-gray-700">
            确定要对账号 "<strong>{{ pendingAccount?.name }}</strong>" 进行扩容操作吗？
          </p>

          <n-radio-group v-model:value="selectedDataSource">
            <n-space vertical>
              <n-radio value="internal">
                <span class="font-medium">系统内部数据源</span>
                <div class="text-sm text-gray-500 mt-1">使用系统内置的数据源进行扩容</div>
              </n-radio>
              <n-radio value="third-party">
                <span class="font-medium">第三方接口</span>
                <div class="text-sm text-gray-500 mt-1">使用第三方API获取数据源</div>
              </n-radio>
            </n-space>
          </n-radio-group>

          <n-input
            v-if="selectedDataSource === 'third-party'"
            v-model:value="thirdPartyUrl"
            placeholder="请输入第三方接口地址"
            clearable
          />

          <div class="flex justify-end space-x-2 mt-4">
            <n-button @click="showDataSourceDialog = false">取消</n-button>
            <n-button type="primary" @click="confirmDataSourceSelection">确定扩容</n-button>
          </div>
        </div>
      </n-card>
    </n-modal>

    <!-- 提取内容弹窗 -->
    <n-modal v-model:show="showExtractDialog" title="提取数据" size="large" style="width: 800px; max-width: 95vw;">
      <n-card title="提取扩容数据" size="small">
        <div class="space-y-4">
          <!-- 加载状态 -->
          <div v-if="loadingExtractData" class="flex items-center justify-center py-12">
            <n-spin size="large">
              <template #description>
                <span class="text-gray-500">获取数据中...</span>
              </template>
            </n-spin>
          </div>

          <!-- 数据加载完成后的内容 -->
          <div v-else class="flex flex-col gap-2">
            <!-- 显示模式选择 -->
            <div>
              <n-radio-group v-model:value="selectedDisplayMode">
                <n-space>
                  <n-radio value="links-only">
                    <span class="font-medium">仅链接</span>
                    <!-- <div class="text-sm text-gray-500 mt-1">只显示链接，一行一个</div> -->
                  </n-radio>
                  <n-radio value="title-link">
                    <span class="font-medium">标题|链接</span>
                    <!-- <div class="text-sm text-gray-500 mt-1">显示标题连接用|连接，一个一行</div> -->
                  </n-radio>
                  <n-radio value="title-newline-link">
                    <span class="font-medium">标题/n链接</span>
                    <!-- <div class="text-sm text-gray-500 mt-1">标题和链接，两行一个</div> -->
                  </n-radio>
                </n-space>
              </n-radio-group>
            </div>

            <!-- 文本显示区域 -->
            <div>
              <n-input
                v-model:value="extractedText"
                type="textarea"
                :rows="15"
                readonly
                placeholder="选择显示模式后，内容将显示在此处"
              />
            </div>

            <!-- 操作按钮 -->
            <div class="flex justify-end space-x-2">
              <n-button @click="copyToClipboard">复制</n-button>
              <n-button type="primary" @click="showExtractDialog = false">关闭</n-button>
            </div>
          </div>
        </div>
      </n-card>
    </n-modal>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: ['auth']
})

import { ref, onMounted, computed, h, watch } from 'vue'
import { useTaskApi } from '~/composables/useApi'
import { useNotification, useDialog } from 'naive-ui'

// 响应式数据
const expansionAccounts = ref([])
const loading = ref(true)
const expandingAccountId = ref(null)
const drawActive = ref(false) // 侧边栏激活
const qrCode = ref("https://app.fengtuiwl.com/#/pages/login/reg?p=22112503")
const showImageModal = ref(false) // 图片模态框
const showDataSourceDialog = ref(false) // 数据源选择弹窗
const selectedDataSource = ref('internal') // internal or third-party
const thirdPartyUrl = ref('https://so.252035.xyz/')
const pendingAccount = ref<any>(null) // 待处理的账号
const showExtractDialog = ref(false) // 提取内容弹窗
const selectedDisplayMode = ref('links-only') // 显示模式: links-only, title-link, title-newline-link
const extractedText = ref('') // 提取的文本内容
const currentExtractAccount = ref<any>(null) // 当前提取的账号
const loadingExtractData = ref(false) // 提取数据加载状态
const extractedResources = ref([]) // 保存获取到的资源数据

// API实例
const taskApi = useTaskApi()
const notification = useNotification()


// 获取支持扩容的账号列表
const fetchExpansionAccounts = async () => {
  loading.value = true
  try {
    const response = await taskApi.getExpansionAccounts()  as any
    expansionAccounts.value = response.accounts || []
  } catch (error) {
    console.error('获取扩容账号列表失败:', error)
    notification.error({
      title: '失败',
      content: '获取扩容账号列表失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}


// 处理扩容操作
const handleExpansion = async (account) => {
  pendingAccount.value = account
  showDataSourceDialog.value = true
}

// 确认数据源选择
const confirmDataSourceSelection = async () => {
  if (!pendingAccount.value) return

  showDataSourceDialog.value = false
  expandingAccountId.value = pendingAccount.value.id

  try {
    const dataSource = selectedDataSource.value === 'internal'
      ? { type: 'internal' }
      : { type: 'third-party', url: thirdPartyUrl.value }

    const response = await taskApi.createExpansionTask({
      pan_account_id: pendingAccount.value.id,
      description: `对 ${pendingAccount.value.name} 账号进行扩容操作`,
      dataSource
    })

    // 启动任务
    await taskApi.startTask(response.task_id)

    notification.success({
      title: '成功',
      content: '扩容任务已创建并启动',
      duration: 3000
    })

    navigateTo('/admin/tasks')
    // 刷新数据
    // await Promise.all([
    //   fetchExpansionAccounts(),
    //   fetchExpansionTasks()
    // ])
  } catch (error) {
    console.error('创建扩容任务失败:', error)
    notification.error({
      title: '失败',
      content: '创建扩容任务失败: ' + (error.message || '未知错误'),
      duration: 3000
    })
  } finally {
    expandingAccountId.value = null
    pendingAccount.value = null
  }
}



// 获取平台图标
const getPlatformIcon = (platformName) => {
  const defaultIcons = {
    '夸克网盘': '<i class="fas fa-cloud text-blue-600"></i>',
    '其他': '<i class="fas fa-cloud text-gray-500"></i>'
  }
  return defaultIcons[platformName] || defaultIcons['其他']
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化容量
const formatCapacity = (used, total) => {
  if (!total) return '-'
  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }
  return `${formatBytes(used || 0)} / ${formatBytes(total)}`
}

// 处理提取操作
const handleExtract = async (account) => {
  currentExtractAccount.value = account
  showExtractDialog.value = true
  loadingExtractData.value = true
  extractedText.value = ''

  try {
    // 获取账号的扩容任务输出数据
    const response = await taskApi.getExpansionOutput(account.id)
    const resources = response.output_data?.transferred_resources || []

    // 保存获取到的数据
    extractedResources.value = resources

    // 根据当前选择的模式格式化文本
    formatExtractedText(resources, selectedDisplayMode.value)

  } catch (error) {
    console.error('获取提取数据失败:', error)
    notification.error({
      title: '失败',
      content: '获取提取数据失败: ' + (error.data?.message || '未知错误'),
      duration: 3000
    })

    // 如果获取失败，使用模拟数据
    const mockData = [
      { title: "示例电影1", url: "https://example.com/1" },
      { title: "示例电影2", url: "https://example.com/2" },
      { title: "示例电影3", url: "https://example.com/3" }
    ]
    extractedResources.value = mockData
    formatExtractedText(mockData, selectedDisplayMode.value)
  } finally {
    loadingExtractData.value = false
  }
}

// 过滤标题文本，移除换行、| 和不可见字符
const cleanTitle = (title) => {
  if (!title) return ''
  return title
    .replace(/[\r\n\t]/g, ' ') // 移除换行符和制表符，替换为空格
    .replace(/[|]/g, ' ') // 移除|符号，替换为空格
    .replace(/[\x00-\x1F\x7F-\x9F]/g, '') // 移除不可见字符
    .replace(/\s+/g, ' ') // 多个空格合并为一个
    .trim() // 移除首尾空格
}

// 格式化提取的文本
const formatExtractedText = (resources, mode) => {
  if (!resources || resources.length === 0) {
    extractedText.value = '暂无数据'
    return
  }

  let text = ''

  switch (mode) {
    case 'links-only':
      // 仅链接，一行一个
      text = resources.map(item => item.url).join('\n')
      break
    case 'title-link':
      // 标题|链接，一个一行
      text = resources.map(item => `${cleanTitle(item.title)}|${item.url}`).join('\n')
      break
    case 'title-newline-link':
      // 标题/n链接，两行一个
      text = resources.map(item => `${cleanTitle(item.title)}\n${item.url}`).join('\n')
      break
    default:
      text = '请选择显示模式'
  }

  extractedText.value = text
}

// 复制到剪贴板
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(extractedText.value)
    notification.success({
      title: '成功',
      content: '内容已复制到剪贴板',
      duration: 3000
    })
  } catch (error) {
    console.error('复制失败:', error)
    notification.error({
      title: '失败',
      content: '复制失败，请手动选择文本复制',
      duration: 3000
    })
  }
}

// 监听显示模式变化
watch(selectedDisplayMode, (newMode) => {
  if (extractedResources.value && extractedResources.value.length > 0) {
    // 使用已保存的数据重新格式化文本
    formatExtractedText(extractedResources.value, newMode)
  }
})

// 页面加载
onMounted(async () => {
  await fetchExpansionAccounts()
})
</script>

<style scoped>
.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}
</style>
