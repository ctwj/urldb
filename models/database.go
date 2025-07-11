package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

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

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// 创建表
	if err := createTables(); err != nil {
		return err
	}

	log.Println("数据库连接成功")
	return nil
}

// createTables 创建数据库表
func createTables() error {
	// 创建pan表
	panTable := `
	CREATE TABLE IF NOT EXISTS pan (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64) DEFAULT NULL,
		key INTEGER DEFAULT NULL,
		icon VARCHAR(128) DEFAULT NULL,
		remark VARCHAR(64) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建分类表
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建标签表
	tagTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建资源表 - 更新后的结构
	resourceTable := `
	CREATE TABLE IF NOT EXISTS resources (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		url VARCHAR(128),
		pan_id INTEGER REFERENCES pan(id) ON DELETE SET NULL,
		quark_url VARCHAR(500),
		file_size VARCHAR(100),
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		view_count INTEGER DEFAULT 0,
		is_valid BOOLEAN DEFAULT true,
		is_public BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建资源标签关联表
	resourceTagTable := `
	CREATE TABLE IF NOT EXISTS resource_tags (
		id SERIAL PRIMARY KEY,
		resource_id INTEGER NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(resource_id, tag_id)
	);`

	if _, err := DB.Exec(panTable); err != nil {
		return err
	}

	if _, err := DB.Exec(categoryTable); err != nil {
		return err
	}

	if _, err := DB.Exec(tagTable); err != nil {
		return err
	}

	if _, err := DB.Exec(resourceTable); err != nil {
		return err
	}

	if _, err := DB.Exec(resourceTagTable); err != nil {
		return err
	}

	// 创建cks表
	cksTable := `
	CREATE TABLE IF NOT EXISTS cks (
		id SERIAL PRIMARY KEY,
		pan_id INTEGER NOT NULL REFERENCES pan(id) ON DELETE CASCADE,
		idx INTEGER DEFAULT NULL,
		ck TEXT,
		is_valid BOOLEAN DEFAULT true,
		space BIGINT DEFAULT 0,
		left_space BIGINT DEFAULT 0,
		remark VARCHAR(64) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建待处理资源表
	readyResourceTable := `
	CREATE TABLE IF NOT EXISTS ready_resource (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255),
		url VARCHAR(500) NOT NULL,
		category VARCHAR(100) DEFAULT NULL,
		tags VARCHAR(500) DEFAULT NULL,
		img VARCHAR(500) DEFAULT NULL,
		source VARCHAR(100) DEFAULT NULL,
		extra TEXT DEFAULT NULL,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		ip VARCHAR(45) DEFAULT NULL
	);`

	// 创建搜索统计表
	searchStatTable := `
	CREATE TABLE IF NOT EXISTS search_stats (
		id SERIAL PRIMARY KEY,
		keyword VARCHAR(255) NOT NULL,
		count INTEGER DEFAULT 1,
		date DATE NOT NULL,
		ip VARCHAR(45),
		user_agent VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建系统配置表
	systemConfigTable := `
	CREATE TABLE IF NOT EXISTS system_configs (
		id SERIAL PRIMARY KEY,
		site_title VARCHAR(200) NOT NULL DEFAULT '网盘资源管理系统',
		site_description VARCHAR(500),
		keywords VARCHAR(500),
		author VARCHAR(100),
		copyright VARCHAR(200),
		auto_process_ready_resources BOOLEAN DEFAULT false,
		auto_process_interval INTEGER DEFAULT 30,
		auto_transfer_enabled BOOLEAN DEFAULT false,
		auto_fetch_hot_drama_enabled BOOLEAN DEFAULT false,
		page_size INTEGER DEFAULT 100,
		maintenance_mode BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建热播剧表
	hotDramaTable := `
	CREATE TABLE IF NOT EXISTS hot_dramas (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		rating DECIMAL(3,1) DEFAULT 0.0,
		year VARCHAR(10),
		directors VARCHAR(500),
		actors VARCHAR(1000),
		category VARCHAR(50),
		sub_type VARCHAR(50),
		source VARCHAR(50) DEFAULT 'douban',
		douban_id VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(panTable); err != nil {
		return err
	}

	if _, err := DB.Exec(cksTable); err != nil {
		return err
	}

	if _, err := DB.Exec(readyResourceTable); err != nil {
		return err
	}

	if _, err := DB.Exec(searchStatTable); err != nil {
		return err
	}

	if _, err := DB.Exec(systemConfigTable); err != nil {
		return err
	}

	if _, err := DB.Exec(hotDramaTable); err != nil {
		return err
	}

	// 插入默认分类
	insertDefaultCategories := `
	INSERT INTO categories (name, description) VALUES 
	('电影', '电影资源'),
	('电视剧', '电视剧资源'),
	('动漫', '动漫资源'),
	('音乐', '音乐资源'),
	('软件', '软件资源'),
	('PC游戏', 'PC游戏资源'),
	('手机游戏', '手机游戏'),
	('文档', '文档资源'),
	('短剧', '短剧'),
	('学习资源', '学习资源'),
	('视频教程', '学习资源'),
	('其他', '其他资源')
	ON CONFLICT (name) DO NOTHING;`

	// 插入默认网盘平台
	insertDefaultPans := `
	INSERT INTO pan (name, key, icon, remark) VALUES 
	('baidu', 1, '<i class="fas fa-cloud text-blue-500"></i>', '百度网盘'),
	('aliyun', 2, '<i class="fas fa-cloud text-orange-500"></i>', '阿里云盘'),
	('quark', 3, '<i class="fas fa-atom text-purple-500"></i>', '夸克网盘'),
	('xunlei', 5, '<i class="fas fa-bolt text-yellow-500"></i>', '迅雷云盘'),
	('lanzou', 7, '<i class="fas fa-cloud text-blue-400"></i>', '蓝奏云'),
	('123', 8, '<i class="fas fa-cloud text-red-500"></i>', '123云盘'),
	('ctfile', 9, '<i class="fas fa-folder text-yellow-600"></i>', '城通网盘'),
	('115', 10, '<i class="fas fa-cloud-upload-alt text-green-600"></i>', '115网盘'),
	('magnet', 11, '<i class="fas fa-magnet text-red-600"></i>', '磁力链接'),
	('uc', 12, '<i class="fas fa-cloud-download-alt text-purple-600"></i>', 'UC网盘'),
	('other', 13, '<i class="fas fa-cloud text-gray-500"></i>', '其他')
	ON CONFLICT (name) DO NOTHING;`

	if _, err := DB.Exec(insertDefaultCategories); err != nil {
		return err
	}

	if _, err := DB.Exec(insertDefaultPans); err != nil {
		return err
	}

	return nil
}
