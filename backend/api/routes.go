package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, handler *Handler) {
	// API路由组
	api := r.Group("/api")
	{
		// 双色球相关路由
		ssq := api.Group("/shuangseqiu")
		{
			ssq.GET("/list", handler.GetShuangseqiuList)             // 获取列表
			ssq.POST("/fetch", handler.FetchShuangseqiu)             // 获取最新数据
			ssq.POST("/fetch-history", handler.FetchShuangseqiuHistory) // 批量获取历史数据
			ssq.GET("/statistics", handler.GetShuangseqiuStatistics) // 统计数据
			ssq.GET("/trend", handler.GetShuangseqiuTrend)          // 走势图数据
			ssq.GET("/recommend", handler.GetShuangseqiuRecommendation) // 智能推荐
		}

		// 大乐透相关路由
		dlt := api.Group("/daletou")
		{
			dlt.GET("/list", handler.GetDaletouList)             // 获取列表
			dlt.POST("/fetch", handler.FetchDaletou)             // 获取最新数据
			dlt.POST("/fetch-history", handler.FetchDaletouHistory) // 批量获取历史数据
			dlt.GET("/statistics", handler.GetDaletouStatistics) // 统计数据
			dlt.GET("/trend", handler.GetDaletouTrend)          // 走势图数据
			dlt.GET("/recommend", handler.GetDaletouRecommendation) // 智能推荐
		}

		// 任务相关路由
		api.GET("/task/:id", handler.GetTask) // 获取任务状态
	}
}
