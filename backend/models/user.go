package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password     string         `json:"-" gorm:"size:255;not null;comment:密码"`
	RealName     string         `json:"real_name" gorm:"size:50;comment:真实姓名"`
	Phone        string         `json:"phone" gorm:"size:20;comment:手机号"`
	Email        string         `json:"email" gorm:"size:100;comment:邮箱"`
	IDCard       string         `json:"id_card" gorm:"size:18;comment:身份证号"`
	Gender       int            `json:"gender" gorm:"default:1;comment:性别 1男 2女"`
	Birthday     *time.Time     `json:"birthday" gorm:"comment:生日"`
	Avatar       string         `json:"avatar" gorm:"size:255;comment:头像"`
	Status       int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	LastLoginAt  *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt    time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Phone    string `json:"phone" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
}

type UserUpdateRequest struct {
	RealName string     `json:"real_name"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	IDCard   string     `json:"id_card"`
	Gender   int        `json:"gender"`
	Birthday *time.Time `json:"birthday"`
	Avatar   string     `json:"avatar"`
}

type UserPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}
