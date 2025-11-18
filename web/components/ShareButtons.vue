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
  return props.title && props.title !== 'undefined' ? props.title : '精彩资源分享'
})

const shareDescription = computed(() => {
  return props.description && props.description !== 'undefined' ? props.description : '发现更多优质资源，尽在urlDB'
})

const shareTags = computed(() => {
  if (props.tags && Array.isArray(props.tags) && props.tags.length > 0) {
    return props.tags.filter(tag => tag && tag !== 'undefined').slice(0, 3).join(',') || '资源分享,网盘,urldb'
  }
  return '资源分享,网盘,urldb'
})

// 获取完整URL - 使用运行时配置
const getFullUrl = () => {
  const config = useRuntimeConfig()

  if (props.url) {
    // 如果props.url已经是完整URL，则直接返回
    if (props.url.startsWith('http://') || props.url.startsWith('https://')) {
      return props.url
    }
    // 否则拼接站点URL
    let siteUrl = config.public.siteUrl
    if (!siteUrl || siteUrl === 'https://yourdomain.com') {
      // 优先在客户端使用当前页面的origin
      if (typeof window !== 'undefined') {
        siteUrl = window.location.origin
      } else {
        // 在服务端渲染时使用默认值
        siteUrl = process.env.NUXT_PUBLIC_SITE_URL || 'https://yourdomain.com'
      }
    }
    return `${siteUrl}${props.url.startsWith('/') ? props.url : '/' + props.url}`
  }

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
    shareElement.setAttribute('data-sites', 'facebook,twitter,reddit')
    shareElement.setAttribute('data-title', shareTitle.value)
    shareElement.setAttribute('data-description', shareDescription.value)
    shareElement.setAttribute('data-url', getFullUrl())
    shareElement.setAttribute('data-image', '') // 设置默认图片
    shareElement.setAttribute('data-pics', '') // 设置图片（QQ空间使用）
    shareElement.setAttribute('data-via', '') // Twitter via
    shareElement.setAttribute('data-wechat-qrcode-title', '微信扫一扫：分享')
    shareElement.setAttribute('data-wechat-qrcode-helper', '<p>微信里点"发现"，扫一下</p><p>二维码便可将本文分享至朋友圈。</p>')

    socialShareElement.value.appendChild(shareElement)

    // 初始化 social-share - 等待一段时间确保库已完全加载
    setTimeout(() => {
      console.log('检查 SocialShare 对象:', window.SocialShare)
      console.log('检查 social-share 元素:', shareElement)

      // 尝试使用 social-share.js 的正确初始化方式
      if (window.socialShare) {
        try {
          // 传入选择器来初始化
          window.socialShare('.social-share')
          console.log('socialShare() 函数调用成功')
        } catch (error) {
          console.error('socialShare 初始化失败:', error)
          // 如果上面失败，尝试另一种方式
          try {
            if (typeof window.socialShare === 'function') {
              window.socialShare()
              console.log('socialShare 全局调用成功')
            }
          } catch (error2) {
            console.error('socialShare 全局调用也失败:', error2)
          }
        }
      } else if (window.SocialShare) {
        try {
          window.SocialShare.init()
          console.log('SocialShare.init() 调用成功')
        } catch (error) {
          console.error('SocialShare 初始化失败:', error)
        }
      } else {
        console.error('SocialShare 对象不存在，库可能未正确加载')
      }
    }, 300)
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
      // 如果CDN加载失败，尝试备用链接
      const backupLink = document.createElement('link')
      backupLink.rel = 'stylesheet'
      backupLink.href = 'https://unpkg.com/social-share.js@1.0.16/dist/css/share.min.css'
      backupLink.onload = () => {
        console.log('备用 social-share.css 加载完成')
      }
      backupLink.onerror = () => {
        console.error('备用 social-share.css 也加载失败')
      }
      document.head.appendChild(backupLink)
    }
    document.head.appendChild(link)
  }

  if (!window.socialShare) {
    console.log('开始加载 social-share.js...')
    const script = document.createElement('script')
    script.src = 'https://cdn.jsdelivr.net/npm/social-share.js@1.0.16/dist/js/social-share.min.js'
    script.onload = () => {
      console.log('social-share.js 加载完成，检查全局对象:', window.socialShare)
      // 加载完成后初始化
      nextTick(() => {
        setTimeout(() => {
          initSocialShare()
        }, 300) // 稍微增加等待时间确保CSS和JS都完全加载
      })
    }
    script.onerror = () => {
      console.error('social-share.js 加载失败')
      // 如果CDN加载失败，尝试备用链接
      const backupScript = document.createElement('script')
      backupScript.src = 'https://unpkg.com/social-share.js@1.0.16/dist/js/social-share.min.js'
      backupScript.onload = () => {
        console.log('备用 social-share.js 加载完成，检查全局对象:', window.socialShare)
        nextTick(() => {
          setTimeout(() => {
            initSocialShare()
          }, 300)
        })
      }
      backupScript.onerror = () => {
        console.error('备用 social-share.js 也加载失败')
        // 如果无法加载外部库，创建基本分享按钮
        createFallbackShareButtons()
      }
      document.head.appendChild(backupScript)
    }
    document.head.appendChild(script)
  } else {
    // 如果已经加载过，直接初始化
    console.log('socialShare 已存在，直接初始化')
    initSocialShare()
  }
}

// 创建备选分享按钮，当social-share.js无法加载时使用
const createFallbackShareButtons = () => {
  if (typeof window === 'undefined' || !socialShareElement.value) return

  // 清空容器
  socialShareElement.value.innerHTML = ''

  // 创建包含基本分享功能的按钮
  const shareContainer = document.createElement('div')
  shareContainer.className = 'fallback-share-buttons'

  const fullUrl = getFullUrl()
  const encodedUrl = encodeURIComponent(fullUrl)
  const encodedTitle = encodeURIComponent(shareTitle.value)
  const encodedDesc = encodeURIComponent(shareDescription.value)

  // Facebook分享链接
  const facebookLink = document.createElement('a')
  facebookLink.href = `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}&t=${encodedTitle}`
  facebookLink.target = '_blank'
  facebookLink.innerHTML = '<i class="fa fa-facebook" style="font-size: 20px; color: #1877f2;"></i>'
  facebookLink.style.display = 'inline-block'
  facebookLink.style.margin = '0 3px'
  facebookLink.title = '分享到Facebook'

  // Twitter分享链接
  const twitterLink = document.createElement('a')
  twitterLink.href = `https://twitter.com/intent/tweet?url=${encodedUrl}&text=${encodedTitle}`
  twitterLink.target = '_blank'
  twitterLink.innerHTML = '<i class="fa fa-twitter" style="font-size: 20px; color: #1da1f2;"></i>'
  twitterLink.style.display = 'inline-block'
  twitterLink.style.margin = '0 3px'
  twitterLink.title = '分享到Twitter'

  // Reddit分享链接
  const redditLink = document.createElement('a')
  redditLink.href = `https://www.reddit.com/submit?url=${encodedUrl}&title=${encodedTitle}`
  redditLink.target = '_blank'
  redditLink.innerHTML = '<i class="fa fa-reddit" style="font-size: 20px; color: #ff4500;"></i>'
  redditLink.style.display = 'inline-block'
  redditLink.style.margin = '0 3px'
  redditLink.title = '分享到Reddit'

  // 添加到容器
  shareContainer.appendChild(facebookLink)
  shareContainer.appendChild(twitterLink)
  shareContainer.appendChild(redditLink)

  socialShareElement.value.appendChild(shareContainer)
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

/* 备选分享按钮样式 */
.fallback-share-buttons {
  display: flex;
  gap: 6px;
  align-items: center;
}

.fallback-share-buttons a {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  transition: all 0.2s ease;
  text-decoration: none;
}

.fallback-share-buttons a:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 响应式设计 */
@media (max-width: 640px) {
  .social-share-wrapper .social-share-icon,
  .fallback-share-buttons a {
    width: 26px !important;
    height: 26px !important;
  }
}
</style>