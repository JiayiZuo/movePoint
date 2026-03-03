package services

import (
	"encoding/json"
	"fmt"
	"time"

	"movePoint/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUserProfile 获取用户个人信息
func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.Select("id", "username", "email", "weight", "height", "birth_date", "avatar_url", "bio", "achievements", "created_at").
		Where("id = ?", userID).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// UpdateUserProfile 更新用户个人信息
func (s *UserService) UpdateUserProfile(userID uint, updates map[string]interface{}) error {
	// 过滤允许更新的字段
	allowedFields := []string{"weight", "height", "birth_date", "avatar_url", "bio"}
	filteredUpdates := make(map[string]interface{})

	for key, value := range updates {
		for _, allowed := range allowedFields {
			if key == allowed {
				filteredUpdates[key] = value
				break
			}
		}
	}

	if len(filteredUpdates) == 0 {
		return fmt.Errorf("没有有效的更新字段")
	}

	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(filteredUpdates)
	return result.Error
}

// GetUserStats 获取用户统计数据
func (s *UserService) GetUserStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取总攀岩次数
	var totalSessions int64
	if err := s.db.Model(&models.ClimbingRecord{}).
		Where("user_id = ?", userID).
		Count(&totalSessions).Error; err != nil {
		return nil, err
	}
	stats["total_sessions"] = totalSessions

	// 获取总攀岩时长
	var totalDuration struct {
		Total int
	}
	if err := s.db.Model(&models.ClimbingRecord{}).
		Select("SUM(duration) as total").
		Where("user_id = ?", userID).
		Scan(&totalDuration).Error; err != nil {
		return nil, err
	}
	stats["total_duration"] = totalDuration.Total

	// 获取最高难度
	var highestGrade struct {
		Grade string
	}
	if err := s.db.Model(&models.ClimbingRecord{}).
		Select("grade").
		Where("user_id = ? AND success = true", userID).
		Order("grade DESC").
		First(&highestGrade).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	stats["highest_grade"] = highestGrade.Grade

	// 获取最近活动时间
	var lastActivity struct {
		LastTime time.Time
	}
	if err := s.db.Model(&models.ClimbingRecord{}).
		Select("MAX(start_time) as last_time").
		Where("user_id = ?", userID).
		Scan(&lastActivity).Error; err != nil {
		return nil, err
	}
	stats["last_activity"] = lastActivity.LastTime

	// 获取本周活动次数
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	var weeklySessions int64
	if err := s.db.Model(&models.ClimbingRecord{}).
		Where("user_id = ? AND start_time >= ?", userID, startOfWeek).
		Count(&weeklySessions).Error; err != nil {
		return nil, err
	}
	stats["weekly_sessions"] = weeklySessions

	return stats, nil
}

// GetUserAchievements 获取用户成就
func (s *UserService) GetUserAchievements(userID uint) ([]models.Achievement, error) {
	var user models.User
	if err := s.db.Select("achievements").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	var achievements []models.Achievement
	if user.Achievements != "" {
		if err := json.Unmarshal([]byte(user.Achievements), &achievements); err != nil {
			return nil, err
		}
	}

	return achievements, nil
}

// UpdateUserAchievements 更新用户成就
func (s *UserService) UpdateUserAchievements(userID uint, achievements []models.Achievement) error {
	achievementsJSON, err := json.Marshal(achievements)
	if err != nil {
		return err
	}

	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("achievements", string(achievementsJSON)).Error
}

// CheckAndUpdateAchievements 检查并更新用户成就
func (s *UserService) CheckAndUpdateAchievements(userID uint) error {
	// 获取用户当前的成就
	achievements, err := s.GetUserAchievements(userID)
	if err != nil {
		return err
	}

	// 如果用户还没有任何成就，初始化一些默认成就
	if len(achievements) == 0 {
		achievements = s.getDefaultAchievements()
	}

	// 获取用户统计数据
	stats, err := s.GetUserStats(userID)
	if err != nil {
		return err
	}

	// 检查每个成就的完成情况
	updated := false
	for i := range achievements {
		if achievements[i].Completed {
			continue // 已经完成的成就跳过
		}

		// 根据成就ID检查进度
		progress := s.checkAchievementProgress(achievements[i].ID, stats)
		if progress > achievements[i].Progress {
			achievements[i].Progress = progress
			if progress >= 100 {
				achievements[i].Completed = true
				achievements[i].UnlockedAt = time.Now()
			}
			updated = true
		}
	}

	// 如果有更新，保存回数据库
	if updated {
		return s.UpdateUserAchievements(userID, achievements)
	}

	return nil
}

// getDefaultAchievements 获取默认成就列表
func (s *UserService) getDefaultAchievements() []models.Achievement {
	return []models.Achievement{
		{
			ID:          "first_climb",
			Name:        "初试攀岩",
			Description: "完成第一次攀岩记录",
			Icon:        "🎯",
			Progress:    0,
			Completed:   false,
		},
		{
			ID:          "weekly_regular",
			Name:        "每周一爬",
			Description: "一周内攀岩3次",
			Icon:        "📅",
			Progress:    0,
			Completed:   false,
		},
		{
			ID:          "v4_climber",
			Name:        "V4征服者",
			Description: "成功完成一条V4难度的线路",
			Icon:        "🏆",
			Progress:    0,
			Completed:   false,
		},
		{
			ID:          "endurance_master",
			Name:        "耐力大师",
			Description: "单次攀岩时长超过2小时",
			Icon:        "⏱️",
			Progress:    0,
			Completed:   false,
		},
		{
			ID:          "social_climber",
			Name:        "社交攀岩者",
			Description: "分享10条攀岩记录到社区",
			Icon:        "👥",
			Progress:    0,
			Completed:   false,
		},
	}
}

// checkAchievementProgress 检查成就进度
func (s *UserService) checkAchievementProgress(achievementID string, stats map[string]interface{}) float64 {
	switch achievementID {
	case "first_climb":
		if total, ok := stats["total_sessions"].(int64); ok && total > 0 {
			return 100
		}
	case "weekly_regular":
		if weekly, ok := stats["weekly_sessions"].(int64); ok {
			progress := float64(weekly) / 3.0 * 100
			if progress > 100 {
				return 100
			}
			return progress
		}
	case "v4_climber":
		// 这里需要更复杂的逻辑来检查用户是否完成了V4难度
		// 简化处理：假设用户最高难度是V4
		if highest, ok := stats["highest_grade"].(string); ok && highest == "V4" {
			return 100
		}
	case "endurance_master":
		// 需要查询单次最长攀岩时长
		// 简化处理：假设用户有一次超过2小时的记录
		return 0 // 实际实现需要查询数据库
	case "social_climber":
		// 需要查询分享到社区的记录数量
		return 0 // 实际实现需要查询数据库
	}
	return 0
}
