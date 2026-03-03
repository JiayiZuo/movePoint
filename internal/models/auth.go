package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username     string  `json:"username" binding:"required,min=3,max=20"`
	Email        string  `json:"email" binding:"required,email"`
	Password     string  `json:"password" binding:"required,min=6"`
	BirthDate    string  `json:"birth_date"` // 格式: "2006-01-02"
	Weight       float64 `json:"weight"`
	Height       float64 `json:"height"`
	AvatarURL    string  `json:"avatar_url"`
	Bio          string  `json:"bio"`
	Achievements string  `json:"achievements"` // JSON 字符串格式
}

// AuthResponse 认证响应结构体
type AuthResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// WeChatLoginRequest 微信登录请求结构体
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// WeChatUserInfo 微信用户信息
type WeChatUserInfo struct {
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"headimgurl"`
	Gender    int    `json:"sex"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	UnionID   string `json:"unionid,omitempty"`
}

// WeChatAccessTokenResponse 微信获取access_token响应
type WeChatAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode,omitempty"`
	ErrMsg       string `json:"errmsg,omitempty"`
}

// WeChatUserInfoResponse 微信获取用户信息响应
type WeChatUserInfoResponse struct {
	OpenID    string `json:"openid"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"headimgurl"`
	Gender    int    `json:"sex"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	UnionID   string `json:"unionid,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
	ErrMsg    string `json:"errmsg,omitempty"`
}

// HashPassword 使用bcrypt加密密码
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 检查密码是否匹配
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// BeforeCreate Gorm钩子，在创建用户前自动加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		if err := u.HashPassword(u.Password); err != nil {
			return err
		}
	}
	return nil
}
