package controllers

import (
	"hospitalRegistration/database"
	"hospitalRegistration/models"
	"hospitalRegistration/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var admin models.Admin
	if result := database.DB.Where("username = ?", req.Username).First(&admin); result.Error != nil {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	if admin.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	now := time.Now()
	admin.LastLoginAt = &now
	database.DB.Save(&admin)

	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"admin": gin.H{
			"id":        admin.ID,
			"username":  admin.Username,
			"real_name": admin.RealName,
			"role":      admin.Role,
		},
	})
}

func GetStatistics(c *gin.Context) {
	var totalUsers int64
	var totalDoctors int64
	var totalAppointments int64
	var todayAppointments int64
	var completedAppointments int64
	var totalConsultations int64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.Doctor{}).Count(&totalDoctors)
	database.DB.Model(&models.Appointment{}).Count(&totalAppointments)
	database.DB.Model(&models.Appointment{}).Where("DATE(date) = ?", time.Now().Format("2006-01-02")).Count(&todayAppointments)
	database.DB.Model(&models.Appointment{}).Where("status = 2").Count(&completedAppointments)
	database.DB.Model(&models.Consultation{}).Count(&totalConsultations)

	type RatingStats struct {
		AvgRating   float64 `json:"avg_rating"`
		RatingCount int64   `json:"rating_count"`
	}
	var stats RatingStats
	database.DB.Model(&models.Doctor{}).Select("AVG(rating) as avg_rating, SUM(rating_count) as rating_count").Scan(&stats)
	avgRating := stats.AvgRating

	var appointmentRate float64
	if totalAppointments > 0 {
		appointmentRate = float64(completedAppointments) / float64(totalAppointments) * 100
	}

	utils.Success(c, gin.H{
		"total_users":          totalUsers,
		"total_doctors":        totalDoctors,
		"total_appointments":   totalAppointments,
		"today_appointments":   todayAppointments,
		"completed_appointments": completedAppointments,
		"total_consultations": totalConsultations,
		"appointment_rate":    appointmentRate,
		"avg_rating":          avgRating,
	})
}

func GetAppointmentStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	type DailyStats struct {
		Date      string `json:"date"`
		Total     int64  `json:"total"`
		Completed int64  `json:"completed"`
		Cancelled int64  `json:"cancelled"`
	}

	var dailyStats []DailyStats
	database.DB.Model(&models.Appointment{}).
		Select("DATE(date) as date, COUNT(*) as total, SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as completed, SUM(CASE WHEN status = 3 THEN 1 ELSE 0 END) as cancelled").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(date)").
		Order("date ASC").
		Find(&dailyStats)

	utils.Success(c, dailyStats)
}

func GetAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	doctorID := c.Query("doctor_id")
	userID := c.Query("user_id")
	date := c.Query("date")

	offset := (page - 1) * pageSize

	var appointments []models.Appointment
	var total int64

	query := database.DB.Model(&models.Appointment{})
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}
	if doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Count(&total)
	query.Preload("User").Preload("Doctor").Preload("Doctor.Department").Preload("Schedule").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&appointments)

	utils.Page(c, total, page, pageSize, appointments)
}

func GetConsultations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	doctorID := c.Query("doctor_id")

	offset := (page - 1) * pageSize

	var consultations []models.Consultation
	var total int64

	query := database.DB.Model(&models.Consultation{})
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}
	if doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	query.Count(&total)
	query.Preload("User").Preload("Doctor").Preload("Appointment").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&consultations)

	utils.Page(c, total, page, pageSize, consultations)
}

func GetMedicinesAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var medicines []models.Medicine
	var total int64

	query := database.DB.Model(&models.Medicine{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&medicines)

	utils.Page(c, total, page, pageSize, medicines)
}

func CreateMedicine(c *gin.Context) {
	var medicine models.Medicine
	if err := c.ShouldBindJSON(&medicine); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	medicine.Status = 1
	if result := database.DB.Create(&medicine); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", medicine)
}

func UpdateMedicine(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var medicine models.Medicine
	if result := database.DB.First(&medicine, id); result.Error != nil {
		utils.NotFound(c, "药品不存在")
		return
	}

	if err := c.ShouldBindJSON(&medicine); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&medicine); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", medicine)
}

func DeleteMedicine(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Medicine{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetSchedules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	doctorID := c.Query("doctor_id")
	date := c.Query("date")

	offset := (page - 1) * pageSize

	var schedules []models.Schedule
	var total int64

	query := database.DB.Model(&models.Schedule{})
	if doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Count(&total)
	query.Preload("Doctor").Preload("Doctor.Department").
		Order("date DESC, period ASC").
		Offset(offset).Limit(pageSize).
		Find(&schedules)

	utils.Page(c, total, page, pageSize, schedules)
}

func CreateScheduleAdmin(c *gin.Context) {
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
	if result := database.DB.Where("doctor_id = ? AND date = ? AND period = ?", req.DoctorID, scheduleDate, req.Period).First(&existingSchedule); result.Error == nil {
		utils.Error(c, 400, "该时段已存在排班")
		return
	}

	schedule := models.Schedule{
		DoctorID:    req.DoctorID,
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

func UpdateScheduleAdmin(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("id"))

	var schedule models.Schedule
	if result := database.DB.First(&schedule, scheduleID); result.Error != nil {
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

func DeleteScheduleAdmin(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("id"))

	var schedule models.Schedule
	if result := database.DB.First(&schedule, scheduleID); result.Error != nil {
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

func GetDoctorsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	departmentID := c.Query("department_id")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var doctors []models.Doctor
	var total int64

	query := database.DB.Model(&models.Doctor{})
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}
	if keyword != "" {
		query = query.Where("real_name LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Preload("Department").Preload("Hospital").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&doctors)

	utils.Page(c, total, page, pageSize, doctors)
}

func CreateDoctor(c *gin.Context) {
	var doctor models.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingDoctor models.Doctor
	if result := database.DB.Where("username = ?", doctor.Username).First(&existingDoctor); result.Error == nil {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	doctor.Password = hashedPassword
	doctor.Status = 1
	doctor.Rating = 5.0
	doctor.RatingCount = 0

	if result := database.DB.Create(&doctor); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功，默认密码为123456", doctor)
}

func UpdateDoctor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var doctor models.Doctor
	if result := database.DB.First(&doctor, id); result.Error != nil {
		utils.NotFound(c, "医生不存在")
		return
	}

	var req models.DoctorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.DepartmentID > 0 {
		updates["department_id"] = req.DepartmentID
	}
	if req.HospitalID > 0 {
		updates["hospital_id"] = req.HospitalID
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Specialty != "" {
		updates["specialty"] = req.Specialty
	}
	if req.Introduction != "" {
		updates["introduction"] = req.Introduction
	}
	if req.ConsultationFee > 0 {
		updates["consultation_fee"] = req.ConsultationFee
	}

	if result := database.DB.Model(&doctor).Updates(updates); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", doctor)
}

func DeleteDoctor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Doctor{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users)

	utils.Page(c, total, page, pageSize, users)
}

func UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user.Status = req.Status
	if result := database.DB.Save(&user); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func GetDepartmentsAdmin(c *gin.Context) {
	hospitalID := c.Query("hospital_id")

	var departments []models.Department
	query := database.DB.Where("status = 1")
	if hospitalID != "" {
		query = query.Where("hospital_id = ?", hospitalID)
	}
	query.Order("sort ASC, created_at DESC").Find(&departments)

	utils.Success(c, departments)
}

func CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	department.Status = 1
	if result := database.DB.Create(&department); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", department)
}

func UpdateDepartment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var department models.Department
	if result := database.DB.First(&department, id); result.Error != nil {
		utils.NotFound(c, "科室不存在")
		return
	}

	if err := c.ShouldBindJSON(&department); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&department); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", department)
}

func DeleteDepartment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Department{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetMedicalProjectsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var projects []models.MedicalProject
	var total int64

	query := database.DB.Model(&models.MedicalProject{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&projects)

	utils.Page(c, total, page, pageSize, projects)
}

func CreateMedicalProject(c *gin.Context) {
	var project models.MedicalProject
	if err := c.ShouldBindJSON(&project); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	project.Status = 1
	if result := database.DB.Create(&project); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", project)
}

func UpdateMedicalProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.MedicalProject
	if result := database.DB.First(&project, id); result.Error != nil {
		utils.NotFound(c, "诊疗项目不存在")
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&project); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", project)
}

func DeleteMedicalProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.MedicalProject{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetHospitalsAdmin(c *gin.Context) {
	var hospitals []models.Hospital
	database.DB.Order("sort ASC, created_at DESC").Find(&hospitals)
	utils.Success(c, hospitals)
}

func CreateHospital(c *gin.Context) {
	var hospital models.Hospital
	if err := c.ShouldBindJSON(&hospital); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	hospital.Status = 1
	if result := database.DB.Create(&hospital); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", hospital)
}

func UpdateHospital(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var hospital models.Hospital
	if result := database.DB.First(&hospital, id); result.Error != nil {
		utils.NotFound(c, "医院不存在")
		return
	}

	if err := c.ShouldBindJSON(&hospital); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&hospital); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", hospital)
}

func DeleteHospital(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Hospital{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetCarouselsAdmin(c *gin.Context) {
	var carousels []models.Carousel
	database.DB.Order("sort ASC, created_at DESC").Find(&carousels)
	utils.Success(c, carousels)
}

func CreateCarousel(c *gin.Context) {
	var carousel models.Carousel
	if err := c.ShouldBindJSON(&carousel); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	carousel.Status = 1
	if result := database.DB.Create(&carousel); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", carousel)
}

func UpdateCarousel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var carousel models.Carousel
	if result := database.DB.First(&carousel, id); result.Error != nil {
		utils.NotFound(c, "轮播图不存在")
		return
	}

	if err := c.ShouldBindJSON(&carousel); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&carousel); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", carousel)
}

func DeleteCarousel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Carousel{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func GetAnnouncementsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	announcementType := c.Query("type")

	offset := (page - 1) * pageSize

	var announcements []models.Announcement
	var total int64

	query := database.DB.Model(&models.Announcement{})
	if announcementType != "" {
		typeInt, _ := strconv.Atoi(announcementType)
		query = query.Where("type = ?", typeInt)
	}

	query.Count(&total)
	query.Order("is_top DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&announcements)

	utils.Page(c, total, page, pageSize, announcements)
}

func CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Create(&announcement); result.Error != nil {
		utils.InternalServerError(c, "创建失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var announcement models.Announcement
	if result := database.DB.First(&announcement, id); result.Error != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if result := database.DB.Save(&announcement); result.Error != nil {
		utils.InternalServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", announcement)
}

func DeleteAnnouncement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if result := database.DB.Delete(&models.Announcement{}, id); result.Error != nil {
		utils.InternalServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}
