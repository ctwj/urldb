/**
 * @name test_upload_plugin
 * @display_name 测试上传插件
 * @version 1.0.0
 * @description 这是一个通过前端上传测试的插件
 * @author Test Author
 * @category utility
 * @dependencies []
 * @permissions []
 * @hooks [load]
 */

function load() {
  console.log('测试上传插件已成功加载');
  return {
    success: true,
    message: '测试上传插件加载成功'
  };
}

if (typeof module !== 'undefined' && module.exports) {
  module.exports = { load };
}