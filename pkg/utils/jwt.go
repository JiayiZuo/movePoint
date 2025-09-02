package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT声明结构体
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID uint, username, email string) (string, error) {
	// 从环境变量获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-key" // 生产环境不应使用默认值
	}

	// 设置令牌过期时间
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时后过期

	// 创建声明
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "climbing-app",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	// 从环境变量获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-key" // 生产环境不应使用默认值
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
