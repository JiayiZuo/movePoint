package handlers

import (
	"net/http"
	"time"

	"movePoint/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	service *services.AnalysisService
}

func NewAnalysisHandler(service *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{service: service}
}

// GetClimbingAnalysis 获取攀岩分析
func (h *AnalysisHandler) GetClimbingAnalysis(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 解析时间范围参数 (默认最近90天)
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")

	var from, to time.Time

	if fromStr == "" {
		from = time.Now().AddDate(0, -3, 0) // 默认3个月前
	} else {
		from, _ = time.Parse("2006-01-02", fromStr)
	}

	if toStr == "" {
		to = time.Now() // 默认到现在
	} else {
		to, _ = time.Parse("2006-01-02", toStr)
	}

	analysis, err := h.service.GetClimbingAnalysis(userID.(uint), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分析数据失败"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
