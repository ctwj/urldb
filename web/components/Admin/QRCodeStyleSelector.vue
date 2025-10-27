<template>
  <n-modal
    :show="show"
    :mask-closable="true"
    preset="card"
    title="选择二维码样式"
    style="width: 90%; max-width: 900px;"
    @close="handleClose"
    @update:show="handleShowUpdate"
  >
    <div class="qr-style-selector">
      <!-- 样式选择区域 -->
      <div class="styles-section">
        <div class="styles-grid">
          <div
            v-for="preset in allQrCodePresets"
            :key="preset.name"
            class="style-item"
            :class="{ active: selectedPreset?.name === preset.name }"
            @click="selectPreset(preset)"
          >
            <div class="qr-preview" :style="preset.style">
              <div :ref="el => setQRContainer(el, preset.name)" v-if="preset"></div>
              <!-- 选中状态指示器 -->
              <div v-if="selectedPreset?.name === preset.name" class="selected-indicator">
                <i class="fas fa-check"></i>
              </div>
            </div>
            <div class="style-name">{{ preset.name }}</div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <n-button @click="handleClose">取消</n-button>
        <n-button type="primary" @click="confirmSelection">
          确认选择
        </n-button>
      </div>
    </div>
  </n-modal>
</template>

<script setup lang="ts">
import QRCodeStyling from 'qr-code-styling'
import { allQrCodePresets, type Preset } from '~/components/QRCode/presets'

// Props
const props = defineProps<{
  show: boolean
  currentStyle?: string
}>()

// Emits
const emit = defineEmits<{
  'update:show': [value: boolean]
  'select': [preset: Preset]
}>()

// 示例数据
const sampleData = ref('https://pan.l9.lc')

// 当前选中的预设
const selectedPreset = ref<Preset | null>(null)

// QR码实例映射
const qrInstances = ref<Map<string, QRCodeStyling>>(new Map())

// 监听显示状态变化
watch(() => props.show, (newShow) => {
  if (newShow) {
    // 查找当前样式对应的预设
    const currentPreset = allQrCodePresets.find(preset => preset.name === props.currentStyle)
    selectedPreset.value = currentPreset || allQrCodePresets[0] // 默认选择 Plain

    // 延迟渲染QR码，确保DOM已经准备好
    nextTick(() => {
      renderAllQRCodes()
    })
  }
})

// 设置QR码容器
const setQRContainer = (el: HTMLElement, presetName: string) => {
  if (el) {
    // 先清空容器内容，防止重复添加
    el.innerHTML = ''

    const preset = allQrCodePresets.find(p => p.name === presetName)
    if (preset) {
      const qrInstance = new QRCodeStyling({
        data: sampleData.value,
        ...preset,
        width: 80,
        height: 80
      })
      qrInstance.append(el)

      // 保存实例引用
      qrInstances.value.set(presetName, qrInstance)
    }
  }
}

// 渲染所有QR码
const renderAllQRCodes = () => {
  // 这个函数会在 setQRContainer 中被调用
  // 这里不需要额外操作，因为每个组件都会自己渲染
}

// 选择预设
const selectPreset = (preset: Preset) => {
  selectedPreset.value = preset
}

// 确认选择
const confirmSelection = () => {
  if (selectedPreset.value) {
    emit('select', selectedPreset.value)
    handleClose()
  }
}

// 处理显示状态更新
const handleShowUpdate = (value: boolean) => {
  emit('update:show', value)
}

// 关闭弹窗
const handleClose = () => {
  emit('update:show', false)
}
</script>

<style scoped>
.qr-style-selector {
  padding: 20px;
}

/* 样式选择区域 */
.styles-section h3 {
  margin-bottom: 20px;
  color: var(--color-text-1);
  font-size: 18px;
  font-weight: 600;
}

.styles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
  max-height: 400px;
  overflow-y: auto;
  padding: 10px;
}

.style-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 15px;
  border: 2px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--color-card-bg);
}

.style-item:hover {
  border-color: var(--color-primary-soft);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.style-item.active {
  border-color: var(--color-primary);
  background: var(--color-primary-soft);
  box-shadow: 0 0 0 2px rgba(24, 160, 88, 0.2);
}

.qr-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  border-radius: 6px;
  transition: all 0.3s ease;
  position: relative;
}

/* 选中状态指示器 */
.selected-indicator {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 24px;
  height: 24px;
  background: var(--color-primary);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  z-index: 10;
}

.style-name {
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-2);
  text-align: center;
}

.style-item.active .style-name {
  color: var(--color-primary);
  font-weight: 600;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--color-border);
}

/* 暗色主题适配 */
.dark .styles-grid {
  background: var(--color-dark-bg);
}

.dark .style-item {
  background: var(--color-dark-card);
}

.dark .style-item:hover {
  background: var(--color-dark-card-hover);
}

.dark .style-item.active {
  background: rgba(24, 160, 88, 0.1);
}

/* 响应式 */
@media (max-width: 768px) {
  .qr-style-selector {
    padding: 15px;
  }

  .styles-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 15px;
    max-height: 300px;
  }

  .style-item {
    padding: 10px;
  }
}

/* 滚动条样式 */
.styles-grid::-webkit-scrollbar {
  width: 6px;
}

.styles-grid::-webkit-scrollbar-track {
  background: var(--color-border);
  border-radius: 3px;
}

.styles-grid::-webkit-scrollbar-thumb {
  background: var(--color-text-3);
  border-radius: 3px;
}

.styles-grid::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-2);
}
</style>