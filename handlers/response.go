package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "操作成功",
		Data:    data,
		Code:    200,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, message string, code int) {
	if code == 0 {
		code = 500
	}
	c.JSON(code, Response{
		Success: false,
		Message: message,
		Data:    nil,
		Code:    code,
	})
}

// ListResponse 列表响应
func ListResponse(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "获取成功",
		Data: gin.H{
			"list":  data,
			"total": total,
		},
		Code: 200,
	})
}

// PageResponse 分页响应
func PageResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "获取成功",
		Data: gin.H{
			"list":  data,
			"total": total,
			"page":  page,
			"limit": limit,
		},
		Code: 200,
	})
}
