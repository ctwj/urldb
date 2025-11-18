package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"time"
)

// ReportToResponse 将举报实体转换为响应对象
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