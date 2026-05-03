package controllers

import (
	"hospitalRegistration/database"
	"hospitalRegistration/models"
	"hospitalRegistration/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DoctorLogin(c *gin.Context) {
	var req models.DoctorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var doctor models.Doctor
	if result := database.DB.Where("username = ?", req.Username).First(&doctor); result.Error != nil {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	if doctor.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, doctor.Password) {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	now := time.Now()
	doctor.LastLoginAt = &now
	database.DB.Save(&doctor)

	token, err := utils.GenerateToken(doctor.ID, doctor.Username, "doctor")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"doctor": gin.H{
			"id":         doctor.ID,
			"username":   doctor.Username,
			"real_name":  doctor.RealName,
			"phone":      doctor.Phone,
			"email":      doctor.Email,
			"avatar":     doctor.Avatar,
			"title":      doctor.Title,
			"department": doctor.Department,
		},
	})
}

func GetDoctorInfo(c *gin.Context) {
	doctorID := c.GetUint("user_id")

	var doctor models.Doctor
	if result := database.DB.Preload("Department").Preload("Hospital").First(&doctor, doctorID); result.Error != nil {
		utils.NotFound(c, "医生不存在")
		return
	}

	utils.Success(c, doctor)
}

func GetDoctorAppointments(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	date := c.Query("date")

	offset := (page - 1) * pageSize

	var appointments []models.Appointment
	var total int64

	query := database.DB.Model(&models.Appointment{}).Where("doctor_id = ?", doctorID)
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Count(&total)
	query.Preload("User").Preload("Schedule").
		Order("date ASC, period ASC, created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&appointments)

	utils.Page(c, total, page, pageSize, appointments)
}

func ConfirmAppointment(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	appointmentID, _ := strconv.Atoi(c.Param("id"))

	var appointment models.Appointment
	if result := database.DB.Where("id = ? AND doctor_id = ?", appointmentID, doctorID).First(&appointment); result.Error != nil {
		utils.NotFound(c, "预约不存在")
		return
	}

	if appointment.Status != 0 {
		utils.Error(c, 400, "该预约状态不可确认")
		return
	}

	appointment.Status = 1
	if result := database.DB.Save(&appointment); result.Error != nil {
		utils.InternalServerError(c, "确认失败")
		return
	}

	utils.SuccessWithMessage(c, "确认成功", nil)
}

func StartConsultation(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	appointmentID, _ := strconv.Atoi(c.Param("id"))

	var appointment models.Appointment
	if result := database.DB.Where("id = ? AND doctor_id = ?", appointmentID, doctorID).First(&appointment); result.Error != nil {
		utils.NotFound(c, "预约不存在")
		return
	}

	if appointment.Status != 1 {
		utils.Error(c, 400, "该预约状态不可开始就诊")
		return
	}

	var existingConsultation models.Consultation
	if result := database.DB.Where("appointment_id = ?", appointmentID).First(&existingConsultation); result.Error == nil {
		utils.Success(c, existingConsultation)
		return
	}

	consultationNo := "CS" + time.Now().Format("20060102150405") + strconv.Itoa(1000+int(appointmentID))
	now := time.Now()

	consultation := models.Consultation{
		ConsultationNo: consultationNo,
		AppointmentID:  appointment.ID,
		UserID:         appointment.UserID,
		DoctorID:       appointment.DoctorID,
		ChiefComplaint: appointment.ChiefComplaint,
		Status:         0,
		ConsultationTime: &now,
	}

	if result := database.DB.Create(&consultation); result.Error != nil {
		utils.InternalServerError(c, "开始就诊失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "开始就诊成功", consultation)
}

func UpdateConsultation(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	consultationID, _ := strconv.Atoi(c.Param("id"))

	var req models.ConsultationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var consultation models.Consultation
	if result := database.DB.Where("id = ? AND doctor_id = ?", consultationID, doctorID).First(&consultation); result.Error != nil {
		utils.NotFound(c, "就诊记录不存在")
		return
	}

	if consultation.Status == 2 {
		utils.Error(c, 400, "该就诊已完成，不可修改")
		return
	}

	updates := make(map[string]interface{})
	if req.ChiefComplaint != "" {
		updates["chief_complaint"] = req.ChiefComplaint
	}
	if req.PresentIllness != "" {
		updates["present_illness"] = req.PresentIllness
	}
	if req.PastHistory != "" {
		updates["past_history"] = req.PastHistory
	}
	if req.PhysicalExam != "" {
		updates["physical_exam"] = req.PhysicalExam
	}
	if req.AuxiliaryExam != "" {
		updates["auxiliary_exam"] = req.AuxiliaryExam
	}
	if req.Diagnosis != "" {
		updates["diagnosis"] = req.Diagnosis
	}
	if req.TreatmentPlan != "" {
		updates["treatment_plan"] = req.TreatmentPlan
	}
	if req.Prescription != "" {
		updates["prescription"] = req.Prescription
	}
	if req.Advice != "" {
		updates["advice"] = req.Advice
	}

	if result := database.DB.Model(&consultation).Updates(updates); result.Error != nil {
		utils.InternalServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func CompleteConsultation(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	consultationID, _ := strconv.Atoi(c.Param("id"))

	var consultation models.Consultation
	if result := database.DB.Where("id = ? AND doctor_id = ?", consultationID, doctorID).First(&consultation); result.Error != nil {
		utils.NotFound(c, "就诊记录不存在")
		return
	}

	if consultation.Status != 0 {
		utils.Error(c, 400, "该就诊状态不可完成")
		return
	}

	now := time.Now()
	consultation.Status = 1
	consultation.CompleteTime = &now

	if result := database.DB.Save(&consultation); result.Error != nil {
		utils.InternalServerError(c, "完成失败")
		return
	}

	var appointment models.Appointment
	database.DB.First(&appointment, consultation.AppointmentID)
	appointment.Status = 2
	database.DB.Save(&appointment)

	utils.SuccessWithMessage(c, "就诊完成", nil)
}

func GetDoctorOwnSchedules(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	date := c.Query("date")

	offset := (page - 1) * pageSize

	var schedules []models.Schedule
	var total int64

	query := database.DB.Model(&models.Schedule{}).Where("doctor_id = ?", doctorID)
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Count(&total)
	query.Order("date DESC, period ASC").
		Offset(offset).Limit(pageSize).
		Find(&schedules)

	utils.Page(c, total, page, pageSize, schedules)
}

func CreateSchedule(c *gin.Context) {
	doctorID := c.GetUint("user_id")

	var req models.ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	scheduleDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}

	var existingSchedule models.Schedule
	if result := database.DB.Where("doctor_id = ? AND date = ? AND period = ?", doctorID, scheduleDate, req.Period).First(&existingSchedule); result.Error == nil {
		utils.Error(c, 400, "该时段已存在排班")
		return
	}

	schedule := models.Schedule{
		DoctorID:    doctorID,
		Date:        scheduleDate,
		Period:      req.Period,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		TotalCount:  req.TotalCount,
		RemainCount: req.TotalCount,
		Location:    req.Location,
		Status:      1,
	}

	if result := database.DB.Create(&schedule); result.Error != nil {
		utils.InternalServerError(c, "创建排班失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", schedule)
}

func UpdateSchedule(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	scheduleID, _ := strconv.Atoi(c.Param("id"))

	var schedule models.Schedule
	if result := database.DB.Where("id = ? AND doctor_id = ?", scheduleID, doctorID).First(&schedule); result.Error != nil {
		utils.NotFound(c, "排班不存在")
		return
	}

	var req models.ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.TotalCount > 0 {
		usedCount := schedule.TotalCount - schedule.RemainCount
		if req.TotalCount < usedCount {
			utils.Error(c, 400, "总号源数不能小于已使用号源数")
			return
		}
		schedule.TotalCount = req.TotalCount
		schedule.RemainCount = req.TotalCount - usedCount
	}

	if req.StartTime != "" {
		schedule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		schedule.EndTime = req.EndTime
	}
	if req.Location != "" {
		schedule.Location = req.Location
	}

	if result := database.DB.Save(&schedule); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", schedule)
}

func CancelSchedule(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	scheduleID, _ := strconv.Atoi(c.Param("id"))

	var schedule models.Schedule
	if result := database.DB.Where("id = ? AND doctor_id = ?", scheduleID, doctorID).First(&schedule); result.Error != nil {
		utils.NotFound(c, "排班不存在")
		return
	}

	if schedule.Status == -1 {
		utils.Error(c, 400, "该排班已取消")
		return
	}

	usedCount := schedule.TotalCount - schedule.RemainCount
	if usedCount > 0 {
		utils.Error(c, 400, "该排班已有预约，不可取消")
		return
	}

	schedule.Status = -1
	if result := database.DB.Save(&schedule); result.Error != nil {
		utils.InternalServerError(c, "取消失败")
		return
	}

	utils.SuccessWithMessage(c, "取消成功", nil)
}

func GetDoctorConsultations(c *gin.Context) {
	doctorID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var consultations []models.Consultation
	var total int64

	query := database.DB.Model(&models.Consultation{}).Where("doctor_id = ?", doctorID)
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)
	query.Preload("User").Preload("Appointment").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&consultations)

	utils.Page(c, total, page, pageSize, consultations)
}
