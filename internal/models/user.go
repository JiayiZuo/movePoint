// internal/models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;type:int unsigned" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 基础信息
	Username string `gorm:"column:username;uniqueIndex:idx_username,length:191;not null" json:"username"`
	Email    string `gorm:"column:email;uniqueIndex:idx_email,length:191;not null" json:"email"`
	Password string `gorm:"column:password;not null" json:"-"`

	// 身体数据
	Weight       float64    `gorm:"column:weight" json:"weight"`
	Height       float64    `gorm:"column:height" json:"height"`
	BirthDate    *time.Time `gorm:"column:birth_date" json:"birth_date"`
	AvatarURL    string     `gorm:"column:avatar_url" json:"avatar_url"`
	Bio          string     `gorm:"column:bio;type:text" json:"bio"`
	Achievements string     `gorm:"column:achievements;type:text" json:"achievements"`

	// 微信小程序登录相关字段（⚠️ 使用 column 显式指定列名）
	WeChatOpenID     *string `gorm:"column:wechat_openid;uniqueIndex:idx_wechat_openid,length:191" json:"wechat_openid,omitempty"`
	WeChatUnionID    *string `gorm:"column:wechat_unionid;uniqueIndex:idx_wechat_unionid,length:191" json:"wechat_unionid,omitempty"`
	WeChatSessionKey *string `gorm:"column:wechat_session_key;type:varchar(255)" json:"-"`
	WeChatPhone      *string `gorm:"column:wechat_phone;uniqueIndex:idx_wechat_phone,length:20" json:"phone,omitempty"`
	WeChatNickname   *string `gorm:"column:wechat_nickname;length:191" json:"nickname,omitempty"`
	WeChatGender     int     `gorm:"column:wechat_gender;default:0" json:"gender"`
	WeChatCountry    *string `gorm:"column:wechat_country;length:100" json:"country,omitempty"`
	WeChatProvince   *string `gorm:"column:wechat_province;length:100" json:"province,omitempty"`
	WeChatCity       *string `gorm:"column:wechat_city;length:100" json:"city,omitempty"`
	WeChatAvatar     *string `gorm:"column:wechat_avatar;type:text" json:"wechat_avatar,omitempty"`

	// 最后登录时间
	LastLoginAt *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`

	// 关联
	ClimbingRecords []ClimbingRecord `json:"climbing_records,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
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
