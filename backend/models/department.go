package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:科室名称"`
	Code        string         `json:"code" gorm:"size:50;uniqueIndex;comment:科室编码"`
	Description string         `json:"description" gorm:"type:text;comment:科室描述"`
	HospitalID  uint           `json:"hospital_id" gorm:"comment:医院ID"`
	Hospital    *Hospital      `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"`
	ParentID    *uint          `json:"parent_id" gorm:"comment:父级ID"`
	Parent      *Department    `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Sort        int            `json:"sort" gorm:"default:0;comment:排序"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
