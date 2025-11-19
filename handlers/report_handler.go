package handlers

import (
	"net/http"
	"strconv"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ReportHandler struct {
	reportRepo    repo.ReportRepository
	resourceRepo  repo.ResourceRepository
	validate      *validator.Validate
}

func NewReportHandler(reportRepo repo.ReportRepository, resourceRepo repo.ResourceRepository) *ReportHandler {
	return &ReportHandler{
		reportRepo:   reportRepo,
		resourceRepo: resourceRepo,
		validate:     validator.New(),
	}
}

// CreateReport 创建举报
// @Summary 创建举报
// @Description 提交资源举报
// @Tags Report
// @Accept json
// @Produce json
// @Param request body dto.ReportCreateRequest true "举报信息"
// @Success 200 {object} Response{data=dto.ReportResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /reports [post]
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req dto.ReportCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, "参数验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 创建举报实体
	report := &entity.Report{
		ResourceKey: req.ResourceKey,
		Reason:      req.Reason,
		Description: req.Description,
		Contact:     req.Contact,
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		Status:      "pending", // 默认为待处理
	}

	// 保存到数据库
	if err := h.reportRepo.Create(report); err != nil {
		ErrorResponse(c, "创建举报失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := converter.ReportToResponse(report)
	SuccessResponse(c, response)
}

// GetReport 获取举报详情
// @Summary 获取举报详情
// @Description 根据ID获取举报详情
// @Tags Report
// @Produce json
// @Param id path int true "举报ID"
// @Success 200 {object} Response{data=dto.ReportResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /reports/{id} [get]
func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	report, err := h.reportRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, "举报不存在", http.StatusNotFound)
			return
		}
		ErrorResponse(c, "获取举报失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := converter.ReportToResponse(report)
	SuccessResponse(c, response)
}

// ListReports 获取举报列表
// @Summary 获取举报列表
// @Description 获取举报列表（支持分页和状态筛选）
// @Tags Report
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "处理状态"
// @Success 200 {object} Response{data=object{items=[]dto.ReportResponse,total=int}}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /reports [get]
func (h *ReportHandler) ListReports(c *gin.Context) {
	var req dto.ReportListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 验证参数
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, "参数验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	reports, total, err := h.reportRepo.List(req.Status, req.Page, req.PageSize)
	if err != nil {
		ErrorResponse(c, "获取举报列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取每个举报关联的资源
	var reportResponses []*dto.ReportResponse
	for _, report := range reports {
		// 通过资源key查找关联的资源
		resources, err := h.getResourcesByResourceKey(report.ResourceKey)
		if err != nil {
			// 如果获取资源失败，仍然返回基本的举报信息
			reportResponses = append(reportResponses, converter.ReportToResponse(report))
		} else {
			// 使用包含资源详情的转换函数
			response := converter.ReportToResponseWithResources(report, resources)
			reportResponses = append(reportResponses, response)
		}
	}

	PageResponse(c, reportResponses, total, req.Page, req.PageSize)
}

// getResourcesByResourceKey 根据资源key获取关联的资源列表
func (h *ReportHandler) getResourcesByResourceKey(resourceKey string) ([]*entity.Resource, error) {
	// 从资源仓库获取与key关联的所有资源
	resources, err := h.resourceRepo.FindByResourceKey(resourceKey)
	if err != nil {
		return nil, err
	}

	// 将 []entity.Resource 转换为 []*entity.Resource
	var resourcePointers []*entity.Resource
	for i := range resources {
		resourcePointers = append(resourcePointers, &resources[i])
	}

	return resourcePointers, nil
}

// UpdateReport 更新举报状态
// @Summary 更新举报状态
// @Description 更新举报处理状态
// @Tags Report
// @Accept json
// @Produce json
// @Param id path int true "举报ID"
// @Param request body dto.ReportUpdateRequest true "更新信息"
// @Success 200 {object} Response{data=dto.ReportResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /reports/{id} [put]
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.ReportUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, "参数验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 获取当前举报
	_, err = h.reportRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, "举报不存在", http.StatusNotFound)
			return
		}
		ErrorResponse(c, "获取举报失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 更新状态
	processedBy := uint(0) // 从上下文获取当前用户ID，如果存在的话
	if currentUser := c.GetUint("user_id"); currentUser > 0 {
		processedBy = currentUser
	}

	if err := h.reportRepo.UpdateStatus(uint(id), req.Status, &processedBy, req.Note); err != nil {
		ErrorResponse(c, "更新举报状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取更新后的举报信息
	updatedReport, err := h.reportRepo.GetByID(uint(id))
	if err != nil {
		ErrorResponse(c, "获取更新后举报信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, converter.ReportToResponse(updatedReport))
}

// DeleteReport 删除举报
// @Summary 删除举报
// @Description 删除举报记录
// @Tags Report
// @Produce json
// @Param id path int true "举报ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /reports/{id} [delete]
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	if err := h.reportRepo.Delete(uint(id)); err != nil {
		ErrorResponse(c, "删除举报失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, nil)
}

// GetReportByResource 获取某个资源的举报列表
// @Summary 获取资源举报列表
// @Description 获取某个资源的所有举报记录
// @Tags Report
// @Produce json
// @Param resource_key path string true "资源Key"
// @Success 200 {object} Response{data=[]dto.ReportResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /reports/resource/{resource_key} [get]
func (h *ReportHandler) GetReportByResource(c *gin.Context) {
	resourceKey := c.Param("resource_key")
	if resourceKey == "" {
		ErrorResponse(c, "资源Key不能为空", http.StatusBadRequest)
		return
	}

	reports, err := h.reportRepo.GetByResourceKey(resourceKey)
	if err != nil {
		ErrorResponse(c, "获取资源举报失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, converter.ReportsToResponse(reports))
}

// RegisterReportRoutes 注册举报相关路由
func RegisterReportRoutes(router *gin.RouterGroup, reportRepo repo.ReportRepository, resourceRepo repo.ResourceRepository) {
	handler := NewReportHandler(reportRepo, resourceRepo)

	reports := router.Group("/reports")
	{
		reports.POST("", handler.CreateReport)                    // 创建举报
		reports.GET("/:id", handler.GetReport)                    // 获取举报详情
		reports.GET("", handler.ListReports)                      // 获取举报列表
		reports.PUT("/:id", handler.UpdateReport)                 // 更新举报状态
		reports.DELETE("/:id", handler.DeleteReport)              // 删除举报
		reports.GET("/resource/:resource_key", handler.GetReportByResource) // 获取资源举报列表
	}
}