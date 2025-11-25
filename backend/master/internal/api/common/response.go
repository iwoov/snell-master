package common

import "github.com/gin-gonic/gin"

// Response 统一响应结构。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Message: "success", Data: data})
}

// Created 返回 201。
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Response{Code: 0, Message: "success", Data: data})
}

// Fail 返回错误响应。
func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Response{Code: status, Message: message})
}
