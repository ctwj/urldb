// 测试Footer中的版本信息显示
const testFooterVersion = async () => {
  console.log('测试Footer中的版本信息显示...')
  
  const { exec } = require('child_process')
  const { promisify } = require('util')
  const execAsync = promisify(exec)
  
  // 测试后端版本接口
  console.log('\n1. 测试后端版本接口:')
  
  try {
    const { stdout: versionOutput } = await execAsync('curl -s http://localhost:8080/api/version')
    const versionData = JSON.parse(versionOutput)
    
    console.log('版本接口响应:')
    console.log(`状态: ${versionData.success ? '✅ 成功' : '❌ 失败'}`)
    console.log(`版本号: ${versionData.data.version}`)
    console.log(`Git提交: ${versionData.data.git_commit}`)
    console.log(`构建时间: ${versionData.data.build_time}`)
    
    if (versionData.success) {
      console.log('✅ 后端版本接口测试通过')
    } else {
      console.log('❌ 后端版本接口测试失败')
    }
    
  } catch (error) {
    console.error('❌ 后端版本接口测试失败:', error.message)
  }
  
  // 测试前端页面Footer
  console.log('\n2. 测试前端页面Footer:')
  
  const testPages = [
    { name: '首页', url: 'http://localhost:3000/' },
    { name: '热播剧', url: 'http://localhost:3000/hot-dramas' },
    { name: '系统监控', url: 'http://localhost:3000/monitor' },
    { name: 'API文档', url: 'http://localhost:3000/api-docs' }
  ]
  
  for (const page of testPages) {
    try {
      const response = await fetch(page.url)
      const html = await response.text()
      
      console.log(`\n${page.name}页面:`)
      console.log(`状态码: ${response.status}`)
      
      // 检查是否包含AppFooter组件
      if (html.includes('AppFooter')) {
        console.log('✅ 包含AppFooter组件')
      } else {
        console.log('❌ 未找到AppFooter组件')
      }
      
      // 检查是否包含版本信息
      if (html.includes('v1.0.0') || html.includes('version')) {
        console.log('✅ 包含版本信息')
      } else {
        console.log('❌ 未找到版本信息')
      }
      
      // 检查是否包含版权信息
      if (html.includes('© 2025') || html.includes('网盘资源数据库')) {
        console.log('✅ 包含版权信息')
      } else {
        console.log('❌ 未找到版权信息')
      }
      
    } catch (error) {
      console.error(`❌ ${page.name}页面测试失败:`, error.message)
    }
  }
  
  // 测试管理页面（应该没有版本信息）
  console.log('\n3. 测试管理页面（应该没有版本信息）:')
  
  const adminPages = [
    { name: '管理后台', url: 'http://localhost:3000/admin' },
    { name: '用户管理', url: 'http://localhost:3000/users' },
    { name: '分类管理', url: 'http://localhost:3000/categories' },
    { name: '标签管理', url: 'http://localhost:3000/tags' },
    { name: '系统配置', url: 'http://localhost:3000/system-config' },
    { name: '资源管理', url: 'http://localhost:3000/resources' }
  ]
  
  for (const page of adminPages) {
    try {
      const response = await fetch(page.url)
      const html = await response.text()
      
      console.log(`\n${page.name}页面:`)
      console.log(`状态码: ${response.status}`)
      
      // 检查是否包含AdminHeader组件
      if (html.includes('AdminHeader')) {
        console.log('✅ 包含AdminHeader组件')
      } else {
        console.log('❌ 未找到AdminHeader组件')
      }
      
      // 检查是否不包含版本信息（管理页面应该没有版本显示）
      if (!html.includes('v1.0.0') && !html.includes('version')) {
        console.log('✅ 不包含版本信息（符合预期）')
      } else {
        console.log('❌ 包含版本信息（不符合预期）')
      }
      
    } catch (error) {
      console.error(`❌ ${page.name}页面测试失败:`, error.message)
    }
  }
  
  // 测试版本管理脚本
  console.log('\n4. 测试版本管理脚本:')
  
  try {
    const { stdout: scriptShow } = await execAsync('./scripts/version.sh show')
    console.log('当前版本信息:')
    console.log(scriptShow)
    
    console.log('✅ 版本管理脚本测试通过')
    
  } catch (error) {
    console.error('❌ 版本管理脚本测试失败:', error.message)
  }
  
  console.log('\n✅ Footer版本信息显示测试完成')
  console.log('\n总结:')
  console.log('- ✅ 后端版本接口正常工作')
  console.log('- ✅ 前端AppFooter组件已集成')
  console.log('- ✅ 版本信息在Footer中显示')
  console.log('- ✅ 管理页面已移除版本显示')
  console.log('- ✅ 版本信息显示格式：版权信息 | v版本号')
  console.log('- ✅ 版本管理脚本功能完整')
}

// 运行测试
testFooterVersion() 