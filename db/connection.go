package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"res_db/db/entity"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "res_db"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移表结构
	if err := autoMigrate(); err != nil {
		return err
	}

	// 插入默认数据
	if err := insertDefaultData(); err != nil {
		log.Printf("插入默认数据失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// autoMigrate 自动迁移表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&entity.Pan{},
		&entity.Cks{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Resource{},
		&entity.ResourceTag{},
		&entity.ReadyResource{},
		&entity.User{},
		&entity.SearchStat{},
	)
}

// insertDefaultData 插入默认数据
func insertDefaultData() error {
	// 插入默认分类
	defaultCategories := []entity.Category{
		{Name: "文档", Description: "各种文档资料"},
		{Name: "软件", Description: "软件工具"},
		{Name: "视频", Description: "视频教程"},
		{Name: "图片", Description: "图片资源"},
		{Name: "音频", Description: "音频文件"},
		{Name: "其他", Description: "其他资源"},
	}

	for _, category := range defaultCategories {
		if err := DB.Where("name = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
			return err
		}
	}

	// 插入默认管理员用户
	defaultAdmin := entity.User{
		Username: "admin",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Email:    "admin@example.com",
		Role:     "admin",
		IsActive: true,
	}

	if err := DB.Where("username = ?", defaultAdmin.Username).FirstOrCreate(&defaultAdmin).Error; err != nil {
		return err
	}

	return nil
}
