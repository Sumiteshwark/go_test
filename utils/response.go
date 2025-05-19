package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse formats a successful API response
func SuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"status":  "success",
		"data":    data,
		"message": nil,
	})
}

// ErrorResponse formats an error API response
func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"status":  "error",
		"data":    nil,
		"message": message,
	})
}
