package helper

import (
	"github.com/gin-gonic/gin"
)

func FormatResponse(c *gin.Context, status string, httpStatus int, message any, data any, meta any) {
	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
