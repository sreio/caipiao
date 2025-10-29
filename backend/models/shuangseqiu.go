package models

import (
	"time"
)

// Shuangseqiu 双色球模型
type Shuangseqiu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Issue     string    `gorm:"uniqueIndex;not null" json:"issue"`      // 期号
	RedBall1  int       `gorm:"not null" json:"red_ball_1"`             // 红球1
	RedBall2  int       `gorm:"not null" json:"red_ball_2"`             // 红球2
	RedBall3  int       `gorm:"not null" json:"red_ball_3"`             // 红球3
	RedBall4  int       `gorm:"not null" json:"red_ball_4"`             // 红球4
	RedBall5  int       `gorm:"not null" json:"red_ball_5"`             // 红球5
	RedBall6  int       `gorm:"not null" json:"red_ball_6"`             // 红球6
	BlueBall  int       `gorm:"not null" json:"blue_ball"`              // 蓝球
	DrawDate  time.Time `gorm:"not null" json:"draw_date"`              // 开奖日期
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Shuangseqiu) TableName() string {
	return "shuangseqiu"
}

// ShuangseqiuStatistics 双色球统计
type ShuangseqiuStatistics struct {
	Number int   `json:"number"` // 号码
	Count  int64 `json:"count"`  // 出现次数
}

