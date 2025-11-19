<template>
  <div v-if="!systemConfig?.maintenance_mode" class="min-h-screen bg-gray-50 dark:bg-slate-900 text-gray-800 dark:text-slate-100">
    <!-- 主要内容区域 -->
    <div class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
      <!-- 头部导航 -->
      <div class="header-container bg-slate-800 dark:bg-slate-800 text-white dark:text-slate-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
        <h1 class="flex items-center justify-center gap-3">
          <img
            v-if="systemConfig?.site_logo"
            :src="getImageUrl(systemConfig.site_logo)"
            :alt="systemConfig?.site_title || 'Logo'"
            class="h-8 w-auto object-contain"
            @error="handleLogoError"
          />
          <img
            v-else
            src="/assets/images/logo.webp"
            alt="Logo"
            class="h-8 w-auto object-contain"
          />
          <NuxtLink to="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline text-2xl sm:text-3xl font-bold">
            {{ systemConfig?.site_title || '老九网盘资源数据库' }}
          </NuxtLink>
        </h1>

        <!-- 右侧导航按钮 -->
        <nav class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
          <NuxtLink to="/" class="hidden sm:flex">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-arrow-left text-xs"></i>
              <span class="ml-1">返回首页</span>
            </n-button>
          </NuxtLink>
        </nav>
      </div>

      <!-- 资源详情内容 -->
      <div v-if="resourcesData" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 主内容区域 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 资源主卡片 -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg overflow-hidden">
            <div class="flex flex-col sm:flex-row gap-6 p-6">
              <!-- 封面图片 -->
              <div class="flex-shrink-0 mx-auto sm:mx-0">
                <n-image
                  :src="getResourceImageUrl(mainResource)"
                  :alt="mainResource?.title || '资源图片'"
                  width="160"
                  class="rounded-xl object-cover border-2 border-gray-200 dark:border-slate-600 shadow-md hover:shadow-xl transition-all duration-300 w-40 h-56"
                  @error="handleResourceImageError"
                />
              </div>

              <!-- 资源信息 -->
              <div class="flex-1 space-y-4">
                <!-- 标题和操作按钮 -->
                <div class="flex flex-col gap-2">
                  <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                    <h1 class="text-xl sm:text-2xl font-bold text-gray-900 dark:text-gray-100 break-words leading-tight flex-1" v-html="mainResource?.title_highlight || mainResource?.title">
                    </h1>

                    <!-- 操作按钮组 -->
                    <div class="flex items-center gap-2 flex-shrink-0">
                      <button
                        class="px-3 py-1.5 text-xs font-medium rounded-lg border border-orange-200 dark:border-orange-400/30 bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400 hover:bg-orange-100 dark:hover:bg-orange-500/20 transition-colors flex items-center gap-1"
                        @click="showReportModal = true"
                        title="举报资源失效"
                      >
                        <i class="fas fa-exclamation-circle"></i>
                        <span class="hidden sm:inline">举报</span>
                      </button>
                      <button
                        class="px-3 py-1.5 text-xs font-medium rounded-lg border border-purple-200 dark:border-purple-400/30 bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400 hover:bg-purple-100 dark:hover:bg-purple-500/20 transition-colors flex items-center gap-1"
                        @click="showCopyrightModal = true"
                        title="版权申述"
                      >
                        <i class="fas fa-balance-scale"></i>
                        <span class="hidden sm:inline">申述</span>
                      </button>
                    </div>
                  </div>

                  <!-- 时间和浏览次数 -->
                  <div class="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
                    <span class="flex items-center gap-1">
                      <i class="fas fa-calendar-alt text-blue-500"></i>
                      {{ formatDate(mainResource?.updated_at) }}
                    </span>
                    <span class="flex items-center gap-1">
                      <i class="fas fa-eye text-green-500"></i>
                      {{ mainResource?.view_count || 0 }}次浏览
                    </span>
                  </div>
                </div>

                <!-- 标签 -->
                <div v-if="mainResource?.tags && mainResource.tags.length > 0" class="flex flex-wrap gap-2">
                  <template v-for="tag in mainResource.tags" :key="tag.id">
                    <span class="resource-tag inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 dark:bg-blue-500/10 text-blue-700 dark:text-blue-100 border border-blue-200 dark:border-blue-400/30">
                      <i class="fas fa-tag mr-1 text-blue-500 dark:text-blue-300 text-xs"></i>
                      {{ tag.name || '未知标签' }}
                    </span>
                  </template>
                </div>

                <!-- 描述 -->
                <div v-if="mainResource?.description" class="text-gray-600 dark:text-gray-300 text-sm leading-relaxed bg-gray-50 dark:bg-slate-700/50 rounded-lg p-4 border border-gray-100 dark:border-slate-600" v-html="mainResource.description_highlight || mainResource.description">
                </div>

                <!-- 基本信息 -->
                <div v-if="mainResource?.file_size || mainResource?.author" class="flex flex-wrap gap-4 text-xs text-gray-500 dark:text-gray-400">
                  <span v-if="mainResource?.file_size" class="flex items-center gap-1">
                    <i class="fas fa-file text-purple-500"></i>
                    {{ mainResource.file_size }}
                  </span>
                  <span v-if="mainResource?.author" class="flex items-center gap-1">
                    <i class="fas fa-user text-orange-500"></i>
                    {{ mainResource.author }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 网盘资源链接列表 -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg overflow-hidden">
            <div class="p-6">
              <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-3">
                  <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                    <i class="fas fa-cloud-download-alt text-blue-500"></i>
                    网盘资源 ({{ resourcesData?.resources?.length || 0 }})
                  </h3>

                  <!-- 检测状态总览 -->
                  <div v-if="resourcesData?.resources?.length > 0" class="flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium" :class="detectionStatus.text + ' bg-opacity-10 ' + detectionStatus.text.replace('text', 'bg')">
                    <i :class="detectionStatus.icon"></i>
                    <span>{{ detectionStatus.label }}</span>
                    <span v-if="detectionStatus.detectedCount > 0" class="ml-1 opacity-75">({{ detectionStatus.detectedCount }}已检测)</span>
                  </div>
                </div>

                <!-- 分享按钮 -->
                <ShareButtons
                  :title="mainResource?.title"
                  :description="mainResource?.description"
                  :url="getResourceUrl"
                  :tags="mainResource?.tags?.map(tag => tag.name)"
                />
              </div>

              <div class="space-y-3">
                <div
                  v-for="(resource, index) in resourcesData?.resources"
                  :key="resource.id"
                  class="flex items-center justify-between p-4 border rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors"
                  :class="resource.is_valid
                    ? 'border-gray-200 dark:border-slate-600'
                    : 'border-red-200 dark:border-red-400 bg-red-50/50 dark:bg-red-500/5'"
                >
                  <!-- 左侧：平台信息 -->
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-lg bg-gray-100 dark:bg-slate-700 flex items-center justify-center">
                      <span v-html="resource.pan?.icon" class="text-lg"></span>
                    </div>
                    <div>
                      <div class="font-medium text-gray-900 dark:text-gray-100">{{ resource.pan?.remark || '未知平台' }}</div>
                      <div class="flex items-center gap-2 mt-1">
                        <span
                          class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
                          :class="resource.is_valid
                            ? 'bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-100'
                            : 'bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-100'"
                        >
                          <span class="w-1.5 h-1.5 rounded-full mr-1" :class="resource.is_valid ? 'bg-green-500' : 'bg-red-500'"></span>
                          {{ resource.is_valid ? '有效' : '无效' }}
                        </span>
                        <span v-if="resource.save_url" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-100">
                          已转存
                        </span>
                      </div>
                    </div>
                  </div>

                  <!-- 右侧：检测状态和操作按钮 -->
                  <div class="flex items-center gap-2">
                    <!-- 检测状态标签（放在最前面） -->
                    <div v-if="detectionMethods[resource.id]" class="flex items-center gap-1">
                      <!-- 检测方法标识 -->
                      <span
                        class="px-1.5 py-0.5 rounded text-xs font-medium"
                        :class="getDetectionMethodClass(detectionMethods[resource.id])"
                        :title="getDetectionMethodTitle(detectionMethods[resource.id], resource)"
                      >
                        {{ getDetectionMethodLabel(detectionMethods[resource.id]) }}
                      </span>

                      <!-- 不支持检测的三角感叹号提示 -->
                      <span v-if="detectionMethods[resource.id] === 'unsupported'" class="text-amber-600 dark:text-amber-400" title="当前网盘暂不支持自动检测，建议您点击链接自行验证">
                        <i class="fas fa-exclamation-triangle"></i>
                      </span>
                    </div>

                    <!-- 检测中状态 -->
                    <div v-if="isDetecting && !detectionResults[resource.id]" class="flex items-center gap-2 px-3 py-2 text-sm text-blue-600 dark:text-blue-400">
                      <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                      <span>检测中... ({{ Object.keys(detectionResults).length }}/{{ resourcesData?.resources?.length }})</span>
                    </div>

                    <!-- 检测完成后的按钮 -->
                    <template v-else>
                      <!-- 重新检测按钮 -->
                      <button
                        class="px-3 py-1.5 text-xs font-medium rounded-lg border border-green-200 dark:border-green-400/30 bg-green-50 dark:bg-green-500/10 text-green-600 dark:text-green-400 hover:bg-green-100 dark:hover:bg-green-500/20 transition-colors flex items-center gap-1 disabled:opacity-50"
                        :disabled="isDetecting"
                        @click="smartDetectResourceValidity(true)"
                        title="重新检测链接有效性"
                      >
                        <i class="fas" :class="isDetecting ? 'fa-spinner fa-spin' : 'fa-sync-alt'"></i>
                        <span class="hidden sm:inline">{{ isDetecting ? '检测中' : '重测' }}</span>
                      </button>

                      <!-- 获取链接按钮 -->
                      <button
                        class="px-4 py-2 text-sm font-medium rounded-lg border border-blue-200 dark:border-blue-400/30 bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-500/20 transition-colors flex items-center gap-2 disabled:opacity-50"
                        :disabled="resource.forbidden || loadingStates[resource.id]"
                        @click="toggleLink(resource)"
                      >
                        <i class="fas" :class="loadingStates[resource.id] ? 'fa-spinner fa-spin' : 'fa-external-link-alt'"></i>
                        {{ resource.forbidden ? '受限' : (loadingStates[resource.id] ? '获取中' : '获取链接') }}
                      </button>

                      <!-- 复制转存链接按钮 -->
                      <button
                        v-if="resource.save_url && !resource.forbidden"
                        class="p-2 text-sm rounded-lg border border-gray-200 dark:border-slate-600 bg-white dark:bg-slate-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-slate-600 transition-colors"
                        @click="copyToClipboard(resource.save_url)"
                        title="复制转存链接"
                      >
                        <i class="fas fa-copy"></i>
                      </button>
                    </template>
                  </div>
                </div>
              </div>

              <!-- 违禁词提示 -->
              <div v-if="resourcesData?.resources?.some(r => r.forbidden)" class="mt-4 p-4 bg-yellow-50 dark:bg-yellow-500/10 rounded-lg border border-yellow-200 dark:border-yellow-400/30">
                <div class="flex items-start gap-2">
                  <i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400 mt-0.5"></i>
                  <div class="text-sm text-yellow-800 dark:text-yellow-200">
                    部分资源包含受限内容，无法正常访问
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧边栏 -->
        <div class="lg:col-span-1 space-y-6">
          <!-- 相关资源 -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg overflow-hidden">
            <div class="p-6">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
                <i class="fas fa-fire text-orange-500"></i>
                相关资源
              </h3>

              <!-- 相关资源列表 -->
              <div v-if="isRelatedResourcesLoading" class="space-y-3">
                <div v-for="i in 5" :key="i" class="animate-pulse">
                  <div class="flex items-center gap-3 p-2 rounded-lg">
                    <div class="w-5 h-5 bg-gray-200 dark:bg-gray-700 rounded-full flex-shrink-0"></div>
                    <div class="flex-1 space-y-2">
                      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded"></div>
                      <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
                    </div>
                  </div>
                </div>
              </div>

              <div v-else-if="displayRelatedResources.length > 0" class="space-y-3">
                <a
                  v-for="(resource, index) in displayRelatedResources"
                  :key="resource.id"
                  :href="`/r/${resource.key}`"
                  class="group block cursor-pointer"
                  @click.prevent="navigateToResource(resource.key)"
                >
                  <div class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                    <!-- 序号 -->
                    <div class="flex-shrink-0 flex items-center justify-center">
                      <div
                        class="w-5 h-5 rounded-full flex items-center justify-center text-xs font-medium"
                        :class="index < 3
                          ? 'bg-blue-500 text-white'
                          : 'bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300'"
                      >
                        {{ index + 1 }}
                      </div>
                    </div>

                    <!-- 资源信息 -->
                    <div class="flex-1 min-w-0">
                      <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                        {{ resource.title }}
                      </h4>
                      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1 line-clamp-1">
                        {{ resource.description }}
                      </p>
                    </div>
                  </div>
                </a>
              </div>

              <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
                <i class="fas fa-inbox text-3xl mb-2"></i>
                <p class="text-sm">暂无相关资源</p>
              </div>
            </div>
          </div>

          <!-- 热门资源 -->
          <div class="bg-white dark:bg-slate-800 rounded-xl shadow-lg overflow-hidden">
            <div class="p-6">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4 flex items-center gap-2">
                <i class="fas fa-trending-up text-red-500"></i>
                热门资源
              </h3>

              <!-- 热门资源列表 -->
              <div v-if="hotResourcesLoading" class="space-y-3">
                <div v-for="i in 10" :key="i" class="animate-pulse">
                  <div class="flex items-center gap-3 p-2 rounded-lg">
                    <div class="w-5 h-5 bg-gray-200 dark:bg-gray-700 rounded-full flex-shrink-0"></div>
                    <div class="flex-1 space-y-2">
                      <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded"></div>
                      <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
                    </div>
                  </div>
                </div>
              </div>

              <div v-else-if="hotResources.length > 0" class="space-y-3">
                <a
                  v-for="(resource, index) in hotResources"
                  :key="resource.id"
                  :href="`/r/${resource.key}`"
                  class="group block cursor-pointer"
                  @click.prevent="navigateToResource(resource.key)"
                >
                  <div class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700/50 transition-colors">
                    <!-- 排名标识 -->
                    <div class="flex-shrink-0 flex items-center justify-center">
                      <div
                        class="w-5 h-5 rounded-full flex items-center justify-center text-xs font-medium"
                        :class="index < 3
                          ? 'bg-red-500 text-white'
                          : 'bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-300'"
                      >
                        {{ index + 1 }}
                      </div>
                    </div>

                    <!-- 资源信息 -->
                    <div class="flex-1 min-w-0">
                      <div class="flex items-center justify-between">
                        <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100 line-clamp-1 group-hover:text-red-600 dark:group-hover:text-red-400 transition-colors flex-1">
                          {{ resource.title }}
                        </h4>
                        <!-- 排名皇冠 -->
                        <div class="flex items-center gap-1 flex-shrink-0 ml-2">
                          <i v-if="index === 0" class="fas fa-crown text-yellow-500 text-xs" title="第一名"></i>
                          <i v-else-if="index === 1" class="fas fa-crown text-gray-400 text-xs" title="第二名"></i>
                          <i v-else-if="index === 2" class="fas fa-crown text-orange-600 text-xs" title="第三名"></i>
                        </div>
                      </div>
                      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1 line-clamp-1">
                        {{ resource.description }}
                      </p>
                    </div>
                  </div>
                </a>
              </div>

              <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
                <i class="fas fa-fire-alt text-3xl mb-2"></i>
                <p class="text-sm">暂无热门资源</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 404 状态 -->
      <div v-else-if="!resourcesData" class="flex flex-col items-center justify-center py-20">
        <div class="text-center space-y-4">
          <i class="fas fa-search text-6xl text-gray-300 dark:text-gray-600"></i>
          <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">资源不存在</h2>
          <p class="text-gray-600 dark:text-gray-400">抱歉，您访问的资源不存在或已被删除</p>
          <NuxtLink
            to="/"
            class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-all duration-200"
          >
            <i class="fas fa-home"></i>
            返回首页
          </NuxtLink>
        </div>
      </div>
    </div>
    </div>

    <!-- 链接模态框 -->
    <QrCodeModal
      :visible="showLinkModal"
      :url="selectedResource?.url"
      :save_url="selectedResource?.save_url"
      :loading="selectedResource?.loading"
      :linkType="selectedResource?.linkType"
      :platform="selectedResource?.platform"
      :message="selectedResource?.message"
      :error="selectedResource?.error"
      :forbidden="selectedResource?.forbidden"
      :forbidden_words="selectedResource?.forbidden_words"
      @close="showLinkModal = false"
    />

    <!-- 举报模态框 -->
    <ReportModal
      :visible="showReportModal"
      :resource-key="resourceKey"
      @close="showReportModal = false"
      @submitted="handleReportSubmitted"
    />

    <!-- 版权申述模态框 -->
    <CopyrightModal
      :visible="showCopyrightModal"
      :resource-key="resourceKey"
      @close="showCopyrightModal = false"
      @submitted="handleCopyrightSubmitted"
    />

    <!-- 页脚 -->
    <AppFooter />

    <!-- 悬浮按钮组件 -->
    <FloatButtons />
  </div>

  <!-- 维护模式 -->
  <div v-if="systemConfig?.maintenance_mode" class="fixed inset-0 z-[1000000] flex items-center justify-center bg-gradient-to-br from-yellow-100/80 via-gray-900/90 to-yellow-200/80 backdrop-blur-sm">
    <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-2xl px-8 py-10 flex flex-col items-center max-w-xs w-full border border-yellow-200 dark:border-yellow-700">
      <i class="fas fa-tools text-yellow-500 text-5xl mb-6 animate-bounce-slow"></i>
      <h3 class="text-2xl font-extrabold text-yellow-600 dark:text-yellow-400 mb-2 tracking-wide drop-shadow">系统维护中</h3>
      <p class="text-base text-gray-600 dark:text-gray-300 mb-6 text-center leading-relaxed">
        我们正在进行系统升级和维护，预计很快恢复服务。<br>
        请稍后再试，感谢您的理解与支持！
      </p>
      <div class="flex space-x-1 mt-2">
        <span class="w-2 h-2 bg-yellow-400 rounded-full animate-blink"></span>
        <span class="w-2 h-2 bg-yellow-500 rounded-full animate-blink delay-200"></span>
        <span class="w-2 h-2 bg-yellow-600 rounded-full animate-blink delay-400"></span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 获取运行时配置
const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()

// 获取资源key参数
const resourceKey = computed(() => route.params.key as string)

// 导入API
import { useResourceApi, usePublicSystemConfigApi } from '~/composables/useApi'

const resourceApi = useResourceApi()
const publicSystemConfigApi = usePublicSystemConfigApi()

// 响应式数据
const showLinkModal = ref(false)
const showReportModal = ref(false)
const showCopyrightModal = ref(false)
const selectedResource = ref<any>(null)
const loadingStates = ref<Record<number, boolean>>({})
const isDetecting = ref(false)
const detectionResults = ref<Record<number, boolean>>({})
const detectionErrors = ref<Record<number, string>>({})
const detectionMethods = ref<Record<number, string>>({})
const detectionNotes = ref<Record<number, string>>({})
const relatedResources = ref<any[]>([])
const relatedResourcesLoading = ref(true)
const hotResources = ref<any[]>([])
const hotResourcesLoading = ref(true)

// 获取系统配置
const { data: systemConfigData } = await useAsyncData('systemConfig', () => publicSystemConfigApi.getPublicSystemConfig())
const systemConfig = computed(() => (systemConfigData.value as any)?.data || { site_title: '老九网盘资源数据库' })

// 获取资源数据
const { data: resourcesData, error: resourcesError } = await useAsyncData(
  `resources-${resourceKey.value}`,
  () => resourceApi.getResourcesByKey(resourceKey.value),
  {
    server: true,
    default: () => null
  }
)

// 获取相关资源（服务端渲染，用于SEO优化）
const { data: relatedResourcesData } = await useAsyncData(
  `related-resources-${resourceKey.value}`,
  () => {
    const params = {
      key: resourceKey.value,
      limit: 5
    }
    return resourceApi.getRelatedResources(params)
  },
  {
    server: true,
    default: () => ({ data: [] })
  }
)

// 主要资源信息
const mainResource = computed(() => {
  const resources = resourcesData.value?.resources
  return resources && resources.length > 0 ? resources[0] : null
})

// 生成完整的资源URL
const getResourceUrl = computed(() => {
  const key = mainResource.value?.key
  if (!key) return ''

  // 在客户端直接使用当前页面的origin
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/r/${key}`
  }

  // 在服务端渲染时返回相对路径（客户端会自动补全）
  return `/r/${key}`
})

// 服务端相关资源处理（去重）
const serverRelatedResources = computed(() => {
  const resources = Array.isArray(relatedResourcesData.value?.data) ? relatedResourcesData.value.data : []

  // 根据key去重，避免显示重复资源
  const uniqueResources = resources.filter((resource, index, self) =>
    index === self.findIndex((r) => r.key === resource.key)
  )

  return uniqueResources.slice(0, 5) // 最多显示5个相关资源
})

// 合并服务端和客户端相关资源，优先显示服务端数据，支持SEO
const displayRelatedResources = computed(() => {
  // 如果有客户端数据（可能是更新的数据），使用客户端数据
  if (relatedResources.value.length > 0) {
    return relatedResources.value
  }

  // 否则使用服务端数据，确保SEO友好
  return serverRelatedResources.value
})

// 相关资源加载状态
const isRelatedResourcesLoading = computed(() => {
  // 如果有服务端数据，不显示加载状态
  if (serverRelatedResources.value.length > 0) {
    return false
  }

  // 否则显示客户端加载状态
  return relatedResourcesLoading.value
})

// 检测状态
const detectionStatus = computed(() => {
  if (isDetecting.value) {
    return {
      icon: 'fas fa-spinner fa-spin text-blue-600',
      text: 'text-blue-600',
      label: '检测中',
      detectedCount: 0
    }
  }

  const resources = resourcesData.value?.resources
  if (!resources || resources.length === 0) {
    return {
      icon: 'fas fa-question-circle text-gray-400',
      text: 'text-gray-400',
      label: '未知状态',
      detectedCount: 0
    }
  }

  // 统计有效和已检测的资源（只统计真正有检测结果的）
  const detectedResources = resources.filter(r => detectionResults.value[r.id] !== undefined)
  const validCount = detectedResources.filter(r => detectionResults.value[r.id] === true).length

  const totalCount = detectedResources.length
  const undetectedCount = resources.length - detectedResources.length

  // 如果没有检测任何资源，显示未检测状态
  if (detectedResources.length === 0) {
    return {
      icon: 'fas fa-question-circle text-gray-400',
      text: 'text-gray-400',
      label: '未检测',
      detectedCount: 0
    }
  }

  // 基于已检测的资源显示状态
  if (validCount === totalCount) {
    return {
      icon: 'fas fa-check-circle text-green-600',
      text: 'text-green-600',
      label: '全部有效',
      detectedCount: detectedResources.length
    }
  } else if (validCount === 0) {
    return {
      icon: 'fas fa-times-circle text-red-600',
      text: 'text-red-600',
      label: '全部无效',
      detectedCount: detectedResources.length
    }
  } else {
    return {
      icon: 'fas fa-exclamation-triangle text-orange-600',
      text: 'text-orange-600',
      label: `${validCount}/${totalCount} 有效`,
      detectedCount: detectedResources.length
    }
  }
})

// 图片URL处理
const { getImageUrl } = useImageUrl()

// Logo错误处理
const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/assets/images/logo.webp'
}

// 获取资源图片URL
const getResourceImageUrl = (resource: any) => {
  if (!resource) return '/assets/images/cover1.webp'

  if (resource.image_url) {
    return getImageUrl(resource.image_url)
  }

  if (resource.cover) {
    return getImageUrl(resource.cover)
  }

  const randomNum = Math.floor(Math.random() * 8) + 1
  return `/assets/images/cover${randomNum}.webp`
}

// 处理资源图片加载错误
const handleResourceImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  const randomNum = Math.floor(Math.random() * 8) + 1
  img.src = `/assets/images/cover${randomNum}.webp`
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 检测资源有效性
const detectResourceValidity = async () => {
  if (!resourcesData.value?.resources) return

  isDetecting.value = true
  detectionResults.value = {} // 重置检测结果
  detectionErrors.value = {} // 重置错误信息
  detectionMethods.value = {} // 重置检测方法
  detectionNotes.value = {} // 重置检测提示

  try {
    // 提取所有资源ID
    const resourceIds = resourcesData.value.resources.map(r => r.id)

    // 批量检测所有资源
    const response = await resourceApi.batchCheckResourceValidity(resourceIds) as any

    // 处理检测结果
    if (response && response.results) {
      response.results.forEach((result: any) => {
        // 只有真正进行了检测的资源才设置检测结果
        if (result.detection_method !== 'unsupported') {
          detectionResults.value[result.resource_id] = result.is_valid
        }
        detectionMethods.value[result.resource_id] = result.detection_method || 'unknown'

        // 保存错误信息
        if (result.error) {
          detectionErrors.value[result.resource_id] = result.error
        }

        // 保存检测提示
        if (result.note) {
          detectionNotes.value[result.resource_id] = result.note
        }

        // 显示缓存状态
        if (result.cached) {
          console.log(`资源 ${result.resource_id} 使用缓存检测结果`)
        }
      })
    }
  } finally {
    isDetecting.value = false
  }
}

// 智能检测：避免频繁重复检测
const smartDetectResourceValidity = async (force = false) => {
  if (!resourcesData.value?.resources) return

  // 如果正在检测，不重复执行
  if (isDetecting.value) {
    console.log('检测正在进行中，跳过重复请求')
    return
  }

  // 如果不是强制检测且已有检测结果，跳过
  if (!force && Object.keys(detectionResults.value).length > 0) {
    const lastDetectionTime = lastDetectionTimestamp.value
    const now = Date.now()
    const timeSinceLastDetection = now - lastDetectionTime

    // 5分钟内不重复检测
    if (timeSinceLastDetection < 5 * 60 * 1000) {
      console.log('距离上次检测时间不足5分钟，跳过重复检测')
      return
    }
  }

  await detectResourceValidity()
  lastDetectionTimestamp.value = Date.now()
}

// 添加最后检测时间戳
const lastDetectionTimestamp = ref(0)

// 切换链接显示
const toggleLink = async (resource: any) => {
  if (resource.forbidden) {
    selectedResource.value = {
      ...resource,
      forbidden: true,
      error: '该资源包含受限内容，无法访问',
      forbidden_words: resource.forbidden_words || []
    }
    showLinkModal.value = true
    return
  }

  loadingStates.value[resource.id] = true
  selectedResource.value = { ...resource, loading: true }
  showLinkModal.value = true

  try {
    const linkData = await resourceApi.getResourceLink(resource.id) as any

    selectedResource.value = {
      ...resource,
      url: linkData.url,
      save_url: linkData.type === 'transferred' ? linkData.url : resource.save_url,
      loading: false,
      linkType: linkData.type,
      platform: linkData.platform,
      message: linkData.message
    }
  } catch (error: any) {
    console.error('获取资源链接失败:', error)
    selectedResource.value = {
      ...resource,
      loading: false,
      error: '检测有效性失败，请自行验证'
    }
  } finally {
    loadingStates.value[resource.id] = false
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    // 显示复制成功提示
    if (process.client) {
      const notification = useNotification()
      if (notification) {
        notification.success({
          content: '已复制到剪贴板',
          duration: 2000
        })
      }
    }
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 处理举报提交
const handleReportSubmitted = () => {
  showReportModal.value = false
  if (process.client) {
    const notification = useNotification()
    if (notification) {
      notification.success({
        content: '举报已提交，感谢您的反馈',
        duration: 3000
      })
    }
  }
}

// 处理版权申述提交
const handleCopyrightSubmitted = () => {
  showCopyrightModal.value = false
  if (process.client) {
    const notification = useNotification()
    if (notification) {
      notification.success({
        content: '版权申述已提交，我们会尽快处理',
        duration: 3000
      })
    }
  }
}

// 获取相关资源（客户端更新，用于交互优化）
const fetchRelatedResources = async () => {
  if (!mainResource.value) return

  // 如果已经有服务端数据，跳过客户端获取（SEO优化）
  if (serverRelatedResources.value.length > 0) {
    console.log('使用服务端相关资源数据，跳过客户端获取')
    return
  }

  try {
    relatedResourcesLoading.value = true

    // 使用新的相关资源API，基于资源key查找相关资源
    const params: any = {
      key: resourceKey.value, // 使用当前资源的key
      limit: 5,
    }

    const response = await resourceApi.getRelatedResources(params) as any

    // 处理响应数据
    const resources = Array.isArray(response?.data) ? response.data : []

    // 根据key去重，避免显示重复资源
    const uniqueResources = resources.filter((resource, index, self) =>
      index === self.findIndex((r) => r.key === resource.key)
    )

    relatedResources.value = uniqueResources.slice(0, 5) // 最多显示5个相关资源

    console.log('获取相关资源成功:', {
      source: response?.source,
      count: relatedResources.value.length,
      params: params
    })

  } catch (error) {
    console.error('获取相关资源失败:', error)
    relatedResources.value = []
  } finally {
    relatedResourcesLoading.value = false
  }
}

// 获取热门资源
const fetchHotResources = async () => {
  try {
    hotResourcesLoading.value = true

    // 使用专门的热门资源API，保持10个热门资源
    const params = {
      limit: 10
    }

    const response = await resourceApi.getHotResources(params) as any

    // 处理响应数据
    const resources = Array.isArray(response?.data) ? response.data : []

    hotResources.value = resources.slice(0, 10)

    console.log('获取热门资源成功:', {
      count: hotResources.value.length,
      params: params
    })

  } catch (error) {
    console.error('获取热门资源失败:', error)
    hotResources.value = []
  } finally {
    hotResourcesLoading.value = false
  }
}

// 导航到资源详情页
const navigateToResource = (key: string) => {
  navigateTo(`/r/${key}`)
}

// 检测方法相关函数
const getDetectionMethodLabel = (method: string) => {
  const labels: Record<string, string> = {
    'quark_deep': '深度检测',
    'unsupported': '未检测',
    'unknown': '未知方法',
    'error': '检测错误'
  }
  return labels[method] || '未知'
}

const getDetectionMethodClass = (method: string) => {
  const classes: Record<string, string> = {
    'quark_deep': 'bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-100',
    'unsupported': 'bg-gray-100 text-gray-600 dark:bg-gray-500/20 dark:text-gray-300',
    'unknown': 'bg-gray-100 text-gray-700 dark:bg-gray-500/20 dark:text-gray-100',
    'error': 'bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-100'
  }
  return classes[method] || 'bg-gray-100 text-gray-700 dark:bg-gray-500/20 dark:text-gray-100'
}

const getDetectionMethodTitle = (method: string, resource: any) => {
  const titles: Record<string, string> = {
    'quark_deep': '使用深度检测（通过实际网盘操作验证）',
    'unsupported': `${resource.pan?.remark || '未知网盘'} 暂不支持深度检测`,
    'unknown': '检测方法未知',
    'error': '检测过程中发生错误'
  }
  return titles[method] || '未知检测方法'
}

// 页面加载完成后
onMounted(() => {
  // 获取相关资源
  nextTick(() => {
    fetchRelatedResources()
  })

  // 获取热门资源
  nextTick(() => {
    fetchHotResources()
  })
})

// 设置页面SEO
const { initSystemConfig, setPageSeo } = useGlobalSeo()
const { generateOgImageUrl } = useSeo()

// 动态生成页面SEO信息
const pageSeo = computed(() => {
  const resource = mainResource.value
  if (!resource) return null

  const title = resource.title || '资源详情'
  const description = resource.description
    ? resource.description.substring(0, 160)
    : `${resource.title} - 多网盘资源下载，支持百度网盘、阿里云盘、夸克网盘等多个平台`

  const keywords = [
    ...(resource.tags?.map((tag: any) => tag.name) || []),
    resource.title,
    '网盘资源',
    '资源下载',
    ...(resource.pan?.name ? [resource.pan.name] : [])
  ].join(', ')

  return { title, description, keywords }
})

// 更新页面SEO的函数
const updatePageSeo = () => {
  if (!pageSeo.value) return

  const { title, description, keywords } = pageSeo.value

  // 设置基本SEO
  setPageSeo(title, {
    description,
    keywords
  })

  // 设置更详细的HTML元数据
  let canonicalUrl
  if (typeof window !== 'undefined') {
    canonicalUrl = `${window.location.origin}/r/${resourceKey.value}`
  } else {
    // 在服务端渲染时使用相对路径
    canonicalUrl = `/r/${resourceKey.value}`
  }

  // 生成动态OG图片URL（使用新的key参数格式）
  const ogImageUrl = generateOgImageUrl(resourceKey.value, '', 'blue')

  useHead({
    htmlAttrs: {
      lang: 'zh-CN'
    },
    link: [
      {
        rel: 'canonical',
        href: canonicalUrl
      }
    ],
    meta: [
      // Open Graph
      {
        property: 'og:title',
        content: title
      },
      {
        property: 'og:description',
        content: description
      },
      {
        property: 'og:url',
        content: canonicalUrl
      },
      {
        property: 'og:type',
        content: 'website'
      },
      {
        property: 'og:image',
        content: ogImageUrl
      },
      {
        property: 'og:site_name',
        content: systemConfig.value?.site_title || '老九网盘资源数据库'
      },
      // Twitter Card
      {
        name: 'twitter:card',
        content: 'summary_large_image'
      },
      {
        name: 'twitter:title',
        content: title
      },
      {
        name: 'twitter:description',
        content: description
      },
      {
        name: 'twitter:image',
        content: ogImageUrl
      },
      // 其他元数据
      {
        name: 'robots',
        content: 'index, follow'
      },
      {
        name: 'author',
        content: mainResource.value?.author || systemConfig.value?.site_title || '老九网盘资源数据库'
      },
      {
        name: 'revisit-after',
        content: '1 days'
      }
    ],
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          "@context": "https://schema.org",
          "@type": "WebPage",
          "name": title,
          "description": description,
          "url": canonicalUrl,
          "image": ogImageUrl,
          "mainEntity": {
            "@type": "SoftwareApplication" || "DigitalDocument",
            "name": title,
            "description": description,
            "author": {
              "@type": "Person" || "Organization",
              "name": mainResource.value?.author || '未知'
            },
            "dateModified": mainResource.value?.updated_at,
            "keywords": keywords,
            "image": ogImageUrl,
            "offers": {
              "@type": "Offer",
              "price": "0",
              "priceCurrency": "CNY"
            }
          },
          "publisher": {
            "@type": "Organization",
            "name": systemConfig.value?.site_title || '老九网盘资源数据库'
          },
          "relatedContent": displayRelatedResources.value.map((resource, index) => ({
            "@type": "SoftwareApplication" || "DigitalDocument",
            "position": index + 1,
            "name": resource.title,
            "description": resource.description?.substring(0, 160) || '',
            "url": `${baseUrl}/r/${resource.key}`,
            "dateModified": resource.updated_at,
            "keywords": resource.tags?.map((tag: any) => tag.name).join(', ') || '',
            "image": generateOgImageUrl(resource.key, '', 'green')
          }))
        })
      }
    ]
  })
}

onBeforeMount(async () => {
  await initSystemConfig()
  updatePageSeo()
})

// 监听资源数据变化，更新SEO
watch([mainResource, systemConfig], () => {
  nextTick(() => {
    updatePageSeo()
  })
}, { deep: true })

// 错误处理
watch(resourcesError, (error) => {
  if (error) {
    console.error('获取资源数据失败:', error)
  }
})
</script>

<style scoped>
.header-container{
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom,
      rgba(0,0,0,0.1) 0%,
      rgba(0,0,0,0.25) 100%
  );
}

.resource-tag {
  transition: all 0.2s ease;
}

.resource-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

@keyframes bounce-slow {
  0%, 100% { transform: translateY(0);}
  50% { transform: translateY(-12px);}
}
.animate-bounce-slow {
  animation: bounce-slow 1.6s infinite;
}
@keyframes blink {
  0%, 80%, 100% { opacity: 0.2;}
  40% { opacity: 1;}
}
.animate-blink {
  animation: blink 1.4s infinite both;
}
.animate-blink.delay-200 { animation-delay: 0.2s; }
.animate-blink.delay-400 { animation-delay: 0.4s; }
</style>