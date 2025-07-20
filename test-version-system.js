// 测试版本系统
const testVersionSystem = async () => {
  console.log('测试版本系统...')
  
  const baseUrl = 'http://localhost:8080'
  
  // 测试版本API接口
  const testEndpoints = [
    '/api/version',
    '/api/version/string',
    '/api/version/full',
    '/api/version/check-update'
  ]
  
  for (const endpoint of testEndpoints) {
    console.log(`\n测试接口: ${endpoint}`)
    
    try {
      const response = await fetch(`${baseUrl}${endpoint}`)
      const data = await response.json()
      
      console.log(`状态码: ${response.status}`)
      console.log(`响应:`, JSON.stringify(data, null, 2))
      
      if (data.success) {
        console.log('✅ 接口测试通过')
      } else {
        console.log('❌ 接口测试失败')
      }
      
    } catch (error) {
      console.error(`❌ 请求失败:`, error.message)
    }
  }
  
  // 测试版本管理脚本
  console.log('\n测试版本管理脚本...')
  
  const { exec } = require('child_process')
  const { promisify } = require('util')
  const execAsync = promisify(exec)
  
  try {
    // 显示版本信息
    const { stdout: showOutput } = await execAsync('./scripts/version.sh show')
    console.log('版本信息:')
    console.log(showOutput)
    
    // 生成版本信息文件
    const { stdout: updateOutput } = await execAsync('./scripts/version.sh update')
    console.log('生成版本信息文件:')
    console.log(updateOutput)
    
    console.log('✅ 版本管理脚本测试通过')
    
  } catch (error) {
    console.error('❌ 版本管理脚本测试失败:', error.message)
  }
  
  // 测试前端版本页面
  console.log('\n测试前端版本页面...')
  
  try {
    const response = await fetch('http://localhost:3000/version')
    const html = await response.text()
    
    console.log(`状态码: ${response.status}`)
    
    if (html.includes('版本信息') && html.includes('VersionInfo')) {
      console.log('✅ 前端版本页面测试通过')
    } else {
      console.log('❌ 前端版本页面测试失败')
    }
    
  } catch (error) {
    console.error('❌ 前端版本页面测试失败:', error.message)
  }
  
  console.log('\n✅ 版本系统测试完成')
}

// 运行测试
testVersionSystem() 