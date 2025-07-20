// 测试GitHub版本系统
const testGitHubVersion = async () => {
  console.log('测试GitHub版本系统...')
  
  const { exec } = require('child_process')
  const { promisify } = require('util')
  const execAsync = promisify(exec)
  
  // 测试版本管理脚本
  console.log('\n1. 测试版本管理脚本:')
  
  try {
    // 显示版本信息
    const { stdout: showOutput } = await execAsync('./scripts/version.sh show')
    console.log('版本信息:')
    console.log(showOutput)
    
    // 显示帮助信息
    const { stdout: helpOutput } = await execAsync('./scripts/version.sh help')
    console.log('帮助信息:')
    console.log(helpOutput)
    
    console.log('✅ 版本管理脚本测试通过')
    
  } catch (error) {
    console.error('❌ 版本管理脚本测试失败:', error.message)
  }
  
  // 测试版本API接口
  console.log('\n2. 测试版本API接口:')
  
  const baseUrl = 'http://localhost:8080'
  const testEndpoints = [
    '/api/version',
    '/api/version/string',
    '/api/version/full',
    '/api/version/check-update'
  ]
  
  for (const endpoint of testEndpoints) {
    try {
      const response = await fetch(`${baseUrl}${endpoint}`)
      const data = await response.json()
      
      console.log(`\n接口: ${endpoint}`)
      console.log(`状态码: ${response.status}`)
      console.log(`响应:`, JSON.stringify(data, null, 2))
      
      if (data.success) {
        console.log('✅ 接口测试通过')
      } else {
        console.log('❌ 接口测试失败')
      }
      
    } catch (error) {
      console.error(`❌ 接口 ${endpoint} 测试失败:`, error.message)
    }
  }
  
  // 测试GitHub版本检查
  console.log('\n3. 测试GitHub版本检查:')
  
  try {
    const response = await fetch('https://api.github.com/repos/ctwj/urldb/releases/latest')
    const data = await response.json()
    
    console.log('GitHub API响应:')
    console.log(`状态码: ${response.status}`)
    console.log(`最新版本: ${data.tag_name || 'N/A'}`)
    console.log(`发布日期: ${data.published_at || 'N/A'}`)
    
    if (data.tag_name) {
      console.log('✅ GitHub版本检查测试通过')
    } else {
      console.log('⚠️  GitHub上暂无Release')
    }
    
  } catch (error) {
    console.error('❌ GitHub版本检查测试失败:', error.message)
  }
  
  // 测试前端版本页面
  console.log('\n4. 测试前端版本页面:')
  
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
  
  // 测试Git标签
  console.log('\n5. 测试Git标签:')
  
  try {
    const { stdout: tagOutput } = await execAsync('git tag -l')
    console.log('当前Git标签:')
    console.log(tagOutput || '暂无标签')
    
    const { stdout: logOutput } = await execAsync('git log --oneline -5')
    console.log('最近5次提交:')
    console.log(logOutput)
    
    console.log('✅ Git标签测试通过')
    
  } catch (error) {
    console.error('❌ Git标签测试失败:', error.message)
  }
  
  console.log('\n✅ GitHub版本系统测试完成')
}

// 运行测试
testGitHubVersion() 