package database

import (
	"caipiao/backend/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	// 确保数据目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&models.Shuangseqiu{},
		&models.Daletou{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 创建性能优化索引
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// createIndexes 创建性能优化索引
func createIndexes(db *gorm.DB) error {
	log.Println("开始创建性能优化索引...")

	// 双色球表索引
	indexes := []string{
		// 开奖日期索引（倒序，常用于列表查询）
		"CREATE INDEX IF NOT EXISTS idx_ssq_draw_date ON shuangseqiu(draw_date DESC)",
		// 创建时间索引
		"CREATE INDEX IF NOT EXISTS idx_ssq_created_at ON shuangseqiu(created_at DESC)",

		// 大乐透表索引
		// 开奖日期索引（倒序，常用于列表查询）
		"CREATE INDEX IF NOT EXISTS idx_dlt_draw_date ON daletou(draw_date DESC)",
		// 创建时间索引
		"CREATE INDEX IF NOT EXISTS idx_dlt_created_at ON daletou(created_at DESC)",
	}

	for _, sql := range indexes {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("执行索引创建失败 [%s]: %v", sql, err)
		}
	}

	log.Println("性能优化索引创建完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
