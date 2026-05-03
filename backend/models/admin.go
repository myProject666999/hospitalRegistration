package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password    string         `json:"-" gorm:"size:255;not null;comment:密码"`
	RealName    string         `json:"real_name" gorm:"size:50;comment:真实姓名"`
	Phone       string         `json:"phone" gorm:"size:20;comment:手机号"`
	Email       string         `json:"email" gorm:"size:100;comment:邮箱"`
	Role        int            `json:"role" gorm:"default:1;comment:角色 1超级管理员 2普通管理员"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	LastLoginAt *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
