package models

import (
	"time"

	"gorm.io/gorm"
)

type Carousel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;comment:轮播图标题"`
	Image     string         `json:"image" gorm:"size:255;not null;comment:轮播图图片"`
	Link      string         `json:"link" gorm:"size:255;comment:跳转链接"`
	LinkType  int            `json:"link_type" gorm:"default:1;comment:链接类型 1内部链接 2外部链接"`
	Position  int            `json:"position" gorm:"default:1;comment:位置 1首页轮播"`
	Sort      int            `json:"sort" gorm:"default:0;comment:排序"`
	Status    int            `json:"status" gorm:"default:1;comment:状态 1显示 0隐藏"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
