package services

import (
	"errors"
	"fmt"
	"time"

	"movePoint/internal/models"
	"movePoint/pkg/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Register 用户注册
func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// 检查邮箱是否已存在
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 检查用户名是否已存在
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已被使用")
	}

	// 处理生日字段 - 如果是空字符串或无效日期，设置为 nil
	var birthDatePtr *time.Time
	if req.BirthDate != "" {
		// 解析日期字符串
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			// 如果日期格式不正确，可以记录日志但继续注册流程
			fmt.Printf("警告: 生日格式不正确: %s\n", req.BirthDate)
			// 设置为 nil 而不是无效日期
			birthDatePtr = nil
		} else {
			birthDatePtr = &birthDate
		}
	}

	// 创建用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // BeforeCreate钩子会自动加密
		BirthDate: birthDatePtr,
		// 可以设置其他字段的默认值
		Weight:       0,
		Height:       0,
		AvatarURL:    "",
		Bio:          "",
		Achievements: "[]", // 默认空数组的 JSON 字符串
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// 返回认证响应
	response := &models.AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return response, nil
}

// Login 用户登录
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// 根据邮箱查找用户
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// 返回认证响应
	response := &models.AuthResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return response, nil
}
