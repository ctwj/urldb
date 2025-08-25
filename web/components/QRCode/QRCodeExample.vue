<template>
  <div class="qr-example">
    <h1>二维码组件使用示例</h1>
    
    <!-- 纯显示组件示例 -->
    <section class="example-section">
      <h2>1. 纯显示组件（支持预设）</h2>
      <div class="qr-container">
        <QRCodeDisplay
          ref="qrDisplayRef"
          :data="qrData"
          :preset="selectedPreset"
          :width="qrSize"
          :height="qrSize"
          :foreground-color="foregroundColor"
          :background-color="backgroundColor"
          :dot-type="dotType"
        />
      </div>
      
      <div class="controls">
        <div class="control-group">
          <label>预设:</label>
          <select v-model="selectedPresetName" @change="onPresetChange">
            <option value="">自定义</option>
            <option v-for="preset in presets" :key="preset.name" :value="preset.name">
              {{ preset.name }}
            </option>
          </select>
        </div>

        <div class="control-group">
          <label>内容:</label>
          <input v-model="qrData" type="text" placeholder="输入二维码内容" />
        </div>
        
        <div class="control-group">
          <label>尺寸:</label>
          <input v-model.number="qrSize" type="range" min="100" max="400" />
          <span>{{ qrSize }}px</span>
        </div>
        
        <div class="control-group">
          <label>前景色:</label>
          <input v-model="foregroundColor" type="color" />
        </div>
        
        <div class="control-group">
          <label>背景色:</label>
          <input v-model="backgroundColor" type="color" />
        </div>
        
        <div class="control-group">
          <label>点样式:</label>
          <select v-model="dotType">
            <option value="square">方形</option>
            <option value="dots">圆点</option>
            <option value="rounded">圆角</option>
            <option value="classy">经典</option>
            <option value="classy-rounded">经典圆角</option>
            <option value="extra-rounded">超圆角</option>
          </select>
        </div>
        
        <div class="button-group">
          <button @click="downloadAsPNG">下载 PNG</button>
          <button @click="downloadAsSVG">下载 SVG</button>
          <button @click="randomizeStyle">随机样式</button>
        </div>
      </div>
    </section>

    <!-- 完整功能组件示例 -->
    <section class="example-section">
      <h2>2. 完整功能组件（支持自定义Logo）</h2>
      <SimpleQRCode 
        :initial-data="'https://example.com'"
        :initial-preset="'Colorful'"
      />
      <div class="feature-note">
        <p>💡 <strong>新功能:</strong> 现在可以自定义Logo了！</p>
        <ul>
          <li>选择"自定义"预设，然后输入Logo图片URL</li>
          <li>调整Logo边距大小</li>
          <li>支持PNG、SVG、JPG等格式的图片</li>
          <li>选择预设时会自动使用预设的Logo</li>
        </ul>
      </div>
    </section>

    <!-- 预设展示 -->
    <section class="example-section">
      <h2>3. 预设样式展示</h2>
      <div class="preset-grid">
        <div 
          v-for="preset in presets" 
          :key="preset.name"
          class="preset-item"
          @click="selectPreset(preset.name)"
        >
          <div class="preset-qr">
            <QRCodeDisplay
              :data="'https://example.com'"
              :preset="preset"
              :width="120"
              :height="120"
            />
          </div>
          <div class="preset-name">{{ preset.name }}</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import QRCodeDisplay from './Display.vue'
import SimpleQRCode from './Simple.vue'
import { allQrCodePresets, findPresetByName, getRandomPreset } from './presets'
import { createRandomColor } from './color'

// 响应式数据
const qrData = ref('https://example.com')
const qrSize = ref(200)
const foregroundColor = ref('#000000')
const backgroundColor = ref('#FFFFFF')
const dotType = ref('rounded')
const selectedPresetName = ref('')

// 组件引用
const qrDisplayRef = ref()

// 预设相关
const presets = allQrCodePresets

const selectedPreset = computed(() => {
  if (!selectedPresetName.value) return null
  return findPresetByName(selectedPresetName.value) || null
})

// 预设变化处理
const onPresetChange = () => {
  if (selectedPresetName.value) {
    const preset = findPresetByName(selectedPresetName.value)
    if (preset) {
      foregroundColor.value = preset.dotsOptions.color
      backgroundColor.value = preset.backgroundOptions.color
      dotType.value = preset.dotsOptions.type
      qrSize.value = preset.width
    }
  }
}

// 选择预设
const selectPreset = (presetName: string) => {
  selectedPresetName.value = presetName
  onPresetChange()
}

// 随机样式
const randomizeStyle = () => {
  const randomPreset = getRandomPreset()
  selectedPresetName.value = randomPreset.name
  foregroundColor.value = createRandomColor()
  backgroundColor.value = createRandomColor()
  dotType.value = ['square', 'dots', 'rounded', 'classy', 'classy-rounded', 'extra-rounded'][
    Math.floor(Math.random() * 6)
  ]
  qrSize.value = Math.floor(Math.random() * 200) + 150
}

// 下载方法
const downloadAsPNG = async () => {
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

const downloadAsSVG = async () => {
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
</script>

<style scoped>
.qr-example {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

h1 {
  text-align: center;
  color: #1f2937;
  margin-bottom: 2rem;
}

.example-section {
  margin-bottom: 3rem;
  padding: 2rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: white;
}

h2 {
  color: #374151;
  margin-bottom: 1rem;
}

.qr-container {
  display: flex;
  justify-content: center;
  margin-bottom: 1rem;
  padding: 2rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.control-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.control-group label {
  font-weight: 500;
  color: #374151;
}

.control-group input,
.control-group select {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 14px;
}

.control-group input[type="color"] {
  width: 50px;
  height: 40px;
  padding: 0;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.control-group input[type="range"] {
  height: 6px;
  border-radius: 3px;
  background: #e5e7eb;
  outline: none;
  cursor: pointer;
}

.control-group input[type="range"]::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
}

.button-group {
  display: flex;
  gap: 0.5rem;
  grid-column: 1 / -1;
}

.button-group button {
  flex: 1;
  padding: 0.75rem 1rem;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.button-group button:hover {
  background: #2563eb;
}

.button-group button:last-child {
  background: #10b981;
}

.button-group button:last-child:hover {
  background: #059669;
}

/* 预设网格 */
.preset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1.5rem;
  margin-top: 1rem;
}

.preset-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-item:hover {
  border-color: #3b82f6;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.preset-qr {
  margin-bottom: 0.5rem;
  padding: 0.5rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.preset-name {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  text-align: center;
}

.feature-note {
  margin-top: 1rem;
  padding: 1rem;
  background: #f0f9ff;
  border: 1px solid #0ea5e9;
  border-radius: 8px;
  color: #0c4a6e;
}

.feature-note p {
  margin: 0 0 0.5rem 0;
  font-weight: 500;
}

.feature-note ul {
  margin: 0;
  padding-left: 1.5rem;
}

.feature-note li {
  margin-bottom: 0.25rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .qr-example {
    padding: 1rem;
  }
  
  .example-section {
    padding: 1rem;
  }
  
  .controls {
    grid-template-columns: 1fr;
  }

  .preset-grid {
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 1rem;
  }
}
</style> 