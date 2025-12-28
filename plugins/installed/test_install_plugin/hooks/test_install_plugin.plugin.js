/**
 * @name test_install_plugin
 * @display_name 测试安装插件
 * @version 1.0.0
 * @description 这是一个用于测试安装功能的插件
 * @author Test Author
 * @category utility
 * @dependencies []
 * @permissions []
 * @hooks [load]
 */

// 插件主函数
function load() {
  console.log('测试安装插件已加载');
  return {
    success: true,
    message: '测试安装插件加载成功'
  };
}

// 导出函数
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { load };
}