package handlers

import (
	"net/http"

	"movePoint/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 获取用户个人信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	user, err := h.userService.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile 更新用户个人信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.userService.UpdateUserProfile(userID.(uint), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// GetStats 获取用户统计数据
func (h *UserHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	stats, err := h.userService.GetUserStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAchievements 获取用户成就
func (h *UserHandler) GetAchievements(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	achievements, err := h.userService.GetUserAchievements(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户成就失败"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// CheckAchievements 检查并更新用户成就
func (h *UserHandler) CheckAchievements(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	if err := h.userService.CheckAndUpdateAchievements(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查成就失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成就检查完成"})
}
