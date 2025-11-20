package services

import (
	"caipiao/database"
	"caipiao/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

// LotteryService 彩票服务
type LotteryService struct {
	shuangseqiuURL string
	daletouURL     string
	client         *resty.Client
}

// NewLotteryService 创建彩票服务实例
func NewLotteryService(shuangseqiuURL, daletouURL string) *LotteryService {
	// 创建 Resty 客户端，模拟完整的浏览器请求
	client := resty.New().
		SetTimeout(30*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(2*time.Second). // 增加重试等待时间，避免被限流
		SetRetryMaxWaitTime(10*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetHeader("Connection", "keep-alive").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("Pragma", "no-cache").
		SetHeader("Sec-Fetch-Dest", "empty").
		SetHeader("Sec-Fetch-Mode", "cors").
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`).
		SetHeader("sec-ch-ua-mobile", "?0").
		SetHeader("sec-ch-ua-platform", `"macOS"`)

	return &LotteryService{
		shuangseqiuURL: shuangseqiuURL,
		daletouURL:     daletouURL,
		client:         client,
	}
}

// ShuangseqiuResponse 双色球API响应（官方接口）
type ShuangseqiuResponse struct {
	State   int                 `json:"state"`
	Message string              `json:"message"`
	Total   int                 `json:"total"`
	Tflag   int                 `json:"Tflag"`
	Result  []ShuangseqiuResult `json:"result"`
}

// ShuangseqiuResult 双色球结果
type ShuangseqiuResult struct {
	Name        string       `json:"name"`        // 游戏名称
	Code        string       `json:"code"`        // 期号
	DetailsLink string       `json:"detailsLink"` // 详情链接
	VideoLink   string       `json:"videoLink"`   // 视频链接
	Date        string       `json:"date"`        // 开奖日期
	Week        string       `json:"week"`        // 星期
	Red         string       `json:"red"`         // 红球号码（逗号分隔）
	Blue        string       `json:"blue"`        // 蓝球号码
	Blue2       string       `json:"blue2"`       // 蓝球2（如果有）
	Sales       string       `json:"sales"`       // 销售额
	PoolMoney   string       `json:"poolmoney"`   // 奖池金额
	Content     string       `json:"content"`     // 详细内容
	AddContent  string       `json:"addcontent"`  // 附加内容
	Msg         string       `json:"msg"`         // 消息
	Z2Add       string       `json:"z2add"`       // 附加信息
	M2Add       string       `json:"m2add"`       // 附加信息
	Prizegrades []PrizeGrade `json:"prizegrades"` // 奖级信息
}

// PrizeGrade 奖级信息
type PrizeGrade struct {
	Type      int    `json:"type"`      // 奖级类型
	TypeNum   string `json:"typenum"`   // 中奖注数
	TypeMoney string `json:"typemoney"` // 单注奖金
}

// DaletouResponse 大乐透API响应（官方接口）
type DaletouResponse struct {
	Success   bool         `json:"success"`
	ErrorCode string       `json:"errorCode"`
	Message   string       `json:"message"`
	Value     DaletouValue `json:"value"`
}

// DaletouValue 大乐透数据
type DaletouValue struct {
	Total int             `json:"total"`
	Num   int             `json:"num"`
	Pages int             `json:"pages"`
	List  []DaletouResult `json:"list"`
}

// DaletouResult 大乐透结果
type DaletouResult struct {
	LotteryDrawNum          string        `json:"lotteryDrawNum"`          // 期号
	LotteryDrawTime         string        `json:"lotteryDrawTime"`         // 开奖时间
	LotteryDrawResult       string        `json:"lotteryDrawResult"`       // 开奖结果
	LotteryUnsortDrawresult string        `json:"lotteryUnsortDrawresult"` // 未排序结果
	LotterySaleAmount       string        `json:"lotterySaleAmount"`       // 销售额
	LotteryPoolAmount       string        `json:"lotteryPoolAmount"`       // 奖池金额
	LotteryDrawStatus       int           `json:"lotteryDrawStatus"`       // 开奖状态
	MatchList               []interface{} `json:"matchList"`               // 比赛列表（可能是数组或null）
	Remark                  string        `json:"remark"`                  // 备注
}

// FetchShuangseqiu 从外部API获取双色球数据
func (s *LotteryService) FetchShuangseqiu(issue string) (*models.Shuangseqiu, error) {
	// 构建请求参数
	params := map[string]string{
		"name":       "ssq",
		"issueCount": "1",
	}
	if issue != "" {
		params["issueStart"] = issue
		params["issueEnd"] = issue
	}

	// 使用 Resty 发送请求并自动解析 JSON
	var apiResp ShuangseqiuResponse
	resp, err := s.client.R().
		SetHeader("Referer", "https://www.cwl.gov.cn/").
		SetHeader("Origin", "https://www.cwl.gov.cn").
		SetQueryParams(params).
		SetResult(&apiResp).
		Get(s.shuangseqiuURL)

	if err != nil {
		log.Printf("双色球API请求失败: %v", err)
		return nil, fmt.Errorf("请求API失败: %v", err)
	}

	if !resp.IsSuccess() {
		log.Printf("双色球API返回错误状态码: %d", resp.StatusCode())
		return nil, fmt.Errorf("API请求失败: HTTP %d", resp.StatusCode())
	}

	log.Printf("双色球API响应: State=%d, Message=%s, ResultCount=%d", apiResp.State, apiResp.Message, len(apiResp.Result))

	if apiResp.State != 0 || len(apiResp.Result) == 0 {
		return nil, fmt.Errorf("API返回错误: %s (state=%d, results=%d)", apiResp.Message, apiResp.State, len(apiResp.Result))
	}

	// 解析第一条结果
	result := apiResp.Result[0]
	log.Printf("双色球API原始数据: 期号=%s, 红球=%s, 蓝球=%s, 日期字符串='%s'", result.Code, result.Red, result.Blue, result.Date)

	// 解析红球（逗号分隔的字符串）
	redBalls := []int{}
	for _, s := range strings.Split(result.Red, ",") {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			log.Printf("解析红球失败: %s, 错误: %v", s, err)
			continue
		}
		redBalls = append(redBalls, num)
	}

	// 验证红球数量
	if len(redBalls) != 6 {
		return nil, fmt.Errorf("红球数量错误: 期望6个，实际%d个", len(redBalls))
	}

	// 解析蓝球
	blueBall, err := strconv.Atoi(strings.TrimSpace(result.Blue))
	if err != nil {
		return nil, fmt.Errorf("解析蓝球失败: %v", err)
	}

	// 解析日期，支持多种格式
	// 先清理日期字符串，去掉括号中的内容（如：2024-08-20(二) -> 2024-08-20）
	dateStr := result.Date
	if idx := strings.Index(dateStr, "("); idx != -1 {
		dateStr = strings.TrimSpace(dateStr[:idx])
	}

	var drawDate time.Time
	dateFormats := []string{
		"2006-01-02",          // 2024-01-01
		"2006-01-02 15:04:05", // 2024-01-01 20:30:00
		"20060102",            // 20240101
	}

	parsed := false
	for _, format := range dateFormats {
		if d, err := time.Parse(format, dateStr); err == nil {
			drawDate = d
			parsed = true
			break
		}
	}

	if !parsed {
		log.Printf("解析日期失败: 原始='%s', 清理后='%s', 使用当前日期", result.Date, dateStr)
		drawDate = time.Now().Truncate(24 * time.Hour) // 使用当前日期（不含时间）
	}

	ssq := &models.Shuangseqiu{
		Issue:    result.Code,
		RedBall1: redBalls[0],
		RedBall2: redBalls[1],
		RedBall3: redBalls[2],
		RedBall4: redBalls[3],
		RedBall5: redBalls[4],
		RedBall6: redBalls[5],
		BlueBall: blueBall,
		DrawDate: drawDate,
	}

	log.Printf("双色球数据解析成功: %+v", ssq)
	return ssq, nil
}

// FetchDaletou 从外部API获取大乐透数据
func (s *LotteryService) FetchDaletou(issue string) (*models.Daletou, error) {
	// 构建请求参数
	params := map[string]string{
		"gameNo":     "85", // 大乐透游戏编号
		"provinceId": "0",
		"pageSize":   "1",
		"isVerify":   "1",
		"pageNo":     "1",
	}
	if issue != "" {
		params["issueStart"] = issue
		params["issueEnd"] = issue
	}

	// 使用 Resty 发送请求并自动解析 JSON
	var apiResp DaletouResponse
	resp, err := s.client.R().
		SetHeader("Referer", "https://www.sporttery.cn/").
		SetHeader("Origin", "https://www.sporttery.cn").
		SetQueryParams(params).
		SetResult(&apiResp).
		Get(s.daletouURL)

	if err != nil {
		log.Printf("大乐透API请求失败: %v", err)
		return nil, fmt.Errorf("请求API失败: %v", err)
	}

	if !resp.IsSuccess() {
		log.Printf("大乐透API返回错误状态码: %d, 响应内容: %s", resp.StatusCode(), string(resp.Body()))
		return nil, fmt.Errorf("API请求失败: HTTP %d", resp.StatusCode())
	}

	log.Printf("大乐透API响应: Success=%v, ErrorCode=%s, Message=%s, ListCount=%d",
		apiResp.Success, apiResp.ErrorCode, apiResp.Message, len(apiResp.Value.List))

	if !apiResp.Success {
		return nil, fmt.Errorf("API返回错误: %s (code=%s)", apiResp.Message, apiResp.ErrorCode)
	}

	if len(apiResp.Value.List) == 0 {
		return nil, fmt.Errorf("API未返回数据")
	}

	// 解析第一条结果
	result := apiResp.Value.List[0]
	log.Printf("大乐透数据: 期号=%s, 开奖结果=%s, 时间=%s",
		result.LotteryDrawNum, result.LotteryDrawResult, result.LotteryDrawTime)

	// 解析开奖结果
	// 支持两种格式：
	// 格式1: "01 02 03 04 05 # 01 02" (带#分隔符)
	// 格式2: "01 02 03 04 05 06 07" (7个数字，前5个是前区，后2个是后区)
	var frontBalls, backBalls []int

	if strings.Contains(result.LotteryDrawResult, "#") {
		// 格式1: 带 # 分隔符
		parts := strings.Split(result.LotteryDrawResult, "#")
		if len(parts) != 2 {
			return nil, fmt.Errorf("开奖结果格式错误: 期望包含#分隔符且分为2部分，实际值=%s", result.LotteryDrawResult)
		}

		// 解析前区
		for _, s := range strings.Fields(strings.TrimSpace(parts[0])) {
			num, err := strconv.Atoi(s)
			if err != nil {
				log.Printf("解析前区球号失败: %s, 错误: %v", s, err)
				continue
			}
			frontBalls = append(frontBalls, num)
		}

		// 解析后区
		for _, s := range strings.Fields(strings.TrimSpace(parts[1])) {
			num, err := strconv.Atoi(s)
			if err != nil {
				log.Printf("解析后区球号失败: %s, 错误: %v", s, err)
				continue
			}
			backBalls = append(backBalls, num)
		}
	} else {
		// 格式2: 无 # 分隔符，直接解析所有数字
		allBalls := []int{}
		for _, s := range strings.Fields(strings.TrimSpace(result.LotteryDrawResult)) {
			num, err := strconv.Atoi(s)
			if err != nil {
				log.Printf("解析球号失败: %s, 错误: %v", s, err)
				continue
			}
			allBalls = append(allBalls, num)
		}

		// 验证总数量
		if len(allBalls) != 7 {
			return nil, fmt.Errorf("开奖结果数量错误: 期望7个数字(前5+后2)，实际%d个 (原始数据: %s)",
				len(allBalls), result.LotteryDrawResult)
		}

		// 前5个是前区，后2个是后区
		frontBalls = allBalls[:5]
		backBalls = allBalls[5:]
		log.Printf("解析格式2: 前区=%v, 后区=%v", frontBalls, backBalls)
	}

	// 验证前区数量
	if len(frontBalls) != 5 {
		return nil, fmt.Errorf("前区球号数量错误: 期望5个，实际%d个", len(frontBalls))
	}

	// 验证后区数量
	if len(backBalls) != 2 {
		return nil, fmt.Errorf("后区球号数量错误: 期望2个，实际%d个", len(backBalls))
	}

	// 解析日期，支持多种格式
	var drawDate time.Time

	// 尝试格式1: "2024-01-01 20:30:00"
	drawDate, err = time.Parse("2006-01-02 15:04:05", result.LotteryDrawTime)
	if err != nil {
		// 尝试格式2: "2024-01-01"
		drawDate, err = time.Parse("2006-01-02", result.LotteryDrawTime)
		if err != nil {
			log.Printf("解析日期失败: %s, 错误: %v", result.LotteryDrawTime, err)
			drawDate = time.Now() // 使用当前时间作为默认值
		}
	}

	dlt := &models.Daletou{
		Issue:      result.LotteryDrawNum,
		FrontBall1: frontBalls[0],
		FrontBall2: frontBalls[1],
		FrontBall3: frontBalls[2],
		FrontBall4: frontBalls[3],
		FrontBall5: frontBalls[4],
		BackBall1:  backBalls[0],
		BackBall2:  backBalls[1],
		DrawDate:   drawDate,
	}

	log.Printf("大乐透数据解析成功: %+v", dlt)
	return dlt, nil
}

// SaveShuangseqiu 保存双色球数据
func (s *LotteryService) SaveShuangseqiu(ssq *models.Shuangseqiu) error {
	db := database.GetDB()

	// 检查期号是否已存在
	var existing models.Shuangseqiu
	err := db.Where("issue = ?", ssq.Issue).First(&existing).Error
	if err == nil {
		// 记录已存在
		return fmt.Errorf("期号 %s 的数据已存在", ssq.Issue)
	}

	// 如果是其他错误（非记录不存在），返回错误
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 保存新数据
	return db.Create(ssq).Error
}

// SaveDaletou 保存大乐透数据
func (s *LotteryService) SaveDaletou(dlt *models.Daletou) error {
	db := database.GetDB()

	// 检查期号是否已存在
	var existing models.Daletou
	err := db.Where("issue = ?", dlt.Issue).First(&existing).Error
	if err == nil {
		// 记录已存在
		return fmt.Errorf("期号 %s 的数据已存在", dlt.Issue)
	}

	// 如果是其他错误（非记录不存在），返回错误
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 保存新数据
	return db.Create(dlt).Error
}

// BatchResult 批量获取结果
type BatchResult struct {
	Total   int    `json:"total"`   // 总数
	Success int    `json:"success"` // 成功
	Skipped int    `json:"skipped"` // 跳过（已存在）
	Failed  int    `json:"failed"`  // 失败
	TaskID  string `json:"task_id"` // 任务ID（异步时使用）
}

// GetLatestIssue 获取数据库中最新的期号
func (s *LotteryService) GetLatestIssue(lotteryType string) (string, error) {
	db := database.GetDB()
	var issue string

	if lotteryType == "shuangseqiu" {
		var ssq models.Shuangseqiu
		err := db.Order("issue DESC").First(&ssq).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return "", nil // 数据库为空
			}
			return "", err
		}
		issue = ssq.Issue
	} else if lotteryType == "daletou" {
		var dlt models.Daletou
		err := db.Order("issue DESC").First(&dlt).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return "", nil // 数据库为空
			}
			return "", err
		}
		issue = dlt.Issue
	}

	return issue, nil
}

// FetchShuangseqiuHistoryAsync 异步批量获取双色球历史数据
func (s *LotteryService) FetchShuangseqiuHistoryAsync(count int) (*BatchResult, error) {
	// 创建任务
	tm := GetTaskManager()
	task := tm.CreateTask("shuangseqiu")

	// 初始化任务状态
	tm.UpdateTask(task.ID, func(t *TaskInfo) {
		t.Status = TaskStatusRunning
		t.Total = count
	})

	// 异步执行
	go func() {
		result, err := s.FetchShuangseqiuHistory(count, task.ID)
		endTime := time.Now()

		tm.UpdateTask(task.ID, func(t *TaskInfo) {
			t.EndTime = &endTime
			if err != nil {
				t.Status = TaskStatusFailed
				t.Message = err.Error()
			} else {
				t.Status = TaskStatusCompleted
				t.Total = result.Total
				t.Success = result.Success
				t.Skipped = result.Skipped
				t.Failed = result.Failed
				t.Result = result
			}
		})
	}()

	return &BatchResult{
		TaskID: task.ID,
	}, nil
}

// FetchShuangseqiuHistory 批量获取双色球历史数据
// taskID 可选，如果提供则实时更新任务进度
func (s *LotteryService) FetchShuangseqiuHistory(count int, taskID ...string) (*BatchResult, error) {
	if count <= 0 {
		count = 100 // 默认获取100期
	}
	if count > 2000 {
		count = 2000 // 最多2000期
	}

	result := &BatchResult{
		Total: count,
	}

	log.Printf("开始批量获取双色球历史数据，共 %d 期", count)

	// 构建请求参数
	params := map[string]string{
		"name":       "ssq",
		"issueCount": fmt.Sprintf("%d", count),
	}

	// 使用 Resty 发送请求并自动解析 JSON
	var apiResp ShuangseqiuResponse
	resp, err := s.client.R().
		SetHeader("Referer", "https://www.cwl.gov.cn/").
		SetHeader("Origin", "https://www.cwl.gov.cn").
		SetQueryParams(params).
		SetResult(&apiResp).
		Get(s.shuangseqiuURL)

	if err != nil {
		log.Printf("双色球批量API请求失败: %v", err)
		return result, fmt.Errorf("请求API失败: %v", err)
	}

	if !resp.IsSuccess() {
		log.Printf("双色球批量API返回错误状态码: %d", resp.StatusCode())
		return result, fmt.Errorf("API请求失败: HTTP %d", resp.StatusCode())
	}

	if apiResp.State != 0 || len(apiResp.Result) == 0 {
		return result, fmt.Errorf("API返回错误: %s (state=%d)", apiResp.Message, apiResp.State)
	}

	log.Printf("成功获取 %d 期双色球数据，开始处理...", len(apiResp.Result))

	// 处理每条数据
	for i, item := range apiResp.Result {
		// 控制请求频率，每10条休息一下
		if i > 0 && i%10 == 0 {
			log.Printf("已处理 %d/%d 期，休息1秒...", i, len(apiResp.Result))
			time.Sleep(1 * time.Second)
		}

		// 解析红球
		redBalls := []int{}
		for _, s := range strings.Split(item.Red, ",") {
			num, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				log.Printf("期号 %s: 解析红球失败: %s", item.Code, s)
				result.Failed++
				continue
			}
			redBalls = append(redBalls, num)
		}

		if len(redBalls) != 6 {
			log.Printf("期号 %s: 红球数量错误，期望6个，实际%d个", item.Code, len(redBalls))
			result.Failed++
			continue
		}

		// 解析蓝球
		blueBall, err := strconv.Atoi(strings.TrimSpace(item.Blue))
		if err != nil {
			log.Printf("期号 %s: 解析蓝球失败: %v", item.Code, err)
			result.Failed++
			continue
		}

		// 解析日期，支持多种格式
		// 先清理日期字符串，去掉括号中的内容（如：2024-08-20(二) -> 2024-08-20）
		dateStr := item.Date
		if idx := strings.Index(dateStr, "("); idx != -1 {
			dateStr = strings.TrimSpace(dateStr[:idx])
		}

		var drawDate time.Time
		dateFormats := []string{
			"2006-01-02",          // 2024-01-01
			"2006-01-02 15:04:05", // 2024-01-01 20:30:00
			"20060102",            // 20240101
		}

		parsed := false
		for _, format := range dateFormats {
			if d, err := time.Parse(format, dateStr); err == nil {
				drawDate = d
				parsed = true
				break
			}
		}

		if !parsed {
			log.Printf("期号 %s: 解析日期失败: 原始='%s', 清理后='%s'", item.Code, item.Date, dateStr)
			// 使用当前日期（不含时间）
			drawDate = time.Now().Truncate(24 * time.Hour)
		}

		ssq := &models.Shuangseqiu{
			Issue:    item.Code,
			RedBall1: redBalls[0],
			RedBall2: redBalls[1],
			RedBall3: redBalls[2],
			RedBall4: redBalls[3],
			RedBall5: redBalls[4],
			RedBall6: redBalls[5],
			BlueBall: blueBall,
			DrawDate: drawDate,
		}

		// 保存数据
		if err := s.SaveShuangseqiu(ssq); err != nil {
			if strings.Contains(err.Error(), "已存在") {
				result.Skipped++
			} else {
				log.Printf("期号 %s: 保存失败: %v", item.Code, err)
				result.Failed++
			}
		} else {
			result.Success++
		}

		// 如果有任务ID，实时更新任务进度
		if len(taskID) > 0 && taskID[0] != "" {
			tm := GetTaskManager()
			tm.UpdateTask(taskID[0], func(t *TaskInfo) {
				t.Success = result.Success
				t.Skipped = result.Skipped
				t.Failed = result.Failed
				// 计算进度
				processed := result.Success + result.Skipped + result.Failed
				if t.Total > 0 {
					t.Progress = int(float64(processed) / float64(t.Total) * 100)
				}
			})
		}
	}

	log.Printf("双色球批量获取完成: 总计=%d, 成功=%d, 跳过=%d, 失败=%d",
		result.Total, result.Success, result.Skipped, result.Failed)

	return result, nil
}

// FetchDaletouHistoryAsync 异步批量获取大乐透历史数据
func (s *LotteryService) FetchDaletouHistoryAsync(count int) (*BatchResult, error) {
	// 创建任务
	tm := GetTaskManager()
	task := tm.CreateTask("daletou")

	// 初始化任务状态
	tm.UpdateTask(task.ID, func(t *TaskInfo) {
		t.Status = TaskStatusRunning
		t.Total = count
	})

	// 异步执行
	go func() {
		result, err := s.FetchDaletouHistory(count, task.ID)
		endTime := time.Now()

		tm.UpdateTask(task.ID, func(t *TaskInfo) {
			t.EndTime = &endTime
			if err != nil {
				t.Status = TaskStatusFailed
				t.Message = err.Error()
			} else {
				t.Status = TaskStatusCompleted
				t.Total = result.Total
				t.Success = result.Success
				t.Skipped = result.Skipped
				t.Failed = result.Failed
				t.Result = result
			}
		})
	}()

	return &BatchResult{
		TaskID: task.ID,
	}, nil
}

// FetchDaletouHistory 批量获取大乐透历史数据
// taskID 可选，如果提供则实时更新任务进度
func (s *LotteryService) FetchDaletouHistory(count int, taskID ...string) (*BatchResult, error) {
	if count <= 0 {
		count = 100 // 默认获取100期
	}
	if count > 2000 {
		count = 2000 // 最多2000期
	}

	result := &BatchResult{
		Total: count,
	}

	log.Printf("开始批量获取大乐透历史数据，共 %d 期", count)

	// 大乐透API使用页码，每页30条
	pageSize := 30
	totalPages := (count + pageSize - 1) / pageSize

	for page := 1; page <= totalPages; page++ {
		log.Printf("正在获取第 %d/%d 页...", page, totalPages)

		// 构建请求参数
		params := map[string]string{
			"gameNo":     "85",
			"provinceId": "0",
			"isVerify":   "1",
			"pageNo":     fmt.Sprintf("%d", page),
			"pageSize":   fmt.Sprintf("%d", pageSize),
		}

		// 使用 Resty 发送请求
		var apiResp DaletouResponse
		resp, err := s.client.R().
			SetHeader("Referer", "https://www.sporttery.cn/").
			SetHeader("Origin", "https://www.sporttery.cn").
			SetQueryParams(params).
			SetResult(&apiResp).
			Get(s.daletouURL)

		if err != nil {
			log.Printf("大乐透批量API请求失败(第%d页): %v", page, err)
			continue
		}

		if !resp.IsSuccess() {
			log.Printf("大乐透批量API返回错误状态码(第%d页): %d, 响应内容: %s", page, resp.StatusCode(), string(resp.Body()))
			continue
		}

		log.Printf("大乐透批量API响应(第%d页): Success=%v, ErrorCode=%s, ListCount=%d",
			page, apiResp.Success, apiResp.ErrorCode, len(apiResp.Value.List))

		if !apiResp.Success || len(apiResp.Value.List) == 0 {
			log.Printf("大乐透批量API第%d页无数据 (Success=%v, ErrorCode=%s, Message=%s)",
				page, apiResp.Success, apiResp.ErrorCode, apiResp.Message)
			break
		}

		log.Printf("第 %d 页获取到 %d 期数据", page, len(apiResp.Value.List))

		// 处理每条数据
		for i, item := range apiResp.Value.List {
			// 控制请求频率
			if i > 0 && i%10 == 0 {
				time.Sleep(500 * time.Millisecond)
			}

			// 解析开奖结果
			var frontBalls, backBalls []int

			if strings.Contains(item.LotteryDrawResult, "#") {
				// 格式1: 带 # 分隔符
				parts := strings.Split(item.LotteryDrawResult, "#")
				if len(parts) != 2 {
					log.Printf("期号 %s: 开奖结果格式错误", item.LotteryDrawNum)
					result.Failed++
					continue
				}

				for _, s := range strings.Fields(strings.TrimSpace(parts[0])) {
					num, _ := strconv.Atoi(s)
					frontBalls = append(frontBalls, num)
				}

				for _, s := range strings.Fields(strings.TrimSpace(parts[1])) {
					num, _ := strconv.Atoi(s)
					backBalls = append(backBalls, num)
				}
			} else {
				// 格式2: 无 # 分隔符
				allBalls := []int{}
				for _, s := range strings.Fields(strings.TrimSpace(item.LotteryDrawResult)) {
					num, _ := strconv.Atoi(s)
					allBalls = append(allBalls, num)
				}

				if len(allBalls) != 7 {
					log.Printf("期号 %s: 球号数量错误", item.LotteryDrawNum)
					result.Failed++
					continue
				}

				frontBalls = allBalls[:5]
				backBalls = allBalls[5:]
			}

			if len(frontBalls) != 5 || len(backBalls) != 2 {
				log.Printf("期号 %s: 球号数量错误", item.LotteryDrawNum)
				result.Failed++
				continue
			}

			// 解析日期
			drawDate, err := time.Parse("2006-01-02 15:04:05", item.LotteryDrawTime)
			if err != nil {
				drawDate, err = time.Parse("2006-01-02", item.LotteryDrawTime)
				if err != nil {
					drawDate = time.Now()
				}
			}

			dlt := &models.Daletou{
				Issue:      item.LotteryDrawNum,
				FrontBall1: frontBalls[0],
				FrontBall2: frontBalls[1],
				FrontBall3: frontBalls[2],
				FrontBall4: frontBalls[3],
				FrontBall5: frontBalls[4],
				BackBall1:  backBalls[0],
				BackBall2:  backBalls[1],
				DrawDate:   drawDate,
			}

			// 保存数据
			if err := s.SaveDaletou(dlt); err != nil {
				if strings.Contains(err.Error(), "已存在") {
					result.Skipped++
				} else {
					log.Printf("期号 %s: 保存失败: %v", item.LotteryDrawNum, err)
					result.Failed++
				}
			} else {
				result.Success++
			}

			// 如果有任务ID，实时更新任务进度
			if len(taskID) > 0 && taskID[0] != "" {
				tm := GetTaskManager()
				tm.UpdateTask(taskID[0], func(t *TaskInfo) {
					t.Success = result.Success
					t.Skipped = result.Skipped
					t.Failed = result.Failed
					// 计算进度
					processed := result.Success + result.Skipped + result.Failed
					if t.Total > 0 {
						t.Progress = int(float64(processed) / float64(t.Total) * 100)
					}
				})
			}
		}

		// 页面之间休息一下
		if page < totalPages {
			log.Printf("休息2秒后继续...")
			time.Sleep(2 * time.Second)
		}
	}

	log.Printf("大乐透批量获取完成: 总计=%d, 成功=%d, 跳过=%d, 失败=%d",
		result.Total, result.Success, result.Skipped, result.Failed)

	return result, nil
}

// GetShuangseqiuList 获取双色球列表
func (s *LotteryService) GetShuangseqiuList(page, pageSize int, issue string) ([]models.Shuangseqiu, int64, error) {
	db := database.GetDB()
	var list []models.Shuangseqiu
	var total int64

	query := db.Model(&models.Shuangseqiu{})
	if issue != "" {
		query = query.Where("issue LIKE ?", "%"+issue+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	// 按照期号倒序排序（最新的在前面）
	err := query.Order("issue DESC").Offset(offset).Limit(pageSize).Find(&list).Error

	return list, total, err
}

// GetDaletouList 获取大乐透列表
func (s *LotteryService) GetDaletouList(page, pageSize int, issue string) ([]models.Daletou, int64, error) {
	db := database.GetDB()
	var list []models.Daletou
	var total int64

	query := db.Model(&models.Daletou{})
	if issue != "" {
		query = query.Where("issue LIKE ?", "%"+issue+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	// 按照期号倒序排序（最新的在前面）
	err := query.Order("issue DESC").Offset(offset).Limit(pageSize).Find(&list).Error

	return list, total, err
}

// GetShuangseqiuStatistics 获取双色球统计数据
func (s *LotteryService) GetShuangseqiuStatistics(ballType string) ([]models.ShuangseqiuStatistics, error) {
	db := database.GetDB()
	var statistics []models.ShuangseqiuStatistics

	if ballType == "red" {
		// 统计红球
		for i := 1; i <= 6; i++ {
			field := fmt.Sprintf("red_ball%d", i)
			for num := 1; num <= 33; num++ {
				var count int64
				db.Model(&models.Shuangseqiu{}).Where(field+" = ?", num).Count(&count)
				if count > 0 {
					// 查找是否已经存在
					found := false
					for j := range statistics {
						if statistics[j].Number == num {
							statistics[j].Count += count
							found = true
							break
						}
					}
					if !found {
						statistics = append(statistics, models.ShuangseqiuStatistics{
							Number: num,
							Count:  count,
						})
					}
				}
			}
		}
	} else {
		// 统计蓝球
		for num := 1; num <= 16; num++ {
			var count int64
			db.Model(&models.Shuangseqiu{}).Where("blue_ball = ?", num).Count(&count)
			if count > 0 {
				statistics = append(statistics, models.ShuangseqiuStatistics{
					Number: num,
					Count:  count,
				})
			}
		}
	}

	return statistics, nil
}

// GetDaletouStatistics 获取大乐透统计数据
func (s *LotteryService) GetDaletouStatistics(ballType string) ([]models.DaletouStatistics, error) {
	db := database.GetDB()
	var statistics []models.DaletouStatistics

	if ballType == "front" {
		// 统计前区
		for i := 1; i <= 5; i++ {
			field := fmt.Sprintf("front_ball%d", i)
			for num := 1; num <= 35; num++ {
				var count int64
				db.Model(&models.Daletou{}).Where(field+" = ?", num).Count(&count)
				if count > 0 {
					found := false
					for j := range statistics {
						if statistics[j].Number == num {
							statistics[j].Count += count
							found = true
							break
						}
					}
					if !found {
						statistics = append(statistics, models.DaletouStatistics{
							Number: num,
							Count:  count,
						})
					}
				}
			}
		}
	} else {
		// 统计后区
		for i := 1; i <= 2; i++ {
			field := fmt.Sprintf("back_ball%d", i)
			for num := 1; num <= 12; num++ {
				var count int64
				db.Model(&models.Daletou{}).Where(field+" = ?", num).Count(&count)
				if count > 0 {
					found := false
					for j := range statistics {
						if statistics[j].Number == num {
							statistics[j].Count += count
							found = true
							break
						}
					}
					if !found {
						statistics = append(statistics, models.DaletouStatistics{
							Number: num,
							Count:  count,
						})
					}
				}
			}
		}
	}

	return statistics, nil
}
