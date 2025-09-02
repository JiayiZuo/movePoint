package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简单的认证中间件 (实际应用中应使用JWT等)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 简单检查Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		//token := parts[1]

		// 这里应该验证token的有效性并提取用户信息
		// 简化处理: 假设token是用户ID
		// userID, err := validateToken(token)

		// 示例: 直接将token作为用户ID (仅用于演示)
		// 实际应用中应使用JWT等安全方案
		c.Set("userID", uint(1)) // 示例用户ID

		c.Next()
	}
}
