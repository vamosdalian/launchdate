package api

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	CodeSuccess = 0
	CodeFailed  = 1
)

func (h *Handler) Json(c *gin.Context, payload any) {
	c.JSON(200, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    payload,
	})
}

func (h *Handler) Error(c *gin.Context, msg string) {
	c.JSON(200, Response{
		Code:    CodeFailed,
		Message: msg,
	})
}

func (h *Handler) Success(c *gin.Context, msg string) {
	c.JSON(200, Response{
		Code:    CodeSuccess,
		Message: msg,
	})
}
