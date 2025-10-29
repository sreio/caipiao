package services

import (
	"caipiao/backend/database"
	"caipiao/backend/models"
)

// TrendService 走势图服务
type TrendService struct{}

// NewTrendService 创建走势图服务
func NewTrendService() *TrendService {
	return &TrendService{}
}

// ShuangseqiuTrendData 双色球走势数据
type ShuangseqiuTrendData struct {
	Issues      []string                `json:"issues"`       // 期号列表
	DrawDates   []string                `json:"draw_dates"`   // 开奖日期列表
	RedBalls    [][]int                 `json:"red_balls"`    // 每期的红球号码
	BlueBalls   []int                   `json:"blue_balls"`   // 每期的蓝球号码
	RedFreq     map[int]int             `json:"red_freq"`     // 红球频率统计
	BlueFreq    map[int]int             `json:"blue_freq"`    // 蓝球频率统计
	RedMissing  map[int]int             `json:"red_missing"`  // 红球遗漏值
	BlueMissing map[int]int             `json:"blue_missing"` // 蓝球遗漏值
	Records     []models.Shuangseqiu    `json:"records"`      // 原始记录
}

// DaletouTrendData 大乐透走势数据
type DaletouTrendData struct {
	Issues       []string             `json:"issues"`        // 期号列表
	DrawDates    []string             `json:"draw_dates"`    // 开奖日期列表
	FrontBalls   [][]int              `json:"front_balls"`   // 每期的前区号码
	BackBalls    [][]int              `json:"back_balls"`    // 每期的后区号码
	FrontFreq    map[int]int          `json:"front_freq"`    // 前区频率统计
	BackFreq     map[int]int          `json:"back_freq"`     // 后区频率统计
	FrontMissing map[int]int          `json:"front_missing"` // 前区遗漏值
	BackMissing  map[int]int          `json:"back_missing"`  // 后区遗漏值
	Records      []models.Daletou     `json:"records"`       // 原始记录
}

// GetShuangseqiuTrend 获取双色球走势数据
func (s *TrendService) GetShuangseqiuTrend(limit int) (*ShuangseqiuTrendData, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	db := database.GetDB()
	var records []models.Shuangseqiu

	// 获取最近N期数据，按期号倒序（最新在上）
	err := db.Order("issue DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, err
	}

	trend := &ShuangseqiuTrendData{
		Issues:      make([]string, 0),
		DrawDates:   make([]string, 0),
		RedBalls:    make([][]int, 0),
		BlueBalls:   make([]int, 0),
		RedFreq:     make(map[int]int),
		BlueFreq:    make(map[int]int),
		RedMissing:  make(map[int]int),
		BlueMissing: make(map[int]int),
		Records:     records,
	}

	// 初始化遗漏值
	for i := 1; i <= 33; i++ {
		trend.RedMissing[i] = 0
	}
	for i := 1; i <= 16; i++ {
		trend.BlueMissing[i] = 0
	}

	// 处理每期数据
	for _, record := range records {
		trend.Issues = append(trend.Issues, record.Issue)
		trend.DrawDates = append(trend.DrawDates, record.DrawDate.Format("2006-01-02"))

		// 红球
		redBalls := []int{
			record.RedBall1,
			record.RedBall2,
			record.RedBall3,
			record.RedBall4,
			record.RedBall5,
			record.RedBall6,
		}
		trend.RedBalls = append(trend.RedBalls, redBalls)

		// 蓝球
		trend.BlueBalls = append(trend.BlueBalls, record.BlueBall)

		// 统计频率
		for _, num := range redBalls {
			trend.RedFreq[num]++
		}
		trend.BlueFreq[record.BlueBall]++

		// 更新遗漏值
		// 先让所有号码遗漏值+1
		for i := 1; i <= 33; i++ {
			trend.RedMissing[i]++
		}
		for i := 1; i <= 16; i++ {
			trend.BlueMissing[i]++
		}
		// 出现的号码遗漏值归零
		for _, num := range redBalls {
			trend.RedMissing[num] = 0
		}
		trend.BlueMissing[record.BlueBall] = 0
	}

	return trend, nil
}

// GetDaletouTrend 获取大乐透走势数据
func (s *TrendService) GetDaletouTrend(limit int) (*DaletouTrendData, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	db := database.GetDB()
	var records []models.Daletou

	// 获取最近N期数据，按期号倒序（最新在上）
	err := db.Order("issue DESC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, err
	}

	trend := &DaletouTrendData{
		Issues:       make([]string, 0),
		DrawDates:    make([]string, 0),
		FrontBalls:   make([][]int, 0),
		BackBalls:    make([][]int, 0),
		FrontFreq:    make(map[int]int),
		BackFreq:     make(map[int]int),
		FrontMissing: make(map[int]int),
		BackMissing:  make(map[int]int),
		Records:      records,
	}

	// 初始化遗漏值
	for i := 1; i <= 35; i++ {
		trend.FrontMissing[i] = 0
	}
	for i := 1; i <= 12; i++ {
		trend.BackMissing[i] = 0
	}

	// 处理每期数据
	for _, record := range records {
		trend.Issues = append(trend.Issues, record.Issue)
		trend.DrawDates = append(trend.DrawDates, record.DrawDate.Format("2006-01-02"))

		// 前区
		frontBalls := []int{
			record.FrontBall1,
			record.FrontBall2,
			record.FrontBall3,
			record.FrontBall4,
			record.FrontBall5,
		}
		trend.FrontBalls = append(trend.FrontBalls, frontBalls)

		// 后区
		backBalls := []int{
			record.BackBall1,
			record.BackBall2,
		}
		trend.BackBalls = append(trend.BackBalls, backBalls)

		// 统计频率
		for _, num := range frontBalls {
			trend.FrontFreq[num]++
		}
		for _, num := range backBalls {
			trend.BackFreq[num]++
		}

		// 更新遗漏值
		// 先让所有号码遗漏值+1
		for i := 1; i <= 35; i++ {
			trend.FrontMissing[i]++
		}
		for i := 1; i <= 12; i++ {
			trend.BackMissing[i]++
		}
		// 出现的号码遗漏值归零
		for _, num := range frontBalls {
			trend.FrontMissing[num] = 0
		}
		for _, num := range backBalls {
			trend.BackMissing[num] = 0
		}
	}

	return trend, nil
}

