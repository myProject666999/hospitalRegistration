package controllers

import (
	"hospitalRegistration/database"
	"hospitalRegistration/models"
	"hospitalRegistration/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCarousels(c *gin.Context) {
	var carousels []models.Carousel
	database.DB.Where("status = 1").Order("sort ASC, created_at DESC").Find(&carousels)
	utils.Success(c, carousels)
}

func GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	announcementType := c.Query("type")

	offset := (page - 1) * pageSize

	var announcements []models.Announcement
	var total int64

	query := database.DB.Model(&models.Announcement{}).Where("status = 1")
	if announcementType != "" {
		typeInt, _ := strconv.Atoi(announcementType)
		query = query.Where("type = ?", typeInt)
	}

	query.Count(&total)
	query.Order("is_top DESC, sort ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&announcements)

	utils.Page(c, total, page, pageSize, announcements)
}

func GetAnnouncementDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var announcement models.Announcement
	if result := database.DB.First(&announcement, id); result.Error != nil {
		utils.NotFound(c, "公告不存在")
		return
	}

	announcement.ViewCount++
	database.DB.Save(&announcement)

	utils.Success(c, announcement)
}

func GetMedicines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var medicines []models.Medicine
	var total int64

	query := database.DB.Model(&models.Medicine{}).Where("status = 1")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("sort ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&medicines)

	utils.Page(c, total, page, pageSize, medicines)
}

func GetMedicineDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var medicine models.Medicine
	if result := database.DB.First(&medicine, id); result.Error != nil {
		utils.NotFound(c, "药品不存在")
		return
	}

	utils.Success(c, medicine)
}

func GetMedicalProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var projects []models.MedicalProject
	var total int64

	query := database.DB.Model(&models.MedicalProject{}).Where("status = 1")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("sort ASC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&projects)

	utils.Page(c, total, page, pageSize, projects)
}

func GetMedicalProjectDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var project models.MedicalProject
	if result := database.DB.First(&project, id); result.Error != nil {
		utils.NotFound(c, "诊疗项目不存在")
		return
	}

	utils.Success(c, project)
}

func GetHospitals(c *gin.Context) {
	var hospitals []models.Hospital
	database.DB.Where("status = 1").Order("sort ASC, created_at DESC").Find(&hospitals)
	utils.Success(c, hospitals)
}

func GetDepartments(c *gin.Context) {
	hospitalID := c.Query("hospital_id")

	var departments []models.Department
	query := database.DB.Where("status = 1")
	if hospitalID != "" {
		query = query.Where("hospital_id = ?", hospitalID)
	}
	query.Order("sort ASC, created_at DESC").Find(&departments)

	utils.Success(c, departments)
}

func GetDoctors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	departmentID := c.Query("department_id")
	hospitalID := c.Query("hospital_id")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	var doctors []models.Doctor
	var total int64

	query := database.DB.Model(&models.Doctor{}).Where("status = 1")
	if departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}
	if hospitalID != "" {
		query = query.Where("hospital_id = ?", hospitalID)
	}
	if keyword != "" {
		query = query.Where("real_name LIKE ? OR title LIKE ? OR specialty LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Preload("Department").Preload("Hospital").
		Order("rating DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&doctors)

	utils.Page(c, total, page, pageSize, doctors)
}

func GetDoctorDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var doctor models.Doctor
	if result := database.DB.Preload("Department").Preload("Hospital").First(&doctor, id); result.Error != nil {
		utils.NotFound(c, "医生不存在")
		return
	}

	utils.Success(c, doctor)
}

func GetDoctorSchedules(c *gin.Context) {
	doctorID, _ := strconv.Atoi(c.Param("id"))
	date := c.Query("date")

	var schedules []models.Schedule
	query := database.DB.Where("doctor_id = ?", doctorID)
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Order("date ASC, period ASC").Find(&schedules)

	utils.Success(c, schedules)
}
