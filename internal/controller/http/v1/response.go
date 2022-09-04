package v1

import (
	"github.com/gin-gonic/gin"
	"tomokari/internal/usecase"
)

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func errorResponse(c *gin.Context, code usecase.Status, message string) {
	c.AbortWithStatusJSON(int(code), response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func responseWithData(c *gin.Context, code usecase.Status, data interface{}, message string) {
	c.JSON(int(code), response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
