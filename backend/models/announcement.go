package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null;comment:公告标题"`
	Content   string         `json:"content" gorm:"type:text;not null;comment:公告内容"`
	Summary   string         `json:"summary" gorm:"size:500;comment:公告摘要"`
	Cover     string         `json:"cover" gorm:"size:255;comment:封面图片"`
	Type      int            `json:"type" gorm:"default:1;comment:类型 1公告 2新闻 3通知"`
	IsTop     int            `json:"is_top" gorm:"default:0;comment:是否置顶 0否 1是"`
	Status    int            `json:"status" gorm:"default:1;comment:状态 1发布 0草稿"`
	ViewCount int            `json:"view_count" gorm:"default:0;comment:浏览量"`
	Sort      int            `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
