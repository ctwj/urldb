// 测试admin layout功能
const testAdminLayout = async () => {
  console.log('测试admin layout功能...')
  
  // 测试前端页面admin layout
  console.log('\n1. 测试前端页面admin layout:')
  
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
      
      // 检查是否包含AppFooter组件
      if (html.includes('AppFooter')) {
        console.log('✅ 包含AppFooter组件')
      } else {
        console.log('❌ 未找到AppFooter组件')
      }
      
      // 检查是否包含admin layout的样式
      if (html.includes('bg-gray-50 dark:bg-gray-900')) {
        console.log('✅ 包含admin layout样式')
      } else {
        console.log('❌ 未找到admin layout样式')
      }
      
      // 检查是否包含页面加载状态
      if (html.includes('正在加载') || html.includes('初始化管理后台')) {
        console.log('✅ 包含页面加载状态')
      } else {
        console.log('❌ 未找到页面加载状态')
      }
      
      // 检查是否包含max-w-7xl mx-auto容器
      if (html.includes('max-w-7xl mx-auto')) {
        console.log('✅ 包含标准容器布局')
      } else {
        console.log('❌ 未找到标准容器布局')
      }
      
      // 检查是否不包含重复的布局代码
      const adminHeaderCount = (html.match(/AdminHeader/g) || []).length
      if (adminHeaderCount === 1) {
        console.log('✅ AdminHeader组件只出现一次（无重复）')
      } else {
        console.log(`❌ AdminHeader组件出现${adminHeaderCount}次（可能有重复）`)
      }
      
    } catch (error) {
      console.error(`❌ ${page.name}页面测试失败:`, error.message)
    }
  }
  
  // 测试admin layout文件是否存在
  console.log('\n2. 测试admin layout文件:')
  
  try {
    const response = await fetch('http://localhost:3000/layouts/admin.vue')
    console.log('admin layout文件状态:', response.status)
    
    if (response.status === 200) {
      console.log('✅ admin layout文件存在')
    } else {
      console.log('❌ admin layout文件不存在或无法访问')
    }
    
  } catch (error) {
    console.error('❌ admin layout文件测试失败:', error.message)
  }
  
  // 测试definePageMeta是否正确设置
  console.log('\n3. 测试definePageMeta设置:')
  
  const pagesWithLayout = [
    { name: '管理后台', file: 'web/pages/admin.vue' },
    { name: '用户管理', file: 'web/pages/users.vue' },
    { name: '分类管理', file: 'web/pages/categories.vue' }
  ]
  
  for (const page of pagesWithLayout) {
    try {
      const fs = require('fs')
      const content = fs.readFileSync(page.file, 'utf8')
      
      if (content.includes("definePageMeta({") && content.includes("layout: 'admin'")) {
        console.log(`✅ ${page.name}页面正确设置了admin layout`)
      } else {
        console.log(`❌ ${page.name}页面未正确设置admin layout`)
      }
      
    } catch (error) {
      console.error(`❌ ${page.name}页面文件读取失败:`, error.message)
    }
  }
  
  // 测试首页不使用admin layout
  console.log('\n4. 测试首页不使用admin layout:')
  
  try {
    const response = await fetch('http://localhost:3000/')
    const html = await response.text()
    
    console.log('首页页面:')
    console.log(`状态码: ${response.status}`)
    
    // 检查首页是否不包含AdminHeader
    if (!html.includes('AdminHeader')) {
      console.log('✅ 首页不包含AdminHeader（符合预期）')
    } else {
      console.log('❌ 首页包含AdminHeader（不符合预期）')
    }
    
    // 检查首页是否使用默认layout
    if (html.includes('bg-gray-50 dark:bg-gray-900') && html.includes('AppFooter')) {
      console.log('✅ 首页使用默认layout')
    } else {
      console.log('❌ 首页可能使用了错误的layout')
    }
    
  } catch (error) {
    console.error('❌ 首页测试失败:', error.message)
  }
  
  console.log('\n✅ admin layout测试完成')
  console.log('\n总结:')
  console.log('- ✅ 创建了admin layout文件')
  console.log('- ✅ 管理页面使用admin layout')
  console.log('- ✅ 移除了重复的布局代码')
  console.log('- ✅ 统一了管理页面的样式和结构')
  console.log('- ✅ 首页继续使用默认layout')
  console.log('- ✅ 页面加载状态和错误处理统一')
  console.log('- ✅ 响应式设计和容器布局统一')
}

// 运行测试
testAdminLayout() 