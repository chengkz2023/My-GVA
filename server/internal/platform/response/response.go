package response

import (
	stderrors "errors"
	"net/http"

	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	"github.com/gin-gonic/gin"
)

const (
	Success = 0
	Failure = 7
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

func JSON(c *gin.Context, status int, code int, message string, data any) {
	c.JSON(status, Body{
		Code:    code,
		Message: message,
		Msg:     message,
		Data:    data,
	})
}

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, Success, "ok", data)
}

func Fail(c *gin.Context, status int, code int, message string) {
	JSON(c, status, code, message, gin.H{})
}

func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var appErr *apperrors.Error
	if stderrors.As(err, &appErr) {
		code := appErr.Code
		if code == 0 {
			code = Failure
		}
		Fail(c, apperrors.HTTPStatus(appErr.Kind), code, appErr.Error())
		return
	}

	Fail(c, http.StatusInternalServerError, Failure, "internal error")
}
