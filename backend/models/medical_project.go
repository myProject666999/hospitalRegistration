package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalProject struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:200;not null;comment:项目名称"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;comment:项目编码"`
	Category    string         `json:"category" gorm:"size:100;comment:项目分类"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);default:0;comment:价格"`
	Description string         `json:"description" gorm:"type:text;comment:项目描述"`
	Notice      string         `json:"notice" gorm:"type:text;comment:注意事项"`
	Duration    string         `json:"duration" gorm:"size:50;comment:检查时长"`
	Preparation string         `json:"preparation" gorm:"type:text;comment:准备事项"`
	Image       string         `json:"image" gorm:"size:255;comment:项目图片"`
	Sort        int            `json:"sort" gorm:"default:0;comment:排序"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
