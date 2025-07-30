<template>
  <div>
    <!-- 如果是禁止的APP，显示禁止页面 -->
    <div v-if="isForbiddenApp" class="forbidden-page">
      <div class="top-bar-guidance">
        <p class="top-bar-guidance-text">请按提示在手机 浏览器 打开<img src="/assets/images/3dian.png" class="icon-safari"></p>
        <p class="top-bar-guidance-text">苹果设备<img src="/assets/images/iphone.png" class="icon-safari">↗↗↗</p>
        <p class="top-bar-guidance-text">安卓设备<img src="/assets/images/android.png" class="icon-safari">↗↗↗</p>
      </div>

      <div id="contens">
        <p><br/><br/></p>
        <p>1.本站不支持 微信,QQ等APP 内访问</p>
        <p><br/></p>
        <p>2.请按提示在手机 浏览器 打开</p>
        <p v-if="isIOS"><br/>3.苹果设备请在Safari浏览器中打开</p>
        <p v-else><br/>3.安卓设备请在Chrome或其他浏览器中打开</p>
      </div>

      <p><br/><br/></p>
      <div class="app-download-tip">
        <span class="guidance-desc">{{ currentUrl }}</span>
      </div>
      <p><br/></p>
      <div class="app-download-tip">
        <span class="guidance-desc">点击右上角···图标 or 复制网址自行打开</span>
      </div>
    </div>
    
    <!-- 如果不是禁止的APP，显示正常页面 -->
    <NuxtPage v-else />
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isForbiddenApp = ref(false)
const isIOS = ref(false)
const currentUrl = ref('')

// 在渲染前就检测用户代理
if (process.client) {
  const ua = navigator.userAgent
  isForbiddenApp.value = ['QQ/', 'MicroMessenger', 'WeiBo', 'DingTalk', 'Mail'].some(it => ua.includes(it))
  
  if (isForbiddenApp.value) {
    currentUrl.value = window.location.href
    
    // 检测设备类型
    const userAgent = navigator.userAgent.toLowerCase()
    isIOS.value = /iphone|ipad|ipod/.test(userAgent) || 
                  (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
    
    // 修改页面标题
    document.title = '请在浏览器中打开'
  }
}
</script>

<style scoped>
.forbidden-page {
  min-height: 100vh;
}

.top-bar-guidance {
  font-size: 15px;
  color: #fff;
  height: 70%;
  line-height: 1.2;
  padding-left: 20px;
  padding-top: 20px;
  background: url('/assets/images/banner.png') center top/cover no-repeat;
}

.top-bar-guidance p {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.top-bar-guidance .icon-safari {
  width: 25px;
  height: 25px;
  vertical-align: middle;
  margin: 0 .2em;
}

.top-bar-guidance-text {
  display: flex;
  justify-items: center;
  word-wrap: nowrap;
}

.top-bar-guidance-text img {
  display: inline-block;
  width: 25px;
  height: 25px;
  vertical-align: middle;
  margin: 0 .2em;
}

#contens {
  font-weight: bold;
  color: #2466f4;
  text-align: center;
  font-size: 20px;
  margin-bottom: 125px;
}

.app-download-tip {
  margin: 0 auto;
  width: 290px;
  text-align: center;
  font-size: 15px;
  color: #2466f4;
  background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAAcAQMAAACak0ePAAAABlBMVEUAAAAdYfh+GakkAAAAAXRSTlMAQObYZgAAAA5JREFUCNdjwA8acEkAAAy4AIE4hQq/AAAAAElFTkSuQmCC) left center/auto 15px repeat-x;
}

.app-download-tip .guidance-desc {
  background-color: #fff;
  padding: 0 5px;
}
</style> 