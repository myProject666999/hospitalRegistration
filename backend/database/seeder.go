package database

import (
	"hospitalRegistration/models"
	"hospitalRegistration/utils"
	"log"

	"gorm.io/gorm"
)

func SeedDatabase() {
	if err := seedAdmin(); err != nil {
		log.Printf("管理员种子数据失败: %v", err)
	}
}

func seedAdmin() error {
	var admin models.Admin
	result := DB.Where("username = ?", "admin").First(&admin)
	if result.Error != gorm.ErrRecordNotFound {
		log.Println("管理员账号已存在")
		return nil
	}

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin = models.Admin{
		Username: "admin",
		Password: hashedPassword,
		RealName: "超级管理员",
		Role:     1,
		Status:   1,
	}

	if result := DB.Create(&admin); result.Error != nil {
		return result.Error
	}

	log.Println("默认管理员账号创建成功: 用户名 admin, 密码 admin123")
	return nil
}
