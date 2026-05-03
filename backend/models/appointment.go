package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	AppointmentNo  string         `json:"appointment_no" gorm:"uniqueIndex;size:50;not null;comment:预约单号"`
	UserID         uint           `json:"user_id" gorm:"not null;index;comment:用户ID"`
	User           *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DoctorID       uint           `json:"doctor_id" gorm:"not null;index;comment:医生ID"`
	Doctor         *Doctor        `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	ScheduleID     uint           `json:"schedule_id" gorm:"not null;index;comment:排班ID"`
	Schedule       *Schedule      `json:"schedule,omitempty" gorm:"foreignKey:ScheduleID"`
	Date           time.Time      `json:"date" gorm:"type:date;not null;comment:预约日期"`
	Period         int            `json:"period" gorm:"not null;comment:时段 1上午 2下午 3晚上"`
	TimeSlot       string         `json:"time_slot" gorm:"size:50;comment:预约时段"`
	PatientName    string         `json:"patient_name" gorm:"size:50;not null;comment:就诊人姓名"`
	PatientPhone   string         `json:"patient_phone" gorm:"size:20;comment:就诊人电话"`
	PatientIDCard  string         `json:"patient_id_card" gorm:"size:18;comment:就诊人身份证"`
	PatientGender  int            `json:"patient_gender" gorm:"default:1;comment:就诊人性别 1男 2女"`
	PatientAge     int            `json:"patient_age" gorm:"comment:就诊人年龄"`
	ChiefComplaint string         `json:"chief_complaint" gorm:"size:500;comment:主诉"`
	ConsultationFee float64       `json:"consultation_fee" gorm:"type:decimal(10,2);default:0;comment:挂号费"`
	Status         int            `json:"status" gorm:"default:0;comment:状态 0待确认 1已确认 2已就诊 3已取消 4已过期"`
	CancelReason   string         `json:"cancel_reason" gorm:"size:500;comment:取消原因"`
	QueuedNumber   int            `json:"queued_number" gorm:"comment:排队号"`
	Location       string         `json:"location" gorm:"size:100;comment:就诊地点"`
	CreatedAt      time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type AppointmentCreateRequest struct {
	DoctorID       uint   `json:"doctor_id" binding:"required"`
	ScheduleID     uint   `json:"schedule_id" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Period         int    `json:"period" binding:"required"`
	TimeSlot       string `json:"time_slot"`
	PatientName    string `json:"patient_name" binding:"required"`
	PatientPhone   string `json:"patient_phone" binding:"required"`
	PatientIDCard  string `json:"patient_id_card"`
	PatientGender  int    `json:"patient_gender"`
	PatientAge     int    `json:"patient_age"`
	ChiefComplaint string `json:"chief_complaint"`
}
