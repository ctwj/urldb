<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">SEO管理</h1>
        <p class="text-gray-600 dark:text-gray-400">搜索引擎优化管理</p>
      </div>
    </template>

    <!-- 内容区 -->
    <template #content>
      <div class="config-content h-full">
        <!-- Tab导航 -->
        <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="google-index" tab="Google Index">
          <div class="tab-content-container">
            <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
              <!-- Google索引配置 -->
              <div class="mb-6">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Google索引配置</h3>
                <p class="text-gray-600 dark:text-gray-400">配置Google Search Console API和索引相关设置</p>
              </div>

        <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h4 class="font-medium text-gray-900 dark:text-white mb-2">Google索引功能</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                开启后系统将自动检查和提交URL到Google索引
              </p>
            </div>
            <n-switch
              v-model:value="googleIndexConfig.enabled"
              @update:value="updateGoogleIndexConfig"
              :loading="configLoading"
              size="large"
            >
              <template #checked>已开启</template>
              <template #unchecked>已关闭</template>
            </n-switch>
          </div>

          <!-- 配置详情 -->
          <div class="border-t border-gray-200 dark:border-gray-600 pt-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">站点URL</label>
                <n-input
                  :value="systemConfig?.site_url || '站点URL未配置'"
                  :disabled="true"
                  placeholder="请先在站点配置中设置站点URL"
                >
                  <template #prefix>
                    <i class="fas fa-globe text-gray-400"></i>
                  </template>
                </n-input>
                <!-- 所有权验证按钮 -->
                <div class="mt-3">
                  <n-button
                    type="info"
                    size="small"
                    ghost
                    @click="showVerificationModal = true"
                  >
                    <template #icon>
                      <i class="fas fa-shield-alt"></i>
                    </template>
                    所有权验证
                  </n-button>
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">凭据文件路径</label>
                <div class="flex flex-col space-y-2">
                  <n-input
                    v-model:value="googleIndexConfig.credentialsFile"
                    placeholder="点击上传按钮选择文件"
                    :disabled="true"
                  />
                  <div class="flex space-x-2">
                    <!-- 申请凭据按钮 -->
                    <n-button
                      size="small"
                      type="warning"
                      ghost
                      @click="showCredentialsGuide = true"
                    >
                      <template #icon>
                        <i class="fas fa-question-circle"></i>
                      </template>
                      申请凭据
                    </n-button>
                    <!-- 上传按钮 -->
                    <n-button
                      size="small"
                      type="primary"
                      ghost
                      @click="selectCredentialsFile"
                    >
                      <template #icon>
                        <i class="fas fa-upload"></i>
                      </template>
                      上传凭据
                    </n-button>
                    <!-- 验证按钮 -->
                    <n-button
                      size="small"
                      type="info"
                      ghost
                      @click="validateCredentials"
                      :loading="validatingCredentials"
                      :disabled="!googleIndexConfig.credentialsFile"
                    >
                      <template #icon>
                        <i class="fas fa-check-circle"></i>
                      </template>
                      验证凭据
                    </n-button>
                  </div>
                  <!-- 隐藏的文件输入 -->
                  <input
                    type="file"
                    ref="credentialsFileInput"
                    accept=".json"
                    @change="handleCredentialsFileSelect"
                    style="display: none;"
                  />
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">检查间隔(分钟)</label>
                <n-input-number
                  v-model:value="googleIndexConfig.checkInterval"
                  :min="1"
                  :max="1440"
                  @update:value="updateGoogleIndexConfig"
                  :disabled="!googleIndexConfig.credentialsFile"
                  style="width: 100%"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">批处理大小</label>
                <n-input-number
                  v-model:value="googleIndexConfig.batchSize"
                  :min="1"
                  :max="1000"
                  @update:value="updateGoogleIndexConfig"
                  :disabled="!googleIndexConfig.credentialsFile"
                  style="width: 100%"
                />
              </div>
            </div>
          </div>

          <!-- 凭据状态 -->
          <div v-if="credentialsStatus" class="mt-4 p-3 rounded-lg border"
            :class="{
              'bg-green-50 border-green-200 text-green-700 dark:bg-green-900/20 dark:border-green-800 dark:text-green-300': credentialsStatus === 'valid',
              'bg-yellow-50 border-yellow-200 text-yellow-700 dark:bg-yellow-900/20 dark:border-yellow-800 dark:text-yellow-300': credentialsStatus === 'invalid',
              'bg-blue-50 border-blue-200 text-blue-700 dark:bg-blue-900/20 dark:border-blue-800 dark:text-blue-300': credentialsStatus === 'verifying'
            }"
          >
            <div class="flex items-center">
              <i
                :class="{
                  'fas fa-check-circle text-green-500 dark:text-green-400': credentialsStatus === 'valid',
                  'fas fa-exclamation-circle text-yellow-500 dark:text-yellow-400': credentialsStatus === 'invalid',
                  'fas fa-spinner fa-spin text-blue-500 dark:text-blue-400': credentialsStatus === 'verifying'
                }"
                class="mr-2"
              ></i>
              <span>{{ credentialsStatusMessage }}</span>
            </div>
          </div>
        </div>

        <!-- Google索引统计 -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                <i class="fas fa-link text-blue-600 dark:text-blue-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">总URL数</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ googleIndexStats.totalURLs || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                <i class="fas fa-check-circle text-green-600 dark:text-green-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">已索引</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ googleIndexStats.indexedURLs || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
                <i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">错误数</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ googleIndexStats.errorURLs || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <i class="fas fa-tasks text-purple-600 dark:text-purple-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">总任务数</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ googleIndexStats.totalTasks || 0 }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex flex-wrap gap-3 mb-6">
          <n-button
            type="primary"
            @click="manualCheckURLs"
            :loading="manualCheckLoading"
            size="large"
          >
            <template #icon>
              <i class="fas fa-search"></i>
            </template>
            手动检查URL
          </n-button>

          <n-button
            type="success"
            @click="submitSitemap"
            :loading="submitSitemapLoading"
            size="large"
          >
            <template #icon>
              <i class="fas fa-upload"></i>
            </template>
            提交网站地图
          </n-button>

          <n-button
            type="info"
            @click="refreshGoogleIndexStatus"
            size="large"
          >
            <template #icon>
              <i class="fas fa-sync-alt"></i>
            </template>
            刷新状态
          </n-button>
        </div>

        <!-- 任务列表 -->
        <div>
          <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4">索引任务列表</h4>
          <n-data-table
            :columns="googleIndexTaskColumns"
            :data="googleIndexTasks"
            :pagination="googleIndexPagination"
            :loading="tasksLoading"
            :bordered="false"
            striped
          />
        </div>
      </div>
    </div>
  </n-tab-pane>

  <n-tab-pane name="sitemap" tab="Sitemap管理">
    <div class="tab-content-container">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <!-- Sitemap配置 -->
        <div class="mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Sitemap配置</h3>
          <p class="text-gray-600 dark:text-gray-400">管理网站的Sitemap生成和配置</p>
        </div>

        <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h4 class="font-medium text-gray-900 dark:text-white mb-2">自动生成Sitemap</h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                开启后系统将定期自动生成Sitemap文件
              </p>
            </div>
            <n-switch
              v-model:value="sitemapConfig.autoGenerate"
              @update:value="updateSitemapConfig"
              :loading="configLoading"
              size="large"
            >
              <template #checked>已开启</template>
              <template #unchecked>已关闭</template>
            </n-switch>
          </div>

          <!-- 配置详情 -->
          <div class="border-t border-gray-200 dark:border-gray-600 pt-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">站点URL</label>
                <n-input
                  :value="systemConfig?.site_url || '站点URL未配置'"
                  :disabled="true"
                  placeholder="请先在站点配置中设置站点URL"
                >
                  <template #prefix>
                    <i class="fas fa-globe text-gray-400"></i>
                  </template>
                </n-input>
              </div>
              <div class="flex flex-col justify-end">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">最后生成时间</label>
                <div class="text-sm text-gray-600 dark:text-gray-400">
                  {{ sitemapConfig.lastGenerate || '尚未生成' }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sitemap统计 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
                <i class="fas fa-link text-blue-600 dark:text-blue-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">资源总数</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.total_resources || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
                <i class="fas fa-sitemap text-green-600 dark:text-green-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">页面数量</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.total_pages || 0 }}</p>
              </div>
            </div>
          </div>

          <div class="bg-purple-50 dark:bg-purple-900/20 rounded-lg p-4">
            <div class="flex items-center">
              <div class="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
                <i class="fas fa-history text-purple-600 dark:text-purple-400"></i>
              </div>
              <div class="ml-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">最后更新</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ sitemapStats.last_generate || 'N/A' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex flex-wrap gap-3 mb-6">
          <n-button
            type="primary"
            @click="generateSitemap"
            :loading="isGenerating"
            size="large"
          >
            <template #icon>
              <i class="fas fa-cog"></i>
            </template>
            生成Sitemap
          </n-button>

          <n-button
            type="success"
            @click="viewSitemap"
            size="large"
          >
            <template #icon>
              <i class="fas fa-external-link-alt"></i>
            </template>
            查看Sitemap
          </n-button>

          <n-button
            type="info"
            @click="refreshSitemapStatus"
            size="large"
          >
            <template #icon>
              <i class="fas fa-sync-alt"></i>
            </template>
            刷新状态
          </n-button>
        </div>

        <!-- 生成状态 -->
        <div v-if="generateStatus" class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <div class="flex items-center">
            <i class="fas fa-info-circle text-blue-500 dark:text-blue-400 mr-2"></i>
            <span class="text-blue-700 dark:text-blue-300">{{ generateStatus }}</span>
          </div>
        </div>
      </div>
    </div>
  </n-tab-pane>

  <n-tab-pane name="site-submit" tab="站点提交">
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
  </n-tab-pane>

  <n-tab-pane name="link-building" tab="外链建设">
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
            <n-button type="primary" @click="addNewLink">
              <template #icon>
                <i class="fas fa-plus"></i>
              </template>
              添加外链
            </n-button>
          </div>

          <n-data-table
            :columns="linkColumns"
            :data="linkList"
            :pagination="linkPagination"
            :loading="linkLoading"
            :bordered="false"
            striped
          />
        </div>
      </div>
    </div>
  </n-tab-pane>

      </n-tabs>
    </div>
  </template>
</AdminPageLayout>

<!-- URL检查模态框 -->
<n-modal v-model:show="urlCheckModal.show" preset="card" title="手动检查URL" style="max-width: 600px;">
  <div class="space-y-4">
    <p class="text-gray-600 dark:text-gray-400">输入要检查索引状态的URL，每行一个</p>
    <n-input
      v-model:value="urlCheckModal.urls"
      type="textarea"
      :autosize="{ minRows: 4, maxRows: 8 }"
      placeholder="https://yoursite.com/page1&#10;https://yoursite.com/page2"
    />
    <div class="flex justify-end space-x-2">
      <n-button @click="urlCheckModal.show = false">取消</n-button>
      <n-button type="primary" @click="confirmManualCheckURLs" :loading="manualCheckLoading">确认</n-button>
    </div>
  </div>
</n-modal>

<!-- 所有权验证模态框 -->
<n-modal v-model:show="showVerificationModal" preset="card" title="站点所有权验证" style="max-width: 600px;">
  <div class="space-y-6">
    <p class="text-gray-600 dark:text-gray-400">
      为了验证网站所有权，请将以下验证字符串添加到网站的HTML头部中：
    </p>

    <div class="space-y-4">
      <div class="bg-gray-50 dark:bg-gray-700/50 p-4 rounded-lg">
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">验证字符串</label>
        <n-input
          v-model:value="verificationCode"
          type="textarea"
          :autosize="{ minRows: 2, maxRows: 4 }"
          placeholder="输入Google Search Console或其他搜索引擎提供的验证字符串"
        />
      </div>

      <div class="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg border border-blue-200 dark:border-blue-800">
        <label class="block text-sm font-medium text-blue-700 dark:text-blue-300 mb-2">在页面中添加的代码</label>
        <div class="bg-white dark:bg-gray-800 p-3 rounded border">
          <code class="text-sm text-gray-800 dark:text-gray-200">
            &lt;meta name="google-site-verification" content="<span class="text-blue-600 dark:text-blue-400">{{ verificationCode || '验证字符串' }}</span>" /&gt;
          </code>
        </div>
      </div>
    </div>

    <div class="flex justify-end space-x-3 pt-2">
      <n-button @click="showVerificationModal = false">取消</n-button>
      <n-button type="primary" @click="saveVerificationCode" :loading="verificationCodeSaving">
        保存
      </n-button>
    </div>
  </div>
</n-modal>

<!-- 申请凭据说明抽屉 -->
<n-drawer v-model:show="showCredentialsGuide" :width="600" placement="right">
  <n-drawer-content title="如何申请Google Search Console API凭据" closable>
    <div class="space-y-6">
      <!-- 步骤1 -->
      <div class="border-l-4 border-blue-500 pl-4">
        <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          <i class="fas fa-number-1 text-blue-500 mr-2"></i>创建Google Cloud项目
        </h4>
        <p class="text-gray-600 dark:text-gray-400 mb-3">
          首先需要在Google Cloud Console中创建一个新项目或选择现有项目。
        </p>
        <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <li>访问 <a href="https://console.cloud.google.com/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Google Cloud Console</a></li>
          <li>点击顶部的项目选择器</li>
          <li>点击"新建项目"或选择现有项目</li>
          <li>输入项目名称，点击"创建"</li>
        </ol>
      </div>

      <!-- 步骤2 -->
      <div class="border-l-4 border-green-500 pl-4">
        <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          <i class="fas fa-number-2 text-green-500 mr-2"></i>启用Search Console API
        </h4>
        <p class="text-gray-600 dark:text-gray-400 mb-3">
          在项目中启用Google Search Console API。
        </p>
        <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <li>在导航菜单中选择"API和服务" > "库"</li>
          <li>搜索"Google Search Console API"</li>
          <li>点击搜索结果中的"Google Search Console API"</li>
          <li>点击"启用"按钮</li>
        </ol>
      </div>

      <!-- 步骤3 -->
      <div class="border-l-4 border-yellow-500 pl-4">
        <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          <i class="fas fa-number-3 text-yellow-500 mr-2"></i>创建服务账号
        </h4>
        <p class="text-gray-600 dark:text-gray-400 mb-3">
          创建服务账号并生成JSON密钥文件。
        </p>
        <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <li>在导航菜单中选择"API和服务" > "凭据"</li>
          <li>点击"创建凭据" > "服务账号"</li>
          <li>输入服务账号名称（如：google-index-api）</li>
          <li>点击"创建并继续"</li>
          <li>在角色选择中，选择"项目" > "编辑者"</li>
          <li>点击"继续"，然后点击"完成"</li>
        </ol>
      </div>

      <!-- 步骤4 -->
      <div class="border-l-4 border-purple-500 pl-4">
        <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          <i class="fas fa-number-4 text-purple-500 mr-2"></i>生成JSON密钥
        </h4>
        <p class="text-gray-600 dark:text-gray-400 mb-3">
          为服务账号生成JSON格式的密钥文件。
        </p>
        <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <li>在服务账号列表中找到刚创建的服务账号</li>
          <li>点击服务账号名称进入详情页面</li>
          <li>切换到"密钥"标签页</li>
          <li>点击"添加密钥" > "创建新密钥"</li>
          <li>选择"JSON"作为密钥类型</li>
          <li>点击"创建"，JSON文件将自动下载</li>
        </ol>
      </div>

      <!-- 步骤5 -->
      <div class="border-l-4 border-red-500 pl-4">
        <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          <i class="fas fa-number-5 text-red-500 mr-2"></i>验证网站所有权
        </h4>
        <p class="text-gray-600 dark:text-gray-400 mb-3">
          在Google Search Console中验证网站并添加服务账号权限。
        </p>
        <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
          <li>访问 <a href="https://search.google.com/search-console/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Google Search Console</a></li>
          <li>添加属性并验证网站所有权</li>
          <li>在设置中找到"用户和权限"</li>
          <li>点击"添加用户"，输入服务账号的邮箱地址</li>
          <li>授予"所有者"或"完整"权限</li>
        </ol>
      </div>

      <!-- 注意事项 -->
      <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
        <h5 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2">
          <i class="fas fa-exclamation-triangle mr-2"></i>重要注意事项
        </h5>
        <ul class="space-y-1 text-sm text-yellow-700 dark:text-yellow-300">
          <li>• 请妥善保管下载的JSON密钥文件，不要泄露给他人</li>
          <li>• 服务账号邮箱地址通常格式为：xxx@xxx.iam.gserviceaccount.com</li>
          <li>• API配额有限制，请合理使用避免超出限制</li>
          <li>• 确保网站已在Search Console中验证所有权</li>
        </ul>
      </div>

      <!-- 完成按钮 -->
      <div class="flex justify-end">
        <n-button type="primary" @click="showCredentialsGuide = false">
          我已了解
        </n-button>
      </div>
    </div>
  </n-drawer-content>
</n-drawer>
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'

// SEO管理页面
definePageMeta({
  layout: 'admin'
})

import { useMessage } from 'naive-ui'
import { useApi } from '~/composables/useApi'
import { ref, watch } from 'vue'

// 获取消息组件
const message = useMessage()

// 当前激活的Tab
const activeTab = ref('site-submit')

// 最后提交时间
const lastSubmitTime = ref({
  baidu: '',
  google: '',
  bing: '',
  sogou: '',
  shenma: '',
  so360: ''
})

// 外链统计
const linkStats = ref({
  total: 156,
  valid: 142,
  pending: 8,
  invalid: 6
})

// 外链列表
const linkList = ref([
  {
    id: 1,
    url: 'https://example1.com',
    title: '示例外链1',
    status: 'valid',
    domain: 'example1.com',
    created_at: '2024-01-15'
  },
  {
    id: 2,
    url: 'https://example2.com',
    title: '示例外链2',
    status: 'pending',
    domain: 'example2.com',
    created_at: '2024-01-16'
  },
  {
    id: 3,
    url: 'https://example3.com',
    title: '示例外链3',
    status: 'invalid',
    domain: 'example3.com',
    created_at: '2024-01-17'
  }
])

const linkLoading = ref(false)

// 分页配置
const linkPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    linkPagination.value.page = page
    loadLinkList()
  },
  onUpdatePageSize: (pageSize: number) => {
    linkPagination.value.pageSize = pageSize
    linkPagination.value.page = 1
    loadLinkList()
  }
})

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
          onClick: () => editLink(row)
        }, '编辑'),
        h('button', {
          class: 'text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300',
          onClick: () => deleteLink(row)
        }, '删除')
      ])
    }
  }
]

// 提交到百度
const submitToBaidu = () => {
  // 模拟提交
  lastSubmitTime.value.baidu = new Date().toLocaleString('zh-CN')
  message.success('已提交到百度')
}

// 提交到谷歌
const submitToGoogle = () => {
  // 模拟提交
  lastSubmitTime.value.google = new Date().toLocaleString('zh-CN')
  message.success('已提交到谷歌')
}

// 提交到必应
const submitToBing = () => {
  // 模拟提交
  lastSubmitTime.value.bing = new Date().toLocaleString('zh-CN')
  message.success('已提交到必应')
}

// 提交到搜狗
const submitToSogou = () => {
  // 模拟提交
  lastSubmitTime.value.sogou = new Date().toLocaleString('zh-CN')
  message.success('已提交到搜狗')
}

// 提交到神马搜索
const submitToShenma = () => {
  // 模拟提交
  lastSubmitTime.value.shenma = new Date().toLocaleString('zh-CN')
  message.success('已提交到神马搜索')
}

// 提交到360搜索
const submitTo360 = () => {
  // 模拟提交
  lastSubmitTime.value.so360 = new Date().toLocaleString('zh-CN')
  message.success('已提交到360搜索')
}

// 批量提交
const submitToAll = () => {
  // 模拟批量提交
  lastSubmitTime.value.baidu = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.google = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.bing = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.sogou = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.shenma = new Date().toLocaleString('zh-CN')
  lastSubmitTime.value.so360 = new Date().toLocaleString('zh-CN')
  message.success('已批量提交到所有搜索引擎')
}

// 加载外链列表
const loadLinkList = () => {
  // 模拟加载数据
  linkLoading.value = true
  setTimeout(() => {
    linkLoading.value = false
  }, 1000)
}

// 添加新外链
const addNewLink = () => {
  message.info('添加外链功能开发中')
}

// 编辑外链
const editLink = (row: any) => {
  message.info(`编辑外链: ${row.title}`)
}

// 删除外链
const deleteLink = (row: any) => {
  message.warning(`删除外链: ${row.title}`)
}

// Sitemap管理相关
const sitemapConfig = ref({
  autoGenerate: false,
  lastGenerate: '',
  lastUpdate: ''
})

const sitemapStats = ref({
  total_resources: 0,
  total_pages: 0,
  last_generate: ''
})

const configLoading = ref(false)
const isGenerating = ref(false)
const generateStatus = ref('')

// 获取Sitemap配置
const loadSitemapConfig = async () => {
  try {
    const sitemapApi = useSitemapApi()
    const response = await sitemapApi.getSitemapConfig()
    if (response) {
      sitemapConfig.value = response
    }
  } catch (error) {
    message.error('获取Sitemap配置失败')
  }
}

// 更新Sitemap配置
const updateSitemapConfig = async (value: boolean) => {
  configLoading.value = true
  try {
    const sitemapApi = useSitemapApi()
    await sitemapApi.updateSitemapConfig({
      autoGenerate: value,
      lastGenerate: sitemapConfig.value.lastGenerate,
      lastUpdate: new Date().toISOString()
    })
    message.success(value ? '自动生成功能已开启' : '自动生成功能已关闭')
  } catch (error) {
    message.error('更新配置失败')
    // 恢复之前的值
    sitemapConfig.value.autoGenerate = !value
  } finally {
    configLoading.value = false
  }
}

// 生成Sitemap
const generateSitemap = async () => {
  isGenerating.value = true
  generateStatus.value = '正在启动生成任务...'

  try {
    // 使用已经加载的系统配置
    const siteUrl = systemConfig.value?.site_url || ''
    if (!siteUrl) {
      message.warning('请先在站点配置中设置站点URL，然后再生成Sitemap')
      isGenerating.value = false
      return
    }

    const sitemapApi = useSitemapApi()
    const response = await sitemapApi.generateSitemap({ site_url: siteUrl })

    if (response) {
      generateStatus.value = response.message || '生成任务已启动'
      message.success(`Sitemap生成任务已启动，使用站点URL: ${siteUrl}`)
      // 更新统计信息
      sitemapStats.value.total_resources = response.total_resources || 0
      sitemapStats.value.total_pages = response.total_pages || 0
    }
  } catch (error: any) {
    generateStatus.value = '生成失败: ' + (error.message || '未知错误')
    message.error('Sitemap生成失败')
  } finally {
    isGenerating.value = false
  }
}

// 刷新Sitemap状态
const refreshSitemapStatus = async () => {
  try {
    const sitemapApi = useSitemapApi()
    const response = await sitemapApi.getSitemapStatus()
    if (response) {
      sitemapStats.value = response
      generateStatus.value = '状态已刷新'
    }
  } catch (error) {
    message.error('刷新状态失败')
  }
}

// 查看Sitemap
const viewSitemap = () => {
  window.open('/sitemap.xml', '_blank')
}

// 获取系统配置
const systemConfig = ref<any>(null)

// Google索引配置
const googleIndexConfig = ref({
  enabled: false,
  siteURL: '',
  credentialsFile: '',
  checkInterval: 60,
  batchSize: 100,
  concurrency: 5
})

// 凭据验证相关
const validatingCredentials = ref(false)
const credentialsStatus = ref<string | null>(null) // 'valid', 'invalid', 'verifying', null
const credentialsStatusMessage = ref('')
const credentialsFileInput = ref<HTMLInputElement | null>(null)

// 申请凭据抽屉显示状态
const showCredentialsGuide = ref(false)

// 所有权验证相关
const showVerificationModal = ref(false)
const verificationCode = ref('')
const verificationCodeSaving = ref(false)

// 获取已保存的验证代码
const loadVerificationCode = async () => {
  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexConfigByKey('google_site_verification')
    if (response?.data?.value) {
      verificationCode.value = response.data.value
    }
  } catch (error) {
    // 如果配置不存在，使用空值
    verificationCode.value = ''
  }
}

// Google索引统计
const googleIndexStats = ref({
  totalURLs: 0,
  indexedURLs: 0,
  errorURLs: 0,
  totalTasks: 0,
  runningTasks: 0,
  completedTasks: 0,
  failedTasks: 0
})

// Google索引任务列表
const googleIndexTasks = ref([])
const tasksLoading = ref(false)

// 分页配置
const googleIndexPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  itemCount: 0,
  onChange: (page: number) => {
    googleIndexPagination.value.page = page
    loadGoogleIndexTasks()
  },
  onUpdatePageSize: (pageSize: number) => {
    googleIndexPagination.value.pageSize = pageSize
    googleIndexPagination.value.page = 1
    loadGoogleIndexTasks()
  }
})

// 模态框状态
const urlCheckModal = ref({
  show: false,
  urls: ''
})

// 加载状态
const manualCheckLoading = ref(false)
const submitSitemapLoading = ref(false)

// Google索引任务表格列
const googleIndexTaskColumns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '标题',
    key: 'title',
    width: 200
  },
  {
    title: '类型',
    key: 'type',
    width: 120,
    render: (row: any) => {
      const typeMap = {
        status_check: { text: '状态检查', class: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' },
        sitemap_submit: { text: '网站地图', class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
        url_indexing: { text: 'URL索引', class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' }
      }
      const type = typeMap[row.type as keyof typeof typeMap] || { text: row.type, class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' }
      return h('span', {
        class: `px-2 py-1 text-xs font-medium rounded ${type.class}`
      }, type.text)
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: any) => {
      const statusMap = {
        pending: { text: '待处理', class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' },
        running: { text: '运行中', class: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' },
        completed: { text: '完成', class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
        failed: { text: '失败', class: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' }
      }
      const status = statusMap[row.status as keyof typeof statusMap] || { text: row.status, class: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' }
      return h('span', {
        class: `px-2 py-1 text-xs font-medium rounded ${status.class}`
      }, status.text)
    }
  },
  {
    title: '总项目',
    key: 'totalItems',
    width: 100
  },
  {
    title: '成功/失败',
    key: 'progress',
    width: 120,
    render: (row: any) => {
      return h('span', `${row.successItems} / ${row.failedItems}`)
    }
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 150,
    render: (row: any) => {
      return row.createdAt ? new Date(row.createdAt).toLocaleString('zh-CN') : 'N/A'
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: any) => {
      return h('div', { class: 'space-x-2' }, [
        h('button', {
          class: 'text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 text-sm',
          onClick: () => viewTaskItems(row.id)
        }, '详情'),
        h('button', {
          class: 'text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-300 text-sm',
          disabled: row.status !== 'pending' && row.status !== 'running',
          onClick: () => startTask(row.id)
        }, '启动')
      ].filter(btn => !btn.props?.disabled))
    }
  }
]

// 加载系统配置
const loadSystemConfig = async () => {
  try {
    const { useSystemConfigStore } = await import('~/stores/systemConfig')
    const systemConfigStore = useSystemConfigStore()
    await systemConfigStore.initConfig(true, true)
    systemConfig.value = systemConfigStore.config
  } catch (error) {
    console.error('获取系统配置失败:', error)
  }
}

// 加载Google索引配置
const loadGoogleIndexConfig = async () => {
  try {
    const api = useApi()
    const configs = await api.googleIndexApi.getGoogleIndexConfig()
    if (configs) {
      // 查找general配置
      const generalConfig = configs.find((c: any) => c.group === 'general')
      const authConfig = configs.find((c: any) => c.group === 'auth')

      if (generalConfig) {
        const configData = JSON.parse(generalConfig.value)
        googleIndexConfig.value.enabled = configData.enabled || false
        googleIndexConfig.value.siteURL = configData.siteURL || ''
        googleIndexConfig.value.checkInterval = configData.checkInterval || 60
        googleIndexConfig.value.batchSize = configData.batchSize || 100
        googleIndexConfig.value.concurrency = configData.concurrency || 5
      }

      if (authConfig) {
        const authData = JSON.parse(authConfig.value)
        googleIndexConfig.value.credentialsFile = authData.credentialsFile || ''
      }
    }
  } catch (error) {
    console.error('获取Google索引配置失败:', error)
  }
}


// 选择凭据文件
const selectCredentialsFile = () => {
  if (credentialsFileInput.value) {
    credentialsFileInput.value.click()
  }
}

// 处理凭据文件选择
const handleCredentialsFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) {
    return
  }

  // 验证文件类型
  if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
    message.error('请上传JSON格式的凭据文件')
    return
  }

  // 验证文件大小 (2MB限制)
  if (file.size > 2 * 1024 * 1024) {
    message.error('文件大小不能超过2MB')
    return
  }

  // 上传文件
  try {
    const api = useApi()
    const response = await api.googleIndexApi.uploadCredentials(file)

    if (response?.filePath) {
      googleIndexConfig.value.credentialsFile = response.filePath
      message.success('凭据文件上传成功，请验证凭据')

      // 清空文件输入以允许重新选择相同文件
      if (credentialsFileInput.value) {
        credentialsFileInput.value.value = ''
      }
    } else {
      message.error('上传响应格式错误')
    }
  } catch (error: any) {
    console.error('凭据文件上传失败:', error)
    message.error('凭据文件上传失败: ' + (error?.message || '未知错误'))
  }
}

// 验证凭据
const validateCredentials = async () => {
  if (!googleIndexConfig.value.credentialsFile) {
    message.warning('请先上传凭据文件')
    return
  }

  validatingCredentials.value = true
  credentialsStatus.value = 'verifying'
  credentialsStatusMessage.value = '正在验证凭据...'

  try {
    const api = useApi()
    const response = await api.googleIndexApi.validateCredentials({
      credentialsFile: googleIndexConfig.value.credentialsFile
    })

    if (response?.valid) {
      credentialsStatus.value = 'valid'
      credentialsStatusMessage.value = '凭据验证成功！凭据有效，可正常使用Google索引功能。'
      message.success('凭据验证成功')
    } else {
      credentialsStatus.value = 'invalid'
      credentialsStatusMessage.value = '凭据验证失败：' + (response?.message || '凭据无效或权限不足')
      message.error('凭据验证失败')
    }
  } catch (error: any) {
    console.error('凭据验证失败:', error)
    credentialsStatus.value = 'invalid'
    credentialsStatusMessage.value = '凭据验证失败：' + (error?.message || '网络错误或服务器异常')
    message.error('凭据验证失败: ' + (error?.message || '网络错误'))
  } finally {
    validatingCredentials.value = false
  }
}

// 检查凭据是否有效
const checkCredentialsValid = async (): Promise<boolean> => {
  if (!googleIndexConfig.value.credentialsFile) {
    return false
  }

  try {
    const api = useApi()
    const response = await api.googleIndexApi.validateCredentials({
      credentialsFile: googleIndexConfig.value.credentialsFile
    })

    return response?.valid || false
  } catch (error) {
    console.error('检查凭据有效性时出错:', error)
    return false
  }
}

// 更新Google索引配置
const updateGoogleIndexConfig = async () => {
  // 如果启用功能，先检查凭据是否有效
  if (googleIndexConfig.value.enabled && googleIndexConfig.value.credentialsFile) {
    const isValid = await checkCredentialsValid()
    if (!isValid) {
      message.warning('凭据未通过验证，无法启用Google索引功能')
      googleIndexConfig.value.enabled = false
      return
    }
  }

  configLoading.value = true
  try {
    const api = useApi()

    // 更新general配置
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'general',
      key: 'general',
      value: JSON.stringify({
        enabled: googleIndexConfig.value.enabled,
        siteURL: systemConfig.value?.site_url || googleIndexConfig.value.siteURL, // 使用系统配置的站点URL
        checkInterval: googleIndexConfig.value.checkInterval,
        batchSize: googleIndexConfig.value.batchSize,
        concurrency: googleIndexConfig.value.concurrency || 5
      })
    })

    // 更新auth配置
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'auth',
      key: 'credentials_file',
      value: JSON.stringify({
        credentialsFile: googleIndexConfig.value.credentialsFile
      })
    })

    // 更新schedule配置
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'schedule',
      key: 'schedule',
      value: JSON.stringify({
        retryAttempts: 3,
        retryDelay: 2,
        checkSchedule: '@daily'
      })
    })

    // 更新sitemap配置
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'sitemap',
      key: 'sitemap',
      value: JSON.stringify({
        sitemapPath: '/sitemap.xml'
      })
    })

    message.success('Google索引配置已更新')
  } catch (error) {
    console.error('更新Google索引配置失败:', error)
    message.error('更新配置失败')
  } finally {
    configLoading.value = false
  }
}

// 刷新Google索引状态
const refreshGoogleIndexStatus = async () => {
  try {
    const api = useApi()

    // 加载统计信息
    const statsResponse = await api.googleIndexApi.getGoogleIndexStatus()
    if (statsResponse) {
      googleIndexStats.value = statsResponse
    }

    // 加载任务列表
    await loadGoogleIndexTasks()
  } catch (error) {
    console.error('刷新Google索引状态失败:', error)
    message.error('刷新状态失败')
  }
}

// 加载Google索引任务列表
const loadGoogleIndexTasks = async () => {
  tasksLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexTasks({
      page: googleIndexPagination.value.page,
      pageSize: googleIndexPagination.value.pageSize
    })
    if (response) {
      googleIndexTasks.value = response.tasks || []
      googleIndexPagination.value.itemCount = response.total || 0
    }
  } catch (error) {
    console.error('加载Google索引任务列表失败:', error)
    message.error('加载任务列表失败')
  } finally {
    tasksLoading.value = false
  }
}

// 手动检查URL
const manualCheckURLs = () => {
  urlCheckModal.value.show = true
  urlCheckModal.value.urls = ''
}

// 确认手动检查URL
const confirmManualCheckURLs = async () => {
  const urls = urlCheckModal.value.urls.split('\n').filter(url => url.trim() !== '')
  if (urls.length === 0) {
    message.warning('请至少输入一个URL')
    return
  }

  manualCheckLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.createGoogleIndexTask({
      title: `手动URL检查任务 - ${new Date().toLocaleString('zh-CN')}`,
      type: 'status_check',
      description: `手动检查 ${urls.length} 个URL的索引状态`,
      URLs: urls  // 注意：后端API中字段名为URLs而不是urls
    })
    if (response) {
      message.success('URL检查任务已创建')
      urlCheckModal.value.show = false
      await refreshGoogleIndexStatus()
    }
  } catch (error) {
    console.error('手动检查URL失败:', error)
    message.error('手动检查URL失败')
  } finally {
    manualCheckLoading.value = false
  }
}

// 提交网站地图
const submitSitemap = async () => {
  const siteUrl = systemConfig.value?.site_url || ''
  if (!siteUrl) {
    message.warning('请先在站点配置中设置站点URL')
    return
  }

  const sitemapUrl = siteUrl + '/sitemap.xml'

  submitSitemapLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.createGoogleIndexTask({
      title: `网站地图提交任务 - ${new Date().toLocaleString('zh-CN')}`,
      type: 'sitemap_submit',
      description: `提交网站地图: ${sitemapUrl}`,
      SitemapURL: sitemapUrl  // 注意：后端API中字段名为SitemapURL而不是sitemapURL
    })
    if (response) {
      message.success('网站地图提交任务已创建')
      await refreshGoogleIndexStatus()
    }
  } catch (error) {
    console.error('提交网站地图失败:', error)
    message.error('提交网站地图失败')
  } finally {
    submitSitemapLoading.value = false
  }
}

// 查看任务详情
const viewTaskItems = async (taskId: number) => {
  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexTaskItems(taskId)
    if (response) {
      // 在新窗口中打开任务详情
      const items = response.items || []
      const content = `任务 ${taskId} 详情:\n\n` +
        items.map((item: any) =>
          `URL: ${item.URL}\n状态: ${item.status}\n索引状态: ${item.indexStatus}\n错误信息: ${item.errorMessage}\n---\n`
        ).join('')
      alert(content)
    }
  } catch (error) {
    console.error('获取任务详情失败:', error)
    message.error('获取任务详情失败')
  }
}

// 监听验证弹窗显示状态
watch(showVerificationModal, (show) => {
  if (show) {
    loadVerificationCode()
  }
})

// 保存验证字符串
const saveVerificationCode = async () => {
  if (!verificationCode.value.trim()) {
    message.warning('请输入验证字符串')
    return
  }

  verificationCodeSaving.value = true
  try {
    const api = useApi()
    // 使用分组配置API保存验证字符串
    await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'verification',
      key: 'google_site_verification',
      value: JSON.stringify({
        code: verificationCode.value.trim()
      })
    })

    message.success('验证字符串保存成功')
    showVerificationModal.value = false
  } catch (error) {
    console.error('保存验证字符串失败:', error)
    message.error('保存验证字符串失败')
  } finally {
    verificationCodeSaving.value = false
  }
}

// 启动任务
const startTask = async (taskId: number) => {
  try {
    const api = useApi()
    const response = await api.googleIndexApi.startGoogleIndexTask(taskId)
    if (response) {
      message.success('任务已启动')
      await loadGoogleIndexTasks()
    }
  } catch (error) {
    console.error('启动任务失败:', error)
    message.error('启动任务失败')
  }
}

// 初始化
onMounted(async () => {
  await loadSystemConfig()
  loadLinkList()
  await loadSitemapConfig()
  await refreshSitemapStatus()
  await loadGoogleIndexConfig()
  await refreshGoogleIndexStatus()
})
</script>

<style scoped>
/* SEO管理页面样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}

.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>