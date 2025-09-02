package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse 认证响应结构体
type AuthResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
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
