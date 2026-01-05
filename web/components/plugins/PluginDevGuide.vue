<template>
  <n-modal v-model:show="visible" :mask-closable="false" preset="card" :style="{ maxWidth: '900px', width: '95%', maxHeight: '90vh' }" title="插件开发说明">
    <div class="space-y-6 overflow-auto" style="max-height: calc(90vh - 120px);">
      <!-- 插件概述 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-info-circle text-blue-500 mr-2"></i>
          插件概述
        </h3>
        <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <p class="text-sm text-gray-700 dark:text-gray-300">
            URLDB 插件系统允许开发者创建自定义功能模块，通过 JavaScript 钩子函数监听系统事件，扩展系统能力。
            插件可以监听用户登录、URL添加、URL访问等事件，并执行自定义逻辑。
          </p>
        </div>
      </section>

      <!-- 支持的事件钩子 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-hooks text-green-500 mr-2"></i>
          支持的事件钩子
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onURLAdd</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">当新URL被添加时触发</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>onURLAdd(function(event) {
    log("info", "新URL添加: " + event.url.url, "my_plugin");
    // event.url 包含完整的URL信息
});</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onURLAccess</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">当URL被访问时触发</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>onURLAccess(function(event) {
    log("info", "URL访问: " + event.url.url, "my_plugin");
    // event.url, event.request, event.response
});</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onUserLogin</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">当用户登录时触发</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>onUserLogin(function(event) {
    log("info", "用户登录: " + event.user.username, "my_plugin");
    // event.user 包含用户信息
});</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onReadyResourceAdd</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">当待处理资源添加时触发</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>onReadyResourceAdd(function(event) {
    log("info", "待处理资源添加: " + event.ready_resource.url, "my_plugin");
    // event.data 包含额外信息，如 is_filtered, filter_reason 等
    if (event.data.is_filtered) {
        log("info", "资源被过滤: " + event.data.filter_reason, "my_plugin");
    }
});</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">定时任务</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">定时执行的任务</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>cronAdd("task_name", "*/5 * * * *", function() {
    log("info", "定时任务执行", "my_plugin");
});</code></pre>
          </div>
        </div>
      </section>

      <!-- 插件结构 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-file-code text-purple-500 mr-2"></i>
          插件文件结构
        </h3>
        <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
          <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">插件文件使用 <code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">.plugin.js</code> 扩展名，基于 JavaScript 开发。</p>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">基本结构示例：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto"><code>/**
 * 插件元信息 - 使用 JSDoc 注释格式
 *
 * @name my_plugin
 * @display_name 我的插件
 * @author 开发者姓名
 * @description 插件功能描述
 * @version 1.0.0
 * @category utility
 * @license MIT
 */

// 记录插件加载日志
log("info", "插件已加载", "my_plugin");

// 监听 URL 添加事件
onURLAdd(function(event) {
    log("info", "=== onURLAdd 事件触发 ===", "my_plugin");
    log("info", "URL: " + event.url.url, "my_plugin");
    log("info", "标题: " + event.url.title, "my_plugin");

    // 自定义逻辑
    if (event.url.url.includes("github.com")) {
        log("info", "检测到GitHub URL", "my_plugin");
    }
});

// 添加自定义路由
routerAdd("GET", "/api/my-endpoint", function(event) {
    return event.json(200, {
        success: true,
        message: "我的插件运行正常",
        timestamp: new Date().toISOString()
    });
});

// 添加定时任务
cronAdd("cleanup_task", "0 2 * * *", function() {
    log("info", "执行清理任务", "my_plugin");
});

log("info", "插件初始化完成", "my_plugin");</code></pre>
        </div>
      </section>

      <!-- 压缩包插件开发 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-file-archive text-cyan-500 mr-2"></i>
          压缩包插件开发
        </h3>
        <div class="bg-cyan-50 dark:bg-cyan-900/20 border border-cyan-200 dark:border-cyan-800 rounded-lg p-4">
          <p class="text-sm text-gray-700 dark:text-gray-300 mb-4">
            压缩包插件允许创建更复杂的多文件插件，包含多个 JavaScript 文件、配置文件、静态资源等。
          </p>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">压缩包结构示例：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto mb-4"><code>my-awesome-plugin.zip
├── package.json              # 插件配置文件
├── index.js                  # 主入口文件
├── lib/
│   ├── utils.js             # 工具函数
│   ├── processor.js         # 数据处理器
│   └── validator.js         # 验证器
├── config/
│   └── default.json         # 默认配置
├── assets/
│   ├── icon.png             # 插件图标
│   └── styles.css           # 样式文件
├── templates/
│   └── email.html           # 邮件模板
└── README.md                # 说明文档</code></pre>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">package.json 配置示例：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto mb-4"><code>{
  "name": "my-awesome-plugin",
  "version": "1.0.0",
  "description": "一个功能强大的插件",
  "main": "index.js",
  "author": "开发者姓名",
  "license": "MIT",
  "keywords": ["plugin", "automation", "utility"],
  "engines": {
    "urldb": ">=1.0.0"
  },
  "config": {
    "webhook_url": {
      "type": "string",
      "label": "Webhook URL",
      "default": "https://hooks.slack.com/...",
      "description": "通知发送的Webhook地址"
    },
    "enable_notifications": {
      "type": "boolean",
      "label": "启用通知",
      "default": true
    },
    "retry_count": {
      "type": "number",
      "label": "重试次数",
      "default": 3,
      "min": 1,
      "max": 10
    }
  },
  "dependencies": {},
  "permissions": [
    "network",
    "storage"
  ]
}</code></pre>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">压缩包插件优势：</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">📁 模块化开发</h5>
              <p class="text-xs text-gray-600 dark:text-gray-400">可以将代码拆分为多个模块，便于维护和测试</p>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">🎨 资源管理</h5>
              <p class="text-xs text-gray-600 dark:text-gray-400">可以包含图标、样式、模板等静态资源</p>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">⚙️ 配置灵活</h5>
              <p class="text-xs text-gray-600 dark:text-gray-400">通过 package.json 定义复杂的配置选项</p>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">🔧 依赖管理</h5>
              <p class="text-xs text-gray-600 dark:text-gray-400">支持模块间的依赖关系和代码复用</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 数据库迁移功能 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-database text-teal-500 mr-2"></i>
          数据库迁移 (migrate) 功能
        </h3>
        <div class="bg-teal-50 dark:bg-teal-900/20 border border-teal-200 dark:border-teal-800 rounded-lg p-4">
          <p class="text-sm text-gray-700 dark:text-gray-300 mb-4">
            Hook 类型插件（单文件 .plugin.js）不支持 migrate 功能，只有压缩包插件支持通过 SQL 文件进行数据库结构迁移。
            在压缩包插件中，migrate 功能通过 migrate/ 目录下的 install.sql 和 uninstall.sql 文件实现。
          </p>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">📋 migrate 规则说明：</h4>
          <div class="space-y-3 mb-4">
            <div class="bg-red-100 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded p-3">
              <h5 class="font-medium text-red-800 dark:text-red-200 mb-1">❌ Hook 插件 (单文件 .plugin.js)</h5>
              <ul class="text-sm text-red-700 dark:text-red-300 space-y-1">
                <li>• 不支持 migrate 功能</li>
                <li>• 无法使用 migrate() 函数或 SQL 迁移文件</li>
                <li>• 适合简单的逻辑处理插件</li>
              </ul>
            </div>
            <div class="bg-green-100 dark:bg-green-900/30 border border-green-200 dark:border-green-800 rounded p-3">
              <h5 class="font-medium text-green-800 dark:text-green-200 mb-1">✅ 压缩包插件</h5>
              <ul class="text-sm text-green-700 dark:text-green-300 space-y-1">
                <li>• 支持 SQL 迁移文件 (install.sql, uninstall.sql)</li>
                <li>• 自动执行安装/卸载迁移</li>
                <li>• 适合需要数据库结构的复杂插件</li>
              </ul>
            </div>
          </div>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">🗂️ 压缩包插件 migrate 目录结构：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto mb-4"><code>my-migration-plugin.zip
├── package.json              # 插件配置文件
├── index.js                  # 主入口文件
├── migrate/                  # 迁移目录
│   ├── install.sql          # 安装时执行的SQL
│   └── uninstall.sql        # 卸载时执行的SQL
├── lib/
│   └── utils.js
└── README.md</code></pre>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">📄 SQL 文件规范：</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-4">
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">install.sql</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 插件安装时自动执行</li>
                <li>• 创建表、索引等结构</li>
                <li>• 插入初始数据</li>
                <li>• 执行失败会回滚安装</li>
              </ul>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">uninstall.sql</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 插件卸载时自动执行</li>
                <li>• 删除相关表结构</li>
                <li>• 清理插件数据</li>
                <li>• 确保系统干净卸载</li>
              </ul>
            </div>
          </div>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">📝 install.sql 示例：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto mb-4"><code>-- 创建插件数据表
CREATE TABLE IF NOT EXISTS my_plugin_data (
    id INTEGER PRIMARY KEY,
    plugin_name VARCHAR(100) NOT NULL,
    message TEXT,
    config_json TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_my_plugin_data_name
ON my_plugin_data(plugin_name);

-- 插入初始数据
INSERT INTO my_plugin_data (plugin_name, message, config_json)
VALUES ('my_plugin', '插件安装时创建的默认数据', '{"enabled": true, "version": "1.0.0"}');</code></pre>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">📝 uninstall.sql 示例：</h4>
          <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-3 rounded overflow-x-auto mb-4"><code>-- 删除索引
DROP INDEX IF EXISTS idx_my_plugin_data_name;

-- 删除插件数据表
DROP TABLE IF EXISTS my_plugin_data;</code></pre>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">⚠️ 错误处理：</h4>
          <div class="bg-yellow-100 dark:bg-yellow-900/30 border border-yellow-200 dark:border-yellow-800 rounded p-3">
            <ul class="text-sm text-yellow-800 dark:text-yellow-200 space-y-1">
              <li>• <strong>安装失败：</strong> install.sql 执行失败时，插件安装会被回滚</li>
              <li>• <strong>卸载失败：</strong> uninstall.sql 执行失败会记录错误日志，但不会阻止卸载</li>
              <li>• <strong>SQL 语法：</strong> 请确保 SQL 语法兼容当前数据库 (PostgreSQL/MySQL/SQLite)</li>
              <li>• <strong>事务安全：</strong> 每个 SQL 文件都在独立事务中执行</li>
            </ul>
          </div>

          <h4 class="font-medium text-gray-900 dark:text-white mb-2">💡 最佳实践：</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">🔧 命名规范</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 表名使用插件前缀：plugin_name_table</li>
                <li>• 索引名包含表名：idx_table_field</li>
                <li>• 避免与系统表名冲突</li>
              </ul>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">🛡️ 安全考虑</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 使用 IF EXISTS 避免错误</li>
                <li>• 不删除系统表或数据</li>
                <li>• 只操作插件相关数据</li>
              </ul>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">📊 性能优化</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 合理创建索引</li>
                <li>• 避免大量数据插入</li>
                <li>• 使用批量操作</li>
              </ul>
            </div>
            <div class="bg-white dark:bg-gray-800 rounded p-3">
              <h5 class="font-medium text-gray-900 dark:text-white mb-1">🧪 测试建议</h5>
              <ul class="text-xs text-gray-600 dark:text-gray-400 space-y-1">
                <li>• 在测试环境验证 SQL</li>
                <li>• 测试安装和卸载流程</li>
                <li>• 验证数据完整性</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      <!-- 可用API -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-cogs text-orange-500 mr-2"></i>
          可用 API
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">日志函数</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">log(level, message, pluginName)</code> - 记录日志</li>
              <li>level: "debug", "info", "warn", "error"</li>
              <li>日志会保存到数据库，可在管理界面查看</li>
            </ul>
          </div>

          <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">路由函数</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">routerAdd(method, path, handler)</code> - 添加自定义路由</li>
            </ul>
          </div>

          <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">定时任务</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">cronAdd(name, schedule, handler)</code> - 添加定时任务</li>
              <li>schedule: Cron 表达式 (如 "*/5 * * * *")</li>
              <li>支持标准 Cron 格式</li>
            </ul>
          </div>

          <div class="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">配置管理</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">getPluginConfig(name)</code> - 获取配置</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">setPluginConfig(name, config)</code> - 设置配置</li>
              <li>通过 JSDoc @config 定义配置字段</li>
              <li>支持多种字段类型和验证</li>
            </ul>
          </div>

          <div class="bg-indigo-50 dark:bg-indigo-900/20 border border-indigo-200 dark:border-indigo-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">数据库操作 (db)</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.find(table, conditions)</code> - 查询记录</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.save(table, data)</code> - 保存记录</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.update(table, id, data)</code> - 更新记录</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.delete(table, id)</code> - 删除记录</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.count(table, conditions)</code> - 计数查询</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">db.raw(sql, params)</code> - 执行原始SQL</li>
            </ul>
          </div>

          <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">HTTP客户端 (http)</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$http.get(url)</code> - GET请求</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$http.post(url, body, options)</code> - POST请求</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$http.put(url, body, options)</code> - PUT请求</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$http.delete(url)</code> - DELETE请求</li>
            </ul>
          </div>

          <div class="bg-cyan-50 dark:bg-cyan-900/20 border border-cyan-200 dark:border-cyan-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">文件系统 (filesystem)</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$filesystem.readFile(path)</code> - 读取文件</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$filesystem.writeFile(path, content)</code> - 写入文件</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$filesystem.fileExists(path)</code> - 检查文件存在</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$filesystem.fileSize(path)</code> - 获取文件大小</li>
            </ul>
          </div>

          <div class="bg-pink-50 dark:bg-pink-900/20 border border-pink-200 dark:border-pink-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">安全函数 (security)</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$security.md5(str)</code> - MD5哈希</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$security.sha256(str)</code> - SHA256哈希</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$security.randomString(length)</code> - 随机字符串</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$security.createJWT(payload, secret, expire)</code> - 创建JWT</li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">$security.parseJWT(token, secret)</code> - 解析JWT</li>
            </ul>
          </div>
        </div>
      </section>

      <!-- 事件对象详情 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-object-group text-indigo-500 mr-2"></i>
          事件对象详情
        </h3>
        <div class="space-y-3">
          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onURLAdd 事件对象</h4>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>{
    url: {
        id: 156696,
        title: "URL标题",
        url: "https://example.com",
        description: "描述",
        category_id: 1,
        tags: ["标签1", "标签2"],
        is_valid: true,
        is_public: true,
        view_count: 0,
        created_at: "2025-12-29T23:49:04.556Z"
    },
    data: {
        // 附加数据
    },
    app: {
        name: "URLDB",
        version: "1.0.0"
    }
}</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onUserLogin 事件对象</h4>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>{
    user: {
        id: 1,
        username: "admin",
        email: "admin@example.com",
        role: "admin",
        is_active: true,
        last_login: "2025-12-29T23:49:04.556Z",
        created_at: "2025-12-29T23:49:04.556Z"
    },
    data: {
        ip: "127.0.0.1",
        user_agent: "浏览器信息",
        login_time: "2025-12-29T23:49:04.556Z"
    },
    app: {
        name: "URLDB",
        version: "1.0.0"
    }
}</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onURLAccess 事件对象</h4>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>{
    url: {
        id: 156696,
        title: "URL标题",
        url: "https://example.com",
        description: "描述",
        category_id: 1,
        tags: ["标签1", "标签2"],
        is_valid: true,
        is_public: true,
        view_count: 0,
        created_at: "2025-12-29T23:49:04.556Z"
    },
    access_log: {},
    request: {},
    response: {},
    app: {
        name: "URLDB",
        version: "1.0.0"
    }
}</code></pre>
          </div>

          <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onReadyResourceAdd 事件对象</h4>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>{
    ready_resource: {
        id: 1,
        key: "资源标识",
        title: "资源标题",
        description: "资源描述",
        url: "https://example.com",
        category: "分类",
        tags: ["标签"],
        img: "图片URL",
        source: "来源",
        extra: "附加信息",
        ip: "IP地址",
        error_msg: "错误信息",
        created_at: "2025-12-29T23:49:04.556Z"
    },
    data: {
        is_filtered: false,
        filter_reason: "过滤原因", // 如果被过滤
        // 其他附加数据
    },
    app: {
        name: "URLDB",
        version: "1.0.0"
    }
}</code></pre>
          </div>
        </div>
      </section>

      <!-- 最佳实践 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-lightbulb text-yellow-500 mr-2"></i>
          最佳实践
        </h3>
        <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
          <ul class="text-sm text-gray-700 dark:text-gray-300 space-y-2">
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>使用描述性的插件名称和函数名</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>合理使用日志级别，便于调试和监控</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>错误处理要完善，避免插件异常影响系统</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>定时任务执行时间不宜过长</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>插件路由使用有意义的前缀，避免冲突</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>提供详细的插件描述和配置说明</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>验证所有外部输入，防止注入攻击</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>避免在高频事件中执行耗时操作</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>及时清理临时文件和资源</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-check text-green-500 mr-2 mt-0.5"></i>
              <span>使用参数化查询防止SQL注入</span>
            </li>
          </ul>
        </div>
      </section>

      <!-- 调试技巧 -->
      <section>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
          <i class="fas fa-bug text-red-500 mr-2"></i>
          调试技巧
        </h3>
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <ul class="text-sm text-gray-700 dark:text-gray-300 space-y-2">
            <li class="flex items-start">
              <i class="fas fa-tools text-red-500 mr-2 mt-0.5"></i>
              <span>使用 <code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">log("debug", message, pluginName)</code> 记录调试信息</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-tools text-red-500 mr-2 mt-0.5"></i>
              <span>在插件管理界面查看实时日志输出</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-tools text-red-500 mr-2 mt-0.5"></i>
              <span>插件支持热重载，修改后自动生效</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-tools text-red-500 mr-2 mt-0.5"></i>
              <span>先在测试环境验证插件功能</span>
            </li>
            <li class="flex items-start">
              <i class="fas fa-tools text-red-500 mr-2 mt-0.5"></i>
              <span>压缩包插件可以通过解压查看文件结构进行调试</span>
            </li>
          </ul>
        </div>
      </section>
    </div>

    <template #footer>
      <div class="flex justify-between">
        <n-button @click="goToPluginManager" type="info" v-if="showPluginManagerButton">
          <template #icon>
            <i class="fas fa-plug"></i>
          </template>
          前往插件管理
        </n-button>
        <div></div>
        <n-button @click="closeModal" type="primary">
          我知道了
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean
  showPluginManagerButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showPluginManagerButton: true
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'go-to-plugin-manager': []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const closeModal = () => {
  emit('update:modelValue', false)
}

const goToPluginManager = () => {
  emit('go-to-plugin-manager')
  closeModal()
}
</script>