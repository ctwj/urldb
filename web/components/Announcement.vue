<template>
  <div v-if="shouldShowAnnouncement" class="announcement-container px-3 py-1">
    <!-- 桌面端：显示完整公告内容 -->
    <div v-if="!isMobile" class="flex items-center justify-between min-h-[24px]">
      <div class="flex items-center gap-2 flex-1 overflow-hidden">
        <i class="fas fa-bullhorn text-blue-600 dark:text-blue-400 text-sm flex-shrink-0"></i>
        <div class="announcement-content overflow-hidden">
          <div class="announcement-item active">
            <span class="text-sm text-gray-700 dark:text-gray-300 truncate" v-html="validAnnouncements[currentIndex].content"></span>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 flex-shrink-0 ml-2">
        <span>{{ (currentIndex + 1) }}/{{ validAnnouncements.length }}</span>
        <button @click="nextAnnouncement" class="hover:text-blue-500 transition-colors">
          <i class="fas fa-chevron-right text-xs"></i>
        </button>
      </div>
    </div>

    <!-- 移动端：使用 Marquee 滚动显示 -->
    <div v-else class="flex items-center gap-2 min-h-[24px]">
      <i class="fas fa-bullhorn text-blue-600 dark:text-blue-400 text-sm flex-shrink-0"></i>
      <div class="flex-1 overflow-hidden">
        <n-marquee
          :speed="30"
          :delay="0"
          :loop="true"
          :auto-play="true"
          :pause-on-hover="true"
        >
          <span
            v-for="(announcement, index) in validAnnouncements"
            :key="index"
            class="text-sm text-gray-700 dark:text-gray-300 inline-block mx-4"
            v-html="announcement.content"
          ></span>
        </n-marquee>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 使用系统配置store获取公告数据
import { useSystemConfigStore } from '~/stores/systemConfig'
const systemConfigStore = useSystemConfigStore()
await systemConfigStore.initConfig(false, false)
const systemConfig = computed(() => systemConfigStore.config)

interface AnnouncementItem {
  content: string
  enabled: boolean
}

// 移动端检测
const isMobile = ref(false)

// 检测是否为移动端
const checkMobile = () => {
  if (process.client) {
    // 检测屏幕宽度
    isMobile.value = window.innerWidth < 768

    // 也可以使用用户代理检测
    const userAgent = navigator.userAgent.toLowerCase()
    const mobileKeywords = ['android', 'iphone', 'ipad', 'ipod', 'blackberry', 'windows phone']
    const isMobileDevice = mobileKeywords.some(keyword => userAgent.includes(keyword))

    // 结合屏幕宽度和设备类型判断
    isMobile.value = isMobile.value || isMobileDevice
  }
}

const currentIndex = ref(0)
const interval = ref<NodeJS.Timeout | null>(null)

// 计算有效公告（开启状态且有内容的公告）
const validAnnouncements = computed(() => {
  if (!systemConfig.value?.announcements) return []

  const announcements = Array.isArray(systemConfig.value.announcements)
    ? systemConfig.value.announcements
    : JSON.parse(systemConfig.value.announcements || '[]')

  return announcements.filter((item: AnnouncementItem) =>
    item.enabled && item.content && item.content.trim()
  )
})

// 判断是否应该显示公告
const shouldShowAnnouncement = computed(() => {
  return systemConfig.value?.enable_announcements && validAnnouncements.value.length > 0
})

// 自动切换公告
const startAutoSwitch = () => {
  if (validAnnouncements.value.length > 1) {
    interval.value = setInterval(() => {
      currentIndex.value = (currentIndex.value + 1) % validAnnouncements.value.length
    }, 4000) // 每4秒切换一次
  }
}

// 手动切换到下一条公告
const nextAnnouncement = () => {
  currentIndex.value = (currentIndex.value + 1) % validAnnouncements.value.length
}

// 监听公告数据变化，重新开始自动切换（仅桌面端）
watch(() => validAnnouncements.value.length, (newLength) => {
  if (newLength > 0) {
    currentIndex.value = 0
    stopAutoSwitch()
    if (!isMobile.value) {
      startAutoSwitch()
    }
  }
})

// 清理定时器
const stopAutoSwitch = () => {
  if (interval.value) {
    clearInterval(interval.value)
    interval.value = null
  }
}

onMounted(() => {
  // 初始化移动端检测
  checkMobile()

  // 监听窗口大小变化
  if (process.client) {
    window.addEventListener('resize', checkMobile)
  }

  if (shouldShowAnnouncement.value && !isMobile.value) {
    // 桌面端才启动自动切换
    startAutoSwitch()
  }
})

onUnmounted(() => {
  stopAutoSwitch()

  // 清理事件监听器
  if (process.client) {
    window.removeEventListener('resize', checkMobile)
  }
})
</script>

<style scoped>
.announcement-content {
  position: relative;
  height: 20px; /* 固定高度 */
}

.announcement-item {
  opacity: 0;
  transform: translateY(5px);
  transition: all 0.5s ease-in-out;
}

.announcement-item.active {
  opacity: 1;
  transform: translateY(0);
}

/* 移动端 Marquee 样式优化 */
@media (max-width: 767px) {
  .announcement-container {
    background: linear-gradient(90deg, transparent 0%, rgba(59, 130, 246, 0.05) 50%, transparent 100%);
    border-radius: 6px;
  }
}

/* Marquee 内文字样式 */
:deep(.n-marquee) {
  --n-bezier: cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.n-marquee__content) {
  display: flex;
  align-items: center;
  min-height: 20px;
}

/* 暗色主题适配 */
.dark-theme .announcement-container {
  background: transparent;
}

.dark .announcement-container {
  background: linear-gradient(90deg, transparent 0%, rgba(59, 130, 246, 0.1) 50%, transparent 100%);
}
</style>