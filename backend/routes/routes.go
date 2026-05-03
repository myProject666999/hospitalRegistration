package routes

import (
	"hospitalRegistration/controllers"
	"hospitalRegistration/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/register", controllers.UserRegister)
		api.POST("/user/login", controllers.UserLogin)
		api.POST("/doctor/login", controllers.DoctorLogin)
		api.POST("/admin/login", controllers.AdminLogin)

		api.GET("/carousels", controllers.GetCarousels)
		api.GET("/announcements", controllers.GetAnnouncements)
		api.GET("/announcements/:id", controllers.GetAnnouncementDetail)
		api.GET("/medicines", controllers.GetMedicines)
		api.GET("/medicines/:id", controllers.GetMedicineDetail)
		api.GET("/medical-projects", controllers.GetMedicalProjects)
		api.GET("/medical-projects/:id", controllers.GetMedicalProjectDetail)
		api.GET("/hospitals", controllers.GetHospitals)
		api.GET("/departments", controllers.GetDepartments)
		api.GET("/doctors", controllers.GetDoctors)
		api.GET("/doctors/:id", controllers.GetDoctorDetail)
		api.GET("/doctors/:id/schedules", controllers.GetDoctorSchedules)
	}

	user := api.Group("/user").Use(middleware.JWTAuth(), middleware.RoleAuth("user"))
	{
		user.GET("/info", controllers.GetUserInfo)
		user.PUT("/info", controllers.UpdateUserInfo)
		user.PUT("/password", controllers.UpdateUserPassword)

		user.GET("/appointments", controllers.GetUserAppointments)
		user.POST("/appointments", controllers.CreateAppointment)
		user.DELETE("/appointments/:id", controllers.CancelAppointment)

		user.GET("/consultations", controllers.GetUserConsultations)
		user.POST("/consultations/:id/comment", controllers.CommentConsultation)

		user.GET("/favorites", controllers.GetUserFavorites)
		user.POST("/favorites", controllers.AddFavorite)
		user.DELETE("/favorites/:id", controllers.RemoveFavorite)
	}

	doctor := api.Group("/doctor").Use(middleware.JWTAuth(), middleware.RoleAuth("doctor"))
	{
		doctor.GET("/info", controllers.GetDoctorInfo)

		doctor.GET("/appointments", controllers.GetDoctorAppointments)
		doctor.PUT("/appointments/:id/confirm", controllers.ConfirmAppointment)
		doctor.POST("/appointments/:id/consultation", controllers.StartConsultation)

		doctor.GET("/consultations", controllers.GetDoctorConsultations)
		doctor.PUT("/consultations/:id", controllers.UpdateConsultation)
		doctor.PUT("/consultations/:id/complete", controllers.CompleteConsultation)

		doctor.GET("/schedules", controllers.GetDoctorOwnSchedules)
		doctor.POST("/schedules", controllers.CreateSchedule)
		doctor.PUT("/schedules/:id", controllers.UpdateSchedule)
		doctor.DELETE("/schedules/:id", controllers.CancelSchedule)
	}

	admin := api.Group("/admin").Use(middleware.JWTAuth(), middleware.RoleAuth("admin"))
	{
		admin.GET("/statistics", controllers.GetStatistics)
		admin.GET("/statistics/appointments", controllers.GetAppointmentStatistics)

		admin.GET("/appointments", controllers.GetAppointments)
		admin.GET("/consultations", controllers.GetConsultations)

		admin.GET("/medicines", controllers.GetMedicinesAdmin)
		admin.POST("/medicines", controllers.CreateMedicine)
		admin.PUT("/medicines/:id", controllers.UpdateMedicine)
		admin.DELETE("/medicines/:id", controllers.DeleteMedicine)

		admin.GET("/schedules", controllers.GetSchedules)
		admin.POST("/schedules", controllers.CreateScheduleAdmin)
		admin.PUT("/schedules/:id", controllers.UpdateScheduleAdmin)
		admin.DELETE("/schedules/:id", controllers.DeleteScheduleAdmin)

		admin.GET("/doctors", controllers.GetDoctorsAdmin)
		admin.POST("/doctors", controllers.CreateDoctor)
		admin.PUT("/doctors/:id", controllers.UpdateDoctor)
		admin.DELETE("/doctors/:id", controllers.DeleteDoctor)

		admin.GET("/users", controllers.GetUsers)
		admin.PUT("/users/:id/status", controllers.UpdateUserStatus)

		admin.GET("/departments", controllers.GetDepartmentsAdmin)
		admin.POST("/departments", controllers.CreateDepartment)
		admin.PUT("/departments/:id", controllers.UpdateDepartment)
		admin.DELETE("/departments/:id", controllers.DeleteDepartment)

		admin.GET("/medical-projects", controllers.GetMedicalProjectsAdmin)
		admin.POST("/medical-projects", controllers.CreateMedicalProject)
		admin.PUT("/medical-projects/:id", controllers.UpdateMedicalProject)
		admin.DELETE("/medical-projects/:id", controllers.DeleteMedicalProject)

		admin.GET("/hospitals", controllers.GetHospitalsAdmin)
		admin.POST("/hospitals", controllers.CreateHospital)
		admin.PUT("/hospitals/:id", controllers.UpdateHospital)
		admin.DELETE("/hospitals/:id", controllers.DeleteHospital)

		admin.GET("/carousels", controllers.GetCarouselsAdmin)
		admin.POST("/carousels", controllers.CreateCarousel)
		admin.PUT("/carousels/:id", controllers.UpdateCarousel)
		admin.DELETE("/carousels/:id", controllers.DeleteCarousel)

		admin.GET("/announcements", controllers.GetAnnouncementsAdmin)
		admin.POST("/announcements", controllers.CreateAnnouncement)
		admin.PUT("/announcements/:id", controllers.UpdateAnnouncement)
		admin.DELETE("/announcements/:id", controllers.DeleteAnnouncement)
	}
}
