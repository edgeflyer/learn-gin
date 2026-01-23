package middleware

import (
	"fmt"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户id
		// gin的c.ClientIP()会自动处理X-Forwarded-For等代理头，比较准
		ip := c.ClientIP()

		// 构造Redis Key
		key := fmt.Sprintf("limiter:login:ip:%s", ip)

		// 操作redis
		// 使用request的context，这样请求了redis操作也会取消
		ctx := c.Request.Context()

		// INCR命令：原子加1
		// 如果key不存在，redis会自动先把它变成0，再加1，返回1
		count, err := redis.RDB.Incr(ctx, key).Result()
		if err != nil {
			// 如果redis连不上，为了不影响业务，通常选择放行
			// 这叫fail open策略
			logger.Log.Error("redis限流报错")
			c.Next()
			return
		}

		// 如果是第一次访问，设置过期时间
		// 比如window是60喵，哪60秒之后，这个key就自动消失，计数清零
		if count == 1 {
			redis.RDB.Expire(ctx, key, window)
		}

		// 判断是否超限
		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,gin.H{
				"code": 429,
				"msg": "请求频繁，请稍后重试",
			})
			return
		}

		c.Next()
	}
}