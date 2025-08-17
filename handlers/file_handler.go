package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/converter"
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileRepo         repo.FileRepository
	systemConfigRepo repo.SystemConfigRepository
	userRepo         repo.UserRepository
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileRepo repo.FileRepository, systemConfigRepo repo.SystemConfigRepository, userRepo repo.UserRepository) *FileHandler {
	return &FileHandler{
		fileRepo:         fileRepo,
		systemConfigRepo: systemConfigRepo,
		userRepo:         userRepo,
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		ErrorResponse(c, "用户未登录", http.StatusUnauthorized)
		return
	}

	// 从数据库获取用户信息
	currentUser, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		ErrorResponse(c, "用户不存在", http.StatusUnauthorized)
		return
	}

	// 获取文件哈希值
	fileHash := c.PostForm("file_hash")
	
	// 如果提供了文件哈希，先检查是否已存在
	if fileHash != "" {
		existingFile, err := h.fileRepo.FindByHash(fileHash)
		if err == nil && existingFile != nil {
			// 文件已存在，直接返回已存在的文件信息
			utils.Info("文件已存在，跳过上传 - Hash: %s, 文件名: %s", fileHash, existingFile.OriginalName)
			
			response := dto.FileUploadResponse{
				File:       converter.FileToResponse(existingFile),
				Message:    "文件已存在，极速上传成功",
				Success:    true,
				IsDuplicate: true,
			}
			
			SuccessResponse(c, response)
			return
		}
	}

	// 获取上传目录配置（从环境变量或使用默认值）
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads" // 默认值
	}

	// 创建年月子文件夹
	now := time.Now()
	yearMonth := now.Format("200601") // 格式：202508
	monthlyDir := filepath.Join(uploadDir, yearMonth)

	// 确保年月目录存在
	if err := os.MkdirAll(monthlyDir, 0755); err != nil {
		ErrorResponse(c, "创建年月目录失败", http.StatusInternalServerError)
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		ErrorResponse(c, "获取上传文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 检查文件大小（5MB）
	maxFileSize := int64(5 * 1024 * 1024) // 5MB
	if header.Size > maxFileSize {
		ErrorResponse(c, "文件大小不能超过5MB", http.StatusBadRequest)
		return
	}

	// 检查文件类型，只允许图片
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
		"image/svg+xml",
	}

	contentType := header.Header.Get("Content-Type")
	isAllowedType := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowedType = true
			break
		}
	}

	if !isAllowedType {
		ErrorResponse(c, "只支持图片格式文件 (JPEG, PNG, GIF, WebP, BMP, SVG)", http.StatusBadRequest)
		return
	}

	// 生成随机文件名
	fileName := h.generateRandomFileName(header.Filename)
	filePath := filepath.Join(monthlyDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		ErrorResponse(c, "创建文件失败", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		ErrorResponse(c, "保存文件失败", http.StatusInternalServerError)
		return
	}

	// 计算文件哈希值（如果前端没有提供）
	if fileHash == "" {
		fileHash, err = h.calculateFileHash(filePath)
		if err != nil {
			ErrorResponse(c, "计算文件哈希失败", http.StatusInternalServerError)
			return
		}
	}

	// 再次检查文件是否已存在（使用计算出的哈希）
	existingFile, err := h.fileRepo.FindByHash(fileHash)
	if err == nil && existingFile != nil {
		// 文件已存在，删除刚上传的文件，返回已存在的文件信息
		os.Remove(filePath)
		utils.Info("文件已存在，跳过上传 - Hash: %s, 文件名: %s", fileHash, existingFile.OriginalName)
		
		response := dto.FileUploadResponse{
			File:       converter.FileToResponse(existingFile),
			Message:    "文件已存在，极速上传成功",
			Success:    true,
			IsDuplicate: true,
		}
		
		SuccessResponse(c, response)
		return
	}

	// 获取文件类型
	fileType := h.getFileType(header.Filename)
	mimeType := header.Header.Get("Content-Type")

	// 获取是否公开
	isPublic := true
	if isPublicStr := c.PostForm("is_public"); isPublicStr != "" {
		if isPublicBool, err := strconv.ParseBool(isPublicStr); err == nil {
			isPublic = isPublicBool
		}
	}

	// 构建访问URL（使用绝对路径，不包含域名）
	accessURL := fmt.Sprintf("/uploads/%s/%s", yearMonth, fileName)

	// 创建文件记录
	fileEntity := &entity.File{
		OriginalName: header.Filename,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     header.Size,
		FileType:     fileType,
		MimeType:     mimeType,
		FileHash:     fileHash,
		AccessURL:    accessURL,
		UserID:       currentUser.ID,
		Status:       entity.FileStatusActive,
		IsPublic:     isPublic,
		IsDeleted:    false,
	}

	// 保存到数据库
	if err := h.fileRepo.Create(fileEntity); err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		ErrorResponse(c, "保存文件记录失败", http.StatusInternalServerError)
		return
	}

	// 返回响应
	response := dto.FileUploadResponse{
		File:    converter.FileToResponse(fileEntity),
		Message: "文件上传成功",
		Success: true,
	}

	SuccessResponse(c, response)
}

// GetFileList 获取文件列表
func (h *FileHandler) GetFileList(c *gin.Context) {
	var req dto.FileListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 获取当前用户ID和角色
	userIDInterface, exists := c.Get("user_id")
	roleInterface, _ := c.Get("role")

	var userID uint
	var role string
	if exists {
		userID = userIDInterface.(uint)
	}
	if roleInterface != nil {
		role = roleInterface.(string)
	}

	// 根据用户角色决定查询范围
	var files []entity.File
	var total int64
	var err error

	if exists && role == "admin" {
		// 管理员可以查看所有文件
		files, total, err = h.fileRepo.SearchFiles(req.Search, req.FileType, req.Status, req.UserID, req.Page, req.PageSize)
	} else if userID > 0 {
		// 普通用户只能查看自己的文件
		files, total, err = h.fileRepo.SearchFiles(req.Search, req.FileType, req.Status, userID, req.Page, req.PageSize)
	} else {
		// 未登录用户只能查看公开文件
		files, total, err = h.fileRepo.FindPublicFiles(req.Page, req.PageSize)
	}

	if err != nil {
		ErrorResponse(c, "获取文件列表失败", http.StatusInternalServerError)
		return
	}

	response := converter.FileListToResponse(files, total, req.Page, req.PageSize)
	SuccessResponse(c, response)
}

// DeleteFiles 删除文件
func (h *FileHandler) DeleteFiles(c *gin.Context) {
	var req dto.FileDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 获取当前用户ID和角色
	userIDInterface, exists := c.Get("user_id")
	roleInterface, _ := c.Get("role")
	if !exists {
		ErrorResponse(c, "用户未登录", http.StatusUnauthorized)
		return
	}

	userID := userIDInterface.(uint)
	role := ""
	if roleInterface != nil {
		role = roleInterface.(string)
	}

	// 检查权限
	if role != "admin" {
		// 普通用户只能删除自己的文件
		for _, id := range req.IDs {
			file, err := h.fileRepo.FindByID(id)
			if err != nil {
				ErrorResponse(c, "文件不存在", http.StatusNotFound)
				return
			}
			if file.UserID != userID {
				ErrorResponse(c, "没有权限删除此文件", http.StatusForbidden)
				return
			}
		}
	}

	// 获取要删除的文件信息
	var filesToDelete []entity.File
	for _, id := range req.IDs {
		file, err := h.fileRepo.FindByID(id)
		if err != nil {
			ErrorResponse(c, "文件不存在", http.StatusNotFound)
			return
		}
		filesToDelete = append(filesToDelete, *file)
	}

	// 删除本地文件
	for _, file := range filesToDelete {
		if err := os.Remove(file.FilePath); err != nil {
			utils.Error("删除本地文件失败: %s, 错误: %v", file.FilePath, err)
			// 继续删除其他文件，不因为单个文件删除失败而中断
		}
	}

	// 删除数据库记录
	for _, id := range req.IDs {
		if err := h.fileRepo.Delete(id); err != nil {
			utils.Error("删除文件记录失败: ID=%d, 错误: %v", id, err)
			// 继续删除其他文件，不因为单个文件删除失败而中断
		}
	}

	SuccessResponse(c, gin.H{"message": "文件删除成功"})
}

// UpdateFile 更新文件信息
func (h *FileHandler) UpdateFile(c *gin.Context) {
	var req dto.FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "请求参数错误", http.StatusBadRequest)
		return
	}

	// 获取当前用户ID和角色
	userIDInterface, exists := c.Get("user_id")
	roleInterface, _ := c.Get("role")
	if !exists {
		ErrorResponse(c, "用户未登录", http.StatusUnauthorized)
		return
	}

	userID := userIDInterface.(uint)
	role := ""
	if roleInterface != nil {
		role = roleInterface.(string)
	}

	// 查找文件
	file, err := h.fileRepo.FindByID(req.ID)
	if err != nil {
		ErrorResponse(c, "文件不存在", http.StatusNotFound)
		return
	}

	// 检查权限
	if role != "admin" && userID != file.UserID {
		ErrorResponse(c, "没有权限修改此文件", http.StatusForbidden)
		return
	}

	// 更新文件信息
	if req.IsPublic != nil {
		if err := h.fileRepo.UpdateFilePublic(req.ID, *req.IsPublic); err != nil {
			ErrorResponse(c, "更新文件状态失败", http.StatusInternalServerError)
			return
		}
	}

	if req.Status != "" {
		if err := h.fileRepo.UpdateFileStatus(req.ID, req.Status); err != nil {
			ErrorResponse(c, "更新文件状态失败", http.StatusInternalServerError)
			return
		}
	}

	SuccessResponse(c, gin.H{"message": "文件更新成功"})
}

// generateRandomFileName 生成随机文件名
func (h *FileHandler) generateRandomFileName(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)

	// 生成随机字符串
	bytes := make([]byte, 16)
	rand.Read(bytes)
	randomStr := fmt.Sprintf("%x", bytes)

	// 添加时间戳
	timestamp := time.Now().Unix()

	return fmt.Sprintf("%d_%s%s", timestamp, randomStr, ext)
}

// getFileType 获取文件类型
func (h *FileHandler) getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// 图片类型
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return "image"
		}
	}

	return "image" // 默认返回image，因为只支持图片格式
}

// calculateFileHash 计算文件哈希值
func (h *FileHandler) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
