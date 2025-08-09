package main

import (
	"log"
	"os"

	"github.com/ctwj/urldb/scheduler"
	"github.com/ctwj/urldb/utils"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/handlers"
	"github.com/ctwj/urldb/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	// 创建全局调度器
	scheduler := scheduler.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
		repoManager.CksRepository,
		repoManager.TagRepository,
		repoManager.CategoryRepository,
	)

	// 确保默认配置存在
	// _, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	// if err != nil {
	// 	utils.Error("初始化默认配置失败: %v", err)
	// } else {
	// 	utils.Info("默认配置初始化完成")
	// }

	// 检查系统配置，决定是否启动各种自动任务
	autoProcessReadyResources, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoProcessReadyResources)
	if err != nil {
		utils.Error("获取自动处理待处理资源配置失败: %v", err)
	} else if autoProcessReadyResources {
		scheduler.StartReadyResourceScheduler()
		utils.Info("已启动待处理资源自动处理任务")
	} else {
		utils.Info("系统配置中自动处理待处理资源功能已禁用，跳过启动定时任务")
	}

	autoFetchHotDramaEnabled, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoFetchHotDramaEnabled)
	if err != nil {
		utils.Error("获取自动拉取热播剧配置失败: %v", err)
	} else if autoFetchHotDramaEnabled {
		scheduler.StartHotDramaScheduler()
		utils.Info("已启动热播剧自动拉取任务")
	} else {
		utils.Info("系统配置中自动拉取热播剧功能已禁用，跳过启动定时任务")
	}

	// autoTransferEnabled, err := repoManager.SystemConfigRepository.GetConfigBool(entity.ConfigKeyAutoTransferEnabled)
	// if err != nil {
	// 	utils.Error("获取自动转存配置失败: %v", err)
	// } else if autoTransferEnabled {
	// 	scheduler.StartAutoTransferScheduler()
	// 	utils.Info("已启动自动转存任务")
	// } else {
	// 	utils.Info("系统配置中自动转存功能已禁用，跳过启动定时任务")
	// }

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

	// 设置公开API中间件的Repository管理器
	middleware.SetRepositoryManager(repoManager)

	// 创建公开API处理器
	publicAPIHandler := handlers.NewPublicAPIHandler()

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
		api.POST("/system/config/toggle-auto-process", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ToggleAutoProcess)
		api.GET("/public/system-config", handlers.GetPublicSystemConfig)

		// 热播剧管理路由（查询接口无需认证）
		api.GET("/hot-dramas", handlers.GetHotDramaList)
		api.GET("/hot-dramas/:id", handlers.GetHotDramaByID)
		api.POST("/hot-dramas", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateHotDrama)
		api.PUT("/hot-dramas/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateHotDrama)
		api.DELETE("/hot-dramas/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteHotDrama)

		// 调度器管理路由（查询接口无需认证）
		api.GET("/scheduler/status", handlers.GetSchedulerStatus)
		api.GET("/scheduler/hot-drama/names", handlers.FetchHotDramaNames)
		api.POST("/scheduler/hot-drama/start", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StartHotDramaScheduler)
		api.POST("/scheduler/hot-drama/stop", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StopHotDramaScheduler)
		api.POST("/scheduler/hot-drama/trigger", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.TriggerHotDramaScheduler)

		// 待处理资源自动处理管理路由
		api.POST("/scheduler/ready-resource/start", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StartReadyResourceScheduler)
		api.POST("/scheduler/ready-resource/stop", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StopReadyResourceScheduler)
		api.POST("/scheduler/ready-resource/trigger", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.TriggerReadyResourceScheduler)

		// 自动转存管理路由
		api.POST("/scheduler/auto-transfer/start", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StartAutoTransferScheduler)
		api.POST("/scheduler/auto-transfer/stop", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.StopAutoTransferScheduler)
		api.POST("/scheduler/auto-transfer/trigger", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.TriggerAutoTransferScheduler)

		// 版本管理路由
		api.GET("/version", handlers.GetVersion)
		api.GET("/version/string", handlers.GetVersionString)
		api.GET("/version/full", handlers.GetFullVersionInfo)
		api.GET("/version/check-update", handlers.CheckUpdate)
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Info("服务器启动在端口 %s", port)
	r.Run(":" + port)
}
