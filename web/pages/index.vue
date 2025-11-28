  <template>
  <div v-if="!systemConfig?.maintenance_mode" class="min-h-screen bg-gray-50 dark:bg-slate-900 text-gray-800 dark:text-slate-100 flex flex-col">
    <!-- 全局加载状态 -->
    <div v-if="pageLoading" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-xl">
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <div class="text-center">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">正在加载...</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">请稍候，正在初始化系统</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <main class="flex-1 p-3 sm:p-5">
      <div class="max-w-7xl mx-auto">
      <!-- 头部 -->
      <header class="header-container bg-slate-800 dark:bg-slate-800 text-white dark:text-slate-100 rounded-lg shadow-lg p-4 sm:p-8 mb-4 sm:mb-8 text-center relative">
        <h1 class="text-2xl sm:text-3xl font-bold mb-4 flex items-center justify-center gap-3">
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
          <a href="/" class="text-white hover:text-gray-200 dark:hover:text-gray-300 no-underline">
            {{ systemConfig?.site_title || '老九网盘资源数据库' }}
          </a>
        </h1>
        
        <nav aria-label="主导航" class="mt-4 flex flex-col sm:flex-row justify-center gap-2 sm:gap-2 right-4 top-0 absolute">
          <NuxtLink to="/hot-dramas" class="hidden sm:flex" title="浏览热门影视剧资源">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-film text-xs" aria-hidden="true"></i> 热播剧
            </n-button>
          </NuxtLink>
          <NuxtLink to="/monitor" class="hidden sm:flex" title="查看系统运行状态和统计数据">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-chart-line text-xs" aria-hidden="true"></i> 系统监控
            </n-button>
          </NuxtLink>
          <NuxtLink to="/api-docs" class="hidden sm:flex" title="查看API接口文档">
            <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
              <i class="fas fa-book text-xs" aria-hidden="true"></i> API文档
            </n-button>
          </NuxtLink>
          <ClientOnly>
            <NuxtLink v-if="authInitialized && !userStore.isAuthenticated" to="/login" class="sm:flex" title="登录账号">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-sign-in-alt text-xs" aria-hidden="true"></i> 登录
              </n-button>
            </NuxtLink>
            <NuxtLink v-if="authInitialized && userStore.isAuthenticated && userStore.user?.role === 'admin'" to="/admin" class="hidden sm:flex" title="进入管理后台">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-user-shield text-xs" aria-hidden="true"></i> 管理后台
              </n-button>
            </NuxtLink>
            <NuxtLink v-if="authInitialized && userStore.isAuthenticated && userStore.user?.role !== 'admin'" to="/user" class="hidden sm:flex" title="进入用户中心">
              <n-button size="tiny" type="tertiary" round ghost class="!px-2 !py-1 !text-xs !text-white dark:!text-white !border-white/30 hover:!border-white">
                <i class="fas fa-user text-xs" aria-hidden="true"></i> 用户中心
              </n-button>
            </NuxtLink>
          </ClientOnly>
        </nav>
      </header>

      <!-- 公告信息 -->
      <aside aria-label="公告信息" class="w-full max-w-3xl mx-auto mb-2 px-2 sm:px-0">
        <Announcement />
      </aside>

      <!-- 搜索区域 -->
      <section aria-label="搜索功能" class="w-full max-w-3xl mx-auto mb-4 sm:mb-8 px-2 sm:px-0">
        <h2 class="sr-only">搜索网盘资源</h2>
        <ClientOnly>
          <div class="relative">
            <n-input 
              round 
              placeholder="搜索资源名称、关键词..." 
              v-model:value="searchQuery" 
              @blur="handleSearch" 
              @keyup.enter="handleSearch" 
              clearable
              aria-label="搜索资源"
            >
                <template #prefix>
                <i class="fas fa-search text-gray-400" aria-hidden="true"></i>
                </template>
              </n-input>
          </div>
        </ClientOnly>

        <!-- 平台类型筛选 -->
        <nav aria-label="平台筛选" class="mt-3 flex flex-wrap gap-2" id="platformFilters">
          <a 
            :href="`/?search=${$route.query.search || ''}&platform=`"
            class="px-2 py-1 text-xs rounded-full bg-slate-800 dark:bg-gray-700 text-white dark:text-gray-100 hover:bg-slate-700 dark:hover:bg-gray-600 transition-colors"
            :class="{ 'active-filter': !selectedPlatform }"
            :aria-current="!selectedPlatform ? 'page' : undefined"
            title="显示所有平台的资源"
          >
            全部
          </a>
          <a 
            v-for="platform in platforms" 
            :key="platform.id"
            :href="`/?search=${$route.query.search || ''}&platform=${platform.id}`"
            class="px-2 py-1 text-xs rounded-full bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
            :class="{ 'active-filter': selectedPlatform === platform.id }"
            :aria-current="selectedPlatform === platform.id ? 'page' : undefined"
            :title="`筛选${platform.name}平台的资源`"
          >
            <span v-html="platform.icon" aria-hidden="true"></span> {{ platform.name }}
          </a>
        </nav>
        
        <!-- 统计信息 -->
        <div class="flex justify-between mt-3 text-sm text-gray-600 dark:text-gray-300 px-2" role="status" aria-live="polite">
          <div class="flex items-center">
            <i class="fas fa-calendar-day text-pink-600 mr-1" aria-hidden="true"></i>
            <span class="sr-only">今日新增资源数量：</span>
            今日资源: 
            <span v-if="statsLoading" class="font-medium text-pink-600 ml-1">
              <i class="fas fa-spinner fa-spin text-xs" aria-hidden="true"></i>
              <span class="sr-only">加载中</span>
            </span>
            <span v-else class="font-medium text-pink-600 ml-1 count-up" :data-target="safeStats?.today_resources || 0" :aria-label="`今日新增${safeStats?.today_resources || 0}个资源`">0</span>
          </div>
          <div class="flex items-center">
            <i class="fas fa-database text-blue-600 mr-1" aria-hidden="true"></i>
            <span class="sr-only">总资源数量：</span>
            总资源数: 
            <span v-if="statsLoading" class="font-medium text-blue-600 ml-1">
              <i class="fas fa-spinner fa-spin text-xs" aria-hidden="true"></i>
              <span class="sr-only">加载中</span>
            </span>
            <span v-else class="font-medium text-blue-600 ml-1 count-up" :data-target="safeStats?.total_resources || 0" :aria-label="`共${safeStats?.total_resources || 0}个资源`">0</span>
          </div>
        </div>
      </section>

      <!-- 资源列表 -->
      <section aria-label="资源列表" class="overflow-x-auto bg-white dark:bg-slate-800 rounded-lg shadow-lg shadow-slate-900/10 dark:shadow-slate-900/50">
        <h2 class="sr-only">{{ searchQuery ? `"${searchQuery}" 的搜索结果` : '最新网盘资源列表' }}</h2>
        <table class="w-full min-w-full" role="table" aria-label="网盘资源列表">
          <caption class="sr-only">
            {{ searchQuery ? `搜索"${searchQuery}"找到${safeResources.length}个资源` : `最新网盘资源，共${safeResources.length}个` }}
          </caption>
          <thead>
            <tr class="bg-slate-800 dark:bg-slate-700 text-white dark:text-slate-100">
              <th scope="col" class="text-left text-xs sm:text-sm w-20 pl-2 sm:pl-3">
                <div class="flex items-center">
                  <i class="fas fa-image mr-1 text-gray-300 dark:text-slate-300" aria-hidden="true"></i> 
                  <span>封面</span>
                </div>
              </th>
              <th scope="col" class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm">
                <div class="flex items-center">
                  <i class="fas fa-cloud mr-1 text-gray-300 dark:text-slate-300" aria-hidden="true"></i> 
                  <span>文件名</span>
                </div>
              </th>
              <th scope="col" class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-24">链接</th>
              <th scope="col" class="px-2 sm:px-6 py-3 sm:py-4 text-left text-xs sm:text-sm hidden sm:table-cell w-32">更新时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-slate-700">
            <tr v-if="safeLoading" class="text-center py-8">
              <td colspan="1" class="text-gray-500 dark:text-gray-400 sm:hidden">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
              <td colspan="4" class="text-gray-500 dark:text-gray-400 hidden sm:table-cell">
                <i class="fas fa-spinner fa-spin mr-2"></i>加载中...
              </td>
            </tr>
            <tr v-else-if="safeResources.length === 0" class="text-center py-12">
              <td colspan="4" class="text-gray-500 dark:text-slate-500">
                <div class="flex flex-col items-center justify-center space-y-4">
                  <img 
                    src="/assets/svg/empty.svg" 
                    alt="暂无数据" 
                    class="!w-64 !h-64 sm:w-64 sm:h-64 opacity-60 dark:opacity-40"
                  />
                  <div class="text-center">
                    <p class="text-lg font-medium text-gray-600 dark:text-gray-400 mb-2">
                      {{ searchQuery ? '没有找到相关资源' : '暂无资源数据' }}
                    </p>
                    <p class="text-sm text-gray-500 dark:text-gray-500">
                      {{ searchQuery ? '请尝试其他关键词或清除搜索条件' : '资源正在整理中，请稍后再来查看' }}
                    </p>
                  </div>
                </div>
              </td>
            </tr>
            <tr
              v-for="(resource, index) in safeResources"
              :key="resource.id"
              :class="isUpdatedToday(resource.updated_at) ? 'hover:bg-pink-50 dark:hover:bg-pink-500/10 bg-pink-50/30 dark:bg-pink-500/5 cursor-pointer' : 'hover:bg-gray-50 dark:hover:bg-slate-700/50 cursor-pointer'"
              :data-index="index"
              @click="navigateToDetail(resource.key)"
            >
              <td class="text-xs sm:text-sm w-20 pl-2 sm:pl-3">
                <div class="flex justify-center">
                  <ClientOnly>
                    <n-image
                      :src="getResourceImageUrl(resource)"
                      :alt="`${resource.title} - ${getPlatformName(resource.pan_id)} 网盘资源封面图`"
                      :title="resource.title"
                      width="80"
                      class="rounded object-cover border border-gray-200 dark:border-slate-600 h-auto"
                      lazy
                      @error="handleResourceImageError"
                    />
                    <template #placeholder>
                      <div class="w-[80px] h-[80px] rounded bg-gray-200 dark:bg-slate-700 flex items-center justify-center">
                        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                      </div>
                    </template>
                    <template #fallback>
                      <div class="w-[80px] h-[80px] rounded bg-gray-200 dark:bg-slate-700 animate-pulse"></div>
                    </template>
                  </ClientOnly>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm">
                <div class="flex items-start">
                  <span class="mr-2 flex-shrink-0" v-html="getPlatformIcon(resource.pan_id || 0)" aria-hidden="true"></span>
                  <div class="flex-1 min-w-0">
                    <h3 class="break-words font-medium text-base" v-html="resource.title_highlight || resource.title"></h3>
                    <!-- 显示标签 -->
                    <div v-if="resource.tags && resource.tags.length > 0" class="mt-1 flex flex-wrap gap-1" role="list" aria-label="资源标签">
                      <template v-for="(tag, index) in resource.tags" :key="tag.id">
                        <span
                          class="resource-tag inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-500/20 text-blue-800 dark:text-blue-100 border dark:border-blue-400/30"
                          :title="`标签: ${tag.name}`"
                          role="listitem"
                        >
                          <i class="fas fa-tag mr-1 dark:text-blue-200" aria-hidden="true"></i>
                          <span>{{ tag.name || '未知标签' }}</span>
                        </span>
                      </template>
                    </div>
                    <!-- 显示描述 -->
                    <div v-if="resource.description_highlight || resource.description" class="text-xs text-gray-600 dark:text-slate-400 mt-1 break-words line-clamp-2" v-html="resource.description_highlight || resource.description">
                    </div>
                  </div>
                </div>
                <div class="sm:hidden mt-2 space-y-2">
                  <!-- 移动端时间和链接按钮一行显示 -->
                  <div class="flex items-center gap-2">
                    <div class="flex-1 min-w-0">
                      <div class="text-xs text-gray-500 dark:text-slate-400 truncate" :title="resource.updated_at">
                        <span v-html="formatRelativeTime(resource.updated_at)"></span>
                      </div>
                    </div>
                    <div class="flex-1 flex justify-end">
                      <NuxtLink
                        :to="`/r/${resource.key}`"
                        class="mobile-link-btn flex items-center gap-1 text-xs no-underline"
                        :title="`查看 ${resource.title} 的详细信息和下载链接`"
                        :aria-label="`查看 ${resource.title} 详情`"
                        @click.stop
                      >
                        <i class="fas fa-eye" aria-hidden="true"></i> 查看详情
                      </NuxtLink>
                    </div>
                  </div>
                </div>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm hidden sm:table-cell w-32">
                <NuxtLink
                  :to="`/r/${resource.key}`"
                  class="text-blue-600 hover:text-blue-800 flex items-center gap-1 show-link-btn"
                  :title="`查看 ${resource.title} 的详细信息和下载链接`"
                  :aria-label="`查看 ${resource.title} 详情`"
                  @click.stop
                >
                  <i class="fas fa-eye" aria-hidden="true"></i> 查看详情
                </NuxtLink>
              </td>
              <td class="px-2 sm:px-6 py-2 sm:py-4 text-xs sm:text-sm text-gray-500 hidden sm:table-cell w-32" :title="resource.updated_at">
                <span v-html="formatRelativeTime(resource.updated_at)"></span>
              </td>
            </tr>
          </tbody>
        </table>
      </section>

    </div>

    </main>

    <!-- 二维码模态框 -->
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

    <!-- 页脚 -->
    <footer>
      <AppFooter />
    </footer>

    <!-- 悬浮按钮组件 -->
    <FloatButtons />
  </div>
  <div v-if="systemConfig?.maintenance_mode" class="fixed inset-0 z-[1000000] flex items-center justify-center bg-gradient-to-br from-yellow-100/80 via-gray-900/90 to-yellow-200/80 backdrop-blur-sm">
    <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-2xl px-8 py-10 flex flex-col items-center max-w-xs w-full border border-yellow-200 dark:border-yellow-700">
      <i class="fas fa-tools text-yellow-500 text-5xl mb-6 animate-bounce-slow"></i>
      <h3 class="text-2xl font-extrabold text-yellow-600 dark:text-yellow-400 mb-2 tracking-wide drop-shadow">系统维护中</h3>
      <p class="text-base text-gray-600 dark:text-gray-300 mb-6 text-center leading-relaxed">
        我们正在进行系统升级和维护，预计很快恢复服务。<br>
        请稍后再试，感谢您的理解与支持！
      </p>
      <!-- 动态点点动画 -->
      <div class="flex space-x-1 mt-2">
        <span class="w-2 h-2 bg-yellow-400 rounded-full animate-blink"></span>
        <span class="w-2 h-2 bg-yellow-500 rounded-full animate-blink delay-200"></span>
        <span class="w-2 h-2 bg-yellow-600 rounded-full animate-blink delay-400"></span>
      </div>
    </div>
  </div>

  <!-- 开发环境缓存信息组件 -->
  <SystemConfigCacheInfo />
</template>

<script setup lang="ts">
// 获取运行时配置
const config = useRuntimeConfig()

import { useResourceApi, useStatsApi, usePanApi, useSearchStatsApi } from '~/composables/useApi'
import SystemConfigCacheInfo from '~/components/SystemConfigCacheInfo.vue'

const resourceApi = useResourceApi()
const statsApi = useStatsApi()
const panApi = usePanApi()

// 路由参数已通过自动导入提供，直接使用
const route = useRoute()
const router = useRouter()

// 使用系统配置Store（带缓存支持）
const { useSystemConfigStore } = await import('~/stores/systemConfig')
const systemConfigStore = useSystemConfigStore()

// 初始化系统配置（会自动使用缓存）
await systemConfigStore.initConfig()

// 检查并自动刷新即将过期的缓存
await systemConfigStore.checkAndRefreshCache()

// 获取平台名称的辅助函数
const getPlatformName = (platformId: string) => {
  if (!platformId) return ''
  const platformList = (platforms.value || []) as any[]
  const platform = platformList.find((p: any) => p.id == platformId)
  return platform?.name || ''
}

// 动态生成页面标题和meta信息 - 使用缓存的系统配置
const pageTitle = computed(() => {
  try {
    const config = systemConfigStore.config
    const siteTitle = config?.site_title || '老九网盘资源数据库'
    const searchKeyword = (route.query?.search) ? route.query.search as string : ''
    const platformId = (route.query?.platform) ? route.query.platform as string : ''
    const platformName = getPlatformName(platformId)
    let title = siteTitle

    // 根据搜索条件组合标题
    if (searchKeyword && platformName) {
      title = `${searchKeyword} - ${platformName} - ${siteTitle}`
    } else if (searchKeyword) {
      title = `${searchKeyword} - 搜索结果 - ${siteTitle}`
    } else if (platformName) {
      title = `${platformName} - ${siteTitle}`
    } else {
      title = `${siteTitle} - 首页`
    }

    return title
  } catch (error) {
    console.error('pageTitle computed error:', error)
    return '老九网盘资源数据库 - 首页'
  }
})

const pageDescription = computed(() => {
  try {
    const config = systemConfigStore.config
    const baseDescription = config?.site_description || '老九网盘资源管理系统， 一个现代化的网盘资源数据库，支持多网盘自动化转存分享，支持百度网盘，阿里云盘，夸克网盘， 天翼云盘，迅雷云盘，123云盘，115网盘，UC网盘'

    const searchKeyword = (route.query && route.query.search) ? route.query.search as string : ''
    const platformId = (route.query && route.query.platform) ? route.query.platform as string : ''
    const platformName = getPlatformName(platformId)

    let description = baseDescription

    // 根据搜索条件优化描述
    if (searchKeyword && platformName) {
      description = `在${platformName}中搜索"${searchKeyword}"的相关资源。${baseDescription}提供海量${searchKeyword}资源下载，支持多网盘平台。`
    } else if (searchKeyword) {
      description = `搜索"${searchKeyword}"的相关资源。${baseDescription}提供海量${searchKeyword}资源下载，支持百度网盘、阿里云盘、夸克网盘等多个平台。`
    } else if (platformName) {
      description = `${platformName}资源专区。${baseDescription}专门收录${platformName}平台的优质资源，每日更新。`
    }

    return description
  } catch (error) {
    console.error('pageDescription computed error:', error)
    return '老九网盘资源管理系统， 一个现代化的网盘资源数据库，支持多网盘自动化转存分享，支持百度网盘，阿里云盘，夸克网盘， 天翼云盘，迅雷云盘，123云盘，115网盘，UC网盘'
  }
})

const pageKeywords = computed(() => {
  try {
    const config = systemConfigStore.config
    const baseKeywords = config?.keywords || '网盘资源,资源管理,数据库'

    const searchKeyword = (route.query && route.query.search) ? route.query.search as string : ''
    const platformId = (route.query && route.query.platform) ? route.query.platform as string : ''
    const platformName = getPlatformName(platformId)

    let keywords = baseKeywords

    // 根据搜索条件添加关键词
    if (searchKeyword) {
      keywords = `${searchKeyword},${baseKeywords},${searchKeyword}下载,${searchKeyword}资源`
    }

    if (platformName) {
      keywords = keywords ? `${platformName},${keywords}` : platformName
    }

    return keywords
  } catch (error) {
    console.error('pageKeywords computed error:', error)
    return '网盘资源,资源管理,数据库'
  }
})

// 设置页面SEO
const { initSystemConfig, setPageSeo, systemConfig: seoSystemConfig } = useGlobalSeo()

// 更新页面SEO的函数 - 合并所有SEO设置到一个函数中
const updatePageSeo = () => {
  // 使用动态计算的标题，而不是默认的"首页"
  setPageSeo(pageTitle.value, {
    description: pageDescription.value,
    keywords: pageKeywords.value,
    ogImage: '/assets/images/og.webp'  // 使用默认的OG图片
  })

  // 设置HTML属性和canonical链接
  const config = useRuntimeConfig()
  const baseUrl = config.public.siteUrl || 'https://yourdomain.com' // 从环境变量获取
  const params = new URLSearchParams()
  if (route.query?.search) params.set('search', route.query.search as string)
  if (route.query?.platform) params.set('platform', route.query.platform as string)
  const queryString = params.toString()
  const canonicalUrl = queryString ? `${baseUrl}?${queryString}` : baseUrl

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
      // Open Graph 标签
      { property: 'og:type', content: 'website' },
      { property: 'og:url', content: canonicalUrl },
      { property: 'og:title', content: pageTitle.value },
      { property: 'og:description', content: pageDescription.value },
      { property: 'og:image', content: `${baseUrl}/assets/images/og.webp` },
      { property: 'og:image:width', content: '1200' },
      { property: 'og:image:height', content: '630' },
      { property: 'og:site_name', content: (seoSystemConfig.value && seoSystemConfig.value.site_title) || '老九网盘资源数据库' },
      { property: 'og:locale', content: 'zh_CN' },
      
      // Twitter Card 标签
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: pageTitle.value },
      { name: 'twitter:description', content: pageDescription.value },
      { name: 'twitter:image', content: `${baseUrl}/assets/images/og.webp` },
      
      // 额外的 SEO meta 标签
      { name: 'author', content: (seoSystemConfig.value && seoSystemConfig.value.site_title) || '老九网盘资源数据库' },
      { name: 'robots', content: 'index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1' },
      { name: 'googlebot', content: 'index, follow' },
      { name: 'bingbot', content: 'index, follow' }
    ],
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          "@context": "https://schema.org",
          "@type": "WebSite",
          "@id": `${baseUrl}#website`,
          "url": baseUrl,
          "name": (seoSystemConfig.value && seoSystemConfig.value.site_title) || '老九网盘资源数据库',
          "description": pageDescription.value,
          "image": `${baseUrl}/assets/images/og.webp`,
          "potentialAction": {
            "@type": "SearchAction",
            "target": {
              "@type": "EntryPoint",
              "urlTemplate": `${baseUrl}?search={search_term_string}`
            },
            "query-input": "required name=search_term_string"
          },
          "publisher": {
            "@type": "Organization",
            "@id": `${baseUrl}#organization`,
            "name": (seoSystemConfig.value && seoSystemConfig.value.site_title) || '老九网盘资源数据库',
            "url": baseUrl,
            "logo": {
              "@type": "ImageObject",
              "url": `${baseUrl}/assets/images/logo.webp`
            }
          }
        })
      }
    ]
  })
}

onBeforeMount(async () => {
  await initSystemConfig()
  updatePageSeo()
})

// 监听路由变化和系统配置数据，当搜索条件或配置改变时更新SEO
watch(
  () => [route.query?.search, route.query?.platform, systemConfigStore.config],
  () => {
    // 使用nextTick确保响应式数据已更新
    nextTick(() => {
      updatePageSeo()
    })
  },
  { deep: true }
)

// 响应式数据
const showLinkModal = ref(false)
const selectedResource = ref<any>(null)
const pageLoading = ref(false)

// 使用 ClientOnly 包装器来处理认证状态
const authInitialized = ref(false)

// 用户状态管理
const userStore = useUserStore()

// 图片URL处理
const { getImageUrl } = useImageUrl()

// Logo错误处理
const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/assets/images/logo.webp'
}

// 获取资源图片URL，如果没有则返回随机默认封面
const getResourceImageUrl = (resource: any) => {
  // console.log('Resource data:', resource)
  // 如果资源有图片，使用资源图片（优先检查image_url，其次检查cover）
  if (resource.image_url) {
    return getImageUrl(resource.image_url)
  }

  if (resource.cover) {
    return getImageUrl(resource.cover)
  }

  // 否则随机选择默认封面图片 (cover1.webp 到 cover8.webp)
  const randomNum = Math.floor(Math.random() * 8) + 1
  return `/assets/images/cover${randomNum}.webp`
}

// 处理资源图片加载错误
const handleResourceImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  // 重新设置一个随机的默认封面图片
  const randomNum = Math.floor(Math.random() * 8) + 1
  img.src = `/assets/images/cover${randomNum}.webp`
}

// 使用 useAsyncData 获取资源数据
const { data: resourcesData, pending, refresh } = await useAsyncData(
  () => `resources-1-${route.query?.search || ''}-${route.query?.platform || ''}`,
  async () => {
    // 如果有搜索关键词，使用带搜索参数的资源接口（后端会优先使用Meilisearch）
    if (route.query?.search) {
      return await resourceApi.getResources({
        page: 1,
        page_size: 50,
        search: route.query.search as string,
        pan_id: route.query.platform as string || ''
      })
    } else {
      // 没有搜索关键词时，使用普通资源接口获取最新数据
      return await resourceApi.getResources({
        page: 1,
        page_size: 50,
        pan_id: route.query?.platform as string || ''
      })
    }
  }
)

// 获取统计数据 - 客户端渲染，30分钟缓存
const getCacheKey = (prefix: string) => {
  const cacheMinutes = Math.floor(Date.now() / (30 * 60 * 1000)) // 30分钟缓存
  return `${prefix}-${cacheMinutes}`
}

const { data: statsData, error: statsError } = await useAsyncData(
  getCacheKey('stats'),
  () => statsApi.getStats(),
  {
    server: false, // 客户端渲染
    lazy: true // 延迟加载，不阻塞首屏
  }
)

// 获取平台数据 - 30分钟缓存
const { data: platformsData, error: platformsError } = await useAsyncData(
  getCacheKey('platforms'),
  () => panApi.getPans(),
  {
    lazy: true // 延迟加载
  }
)

// 系统配置已在顶部获取，这里处理错误
const systemConfigError = ref(null)

// 错误处理
const notification = ref()

// 监听错误
watch(statsError, (error) => {
  if (error && process.client) {
    console.error('获取统计数据失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取统计数据失败',
      duration: 5000
    })
  }
})

watch(platformsError, (error) => {
  if (error && process.client) {
    console.error('获取平台数据失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取平台数据失败',
      duration: 5000
    })
  }
})

watch(systemConfigError, (error) => {
  if (error && process.client) {
    console.error('获取系统配置失败:', error)
    notification.value = useNotification()
    notification.value.error({
      content: error.message || '获取系统配置失败',
      duration: 5000
    })
  }
})

// 从 SSR 数据中获取值
const safeResources = computed(() => {
  const data = resourcesData.value as any
  let resources: any[] = []

  // 处理嵌套的data结构：{data: {data: [...], total: ...}}
  if (data?.data?.data && Array.isArray(data.data.data)) {
    resources = data.data.data
    console.log('第一层嵌套资源:', resources)
  }
  // 处理直接的data结构：{data: [...], total: ...}
  else if (data?.data && Array.isArray(data.data)) {
    resources = data.data
    // console.log('第二层嵌套资源:', resources)
  }
  // 处理直接的数组结构
  else if (Array.isArray(data)) {
    resources = data
    // console.log('直接数组结构:', data)
  }

  // 根据 key 字段去重
  if (resources.length > 0) {
    const keyMap = new Map()
    const deduplicatedResources: any[] = []

    for (const resource of resources) {
      const key = resource.key
      if (!keyMap.has(key)) {
        keyMap.set(key, true)
        deduplicatedResources.push(resource)
      }
    }

    // console.log(`去重前: ${resources.length} 个资源, 去重后: ${deduplicatedResources.length} 个资源`)
    return deduplicatedResources
  }

  // console.log('未匹配到任何数据结构')
  return []
})
// 统计数据加载状态
const statsLoading = computed(() => !statsData.value)

// 安全的统计数据，提供默认值
const safeStats = computed(() => {
  const data = statsData.value as any
  // 如果数据还在加载中，返回 null 以便显示加载状态
  if (!data) return null
  
  // 数据加载完成，返回实际数据
  return data || { total_resources: 0, total_categories: 0, total_tags: 0, total_views: 0, today_resources: 0 }
})
const platforms = computed(() => (platformsData.value as any) || [])
const systemConfig = computed(() => systemConfigStore.config || { site_title: '老九网盘资源数据库' })
const safeLoading = computed(() => pending.value)


// 从路由参数获取当前状态
const searchQuery = ref(route.query?.search as string || '')
const selectedPlatform = computed(() => route.query?.platform as string || '')

// 记录搜索统计的函数
const recordSearchStats = (keyword: string) => {
  if (!keyword || keyword.trim().length === 0) {
    // console.log('搜索关键词为空，跳过统计记录')
    return
  }
  
  const trimmedKeyword = keyword.trim()
  // console.log('记录搜索统计:', trimmedKeyword)
  
  // 延迟执行，确保页面完全加载
  setTimeout(() => {
    const searchStatsApi = useSearchStatsApi()
    searchStatsApi.recordSearch({ keyword: trimmedKeyword }).catch(err => {
      console.error('记录搜索统计失败:', err)
    })
  }, 0)
}

const handleSearch = () => {
  const params = new URLSearchParams()
  if (searchQuery.value) params.set('search', searchQuery.value)
  if (selectedPlatform.value) params.set('platform', selectedPlatform.value)
  window.location.href = `/?${params.toString()}`
}

// 监听统计数据加载完成，触发数字动画
watch(statsData, (newData) => {
  if (newData && process.client) {
    // 使用 nextTick 确保 DOM 已更新
    nextTick(() => {
      animateCounters()
    })
  }
})

// 初始化认证状态
onMounted(() => {
  // 初始化认证状态
  authInitialized.value = true

  // 如果统计数据已经加载（SSR 或缓存），立即执行动画
  if (statsData.value) {
    animateCounters()
  }

  // 页面挂载完成时，如果有搜索关键词，记录搜索统计
  if (process.client && route.query?.search) {
    const searchKeyword = route.query.search as string
    recordSearchStats(searchKeyword)
  } else {
    console.log('无搜索参数，跳过统计记录')
  }
})



// 获取平台名称
const getPlatformIcon = (panId: string | number) => {
  const platform = (platforms.value as any).find((p: any) => p.id == panId)
  return platform?.icon || '未知平台'
}

// 注意：链接访问统计已整合到 getResourceLink API 中

// 导航到详情页
const navigateToDetail = (key: string) => {
  router.push(`/r/${key}`)
}

// 切换链接显示（保留用于其他可能的用途）
const toggleLink = async (resource: any) => {
  navigateToDetail(resource.key)
}

// 复制到剪贴板
const copyToClipboard = async (text: any) => {
  try {
    await navigator.clipboard.writeText(text)
    if (process.client) {
      const button = document.querySelector('.show-link-btn')
      if (button) {
        const originalText = button.innerHTML
        button.innerHTML = '<i class="fas fa-check"></i> 已复制'
        button.classList.add('bg-green-600')
        setTimeout(() => {
          button.innerHTML = originalText
          button.classList.remove('bg-green-600')
        }, 2000)
      }
    }
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 格式化相对时间
const formatRelativeTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)
  const diffWeek = Math.floor(diffDay / 7)
  const diffMonth = Math.floor(diffDay / 30)
  const diffYear = Math.floor(diffDay / 365)

  const isToday = date.toDateString() === now.toDateString()

  // 处理今天更新的情况
  if (isToday) {
    if (diffMin < 1) {
      return '<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>刚刚更新</span>'
    } else if (diffHour < 1) {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffMin}分钟前</span>`
    } else {
      return `<span class="text-pink-600 font-medium flex items-center"><i class="fas fa-circle-dot text-xs mr-1 animate-pulse"></i>${diffHour}小时前</span>`
    }
  }

  // 处理昨天更新的情况 - 显示具体时间
  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  const isYesterday = date.toDateString() === yesterday.toDateString()

  if (isYesterday) {
    if (diffHour < 24) {
      // 昨天但不足24小时
      if (diffHour < 1) {
        return `<span class="text-gray-600">${diffMin}分钟前</span>`
      } else {
        return `<span class="text-gray-600">${diffHour}小时前</span>`
      }
    } else {
      // 超过24小时但仍然是昨天
      return `<span class="text-gray-600">${diffDay}天前</span>`
    }
  }

  // 处理其他情况
  if (diffDay < 7) {
    return `<span class="text-gray-600">${diffDay}天前</span>`
  } else if (diffWeek < 4) {
    return `<span class="text-gray-600">${diffWeek}周前</span>`
  } else if (diffMonth < 12) {
    return `<span class="text-gray-600">${diffMonth}个月前</span>`
  } else {
    return `<span class="text-gray-600">${diffYear}年前</span>`
  }
}

// 检查是否为今天更新
const isUpdatedToday = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  return date.toDateString() === now.toDateString()
}

// 数字动画效果
const animateCounters = () => {
  if (!process.client) return
  
  const counters = document.querySelectorAll('.count-up')
  const speed = 200
  
  counters.forEach((counter) => {
    const target = parseInt(counter.getAttribute('data-target') || '0')
    const increment = Math.ceil(target / speed)
    let count = 0
    
    const updateCount = () => {
      if (count < target) {
        count += increment
        if (count > target) count = target
        counter.textContent = count.toString()
        setTimeout(updateCount, 1)
      } else {
        counter.textContent = target.toString()
      }
    }
    
    updateCount()
  })
}






</script>

<style scoped>
/* 屏幕阅读器专用样式 - 隐藏但对辅助技术可见 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.active-filter {
  @apply bg-slate-800 text-white;
}

.count-up {
  transition: all 0.3s ease;
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: .5;
  }
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
.header-container{
  background: url(/assets/images/banner.webp) center top/cover no-repeat,
  linear-gradient(
      to bottom, 
      rgba(0,0,0,0.1) 0%, 
      rgba(0,0,0,0.25) 100%
  );
}

/* 文本截断样式 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  word-wrap: break-word;
  word-break: break-word;
}

/* 表格单元格内容溢出控制 */
table td {
  overflow: hidden;
  word-wrap: break-word;
  word-break: break-word;
}

/* 确保flex容器不会溢出 */
.min-w-0 {
  min-width: 0;
}

/* 标签样式优化 */
.resource-tag {
  transition: all 0.2s ease;
}

.resource-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

/* 移动端按钮专用样式 */
.mobile-link-btn {
  border: 1px solid transparent;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 500;
  font-size: 11px;
  line-height: 1.2;
  transition: all 0.3s ease;
  min-height: 28px;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
}

.mobile-link-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
  transition: left 0.5s;
}

.mobile-link-btn:hover::before {
  left: 100%;
}

.mobile-link-btn:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  transform: translateY(-1px);
}

.mobile-link-btn:active {
  background: linear-gradient(135deg, #1d4ed8 0%, #172554 100%);
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.4);
}

.mobile-link-btn:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.6);
}
</style> 