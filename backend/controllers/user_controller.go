package controllers

import (
	"fmt"
	"hospitalRegistration/database"
	"hospitalRegistration/models"
	"hospitalRegistration/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingUser models.User
	if result := database.DB.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.Error(c, 400, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   1,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		utils.InternalServerError(c, "注册失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功", nil)
}

func UserLogin(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Error(c, 400, "用户名或密码错误")
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"phone":     user.Phone,
			"email":     user.Email,
			"avatar":    user.Avatar,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
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
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
	}
	if req.Gender > 0 {
		updates["gender"] = req.Gender
	}
	if req.Birthday != nil {
		updates["birthday"] = req.Birthday
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if result := database.DB.Model(&user).Updates(updates); result.Error != nil {
		utils.InternalServerError(c, "更新失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func UpdateUserPassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.UserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		utils.Error(c, 400, "原密码错误")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalServerError(c, "密码加密失败")
		return
	}

	user.Password = hashedPassword
	if result := database.DB.Save(&user); result.Error != nil {
		utils.InternalServerError(c, "更新密码失败")
		return
	}

	utils.SuccessWithMessage(c, "密码更新成功", nil)
}

func GetUserAppointments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	var appointments []models.Appointment
	var total int64

	query := database.DB.Model(&models.Appointment{}).Where("user_id = ?", userID)
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)
	query.Preload("Doctor").Preload("Doctor.Department").Preload("Schedule").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&appointments)

	utils.Page(c, total, page, pageSize, appointments)
}

func CreateAppointment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.AppointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var schedule models.Schedule
	if result := database.DB.First(&schedule, req.ScheduleID); result.Error != nil {
		utils.NotFound(c, "排班不存在")
		return
	}

	if schedule.Status != 1 || schedule.RemainCount <= 0 {
		utils.Error(c, 400, "该时段已约满")
		return
	}

	var existingAppointment models.Appointment
	if result := database.DB.Where("user_id = ? AND schedule_id = ? AND status IN (0, 1)", userID, req.ScheduleID).First(&existingAppointment); result.Error == nil {
		utils.Error(c, 400, "您已预约该时段")
		return
	}

	appointmentDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.BadRequest(c, "日期格式错误")
		return
	}

	appointmentNo := fmt.Sprintf("APT%s%d", time.Now().Format("20060102150405"), userID)

	var doctor models.Doctor
	database.DB.First(&doctor, req.DoctorID)

	appointment := models.Appointment{
		AppointmentNo:    appointmentNo,
		UserID:           userID,
		DoctorID:         req.DoctorID,
		ScheduleID:       req.ScheduleID,
		Date:             appointmentDate,
		Period:           req.Period,
		TimeSlot:         req.TimeSlot,
		PatientName:      req.PatientName,
		PatientPhone:     req.PatientPhone,
		PatientIDCard:    req.PatientIDCard,
		PatientGender:    req.PatientGender,
		PatientAge:       req.PatientAge,
		ChiefComplaint:   req.ChiefComplaint,
		ConsultationFee:  doctor.ConsultationFee,
		Status:           0,
		Location:         schedule.Location,
	}

	tx := database.DB.Begin()

	if result := tx.Create(&appointment); result.Error != nil {
		tx.Rollback()
		utils.InternalServerError(c, "创建预约失败: "+result.Error.Error())
		return
	}

	schedule.RemainCount--
	if schedule.RemainCount == 0 {
		schedule.Status = 0
	}
	if result := tx.Save(&schedule); result.Error != nil {
		tx.Rollback()
		utils.InternalServerError(c, "更新排班失败")
		return
	}

	tx.Commit()

	utils.SuccessWithMessage(c, "预约成功", appointment)
}

func CancelAppointment(c *gin.Context) {
	userID := c.GetUint("user_id")
	appointmentID, _ := strconv.Atoi(c.Param("id"))

	var appointment models.Appointment
	if result := database.DB.Where("id = ? AND user_id = ?", appointmentID, userID).First(&appointment); result.Error != nil {
		utils.NotFound(c, "预约不存在")
		return
	}

	if appointment.Status != 0 && appointment.Status != 1 {
		utils.Error(c, 400, "该预约状态不可取消")
		return
	}

	appointment.Status = 3
	if result := database.DB.Save(&appointment); result.Error != nil {
		utils.InternalServerError(c, "取消失败")
		return
	}

	var schedule models.Schedule
	if result := database.DB.First(&schedule, appointment.ScheduleID); result.Error == nil {
		schedule.RemainCount++
		if schedule.RemainCount > 0 {
			schedule.Status = 1
		}
		database.DB.Save(&schedule)
	}

	utils.SuccessWithMessage(c, "取消成功", nil)
}

func GetUserFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	favoriteType := c.Query("type")

	var favorites []models.Favorite
	query := database.DB.Where("user_id = ?", userID)
	if favoriteType != "" {
		typeInt, _ := strconv.Atoi(favoriteType)
		query = query.Where("type = ?", typeInt)
	}

	query.Find(&favorites)

	for i := range favorites {
		if favorites[i].Type == 1 {
			var doctor models.Doctor
			database.DB.Preload("Department").First(&doctor, favorites[i].TargetID)
			favorites[i].Doctor = &doctor
		} else if favorites[i].Type == 2 {
			var medicine models.Medicine
			database.DB.First(&medicine, favorites[i].TargetID)
			favorites[i].Medicine = &medicine
		}
	}

	utils.Success(c, favorites)
}

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existingFavorite models.Favorite
	if result := database.DB.Where("user_id = ? AND type = ? AND target_id = ?", userID, req.Type, req.TargetID).First(&existingFavorite); result.Error == nil {
		utils.Error(c, 400, "已收藏")
		return
	}

	favorite := models.Favorite{
		UserID:   userID,
		Type:     req.Type,
		TargetID: req.TargetID,
	}

	if result := database.DB.Create(&favorite); result.Error != nil {
		utils.InternalServerError(c, "收藏失败: "+result.Error.Error())
		return
	}

	utils.SuccessWithMessage(c, "收藏成功", nil)
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	favoriteID, _ := strconv.Atoi(c.Param("id"))

	var favorite models.Favorite
	if result := database.DB.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favorite); result.Error != nil {
		utils.NotFound(c, "收藏不存在")
		return
	}

	if result := database.DB.Delete(&favorite); result.Error != nil {
		utils.InternalServerError(c, "取消收藏失败")
		return
	}

	utils.SuccessWithMessage(c, "取消收藏成功", nil)
}

func GetUserConsultations(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var consultations []models.Consultation
	var total int64

	database.DB.Model(&models.Consultation{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).
		Preload("Doctor").Preload("Doctor.Department").
		Preload("Appointment").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&consultations)

	utils.Page(c, total, page, pageSize, consultations)
}

func CommentConsultation(c *gin.Context) {
	userID := c.GetUint("user_id")
	consultationID, _ := strconv.Atoi(c.Param("id"))

	var req models.ConsultationCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var consultation models.Consultation
	if result := database.DB.Where("id = ? AND user_id = ?", consultationID, userID).First(&consultation); result.Error != nil {
		utils.NotFound(c, "就诊记录不存在")
		return
	}

	if consultation.Status == 2 {
		utils.Error(c, 400, "已评价，不可重复评价")
		return
	}

	now := time.Now()
	consultation.Rating = req.Rating
	consultation.Comment = req.Comment
	consultation.CommentTime = &now
	consultation.Status = 2

	if result := database.DB.Save(&consultation); result.Error != nil {
		utils.InternalServerError(c, "评价失败")
		return
	}

	var doctor models.Doctor
	if result := database.DB.First(&doctor, consultation.DoctorID); result.Error == nil {
		oldTotal := float64(doctor.Rating) * float64(doctor.RatingCount)
		doctor.RatingCount++
		newTotal := oldTotal + float64(req.Rating)
		doctor.Rating = newTotal / float64(doctor.RatingCount)
		database.DB.Save(&doctor)
	}

	utils.SuccessWithMessage(c, "评价成功", nil)
}
