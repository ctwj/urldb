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
          <!-- Sitemap管理 Tab -->
          <n-tab-pane name="sitemap" tab="Sitemap管理">
            <SitemapTab
              :system-config="systemConfig"
              :sitemap-config="sitemapConfig"
              :sitemap-stats="sitemapStats"
              :config-loading="configLoading"
              :is-generating="isGenerating"
              :generate-status="generateStatus"
              @update:sitemap-config="updateSitemapConfig"
              @refresh-status="refreshSitemapStatus"
            />
          </n-tab-pane>

          <!-- Google Index Tab -->
          <n-tab-pane name="google-index" tab="Google Index">
            <GoogleIndexTab
              :system-config="systemConfig"
              :google-index-config="googleIndexConfig"
              :tasks="googleIndexTasks"
              :credentials-status="credentialsStatus"
              :credentials-status-message="credentialsStatusMessage"
              :config-loading="configLoading"
              :manual-check-loading="manualCheckLoading"
              :manual-submit-loading="manualSubmitLoading"
              :submit-sitemap-loading="submitSitemapLoading"
              :tasks-loading="tasksLoading"
              :diagnose-loading="diagnoseLoading"
              :pagination="googleIndexPagination"
              @update:google-index-config="updateGoogleIndexConfig"
              @show-verification="showVerificationModal = true"
              @show-credentials-guide="showCredentialsGuide = true"
              @select-credentials-file="selectCredentialsFile"
              @manual-check-urls="manualCheckURLs"
              @manual-submit-urls="manualSubmitURLs"
              @refresh-status="refreshGoogleIndexStatus"
              @diagnose-permissions="diagnosePermissions"
              @view-task-items="viewTaskItems"
              @start-task="startTask"
            />
          </n-tab-pane>

          <!-- Bing索引 Tab -->
          <n-tab-pane name="bing-index" tab="Bing索引">
            <BingTab
              :system-config="systemConfig"
              :bing-index-config="bingIndexConfig"
              :submit-history="bingSubmitHistory"
              :last-submit-status="bingLastSubmitStatus"
              :last-submit-time="bingLastSubmitTime"
              :config-loading="configLoading"
              :submit-sitemap-loading="bingSubmitSitemapLoading"
              :batch-submit-loading="bingBatchSubmitLoading"
              :status-loading="bingStatusLoading"
              :history-loading="bingHistoryLoading"
              :pagination="bingPagination"
              @update:bing-index-config="updateBingIndexConfig"
              @refresh-status="refreshBingStatus"
            />
          </n-tab-pane>

          <!-- 站点提交 Tab -->
          <n-tab-pane name="site-submit" tab="站点提交">
            <SiteSubmitTab
              :last-submit-time="lastSubmitTime"
              @update:last-submit-time="updateLastSubmitTime"
            />
          </n-tab-pane>

          <!-- 外链建设 Tab -->
          <n-tab-pane name="link-building" tab="外链建设">
            <LinkBuildingTab
              :link-stats="linkStats"
              :link-list="linkList"
              :loading="linkLoading"
              :pagination="linkPagination"
              @add-new-link="addNewLink"
              @edit-link="editLink"
              @delete-link="deleteLink"
              @load-link-list="loadLinkList"
            />
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

  <!-- URL提交模态框 -->
  <n-modal v-model:show="urlSubmitModal.show" preset="card" title="手动提交URL到Google索引" style="max-width: 600px;">
    <div class="space-y-4">
      <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">重要说明</h3>
            <div class="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
              <ul class="list-disc list-inside space-y-1">
                <li>此功能将直接向Google提交URL索引请求</li>
                <li>Google Indexing API有每日配额限制（建议不超过100个URL/天）</li>
                <li>提交成功不代表立即被索引，Google仍会根据页面质量决定</li>
                <li>请确保URL可正常访问且内容符合Google质量指南</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <p class="text-gray-600 dark:text-gray-400">输入要提交到Google索引的URL，每行一个：</p>
      <n-input
        v-model:value="urlSubmitModal.urls"
        type="textarea"
        :autosize="{ minRows: 4, maxRows: 8 }"
        placeholder="https://yoursite.com/page1&#10;https://yoursite.com/page2"
      />

      <div class="flex justify-between items-center">
        <div class="text-sm text-gray-500">
          <i class="fas fa-info-circle mr-1"></i>
          提交后需要等待Google处理，结果可在任务列表中查看
        </div>
        <div class="flex space-x-2">
          <n-button @click="urlSubmitModal.show = false">取消</n-button>
          <n-button type="primary" @click="confirmManualSubmitURLs" :loading="urlSubmitLoading">
            确认提交
          </n-button>
        </div>
      </div>
    </div>
  </n-modal>

  <!-- 所有权验证模态框 -->
  <n-modal v-model:show="showVerificationModal" preset="card" title="站点所有权验证" style="max-width: 600px;">
    <div class="space-y-6">
      <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="fas fa-info-circle text-blue-500 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800 dark:text-blue-200">DNS方式验证</h3>
            <div class="mt-2 text-sm text-blue-700 dark:text-blue-300">
              <p>推荐使用DNS方式验证站点所有权，这是最安全和可靠的方法：</p>
              <ol class="list-decimal list-inside mt-2 space-y-1">
                <li>登录您的域名注册商或DNS管理平台</li>
                <li>添加一条TXT记录</li>
                <li>在Google Search Console中输入您的验证字符串</li>
                <li>验证DNS TXT记录是否生效</li>
              </ol>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">注意事项</h3>
            <div class="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
              <ul class="list-disc list-inside space-y-1">
                <li>DNS验证比HTML标签更安全，不易被其他网站复制</li>
                <li>验证成功后，Google会自动检测您的站点所有权</li>
                <li>如果您的域名服务商不支持TXT记录，请联系客服寻求帮助</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end pt-2">
        <n-button type="primary" @click="showVerificationModal = false">
          确定
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
        <div class="border-l-4 border-orange-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-3 text-orange-500 mr-2"></i>启用Indexing API
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            <strong class="text-orange-600">重要：</strong>除了Search Console API，还需要启用Indexing API才能提交URL到Google索引。
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>在导航菜单中选择"API和服务" > "库"</li>
            <li>搜索"Indexing API"或"Google Indexing API"</li>
            <li>点击搜索结果中的"Indexing API"</li>
            <li>点击"启用"按钮</li>
            <li class="text-orange-600 font-medium">⚠️ 如果找不到Indexing API，请确保项目已启用Google Search Console API</li>
          </ol>
          <div class="bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded p-3 mt-3">
            <p class="text-sm text-orange-700 dark:text-orange-300">
              <strong>为什么需要两个API？</strong><br>
              • Search Console API：用于检查URL索引状态和获取站点数据<br>
              • Indexing API：用于主动提交URL到Google索引队列
            </p>
          </div>
        </div>

        <!-- 步骤4 -->
        <div class="border-l-4 border-yellow-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-4 text-yellow-500 mr-2"></i>创建服务账号
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

        <!-- 步骤6 -->
        <div class="border-l-4 border-red-500 pl-4">
          <h4 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            <i class="fas fa-number-6 text-red-500 mr-2"></i>验证网站所有权并配置权限
          </h4>
          <p class="text-gray-600 dark:text-gray-400 mb-3">
            在Google Search Console中验证网站并添加服务账号权限。这是最关键的一步！
          </p>
          <ol class="list-decimal list-inside space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>访问 <a href="https://search.google.com/search-console/" target="_blank" class="text-blue-600 hover:text-blue-800 underline">Google Search Console</a></li>
            <li><strong class="text-red-600">如果尚未验证网站所有权：</strong>
              <ul class="list-disc list-inside ml-4 mt-1 space-y-1">
                <li>点击"添加属性"</li>
                <li>选择"网址前缀"（推荐）或"网域"</li>
                <li>输入您的网站URL（如：https://pan.l9.lc）</li>
                <li>选择验证方法（DNS记录、HTML文件上传、HTML标签或Google Analytics）</li>
                <li>按照指示完成验证</li>
              </ul>
            </li>
            <li><strong class="text-green-600">添加服务账号权限：</strong>
              <ul class="list-disc list-inside ml-4 mt-1 space-y-1">
                <li>选择已验证的网站属性</li>
                <li>在左侧菜单中点击"设置" ⚙️</li>
                <li>选择"用户和权限"</li>
                <li>点击右上角的"添加用户"</li>
                <li>输入服务账号邮箱（格式：xxx@xxx.iam.gserviceaccount.com）</li>
                <li><strong class="text-orange-600">权限选择：</strong>
                  <ul class="list-circle list-inside ml-4 mt-1">
                    <li>✅ 推荐："所有者" - 完全访问权限</li>
                    <li>⚠️ 可选："完整" - 大部分功能权限</li>
                    <li>❌ 不推荐："受限" - 功能受限</li>
                  </ul>
                </li>
                <li>点击"添加"完成授权</li>
              </ul>
            </li>
          </ol>
          <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded p-3 mt-3">
            <p class="text-sm text-green-700 dark:text-green-300">
              <strong>✅ 验证权限配置成功：</strong><br>
              • 添加权限后等待5-10分钟生效<br>
              • 使用上方的"权限诊断"按钮验证配置<br>
              • 如果诊断显示"可访问站点数: 1"，说明配置成功
            </p>
          </div>
        </div>

        <!-- 故障排除 -->
        <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <h5 class="font-semibold text-blue-800 dark:text-blue-200 mb-3">
            <i class="fas fa-tools mr-2"></i>故障排除
          </h5>
          <div class="space-y-3 text-sm text-blue-700 dark:text-blue-300">
            <div>
              <strong class="text-blue-600">❌ 诊断显示"可访问站点数: 0"</strong>
              <ul class="list-disc list-inside ml-4 mt-1 space-y-1">
                <li>确认服务账号邮箱输入正确</li>
                <li>确认已授予"所有者"或"完整"权限</li>
                <li>等待权限生效（可能需要10-15分钟）</li>
                <li>检查网站所有权验证是否有效</li>
              </ul>
            </div>
            <div>
              <strong class="text-blue-600">❌ 提交URL失败，权限错误</strong>
              <ul class="list-disc list-inside ml-4 mt-1 space-y-1">
                <li>确认已启用Indexing API（步骤3）</li>
                <li>确认网站URL格式正确（建议使用https://example.com/）</li>
                <li>检查API配额是否超限</li>
              </ul>
            </div>
            <div>
              <strong class="text-green-600">✅ 使用权限诊断工具</strong>
              <p>点击上方的"权限诊断"按钮，系统会自动检查所有配置并提供详细建议。</p>
            </div>
          </div>
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
            <li>• Indexing API有严格的速率限制，不要频繁提交</li>
            <li>• 权限更改后需要等待几分钟才能生效</li>
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

  <!-- 任务详情模态框 -->
  <n-modal v-model:show="taskDetailModal.show" preset="card" title="任务详情" style="max-width: 900px;">
    <div class="space-y-4">
      <!-- 任务头部信息 -->
      <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold">任务 #{{ taskDetailModal.taskId }}</h3>
          <n-tag :type="taskDetailModal.items.length > 0 ? 'info' : 'default'">
            {{ taskDetailModal.items.length }} 个任务项
          </n-tag>
        </div>
      </div>

      <!-- 任务项列表 -->
      <div v-if="taskDetailModal.loading" class="flex justify-center py-8">
        <n-spin size="large" />
      </div>

      <div v-else-if="taskDetailModal.items.length === 0" class="text-center py-8 text-gray-500">
        暂无任务项
      </div>

      <div v-else class="space-y-4 max-h-96 overflow-y-auto">
        <div
          v-for="(item, index) in taskDetailModal.items"
          :key="item.id"
          class="border rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
          <!-- 任务项头部 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-2">
              <span class="text-sm font-medium text-gray-500">#{{ index + 1 }}</span>
              <n-tag
                :type="item.status === 'success' ? 'success' : item.status === 'failed' ? 'error' : 'default'"
                size="small"
              >
                {{ item.status }}
              </n-tag>
            </div>
            <div class="text-sm text-gray-500">
              {{ item.completedAt ? new Date(item.completedAt).toLocaleString('zh-CN') : '未完成' }}
            </div>
          </div>

          <!-- URL信息 -->
          <div class="mb-3">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">URL:</div>
            <div class="p-2 bg-gray-100 dark:bg-gray-700 rounded text-sm break-all">
              {{ item.URL || 'N/A' }}
            </div>
          </div>

          <!-- 详细状态信息 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <!-- 索引状态 -->
            <div>
              <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">索引状态:</div>
              <n-tag
                :type="formatIndexStatus(item.indexStatus).color"
                size="small"
                v-if="item.indexStatus"
              >
                {{ formatIndexStatus(item.indexStatus).text }}
              </n-tag>
              <span v-else class="text-gray-500">未知</span>
            </div>

            <!-- 移动友好 -->
            <div>
              <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">移动友好:</div>
              <n-tag
                :type="item.mobileFriendly ? 'success' : 'error'"
                size="small"
              >
                {{ item.mobileFriendly ? '是' : '否' }}
              </n-tag>
            </div>

            <!-- HTTP状态码 -->
            <div>
              <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">HTTP状态码:</div>
              <n-tag
                :type="item.statusCode >= 200 && item.statusCode < 300 ? 'success' : 'error'"
                size="small"
              >
                {{ item.statusCode || 'N/A' }}
              </n-tag>
            </div>

            <!-- 最后抓取时间 -->
            <div>
              <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">最后抓取:</div>
              <div class="text-sm">
                {{ item.lastCrawled ? new Date(item.lastCrawled).toLocaleString('zh-CN') : '从未抓取' }}
              </div>
            </div>
          </div>

          <!-- 错误信息 -->
          <div v-if="item.errorMessage" class="mt-3">
            <div class="text-sm font-medium text-red-600 dark:text-red-400 mb-1">错误信息:</div>
            <div class="p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-sm text-red-700 dark:text-red-300">
              {{ item.errorMessage }}
            </div>
          </div>

          <!-- 检查结果详情 -->
          <div v-if="item.inspectResult" class="mt-3">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">检查结果详情:</div>
            <div class="p-2 bg-gray-50 dark:bg-gray-800 rounded text-xs font-mono">
              <pre>{{ JSON.stringify(JSON.parse(item.inspectResult), null, 2) }}</pre>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-end space-x-2 pt-4 border-t">
        <n-button @click="taskDetailModal.show = false">关闭</n-button>
        <n-button
          v-if="taskDetailModal.items.some(item => item.status === 'failed')"
          type="primary"
          @click="retryFailedItems(taskDetailModal.taskId)"
        >
          重试失败项
        </n-button>
      </div>
    </div>
  </n-modal>

  <!-- 隐藏的文件输入 -->
  <input
    type="file"
    ref="credentialsFileInput"
    accept=".json"
    @change="handleCredentialsFileSelect"
    style="display: none;"
  />
</template>

<script setup lang="ts">
import AdminPageLayout from '~/components/AdminPageLayout.vue'
import GoogleIndexTab from '~/components/Admin/GoogleIndexTab.vue'
import SitemapTab from '~/components/Admin/SitemapTab.vue'
import SiteSubmitTab from '~/components/Admin/SiteSubmitTab.vue'
import BingTab from '~/components/Admin/BingTab.vue'
import LinkBuildingTab from '~/components/Admin/LinkBuildingTab.vue'

// SEO管理页面
definePageMeta({
  layout: 'admin'
})

import { useMessage, useDialog } from 'naive-ui'
import { useApi } from '~/composables/useApi'
import { ref, onMounted, watch, h } from 'vue'

// 获取消息组件
const message = useMessage()
const dialog = useDialog()

// 当前激活的Tab - 默认显示 Sitemap管理
const activeTab = ref('sitemap')

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

// Bing索引配置
const bingIndexConfig = ref({
  enabled: false,
  submitInterval: 60,
  batchSize: 5,
  retryCount: 3
})

// 凭据验证相关
const credentialsStatus = ref<string | null>(null)
const credentialsStatusMessage = ref('')
const credentialsFileInput = ref<HTMLInputElement | null>(null)

// 申请凭据抽屉显示状态
const showCredentialsGuide = ref(false)

// 所有权验证相关
const showVerificationModal = ref(false)


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

// 任务详情模态框状态
const taskDetailModal = ref({
  show: false,
  taskId: 0,
  items: [],
  loading: false
})

// URL提交模态框状态
const urlSubmitModal = ref({
  show: false,
  urls: ''
})

// URL提交加载状态
const urlSubmitLoading = ref(false)

// 加载状态
const configLoading = ref(false)
const manualCheckLoading = ref(false)
const manualSubmitLoading = ref(false)
const diagnoseLoading = ref(false)
const submitSitemapLoading = ref(false)

// Bing相关状态
const bingSubmitSitemapLoading = ref(false)
const bingBatchSubmitLoading = ref(false)
const bingStatusLoading = ref(false)
const bingHistoryLoading = ref(false)
const bingLastSubmitStatus = ref('')
const bingLastSubmitTime = ref('')
const bingSubmitHistory = ref([])

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

const isGenerating = ref(false)
const generateStatus = ref('')

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

// Bing分页配置
const bingPagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    bingPagination.value.page = page
    loadBingSubmitHistory()
  },
  onUpdatePageSize: (pageSize: number) => {
    bingPagination.value.pageSize = pageSize
    bingPagination.value.page = 1
    loadBingSubmitHistory()
  }
})

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
    console.log('开始加载 Google 索引配置...')
    const api = useApi()
    const configs = await api.googleIndexApi.getGoogleIndexConfig()
    console.log('获取到的配置:', configs)
    if (configs) {
      // 查找general配置
      const generalConfig = configs.find((c: any) => c.group === 'general')
      const authConfig = configs.find((c: any) => c.group === 'auth')
      console.log('找到的配置 - general:', generalConfig, 'auth:', authConfig)

      let newConfig = { ...googleIndexConfig.value }

      if (generalConfig) {
        const configData = JSON.parse(generalConfig.value)
        newConfig.enabled = configData.enabled || false
        newConfig.siteURL = configData.siteURL || ''
        newConfig.checkInterval = configData.checkInterval || 60
        newConfig.batchSize = configData.batchSize || 100
        newConfig.concurrency = configData.concurrency || 5
      }

      if (authConfig) {
        console.log('解析 auth 配置:', authConfig.value)
        const authData = JSON.parse(authConfig.value)
        console.log('解析后的 authData:', authData)
        newConfig.credentialsFile = authData.credentialsFile || authData.credentials_file || ''
        console.log('设置凭据文件路径:', newConfig.credentialsFile)
      }

      // 强制触发响应式更新
      googleIndexConfig.value = newConfig
      console.log('最终配置:', googleIndexConfig.value)
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

    // 检查API是否成功（success字段为true）且包含有效的文件路径
    if (response?.success === true && response?.file_path) {
      console.log('上传成功，文件路径:', response.file_path)
      // 统一路径格式为 Unix 格式
      const normalizedPath = response.file_path.replace(/\\/g, '/')
      console.log('标准化路径:', normalizedPath)
      // 强制触发响应式更新
      googleIndexConfig.value = {
        ...googleIndexConfig.value,
        credentialsFile: normalizedPath
      }
      console.log('更新后的 googleIndexConfig:', googleIndexConfig.value)
      message.success(response.message || '凭据文件上传成功，请验证凭据')

      // 清空文件输入以允许重新选择相同文件
      if (credentialsFileInput.value) {
        credentialsFileInput.value.value = ''
      }

      // 上传成功后立即更新后端配置并重新加载配置
      try {
        const configData = {
          group: 'auth',
          key: 'credentials_file',
          value: JSON.stringify({
            credentials_file: googleIndexConfig.value.credentialsFile.replace(/\\/g, '/')
          })
        }
        console.log('更新后端配置，发送数据:', JSON.stringify(configData, null, 2))

        const updateResponse = await api.googleIndexApi.updateGoogleIndexGroupConfig(configData)
        console.log('后端配置更新响应:', updateResponse)

        // 等待一下确保后端处理完成
        await new Promise(resolve => setTimeout(resolve, 500))

        // 重新加载配置以确保UI状态与后端同步
        console.log('重新加载配置...')
        await loadGoogleIndexConfig()
        console.log('配置重新加载完成')
      } catch (configError) {
        console.error('更新配置失败:', configError)
        message.error('配置更新失败，但文件已上传')

        // 即使配置更新失败，也尝试刷新状态
        setTimeout(async () => {
          console.log('延迟重新加载配置...')
          await loadGoogleIndexConfig()
        }, 1000)
      }
    } else {
      // 如果API调用成功但返回的数据有问题，或者API调用失败
      message.error(response?.message || '上传响应格式错误')
    }
  } catch (error: any) {
    console.error('凭据文件上传失败:', error)
    message.error('凭据文件上传失败: ' + (error?.message || '未知错误'))
  }
}

// 更新Google索引配置
const updateGoogleIndexConfig = async (value: boolean) => {
  configLoading.value = true
  try {
    const api = useApi()

    // 先更新本地状态
    googleIndexConfig.value.enabled = value

    // 更新general配置
    const response = await api.googleIndexApi.updateGoogleIndexGroupConfig({
      group: 'general',
      key: 'general',
      value: JSON.stringify({
        enabled: value,
        siteURL: systemConfig.value?.site_url || googleIndexConfig.value.siteURL,
        checkInterval: googleIndexConfig.value.checkInterval,
        batchSize: googleIndexConfig.value.batchSize,
        concurrency: googleIndexConfig.value.concurrency || 5
      })
    })

    message.success('Google索引配置已更新')
    
    // 延迟重新加载配置以验证后端状态（在后台进行，不阻塞UI）
    setTimeout(async () => {
      try {
        await loadGoogleIndexConfig()
      } catch (error) {
        console.error('重新加载配置失败:', error)
      }
    }, 1000)
  } catch (error) {
    console.error('更新Google索引配置失败:', error)
    message.error('更新配置失败')
    // 失败时恢复原状态
    googleIndexConfig.value.enabled = !value
  } finally {
    configLoading.value = false
  }
}

// 刷新Google索引状态
const refreshGoogleIndexStatus = async () => {
  try {
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

// 手动提交URL
const manualSubmitURLs = () => {
  urlSubmitModal.value.show = true
  urlSubmitModal.value.urls = ''
}

// 诊断Google API权限
const diagnosePermissions = async () => {
  diagnoseLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.diagnosePermissions({})

    if (response?.diagnosis) {
      const diagnosis = response.diagnosis

      // 创建诊断结果对话框
      dialog.create({
        title: 'Google API 权限诊断结果',
        style: {
          width: '800px',
          maxWidth: '90vw'
        },
        content: () => h('div', { class: 'space-y-6 p-4' }, [
          // 凭据信息
          h('div', { class: 'bg-gray-50 dark:bg-gray-700/50 p-4 rounded-lg' }, [
            h('h4', { class: 'font-semibold text-gray-900 dark:text-white mb-3' }, '📋 凭据信息'),
            h('div', { class: 'grid grid-cols-2 gap-4 text-sm' }, [
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, '服务账号: '),
                h('span', { class: 'font-mono text-gray-900 dark:text-white' }, diagnosis.credentials.service_account || 'N/A')
              ]),
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, '项目ID: '),
                h('span', { class: 'font-mono text-gray-900 dark:text-white' }, diagnosis.credentials.project_id || 'N/A')
              ]),
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, '凭据类型: '),
                h('span', { class: 'text-gray-900 dark:text-white' }, diagnosis.credentials.type || 'N/A')
              ]),
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, '文件状态: '),
                h('span', { class: diagnosis.credentials.file_exists ? 'text-green-600' : 'text-red-600' },
                   diagnosis.credentials.file_exists ? '✅ 存在' : '❌ 不存在')
              ])
            ])
          ]),

          // API访问状态
          h('div', { class: 'bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg' }, [
            h('h4', { class: 'font-semibold text-gray-900 dark:text-white mb-3' }, '🔌 API访问状态'),
            h('div', { class: 'grid grid-cols-2 gap-4 text-sm' }, [
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, '可访问站点数: '),
                h('span', { class: `font-bold ${diagnosis.api_access.sites_count > 0 ? 'text-green-600' : 'text-red-600'}` },
                   diagnosis.api_access.sites_count)
              ]),
              h('div', [
                h('span', { class: 'text-gray-600 dark:text-gray-400' }, 'Search Console: '),
                h('span', { class: diagnosis.api_access.search_console_enabled ? 'text-green-600' : 'text-red-600' },
                   diagnosis.api_access.search_console_enabled ? '✅ 已启用' : '❌ 未启用')
              ])
            ])
          ]),

          // 站点测试结果
          h('div', { class: 'bg-yellow-50 dark:bg-yellow-900/20 p-4 rounded-lg' }, [
            h('h4', { class: 'font-semibold text-gray-900 dark:text-white mb-3' }, '🔍 站点访问测试'),
            ...diagnosis.site_tests.map((test: any) =>
              h('div', { class: 'mb-3 p-3 border border-yellow-200 dark:border-yellow-800 rounded' }, [
                h('div', { class: 'font-mono text-sm mb-2' }, test.site_format),
                h('div', { class: 'grid grid-cols-2 gap-4 text-sm' }, [
                  h('div', [
                    h('span', { class: 'text-gray-600' }, '站点访问: '),
                    h('span', { class: test.site_access ? 'text-green-600' : 'text-red-600' },
                       test.site_access ? '✅ 成功' : '❌ 失败')
                  ]),
                  h('div', [
                    h('span', { class: 'text-gray-600' }, 'URL检查: '),
                    h('span', { class: test.url_inspect ? 'text-green-600' : 'text-red-600' },
                       test.url_inspect ? '✅ 成功' : '❌ 失败')
                  ])
                ]),
                test.site_error && h('div', { class: 'text-red-600 text-xs mt-1' }, test.site_error),
                test.inspect_error && h('div', { class: 'text-red-600 text-xs mt-1' }, test.inspect_error)
              ])
            )
          ]),

          // 建议和解决方案
          h('div', { class: 'bg-orange-50 dark:bg-orange-900/20 p-4 rounded-lg' }, [
            h('h4', { class: 'font-semibold text-gray-900 dark:text-white mb-3' }, '💡 建议和解决方案'),
            ...diagnosis.recommendations.map((rec: string) =>
              h('div', { class: 'text-sm text-gray-700 dark:text-gray-300 mb-2 leading-relaxed' }, rec)
            )
          ])
        ]),
        positiveText: '关闭',
        onPositiveClick: () => {
          dialog.destroyAll()
        }
      })

      message.success('权限诊断完成')
    }
  } catch (error: any) {
    console.error('权限诊断失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '权限诊断失败'
    message.error('权限诊断失败: ' + errorMsg)
  } finally {
    diagnoseLoading.value = false
  }
}

// 确认手动提交URL
const confirmManualSubmitURLs = async () => {
  const urls = urlSubmitModal.value.urls.split('\n').filter(url => url.trim() !== '')
  if (urls.length === 0) {
    message.warning('请至少输入一个URL')
    return
  }

  urlSubmitLoading.value = true
  try {
    const api = useApi()
    const response = await api.googleIndexApi.submitURLsToIndex({
      urls: urls
    })
    if (response) {
      message.success('URL提交任务已创建，正在后台处理')
      urlSubmitModal.value.show = false
      await refreshGoogleIndexStatus()
    }
  } catch (error: any) {
    console.error('手动提交URL失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '手动提交URL失败'
    message.error('手动提交URL失败: ' + errorMsg)
  } finally {
    urlSubmitLoading.value = false
  }
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
      URLs: urls
    })
    if (response) {
      message.success('URL检查任务已创建')
      urlCheckModal.value.show = false
      await refreshGoogleIndexStatus()
    }
  } catch (error: any) {
    console.error('手动检查URL失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '手动检查URL失败'
    message.error('手动检查URL失败: ' + errorMsg)
  } finally {
    manualCheckLoading.value = false
  }
}

// 查看任务详情
const viewTaskItems = async (taskId: number) => {
  taskDetailModal.value.show = true
  taskDetailModal.value.taskId = taskId
  taskDetailModal.value.loading = true
  taskDetailModal.value.items = []

  try {
    const api = useApi()
    const response = await api.googleIndexApi.getGoogleIndexTaskItems(taskId)
    if (response) {
      taskDetailModal.value.items = response.items || []
    }
  } catch (error: any) {
    console.error('获取任务详情失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '获取任务详情失败'
    message.error('获取任务详情失败: ' + errorMsg)
  } finally {
    taskDetailModal.value.loading = false
  }
}

// 格式化索引状态
const formatIndexStatus = (status: string) => {
  const statusMap: Record<string, { text: string; color: string }> = {
    'SUBMITTED': { text: '已提交', color: 'blue' },
    'INDEXING_ALLOWED': { text: '允许索引', color: 'green' },
    'INDEXING_BLOCKED': { text: '索引被阻止', color: 'red' },
    'BLOCKED_BY_ROBOTS_TXT': { text: '被robots.txt阻止', color: 'orange' },
    'PAGE_WITH_REDIRECT': { text: '页面重定向', color: 'orange' },
    'NOT_FOUND': { text: '页面未找到', color: 'red' }
  }

  const statusInfo = statusMap[status] || { text: status || '未知', color: 'gray' }
  return statusInfo
}

// 获取状态颜色类
const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    'success': 'text-green-600',
    'failed': 'text-red-600',
    'pending': 'text-gray-600',
    'processing': 'text-blue-600'
  }
  return colorMap[status] || 'text-gray-600'
}

// 重试失败的任务项
const retryFailedItems = async (taskId: number) => {
  try {
    const api = useApi()
    // 这里可以调用重试API，暂时重新启动任务
    await api.googleIndexApi.startGoogleIndexTask(taskId)
    message.success('已重新启动任务')
    taskDetailModal.value.show = false
    await loadGoogleIndexTasks()
  } catch (error: any) {
    console.error('重试失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '重试失败'
    message.error('重试失败: ' + errorMsg)
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
  } catch (error: any) {
    console.error('启动任务失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '启动任务失败'
    message.error('启动任务失败: ' + errorMsg)
  }
}

// 获取Sitemap配置
const loadSitemapConfig = async () => {
  try {
    const api = useApi()
    const response = await api.sitemapApi.getSitemapConfig()
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
    const api = useApi()
    await api.sitemapApi.updateSitemapConfig({
      autoGenerate: value,
      lastGenerate: sitemapConfig.value.lastGenerate,
      lastUpdate: new Date().toISOString()
    })
    message.success(value ? '自动生成功能已开启' : '自动生成功能已关闭')

    // 重新加载配置以同步前端状态
    await loadSitemapConfig()
  } catch (error) {
    message.error('更新配置失败')
  } finally {
    configLoading.value = false
  }
}

// 刷新Sitemap状态
const refreshSitemapStatus = async () => {
  try {
    const api = useApi()
    const response = await api.sitemapApi.getSitemapStatus()
    if (response) {
      sitemapStats.value = response
      generateStatus.value = '状态已刷新'
    }
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '刷新状态失败'
    message.error('刷新状态失败: ' + errorMsg)
  }
}

// 更新最后提交时间
const updateLastSubmitTime = (engine: string, time: string) => {
  lastSubmitTime.value[engine as keyof typeof lastSubmitTime.value] = time
}

// 加载外链列表
const loadLinkList = () => {
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

// 加载Bing索引配置
const loadBingIndexConfig = async () => {
  try {
    console.log('开始加载 Bing 索引配置...')
    const api = useApi()
    const response = await api.bingApi.getConfig()

    if (response?.success && response?.data) {
      bingIndexConfig.value = response.data
      console.log('Bing索引配置加载完成:', response.data)
    }
  } catch (error) {
    console.error('获取Bing索引配置失败:', error)
  }
}

// 更新Bing索引配置
const updateBingIndexConfig = async (value: boolean) => {
  configLoading.value = true
  try {
    const api = useApi()

    // 先更新本地状态
    bingIndexConfig.value.enabled = value

    // 调用后端API保存配置
    const response = await api.bingApi.updateConfig({
      enabled: value
    })

    if (response?.success) {
      message.success(response.message || 'Bing索引配置已更新')
      console.log('Bing索引配置更新成功:', response)
      
      // 延迟重新加载配置以验证后端状态（在后台进行，不阻塞UI）
      setTimeout(async () => {
        try {
          await loadBingIndexConfig()
        } catch (error) {
          console.error('重新加载配置失败:', error)
        }
      }, 1000)
    } else {
      message.error(response?.message || '更新配置失败')
      // 失败时恢复原状态
      bingIndexConfig.value.enabled = !value
    }
  } catch (error: any) {
    console.error('更新Bing索引配置失败:', error)
    const errorMsg = error?.response?.data?.message || error?.message || '更新配置失败'
    message.error('更新配置失败: ' + errorMsg)
    // 失败时恢复原状态
    bingIndexConfig.value.enabled = !value
  } finally {
    configLoading.value = false
  }
}

// 刷新Bing状态
const refreshBingStatus = async () => {
  try {
    await loadBingIndexConfig()
  } catch (error) {
    console.error('刷新Bing状态失败:', error)
    message.error('刷新状态失败')
  }
}

// 初始化
onMounted(async () => {
  await loadSystemConfig()
  await loadGoogleIndexConfig()
  await refreshGoogleIndexStatus()
  await loadSitemapConfig()
  await refreshSitemapStatus()
  await loadBingIndexConfig()
  loadLinkList()
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
</style>