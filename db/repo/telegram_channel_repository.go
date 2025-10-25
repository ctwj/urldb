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
	CleanupDuplicateChannels() error
	FindActiveChannelsByTypes(chatTypes []string) ([]entity.TelegramChannel, error)
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

// FindActiveChannelsByTypes 根据多个类型查找活跃频道/群组
func (r *TelegramChannelRepositoryImpl) FindActiveChannelsByTypes(chatTypes []string) ([]entity.TelegramChannel, error) {
	var channels []entity.TelegramChannel
	err := r.db.Where("chat_type IN (?) AND is_active = ?", chatTypes, true).Find(&channels).Error
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

	// 先获取所有活跃且启用推送的频道
	err := r.db.Where("is_active = ? AND push_enabled = ?", true, true).Find(&channels).Error
	if err != nil {
		return nil, err
	}

	// 在内存中过滤出需要推送的频道（更可靠的跨数据库方案）
	var dueChannels []entity.TelegramChannel
	now := time.Now()

	// 用于去重的map，以chat_id为键
	seenChatIDs := make(map[int64]bool)

	for _, channel := range channels {
		// 检查是否已经处理过这个chat_id（去重）
		if seenChatIDs[channel.ChatID] {
			continue
		}

		// 如果从未推送过，或者距离上次推送已超过推送频率小时
		isDue := false
		if channel.LastPushAt == nil {
			isDue = true
		} else {
			// 计算下次推送时间：上次推送时间 + 推送频率分钟
			nextPushTime := channel.LastPushAt.Add(time.Duration(channel.PushFrequency) * time.Minute)
			if now.After(nextPushTime) {
				isDue = true
			}
		}

		if isDue {
			dueChannels = append(dueChannels, channel)
			seenChatIDs[channel.ChatID] = true // 标记此chat_id已处理
		}
	}

	return dueChannels, nil
}

// CleanupDuplicateChannels 清理重复的频道记录，保留ID最小的记录
func (r *TelegramChannelRepositoryImpl) CleanupDuplicateChannels() error {
	// 使用SQL查询找出重复的chat_id，并删除除了ID最小外的所有记录
	query := `
		DELETE t1 FROM telegram_channels t1
		INNER JOIN (
			SELECT chat_id, MIN(id) as min_id
			FROM telegram_channels
			GROUP BY chat_id
			HAVING COUNT(*) > 1
		) t2 ON t1.chat_id = t2.chat_id
		WHERE t1.id > t2.min_id
	`

	return r.db.Exec(query).Error
}
