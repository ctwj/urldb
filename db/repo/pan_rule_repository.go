package repo

import (
	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// PanRuleRepository PanRule的Repository接口
type PanRuleRepository interface {
	BaseRepository[entity.PanRule]
	FindByPanID(panID uint) ([]entity.PanRule, error)
	FindEnabledRules() ([]entity.PanRule, error)
}

// PanRuleRepositoryImpl PanRule的Repository实现
type PanRuleRepositoryImpl struct {
	BaseRepositoryImpl[entity.PanRule]
}

// NewPanRuleRepository 创建PanRule Repository
func NewPanRuleRepository(db *gorm.DB) PanRuleRepository {
	return &PanRuleRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.PanRule]{db: db},
	}
}

// FindByPanID 根据网盘ID查找规则
func (r *PanRuleRepositoryImpl) FindByPanID(panID uint) ([]entity.PanRule, error) {
	var rules []entity.PanRule
	err := r.db.Where("pan_id = ?", panID).Find(&rules).Error
	return rules, err
}

// FindEnabledRules 查找所有启用的规则
func (r *PanRuleRepositoryImpl) FindEnabledRules() ([]entity.PanRule, error) {
	var rules []entity.PanRule
	err := r.db.Preload("Pan").Where("enabled = ?", true).Order("priority ASC").Find(&rules).Error
	return rules, err
}