import type { DrawType, Options as StyledQRCodeProps } from 'qr-code-styling'

export interface CustomStyleProps {
  borderRadius?: string
  background?: string
}

export type PresetAttributes = {
  style: CustomStyleProps
  name: string
}

export type Preset = Omit<
  Required<StyledQRCodeProps>,
  'shape' | 'qrOptions' | 'nodeCanvas' | 'jsdom'
> &
  PresetAttributes

const defaultPresetOptions = {
  backgroundOptions: {
    color: 'transparent'
  },
  imageOptions: {
    margin: 0,
    hideBackgroundDots: false,
    imageSize: 0.4,
    crossOrigin: undefined
  },
  width: 200,
  height: 200,
  margin: 0,
  type: 'svg' as DrawType
}

// 预设样式配置
export const plainPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Plain',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23000',
  dotsOptions: { color: '#000000', type: 'square' },
  cornersSquareOptions: { color: '#000000', type: 'square' },
  cornersDotOptions: { color: '#000000', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '0px', background: '#FFFFFF' }
}

export const roundedPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Rounded',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23000',
  dotsOptions: { color: '#000000', type: 'rounded' },
  cornersSquareOptions: { color: '#000000', type: 'extra-rounded' },
  cornersDotOptions: { color: '#000000', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#FFFFFF' }
}

export const colorfulPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Colorful',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%233B82F6',
  dotsOptions: { color: '#3B82F6', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#EF4444', type: 'extra-rounded' },
  cornersDotOptions: { color: '#10B981', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '16px', background: '#F8FAFC' }
}

export const darkPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Dark',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FFF',
  dotsOptions: { color: '#FFFFFF', type: 'classy' },
  cornersSquareOptions: { color: '#FFFFFF', type: 'square' },
  cornersDotOptions: { color: '#FFFFFF', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '8px', background: '#1F2937' }
}

export const gradientPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Gradient',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%238B5CF6',
  dotsOptions: { color: '#8B5CF6', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#EC4899', type: 'extra-rounded' },
  cornersDotOptions: { color: '#F59E0B', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '20px', background: '#FEF3C7' }
}

export const minimalPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Minimal',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%236B7280',
  dotsOptions: { color: '#6B7280', type: 'dots' },
  cornersSquareOptions: { color: '#6B7280', type: 'dot' },
  cornersDotOptions: { color: '#6B7280', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '4px', background: '#F9FAFB' }
}

export const techPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Tech',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%2300D4FF',
  dotsOptions: { color: '#00D4FF', type: 'classy' },
  cornersSquareOptions: { color: '#00D4FF', type: 'square' },
  cornersDotOptions: { color: '#00D4FF', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '0px', background: '#000000' }
}

export const naturePreset: Preset = {
  ...defaultPresetOptions,
  name: 'Nature',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23059669',
  dotsOptions: { color: '#059669', type: 'rounded' },
  cornersSquareOptions: { color: '#059669', type: 'extra-rounded' },
  cornersDotOptions: { color: '#10B981', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '24px', background: '#ECFDF5' }
}

export const warmPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Warm',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23DC2626',
  dotsOptions: { color: '#DC2626', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#EA580C', type: 'extra-rounded' },
  cornersDotOptions: { color: '#F59E0B', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '16px', background: '#FEF2F2' }
}

export const coolPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Cool',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%231E40AF',
  dotsOptions: { color: '#1E40AF', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#7C3AED', type: 'extra-rounded' },
  cornersDotOptions: { color: '#EC4899', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '20px', background: '#EFF6FF' }
}

// 原项目预设
export const padletPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Padlet',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%237ABE4A',
  dotsOptions: { color: '#7ABE4A', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#ed457e', type: 'extra-rounded' },
  cornersDotOptions: { color: '#ed457e', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '24px', background: '#000000' }
}

export const vercelLightPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Vercel Light',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:logo-vercel.svg?color=%23000',
  dotsOptions: { color: '#000000', type: 'classy' },
  cornersSquareOptions: { color: '#000000', type: 'square' },
  cornersDotOptions: { color: '#000000', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '0px', background: '#FFFFFF' }
}

export const vercelDarkPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Vercel Dark',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:logo-vercel.svg?color=%23FFF',
  dotsOptions: { color: '#FFFFFF', type: 'classy' },
  cornersSquareOptions: { color: '#FFFFFF', type: 'square' },
  cornersDotOptions: { color: '#FFFFFF', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '0px', background: '#000000' }
}

export const supabaseGreenPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Supabase Green',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/logos:supabase-icon.svg',
  dotsOptions: { color: '#3ecf8e', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#3ecf8e', type: 'square' },
  cornersDotOptions: { color: '#3ecf8e', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#000000' }
}

export const supabasePurplePreset: Preset = {
  ...defaultPresetOptions,
  name: 'Supabase Purple',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%237700ff',
  dotsOptions: { color: '#7700ff', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#7700ff', type: 'square' },
  cornersDotOptions: { color: '#7700ff', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#000000' }
}

export const uiliciousPreset: Preset = {
  ...defaultPresetOptions,
  name: 'UIlicious',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FF6B6B',
  dotsOptions: { color: '#FF6B6B', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#FF6B6B', type: 'extra-rounded' },
  cornersDotOptions: { color: '#FF6B6B', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '24px', background: '#FFFFFF' }
}

export const viteConf2023Preset: Preset = {
  ...defaultPresetOptions,
  name: 'ViteConf 2023',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23646CFF',
  dotsOptions: { color: '#646CFF', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#646CFF', type: 'square' },
  cornersDotOptions: { color: '#646CFF', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#000000' }
}

export const vueJsPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Vue.js',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%2342D392',
  dotsOptions: { color: '#42D392', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#42D392', type: 'square' },
  cornersDotOptions: { color: '#42D392', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#000000' }
}

export const vuei18nPreset: Preset = {
  ...defaultPresetOptions,
  name: 'Vue i18n',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FF6B6B',
  dotsOptions: { color: '#FF6B6B', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#FF6B6B', type: 'square' },
  cornersDotOptions: { color: '#FF6B6B', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '12px', background: '#000000' }
}

export const lyqhtPreset: Preset = {
  ...defaultPresetOptions,
  name: 'LYQHT',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23FF6B6B',
  dotsOptions: { color: '#FF6B6B', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#FF6B6B', type: 'extra-rounded' },
  cornersDotOptions: { color: '#FF6B6B', type: 'square' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '24px', background: '#000000' }
}

export const pejuangKodePreset: Preset = {
  ...defaultPresetOptions,
  name: 'Pejuang Kode',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23252f3f',
  dotsOptions: { color: '#252f3f', type: 'classy-rounded' },
  cornersSquareOptions: { color: '#252f3f', type: 'dot' },
  cornersDotOptions: { color: '#f05252', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '22px', background: '#ffffff' }
}

export const geeksHackingPreset: Preset = {
  ...defaultPresetOptions,
  name: 'GeeksHacking',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23cebe2c',
  dotsOptions: { color: '#cebe2c', type: 'classy' },
  cornersSquareOptions: { color: '#ced043', type: 'dot' },
  cornersDotOptions: { color: '#ced043', type: 'dot' },
  imageOptions: { margin: 2 },
  style: { borderRadius: '28px', background: '#000000' }
}

export const spDigitalPreset: Preset = {
  ...defaultPresetOptions,
  name: 'SP Digital',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%232196b0',
  dotsOptions: { color: '#2196b0', type: 'extra-rounded' },
  cornersSquareOptions: { color: '#2196b0', type: 'dot' },
  cornersDotOptions: { color: '#11b2b1', type: 'dot' },
  imageOptions: { margin: 2 },
  style: { borderRadius: '28px', background: '#ffffff' }
}

export const govtechStackCommunityPreset: Preset = {
  ...defaultPresetOptions,
  name: 'GovTech - Stack Community',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/ion:qr-code-outline.svg?color=%23000000',
  dotsOptions: { color: '#000000', type: 'square' },
  cornersSquareOptions: { color: '#000000', type: 'square' },
  cornersDotOptions: { color: '#000000', type: 'square' },
  imageOptions: { margin: 0 },
  style: { borderRadius: '24px', background: '#ffffff' }
}

export const qqGroupPreset: Preset = {
  ...defaultPresetOptions,
  name: 'QQ Group',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/simple-icons:qq.svg?color=%2371cdfc',
  dotsOptions: { color: '#71cdfc', type: 'dots' },
  cornersSquareOptions: { color: '#71cdfc', type: 'dot' },
  cornersDotOptions: { color: '#71cdfc', type: 'dot' },
  imageOptions: { margin: 8 },
  style: { borderRadius: '24px', background: '#ffffff' }
}

export const wechatGroupPreset: Preset = {
  ...defaultPresetOptions,
  name: 'WeChat Group',
  data: 'https://pan.l9.lc',
  image: 'https://api.iconify.design/simple-icons:wechat.svg?color=%23000000',
  dotsOptions: { color: '#000000', type: 'rounded' },
  cornersSquareOptions: { color: '#000000', type: 'rounded' },
  cornersDotOptions: { color: '#000000', type: 'rounded' },
  imageOptions: { margin: 8 },
  margin: 4,
  style: { borderRadius: '24px', background: '#ffffff' }
}



  // 预设列表
export const builtInPresets: Preset[] = [
  // 我们的自定义预设
  plainPreset,
  roundedPreset,
  colorfulPreset,
  darkPreset,
  gradientPreset,
  minimalPreset,
  techPreset,
  naturePreset,
  warmPreset,
  coolPreset,
  // 原项目预设
  padletPreset,
  vercelLightPreset,
  vercelDarkPreset,
  supabaseGreenPreset,
  supabasePurplePreset,
  uiliciousPreset,
  viteConf2023Preset,
  vueJsPreset,
  vuei18nPreset,
  lyqhtPreset,
  pejuangKodePreset,
  geeksHackingPreset,
  spDigitalPreset,
  govtechStackCommunityPreset,
  // 社交应用预设
  qqGroupPreset,
  wechatGroupPreset
]

// 默认预设
export const defaultPreset: Preset = builtInPresets[0]

// 获取所有预设
export const allQrCodePresets: Preset[] = builtInPresets

// 根据名称查找预设
export function findPresetByName(name: string): Preset | undefined {
  return allQrCodePresets.find(preset => preset.name === name)
}

// 随机获取预设
export function getRandomPreset(): Preset {
  return allQrCodePresets[Math.floor(Math.random() * allQrCodePresets.length)]
} 