package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RDB是一个全局变量，到处给其他包使用
var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})

	// 测试连接
	// context.Background()是处理超市和取消的标准方式
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}

	fmt.Println("Redis连接成功")
}