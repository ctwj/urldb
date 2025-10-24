<template>
  <div v-if="enableFloatButtons" class="float-right float-buttons">

   <!-- 返回顶部按钮 -->
   <a class="float-btn ontop fade show" @click="scrollToTop" title="返回顶部">
     <i class="fa fa-angle-up em12" style="color: #3498db;"></i>
   </a>

   <!-- 扫一扫在手机上体验 -->
   <span class="float-btn qrcode-btn hover-show">
     <i class="fa fa-qrcode" style="color: #27ae60;"></i>
     <div class="hover-show-con dropdown-menu qrcode-btn" style="width: 150px; height: auto;">
       <div class="qrcode" data-size="100">
         <n-qr-code :value="currentUrl" :size="100" :bordered="false" />
       </div>
       <div class="mt6 px12 muted-color">扫一扫在手机上体验</div>
     </div>
   </span>

   <!-- Telegram -->
   <span v-if="telegramQrImage" class="newadd-btns hover-show float-btn telegram-btn">
     <i class="fab fa-telegram-plane" style="color: #0088cc;"></i>
     <div class="hover-show-con dropdown-menu drop-newadd newadd-btns" style="width: 200px; height: auto;">
       <div class="image" data-size="100">
         <n-image :src="getImageUrl(telegramQrImage)" width="180" height="180" />
       </div>
       <div class="mt6 px12 muted-color text-center">扫码加入Telegram群组</div>
     </div>
   </span>

   <!-- 微信公众号 -->
   <a v-if="wechatSearchImage" class="float-btn service-wechat hover-show nowave" title="扫码添加微信" href="javascript:;">
     <i class="fab fa-weixin" style="color: #07c160;"></i>
     <div class="hover-show-con dropdown-menu service-wechat" style="width: 300px; height: auto;">
        <!-- <div class="image" data-size="100"> -->
         <n-image :src="getImageUrl(wechatSearchImage)" class="radius4" />
        <!-- </div> -->
     </div>
   </a>

  </div>
</template>

<script setup lang="ts">
// 导入系统配置store
import { useSystemConfigStore } from '~/stores/systemConfig'

// 获取系统配置store
const systemConfigStore = useSystemConfigStore()

// 使用图片URL composable
const { getImageUrl } = useImageUrl()

// 计算属性：是否启用浮动按钮
const enableFloatButtons = computed(() => {
  return systemConfigStore.config?.enable_float_buttons || false
})

// 计算属性：微信搜一搜图片
const wechatSearchImage = computed(() => {
  return systemConfigStore.config?.wechat_search_image || ''
})

// 计算属性：Telegram二维码图片
const telegramQrImage = computed(() => {
  return systemConfigStore.config?.telegram_qr_image || ''
})

// 滚动到顶部
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// 获取当前页面URL
const currentUrl = computed(() => {
  if (typeof window !== 'undefined') {
    return window.location.href
  }
  return ''
})
</script>

<style scoped>
/* 悬浮按钮容器 */
.float-right {
  position: fixed;
  right: 20px;
  bottom: 60px;
  z-index: 1030;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-bottom: env(safe-area-inset-bottom);
}

/* 悬浮按钮基础样式 */
.float-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(200, 200, 200, 0.4);
  border-radius: 50%;
  color: #666;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.float-btn:hover {
  background: rgba(200, 200, 200, 0.6);
  transform: translateY(-2px);
}

/* 返回顶部按钮特殊样式 */
.float-btn.ontop {
  opacity: 0;
  transform: translateY(10px);
  visibility: hidden;
}

.ontop.show {
  opacity: 1;
  transform: translateY(0);
  visibility: visible;
}

/* 悬浮菜单 */
.hover-show-con {
  position: absolute;
  right: 50px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 10px;
  opacity: 0;
  visibility: hidden;
  transform: translateY(10px);
  transition: all 0.3s ease;
  z-index: 1001;
}

.hover-show:hover .hover-show-con {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

/* 按钮位置样式 */
.hover-show-con.qrcode-btn {
  top: -60px;
}

.hover-show-con.newadd-btns {
  top: -100px;
}

.hover-show-con.service-wechat {
  top: -100px;
}

/* 图片容器 */
.image {
  text-align: center;
  padding: 5px;
}

.image .n-image {
  border-radius: 8px;
  overflow: hidden;
}

/* 居中文字 */
.text-center {
  text-align: center;
  margin-top: 8px;
  color: #666;
  font-size: 12px;
}

/* 悬浮菜单中的文字 */
.hover-show-con .muted-color {
  font-size: 12px;
}

/* 二维码容器 */
.qr-container {
  height: 200px;
  width: 200px;
  background-color: #F5F5F5;
}

.n-qr-code {
  padding: 0 !important;
}

/* 响应式 */
@media (max-width: 768px) {
  .float-right {
    right: 10px;
    gap: 8px;
  }

  .float-btn {
    width: 36px;
    height: 36px;
    font-size: 16px;
  }

  .hover-show-con {
    right: 46px;
    min-width: 140px;
  }

  /* 小屏下隐藏二维码和Telegram按钮 */
  .float-btn.qrcode-btn {
    display: none;
  }
}
float-buttons {
  font-size: 8px;
}
.dropdown-menu {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>