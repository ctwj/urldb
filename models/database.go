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
		ck TEXT,
		is_valid BOOLEAN DEFAULT true,
		space BIGINT DEFAULT 0,
		left_space BIGINT DEFAULT 0,
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
		t VARCHAR(64) DEFAULT NULL,
		idx INTEGER DEFAULT NULL,
		ck TEXT,
		remark VARCHAR(64) NOT NULL
	);`

	// 创建待处理资源表
	readyResourceTable := `
	CREATE TABLE IF NOT EXISTS ready_resource (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255),
		url VARCHAR(500) NOT NULL,
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

	// 插入默认分类
	insertDefaultCategories := `
	INSERT INTO categories (name, description) VALUES 
	('文档', '各种文档资料'),
	('软件', '软件工具'),
	('视频', '视频教程'),
	('图片', '图片资源'),
	('音频', '音频文件'),
	('其他', '其他资源')
	ON CONFLICT (name) DO NOTHING;`

	if _, err := DB.Exec(insertDefaultCategories); err != nil {
		log.Printf("插入默认分类失败: %v", err)
	}

	return nil
}
