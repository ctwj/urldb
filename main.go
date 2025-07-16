// @title 网盘资源管理系统公开API
// @version 1.0
// @description 网盘资源管理系统的公开API接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/public

// @securityDefinitions.apikey ApiTokenAuth
// @in header
// @name X-API-Token
// @description API Token认证

package main

import (
	"log"
	"os"
	"res_db/utils"

	"res_db/db"
	"res_db/db/repo"
	"res_db/handlers"
	"res_db/middleware"

	_ "res_db/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建Repository管理器
	repoManager := repo.NewRepositoryManager(db.DB)

	// 创建全局调度器
	scheduler := utils.GetGlobalScheduler(
		repoManager.HotDramaRepository,
		repoManager.ReadyResourceRepository,
		repoManager.ResourceRepository,
		repoManager.SystemConfigRepository,
		repoManager.PanRepository,
	)

	// 检查系统配置，决定是否启动待处理资源自动处理任务
	systemConfig, err := repoManager.SystemConfigRepository.GetOrCreateDefault()
	if err != nil {
		log.Printf("获取系统配置失败: %v", err)
	} else {
		// 检查是否启动待处理资源自动处理任务
		if systemConfig.AutoProcessReadyResources {
			scheduler.StartReadyResourceScheduler()
			log.Println("已启动待处理资源自动处理任务")
		} else {
			log.Println("系统配置中自动处理待处理资源功能已禁用，跳过启动定时任务")
		}

		// 检查是否启动热播剧自动拉取任务
		if systemConfig.AutoFetchHotDramaEnabled {
			scheduler.StartHotDramaScheduler()
			log.Println("已启动热播剧自动拉取任务")
		} else {
			log.Println("系统配置中自动拉取热播剧功能已禁用，跳过启动定时任务")
		}
	}

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
			// 单个添加资源
			publicAPI.POST("/resources/add", publicAPIHandler.AddSingleResource)
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

		// 标签管理
		api.GET("/tags", handlers.GetTags)
		api.POST("/tags", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateTag)
		api.PUT("/tags/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateTag)
		api.DELETE("/tags/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteTag)
		api.GET("/tags/:id", handlers.GetTagByID)
		api.GET("/categories/:categoryId/tags", handlers.GetTagsByCategory)

		// 待处理资源管理
		api.GET("/ready-resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.GetReadyResources)
		api.POST("/ready-resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateReadyResource)
		api.POST("/ready-resources/batch", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.BatchCreateReadyResources)
		api.POST("/ready-resources/text", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.CreateReadyResourcesFromText)
		api.DELETE("/ready-resources/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.DeleteReadyResource)
		api.DELETE("/ready-resources", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ClearReadyResources)

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
		api.POST("/search-stats/record", handlers.RecordSearch)

		// 系统配置路由
		api.GET("/system/config", handlers.GetSystemConfig)
		api.POST("/system/config", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.UpdateSystemConfig)

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
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 注册Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	r.Run(":" + port)
}
