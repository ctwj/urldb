/// <reference path="../pb_data/types.d.ts" />

//  * // 非必填项，使用 @optional
//  * // 默认值，使用 @default
//  * // @field {类型} 名字 表单名字 placeholder [@optional] [@default]


/**
 * config_demo 钩子
 * 创建时间: 2025-12-25 23:23:19
 *
 * @name config_demo
 * @display_name 配置演示插件
 * @author URLDB开发团队
 * @description 演示插件配置功能的示例插件，包含Webhook通知、日志级别设置等配置选项
 * @version 1.0.1
 * @category demo
 * @license MIT
 *
 * @config
 * @field {string} webhook_url Webhook URL "通知发送的Webhook地址" @default "https://hooks.slack.com/services/YOUR/DEFAULT/WEBHOOK"
 * @field {boolean} enable_notification 启用通知 "是否启用通知功能" @default true
 * @field {number} retry_count 重试次数 "通知失败时的重试次数" @default 3
 * @field {select} log_level 日志级别 "日志输出级别" ["debug", "info", "warn", "error"] @default "info"
 * @field {text} custom_message 自定义消息 "自定义通知消息内容" @optional @default "这是来自 config_demo 插件的默认消息"
 * @config
 */

// 提取的配置处理函数
function processConfigDemo() {
    try {
        // 获取插件配置
        const config = getPluginConfig("config_demo");

        // 最简化处理，避免所有console.log
        if (config) {
            return {
                success: true,
                config: config,
                timestamp: new Date().toISOString()
            };
        } else {
            return {
                success: false,
                error: "未找到插件配置",
                timestamp: new Date().toISOString()
            };
        }
    } catch (error) {
        return {
            success: false,
            error: error.message,
            timestamp: new Date().toISOString()
        };
    }
}

// 示例：监听 URL 添加事件
onURLAdd(function(event) {
    log("info", "=== config_demo onURLAdd 事件触发 ===", "config_demo");
    log("info", "URL ID: " + event.url.id, "config_demo");
    log("info", "URL Title: " + event.url.title, "config_demo");
    log("info", "URL: " + event.url.url, "config_demo");

    // 在这里添加你的自定义逻辑
    // 例如：自动分类、标签提取、通知等
    if (event.url.url && event.url.url.includes("github.com")) {
        log("info", "检测到GitHub URL，建议分类为: 开发工具", "config_demo");
    }

    log("info", "=== config_demo onURLAdd 事件处理完成 ===", "config_demo");
});

// 示例：监听用户登录事件
onUserLogin(function(event) {
    log("info", "=== config_demo onUserLogin 事件触发 ===", "config_demo");
    log("info", "用户ID: " + event.user.id, "config_demo");
    log("info", "用户名: " + event.user.username, "config_demo");
    log("info", "邮箱: " + event.user.email, "config_demo");

    // 在这里添加登录后处理逻辑
    // 例如：日志记录、欢迎消息、权限检查等
    log("info", "欢迎 " + event.user.username + " 登录系统！", "config_demo");
    log("info", "=== config_demo onUserLogin 事件处理完成 ===", "config_demo");
});

// 示例：添加自定义路由 - 获取配置信息
routerAdd("GET", "/api/config-demo", (e) => {
    const result = processConfigDemo();

    // 尝试获取随机资源作为附加信息
    let randomResource = null;
    try {
        const randomResult = db.raw(`
            SELECT id, name, url, category
            FROM resources
            WHERE deleted_at IS NULL OR deleted_at = '' OR deleted_at IS NULL
            ORDER BY RANDOM()
            LIMIT 1
        `);

        if (randomResult && randomResult.length > 0) {
            randomResource = {
                id: randomResult[0].id,
                name: randomResult[0].name || randomResult[0].url || "未命名资源",
                url: randomResult[0].url,
                category: randomResult[0].category || "未分类"
            };
        }
    } catch (error) {
        // 静默处理错误，不影响主要功能
    }

    return e.json(200, {
        message: "来自 config_demo 插件的自定义 API",
        data: result,
        random_resource: randomResource,
        timestamp: new Date().toISOString()
    });
});

// 添加新的路由 - 手动触发配置处理
routerAdd("POST", "/api/config-demo/refresh", (e) => {
    log("info", "手动触发配置处理", "config_demo");

    const result = processConfigDemo();

    return e.json(200, {
        message: "配置处理完成",
        data: result,
        timestamp: new Date().toISOString()
    });
});

// 临时解决方案：在现有refresh API中添加随机资源功能
routerAdd("GET", "/api/config-demo/refresh", (e) => {
    log("info", "获取随机资源（临时实现）", "config_demo");

    try {
        // 尝试获取随机资源
        let randomResource = null;

        try {
            const randomResult = db.raw(`
                SELECT id, name, url, category
                FROM resources
                WHERE deleted_at IS NULL OR deleted_at = '' OR deleted_at IS NULL
                ORDER BY RANDOM()
                LIMIT 1
            `);

            if (randomResult && randomResult.length > 0) {
                randomResource = randomResult[0];
            }
        } catch (error) {
            log("warn", "随机查询失败: " + error.message, "config_demo");

            // 备选方案：获取第一个资源
            try {
                const allResources = db.find("resources", {});
                if (allResources && allResources.length > 0) {
                    const activeResources = allResources.filter(resource =>
                        !resource.deleted_at || resource.deleted_at === ''
                    );
                    if (activeResources.length > 0) {
                        randomResource = activeResources[0];
                    }
                }
            } catch (fallbackError) {
                log("warn", "备选方案也失败: " + fallbackError.message, "config_demo");
            }
        }

        if (randomResource) {
            return e.json(200, {
                success: true,
                message: "成功获取随机资源",
                data: {
                    id: randomResource.id,
                    name: randomResource.name || randomResource.url || "未命名资源",
                    url: randomResource.url,
                    category: randomResource.category || "未分类",
                    selected_at: new Date().toISOString()
                },
                timestamp: new Date().toISOString()
            });
        } else {
            return e.json(404, {
                success: false,
                message: "未找到任何资源记录",
                suggestion: "请确保数据库中有资源记录",
                timestamp: new Date().toISOString()
            });
        }

    } catch (error) {
        return e.json(500, {
            success: false,
            error: error.message,
            message: "获取随机资源时发生错误",
            timestamp: new Date().toISOString()
        });
    }
});

// 添加新的路由 - 获取配置摘要
routerAdd("GET", "/api/config-demo/summary", (e) => {
    const config = getPluginConfig("config_demo");

    if (!config) {
        return e.json(404, {
            error: "未找到插件配置",
            timestamp: new Date().toISOString()
        });
    }

    return e.json(200, {
        plugin_name: "config_demo",
        display_name: "配置演示插件",
        version: "1.0.1",
        webhook_configured: config.webhook_url && config.webhook_url !== "https://hooks.slack.com/services/YOUR/DEFAULT/WEBHOOK",
        notification_enabled: config.enable_notification || false,
        log_level: config.log_level || "info",
        retry_count: config.retry_count || 0,
        timestamp: new Date().toISOString()
    });
});

// 添加新的路由 - 测试所有绑定函数
routerAdd("GET", "/api/bindings-test", (e) => {
    log("info", "=== 开始测试所有绑定函数 ===", "config_demo");

    const testResults = {
        timestamp: new Date().toISOString(),
        tests: {},
        summary: { passed: 0, failed: 0, total: 0 }
    };

    try {
        // 1. 测试安全函数 (security)
        log("info", "测试安全函数...", "config_demo");
        try {
            const securityTests = {
                md5: $security.md5("test"),
                sha256: $security.sha256("test"),
                sha512: $security.sha512("test"),
                randomString: $security.randomString(10),
                jwt: null,
                encryption: null
            };

            // 测试 JWT
            const jwtPayload = { user: "test", exp: Math.floor(Date.now() / 1000) + 3600 };
            const jwtToken = $security.createJWT(jwtPayload, "test-secret", 3600);
            securityTests.jwt = {
                created: jwtToken,
                verified: $security.parseJWT(jwtToken, "test-secret")
            };

            // 测试加密
            const encrypted = $security.encrypt("hello world", "test-key");
            const decrypted = $security.decrypt(encrypted, "test-key");
            securityTests.encryption = {
                encrypted: encrypted,
                decrypted: decrypted
            };

            testResults.tests.security = { status: "success", data: securityTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.security = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 2. 测试操作系统函数 (os)
        log("info", "测试操作系统函数...", "config_demo");
        try {
            const osTests = {
                env: $os.getenv("PATH") || "PATH_NOT_FOUND",
                platform: os.platform(),
                tempDir: $os.tempDir(),
                args: Array.from($os.args || []).slice(0, 3) // 只显示前3个参数
            };

            testResults.tests.os = { status: "success", data: osTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.os = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 3. 测试文件系统函数 (fs)
        log("info", "测试文件系统函数...", "config_demo");
        try {
            const testContent = "Hello from config_demo plugin test!";
            const testFilePath = $os.tempDir() + "/urldb_test_" + Date.now() + ".txt";

            // 测试写入文件
            $filesystem.writeFile(testFilePath, testContent);

            // 测试读取文件
            const readContent = $filesystem.readFile(testFilePath);

            const fsTests = {
                testFile: testFilePath,
                originalContent: testContent,
                readContent: readContent,
                contentMatch: readContent === testContent
            };

            testResults.tests.filesystem = { status: "success", data: fsTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.filesystem = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 4. 测试文件路径函数 (filepath)
        log("info", "测试文件路径函数...", "config_demo");
        try {
            const filepathTests = {
                join: filepath.join("path", "to", "file.txt"),
                base: filepath.base("/path/to/file.txt"),
                dir: filepath.dir("/path/to/file.txt"),
                ext: filepath.ext("/path/to/file.txt"),
                isAbs: filepath.isAbs("/path/to/file.txt"),
                clean: filepath.clean("path//to///./file.txt")
            };

            testResults.tests.filepath = { status: "success", data: filepathTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.filepath = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 5. 测试表单验证函数 (forms)
        log("info", "测试表单验证函数...", "config_demo");
        try {
            const testData = {
                name: "测试用户",
                email: "test@example.com",
                age: 25,
                phone: "1234567890"
            };

            const validationRules = {
                name: { required: true, minLength: 2, maxLength: 50 },
                email: { required: true, isEmail: true },
                age: { required: true, minLength: 18, maxLength: 120 },
                phone: { pattern: "^[0-9]{10,11}$" }
            };

            const validationResult = forms.validate(testData, validationRules);

            const formsTests = {
                testData: testData,
                validationRules: validationRules,
                validationResult: validationResult
            };

            testResults.tests.forms = { status: "success", data: formsTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.forms = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 6. 测试HTTP客户端函数 (http)
        log("info", "测试HTTP客户端函数...", "config_demo");
        try {
            // 测试 GET 请求
            const httpGetResult = http.get("https://httpbin.org/get");

            // 测试 POST 请求
            const httpPostResult = http.post("https://httpbin.org/post", {
                plugin: "config_demo",
                test: true,
                timestamp: new Date().toISOString()
            });

            const httpTests = {
                getRequest: {
                    url: "https://httpbin.org/get",
                    status: httpGetResult.status || httpGetResult.statusCode,
                    success: (httpGetResult.status || httpGetResult.statusCode) >= 200 && (httpGetResult.status || httpGetResult.statusCode) < 300
                },
                postRequest: {
                    url: "https://httpbin.org/post",
                    status: httpPostResult.status || httpPostResult.statusCode,
                    success: (httpPostResult.status || httpPostResult.statusCode) >= 200 && (httpPostResult.status || httpPostResult.statusCode) < 300
                }
            };

            testResults.tests.http = { status: "success", data: httpTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.http = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 7. 测试数据库函数 (db)
        log("info", "测试数据库函数...", "config_demo");
        try {
            const dbTests = {
                // 测试查询resources表
                resourcesQuery: null,
                resourcesCount: null,
                testRaw: null,
                randomResource: null
            };

            // 尝试查询resources表
            try {
                const resourcesResult = db.find("resources", {});
                dbTests.resourcesQuery = {
                    success: true,
                    count: Array.isArray(resourcesResult) ? resourcesResult.length : 0,
                    sample: resourcesResult ? resourcesResult.slice(0, 3) : []
                };
            } catch (dbError) {
                dbTests.resourcesQuery = { success: false, error: dbError.message };
            }

            // 尝试resources计数查询
            try {
                const countResult = db.count("resources", {});
                dbTests.resourcesCount = { success: true, count: countResult };
            } catch (countError) {
                dbTests.resourcesCount = { success: false, error: countError.message };
            }

            // 尝试原始SQL查询 - 获取随机资源
            try {
                const randomResult = db.raw(`
                    SELECT id, title, url, category
                    FROM resources
                    WHERE deleted_at IS NULL OR deleted_at = '' OR deleted_at IS NULL
                    ORDER BY RANDOM()
                    LIMIT 1
                `);
                dbTests.randomResource = {
                    success: true,
                    result: randomResult && randomResult.length > 0 ? randomResult[0] : null
                };
            } catch (randomError) {
                dbTests.randomResource = { success: false, error: randomError.message };
            }

            // 尝试原始SQL查询
            try {
                const rawResult = db.raw("SELECT 1 as test_value");
                dbTests.testRaw = { success: true, result: rawResult };
            } catch (rawError) {
                dbTests.testRaw = { success: false, error: rawError.message };
            }

            testResults.tests.database = { status: "success", data: dbTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.database = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        // 8. 测试邮件函数 (mails)
        log("info", "测试邮件函数...", "config_demo");
        try {
            const mailResult = mails.send(
                "test@example.com",
                "测试邮件 - config_demo插件",
                "这是一封来自config_demo插件的测试邮件，发送时间: " + new Date().toISOString()
            );

            const mailTests = {
                sendResult: mailResult,
                timestamp: new Date().toISOString()
            };

            testResults.tests.mails = { status: "success", data: mailTests };
            testResults.summary.passed++;
        } catch (error) {
            testResults.tests.mails = { status: "error", error: error.message };
            testResults.summary.failed++;
        }
        testResults.summary.total++;

        log("info", `=== 绑定函数测试完成 === 成功: ${testResults.summary.passed}, 失败: ${testResults.summary.failed}, 总计: ${testResults.summary.total}`, "config_demo");

    } catch (error) {
        log("error", "测试过程中发生严重错误: " + error.message, "config_demo");
        testResults.error = error.message;
    }

    return e.json(200, {
        message: "绑定函数测试完成",
        results: testResults,
        timestamp: new Date().toISOString()
    });
});

// 添加调试路由 - 测试绑定函数可用性
routerAdd("GET", "/api/bindings-test/debug", (e) => {
    try {
        const debugResults = {
            timestamp: new Date().toISOString(),
            availableFunctions: {},
            errors: {}
        };

        // 测试 $security 对象
        try {
            debugResults.availableFunctions.$security = {
                md5: typeof $security.md5,
                sha256: typeof $security.sha256,
                randomString: typeof $security.randomString,
                createJWT: typeof $security.createJWT,
                parseJWT: typeof $security.parseJWT,
                encrypt: typeof $security.encrypt,
                decrypt: typeof $security.decrypt
            };
        } catch (error) {
            debugResults.errors.$security = error.message;
        }

        // 测试 $os 对象
        try {
            debugResults.availableFunctions.$os = {
                getenv: typeof $os.getenv,
                tempDir: typeof $os.tempDir,
                platform: typeof $os.platform,
                args: typeof $os.args
            };
        } catch (error) {
            debugResults.errors.$os = error.message;
        }

        // 测试 $filesystem 对象
        try {
            debugResults.availableFunctions.$filesystem = {
                readFile: typeof $filesystem.readFile,
                writeFile: typeof $filesystem.writeFile,
                fileExists: typeof $filesystem.fileExists,
                fileSize: typeof $filesystem.fileSize
            };
        } catch (error) {
            debugResults.errors.$filesystem = error.message;
        }

        // 测试 $filepath 对象
        try {
            debugResults.availableFunctions.$filepath = {
                join: typeof $filepath.join,
                base: typeof $filepath.base,
                dir: typeof $filepath.dir,
                ext: typeof $filepath.ext
            };
        } catch (error) {
            debugResults.errors.$filepath = error.message;
        }

        // 测试 $forms 对象
        try {
            debugResults.availableFunctions.$forms = {
                validate: typeof $forms.validate,
                create: typeof $forms.create
            };
        } catch (error) {
            debugResults.errors.$forms = error.message;
        }

        // 测试 $db 对象
        try {
            debugResults.availableFunctions.$db = {
                find: typeof $db.find,
                save: typeof $db.save,
                update: typeof $db.update,
                delete: typeof $db.delete,
                count: typeof $db.count,
                raw: typeof $db.raw
            };
        } catch (error) {
            debugResults.errors.$db = error.message;
        }

        // 测试 $http 对象
        try {
            debugResults.availableFunctions.$http = {
                get: typeof $http.get,
                post: typeof $http.post,
                put: typeof $http.put,
                delete: typeof $http.delete,
                send: typeof $http.send
            };
        } catch (error) {
            debugResults.errors.$http = error.message;
        }

        return e.json(200, {
            message: "绑定函数调试信息",
            data: debugResults
        });
    } catch (error) {
        return e.json(500, {
            error: error.message,
            timestamp: new Date().toISOString()
        });
    }
});

// 添加简单的快速测试路由
routerAdd("GET", "/api/bindings-test/simple", (e) => {
    try {
        const simpleResults = {
            timestamp: new Date().toISOString(),
            security: {
                hash: security.hash("quick test"),
                md5: $security.md5("quick test"),
                sha256: $security.sha256("quick test")
            },
            system: {
                platform: os.platform(),
                env_path: $os.getenv("PATH") ? "available" : "not available",
                tempDir: $os.tempDir()
            },
            filesystem: {
                tempDir: $os.tempDir(),
                canWrite: "tested"
            },
            filepath: {
                join: $filepath.join("path", "to", "file.txt"),
                base: $filepath.base("/path/to/file.txt")
            },
            forms: {
                simpleValidation: $forms.validate(
                    { name: "test", email: "test@test.com" },
                    { name: { required: true }, email: { isEmail: true } }
                )
            }
        };

        return e.json(200, {
            message: "简单绑定函数测试",
            data: simpleResults
        });
    } catch (error) {
        return e.json(500, {
            error: error.message,
            timestamp: new Date().toISOString()
        });
    }
});

// 示例：添加定时任务 - 重新启用，每1分钟执行一次
cronAdd("config_demo_task", "*/1 * * * *", () => {
    log("info", "执行定时任务: config_demo - 每1分钟执行一次", "config_demo");

    try {
        // 调用提取的函数
        const result = processConfigDemo();

        // 记录执行结果
        if (result && result.success) {
            log("info", "定时任务执行成功，配置已处理", "config_demo");
        } else {
            log("error", "定时任务执行失败: " + ((result && result.error) || "未知错误"), "config_demo");
        }
    } catch (error) {
        log("error", "定时任务执行异常: " + (error.message || error), "config_demo");
    }
});

// 注意：Hook 类型插件不支持 migrate 功能
// 如需数据库迁移，请使用 migrations 目录或压缩包插件的 migrate 目录

// 添加迁移测试 API
routerAdd("GET", "/api/config-demo/migration-test", (e) => {
    return e.json(200, {
        message: "迁移注册测试 API",
        info: "检查日志以确认迁移注册情况",
        note: "迁移已注册在插件加载时，请查看启动日志",
        timestamp: new Date().toISOString()
    });
});
// 触发热重载: 2026年 1月 3日 星期六 23时55分02秒 CST
