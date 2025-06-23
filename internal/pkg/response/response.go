package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// APIResponse 标准 API 响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// APIError 标准 API 错误响应结构
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// sendError 统一发送错误响应
func sendError(c *gin.Context, statusCode int, code int, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, APIError{
		Code:    code,
		Message: message,
		Error:   errMsg,
	})
}

// 统一的成功响应
func Success(c *gin.Context, data interface{}) {
	resp := APIResponse{
		Code:    0,
		Message: "Success",
		Data:    data,
	}
	c.JSON(http.StatusOK, resp)
}

// ValidationError 发送校验失败错误
func ValidationError(c *gin.Context, message string) {
	if message == "" {
		message = "请求参数不合法"
	}
	sendError(c, http.StatusBadRequest, 400, message, nil)
}

// InternalError 发送服务器内部错误
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器发生未知错误"
	}
	sendError(c, http.StatusInternalServerError, 500, message, nil)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "请求的资源不存在"
	}
	sendError(c, http.StatusNotFound, 404, message, nil)
}

// Created 发送创建成功的响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data:    data,
	})
}
