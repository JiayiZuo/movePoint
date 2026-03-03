// cmd/migrate/main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"movePoint/internal/models" // 替换为你的实际模块名
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  警告：未找到 .env 文件，使用默认配置")
	}

	// 获取数据库连接字符串
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 默认配置，根据你的实际情况修改
		dsn = "root:123456@tcp(127.0.0.1:3306)/movepoint?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	log.Println("✅ 数据库连接成功")

	// 自动迁移所有模型
	log.Println("🔄 开始创建数据表...")

	err = db.AutoMigrate(
		// 在这里列出所有模型
		&models.User{},
		// &models.ClimbingRecord{},  // 添加其他模型
		// &models.Activity{},
		// &models.Comment{},
	)

	if err != nil {
		log.Fatalf("❌ 数据表创建失败: %v", err)
	}

	log.Println("✅ 所有数据表创建成功！")

	// 列出已创建的表
	var tables []string
	db.Raw("SHOW TABLES").Scan(&tables)
	log.Println("📋 已创建的表:", tables)
}
