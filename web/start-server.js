#!/usr/bin/env node

// 启动 Nitro 服务器的脚本
import { createNitro } from 'nitropack'
import { listen } from 'nitropack/dist/runtime/listen'

const nitro = await createNitro({
  preset: 'node-listener',
  rootDir: process.cwd(),
  logLevel: 'info'
})

const listener = await listen(nitro)
console.log(`🚀 Server listening on ${listener.url}`) 