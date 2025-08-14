<template>
  <n-modal :show="visible" @update:show="closeModal" preset="card" title="链接二维码" class="max-w-sm">
    <div class="text-center">
      <!-- 加载状态 -->
      <div v-if="loading" class="space-y-4">
        <div class="flex flex-col items-center justify-center py-8">
          <n-spin size="large" />
          <p class="text-sm text-gray-600 dark:text-gray-400 mt-4">正在获取链接...</p>
        </div>
      </div>
      
      <!-- 错误状态 -->
      <div v-else-if="error" class="space-y-4">
        <n-alert type="error" :show-icon="false">
          <template #icon>
            <i class="fas fa-exclamation-triangle text-red-500 mr-2"></i>
          </template>
          {{ error }}
        </n-alert>
        <n-card size="small">
          <p class="text-sm text-gray-700 dark:text-gray-300 break-all">{{ url }}</p>
        </n-card>
        <div class="flex gap-2">
          <n-button type="primary" @click="openLink" class="flex-1">
            <template #icon>
              <i class="fas fa-external-link-alt"></i>
            </template>
            跳转
          </n-button>
          <n-button type="success" @click="copyUrl" class="flex-1">
            <template #icon>
              <i class="fas fa-copy"></i>
            </template>
            复制
          </n-button>
        </div>
      </div>
      
      <!-- 正常显示 -->
      <div v-else>
        <!-- 移动端：所有链接都显示链接文本和操作按钮 -->
        <div v-if="isMobile" class="space-y-4">
          <!-- 显示链接状态信息 -->
          <n-alert v-if="message" type="info" :show-icon="false">
            <template #icon>
              <i class="fas fa-info-circle text-blue-500 mr-2"></i>
            </template>
            {{ message }}
          </n-alert>
          
          <n-card size="small">
            <p class="text-sm text-gray-700 dark:text-gray-300 break-all">{{ url }}</p>
          </n-card>
          <div class="flex gap-2">
            <n-button type="primary" @click="openLink" class="flex-1">
              <template #icon>
                <i class="fas fa-external-link-alt"></i>
              </template>
              跳转
            </n-button>
            <n-button type="success" @click="copyUrl" class="flex-1">
              <template #icon>
                <i class="fas fa-copy"></i>
              </template>
              复制
            </n-button>
          </div>
        </div>
      
      <!-- PC端：根据链接类型显示不同内容 -->
      <div v-else class="space-y-4">
        <!-- 显示链接状态信息 -->
        <n-alert v-if="message" type="info" :show-icon="false">
          <template #icon>
            <i class="fas fa-info-circle text-blue-500 mr-2"></i>
          </template>
          {{ message }}
        </n-alert>
        
        <!-- 夸克链接：只显示二维码 -->
        <div v-if="isQuarkLink" class="space-y-4">
          <div class="flex justify-center">
            <n-qr-code 
              :value="save_url || url" 
              :size="size" 
              :color="color"
              :background-color="backgroundColor"
              />
          </div>
          <div class="text-center">
            <n-button type="primary" @click="closeModal">
              <template #icon>
                <i class="fas fa-check"></i>
              </template>
              确认
            </n-button>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">请使用手机扫码操作</p>
          </div>
        </div>
        
        <!-- 其他链接：同时显示链接和二维码 -->
        <div v-else class="space-y-4">
          <div class="mb-4">
            <div class="flex justify-center">
              <n-qr-code :value="save_url || url" 
                :size="size"
                :color="color"
                :background-color="backgroundColor"
                />
            </div>
          </div>
          
          <div class="mb-4">
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">扫描二维码访问链接</p>
            <n-card size="small">
              <p class="text-xs text-gray-700 dark:text-gray-300 break-all">{{ url }}</p>
            </n-card>
          </div>
          
          <div class="flex gap-2">
            <n-button type="primary" @click="copyUrl" class="flex-1">
              <template #icon>
                <i class="fas fa-copy"></i>
              </template>
              复制链接
            </n-button>
            <n-button type="success" @click="downloadQrCode" class="flex-1">
              <template #icon>
                <i class="fas fa-download"></i>
              </template>
              下载二维码
            </n-button>
          </div>
        </div>
      </div>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
interface Props {
  visible: boolean
  save_url?: string
  url?: string
  loading?: boolean
  linkType?: string
  platform?: string
  message?: string
  error?: string
}

interface Emits {
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  url: ''
})
const emit = defineEmits<Emits>()

const size = ref(200)
const color = ref('#409eff')
const backgroundColor = ref('#F5F5F5')

// 检测是否为移动设备
const isMobile = ref(false)

// 检测设备类型
const detectDevice = () => {
  if (process.client) {
    isMobile.value = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
  }
}

// 判断是否为夸克链接
const isQuarkLink = computed(() => {
  return (props.url.includes('pan.quark.cn') || props.url.includes('quark.cn')) && !!props.save_url
})

// 关闭模态框
const closeModal = () => {
  emit('close')
}

// 复制链接
const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(props.url)
    // 可以添加一个简单的提示
    const button = event?.target as HTMLButtonElement
    if (button) {
      const originalText = button.innerHTML
      button.innerHTML = '<i class="fas fa-check"></i> 已复制'
      button.classList.add('bg-green-600')
      setTimeout(() => {
        button.innerHTML = originalText
        button.classList.remove('bg-green-600')
      }, 2000)
    }
  } catch (error) {
    console.error('复制失败:', error)
  }
}

// 跳转到链接
const openLink = () => {
  if (process.client) {
    window.open(props.url, '_blank')
  }
}

// 下载二维码
const downloadQrCode = () => {
  // 使用 Naive UI 的二维码组件，需要获取 DOM 元素
  const qrElement = document.querySelector('.n-qr-code canvas') as HTMLCanvasElement
  if (!qrElement) return
  
  try {
    const link = document.createElement('a')
    link.download = 'qrcode.png'
    link.href = qrElement.toDataURL()
    link.click()
  } catch (error) {
    console.error('下载失败:', error)
  }
}

// 监听visible变化
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    detectDevice()
  }
})
</script>

<style scoped>
/* 可以添加一些动画效果 */
.n-modal {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style> 