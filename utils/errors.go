package utils

import "fmt"

// ErrorType 错误类型枚举
type ErrorType string

const (
	// ErrorTypeUnsupportedLink 不支持的链接
	ErrorTypeUnsupportedLink ErrorType = "UNSUPPORTED_LINK"
	// ErrorTypeInvalidLink 无效链接
	ErrorTypeInvalidLink ErrorType = "INVALID_LINK"
	// ErrorTypeNoAccount 没有可用账号
	ErrorTypeNoAccount ErrorType = "NO_ACCOUNT"
	// ErrorTypeNoValidAccount 没有有效账号
	ErrorTypeNoValidAccount ErrorType = "NO_VALID_ACCOUNT"
	// ErrorTypeServiceCreation 服务创建失败
	ErrorTypeServiceCreation ErrorType = "SERVICE_CREATION_FAILED"
	// ErrorTypeTransferFailed 转存失败
	ErrorTypeTransferFailed ErrorType = "TRANSFER_FAILED"
	// ErrorTypeTagProcessing 标签处理失败
	ErrorTypeTagProcessing ErrorType = "TAG_PROCESSING_FAILED"
	// ErrorTypeCategoryProcessing 分类处理失败
	ErrorTypeCategoryProcessing ErrorType = "CATEGORY_PROCESSING_FAILED"
	// ErrorTypeResourceSave 资源保存失败
	ErrorTypeResourceSave ErrorType = "RESOURCE_SAVE_FAILED"
	// ErrorTypePlatformNotFound 平台未找到
	ErrorTypePlatformNotFound ErrorType = "PLATFORM_NOT_FOUND"
	// ErrorTypeLinkCheckFailed 链接检查失败
	ErrorTypeLinkCheckFailed ErrorType = "LINK_CHECK_FAILED"
)

// ResourceError 资源处理错误
type ResourceError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	URL     string    `json:"url,omitempty"`
	Details string    `json:"details,omitempty"`
}

// Error 实现error接口
func (e *ResourceError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// NewResourceError 创建新的资源错误
func NewResourceError(errorType ErrorType, message string, url string, details string) *ResourceError {
	return &ResourceError{
		Type:    errorType,
		Message: message,
		URL:     url,
		Details: details,
	}
}

// NewUnsupportedLinkError 创建不支持的链接错误
func NewUnsupportedLinkError(url string) *ResourceError {
	return NewResourceError(ErrorTypeUnsupportedLink, "不支持的链接地址", url, "")
}

// NewInvalidLinkError 创建无效链接错误
func NewInvalidLinkError(url string, details string) *ResourceError {
	return NewResourceError(ErrorTypeInvalidLink, "链接无效", url, details)
}

// NewNoAccountError 创建没有账号错误
func NewNoAccountError(platform string) *ResourceError {
	return NewResourceError(ErrorTypeNoAccount, "没有可用的网盘账号", "", fmt.Sprintf("平台: %s", platform))
}

// NewNoValidAccountError 创建没有有效账号错误
func NewNoValidAccountError(platform string) *ResourceError {
	return NewResourceError(ErrorTypeNoValidAccount, "没有有效的网盘账号", "", fmt.Sprintf("平台: %s", platform))
}

// NewServiceCreationError 创建服务创建失败错误
func NewServiceCreationError(url string, details string) *ResourceError {
	return NewResourceError(ErrorTypeServiceCreation, "创建网盘服务失败", url, details)
}

// NewTransferFailedError 创建转存失败错误
func NewTransferFailedError(url string, details string) *ResourceError {
	return NewResourceError(ErrorTypeTransferFailed, "网盘信息获取失败", url, details)
}

// NewTagProcessingError 创建标签处理失败错误
func NewTagProcessingError(details string) *ResourceError {
	return NewResourceError(ErrorTypeTagProcessing, "处理标签失败", "", details)
}

// NewCategoryProcessingError 创建分类处理失败错误
func NewCategoryProcessingError(details string) *ResourceError {
	return NewResourceError(ErrorTypeCategoryProcessing, "处理分类失败", "", details)
}

// NewResourceSaveError 创建资源保存失败错误
func NewResourceSaveError(url string, details string) *ResourceError {
	return NewResourceError(ErrorTypeResourceSave, "资源保存失败", url, details)
}

// NewPlatformNotFoundError 创建平台未找到错误
func NewPlatformNotFoundError(platform string) *ResourceError {
	return NewResourceError(ErrorTypePlatformNotFound, "未找到对应的平台ID", "", fmt.Sprintf("平台: %s", platform))
}

// NewLinkCheckError 创建链接检查失败错误
func NewLinkCheckError(url string, details string) *ResourceError {
	return NewResourceError(ErrorTypeLinkCheckFailed, "链接检查失败", url, details)
}

// IsResourceError 检查是否为资源错误
func IsResourceError(err error) bool {
	_, ok := err.(*ResourceError)
	return ok
}

// GetResourceError 获取资源错误
func GetResourceError(err error) *ResourceError {
	if resourceErr, ok := err.(*ResourceError); ok {
		return resourceErr
	}
	return nil
}

// GetErrorType 获取错误类型
func GetErrorType(err error) ErrorType {
	if resourceErr := GetResourceError(err); resourceErr != nil {
		return resourceErr.Type
	}
	return ""
}

// IsRetryableError 检查是否为可重试的错误
func IsRetryableError(err error) bool {
	errorType := GetErrorType(err)
	switch errorType {
	case ErrorTypeNoAccount, ErrorTypeNoValidAccount, ErrorTypeTransferFailed, ErrorTypeLinkCheckFailed:
		return true
	default:
		return false
	}
}

// GetErrorSummary 获取错误摘要
func GetErrorSummary(err error) string {
	if resourceErr := GetResourceError(err); resourceErr != nil {
		return fmt.Sprintf("%s: %s", resourceErr.Type, resourceErr.Message)
	}
	return err.Error()
}
