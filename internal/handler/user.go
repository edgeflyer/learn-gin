package handler

import (
	"encoding/json"
	"fmt"
	"learn-gin/internal/logger"
	"learn-gin/internal/redis"
	"learn-gin/internal/response"
	"learn-gin/internal/service"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	srv service.UserService
}

func NewUserHandler(srv service.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	// 定义请求参数结构体
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError)
		return
	}

	if err := h.srv.Register(c, req.Username, req.Password); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError)
		return
	}
	token, err := h.srv.Login(c, req.Username, req.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, gin.H{"token": token})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userIDraw, exists := c.Get("userID")
	if !exists {
		response.FailByError(c, response.AuthError)
		return
	}

	userID := userIDraw.(int64)
	userProfile, err := h.srv.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	
	response.OK(c, userProfile)
}

func (h *UserHandler)SendVerificationCode(c *gin.Context) {
	type Request struct {
		Mail string `json:"mail" binding:"required"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("获取用户邮箱失败")
		response.Fail(c, response.ParamError)
		return
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 业务状态存储
	key := fmt.Sprintf("otp:login:%s", req.Mail)

	ctx := c.Request.Context()
	err := redis.RDB.Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		logger.Log.Error("redis插入错误")
		response.Fail(c, response.RedisErr)
		return
	}

	// 生产者：把任务丢进队列
	taskPayload := map[string]string {
		"mail": req.Mail,
		"code": code,
	}
	// 序列化成json字符串
	taskData, _ := json.Marshal(taskPayload)

	//Rpush:从列表右侧推入(queue)
	err = redis.RDB.RPush(ctx, "queue:mail:send", taskData).Err()
	if err != nil {
		// 入队失败
		response.Fail(c, response.WorkCommitError)
		return
	}

	fmt.Printf("[模拟发送验证码：] 邮箱:%s, 验证码:%s\n", req.Mail, code)
	response.OK(c, gin.H{
		"msg": "验证码发送请求已受理",
		"code": code,
	})
}
func (h *UserHandler)VerifyCode(c *gin.Context) {
	type Request struct {
		Mail string `json:"mail" binding:"required"`
		Code string `json:"code" binding:"required"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("解码请求出错")
		response.Fail(c, response.ParamError)
		return
	}

	key := fmt.Sprintf("otp:login:%s", req.Mail)
	ctx := c.Request.Context()
	saveCode, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		logger.Log.Error("redis中无该验证码")
		response.Fail(c, response.RedisRecordNotfound)
		return
	}

	if saveCode != req.Code {
		logger.Log.Error("验证码错误")
		response.Fail(c, response.VerifyCodeError)
		return
	}

	redis.RDB.Del(ctx, key)

	response.OK(c, gin.H{"msg": "登录成功"})
}