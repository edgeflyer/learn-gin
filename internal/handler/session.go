package handler

import (
	"learn-gin/internal/response"
	"learn-gin/internal/service"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	srv service.SessionService
}

func NewSessionHandler(srv service.SessionService) *SessionHandler {
	return &SessionHandler{srv: srv}
}

func (h *SessionHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError)
		return
	}

	res, err := h.srv.Login(c, req.Username, req.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{
		"access_token": res.AccessToken,
		"refresh_token": res.RefreshToken,
	})
}

func (h *SessionHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError)
		return
	}

	access, err := h.srv.Refresh(c, req.RefreshToken)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{
		"accesss_token": access,
	})
}

func (h *SessionHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamError)
		return
	}
	
	if err := h.srv.Logout(c, req.RefreshToken); err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, nil)
}
