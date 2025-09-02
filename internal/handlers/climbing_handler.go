package handlers

import (
	"net/http"
	"strconv"
	"time"

	"movePoint/internal/models"
	"movePoint/internal/services"

	"github.com/gin-gonic/gin"
)

type ClimbingHandler struct {
	service *services.ClimbingService
}

func NewClimbingHandler(service *services.ClimbingService) *ClimbingHandler {
	return &ClimbingHandler{service: service}
}

// CreateRecord 创建攀岩记录
func (h *ClimbingHandler) CreateRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var record models.ClimbingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := h.service.CreateRecord(userID.(uint), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建记录失败"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// GetRecords 获取用户记录列表
func (h *ClimbingHandler) GetRecords(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 解析时间范围参数
	var from, to time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		from, _ = time.Parse("2006-01-02", fromStr)
	}
	if toStr := c.Query("to"); toStr != "" {
		to, _ = time.Parse("2006-01-02", toStr)
	}

	records, total, err := h.service.GetUserRecords(userID.(uint), page, limit, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  records,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRecord 获取单个记录
func (h *ClimbingHandler) GetRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	record, err := h.service.GetRecordByID(userID.(uint), uint(recordID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateRecord 更新记录
func (h *ClimbingHandler) UpdateRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	var record models.ClimbingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	record.ID = uint(recordID)
	if err := h.service.UpdateRecord(userID.(uint), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新记录失败"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord 删除记录
func (h *ClimbingHandler) DeleteRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	recordID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	if err := h.service.DeleteRecord(userID.(uint), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "记录删除成功"})
}
