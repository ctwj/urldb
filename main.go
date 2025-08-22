package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/handlers"
	"github.com/ctwj/urldb/middleware"
	"github.com/ctwj/urldb/scheduler"
	"github.com/ctwj/urldb/services"
	"github.com/ctwj/urldb/task"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 检查命令行参数
	if len(os.Args) > 1 && os.Args[1] == "version" {
		versionInfo := utils.GetVersionInfo()
		fmt.Printf("版本: v%s\n", versionInfo.Version)
		fmt.Printf("构建时间: %s\n", versionInfo.BuildTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("Git提交: %s\n", versionInfo.GitCommit)
		fmt.Printf("Git分支: %s\n", versionInfo.GitBranch)
		fmt.Printf("Go版本: %s\n", versionInfo.GoVersion)
		fmt.Printf("平台: %s/%s\n", versionInfo.Platform, versionInfo.Arch)
		return
	}

	// 初始化日志系统
	if err := utils.InitLogger(nil); err != nil {
		log.Fatal("初始化日志系统失败:", err)
	}
	defer utils.GetLogger().Close()

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		utils.Info("未找到.env文件，使用默认配置")
	}

	// 初始化时区设置
	utils.InitTimezone()

	// 设置Gin运行模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		// 如果没有设置GIN_MODE，根据环境判断
		if os.Getenv("ENV") == "production" {
			gin.SetMode(gin.ReleaseMode)
			utils.Info("设置Gin为Release模式")
		} else {
			gin.SetMode(gin.DebugMode)
			utils.Info("设置Gin为Debug模式")
		}
	} else {
		// 如果已经设置了GIN_MODE，根据值设置模式
		switch ginMode {
		case "release":
			gin.SetMode(gin.ReleaseMode)
			utils.Info("设置Gin为Release模式 (来自环境变量)")
		case "debug":
			gin.SetMode(gin.DebugMode)
			utils.Info("设置Gin为Debug模式 (来自环境变量)")
		case "test":
			gin.SetMode(gin.TestMode)
			utils.Info("设置Gin为Test模式 (来自环境变量)")
		default:
			gin.SetMode(gin.DebugMode)
			utils.Info("未知的GIN_MODE值: %s，使用Debug模式", ginMode)
		}
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		utils.Fatal("数据库连接失败: %v", err)
	}

	// 创建Repository管理器
	repoManager := repo.NewRepositoryManager(db.DB)

	// 创建任务管理器
	taskManager := task.NewTaskManager(repoManager)

	// 注册转存任务处理器
	transferProcessor := task.NewTransferProcessor(repoManager)
	taskManager.RegisterProcessor(transferProcessor)

	// 初始化Meilisearch管理器
	meilisearchManager := services.NewMeilisearchManager(repoManager)
	if err := meilisearchManager.Initialize(); err != nil {
		utils.Error("初始化Meilisearch管理器失败: %v", err)
	}

	// 恢复运行中的任务（服务器重启后）
	if err := taskManager.RecoverRunningTasks(); err != nil {
		utils.Error("恢复运行中任务失败: %v", err)
	} else {
		utils.Info("运行中任务恢复完成")
	}

	utils.Info("任务管理器初始化完成")

	// 创建Gin实例
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 将Repository管理器注入到handlers中
	handlers.SetRepositoryManager(repoManager)

	// 设置Meilisearch管理器到handlers中
	handlers.SetMeilisearchManager(meilisearchManager)

	// 设置全局调度器的Meilisearch管理器
	scheduler.SetGlobalMeilisearchManager(meilisearchManager)

	// 设置公开API中间件的Repository管理器
	middleware.SetRepositoryManager(repoManager)

	// 创建公开API处理器
	publicAPIHandler := handlers.NewPublicAPIHandler()

	// 创建任务处理器
	taskHandler := handlers.NewTaskHandler(repoManager, taskManager)

	// 创建文件处理器
	fileHandler := handlers.NewFileHandler(repoManager.FileRepository, repoManager.SystemConfigRepository, repoManager.UserRepository)

	// 创建Meilisearch处理器
	meilisearchHandler := handlers.NewMeilisearchHandler(meilisearchManager)

	// API路由
	api := r.Group("/api")
	{
		// 公开API路由（需要API Token认证）
		publicAPI := api.Group("/public")
		publicAPI.Use(middleware.PublicAPIAuth())
		{
			// 批量添加资源
			publicAPI.POST("/resources/batch-add", publicAPIHandler.AddBatchResources)
			// 资源搜索
			publicAPI.GET("/resources/search", publicAPIHandler.SearchResources)
			// 热门剧
			publicAPI.GET("/hot-dramas", publicAPIHandler.GetHotDramas)
		}

		// 认证路由
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/register", handlers.Register)
		api.GET("/auth/profile", middleware.AuthMiddleware(), handlers.GetProfile)

		// 资源管理
		api.GET("/resources", handlers.GetResources)
		api.POST("/resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateResource)
		api.PUT("/resources/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateResource)
		api.DELETE("/resources/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteResource)
		api.GET("/resources/:id", handlers.GetResourceByID)
		api.GET("/resources/check-exists", handlers.CheckResourceExists)
		api.POST("/resources/:id/view", handlers.IncrementResourceViewCount)
		api.GET("/resources/:id/link", handlers.GetResourceLink)
		api.DELETE("/resources/batch", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.BatchDeleteResources)

		// 分类管理
		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateCategory)
		api.PUT("/categories/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateCategory)
		api.DELETE("/categories/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteCategory)

		// 搜索
		api.GET("/search", handlers.SearchResources)

		// 统计
		api.GET("/stats", handlers.GetStats)
		api.GET("/performance", handlers.GetPerformanceStats)
		api.GET("/stats/views-trend", handlers.GetViewsTrend)
		api.GET("/stats/searches-trend", handlers.GetSearchesTrend)
		api.GET("/system/info", handlers.GetSystemInfo)

		// 平台管理
		api.GET("/pans", handlers.GetPans)
		api.POST("/pans", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreatePan)
		api.PUT("/pans/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdatePan)
		api.DELETE("/pans/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeletePan)
		api.GET("/pans/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetPan)

		// Cookie管理
		api.GET("/cks", handlers.GetCks)
		api.POST("/cks", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateCks)
		api.PUT("/cks/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateCks)
		api.DELETE("/cks/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteCks)
		api.GET("/cks/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetCksByID)
		api.POST("/cks/:id/refresh-capacity", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.RefreshCapacity)

		// 标签管理
		api.GET("/tags", handlers.GetTags)
		api.POST("/tags", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateTag)
		api.PUT("/tags/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateTag)
		api.DELETE("/tags/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteTag)
		api.GET("/tags/:id", handlers.GetTagByID)
		api.GET("/categories/:categoryId/tags", handlers.GetTagsByCategory)

		// 待处理资源管理
		api.GET("/ready-resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetReadyResources)
		api.POST("/ready-resources/batch", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.BatchCreateReadyResources)
		api.POST("/ready-resources/text", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateReadyResourcesFromText)
		api.DELETE("/ready-resources/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteReadyResource)
		api.DELETE("/ready-resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ClearReadyResources)
		api.GET("/ready-resources/key/:key", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetReadyResourcesByKey)
		api.DELETE("/ready-resources/key/:key", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteReadyResourcesByKey)
		api.GET("/ready-resources/errors", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetReadyResourcesWithErrors)
		api.POST("/ready-resources/:id/clear-error", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ClearErrorMsg)
		api.POST("/ready-resources/retry-failed", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.RetryFailedResources)
		api.POST("/ready-resources/batch-restore", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.BatchRestoreToReadyPool)
		api.POST("/ready-resources/batch-restore-by-query", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.BatchRestoreToReadyPoolByQuery)
		api.POST("/ready-resources/clear-all-errors-by-query", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ClearAllErrorsByQuery)

		// 用户管理（仅管理员）
		api.GET("/users", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetUsers)
		api.POST("/users", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateUser)
		api.PUT("/users/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateUser)
		api.PUT("/users/:id/password", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ChangePassword)
		api.DELETE("/users/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteUser)

		// 搜索统计路由
		api.GET("/search-stats", handlers.GetSearchStats)
		api.GET("/search-stats/hot-keywords", handlers.GetHotKeywords)
		api.GET("/search-stats/daily", handlers.GetDailyStats)
		api.GET("/search-stats/trend", handlers.GetSearchTrend)
		api.GET("/search-stats/keyword/:keyword/trend", handlers.GetKeywordTrend)
		api.POST("/search-stats", handlers.RecordSearch)
		api.POST("/search-stats/record", handlers.RecordSearch)
		api.GET("/search-stats/summary", handlers.GetSearchStatsSummary)

		// 系统配置路由
		api.GET("/system/config", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetSystemConfig)
		api.POST("/system/config", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateSystemConfig)
		api.GET("/system/config/status", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetConfigStatus)
		api.POST("/system/config/toggle-auto-process", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ToggleAutoProcess)
		api.GET("/public/system-config", handlers.GetPublicSystemConfig)

		// 热播剧管理路由（查询接口无需认证）
		api.GET("/hot-dramas", handlers.GetHotDramaList)
		api.GET("/hot-dramas/:id", handlers.GetHotDramaByID)
		api.POST("/hot-dramas", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateHotDrama)
		api.PUT("/hot-dramas/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateHotDrama)
		api.DELETE("/hot-dramas/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteHotDrama)

		// 任务管理路由
		api.POST("/tasks/transfer", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.CreateBatchTransferTask)
		api.GET("/tasks", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.GetTasks)
		api.GET("/tasks/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.GetTaskStatus)
		api.POST("/tasks/:id/start", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.StartTask)
		api.POST("/tasks/:id/stop", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.StopTask)
		api.POST("/tasks/:id/pause", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.PauseTask)
		api.DELETE("/tasks/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.DeleteTask)
		api.GET("/tasks/:id/items", middleware.AuthMiddleware(), middleware.AdminMiddleware(), taskHandler.GetTaskItems)

		// 版本管理路由
		api.GET("/version", handlers.GetVersion)
		api.GET("/version/string", handlers.GetVersionString)
		api.GET("/version/full", handlers.GetFullVersionInfo)
		api.GET("/version/check-update", handlers.CheckUpdate)

		// Meilisearch管理路由
		api.GET("/meilisearch/status", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetStatus)
		api.GET("/meilisearch/unsynced-count", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetUnsyncedCount)
		api.GET("/meilisearch/unsynced", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetUnsyncedResources)
		api.GET("/meilisearch/synced", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetSyncedResources)
		api.GET("/meilisearch/resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetAllResources)
		api.POST("/meilisearch/sync-all", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.SyncAllResources)
		api.GET("/meilisearch/sync-progress", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.GetSyncProgress)
		api.POST("/meilisearch/stop-sync", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.StopSync)
		api.POST("/meilisearch/clear-index", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.ClearIndex)
		api.POST("/meilisearch/test-connection", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.TestConnection)
		api.POST("/meilisearch/update-settings", middleware.AuthMiddleware(), middleware.AdminMiddleware(), meilisearchHandler.UpdateIndexSettings)

		// 文件上传相关路由
		api.POST("/files/upload", middleware.AuthMiddleware(), fileHandler.UploadFile)
		api.GET("/files", middleware.AuthMiddleware(), fileHandler.GetFileList)
		api.DELETE("/files", middleware.AuthMiddleware(), fileHandler.DeleteFiles)
		api.PUT("/files", middleware.AuthMiddleware(), fileHandler.UpdateFile)
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 添加CORS头到静态文件
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		}
		c.Next()
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Info("服务器启动在端口 %s", port)
	r.Run(":" + port)
}
