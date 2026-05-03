package models

import (
	"time"

	"gorm.io/gorm"
)

type Hospital struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:200;not null;comment:医院名称"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;comment:医院编码"`
	Logo        string         `json:"logo" gorm:"size:255;comment:医院logo"`
	Address     string         `json:"address" gorm:"size:500;comment:医院地址"`
	Phone       string         `json:"phone" gorm:"size:50;comment:联系电话"`
	Email       string         `json:"email" gorm:"size:100;comment:邮箱"`
	Level       string         `json:"level" gorm:"size:50;comment:医院等级"`
	Type        string         `json:"type" gorm:"size:50;comment:医院类型"`
	Description string         `json:"description" gorm:"type:text;comment:医院简介"`
	Latitude    float64        `json:"latitude" gorm:"type:decimal(10,7);comment:纬度"`
	Longitude   float64        `json:"longitude" gorm:"type:decimal(10,7);comment:经度"`
	Sort        int            `json:"sort" gorm:"default:0;comment:排序"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
