# QRCode 组件库

基于原 Mini QR 项目提取的二维码显示组件，支持预设样式和自定义配置。

## 功能特性

- 🎨 **预设样式**：内置 26 种精美预设样式
- 🔧 **自定义配置**：支持颜色、点样式、尺寸等自定义
- 🖼️ **自定义Logo**：支持自定义Logo图片和边距调整
- 📱 **响应式设计**：适配移动端和桌面端
- 🖼️ **多格式导出**：支持 PNG、SVG、JPG 格式
- 🎲 **随机样式**：一键生成随机样式
- 🔧 **TypeScript 支持**：完整的类型定义

## 组件说明

### 1. QRCodeDisplay.vue - 纯显示组件

只负责显示二维码，支持预设和自定义配置。

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | string | - | 二维码内容（必需） |
| preset | Preset | null | 预设样式配置 |
| width | number | 200 | 二维码宽度 |
| height | number | 200 | 二维码高度 |
| foregroundColor | string | '#000000' | 前景色 |
| backgroundColor | string | '#FFFFFF' | 背景色 |
| dotType | DotType | 'rounded' | 点样式类型 |
| cornerSquareType | CornerSquareType | 'extra-rounded' | 角点样式类型 |
| cornerDotType | CornerDotType | 'dot' | 角点类型 |
| errorCorrectionLevel | 'L' \| 'M' \| 'Q' \| 'H' | 'Q' | 纠错级别 |
| margin | number | 0 | 边距 |
| type | DrawType | 'svg' | 渲染类型 |
| borderRadius | string | '0px' | 容器圆角 |
| background | string | 'transparent' | 容器背景色 |
| customImage | string | undefined | 自定义Logo图片URL |
| customImageOptions | object | undefined | 自定义Logo配置选项 |

#### 方法

| 方法 | 返回值 | 说明 |
|------|--------|------|
| downloadPNG() | Promise<string> | 获取 PNG 格式的 dataURL |
| downloadSVG() | Promise<string> | 获取 SVG 格式的 dataURL |
| downloadJPG() | Promise<string> | 获取 JPG 格式的 dataURL |

#### 使用示例

```vue
<template>
  <div>
    <!-- 使用预设 -->
    <QRCodeDisplay
      ref="qrRef"
      :data="qrData"
      :preset="selectedPreset"
      :width="200"
      :height="200"
    />
    
    <!-- 使用自定义Logo -->
    <QRCodeDisplay
      :data="qrData"
      :custom-image="customLogoUrl"
      :custom-image-options="{ margin: 8 }"
      :width="200"
      :height="200"
    />
    
    <button @click="downloadQR">下载二维码</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QRCodeDisplay, findPresetByName } from '@/components/QRCode'

const qrData = ref('https://example.com')
const selectedPreset = findPresetByName('Colorful')
const customLogoUrl = ref('https://api.iconify.design/ion:logo-github.svg?color=%23000')
const qrRef = ref()

const downloadQR = async () => {
  try {
    const dataURL = await qrRef.value.downloadPNG()
    const link = document.createElement('a')
    link.download = 'qrcode.png'
    link.href = dataURL
    link.click()
  } catch (error) {
    console.error('下载失败:', error)
  }
}
</script>
```

### 2. SimpleQRCode.vue - 完整功能组件

包含配置界面和二维码显示，内置预设选择功能和自定义Logo支持。

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| initialData | string | 'https://example.com' | 初始二维码内容 |
| initialSize | number | 200 | 初始尺寸 |
| initialForegroundColor | string | '#000000' | 初始前景色 |
| initialBackgroundColor | string | '#FFFFFF' | 初始背景色 |
| initialPreset | string | '' | 初始预设名称 |

#### 功能特性

- **预设选择**：从内置预设中选择样式
- **自定义配置**：调整内容、尺寸、颜色、点样式
- **自定义Logo**：输入Logo图片URL，支持PNG、SVG、JPG等格式
- **Logo边距调整**：控制Logo与二维码的间距
- **随机样式**：一键生成随机样式
- **下载功能**：支持PNG和SVG格式下载
- **响应式设计**：适配移动端和桌面端

#### 自定义Logo使用说明

1. 选择"自定义"预设
2. 在"Logo URL"输入框中输入图片URL
3. 调整"Logo边距"滑块控制间距
4. 点击"清除Logo"按钮可移除自定义Logo
5. 选择预设时会自动使用预设的Logo

#### 使用示例

```vue
<template>
  <div>
    <SimpleQRCode 
      :initial-data="'https://example.com'"
      :initial-preset="'Colorful'"
    />
  </div>
</template>

<script setup lang="ts">
import { SimpleQRCode } from '@/components/QRCode'
</script>
```

## 预设样式

### 内置预设

#### 自定义预设
| 预设名称 | 描述 | 特点 |
|----------|------|------|
| Plain | 简洁 | 黑白方形，经典样式 |
| Rounded | 圆角 | 圆角设计，现代感 |
| Colorful | 多彩 | 蓝红绿配色，活力十足 |
| Dark | 暗色 | 白点黑底，科技感 |
| Gradient | 渐变 | 紫粉橙配色，温暖 |
| Minimal | 极简 | 灰色圆点，简约 |
| Tech | 科技 | 青色科技风 |
| Nature | 自然 | 绿色生态风 |
| Warm | 温暖 | 红橙黄暖色调 |
| Cool | 冷色 | 蓝紫粉冷色调 |

#### 原项目预设
| 预设名称 | 描述 | 特点 |
|----------|------|------|
| Padlet | Padlet 风格 | 绿色圆角设计 |
| Vercel Light | Vercel 浅色 | 简洁现代风格 |
| Vercel Dark | Vercel 深色 | 科技感设计 |
| Supabase Green | Supabase 绿色 | 数据库风格 |
| Supabase Purple | Supabase 紫色 | 优雅设计 |
| UIlicious | UI 测试风格 | 红色圆角设计 |
| ViteConf 2023 | Vite 会议主题 | 紫色科技风 |
| Vue.js | Vue.js 主题 | 绿色框架风格 |
| Vue i18n | Vue 国际化 | 红色设计 |
| LYQHT | 项目作者主题 | 红色圆角设计 |
| Pejuang Kode | Pejuang Kode 主题 | 深蓝红配色 |
| GeeksHacking | GeeksHacking 主题 | 黄色经典设计 |
| SP Digital | SP Digital 主题 | 蓝色圆角设计 |
| GovTech - Stack Community | GovTech 社区主题 | 黑白简约设计 |
| QQ Group | QQ群聊主题 | 蓝紫渐变圆形设计 |
| WeChat Group | 微信群聊主题 | 经典黑白方形设计 |

### 使用预设

```typescript
import { allQrCodePresets, findPresetByName, getRandomPreset } from '@/components/QRCode'

// 获取所有预设
const presets = allQrCodePresets

// 根据名称查找预设
const colorfulPreset = findPresetByName('Colorful')

// 随机获取预设
const randomPreset = getRandomPreset()
```

## 样式类型

### 点样式 (DotType)
- `square` - 方形
- `dots` - 圆点
- `rounded` - 圆角
- `classy` - 经典
- `classy-rounded` - 经典圆角
- `extra-rounded` - 超圆角

### 角点样式 (CornerSquareType)
- `square` - 方形
- `extra-rounded` - 超圆角
- `dot` - 圆点

### 角点类型 (CornerDotType)
- `square` - 方形
- `dot` - 圆点

### 纠错级别
- `L` - 低 (7%)
- `M` - 中 (15%)
- `Q` - 高 (25%)
- `H` - 最高 (30%)

## 工具函数

### 颜色工具

```typescript
import { createRandomColor, getRandomItemInArray } from '@/components/QRCode'

// 生成随机颜色
const randomColor = createRandomColor()

// 从数组中随机选择
const randomItem = getRandomItemInArray(['a', 'b', 'c'])
```

## 完整示例

查看 `QRCodeExample.vue` 文件，了解完整的使用示例，包括：

- 预设选择和切换
- 自定义样式配置
- 随机样式生成
- 多格式下载
- 预设样式展示

## 文件结构

```
src/components/QRCode/
├── QRCodeDisplay.vue    # 纯显示组件
├── SimpleQRCode.vue     # 完整功能组件
├── QRCodeExample.vue    # 使用示例
├── presets.ts          # 预设配置
├── color.ts            # 颜色工具
├── index.ts            # 导出文件
└── README.md           # 说明文档
```

## 依赖

- Vue 3
- qr-code-styling
- TypeScript

## 许可证

基于原 Mini QR 项目的 GPL v3 许可证。 