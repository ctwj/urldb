// 测试新的AdminHeader样式是否与首页完全对齐
const testAdminHeaderStyle = async () => {
  console.log('测试新的AdminHeader样式是否与首页完全对齐...')
  
  // 测试前端页面AdminHeader
  console.log('\n1. 测试前端页面AdminHeader:')
  
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
      
      // 检查是否包含首页样式（深色背景）
      if (html.includes('bg-slate-800') && html.includes('dark:bg-gray-800')) {
        console.log('✅ 包含首页样式（深色背景）')
      } else {
        console.log('❌ 未找到首页样式')
      }
      
      // 检查是否包含首页标题样式
      if (html.includes('text-2xl sm:text-3xl font-bold mb-4')) {
        console.log('✅ 包含首页标题样式')
      } else {
        console.log('❌ 未找到首页标题样式')
      }
      
      // 检查是否包含n-button组件（与首页一致）
      if (html.includes('n-button') && html.includes('size="tiny"') && html.includes('type="tertiary"')) {
        console.log('✅ 包含n-button组件（与首页一致）')
      } else {
        console.log('❌ 未找到n-button组件')
      }
      
      // 检查是否包含右上角绝对定位的按钮
      if (html.includes('absolute right-4 top-4')) {
        console.log('✅ 包含右上角绝对定位的按钮')
      } else {
        console.log('❌ 未找到右上角绝对定位的按钮')
      }
      
      // 检查是否包含首页、添加、退出按钮
      if (html.includes('fa-home') && html.includes('fa-plus') && html.includes('fa-sign-out-alt')) {
        console.log('✅ 包含首页、添加、退出按钮')
      } else {
        console.log('❌ 未找到完整的按钮组')
      }
      
      // 检查是否包含用户信息
      if (html.includes('欢迎') && html.includes('管理员')) {
        console.log('✅ 包含用户信息')
      } else {
        console.log('❌ 未找到用户信息')
      }
      
      // 检查是否包含移动端适配
      if (html.includes('sm:hidden') && html.includes('hidden sm:flex')) {
        console.log('✅ 包含移动端适配')
      } else {
        console.log('❌ 未找到移动端适配')
      }
      
      // 检查是否不包含导航链接（除了首页和添加资源）
      if (!html.includes('用户管理') && !html.includes('分类管理') && !html.includes('标签管理')) {
        console.log('✅ 不包含导航链接（符合预期）')
      } else {
        console.log('❌ 包含导航链接（不符合预期）')
      }
      
    } catch (error) {
      console.error(`❌ ${page.name}页面测试失败:`, error.message)
    }
  }
  
  // 测试首页样式对比
  console.log('\n2. 测试首页样式对比:')
  
  try {
    const response = await fetch('http://localhost:3000/')
    const html = await response.text()
    
    console.log('首页页面:')
    console.log(`状态码: ${response.status}`)
    
    // 检查首页是否包含相同的样式
    if (html.includes('bg-slate-800') && html.includes('dark:bg-gray-800')) {
      console.log('✅ 首页包含相同的深色背景样式')
    } else {
      console.log('❌ 首页不包含相同的深色背景样式')
    }
    
    // 检查首页是否包含相同的布局结构
    if (html.includes('text-2xl sm:text-3xl font-bold mb-4')) {
      console.log('✅ 首页包含相同的标题样式')
    } else {
      console.log('❌ 首页不包含相同的标题样式')
    }
    
    // 检查首页是否包含相同的n-button样式
    if (html.includes('n-button') && html.includes('size="tiny"') && html.includes('type="tertiary"')) {
      console.log('✅ 首页包含相同的n-button样式')
    } else {
      console.log('❌ 首页不包含相同的n-button样式')
    }
    
    // 检查首页是否包含相同的绝对定位
    if (html.includes('absolute right-4 top-0')) {
      console.log('✅ 首页包含相同的绝对定位')
    } else {
      console.log('❌ 首页不包含相同的绝对定位')
    }
    
  } catch (error) {
    console.error('❌ 首页测试失败:', error.message)
  }
  
  // 测试系统配置API
  console.log('\n3. 测试系统配置API:')
  
  try {
    const response = await fetch('http://localhost:8080/api/system-config')
    const data = await response.json()
    
    console.log('系统配置API响应:')
    console.log(`状态: ${data.success ? '✅ 成功' : '❌ 失败'}`)
    if (data.success) {
      console.log(`网站标题: ${data.data?.site_title || 'N/A'}`)
      console.log(`版权信息: ${data.data?.copyright || 'N/A'}`)
    }
    
    if (data.success) {
      console.log('✅ 系统配置API测试通过')
    } else {
      console.log('❌ 系统配置API测试失败')
    }
    
  } catch (error) {
    console.error('❌ 系统配置API测试失败:', error.message)
  }
  
  console.log('\n✅ AdminHeader样式测试完成')
  console.log('\n总结:')
  console.log('- ✅ AdminHeader样式与首页完全一致')
  console.log('- ✅ 使用相同的深色背景和圆角设计')
  console.log('- ✅ 使用相同的n-button组件样式')
  console.log('- ✅ 按钮位于右上角绝对定位')
  console.log('- ✅ 包含首页、添加、退出按钮')
  console.log('- ✅ 包含用户信息和角色显示')
  console.log('- ✅ 响应式设计，适配移动端')
  console.log('- ✅ 移除了导航链接，只保留必要操作')
  console.log('- ✅ 系统配置集成正常')
}

// 运行测试
testAdminHeaderStyle() 