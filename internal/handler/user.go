package handler

import (
	"learn-gin/internal/response"
	"learn-gin/internal/service"

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
	userID, exists := c.Get("userID")
	if !exists {
		response.FailByError(c, response.AuthError)
		return
	}

	username, _ := c.Get("username")

	response.OK(c, gin.H{
		"user_id": userID,
		"username": username,
		"message": "身份验证通过",
	})
}