<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <label class="text-lg font-semibold text-gray-800 dark:text-gray-200">公告配置</label>
        <span class="text-xs text-gray-500 dark:text-gray-400">开启后可在网站显示公告信息</span>
      </div>
      <n-button v-if="modelValue.enable_announcements" @click="addAnnouncement" type="primary" size="small">
        <template #icon>
          <i class="fas fa-plus"></i>
        </template>
        添加公告
      </n-button>
    </div>
    <n-switch v-model:value="enableAnnouncements" />

    <!-- 公告列表 -->
    <div v-if="modelValue.enable_announcements && modelValue.announcements && modelValue.announcements.length > 0" class="announcement-list space-y-3">
      <div v-for="(announcement, index) in modelValue.announcements" :key="index" class="announcement-item border rounded-lg p-3 bg-gray-50 dark:bg-gray-800">
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center space-x-3">
            <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100">公告 {{ index + 1 }}</h4>
            <n-switch v-model:value="announcement.enabled" size="small" />
          </div>
          <div class="flex items-center space-x-1">
            <n-button text @click="moveAnnouncementUp(index)" :disabled="index === 0" size="small">
              <template #icon>
                <i class="fas fa-arrow-up"></i>
              </template>
            </n-button>
            <n-button text @click="moveAnnouncementDown(index)" :disabled="index === modelValue.announcements.length - 1" size="small">
              <template #icon>
                <i class="fas fa-arrow-down"></i>
              </template>
            </n-button>
            <n-button text @click="removeAnnouncement(index)" type="error" size="small">
              <template #icon>
                <i class="fas fa-trash"></i>
              </template>
            </n-button>
          </div>
        </div>

        <div class="space-y-2">
          <n-input
            v-model:value="announcement.content"
            placeholder="公告内容，支持HTML标签"
            type="textarea"
            :rows="2"
            size="small"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 公告接口
interface Announcement {
  content: string
  enabled: boolean
}

// 配置数据接口
interface ConfigData {
  enable_announcements: boolean
  announcements: Announcement[]
}

// Props
const props = defineProps<{
  modelValue: ConfigData
}>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: ConfigData]
}>()

// 计算属性用于双向绑定
const enableAnnouncements = computed({
  get: () => props.modelValue.enable_announcements,
  set: (value: boolean) => {
    emit('update:modelValue', {
      enable_announcements: value,
      announcements: props.modelValue.announcements
    })
  }
})

console.log(props.modelValue)

// 更新数据
const updateValue = (newValue: ConfigData) => {
  emit('update:modelValue', newValue)
}

// 添加公告
const addAnnouncement = () => {
  const newAnnouncements = [...props.modelValue.announcements, {
    content: '',
    enabled: true
  }]
  updateValue({
    enable_announcements: props.modelValue.enable_announcements,
    announcements: newAnnouncements
  })
}

// 删除公告
const removeAnnouncement = (index: number) => {
  const currentAnnouncements = Array.isArray(props.modelValue.announcements) ? props.modelValue.announcements : []
  const newAnnouncements = currentAnnouncements.filter((_, i) => i !== index)
  updateValue({
    enable_announcements: props.modelValue.enable_announcements,
    announcements: newAnnouncements
  })
}

// 上移公告
const moveAnnouncementUp = (index: number) => {
  if (index > 0) {
    const currentAnnouncements = Array.isArray(props.modelValue.announcements) ? props.modelValue.announcements : []
    const newAnnouncements = [...currentAnnouncements]
    const temp = newAnnouncements[index]
    newAnnouncements[index] = newAnnouncements[index - 1]
    newAnnouncements[index - 1] = temp
    updateValue({
      enable_announcements: props.modelValue.enable_announcements,
      announcements: newAnnouncements
    })
  }
}

// 下移公告
const moveAnnouncementDown = (index: number) => {
  const currentAnnouncements = Array.isArray(props.modelValue.announcements) ? props.modelValue.announcements : []
  if (index < currentAnnouncements.length - 1) {
    const newAnnouncements = [...currentAnnouncements]
    const temp = newAnnouncements[index]
    newAnnouncements[index] = newAnnouncements[index + 1]
    newAnnouncements[index + 1] = temp
    updateValue({
      enable_announcements: props.modelValue.enable_announcements,
      announcements: newAnnouncements
    })
  }
}
</script>

<style scoped>
.announcement-list {
  max-height: 400px;
  overflow-y: auto;
}
</style>