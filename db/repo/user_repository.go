package repo

import (
	"github.com/ctwj/urldb/db/entity"

	"gorm.io/gorm"
)

// UserRepository 用户Repository接口
type UserRepository interface {
	BaseRepository[entity.User]
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	UpdateLastLogin(id uint) error
	FindByRole(role string) ([]entity.User, error)
}

// UserRepositoryImpl 用户Repository实现
type UserRepositoryImpl struct {
	BaseRepositoryImpl[entity.User]
}

// NewUserRepository 创建用户Repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entity.User]{db: db},
	}
}

// FindByUsername 根据用户名查找
func (r *UserRepositoryImpl) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找
func (r *UserRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepositoryImpl) UpdateLastLogin(id uint) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).
		UpdateColumn("last_login", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// FindByRole 根据角色查找用户
func (r *UserRepositoryImpl) FindByRole(role string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	return users, err
}
