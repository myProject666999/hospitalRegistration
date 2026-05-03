package models

import (
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_user_type;comment:用户ID"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type      int            `json:"type" gorm:"not null;uniqueIndex:idx_user_type;comment:收藏类型 1医生 2药品 3公告"`
	TargetID  uint           `json:"target_id" gorm:"not null;uniqueIndex:idx_user_type;comment:目标ID"`
	Doctor    *Doctor        `json:"doctor,omitempty" gorm:"-"`
	Medicine  *Medicine      `json:"medicine,omitempty" gorm:"-"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type FavoriteRequest struct {
	Type     int `json:"type" binding:"required,min=1,max=3"`
	TargetID uint `json:"target_id" binding:"required"`
}
