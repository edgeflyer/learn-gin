package middleware

import (
	"learn-gin/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger接收Gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		c.Next()

		// 计算耗时
		cost := time.Since(start)
		// 记录日志
		logger.Log.Info("request",
            zap.Int("status", c.Writer.Status()),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.String("ip", c.ClientIP()),
            zap.String("user-agent", c.Request.UserAgent()),
            zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
            zap.Duration("cost", cost),
        )
	}
}