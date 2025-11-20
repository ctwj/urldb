package handlers

import (
	"net/http"
	"strconv"
	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CopyrightClaimHandler struct {
	copyrightClaimRepo repo.CopyrightClaimRepository
	resourceRepo       repo.ResourceRepository
	validate           *validator.Validate
}

func NewCopyrightClaimHandler(copyrightClaimRepo repo.CopyrightClaimRepository, resourceRepo repo.ResourceRepository) *CopyrightClaimHandler {
	return &CopyrightClaimHandler{
		copyrightClaimRepo: copyrightClaimRepo,
		resourceRepo:       resourceRepo,
		validate:           validator.New(),
	}
}

// CreateCopyrightClaim 创建版权申述
// @Summary 创建版权申述
// @Description 提交资源版权申述
// @Tags CopyrightClaim
// @Accept json
// @Produce json
// @Param request body dto.CopyrightClaimCreateRequest true "版权申述信息"
// @Success 200 {object} Response{data=dto.CopyrightClaimResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims [post]
func (h *CopyrightClaimHandler) CreateCopyrightClaim(c *gin.Context) {
	var req dto.CopyrightClaimCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, "参数验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 创建版权申述实体
	claim := &entity.CopyrightClaim{
		ResourceKey:  req.ResourceKey,
		Identity:     req.Identity,
		ProofType:    req.ProofType,
		Reason:       req.Reason,
		ContactInfo:  req.ContactInfo,
		ClaimantName: req.ClaimantName,
		ProofFiles:   req.ProofFiles,
		UserAgent:    req.UserAgent,
		IPAddress:    req.IPAddress,
		Status:       "pending", // 默认为待处理
	}

	// 保存到数据库
	if err := h.copyrightClaimRepo.Create(claim); err != nil {
		ErrorResponse(c, "创建版权申述失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := converter.CopyrightClaimToResponse(claim)
	SuccessResponse(c, response)
}

// GetCopyrightClaim 获取版权申述详情
// @Summary 获取版权申述详情
// @Description 根据ID获取版权申述详情
// @Tags CopyrightClaim
// @Produce json
// @Param id path int true "版权申述ID"
// @Success 200 {object} Response{data=dto.CopyrightClaimResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims/{id} [get]
func (h *CopyrightClaimHandler) GetCopyrightClaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	claim, err := h.copyrightClaimRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, "版权申述不存在", http.StatusNotFound)
			return
		}
		ErrorResponse(c, "获取版权申述失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, converter.CopyrightClaimToResponse(claim))
}

// ListCopyrightClaims 获取版权申述列表
// @Summary 获取版权申述列表
// @Description 获取版权申述列表（支持分页和状态筛选）
// @Tags CopyrightClaim
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "处理状态"
// @Success 200 {object} Response{data=object{items=[]dto.CopyrightClaimResponse,total=int}}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims [get]
func (h *CopyrightClaimHandler) ListCopyrightClaims(c *gin.Context) {
	var req dto.CopyrightClaimListRequest
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

	claims, total, err := h.copyrightClaimRepo.List(req.Status, req.Page, req.PageSize)
	if err != nil {
		ErrorResponse(c, "获取版权申述列表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 转换为包含资源信息的响应
	var responses []*dto.CopyrightClaimResponse
	for _, claim := range claims {
		// 查询关联的资源信息
		resources, err := h.getResourcesByResourceKey(claim.ResourceKey)
		if err != nil {
			// 如果查询资源失败，使用空资源列表
			responses = append(responses, converter.CopyrightClaimToResponse(claim))
		} else {
			// 使用包含资源详情的转换函数
			responses = append(responses, converter.CopyrightClaimToResponseWithResources(claim, resources))
		}
	}

	PageResponse(c, responses, total, req.Page, req.PageSize)
}

// getResourcesByResourceKey 根据资源key获取关联的资源列表
func (h *CopyrightClaimHandler) getResourcesByResourceKey(resourceKey string) ([]*entity.Resource, error) {
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

// UpdateCopyrightClaim 更新版权申述状态
// @Summary 更新版权申述状态
// @Description 更新版权申述处理状态
// @Tags CopyrightClaim
// @Accept json
// @Produce json
// @Param id path int true "版权申述ID"
// @Param request body dto.CopyrightClaimUpdateRequest true "更新信息"
// @Success 200 {object} Response{data=dto.CopyrightClaimResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims/{id} [put]
func (h *CopyrightClaimHandler) UpdateCopyrightClaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	var req dto.CopyrightClaimUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, "参数验证失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 获取当前版权申述
	_, err = h.copyrightClaimRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ErrorResponse(c, "版权申述不存在", http.StatusNotFound)
			return
		}
		ErrorResponse(c, "获取版权申述失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 更新状态
	processedBy := uint(0) // 从上下文获取当前用户ID，如果存在的话
	if currentUser := c.GetUint("user_id"); currentUser > 0 {
		processedBy = currentUser
	}

	if err := h.copyrightClaimRepo.UpdateStatus(uint(id), req.Status, &processedBy, req.Note); err != nil {
		ErrorResponse(c, "更新版权申述状态失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取更新后的版权申述信息
	updatedClaim, err := h.copyrightClaimRepo.GetByID(uint(id))
	if err != nil {
		ErrorResponse(c, "获取更新后版权申述信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, converter.CopyrightClaimToResponse(updatedClaim))
}

// DeleteCopyrightClaim 删除版权申述
// @Summary 删除版权申述
// @Description 删除版权申述记录
// @Tags CopyrightClaim
// @Produce json
// @Param id path int true "版权申述ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims/{id} [delete]
func (h *CopyrightClaimHandler) DeleteCopyrightClaim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的ID", http.StatusBadRequest)
		return
	}

	if err := h.copyrightClaimRepo.Delete(uint(id)); err != nil {
		ErrorResponse(c, "删除版权申述失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, nil)
}

// GetCopyrightClaimByResource 获取某个资源的版权申述列表
// @Summary 获取资源版权申述列表
// @Description 获取某个资源的所有版权申述记录
// @Tags CopyrightClaim
// @Produce json
// @Param resource_key path string true "资源Key"
// @Success 200 {object} Response{data=[]dto.CopyrightClaimResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /copyright-claims/resource/{resource_key} [get]
func (h *CopyrightClaimHandler) GetCopyrightClaimByResource(c *gin.Context) {
	resourceKey := c.Param("resource_key")
	if resourceKey == "" {
		ErrorResponse(c, "资源Key不能为空", http.StatusBadRequest)
		return
	}

	claims, err := h.copyrightClaimRepo.GetByResourceKey(resourceKey)
	if err != nil {
		ErrorResponse(c, "获取资源版权申述失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, converter.CopyrightClaimsToResponse(claims))
}

// RegisterCopyrightClaimRoutes 注册版权申述相关路由
func RegisterCopyrightClaimRoutes(router *gin.RouterGroup, copyrightClaimRepo repo.CopyrightClaimRepository, resourceRepo repo.ResourceRepository) {
	handler := NewCopyrightClaimHandler(copyrightClaimRepo, resourceRepo)

	claims := router.Group("/copyright-claims")
	{
		claims.POST("", handler.CreateCopyrightClaim)                    // 创建版权申述
		claims.GET("/:id", handler.GetCopyrightClaim)                    // 获取版权申述详情
		claims.GET("", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handler.ListCopyrightClaims)                      // 获取版权申述列表
		claims.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handler.UpdateCopyrightClaim)                 // 更新版权申述状态
		claims.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handler.DeleteCopyrightClaim)              // 删除版权申述
		claims.GET("/resource/:resource_key", handler.GetCopyrightClaimByResource) // 获取资源版权申述列表
	}
}