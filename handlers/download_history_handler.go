package handlers

import (
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// GetDownloadHistory 获取当前用户的下载历史（分页，按最近下载时间倒序）
func GetDownloadHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	histories, total, err := repoManager.DownloadHistoryRepository.ListByUser(userID, page, pageSize)
	if err != nil {
		utils.Error("GetDownloadHistory - 查询下载历史失败 - 用户ID: %d, Error: %v", userID, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// 批量加载关联资源信息
	resourceIDs := make([]uint, 0, len(histories))
	for _, h := range histories {
		resourceIDs = append(resourceIDs, h.ResourceID)
	}
	resourceMap := make(map[uint]resourceBrief)
	if len(resourceIDs) > 0 {
		resources, err := repoManager.ResourceRepository.FindByIDs(resourceIDs)
		if err != nil {
			utils.Error("GetDownloadHistory - 查询关联资源失败 - 用户ID: %d, Error: %v", userID, err)
		} else {
			panMap := loadPanNameMap()
			for i := range resources {
				res := &resources[i]
				platform := ""
				if res.PanID != nil {
					platform = panMap[*res.PanID]
				}
				resourceMap[res.ID] = resourceBrief{
					Key:         res.Key,
					Title:       res.Title,
					Description: res.Description,
					Platform:    platform,
				}
			}
		}
	}

	items := make([]gin.H, 0, len(histories))
	for _, h := range histories {
		item := gin.H{
			"id":             h.ID,
			"resource_id":    h.ResourceID,
			"download_count": h.DownloadCount,
			"first_download": h.CreatedAt,
			"last_download":  h.UpdatedAt,
		}
		if brief, ok := resourceMap[h.ResourceID]; ok {
			item["resource_key"] = brief.Key
			item["title"] = brief.Title
			item["description"] = brief.Description
			item["platform"] = brief.Platform
		}
		items = append(items, item)
	}

	SuccessResponse(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDownloadHistoryStats 获取当前用户下载历史统计
func GetDownloadHistoryStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	total, today, thisWeek, thisMonth, err := repoManager.DownloadHistoryRepository.StatsByUser(userID)
	if err != nil {
		utils.Error("GetDownloadHistoryStats - 统计下载历史失败 - 用户ID: %d, Error: %v", userID, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, gin.H{
		"total":      total,
		"today":      today,
		"thisWeek":   thisWeek,
		"thisMonth":  thisMonth,
	})
}

// ClearDownloadHistory 清空当前用户全部下载历史
func ClearDownloadHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	affected, err := repoManager.DownloadHistoryRepository.DeleteByUser(userID)
	if err != nil {
		utils.Error("ClearDownloadHistory - 清空下载历史失败 - 用户ID: %d, Error: %v", userID, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Info("ClearDownloadHistory - 清空下载历史 - 用户ID: %d, 行数: %d", userID, affected)
	SuccessResponse(c, gin.H{"message": "下载历史已清空", "affected_rows": affected})
}

// DeleteDownloadHistoryItem 删除当前用户的一条下载历史
func DeleteDownloadHistoryItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	affected, err := repoManager.DownloadHistoryRepository.DeleteByIDAndUser(uint(id), userID)
	if err != nil {
		utils.Error("DeleteDownloadHistoryItem - 删除下载记录失败 - 用户ID: %d, 记录ID: %d, Error: %v", userID, id, err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		ErrorResponse(c, "记录不存在", http.StatusNotFound)
		return
	}

	SuccessResponse(c, gin.H{"message": "删除成功", "affected_rows": affected})
}

// resourceBrief 下载历史列表中展示的资源摘要
type resourceBrief struct {
	Key         string
	Title       string
	Description string
	Platform    string
}

// loadPanNameMap 加载平台ID到展示名的映射（平台表数据量小，一次性加载）
func loadPanNameMap() map[uint]string {
	panMap := make(map[uint]string)
	pans, err := repoManager.PanRepository.FindAll()
	if err != nil {
		utils.Error("loadPanNameMap - 查询平台失败: %v", err)
		return panMap
	}
	for _, pan := range pans {
		name := pan.Remark
		if name == "" {
			name = pan.Name
		}
		panMap[pan.ID] = name
	}
	return panMap
}
