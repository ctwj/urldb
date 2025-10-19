<template>
  <div class="space-y-4">
    <div class="flex items-center space-x-2">
      <label class="text-lg font-semibold text-gray-800 dark:text-gray-200">浮动按钮配置</label>
      <span class="text-xs text-gray-500 dark:text-gray-400">开启后显示右下角浮动按钮</span>
    </div>
    <n-switch v-model:value="enableFloatButtons" />

    <!-- 浮动按钮设置 -->
    <div v-if="modelValue.enable_float_buttons" class="float-buttons-config space-y-4">
      <!-- 微信搜一搜图片 -->
      <div class="space-y-2">
        <div class="flex items-center space-x-2">
          <label class="text-base font-semibold text-gray-800 dark:text-gray-200">微信搜一搜图片</label>
          <span class="text-xs text-gray-500 dark:text-gray-400">选择微信搜一搜的二维码图片</span>
        </div>
        <div class="flex items-center space-x-4">
          <div v-if="modelValue.wechat_search_image" class="flex-shrink-0">
            <n-image
              :src="getImageUrl(modelValue.wechat_search_image)"
              alt="微信搜一搜"
              width="80"
              height="80"
              object-fit="cover"
              class="rounded-lg border"
            />
          </div>
          <div class="flex-1">
            <n-button type="primary" @click="openWechatSelector">
              <template #icon>
                <i class="fas fa-image"></i>
              </template>
              {{ modelValue.wechat_search_image ? '更换图片' : '选择图片' }}
            </n-button>
            <n-button v-if="modelValue.wechat_search_image" @click="clearWechatImage" class="ml-2">
              <template #icon>
                <i class="fas fa-times"></i>
              </template>
              清除
            </n-button>
          </div>
        </div>
      </div>

      <!-- Telegram二维码 -->
      <div class="space-y-2">
        <div class="flex items-center space-x-2">
          <label class="text-base font-semibold text-gray-800 dark:text-gray-200">Telegram二维码</label>
          <span class="text-xs text-gray-500 dark:text-gray-400">选择Telegram群组的二维码图片</span>
        </div>
        <div class="flex items-center space-x-4">
          <div v-if="modelValue.telegram_qr_image" class="flex-shrink-0">
            <n-image
              :src="getImageUrl(modelValue.telegram_qr_image)"
              alt="Telegram二维码"
              width="80"
              height="80"
              object-fit="cover"
              class="rounded-lg border"
            />
          </div>
          <div class="flex-1">
            <n-button type="primary" @click="openTelegramSelector">
              <template #icon>
                <i class="fas fa-image"></i>
              </template>
              {{ modelValue.telegram_qr_image ? '更换图片' : '选择图片' }}
            </n-button>
            <n-button v-if="modelValue.telegram_qr_image" @click="clearTelegramImage" class="ml-2">
              <template #icon>
                <i class="fas fa-times"></i>
              </template>
              清除
            </n-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

// 使用图片URL composable
const { getImageUrl } = useImageUrl()

<script setup lang="ts">
import ImageSelectorModal from '~/components/Admin/ImageSelectorModal.vue'
// 配置数据接口
interface ConfigData {
  enable_float_buttons: boolean
  wechat_search_image: string
  telegram_qr_image: string
}

// Props
const props = defineProps<{
  modelValue: ConfigData
}>()

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: ConfigData]
  'openWechatSelector': []
  'openTelegramSelector': []
}>()

// 计算属性用于双向绑定
const enableFloatButtons = computed({
  get: () => props.modelValue.enable_float_buttons,
  set: (value: boolean) => {
    emit('update:modelValue', {
      enable_float_buttons: value,
      wechat_search_image: props.modelValue.wechat_search_image,
      telegram_qr_image: props.modelValue.telegram_qr_image
    })
  }
})

// 使用图片URL composable
const { getImageUrl } = useImageUrl()

// 选择器状态
const showWechatSelector = ref(false)
const showTelegramSelector = ref(false)

// 清除微信图片
const clearWechatImage = () => {
  emit('update:modelValue', {
    enable_float_buttons: props.modelValue.enable_float_buttons,
    wechat_search_image: '',
    telegram_qr_image: props.modelValue.telegram_qr_image
  })
}

// 清除Telegram图片
const clearTelegramImage = () => {
  emit('update:modelValue', {
    enable_float_buttons: props.modelValue.enable_float_buttons,
    wechat_search_image: props.modelValue.wechat_search_image,
    telegram_qr_image: ''
  })
}

// 打开微信选择器
const openWechatSelector = () => {
  emit('openWechatSelector')
}

// 打开Telegram选择器
const openTelegramSelector = () => {
  emit('openTelegramSelector')
}

// 处理微信图片选择
const handleWechatImageSelect = (file: any) => {
  emit('update:modelValue', {
    enable_float_buttons: props.modelValue.enable_float_buttons,
    wechat_search_image: file.access_url,
    telegram_qr_image: props.modelValue.telegram_qr_image
  })
}

// 处理Telegram图片选择
const handleTelegramImageSelect = (file: any) => {
  emit('update:modelValue', {
    enable_float_buttons: props.modelValue.enable_float_buttons,
    wechat_search_image: props.modelValue.wechat_search_image,
    telegram_qr_image: file.access_url
  })
}
</script>

