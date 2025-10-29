package services

import (
	"sync"
	"time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 等待中
	TaskStatusRunning   TaskStatus = "running"   // 运行中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 失败
)

// TaskInfo 任务信息
type TaskInfo struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`         // shuangseqiu 或 daletou
	Status      TaskStatus  `json:"status"`
	Progress    int         `json:"progress"`     // 进度百分比
	Total       int         `json:"total"`        // 总数
	Success     int         `json:"success"`      // 成功
	Skipped     int         `json:"skipped"`      // 跳过
	Failed      int         `json:"failed"`       // 失败
	Message     string      `json:"message"`      // 消息
	StartTime   time.Time   `json:"start_time"`   // 开始时间
	EndTime     *time.Time  `json:"end_time"`     // 结束时间
	Result      interface{} `json:"result"`       // 结果
}

// TaskManager 任务管理器
type TaskManager struct {
	tasks map[string]*TaskInfo
	mutex sync.RWMutex
}

var (
	taskManager *TaskManager
	once        sync.Once
)

// GetTaskManager 获取任务管理器单例
func GetTaskManager() *TaskManager {
	once.Do(func() {
		taskManager = &TaskManager{
			tasks: make(map[string]*TaskInfo),
		}
	})
	return taskManager
}

// CreateTask 创建任务
func (tm *TaskManager) CreateTask(taskType string) *TaskInfo {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	task := &TaskInfo{
		ID:        generateTaskID(),
		Type:      taskType,
		Status:    TaskStatusPending,
		StartTime: time.Now(),
	}

	tm.tasks[task.ID] = task
	return task
}

// GetTask 获取任务
func (tm *TaskManager) GetTask(id string) *TaskInfo {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return tm.tasks[id]
}

// UpdateTask 更新任务
func (tm *TaskManager) UpdateTask(id string, updater func(*TaskInfo)) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if task, ok := tm.tasks[id]; ok {
		updater(task)
	}
}

// CleanOldTasks 清理旧任务（保留最近1小时的）
func (tm *TaskManager) CleanOldTasks() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	now := time.Now()
	for id, task := range tm.tasks {
		if task.EndTime != nil && now.Sub(*task.EndTime) > time.Hour {
			delete(tm.tasks, id)
		}
	}
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

