package models

import (
	"time"

	"gorm.io/gorm"
)

type Consultation struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ConsultationNo  string         `json:"consultation_no" gorm:"uniqueIndex;size:50;not null;comment:就诊单号"`
	AppointmentID   uint           `json:"appointment_id" gorm:"not null;uniqueIndex;comment:预约ID"`
	Appointment     *Appointment   `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID"`
	UserID          uint           `json:"user_id" gorm:"not null;index;comment:用户ID"`
	User            *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DoctorID        uint           `json:"doctor_id" gorm:"not null;index;comment:医生ID"`
	Doctor          *Doctor        `json:"doctor,omitempty" gorm:"foreignKey:DoctorID"`
	ChiefComplaint  string         `json:"chief_complaint" gorm:"size:500;comment:主诉"`
	PresentIllness  string         `json:"present_illness" gorm:"type:text;comment:现病史"`
	PastHistory     string         `json:"past_history" gorm:"type:text;comment:既往史"`
	PhysicalExam    string         `json:"physical_exam" gorm:"type:text;comment:体格检查"`
	AuxiliaryExam   string         `json:"auxiliary_exam" gorm:"type:text;comment:辅助检查"`
	Diagnosis       string         `json:"diagnosis" gorm:"type:text;comment:诊断"`
	TreatmentPlan   string         `json:"treatment_plan" gorm:"type:text;comment:治疗方案"`
	Prescription    string         `json:"prescription" gorm:"type:text;comment:处方"`
	Advice          string         `json:"advice" gorm:"type:text;comment:医嘱"`
	Rating          int            `json:"rating" gorm:"comment:评分 1-5"`
	Comment         string         `json:"comment" gorm:"type:text;comment:评价内容"`
	CommentTime     *time.Time     `json:"comment_time" gorm:"comment:评价时间"`
	Status          int            `json:"status" gorm:"default:0;comment:状态 0就诊中 1已完成 2已评价"`
	ConsultationTime *time.Time    `json:"consultation_time" gorm:"comment:就诊时间"`
	CompleteTime    *time.Time     `json:"complete_time" gorm:"comment:完成时间"`
	CreatedAt       time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type ConsultationUpdateRequest struct {
	ChiefComplaint string `json:"chief_complaint"`
	PresentIllness string `json:"present_illness"`
	PastHistory    string `json:"past_history"`
	PhysicalExam   string `json:"physical_exam"`
	AuxiliaryExam  string `json:"auxiliary_exam"`
	Diagnosis      string `json:"diagnosis"`
	TreatmentPlan  string `json:"treatment_plan"`
	Prescription   string `json:"prescription"`
	Advice         string `json:"advice"`
}

type ConsultationCommentRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}
