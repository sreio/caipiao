package services

import (
	"caipiao/database"
	"caipiao/models"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// RecommendService 推荐服务
type RecommendService struct{}

// NewRecommendService 创建推荐服务实例
func NewRecommendService() *RecommendService {
	return &RecommendService{}
}

// NumberFrequency 号码频率统计
type NumberFrequency struct {
	Number int `json:"number"`
	Count  int `json:"count"`
}

// Recommendation 推荐结果
type Recommendation struct {
	Numbers    []int    `json:"numbers"`     // 推荐号码
	BlueNumber int      `json:"blue_number"` // 蓝球（双色球）或后区号码
	Confidence int      `json:"confidence"`  // 置信度(0-100)
	Basis      []string `json:"basis"`       // 推荐依据
}

// ShuangseqiuRecommendation 双色球推荐结果
type ShuangseqiuRecommendation struct {
	RedBalls  []int    `json:"red_balls"`  // 推荐红球（6个）
	BlueBall  int      `json:"blue_ball"`  // 推荐蓝球
	Basis     []string `json:"basis"`      // 推荐依据
	HotBalls  []int    `json:"hot_balls"`  // 热门号码
	ColdBalls []int    `json:"cold_balls"` // 冷门号码
}

// DaletouRecommendation 大乐透推荐结果
type DaletouRecommendation struct {
	FrontBalls []int    `json:"front_balls"` // 推荐前区（5个）
	BackBalls  []int    `json:"back_balls"`  // 推荐后区（2个）
	Basis      []string `json:"basis"`       // 推荐依据
	HotBalls   []int    `json:"hot_balls"`   // 热门号码
	ColdBalls  []int    `json:"cold_balls"`  // 冷门号码
}

// GenerateShuangseqiuRecommendation 生成双色球推荐
func (s *RecommendService) GenerateShuangseqiuRecommendation(count int) ([]ShuangseqiuRecommendation, error) {
	if count <= 0 || count > 10 {
		count = 5 // 默认生成5注
	}

	// 获取最近100期数据
	var results []models.Shuangseqiu
	if err := database.DB.Order("draw_date DESC").Limit(100).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}

	if len(results) < 30 {
		return nil, fmt.Errorf("数据不足，至少需要30期数据")
	}

	// 分析热号和冷号
	hotReds := s.analyzeHotNumbers(results, "red", 10)
	coldReds := s.analyzeColdNumbers(results, "red", 10)
	hotBlues := s.analyzeHotNumbers(results, "blue", 5)

	recommendations := make([]ShuangseqiuRecommendation, count)

	for i := 0; i < count; i++ {
		var redBalls []int
		var basis []string

		// 策略1: 热号为主 (前3注)
		if i < 3 {
			redBalls = s.generateFromHot(hotReds, 6)
			basis = append(basis, "基于近期高频号码")
		} else if i < 6 {
			// 策略2: 冷热结合 (4-6注)
			redBalls = s.combineHotAndCold(hotReds, coldReds, 6)
			basis = append(basis, "冷热号码均衡组合")
		} else {
			// 策略3: 冷号为主 (7-10注)
			redBalls = s.generateFromCold(coldReds, 6)
			basis = append(basis, "基于长期遗漏号码")
		}

		// 排序红球
		sort.Ints(redBalls)

		// 选择蓝球（从热门蓝球中随机选择）
		blueBall := hotBlues[rand.Intn(len(hotBlues))]

		recommendations[i] = ShuangseqiuRecommendation{
			RedBalls:  redBalls,
			BlueBall:  blueBall,
			Basis:     basis,
			HotBalls:  hotReds,
			ColdBalls: coldReds,
		}
	}

	return recommendations, nil
}

// GenerateDaletouRecommendation 生成大乐透推荐
func (s *RecommendService) GenerateDaletouRecommendation(count int) ([]DaletouRecommendation, error) {
	if count <= 0 || count > 10 {
		count = 5 // 默认生成5注
	}

	// 获取最近100期数据
	var results []models.Daletou
	if err := database.DB.Order("draw_date DESC").Limit(100).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}

	if len(results) < 30 {
		return nil, fmt.Errorf("数据不足，至少需要30期数据")
	}

	// 分析热号和冷号
	hotFronts := s.analyzeHotNumbersDLT(results, "front", 10)
	coldFronts := s.analyzeColdNumbersDLT(results, "front", 10)
	hotBacks := s.analyzeHotNumbersDLT(results, "back", 5)

	recommendations := make([]DaletouRecommendation, count)

	for i := 0; i < count; i++ {
		var frontBalls []int
		var basis []string

		// 策略与双色球类似
		if i < 3 {
			frontBalls = s.generateFromHot(hotFronts, 5)
			basis = append(basis, "基于近期高频号码")
		} else if i < 6 {
			frontBalls = s.combineHotAndCold(hotFronts, coldFronts, 5)
			basis = append(basis, "冷热号码均衡组合")
		} else {
			frontBalls = s.generateFromCold(coldFronts, 5)
			basis = append(basis, "基于长期遗漏号码")
		}

		// 排序前区
		sort.Ints(frontBalls)

		// 选择后区（从热门后区中随机选择2个）
		backBalls := s.selectRandom(hotBacks, 2)
		sort.Ints(backBalls)

		recommendations[i] = DaletouRecommendation{
			FrontBalls: frontBalls,
			BackBalls:  backBalls,
			Basis:      basis,
			HotBalls:   hotFronts,
			ColdBalls:  coldFronts,
		}
	}

	return recommendations, nil
}

// analyzeHotNumbers 分析热门号码（双色球）
func (s *RecommendService) analyzeHotNumbers(results []models.Shuangseqiu, ballType string, topN int) []int {
	frequency := make(map[int]int)

	for _, r := range results {
		if ballType == "red" {
			frequency[r.RedBall1]++
			frequency[r.RedBall2]++
			frequency[r.RedBall3]++
			frequency[r.RedBall4]++
			frequency[r.RedBall5]++
			frequency[r.RedBall6]++
		} else if ballType == "blue" {
			frequency[r.BlueBall]++
		}
	}

	// 转换为切片并排序
	var freqList []NumberFrequency
	for num, count := range frequency {
		freqList = append(freqList, NumberFrequency{Number: num, Count: count})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count > freqList[j].Count
	})

	// 返回前N个热门号码
	var hotNumbers []int
	for i := 0; i < topN && i < len(freqList); i++ {
		hotNumbers = append(hotNumbers, freqList[i].Number)
	}

	return hotNumbers
}

// analyzeColdNumbers 分析冷门号码（双色球）
func (s *RecommendService) analyzeColdNumbers(results []models.Shuangseqiu, ballType string, topN int) []int {
	frequency := make(map[int]int)
	maxNum := 33 // 红球最大号码

	if ballType == "blue" {
		maxNum = 16
	}

	// 初始化所有号码
	for i := 1; i <= maxNum; i++ {
		frequency[i] = 0
	}

	// 统计出现次数
	for _, r := range results {
		if ballType == "red" {
			frequency[r.RedBall1]++
			frequency[r.RedBall2]++
			frequency[r.RedBall3]++
			frequency[r.RedBall4]++
			frequency[r.RedBall5]++
			frequency[r.RedBall6]++
		} else if ballType == "blue" {
			frequency[r.BlueBall]++
		}
	}

	// 转换为切片并排序（出现次数少的在前）
	var freqList []NumberFrequency
	for num, count := range frequency {
		freqList = append(freqList, NumberFrequency{Number: num, Count: count})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count < freqList[j].Count
	})

	// 返回前N个冷门号码
	var coldNumbers []int
	for i := 0; i < topN && i < len(freqList); i++ {
		coldNumbers = append(coldNumbers, freqList[i].Number)
	}

	return coldNumbers
}

// analyzeHotNumbersDLT 分析热门号码（大乐透）
func (s *RecommendService) analyzeHotNumbersDLT(results []models.Daletou, ballType string, topN int) []int {
	frequency := make(map[int]int)

	for _, r := range results {
		if ballType == "front" {
			frequency[r.FrontBall1]++
			frequency[r.FrontBall2]++
			frequency[r.FrontBall3]++
			frequency[r.FrontBall4]++
			frequency[r.FrontBall5]++
		} else if ballType == "back" {
			frequency[r.BackBall1]++
			frequency[r.BackBall2]++
		}
	}

	var freqList []NumberFrequency
	for num, count := range frequency {
		freqList = append(freqList, NumberFrequency{Number: num, Count: count})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count > freqList[j].Count
	})

	var hotNumbers []int
	for i := 0; i < topN && i < len(freqList); i++ {
		hotNumbers = append(hotNumbers, freqList[i].Number)
	}

	return hotNumbers
}

// analyzeColdNumbersDLT 分析冷门号码（大乐透）
func (s *RecommendService) analyzeColdNumbersDLT(results []models.Daletou, ballType string, topN int) []int {
	frequency := make(map[int]int)
	maxNum := 35 // 前区最大号码

	if ballType == "back" {
		maxNum = 12
	}

	for i := 1; i <= maxNum; i++ {
		frequency[i] = 0
	}

	for _, r := range results {
		if ballType == "front" {
			frequency[r.FrontBall1]++
			frequency[r.FrontBall2]++
			frequency[r.FrontBall3]++
			frequency[r.FrontBall4]++
			frequency[r.FrontBall5]++
		} else if ballType == "back" {
			frequency[r.BackBall1]++
			frequency[r.BackBall2]++
		}
	}

	var freqList []NumberFrequency
	for num, count := range frequency {
		freqList = append(freqList, NumberFrequency{Number: num, Count: count})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count < freqList[j].Count
	})

	var coldNumbers []int
	for i := 0; i < topN && i < len(freqList); i++ {
		coldNumbers = append(coldNumbers, freqList[i].Number)
	}

	return coldNumbers
}

// generateFromHot 从热门号码中生成
func (s *RecommendService) generateFromHot(hotNumbers []int, count int) []int {
	rand.Seed(time.Now().UnixNano())
	
	// 从热门号码中随机选择
	selected := make(map[int]bool)
	var result []int

	for len(result) < count && len(result) < len(hotNumbers) {
		idx := rand.Intn(len(hotNumbers))
		num := hotNumbers[idx]
		if !selected[num] {
			selected[num] = true
			result = append(result, num)
		}
	}

	return result
}

// generateFromCold 从冷门号码中生成
func (s *RecommendService) generateFromCold(coldNumbers []int, count int) []int {
	rand.Seed(time.Now().UnixNano())
	
	selected := make(map[int]bool)
	var result []int

	for len(result) < count && len(result) < len(coldNumbers) {
		idx := rand.Intn(len(coldNumbers))
		num := coldNumbers[idx]
		if !selected[num] {
			selected[num] = true
			result = append(result, num)
		}
	}

	return result
}

// combineHotAndCold 冷热结合生成
func (s *RecommendService) combineHotAndCold(hotNumbers, coldNumbers []int, count int) []int {
	rand.Seed(time.Now().UnixNano())
	
	selected := make(map[int]bool)
	var result []int

	// 一半热号，一半冷号
	hotCount := count / 2
	coldCount := count - hotCount

	// 选择热号
	for len(result) < hotCount && len(hotNumbers) > 0 {
		idx := rand.Intn(len(hotNumbers))
		num := hotNumbers[idx]
		if !selected[num] {
			selected[num] = true
			result = append(result, num)
		}
	}

	// 选择冷号
	for len(result) < count && len(coldNumbers) > 0 {
		idx := rand.Intn(len(coldNumbers))
		num := coldNumbers[idx]
		if !selected[num] {
			selected[num] = true
			result = append(result, num)
		}
	}

	return result
}

// selectRandom 随机选择指定数量的号码
func (s *RecommendService) selectRandom(numbers []int, count int) []int {
	rand.Seed(time.Now().UnixNano())
	
	selected := make(map[int]bool)
	var result []int

	for len(result) < count && len(result) < len(numbers) {
		idx := rand.Intn(len(numbers))
		num := numbers[idx]
		if !selected[num] {
			selected[num] = true
			result = append(result, num)
		}
	}

	return result
}
