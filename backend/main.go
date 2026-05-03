package main

import (
	"fmt"
	"hospitalRegistration/config"
	"hospitalRegistration/database"
	"hospitalRegistration/middleware"
	"hospitalRegistration/models"
	"hospitalRegistration/routes"
	"hospitalRegistration/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	utils.InitLogger()
	defer utils.Logger.Sync()

	database.InitDatabase()

	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Doctor{},
		&models.Admin{},
		&models.Department{},
		&models.Hospital{},
		&models.Medicine{},
		&models.Announcement{},
		&models.Carousel{},
		&models.Schedule{},
		&models.Appointment{},
		&models.Consultation{},
		&models.Favorite{},
		&models.MedicalProject{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功")

	database.SeedDatabase()

	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(utils.GinLogger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
