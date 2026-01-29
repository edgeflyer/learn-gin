package main

import (
	"context"
	"learn-gin/internal/config"
	"learn-gin/internal/db"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"learn-gin/internal/router"
	"learn-gin/internal/worker"
	"learn-gin/pb"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var limiterClient pb.RateLimiterClient

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
	go worker.StartSMSWorker()
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
	r.Use(RateLimitMiddleware())

	logger.Log.Info("服务启动成功", zap.String("port", config.Conf.Server.Port))
	if err := r.Run(":" + config.Conf.Server.Port); err != nil {
		logger.Log.Fatal("服务启动失败", zap.Error(err))
	}
}

func initGRPCClient() {
	//建立长连接（在生产环境建议增加断线重连机制）
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接grpc服务端: %v", err)
	}
	// 创建全局的客户端对象
	limiterClient = pb.NewRateLimiterClient(conn)
}

// grpc中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 构造请求数据，把客户端ip发给服务端
		req := &pb.CheckRequest{
			Key: c.ClientIP(),
			Limit: 5, // 设定每秒能访问五次
		}

		// 设置超时时间，防止grpc server挂了导致gin一直阻塞
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*500)
		defer cancel()

		// 远程调用check方法
		resp, err := limiterClient.Check(ctx, req)

		if err != nil {
			log.Printf("grpc调用失败: %v", err)
			c.Next() // 如果服务挂了，通常降级，让服务通过
			return
		}

		if !resp.Allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg": resp.Reason,
				"detail": "Wait for 10s",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}