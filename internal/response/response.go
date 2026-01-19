package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int `json:"code"` // 业务码
	Msg string `json:"msg"` // 提示信息
	Data any `json:"data,omitempty"` // 数据
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg: "success",
		Data: data,
	})
}

func Fail(c *gin.Context, err error) {
	var (
		// 默认错误信息
		code = ServerError.Code
		msg = err.Error()
		httpCode = http.StatusInternalServerError
	)

	// 断言类型，判断这个err是不是定义的AppError
	var appErr *AppError

	if errors.As(err, &appErr) {
		code = appErr.Code
		msg = appErr.Msg
		httpCode = appErr.HttpCode
	}

	c.JSON(httpCode, Response{
		Code: code,
		Msg: msg,
		Data: nil,
	})
}

func FailByError(c *gin.Context, appErr *AppError) {
	c.JSON(appErr.HttpCode, Response{
		Code: appErr.Code,
		Msg: appErr.Msg,
		Data: nil,
	})
}