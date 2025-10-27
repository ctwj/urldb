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
import { onMounted, ref, watch, nextTick, computed, onUnmounted } from 'vue'
import type { Preset } from './presets'
import { imageLoader } from './image-utils'

// 防抖函数
const debounce = (fn: Function, delay: number) => {
  let timeoutId: NodeJS.Timeout
  return (...args: any[]) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

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
  className?: string
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
    const style = {
      borderRadius: props.preset.style.borderRadius || '0px',
      background: props.preset.style.background || 'transparent',
      padding: '16px'
    }

    // 如果预设有className，添加到样式中
    if (props.preset.style.className) {
      return {
        ...style,
        class: props.preset.style.className
      }
    }
    return style
  }

  const style = {
    borderRadius: props.borderRadius,
    background: props.background,
    padding: '16px'
  }

  // 如果props有className，添加到样式中
  if (props.className) {
    return {
      ...style,
      class: props.className
    }
  }

  return style
})

// 生成配置键，用于缓存
const generateConfigKey = () => {
  if (props.preset) {
    return `${props.preset.name}-${props.data}-${props.width}-${props.height}-${props.customImage || props.preset.image}-${props.errorCorrectionLevel}`
  }
  return `${props.data}-${props.width}-${props.height}-${props.foregroundColor}-${props.backgroundColor}-${props.customImage}-${props.dotType}-${props.cornerSquareType}-${props.cornerDotType}-${props.errorCorrectionLevel}-${props.margin}-${props.type}`
}

// 获取当前配置
const getCurrentConfig = () => {
  const configKey = generateConfigKey()

  // 如果配置未变化，返回缓存的配置
  if (lastConfig && configKey === lastConfigKey) {
    return lastConfig
  }

  let config: any

  if (props.preset) {
    config = {
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
  } else {
    config = {
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

  // 缓存配置
  lastConfig = config
  lastConfigKey = configKey

  return config
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

// 配置对象缓存
let lastConfig: any = null
let lastConfigKey = ''

// 监听关键 props 变化
watch([
  () => props.data,
  () => props.preset,
  () => props.width,
  () => props.height,
  () => props.foregroundColor,
  () => props.backgroundColor,
  () => props.customImage,
  () => props.customImageOptions,
  () => props.dotType,
  () => props.cornerSquareType,
  () => props.cornerDotType,
  () => props.errorCorrectionLevel,
  () => props.margin,
  () => props.type
], async () => {
  // 预加载新图片
  const config = getCurrentConfig()
  if (config.image) {
    try {
      await imageLoader.preloadImage(config.image)
    } catch (error) {
      console.warn('Failed to preload QR code image:', error)
    }
  }

  nextTick(() => {
    debouncedUpdateQRCode()
  })
})

// 防抖更新，避免频繁重绘
const debouncedUpdateQRCode = debounce(updateQRCode, 50)

// 组件挂载
onMounted(async () => {
  // 预加载当前配置中的图片
  const config = getCurrentConfig()
  if (config.image) {
    try {
      await imageLoader.preloadImage(config.image)
    } catch (error) {
      console.warn('Failed to preload QR code image:', error)
    }
  }

  initQRCode()
})

// 组件卸载时清理
onUnmounted(() => {
  if (qrCodeInstance) {
    // 清理 QRCode 实例
    qrCodeInstance = null
  }
  // 清空缓存
  lastConfig = null
  lastConfigKey = ''
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