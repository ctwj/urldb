<template>
  <div class="simple-qr-code">
    <!-- 二维码显示区域 -->
    <div class="qr-display">
      <QRCodeDisplay
        ref="qrDisplayRef"
        :data="content"
        :preset="selectedPreset"
        :foreground-color="foregroundColor"
        :background-color="backgroundColor"
        :dot-type="dotType"
        :width="size"
        :height="size"
        :custom-image="customLogoUrl"
        :custom-image-options="{ margin: logoMargin }"
      />
    </div>

    <!-- 基础配置 -->
    <div class="qr-config">
      <!-- 预设选择 -->
      <div class="input-group">
        <label>预设样式:</label>
        <select v-model="selectedPresetName" class="preset-select">
          <option value="">自定义</option>
          <option v-for="preset in presets" :key="preset.name" :value="preset.name">
            {{ preset.name }}
          </option>
        </select>
      </div>

      <div class="input-group">
        <label>内容:</label>
        <input 
          v-model="content" 
          type="text" 
          placeholder="输入二维码内容"
          class="content-input"
        />
      </div>
      
      <div class="input-group">
        <label>尺寸:</label>
        <input 
          v-model.number="size" 
          type="range" 
          min="100" 
          max="300" 
          class="size-slider"
        />
        <span>{{ size }}px</span>
      </div>
      
      <div class="input-group">
        <label>前景色:</label>
        <input v-model="foregroundColor" type="color" class="color-picker" />
      </div>
      
      <div class="input-group">
        <label>背景色:</label>
        <input v-model="backgroundColor" type="color" class="color-picker" />
      </div>
      
      <div class="input-group">
        <label>点样式:</label>
        <select v-model="dotType" class="style-select">
          <option value="square">方形</option>
          <option value="dots">圆点</option>
          <option value="rounded">圆角</option>
          <option value="classy">经典</option>
          <option value="classy-rounded">经典圆角</option>
          <option value="extra-rounded">超圆角</option>
        </select>
      </div>

      <!-- 自定义Logo -->
      <div class="input-group">
        <label>Logo URL:</label>
        <input 
          v-model="customLogoUrl" 
          type="text" 
          placeholder="输入Logo图片URL (可选)"
          class="content-input"
        />
      </div>

      <div class="input-group">
        <label>Logo边距:</label>
        <input 
          v-model.number="logoMargin" 
          type="range" 
          min="0" 
          max="20" 
          class="size-slider"
        />
        <span>{{ logoMargin }}px</span>
      </div>

      <div class="input-group" v-if="customLogoUrl">
        <button @click="clearCustomLogo" class="clear-btn">清除Logo</button>
      </div>

      <!-- 随机样式按钮 -->
      <div class="input-group">
        <button @click="randomizeStyle" class="random-btn">随机样式</button>
      </div>
      
      <div class="button-group">
        <button @click="downloadPNG" class="download-btn">下载 PNG</button>
        <button @click="downloadSVG" class="download-btn">下载 SVG</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DotType } from 'qr-code-styling'
import { ref, watch, computed, onMounted } from 'vue'
import QRCodeDisplay from './Display.vue'
import { allQrCodePresets, findPresetByName, getRandomPreset } from './presets'
import { createRandomColor } from './color'
import { preloadCommonLogos } from './image-utils'

// Props
interface Props {
  initialData?: string
  initialSize?: number
  initialForegroundColor?: string
  initialBackgroundColor?: string
  initialPreset?: string
}

const props = withDefaults(defineProps<Props>(), {
  initialData: 'https://example.com',
  initialSize: 200,
  initialForegroundColor: '#000000',
  initialBackgroundColor: '#FFFFFF',
  initialPreset: ''
})

// 响应式数据
const content = ref(props.initialData)
const size = ref(props.initialSize)
const foregroundColor = ref(props.initialForegroundColor)
const backgroundColor = ref(props.initialBackgroundColor)
const dotType = ref<DotType>('rounded')
const selectedPresetName = ref(props.initialPreset)
const customLogoUrl = ref('')
const logoMargin = ref(8)

// 组件引用
const qrDisplayRef = ref()

// 预设相关
const presets = allQrCodePresets

const selectedPreset = computed(() => {
  if (!selectedPresetName.value) return undefined
  return findPresetByName(selectedPresetName.value) || undefined
})

// 随机样式
const randomizeStyle = () => {
  const randomPreset = getRandomPreset()
  selectedPresetName.value = randomPreset.name
  foregroundColor.value = createRandomColor()
  backgroundColor.value = createRandomColor()
  dotType.value = ['square', 'dots', 'rounded', 'classy', 'classy-rounded', 'extra-rounded'][
    Math.floor(Math.random() * 6)
  ] as DotType
  size.value = Math.floor(Math.random() * 200) + 150
}

// 清除自定义logo
const clearCustomLogo = () => {
  customLogoUrl.value = ''
}

// 下载 PNG
const downloadPNG = async () => {
  try {
    const dataURL = await qrDisplayRef.value?.downloadPNG()
    const link = document.createElement('a')
    link.download = 'qrcode.png'
    link.href = dataURL
    link.click()
  } catch (error) {
    console.error('下载失败:', error)
  }
}

// 下载 SVG
const downloadSVG = async () => {
  try {
    const dataURL = await qrDisplayRef.value?.downloadSVG()
    const link = document.createElement('a')
    link.download = 'qrcode.svg'
    link.href = dataURL
    link.click()
  } catch (error) {
    console.error('下载失败:', error)
  }
}

// 监听预设变化
watch(selectedPresetName, (newPresetName) => {
  if (newPresetName) {
    const preset = findPresetByName(newPresetName)
    if (preset) {
      // 应用预设样式
      foregroundColor.value = preset.dotsOptions.color || '#000000'
      backgroundColor.value = preset.backgroundOptions.color || '#FFFFFF'
      dotType.value = preset.dotsOptions.type || 'rounded'
      size.value = preset.width || 200
      // 清除自定义logo，使用预设的logo
      customLogoUrl.value = ''
    }
  } else {
    // 选择自定义时，保持当前设置
  }
})

// 组件挂载时预加载常用Logo
onMounted(async () => {
  await preloadCommonLogos()
})
</script>

<style scoped>
.simple-qr-code {
  display: flex;
  gap: 2rem;
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.qr-display {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1rem;
}

.qr-config {
  flex: 1;
  max-width: 300px;
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.input-group {
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.input-group label {
  min-width: 60px;
  font-size: 14px;
  color: #374151;
}

.content-input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
}

.content-input:focus {
  outline: none;
  border-color: #3b82f6;
}

.preset-select {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.preset-select:focus {
  outline: none;
  border-color: #3b82f6;
}

.size-slider {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: #e5e7eb;
  outline: none;
  cursor: pointer;
}

.size-slider::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
}

.color-picker {
  width: 40px;
  height: 32px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.style-select {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
  background: white;
}

.random-btn {
  flex: 1;
  padding: 0.5rem 1rem;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.random-btn:hover {
  background: #059669;
}

.clear-btn {
  flex: 1;
  padding: 0.5rem 1rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.clear-btn:hover {
  background: #dc2626;
}

.button-group {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

.download-btn {
  flex: 1;
  padding: 0.5rem 1rem;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.download-btn:hover {
  background: #2563eb;
}

/* 响应式设计 */
@media (max-width: 600px) {
  .simple-qr-code {
    flex-direction: column;
    gap: 1rem;
  }
  
  .qr-config {
    max-width: none;
  }
}
</style> 