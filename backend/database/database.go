package database

import (
	"hospitalRegistration/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := config.AppConfig.Database.DSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")
}
