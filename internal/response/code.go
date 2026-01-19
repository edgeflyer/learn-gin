package response

import "net/http"

type AppError struct {
	Code int // 业务错误码
	Msg string // 错误信息(给用户看)
	HttpCode int // HTTP状态码
}

func (e *AppError) Error() string {
	return e.Msg
}

// 错误类型
var (
	Success = NewError(0, "success", http.StatusOK)
	ServerError = NewError(10001, "服务器内部错误", http.StatusInternalServerError)
	ParamError = NewError(10002, "参数错误", http.StatusBadRequest)
	AuthError = NewError(10003, "未授权", http.StatusUnauthorized)
	NotFound = NewError(10004, "资源不存在", http.StatusNotFound)
	UserExists = NewError(10005, "用户已存在", http.StatusConflict)
	UserNotFound = NewError(10006, "用户不存在", http.StatusConflict) // 这个状态码应该写什么？
)

func NewError(code int, msg string, httpCode int) *AppError {
	return &AppError{
		Code: code,
		Msg: msg,
		HttpCode: httpCode,
	}
}