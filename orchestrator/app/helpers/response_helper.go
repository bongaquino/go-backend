package helpers

import (
	"github.com/gin-gonic/gin"
)

func FormatResponse(c *gin.Context, status string, httpStatus int, message string, data interface{}, meta interface{}) {
	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
