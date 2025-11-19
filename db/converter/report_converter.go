package converter

import (
	"time"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
)

// ReportToResponseWithResources 将举报实体和关联资源转换为响应对象
func ReportToResponseWithResources(report *entity.Report, resources []*entity.Resource) *dto.ReportResponse {
	if report == nil {
		return nil
	}

	// 转换关联的资源信息
	var resourceInfos []dto.ResourceInfo
	for _, resource := range resources {
		categoryName := ""
		if resource.Category.ID != 0 {
			categoryName = resource.Category.Name
		}

		panName := ""
		if resource.Pan.ID != 0 {
			panName = resource.Pan.Name
		}

		resourceInfo := dto.ResourceInfo{
			ID:          resource.ID,
			Title:       resource.Title,
			Description: resource.Description,
			URL:         resource.URL,
			SaveURL:     resource.SaveURL,
			FileSize:    resource.FileSize,
			Category:    categoryName,
			PanName:     panName,
			ViewCount:   resource.ViewCount,
			IsValid:     resource.IsValid,
			CreatedAt:   resource.CreatedAt.Format(time.RFC3339),
		}
		resourceInfos = append(resourceInfos, resourceInfo)
	}

	return &dto.ReportResponse{
		ID:          report.ID,
		ResourceKey: report.ResourceKey,
		Reason:      report.Reason,
		Description: report.Description,
		Contact:     report.Contact,
		UserAgent:   report.UserAgent,
		IPAddress:   report.IPAddress,
		Status:      report.Status,
		Note:        report.Note,
		CreatedAt:   report.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   report.UpdatedAt.Format(time.RFC3339),
		Resources:   resourceInfos,
	}
}

// ReportToResponse 将举报实体转换为响应对象（不包含资源详情）
func ReportToResponse(report *entity.Report) *dto.ReportResponse {
	if report == nil {
		return nil
	}

	return &dto.ReportResponse{
		ID:          report.ID,
		ResourceKey: report.ResourceKey,
		Reason:      report.Reason,
		Description: report.Description,
		Contact:     report.Contact,
		UserAgent:   report.UserAgent,
		IPAddress:   report.IPAddress,
		Status:      report.Status,
		Note:        report.Note,
		CreatedAt:   report.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   report.UpdatedAt.Format(time.RFC3339),
		Resources:   []dto.ResourceInfo{}, // 空的资源列表
	}
}

// ReportsToResponse 将举报实体列表转换为响应对象列表
func ReportsToResponse(reports []*entity.Report) []*dto.ReportResponse {
	var responses []*dto.ReportResponse
	for _, report := range reports {
		responses = append(responses, ReportToResponse(report))
	}
	return responses
}