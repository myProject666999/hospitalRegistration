package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:200;not null;comment:药品名称"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;comment:药品编码"`
	Category    string         `json:"category" gorm:"size:100;comment:药品分类"`
	Spec        string         `json:"spec" gorm:"size:100;comment:规格"`
	Unit        string         `json:"unit" gorm:"size:20;comment:单位"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);default:0;comment:价格"`
	Manufacturer string        `json:"manufacturer" gorm:"size:200;comment:生产厂家"`
	Usage       string         `json:"usage" gorm:"size:500;comment:用法用量"`
	Indication  string         `json:"indication" gorm:"type:text;comment:适应症"`
	Contraindication string    `json:"contraindication" gorm:"type:text;comment:禁忌症"`
	AdverseReaction string     `json:"adverse_reaction" gorm:"type:text;comment:不良反应"`
	Image       string         `json:"image" gorm:"size:255;comment:药品图片"`
	Description string         `json:"description" gorm:"type:text;comment:药品描述"`
	Stock       int            `json:"stock" gorm:"default:0;comment:库存"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
