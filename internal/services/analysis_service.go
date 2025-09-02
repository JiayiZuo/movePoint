package services

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"movePoint/internal/models"
)

type AnalysisService struct {
	db *gorm.DB
}

func NewAnalysisService(db *gorm.DB) *AnalysisService {
	return &AnalysisService{db: db}
}

// AnalysisData 分析结果数据结构
type AnalysisData struct {
	Summary struct {
		TotalSessions int     `json:"total_sessions"`
		TotalDuration int     `json:"total_duration"` // 分钟
		TotalCalories float64 `json:"total_calories"`
		HighestGrade  string  `json:"highest_grade"`
	} `json:"summary"`

	GradeDistribution  map[string]GradeStats `json:"grade_distribution"`
	SuccessRateByGrade map[string]float64    `json:"success_rate_by_grade"`
	MonthlyTrends      []MonthlyStat         `json:"monthly_trends"`
	// 可以添加更多分析维度...
}

type GradeStats struct {
	Attempts int `json:"attempts"`
	Success  int `json:"success"`
}

type MonthlyStat struct {
	Month    string `json:"month"` // YYYY-MM
	Sessions int    `json:"sessions"`
	Duration int    `json:"duration"`
}

// GetClimbingAnalysis 获取用户攀岩数据分析
func (s *AnalysisService) GetClimbingAnalysis(userID uint, from, to time.Time) (*AnalysisData, error) {
	var records []models.ClimbingRecord

	// 查询时间范围内的记录
	query := s.db.Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, from, to)
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	// 生成分析数据
	analysis := s.analyzeRecords(records)

	// 缓存分析结果 (可选)
	go s.cacheAnalysis(userID, analysis)

	return analysis, nil
}

// analyzeRecords 分析记录数据
func (s *AnalysisService) analyzeRecords(records []models.ClimbingRecord) *AnalysisData {
	var data AnalysisData
	gradeStats := make(map[string]GradeStats)
	monthlyStats := make(map[string]MonthlyStat)

	// 初始化分析
	data.Summary.TotalSessions = len(records)
	data.GradeDistribution = make(map[string]GradeStats)
	data.SuccessRateByGrade = make(map[string]float64)

	for _, record := range records {
		// 汇总统计
		data.Summary.TotalDuration += record.Duration
		data.Summary.TotalCalories += record.Calories

		// 更新最高难度
		if isHigherGrade(record.Grade, data.Summary.HighestGrade) {
			data.Summary.HighestGrade = record.Grade
		}

		// 按难度等级统计
		if record.Grade != "" {
			stats := gradeStats[record.Grade]
			stats.Attempts++
			if record.Success {
				stats.Success++
			}
			gradeStats[record.Grade] = stats
		}

		// 按月统计
		month := record.StartTime.Format("2006-01")
		stat := monthlyStats[month]
		stat.Month = month
		stat.Sessions++
		stat.Duration += record.Duration
		monthlyStats[month] = stat
	}

	// 处理难度分布数据
	for grade, stats := range gradeStats {
		data.GradeDistribution[grade] = stats
		if stats.Attempts > 0 {
			data.SuccessRateByGrade[grade] = float64(stats.Success) / float64(stats.Attempts) * 100
		}
	}

	// 处理月度趋势数据
	for _, stat := range monthlyStats {
		data.MonthlyTrends = append(data.MonthlyTrends, stat)
	}

	return &data
}

// cacheAnalysis 缓存分析结果
func (s *AnalysisService) cacheAnalysis(userID uint, data *AnalysisData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to marshal analysis data: %v\n", err)
		return
	}

	analysis := models.ClimbingAnalysis{
		UserID: userID,
		Date:   time.Now(),
		Data:   string(jsonData),
	}

	s.db.Create(&analysis)
}

// isHigherGrade 比较两个难度等级
func isHigherGrade(grade1, grade2 string) bool {
	// 这里需要实现攀岩难度的比较逻辑
	// 简化处理: 只比较字符串长度 (实际应解析难度等级)
	if grade2 == "" {
		return true
	}
	return len(grade1) > len(grade2)
}
