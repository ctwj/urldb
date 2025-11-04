package repo

import (
	"github.com/ctwj/urldb/db/entity"
	"gorm.io/gorm"
)

// PluginDataRepository 插件数据Repository接口
type PluginDataRepository interface {
	BaseRepository[entity.PluginData]
	FindByPluginAndKey(pluginName, dataType, key string) (*entity.PluginData, error)
	FindByPluginAndType(pluginName, dataType string) ([]entity.PluginData, error)
	DeleteByPluginAndKey(pluginName, dataType, key string) error
	DeleteByPluginAndType(pluginName, dataType string) error
	DeleteExpired() (int64, error)
}

// PluginDataRepositoryImpl 插件数据Repository实现
type PluginDataRepositoryImpl struct {
	BaseRepositoryImpl[entity.PluginData]
}

// NewPluginDataRepository 创建插件数据Repository
func NewPluginDataRepository(db *gorm.DB) PluginDataRepository {
	return &PluginDataRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.PluginData]{db: db},
	}
}

// FindByPluginAndKey 根据插件名称、数据类型和键查找数据
func (r *PluginDataRepositoryImpl) FindByPluginAndKey(pluginName, dataType, key string) (*entity.PluginData, error) {
	var data entity.PluginData
	err := r.db.Where("plugin_name = ? AND data_type = ? AND data_key = ?", pluginName, dataType, key).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// FindByPluginAndType 根据插件名称和数据类型查找数据
func (r *PluginDataRepositoryImpl) FindByPluginAndType(pluginName, dataType string) ([]entity.PluginData, error) {
	var data []entity.PluginData
	err := r.db.Where("plugin_name = ? AND data_type = ?", pluginName, dataType).Find(&data).Error
	return data, err
}

// DeleteByPluginAndKey 根据插件名称、数据类型和键删除数据
func (r *PluginDataRepositoryImpl) DeleteByPluginAndKey(pluginName, dataType, key string) error {
	return r.db.Where("plugin_name = ? AND data_type = ? AND data_key = ?", pluginName, dataType, key).Delete(&entity.PluginData{}).Error
}

// DeleteByPluginAndType 根据插件名称和数据类型删除数据
func (r *PluginDataRepositoryImpl) DeleteByPluginAndType(pluginName, dataType string) error {
	return r.db.Where("plugin_name = ? AND data_type = ?", pluginName, dataType).Delete(&entity.PluginData{}).Error
}

// DeleteExpired 删除过期数据
func (r *PluginDataRepositoryImpl) DeleteExpired() (int64, error) {
	result := r.db.Where("expires_at IS NOT NULL AND expires_at < NOW()").Delete(&entity.PluginData{})
	return result.RowsAffected, result.Error
}