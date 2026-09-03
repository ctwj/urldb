package repo

import (
	"errors"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"

	"gorm.io/gorm"
)

// DownloadHistoryRepository 用户下载历史Repository接口
type DownloadHistoryRepository interface {
	BaseRepository[entity.DownloadHistory]
	// Record 记录一次下载：同一用户同一资源仅保留一条，重复下载累加次数并刷新时间
	Record(userID, resourceID uint) error
	// ListByUser 按最近下载时间倒序分页查询
	ListByUser(userID uint, page, pageSize int) ([]*entity.DownloadHistory, int64, error)
	// StatsByUser 返回总数、今日、最近7天、最近30天下载人数（按资源去重）
	StatsByUser(userID uint) (total, today, thisWeek, thisMonth int64, err error)
	// DeleteByUser 清空指定用户全部下载历史，返回受影响行数
	DeleteByUser(userID uint) (int64, error)
	// DeleteByIDAndUser 删除指定用户的一条下载历史，返回受影响行数
	DeleteByIDAndUser(id, userID uint) (int64, error)
}

// DownloadHistoryRepositoryImpl 用户下载历史Repository实现
type DownloadHistoryRepositoryImpl struct {
	BaseRepositoryImpl[entity.DownloadHistory]
}

// NewDownloadHistoryRepository 创建用户下载历史Repository
func NewDownloadHistoryRepository(db *gorm.DB) DownloadHistoryRepository {
	return &DownloadHistoryRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.DownloadHistory]{db: db},
	}
}

// Record 记录一次下载
func (r *DownloadHistoryRepositoryImpl) Record(userID, resourceID uint) error {
	var history entity.DownloadHistory
	err := r.GetDB().Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&history).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			history = entity.DownloadHistory{
				UserID:        userID,
				ResourceID:    resourceID,
				DownloadCount: 1,
			}
			return r.GetDB().Create(&history).Error
		}
		return err
	}

	return r.GetDB().Model(&entity.DownloadHistory{}).Where("id = ?", history.ID).Updates(map[string]interface{}{
		"download_count": gorm.Expr("download_count + 1"),
		"updated_at":     utils.GetCurrentTime(),
	}).Error
}

// ListByUser 按最近下载时间倒序分页查询
func (r *DownloadHistoryRepositoryImpl) ListByUser(userID uint, page, pageSize int) ([]*entity.DownloadHistory, int64, error) {
	var histories []*entity.DownloadHistory
	var total int64

	query := r.GetDB().Model(&entity.DownloadHistory{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("updated_at DESC").Find(&histories).Error
	return histories, total, err
}

// StatsByUser 返回总数、今日、最近7天、最近30天下载数
func (r *DownloadHistoryRepositoryImpl) StatsByUser(userID uint) (total, today, thisWeek, thisMonth int64, err error) {
	now := utils.GetCurrentTime()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -6)
	monthStart := todayStart.AddDate(0, 0, -29)

	if err = r.GetDB().Model(&entity.DownloadHistory{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return
	}
	if err = r.GetDB().Model(&entity.DownloadHistory{}).Where("user_id = ? AND updated_at >= ?", userID, todayStart).Count(&today).Error; err != nil {
		return
	}
	if err = r.GetDB().Model(&entity.DownloadHistory{}).Where("user_id = ? AND updated_at >= ?", userID, weekStart).Count(&thisWeek).Error; err != nil {
		return
	}
	err = r.GetDB().Model(&entity.DownloadHistory{}).Where("user_id = ? AND updated_at >= ?", userID, monthStart).Count(&thisMonth).Error
	return
}

// DeleteByUser 清空指定用户全部下载历史
func (r *DownloadHistoryRepositoryImpl) DeleteByUser(userID uint) (int64, error) {
	result := r.GetDB().Where("user_id = ?", userID).Delete(&entity.DownloadHistory{})
	return result.RowsAffected, result.Error
}

// DeleteByIDAndUser 删除指定用户的一条下载历史
func (r *DownloadHistoryRepositoryImpl) DeleteByIDAndUser(id, userID uint) (int64, error) {
	result := r.GetDB().Where("id = ? AND user_id = ?", id, userID).Delete(&entity.DownloadHistory{})
	return result.RowsAffected, result.Error
}
