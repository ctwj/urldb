package repo

import (
	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// ApiApplicationRepository API 开通申请单 Repository 接口（016-api-access-application）
type ApiApplicationRepository interface {
	BaseRepository[entity.ApiApplication]
	// FindLatestByUser 用户最近一条申请（个人中心状态展示；无申请返回 nil,nil）
	FindLatestByUser(userID uint) (*entity.ApiApplication, error)
	// ExistsPendingByUser 是否存在待审核申请（重复提交拦截，FR-003）
	ExistsPendingByUser(userID uint) (bool, error)
	// ListByStatus 管理端分页列表（status 为空或 "all" 时不过滤；keyword 匹配邮箱/用户名/用途；created_at 倒序）
	ListByStatus(status, keyword string, page, pageSize int) ([]*entity.ApiApplication, int64, error)
	// CountByStatus 按状态计数（管理端统计卡）
	CountByStatus() (map[string]int64, error)
}

// ApiApplicationRepositoryImpl API 开通申请单 Repository 实现
type ApiApplicationRepositoryImpl struct {
	BaseRepositoryImpl[entity.ApiApplication]
}

// NewApiApplicationRepository 创建 API 开通申请单 Repository
func NewApiApplicationRepository(db *gorm.DB) ApiApplicationRepository {
	return &ApiApplicationRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.ApiApplication]{db: db},
	}
}

// FindLatestByUser 用户最近一条申请
func (r *ApiApplicationRepositoryImpl) FindLatestByUser(userID uint) (*entity.ApiApplication, error) {
	var app entity.ApiApplication
	err := r.GetDB().Where("user_id = ?", userID).Order("id DESC").First(&app).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

// ExistsPendingByUser 是否存在待审核申请
func (r *ApiApplicationRepositoryImpl) ExistsPendingByUser(userID uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&entity.ApiApplication{}).
		Where("user_id = ? AND status = ?", userID, entity.ApiApplicationStatusPending).
		Count(&count).Error
	return count > 0, err
}

// ListByStatus 管理端分页列表
func (r *ApiApplicationRepositoryImpl) ListByStatus(status, keyword string, page, pageSize int) ([]*entity.ApiApplication, int64, error) {
	var list []*entity.ApiApplication
	var total int64

	query := r.GetDB().Model(&entity.ApiApplication{})
	if status != "" && status != "all" {
		query = query.Where("api_applications.status = ?", status)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		query = query.Joins("JOIN users ON users.id = api_applications.user_id").
			Where("api_applications.email ILIKE ? OR users.username ILIKE ? OR api_applications.purpose ILIKE ?", kw, kw, kw)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("api_applications.created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// CountByStatus 按状态计数
func (r *ApiApplicationRepositoryImpl) CountByStatus() (map[string]int64, error) {
	var rows []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if err := r.GetDB().Model(&entity.ApiApplication{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]int64, len(rows))
	for _, row := range rows {
		result[row.Status] = row.Count
	}
	return result, nil
}
