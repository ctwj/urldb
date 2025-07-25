<template>
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
</template>

<script setup>
import { ref, onMounted } from 'vue'

// 设置使用 forbidden layout
definePageMeta({
  layout: 'single',
  ssr: false
})

const currentUrl = ref('')
const isIOS = ref(false)

onMounted(() => {
  currentUrl.value = window.location.href
  
  // 检测设备类型
  const userAgent = navigator.userAgent.toLowerCase()
  isIOS.value = /iphone|ipad|ipod/.test(userAgent) || 
                (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
})
</script>

<style scoped>
#browser img{
    width:50px;
}
#browser{
    margin: 0px 10px;
    text-align:center;
}
#contens{
    font-weight: bold;
    color: #2466f4;
    margin:-285px 0px 10px;
    text-align:center;
    font-size:20px;
    margin-bottom: 125px;
}
.top-bar-guidance{font-size:15px;color:#fff;height:70%;line-height:1.2;padding-left:20px;padding-top:20px;background:url(/assets/images/banner.png) center top/contain no-repeat}
.top-bar-guidance p{white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.top-bar-guidance .icon-safari{width:25px;height:25px;vertical-align:middle;margin:0 .2em}
.app-download-tip{margin:0 auto;width:290px;text-align:center;font-size:15px;color:#2466f4;background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAAcAQMAAACak0ePAAAABlBMVEUAAAAdYfh+GakkAAAAAXRSTlMAQObYZgAAAA5JREFUCNdjwA8acEkAAAy4AIE4hQq/AAAAAElFTkSuQmCC) left center/auto 15px repeat-x}
.app-download-tip .guidance-desc{background-color:#fff;padding:0 5px}
.app-download-btn{display:block;width:214px;height:40px;line-height:40px;margin:18px auto 0 auto;text-align:center;font-size:18px;color:#2466f4;border-radius:20px;border:.5px #2466f4 solid;text-decoration:none}
.top-bar-guidance-text{
  display: felx;
  justify-items: center;
  word-wrap: nowrap;
}
.top-bar-guidance-text img{
  display: inline-block;
  width: 25px;
  height: 25px;
  vertical-align: middle;
  margin: 0 .2em;
}
</style> 