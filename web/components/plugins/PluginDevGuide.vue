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
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">onUserLogin</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">当用户登录时触发</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>onUserLogin(function(event) {
    log("info", "用户登录: " + event.user.username, "my_plugin");
    // event.user 包含用户信息
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
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">定时任务</h4>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">定时执行的任务</p>
            <pre class="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded overflow-x-auto"><code>cron("task_name", "*/5 * * * *", function() {
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
router.get("/api/my-endpoint", function() {
    return {
        success: true,
        message: "我的插件运行正常",
        timestamp: new Date().toISOString()
    };
});

// 添加定时任务
cron("cleanup_task", "0 2 * * *", function() {
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
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">router.get(path, handler)</code></li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">router.post(path, handler)</code></li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">router.put(path, handler)</code></li>
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">router.delete(path, handler)</code></li>
            </ul>
          </div>

          <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">定时任务</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">cron(name, schedule, handler)</code></li>
              <li>schedule: Cron 表达式 (如 "*/5 * * * *")</li>
              <li>支持标准 Cron 格式</li>
            </ul>
          </div>

          <div class="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800 rounded-lg p-4">
            <h4 class="font-medium text-gray-900 dark:text-white mb-2">配置管理</h4>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li><code class="bg-gray-200 dark:bg-gray-700 px-1 rounded">getPluginConfig(name)</code> - 获取配置</li>
              <li>通过 JSDoc @config 定义配置字段</li>
              <li>支持多种字段类型和验证</li>
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