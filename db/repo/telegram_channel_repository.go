package repo

import (
	"time"

	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

type TelegramChannelRepository interface {
	BaseRepository[entity.TelegramChannel]
	FindActiveChannels() ([]entity.TelegramChannel, error)
	FindByChatID(chatID int64) (*entity.TelegramChannel, error)
	FindByChatType(chatType string) ([]entity.TelegramChannel, error)
	UpdateLastPushAt(id uint, lastPushAt time.Time) error
	FindDueForPush() ([]entity.TelegramChannel, error)
}

type TelegramChannelRepositoryImpl struct {
	BaseRepositoryImpl[entity.TelegramChannel]
}

func NewTelegramChannelRepository(db *gorm.DB) TelegramChannelRepository {
	return &TelegramChannelRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.TelegramChannel]{db: db},
	}
}

// 实现基类方法
func (r *TelegramChannelRepositoryImpl) Create(entity *entity.TelegramChannel) error {
	return r.db.Create(entity).Error
}

func (r *TelegramChannelRepositoryImpl) Update(entity *entity.TelegramChannel) error {
	return r.db.Save(entity).Error
}

func (r *TelegramChannelRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.TelegramChannel{}, id).Error
}

func (r *TelegramChannelRepositoryImpl) FindByID(id uint) (*entity.TelegramChannel, error) {
	var channel entity.TelegramChannel
	err := r.db.First(&channel, id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *TelegramChannelRepositoryImpl) FindAll() ([]entity.TelegramChannel, error) {
	var channels []entity.TelegramChannel
	err := r.db.Order("created_at desc").Find(&channels).Error
	return channels, err
}

// FindActiveChannels 查找活跃的频道/群组
func (r *TelegramChannelRepositoryImpl) FindActiveChannels() ([]entity.TelegramChannel, error) {
	var channels []entity.TelegramChannel
	err := r.db.Where("is_active = ? AND push_enabled = ?", true, true).Order("created_at desc").Find(&channels).Error
	return channels, err
}

// FindByChatID 根据 ChatID 查找频道/群组
func (r *TelegramChannelRepositoryImpl) FindByChatID(chatID int64) (*entity.TelegramChannel, error) {
	var channel entity.TelegramChannel
	err := r.db.Where("chat_id = ?", chatID).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// FindByChatType 根据类型查找频道/群组
func (r *TelegramChannelRepositoryImpl) FindByChatType(chatType string) ([]entity.TelegramChannel, error) {
	var channels []entity.TelegramChannel
	err := r.db.Where("chat_type = ?", chatType).Order("created_at desc").Find(&channels).Error
	return channels, err
}

// UpdateLastPushAt 更新最后推送时间
func (r *TelegramChannelRepositoryImpl) UpdateLastPushAt(id uint, lastPushAt time.Time) error {
	return r.db.Model(&entity.TelegramChannel{}).Where("id = ?", id).Update("last_push_at", lastPushAt).Error
}

// FindDueForPush 查找需要推送的频道/群组
func (r *TelegramChannelRepositoryImpl) FindDueForPush() ([]entity.TelegramChannel, error) {
	var channels []entity.TelegramChannel
	// 查找活跃、启用推送的频道，且距离上次推送已超过推送频率小时的记录
	err := r.db.Where("is_active = ? AND push_enabled = ? AND (last_push_at IS NULL OR last_push_at < DATE_SUB(NOW(), INTERVAL push_frequency HOUR))",
		true, true).Find(&channels).Error
	return channels, err
}
