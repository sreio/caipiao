package api

import (
	"caipiao/backend/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	lotteryService *services.LotteryService
	trendService   *services.TrendService
}

// NewHandler 创建处理器实例
func NewHandler(lotteryService *services.LotteryService, trendService *services.TrendService) *Handler {
	return &Handler{
		lotteryService: lotteryService,
		trendService:   trendService,
	}
}

// Response 统一响应格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"page_size"`
}

// GetShuangseqiuList 获取双色球列表
func (h *Handler) GetShuangseqiuList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	issue := c.Query("issue")

	list, total, err := h.lotteryService.GetShuangseqiuList(page, pageSize, issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: PageResponse{
			List:  list,
			Total: total,
			Page:  page,
			Size:  pageSize,
		},
	})
}

// GetDaletouList 获取大乐透列表
func (h *Handler) GetDaletouList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	issue := c.Query("issue")

	list, total, err := h.lotteryService.GetDaletouList(page, pageSize, issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: PageResponse{
			List:  list,
			Total: total,
			Page:  page,
			Size:  pageSize,
		},
	})
}

// FetchShuangseqiu 获取并保存双色球数据
func (h *Handler) FetchShuangseqiu(c *gin.Context) {
	issue := c.Query("issue")

	ssq, err := h.lotteryService.FetchShuangseqiu(issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "获取数据失败: " + err.Error(),
		})
		return
	}

	if err := h.lotteryService.SaveShuangseqiu(ssq); err != nil {
		// 检查是否是数据已存在的错误
		if strings.Contains(err.Error(), "数据已存在") {
			c.JSON(http.StatusOK, Response{
				Code: 1,
				Msg:  err.Error(),
				Data: ssq,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "保存数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "数据获取成功",
		Data: ssq,
	})
}

// FetchDaletou 获取并保存大乐透数据
func (h *Handler) FetchDaletou(c *gin.Context) {
	issue := c.Query("issue")

	dlt, err := h.lotteryService.FetchDaletou(issue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "获取数据失败: " + err.Error(),
		})
		return
	}

	if err := h.lotteryService.SaveDaletou(dlt); err != nil {
		// 检查是否是数据已存在的错误
		if strings.Contains(err.Error(), "数据已存在") {
			c.JSON(http.StatusOK, Response{
				Code: 1,
				Msg:  err.Error(),
				Data: dlt,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "保存数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "数据获取成功",
		Data: dlt,
	})
}

// GetShuangseqiuStatistics 获取双色球统计
func (h *Handler) GetShuangseqiuStatistics(c *gin.Context) {
	ballType := c.DefaultQuery("type", "red") // red or blue

	stats, err := h.lotteryService.GetShuangseqiuStatistics(ballType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: stats,
	})
}

// GetDaletouStatistics 获取大乐透统计
func (h *Handler) GetDaletouStatistics(c *gin.Context) {
	ballType := c.DefaultQuery("type", "front") // front or back

	stats, err := h.lotteryService.GetDaletouStatistics(ballType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: stats,
	})
}

// FetchShuangseqiuHistory 批量获取双色球历史数据
func (h *Handler) FetchShuangseqiuHistory(c *gin.Context) {
	countStr := c.DefaultQuery("count", "100")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 100
	}

	// 检查是否需要异步处理（数量大于100时）
	asyncStr := c.DefaultQuery("async", "false")
	async := asyncStr == "true" || count > 100

	if async {
		// 异步处理
		result, err := h.lotteryService.FetchShuangseqiuHistoryAsync(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code: -1,
				Msg:  "创建任务失败: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "任务已创建，请通过任务ID查询进度",
			Data: result,
		})
	} else {
		// 同步处理（不传taskID）
		result, err := h.lotteryService.FetchShuangseqiuHistory(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code: -1,
				Msg:  "批量获取失败: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "批量获取完成",
			Data: result,
		})
	}
}

// FetchDaletouHistory 批量获取大乐透历史数据
func (h *Handler) FetchDaletouHistory(c *gin.Context) {
	countStr := c.DefaultQuery("count", "100")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 100
	}

	// 检查是否需要异步处理（数量大于100时）
	asyncStr := c.DefaultQuery("async", "false")
	async := asyncStr == "true" || count > 100

	if async {
		// 异步处理
		result, err := h.lotteryService.FetchDaletouHistoryAsync(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code: -1,
				Msg:  "创建任务失败: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "任务已创建，请通过任务ID查询进度",
			Data: result,
		})
	} else {
		// 同步处理（不传taskID）
		result, err := h.lotteryService.FetchDaletouHistory(count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code: -1,
				Msg:  "批量获取失败: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			Code: 0,
			Msg:  "批量获取完成",
			Data: result,
		})
	}
}

// GetTask 获取任务状态
func (h *Handler) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	tm := services.GetTaskManager()
	task := tm.GetTask(taskID)

	if task == nil {
		c.JSON(http.StatusNotFound, Response{
			Code: -1,
			Msg:  "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: task,
	})
}

// GetShuangseqiuTrend 获取双色球走势数据
func (h *Handler) GetShuangseqiuTrend(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	trend, err := h.trendService.GetShuangseqiuTrend(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "获取走势数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: trend,
	})
}

// GetDaletouTrend 获取大乐透走势数据
func (h *Handler) GetDaletouTrend(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	trend, err := h.trendService.GetDaletouTrend(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code: -1,
			Msg:  "获取走势数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: trend,
	})
}
