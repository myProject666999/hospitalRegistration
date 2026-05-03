package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password      string         `json:"-" gorm:"size:255;not null;comment:密码"`
	RealName      string         `json:"real_name" gorm:"size:50;not null;comment:真实姓名"`
	Phone         string         `json:"phone" gorm:"size:20;comment:手机号"`
	Email         string         `json:"email" gorm:"size:100;comment:邮箱"`
	Avatar        string         `json:"avatar" gorm:"size:255;comment:头像"`
	DepartmentID  uint           `json:"department_id" gorm:"comment:科室ID"`
	Department    *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	HospitalID    uint           `json:"hospital_id" gorm:"comment:医院ID"`
	Hospital      *Hospital      `json:"hospital,omitempty" gorm:"foreignKey:HospitalID"`
	Title         string         `json:"title" gorm:"size:50;comment:职称"`
	Specialty     string         `json:"specialty" gorm:"size:500;comment:专长"`
	Introduction  string         `json:"introduction" gorm:"type:text;comment:简介"`
	ConsultationFee float64      `json:"consultation_fee" gorm:"type:decimal(10,2);default:0;comment:挂号费"`
	Rating        float64        `json:"rating" gorm:"type:decimal(3,1);default:5.0;comment:评分"`
	RatingCount   int            `json:"rating_count" gorm:"default:0;comment:评分次数"`
	Status        int            `json:"status" gorm:"default:1;comment:状态 1正常 0禁用"`
	LastLoginAt   *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt     time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type DoctorLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DoctorUpdateRequest struct {
	RealName        string  `json:"real_name"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	Avatar          string  `json:"avatar"`
	DepartmentID    uint    `json:"department_id"`
	HospitalID      uint    `json:"hospital_id"`
	Title           string  `json:"title"`
	Specialty       string  `json:"specialty"`
	Introduction    string  `json:"introduction"`
	ConsultationFee float64 `json:"consultation_fee"`
	Status          int     `json:"status"`
}
