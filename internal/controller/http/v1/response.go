package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
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

func errorResponse2(res http.ResponseWriter, code usecase.Status, message string) {
	res.WriteHeader(int(code))
	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(response{
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

func responseWithData2(res http.ResponseWriter, code usecase.Status, data interface{}, message string) {
	res.WriteHeader(int(code))
	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
