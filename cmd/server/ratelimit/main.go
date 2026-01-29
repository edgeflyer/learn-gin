package main

import (
	"context"
	"fmt"
	"learn-gin/pb"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)


var rdb *redis.Client
var ctx = context.Background()

type server struct {
	pb.UnimplementedRateLimiterServer
}

// 实现业务逻辑（check方法）
func (s *server) Check(c context.Context, in *pb.CheckRequest) (*pb.CheckResponse, error) {
	key := fmt.Sprintf("limit:%s", in.Key)
	limit := int64(in.Limit)

	count, err := rdb.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return &pb.CheckResponse{Allowed: false, Reason: "Redis 错误"}, nil
	}

	if count >= limit {
		return &pb.CheckResponse{Allowed: false, Reason: "请求过于频繁，请稍后再试"}, nil
	}

	// 计数增加，设置10秒过期（模拟滑动窗口）
	pipe := rdb.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Second*10)
	_, _ = pipe.Exec(ctx)

	log.Printf("限流检查：Key=%s, Current=%d, Limit=%d", in.Key, count+1, limit)
	return &pb.CheckResponse{Allowed: true}, nil
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("无法监听端口")
	}
	s := grpc.NewServer()
	pb.RegisterRateLimiterServer(s, &server{})
	log.Println("grpc限流服务启动在:50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务停止")
	}
}