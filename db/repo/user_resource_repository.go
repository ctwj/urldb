package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"gorm.io/gorm"
)

// UserResourceRepository 用户上传资源Repository接口（015-user-resource-upload）
type UserResourceRepository interface {
	BaseRepository[entity.UserResource]
	// FindByIDAndUser 按ID与归属查询（越权隔离，非本人/不存在统一返回未找到）
	FindByIDAndUser(id, userID uint) (*entity.UserResource, error)
	// FindByUserAndURL 本人内按URL查重（FR-008）
	FindByUserAndURL(userID uint, url string) (*entity.UserResource, error)
	// CountTodayByUser 当日提交计数（FR-009 频率限制，Asia/Shanghai 自然日）
	CountTodayByUser(userID uint) (int64, error)
	// ListByUser 分页查询，支持状态与关键词过滤
	ListByUser(userID uint, page, pageSize int, status, keyword string) ([]*entity.UserResource, int64, error)
	// GetStatsByUser 五态计数（前端统计卡）
	GetStatsByUser(userID uint) (stats map[string]int64, err error)
}

// UserResourceRepositoryImpl 用户上传资源Repository实现
type UserResourceRepositoryImpl struct {
	BaseRepositoryImpl[entity.UserResource]
}

// NewUserResourceRepository 创建用户上传资源Repository
func NewUserResourceRepository(db *gorm.DB) UserResourceRepository {
	return &UserResourceRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.UserResource]{db: db},
	}
}

// FindByIDAndUser 按ID与归属查询
func (r *UserResourceRepositoryImpl) FindByIDAndUser(id, userID uint) (*entity.UserResource, error) {
	var ur entity.UserResource
	err := r.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&ur).Error
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

// FindByUserAndURL 本人内按URL查重
func (r *UserResourceRepositoryImpl) FindByUserAndURL(userID uint, url string) (*entity.UserResource, error) {
	var ur entity.UserResource
	err := r.GetDB().Where("user_id = ? AND url = ?", userID, url).First(&ur).Error
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

// CountTodayByUser 当日提交计数
func (r *UserResourceRepositoryImpl) CountTodayByUser(userID uint) (int64, error) {
	now := utils.GetCurrentTime()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var count int64
	err := r.GetDB().Model(&entity.UserResource{}).
		Where("user_id = ? AND created_at >= ?", userID, todayStart).
		Count(&count).Error
	return count, err
}

// ListByUser 分页查询，支持状态与关键词过滤
func (r *UserResourceRepositoryImpl) ListByUser(userID uint, page, pageSize int, status, keyword string) ([]*entity.UserResource, int64, error) {
	var list []*entity.UserResource
	var total int64

	query := r.GetDB().Model(&entity.UserResource{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Pan").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

// GetStatsByUser 五态计数
func (r *UserResourceRepositoryImpl) GetStatsByUser(userID uint) (map[string]int64, error) {
	type statusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var counts []statusCount
	err := r.GetDB().Model(&entity.UserResource{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	stats := map[string]int64{
		entity.UserResourceStatusPending:    0,
		entity.UserResourceStatusValid:      0,
		entity.UserResourceStatusInvalid:    0,
		entity.UserResourceStatusProcessing: 0,
		entity.UserResourceStatusPublished:  0,
	}
	var total int64
	for _, c := range counts {
		stats[c.Status] = c.Count
		total += c.Count
	}
	stats["total"] = total
	return stats, nil
}
