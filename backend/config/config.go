package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

type JWTConfig struct {
	Secret string
	Expire int
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var AppConfig *Config

func InitConfig() {
	_, currentFile, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(currentFile)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %s", err))
	}
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}
