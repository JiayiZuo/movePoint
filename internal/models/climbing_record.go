package models

import (
	"gorm.io/gorm"
	"time"
)

// ClimbingType 攀岩类型枚举
type ClimbingType string

const (
	Bouldering    ClimbingType = "bouldering"     // 抱石
	SportClimbing ClimbingType = "sport_climbing" // 难度攀登
)

// AttemptRange 尝试次数枚举
type AttemptRange string

const (
	Flash     AttemptRange = "flash"  // 一次完成
	TwoThree  AttemptRange = "2-3"    // 2-3次
	FourSix   AttemptRange = "4-6"    // 4-6次
	SevenPlus AttemptRange = "7+"     // 7次以上
	Failed    AttemptRange = "failed" // 未完成
)

type ClimbingRecord struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UserID uint `gorm:"not null;index" json:"user_id"`

	// 基本记录信息
	Type      ClimbingType `gorm:"type:varchar(20);not null" json:"type"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Duration  int          `json:"duration"` // 单位: 分钟，由StartTime和EndTime计算得出

	// 线路信息
	Grade    string       `gorm:"type:varchar(10)" json:"grade"` // 如 "V4", "5.11a"
	Color    string       `gorm:"type:varchar(20)" json:"color"` // 抱石垫颜色
	Attempts AttemptRange `gorm:"type:varchar(10)" json:"attempts"`
	Success  bool         `json:"success"`                                     // 是否成功完成
	Rating   int          `gorm:"check:rating>=1 AND rating<=5" json:"rating"` // 1-5星评分

	// 位置和媒体
	Location  string `gorm:"type:varchar(255)" json:"location"`
	Notes     string `gorm:"type:text" json:"notes"`
	MediaURLs string `gorm:"type:text" json:"media_urls"` // JSON数组存储多个媒体URL

	// 计算字段
	Calories float64 `json:"calories"` // 估算的热量消耗
}
