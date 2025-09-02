package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey;type:int unsigned" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Username string `gorm:"uniqueIndex:idx_username,length:191;not null" json:"username"`
	Email    string `gorm:"uniqueIndex:idx_email,length:191;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	Weight       float64    `json:"weight"`
	Height       float64    `json:"height"`
	BirthDate    *time.Time `json:"birth_date"`
	AvatarURL    string     `json:"avatar_url"`
	Bio          string     `gorm:"type:text" json:"bio"`
	Achievements string     `gorm:"type:text" json:"achievements"`

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
