package helper

import (
	"github.com/gin-gonic/gin"
)

func FormatResponse(ctx *gin.Context, status string, httpStatus int, message any, data any, meta any) {
	ctx.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
