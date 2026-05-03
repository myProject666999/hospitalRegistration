package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	DoctorID      uint           `json:"doctor_id" gorm:"not null;index;comment:医生ID"`
	Doctor        *Doctor        `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	Date          time.Time      `json:"date" gorm:"type:date;not null;comment:日期"`
	Period        int            `json:"period" gorm:"not null;comment:时段 1上午 2下午 3晚上"`
	StartTime     string         `json:"start_time" gorm:"size:10;comment:开始时间"`
	EndTime       string         `json:"end_time" gorm:"size:10;comment:结束时间"`
	TotalCount    int            `json:"total_count" gorm:"default:0;comment:总号源数"`
	RemainCount   int            `json:"remain_count" gorm:"default:0;comment:剩余号源数"`
	Location      string         `json:"location" gorm:"size:100;comment:出诊地点"`
	Status        int            `json:"status" gorm:"default:1;comment:状态 1可预约 0已约满 -1已取消"`
	CreatedAt     time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type ScheduleCreateRequest struct {
	DoctorID    uint   `json:"doctor_id" binding:"required"`
	Date        string `json:"date" binding:"required"`
	Period      int    `json:"period" binding:"required"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	TotalCount  int    `json:"total_count" binding:"required"`
	Location    string `json:"location"`
}
