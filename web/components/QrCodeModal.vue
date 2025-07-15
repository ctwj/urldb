<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeModal">
    <div class="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-sm w-full mx-4" @click.stop>
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {{ isQuarkLink ? '夸克网盘链接' : '链接二维码' }}
        </h3>
        <button 
          @click="closeModal"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <i class="fas fa-times text-xl"></i>
        </button>
      </div>
      
      <div class="text-center">
        <!-- 移动端：所有链接都显示链接文本和操作按钮 -->
        <div v-if="isMobile" class="space-y-4">
          <div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
            <p class="text-sm text-gray-700 dark:text-gray-300 break-all">{{ url }}</p>
          </div>
          <div class="flex gap-2">
            <button 
              @click="openLink"
              class="flex-1 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors text-sm flex items-center justify-center gap-2"
            >
              <i class="fas fa-external-link-alt"></i> 跳转
            </button>
            <button 
              @click="copyUrl"
              class="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors text-sm flex items-center justify-center gap-2"
            >
              <i class="fas fa-copy"></i> 复制
            </button>
          </div>
        </div>
        
        <!-- PC端：根据链接类型显示不同内容 -->
        <div v-else class="space-y-4">
          <!-- 夸克链接：只显示二维码 -->
          <div v-if="isQuarkLink" class="space-y-4">
            <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg">
              <canvas ref="qrCanvas" class="mx-auto"></canvas>
            </div>
            <div class="text-center">
              <button 
                @click="closeModal"
                class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors text-sm flex items-center justify-center gap-2 mx-auto"
              >
                <i class="fas fa-check"></i> 确认
              </button>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">请使用手机扫码操作</p>
            </div>
          </div>
          
          <!-- 其他链接：同时显示链接和二维码 -->
          <div v-else class="space-y-4">
            <div class="mb-4">
              <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg">
                <canvas ref="qrCanvas" class="mx-auto"></canvas>
              </div>
            </div>
            
            <div class="mb-4">
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">扫描二维码访问链接</p>
              <div class="bg-gray-50 dark:bg-gray-700 p-3 rounded border">
                <p class="text-xs text-gray-700 dark:text-gray-300 break-all">{{ url }}</p>
              </div>
            </div>
            
            <div class="flex gap-2">
              <button 
                @click="copyUrl"
                class="flex-1 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors text-sm flex items-center justify-center gap-2"
              >
                <i class="fas fa-copy"></i>
                复制链接
              </button>
              <button 
                @click="downloadQrCode"
                class="flex-1 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors text-sm flex items-center justify-center gap-2"
              >
                <i class="fas fa-download"></i>
                下载二维码
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import QRCode from 'qrcode'

interface Props {
  visible: boolean
  url: string
}

interface Emits {
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const qrCanvas = ref<HTMLCanvasElement>()

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
  return props.url.includes('pan.quark.cn') || props.url.includes('quark.cn')
})

// 生成二维码
const generateQrCode = async () => {
  if (!qrCanvas.value || !props.url) return
  
  try {
    await QRCode.toCanvas(qrCanvas.value, props.url, {
      width: 200,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      }
    })
  } catch (error) {
    console.error('生成二维码失败:', error)
  }
}

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
  window.open(props.url, '_blank')
}

// 下载二维码
const downloadQrCode = () => {
  if (!qrCanvas.value) return
  
  try {
    const link = document.createElement('a')
    link.download = 'qrcode.png'
    link.href = qrCanvas.value.toDataURL()
    link.click()
  } catch (error) {
    console.error('下载失败:', error)
  }
}

// 监听visible变化，生成二维码
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    detectDevice()
    nextTick(() => {
      // PC端生成二维码（包括夸克链接）
      if (!isMobile.value) {
        generateQrCode()
      }
    })
  }
})

// 监听url变化，重新生成二维码
watch(() => props.url, () => {
  if (props.visible && !isMobile.value) {
    nextTick(() => {
      generateQrCode()
    })
  }
})
</script>

<style scoped>
/* 可以添加一些动画效果 */
.fixed {
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