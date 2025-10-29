package models

import (
	"time"
)

// Daletou 大乐透模型
type Daletou struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Issue      string    `gorm:"uniqueIndex;not null" json:"issue"`   // 期号
	FrontBall1 int       `gorm:"not null" json:"front_ball_1"`        // 前区1
	FrontBall2 int       `gorm:"not null" json:"front_ball_2"`        // 前区2
	FrontBall3 int       `gorm:"not null" json:"front_ball_3"`        // 前区3
	FrontBall4 int       `gorm:"not null" json:"front_ball_4"`        // 前区4
	FrontBall5 int       `gorm:"not null" json:"front_ball_5"`        // 前区5
	BackBall1  int       `gorm:"not null" json:"back_ball_1"`         // 后区1
	BackBall2  int       `gorm:"not null" json:"back_ball_2"`         // 后区2
	DrawDate   time.Time `gorm:"not null" json:"draw_date"`           // 开奖日期
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Daletou) TableName() string {
	return "daletou"
}

// DaletouStatistics 大乐透统计
type DaletouStatistics struct {
	Number int   `json:"number"` // 号码
	Count  int64 `json:"count"`  // 出现次数
}

