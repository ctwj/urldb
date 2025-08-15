package converter

import (
	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
)

// FileToResponse 将文件实体转换为响应DTO
func FileToResponse(file *entity.File) dto.FileResponse {
	response := dto.FileResponse{
		ID:           file.ID,
		CreatedAt:    utils.FormatTime(file.CreatedAt, "2006-01-02 15:04:05"),
		UpdatedAt:    utils.FormatTime(file.UpdatedAt, "2006-01-02 15:04:05"),
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		MimeType:     file.MimeType,
		AccessURL:    file.AccessURL,
		UserID:       file.UserID,
		Status:       file.Status,
		IsPublic:     file.IsPublic,
		IsDeleted:    file.IsDeleted,
	}

	// 添加用户名
	if file.User.ID > 0 {
		response.User = file.User.Username
	}

	return response
}

// FilesToResponse 将文件实体列表转换为响应DTO列表
func FilesToResponse(files []entity.File) []dto.FileResponse {
	var responses []dto.FileResponse
	for _, file := range files {
		responses = append(responses, FileToResponse(&file))
	}
	return responses
}

// FileListToResponse 将文件列表转换为列表响应
func FileListToResponse(files []entity.File, total int64, page, pageSize int) dto.FileListResponse {
	return dto.FileListResponse{
		Files: FilesToResponse(files),
		Total: total,
		Page:  page,
		Size:  pageSize,
	}
}
