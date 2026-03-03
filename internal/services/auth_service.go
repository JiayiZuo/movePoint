package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// Login 用户登录（微信小程序登录）
func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// 获取微信小程序配置
	appID := utils.GetEnv("WECHAT_MINIAPP_APPID", "")
	appSecret := utils.GetEnv("WECHAT_MINIAPP_SECRET", "")

	if appID == "" || appSecret == "" {
		return nil, errors.New("微信小程序配置缺失")
	}

	// 向微信服务器请求获取session_key和openid
	weChatURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, req.Code)

	resp, err := http.Get(weChatURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信服务器失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信响应失败: %v", err)
	}

	// 解析微信返回的数据
	var weChatResp struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid,omitempty"`
		ErrCode    int    `json:"errcode,omitempty"`
		ErrMsg     string `json:"errmsg,omitempty"`
	}

	if err := json.Unmarshal(body, &weChatResp); err != nil {
		return nil, fmt.Errorf("解析微信响应失败: %v", err)
	}

	// 检查微信返回的错误
	if weChatResp.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", weChatResp.ErrMsg)
	}

	// 检查用户是否已存在（通过OpenID）
	var existingUser models.User
	if err := s.db.Where("wechat_openid = ?", weChatResp.OpenID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，创建新用户
			user := models.User{
				Username:     fmt.Sprintf("wx_user_%s", weChatResp.OpenID[:8]),     // 使用OpenID的一部分作为用户名
				Email:        fmt.Sprintf("%s@wechat.user", weChatResp.OpenID[:8]), // 使用OpenID的一部分作为邮箱
				Password:     "",                                                   // 微信登录用户不需要密码
				WeChatOpenID: &weChatResp.OpenID,                                   // 设置微信OpenID
				// 可以设置其他字段的默认值
				Weight:       0,
				Height:       0,
				AvatarURL:    "",
				Bio:          "通过微信登录的用户",
				Achievements: "[]", // 默认空数组的 JSON 字符串
			}

			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("创建用户失败: %v", err)
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
		} else {
			return nil, fmt.Errorf("查询用户失败: %v", err)
		}
	} else {
		// 用户已存在，直接生成令牌
		token, err := utils.GenerateJWT(existingUser.ID, existingUser.Username, existingUser.Email)
		if err != nil {
			return nil, err
		}

		// 返回认证响应
		response := &models.AuthResponse{
			UserID:   existingUser.ID,
			Username: existingUser.Username,
			Email:    existingUser.Email,
			Token:    token,
		}

		return response, nil
	}
}

// WeChatLogin 微信登录（保留原方法，但实际业务逻辑已在Login中实现）
func (s *AuthService) WeChatLogin(req *models.WeChatLoginRequest) (*models.AuthResponse, error) {
	// 直接复用Login方法的逻辑
	loginReq := &models.LoginRequest{Code: req.Code}
	return s.Login(loginReq)
}
