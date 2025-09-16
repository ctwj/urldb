<template>
  <AdminPageLayout>
    <!-- 页面头部 - 标题 -->
    <template #page-header>
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">机器人管理</h1>
        <p class="text-gray-600 dark:text-gray-400">管理各平台的机器人配置和自动回复功能</p>
      </div>
    </template>

    <!-- 内容区 - 配置表单 -->
    <template #content>
      <div class="config-content h-full">
        <!-- 顶部Tabs -->
        <n-tabs
          v-model:value="activeTab"
          type="line"
          animated
          class="mb-6"
        >
          <n-tab-pane name="qq" tab="QQ机器人">
            <QqBotTab />
          </n-tab-pane>

        <n-tab-pane name="wechat" tab="微信公众号">
          <div class="tab-content-container">
            <div class="text-center py-12">
              <i class="fas fa-lock text-4xl text-gray-400 mb-4"></i>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能暂未开放</h3>
              <p class="text-gray-500 dark:text-gray-400">微信公众号机器人功能正在开发中，敬请期待</p>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="telegram" tab="Telegram机器人">
          <div class="tab-content-container">
            <TelegramBotTab />
          </div>
        </n-tab-pane>

        <n-tab-pane name="wechat_open" tab="微信开放平台">
          <div class="tab-content-container">
            <div class="text-center py-12">
              <i class="fas fa-lock text-4xl text-gray-400 mb-4"></i>
              <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">功能暂未开放</h3>
              <p class="text-gray-500 dark:text-gray-400">微信开放平台机器人功能正在开发中，敬请期待</p>
            </div>
          </div>
        </n-tab-pane>
        </n-tabs>

      </div>
    </template>
  </AdminPageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AdminPageLayout from '~/components/AdminPageLayout.vue'
import TelegramBotTab from '~/components/TelegramBotTab.vue'
import QqBotTab from '~/components/QqBotTab.vue'

// 设置页面布局
definePageMeta({
  layout: 'admin',
  ssr: false
})

const activeTab = ref('qq')

// 页面加载时获取配置
onMounted(async () => {
  console.log('机器人管理页面已加载')
})
</script>

<style scoped>
/* 机器人管理页面样式 */

.config-content {
  padding: 8px;
  background-color: var(--color-white, #ffffff);
}

.dark .config-content {
  background-color: var(--color-dark-bg, #1f2937);
}

/* tab内容容器 - 个别内容滚动 */
.tab-content-container {
  height: calc(100vh - 240px);
  overflow-y: auto;
  padding-bottom: 1rem;
}
</style>