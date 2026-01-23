package main

import (
	"learn-gin/internal/config"
	"learn-gin/internal/db"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"learn-gin/internal/router"
	"log"

	"go.uber.org/zap"
)

// @title gin用户中心API
// @version 1.0
// @description 这是一个基于Gin的go web项目
// termsOfService http://swagger.io.terms/

// @contact.name API support
// @contact.url
func main() {
	config.Init()

	redis.Init()
	logger.Init()
	defer logger.Log.Sync()
	logger.Log.Info("服务正在启动")
	database, err := db.Init(config.Conf.Database.Dsn)
	if err != nil {
		log.Fatalf("无法连接数据库：%v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		logger.Log.Fatal("获取底层sql.DB失败")
	}
	defer sqlDB.Close()

	// 初始化路由
	r := router.NewRouter(database)

	logger.Log.Info("服务启动成功", zap.String("port", config.Conf.Server.Port))
	if err := r.Run(":" + config.Conf.Server.Port); err != nil {
		logger.Log.Fatal("服务启动失败", zap.Error(err))
	}
}