package middleware

import (
	"fmt"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义Lua脚本
const limitScript = `
local count = redis.call('INCR', KEYS[1])
if count == 1 then
	redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return count
`

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
		count, err := redis.RDB.Eval(ctx, limitScript, []string{key}, window.Seconds()).Int()
		if err != nil {
			// 如果redis连不上，为了不影响业务，通常选择放行
			// 这叫fail open策略
			logger.Log.Error("redis限流报错")
			c.Next()
			return
		}

		// 判断是否超限
		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,gin.H{
				"code": 429,
				"msg": "请求频繁，请稍后重试",
			})
			return
		}

		c.Next()
	}
}