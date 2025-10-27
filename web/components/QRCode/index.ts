// 导出组件
export { default as QRCodeDisplay } from './Display.vue'
export { default as SimpleQRCode } from './Simple.vue'
export { default as QRCodeExample } from './QRCodeExample.vue'

// 导出预设和工具
export * from './presets'
export * from './color'
export * from './image-utils'

// 导出类型
export type { Preset, PresetAttributes, CustomStyleProps } from './presets' 