package models

import "time"

// ClimbingAnalysis 用于存储用户的分析结果（可缓存）
type ClimbingAnalysis struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	UserID uint      `gorm:"not null;index" json:"user_id"`
	Date   time.Time `json:"date"`                  // 分析日期
	Data   string    `gorm:"type:text" json:"data"` // JSON格式的分析数据
}
