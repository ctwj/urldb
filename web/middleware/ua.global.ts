import { defineNuxtRouteMiddleware, useRequestEvent } from 'nuxt/app'
import { getHeader, getRequestURL, setResponseStatus, setResponseHeader, send } from 'h3'

export default defineNuxtRouteMiddleware((to, from) => {
  // 只在服务端执行
  if (!process.server) return

  const event = useRequestEvent()
  if (!event) return

  const userAgent = getHeader(event, 'user-agent') || ''
  const isForbiddenApp = ['QQ/', 'MicroMessenger', 'WeiBo', 'DingTalk', 'Mail'].some(it => userAgent.includes(it))
  
  if (isForbiddenApp) {
    // 获取当前 URL
    const currentUrl = getRequestURL(event).href
    
    // 设置响应头
    setResponseStatus(event, 200)
    setResponseHeader(event, 'Content-Type', 'text/html; charset=utf-8')
    
    // 直接返回 HTML 响应
    return send(event, generateForbiddenPage(currentUrl, userAgent))
  }
})

// 生成禁止访问页面的函数
function generateForbiddenPage(url: string, userAgent: string) {
  return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>请在浏览器中打开</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #f8f9fa;
        }
        .forbidden-page {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        .top-bar-guidance {
            font-size: 15px;
            color: #fff;
            height: 70%;
            line-height: 1.2;
            padding-left: 20px;
            padding-top: 20px;
            background: url('/assets/images/banner.png') center right/cover no-repeat;
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
</head>
<body>
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
            <span class="guidance-desc" id="current-url">${url}</span>
        </div>
        <p><br/></p>
        <div class="app-download-tip">
            <span class="guidance-desc">点击右上角···图标 or 复制网址自行打开</span>
        </div>
    </div>
    

</body>
</html>`
}



