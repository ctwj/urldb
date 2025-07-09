package main

import (
	"log"
	"os"

	"res_db/handlers"
	"res_db/models"

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
	if err := models.InitDB(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 创建Gin实例
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

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
