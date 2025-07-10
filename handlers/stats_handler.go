package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	// 获取资源总数
	totalResources, err := repoManager.ResourceRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取分类总数
	totalCategories, err := repoManager.CategoryRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取标签总数
	totalTags, err := repoManager.TagRepository.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计算总浏览次数
	var totalViews int64
	for _, resource := range totalResources {
		totalViews += int64(resource.ViewCount)
	}

	stats := map[string]interface{}{
		"total_resources":  len(totalResources),
		"total_categories": len(totalCategories),
		"total_tags":       len(totalTags),
		"total_views":      totalViews,
	}

	c.JSON(http.StatusOK, stats)
}
