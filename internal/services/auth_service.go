package services

import (
	"errors"

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

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // BeforeCreate钩子会自动加密
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
