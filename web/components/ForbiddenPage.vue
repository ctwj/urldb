<template>
  <div class="forbidden-page">
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
      <p id="device-tip"><br/>3.请在浏览器中打开</p>
    </div>

    <p><br/><br/></p>
    <div class="app-download-tip">
      <span class="guidance-desc" id="current-url"></span>
    </div>
    <p><br/></p>
    <div class="app-download-tip">
      <span class="guidance-desc">点击右上角···图标 or 复制网址自行打开</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

// 组件属性
interface Props {
  currentUrl?: string
  isIOS?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  currentUrl: '',
  isIOS: false
})

// 在组件挂载时设置内容
onMounted(() => {
  const currentUrlElement = document.getElementById('current-url')
  const deviceTipElement = document.getElementById('device-tip')
  
  if (currentUrlElement) {
    currentUrlElement.textContent = props.currentUrl || window.location.href
  }
  
  if (deviceTipElement) {
    const deviceText = props.isIOS ? '苹果设备请在Safari浏览器中打开' : '安卓设备请在Chrome或其他浏览器中打开'
    deviceTipElement.innerHTML = `<br/>3.${deviceText}`
  }
})
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