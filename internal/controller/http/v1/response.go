package v1

import (
	"github.com/gin-gonic/gin"
	"tomokari/internal/usecase"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code usecase.Status, msg string) {
	c.AbortWithStatusJSON(int(code), response{msg})
}

func responseWithData(c *gin.Context, code usecase.Status, data interface{}) {
	c.JSON(int(code), data)
}
