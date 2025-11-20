package main

import (
	"context"
	"fmt"
	"log"

	"caipiao/config"
	"caipiao/database"
	"caipiao/services"
)

// App 应用结构
type App struct {
	ctx              context.Context
	lotteryService   *services.LotteryService
	trendService     *services.TrendService
	recommendService *services.RecommendService
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Startup 应用启动时调用
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("彩票助手启动中...")

	// 加载配置
	cfg := config.GetConfig()

	// 初始化数据库（使用用户目录）
	dbPath := getDatabasePath()
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建服务
	a.lotteryService = services.NewLotteryService(cfg.API.ShuangseqiuURL, cfg.API.DaletouURL)
	a.trendService = services.NewTrendService()
	a.recommendService = services.NewRecommendService()

	log.Println("彩票助手启动完成！")
}

// Shutdown 应用关闭时调用
func (a *App) Shutdown(ctx context.Context) {
	log.Println("彩票助手正在关闭...")
}

// GetShuangseqiuList 获取双色球列表
func (a *App) GetShuangseqiuList(page, pageSize int, issue string) (interface{}, error) {
	return a.lotteryService.GetShuangseqiuList(page, pageSize, issue)
}

// FetchShuangseqiu 获取双色球数据
func (a *App) FetchShuangseqiu(issue string) (interface{}, error) {
	ssq, err := a.lotteryService.FetchShuangseqiu(issue)
	if err != nil {
		return nil, err
	}
	return ssq, a.lotteryService.SaveShuangseqiu(ssq)
}

// GetShuangseqiuStatistics 获取双色球统计
func (a *App) GetShuangseqiuStatistics(ballType string) (interface{}, error) {
	return a.lotteryService.GetShuangseqiuStatistics(ballType)
}

// GetShuangseqiuTrend 获取双色球走势
func (a *App) GetShuangseqiuTrend(limit int) (interface{}, error) {
	return a.trendService.GetShuangseqiuTrend(limit)
}

// GetShuangseqiuRecommendation 获取双色球推荐
func (a *App) GetShuangseqiuRecommendation(count int) (interface{}, error) {
	return a.recommendService.GenerateShuangseqiuRecommendation(count)
}

// GetDaletouList 获取大乐透列表
func (a *App) GetDaletouList(page, pageSize int, issue string) (interface{}, error) {
	return a.lotteryService.GetDaletouList(page, pageSize, issue)
}

// FetchDaletou 获取大乐透数据
func (a *App) FetchDaletou(issue string) (interface{}, error) {
	dlt, err := a.lotteryService.FetchDaletou(issue)
	if err != nil {
		return nil, err
	}
	return dlt, a.lotteryService.SaveDaletou(dlt)
}

// GetDaletouStatistics 获取大乐透统计
func (a *App) GetDaletouStatistics(ballType string) (interface{}, error) {
	return a.lotteryService.GetDaletouStatistics(ballType)
}

// GetDaletouTrend 获取大乐透走势
func (a *App) GetDaletouTrend(limit int) (interface{}, error) {
	return a.trendService.GetDaletouTrend(limit)
}

// GetDaletouRecommendation 获取大乐透推荐
func (a *App) GetDaletouRecommendation(count int) (interface{}, error) {
	return a.recommendService.GenerateDaletouRecommendation(count)
}
