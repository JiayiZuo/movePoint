package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	// 用户基本信息
	Weight    float64   `json:"weight"` // 体重(kg)，用于计算热量消耗
	Height    float64   `json:"height"` // 身高(cm)
	BirthDate time.Time `json:"birth_date"`
	AvatarURL string    `json:"avatar_url"`           // 头像URL
	Bio       string    `gorm:"type:text" json:"bio"` // 个人简介

	// 成就系统
	Achievements string `gorm:"type:text" json:"achievements"` // JSON格式存储成就数据

	// 关联关系
	ClimbingRecords []ClimbingRecord `json:"climbing_records,omitempty"`
}

// Achievement 成就结构
type Achievement struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	UnlockedAt  time.Time `json:"unlocked_at"`
	Progress    float64   `json:"progress"` // 0-100表示进度
	Completed   bool      `json:"completed"`
}
