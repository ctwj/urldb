package main

import (
	"log"
	"os"

	"res_db/db"
	"res_db/db/repo"
	"res_db/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// API路由
	api := r.Group("/api")
	{
		// 资源管理
		api.GET("/resources", handlers.GetResources)
		api.POST("/resources", handlers.CreateResource)
		api.PUT("/resources/:id", handlers.UpdateResource)
		api.DELETE("/resources/:id", handlers.DeleteResource)
		api.GET("/resources/:id", handlers.GetResourceByID)

		// 分类管理
		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.PUT("/categories/:id", handlers.UpdateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		// 搜索
		api.GET("/search", handlers.SearchResources)

		// 统计
		api.GET("/stats", handlers.GetStats)

		// 平台管理
		api.GET("/pans", handlers.GetPans)
		api.POST("/pans", handlers.CreatePan)
		api.PUT("/pans/:id", handlers.UpdatePan)
		api.DELETE("/pans/:id", handlers.DeletePan)
		api.GET("/pans/:id", handlers.GetPan)

		// Cookie管理
		api.GET("/cks", handlers.GetCks)
		api.POST("/cks", handlers.CreateCks)
		api.PUT("/cks/:id", handlers.UpdateCks)
		api.DELETE("/cks/:id", handlers.DeleteCks)
		api.GET("/cks/:id", handlers.GetCksByID)

		// 标签管理
		api.GET("/tags", handlers.GetTags)
		api.POST("/tags", handlers.CreateTag)
		api.PUT("/tags/:id", handlers.UpdateTag)
		api.DELETE("/tags/:id", handlers.DeleteTag)
		api.GET("/tags/:id", handlers.GetTagByID)
		api.GET("/resources/:id/tags", handlers.GetResourceTags)

		// 待处理资源管理
		api.GET("/ready-resources", handlers.GetReadyResources)
		api.POST("/ready-resources", handlers.CreateReadyResource)
		api.POST("/ready-resources/batch", handlers.BatchCreateReadyResources)
		api.POST("/ready-resources/text", handlers.CreateReadyResourcesFromText)
		api.DELETE("/ready-resources/:id", handlers.DeleteReadyResource)
		api.DELETE("/ready-resources", handlers.ClearReadyResources)
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	r.Run(":" + port)
}
