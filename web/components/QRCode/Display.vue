<template>
  <div class="qr-code-display" :style="containerStyle">
    <div ref="qrCodeContainer" class="qr-wrapper" />
  </div>
</template>

<script setup lang="ts">
import type {
  CornerDotType,
  CornerSquareType,
  DotType,
  DrawType,
  Options as StyledQRCodeProps
} from 'qr-code-styling'
import QRCodeStyling from 'qr-code-styling'
import { onMounted, ref, watch, nextTick, computed } from 'vue'
import type { Preset } from './presets'

// Props
interface Props {
  data: string
  width?: number
  height?: number
  foregroundColor?: string
  backgroundColor?: string
  dotType?: DotType
  cornerSquareType?: CornerSquareType
  cornerDotType?: CornerDotType
  errorCorrectionLevel?: 'L' | 'M' | 'Q' | 'H'
  margin?: number
  type?: DrawType
  preset?: Preset
  borderRadius?: string
  background?: string
  customImage?: string
  customImageOptions?: {
    margin?: number
    hideBackgroundDots?: boolean
    imageSize?: number
    crossOrigin?: string
  }
}

const props = withDefaults(defineProps<Props>(), {
  width: 200,
  height: 200,
  foregroundColor: '#000000',
  backgroundColor: '#FFFFFF',
  dotType: 'rounded',
  cornerSquareType: 'extra-rounded',
  cornerDotType: 'dot',
  errorCorrectionLevel: 'Q',
  margin: 0,
  type: 'svg',
  borderRadius: '0px',
  background: 'transparent'
})

// DOM 引用
const qrCodeContainer = ref<HTMLElement>()

// QR Code 实例
let qrCodeInstance: QRCodeStyling | null = null

// 计算容器样式
const containerStyle = computed(() => {
  if (props.preset) {
    return {
      borderRadius: props.preset.style.borderRadius || '0px',
      background: props.preset.style.background || 'transparent',
      padding: '16px'
    }
  }
  return {
    borderRadius: props.borderRadius,
    background: props.background,
    padding: '16px'
  }
})

// 获取当前配置
const getCurrentConfig = () => {
  if (props.preset) {
    return {
      data: props.data,
      width: props.preset.width,
      height: props.preset.height,
      type: props.preset.type,
      margin: props.preset.margin,
      image: props.customImage || props.preset.image,
      imageOptions: {
        margin: (props.customImageOptions || props.preset.imageOptions)?.margin ?? 0,
        hideBackgroundDots: (props.customImageOptions || props.preset.imageOptions)?.hideBackgroundDots ?? true,
        imageSize: (props.customImageOptions || props.preset.imageOptions)?.imageSize ?? 0.3,
        crossOrigin: (props.customImageOptions || props.preset.imageOptions)?.crossOrigin ?? undefined
      },
      dotsOptions: props.preset.dotsOptions,
      backgroundOptions: props.preset.backgroundOptions,
      cornersSquareOptions: props.preset.cornersSquareOptions,
      cornersDotOptions: props.preset.cornersDotOptions,
      qrOptions: {
        errorCorrectionLevel: props.errorCorrectionLevel
      }
    }
  }

  return {
    data: props.data,
    width: props.width,
    height: props.height,
    type: props.type,
    margin: props.margin,
    image: props.customImage,
    imageOptions: {
      margin: props.customImageOptions?.margin ?? 0,
      hideBackgroundDots: props.customImageOptions?.hideBackgroundDots ?? false,
      imageSize: props.customImageOptions?.imageSize ?? 0.4,
      crossOrigin: props.customImageOptions?.crossOrigin ?? undefined
    },
    dotsOptions: {
      color: props.foregroundColor,
      type: props.dotType
    },
    backgroundOptions: {
      color: props.backgroundColor
    },
    cornersSquareOptions: {
      color: props.foregroundColor,
      type: props.cornerSquareType
    },
    cornersDotOptions: {
      color: props.foregroundColor,
      type: props.cornerDotType
    },
    qrOptions: {
      errorCorrectionLevel: props.errorCorrectionLevel
    }
  }
}

// 初始化 QR Code
const initQRCode = () => {
  if (!qrCodeContainer.value) return

  const config = getCurrentConfig()
  qrCodeInstance = new QRCodeStyling(config)
  qrCodeInstance.append(qrCodeContainer.value)
}

// 更新 QR Code
const updateQRCode = () => {
  if (!qrCodeInstance) return

  const config = getCurrentConfig()
  qrCodeInstance.update(config)
}

// 暴露方法给父组件
const downloadPNG = async (): Promise<string> => {
  if (!qrCodeInstance) throw new Error('QR Code not initialized')
  return await qrCodeInstance.getDataURL('png')
}

const downloadSVG = async (): Promise<string> => {
  if (!qrCodeInstance) throw new Error('QR Code not initialized')
  return await qrCodeInstance.getDataURL('svg')
}

const downloadJPG = async (): Promise<string> => {
  if (!qrCodeInstance) throw new Error('QR Code not initialized')
  return await qrCodeInstance.getDataURL('jpeg')
}

// 暴露方法
defineExpose({
  downloadPNG,
  downloadSVG,
  downloadJPG
})

// 监听 props 变化
watch(
  () => props,
  () => {
    nextTick(() => {
      updateQRCode()
    })
  },
  { deep: true }
)

// 组件挂载
onMounted(() => {
  initQRCode()
})
</script>

<style scoped>
.qr-code-display {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  transition: all 0.3s ease;
}

.qr-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style> 