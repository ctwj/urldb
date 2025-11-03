package main

import (
	"math/rand"
	"time"

	"github.com/ctwj/urldb/utils"
)

// Resource represents a resource entity
type Resource struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	// Initialize logger
	utils.InitLogger(nil)

	// Start the demo plugin
	utils.Info("Demo plugin started")

	// Run the demo task every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fetchAndLogResource()
		}
	}
}

// fetchAndLogResource fetches a random resource and logs it
func fetchAndLogResource() {
	// Simulate fetching a resource from the database
	resource := fetchRandomResource()

	if resource != nil {
		utils.Info("Fetched resource: ID=%d, Title=%s, URL=%s",
			resource.ID, resource.Title, resource.URL)
	} else {
		utils.Warn("No resources found in database")
	}
}

// fetchRandomResource simulates fetching a random resource from the database
func fetchRandomResource() *Resource {
	// In a real implementation, this would query the actual database
	// For demo purposes, we'll generate a random resource

	// Simulate some resources
	resources := []Resource{
		{ID: 1, Title: "Go语言编程指南", Description: "学习Go语言的完整指南", URL: "https://example.com/go-guide", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "微服务架构设计", Description: "构建可扩展的微服务系统", URL: "https://example.com/microservices", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Title: "容器化部署实践", Description: "Docker和Kubernetes实战", URL: "https://example.com/container-deployment", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 4, Title: "数据库优化技巧", Description: "提升数据库性能的方法", URL: "https://example.com/db-optimization", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 5, Title: "前端框架比较", Description: "React vs Vue vs Angular", URL: "https://example.com/frontend-frameworks", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Return a random resource
	if len(resources) > 0 {
		return &resources[rand.Intn(len(resources))]
	}

	return nil
}