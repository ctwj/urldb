package db

import (
	"fmt"
	"os"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// 自动迁移数据库表结构
	err = DB.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Pan{},
		&entity.Cks{},
		&entity.Tag{},
		&entity.Resource{},
		&entity.ResourceTag{},
		&entity.ReadyResource{},
		&entity.SearchStat{},
		&entity.SystemConfig{},
		&entity.HotDrama{},
	)
	if err != nil {
		utils.Fatal("数据库迁移失败: %v", err)
	}

	// 创建索引以提高查询性能
	createIndexes(DB)

	// 插入默认数据（只在数据库为空时）
	if err := insertDefaultDataIfEmpty(); err != nil {
		utils.Error("插入默认数据失败: %v", err)
	}

	utils.Info("数据库连接成功")
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
		&entity.SystemConfig{},
		&entity.HotDrama{},
	)
}

// createIndexes 创建数据库索引以提高查询性能
func createIndexes(db *gorm.DB) {
	// 资源表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_title ON resources USING gin(to_tsvector('chinese', title))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_description ON resources USING gin(to_tsvector('chinese', description))")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_category_id ON resources(category_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_pan_id ON resources(pan_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_created_at ON resources(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_updated_at ON resources(updated_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_is_valid ON resources(is_valid)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_is_public ON resources(is_public)")

	// 搜索统计表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_search_stats_query ON search_stats(query)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_search_stats_created_at ON search_stats(created_at DESC)")

	// 热播剧表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_title ON hot_dramas(title)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_category ON hot_dramas(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_created_at ON hot_dramas(created_at DESC)")

	// 资源标签关联表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_tags_resource_id ON resource_tags(resource_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_tags_tag_id ON resource_tags(tag_id)")

	utils.Info("数据库索引创建完成")
}

// insertDefaultDataIfEmpty 只在数据库为空时插入默认数据
func insertDefaultDataIfEmpty() error {
	// 检查是否已有数据
	var panCount int64
	if err := DB.Model(&entity.Pan{}).Count(&panCount).Error; err != nil {
		return err
	}

	// 如果pan表已有数据，跳过插入
	if panCount > 0 {
		utils.Info("pan表已有数据，跳过默认数据插入")
		return nil
	}

	utils.Info("pan表为空，开始插入默认数据...")

	// 插入默认分类（使用FirstOrCreate避免重复）
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
			utils.Error("插入分类 %s 失败: %v", category.Name, err)
			// 继续执行，不因为单个分类失败而停止
		}
	}

	// 插入默认网盘平台（使用FirstOrCreate避免重复）
	defaultPans := []entity.Pan{
		{Name: "baidu", Key: 1, Icon: "<i class=\"fas fa-cloud text-blue-500\"></i>", Remark: "百度网盘"},
		{Name: "aliyun", Key: 2, Icon: "<i class=\"fas fa-cloud text-orange-500\"></i>", Remark: "阿里云盘"},
		{Name: "quark", Key: 3, Icon: "<i class=\"fas fa-atom text-purple-500\"></i>", Remark: "夸克网盘"},
		{Name: "tianyi", Key: 4, Icon: "<i class=\"fas fa-cloud text-cyan-500\"></i>", Remark: "天翼云盘"},
		{Name: "xunlei", Key: 5, Icon: "<i class=\"fas fa-bolt text-yellow-500\"></i>", Remark: "迅雷云盘"},
		{Name: "weiyun", Key: 6, Icon: "<i class=\"fas fa-cloud text-green-500\"></i>", Remark: "微云"},
		{Name: "lanzou", Key: 7, Icon: "<i class=\"fas fa-cloud text-blue-400\"></i>", Remark: "蓝奏云"},
		{Name: "123", Key: 8, Icon: "<i class=\"fas fa-cloud text-red-500\"></i>", Remark: "123云盘"},
		{Name: "onedrive", Key: 9, Icon: "<i class=\"fas fa-cloud text-blue-600\"></i>", Remark: "OneDrive"},
		{Name: "google", Key: 10, Icon: "<i class=\"fas fa-cloud text-green-600\"></i>", Remark: "Google云盘"},
		{Name: "ctfile", Key: 11, Icon: "<i class=\"fas fa-folder text-yellow-600\"></i>", Remark: "城通网盘"},
		{Name: "115", Key: 12, Icon: "<i class=\"fas fa-cloud-upload-alt text-green-600\"></i>", Remark: "115网盘"},
		{Name: "magnet", Key: 13, Icon: "<i class=\"fas fa-magnet text-red-600\"></i>", Remark: "磁力链接"},
		{Name: "uc", Key: 14, Icon: "<i class=\"fas fa-cloud-download-alt text-purple-600\"></i>", Remark: "UC网盘"},
		{Name: "other", Key: 15, Icon: "<i class=\"fas fa-cloud text-gray-500\"></i>", Remark: "其他"},
	}

	for _, pan := range defaultPans {
		if err := DB.Where("name = ?", pan.Name).FirstOrCreate(&pan).Error; err != nil {
			utils.Error("插入平台 %s 失败: %v", pan.Name, err)
			// 继续执行，不因为单个平台失败而停止
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

	if err := DB.Create(&defaultAdmin).Error; err != nil {
		return err
	}

	utils.Info("默认数据插入完成")
	return nil
}
