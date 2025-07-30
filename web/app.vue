<template>
  <div>
    <!-- 加载状态 - 在检测完成前显示 -->
    <div id="loading-state" style="display: block;">
      <LoadingState />
    </div>

    <!-- 禁止页面 - 使用CSS控制显示 -->
    <div id="forbidden-page" style="display: none;">
      <ForbiddenPage :current-url="currentUrl" :is-ios="isIOS" />
    </div>
    
    <!-- 正常页面内容 -->
    <div id="normal-page" style="display: none;">
      <NuxtPage />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const currentUrl = ref('')
const isIOS = ref(false)

// 在客户端挂载后检测用户代理并控制显示
onMounted(() => {
  const ua = navigator.userAgent
  const isForbiddenApp = ['QQ/', 'MicroMessenger', 'WeiBo', 'DingTalk', 'Mail'].some(it => ua.includes(it))
  
  // 获取所有页面元素
  const loadingState = document.getElementById('loading-state')
  const forbiddenPage = document.getElementById('forbidden-page')
  const normalPage = document.getElementById('normal-page')
  
  if (isForbiddenApp) {
    // 设置禁止页面需要的属性
    currentUrl.value = window.location.href
    
    // 检测设备类型
    const userAgent = navigator.userAgent.toLowerCase()
    isIOS.value = /iphone|ipad|ipod/.test(userAgent) || 
                  (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
    
    // 隐藏加载状态，显示禁止页面
    if (loadingState && forbiddenPage && normalPage) {
      loadingState.style.display = 'none'
      forbiddenPage.style.display = 'block'
      normalPage.style.display = 'none'
    }
    
    // 修改页面标题
    document.title = '请在浏览器中打开'
  } else {
    // 显示正常页面
    if (loadingState && forbiddenPage && normalPage) {
      // 隐藏加载状态，显示正常页面
      loadingState.style.display = 'none'
      forbiddenPage.style.display = 'none'
      normalPage.style.display = 'block'
    }
  }
})
</script> 