package services

import (
	_ "encoding/json"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"movePoint/internal/models"
)

type ClimbingService struct {
	db *gorm.DB
}

func NewClimbingService(db *gorm.DB) *ClimbingService {
	return &ClimbingService{db: db}
}

// CreateRecord 创建攀岩记录
func (s *ClimbingService) CreateRecord(userID uint, record *models.ClimbingRecord) error {
	// 计算持续时间和热量消耗
	duration := record.EndTime.Sub(record.StartTime)
	record.Duration = int(duration.Minutes())
	record.Calories = s.calculateCalories(userID, record.Type, record.Duration)

	// 设置用户ID
	record.UserID = userID

	// 保存到数据库
	result := s.db.Create(record)
	if result.Error != nil {
		return result.Error
	}

	// 创建记录后检查成就
	userService := NewUserService(s.db)
	go func() {
		err := userService.CheckAndUpdateAchievements(userID)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// GetUserRecords 获取用户的攀岩记录
func (s *ClimbingService) GetUserRecords(userID uint, page, limit int, from, to time.Time) ([]models.ClimbingRecord, int64, error) {
	var records []models.ClimbingRecord
	var total int64

	query := s.db.Where("user_id = ?", userID)

	// 时间范围过滤
	if !from.IsZero() {
		query = query.Where("start_time >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("start_time <= ?", to)
	}

	// 获取总数
	if err := query.Model(&models.ClimbingRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页获取记录
	offset := (page - 1) * limit
	if err := query.Order("start_time DESC").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetRecordByID 根据ID获取记录
func (s *ClimbingService) GetRecordByID(userID, recordID uint) (*models.ClimbingRecord, error) {
	var record models.ClimbingRecord
	result := s.db.Where("user_id = ? AND id = ?", userID, recordID).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, result.Error
	}
	return &record, nil
}

// UpdateRecord 更新记录
func (s *ClimbingService) UpdateRecord(userID uint, record *models.ClimbingRecord) error {
	// 验证记录属于该用户
	var existing models.ClimbingRecord
	result := s.db.Where("user_id = ? AND id = ?", userID, record.ID).First(&existing)
	if result.Error != nil {
		return result.Error
	}

	// 重新计算持续时间和热量消耗
	if !record.StartTime.IsZero() && !record.EndTime.IsZero() {
		duration := record.EndTime.Sub(record.StartTime)
		record.Duration = int(duration.Minutes())
		record.Calories = s.calculateCalories(userID, record.Type, record.Duration)
	}

	result = s.db.Model(&existing).Updates(record)
	return result.Error
}

// DeleteRecord 删除记录
func (s *ClimbingService) DeleteRecord(userID, recordID uint) error {
	result := s.db.Where("user_id = ? AND id = ?", userID, recordID).Delete(&models.ClimbingRecord{})
	return result.Error
}

// calculateCalories 估算热量消耗 (简化算法)
func (s *ClimbingService) calculateCalories(userID uint, climbingType models.ClimbingType, duration int) float64 {
	// 在实际应用中，这里应该查询用户体重并结合运动类型计算
	// 简化处理: 假设平均消耗为 8-12 kcal/min
	baseRate := 10.0 // 平均10 kcal/min

	// 根据攀岩类型调整
	if climbingType == models.Bouldering {
		baseRate = 12.0 // 抱石强度更大
	} else if climbingType == models.SportClimbing {
		baseRate = 8.0 // 难度攀登更持久但强度略低
	}

	return baseRate * float64(duration)
}
