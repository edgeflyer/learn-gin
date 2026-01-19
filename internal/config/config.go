package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT JWTConfig `mapstructure:"jwt"`
	Log LogConfig `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Dsn string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int `mapstructure:"expire"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxAge     int    `mapstructure:"max_age"`
    MaxBackups int    `mapstructure:"max_backups"`
}

func Init() {
	v := viper.New()
	
	v.SetConfigName("config") // 文件名
	v.SetConfigType("yaml") // 文件类型
	v.AddConfigPath("./config") // 查找路径（根目录下的config文件夹）

	// 支持SERVER_PORT覆盖server.port
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("读取被指文件失败：%v", err)
	}

	if err := v.Unmarshal(&Conf); err != nil {
		log.Fatalf("配置解析失败：%v", err)
	}

	if Conf.JWT.Secret == "" {
		log.Fatalf("配置缺失：jwt.secret")
	}
	if Conf.Database.Driver == "" {
		log.Fatal("配置缺失：database.dsn")
	}

	log.Println("配置加载成功")
}