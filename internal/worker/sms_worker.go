package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"time"
)

// 启动邮箱消费者
// 这是一个死循环，要在main里的goroutine中运行
func StartSMSWorker() {
	logger.Log.Info("[worker] 正在监听邮箱任务队列")

	for {
		// 从队列中阻塞读取
		// BLPop(ctx, 超时时间, key...)
		// timeout = 0 表示无限等待，知道队列里有东西
		// result是一个切片，[0]=kye名字, [1]=取出的value
		result, err := redis.RDB.BLPop(context.Background(), 0, "queue:sms:send").Result()
		if err != nil {
			// 如果redis连接断了，休息一下重试，防止cup空转
			logger.Log.Error("[worker] redis读取错误")
			time.Sleep(3 * time.Second)
			continue
		}

		// 处理任务
		rawJSON := result[1]
		// 解析json
		var task map[string]string
		if err := json.Unmarshal([]byte(rawJSON), &task); err != nil {
			logger.Log.Error("[worker] 任务格式错误")
			continue
		}

		mail := task["mail"]
		code := task["code"]

		fmt.Printf("mail: %s, code: %s\n", mail, code)

		// 模拟发送过程
		logger.Log.Info("[worker] 正在发送邮箱验证码")
		time.Sleep(2 * time.Second)
		logger.Log.Info("[worker] 短信发送成功")
	}
}