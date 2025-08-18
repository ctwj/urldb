<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">站点配置</h1>
        <p class="text-gray-600 dark:text-gray-400">管理网站基本信息和设置</p>
      </div>
      <n-button type="primary" @click="saveConfig" :loading="saving">
        <template #icon>
          <i class="fas fa-save"></i>
        </template>
        保存配置
      </n-button>
    </div>

        <!-- 配置表单 -->
    <n-card>
      <!-- 顶部Tabs -->
      <n-tabs
        v-model:value="activeTab"
        type="line"
        animated
        class="mb-6"
      >
        <n-tab-pane name="basic" tab="基本信息">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 网站标题 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站标题</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站的主要标识，显示在浏览器标签页和搜索结果中</span>
                </div>
                <n-input
                  v-model:value="configForm.site_title"
                  placeholder="请输入网站标题"
                />
              </div>

              <!-- 网站描述 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站描述</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站的简要介绍，用于SEO和社交媒体分享</span>
                </div>
                <n-input
                  v-model:value="configForm.site_description"
                  placeholder="请输入网站描述"
                />
              </div>

              <!-- 关键词 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">关键词</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">用于SEO优化，多个关键词用逗号分隔</span>
                </div>
                <n-input
                  v-model:value="configForm.keywords"
                  placeholder="请输入关键词，用逗号分隔"
                />
              </div>

              <!-- 网站Logo -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">网站Logo</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">选择网站Logo图片，建议使用正方形图片</span>
                </div>
                <div class="flex items-center space-x-4">
                  <div v-if="configForm.site_logo" class="flex-shrink-0">
                    <n-image
                      :src="getImageUrl(configForm.site_logo)"
                      alt="网站Logo"
                      width="80"
                      height="80"
                      object-fit="cover"
                      class="rounded-lg border"
                    />
                  </div>
                  <div class="flex-1">
                    <n-button type="primary" @click="openLogoSelector">
                      <template #icon>
                        <i class="fas fa-image"></i>
                      </template>
                      {{ configForm.site_logo ? '更换Logo' : '选择Logo' }}
                    </n-button>
                    <n-button v-if="configForm.site_logo" @click="clearLogo" class="ml-2">
                      <template #icon>
                        <i class="fas fa-times"></i>
                      </template>
                      清除
                    </n-button>
                  </div>
                </div>
              </div>

              <!-- 版权信息 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">版权信息</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">网站底部的版权声明信息</span>
                </div>
                <n-input
                  v-model:value="configForm.copyright"
                  placeholder="请输入版权信息"
                />
              </div>
            </div>
          </n-form>
        </n-tab-pane>



        <n-tab-pane name="security" tab="安全设置">
          
          <n-form
            ref="formRef"
            :model="configForm"
            :rules="rules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <div class="space-y-6">
              <!-- 维护模式 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">维护模式</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后网站将显示维护页面，暂停用户访问</span>
                </div>
                <n-switch v-model:value="configForm.maintenance_mode" />
              </div>

              <!-- 违禁词 -->
              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <label class="text-base font-semibold text-gray-800 dark:text-gray-200">违禁词</label>
                    <span class="text-xs text-gray-500 dark:text-gray-400">包含这些词汇的资源将被过滤，多个词汇用逗号分隔</span>
                  </div>
                  <a 
                    href="https://raw.githubusercontent.com/ctwj/urldb/refs/heads/main/db/forbidden.txt" 
                    target="_blank" 
                    class="text-xs text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 underline"
                  >
                    开源违禁词
                  </a>
                </div>
                <n-input
                  v-model:value="configForm.forbidden_words"
                  placeholder="请输入违禁词，用逗号分隔"
                  type="textarea"
                  :rows="4"
                />
              </div>

              <!-- 开启注册 -->
              <div class="space-y-2">
                <div class="flex items-center space-x-2">
                  <label class="text-base font-semibold text-gray-800 dark:text-gray-200">开启注册</label>
                  <span class="text-xs text-gray-500 dark:text-gray-400">开启后用户才能注册新账号，关闭后注册页面将显示"当前系统已关闭注册功能"</span>
                </div>
                <n-switch v-model:value="configForm.enable_register" />
              </div>
            </div>
          </n-form>
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- Logo选择模态框 -->
    <n-modal v-model:show="showLogoSelector" preset="card" title="选择Logo图片" style="width: 90vw; max-width: 1200px; max-height: 80vh;">
      <div class="space-y-4">
        <!-- 搜索 -->
        <div class="flex gap-4">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索文件名..."
            @keyup.enter="handleSearch"
            class="flex-1"
            clearable
          >
            <template #prefix>
              <i class="fas fa-search"></i>
            </template>
          </n-input>
          
          <n-button type="primary" @click="handleSearch" class="w-20">
            <template #icon>
              <i class="fas fa-search"></i>
            </template>
            搜索
          </n-button>
        </div>

        <!-- 文件列表 -->
        <div v-if="loading" class="flex items-center justify-center py-8">
          <n-spin size="large" />
        </div>

        <div v-else-if="fileList.length === 0" class="text-center py-8">
          <i class="fas fa-file-upload text-4xl text-gray-400 mb-4"></i>
          <p class="text-gray-500">暂无图片文件</p>
        </div>

        <div v-else class="file-grid">
          <div 
            v-for="file in fileList" 
            :key="file.id" 
            class="file-item cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800 rounded-lg p-3 transition-colors"
            :class="{ 'bg-blue-50 dark:bg-blue-900/20 border-2 border-blue-300 dark:border-blue-600': selectedFileId === file.id }"
            @click="selectFile(file)"
          >
            <div class="image-preview">
              <n-image
                :src="getImageUrl(file.access_url)"
                :alt="file.original_name"
                :lazy="false"
                object-fit="cover"
                class="preview-image rounded"
                @error="handleImageError"
                @load="handleImageLoad"
              />

              <div class="image-info mt-2">
                <div class="file-name text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                  {{ file.original_name }}
                </div>
                <div class="file-size text-xs text-gray-500 dark:text-gray-400">
                  {{ formatFileSize(file.file_size) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <n-pagination
            v-model:page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-count="Math.ceil(pagination.total / pagination.pageSize)"
            :page-sizes="pagination.pageSizes"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </div>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showLogoSelector = false">取消</n-button>
          <n-button 
            type="primary" 
            @click="confirmSelection"
            :disabled="!selectedFileId"
          >
            确认选择
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})


import { useImageUrl } from '~/composables/useImageUrl'
import { useConfigChangeDetection } from '~/composables/useConfigChangeDetection'

const notification = useNotification()
const { getImageUrl } = useImageUrl()
const formRef = ref()
const saving = ref(false)
const activeTab = ref('basic')

// Logo选择器相关数据
const showLogoSelector = ref(false)
const loading = ref(false)
const fileList = ref<any[]>([])
const selectedFileId = ref<number | null>(null)
const searchKeyword = ref('')

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0,
  pageSizes: [10, 20, 50, 100]
})

// 配置表单数据类型
interface SiteConfigForm {
  site_title: string
  site_description: string
  keywords: string
  copyright: string
  site_logo: string
  maintenance_mode: boolean
  enable_register: boolean
  forbidden_words: string
  enable_sitemap: boolean
  sitemap_update_frequency: string
}

// 使用配置改动检测
const {
  setOriginalConfig,
  updateCurrentConfig,
  getChangedConfig,
  hasChanges,
  getChangedDetails,
  updateOriginalConfig,
  saveConfig: saveConfigWithDetection
} = useConfigChangeDetection<SiteConfigForm>({
  debug: true,
  // 字段映射：前端字段名 -> 后端字段名
  fieldMapping: {
    site_title: 'site_title',
    site_description: 'site_description',
    keywords: 'keywords',
    copyright: 'copyright',
    site_logo: 'site_logo',
    maintenance_mode: 'maintenance_mode',
    enable_register: 'enable_register',
    forbidden_words: 'forbidden_words',
    enable_sitemap: 'enable_sitemap',
    sitemap_update_frequency: 'sitemap_update_frequency'
  }
})

// 配置表单数据
const configForm = ref<SiteConfigForm>({
  site_title: '',
  site_description: '',
  keywords: '',
  copyright: '',
  site_logo: '',
  maintenance_mode: false,
  enable_register: false,
  forbidden_words: '',
  enable_sitemap: false,
  sitemap_update_frequency: 'daily'
})



// 表单验证规则
const rules = {
  site_title: {
    required: true,
    message: '请输入网站标题',
    trigger: 'blur'
  },
  site_description: {
    required: true,
    message: '请输入网站描述',
    trigger: 'blur'
  }
}

// 获取系统配置
const fetchConfig = async () => {
  try {
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    const response = await systemConfigApi.getSystemConfig() as any
    
    if (response) {
      const configData = {
        site_title: response.site_title || '',
        site_description: response.site_description || '',
        keywords: response.keywords || '',
        copyright: response.copyright || '',
        site_logo: response.site_logo || '',
        maintenance_mode: response.maintenance_mode || false,
        enable_register: response.enable_register || false,
        forbidden_words: response.forbidden_words || '',
        enable_sitemap: response.enable_sitemap || false,
        sitemap_update_frequency: response.sitemap_update_frequency || 'daily'
      }
      
      // 设置表单数据和原始数据
      configForm.value = { ...configData }
      setOriginalConfig(configData)
    }
  } catch (error) {
    console.error('获取系统配置失败:', error)
    notification.error({
      content: '获取系统配置失败',
      duration: 3000
    })
  }
}



// 保存配置
const saveConfig = async () => {
  try {
    await formRef.value?.validate()
    
    saving.value = true
    
    // 更新当前配置数据
    updateCurrentConfig({
      site_title: configForm.value.site_title,
      site_description: configForm.value.site_description,
      keywords: configForm.value.keywords,
      copyright: configForm.value.copyright,
      site_logo: configForm.value.site_logo,
      maintenance_mode: configForm.value.maintenance_mode,
      enable_register: configForm.value.enable_register,
      forbidden_words: configForm.value.forbidden_words,
      enable_sitemap: configForm.value.enable_sitemap,
      sitemap_update_frequency: configForm.value.sitemap_update_frequency
    })
    
    const { useSystemConfigApi } = await import('~/composables/useApi')
    const systemConfigApi = useSystemConfigApi()
    
    // 使用通用保存函数
    const result = await saveConfigWithDetection(
      systemConfigApi.updateSystemConfig,
      {
        onlyChanged: true,
        includeAllFields: true
      },
      // 成功回调
      async () => {
        notification.success({
          content: '站点配置保存成功',
          duration: 3000
        })
        
        // 刷新系统配置状态，确保顶部导航同步更新
        const { useSystemConfigStore } = await import('~/stores/systemConfig')
        const systemConfigStore = useSystemConfigStore()
        await systemConfigStore.initConfig(true, true)
      },
      // 错误回调
      (error) => {
        console.error('保存站点配置失败:', error)
        notification.error({
          content: '保存站点配置失败',
          duration: 3000
        })
      }
    )
    
    // 如果没有改动，显示提示
    if (result && result.message === '没有检测到任何改动') {
      notification.info({
        content: '没有检测到任何改动',
        duration: 3000
      })
    }
  } finally {
    saving.value = false
  }
}

// Logo选择器方法
const openLogoSelector = () => {
  showLogoSelector.value = true
  loadFileList()
}

const clearLogo = () => {
  configForm.value.site_logo = ''
}

const loadFileList = async () => {
  try {
    loading.value = true
    const { useFileApi } = await import('~/composables/useFileApi')
    const fileApi = useFileApi()
    
    const response = await fileApi.getFileList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
      search: searchKeyword.value,
      fileType: 'image', // 只获取图片文件
      status: 'active'   // 只获取正常状态的文件
    }) as any

    if (response && response.data) {
      fileList.value = response.data.files || []
      pagination.value.total = response.data.total || 0
      console.log('获取到的图片文件:', fileList.value) // 调试信息
      
      // 添加图片URL处理调试
      fileList.value.forEach(file => {
        console.log('图片文件详情:', {
          id: file.id,
          name: file.original_name,
          accessUrl: file.access_url,
          processedUrl: getImageUrl(file.access_url),
          fileType: file.file_type,
          mimeType: file.mime_type
        })
      })
    }
  } catch (error) {
    console.error('获取文件列表失败:', error)
    notification.error({
      content: '获取文件列表失败',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadFileList()
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  loadFileList()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize
  pagination.value.page = 1
  loadFileList()
}

const selectFile = (file: any) => {
  selectedFileId.value = file.id
}

const confirmSelection = () => {
  if (selectedFileId.value) {
    const file = fileList.value.find(f => f.id === selectedFileId.value)
    if (file) {
      configForm.value.site_logo = file.access_url
      showLogoSelector.value = false
      selectedFileId.value = null
    }
  }
}

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(1) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(1) + ' GB'
}

const handleImageError = (event: any) => {
  console.error('图片加载失败:', event)
}

const handleImageLoad = (event: any) => {
  console.log('图片加载成功:', event)
}

// 页面加载时获取配置
onMounted(() => {
  fetchConfig()
})


</script>

<style scoped>
/* 自定义样式 */
.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  max-height: 400px;
  overflow-y: auto;
}

.file-item {
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.file-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.preview-image {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
}



.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}
</style> 