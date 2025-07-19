// docsify 配置文件
window.$docsify = {
  name: 'URL数据库管理系统',
  repo: 'https://github.com/ctwj/urldb',
  loadSidebar: true,
  subMaxLevel: 3,
  auto2top: true,
  search: {
    maxAge: 86400000,
    paths: 'auto',
    placeholder: '搜索文档...',
    noData: '找不到结果',
    depth: 6
  },
  copyCode: {
    buttonText: '复制',
    errorText: '错误',
    successText: '已复制'
  },
  pagination: {
    previousText: '上一页',
    nextText: '下一页',
    crossChapter: true,
    crossChapterText: true,
  },
  plugins: [
    function(hook, vm) {
      hook.beforeEach(function (html) {
        // 添加页面标题
        var url = '#' + vm.route.path;
        var title = vm.route.path === '/' ? '首页' : vm.route.path.replace('/', '');
        return html + '\n\n---\n\n' + 
               '<div style="text-align: center; color: #666; font-size: 14px;">' +
               '最后更新: ' + new Date().toLocaleDateString('zh-CN') +
               '</div>';
      });
    }
  ]
}; 