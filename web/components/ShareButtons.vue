<template>
  <div class="share-container">
    <!-- 直接显示分享按钮 -->
    <div
      ref="socialShareElement"
      class="social-share-wrapper"
    ></div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  url: {
    type: String,
    default: ''
  },
  tags: {
    type: Array,
    default: () => []
  }
})

const route = useRoute()

// 响应式数据
const socialShareElement = ref(null)

// 计算属性 - 避免在SSR中访问客户端API
const shareTitle = computed(() => {
  return props.title || '精彩资源分享'
})

const shareDescription = computed(() => {
  return props.description || '发现更多优质资源，尽在urlDB'
})

const shareTags = computed(() => {
  return props.tags?.slice(0, 3).join(',') || '资源分享,网盘,urldb'
})

// 获取完整URL - 仅在客户端调用
const getFullUrl = () => {
  if (props.url) return props.url
  if (typeof window !== 'undefined') {
    return `${window.location.origin}${route.fullPath}`
  }
  return route.fullPath
}

// 初始化 social-share - 仅在客户端调用
const initSocialShare = () => {
  if (typeof window === 'undefined') return

  if (socialShareElement.value) {
    // 清空容器
    socialShareElement.value.innerHTML = ''

    // 创建 social-share 元素
    const shareElement = document.createElement('div')
    shareElement.className = 'social-share'
    shareElement.setAttribute('data-sites', 'weibo,qq,wechat,qzone,twitter,telegram')
    shareElement.setAttribute('data-title', shareTitle.value)
    shareElement.setAttribute('data-description', shareDescription.value)
    shareElement.setAttribute('data-url', getFullUrl())
    shareElement.setAttribute('data-twitter', shareTags.value)
    shareElement.setAttribute('data-wechat-qrcode-title', '微信扫一扫：分享')
    shareElement.setAttribute('data-wechat-qrcode-helper', '<p>微信里点"发现"，扫一下</p><p>二维码便可将本文分享至朋友圈。</p>')

    socialShareElement.value.appendChild(shareElement)

    // 初始化 social-share - 等待一段时间确保库已完全加载
    setTimeout(() => {
      console.log('检查 SocialShare 对象:', window.SocialShare)
      console.log('检查 social-share 元素:', shareElement)

      // 尝试多种初始化方式
      if (window.SocialShare) {
        if (typeof window.SocialShare.init === 'function') {
          window.SocialShare.init()
          console.log('SocialShare.init() 调用成功')
        } else if (typeof window.SocialShare === 'function') {
          window.SocialShare()
          console.log('SocialShare() 函数调用成功')
        } else {
          console.log('SocialShare 对象存在但不是函数:', typeof window.SocialShare)
          // 尝试手动初始化
          try {
            const socialShareElements = document.querySelectorAll('.social-share')
            console.log('找到 social-share 元素:', socialShareElements.length)
            if (socialShareElements.length > 0) {
              // 检查是否已经生成了分享按钮
              const generatedButtons = socialShareElements[0].querySelectorAll('.social-share-icon')
              console.log('已生成的分享按钮:', generatedButtons.length)
            }
          } catch (e) {
            console.error('手动检查失败:', e)
          }
        }
      } else if (window.socialShare) {
        // 尝试使用 socialShare 变量
        console.log('找到 socialShare 全局变量，尝试初始化')
        console.log('socialShare 对象类型:', typeof window.socialShare)
        console.log('socialShare 对象内容:', window.socialShare)

        if (typeof window.socialShare.init === 'function') {
          try {
            window.socialShare.init()
            console.log('socialShare.init() 调用成功')
          } catch (error) {
            console.error('socialShare.init() 调用失败:', error)
          }
        } else if (typeof window.socialShare === 'function') {
          try {
            // social-share.js 需要传入选择器作为参数
            window.socialShare('.social-share')
            console.log('socialShare() 函数调用成功')
          } catch (error) {
            console.error('socialShare() 调用失败:', error)
            // 尝试不带参数调用
            try {
              window.socialShare()
              console.log('socialShare() 无参数调用成功')
            } catch (error2) {
              console.error('socialShare() 无参数调用也失败:', error2)
            }
          }
        } else {
          console.log('socialShare 对象存在但不是函数:', typeof window.socialShare)
          console.log('socialShare 对象的属性:', Object.keys(window.socialShare || {}))
        }
      } else {
        console.error('SocialShare 对象不存在，检查库是否正确加载')
        // 检查是否有其他全局变量
        console.log('可用全局变量:', Object.keys(window).filter(key => key.toLowerCase().includes('social')))
      }
    }, 500)
  }
}

// 动态加载 social-share.js 和 CSS - 仅在客户端调用
const loadSocialShare = () => {
  if (typeof window === 'undefined') return

  // 加载 CSS 文件
  if (!document.querySelector('link[href*="social-share.min.css"]')) {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = 'https://cdn.jsdelivr.net/npm/social-share.js@1.0.16/dist/css/share.min.css'
    link.onload = () => {
      console.log('social-share.css 加载完成')
    }
    link.onerror = () => {
      console.error('social-share.css 加载失败')
    }
    document.head.appendChild(link)
  }

  if (!window.SocialShare) {
    console.log('开始加载 social-share.js...')
    const script = document.createElement('script')
    script.src = 'https://cdn.jsdelivr.net/npm/social-share.js@1.0.16/dist/js/social-share.min.js'
    script.onload = () => {
      console.log('social-share.js 加载完成，检查全局对象:', window.SocialShare)
      // 加载完成后初始化
      nextTick(() => {
        setTimeout(() => {
          initSocialShare()
        }, 200) // 增加等待时间确保CSS和JS都完全加载
      })
    }
    script.onerror = () => {
      console.error('social-share.js 加载失败')
    }
    document.head.appendChild(script)
  } else {
    // 如果已经加载过，直接初始化
    console.log('SocialShare 已存在，直接初始化')
    initSocialShare()
  }
}

// 组件挂载时直接初始化 - 仅在客户端执行
onMounted(() => {
  if (typeof window !== 'undefined') {
    // 页面加载完成后直接初始化 social-share
    nextTick(() => {
      loadSocialShare()
    })
  }
})
</script>

<style scoped>
.share-container {
  position: relative;
  display: inline-block;
}

/* social-share.js 样式适配 */
.social-share-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

/* social-share.js 默认样式覆盖 */
.social-share-wrapper .social-share {
  display: flex !important;
  flex-wrap: wrap;
  gap: 6px;
}

.social-share-wrapper .social-share-icon {
  width: 28px !important;
  height: 28px !important;
  margin: 0 !important;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.social-share-wrapper .social-share-icon:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 暗色模式下的 social-share 图标 */
.dark .social-share-wrapper .social-share-icon {
  filter: brightness(0.9);
}

/* 响应式设计 */
@media (max-width: 640px) {
  .social-share-wrapper .social-share-icon {
    width: 26px !important;
    height: 26px !important;
  }
}
</style>