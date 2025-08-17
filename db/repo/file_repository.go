package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// FileRepository 文件Repository接口
type FileRepository interface {
	BaseRepository[entity.File]
	FindByFileName(fileName string) (*entity.File, error)
	FindByHash(fileHash string) (*entity.File, error)
	FindByUserID(userID uint, page, pageSize int) ([]entity.File, int64, error)
	FindPublicFiles(page, pageSize int) ([]entity.File, int64, error)
	SearchFiles(search string, fileType, status string, userID uint, page, pageSize int) ([]entity.File, int64, error)
	SoftDeleteByIDs(ids []uint) error
	UpdateFileStatus(id uint, status string) error
	UpdateFilePublic(id uint, isPublic bool) error
}

// FileRepositoryImpl 文件Repository实现
type FileRepositoryImpl struct {
	BaseRepositoryImpl[entity.File]
}

// NewFileRepository 创建文件Repository
func NewFileRepository(db *gorm.DB) FileRepository {
	return &FileRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.File]{db: db},
	}
}

// FindByFileName 根据文件名查找文件
func (r *FileRepositoryImpl) FindByFileName(fileName string) (*entity.File, error) {
	var file entity.File
	err := r.db.Where("file_name = ? AND is_deleted = ?", fileName, false).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByUserID 根据用户ID查找文件
func (r *FileRepositoryImpl) FindByUserID(userID uint, page, pageSize int) ([]entity.File, int64, error) {
	var files []entity.File
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := r.db.Model(&entity.File{}).Where("user_id = ? AND is_deleted = ?", userID, false).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取文件列表
	err = r.db.Where("user_id = ? AND is_deleted = ?", userID, false).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}

// FindPublicFiles 查找公开文件
func (r *FileRepositoryImpl) FindPublicFiles(page, pageSize int) ([]entity.File, int64, error) {
	var files []entity.File
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := r.db.Model(&entity.File{}).Where("is_public = ? AND is_deleted = ? AND status = ?", true, false, entity.FileStatusActive).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取文件列表
	err = r.db.Where("is_public = ? AND is_deleted = ? AND status = ?", true, false, entity.FileStatusActive).
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}

// SearchFiles 搜索文件
func (r *FileRepositoryImpl) SearchFiles(search string, fileType, status string, userID uint, page, pageSize int) ([]entity.File, int64, error) {
	var files []entity.File
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&entity.File{}).Where("is_deleted = ?", false)

	// 添加搜索条件
	if search != "" {
		query = query.Where("original_name LIKE ? OR file_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取文件列表
	err = query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&files).Error

	return files, total, err
}

// SoftDeleteByIDs 软删除文件
func (r *FileRepositoryImpl) SoftDeleteByIDs(ids []uint) error {
	return r.db.Model(&entity.File{}).Where("id IN ?", ids).Update("is_deleted", true).Error
}

// UpdateFileStatus 更新文件状态
func (r *FileRepositoryImpl) UpdateFileStatus(id uint, status string) error {
	return r.db.Model(&entity.File{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateFilePublic 更新文件公开状态
func (r *FileRepositoryImpl) UpdateFilePublic(id uint, isPublic bool) error {
	return r.db.Model(&entity.File{}).Where("id = ?", id).Update("is_public", isPublic).Error
}

// FindByHash 根据文件哈希查找文件
func (r *FileRepositoryImpl) FindByHash(fileHash string) (*entity.File, error) {
	var file entity.File
	err := r.db.Where("file_hash = ? AND is_deleted = ?", fileHash, false).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}
