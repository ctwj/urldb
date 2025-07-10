package handlers

import (
	"github.com/gin-gonic/gin"
)

// StandardResponse 标准化响应结构
type StandardResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Error      string          `json:"error,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(200, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Error:   message,
	})
}

// PaginatedResponse 分页响应
func PaginatedResponse(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(200, StandardResponse{
		Success: true,
		Data:    data,
		Pagination: &PaginationInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// SimpleSuccessResponse 简单成功响应
func SimpleSuccessResponse(c *gin.Context, message string) {
	c.JSON(200, StandardResponse{
		Success: true,
		Message: message,
	})
}

// CreatedResponse 创建成功响应
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(201, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
